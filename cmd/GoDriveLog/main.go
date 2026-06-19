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
	flag.Parse()

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
			DashboardSink: displaySink.Submit,
		})
		if closeErr := displaySink.Close(); err == nil {
			err = closeErr
		}
		if err == nil {
			log.Printf("v3 runtime summary: vehicle=%s endpoint=%s sensors=%d logs=%d dashboards=%d", summary.VehicleID, summary.Endpoint, summary.SensorCount, summary.LogCount, summary.DashboardCount)
		} else {
			log.Printf("v3 runtime stopped with error: %v", err)
		}
		errCh <- err
	}()

	window.SetCloseIntercept(func() {
		stop()
		application.Quit()
	})
	window.ShowAndRun()
	stop()
	return <-errCh
}

func runV3HarnessCommand(configPath, vehicleID, pattern string, interval time.Duration) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
			Sink:       displaySink.Submit,
		})
		if closeErr := displaySink.Close(); err == nil {
			err = closeErr
		}
		if err == nil {
			log.Printf("v3 dashboard harness summary: vehicle=%s sensors=%d dashboards=%d pattern=%s interval=%s events=%d", summary.VehicleID, summary.SensorCount, summary.DashboardCount, summary.Pattern, summary.Interval, summary.Events)
		} else {
			log.Printf("v3 dashboard harness stopped with error: %v", err)
		}
		errCh <- err
	}()

	window.SetCloseIntercept(func() {
		stop()
		application.Quit()
	})
	window.ShowAndRun()
	stop()
	return <-errCh
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

func newReader(cfg config.Config) (sensors.Reader, error) {
	if cfg.OBD.MockMode {
		return sensors.NewMockReader(), nil
	}
	return sensors.NewELMOBDReader(cfg.OBD.Address, cfg.OBD.Debug)
}

func sourceName(mock bool) string {
	if mock {
		return "mock"
	}
	return "obd"
}
