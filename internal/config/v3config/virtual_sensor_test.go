package v3config

import "testing"

func TestAllowVirtualSensorForHarnessConfigs(t *testing.T) {
	cfg := validMinimalConfig()
	cfg.Sensors["speed"] = SensorConfig{Type: SensorTypeVirtual, Unit: "km/h", Poll: 250, Min: floatPtr(0), Max: floatPtr(220)}
	if err := Validate(cfg); err != nil {
		t.Fatalf("expected virtual sensor config to validate: %v", err)
	}
	if got := SensorEffectiveValueKind(cfg.Sensors["speed"]); got != ValueKindNumeric {
		t.Fatalf("effective value kind = %q, want numeric", got)
	}
}

func TestRejectVirtualSensorPID(t *testing.T) {
	cfg := validMinimalConfig()
	cfg.Sensors["speed"] = SensorConfig{Type: SensorTypeVirtual, PID: "bad", Unit: "km/h", Poll: 250, Min: floatPtr(0), Max: floatPtr(220)}
	if err := Validate(cfg); err == nil {
		t.Fatalf("expected virtual sensor pid to fail")
	} else {
		assertErrorContains(t, err, "pid must be empty for virtual sensors")
	}
}

func floatPtr(value float64) *float64 {
	return &value
}
