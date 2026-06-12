package vehicle

import (
	"bufio"
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
)

type stubReader struct{}

func (stubReader) Read(ctx context.Context, pid string) (float64, string, error) {
	return 0, "", nil
}

func TestParseEndpointSerialPath(t *testing.T) {
	endpoint, err := ParseEndpoint(v3config.OBDConfig{
		Address: "serial:///dev/ttyUSB0",
		Timeout: 1000,
	})
	if err != nil {
		t.Fatalf("expected serial endpoint to parse: %v", err)
	}
	if endpoint.Scheme != "serial" {
		t.Fatalf("expected serial scheme, got %q", endpoint.Scheme)
	}
	if endpoint.SerialTarget != "/dev/ttyUSB0" {
		t.Fatalf("expected serial target /dev/ttyUSB0, got %q", endpoint.SerialTarget)
	}
	if endpoint.Timeout != time.Second {
		t.Fatalf("expected 1s timeout, got %s", endpoint.Timeout)
	}
}

func TestParseEndpointTCPAddress(t *testing.T) {
	endpoint, err := ParseEndpoint(v3config.OBDConfig{
		Address: "tcp://127.0.0.1:35000",
		Timeout: 1500,
	})
	if err != nil {
		t.Fatalf("expected tcp endpoint to parse: %v", err)
	}
	if endpoint.Scheme != "tcp" {
		t.Fatalf("expected tcp scheme, got %q", endpoint.Scheme)
	}
	if endpoint.TCPAddress != "127.0.0.1:35000" {
		t.Fatalf("expected tcp address, got %q", endpoint.TCPAddress)
	}
	if endpoint.Timeout != 1500*time.Millisecond {
		t.Fatalf("expected 1500ms timeout, got %s", endpoint.Timeout)
	}
}

func TestParseEndpointRejectsUnsupportedScheme(t *testing.T) {
	_, err := ParseEndpoint(v3config.OBDConfig{
		Address: "udp://127.0.0.1:35000",
		Timeout: 1000,
	})
	if err == nil {
		t.Fatalf("expected unsupported scheme to fail")
	}
}

func TestConnectPlanUsesResolvedEndpoint(t *testing.T) {
	var gotTarget string
	connector := Connector{
		NewSerialReader: func(target string) (Reader, error) {
			gotTarget = target
			return stubReader{}, nil
		},
	}
	plan := v3config.RuntimePlan{
		Endpoint: v3config.OBDConfig{
			Address: "serial:///dev/ttyUSB0",
			Timeout: 1000,
		},
	}

	reader, err := connector.ConnectPlan(context.Background(), plan)
	if err != nil {
		t.Fatalf("expected plan endpoint to connect: %v", err)
	}
	if reader == nil {
		t.Fatalf("expected reader")
	}
	if gotTarget != "/dev/ttyUSB0" {
		t.Fatalf("expected serial target from selected runtime plan, got %q", gotTarget)
	}
}

func TestConnectEndpointTCPReader(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()

	connector := Connector{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			if network != "tcp" {
				t.Fatalf("expected tcp network, got %q", network)
			}
			if address != "127.0.0.1:35000" {
				t.Fatalf("expected configured tcp address, got %q", address)
			}
			return client, nil
		},
	}

	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		serverReader := bufio.NewReader(server)
		line, err := serverReader.ReadString('\n')
		if err != nil {
			t.Errorf("server read PID: %v", err)
			return
		}
		if line != "010D\n" {
			t.Errorf("expected PID request, got %q", line)
			return
		}
		if _, err := server.Write([]byte("010D 88 km/h\n")); err != nil {
			t.Errorf("server write response: %v", err)
		}
	}()

	reader, err := connector.ConnectEndpoint(context.Background(), v3config.OBDConfig{
		Address: "tcp://127.0.0.1:35000",
		Timeout: 1000,
	})
	if err != nil {
		t.Fatalf("expected tcp endpoint to connect: %v", err)
	}

	value, unit, err := reader.Read(context.Background(), "010D")
	if err != nil {
		t.Fatalf("expected tcp reader to return value: %v", err)
	}
	if value != 88 {
		t.Fatalf("expected value 88, got %v", value)
	}
	if unit != "km/h" {
		t.Fatalf("expected km/h unit, got %q", unit)
	}

	<-serverDone
}

func TestTCPReaderReadReturnsCancelledContextBeforeBlockingIO(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	reader := NewTCPReader(client, 0)
	done := make(chan error, 1)
	go func() {
		_, _, err := reader.Read(ctx, "010D")
		done <- err
	}()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("expected cancelled context to return before blocking I/O")
	}
}

func TestParseTCPReadingAllowsValueOnlyResponse(t *testing.T) {
	value, unit, err := parseTCPReading("010C", "1234\n")
	if err != nil {
		t.Fatalf("expected value-only response to parse: %v", err)
	}
	if value != 1234 {
		t.Fatalf("expected 1234, got %v", value)
	}
	if unit != "" {
		t.Fatalf("expected empty unit, got %q", unit)
	}
}
