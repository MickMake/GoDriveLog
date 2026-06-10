package main

import (
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func TestActiveSensorsForDisplayAddsRaceDemoDashboardSensors(t *testing.T) {
	cfg := config.Config{
		OBD: config.OBDConfig{Provider: config.OBDProviderRaceDemo},
		Sensors: map[string]config.SensorConfig{
			"oil_temperature": {
				Type:    "obd",
				PID:     "015C",
				Unit:    "C",
				Refresh: 250,
				Min:     0,
				Max:     160,
				Log:     true,
				Display: true,
			},
		},
	}

	active := activeSensorsForDisplay(cfg)
	byKey := runtimeSensorsByKey(active)

	want := map[string]string{
		"rpm":               "010C",
		"speed":             "010D",
		"throttle_position": "0111",
		"engine_load":       "DEMO_ENGINE_LOAD",
		"coolant_temp":      "DEMO_COOLANT_TEMP",
		"oil_temperature":   "015C",
		"oil_pressure":      "DEMO_OIL_PRESSURE",
		"gear":              "DEMO_GEAR",
		"warning_level":     "DEMO_WARNING_LEVEL",
		"engine_failed":     "DEMO_ENGINE_FAILED",
		"requires_reset":    "DEMO_REQUIRES_RESET",
		"battery_voltage":   "DEMO_BATTERY",
	}

	for key, rawPID := range want {
		runtimeSensor, ok := byKey[key]
		if !ok {
			t.Fatalf("missing race-demo dashboard sensor %q in %#v", key, active)
		}
		if runtimeSensor.RawPID != rawPID {
			t.Fatalf("race-demo sensor %q RawPID = %q, want %q", key, runtimeSensor.RawPID, rawPID)
		}
		if !runtimeSensor.Display {
			t.Fatalf("race-demo sensor %q Display = false, want true", key)
		}
	}
}

func TestActiveSensorsForDisplayDoesNotChangeOBDSensors(t *testing.T) {
	cfg := config.Config{
		OBD: config.OBDConfig{Provider: config.OBDProviderOBD},
		Sensors: map[string]config.SensorConfig{
			"rpm": {
				Type:    "obd",
				PID:     "010C",
				Unit:    "rpm",
				Refresh: 250,
				Min:     0,
				Max:     7000,
				Log:     true,
				Display: true,
			},
		},
	}

	active := activeSensorsForDisplay(cfg)
	if len(active) != 1 {
		t.Fatalf("len(active) = %d, want 1: %#v", len(active), active)
	}
	if active[0].Key != "rpm" || active[0].RawPID != "010C" {
		t.Fatalf("active[0] = %#v, want rpm 010C", active[0])
	}
}

func TestAppendMissingRaceDemoDisplaySensorsDoesNotDuplicateExistingKeys(t *testing.T) {
	active := []config.RuntimeSensor{
		{Key: "rpm", RawPID: "010C", Unit: "rpm", Refresh: 250, Log: true, Display: true, Min: 0, Max: 7000},
	}

	merged := appendMissingRaceDemoDisplaySensors(active)
	count := 0
	for _, runtimeSensor := range merged {
		if runtimeSensor.Key == "rpm" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("rpm count = %d, want 1: %#v", count, merged)
	}
}

func runtimeSensorsByKey(runtimeSensors []config.RuntimeSensor) map[string]config.RuntimeSensor {
	byKey := make(map[string]config.RuntimeSensor, len(runtimeSensors))
	for _, runtimeSensor := range runtimeSensors {
		byKey[runtimeSensor.Key] = runtimeSensor
	}
	return byKey
}
