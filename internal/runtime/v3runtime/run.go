package v3runtime

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
	jsonlogger "github.com/MickMake/GoDriveLog/internal/logger"
	"github.com/MickMake/GoDriveLog/internal/sensors"
	"github.com/MickMake/GoDriveLog/internal/vehicle"
)

const defaultEventBuffer = 32

// DashboardSink is the v3 dashboard output boundary for v3.1.0.
// A later display adapter can consume this boundary without reading sensors or
// touching OBD endpoint code directly.
type DashboardSink func([]v3dashboard.Scene) error

// Options controls one runnable v3 command path execution.
type Options struct {
	ConfigPath    string
	VehicleID     string
	RepoRoot      string
	Connector     vehicle.Connector
	DashboardSink DashboardSink
	EventBuffer   int
	Logger        *log.Logger
}

// Summary describes the resolved runnable v3 path.
type Summary struct {
	VehicleID      string
	VehicleName    string
	Endpoint       string
	SensorCount    int
	LogCount       int
	DashboardCount int
}

// Run loads v3 config, resolves the selected vehicle runtime plan, connects the
// selected endpoint, starts the central sensor polling runtime, wires selected
// JSONL log subscribers, exposes the dashboard scene boundary, and shuts down
// cleanly when ctx is cancelled.
func Run(ctx context.Context, opts Options) (Summary, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(opts.ConfigPath) == "" {
		return Summary{}, fmt.Errorf("v3 runtime requires a config path")
	}

	cfg, err := v3config.LoadFile(opts.ConfigPath)
	if err != nil {
		return Summary{}, fmt.Errorf("load v3 config: %w", err)
	}
	plan, err := v3config.Resolve(cfg, opts.VehicleID)
	if err != nil {
		return Summary{}, fmt.Errorf("resolve v3 runtime plan: %w", err)
	}

	connector := opts.Connector
	if connector.DialContext == nil && connector.NewSerialReader == nil {
		connector = vehicle.NewConnector()
	}
	reader, err := connector.ConnectPlan(ctx, plan)
	if err != nil {
		return Summary{}, fmt.Errorf("connect v3 endpoint: %w", err)
	}
	defer closeReader(reader)

	pollingRuntime, err := sensors.NewPollingRuntime(reader, plan.Sensors)
	if err != nil {
		return Summary{}, fmt.Errorf("create v3 sensor runtime: %w", err)
	}

	logSubscribers, err := jsonlogger.NewJSONLSubscribersFromPlan(plan)
	if err != nil {
		return Summary{}, fmt.Errorf("create v3 jsonl subscribers: %w", err)
	}
	defer closeSubscribers(logSubscribers)

	dashboardRuntime, err := newDashboardRuntime(plan, opts.RepoRoot)
	if err != nil {
		return Summary{}, err
	}

	summary := Summary{
		VehicleID:      plan.VehicleID,
		VehicleName:    plan.Vehicle.Name,
		Endpoint:       plan.Endpoint.Address,
		SensorCount:    len(plan.Sensors),
		LogCount:       len(plan.Logs),
		DashboardCount: len(plan.Dashboards),
	}

	logger := opts.Logger
	if logger != nil {
		logger.Printf("v3 runtime starting: vehicle=%s endpoint=%s sensors=%d logs=%d dashboards=%d", summary.VehicleID, summary.Endpoint, summary.SensorCount, summary.LogCount, summary.DashboardCount)
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventBuffer := opts.EventBuffer
	if eventBuffer <= 0 {
		eventBuffer = defaultEventBuffer
	}

	var wg sync.WaitGroup
	var firstErr error
	var firstErrMu sync.Mutex
	recordError := func(err error) {
		if err == nil || isContextDone(err) {
			return
		}
		firstErrMu.Lock()
		defer firstErrMu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	for _, subscriber := range logSubscribers {
		subscriber := subscriber
		events := pollingRuntime.Subscribe(eventBuffer)
		wg.Add(1)
		go func() {
			defer wg.Done()
			runJSONLSubscriberDrain(runCtx, cancel, subscriber, events, recordError)
		}()
	}

	if dashboardRuntime != nil && opts.DashboardSink != nil {
		events := pollingRuntime.Subscribe(eventBuffer)
		wg.Add(1)
		go func() {
			defer wg.Done()
			runDashboardSinkDrain(runCtx, cancel, dashboardRuntime, opts.DashboardSink, events, recordError)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		pollingRuntime.Run(runCtx)
	}()

	wg.Wait()
	firstErrMu.Lock()
	defer firstErrMu.Unlock()
	if firstErr != nil {
		return summary, firstErr
	}
	if logger != nil {
		logger.Printf("v3 runtime stopped: vehicle=%s", summary.VehicleID)
	}
	return summary, nil
}

func newDashboardRuntime(plan v3config.RuntimePlan, repoRoot string) (*v3dashboard.Runtime, error) {
	if len(plan.Dashboards) == 0 {
		return nil, nil
	}
	registry, err := v3assets.Load(plan.Assets, repoRoot)
	if err != nil {
		return nil, fmt.Errorf("load v3 dashboard assets: %w", err)
	}
	dashboardRuntime, err := v3dashboard.NewRuntime(plan, registry)
	if err != nil {
		return nil, fmt.Errorf("create v3 dashboard runtime: %w", err)
	}
	return dashboardRuntime, nil
}

func runJSONLSubscriberDrain(ctx context.Context, cancel context.CancelFunc, subscriber *jsonlogger.JSONLSubscriber, events <-chan sensors.SensorEvent, recordError func(error)) {
	processing := true
	for event := range events {
		if !processing || ctx.Err() != nil {
			processing = false
			continue
		}
		if err := subscriber.Handle(event); err != nil {
			recordError(fmt.Errorf("run v3 jsonl subscriber %q: %w", subscriber.ID, err))
			cancel()
			processing = false
		}
	}
}

func runDashboardSinkDrain(ctx context.Context, cancel context.CancelFunc, runtime *v3dashboard.Runtime, sink DashboardSink, events <-chan sensors.SensorEvent, recordError func(error)) {
	processing := true
	for event := range events {
		if !processing || ctx.Err() != nil {
			processing = false
			continue
		}
		scenes, changed, err := runtime.ApplyEvent(event)
		if err != nil {
			recordError(fmt.Errorf("run v3 dashboard boundary: %w", err))
			cancel()
			processing = false
			continue
		}
		if !changed {
			continue
		}
		if err := sink(scenes); err != nil {
			recordError(fmt.Errorf("run v3 dashboard boundary: %w", err))
			cancel()
			processing = false
		}
	}
}

func isContextDone(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func closeReader(reader vehicle.Reader) {
	closer, ok := reader.(io.Closer)
	if !ok || closer == nil {
		return
	}
	_ = closer.Close()
}

func closeSubscribers(subscribers []*jsonlogger.JSONLSubscriber) {
	for _, subscriber := range subscribers {
		_ = subscriber.Close()
	}
}
