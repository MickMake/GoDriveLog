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
	packageDir := makeDashboardRadialGaugePackageWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, &threshold, nil, false)

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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, &threshold, nil, false)

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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, overshoot, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, overshoot, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", true, nil, overshoot, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, overshoot, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, false)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, true)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, true)
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
	runtime := testRadialMovementRuntimeWithRealism(t, "", false, nil, nil, true)
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
		t.Fatalf("expected in-range peg bounce movement to animate linearly")
	}

	movement := runtime.movements[movementKey("primary", "rpm")]
	if !movement.PegBounceEnabled || movement.PegBounceReboundValue != 0 || movement.PegBounceStopValue != 0 {
		t.Fatalf("expected in-range movement to avoid scheduling peg bounce, got %#v", movement)
	}

	_, changed, err = runtime.Tick(start.Add(200 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected in-range peg bounce tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue <= 2000 || movement.DisplayValue >= 6000 {
		t.Fatalf("expected ordinary in-range movement without stop bounce, got %#v", movement)
	}

	_, changed, err = runtime.Tick(start.Add(320 * time.Millisecond))
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}
	if !changed {
		t.Fatalf("expected in-range movement completion tick to redraw")
	}
	movement = runtime.movements[movementKey("primary", "rpm")]
	if movement.DisplayValue != 6000 || movement.TargetValue != 6000 || movement.Phase != movementPhaseSettled {
		t.Fatalf("expected in-range movement to settle exactly on target, got %#v", movement)
	}
}

func TestRuntimeGaugeMovementEaseOutPolicyAdvancesFurtherThanLinear(t *testing.T) {
	linearRuntime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyLinear, true, nil, nil, false)
	easeOutRuntime := testRadialMovementRuntimeWithRealism(t, v3gauges.MovementPolicyEaseOut, true, nil, nil, false)

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
	return makeDashboardRadialGaugePackageWithRealism(t, policy, false, nil, nil, false)
}

func makeDashboardRadialGaugePackageWithRealism(t *testing.T, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool) string {
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
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(dashboardRadialGaugeYAML(policy, damping, stiction, overshoot, pegBounce)), 0o600); err != nil {
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

func dashboardRadialGaugeYAML(policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool) string {
	realismLines := []string{}
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
	return testRadialMovementRuntimeWithRealism(t, policy, false, nil, nil, false)
}

func testRadialMovementRuntimeWithDamping(t *testing.T, damping bool) *Runtime {
	return testRadialMovementRuntimeWithRealism(t, "", damping, nil, nil, false)
}

func testRadialMovementRuntimeWithRealism(t *testing.T, policy string, damping bool, stiction *float64, overshoot *v3gauges.OvershootConfig, pegBounce bool) *Runtime {
	t.Helper()
	packageDir := makeDashboardRadialGaugePackageWithRealism(t, policy, damping, stiction, overshoot, pegBounce)
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
