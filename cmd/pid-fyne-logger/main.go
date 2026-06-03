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

	"pid-fyne-logger/internal/config"
	jsonlogger "pid-fyne-logger/internal/logger"
	"pid-fyne-logger/internal/sensors"
	"pid-fyne-logger/internal/ui"
)

func main() {
	configPath := flag.String("config", "config.example.json", "path to JSON config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := jsonlogger.NewJSONL(cfg.LogDir)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	var reader sensors.Reader
	if cfg.MockMode {
		reader = sensors.NewMockReader()
	} else {
		reader = &sensors.ELM327Reader{}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.New()
	window := application.NewWindow("PID Logger")
	window.Resize(fyne.NewSize(800, 480))

	dash := ui.NewDashboard(cfg.Sensors)
	status := widget.NewLabel("log: " + logger.ActivePath())
	content := container.NewBorder(nil, status, nil, nil, dash.CanvasObject())
	window.SetContent(content)

	var engineWasRunning atomic.Bool
	for _, sc := range cfg.Sensors {
		sc := sc
		go func() {
			ticker := time.NewTicker(time.Duration(sc.RefreshMS) * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					value, unit, err := reader.Read(ctx, sc.PID)
					if err != nil {
						log.Printf("read %s: %v", sc.PID, err)
						continue
					}

					reading := sensors.Reading{
						Time:   time.Now(),
						PID:    sc.PID,
						Name:   sc.Name,
						Value:  value,
						Unit:   unit,
						Source: sourceName(cfg.MockMode),
					}

					if sc.PID == cfg.EngineStartPID {
						running := value >= cfg.EngineStartThreshold
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

func sourceName(mock bool) string {
	if mock {
		return "mock"
	}
	return "obd"
}
