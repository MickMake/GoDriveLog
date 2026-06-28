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
	packageDir := makeDashboardRadialGaugePackageWithPolicy(t, v3gauges.MovementPolicyLinear)
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
	runtime := testRadialMovementRuntimeWithPolicy(t, v3gauges.MovementPolicyLinear)
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
	runtime := testRadialMovementRuntimeWithPolicy(t, v3gauges.MovementPolicyLinear)
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
	runtime := testRadialMovementRuntimeWithPolicy(t, v3gauges.MovementPolicyLinear)
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
	runtime := testRadialMovementRuntimeWithPolicy(t, v3gauges.MovementPolicyLinear)
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

func TestRuntimeGaugeMovementDefaultsToImmediatePolicy(t *testing.T) {
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

func TestRuntimeGaugeMovementEaseOutPolicyAdvancesFurtherThanLinear(t *testing.T) {
	linearRuntime := testRadialMovementRuntimeWithPolicy(t, v3gauges.MovementPolicyLinear)
	easeOutRuntime := testRadialMovementRuntimeWithPolicy(t, v3gauges.MovementPolicyEaseOut)

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
	if !almostEqual(wheels[0].StripOffset, 25.8) || !almostEqual(wheels[1].StripOffset, 58.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
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
	if !almostEqual(wheels[0].StripOffset, 24.0) || !almostEqual(wheels[1].StripOffset, 40.0) || !almostEqual(wheels[2].StripOffset, 0.0) {
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
	if !(wheels[0].StripOffset >= 24.9 && wheels[0].StripOffset < 25.8) {
		t.Fatalf("expected tens wheel to advance between start and final offsets, got %.2f", wheels[0].StripOffset)
	}
	if !(wheels[1].StripOffset >= 49.0 && wheels[1].StripOffset < 58.0) {
		t.Fatalf("expected ones wheel to advance between start and final offsets, got %.2f", wheels[1].StripOffset)
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
	if !almostEqual(wheels[0].StripOffset, 25.8) || !almostEqual(wheels[1].StripOffset, 58.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
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
	if !(easeOutWheels[1].StripOffset > linearWheels[1].StripOffset) {
		t.Fatalf("expected ease_out to advance further than linear: linear=%.2f ease_out=%.2f", linearWheels[1].StripOffset, easeOutWheels[1].StripOffset)
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
	if !(bellWheels[1].StripOffset < linearWheels[1].StripOffset) {
		t.Fatalf("expected bell to start slower than linear: linear=%.2f bell=%.2f", linearWheels[1].StripOffset, bellWheels[1].StripOffset)
	}

	scenes, changed, err := bellRuntime.Tick(start.Add(210 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected final bell tick to redraw")
	}
	wheels := wheelStripWidgetParts(requireWidget(t, scenes[0], "trip"))
	if !almostEqual(wheels[0].StripOffset, 25.8) || !almostEqual(wheels[1].StripOffset, 58.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
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
	if !almostEqual(finalWheels[0].StripOffset, 40.0) || !almostEqual(finalWheels[1].StripOffset, 400.0) || !almostEqual(finalWheels[2].StripOffset, 4000.0) {
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
	if !(baseWheels[1].StripOffset < 400.0) {
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
	if !almostEqual(finalWheels[0].StripOffset, 40.4) || !almostEqual(finalWheels[1].StripOffset, 404.0) || !almostEqual(finalWheels[2].StripOffset, 4040.0) {
		t.Fatalf("expected carry-drag straddling update to settle exactly on target, got %.2f/%.2f/%.2f", finalWheels[0].StripOffset, finalWheels[1].StripOffset, finalWheels[2].StripOffset)
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
	if !(settleWheels[0].StripOffset > baseWheels[0].StripOffset) {
		t.Fatalf("expected snap_settle to nudge tens wheel beyond target: base=%.2f settle=%.2f", baseWheels[0].StripOffset, settleWheels[0].StripOffset)
	}
	if !(settleWheels[1].StripOffset > baseWheels[1].StripOffset) {
		t.Fatalf("expected snap_settle to nudge ones wheel beyond target: base=%.2f settle=%.2f", baseWheels[1].StripOffset, settleWheels[1].StripOffset)
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
	if !almostEqual(finalWheels[0].StripOffset, 25.8) || !almostEqual(finalWheels[1].StripOffset, 58.0) || !almostEqual(finalWheels[2].StripOffset, 180.0) {
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
			if !almostEqual(wheels[0].StripOffset, 25.8) || !almostEqual(wheels[1].StripOffset, 58.0) || !almostEqual(wheels[2].StripOffset, 180.0) {
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
	if got := gaugePartSequence(widget); got != "layer:bezel,layer:on,layer:glass" {
		t.Fatalf("indicator part sequence = %q", got)
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
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardRadialGaugeYAML(policy)), 0o600); err != nil {
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
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/indicator/check_engine/bezel.png",
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
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardIndicatorGaugeYAML()), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
}

func makeDashboardBarGaugePackage(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/bar/coolant/panel.png",
		"assets/gauges/bar/coolant/level.png",
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
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardBarGaugeYAML()), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return packageDir
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

func dashboardRadialGaugeYAML(policy string) string {
	realismBlock := ""
	if policy != "" {
		realismBlock = "realism:\n  movement_policy: " + policy + "\n"
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

func dashboardIndicatorGaugeYAML() string {
	return `id: dashboard_check_engine_indicator
type: indicator
sensor: check_engine
size:
  width: 48
  height: 48
layers:
  bezel: bezel.png
  on: on.png
  glass: glass.png
`
}

func dashboardBarGaugeYAML() string {
	return `id: dashboard_coolant_bar
type: bar
sensor: coolant_temperature
size:
  width: 120
  height: 220
layers:
  panel: panel.png
  level: level.png
  glass: glass.png
value_map:
  min: 40
  max: 120
  clamp: true
bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [40, 20, 24, 180]
`
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
	t.Helper()
	packageDir := makeDashboardRadialGaugePackageWithPolicy(t, policy)
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
		case PartKindNeedle:
			parts = append(parts, fmt.Sprintf("needle:%.0f", part.Angle))
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
