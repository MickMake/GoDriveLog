package main

import (
	"context"
	"flag"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"GoDriveLog/internal/config"
	jsonlogger "GoDriveLog/internal/logger"
	"GoDriveLog/internal/sensors"
	"GoDriveLog/internal/ui"
)

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	activePIDs := config.ActivePIDs(cfg)

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
	window.Resize(fyne.NewSize(800, 480))

	dash := ui.NewDashboard(cfg.Vehicle.PIDs)
	status := widget.NewLabel("log: " + logger.ActivePath())
	content := container.NewBorder(nil, status, nil, nil, dash.CanvasObject())
	window.SetContent(content)

	for _, runtimePID := range activePIDs {
		runtimePID := runtimePID
		go func() {
			ticker := time.NewTicker(time.Duration(runtimePID.Refresh) * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					value, unit, err := reader.Read(ctx, runtimePID.RawPID)
					if err != nil {
						log.Printf("read %s: %v", runtimePID.RawPID, err)
						if runtimePID.Display.Enabled {
							fyne.Do(func() { dash.SetError(runtimePID.RawPID, err) })
						}
						continue
					}

					if runtimePID.Unit != "" {
						unit = runtimePID.Unit
					}

					reading := sensors.Reading{
						Time:   time.Now(),
						Key:    runtimePID.Key,
						PID:    runtimePID.RawPID,
						Name:   runtimePID.Key,
						Value:  value,
						Unit:   unit,
						Source: sourceName(cfg.MockMode),
					}

					if runtimePID.Log {
						if err := logger.Write(reading); err != nil {
							log.Printf("write log: %v", err)
						}
					}

					if runtimePID.Display.Enabled {
						fyne.Do(func() { dash.Update(reading) })
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
	if cfg.MockMode {
		return sensors.NewMockReader(), nil
	}
	return sensors.NewELMOBDReader(cfg.OBDAddress, cfg.OBDDebug)
}

func sourceName(mock bool) string {
	if mock {
		return "mock"
	}
	return "obd"
}
