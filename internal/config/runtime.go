package config

import (
	"sort"
	"time"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

const staleAfterRefreshMultiplier = 2

type RuntimeSensor struct {
	Key     string
	RawPID  string
	Unit    string
	Refresh int
	Log     bool
	Display bool
	Min     float64
	Max     float64
}

func ActiveSensors(cfg Config) []RuntimeSensor {
	keys := make([]string, 0, len(cfg.Sensors))
	for key := range cfg.Sensors {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	active := make([]RuntimeSensor, 0, len(keys))
	for _, key := range keys {
		sensor := cfg.Sensors[key]
		if sensor.Type != "obd" {
			continue
		}
		if !sensor.Log && !sensor.Display {
			continue
		}

		active = append(active, RuntimeSensor{
			Key:     key,
			RawPID:  sensor.PID,
			Unit:    sensor.Unit,
			Refresh: sensor.Refresh,
			Log:     sensor.Log,
			Display: sensor.Display,
			Min:     sensor.Min,
			Max:     sensor.Max,
		})
	}

	return active
}

func SensorStateDefinitions(runtimeSensors []RuntimeSensor) []sensors.SensorDefinition {
	definitions := make([]sensors.SensorDefinition, 0, len(runtimeSensors))
	for _, runtimeSensor := range runtimeSensors {
		definitions = append(definitions, sensors.SensorDefinition{
			ID:         runtimeSensor.Key,
			Unit:       runtimeSensor.Unit,
			Min:        runtimeSensor.Min,
			Max:        runtimeSensor.Max,
			StaleAfter: staleAfterForRefresh(runtimeSensor.Refresh),
		})
	}
	return definitions
}

func staleAfterForRefresh(refreshMilliseconds int) time.Duration {
	if refreshMilliseconds <= 0 {
		return 0
	}
	return time.Duration(refreshMilliseconds*staleAfterRefreshMultiplier) * time.Millisecond
}
