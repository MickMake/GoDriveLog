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

func dashboardGaugeYAML(count int, format string) string {
	var positions strings.Builder
	for slot := 0; slot < count; slot++ {
		positions.WriteString(fmt.Sprintf("    - [%d, 12]\n", slot*10+2))
	}
	return fmt.Sprintf(`id: dashboard_%d_digit_rpm
type: seven_segment
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

func firstPartCharacter(widget Widget, character string) Part {
	for _, part := range widget.Parts {
		if part.Kind == PartKindCharacter && part.Character == character {
			return part
		}
	}
	return Part{}
}

var _ = v3assets.IndicatorStateOff
