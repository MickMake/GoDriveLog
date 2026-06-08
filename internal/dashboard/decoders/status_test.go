package decoders

import (
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestExecutePropagatesErroredSensorStatusThroughDerivedValues(t *testing.T) {
	values, err := Execute([]config.DashboardDecoderConfig{
		{ID: "rpm_text", Type: config.DashboardDecoderFormatNumber, Sensor: "rpm", Format: "0000"},
		{ID: "rpm_digits", Type: config.DashboardDecoderDigits, Input: "rpm_text"},
		{ID: "rpm_warning", Type: config.DashboardDecoderThreshold, Sensor: "rpm", Thresholds: []config.ThresholdConfig{{At: 0, Value: "normal"}, {At: 2100, Value: "warning"}}},
	}, Inputs{Sensors: map[string]sensors.SensorState{
		"rpm": {ID: "rpm", Value: 2400, Min: 0, Max: 7000, Status: sensors.StatusError, Error: "read failed"},
	}})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	for _, id := range []string{"rpm_text", "rpm_digits", "rpm_warning"} {
		if values[id].Status != sensors.StatusError {
			t.Fatalf("%s status = %q, want %q", id, values[id].Status, sensors.StatusError)
		}
		if values[id].Error != "read failed" {
			t.Fatalf("%s error = %q, want read failed", id, values[id].Error)
		}
	}
	if values["rpm_warning"].Text != "warning" {
		t.Fatalf("rpm_warning = %q, want warning", values["rpm_warning"].Text)
	}
}
