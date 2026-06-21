package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	v3fyneadapter "github.com/MickMake/GoDriveLog/internal/dashboard/adapter/fyne"
	v3harness "github.com/MickMake/GoDriveLog/internal/dashboard/harness"
	"github.com/MickMake/GoDriveLog/internal/dashboard/scenesink"
	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	jsonlogger "github.com/MickMake/GoDriveLog/internal/logger"
	v3runtime "github.com/MickMake/GoDriveLog/internal/runtime/v3runtime"
	"github.com/MickMake/GoDriveLog/internal/sensors"
	"github.com/MickMake/GoDriveLog/internal/ui"
)

const v3SceneGap = 12

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	useV3 := flag.Bool("v3", true, "run the v3 selected-vehicle runtime path")
	useHarness := flag.Bool("harness", false, "run the v3 dashboard harness without OBD; requires --v3")
	vehicleID := flag.String("vehicle", "", "v3 vehicle id; required when the v3 config contains multiple vehicles")
	harnessPattern := flag.String("pattern", v3harness.PatternSweep, "v3 dashboard harness pattern: sweep, heartbeat, or fixed")
	harnessInterval := flag.Duration("interval", 100*time.Millisecond, "v3 dashboard harness update interval, such as 50ms or 100ms")
	duration := flag.Duration("duration", 0, "optional v3 runtime or harness duration, such as 60s; zero runs until interrupted")
	renderer := flag.String("renderer", v3RendererFyne, "v3 renderer backend: fyne or ebiten")
	flag.Parse()

	normalizedRenderer, err := normalizeV3Renderer(*renderer)
	if err != nil {
		log.Fatal(err)
	}
	selectedV3Renderer = normalizedRenderer
	selectedV3Duration = *duration

	if *useHarness && !*useV3 {
		log.Fatal("--harness requires --v3")
	}
	if *useV3 {
		if *useHarness {
			if err := runV3HarnessCommand(*configPath, *vehicleID, *harnessPattern, *harnessInterval); err != nil {
				log.Fatal(err)
			}
			return
		}
		if err := runV3Command(*configPath, *vehicleID); err != nil {
			log.Fatal(err)
		}
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	activeSensors := config.ActiveSensors(cfg)
	stateStore := sensors.NewStateStore(config.SensorStateDefinitions(activeSensors))

	logger, err := jsonlogger.NewJSONL(cfg.Log.Directory)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	reader, err := newReader(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.New()
	window := application.NewWindow("GoDriveLog")
	window.Resize(fyne.NewSize(float32(cfg.Dashboard.Canvas.Width), float32(cfg.Dashboard.Canvas.Height)))

	dash, err := ui.NewDashboardWithConfigPath(cfg.Dashboard, *configPath, stateStore)
	if err != nil {
		log.Fatal(err)
	}
	dash.Start(ctx, time.Duration(cfg.Dashboard.RefreshMS)*time.Millisecond)

	lastLogPath := logger.ActivePath()
	status := widget.NewLabel("log: " + lastLogPath)
	var statusMu sync.Mutex

	updateLogStatus := func() {
		path := logger.ActivePath()

		statusMu.Lock()
		if path == lastLogPath {
			statusMu.Unlock()
			return
		}
		lastLogPath = path
		statusMu.Unlock()

		fyne.Do(func() {
			status.SetText("log: " + path)
		})
	}

	content := container.NewBorder(nil, status, nil, nil, dash.CanvasObject())
	window.SetContent(content)

	for _, runtimeSensor := range activeSensors {
		runtimeSensor := runtimeSensor
		go func() {
			ticker := time.NewTicker(time.Duration(runtimeSensor.Refresh) * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					now := time.Now()
					value, unit, err := reader.Read(ctx, runtimeSensor.RawPID)
					if err != nil {
						stateStore.SetError(runtimeSensor.Key, err, now)
						log.Printf("read %s: %v", runtimeSensor.RawPID, err)
						continue
					}

					if runtimeSensor.Unit != "" {
						unit = runtimeSensor.Unit
					}
					stateStore.SetValue(runtimeSensor.Key, value, unit, now)

					reading := sensors.Reading{
						Time:      now,
						SensorKey: runtimeSensor.Key,
						PID:       runtimeSensor.RawPID,
						Name:      runtimeSensor.Key,
						Value:     value,
						Unit:      unit,
						Source:    sourceName(cfg.OBD.MockMode),
					}

					if runtimeSensor.Log {
						if err := logger.Write(reading); err != nil {
							log.Printf("write log: %v", err)
						} else {
							updateLogStatus()
						}
					}
				}
			}
		}()
	}

	window.SetCloseIntercept(func() {
		cancel()
		window.Close()
	})
	window.ShowAndRun()
}

func runV3Command(configPath, vehicleID string) error {
	if selectedV3Renderer == v3RendererEbiten {
		return runV3EbitenCommand(configPath, vehicleID, selectedV3Duration)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if selectedV3Duration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, selectedV3Duration)
		defer cancel()
	}

	initialSize, err := initialV3WindowSize(configPath, vehicleID)
	if err != nil {
		return err
	}
	adapter, err := v3fyneadapter.New(".")
	if err != nil {
		return err
	}

	application := app.New()
	window := application.NewWindow("GoDriveLog v3")
	window.Resize(initialSize)
	window.SetContent(adapter.CanvasObject())
	shutdown := newFyneShutdown(stop, fyne.Do, application.Quit)

	displaySink, err := newFyneSceneSink(func(scenes []v3runtime.Scene) error {
		return adapter.Update(scenes)
	}, "v3 dashboard adapter update")
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		summary, err := v3runtime.Run(ctx, v3runtime.Options{
			ConfigPath:    configPath,
			VehicleID:     vehicleID,
			Logger:        log.Default(),
			DashboardSink: displaySink.SubmitLatest,
		})
		if closeErr := displaySink.Close(); err == nil {
			err = closeErr
		}
		stats := displaySink.Stats()
		if err == nil || isContextStop(err) {
			log.Printf("v3 runtime summary: vehicle=%s endpoint=%s sensors=%d logs=%d dashboards=%d renderer=%s display_submitted=%d display_rendered=%d display_superseded=%d display_last_render=%s", summary.VehicleID, summary.Endpoint, summary.SensorCount, summary.LogCount, summary.DashboardCount, v3RendererFyne, stats.Submitted, stats.Rendered, stats.Superseded, stats.LastRenderDuration)
		} else {
			log.Printf("v3 runtime stopped with error: %v", err)
		}
		errCh <- err
		shutdown.Quit()
	}()
	go func() {
		<-ctx.Done()
		shutdown.Quit()
	}()

	window.SetCloseIntercept(shutdown.CancelAndQuit)
	window.ShowAndRun()
	shutdown.Cancel()
	return ignoreContextStop(<-errCh)
}

func runV3HarnessCommand(configPath, vehicleID, pattern string, interval time.Duration) error {
	if selectedV3Renderer == v3RendererEbiten {
		return runV3EbitenHarnessCommand(configPath, vehicleID, pattern, interval, selectedV3Duration)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if selectedV3Duration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, selectedV3Duration)
		defer cancel()
	}

	initialSize, err := initialV3WindowSize(configPath, vehicleID)
	if err != nil {
		return err
	}
	adapter, err := v3fyneadapter.New(".")
	if err != nil {
		return err
	}

	application := app.New()
	window := application.NewWindow("GoDriveLog v3 harness")
	window.Resize(initialSize)
	window.SetContent(adapter.CanvasObject())
	shutdown := newFyneShutdown(stop, fyne.Do, application.Quit)

	displaySink, err := newFyneSceneSink(func(scenes []v3harness.Scene) error {
		return adapter.Update(scenes)
	}, "v3 dashboard harness adapter update")
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		summary, err := v3harness.Run(ctx, v3harness.Options{
			ConfigPath: configPath,
			VehicleID:  vehicleID,
			Pattern:    pattern,
			Interval:   interval,
			Logger:     log.Default(),
			Sink:       displaySink.SubmitLatest,
		})
		if closeErr := displaySink.Close(); err == nil {
			err = closeErr
		}
		stats := displaySink.Stats()
		if err == nil || isContextStop(err) {
			log.Printf("v3 dashboard harness summary: vehicle=%s sensors=%d dashboards=%d pattern=%s interval=%s renderer=%s events=%d display_submitted=%d display_rendered=%d display_superseded=%d display_last_render=%s", summary.VehicleID, summary.SensorCount, summary.DashboardCount, summary.Pattern, summary.Interval, v3RendererFyne, summary.Events, stats.Submitted, stats.Rendered, stats.Superseded, stats.LastRenderDuration)
		} else {
			log.Printf("v3 dashboard harness stopped with error: %v", err)
		}
		errCh <- err
		shutdown.Quit()
	}()
	go func() {
		<-ctx.Done()
		shutdown.Quit()
	}()

	window.SetCloseIntercept(shutdown.CancelAndQuit)
	window.ShowAndRun()
	shutdown.Cancel()
	return ignoreContextStop(<-errCh)
}

func newFyneSceneSink(update scenesink.Sink, label string) (*scenesink.LatestSink, error) {
	return scenesink.NewLatestSink(func(scenes []v3runtime.Scene) error {
		var updateErr error
		fyne.DoAndWait(func() {
			updateErr = update(scenes)
		})
		if updateErr == nil {
			log.Printf("%s: scenes=%d", label, len(scenes))
		}
		return updateErr
	})
}

func newReader(cfg config.Config) (sensors.Reader, error) {
	if !cfg.OBD.MockMode {
		return nil, fmt.Errorf("non-mock OBD reader is not available in this command path: %s", cfg.OBD.Address)
	}
	return sensors.NewMockReader(), nil
}

func sourceName(mockMode bool) string {
	if mockMode {
		return "mock"
	}
	return "obd"
}

type fyneShutdown struct {
	cancel     func()
	schedule   func(func())
	quit       func()
	cancelOnce sync.Once
	quitOnce   sync.Once
}

func newFyneShutdown(cancel func(), schedule func(func()), quit func()) *fyneShutdown {
	return &fyneShutdown{cancel: cancel, schedule: schedule, quit: quit}
}

func (s *fyneShutdown) Cancel() {
	if s == nil {
		return
	}
	s.cancelOnce.Do(func() {
		if s.cancel != nil {
			s.cancel()
		}
	})
}

func (s *fyneShutdown) Quit() {
	if s == nil {
		return
	}
	s.quitOnce.Do(func() {
		if s.schedule != nil && s.quit != nil {
			s.schedule(s.quit)
		}
	})
}

func (s *fyneShutdown) CancelAndQuit() {
	s.Cancel()
	s.Quit()
}

func initialV3WindowSize(configPath, vehicleID string) (fyne.Size, error) {
	cfg, err := v3config.LoadFile(configPath)
	if err != nil {
		return fyne.Size{}, fmt.Errorf("load v3 config for initial window size: %w", err)
	}
	plan, err := v3config.Resolve(cfg, vehicleID)
	if err != nil {
		return fyne.Size{}, fmt.Errorf("resolve v3 runtime plan for initial window size: %w", err)
	}
	return selectedDashboardsSize(plan.Dashboards), nil
}

func selectedDashboardsSize(dashboards []v3config.ResolvedDashboard) fyne.Size {
	var width float32
	var height float32
	for index, dashboard := range dashboards {
		if float32(dashboard.Config.Size.Width) > width {
			width = float32(dashboard.Config.Size.Width)
		}
		height += float32(dashboard.Config.Size.Height)
		if index < len(dashboards)-1 {
			height += v3SceneGap
		}
	}
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 480
	}
	return fyne.NewSize(width, height)
}
