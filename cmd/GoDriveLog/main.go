package main

import (
	"context"
	"flag"
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
	"github.com/MickMake/GoDriveLog/internal/config"
	jsonlogger "github.com/MickMake/GoDriveLog/internal/logger"
	v3runtime "github.com/MickMake/GoDriveLog/internal/runtime/v3runtime"
	"github.com/MickMake/GoDriveLog/internal/sensors"
	"github.com/MickMake/GoDriveLog/internal/ui"
)

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	useV3 := flag.Bool("v3", false, "run the v3 selected-vehicle runtime path")
	vehicleID := flag.String("vehicle", "", "v3 vehicle id; required when the v3 config contains multiple vehicles")
	repoRoot := flag.String("repo-root", ".", "repository root for v3 dashboard assets")
	flag.Parse()

	if *useV3 {
		if err := runV3Command(*configPath, *vehicleID, *repoRoot); err != nil {
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

func runV3Command(configPath, vehicleID, repoRoot string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	adapter, err := v3fyneadapter.New(repoRoot)
	if err != nil {
		return err
	}

	application := app.New()
	window := application.NewWindow("GoDriveLog v3")
	window.Resize(fyne.NewSize(800, 480))
	window.SetContent(adapter.CanvasObject())

	errCh := make(chan error, 1)
	go func() {
		summary, err := v3runtime.Run(ctx, v3runtime.Options{
			ConfigPath: configPath,
			VehicleID:  vehicleID,
			RepoRoot:   repoRoot,
			Logger:     log.Default(),
			DashboardSink: func(scenes []v3runtime.Scene) error {
				var updateErr error
				fyne.DoAndWait(func() {
					updateErr = adapter.Update(scenes)
				})
				if updateErr == nil {
					log.Printf("v3 dashboard adapter update: scenes=%d", len(scenes))
				}
				return updateErr
			},
		})
		if err == nil {
			log.Printf("v3 runtime summary: vehicle=%s endpoint=%s sensors=%d logs=%d dashboards=%d", summary.VehicleID, summary.Endpoint, summary.SensorCount, summary.LogCount, summary.DashboardCount)
		}
		errCh <- err
		fyne.Do(func() {
			window.Close()
		})
	}()

	window.SetCloseIntercept(func() {
		stop()
		window.Close()
	})
	window.ShowAndRun()
	stop()
	return <-errCh
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
