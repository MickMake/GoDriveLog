package config

import "testing"

func TestActiveSensorsSeparatesLogAndDisplay(t *testing.T) {
	cfg := Config{Sensors: map[string]SensorConfig{
		"no_log_no_display": {
			Type:    "obd",
			PID:     "0104",
			Unit:    "%",
			Refresh: 500,
			Min:     0,
			Max:     100,
			Log:     false,
			Display: false,
		},
		"log_no_display": {
			Type:    "obd",
			PID:     "010C",
			Unit:    "rpm",
			Refresh: 250,
			Min:     0,
			Max:     7000,
			Log:     true,
			Display: false,
		},
		"no_log_display": {
			Type:    "obd",
			PID:     "010D",
			Unit:    "km/h",
			Refresh: 500,
			Min:     0,
			Max:     160,
			Log:     false,
			Display: true,
		},
		"log_display": {
			Type:    "obd",
			PID:     "0111",
			Unit:    "%",
			Refresh: 500,
			Min:     0,
			Max:     100,
			Log:     true,
			Display: true,
		},
		"virtual": {
			Type:    "virtual",
			Unit:    "bool",
			Refresh: 1000,
			Min:     0,
			Max:     1,
			Log:     true,
			Display: true,
		},
	}}

	active := ActiveSensors(cfg)
	want := map[string]RuntimeSensor{
		"log_no_display": {Log: true, Display: false},
		"no_log_display": {Log: false, Display: true},
		"log_display":    {Log: true, Display: true},
	}

	if len(active) != len(want) {
		t.Fatalf("len(active) = %d, want %d: %#v", len(active), len(want), active)
	}
	for _, runtimeSensor := range active {
		wantSensor, ok := want[runtimeSensor.Key]
		if !ok {
			t.Fatalf("unexpected active sensor %q in %#v", runtimeSensor.Key, active)
		}
		if runtimeSensor.Log != wantSensor.Log || runtimeSensor.Display != wantSensor.Display {
			t.Fatalf("active sensor %q flags = log:%v display:%v, want log:%v display:%v", runtimeSensor.Key, runtimeSensor.Log, runtimeSensor.Display, wantSensor.Log, wantSensor.Display)
		}
		delete(want, runtimeSensor.Key)
	}
	if len(want) != 0 {
		t.Fatalf("missing active sensors: %#v", want)
	}
}

func TestSensorStateDefinitionsOnlyIncludesDisplayedSensors(t *testing.T) {
	runtimeSensors := []RuntimeSensor{
		{Key: "log_no_display", Unit: "rpm", Refresh: 250, Log: true, Display: false, Min: 0, Max: 7000},
		{Key: "no_log_display", Unit: "km/h", Refresh: 500, Log: false, Display: true, Min: 0, Max: 160},
		{Key: "log_display", Unit: "%", Refresh: 500, Log: true, Display: true, Min: 0, Max: 100},
	}

	definitions := SensorStateDefinitions(runtimeSensors)
	if len(definitions) != 2 {
		t.Fatalf("len(definitions) = %d, want 2: %#v", len(definitions), definitions)
	}
	seen := map[string]bool{}
	for _, definition := range definitions {
		seen[definition.ID] = true
	}
	if seen["log_no_display"] {
		t.Fatal("log-only sensor was included in display state definitions")
	}
	if !seen["no_log_display"] || !seen["log_display"] {
		t.Fatalf("display sensors missing from definitions: %#v", seen)
	}
}
