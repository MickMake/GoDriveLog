package main

import (
	"context"
	"flag"
	"log"
	"sync/atomic"
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

type configuredPID struct {
	Key string
	PID config.PIDConfig
}

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

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

	var engineWasRunning atomic.Bool
	for _, item := range pollingPIDs(cfg) {
		item := item
		go func() {
			ticker := time.NewTicker(time.Duration(item.PID.Refresh) * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					value, unit, err := reader.Read(ctx, item.PID.PID)
					if err != nil {
						log.Printf("read %s: %v", item.PID.PID, err)
						fyne.Do(func() { dash.SetError(item.PID.PID, err) })
						continue
					}

					if item.PID.Unit != "" {
						unit = item.PID.Unit
					}

					reading := sensors.Reading{
						Time:   time.Now(),
						PID:    item.PID.PID,
						Name:   item.Key,
						Value:  value,
						Unit:   unit,
						Source: sourceName(cfg.MockMode),
					}

					// TODO: daily rotation replaces engine-start rotation in version 1.5.
					if item.PID.PID == "010C" {
						running := value >= 50
						if running && !engineWasRunning.Load() {
							if err := logger.Rotate("engine-start"); err != nil {
								log.Printf("rotate log: %v", err)
							} else {
								fyne.Do(func() { status.SetText("log: " + logger.ActivePath()) })
							}
						}
						engineWasRunning.Store(running)
					}

					if err := logger.Write(reading); err != nil {
						log.Printf("write log: %v", err)
					}

					fyne.Do(func() { dash.Update(reading) })
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

func pollingPIDs(cfg config.Config) []configuredPID {
	items := make([]configuredPID, 0, len(cfg.Vehicle.PIDs))
	for key, pid := range cfg.Vehicle.PIDs {
		if pid.Type != "obd" {
			continue
		}
		items = append(items, configuredPID{Key: key, PID: pid})
	}
	return items
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
