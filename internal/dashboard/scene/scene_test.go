package scene

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/assets"
	"github.com/MickMake/GoDriveLog/internal/dashboard/decoders"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestEvaluateSortsLayersByZOrder(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()
	dashboard.Layers = []config.DashboardLayerConfig{
		{ID: "front", Z: 10, Blocks: []string{"rpm_display"}},
		{ID: "back", Z: 0, Blocks: []string{"background_panel"}},
	}

	scene, err := Evaluate(dashboard, registry, baseDecoderValues(), nil, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if len(scene.Elements) != 2 {
		t.Fatalf("len(scene.Elements) = %d, want 2", len(scene.Elements))
	}
	if scene.Elements[0].ID != "background_panel" || scene.Elements[1].ID != "rpm_display" {
		t.Fatalf("elements order = %s, %s; want background_panel, rpm_display", scene.Elements[0].ID, scene.Elements[1].ID)
	}
}

func TestEvaluateConditionCanHideBlockFromDecoderValue(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()
	min := 6000.0

	scene, err := Evaluate(dashboard, registry, baseDecoderValues(), nil, Options{Conditions: map[string]Condition{
		"rpm_display": {Decoder: "rpm_value", Min: &min},
	}})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	rpm := findElement(t, scene.Elements, "rpm_display")
	if rpm.Visible {
		t.Fatal("rpm_display visible = true, want false")
	}
}

func TestEvaluateConditionCanShowBlockFromSensorValue(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()
	min := 50.0

	scene, err := Evaluate(dashboard, registry, baseDecoderValues(), map[string]sensors.SensorState{
		"throttle_position": {ID: "throttle_position", Value: 75, Status: sensors.StatusOK},
	}, Options{Conditions: map[string]Condition{
		"throttle_bar": {Sensor: "throttle_position", Min: &min},
	}})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	bar := findElement(t, scene.Elements, "throttle_bar")
	if !bar.Visible {
		t.Fatal("throttle_bar visible = false, want true")
	}
}

func TestEvaluateSpriteFrameResolvesFrameSetAsset(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()

	scene, err := Evaluate(dashboard, registry, baseDecoderValues(), nil, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	bar := findElement(t, scene.Elements, "throttle_bar")
	if !bar.HasFrame {
		t.Fatal("throttle_bar HasFrame = false, want true")
	}
	if bar.Frame.Index != 1 || string(bar.Frame.Data) != "frame-1" {
		t.Fatalf("throttle_bar frame = %#v, want index 1 frame-1", bar.Frame)
	}
}

func TestEvaluateSpriteTextResolvesCharsetGlyphs(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()

	scene, err := Evaluate(dashboard, registry, baseDecoderValues(), nil, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	rpm := findElement(t, scene.Elements, "rpm_display")
	if rpm.Text != "10" {
		t.Fatalf("rpm_display Text = %q, want 10", rpm.Text)
	}
	if len(rpm.Glyphs) != 2 {
		t.Fatalf("len(rpm.Glyphs) = %d, want 2", len(rpm.Glyphs))
	}
	if string(rpm.Glyphs[0].Data) != "glyph-1" || string(rpm.Glyphs[1].Data) != "glyph-0" {
		t.Fatalf("rpm glyph data = %q, %q; want glyph-1, glyph-0", string(rpm.Glyphs[0].Data), string(rpm.Glyphs[1].Data))
	}
}

func TestEvaluateGroupContainsChildElements(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()
	dashboard.Layers = []config.DashboardLayerConfig{{ID: "base", Z: 0, Blocks: []string{"main_cluster"}}}

	scene, err := Evaluate(dashboard, registry, baseDecoderValues(), nil, Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	group := findElement(t, scene.Elements, "main_cluster")
	if len(group.Children) != 2 {
		t.Fatalf("len(group.Children) = %d, want 2", len(group.Children))
	}
	if group.Children[0].ID != "rpm_display" || group.Children[1].ID != "throttle_bar" {
		t.Fatalf("group children = %s, %s; want rpm_display, throttle_bar", group.Children[0].ID, group.Children[1].ID)
	}
}

func TestEvaluateRejectsMissingGlyph(t *testing.T) {
	registry := makeRegistry(t)
	dashboard := baseDashboard()
	values := baseDecoderValues()
	values["rpm_digits"] = decoders.Value{Type: decoders.ValueTypeDigits, Digits: []string{"9"}}

	_, err := Evaluate(dashboard, registry, values, nil, Options{})
	if err == nil {
		t.Fatal("Evaluate returned nil error, want missing glyph error")
	}
}

func makeRegistry(t *testing.T) *assets.Registry {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"assets/background.png":       "background",
		"assets/frames/frame_000.png": "frame-0",
		"assets/frames/frame_001.png": "frame-1",
		"assets/digits/0.png":         "glyph-0",
		"assets/digits/1.png":         "glyph-1",
	}
	for path, content := range files {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	registry, err := assets.Load(baseDashboard(), filepath.Join(root, "dashboard.yaml"))
	if err != nil {
		t.Fatalf("assets.Load returned error: %v", err)
	}
	return registry
}

func baseDashboard() config.DashboardConfig {
	return config.DashboardConfig{
		AssetRoot: "assets",
		Assets: []config.DashboardAssetConfig{
			{ID: "background", Type: config.DashboardAssetImage, Path: "background.png"},
			{ID: "throttle_frames", Type: config.DashboardAssetFrameSet, Frames: []string{"frames/frame_000.png", "frames/frame_001.png"}},
			{ID: "yellow_digits", Type: config.DashboardAssetCharset, Glyphs: map[string]string{"0": "digits/0.png", "1": "digits/1.png"}},
		},
		Blocks: []config.DashboardBlockConfig{
			{ID: "background_panel", Type: config.DashboardBlockImage, Asset: "background", Geometry: config.RectConfig{X: 0, Y: 0, Width: 800, Height: 480}},
			{ID: "rpm_display", Type: config.DashboardBlockSpriteText, Asset: "yellow_digits", Decoder: "rpm_digits", Geometry: config.RectConfig{X: 100, Y: 60, Width: 240, Height: 80}},
			{ID: "throttle_bar", Type: config.DashboardBlockSpriteFrame, Asset: "throttle_frames", Decoder: "throttle_frame", Geometry: config.RectConfig{X: 100, Y: 170, Width: 300, Height: 40}},
			{ID: "main_cluster", Type: config.DashboardBlockGroup, Blocks: []string{"rpm_display", "throttle_bar"}},
		},
		Layers: []config.DashboardLayerConfig{
			{ID: "base", Z: 0, Blocks: []string{"background_panel", "rpm_display", "throttle_bar"}},
		},
	}
}

func baseDecoderValues() map[string]decoders.Value {
	return map[string]decoders.Value{
		"rpm_digits":     {Type: decoders.ValueTypeDigits, Digits: []string{"1", "0"}},
		"rpm_value":      {Type: decoders.ValueTypeNumber, Number: 5000},
		"throttle_frame": {Type: decoders.ValueTypeFrameIndex, Number: 1, FrameIndex: 1},
	}
}

func findElement(t *testing.T, elements []Element, id string) Element {
	t.Helper()
	for _, element := range elements {
		if element.ID == id {
			return element
		}
		child := findChild(element.Children, id)
		if child != nil {
			return *child
		}
	}
	t.Fatalf("element %q not found", id)
	return Element{}
}

func findChild(elements []Element, id string) *Element {
	for _, element := range elements {
		if element.ID == id {
			return &element
		}
		if child := findChild(element.Children, id); child != nil {
			return child
		}
	}
	return nil
}
