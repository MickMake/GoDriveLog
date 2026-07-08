package v3dashboard

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	v3gauges "github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
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

func TestRuntimeRadialGaugeWidgetKeepsPointerMarkersHiddenWhenDisabled(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackageWithPointerMarkersAndRealism(t, "    max: false\n    min: false\n", "", false, nil, nil, false, nil, nil, nil, false)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
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
	if got := gaugePartSequence(widget); got != "layer:background,layer:face,layer:ticks,needle:0,layer:overlay" {
		t.Fatalf("radial part sequence = %q", got)
	}
	if got := countParts(widget, PartKindNeedleMin); got != 0 {
		t.Fatalf("expected no min marker parts, got %d", got)
	}
	if got := countParts(widget, PartKindNeedleMax); got != 0 {
		t.Fatalf("expected no max marker parts, got %d", got)
	}
}

func TestRuntimeRadialGaugeWidgetRendersPointerMarkersAboveNeedleBeforeOverlay(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackageWithPointerMarkersAndRealism(t, "    max: true\n    min: true\n", "", false, nil, nil, false, nil, nil, nil, false)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
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
	if got := gaugePartSequence(widget); got != "layer:background,layer:face,layer:ticks,needle:0,needle_min:0,needle_max:0,layer:overlay" {
		t.Fatalf("radial part sequence = %q", got)
	}
	minMarker := firstPartKind(widget, PartKindNeedleMin)
	maxMarker := firstPartKind(widget, PartKindNeedleMax)
	if minMarker.Layer != "needle_min" || minMarker.AssetPath == "" {
		t.Fatalf("min marker part = %#v", minMarker)
	}
	if maxMarker.Layer != "needle_max" || maxMarker.AssetPath == "" {
		t.Fatalf("max marker part = %#v", maxMarker)
	}
	if minMarker.FacePivot != widget.GaugeFacePivot || minMarker.NeedlePivot != widget.GaugeNeedlePivot {
		t.Fatalf("min marker pivots = face %#v needle %#v", minMarker.FacePivot, minMarker.NeedlePivot)
	}
	if maxMarker.FacePivot != widget.GaugeFacePivot || maxMarker.NeedlePivot != widget.GaugeNeedlePivot {
		t.Fatalf("max marker pivots = face %#v needle %#v", maxMarker.FacePivot, maxMarker.NeedlePivot)
	}
}

func TestRuntimeRadialGaugeWidgetKeepsAveragePointerMarkerHiddenWhenDisabled(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackageWithPointerMarkersAndRealism(t, "    average: false\n", "", false, nil, nil, false, nil, nil, nil, false)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	_, _, err = runtime.ApplyEvent(sensorEventAt("rpm", okState("rpm", 3500, "rpm"), time.Unix(1, 0)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "rpm")
	if got := gaugePartSequence(widget); got != "layer:background,layer:face,layer:ticks,needle:0,layer:overlay" {
		t.Fatalf("radial part sequence = %q", got)
	}
	if got := countParts(widget, PartKindNeedleAverage); got != 0 {
		t.Fatalf("expected no average marker parts, got %d", got)
	}
}

func TestRuntimeRadialGaugeWidgetRendersAveragePointerMarkerAboveNeedleBeforeOverlay(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackageWithPointerMarkersAndRealism(t, "    max: true\n    min: true\n    average: true\n", "", false, nil, nil, false, nil, nil, nil, false)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	start := time.Unix(700, 0)
	_, _, err = runtime.ApplyEvent(sensorEventAt("rpm", okState("rpm", 0, "rpm"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, _, err := runtime.ApplyEvent(sensorEventAt("rpm", okState("rpm", 7000, "rpm"), start.Add(10*time.Second)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	widget := requireWidget(t, scenes[0], "rpm")
	if got := gaugePartSequence(widget); got != "layer:background,layer:face,layer:ticks,needle:135,needle_min:-135,needle_max:135,needle_average:36,layer:overlay" {
		t.Fatalf("radial part sequence = %q", got)
	}
	marker := firstPartKind(widget, PartKindNeedleAverage)
	if marker.Layer != "needle_average" || marker.AssetPath == "" {
		t.Fatalf("average marker part = %#v", marker)
	}
	if marker.FacePivot != widget.GaugeFacePivot || marker.NeedlePivot != widget.GaugeNeedlePivot {
		t.Fatalf("average marker pivots = face %#v needle %#v", marker.FacePivot, marker.NeedlePivot)
	}
}

func TestRuntimeRadialGaugeWidgetIncludesNeedleShadowBeforeNeedle(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackageWithNeedleShadow(t, []int{3, 4}, nil)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
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
	if got := gaugePartSequence(widget); got != "layer:background,layer:face,layer:ticks,needle_shadow:0,needle:0,layer:overlay" {
		t.Fatalf("radial part sequence = %q", got)
	}
	shadow := firstPartKind(widget, PartKindNeedleShadow)
	needle := firstPartKind(widget, PartKindNeedle)
	if !intSlicesEqual(shadow.Position, []int{3, 4}) {
		t.Fatalf("shadow offset = %#v, want [3 4]", shadow.Position)
	}
	if !almostEqual(shadow.Alpha, 0.35) {
		t.Fatalf("shadow alpha = %v", shadow.Alpha)
	}
	if !almostEqual(shadow.Angle, needle.Angle) || !almostEqual(widget.GaugeAngle, needle.Angle) {
		t.Fatalf("shadow/needle/widget angles = %v/%v/%v", shadow.Angle, needle.Angle, widget.GaugeAngle)
	}
}

func TestRuntimeRadialGaugeWidgetCalibrationOffsetChangesOnlyDisplayedAngle(t *testing.T) {
	offset := 12.0
	packageDir := makeDashboardRadialGaugePackageWithCalibrationOffset(t, &offset)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	input := okState("rpm", 3500, "rpm")
	runtime.SetState(input)
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "rpm")
	needle := firstPartKind(widget, PartKindNeedle)
	if !almostEqual(widget.GaugeAngle, 12) {
		t.Fatalf("widget gauge angle = %v, want 12", widget.GaugeAngle)
	}
	if !almostEqual(needle.Angle, 12) {
		t.Fatalf("needle angle = %v, want 12", needle.Angle)
	}
	if stored := runtime.states["rpm"].Value; !almostEqual(stored, 3500) {
		t.Fatalf("stored sensor value = %v, want 3500", stored)
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

func TestRuntimeGaugeMovementLifecycle(t *testing.T) {
	packageDir := makeDashboardRadialGaugePackageWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false, false)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	now := time.Unix(100, 0)
	runtime.clock = func() time.Time { return now }
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		if context.DashboardID != "primary" || context.WidgetID != "rpm" || context.SensorID != "rpm" || context.GaugeType != v3gauges.TypeRadial {
			t.Fatalf("unexpected movement planner context: %#v", context)
		}
		if state.Value != 7000 {
			t.Fatalf("unexpected movement planner state value: %#v", state)
		}
		if !current.HasValue || current.DisplayValue != 3500 || current.TargetValue != 7000 {
			t.Fatalf("unexpected current movement state: %#v", current)
		}
		return 200 * time.Millisecond
	}

	initialState := okState("rpm", 3500, "rpm")
	initialState.UpdatedAt = now
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     initialState,
		Timestamp: now,
		ReadAt:    now,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected initial gauge event to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("did not expect initial static state to activate movement")
	}

	now = now.Add(10 * time.Millisecond)
	targetState := okState("rpm", 7000, "rpm")
	targetState.UpdatedAt = now
	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     targetState,
		Timestamp: now,
		ReadAt:    now,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed target value to start movement")
	}
	if !runtime.HasActiveMovement() {
		t.Fatalf("expected active movement after target change")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.Phase != movementPhaseMoving || movement.PreviousDisplayValue != 3500 || movement.DisplayValue != 3500 || movement.TargetValue != 7000 {
		t.Fatalf("unexpected movement start state: %#v", movement)
	}
	if movement.Duration != 200*time.Millisecond || !movement.StartedAt.Equal(now) {
		t.Fatalf("unexpected movement timing: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; got != 0 {
		t.Fatalf("expected start render to keep previous angle, got %v", got)
	}

	now = now.Add(100 * time.Millisecond)
	scenes, changed, err = runtime.Tick(now)
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected active movement tick to redraw")
	}
	if !runtime.HasActiveMovement() {
		t.Fatalf("expected movement to remain active mid-tick")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.Phase != movementPhaseMoving {
		t.Fatalf("expected moving phase mid-tick, got %#v", movement)
	}
	if math.Abs(movement.DisplayValue-5250) > 0.001 {
		t.Fatalf("expected halfway display value, got %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got-67.5) > 0.001 {
		t.Fatalf("expected halfway angle 67.5, got %v", got)
	}

	now = now.Add(100 * time.Millisecond)
	scenes, changed, err = runtime.Tick(now)
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected settling tick to redraw")
	}
	if !runtime.HasActiveMovement() {
		t.Fatalf("expected settled phase to request one cleanup tick")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.Phase != movementPhaseSettled || movement.DisplayValue != 7000 || movement.TargetValue != 7000 {
		t.Fatalf("unexpected settled state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got-135) > 0.001 {
		t.Fatalf("expected settled angle 135, got %v", got)
	}

	now = now.Add(1 * time.Millisecond)
	scenes, changed, err = runtime.Tick(now)
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected cleanup tick to finish movement lifecycle")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected static phase after cleanup tick")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.Phase != movementPhaseStatic || movement.DisplayValue != 7000 || movement.TargetValue != 7000 {
		t.Fatalf("unexpected final static movement state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got-135) > 0.001 {
		t.Fatalf("expected final static angle 135, got %v", got)
	}

	_, changed, err = runtime.ApplyEvent(sensorEvent("rpm", okState("rpm", 7000, "rpm")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected unchanged target value to skip redraw once static")
	}
}

func TestRuntimeGaugeMovementUsesEventTimestamp(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false, false)
	wallClock := time.Unix(500, 0)
	runtime.clock = func() time.Time { return wallClock }
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 200 * time.Millisecond
	}

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	eventTime := start.Add(10 * time.Millisecond)
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: eventTime,
		ReadAt:    wallClock,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected movement-start event to redraw")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.StartedAt.Equal(eventTime) {
		t.Fatalf("StartedAt = %v, want event timestamp %v", movement.StartedAt, eventTime)
	}
}

func TestRuntimeGaugeMovementFallsBackToReadAt(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false, false)
	wallClock := time.Unix(500, 0)
	runtime.clock = func() time.Time { return wallClock }
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 200 * time.Millisecond
	}

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	readAt := start.Add(10 * time.Millisecond)
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:     sensors.EventKindValueChange,
		SensorID: "rpm",
		State:    okState("rpm", 3500, "rpm"),
		ReadAt:   readAt,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected movement-start event to redraw")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.StartedAt.Equal(readAt) {
		t.Fatalf("StartedAt = %v, want ReadAt %v", movement.StartedAt, readAt)
	}
}

func TestRuntimeGaugeMovementFallsBackToRuntimeClock(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false, false)
	wallClock := time.Unix(500, 0)
	runtime.clock = func() time.Time { return wallClock }
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 200 * time.Millisecond
	}

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:     sensors.EventKindValueChange,
		SensorID: "rpm",
		State:    okState("rpm", 3500, "rpm"),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected movement-start event to redraw")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.StartedAt.Equal(wallClock) {
		t.Fatalf("StartedAt = %v, want runtime wall clock %v", movement.StartedAt, wallClock)
	}
}

func TestRuntimeGaugeMovementRetargetsFromAdvancedDisplayPosition(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false, false)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 200 * time.Millisecond
	}

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	firstTarget := start.Add(10 * time.Millisecond)
	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: firstTarget,
		ReadAt:    firstTarget,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first movement leg to redraw")
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got+135) > 0.001 {
		t.Fatalf("expected first leg to start at angle -135, got %v", got)
	}

	retargetAt := firstTarget.Add(100 * time.Millisecond)
	scenes, changed, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: retargetAt,
		ReadAt:    retargetAt,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected retarget event to redraw")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.Phase != movementPhaseMoving {
		t.Fatalf("expected retargeted movement to remain active, got %#v", movement)
	}
	if math.Abs(movement.PreviousDisplayValue-1750) > 0.001 || math.Abs(movement.DisplayValue-1750) > 0.001 {
		t.Fatalf("expected retarget to continue from advanced halfway display, got %#v", movement)
	}
	if movement.TargetValue != 7000 || !movement.StartedAt.Equal(retargetAt) {
		t.Fatalf("unexpected retarget timing/state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got+67.5) > 0.001 {
		t.Fatalf("expected retarget render to continue at angle -67.5, got %v", got)
	}

	scenes, changed, err = runtime.Tick(retargetAt.Add(200 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final target tick to redraw")
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got-135) > 0.001 {
		t.Fatalf("expected final target angle 135, got %v", got)
	}

	scenes, changed, err = runtime.Tick(retargetAt.Add(201 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected cleanup tick after retarget to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected retargeted movement to return to static after cleanup")
	}
}

func TestRuntimeRadialGaugeMovementDefaultsToImmediateWithoutDamping(t *testing.T) {
	runtime := testRadialMovementRuntime(t)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 200 * time.Millisecond
	}

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected immediate policy update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected default immediate policy to avoid active movement")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.Policy != v3gauges.MovementPolicyImmediate || movement.Phase != movementPhaseStatic || movement.DisplayValue != 3500 {
		t.Fatalf("unexpected immediate policy state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got) > 0.001 {
		t.Fatalf("expected immediate policy to jump to target angle 0, got %v", got)
	}
}

func TestRuntimeRadialGaugeHysteresisDisabledRemainsImmediate(t *testing.T) {
	runtime := testRadialMovementRuntimeWithHysteresis(t, false)
	start := time.Unix(90, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected non-hysteresis update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected disabled hysteresis to avoid active movement")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.HysteresisEnabled || movement.ApproachDirection != 0 || movement.DisplayValue != 3500 {
		t.Fatalf("unexpected disabled hysteresis state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got) > 0.001 {
		t.Fatalf("expected disabled hysteresis to keep target angle 0, got %v", got)
	}
}

func TestRuntimeRadialGaugeHysteresisOffsetsRisingApproach(t *testing.T) {
	runtime := testRadialMovementRuntimeWithHysteresis(t, true)
	start := time.Unix(95, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected rising hysteresis update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected hysteresis-only update to stay immediate")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.HysteresisEnabled || movement.ApproachDirection != 1 || movement.RawTargetValue != 3500 || movement.DisplayValue <= 3500 {
		t.Fatalf("unexpected rising hysteresis state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got-2.7) > 0.001 {
		t.Fatalf("expected rising hysteresis angle 2.7, got %v", got)
	}
}

func TestRuntimeRadialGaugeHysteresisOffsetsFallingApproach(t *testing.T) {
	runtime := testRadialMovementRuntimeWithHysteresis(t, true)
	start := time.Unix(98, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected falling hysteresis update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected hysteresis-only falling update to stay immediate")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.HysteresisEnabled || movement.ApproachDirection != -1 || movement.RawTargetValue != 3500 || movement.DisplayValue >= 3500 {
		t.Fatalf("unexpected falling hysteresis state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got+2.7) > 0.001 {
		t.Fatalf("expected falling hysteresis angle -2.7, got %v", got)
	}
}

func TestRuntimeRadialGaugeHysteresisClampTrueStaysWithinValueMapBounds(t *testing.T) {
	runtime := testRadialMovementRuntimeWithHysteresisAndClamp(t, true, true)
	start := time.Unix(98, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamped hysteresis update to redraw")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 7000 || movement.RawTargetValue != 7000 {
		t.Fatalf("expected clamp=true hysteresis to stay at max value 7000, got %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got-135) > 0.001 {
		t.Fatalf("expected clamp=true hysteresis angle 135, got %v", got)
	}
}

func TestRuntimeRadialGaugeHysteresisClampFalsePreservesOutOfRangeRendering(t *testing.T) {
	runtime := testRadialMovementRuntimeWithHysteresisAndClamp(t, true, false)
	start := time.Unix(98, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected unclamped hysteresis update to redraw")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue <= 7000 || movement.RawTargetValue != 7000 {
		t.Fatalf("expected clamp=false hysteresis to keep out-of-range display value, got %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; got <= 135 {
		t.Fatalf("expected clamp=false hysteresis angle above 135, got %v", got)
	}
	if got := runtime.states["rpm"].Value; got != 7000 {
		t.Fatalf("expected stored source value 7000 to remain unchanged, got %v", got)
	}
}

func TestRuntimeRadialGaugeHysteresisLeavesStoredSourceValueUnchanged(t *testing.T) {
	runtime := testRadialMovementRuntimeWithHysteresis(t, true)
	start := time.Unix(99, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected hysteresis source update to redraw")
	}
	if got := runtime.states["rpm"].Value; got != 3500 {
		t.Fatalf("expected stored source value 3500 to remain unchanged, got %v", got)
	}

	_, changed, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(20 * time.Millisecond),
		ReadAt:    start.Add(20 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected unchanged source value to avoid redraw")
	}
}

func TestRuntimeRadialGaugeDampingAnimatesWithDefaultLinearCurve(t *testing.T) {
	runtime := testRadialMovementRuntimeWithDamping(t, true)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		if !current.DampingEnabled {
			t.Fatalf("expected damping-enabled movement state, got %#v", current)
		}
		return 200 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected configured damping to start active movement")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.Policy != v3gauges.MovementPolicyLinear {
		t.Fatalf("expected damping to default to linear movement curve, got %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got+135) > 0.001 {
		t.Fatalf("expected damping start angle -135, got %v", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected active damping tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if math.Abs(movement.DisplayValue-1750) > 0.001 {
		t.Fatalf("expected damping midpoint display 1750, got %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got+67.5) > 0.001 {
		t.Fatalf("expected damping midpoint angle -67.5, got %v", got)
	}
}

func TestRuntimeRadialGaugeStictionBelowThresholdHoldsDisplay(t *testing.T) {
	threshold := 200.0
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, &threshold, nil, false, false)

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 150, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed || scenes != nil {
		t.Fatalf("expected below-threshold stiction to suppress redraw, changed=%v scenes=%#v", changed, scenes)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected below-threshold stiction to remain static")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 0 || movement.TargetValue != 150 || movement.Phase != movementPhaseStatic {
		t.Fatalf("unexpected held stiction state: %#v", movement)
	}
}

func TestRuntimeRadialGaugeStictionReleasesAboveThreshold(t *testing.T) {
	threshold := 200.0
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, &threshold, nil, false, false)

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 150, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 250, "rpm"),
		Timestamp: start.Add(20 * time.Millisecond),
		ReadAt:    start.Add(20 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || scenes == nil {
		t.Fatalf("expected above-threshold stiction to release and redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected stiction-only release to jump without active damping")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 250 || movement.TargetValue != 250 || movement.Phase != movementPhaseStatic {
		t.Fatalf("unexpected released stiction state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; !(got > -135) {
		t.Fatalf("expected released needle angle to move off baseline, got %v", got)
	}
}

func TestRuntimeRadialGaugeStictionDefaultDisabledDoesNotHoldSmallChanges(t *testing.T) {
	runtime := testRadialMovementRuntime(t)

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 150, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || scenes == nil {
		t.Fatalf("expected default radial behavior to redraw immediately without stiction")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 150 || movement.Phase != movementPhaseStatic {
		t.Fatalf("unexpected default non-stiction state: %#v", movement)
	}
}

func TestRuntimeRadialGaugeOvershootDefaultDisabledDoesNotPassTarget(t *testing.T) {
	runtime := testRadialMovementRuntimeWithDamping(t, true)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.Tick(start.Add(260 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue > 3500 {
		t.Fatalf("expected default damping-only movement to stay at or below target, got %#v", movement)
	}
}

func TestRuntimeRadialGaugeOvershootAnimatesWithoutDamping(t *testing.T) {
	overshoot := &v3gauges.OvershootConfig{}
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, overshoot, false, false)
	start := time.Unix(100, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot-only radial movement to become active")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DampingEnabled {
		t.Fatalf("expected overshoot-only movement to keep damping disabled: %#v", movement)
	}
	if !movement.OvershootEnabled || movement.Duration <= 0 || movement.Phase != movementPhaseMoving {
		t.Fatalf("expected overshoot-only movement to schedule animation: %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(140 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot-only tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue <= 3500 || movement.DisplayValue > 7000 {
		t.Fatalf("expected overshoot-only movement above target and within range, got %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; got <= 0 || got > 135 {
		t.Fatalf("expected overshoot-only angle between target and max stop, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(230 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot-only settle tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 3500 || movement.TargetValue != 3500 || movement.Phase != movementPhaseSettled {
		t.Fatalf("expected overshoot-only movement to settle exactly on target, got %#v", movement)
	}
}

func TestRuntimeRadialGaugeOvershootModerateChangeTriggersWithLowerMinChangeRatio(t *testing.T) {
	minChangeRatio := 0.05
	ratio := 0.18
	overshoot := &v3gauges.OvershootConfig{
		Ratio:          &ratio,
		MinChangeRatio: &minChangeRatio,
	}
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, overshoot, false, false)
	start := time.Unix(100, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected lower overshoot threshold to activate movement")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.OvershootTargetValue <= 2500 {
		t.Fatalf("expected moderate change to schedule overshoot above target, got %#v", movement)
	}
}

func TestRuntimeRadialGaugeOvershootStaysBoundedAndSettlesOnTarget(t *testing.T) {
	overshoot := &v3gauges.OvershootConfig{}
	runtime := testRadialMovementRuntimeWithRealism(t, "", true, nil, overshoot, false, false)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.Tick(start.Add(220 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue <= 3500 || movement.DisplayValue > 7000 {
		t.Fatalf("expected bounded overshoot above target and within range, got %#v", movement)
	}
	if movement.OvershootTargetValue <= 3500 || movement.OvershootTargetValue > 7000 {
		t.Fatalf("unexpected overshoot target %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; got <= 0 || got > 135 {
		t.Fatalf("expected overshoot angle between target and max stop, got %v", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected settle tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.Phase != movementPhaseSettled || movement.DisplayValue != 3500 || movement.TargetValue != 3500 {
		t.Fatalf("unexpected settled overshoot state: %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; math.Abs(got) > 0.001 {
		t.Fatalf("expected settled overshoot angle 0, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(321 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected cleanup tick after overshoot settle")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected overshoot movement to return to static after cleanup")
	}
}

func TestRuntimeRadialPointerMarkersTrackRenderedOvershoot(t *testing.T) {
	overshoot := &v3gauges.OvershootConfig{}
	runtime := testRadialMovementRuntimeWithPointerMarkersAndRealism(t, "    max: true\n    min: true\n    window: 5m\n", "", false, nil, overshoot, false, false)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 3500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	_, changed, err := runtime.Tick(start.Add(220 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot tick to redraw")
	}

	markerState := runtime.pointerMarkers[movementKey("primary", "rpm")]
	if !markerState.Min.Set || !markerState.Max.Set {
		t.Fatalf("expected pointer marker min/max to be set, got %#v", markerState)
	}
	if markerState.Min.NormalizedPosition != 0 {
		t.Fatalf("expected initial zero reading to seed min marker, got %#v", markerState)
	}
	if markerState.Max.NormalizedPosition <= 0.5 || markerState.Max.NormalizedPosition > 1 {
		t.Fatalf("expected overshoot to push max marker past target, got %#v", markerState)
	}
	if got := runtime.states["rpm"].Value; got != 3500 {
		t.Fatalf("expected stored source state to remain at 3500 during overshoot, got %v", got)
	}
	if len(markerState.Samples) < 2 {
		t.Fatalf("expected overshoot movement tick to add a rolling sample, got %#v", markerState)
	}
}

func TestRuntimeRollingPointerMarkersExpireDuringIdleTicks(t *testing.T) {
	runtime := testRadialMovementRuntimeWithPointerMarkersAndRealism(t, "    max: true\n    min: true\n    window: 30m\n", "", false, nil, nil, false, false)
	start := time.Unix(100, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 1400, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	key := movementKey("primary", "rpm")
	initial := runtime.pointerMarkers[key]
	if len(initial.Samples) != 1 {
		t.Fatalf("expected one rolling sample after event, got %#v", initial)
	}
	if !initial.Samples[0].RecordedAt.Equal(start) {
		t.Fatalf("expected initial sample timestamp %v, got %#v", start, initial.Samples[0])
	}

	_, changed, err := runtime.Tick(start.Add(20 * time.Minute))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if changed {
		t.Fatalf("expected idle tick with unchanged render not to redraw")
	}
	state := runtime.pointerMarkers[key]
	if len(state.Samples) != 1 {
		t.Fatalf("expected idle tick not to add samples, got %#v", state)
	}
	if !state.Samples[0].RecordedAt.Equal(start) {
		t.Fatalf("expected idle tick not to refresh sample timestamp, got %#v", state.Samples[0])
	}

	_, changed, err = runtime.Tick(start.Add(31 * time.Minute))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected pruning-only idle tick to redraw when visible markers expire")
	}
	state = runtime.pointerMarkers[key]
	if len(state.Samples) != 0 {
		t.Fatalf("expected rolling sample to expire during idle ticks, got %#v", state)
	}
	if state.Min.Set || state.Max.Set {
		t.Fatalf("expected expired idle rolling markers to unset, got %#v", state)
	}
}

func TestRuntimeAverageOnlyGaugeReportsActiveWhileUnsettled(t *testing.T) {
	runtime := testRadialMovementRuntimeWithPointerMarkersAndRealism(t, "    average: true\n", "", false, nil, nil, false, false)
	start := time.Unix(600, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected seeded average marker at its initial target to remain inactive")
	}

	key := movementKey("primary", "rpm")
	initial := runtime.pointerMarkers[key]
	if !initial.Average.Set || initial.Average.NormalizedPosition != 0 || !initial.LastRenderedPositionSet || initial.LastRenderedPosition != 0 {
		t.Fatalf("expected initial rendered position to seed average marker, got %#v", initial)
	}
	if initial.Min.Set || initial.Max.Set {
		t.Fatalf("expected average-only config not to set min/max markers, got %#v", initial)
	}

	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Second),
		ReadAt:    start.Add(10 * time.Second),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed live gauge value to redraw")
	}
	if !runtime.HasActiveMovement() {
		t.Fatalf("expected unsettled average marker to request runtime ticks")
	}

	updated := runtime.pointerMarkers[key]
	expected := 1 - math.Exp(-1)
	if !updated.Average.Set || math.Abs(updated.Average.NormalizedPosition-expected) > 1e-12 {
		t.Fatalf("expected damped average marker position %v, got %#v", expected, updated)
	}
	if !updated.LastRenderedPositionSet || updated.LastRenderedPosition != 1 {
		t.Fatalf("expected rendered target position 1 to remain tracked, got %#v", updated)
	}
	if updated.Min.Set || updated.Max.Set || len(updated.Samples) != 0 {
		t.Fatalf("expected average-only state without min/max history, got %#v", updated)
	}
}

func TestRuntimeTickAdvancesAveragePointerMarkerWithoutNewEvent(t *testing.T) {
	runtime := testRadialMovementRuntimeWithPointerMarkersAndRealism(t, "    average: true\n", "", false, nil, nil, false, false)
	start := time.Unix(610, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Second),
		ReadAt:    start.Add(10 * time.Second),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	key := movementKey("primary", "rpm")
	before := runtime.pointerMarkers[key]
	_, _, err = runtime.Tick(start.Add(20 * time.Second))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}

	after := runtime.pointerMarkers[key]
	if after.Average.NormalizedPosition <= before.Average.NormalizedPosition || after.Average.NormalizedPosition >= 1 {
		t.Fatalf("expected idle tick to advance average marker toward target, before=%#v after=%#v", before, after)
	}
	if got := runtime.states["rpm"].Value; got != 7000 {
		t.Fatalf("expected stored source state to remain raw 7000, got %v", got)
	}
}

func TestRuntimeAveragePointerMarkerStopsRequiringTicksWhenSettled(t *testing.T) {
	runtime := testRadialMovementRuntimeWithPointerMarkersAndRealism(t, "    average: true\n", "", false, nil, nil, false, false)
	start := time.Unix(620, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Second),
		ReadAt:    start.Add(10 * time.Second),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	if !runtime.HasActiveMovement() {
		t.Fatalf("expected unsettled average marker to begin active")
	}
	_, _, err = runtime.Tick(start.Add(220 * time.Second))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}

	if runtime.HasActiveMovement() {
		t.Fatalf("expected effectively settled average marker to stop requesting ticks")
	}
	state := runtime.pointerMarkers[movementKey("primary", "rpm")]
	if math.Abs(state.Average.NormalizedPosition-state.LastRenderedPosition) > 1e-9 {
		t.Fatalf("expected settled average marker to be within epsilon of target, got %#v", state)
	}
}

func TestRuntimeMinMaxOnlyPointerMarkersDoNotReportActiveMovement(t *testing.T) {
	runtime := testRadialMovementRuntimeWithPointerMarkersAndRealism(t, "    max: true\n    min: true\n", "", false, nil, nil, false, false)
	start := time.Unix(630, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Second),
		ReadAt:    start.Add(10 * time.Second),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	if runtime.HasActiveMovement() {
		t.Fatalf("expected min/max-only pointer markers not to request idle settling ticks")
	}
	state := runtime.pointerMarkers[movementKey("primary", "rpm")]
	if !state.Min.Set || !state.Max.Set || state.Average.Set {
		t.Fatalf("expected min/max-only state to remain unchanged, got %#v", state)
	}
	if got := runtime.states["rpm"].Value; got != 7000 {
		t.Fatalf("expected stored source state to remain raw 7000, got %v", got)
	}
}

func TestRuntimeRadialGaugeOvershootOscillatesAcrossTargetAndSettles(t *testing.T) {
	ratio := 0.20
	minChangeRatio := 0.05
	maxSpanRatio := 0.08
	settleCycles := 1.75
	settleDamping := 4.5
	overshoot := &v3gauges.OvershootConfig{
		Ratio:          &ratio,
		MinChangeRatio: &minChangeRatio,
		MaxSpanRatio:   &maxSpanRatio,
		SettleMode:     v3gauges.OvershootSettleOscillate,
		SettleCycles:   &settleCycles,
		SettleDamping:  &settleDamping,
	}
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, overshoot, false, false)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2500, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	_, changed, err := runtime.Tick(start.Add(140 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected oscillating overshoot travel tick to redraw")
	}
	travelMovement := runtime.movements[movementKey("primary", "rpm")]
	if travelMovement.DisplayValue <= 2500 {
		t.Fatalf("expected travel phase above target, got %#v", travelMovement)
	}

	_, changed, err = runtime.Tick(start.Add(190 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected oscillating overshoot crossback tick to redraw")
	}
	crossbackMovement := runtime.movements[movementKey("primary", "rpm")]
	if crossbackMovement.DisplayValue >= 2500 {
		t.Fatalf("expected oscillating settle phase to cross below target, got %#v", crossbackMovement)
	}

	_, changed, err = runtime.Tick(start.Add(240 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected later oscillating settle tick to redraw")
	}
	laterMovement := runtime.movements[movementKey("primary", "rpm")]
	if laterMovement.DisplayValue <= 2500 {
		t.Fatalf("expected oscillation to swing back above target, got %#v", laterMovement)
	}
	if math.Abs(laterMovement.DisplayValue-2500) >= math.Abs(crossbackMovement.DisplayValue-2500) {
		t.Fatalf("expected later oscillation to diminish in magnitude, got crossback=%#v later=%#v", crossbackMovement, laterMovement)
	}

	_, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected oscillating settle completion tick to redraw")
	}
	finalMovement := runtime.movements[movementKey("primary", "rpm")]
	if finalMovement.DisplayValue != 2500 || finalMovement.TargetValue != 2500 || finalMovement.Phase != movementPhaseSettled {
		t.Fatalf("expected oscillating overshoot to settle exactly on target, got %#v", finalMovement)
	}
}

func TestRuntimeRadialGaugePegBounceDefaultDisabledSettlesImmediatelyAtStop(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, false, false)
	start := time.Unix(100, 0)

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.PegBounceEnabled || movement.Phase != movementPhaseStatic || movement.DisplayValue != 7000 {
		t.Fatalf("expected default stop movement without peg bounce, got %#v", movement)
	}
}

func TestRuntimeRadialGaugePegBounceAtMaxStopSettlesBackToLimit(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, true, false)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected peg bounce at max stop to animate")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.PegBounceEnabled || movement.TargetValue != 7000 || movement.PegBounceStopValue != 7000 || movement.PegBounceReboundValue >= 7000 {
		t.Fatalf("expected max-stop peg bounce to schedule inward rebound, got %#v", movement)
	}
	if movement.SettleDuration < defaultRadialPegBounceMinSettleDuration {
		t.Fatalf("expected visible peg-bounce settle duration, got %#v", movement)
	}
	if movement.PegBounceStopValue-movement.PegBounceReboundValue > (7000*defaultRadialPegBounceSpanRatio)+0.001 {
		t.Fatalf("expected bounded max-stop peg-bounce amplitude, got %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(270 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected max-stop peg bounce settle tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue >= 7000 || movement.DisplayValue <= movement.PegBounceReboundValue {
		t.Fatalf("expected max-stop peg bounce to rebound slightly below the stop, got %#v", movement)
	}
	if got := requireWidget(t, scenes[0], "rpm").GaugeAngle; got >= 135 || got <= 120 {
		t.Fatalf("expected max-stop peg bounce angle near the stop, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected max-stop peg bounce completion tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 7000 || movement.TargetValue != 7000 || movement.Phase != movementPhaseSettled {
		t.Fatalf("expected max-stop peg bounce to settle exactly at the stop, got %#v", movement)
	}
}

func TestRuntimeRadialGaugePegBounceAtMinStopSettlesBackToLimit(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, true, false)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 5000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 0, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	_, changed, err := runtime.Tick(start.Add(270 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected min-stop peg bounce settle tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.PegBounceStopValue != 0 || movement.PegBounceReboundValue <= 0 || movement.DisplayValue <= 0 {
		t.Fatalf("expected min-stop peg bounce to rebound above zero, got %#v", movement)
	}
	if movement.PegBounceReboundValue-movement.PegBounceStopValue > (7000*defaultRadialPegBounceSpanRatio)+0.001 {
		t.Fatalf("expected bounded min-stop peg-bounce amplitude, got %#v", movement)
	}

	_, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected min-stop peg bounce completion tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 0 || movement.TargetValue != 0 || movement.Phase != movementPhaseSettled {
		t.Fatalf("expected min-stop peg bounce to settle exactly at the stop, got %#v", movement)
	}
}

func TestRuntimeRadialGaugePegBounceDoesNotTriggerForInRangeTarget(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, true, false)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 6000, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected in-range peg-bounce-only change to redraw immediately")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.PegBounceEnabled || movement.PegBounceActive || movement.PegBounceReboundValue != 0 || movement.PegBounceStopValue != 0 {
		t.Fatalf("expected in-range movement to avoid scheduling peg bounce, got %#v", movement)
	}
	if movement.DisplayValue != 6000 || movement.TargetValue != 6000 || movement.Phase != movementPhaseStatic || runtime.HasActiveMovement() {
		t.Fatalf("expected in-range peg-bounce-only change to settle immediately with no active movement, got %#v", movement)
	}

	_, changed, err = runtime.Tick(start.Add(200 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if changed {
		t.Fatalf("expected no follow-up animation tick for in-range peg-bounce-only change")
	}
}

func TestRuntimeRadialGaugePegBounceExtendsShortMovementForVisibleSettle(t *testing.T) {
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, true, false)
	start := time.Unix(100, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 120 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 2000, "rpm"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "rpm",
		State:     okState("rpm", 7000, "rpm"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected short planned peg-bounce movement to animate")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if movement.TravelDuration < defaultRadialPegBounceMinTravelDuration {
		t.Fatalf("expected guarded peg-bounce travel duration, got %#v", movement)
	}
	if movement.SettleDuration < defaultRadialPegBounceMinSettleDuration {
		t.Fatalf("expected guarded peg-bounce settle duration, got %#v", movement)
	}
	if movement.Duration != movement.TravelDuration+movement.SettleDuration {
		t.Fatalf("expected peg-bounce duration to match guarded phases, got %#v", movement)
	}
	if movement.Duration <= 120*time.Millisecond {
		t.Fatalf("expected guarded peg-bounce movement to extend short planner duration, got %#v", movement)
	}
}

func TestRuntimeGaugeMovementEaseOutPolicyAdvancesFurtherThanLinear(t *testing.T) {
	linearRuntime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false, false)
	easeOutRuntime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyEaseOut, true, nil, nil, false, false)

	for _, runtime := range []*Runtime{linearRuntime, easeOutRuntime} {
		runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
			return 200 * time.Millisecond
		}
	}

	start := time.Unix(100, 0)
	for _, runtime := range []*Runtime{linearRuntime, easeOutRuntime} {
		_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "rpm",
			State:     okState("rpm", 0, "rpm"),
			Timestamp: start,
			ReadAt:    start,
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
		_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "rpm",
			State:     okState("rpm", 3500, "rpm"),
			Timestamp: start.Add(10 * time.Millisecond),
			ReadAt:    start.Add(10 * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
	}

	_, _, err := linearRuntime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	_, _, err = easeOutRuntime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}

	linearMovement := linearRuntime.movements[movementKey("primary", "rpm")]
	easeOutMovement := easeOutRuntime.movements[movementKey("primary", "rpm")]
	if linearMovement.Policy != v3gauges.MovementPolicyLinear || easeOutMovement.Policy != v3gauges.MovementPolicyEaseOut {
		t.Fatalf("unexpected policies: linear=%#v easeOut=%#v", linearMovement, easeOutMovement)
	}
	if !(easeOutMovement.DisplayValue > linearMovement.DisplayValue) {
		t.Fatalf("expected ease_out to advance further than linear: linear=%#v ease_out=%#v", linearMovement, easeOutMovement)
	}
}

func TestRuntimeLoadsOdometerGaugeWidgetPackage(t *testing.T) {
	packageDir := makeDashboardOdometerGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "trip", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{300, 40}, Scale: 1.5}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("trip_distance", 12.3, "km"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "trip")
	if widget.Type != v3config.WidgetTypeGauge || widget.SensorID != "trip_distance" || widget.GaugeID != "dashboard_trip_odometer" {
		t.Fatalf("odometer widget identity = %#v", widget)
	}
	if widget.Status != sensors.StatusOK || widget.Scale != 1.5 || widget.GaugeMovement != v3gauges.MovementBell {
		t.Fatalf("odometer widget status/scale/movement = %#v", widget)
	}
	if got := gaugePartSequence(widget); got != "layer:panel,wheel_strip:0,wheel_strip:1,wheel_strip:2,layer:glass" {
		t.Fatalf("odometer part sequence = %q", got)
	}
	wheel := firstPartKind(widget, PartKindWheelStrip)
	if wheel.AssetPath == "" || wheel.Window.Width != 12 || wheel.Window.Height != 20 || wheel.StripOffset == 0 {
		t.Fatalf("wheel strip part = %#v", wheel)
	}
}

func TestRuntimeOdometerGaugeInstantMovementJumpsToTarget(t *testing.T) {
	runtime := testOdometerMovementRuntimeWithConfig(t, v3gauges.MovementInstant)

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "trip_distance",
		State:     okState("trip_distance", 12.0, "km"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "trip_distance",
		State:     okState("trip_distance", 12.9, "km"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected instant odometer update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected instant odometer movement to stay static")
	}

	wheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if len(wheels) != 3 {
		t.Fatalf("wheel parts = %d, want 3", len(wheels))
	}
	if !almostEqual(wheels[0].StripOffset, 20.0) || !almostEqual(wheels[1].StripOffset, 40.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
		t.Fatalf("expected exact target offsets, got %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}
}

func TestRuntimeOdometerGaugeLinearMovementAnimatesAndSettlesExactlyOnTarget(t *testing.T) {
	runtime := testOdometerMovementRuntimeWithConfig(t, v3gauges.MovementLinear)

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "trip_distance",
		State:     okState("trip_distance", 12.0, "km"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "trip_distance",
		State:     okState("trip_distance", 12.9, "km"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected linear odometer movement to become active")
	}
	wheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if !almostEqual(wheels[0].StripOffset, 20.0) || !almostEqual(wheels[1].StripOffset, 40.0) || !almostEqual(wheels[2].StripOffset, 0.0) {
		t.Fatalf("expected movement to start from previous offsets, got %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}

	scenes, changed, err = runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected mid-transition tick to redraw")
	}
	wheels = wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if !almostEqual(wheels[0].StripOffset, 20.0) {
		t.Fatalf("expected unchanged tens wheel to stay on its exact slot, got %.2f", wheels[0].StripOffset)
	}
	if !almostEqual(wheels[1].StripOffset, 40.0) {
		t.Fatalf("expected unchanged ones wheel to stay on its exact slot, got %.2f", wheels[1].StripOffset)
	}
	if !(wheels[2].StripOffset >= 90.0 && wheels[2].StripOffset < 180.0) {
		t.Fatalf("expected tenths wheel to advance between start and final offsets, got %.2f", wheels[2].StripOffset)
	}

	scenes, changed, err = runtime.Tick(start.Add(210 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final tick to redraw")
	}
	wheels = wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if !almostEqual(wheels[0].StripOffset, 20.0) || !almostEqual(wheels[1].StripOffset, 40.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
		t.Fatalf("expected final settled offsets, got %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}

	_, changed, err = runtime.Tick(start.Add(211 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected cleanup tick to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected linear odometer transition to return to static")
	}
}

func TestRuntimeOdometerGaugeEaseOutMovementAdvancesFurtherThanLinearAtSameTick(t *testing.T) {
	linearRuntime := testOdometerMovementRuntimeWithConfig(t, v3gauges.MovementLinear)
	easeOutRuntime := testOdometerMovementRuntimeWithConfig(t, v3gauges.MovementEaseOut)

	start := time.Unix(100, 0)
	for _, runtime := range []*Runtime{linearRuntime, easeOutRuntime} {
		_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 12.0, "km"),
			Timestamp: start,
			ReadAt:    start,
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
		_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 12.9, "km"),
			Timestamp: start.Add(10 * time.Millisecond),
			ReadAt:    start.Add(10 * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
	}

	linearScenes, changed, err := linearRuntime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected linear tick to redraw")
	}
	easeOutScenes, changed, err := easeOutRuntime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected ease_out tick to redraw")
	}

	linearWheels := wheelStripWidgetParts(requireWidget(t, linearScenes[0], "trip"))
	easeOutWheels := wheelStripWidgetParts(requireWidget(t, easeOutScenes[0], "trip"))
	if !(easeOutWheels[2].StripOffset > linearWheels[2].StripOffset) {
		t.Fatalf("expected ease_out to advance further than linear on the moving wheel: linear=%.2f ease_out=%.2f", linearWheels[2].StripOffset, easeOutWheels[2].StripOffset)
	}
}

func TestRuntimeOdometerGaugeEaseOutRolloverMapsDigitZeroAfterNine(t *testing.T) {
	runtime := testOdometerMovementRuntimeWithConfig(t, v3gauges.MovementEaseOut)

	start := time.Unix(100, 0)
	_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "trip_distance",
		State:     okState("trip_distance", 0.9, "km"),
		Timestamp: start,
		ReadAt:    start,
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
		Kind:      sensors.EventKindValueChange,
		SensorID:  "trip_distance",
		State:     okState("trip_distance", 1.0, "km"),
		Timestamp: start.Add(10 * time.Millisecond),
		ReadAt:    start.Add(10 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected ease_out rollover update to redraw")
	}

	initialWheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if got := wheelPartSliceDigits(initialWheels[2]); !intSlicesEqual(got, []int{9}) {
		t.Fatalf("expected initial tenths wheel to show digit 9, got digits=%v", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected mid ease_out rollover tick to redraw")
	}

	midWheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if got := wheelPartSliceDigits(midWheels[2]); !intSlicesEqual(got, []int{9, 0}) {
		t.Fatalf("expected ease_out rollover to render digits 9 and 0 during movement, got %v", got)
	}
}

func TestRuntimeOdometerGaugeFractionalSlicesUseSameAdjacentDigitsAcrossMovementModes(t *testing.T) {
	testCases := []struct {
		name         string
		startValue   float64
		targetValue  float64
		expectedPair []int
	}{
		{name: "tenths_1_to_2", startValue: 123.1, targetValue: 123.2, expectedPair: []int{1, 2}},
		{name: "tenths_8_to_9", startValue: 123.8, targetValue: 123.9, expectedPair: []int{8, 9}},
		{name: "tenths_9_to_0", startValue: 123.9, targetValue: 124.0, expectedPair: []int{9, 0}},
	}
	movements := []string{v3gauges.MovementLinear, v3gauges.MovementEaseOut, v3gauges.MovementBell}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, movement := range movements {
				t.Run(movement, func(t *testing.T) {
					runtime := testOdometerMovementRuntimeWithConfig(t, movement)
					start := time.Unix(100, 0)

					_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
						Kind:      sensors.EventKindValueChange,
						SensorID:  "trip_distance",
						State:     okState("trip_distance", tc.startValue, "km"),
						Timestamp: start,
						ReadAt:    start,
					})
					if err != nil {
						t.Fatalf("ApplyEvent failed: %v", err)
					}

					_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
						Kind:      sensors.EventKindValueChange,
						SensorID:  "trip_distance",
						State:     okState("trip_distance", tc.targetValue, "km"),
						Timestamp: start.Add(10 * time.Millisecond),
						ReadAt:    start.Add(10 * time.Millisecond),
					})
					if err != nil {
						t.Fatalf("ApplyEvent failed: %v", err)
					}

					scenes, changed, err := runtime.Tick(start.Add(110 * time.Millisecond))
					if err != nil {
						t.Fatalf("Tick failed: %v", err)
					}
					if !changed {
						t.Fatalf("expected %s mid-transition tick to redraw", movement)
					}

					wheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
					if got := wheelPartSliceDigits(wheels[2]); !intSlicesEqual(got, tc.expectedPair) {
						t.Fatalf("expected %s to render adjacent digits %v, got %v", movement, tc.expectedPair, got)
					}
				})
			}
		})
	}
}

func TestRuntimeOdometerGaugeBellMovementStartsSlowerThanLinearAndSettlesExactlyOnTarget(t *testing.T) {
	linearRuntime := testOdometerMovementRuntimeWithConfig(t, v3gauges.MovementLinear)
	bellRuntime := testOdometerMovementRuntimeWithConfig(t, v3gauges.MovementBell)
	start := time.Unix(100, 0)
	for _, runtime := range []*Runtime{linearRuntime, bellRuntime} {
		_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 12.0, "km"),
			Timestamp: start,
			ReadAt:    start,
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
		_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 12.9, "km"),
			Timestamp: start.Add(10 * time.Millisecond),
			ReadAt:    start.Add(10 * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
	}

	linearScenes, changed, err := linearRuntime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected linear tick to redraw")
	}
	bellScenes, changed, err := bellRuntime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected bell tick to redraw")
	}
	linearWheels := wheelStripWidgetParts(requireWidget(t, linearScenes[0], "trip"))
	bellWheels := wheelStripWidgetParts(requireWidget(t, bellScenes[0], "trip"))
	if !(bellWheels[2].StripOffset < linearWheels[2].StripOffset) {
		t.Fatalf("expected bell to start slower than linear on the moving wheel: linear=%.2f bell=%.2f", linearWheels[2].StripOffset, bellWheels[2].StripOffset)
	}

	scenes, changed, err := bellRuntime.Tick(start.Add(210 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final bell tick to redraw")
	}
	wheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if !almostEqual(wheels[0].StripOffset, 20.0) || !almostEqual(wheels[1].StripOffset, 40.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
		t.Fatalf("expected final bell offsets, got %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}
}

func TestRuntimeOdometerGaugeCarryDragAdvancesHigherWheelNearRolloverAndSettlesExactlyOnTarget(t *testing.T) {
	baseRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, true, false, false)
	carryRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, true, true, false)

	start := time.Unix(100, 0)
	for _, runtime := range []*Runtime{baseRuntime, carryRuntime} {
		_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 19.9, "km"),
			Timestamp: start,
			ReadAt:    start,
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
		_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 20.0, "km"),
			Timestamp: start.Add(10 * time.Millisecond),
			ReadAt:    start.Add(10 * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
	}

	baseScenes, changed, err := baseRuntime.Tick(start.Add(170 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected base rollover tick to redraw")
	}
	carryScenes, changed, err := carryRuntime.Tick(start.Add(170 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected carry-drag rollover tick to redraw")
	}

	baseWheels := wheelStripWidgetParts(requireWidget(t, baseScenes[0], "trip"))
	carryWheels := wheelStripWidgetParts(requireWidget(t, carryScenes[0], "trip"))
	if !(carryWheels[0].StripOffset > baseWheels[0].StripOffset) {
		t.Fatalf("expected carry_drag to advance tens wheel further than base: base=%.2f carry=%.2f", baseWheels[0].StripOffset, carryWheels[0].StripOffset)
	}
	if !(carryWheels[1].StripOffset > baseWheels[1].StripOffset) {
		t.Fatalf("expected carry_drag to advance ones wheel further than base: base=%.2f carry=%.2f", baseWheels[1].StripOffset, carryWheels[1].StripOffset)
	}

	finalScenes, changed, err := carryRuntime.Tick(start.Add(210 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final carry-drag tick to redraw")
	}
	finalWheels := wheelStripWidgetParts(requireWidget(t, finalScenes[0], "trip"))
	if !almostEqual(finalWheels[0].StripOffset, 40.0) || !almostEqual(finalWheels[1].StripOffset, 0.0) || !almostEqual(finalWheels[2].StripOffset, 0.0) {
		t.Fatalf("expected carry-drag to settle exactly on rollover target, got %.2f/%.2f/%.2f", finalWheels[0].StripOffset, finalWheels[1].StripOffset, finalWheels[2].StripOffset)
	}
}

func TestRuntimeOdometerGaugeCarryDragStraddlingUpdatePullsHigherWheelBeforeRolloverAndStillSettlesExactlyOnTarget(t *testing.T) {
	baseRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, true, false, false)
	carryRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, true, true, false)

	start := time.Unix(100, 0)
	for _, runtime := range []*Runtime{baseRuntime, carryRuntime} {
		_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 19.8, "km"),
			Timestamp: start,
			ReadAt:    start,
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
		_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 20.2, "km"),
			Timestamp: start.Add(10 * time.Millisecond),
			ReadAt:    start.Add(10 * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
	}

	baseScenes, changed, err := baseRuntime.Tick(start.Add(100 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected base straddling tick to redraw")
	}
	carryScenes, changed, err := carryRuntime.Tick(start.Add(100 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected carry straddling tick to redraw")
	}

	baseWheels := wheelStripWidgetParts(requireWidget(t, baseScenes[0], "trip"))
	carryWheels := wheelStripWidgetParts(requireWidget(t, carryScenes[0], "trip"))
	if !(baseWheels[1].StripOffset < 200.0) {
		t.Fatalf("expected lower wheel to still be approaching rollover, got base ones offset %.2f", baseWheels[1].StripOffset)
	}
	if !(carryWheels[0].StripOffset > baseWheels[0].StripOffset) {
		t.Fatalf("expected carry_drag to pull tens wheel before rollover on straddling update: base=%.2f carry=%.2f", baseWheels[0].StripOffset, carryWheels[0].StripOffset)
	}

	finalScenes, changed, err := carryRuntime.Tick(start.Add(210 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final carry-drag straddling tick to redraw")
	}
	finalWheels := wheelStripWidgetParts(requireWidget(t, finalScenes[0], "trip"))
	if !almostEqual(finalWheels[0].StripOffset, 40.0) || !almostEqual(finalWheels[1].StripOffset, 0.0) || !almostEqual(finalWheels[2].StripOffset, 40.0) {
		t.Fatalf("expected carry-drag straddling update to settle exactly on target, got %.2f/%.2f/%.2f", finalWheels[0].StripOffset, finalWheels[1].StripOffset, finalWheels[2].StripOffset)
	}
}

func TestRuntimeOdometerGaugeCarryDragSparseMultiRolloverUpdateDoesNotYankHigherWheelAheadAndStillSettlesExactlyOnTarget(t *testing.T) {
	baseRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, true, false, false)
	carryRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, true, true, false)

	start := time.Unix(100, 0)
	for _, runtime := range []*Runtime{baseRuntime, carryRuntime} {
		_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 20.8, "km"),
			Timestamp: start,
			ReadAt:    start,
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
		_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 31.2, "km"),
			Timestamp: start.Add(10 * time.Millisecond),
			ReadAt:    start.Add(10 * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
	}

	baseScenes, changed, err := baseRuntime.Tick(start.Add(50 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected base sparse carry_drag tick to redraw")
	}
	carryScenes, changed, err := carryRuntime.Tick(start.Add(50 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected carry sparse carry_drag tick to redraw")
	}

	baseWheels := wheelStripWidgetParts(requireWidget(t, baseScenes[0], "trip"))
	carryWheels := wheelStripWidgetParts(requireWidget(t, carryScenes[0], "trip"))
	if !almostEqual(carryWheels[0].StripOffset, baseWheels[0].StripOffset) {
		t.Fatalf("expected sparse multi-rollover update to avoid early tens drag: base=%.2f carry=%.2f", baseWheels[0].StripOffset, carryWheels[0].StripOffset)
	}
	if !almostEqual(carryWheels[1].StripOffset, baseWheels[1].StripOffset) {
		t.Fatalf("expected sparse multi-rollover update to avoid early ones drag: base=%.2f carry=%.2f", baseWheels[1].StripOffset, carryWheels[1].StripOffset)
	}

	finalScenes, changed, err := carryRuntime.Tick(start.Add(210 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final sparse carry-drag tick to redraw")
	}
	finalWheels := wheelStripWidgetParts(requireWidget(t, finalScenes[0], "trip"))
	if !almostEqual(finalWheels[0].StripOffset, 60.0) || !almostEqual(finalWheels[1].StripOffset, 20.0) || !almostEqual(finalWheels[2].StripOffset, 40.0) {
		t.Fatalf("expected sparse multi-rollover update to settle exactly on target, got %.2f/%.2f/%.2f", finalWheels[0].StripOffset, finalWheels[1].StripOffset, finalWheels[2].StripOffset)
	}
}

func TestRuntimeOdometerGaugeSnapSettleAddsShortTailAndSettlesExactlyOnTarget(t *testing.T) {
	baseRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, false, false, false)
	settleRuntime := testOdometerMovementRuntimeWithRealism(t, v3gauges.MovementLinear, false, false, true)

	start := time.Unix(100, 0)
	for _, runtime := range []*Runtime{baseRuntime, settleRuntime} {
		_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 12.0, "km"),
			Timestamp: start,
			ReadAt:    start,
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
		_, _, err = runtime.ApplyEvent(sensors.SensorEvent{
			Kind:      sensors.EventKindValueChange,
			SensorID:  "trip_distance",
			State:     okState("trip_distance", 12.9, "km"),
			Timestamp: start.Add(10 * time.Millisecond),
			ReadAt:    start.Add(10 * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("ApplyEvent failed: %v", err)
		}
	}

	baseScenes, changed, err := baseRuntime.Tick(start.Add(230 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected baseline cleanup tick to redraw")
	}
	settleScenes, changed, err := settleRuntime.Tick(start.Add(230 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected snap_settle tail tick to redraw")
	}
	if !settleRuntime.HasActiveMovement() {
		t.Fatalf("expected snap_settle tail to remain active after main travel completes")
	}

	baseWheels := wheelStripWidgetParts(requireWidget(t, baseScenes[0], "trip"))
	settleWheels := wheelStripWidgetParts(requireWidget(t, settleScenes[0], "trip"))
	if !almostEqual(settleWheels[0].StripOffset, baseWheels[0].StripOffset) {
		t.Fatalf("expected snap_settle to leave unchanged tens wheel alone: base=%.2f settle=%.2f", baseWheels[0].StripOffset, settleWheels[0].StripOffset)
	}
	if !almostEqual(settleWheels[1].StripOffset, baseWheels[1].StripOffset) {
		t.Fatalf("expected snap_settle to leave unchanged ones wheel alone: base=%.2f settle=%.2f", baseWheels[1].StripOffset, settleWheels[1].StripOffset)
	}
	if !(settleWheels[2].StripOffset > baseWheels[2].StripOffset) {
		t.Fatalf("expected snap_settle to nudge tenths wheel beyond target: base=%.2f settle=%.2f", baseWheels[2].StripOffset, settleWheels[2].StripOffset)
	}

	finalScenes, changed, err := settleRuntime.Tick(start.Add(270 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final snap_settle tick to redraw")
	}
	finalWheels := wheelStripWidgetParts(requireWidget(t, finalScenes[0], "trip"))
	if !almostEqual(finalWheels[0].StripOffset, 20.0) || !almostEqual(finalWheels[1].StripOffset, 40.0) || !almostEqual(finalWheels[2].StripOffset, 180.0) {
		t.Fatalf("expected snap_settle to settle exactly on target, got %.2f/%.2f/%.2f", finalWheels[0].StripOffset, finalWheels[1].StripOffset, finalWheels[2].StripOffset)
	}

	_, changed, err = settleRuntime.Tick(start.Add(271 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected cleanup tick after snap_settle")
	}
	if settleRuntime.HasActiveMovement() {
		t.Fatalf("expected snap_settle lifecycle to return to static")
	}
}

func TestRuntimeOdometerGaugeRecognizedMovementFallbacksStayInstant(t *testing.T) {
	for _, movement := range []string{v3gauges.MovementSmooth, v3gauges.MovementClick} {
		t.Run(movement, func(t *testing.T) {
			runtime := testOdometerMovementRuntimeWithConfig(t, movement)

			start := time.Unix(100, 0)
			_, _, err := runtime.ApplyEvent(sensors.SensorEvent{
				Kind:      sensors.EventKindValueChange,
				SensorID:  "trip_distance",
				State:     okState("trip_distance", 12.0, "km"),
				Timestamp: start,
				ReadAt:    start,
			})
			if err != nil {
				t.Fatalf("ApplyEvent failed: %v", err)
			}

			scenes, changed, err := runtime.ApplyEvent(sensors.SensorEvent{
				Kind:      sensors.EventKindValueChange,
				SensorID:  "trip_distance",
				State:     okState("trip_distance", 12.9, "km"),
				Timestamp: start.Add(10 * time.Millisecond),
				ReadAt:    start.Add(10 * time.Millisecond),
			})
			if err != nil {
				t.Fatalf("ApplyEvent failed: %v", err)
			}
			if !changed {
				t.Fatalf("expected fallback instant update to redraw")
			}
			if runtime.HasActiveMovement() {
				t.Fatalf("expected fallback instant movement to stay static")
			}
			wheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
			if !almostEqual(wheels[0].StripOffset, 20.0) || !almostEqual(wheels[1].StripOffset, 40.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
				t.Fatalf("expected instant fallback offsets, got %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
			}
		})
	}
}

func TestRuntimeOdometerGaugeWidgetSceneSignatureChangesWithSmoothOffset(t *testing.T) {
	packageDir := makeDashboardOdometerGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "trip", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensorEvent("trip_distance", okState("trip_distance", 12.1, "km")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first odometer event to change rendered output")
	}
	_, changed, err = runtime.ApplyEvent(sensorEvent("trip_distance", okState("trip_distance", 12.2, "km")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed smooth odometer offset to redraw")
	}
}

func TestRuntimeLoadsIndicatorGaugeWidgetPackage(t *testing.T) {
	packageDir := makeDashboardIndicatorGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "check_engine", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{40, 50}, Scale: 1.5}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("check_engine", 1, ""))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "check_engine")
	if widget.Type != v3config.WidgetTypeGauge || widget.SensorID != "check_engine" || widget.GaugeID != "dashboard_check_engine_indicator" {
		t.Fatalf("indicator widget identity = %#v", widget)
	}
	if widget.Status != sensors.StatusOK || widget.Scale != 1.5 {
		t.Fatalf("indicator widget status/scale = %#v", widget)
	}
	if got := gaugePartSequence(widget); got != "layer:bezel,layer:face,layer:on,layer:glass" {
		t.Fatalf("indicator part sequence = %q", got)
	}
}

func TestRuntimeIndicatorThermalFadeDisabledRemainsImmediate(t *testing.T) {
	packageDir := makeDashboardIndicatorGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 320, Height: 240}, Widgets: []v3config.WidgetConfig{{ID: "check_engine", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensorEvent("check_engine", okState("check_engine", 1, "")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected immediate indicator change to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected disabled thermal fade to stay static")
	}
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	on := firstGaugeLayerPart(requireWidget(t, scenes[0], "check_engine"), "on")
	if !almostEqual(on.Alpha, 0) {
		t.Fatalf("expected immediate on layer to use default opaque alpha, got %v", on.Alpha)
	}
}

func TestRuntimeIndicatorThermalFadeOnTransition(t *testing.T) {
	rise, fall := 100, 240
	packageDir := makeDashboardIndicatorGaugePackageWithThermalFade(t, &rise, &fall)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 320, Height: 240}, Widgets: []v3config.WidgetConfig{{ID: "check_engine", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	now := time.Unix(200, 0)
	runtime.clock = func() time.Time { return now }
	runtime.SetState(okState("check_engine", 0, ""))
	if _, err := runtime.Snapshot(); err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("check_engine", okState("check_engine", 1, ""), now))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected thermal fade start to mark the widget changed")
	}
	if !runtime.HasActiveMovement() {
		t.Fatalf("expected thermal fade on transition to stay active")
	}
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	startWidget := requireWidget(t, scenes[0], "check_engine")
	if got := gaugePartSequence(startWidget); got != "layer:bezel,layer:face,layer:off,layer:glass" {
		t.Fatalf("expected fade start to begin from the cold lamp face, got %q", got)
	}

	midScenes, changed, err := runtime.Tick(now.Add(50 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected mid-fade tick to redraw")
	}
	midWidget := requireWidget(t, midScenes[0], "check_engine")
	if got := gaugePartSequence(midWidget); got != "layer:bezel,layer:face,layer:off,layer:on,layer:glass" {
		t.Fatalf("expected mid-fade layers to include off and on, got %q", got)
	}
	on := firstGaugeLayerPart(midWidget, "on")
	if !(on.Alpha > 0 && on.Alpha < 1) {
		t.Fatalf("expected mid-fade alpha in (0,1), got %v", on.Alpha)
	}

	finalScenes, changed, err := runtime.Tick(now.Add(100 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final fade tick to redraw")
	}
	finalWidget := requireWidget(t, finalScenes[0], "check_engine")
	if got := gaugePartSequence(finalWidget); got != "layer:bezel,layer:face,layer:on,layer:glass" {
		t.Fatalf("expected final on layers, got %q", got)
	}
	if !almostEqual(firstGaugeLayerPart(finalWidget, "on").Alpha, 1) {
		t.Fatalf("expected final on alpha to settle at 1, got %v", firstGaugeLayerPart(finalWidget, "on").Alpha)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected thermal fade to stop ticking once settled on")
	}
}

func TestRuntimeIndicatorThermalFadeOffTransition(t *testing.T) {
	rise, fall := 100, 240
	packageDir := makeDashboardIndicatorGaugePackageWithThermalFade(t, &rise, &fall)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 320, Height: 240}, Widgets: []v3config.WidgetConfig{{ID: "check_engine", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	now := time.Unix(300, 0)
	runtime.clock = func() time.Time { return now }
	runtime.SetState(okState("check_engine", 1, ""))
	if _, err := runtime.Snapshot(); err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensorEventAt("check_engine", okState("check_engine", 0, ""), now))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected thermal fade off transition to mark the widget changed")
	}
	if !runtime.HasActiveMovement() {
		t.Fatalf("expected thermal fade off transition to stay active")
	}

	midScenes, changed, err := runtime.Tick(now.Add(120 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected mid-fade-off tick to redraw")
	}
	midWidget := requireWidget(t, midScenes[0], "check_engine")
	if got := gaugePartSequence(midWidget); got != "layer:bezel,layer:face,layer:off,layer:on,layer:glass" {
		t.Fatalf("expected mid-fade-off layers to include off and on, got %q", got)
	}
	on := firstGaugeLayerPart(midWidget, "on")
	if !(on.Alpha > 0 && on.Alpha < 1) {
		t.Fatalf("expected mid-fade-off alpha in (0,1), got %v", on.Alpha)
	}

	finalScenes, changed, err := runtime.Tick(now.Add(240 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final fade-off tick to redraw")
	}
	finalWidget := requireWidget(t, finalScenes[0], "check_engine")
	if got := gaugePartSequence(finalWidget); got != "layer:bezel,layer:face,layer:off,layer:glass" {
		t.Fatalf("expected final off layers, got %q", got)
	}
	if on := firstGaugeLayerPart(finalWidget, "on"); on.AssetPath != "" {
		t.Fatalf("expected on layer to be gone after settling off, got %#v", on)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected thermal fade to stop ticking once settled off")
	}
}

func TestRuntimeLoadsBarGaugeWidgetPackageAndRendersLevelReveal(t *testing.T) {
	packageDir := makeDashboardBarGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{120, 80}, Scale: 1.5}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("coolant_temperature", 80, "c"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "coolant")
	if widget.Type != v3config.WidgetTypeGauge || widget.SensorID != "coolant_temperature" || widget.GaugeID != "dashboard_coolant_bar" {
		t.Fatalf("bar widget identity = %#v", widget)
	}
	if widget.Status != sensors.StatusOK || widget.Scale != 1.5 || widget.GaugeBarMode != "level" || widget.GaugeBarAxis != "vertical" || widget.GaugeBarOrigin != "bottom" {
		t.Fatalf("bar widget status/scale/config = %#v", widget)
	}
	if len(widget.GaugeBarBounds) != 4 || widget.GaugeBarBounds[0] != 40 || widget.GaugeBarBounds[3] != 180 {
		t.Fatalf("bar widget bounds = %#v", widget.GaugeBarBounds)
	}
	if got := gaugePartSequence(widget); got != "layer:panel,bar:level,layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	bar := firstPartKind(widget, PartKindBar)
	if bar.AssetPath == "" || len(bar.Position) != 2 || len(bar.Source) != 2 || bar.Window.Width != 24 || bar.Window.Height != 90 {
		t.Fatalf("bar part = %#v", bar)
	}
	if bar.Position[0] != 40 || bar.Position[1] != 110 {
		t.Fatalf("bar position = %#v, want [40 110]", bar.Position)
	}
	if bar.Source[0] != 40 || bar.Source[1] != 110 {
		t.Fatalf("bar source = %#v, want [40 110]", bar.Source)
	}
}

func TestRuntimeBarGaugeWidgetKeepsPointerMarkersHiddenWhenDisabled(t *testing.T) {
	packageDir := makeDashboardBarGaugePackageWithPointerMarkers(t, "    max: false\n    min: false\n")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("coolant_temperature", 80, "c"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	widget := requireWidget(t, scenes[0], "coolant")
	if got := gaugePartSequence(widget); got != "layer:panel,bar:level,layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	if got := countParts(widget, PartKindMarkerMin); got != 0 {
		t.Fatalf("expected no min marker parts, got %d", got)
	}
	if got := countParts(widget, PartKindMarkerMax); got != 0 {
		t.Fatalf("expected no max marker parts, got %d", got)
	}
}

func TestRuntimeBarGaugeWidgetRendersPointerMarkersMinOnly(t *testing.T) {
	packageDir := makeDashboardBarGaugePackageWithPointerMarkers(t, "    min: true\n")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("coolant_temperature", 80, "c"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	widget := requireWidget(t, scenes[0], "coolant")
	if got := gaugePartSequence(widget); got != "layer:panel,bar:level,marker_min:[40 110],layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	marker := firstPartKind(widget, PartKindMarkerMin)
	if marker.Layer != "marker_min" || marker.AssetPath == "" || !intSlicesEqual(marker.Position, []int{40, 110}) {
		t.Fatalf("min marker part = %#v", marker)
	}
	if got := countParts(widget, PartKindMarkerMax); got != 0 {
		t.Fatalf("expected no max marker parts, got %d", got)
	}
}

func TestRuntimeBarGaugeWidgetRendersPointerMarkersMaxOnly(t *testing.T) {
	packageDir := makeDashboardBarGaugePackageWithPointerMarkers(t, "    max: true\n")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("coolant_temperature", 80, "c"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	widget := requireWidget(t, scenes[0], "coolant")
	if got := gaugePartSequence(widget); got != "layer:panel,bar:level,marker_max:[40 110],layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	marker := firstPartKind(widget, PartKindMarkerMax)
	if marker.Layer != "marker_max" || marker.AssetPath == "" || !intSlicesEqual(marker.Position, []int{40, 110}) {
		t.Fatalf("max marker part = %#v", marker)
	}
	if got := countParts(widget, PartKindMarkerMin); got != 0 {
		t.Fatalf("expected no min marker parts, got %d", got)
	}
}

func TestRuntimeBarGaugeWidgetRendersPointerMarkersAboveBarBeforeGlass(t *testing.T) {
	packageDir := makeDashboardBarGaugePackageWithPointerMarkers(t, "    max: true\n    min: true\n")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	_, _, err = runtime.ApplyEvent(sensorEvent("coolant_temperature", okState("coolant_temperature", 60, "c")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEvent("coolant_temperature", okState("coolant_temperature", 100, "c")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected pointer marker update to redraw")
	}

	widget := requireWidget(t, scenes[0], "coolant")
	if got := gaugePartSequence(widget); got != "layer:panel,bar:level,marker_min:[40 155],marker_max:[40 65],layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	if got := countParts(widget, PartKindMarkerMin); got != 1 {
		t.Fatalf("expected one min marker part, got %d", got)
	}
	if got := countParts(widget, PartKindMarkerMax); got != 1 {
		t.Fatalf("expected one max marker part, got %d", got)
	}
}

func TestRuntimeBarGaugeWidgetKeepsAveragePointerMarkerHiddenWhenDisabled(t *testing.T) {
	packageDir := makeDashboardBarGaugePackageWithPointerMarkers(t, "    average: false\n")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	runtime.SetState(okState("coolant_temperature", 80, "c"))
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "coolant")
	if got := gaugePartSequence(widget); got != "layer:panel,bar:level,layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	if got := countParts(widget, PartKindMarkerAverage); got != 0 {
		t.Fatalf("expected no average marker parts, got %d", got)
	}
}

func TestRuntimeBarGaugeWidgetRendersAveragePointerMarkerAboveBarBeforeGlass(t *testing.T) {
	packageDir := makeDashboardBarGaugePackageWithPointerMarkers(t, "    max: true\n    min: true\n    average: true\n")
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	start := time.Unix(710, 0)
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 120, "c"), start.Add(10*time.Second)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	widget := requireWidget(t, scenes[0], "coolant")
	if got := gaugePartSequence(widget); got != "layer:panel,bar:level,marker_min:[40 200],marker_max:[40 20],marker_average:[40 86],layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	marker := firstPartKind(widget, PartKindMarkerAverage)
	if marker.Layer != "marker_average" || marker.AssetPath == "" || !intSlicesEqual(marker.Position, []int{40, 86}) {
		t.Fatalf("average marker part = %#v", marker)
	}
}

func TestRuntimeBarGaugeWidgetSceneSignatureChangesWithRevealHeight(t *testing.T) {
	packageDir := makeDashboardBarGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{120, 80}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	_, changed, err := runtime.ApplyEvent(sensorEvent("coolant_temperature", okState("coolant_temperature", 80, "c")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first bar event to change rendered output")
	}
	_, changed, err = runtime.ApplyEvent(sensorEvent("coolant_temperature", okState("coolant_temperature", 81, "c")))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed bar reveal to redraw")
	}
}

func TestRuntimeBarGaugeDampingDisabledRemainsImmediate(t *testing.T) {
	runtime := testBarMovementRuntimeWithDamping(t, "")
	start := time.Unix(400, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected immediate bar change to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected disabled bar damping to stay static")
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 90 {
		t.Fatalf("expected immediate reveal height 90, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingAnimatesRisingReveal(t *testing.T) {
	runtime := testBarMovementRuntimeWithDamping(t, "  damping: true\n")
	start := time.Unix(500, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		if !current.DampingEnabled {
			t.Fatalf("expected damping-enabled movement state, got %#v", current)
		}
		return 200 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected rising bar damping to start active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.Policy != v3gauges.MovementPolicyLinear {
		t.Fatalf("expected bar damping to default to linear movement curve, got %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected active bar damping tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-60) > 0.001 {
		t.Fatalf("expected rising bar midpoint display 60, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 45 {
		t.Fatalf("expected rising midpoint reveal height 45, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingAnimatesFallingReveal(t *testing.T) {
	runtime := testBarMovementRuntimeWithDamping(t, "  damping:\n    rise_ms: 100\n    fall_ms: 240\n")
	start := time.Unix(600, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 120, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected falling bar damping to start active movement")
	}

	scenes, changed, err := runtime.Tick(start.Add(130 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected falling bar damping tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-100) > 0.001 {
		t.Fatalf("expected falling bar midpoint display 100, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 135 {
		t.Fatalf("expected falling midpoint reveal height 135, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingSettlesAtFinalReveal(t *testing.T) {
	runtime := testBarMovementRuntimeWithDamping(t, "  damping:\n    rise_ms: 100\n    fall_ms: 240\n")
	start := time.Unix(700, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected damping start to change rendered output")
	}

	scenes, changed, err := runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final damping tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 80 || movement.TargetValue != 80 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected rising bar damping to settle exactly at the target, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 90 {
		t.Fatalf("expected settled reveal height 90, got %d", got)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected settled bar damping to stop ticking")
	}
}

func TestRuntimeBarGaugeOvershootDefaultDisabledStaysImmediate(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "", 0, 100, true)
	start := time.Unix(705, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected immediate bar update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected default bar overshoot-disabled update to stay immediate")
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 144 {
		t.Fatalf("expected immediate reveal height 144, got %d", got)
	}
}

func TestRuntimeBarGaugeOvershootAnimatesRisingReveal(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  overshoot: {}\n", 0, 100, true)
	start := time.Unix(706, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected overshoot-only rising bar movement to become active")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DampingEnabled {
		t.Fatalf("expected overshoot-only bar movement to keep damping disabled: %#v", movement)
	}
	if !movement.OvershootEnabled || movement.Duration <= 0 || movement.Phase != movementPhaseMoving {
		t.Fatalf("expected overshoot-only rising bar movement to schedule animation: %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(140 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot-only rising tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue <= 50 || movement.DisplayValue > 100 {
		t.Fatalf("expected rising overshoot above target and within range, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got <= 90 || got > 180 {
		t.Fatalf("expected rising overshoot reveal height between target and max, got %d", got)
	}
}

func TestRuntimeBarGaugeOvershootAnimatesFallingReveal(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  overshoot: {}\n", 0, 100, true)
	start := time.Unix(707, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 100, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected overshoot-only falling bar movement to become active")
	}

	scenes, changed, err := runtime.Tick(start.Add(140 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot-only falling tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue >= 50 || movement.DisplayValue < 0 {
		t.Fatalf("expected falling overshoot below target and within range, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got >= 90 {
		t.Fatalf("expected falling overshoot reveal height below target, got %d", got)
	}
}

func TestRuntimeBarGaugeOvershootStaysBoundedAndSettlesOnTarget(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  overshoot: {}\n", 0, 100, true)
	start := time.Unix(708, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.Tick(start.Add(220 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected bounded overshoot tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue <= 80 || movement.DisplayValue > 100 {
		t.Fatalf("expected bounded overshoot above target and within range, got %#v", movement)
	}
	if movement.OvershootTargetValue <= 80 || movement.OvershootTargetValue > 100 {
		t.Fatalf("unexpected bounded bar overshoot target %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got <= 144 || got > 180 {
		t.Fatalf("expected bounded overshoot reveal height between target and max, got %d", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot settle tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 80 || movement.TargetValue != 80 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected bounded bar overshoot to settle exactly on target, got %#v", movement)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected settled bar overshoot to stop ticking")
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 144 {
		t.Fatalf("expected settled reveal height 144, got %d", got)
	}
}

func TestRuntimeBarGaugeOvershootSettlesAtFinalReveal(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  overshoot: {}\n", 0, 100, true)
	start := time.Unix(709, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.Tick(start.Add(220 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final overshoot settle tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 50 || movement.TargetValue != 50 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected bar overshoot to settle exactly at final target, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 90 {
		t.Fatalf("expected settled reveal height 90, got %d", got)
	}
}

func TestRuntimeBarGaugeOvershootLeavesStoredSourceValueUnchanged(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  overshoot: {}\n", 0, 100, true)
	start := time.Unix(709, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected overshoot source update to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 50 {
		t.Fatalf("expected stored source value 50 to remain unchanged, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(140 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected active overshoot tick to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 50 {
		t.Fatalf("expected stored source value 50 to remain unchanged during overshoot, got %v", got)
	}
}

func TestRuntimeBarGaugeHysteresisDefaultDisabledStaysImmediate(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "", 0, 100, true)
	start := time.Unix(709, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected non-hysteresis bar update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected disabled bar hysteresis to avoid active movement")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.HasValue {
		t.Fatalf("unexpected disabled bar hysteresis state: %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 90 {
		t.Fatalf("expected disabled bar hysteresis reveal height 90, got %d", got)
	}
}

func TestRuntimeBarGaugeHysteresisOffsetsRisingApproach(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  hysteresis: true\n", 0, 100, true)
	start := time.Unix(710, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected rising bar hysteresis update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected hysteresis-only bar update to stay immediate")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if !movement.HysteresisEnabled || movement.ApproachDirection != 1 || movement.RawTargetValue != 50 || movement.DisplayValue <= 50 {
		t.Fatalf("unexpected rising bar hysteresis state: %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 92 {
		t.Fatalf("expected rising bar hysteresis reveal height 92, got %d", got)
	}
}

func TestRuntimeBarGaugeHysteresisOffsetsFallingApproach(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  hysteresis: true\n", 0, 100, true)
	start := time.Unix(711, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 100, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected falling bar hysteresis update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected hysteresis-only falling bar update to stay immediate")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if !movement.HysteresisEnabled || movement.ApproachDirection != -1 || movement.RawTargetValue != 50 || movement.DisplayValue >= 50 {
		t.Fatalf("unexpected falling bar hysteresis state: %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 88 {
		t.Fatalf("expected falling bar hysteresis reveal height 88, got %d", got)
	}
}

func TestRuntimeBarGaugeHysteresisLeavesStoredSourceValueUnchanged(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  hysteresis: true\n", 0, 100, true)
	start := time.Unix(712, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected bar hysteresis source update to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 50 {
		t.Fatalf("expected stored source value 50 to remain unchanged, got %v", got)
	}

	_, changed, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(20*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected unchanged bar hysteresis source value to avoid redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 50 {
		t.Fatalf("expected stored source value 50 to remain unchanged after repeated input, got %v", got)
	}
}

func TestRuntimeBarGaugeHysteresisClampTrueStaysWithinValueMapBounds(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  hysteresis: true\n", 0, 100, true)
	start := time.Unix(713, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=true bar hysteresis update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected clamp=true bar hysteresis-only update to stay immediate")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != 200 || movement.TargetValue != 100 || movement.DisplayValue != 100 {
		t.Fatalf("expected clamp=true bar hysteresis display target to stay within min/max, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected clamp=true bar hysteresis reveal height 180 at displayed max, got %d", got)
	}
}

func TestRuntimeBarGaugeHysteresisClampFalsePreservesOutOfRangeShiftedTarget(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  hysteresis: true\n", 0, 100, false)
	start := time.Unix(714, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=false bar hysteresis update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected clamp=false bar hysteresis-only update to stay immediate")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != 200 || movement.TargetValue != 201 || movement.DisplayValue != 201 {
		t.Fatalf("expected clamp=false bar hysteresis to preserve out-of-range shifted target, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected clamp=false bar hysteresis reveal height to remain capped by bar geometry, got %d", got)
	}
}

func TestRuntimeBarGaugeHysteresisClampFalseDampingPreservesUnclampedTarget(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  hysteresis: true\n  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, false)
	start := time.Unix(715, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected clamp=false bar hysteresis+damping to start active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != 200 || movement.TargetValue != 201 {
		t.Fatalf("expected clamp=false bar hysteresis+damping to keep unclamped shifted target, got %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=false bar hysteresis+damping midpoint tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-100.5) > 0.001 {
		t.Fatalf("expected clamp=false bar hysteresis+damping midpoint display 100.5, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected clamp=false bar hysteresis+damping reveal height to stay capped by bar geometry, got %d", got)
	}
}

func TestRuntimeBarGaugeHysteresisClampFalseLeavesStoredSourceValueRaw(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  hysteresis: true\n  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, false)
	start := time.Unix(716, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=false bar hysteresis raw source update to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 200 {
		t.Fatalf("expected stored source value 200 to remain raw, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=false bar hysteresis active tick to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 200 {
		t.Fatalf("expected stored source value 200 to remain unchanged during unclamped hysteresis movement, got %v", got)
	}
}

func TestRuntimeBarGaugeStictionBelowThresholdHoldsDisplay(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  stiction: 15\n", 0, 100, true)
	start := time.Unix(717, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed || scenes != nil {
		t.Fatalf("expected below-threshold bar stiction to suppress redraw, changed=%v scenes=%#v", changed, scenes)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected below-threshold bar stiction to remain static")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 40 || movement.TargetValue != 50 || movement.RawTargetValue != 50 || movement.Phase != movementPhaseStatic {
		t.Fatalf("unexpected held bar stiction state: %#v", movement)
	}
}

func TestRuntimeBarGaugeStictionReleasesAboveThreshold(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  stiction: 15\n", 0, 100, true)
	start := time.Unix(718, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 60, "c"), start.Add(20*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || scenes == nil {
		t.Fatalf("expected above-threshold bar stiction to release and redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected stiction-only bar release to jump without active damping")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 60 || movement.TargetValue != 60 || movement.RawTargetValue != 60 || movement.Phase != movementPhaseStatic {
		t.Fatalf("unexpected released bar stiction state: %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 108 {
		t.Fatalf("expected released bar stiction reveal height 108, got %d", got)
	}
}

func TestRuntimeBarGaugeStictionLargeChangeMovesImmediately(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  stiction: 15\n", 0, 100, true)
	start := time.Unix(719, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || scenes == nil {
		t.Fatalf("expected large bar stiction change to redraw immediately")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected large stiction-only bar change to avoid active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 40 || movement.TargetValue != 40 || movement.RawTargetValue != 40 || movement.Phase != movementPhaseStatic {
		t.Fatalf("unexpected large-change bar stiction state: %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 72 {
		t.Fatalf("expected large-change bar stiction reveal height 72, got %d", got)
	}
}

func TestRuntimeBarGaugeStictionDefaultDisabledDoesNotHoldSmallChanges(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "", 0, 100, true)
	start := time.Unix(720, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || scenes == nil {
		t.Fatalf("expected default bar behavior to redraw immediately without stiction")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected default bar behavior to remain immediate")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.HasValue {
		t.Fatalf("unexpected default non-stiction bar state: %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 90 {
		t.Fatalf("expected default non-stiction reveal height 90, got %d", got)
	}
}

func TestRuntimeBarGaugeStictionDampingSettlesAfterRelease(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  stiction: 15\n  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(721, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 60, "c"), start.Add(20*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected bar stiction release with damping to start active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.PreviousDisplayValue != 40 || movement.DisplayValue != 40 || movement.TargetValue != 60 || movement.RawTargetValue != 60 || movement.Duration != 100*time.Millisecond {
		t.Fatalf("unexpected released damped bar stiction state: %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(70 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected damped bar stiction midpoint tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-50) > 0.001 {
		t.Fatalf("expected damped bar stiction midpoint display 50, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 90 {
		t.Fatalf("expected damped bar stiction midpoint reveal height 90, got %d", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(120 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected damped bar stiction settle tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 60 || movement.TargetValue != 60 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected damped bar stiction to settle cleanly, got %#v", movement)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected damped bar stiction to finish active movement after settlement")
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 108 {
		t.Fatalf("expected damped bar stiction settled reveal height 108, got %d", got)
	}
}

func TestRuntimeBarGaugeStictionLeavesStoredSourceValueUnchanged(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  stiction: 15\n", 0, 100, true)
	start := time.Unix(722, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 40, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected held bar stiction update not to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 50 {
		t.Fatalf("expected stored source value 50 to remain raw while display is held, got %v", got)
	}

	_, changed, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 50, "c"), start.Add(20*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected unchanged held bar stiction source value to avoid redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 50 {
		t.Fatalf("expected stored source value 50 to remain unchanged after repeated held input, got %v", got)
	}
}

func TestRuntimeBarGaugePegBounceDefaultDisabledSettlesImmediatelyAtStop(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "", 0, 100, true)
	start := time.Unix(723, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 20, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 100, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || scenes == nil {
		t.Fatalf("expected default bar stop update to redraw")
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected default bar stop update to avoid peg-bounce animation")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.HasValue {
		t.Fatalf("unexpected default peg-bounce-disabled bar state: %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected default stop reveal height 180, got %d", got)
	}
}

func TestRuntimeBarGaugePegBounceAtMaxStopSettlesBackToLimit(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  peg_bounce: true\n", 0, 100, true)
	start := time.Unix(724, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 20, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 100, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected max-stop bar peg bounce to animate")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if !movement.PegBounceEnabled || !movement.PegBounceActive || movement.TargetValue != 100 || movement.PegBounceStopValue != 100 || movement.PegBounceReboundValue >= 100 {
		t.Fatalf("expected max-stop bar peg bounce to schedule inward rebound, got %#v", movement)
	}
	if movement.SettleDuration < defaultRadialPegBounceMinSettleDuration {
		t.Fatalf("expected visible bar peg-bounce settle duration, got %#v", movement)
	}
	if movement.PegBounceStopValue-movement.PegBounceReboundValue > (100*defaultRadialPegBounceSpanRatio)+0.001 {
		t.Fatalf("expected bounded max-stop bar peg-bounce amplitude, got %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(280 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected max-stop bar peg bounce settle tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue >= 100 || movement.DisplayValue <= movement.PegBounceReboundValue {
		t.Fatalf("expected max-stop bar peg bounce to rebound slightly below the stop, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got >= 180 || got <= 176 {
		t.Fatalf("expected max-stop bar peg bounce height near the stop, got %d", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected max-stop bar peg bounce completion tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 100 || movement.TargetValue != 100 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected max-stop bar peg bounce to settle exactly at the stop, got %#v", movement)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected settled max-stop bar peg bounce to stop ticking")
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected settled max-stop reveal height 180, got %d", got)
	}
}

func TestRuntimeBarGaugePegBounceAtMinStopSettlesBackToLimit(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  peg_bounce: true\n", 0, 100, true)
	start := time.Unix(725, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected min-stop bar peg bounce to animate")
	}

	scenes, changed, err := runtime.Tick(start.Add(280 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected min-stop bar peg bounce settle tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.PegBounceStopValue != 0 || movement.PegBounceReboundValue <= 0 || movement.DisplayValue <= 0 {
		t.Fatalf("expected min-stop bar peg bounce to rebound above zero, got %#v", movement)
	}
	if movement.PegBounceReboundValue-movement.PegBounceStopValue > (100*defaultRadialPegBounceSpanRatio)+0.001 {
		t.Fatalf("expected bounded min-stop bar peg-bounce amplitude, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got <= 0 || got >= 4 {
		t.Fatalf("expected min-stop bar peg bounce height just above zero, got %d", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected min-stop bar peg bounce completion tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 0 || movement.TargetValue != 0 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected min-stop bar peg bounce to settle exactly at the stop, got %#v", movement)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected settled min-stop bar peg bounce to stop ticking")
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 0 {
		t.Fatalf("expected settled min-stop reveal height 0, got %d", got)
	}
}

func TestRuntimeBarGaugePegBounceDoesNotTriggerForInRangeTarget(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  peg_bounce: true\n", 0, 100, true)
	start := time.Unix(726, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 20, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	scenes, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || scenes == nil {
		t.Fatalf("expected in-range bar peg-bounce-only change to redraw immediately")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if !movement.PegBounceEnabled || movement.PegBounceActive || movement.PegBounceReboundValue != 0 || movement.PegBounceStopValue != 0 {
		t.Fatalf("expected in-range bar movement to avoid scheduling peg bounce, got %#v", movement)
	}
	if movement.DisplayValue != 80 || movement.TargetValue != 80 || movement.Phase != movementPhaseStatic || runtime.HasActiveMovement() {
		t.Fatalf("expected in-range bar peg-bounce-only change to settle immediately with no active movement, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 144 {
		t.Fatalf("expected in-range bar reveal height 144, got %d", got)
	}

	_, changed, err = runtime.Tick(start.Add(200 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if changed {
		t.Fatalf("expected no follow-up animation tick for in-range bar peg-bounce-only change")
	}
}

func TestRuntimeBarGaugePegBounceLeavesStoredSourceValueUnchanged(t *testing.T) {
	runtime := testBarMovementRuntimeWithRealismAndValueMap(t, "  peg_bounce: true\n", 0, 100, true)
	start := time.Unix(727, 0)
	runtime.movementPlanner = func(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
		return 300 * time.Millisecond
	}

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 20, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected max-stop peg-bounce source update to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 200 {
		t.Fatalf("expected stored source value 200 to remain raw, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(280 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected active max-stop peg-bounce tick to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 200 {
		t.Fatalf("expected stored source value 200 to remain unchanged during peg bounce, got %v", got)
	}
}

func TestRuntimeBarGaugeDampingClampTrueTargetsDisplayedMaxValue(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(710, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected clamped bar damping to start active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != 200 || movement.TargetValue != 100 || movement.DisplayValue != 0 {
		t.Fatalf("expected clamp=true bar damping to animate toward displayed max only, got %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamped bar damping midpoint tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-50) > 0.001 {
		t.Fatalf("expected clamp=true midpoint display 50, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 90 {
		t.Fatalf("expected clamp=true midpoint reveal height 90, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingClampTrueFallsFromDisplayedMaxValue(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(720, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=true rise settle tick to redraw")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 100 || movement.TargetValue != 100 || movement.RawTargetValue != 200 {
		t.Fatalf("expected clamp=true settled state to stay at displayed max, got %#v", movement)
	}

	_, changed, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 80, "c"), start.Add(120*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected clamp=true fall from displayed max to start active movement")
	}

	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.PreviousDisplayValue != 100 || movement.TargetValue != 80 || movement.RawTargetValue != 80 {
		t.Fatalf("expected clamp=true fall to start from displayed max 100, got %#v", movement)
	}

	var scenes []Scene
	scenes, changed, err = runtime.Tick(start.Add(170 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=true falling midpoint tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-90) > 0.001 {
		t.Fatalf("expected clamp=true falling midpoint display 90, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 162 {
		t.Fatalf("expected clamp=true falling midpoint reveal height 162, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingClampTrueTargetsDisplayedMinValue(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(730, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 20, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", -40, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected clamp=true low out-of-range bar damping to start active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != -40 || movement.TargetValue != 0 {
		t.Fatalf("expected clamp=true low out-of-range target to clamp to displayed min, got %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=true falling-to-min midpoint tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-10) > 0.001 {
		t.Fatalf("expected clamp=true falling-to-min midpoint display 10, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 18 {
		t.Fatalf("expected clamp=true falling-to-min reveal height 18, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingClampFalsePreservesRawTargetValue(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, false)
	start := time.Unix(740, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected unclamped bar damping to keep active raw movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != 200 || movement.TargetValue != 200 {
		t.Fatalf("expected clamp=false bar damping to keep raw target semantics, got %#v", movement)
	}

	scenes, changed, err := runtime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected unclamped bar damping midpoint tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if math.Abs(movement.DisplayValue-100) > 0.001 {
		t.Fatalf("expected clamp=false midpoint display 100 from raw target 200, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected clamp=false midpoint reveal height 180 at displayed value 100, got %d", got)
	}

	scenes, changed, err = runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected unclamped bar damping settle tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 200 || movement.TargetValue != 200 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected clamp=false bar damping to settle at raw target 200, got %#v", movement)
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected clamp=false settled reveal height to stay capped by bar geometry, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingClampTrueLeavesStoredSourceValueUnchanged(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(750, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=true source update to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 200 {
		t.Fatalf("expected stored source value 200 to remain unchanged, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(60 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected clamp=true active tick to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 200 {
		t.Fatalf("expected stored source value 200 to remain unchanged during clamped movement, got %v", got)
	}
}

func TestRuntimeBarGaugeDampingClampTrueSameDisplayedTargetDoesNotRestart(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(760, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	firstRetargetAt := start.Add(10 * time.Millisecond)
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), firstRetargetAt))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	secondRawAt := start.Add(60 * time.Millisecond)
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 150, "c"), secondRawAt))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected same-display clamp=true update to keep active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != 150 {
		t.Fatalf("expected latest raw target 150, got %#v", movement)
	}
	if movement.TargetValue != 100 || movement.PreviousDisplayValue != 0 {
		t.Fatalf("expected same visible target to preserve original display retarget state, got %#v", movement)
	}
	if !movement.StartedAt.Equal(firstRetargetAt) || movement.Duration != 100*time.Millisecond || movement.Phase != movementPhaseMoving {
		t.Fatalf("expected same visible target update not to restart active movement, got %#v", movement)
	}
	if math.Abs(movement.DisplayValue-50) > 0.001 {
		t.Fatalf("expected continued midpoint display 50 without restart, got %#v", movement)
	}
}

func TestRuntimeBarGaugeDampingClampTrueRepeatedHighUpdatesStillSettleOnSchedule(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(770, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 150, "c"), start.Add(60*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 180, "c"), start.Add(90*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	scenes, changed, err := runtime.Tick(start.Add(110 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected scheduled settle tick to redraw")
	}
	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.DisplayValue != 100 || movement.TargetValue != 100 || movement.RawTargetValue != 180 || movement.Phase != movementPhaseStatic {
		t.Fatalf("expected repeated same-display clamp=true updates to settle at displayed max on schedule, got %#v", movement)
	}
	if runtime.HasActiveMovement() {
		t.Fatalf("expected repeated same-display clamp=true updates not to extend active movement")
	}
	if got := firstPartKind(requireWidget(t, scenes[0], "coolant"), PartKindBar).Window.Height; got != 180 {
		t.Fatalf("expected settled reveal height 180 at displayed max, got %d", got)
	}
}

func TestRuntimeBarGaugeDampingClampFalseRawRetargetRestartsMovement(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, false)
	start := time.Unix(780, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}

	retargetAt := start.Add(60 * time.Millisecond)
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 150, "c"), retargetAt))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed || !runtime.HasActiveMovement() {
		t.Fatalf("expected clamp=false raw retarget to keep active movement")
	}

	movement := runtime.movements[movementKey("primary", "coolant")]
	if movement.RawTargetValue != 150 || movement.TargetValue != 150 {
		t.Fatalf("expected clamp=false retarget to update raw and display targets, got %#v", movement)
	}
	if math.Abs(movement.PreviousDisplayValue-100) > 0.001 || math.Abs(movement.DisplayValue-100) > 0.001 {
		t.Fatalf("expected clamp=false retarget to restart from current display 100, got %#v", movement)
	}
	if !movement.StartedAt.Equal(retargetAt) || movement.Duration != 100*time.Millisecond || movement.Phase != movementPhaseMoving {
		t.Fatalf("expected clamp=false retarget to restart timing, got %#v", movement)
	}
}

func TestRuntimeBarGaugeDampingClampTrueStoresLatestRawValue(t *testing.T) {
	runtime := testBarMovementRuntimeWithDampingAndValueMap(t, "  damping:\n    rise_ms: 100\n    fall_ms: 100\n", 0, 100, true)
	start := time.Unix(790, 0)

	_, _, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 0, "c"), start))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, _, err = runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 200, "c"), start.Add(10*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	_, changed, err := runtime.ApplyEvent(sensorEventAt("coolant_temperature", okState("coolant_temperature", 150, "c"), start.Add(60*time.Millisecond)))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected latest raw source update to redraw current movement state")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 150 {
		t.Fatalf("expected stored source value 150 after same-display retarget, got %v", got)
	}

	_, changed, err = runtime.Tick(start.Add(80 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected active movement tick to redraw")
	}
	if got := runtime.states["coolant_temperature"].Value; got != 150 {
		t.Fatalf("expected stored source value 150 to remain latest raw value during continued movement, got %v", got)
	}
}

func TestRuntimeLoadsSegmentedGaugeWidgetPackageAndPersistsHysteresis(t *testing.T) {
	packageDir := makeDashboardSegmentedGaugePackage(t)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{120, 80}, Scale: 1.5}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}

	state := okState("rpm", 3500, "rpm")
	state.Min = 0
	state.Max = 7000
	_, changed, err := runtime.ApplyEvent(sensorEvent("rpm", state))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected first segmented event to change rendered output")
	}
	scenes, err := runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	widget := requireWidget(t, scenes[0], "rpm")
	if widget.Type != v3config.WidgetTypeGauge || widget.SensorID != "rpm" || widget.GaugeID != "dashboard_rpm_segmented" {
		t.Fatalf("segmented widget identity = %#v", widget)
	}
	if got := gaugePartSequence(widget); got != "layer:panel,layer:segments,layer:glass" {
		t.Fatalf("segmented part sequence = %q", got)
	}
	if selected := firstGaugeLayerPart(widget, "segments"); selected.AssetPath == "" || !strings.HasSuffix(selected.AssetPath, "rpm_050.png") {
		t.Fatalf("first selected segment = %#v", selected)
	}

	state = okState("rpm", 3080, "rpm")
	state.Min = 0
	state.Max = 7000
	_, changed, err = runtime.ApplyEvent(sensorEvent("rpm", state))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if changed {
		t.Fatalf("expected hysteresis to keep the 050 segment active")
	}
	scenes, err = runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	if selected := firstGaugeLayerPart(requireWidget(t, scenes[0], "rpm"), "segments"); selected.AssetPath == "" || !strings.HasSuffix(selected.AssetPath, "rpm_050.png") {
		t.Fatalf("held segment = %#v", selected)
	}

	state = okState("rpm", 3000, "rpm")
	state.Min = 0
	state.Max = 7000
	_, changed, err = runtime.ApplyEvent(sensorEvent("rpm", state))
	if err != nil {
		t.Fatalf("ApplyEvent failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected 3000 to fall back to the 025 segment")
	}
	scenes, err = runtime.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	if selected := firstGaugeLayerPart(requireWidget(t, scenes[0], "rpm"), "segments"); selected.AssetPath == "" || !strings.HasSuffix(selected.AssetPath, "rpm_025.png") {
		t.Fatalf("dropped segment = %#v", selected)
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
	return makeDashboardRadialGaugePackageWithPolicy(t, "")
}

func makeDashboardRadialGaugePackageWithPolicy(t *testing.T, policy string) string {
	return makeDashboardRadialGaugePackageWithRealism(t, policy, false, nil, nil, false, false)
}

func makeDashboardRadialGaugePackageWithRealism(t *testing.T, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool, hysteresis bool) string {
	return makeDashboardRadialGaugePackageWithPointerMarkersAndRealism(t, "", policy, damping, stiction, overshoot, pegBounce, nil, nil, nil, hysteresis)
}

func makeDashboardRadialGaugePackageWithPointerMarkersAndRealism(t *testing.T, pointerMarkersYAML string, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool, shadowOffset []int, shadowAlpha *float64, calibrationOffset *float64, hysteresis bool) string {
	return makeDashboardRadialGaugePackageWithExtendedRealism(t, pointerMarkersYAML, policy, damping, stiction, overshoot, pegBounce, shadowOffset, shadowAlpha, calibrationOffset, hysteresis)
}

func makeDashboardRadialGaugePackageWithNeedleShadow(t *testing.T, offset []int, alpha *float64) string {
	return makeDashboardRadialGaugePackageWithExtendedRealism(t, "", "", false, nil, nil, false, offset, alpha, nil, false)
}

func makeDashboardRadialGaugePackageWithCalibrationOffset(t *testing.T, calibrationOffset *float64) string {
	return makeDashboardRadialGaugePackageWithExtendedRealism(t, "", "", false, nil, nil, false, nil, nil, calibrationOffset, false)
}

func makeDashboardRadialGaugePackageWithHysteresisAndClamp(t *testing.T, hysteresis bool, clamp bool) string {
	t.Helper()
	packageDir := makeDashboardRadialGaugePackageWithRealism(t, "", false, nil, nil, false, hysteresis)
	if clamp {
		return packageDir
	}
	gaugeYAML := dashboardRadialGaugeYAML("", "", false, nil, nil, false, nil, nil, nil, hysteresis)
	gaugeYAML = strings.Replace(gaugeYAML, "  clamp: true\n", "  clamp: false\n", 1)
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(gaugeYAML), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func makeDashboardRadialGaugePackageWithExtendedRealism(t *testing.T, pointerMarkersYAML string, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool, shadowOffset []int, shadowAlpha *float64, calibrationOffset *float64, hysteresis bool) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/radial/simple_rpm/background.png",
		"assets/gauges/radial/simple_rpm/face.png",
		"assets/gauges/radial/simple_rpm/ticks.png",
		"assets/gauges/radial/simple_rpm/needle.png",
		"assets/gauges/radial/simple_rpm/needle_min.png",
		"assets/gauges/radial/simple_rpm/needle_max.png",
		"assets/gauges/radial/simple_rpm/needle_average.png",
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
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardRadialGaugeYAML(pointerMarkersYAML, policy, damping, stiction, overshoot, pegBounce, shadowOffset, shadowAlpha, calibrationOffset, hysteresis)), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func makeDashboardOdometerGaugePackage(t *testing.T) string {
	return makeDashboardOdometerGaugePackageWithRealism(t, v3gauges.MovementBell, false, false, false)
}

func makeDashboardOdometerGaugePackageWithConfig(t *testing.T, movement string) string {
	return makeDashboardOdometerGaugePackageWithRealism(t, movement, false, false, false)
}

func makeDashboardOdometerGaugePackageWithRealism(t *testing.T, movement string, wraparound bool, carryDrag bool, snapSettle bool) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/odometer/trip/panel.png",
		"assets/gauges/odometer/trip/glass.png",
		"assets/gauges/odometer/trip/digits.png",
		"assets/gauges/odometer/trip/red_digits.png",
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
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "trip")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardOdometerGaugeYAML(movement, wraparound, carryDrag, snapSettle)), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func makeDashboardIndicatorGaugePackage(t *testing.T) string {
	return makeDashboardIndicatorGaugePackageWithThermalFade(t, nil, nil)
}

func makeDashboardIndicatorGaugePackageWithThermalFade(t *testing.T, riseMS *int, fallMS *int) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/indicator/check_engine/bezel.png",
		"assets/gauges/indicator/check_engine/face.png",
		"assets/gauges/indicator/check_engine/off.png",
		"assets/gauges/indicator/check_engine/on.png",
		"assets/gauges/indicator/check_engine/glass.png",
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
	packageDir := filepath.Join(root, "assets", "gauges", "indicator", "check_engine")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardIndicatorGaugeYAML(riseMS, fallMS)), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func makeDashboardBarGaugePackage(t *testing.T) string {
	return makeDashboardBarGaugePackageWithRealism(t, "")
}

func makeDashboardBarGaugePackageWithRealism(t *testing.T, realismYAML string) string {
	return makeDashboardBarGaugePackageWithRealismAndValueMap(t, realismYAML, 40, 120, true)
}

func makeDashboardBarGaugePackageWithRealismAndValueMap(t *testing.T, realismYAML string, min float64, max float64, clamp bool) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/bar/coolant/panel.png",
		"assets/gauges/bar/coolant/level.png",
		"assets/gauges/bar/coolant/marker_min.png",
		"assets/gauges/bar/coolant/marker_max.png",
		"assets/gauges/bar/coolant/marker_average.png",
		"assets/gauges/bar/coolant/glass.png",
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
	packageDir := filepath.Join(root, "assets", "gauges", "bar", "coolant")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardBarGaugeYAMLWithValueMap(realismYAML, min, max, clamp)), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func makeDashboardBarGaugePackageWithPointerMarkers(t *testing.T, pointerMarkersYAML string) string {
	t.Helper()
	realismLines := []string{"  pointer_markers:"}
	for _, line := range strings.Split(strings.TrimSuffix(pointerMarkersYAML, "\n"), "\n") {
		realismLines = append(realismLines, "  "+line)
	}
	return makeDashboardBarGaugePackageWithRealism(t, strings.Join(realismLines, "\n")+"\n")
}

func makeDashboardSegmentedGaugePackage(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/segmented/rpm/panel.png",
		"assets/gauges/segmented/rpm/glass.png",
		"assets/gauges/segmented/rpm/levels/rpm_025.png",
		"assets/gauges/segmented/rpm/levels/rpm_050.png",
		"assets/gauges/segmented/rpm/levels/rpm_100.png",
		"assets/gauges/segmented/rpm/levels/rpm_150.png",
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
	packageDir := filepath.Join(root, "assets", "gauges", "segmented", "rpm")
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardSegmentedGaugeYAML()), 0o600); err != nil {
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

func dashboardRadialGaugeYAML(pointerMarkersYAML string, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool, shadowOffset []int, shadowAlpha *float64, calibrationOffset *float64, hysteresis bool) string {
	realismLines := []string{}
	if strings.TrimSpace(pointerMarkersYAML) != "" {
		realismLines = append(realismLines, "  pointer_markers:")
		for _, line := range strings.Split(strings.TrimSuffix(pointerMarkersYAML, "\n"), "\n") {
			realismLines = append(realismLines, "  "+line)
		}
	}
	if hysteresis {
		realismLines = append(realismLines, "  hysteresis: true")
	}
	if stiction != nil {
		realismLines = append(realismLines, fmt.Sprintf("  stiction: %.0f", *stiction))
	}
	if damping {
		realismLines = append(realismLines, "  damping: true")
	}
	if pegBounce {
		realismLines = append(realismLines, "  peg_bounce: true")
	}
	if overshoot != nil {
		if overshoot.Ratio == nil &&
			overshoot.MinChangeRatio == nil &&
			overshoot.MaxSpanRatio == nil &&
			overshoot.SettleMode == "" &&
			overshoot.SettleCycles == nil &&
			overshoot.SettleDamping == nil &&
			!overshoot.AllowExtremes {
			realismLines = append(realismLines, "  overshoot: {}")
		} else {
			realismLines = append(realismLines, "  overshoot:")
			if overshoot.Ratio != nil {
				realismLines = append(realismLines, fmt.Sprintf("    ratio: %.2f", *overshoot.Ratio))
			}
			if overshoot.MinChangeRatio != nil {
				realismLines = append(realismLines, fmt.Sprintf("    min_change_ratio: %.2f", *overshoot.MinChangeRatio))
			}
			if overshoot.MaxSpanRatio != nil {
				realismLines = append(realismLines, fmt.Sprintf("    max_span_ratio: %.2f", *overshoot.MaxSpanRatio))
			}
			if overshoot.SettleMode != "" {
				realismLines = append(realismLines, "    settle_mode: "+overshoot.SettleMode)
			}
			if overshoot.SettleCycles != nil {
				realismLines = append(realismLines, fmt.Sprintf("    settle_cycles: %.2f", *overshoot.SettleCycles))
			}
			if overshoot.SettleDamping != nil {
				realismLines = append(realismLines, fmt.Sprintf("    settle_damping: %.2f", *overshoot.SettleDamping))
			}
			if overshoot.AllowExtremes {
				realismLines = append(realismLines, "    allow_extremes: true")
			}
		}
	}
	if len(shadowOffset) == 2 {
		realismLines = append(realismLines, "  needle_shadow:")
		realismLines = append(realismLines, fmt.Sprintf("    offset: [%d, %d]", shadowOffset[0], shadowOffset[1]))
		if shadowAlpha != nil {
			realismLines = append(realismLines, fmt.Sprintf("    alpha: %.3f", *shadowAlpha))
		}
	}
	if calibrationOffset != nil {
		realismLines = append(realismLines, fmt.Sprintf("  calibration_offset: %.3f", *calibrationOffset))
	}
	if policy != "" {
		realismLines = append(realismLines, "  movement_policy: "+policy)
	}
	realismBlock := ""
	if len(realismLines) > 0 {
		realismBlock = "realism:\n" + strings.Join(realismLines, "\n") + "\n"
	}
	return `id: dashboard_radial_rpm
type: radial
sensor: rpm
` + realismBlock + `size:
  width: 512
  height: 512
layers:
  background: background.png
  face: face.png
  ticks: ticks.png
  needle: needle.png
  needle_min: needle_min.png
  needle_max: needle_max.png
  needle_average: needle_average.png
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

func dashboardOdometerGaugeYAML(movement string, wraparound bool, carryDrag bool, snapSettle bool) string {
	if strings.TrimSpace(movement) == "" {
		movement = v3gauges.MovementInstant
	}
	realismLines := []string{}
	if wraparound {
		realismLines = append(realismLines, "  wraparound: true")
	}
	if carryDrag {
		realismLines = append(realismLines, "  carry_drag: true")
	}
	if snapSettle {
		realismLines = append(realismLines, "  snap_settle: true")
	}
	realismBlock := ""
	if len(realismLines) > 0 {
		realismBlock = "realism:\n" + strings.Join(realismLines, "\n") + "\n"
	}
	return `id: dashboard_trip_odometer
type: odometer
sensor: trip_distance
` + realismBlock + `size:
  width: 150
  height: 60
layers:
  panel: panel.png
  glass: glass.png
odometer:
  movement: ` + movement + `
  wheels:
    - strip: digits.png
      position: [10, 12]
      window: { width: 12, height: 20 }
    - strip: digits.png
      position: [24, 12]
      window: { width: 12, height: 20 }
    - strip: red_digits.png
      position: [42, 12]
      window: { width: 12, height: 20 }
      role: sub_unit
`
}

func dashboardIndicatorGaugeYAML(riseMS *int, fallMS *int) string {
	realismBlock := ""
	if riseMS != nil || fallMS != nil {
		rise := 120
		fall := 240
		if riseMS != nil {
			rise = *riseMS
		}
		if fallMS != nil {
			fall = *fallMS
		}
		realismBlock = fmt.Sprintf("realism:\n  thermal_fade:\n    rise_ms: %d\n    fall_ms: %d\n", rise, fall)
	}
	return `id: dashboard_check_engine_indicator
type: indicator
sensor: check_engine
` + realismBlock + `size:
  width: 48
  height: 48
layers:
  bezel: bezel.png
  face: face.png
  off: off.png
  on: on.png
  glass: glass.png
`
}

func dashboardBarGaugeYAML(realismYAML string) string {
	return dashboardBarGaugeYAMLWithValueMap(realismYAML, 40, 120, true)
}

func dashboardBarGaugeYAMLWithValueMap(realismYAML string, min float64, max float64, clamp bool) string {
	realismBlock := ""
	if strings.TrimSpace(realismYAML) != "" {
		realismBlock = "realism:\n" + realismYAML
	}
	return fmt.Sprintf(`id: dashboard_coolant_bar
type: bar
sensor: coolant_temperature
%ssize:
  width: 120
  height: 220
layers:
  panel: panel.png
  level: level.png
  marker_min: marker_min.png
  marker_max: marker_max.png
  marker_average: marker_average.png
  glass: glass.png
value_map:
  min: %g
  max: %g
  clamp: %t
bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [40, 20, 24, 180]
`, realismBlock, min, max, clamp)
}

func dashboardSegmentedGaugeYAML() string {
	return `id: dashboard_rpm_segmented
type: segmented
sensor: rpm
size:
  width: 120
  height: 120
layers:
  panel: panel.png
  segments: levels/rpm_{percent:03}.png
  glass: glass.png
`
}

func testRadialMovementRuntime(t *testing.T) *Runtime {
	return testRadialMovementRuntimeWithPolicy(t, "")
}

func testRadialMovementRuntimeWithPolicy(t *testing.T, policy string) *Runtime {
	return testRadialMovementRuntimeWithRealism(t, policy, false, nil, nil, false, false)
}

func testRadialMovementRuntimeWithPointerMarkersAndRealism(t *testing.T, pointerMarkersYAML string, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool, hysteresis bool) *Runtime {
	t.Helper()
	packageDir := makeDashboardRadialGaugePackageWithPointerMarkersAndRealism(t, pointerMarkersYAML, policy, damping, stiction, overshoot, pegBounce, nil, nil, nil, hysteresis)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	return runtime
}

func testRadialMovementRuntimeWithDamping(t *testing.T, damping bool) *Runtime {
	return testRadialMovementRuntimeWithRealism(t, "", damping, nil, nil, false, false)
}

func testRadialMovementRuntimeWithHysteresis(t *testing.T, hysteresis bool) *Runtime {
	return testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, false, hysteresis)
}

func testRadialMovementRuntimeWithHysteresisAndClamp(t *testing.T, hysteresis bool, clamp bool) *Runtime {
	t.Helper()
	packageDir := makeDashboardRadialGaugePackageWithHysteresisAndClamp(t, hysteresis, clamp)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	return runtime
}

func testRadialMovementRuntimeWithRealism(t *testing.T, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool, hysteresis bool) *Runtime {
	t.Helper()
	packageDir := makeDashboardRadialGaugePackageWithRealism(t, policy, damping, stiction, overshoot, pegBounce, hysteresis)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "rpm", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	return runtime
}

func testOdometerMovementRuntimeWithConfig(t *testing.T, movement string) *Runtime {
	return testOdometerMovementRuntimeWithRealism(t, movement, false, false, false)
}

func testOdometerMovementRuntimeWithRealism(t *testing.T, movement string, wraparound bool, carryDrag bool, snapSettle bool) *Runtime {
	t.Helper()
	packageDir := makeDashboardOdometerGaugePackageWithRealism(t, movement, wraparound, carryDrag, snapSettle)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "trip", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	return runtime
}

func testBarMovementRuntimeWithDamping(t *testing.T, realismYAML string) *Runtime {
	return testBarMovementRuntimeWithRealismAndValueMap(t, realismYAML, 40, 120, true)
}

func testBarMovementRuntimeWithDampingAndValueMap(t *testing.T, realismYAML string, min float64, max float64, clamp bool) *Runtime {
	return testBarMovementRuntimeWithRealismAndValueMap(t, realismYAML, min, max, clamp)
}

func testBarMovementRuntimeWithRealismAndValueMap(t *testing.T, realismYAML string, min float64, max float64, clamp bool) *Runtime {
	t.Helper()
	packageDir := makeDashboardBarGaugePackageWithRealismAndValueMap(t, realismYAML, min, max, clamp)
	plan := v3config.RuntimePlan{Dashboards: []v3config.ResolvedDashboard{{ID: "primary", Config: v3config.DashboardConfig{Display: "HDMI-1", Size: v3config.SizeConfig{Width: 1024, Height: 600}, Widgets: []v3config.WidgetConfig{{ID: "coolant", Type: v3config.WidgetTypeGauge, Gauge: packageDir, Position: []int{0, 0}, Scale: 1}}}}}}
	runtime, err := NewRuntime(plan, testAssetRegistry())
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	return runtime
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

func wheelStripWidgetParts(widget Widget) []Part {
	parts := []Part{}
	for _, part := range widget.Parts {
		if part.Kind == PartKindWheelStrip {
			parts = append(parts, part)
		}
	}
	return parts
}

func wheelPartSliceDigits(part Part) []int {
	digits := make([]int, len(part.WheelSlices))
	for index, slice := range part.WheelSlices {
		digits[index] = slice.Digit
	}
	return digits
}

func intSlicesEqual(left []int, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func firstGaugeLayerPart(widget Widget, layer string) Part {
	for _, part := range widget.Parts {
		if part.Kind == PartKindLayer && part.Layer == layer {
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
		case PartKindNeedleShadow:
			parts = append(parts, fmt.Sprintf("needle_shadow:%.0f", part.Angle))
		case PartKindNeedle:
			parts = append(parts, fmt.Sprintf("needle:%.0f", part.Angle))
		case PartKindNeedleMin:
			parts = append(parts, fmt.Sprintf("needle_min:%.0f", part.Angle))
		case PartKindNeedleMax:
			parts = append(parts, fmt.Sprintf("needle_max:%.0f", part.Angle))
		case PartKindNeedleAverage:
			parts = append(parts, fmt.Sprintf("needle_average:%.0f", part.Angle))
		case PartKindMarkerMin:
			parts = append(parts, fmt.Sprintf("marker_min:%v", part.Position))
		case PartKindMarkerMax:
			parts = append(parts, fmt.Sprintf("marker_max:%v", part.Position))
		case PartKindMarkerAverage:
			parts = append(parts, fmt.Sprintf("marker_average:%v", part.Position))
		case PartKindBar:
			parts = append(parts, "bar:"+part.Layer)
		case PartKindWheelStrip:
			parts = append(parts, fmt.Sprintf("wheel_strip:%d", part.Slot))
		default:
			parts = append(parts, part.Kind)
		}
	}
	return strings.Join(parts, ",")
}

var _ = v3assets.IndicatorStateOff

func almostEqual(left float64, right float64) bool {
	return math.Abs(left-right) < 0.001
}
