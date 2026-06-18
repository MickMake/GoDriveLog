package v3dashboard

import (
	"testing"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestDashboardMissingSensorStateUsesMissingStatus(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{0, 0}, Digits: 3}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "speed")
	if widget.Status != sensors.StatusMissing {
		t.Fatalf("Status = %q, want %q", widget.Status, sensors.StatusMissing)
	}
	if countParts(widget, PartKindCharacter) != 0 {
		t.Fatalf("missing sensor rendered live digits: %#v", widget.Parts)
	}
}

func TestDashboardUnavailableStatusesDoNotRenderLiveValues(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{name: "missing", status: sensors.StatusMissing},
		{name: "unsupported", status: sensors.StatusUnsupported},
		{name: "timeout", status: sensors.StatusTimeout},
		{name: "parse", status: sensors.StatusParseError},
		{name: "error", status: sensors.StatusError},
		{name: "stale", status: sensors.StatusStale},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "warning", Type: v3config.WidgetTypeIndicator, Sensor: "warning", Asset: "warning", Position: []int{0, 0}}}}}
			scene, err := dashboard.Render(map[string]sensors.SensorState{"warning": {ID: "warning", Value: 1, Status: test.status, Error: "not available"}})
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}
			widget := requireWidget(t, scene, "warning")
			if widget.Status != test.status || widget.Error != "not available" {
				t.Fatalf("widget = %#v, want status %q and error", widget, test.status)
			}
			if got := statePart(widget); got != v3assets.IndicatorStateUnknown {
				t.Fatalf("indicator state = %q, want unknown", got)
			}
		})
	}
}
