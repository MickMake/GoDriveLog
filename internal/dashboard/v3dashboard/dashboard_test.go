package v3dashboard

import (
	"strings"
	"testing"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestBarDisplayKeepsOKBehaviour(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "temperature", Type: v3config.WidgetTypeBarDisplay, Sensor: "temperature", Asset: "temperature_bar", Position: []int{0, 0}, Cells: 5, Min: floatPtr(0), Max: floatPtr(100), Zones: []v3config.ZoneConfig{{UpTo: 70, Cell: "on"}, {UpTo: 100, Cell: "warning"}}}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"temperature": okState("temperature", 80, "c")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "temperature")
	if got := cellNames(widget); got != "warning,warning,warning,warning,off" {
		t.Fatalf("expected zone-selected warning cells, got %q", got)
	}
}

func TestBarDisplayRendersNoCellsForNonOKStatus(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "temperature", Type: v3config.WidgetTypeBarDisplay, Sensor: "temperature", Asset: "temperature_bar", Position: []int{0, 0}, Cells: 5, Min: floatPtr(0), Max: floatPtr(100)}}}}
	statuses := []string{
		sensors.StatusMissing,
		sensors.StatusUnsupported,
		sensors.StatusTimeout,
		sensors.StatusParseError,
		sensors.StatusError,
		sensors.StatusStale,
	}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			scene, err := dashboard.Render(map[string]sensors.SensorState{"temperature": {ID: "temperature", Value: 0, Status: status, Error: "not available"}})
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}
			widget := requireWidget(t, scene, "temperature")
			if widget.Status != status {
				t.Fatalf("Status = %q, want %q", widget.Status, status)
			}
			if got := countParts(widget, PartKindCell); got != 0 {
				t.Fatalf("expected non-ok bar to render no cell parts, got %d from %#v", got, widget.Parts)
			}
		})
	}
}

func TestIndicatorUsesUnknownForNonOKStatus(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "warning", Type: v3config.WidgetTypeIndicator, Sensor: "warning", Asset: "warning", Position: []int{0, 0}}}}}
	scene, err := dashboard.Render(map[string]sensors.SensorState{"warning": {ID: "warning", Value: 1, Status: sensors.StatusStale}})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	warning := requireWidget(t, scene, "warning")
	if statePart(warning) != v3assets.IndicatorStateUnknown {
		t.Fatalf("expected stale indicator to use unknown state, got %#v", warning)
	}
}

func testAssetRegistry() *v3assets.Registry {
	characters := map[string]v3assets.ImageAsset{}
	for _, ch := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"} {
		characters[ch] = v3assets.ImageAsset{Path: "assets/digits/" + ch + ".png"}
	}
	characters["-"] = v3assets.ImageAsset{Path: "assets/digits/minus.png"}

	frames := map[int]v3assets.ImageAsset{}
	for frame := 0; frame <= 4; frame++ {
		frames[frame] = v3assets.ImageAsset{Path: "assets/throttle/frame_00" + string(rune('0'+frame)) + ".png"}
	}

	return &v3assets.Registry{
		Images: map[string]v3assets.ImageSet{"panel": {ID: "panel", Image: &v3assets.ImageAsset{Path: "assets/panel.png"}}},
		DigitSets: map[string]v3assets.DigitSet{"digits": {ID: "digits", Background: &v3assets.ImageAsset{Path: "assets/digits/back.png"}, Characters: characters, DecimalPoint: &v3assets.ImageAsset{Path: "assets/digits/dp.png"}}},
		BarSets: map[string]v3assets.BarSet{"temperature_bar": {ID: "temperature_bar", Cells: map[string]v3assets.ImageAsset{"off": {Path: "assets/bar/off.png"}, "on": {Path: "assets/bar/on.png"}, "warning": {Path: "assets/bar/warning.png"}}}},
		FrameSets: map[string]v3assets.FrameSet{"throttle_frames": {ID: "throttle_frames", Frames: frames, First: 0, Last: 4}},
		IndicatorSets: map[string]v3assets.IndicatorSet{"warning": {ID: "warning", States: map[string]v3assets.ImageAsset{v3assets.IndicatorStateOff: {Path: "assets/warning/off.png"}, v3assets.IndicatorStateOn: {Path: "assets/warning/on.png"}, v3assets.IndicatorStateUnknown: {Path: "assets/warning/unknown.png"}}}},
	}
}

func okState(id string, value float64, unit string) sensors.SensorState {
	return sensors.SensorState{ID: id, Value: value, Unit: unit, Status: sensors.StatusOK}
}

func floatPtr(value float64) *float64 { return &value }

func requireWidget(t *testing.T, scene Scene, id string) Widget {
	t.Helper()
	for _, widget := range scene.Widgets {
		if widget.ID == id {
			return widget
		}
	}
	t.Fatalf("widget %q not found in %#v", id, scene.Widgets)
	return Widget{}
}

func countParts(widget Widget, kind string) int { return len(partsByKind(widget, kind)) }

func partsByKind(widget Widget, kind string) []Part {
	parts := []Part{}
	for _, part := range widget.Parts {
		if part.Kind == kind {
			parts = append(parts, part)
		}
	}
	return parts
}

func cellNames(widget Widget) string {
	cells := []string{}
	for _, part := range widget.Parts {
		if part.Kind == PartKindCell {
			cells = append(cells, part.Cell)
		}
	}
	return strings.Join(cells, ",")
}

func statePart(widget Widget) string {
	for _, part := range widget.Parts {
		if part.Kind == PartKindState {
			return part.State
		}
	}
	return ""
}
