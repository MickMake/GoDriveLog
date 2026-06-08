package decoders

import (
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestExecuteDecoderTypes(t *testing.T) {
	inputs := Inputs{Sensors: map[string]sensors.SensorState{
		"rpm":      {ID: "rpm", Value: 3500, Unit: "rpm", Min: 0, Max: 7000, Status: sensors.StatusOK},
		"throttle": {ID: "throttle", Value: 0.50, Unit: "%", Min: 0, Max: 1, Status: sensors.StatusOK},
		"warning":  {ID: "warning", Value: 1, Min: 0, Max: 1, Status: sensors.StatusOK},
	}}

	configs := []config.DashboardDecoderConfig{
		{ID: "rpm_norm", Type: config.DashboardDecoderNormalize, Sensor: "rpm"},
		{ID: "rpm_zone", Type: config.DashboardDecoderThreshold, Sensor: "rpm", Thresholds: []config.ThresholdConfig{{At: 0, Value: "low"}, {At: 3000, Value: "mid"}, {At: 6000, Value: "high"}}},
		{ID: "throttle_frame", Type: config.DashboardDecoderFrameIndex, Sensor: "throttle", FrameCount: 11},
		{ID: "rpm_text", Type: config.DashboardDecoderFormatNumber, Sensor: "rpm", Format: "%.0f"},
		{ID: "rpm_digits", Type: config.DashboardDecoderDigits, Input: "rpm_text", Format: "%.0f"},
		{ID: "warning_bool", Type: config.DashboardDecoderBoolean, Sensor: "warning"},
	}

	values, err := Execute(configs, inputs)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if got := values["rpm_norm"].Number; got != 0.5 {
		t.Fatalf("rpm_norm = %v, want 0.5", got)
	}
	if got := values["rpm_zone"].Text; got != "mid" {
		t.Fatalf("rpm_zone = %q, want mid", got)
	}
	if got := values["throttle_frame"].FrameIndex; got != 5 {
		t.Fatalf("throttle_frame = %d, want 5", got)
	}
	if got := values["rpm_text"].Text; got != "3500" {
		t.Fatalf("rpm_text = %q, want 3500", got)
	}
	if got := values["rpm_digits"].Digits; len(got) != 4 || got[0] != "3" || got[3] != "0" {
		t.Fatalf("rpm_digits = %#v, want [3 5 0 0]", got)
	}
	if got := values["warning_bool"].Bool; !got {
		t.Fatalf("warning_bool = false, want true")
	}
}

func TestDecoderErrors(t *testing.T) {
	tests := []struct {
		name    string
		decoder config.DashboardDecoderConfig
		inputs  Inputs
	}{
		{
			name:    "unknown sensor",
			decoder: config.DashboardDecoderConfig{ID: "missing", Type: config.DashboardDecoderNormalize, Sensor: "missing"},
			inputs:  Inputs{Sensors: map[string]sensors.SensorState{}},
		},
		{
			name:    "sensor error",
			decoder: config.DashboardDecoderConfig{ID: "bad", Type: config.DashboardDecoderNormalize, Sensor: "rpm"},
			inputs:  Inputs{Sensors: map[string]sensors.SensorState{"rpm": {ID: "rpm", Status: sensors.StatusError, Error: "read failed"}}},
		},
		{
			name:    "invalid normalize range",
			decoder: config.DashboardDecoderConfig{ID: "flat", Type: config.DashboardDecoderNormalize, Sensor: "rpm"},
			inputs:  Inputs{Sensors: map[string]sensors.SensorState{"rpm": {ID: "rpm", Value: 1, Min: 1, Max: 1, Status: sensors.StatusOK}}},
		},
		{
			name:    "invalid frame count",
			decoder: config.DashboardDecoderConfig{ID: "frame", Type: config.DashboardDecoderFrameIndex, Sensor: "rpm"},
			inputs:  Inputs{Sensors: map[string]sensors.SensorState{"rpm": {ID: "rpm", Value: 1, Min: 0, Max: 1, Status: sensors.StatusOK}}},
		},
		{
			name:    "digits rejects formatted decimal",
			decoder: config.DashboardDecoderConfig{ID: "digits", Type: config.DashboardDecoderDigits, Sensor: "rpm", Format: "%.1f"},
			inputs:  Inputs{Sensors: map[string]sensors.SensorState{"rpm": {ID: "rpm", Value: 12.3, Min: 0, Max: 100, Status: sensors.StatusOK}}},
		},
		{
			name:    "unknown prior input",
			decoder: config.DashboardDecoderConfig{ID: "derived", Type: config.DashboardDecoderBoolean, Input: "missing"},
			inputs:  Inputs{Values: map[string]Value{}},
		},
	}

	registry := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := registry.Decode(tt.decoder, tt.inputs); err == nil {
				t.Fatal("Decode returned nil error, want error")
			}
		})
	}
}
