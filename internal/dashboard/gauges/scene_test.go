package gauges

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestNumericSceneUsesPackageOwnedFormatPositionsAndStaticLayers(t *testing.T) {
	pkg := loadNumericScenePackage(t, 4, "%04.0f")

	scene, err := NumericScene(pkg, Placement{Position: []int{780, 40}, Scale: 1.5}, okGaugeState("rpm", 12))
	if err != nil {
		t.Fatalf("NumericScene returned error: %v", err)
	}

	if scene.PackageID != "test_4_digit_rpm" || scene.SensorID != "rpm" || scene.Type != TypeNumeric {
		t.Fatalf("scene identity = %#v", scene)
	}
	if scene.Text != "0012" || scene.Status != sensors.StatusOK {
		t.Fatalf("scene text/status = %q/%q, want 0012/ok", scene.Text, scene.Status)
	}
	if scene.Position[0] != 780 || scene.Position[1] != 40 || scene.Scale != 1.5 {
		t.Fatalf("scene placement = %#v / %v", scene.Position, scene.Scale)
	}
	if scene.Size.Width != 398 || scene.Size.Height != 150 {
		t.Fatalf("scene size = %#v", scene.Size)
	}
	if got := layerNames(scene); got != "panel,glass" {
		t.Fatalf("static layer names = %q, want panel,glass", got)
	}
	if got := sceneCharacters(scene); got != "0012" {
		t.Fatalf("characters = %q, want 0012", got)
	}
	char := firstCharacterPart(scene, "1")
	if char.Slot != 2 || char.Position[0] != 22 || char.Position[1] != 12 {
		t.Fatalf("expected character 1 at slot 2 position [22,12], got %#v", char)
	}
}

func TestNumericSceneEmitsPanelUnderDigitsAndGlassOverDigits(t *testing.T) {
	pkg := loadNumericScenePackage(t, 4, "%04.0f")
	scene, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 12))
	if err != nil {
		t.Fatalf("NumericScene returned error: %v", err)
	}
	sequence := partLayerSequence(scene)
	if !strings.HasPrefix(sequence, "layer:panel,") {
		t.Fatalf("expected panel underlay first, got %q", sequence)
	}
	if !strings.HasSuffix(sequence, ",layer:glass") {
		t.Fatalf("expected glass overlay last, got %q", sequence)
	}
	if strings.Index(sequence, "layer:glass") < strings.Index(sequence, "character:1") {
		t.Fatalf("expected glass after live digits, got %q", sequence)
	}
}

func TestNumericSceneSupportsTwoThreeFourAndFiveDigitShapes(t *testing.T) {
	for _, count := range []int{2, 3, 4, 5} {
		t.Run(fmt.Sprintf("%d_digits", count), func(t *testing.T) {
			format := fmt.Sprintf("%%0%d.0f", count)
			pkg := loadNumericScenePackage(t, count, format)
			scene, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 12))
			if err != nil {
				t.Fatalf("NumericScene returned error: %v", err)
			}
			if got := countSceneParts(scene, ScenePartKindCharacter); got != count {
				t.Fatalf("expected %d character parts, got %d from %#v", count, got, scene.Parts)
			}
			last := firstCharacterPart(scene, "2")
			if last.Slot != count-1 || last.Position[0] != (count-1)*10+2 {
				t.Fatalf("last digit position = %#v", last)
			}
		})
	}
}

func TestNumericSceneDoesNotRenderLiveDigitsForNonOKStates(t *testing.T) {
	pkg := loadNumericScenePackage(t, 4, "%04.0f")
	statuses := []string{
		sensors.StatusMissing,
		sensors.StatusUnsupported,
		sensors.StatusTimeout,
		sensors.StatusParseError,
		sensors.StatusError,
		sensors.StatusStale,
		sensors.StatusUnknown,
	}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			scene, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "rpm", Status: status, Error: "not live"})
			if err != nil {
				t.Fatalf("NumericScene returned error: %v", err)
			}
			if scene.Status != status || scene.Error != "not live" {
				t.Fatalf("scene status/error = %q/%q", scene.Status, scene.Error)
			}
			if scene.Text != "" {
				t.Fatalf("expected no live text for %q, got %q", status, scene.Text)
			}
			if got := countSceneParts(scene, ScenePartKindCharacter); got != 0 {
				t.Fatalf("expected no live characters for %q, got %d", status, got)
			}
			if got := countSceneParts(scene, ScenePartKindBackground); got != 0 {
				t.Fatalf("expected no digit backgrounds for %q, got %d", status, got)
			}
			if got := countSceneParts(scene, ScenePartKindLayer); got != 2 {
				t.Fatalf("expected static panel/glass layers for %q, got %d", status, got)
			}
			if got := layerNames(scene); got != "panel,glass" {
				t.Fatalf("expected non-ok layers to keep draw order, got %q", got)
			}
		})
	}
}

func TestNumericSceneSignatureChangesWithFormattedOutput(t *testing.T) {
	pkg := loadNumericScenePackage(t, 4, "%04.0f")
	first, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 42.1))
	if err != nil {
		t.Fatalf("NumericScene returned error: %v", err)
	}
	unchanged, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 42.2))
	if err != nil {
		t.Fatalf("NumericScene returned error: %v", err)
	}
	changed, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 43))
	if err != nil {
		t.Fatalf("NumericScene returned error: %v", err)
	}

	if first.Signature() != unchanged.Signature() {
		t.Fatalf("expected unchanged rounded output to keep same signature")
	}
	if first.Signature() == changed.Signature() {
		t.Fatalf("expected changed formatted output to change signature")
	}
}

func TestNumericSceneSignatureIncludesDigitPositionsForNonOKState(t *testing.T) {
	pkg := loadNumericScenePackage(t, 4, "%04.0f")
	first, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout})
	if err != nil {
		t.Fatalf("NumericScene returned error: %v", err)
	}
	pkg.Digits.Positions[2] = []int{999, 12}
	changed, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout})
	if err != nil {
		t.Fatalf("NumericScene returned error: %v", err)
	}
	if first.Signature() == changed.Signature() {
		t.Fatalf("expected non-ok signature to change when digit positions change")
	}
}

func TestNumericSceneRejectsMissingDecimalPointWhenFormatNeedsIt(t *testing.T) {
	pkg := loadNumericScenePackage(t, 4, "%.1f")
	pkg.DigitSet.DecimalPoint = ""

	_, err := NumericScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 12.3))
	if err == nil {
		t.Fatal("expected missing decimal point to fail")
	}
	assertErrorContains(t, err, "decimal_point")
}

func TestRadialSceneUsesPackageOwnedPivotsValueMapAndLayerOrder(t *testing.T) {
	pkg := loadRadialScenePackage(t)

	scene, err := RadialScene(pkg, Placement{Position: []int{120, 80}, Scale: 0.75}, okGaugeState("rpm", 3500))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}

	if scene.PackageID != "simple_radial_rpm" || scene.SensorID != "rpm" || scene.Type != TypeRadial {
		t.Fatalf("scene identity = %#v", scene)
	}
	if scene.Position[0] != 120 || scene.Position[1] != 80 || scene.Scale != 0.75 {
		t.Fatalf("scene placement = %#v / %v", scene.Position, scene.Scale)
	}
	if scene.Size.Width != 512 || scene.Size.Height != 512 {
		t.Fatalf("scene size = %#v", scene.Size)
	}
	if scene.FacePivot.X != 0.5 || scene.FacePivot.Y != 0.55 || scene.NeedlePivot.X != 0.5 || scene.NeedlePivot.Y != 0.9 {
		t.Fatalf("scene pivots = face %#v needle %#v", scene.FacePivot, scene.NeedlePivot)
	}
	if !almostEqual(scene.Angle, 0) {
		t.Fatalf("scene angle = %v, want 0", scene.Angle)
	}
	if got := partLayerSequence(scene); got != "layer:background,layer:face,layer:ticks,needle:0,layer:overlay" {
		t.Fatalf("part sequence = %q", got)
	}
	needle := firstPart(scene, ScenePartKindNeedle)
	if needle.Layer != "needle" || needle.AssetPath == "" || !almostEqual(needle.Angle, 0) {
		t.Fatalf("needle part = %#v", needle)
	}
	if needle.FacePivot != scene.FacePivot || needle.NeedlePivot != scene.NeedlePivot {
		t.Fatalf("needle pivots = face %#v needle %#v", needle.FacePivot, needle.NeedlePivot)
	}
}

func TestRadialSceneClampsAnglesAndChangesSignatureWithAngle(t *testing.T) {
	pkg := loadRadialScenePackage(t)

	minimum, err := RadialScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", -100))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}
	maximum, err := RadialScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 9999))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}

	if !almostEqual(minimum.Angle, -135) || !almostEqual(maximum.Angle, 135) {
		t.Fatalf("clamped angles = %v/%v, want -135/135", minimum.Angle, maximum.Angle)
	}
	if minimum.Signature() == maximum.Signature() {
		t.Fatal("expected signature to change when radial angle changes")
	}
}

func TestRadialSceneDoesNotRenderNeedleForNonOKStates(t *testing.T) {
	pkg := loadRadialScenePackage(t)

	scene, err := RadialScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout, Error: "not live"})
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}
	if scene.Status != sensors.StatusTimeout || scene.Error != "not live" {
		t.Fatalf("scene status/error = %q/%q", scene.Status, scene.Error)
	}
	if got := countSceneParts(scene, ScenePartKindNeedle); got != 0 {
		t.Fatalf("expected no needle for non-ok state, got %d", got)
	}
	if got := partLayerSequence(scene); got != "layer:background,layer:face,layer:ticks,layer:overlay" {
		t.Fatalf("non-ok radial sequence = %q", got)
	}
}

func TestRadialSceneRejectsMissingNeedleLayer(t *testing.T) {
	pkg := loadRadialScenePackage(t)
	delete(pkg.Layers, "needle")

	_, err := RadialScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 10))
	if err == nil {
		t.Fatal("expected missing needle layer to fail")
	}
	assertErrorContains(t, err, "needle")
}

func TestOdometerSceneEmitsSmoothWheelStripOffsets(t *testing.T) {
	pkg := loadOdometerScenePackage(t, "")

	scene, err := OdometerScene(pkg, Placement{Position: []int{50, 20}, Scale: 1.25}, okGaugeState("trip_distance", 12.3))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	if scene.PackageID != "test_trip_odometer" || scene.SensorID != "trip_distance" || scene.Type != TypeOdometer {
		t.Fatalf("scene identity = %#v", scene)
	}
	if scene.Movement != MovementSmooth || scene.Status != sensors.StatusOK {
		t.Fatalf("scene movement/status = %q/%q", scene.Movement, scene.Status)
	}
	if scene.Position[0] != 50 || scene.Position[1] != 20 || scene.Scale != 1.25 {
		t.Fatalf("scene placement = %#v / %v", scene.Position, scene.Scale)
	}
	if got := partLayerSequence(scene); got != "layer:panel,wheel_strip:0,wheel_strip:1,wheel_strip:2,layer:glass" {
		t.Fatalf("part sequence = %q", got)
	}
	wheels := wheelStripParts(scene)
	if len(wheels) != 3 {
		t.Fatalf("wheel parts = %d, want 3", len(wheels))
	}
	if !almostEqual(wheels[0].StripOffset, 24.6) || !almostEqual(wheels[1].StripOffset, 46) || !almostEqual(wheels[2].StripOffset, 60) {
		t.Fatalf("smooth offsets = %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}
	if wheels[2].Role != WheelRoleSubUnit || wheels[2].Source[0] != 2 || wheels[2].Source[1] != 64 {
		t.Fatalf("sub-unit wheel = %#v", wheels[2])
	}
}

func TestOdometerSceneClickMovementSnapsWheelStripOffsets(t *testing.T) {
	pkg := loadOdometerScenePackage(t, "click")

	scene, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 12.9))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	wheels := wheelStripParts(scene)
	if len(wheels) != 3 {
		t.Fatalf("wheel parts = %d, want 3", len(wheels))
	}
	if !almostEqual(wheels[0].StripOffset, 20) || !almostEqual(wheels[1].StripOffset, 40) || !almostEqual(wheels[2].StripOffset, 180) {
		t.Fatalf("click offsets = %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}
}

func TestOdometerSceneDoesNotRenderWheelStripsForNonOKStates(t *testing.T) {
	pkg := loadOdometerScenePackage(t, "")

	scene, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "trip_distance", Status: sensors.StatusTimeout, Error: "not live"})
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}
	if scene.Status != sensors.StatusTimeout || scene.Error != "not live" {
		t.Fatalf("scene status/error = %q/%q", scene.Status, scene.Error)
	}
	if got := countSceneParts(scene, ScenePartKindWheelStrip); got != 0 {
		t.Fatalf("expected no wheel strips for non-ok state, got %d", got)
	}
	if got := partLayerSequence(scene); got != "layer:panel,layer:glass" {
		t.Fatalf("non-ok odometer sequence = %q", got)
	}
}

func TestOdometerSceneSignatureChangesWithSmoothOffset(t *testing.T) {
	pkg := loadOdometerScenePackage(t, "")
	first, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 12.1))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}
	changed, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 12.2))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}
	if first.Signature() == changed.Signature() {
		t.Fatal("expected signature to change when smooth wheel offset changes")
	}
}

func TestIndicatorSceneSelectsOffAndOnLayers(t *testing.T) {
	pkg := loadIndicatorScenePackage(t)

	offScene, err := IndicatorScene(pkg, Placement{Position: []int{16, 24}, Scale: 1.25}, okGaugeState("check_engine", 0))
	if err != nil {
		t.Fatalf("IndicatorScene returned error: %v", err)
	}
	onScene, err := IndicatorScene(pkg, Placement{Position: []int{16, 24}, Scale: 1.25}, okGaugeState("check_engine", 1))
	if err != nil {
		t.Fatalf("IndicatorScene returned error: %v", err)
	}

	if offScene.PackageID != "test_check_engine_indicator" || offScene.SensorID != "check_engine" || offScene.Type != TypeIndicator {
		t.Fatalf("scene identity = %#v", offScene)
	}
	if offScene.Position[0] != 16 || offScene.Position[1] != 24 || offScene.Scale != 1.25 {
		t.Fatalf("scene placement = %#v / %v", offScene.Position, offScene.Scale)
	}
	if got := partLayerSequence(offScene); got != "layer:bezel,layer:face,layer:off,layer:glass" {
		t.Fatalf("off sequence = %q", got)
	}
	if got := partLayerSequence(onScene); got != "layer:bezel,layer:face,layer:on,layer:glass" {
		t.Fatalf("on sequence = %q", got)
	}
	if offScene.Signature() == onScene.Signature() {
		t.Fatal("expected indicator signature to change between off and on")
	}
}

func TestIndicatorSceneUsesBoolTypedValueWhenPresent(t *testing.T) {
	pkg := loadIndicatorScenePackage(t)
	value := true

	scene, err := IndicatorScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "check_engine", Status: sensors.StatusOK, TypedValue: sensors.Value{Kind: sensors.ValueKindBool, Bool: &value}})
	if err != nil {
		t.Fatalf("IndicatorScene returned error: %v", err)
	}
	if got := partLayerSequence(scene); got != "layer:bezel,layer:face,layer:on,layer:glass" {
		t.Fatalf("bool indicator sequence = %q", got)
	}
}

func TestIndicatorSceneRendersOffForNonOKState(t *testing.T) {
	pkg := loadIndicatorScenePackage(t)

	scene, err := IndicatorScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "check_engine", Value: 1, Status: sensors.StatusTimeout, Error: "not live"})
	if err != nil {
		t.Fatalf("IndicatorScene returned error: %v", err)
	}
	if scene.Status != sensors.StatusTimeout || scene.Error != "not live" {
		t.Fatalf("scene status/error = %q/%q", scene.Status, scene.Error)
	}
	if got := partLayerSequence(scene); got != "layer:bezel,layer:face,layer:off,layer:glass" {
		t.Fatalf("non-ok indicator sequence = %q", got)
	}
}

func TestIndicatorSceneWithOnlyOnLayerDrawsNoOffStateLayer(t *testing.T) {
	pkg := loadIndicatorScenePackageWithOnlyOnLayer(t)

	scene, err := IndicatorScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("check_engine", 0))
	if err != nil {
		t.Fatalf("IndicatorScene returned error: %v", err)
	}
	if got := partLayerSequence(scene); got != "layer:bezel,layer:glass" {
		t.Fatalf("off sequence without off layer = %q", got)
	}
}

func TestIndicatorSceneWithOnlyOnLayerDrawsNoStateLayerForNonOKState(t *testing.T) {
	pkg := loadIndicatorScenePackageWithOnlyOnLayer(t)

	scene, err := IndicatorScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "check_engine", Value: 1, Status: sensors.StatusTimeout, Error: "not live"})
	if err != nil {
		t.Fatalf("IndicatorScene returned error: %v", err)
	}
	if scene.Status != sensors.StatusTimeout || scene.Error != "not live" {
		t.Fatalf("scene status/error = %q/%q", scene.Status, scene.Error)
	}
	if got := partLayerSequence(scene); got != "layer:bezel,layer:glass" {
		t.Fatalf("non-ok sequence without off layer = %q", got)
	}
}

func TestBarSceneUsesPackageBoundsAndRevealHeight(t *testing.T) {
	pkg := loadBarScenePackage(t)

	scene, err := BarScene(pkg, Placement{Position: []int{50, 20}, Scale: 1.25}, okGaugeState("coolant_temperature", 80))
	if err != nil {
		t.Fatalf("BarScene returned error: %v", err)
	}

	if scene.PackageID != "test_coolant_bar" || scene.SensorID != "coolant_temperature" || scene.Type != TypeBar {
		t.Fatalf("scene identity = %#v", scene)
	}
	if scene.BarMode != "level" || scene.BarAxis != "vertical" || scene.BarOrigin != "bottom" {
		t.Fatalf("bar config = %#v", scene)
	}
	if len(scene.BarBounds) != 4 || scene.BarBounds[0] != 40 || scene.BarBounds[3] != 180 {
		t.Fatalf("bar bounds = %#v", scene.BarBounds)
	}
	if got := partLayerSequence(scene); got != "layer:panel,bar:level,layer:glass" {
		t.Fatalf("bar part sequence = %q", got)
	}
	bar := firstPart(scene, ScenePartKindBar)
	if bar.Layer != "level" || bar.AssetPath == "" {
		t.Fatalf("bar part = %#v", bar)
	}
	if len(bar.Position) != 2 || bar.Position[0] != 40 || bar.Position[1] != 110 {
		t.Fatalf("bar position = %#v, want [40 110]", bar.Position)
	}
	if len(bar.Source) != 2 || bar.Source[0] != 40 || bar.Source[1] != 110 {
		t.Fatalf("bar source = %#v, want [40 110]", bar.Source)
	}
	if bar.Window.Width != 24 || bar.Window.Height != 90 {
		t.Fatalf("bar window = %#v, want 24x90", bar.Window)
	}
}

func TestBarSceneClampBehaviourUsesValueMapAndDrawableGeometry(t *testing.T) {
	pkg := loadBarScenePackage(t)

	belowMin, err := BarScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("coolant_temperature", 20))
	if err != nil {
		t.Fatalf("BarScene returned error: %v", err)
	}
	if got := countSceneParts(belowMin, ScenePartKindBar); got != 0 {
		t.Fatalf("expected no bar part below min, got %d", got)
	}

	aboveMax, err := BarScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("coolant_temperature", 130))
	if err != nil {
		t.Fatalf("BarScene returned error: %v", err)
	}
	bar := firstPart(aboveMax, ScenePartKindBar)
	if len(bar.Source) != 2 || bar.Source[0] != 40 || bar.Source[1] != 20 {
		t.Fatalf("full bar source = %#v, want [40 20]", bar.Source)
	}
	if len(bar.Position) != 2 || bar.Position[0] != 40 || bar.Position[1] != 20 {
		t.Fatalf("full bar position = %#v, want [40 20]", bar.Position)
	}
	if bar.Window.Width != 24 || bar.Window.Height != 180 {
		t.Fatalf("full bar window = %#v, want 24x180", bar.Window)
	}
}

func TestBarSceneDoesNotRenderLevelForNonOKState(t *testing.T) {
	pkg := loadBarScenePackage(t)

	scene, err := BarScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "coolant_temperature", Status: sensors.StatusTimeout, Error: "not live"})
	if err != nil {
		t.Fatalf("BarScene returned error: %v", err)
	}
	if scene.Status != sensors.StatusTimeout || scene.Error != "not live" {
		t.Fatalf("scene status/error = %q/%q", scene.Status, scene.Error)
	}
	if got := countSceneParts(scene, ScenePartKindBar); got != 0 {
		t.Fatalf("expected no live bar part for non-ok state, got %d", got)
	}
	if got := partLayerSequence(scene); got != "layer:panel,layer:glass" {
		t.Fatalf("expected static bar layers in draw order, got %q", got)
	}
}

func TestBarSceneSignatureChangesWithRevealHeight(t *testing.T) {
	pkg := loadBarScenePackage(t)

	first, err := BarScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("coolant_temperature", 80))
	if err != nil {
		t.Fatalf("BarScene returned error: %v", err)
	}
	changed, err := BarScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("coolant_temperature", 81))
	if err != nil {
		t.Fatalf("BarScene returned error: %v", err)
	}
	if first.Signature() == changed.Signature() {
		t.Fatal("expected different reveal height to change signature")
	}
}

func loadNumericScenePackage(t *testing.T, count int, format string) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "7Seg", "amber", fmt.Sprintf("%d_digit_rpm", count))
	writeGaugeYAML(t, packageDir, numericGaugeYAML(count, format))
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func loadRadialScenePackage(t *testing.T) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "simple_rpm")
	writeGaugeYAML(t, packageDir, radialGaugeYAML())
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func loadOdometerScenePackage(t *testing.T, movement string) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "trip")
	writeGaugeYAML(t, packageDir, odometerGaugeYAML(movement))
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func loadIndicatorScenePackage(t *testing.T) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "indicator", "check_engine")
	writeGaugeYAML(t, packageDir, indicatorGaugeYAML())
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func loadIndicatorScenePackageWithOnlyOnLayer(t *testing.T) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "indicator", "check_engine")
	writeGaugeYAML(t, packageDir, indicatorOnOnlyGaugeYAML())
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func loadBarScenePackage(t *testing.T) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "bar", "coolant")
	writeGaugeYAML(t, packageDir, barGaugeYAML())
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func numericGaugeYAML(count int, format string) string {
	var positions strings.Builder
	for slot := 0; slot < count; slot++ {
		positions.WriteString(fmt.Sprintf("    - [%d, 12]\n", slot*10+2))
	}
	return fmt.Sprintf(`id: test_%d_digit_rpm
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

func radialGaugeYAML() string {
	return `id: simple_radial_rpm
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

func odometerGaugeYAML(movement string) string {
	movementLine := ""
	if movement != "" {
		movementLine = fmt.Sprintf("  movement: %s\n", movement)
	}
	return fmt.Sprintf(`id: test_trip_odometer
type: odometer
sensor: trip_distance
size:
  width: 150
  height: 60
layers:
  panel: panel.png
  glass: glass.png
odometer:
%s  wheels:
    - strip: digits.png
      position: [10, 12]
      window: { width: 12, height: 20 }
    - strip: digits.png
      position: [24, 12]
      window: { width: 12, height: 20 }
    - strip: red_digits.png
      position: [42, 12]
      window: { width: 12, height: 20 }
      offset: [2, 4]
      role: sub_unit
`, movementLine)
}

func indicatorGaugeYAML() string {
	return `id: test_check_engine_indicator
type: indicator
sensor: check_engine
size:
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

func indicatorOnOnlyGaugeYAML() string {
	return `id: test_check_engine_indicator
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

func barGaugeYAML() string {
	return `id: test_coolant_bar
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

func okGaugeState(id string, value float64) sensors.SensorState {
	return sensors.SensorState{ID: id, Value: value, Status: sensors.StatusOK}
}

func countSceneParts(scene Scene, kind string) int {
	count := 0
	for _, part := range scene.Parts {
		if part.Kind == kind {
			count++
		}
	}
	return count
}

func layerNames(scene Scene) string {
	names := []string{}
	for _, part := range scene.Parts {
		if part.Kind == ScenePartKindLayer {
			names = append(names, part.Layer)
		}
	}
	return strings.Join(names, ",")
}

func partLayerSequence(scene Scene) string {
	parts := make([]string, 0, len(scene.Parts))
	for _, part := range scene.Parts {
		switch part.Kind {
		case ScenePartKindLayer:
			parts = append(parts, "layer:"+part.Layer)
		case ScenePartKindCharacter:
			parts = append(parts, "character:"+part.Character)
		case ScenePartKindNeedle:
			parts = append(parts, fmt.Sprintf("needle:%.0f", part.Angle))
		case ScenePartKindBar:
			parts = append(parts, "bar:"+part.Layer)
		case ScenePartKindWheelStrip:
			parts = append(parts, fmt.Sprintf("wheel_strip:%d", part.Slot))
		default:
			parts = append(parts, part.Kind)
		}
	}
	return strings.Join(parts, ",")
}

func wheelStripParts(scene Scene) []ScenePart {
	parts := []ScenePart{}
	for _, part := range scene.Parts {
		if part.Kind == ScenePartKindWheelStrip {
			parts = append(parts, part)
		}
	}
	return parts
}

func sceneCharacters(scene Scene) string {
	var b strings.Builder
	for _, part := range scene.Parts {
		if part.Kind == ScenePartKindCharacter {
			b.WriteString(part.Character)
		}
	}
	return b.String()
}

func firstCharacterPart(scene Scene, character string) ScenePart {
	for _, part := range scene.Parts {
		if part.Kind == ScenePartKindCharacter && part.Character == character {
			return part
		}
	}
	return ScenePart{}
}

func firstPart(scene Scene, kind string) ScenePart {
	for _, part := range scene.Parts {
		if part.Kind == kind {
			return part
		}
	}
	return ScenePart{}
}
