package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/MickMake/GoDriveLog/internal/config"
	jsonlogger "github.com/MickMake/GoDriveLog/internal/logger"
	"github.com/MickMake/GoDriveLog/internal/sensors"
	"github.com/MickMake/GoDriveLog/internal/ui"
)

const instrumentRefreshMS = 50

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	providerOverride := flag.String("sensor-provider", "", "sensor provider override: obd, mock, or race-demo")
	debugStrip := flag.Bool("debug-strip", false, "show machine-readable dashboard debug strip")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if *providerOverride != "" {
		cfg.OBD.Provider = config.NormalizeOBDProvider(*providerOverride, cfg.OBD.MockMode)
	}
	activeSensors := activeSensorsForDisplay(cfg)
	if *debugStrip {
		printDebugActiveSensors(activeSensors)
	}
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
	window.Resize(fyne.NewSize(1920, 480))

	dash, err := ui.NewInstrumentDashboard1920x480WithPNGDigitConfig(stateStore, ui.InstrumentDashboardOptions{
		DebugStrip: *debugStrip,
		DebugSource: sourceName(cfg.OBD.Provider),
		DebugPIDs: runtimeSensorPIDMap(activeSensors),
	}, pngDigitConfigFromDashboard(cfg))
	if err != nil {
		log.Fatal(err)
	}
	dash.Start(ctx, instrumentRefreshMS*time.Millisecond)
	window.SetContent(dash.CanvasObject())

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
						if *debugStrip {
							printDebugSensorRead(runtimeSensor, 0, runtimeSensor.Unit, "error", now)
						}
						continue
					}

					if runtimeSensor.Unit != "" {
						unit = runtimeSensor.Unit
					}
					stateStore.SetValue(runtimeSensor.Key, value, unit, now)
					if *debugStrip {
						printDebugSensorRead(runtimeSensor, value, unit, "ok", now)
					}

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

func pngDigitConfigFromDashboard(cfg config.Config) ui.PNGDigitConfig {
	for _, asset := range cfg.Dashboard.Assets {
		if asset.Type != config.DashboardAssetPNGDigitSet {
			continue
		}
		return ui.PNGDigitConfig{
			AssetRoot: cfg.Dashboard.AssetRoot,
			Glyphs:    asset.Glyphs,
		}
	}
	return ui.PNGDigitConfig{}
}

func runtimeSensorPIDMap(activeSensors []config.RuntimeSensor) map[string]string {
	mapped := make(map[string]string, len(activeSensors))
	for _, runtimeSensor := range activeSensors {
		mapped[runtimeSensor.Key] = runtimeSensor.RawPID
	}
	return mapped
}

func printDebugActiveSensors(activeSensors []config.RuntimeSensor) {
	for _, runtimeSensor := range activeSensors {
		fmt.Fprintf(os.Stderr, "GDLDBG_ACTIVE|key=%s|pid=%s|display=%t|log=%t|refresh_ms=%d\n", runtimeSensor.Key, runtimeSensor.RawPID, runtimeSensor.Display, runtimeSensor.Log, runtimeSensor.Refresh)
	}
}

func printDebugSensorRead(runtimeSensor config.RuntimeSensor, value float64, unit string, status string, updatedAt time.Time) {
	fmt.Fprintf(os.Stderr, "GDLDBG_READ|key=%s|pid=%s|value=%.2f|unit=%s|status=%s|unix_ms=%d\n", runtimeSensor.Key, runtimeSensor.RawPID, value, unit, status, updatedAt.UnixMilli())
}

func activeSensorsForDisplay(cfg config.Config) []config.RuntimeSensor {
	activeSensors := config.ActiveSensors(cfg)
	if cfg.OBD.Provider != config.OBDProviderRaceDemo {
		return activeSensors
	}

	return appendMissingRaceDemoDisplaySensors(activeSensors)
}

func appendMissingRaceDemoDisplaySensors(activeSensors []config.RuntimeSensor) []config.RuntimeSensor {
	seen := make(map[string]bool, len(activeSensors))
	for _, runtimeSensor := range activeSensors {
		seen[runtimeSensor.Key] = true
	}

	for _, runtimeSensor := range []config.RuntimeSensor{
		{Key: "rpm", RawPID: "010C", Unit: "rpm", Refresh: 250, Log: true, Display: true, Min: 0, Max: 7000},
		{Key: "speed", RawPID: "010D", Unit: "km/h", Refresh: 250, Log: true, Display: true, Min: 0, Max: 160},
		{Key: "throttle_position", RawPID: "0111", Unit: "%", Refresh: 250, Log: true, Display: true, Min: 0, Max: 100},
		{Key: "engine_load", RawPID: "DEMO_ENGINE_LOAD", Unit: "%", Refresh: 250, Log: true, Display: true, Min: 0, Max: 100},
		{Key: "coolant_temp", RawPID: "DEMO_COOLANT_TEMP", Unit: "C", Refresh: 250, Log: true, Display: true, Min: -40, Max: 140},
		{Key: "battery_voltage", RawPID: "DEMO_BATTERY", Unit: "V", Refresh: 500, Log: true, Display: true, Min: 0, Max: 16},
		{Key: "oil_temperature", RawPID: "DEMO_OIL_TEMP", Unit: "C", Refresh: 250, Log: true, Display: true, Min: 0, Max: 160},
		{Key: "oil_pressure", RawPID: "DEMO_OIL_PRESSURE", Unit: "kPa", Refresh: 250, Log: true, Display: true, Min: 0, Max: 500},
		{Key: "gear", RawPID: "DEMO_GEAR", Unit: "gear", Refresh: 250, Log: true, Display: true, Min: 0, Max: 6},
		{Key: "warning_level", RawPID: "DEMO_WARNING_LEVEL", Unit: "level", Refresh: 250, Log: true, Display: true, Min: 0, Max: 2},
		{Key: "engine_failed", RawPID: "DEMO_ENGINE_FAILED", Unit: "bool", Refresh: 250, Log: true, Display: true, Min: 0, Max: 1},
		{Key: "requires_reset", RawPID: "DEMO_REQUIRES_RESET", Unit: "bool", Refresh: 250, Log: true, Display: true, Min: 0, Max: 1},
	} {
		if seen[runtimeSensor.Key] {
			continue
		}
		activeSensors = append(activeSensors, runtimeSensor)
		seen[runtimeSensor.Key] = true
	}

	return activeSensors
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
