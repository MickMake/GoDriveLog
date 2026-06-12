package vehicle

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

// Reader is the v3 endpoint-facing reader contract used by the later sensor runtime.
// It intentionally matches the existing sensors.Reader behaviour while keeping the
// selected endpoint decision behind this package.
type Reader interface {
	Read(ctx context.Context, pid string) (float64, string, error)
}

// Connector creates a Reader for a resolved v3 runtime plan endpoint.
// Runtime code should depend on this boundary rather than branching on endpoint type.
type Connector struct {
	DialContext     func(ctx context.Context, network, address string) (net.Conn, error)
	NewSerialReader func(target string) (Reader, error)
}

// Endpoint is the parsed, validated endpoint address from the selected vehicle.
type Endpoint struct {
	Scheme       string
	SerialTarget string
	TCPAddress   string
	Timeout      time.Duration
}

// NewConnector returns the default v3 endpoint connector.
func NewConnector() Connector {
	dialer := net.Dialer{}
	return Connector{
		DialContext: dialer.DialContext,
		NewSerialReader: func(target string) (Reader, error) {
			return sensors.NewELMOBDReader(target, false)
		},
	}
}

// ConnectPlan creates a Reader from a resolved runtime plan.
func (c Connector) ConnectPlan(ctx context.Context, plan v3config.RuntimePlan) (Reader, error) {
	return c.ConnectEndpoint(ctx, plan.Endpoint)
}

// ConnectEndpoint creates a Reader for one selected v3 OBD endpoint.
func (c Connector) ConnectEndpoint(ctx context.Context, cfg v3config.OBDConfig) (Reader, error) {
	endpoint, err := ParseEndpoint(cfg)
	if err != nil {
		return nil, err
	}

	switch endpoint.Scheme {
	case "serial":
		newSerialReader := c.NewSerialReader
		if newSerialReader == nil {
			newSerialReader = NewConnector().NewSerialReader
		}
		return newSerialReader(endpoint.SerialTarget)
	case "tcp":
		dialContext := c.DialContext
		if dialContext == nil {
			dialContext = NewConnector().DialContext
		}
		conn, err := dialContext(ctx, "tcp", endpoint.TCPAddress)
		if err != nil {
			return nil, err
		}
		return NewTCPReader(conn, endpoint.Timeout), nil
	default:
		return nil, fmt.Errorf("unsupported endpoint scheme %q", endpoint.Scheme)
	}
}

// ParseEndpoint validates and normalises the selected vehicle endpoint address.
func ParseEndpoint(cfg v3config.OBDConfig) (Endpoint, error) {
	u, err := url.Parse(strings.TrimSpace(cfg.Address))
	if err != nil {
		return Endpoint{}, fmt.Errorf("endpoint address must parse as URL: %w", err)
	}

	timeout := time.Duration(cfg.Timeout) * time.Millisecond
	endpoint := Endpoint{Scheme: u.Scheme, Timeout: timeout}

	switch u.Scheme {
	case "serial":
		target := strings.TrimSpace(u.Opaque)
		if target == "" {
			target = strings.TrimSpace(u.Path)
		}
		if target == "" {
			return Endpoint{}, fmt.Errorf("serial endpoint must include a non-empty serial path")
		}
		endpoint.SerialTarget = target
		return endpoint, nil
	case "tcp":
		if strings.TrimSpace(u.Hostname()) == "" || strings.TrimSpace(u.Port()) == "" {
			return Endpoint{}, fmt.Errorf("tcp endpoint must include host and port")
		}
		endpoint.TCPAddress = u.Host
		return endpoint, nil
	default:
		return Endpoint{}, fmt.Errorf("endpoint address must use serial:// or tcp://")
	}
}

// TCPReader is a small line-oriented bench/simulator reader.
// Protocol: write "<PID>\n", then read either "<value> [unit]" or "<PID> <value> [unit]".
type TCPReader struct {
	mu      sync.Mutex
	conn    net.Conn
	reader  *bufio.Reader
	timeout time.Duration
}

func NewTCPReader(conn net.Conn, timeout time.Duration) *TCPReader {
	return &TCPReader{
		conn:    conn,
		reader:  bufio.NewReader(conn),
		timeout: timeout,
	}
}

func (r *TCPReader) Read(ctx context.Context, pid string) (float64, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return 0, "", err
	}
	if err := r.applyDeadline(ctx); err != nil {
		return 0, "", err
	}
	if err := ctx.Err(); err != nil {
		return 0, "", err
	}

	if _, err := fmt.Fprintf(r.conn, "%s\n", pid); err != nil {
		return 0, "", err
	}

	line, err := r.reader.ReadString('\n')
	if err != nil {
		return 0, "", err
	}
	return parseTCPReading(pid, line)
}

func (r *TCPReader) Close() error {
	return r.conn.Close()
}

func (r *TCPReader) applyDeadline(ctx context.Context) error {
	deadline := time.Time{}
	if r.timeout > 0 {
		deadline = time.Now().Add(r.timeout)
	}
	if ctxDeadline, ok := ctx.Deadline(); ok && (deadline.IsZero() || ctxDeadline.Before(deadline)) {
		deadline = ctxDeadline
	}
	if !deadline.IsZero() {
		return r.conn.SetDeadline(deadline)
	}
	return nil
}

func parseTCPReading(pid, line string) (float64, string, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return 0, "", fmt.Errorf("empty response for PID %s", pid)
	}

	if strings.EqualFold(fields[0], pid) {
		fields = fields[1:]
	}
	if len(fields) == 0 {
		return 0, "", fmt.Errorf("response for PID %s did not include a value", pid)
	}

	value, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, "", fmt.Errorf("response for PID %s value %q must be numeric: %w", pid, fields[0], err)
	}

	unit := ""
	if len(fields) > 1 {
		unit = strings.Join(fields[1:], " ")
	}
	return value, unit, nil
}
