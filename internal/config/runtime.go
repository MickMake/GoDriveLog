package config

import (
	"sort"

	"github.com/MickMake/GoDriveLog/internal/state"
)

type RuntimeSensor struct {
	Key     string
	RawPID  string
	Unit    string
	Refresh int
	Log     bool
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
		if !sensor.Log {
			continue
		}

		active = append(active, RuntimeSensor{
			Key:     key,
			RawPID:  sensor.PID,
			Unit:    sensor.Unit,
			Refresh: sensor.Refresh,
			Log:     sensor.Log,
			Min:     sensor.Min,
			Max:     sensor.Max,
		})
	}

	return active
}

func SensorStateDefinitions(sensors []RuntimeSensor) []state.SensorDefinition {
	definitions := make([]state.SensorDefinition, 0, len(sensors))
	for _, sensor := range sensors {
		definitions = append(definitions, state.SensorDefinition{
			ID:   sensor.Key,
			Unit: sensor.Unit,
			Min:  sensor.Min,
			Max:  sensor.Max,
		})
	}
	return definitions
}
