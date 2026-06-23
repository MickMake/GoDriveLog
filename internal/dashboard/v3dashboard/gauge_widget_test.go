package v3dashboard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestRuntimeLoadsGaugeWidgetPackageAndRendersSensorState(t *testing.T) {
	packageDir := makeDashboardGaugePackage(t, 4, "%04.0f")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{780, 40}, Scale: 1.5}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("rpm", 12, "rpm"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	if len(scenes) != 1 {
		t.Fatalf("expected one scene, got %d", len(scenes))
	}
	widget := requireWidget(t, scenes[0], "rpm")
	if widget.Type != v3config.WidgetTypeGauge || widget.SensorID != "rpm" || widget.GaugeID != "dashboard_4_digit_rpm" {
		t.Fatalf("gauge widget identity = %#v", widget)
	}
	if widget.Status != sensors.StatusOK || widget.Text != "0012" || widget.Scale != 1.5 {
		t.Fatalf("gauge widget status/text/scale = %#v", widget)
	}
	if got := countParts(widget, PartKindLayer); got != 2 {
		t.Fatalf("expected panel/glass static layers, got %d from %#v", got, widget.Parts)
	}
	if got := gaugePartSequence(widget); !strings.HasPrefix(got, "layer:panel,") || !strings.HasSuffix(got, ",layer:glass") {
		t.Fatalf("expected panel under digits and glass over digits, got %q", got)
	}
	if got := characters(widget); got != "0012" {
		t.Fatalf("expected rendered RPM characters 0012, got %q", got)
	}
	part := firstPartCharacter(widget, "1")
	if part.Slot != 2 || len(part.Position) != 2 || part.Position[0] != 22 || part.Position[1] != 12 {
		t.Fatalf("expected package digit position on character part, got %#v", part)
	}
}

func TestRuntimeGaugeWidgetNonOKStateDoesNotRenderLiveDigits(t *testing.T) {
	packageDir := makeDashboardGaugePackage(t, 4, "%04.0f")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{780, 40}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout, Error: "no response"})
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "rpm")
	if widget.Status != sensors.StatusTimeout || widget.Error != "no response" {
		t.Fatalf("gauge widget status/error = %#v", widget)
	}
	if widget.Text != "" {
		t.Fatalf("expected no live text for non-ok gauge state, got %q", widget.Text)
	}
	if got := countParts(widget, PartKindCharacter); got != 0 {
		t.Fatalf("expected no live character parts for non-ok gauge state, got %d", got)
	}
	if got := countParts(widget, PartKindLayer); got != 2 {
		t.Fatalf("expected static layers to remain, got %d", got)
	}
	if got := gaugePartSequence(widget); got != "layer:panel,layer:glass" {
		t.Fatalf("expected non-ok static layers in draw order, got %q", got)
	}
}

func TestRuntimeGaugeWidgetSceneSignatureChangesWithFormattedOutput(t *testing.T) {
	packageDir := makeDashboardGaugePackage(t, 4, "%04.0f")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{780, 40}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensorEvent("rpm", okState("rpm", 42.1, "rpm")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first gauge event to change rendered output")
	}
	_, changed, err = runtime.ApplyEvent(sensorEvent("rpm", okState("rpm", 42.2, "rpm")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected unchanged formatted gauge output to skip redraw")
	}
	_, changed, err = runtime.ApplyEvent(sensorEvent("rpm", okState("rpm", 43, "rpm")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed formatted gauge output to redraw")
	}
}

func TestRuntimeGaugeWidgetSceneSignatureChangesWithNonOKDigitPositions(t *testing.T) {
	packageDir := makeDashboardGaugePackage(t, 4, "%04.0f")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{780, 40}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	runtime.SetState(sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout})

	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	firstSignature := sceneSignature(scenes[0])

	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardGaugeYAMLWithOffset(4, "%04.0f", 777)), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	scenes, err = runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	if firstSignature == sceneSignature(scenes[0]) {
		t.Fatalf("expected non-ok scene signature to change when package digit positions change")
	}
}

func TestRuntimeLoadsRadialGaugeWidgetPackageAndPreservesAnglePivots(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{120, 80}, Scale: 0.75}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("rpm", 3500, "rpm"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "rpm")
	if widget.Type != v3config.WidgetTypeGauge || widget.SensorID != "rpm" || widget.GaugeID != "dashboard_radial_rpm" {
		t.Fatalf("radial widget identity = %#v", widget)
	}
	if widget.Status != sensors.StatusOK || widget.Scale != 0.75 || widget.GaugeAngle != 0 {
		t.Fatalf("radial widget status/scale/angle = %#v", widget)
	}
	if widget.GaugeFacePivot.X != 0.5 || widget.GaugeFacePivot.Y != 0.55 || widget.GaugeNeedlePivot.X != 0.5 || widget.GaugeNeedlePivot.Y != 0.9 {
		t.Fatalf("radial widget pivots = face %#v needle %#v", widget.GaugeFacePivot, widget.GaugeNeedlePivot)
	}
	if got := gaugePartSequence(widget); got != "layer:background,layer:face,layer:ticks,needle:0,layer:overlay" {
		t.Fatalf("radial part sequence = %q", got)
	}
	needle := firstPartKind(widget, PartKindNeedle)
	if needle.Layer != "needle" || needle.AssetPath == "" || needle.Angle != 0 {
		t.Fatalf("needle part = %#v", needle)
	}
	if needle.FacePivot != widget.GaugeFacePivot || needle.NeedlePivot != widget.GaugeNeedlePivot {
		t.Fatalf("needle pivots = face %#v needle %#v", needle.FacePivot, needle.NeedlePivot)
	}
}

func TestRuntimeRadialGaugeWidgetNonOKStateDoesNotRenderNeedle(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout, Error: "no response"})
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "rpm")
	if widget.Status != sensors.StatusTimeout || widget.Error != "no response" {
		t.Fatalf("radial widget status/error = %#v", widget)
	}
	if got := countParts(widget, PartKindNeedle); got != 0 {
		t.Fatalf("expected no live radial needle for non-ok state, got %d", got)
	}
	if got := gaugePartSequence(widget); got != "layer:background,layer:face,layer:ticks,layer:overlay" {
		t.Fatalf("expected static radial layers in draw order, got %q", got)
	}
}

func TestRuntimeRadialGaugeWidgetSceneSignatureChangesWithAngle(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensorEvent("rpm", okState("rpm", 3500, "rpm")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first radial event to change rendered output")
	}
	_, changed, err = runtime.ApplyEvent(sensorEvent("rpm", okState("rpm", 3500, "rpm")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected unchanged radial angle to skip redraw")
	}
	_, changed, err = runtime.ApplyEvent(sensorEvent("rpm", okState("rpm", 4000, "rpm")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed radial angle to redraw")
	}
}

func makeDashboardGaugePackage(t *testing.T, count int, format string) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/7Seg/7Seg4Digits.png",
		"assets/gauges/7Seg/Glass.png",
		"assets/gauges/7Seg/7SegBack.png",
		"assets/gauges/7Seg/amber/7SegDP.png",
	}
	for digit := 0; digit <= 9; digit++ {
		files = append(files, fmt.Sprintf("assets/gauges/7Seg/amber/7Seg%d.png", digit))
	}
	for _, path := range files {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(path), 0o600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}
	packageDir := filepath.Join(root, "assets", "gauges", "7Seg", "amber", fmt.Sprintf("%d_digit_rpm", count))
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardGaugeYAML(count, format)), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func makeDashboardRadialGaugePackage(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/radial/simple_rpm/background.png",
		"assets/gauges/radial/simple_rpm/face.png",
		"assets/gauges/radial/simple_rpm/ticks.png",
		"assets/gauges/radial/simple_rpm/needle.png",
		"assets/gauges/radial/simple_rpm/overlay.png",
	}
	for _, path := range files {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(path), 0o600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "simple_rpm")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardRadialGaugeYAML()), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func dashboardGaugeYAML(count int, format string) string {
	return dashboardGaugeYAMLWithOffset(count, format, 2)
}

func dashboardGaugeYAMLWithOffset(count int, format string, xOffset int) string {
	var positions strings.Builder
	for slot := 0; slot < count; slot++ {
		positions.WriteString(fmt.Sprintf("    - [%d, 12]\n", slot*10+xOffset))
	}
	return fmt.Sprintf(`id: dashboard_%d_digit_rpm
type: numeric
sensor: rpm
format: %q
size:
  width: 398
  height: 150
layers:
  panel: ../../7Seg4Digits.png
  glass: ../../Glass.png
digit_set:
  background: ../../7SegBack.png
  characters:
    "0": ../7Seg0.png
    "1": ../7Seg1.png
    "2": ../7Seg2.png
    "3": ../7Seg3.png
    "4": ../7Seg4.png
    "5": ../7Seg5.png
    "6": ../7Seg6.png
    "7": ../7Seg7.png
    "8": ../7Seg8.png
    "9": ../7Seg9.png
  decimal_point: ../7SegDP.png
  spacing: 4
digits:
  count: %d
  positions:
%s`, count, format, count, positions.String())
}

func dashboardRadialGaugeYAML() string {
	return `id: dashboard_radial_rpm
type: radial
sensor: rpm
size:
  width: 512
  height: 512
layers:
  background: background.png
  face: face.png
  ticks: ticks.png
  needle: needle.png
  overlay: overlay.png
pivot:
  face: { x: 0.5, y: 0.55 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 7000
  start_angle: -135
  end_angle: 135
  clamp: true
`
}

func firstPartCharacter(widget Widget, character string) Part {
	for _, part := range widget.Parts {
		if part.Kind == PartKindCharacter && part.Character == character {
			return part
		}
	}
	return Part{}
}

func firstPartKind(widget Widget, kind string) Part {
	for _, part := range widget.Parts {
		if part.Kind == kind {
			return part
		}
	}
	return Part{}
}

func gaugePartSequence(widget Widget) string {
	parts := make([]string, 0, len(widget.Parts))
	for _, part := range widget.Parts {
		switch part.Kind {
		case PartKindLayer:
			parts = append(parts, "layer:"+part.Layer)
		case PartKindCharacter:
			parts = append(parts, "character:"+part.Character)
		case PartKindNeedle:
			parts = append(parts, fmt.Sprintf("needle:%.0f", part.Angle))
		default:
			parts = append(parts, part.Kind)
		}
	}
	return strings.Join(parts, ",")
}

var _ = v3assets.IndicatorStateOff
