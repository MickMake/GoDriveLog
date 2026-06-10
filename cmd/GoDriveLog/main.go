package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/internal/config"
	jsonlogger "github.com/MickMake/GoDriveLog/internal/logger"
	"github.com/MickMake/GoDriveLog/internal/sensors"
	"github.com/MickMake/GoDriveLog/internal/ui"
)

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	providerOverride := flag.String("sensor-provider", "", "sensor provider override: obd, mock, or race-demo")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if *providerOverride != "" {
		cfg.OBD.Provider = config.NormalizeOBDProvider(*providerOverride, cfg.OBD.MockMode)
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
						Source:    sourceName(cfg.OBD.Provider),
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

func newReader(cfg config.Config) (sensors.Reader, error) {
	switch cfg.OBD.Provider {
	case config.OBDProviderMock:
		return sensors.NewMockReader(), nil
	case config.OBDProviderRaceDemo:
		return sensors.NewRaceDemoReader(), nil
	case config.OBDProviderOBD:
		return sensors.NewELMOBDReader(cfg.OBD.Address, cfg.OBD.Debug)
	default:
		return nil, fmt.Errorf("unsupported sensor provider %q", cfg.OBD.Provider)
	}
}

func sourceName(provider string) string {
	return provider
}
