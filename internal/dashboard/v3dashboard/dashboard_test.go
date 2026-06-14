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

func TestRuntimeRendersImageDigitBarFrameAndIndicatorFromSensorState(t *testing.T) {
	runtime := testRuntime(t)
	runtime.SetState(okState("speed", 42, "km/h"))
	runtime.SetState(okState("temperature", 45, "c"))
	runtime.SetState(okState("throttle", 50, "%"))
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

	temperature := requireWidget(t, scenes[0], "temperature")
	if got := cellNames(temperature); got != "on,on,on,off,off" {
		t.Fatalf("expected temperature bar cells, got %q from %#v", got, temperature.Parts)
	}

	throttle := requireWidget(t, scenes[0], "throttle")
	if got := framePart(throttle); got != 2 {
		t.Fatalf("expected throttle frame 2, got %d from %#v", got, throttle.Parts)
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

func TestApplyEventSkipsUnchangedFrameGaugeOutput(t *testing.T) {
	runtime := testRuntime(t)

	_, changed, err := runtime.ApplyEvent(sensorEvent("throttle", okState("throttle", 49, "%")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first frame event to change rendered output")
	}

	_, changed, err = runtime.ApplyEvent(sensorEvent("throttle", okState("throttle", 51, "%")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected same mapped frame to skip redraw")
	}

	_, changed, err = runtime.ApplyEvent(sensorEvent("throttle", okState("throttle", 76, "%")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed frame output to redraw")
	}
}

func TestDigitDecimalPointDoesNotConsumeSlot(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{0, 0}, Digits: 3, Format: "%.1f"}}}}

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

func TestDigitDefaultFormatDoesNotRequireDecimalPoint(t *testing.T) {
	registry := testAssetRegistry()
	set := registry.DigitSets["digits"]
	set.DecimalPoint = nil
	registry.DigitSets["digits"] = set

	dashboard := Dashboard{ID: "primary", Assets: registry, Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{0, 0}, Digits: 3}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"speed": okState("speed", 12.3, "km/h")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "speed")
	if widget.Text != "12" {
		t.Fatalf("expected default whole-number text 12, got %q", widget.Text)
	}
	if countParts(widget, PartKindDecimalPoint) != 0 {
		t.Fatalf("expected omitted format not to require decimal point, got %#v", widget.Parts)
	}
}

func TestDigitDecimalPointRendersBeforeForegroundForSlot(t *testing.T) {
	registry := testAssetRegistry()
	set := registry.DigitSets["digits"]
	set.Foreground = &v3assets.ImageAsset{Path: "assets/digits/front.png"}
	registry.DigitSets["digits"] = set

	dashboard := Dashboard{ID: "primary", Assets: registry, Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{0, 0}, Digits: 3, Format: "%.1f"}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"speed": okState("speed", 12.3, "km/h")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "speed")
	got := partKinds(partsForSlot(widget, 1))
	want := []string{PartKindBackground, PartKindCharacter, PartKindDecimalPoint, PartKindForeground}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("expected decimal point below foreground for slot 1, got %v want %v", got, want)
	}
}

func TestBarDisplayClampsAndReverseOnlyChangesFillDirection(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "temperature", Type: v3config.WidgetTypeBarDisplay, Sensor: "temperature", Asset: "temperature_bar", Position: []int{0, 0}, Cells: 5, Min: floatPtr(0), Max: floatPtr(100), Reverse: true}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"temperature": okState("temperature", 40, "c")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "temperature")
	if got := cellNames(widget); got != "off,off,off,on,on" {
		t.Fatalf("expected reverse fill direction only, got %q", got)
	}

	scene, err = dashboard.Render(map[string]sensors.SensorState{"temperature": okState("temperature", 150, "c")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget = requireWidget(t, scene, "temperature")
	if got := cellNames(widget); got != "on,on,on,on,on" {
		t.Fatalf("expected above-max value to fill all cells, got %q", got)
	}
}

func TestBarDisplayUsesZonesBySensorValue(t *testing.T) {
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

func TestBarDisplayUsesOffCellsForNonOKStatus(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "temperature", Type: v3config.WidgetTypeBarDisplay, Sensor: "temperature", Asset: "temperature_bar", Position: []int{0, 0}, Cells: 5, Min: floatPtr(0), Max: floatPtr(100)}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"temperature": {ID: "temperature", Value: 100, Status: sensors.StatusStale}})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "temperature")
	if widget.Status != sensors.StatusStale {
		t.Fatalf("expected stale status to be visible, got %#v", widget)
	}
	if got := cellNames(widget); got != "off,off,off,off,off" {
		t.Fatalf("expected non-ok bar to avoid live filled cells, got %q", got)
	}
}

func TestFrameGaugeClampsToFrameRange(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "throttle", Type: v3config.WidgetTypeFrameGauge, Sensor: "throttle", Asset: "throttle_frames", Position: []int{0, 0}, Min: floatPtr(0), Max: floatPtr(100)}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"throttle": okState("throttle", -10, "%")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if got := framePart(requireWidget(t, scene, "throttle")); got != 0 {
		t.Fatalf("expected below-min frame 0, got %d", got)
	}

	scene, err = dashboard.Render(map[string]sensors.SensorState{"throttle": okState("throttle", 150, "%")})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if got := framePart(requireWidget(t, scene, "throttle")); got != 4 {
		t.Fatalf("expected above-max frame 4, got %d", got)
	}
}

func TestFrameGaugeSkipsLiveFrameForNonOKStatus(t *testing.T) {
	dashboard := Dashboard{ID: "primary", Assets: testAssetRegistry(), Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "throttle", Type: v3config.WidgetTypeFrameGauge, Sensor: "throttle", Asset: "throttle_frames", Position: []int{0, 0}, Min: floatPtr(0), Max: floatPtr(100)}}}}

	scene, err := dashboard.Render(map[string]sensors.SensorState{"throttle": {ID: "throttle", Value: 100, Status: sensors.StatusError, Error: "read failed"}})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	widget := requireWidget(t, scene, "throttle")
	if widget.Status != sensors.StatusError || widget.Error != "read failed" {
		t.Fatalf("expected error status to be visible, got %#v", widget)
	}
	if countParts(widget, PartKindFrame) != 0 {
		t.Fatalf("expected non-ok frame gauge not to render live frame, got %#v", widget.Parts)
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
	dashboard := Dashboard{ID: "primary", Assets: registry, Config: v3config.DashboardConfig{Display: "test", Size: v3config.SizeConfig{Width: 320, Height: 120}, Widgets: []v3config.WidgetConfig{{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{0, 0}, Digits: 4, Format: "%04.0f"}}}}

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
			"speed":       {Type: "obd", PID: "010D", Unit: "km/h", Poll: 250},
			"temperature": {Type: "obd", PID: "0105", Unit: "c", Poll: 500},
			"throttle":    {Type: "obd", PID: "0111", Unit: "%", Poll: 250},
			"warning":     {Type: "obd", PID: "0142", Unit: "bool", Poll: 500},
		},
		Assets: v3config.AssetConfig{
			ImageSets: map[string]v3config.ImageSetConfig{
				"panel": {Image: "assets/panel.png"},
			},
			DigitSets: map[string]v3config.DigitSetConfig{
				"digits": {Characters: characters, DecimalPoint: "assets/digits/dp.png", Background: "assets/digits/back.png"},
			},
			BarSets: map[string]v3config.BarSetConfig{
				"temperature_bar": {
					Cells: map[string]string{
						"off":     "assets/bar/off.png",
						"on":      "assets/bar/on.png",
						"warning": "assets/bar/warning.png",
					},
				},
			},
			FrameSets: map[string]v3config.FrameSetConfig{
				"throttle_frames": {Frames: v3config.FrameRangeConfig{Path: "assets/throttle/frame_%03d.png", First: 0, Last: 4}},
			},
			IndicatorSets: map[string]v3config.IndicatorSetConfig{
				"warning": {States: map[string]string{v3assets.IndicatorStateOff: "assets/warning/off.png", v3assets.IndicatorStateOn: "assets/warning/on.png", v3assets.IndicatorStateUnknown: "assets/warning/unknown.png"}},
			},
		},
		Dashboards: map[string]v3config.DashboardConfig{
			"primary": {
				Display: "HDMI-1",
				Size:    v3config.SizeConfig{Width: 320, Height: 120},
				Widgets: []v3config.WidgetConfig{
					{ID: "panel", Type: v3config.WidgetTypeImage, Asset: "panel", Position: []int{0, 0}},
					{ID: "speed", Type: v3config.WidgetTypeDigitDisplay, Sensor: "speed", Asset: "digits", Position: []int{10, 10}, Digits: 3, Format: "%03.0f"},
					{ID: "temperature", Type: v3config.WidgetTypeBarDisplay, Sensor: "temperature", Asset: "temperature_bar", Position: []int{60, 10}, Cells: 5, Min: floatPtr(0), Max: floatPtr(100)},
					{ID: "throttle", Type: v3config.WidgetTypeFrameGauge, Sensor: "throttle", Asset: "throttle_frames", Position: []int{80, 10}, Min: floatPtr(0), Max: floatPtr(100)},
					{ID: "warning", Type: v3config.WidgetTypeIndicator, Sensor: "warning", Asset: "warning", Position: []int{100, 10}},
				},
			},
			"alternate": {
				Display: "HDMI-2",
				Size:    v3config.SizeConfig{Width: 320, Height: 120},
				Widgets: []v3config.WidgetConfig{{ID: "panel", Type: v3config.WidgetTypeImage, Asset: "panel", Position: []int{0, 0}}},
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

	frames := map[int]v3assets.ImageAsset{}
	for frame := 0; frame <= 4; frame++ {
		frames[frame] = v3assets.ImageAsset{Path: "assets/throttle/frame_00" + string(rune('0'+frame)) + ".png"}
	}

	return &v3assets.Registry{
		Images: map[string]v3assets.ImageSet{
			"panel": {ID: "panel", Image: &v3assets.ImageAsset{Path: "assets/panel.png"}},
		},
		DigitSets: map[string]v3assets.DigitSet{
			"digits": {ID: "digits", Background: &v3assets.ImageAsset{Path: "assets/digits/back.png"}, Characters: characters, DecimalPoint: &v3assets.ImageAsset{Path: "assets/digits/dp.png"}},
		},
		BarSets: map[string]v3assets.BarSet{
			"temperature_bar": {
				ID: "temperature_bar",
				Cells: map[string]v3assets.ImageAsset{
					"off":     {Path: "assets/bar/off.png"},
					"on":      {Path: "assets/bar/on.png"},
					"warning": {Path: "assets/bar/warning.png"},
				},
			},
		},
		FrameSets: map[string]v3assets.FrameSet{
			"throttle_frames": {ID: "throttle_frames", Frames: frames, First: 0, Last: 4},
		},
		IndicatorSets: map[string]v3assets.IndicatorSet{
			"warning": {ID: "warning", States: map[string]v3assets.ImageAsset{v3assets.IndicatorStateOff: {Path: "assets/warning/off.png"}, v3assets.IndicatorStateOn: {Path: "assets/warning/on.png"}, v3assets.IndicatorStateUnknown: {Path: "assets/warning/unknown.png"}}},
		},
	}
}

func okState(id string, value float64, unit string) sensors.SensorState {
	return sensors.SensorState{ID: id, Value: value, Unit: unit, Status: sensors.StatusOK, UpdatedAt: time.Unix(1, 0)}
}

func sensorEvent(id string, state sensors.SensorState) sensors.SensorEvent {
	return sensors.SensorEvent{Kind: sensors.EventKindValueChange, SensorID: id, State: state, Timestamp: state.UpdatedAt, ReadAt: state.UpdatedAt}
}

func floatPtr(value float64) *float64 {
	return &value
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

func partsForSlot(widget Widget, slot int) []Part {
	parts := []Part{}
	for _, part := range widget.Parts {
		if part.Slot == slot {
			parts = append(parts, part)
		}
	}
	return parts
}

func partKinds(parts []Part) []string {
	kinds := make([]string, 0, len(parts))
	for _, part := range parts {
		kinds = append(kinds, part.Kind)
	}
	return kinds
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

func cellNames(widget Widget) string {
	cells := []string{}
	for _, part := range widget.Parts {
		if part.Kind == PartKindCell {
			cells = append(cells, part.Cell)
		}
	}
	return strings.Join(cells, ",")
}

func framePart(widget Widget) int {
	for _, part := range widget.Parts {
		if part.Kind == PartKindFrame {
			return part.Frame
		}
	}
	return -1
}

func statePart(widget Widget) string {
	for _, part := range widget.Parts {
		if part.Kind == PartKindState {
			return part.State
		}
	}
	return ""
}
