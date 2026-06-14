package v3dashboard

import (
	"strings"
	"testing"
	"time"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestRuntimeUsesResolvedSelectedDashboardsOnly(t *testing.T) {
	cfg := testConfig()
	plan, err := v3config.Resolve(cfg, "vw_caddy")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	if runtime.DashboardCount() != 1 {
		t.Fatalf("expected one selected dashboard, got %d", runtime.DashboardCount())
	}
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	if len(scenes) != 1 || scenes[0].DashboardID != "primary" {
		t.Fatalf("expected only selected primary dashboard, got %#v", scenes)
	}
}

func TestRuntimeRendersImageDigitAndIndicatorFromSensorState(t *testing.T) {
	runtime := testRuntime(t)
	runtime.SetState(okState("speed", 42, "km/h"))
	runtime.SetState(okState("warning", 1, "bool"))

	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	if len(scenes) != 1 {
		t.Fatalf("expected one scene, got %d", len(scenes))
	}

	panel := requireWidget(t, scenes[0], "panel")
	if panel.Type != v3config.WidgetTypeImage || countParts(panel, PartKindImage) != 1 {
		t.Fatalf("expected static panel image widget, got %#v", panel)
	}

	speed := requireWidget(t, scenes[0], "speed")
	if speed.Status != sensors.StatusOK || speed.Text != "042" {
		t.Fatalf("expected formatted speed 042 with ok status, got %#v", speed)
	}
	if got := characters(speed); got != "042" {
		t.Fatalf("expected digit characters 042, got %q", got)
	}

	warning := requireWidget(t, scenes[0], "warning")
	if warning.Status != sensors.StatusOK || statePart(warning) != v3assets.IndicatorStateOn {
		t.Fatalf("expected warning indicator on, got %#v", warning)
	}
}

func TestApplyEventSkipsUnchangedFormattedDigitOutput(t *testing.T) {
	runtime := testRuntime(t)

	_, changed, err := runtime.ApplyEvent(sensorEvent("speed", okState("speed", 42.1, "km/h")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first event to change rendered output")
	}

	_, changed, err = runtime.ApplyEvent(sensorEvent("speed", okState("speed", 42.2, "km/h")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected unchanged formatted output to skip redraw")
	}

	_, changed, err = runtime.ApplyEvent(sensorEvent("speed", okState("speed", 43.0, "km/h")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed formatted output to redraw")
	}
}

func TestDigitDecimalPointDoesNotConsumeSlot(t *testing.T) {
	dashboard := Dashboard{
		ID:     "primary",
		Assets: testAssetRegistry(),
		Config: v3config.DashboardConfig{
			Display: "test",
			Size:    v3config.SizeConfig{Width: 320, Height: 120},
			Widgets: []v3config.WidgetConfig{
				{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{0, 0}, Digits: 3, Format: "%.1f"},
			},
		},
	}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"speed": okState("speed", 12.3, "km/h")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "speed")
	if got := characters(widget); got != "123" {
		t.Fatalf("expected decimal separator not to consume slot, got characters %q", got)
	}
	decimal := partsByKind(widget, PartKindDecimalPoint)
	if len(decimal) != 1 || decimal[0].Slot != 1 {
		t.Fatalf("expected decimal point overlay on slot 1, got %#v", decimal)
	}
}

func TestIndicatorUsesUnknownForNonOKStatus(t *testing.T) {
	runtime := testRuntime(t)
	_, changed, err := runtime.ApplyEvent(sensorEvent("warning", sensors.SensorState{ID: "warning", Value: 1, Status: sensors.StatusStale, UpdatedAt: time.Unix(1, 0)}))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected stale status to change indicator output")
	}
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	warning := requireWidget(t, scenes[0], "warning")
	if statePart(warning) != v3assets.IndicatorStateUnknown {
		t.Fatalf("expected stale indicator to use unknown state, got %#v", warning)
	}
}

func TestDigitReportsMissingNonNumericCharacterAsset(t *testing.T) {
	registry := testAssetRegistry()
	delete(registry.DigitSets["digits"].Characters, "-")
	dashboard := Dashboard{
		ID:     "primary",
		Assets: registry,
		Config: v3config.DashboardConfig{
			Display: "test",
			Size:    v3config.SizeConfig{Width: 320, Height: 120},
			Widgets: []v3config.WidgetConfig{
				{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{0, 0}, Digits: 4, Format: "%04.0f"},
			},
		},
	}

	_, err := dashboard.Render(map[string]sensors.SensorState{"speed": okState("speed", -12, "km/h")})
	if err == nil {
		t.Fatalf("expected missing minus glyph to fail")
	}
	if !strings.Contains(err.Error(), "has no character asset for \"-\"") {
		t.Fatalf("expected useful missing character error, got %v", err)
	}
}

func testRuntime(t *testing.T) *Runtime {
	t.Helper()
	plan, err := v3config.Resolve(testConfig(), "vw_caddy")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	return runtime
}

func testConfig() v3config.Config {
	characters := map[string]string{}
	for _, ch := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"} {
		characters[ch] = "assets/digits/" + ch + ".png"
	}
	characters["-"] = "assets/digits/minus.png"
	return v3config.Config{
		Vehicles: map[string]v3config.VehicleConfig{
			"vw_caddy": {
				Name:       "VW Caddy",
				OBD:        v3config.OBDConfig{Address: "tcp://127.0.0.1:35000", Timeout: 1000},
				Dashboards: []string{"primary"},
			},
		},
		Sensors: map[string]v3config.SensorConfig{
			"speed":   {Type: "obd", PID: "010D", Unit: "km/h", Poll: 250},
			"warning": {Type: "obd", PID: "0142", Unit: "bool", Poll: 500},
		},
		Assets: v3config.AssetConfig{
			ImageSets: map[string]v3config.ImageSetConfig{
				"panel": {Image: "assets/panel.png"},
			},
			DigitSets: map[string]v3config.DigitSetConfig{
				"digits": {Characters: characters, DecimalPoint: "assets/digits/dp.png", Background: "assets/digits/back.png"},
			},
			IndicatorSets: map[string]v3config.IndicatorSetConfig{
				"warning": {States: map[string]string{
					v3assets.IndicatorStateOff:     "assets/warning/off.png",
					v3assets.IndicatorStateOn:      "assets/warning/on.png",
					v3assets.IndicatorStateUnknown: "assets/warning/unknown.png",
				}},
			},
		},
		Dashboards: map[string]v3config.DashboardConfig{
			"primary": {
				Display: "HDMI-1",
				Size:    v3config.SizeConfig{Width: 320, Height: 120},
				Widgets: []v3config.WidgetConfig{
					{ID: "panel", Type: v3config.WidgetTypeImage, Asset: "panel", Position: []int{0, 0}},
					{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{10, 10}, Digits: 3, Format: "%03.0f"},
					{ID: "warning", Type: v3config.WidgetTypeIndicator, Sensor: "warning", Asset: "warning", Position: []int{100, 10}},
				},
			},
			"alternate": {
				Display: "HDMI-2",
				Size:    v3config.SizeConfig{Width: 320, Height: 120},
				Widgets: []v3config.WidgetConfig{
					{ID: "panel", Type: v3config.WidgetTypeImage, Asset: "panel", Position: []int{0, 0}},
				},
			},
		},
	}
}

func testAssetRegistry() *v3assets.Registry {
	characters := map[string]v3assets.ImageAsset{}
	for _, ch := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"} {
		characters[ch] = v3assets.ImageAsset{Path: "assets/digits/" + ch + ".png"}
	}
	characters["-"] = v3assets.ImageAsset{Path: "assets/digits/minus.png"}
	return &v3assets.Registry{
		Images: map[string]v3assets.ImageSet{
			"panel": {ID: "panel", Image: &v3assets.ImageAsset{Path: "assets/panel.png"}},
		},
		DigitSets: map[string]v3assets.DigitSet{
			"digits": {
				ID:           "digits",
				Background:   &v3assets.ImageAsset{Path: "assets/digits/back.png"},
				Characters:   characters,
				DecimalPoint: &v3assets.ImageAsset{Path: "assets/digits/dp.png"},
			},
		},
		IndicatorSets: map[string]v3assets.IndicatorSet{
			"warning": {ID: "warning", States: map[string]v3assets.ImageAsset{
				v3assets.IndicatorStateOff:     {Path: "assets/warning/off.png"},
				v3assets.IndicatorStateOn:      {Path: "assets/warning/on.png"},
				v3assets.IndicatorStateUnknown: {Path: "assets/warning/unknown.png"},
			}},
		},
	}
}

func okState(id string, value float64, unit string) sensors.SensorState {
	return sensors.SensorState{ID: id, Value: value, Unit: unit, Status: sensors.StatusOK, UpdatedAt: time.Unix(1, 0)}
}

func sensorEvent(id string, state sensors.SensorState) sensors.SensorEvent {
	return sensors.SensorEvent{Kind: sensors.EventKindValueChange, SensorID: id, State: state, Timestamp: state.UpdatedAt, ReadAt: state.UpdatedAt}
}

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

func countParts(widget Widget, kind string) int {
	return len(partsByKind(widget, kind))
}

func partsByKind(widget Widget, kind string) []Part {
	parts := []Part{}
	for _, part := range widget.Parts {
		if part.Kind == kind {
			parts = append(parts, part)
		}
	}
	return parts
}

func characters(widget Widget) string {
	var b strings.Builder
	for _, part := range widget.Parts {
		if part.Kind == PartKindCharacter {
			b.WriteString(part.Character)
		}
	}
	return b.String()
}

func statePart(widget Widget) string {
	for _, part := range widget.Parts {
		if part.Kind == PartKindState {
			return part.State
		}
	}
	return ""
}
