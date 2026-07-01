package gauges

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

func TestRadialSceneAddsNeedleShadowBeforeNeedleWhenConfigured(t *testing.T) {
	pkg := loadRadialScenePackageWithNeedleShadow(t, []int{3, 4}, nil)

	scene, err := RadialScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 3500))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}

	if got := partLayerSequence(scene); got != "layer:background,layer:face,layer:ticks,needle_shadow:0,needle:0,layer:overlay" {
		t.Fatalf("part sequence = %q", got)
	}
	shadow := firstPart(scene, ScenePartKindNeedleShadow)
	needle := firstPart(scene, ScenePartKindNeedle)
	if shadow.AssetPath != needle.AssetPath || !almostEqual(shadow.Angle, needle.Angle) {
		t.Fatalf("shadow/needle geometry = %#v / %#v", shadow, needle)
	}
	if !intSlicesEqual(shadow.Position, []int{3, 4}) {
		t.Fatalf("shadow offset = %#v, want [3 4]", shadow.Position)
	}
	if !almostEqual(shadow.Alpha, defaultNeedleShadowAlpha) {
		t.Fatalf("shadow alpha = %v, want %v", shadow.Alpha, defaultNeedleShadowAlpha)
	}
	if !almostEqual(scene.Angle, needle.Angle) {
		t.Fatalf("scene angle = %v, want needle angle %v", scene.Angle, needle.Angle)
	}
}

func TestRadialSceneCalibrationOffsetZeroPreservesAngle(t *testing.T) {
	zero := 0.0
	basePkg := loadRadialScenePackage(t)
	offsetPkg := loadRadialScenePackageWithCalibrationOffset(t, &zero)

	baseScene, err := RadialScene(basePkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 3500))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}
	offsetScene, err := RadialScene(offsetPkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 3500))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}

	if !almostEqual(baseScene.Angle, offsetScene.Angle) {
		t.Fatalf("zero calibration offset changed angle: base=%v offset=%v", baseScene.Angle, offsetScene.Angle)
	}
}

func TestRadialSceneCalibrationOffsetAppliesPositiveAndNegativeDegrees(t *testing.T) {
	positiveOffset := 12.0
	negativeOffset := -9.0
	positivePkg := loadRadialScenePackageWithCalibrationOffset(t, &positiveOffset)
	negativePkg := loadRadialScenePackageWithCalibrationOffset(t, &negativeOffset)

	positiveScene, err := RadialScene(positivePkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 3500))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}
	negativeScene, err := RadialScene(negativePkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 3500))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}

	if !almostEqual(positiveScene.Angle, 12) {
		t.Fatalf("positive calibration angle = %v, want 12", positiveScene.Angle)
	}
	if !almostEqual(negativeScene.Angle, -9) {
		t.Fatalf("negative calibration angle = %v, want -9", negativeScene.Angle)
	}
	if !almostEqual(firstPart(positiveScene, ScenePartKindNeedle).Angle, positiveScene.Angle) {
		t.Fatalf("positive needle angle = %v, want scene angle %v", firstPart(positiveScene, ScenePartKindNeedle).Angle, positiveScene.Angle)
	}
	if !almostEqual(firstPart(negativeScene, ScenePartKindNeedle).Angle, negativeScene.Angle) {
		t.Fatalf("negative needle angle = %v, want scene angle %v", firstPart(negativeScene, ScenePartKindNeedle).Angle, negativeScene.Angle)
	}
}

func TestRadialSceneCalibrationOffsetClampsToDialBounds(t *testing.T) {
	positiveOffset := 25.0
	negativeOffset := -25.0
	highPkg := loadRadialScenePackageWithCalibrationOffset(t, &positiveOffset)
	lowPkg := loadRadialScenePackageWithCalibrationOffset(t, &negativeOffset)

	highScene, err := RadialScene(highPkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 9999))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}
	lowScene, err := RadialScene(lowPkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", -100))
	if err != nil {
		t.Fatalf("RadialScene returned error: %v", err)
	}

	if !almostEqual(highScene.Angle, 135) {
		t.Fatalf("high clamped angle = %v, want 135", highScene.Angle)
	}
	if !almostEqual(lowScene.Angle, -135) {
		t.Fatalf("low clamped angle = %v, want -135", lowScene.Angle)
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

func TestOdometerSceneEmitsBellWheelStripOffsetsAtExactTarget(t *testing.T) {
	pkg := loadOdometerScenePackage(t, MovementBell, false, nil)

	scene, err := OdometerScene(pkg, Placement{Position: []int{50, 20}, Scale: 1.25}, okGaugeState("trip_distance", 12.3))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	if scene.PackageID != "test_trip_odometer" || scene.SensorID != "trip_distance" || scene.Type != TypeOdometer {
		t.Fatalf("scene identity = %#v", scene)
	}
	if scene.Movement != MovementBell || scene.Status != sensors.StatusOK {
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
	if !almostEqual(wheels[0].StripOffset, 20) || !almostEqual(wheels[1].StripOffset, 40) || !almostEqual(wheels[2].StripOffset, 60) {
		t.Fatalf("bell offsets = %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}
	if wheels[2].Role != WheelRoleSubUnit || wheels[2].Source[0] != 2 || wheels[2].Source[1] != 64 {
		t.Fatalf("sub-unit wheel = %#v", wheels[2])
	}
}

func TestOdometerSceneUsesDiscreteDigitSlotsForEveryWheel(t *testing.T) {
	pkg := loadFourWheelOdometerScenePackage(t)

	scene, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 123.4))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	wheels := wheelStripParts(scene)
	if len(wheels) != 4 {
		t.Fatalf("wheel parts = %d, want 4", len(wheels))
	}
	if !almostEqual(wheels[0].StripOffset, 20) || !almostEqual(wheels[1].StripOffset, 40) || !almostEqual(wheels[2].StripOffset, 60) || !almostEqual(wheels[3].StripOffset, 80) {
		t.Fatalf("expected exact discrete wheel slots for 123.4, got %.2f/%.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset, wheels[3].StripOffset)
	}
}

func TestOdometerWheelStripOffsetsOnlyChangeAffectedDigits(t *testing.T) {
	pkg := loadFourWheelOdometerScenePackage(t)

	first, err := OdometerWheelStripOffsets(pkg, 123.4)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	second, err := OdometerWheelStripOffsets(pkg, 123.5)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	if !almostEqual(second[0], first[0]) || !almostEqual(second[1], first[1]) || !almostEqual(second[2], first[2]) || !almostEqual(second[3], 100) {
		t.Fatalf("expected 123.4 -> 123.5 to change only sub-unit wheel, got %v -> %v", first, second)
	}

	third, err := OdometerWheelStripOffsets(pkg, 123.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	fourth, err := OdometerWheelStripOffsets(pkg, 124.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	if !almostEqual(fourth[0], third[0]) || !almostEqual(fourth[1], third[1]) || !almostEqual(fourth[2], 80) || !almostEqual(fourth[3], 0) {
		t.Fatalf("expected 123.9 -> 124.0 to change ones and sub-unit only, got %v -> %v", third, fourth)
	}

	fifth, err := OdometerWheelStripOffsets(pkg, 129.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	sixth, err := OdometerWheelStripOffsets(pkg, 130.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	if !almostEqual(sixth[0], fifth[0]) || !almostEqual(sixth[1], 60) || !almostEqual(sixth[2], 0) || !almostEqual(sixth[3], 0) {
		t.Fatalf("expected 129.9 -> 130.0 to change tens, ones, and sub-unit only, got %v -> %v", fifth, sixth)
	}
}

func TestOdometerSceneInstantMovementUsesExactTargetOffsets(t *testing.T) {
	pkg := loadOdometerScenePackage(t, MovementInstant, false, nil)

	scene, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 12.9))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	wheels := wheelStripParts(scene)
	if len(wheels) != 3 {
		t.Fatalf("wheel parts = %d, want 3", len(wheels))
	}
	if !almostEqual(wheels[0].StripOffset, 20) || !almostEqual(wheels[1].StripOffset, 40) || !almostEqual(wheels[2].StripOffset, 180) {
		t.Fatalf("instant offsets = %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}
}

func TestOdometerSceneDoesNotRenderWheelStripsForNonOKStates(t *testing.T) {
	pkg := loadOdometerScenePackage(t, "", false, nil)

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
	pkg := loadOdometerScenePackage(t, "", false, nil)
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

func TestOdometerSceneAlwaysUsesCircularWheelRendering(t *testing.T) {
	withConfig := loadOdometerScenePackage(t, "", true, nil)
	withoutConfig := loadOdometerScenePackage(t, "", false, nil)

	withScene, err := OdometerScene(withConfig, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 10.0))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}
	withoutScene, err := OdometerScene(withoutConfig, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 10.0))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	withWheels := wheelStripParts(withScene)
	withoutWheels := wheelStripParts(withoutScene)
	if len(withWheels) != len(withoutWheels) {
		t.Fatalf("wheel parts mismatch = %d vs %d", len(withWheels), len(withoutWheels))
	}
	for index := range withWheels {
		if !withWheels[index].Wraparound || !withoutWheels[index].Wraparound {
			t.Fatalf("expected all odometer wheels to render with circular wraparound, got with=%#v without=%#v", withWheels, withoutWheels)
		}
		if !almostEqual(withWheels[index].StripOffset, withoutWheels[index].StripOffset) || withWheels[index].Source[1] != withoutWheels[index].Source[1] {
			t.Fatalf("expected wraparound config to be a compatibility no-op, got with=%#v without=%#v", withWheels[index], withoutWheels[index])
		}
	}
}

func TestOdometerInterpolatedWheelOffsetsUseInfiniteForwardRoutingAcrossNineToZero(t *testing.T) {
	pkg := loadOdometerScenePackage(t, MovementLinear, false, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 0.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 1.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	mid, err := OdometerInterpolatedWheelOffsets(pkg, 0.9, 1.0, previousOffsets, targetOffsets, 0.5)
	if err != nil {
		t.Fatalf("OdometerInterpolatedWheelOffsets failed: %v", err)
	}
	if !(mid[2] > previousOffsets[2]) {
		t.Fatalf("expected forward 9 -> 0 to keep moving forward on the virtual strip, got prev=%.2f mid=%.2f target=%.2f", previousOffsets[2], mid[2], targetOffsets[2])
	}
	if !almostEqual(mid[2], 190) {
		t.Fatalf("expected forward 9 -> 0 midpoint to route through virtual slot 10, got %.2f", mid[2])
	}
}

func TestOdometerInterpolatedWheelOffsetsUseInfiniteBackwardRoutingAcrossZeroToNine(t *testing.T) {
	pkg := loadOdometerScenePackage(t, MovementLinear, false, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 1.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 0.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	mid, err := OdometerInterpolatedWheelOffsets(pkg, 1.0, 0.9, previousOffsets, targetOffsets, 0.5)
	if err != nil {
		t.Fatalf("OdometerInterpolatedWheelOffsets failed: %v", err)
	}
	if !(mid[2] < previousOffsets[2]) {
		t.Fatalf("expected backward 0 -> 9 to keep moving backward on the virtual strip, got prev=%.2f mid=%.2f target=%.2f", previousOffsets[2], mid[2], targetOffsets[2])
	}
	if !almostEqual(mid[2], -10) {
		t.Fatalf("expected backward 0 -> 9 midpoint to route through virtual slot -1, got %.2f", mid[2])
	}
}

func TestOdometerWheelSourceMapsVirtualSlotTenToDigitZero(t *testing.T) {
	wheel := OdometerWheel{Window: Size{Width: 12, Height: 20}}
	_, sourceY := odometerWheelSource(wheel, 200)
	if sourceY != 0 {
		t.Fatalf("expected virtual slot 10 to map to source digit 0, got sourceY=%d", sourceY)
	}
}

func TestOdometerWheelSourceMapsVirtualSlotElevenToDigitOne(t *testing.T) {
	wheel := OdometerWheel{Window: Size{Width: 12, Height: 20}}
	_, sourceY := odometerWheelSource(wheel, 220)
	if sourceY != 20 {
		t.Fatalf("expected virtual slot 11 to map to source digit 1, got sourceY=%d", sourceY)
	}
}

func TestOdometerWheelSourceMapsVirtualSlotNegativeOneToDigitNine(t *testing.T) {
	wheel := OdometerWheel{Window: Size{Width: 12, Height: 20}}
	_, sourceY := odometerWheelSource(wheel, -20)
	if sourceY != 180 {
		t.Fatalf("expected virtual slot -1 to map to source digit 9, got sourceY=%d", sourceY)
	}
}

func TestOdometerWheelSlicesRenderAdjacentDigitsForVirtualSlotOnePointFive(t *testing.T) {
	wheel := OdometerWheel{Window: Size{Width: 12, Height: 20}}
	slices := odometerWheelSlices(wheel, 30)
	if got := wheelSliceDigits(slices); !intSlicesEqual(got, []int{1, 2}) {
		t.Fatalf("expected virtual slot 1.5 to render digits 1 and 2, got %v", got)
	}
}

func TestOdometerWheelSlicesRenderAdjacentDigitsForVirtualSlotEightPointFive(t *testing.T) {
	wheel := OdometerWheel{Window: Size{Width: 12, Height: 20}}
	slices := odometerWheelSlices(wheel, 170)
	if got := wheelSliceDigits(slices); !intSlicesEqual(got, []int{8, 9}) {
		t.Fatalf("expected virtual slot 8.5 to render digits 8 and 9, got %v", got)
	}
}

func TestOdometerWheelSlicesRenderAdjacentDigitsForVirtualSlotNinePointFive(t *testing.T) {
	wheel := OdometerWheel{Window: Size{Width: 12, Height: 20}}
	slices := odometerWheelSlices(wheel, 190)
	if got := wheelSliceDigits(slices); !intSlicesEqual(got, []int{9, 0}) {
		t.Fatalf("expected virtual slot 9.5 to render digits 9 and 0, got %v", got)
	}
}

func TestOdometerWheelSlicesRenderAdjacentDigitsForVirtualSlotNegativeHalf(t *testing.T) {
	wheel := OdometerWheel{Window: Size{Width: 12, Height: 20}}
	slices := odometerWheelSlices(wheel, -10)
	if got := wheelSliceDigits(slices); !intSlicesEqual(got, []int{9, 0}) {
		t.Fatalf("expected virtual slot -0.5 to render digits 9 and 0, got %v", got)
	}
}

func TestOdometerCarryDragDisabledKeepsBaseWheelOffsets(t *testing.T) {
	pkg := loadOdometerScenePackageWithRealism(t, MovementLinear, true, false, false, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 19.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 20.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	base, err := OdometerInterpolatedWheelOffsets(pkg, 19.9, 20.0, previousOffsets, targetOffsets, 0.85)
	if err != nil {
		t.Fatalf("OdometerInterpolatedWheelOffsets failed: %v", err)
	}
	adjusted, err := OdometerCarryDragWheelOffsets(pkg, 19.9, 20.0, previousOffsets, targetOffsets, base)
	if err != nil {
		t.Fatalf("OdometerCarryDragWheelOffsets failed: %v", err)
	}
	if !float64SlicesAlmostEqual(base, adjusted) {
		t.Fatalf("expected carry_drag disabled to keep base offsets, got base=%v adjusted=%v", base, adjusted)
	}
}

func TestOdometerCarryDragEnabledAdvancesHigherWheelNearRollover(t *testing.T) {
	pkg := loadOdometerScenePackageWithRealism(t, MovementLinear, true, true, false, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 19.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 20.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	base, err := OdometerInterpolatedWheelOffsets(pkg, 19.9, 20.0, previousOffsets, targetOffsets, 0.9)
	if err != nil {
		t.Fatalf("OdometerInterpolatedWheelOffsets failed: %v", err)
	}
	adjusted, err := OdometerCarryDragWheelOffsets(pkg, 19.9, 20.0, previousOffsets, targetOffsets, base)
	if err != nil {
		t.Fatalf("OdometerCarryDragWheelOffsets failed: %v", err)
	}
	if !(adjusted[0] > base[0] && adjusted[0] < targetOffsets[0]) {
		t.Fatalf("expected carry_drag to advance tens wheel toward target, got base=%v adjusted=%v target=%v", base[0], adjusted[0], targetOffsets[0])
	}
	if !(adjusted[1] > base[1]) {
		t.Fatalf("expected carry_drag to advance ones wheel toward its routed rollover, got base=%v adjusted=%v target=%v", base[1], adjusted[1], targetOffsets[1])
	}
}

func TestOdometerCarryDragStraddlingUpdateStartsBeforeLowerWheelPassesRollover(t *testing.T) {
	pkg := loadOdometerScenePackageWithRealism(t, MovementLinear, true, true, false, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 19.8)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 20.2)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	base, err := OdometerInterpolatedWheelOffsets(pkg, 19.8, 20.2, previousOffsets, targetOffsets, 0.45)
	if err != nil {
		t.Fatalf("OdometerInterpolatedWheelOffsets failed: %v", err)
	}
	adjusted, err := OdometerCarryDragWheelOffsets(pkg, 19.8, 20.2, previousOffsets, targetOffsets, base)
	if err != nil {
		t.Fatalf("OdometerCarryDragWheelOffsets failed: %v", err)
	}
	rolloverOffset, err := odometerDiscreteWheelOffset(pkg.Odometer.Wheels[1], odometerDigitPlaces(pkg.Odometer.Wheels)[1], 20.0)
	if err != nil {
		t.Fatalf("odometerDiscreteWheelOffset failed: %v", err)
	}
	rolloverOffset += 200
	if !(base[1] < rolloverOffset) {
		t.Fatalf("expected lower wheel to still be approaching rollover, got base ones offset %.2f with rollover %.2f", base[1], rolloverOffset)
	}
	if !(adjusted[0] > base[0]) {
		t.Fatalf("expected carry_drag to advance tens wheel before rollover on straddling update, got base=%v adjusted=%v", base[0], adjusted[0])
	}
	if !(adjusted[1] >= base[1]) {
		t.Fatalf("expected lower wheel base offset to stay monotonic, got base=%v adjusted=%v", base[1], adjusted[1])
	}
}

func TestOdometerCarryDragSkipsMultiRolloverSpanForWheelPair(t *testing.T) {
	pkg := loadOdometerScenePackageWithRealism(t, MovementLinear, true, true, false, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 20.8)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 31.2)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	base, err := OdometerInterpolatedWheelOffsets(pkg, 20.8, 31.2, previousOffsets, targetOffsets, 0.2)
	if err != nil {
		t.Fatalf("OdometerInterpolatedWheelOffsets failed: %v", err)
	}
	adjusted, err := OdometerCarryDragWheelOffsets(pkg, 20.8, 31.2, previousOffsets, targetOffsets, base)
	if err != nil {
		t.Fatalf("OdometerCarryDragWheelOffsets failed: %v", err)
	}
	if !float64SlicesAlmostEqual(base, adjusted) {
		t.Fatalf("expected multi-rollover carry_drag span to skip pair drag, got base=%v adjusted=%v", base, adjusted)
	}
}

func TestOdometerSnapSettleDisabledKeepsBaseWheelOffsets(t *testing.T) {
	pkg := loadOdometerScenePackageWithRealism(t, MovementLinear, false, false, false, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 12.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 12.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	base := cloneFloat64s(targetOffsets)
	adjusted, err := OdometerSnapSettleWheelOffsets(pkg, 12.0, 12.9, previousOffsets, targetOffsets, base, 0.35)
	if err != nil {
		t.Fatalf("OdometerSnapSettleWheelOffsets failed: %v", err)
	}
	if !float64SlicesAlmostEqual(base, adjusted) {
		t.Fatalf("expected snap_settle disabled to keep base offsets, got base=%v adjusted=%v", base, adjusted)
	}
}

func TestOdometerSnapSettleEnabledAddsSmallForwardSettleAndReturnsToTarget(t *testing.T) {
	pkg := loadOdometerScenePackageWithRealism(t, MovementLinear, false, false, true, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 12.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 12.9)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	base := cloneFloat64s(targetOffsets)
	adjusted, err := OdometerSnapSettleWheelOffsets(pkg, 12.0, 12.9, previousOffsets, targetOffsets, base, 0.35)
	if err != nil {
		t.Fatalf("OdometerSnapSettleWheelOffsets failed: %v", err)
	}
	if !(adjusted[2] > targetOffsets[2]) {
		t.Fatalf("expected snap_settle to nudge moving wheels past target, got target=%v adjusted=%v", targetOffsets, adjusted)
	}
	if !almostEqual(adjusted[0], targetOffsets[0]) || !almostEqual(adjusted[1], targetOffsets[1]) {
		t.Fatalf("expected unchanged wheels to stay on exact target slots, got target=%v adjusted=%v", targetOffsets, adjusted)
	}

	settled, err := OdometerSnapSettleWheelOffsets(pkg, 12.0, 12.9, previousOffsets, targetOffsets, base, 1)
	if err != nil {
		t.Fatalf("OdometerSnapSettleWheelOffsets failed: %v", err)
	}
	if !float64SlicesAlmostEqual(targetOffsets, settled) {
		t.Fatalf("expected snap_settle to settle exactly on target, got target=%v settled=%v", targetOffsets, settled)
	}
}

func TestOdometerSnapSettleDoesNotOvershootBelowZeroAtLowerBoundary(t *testing.T) {
	pkg := loadOdometerScenePackageWithRealism(t, MovementLinear, false, false, true, nil)
	previousOffsets, err := OdometerWheelStripOffsets(pkg, 1.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	targetOffsets, err := OdometerWheelStripOffsets(pkg, 0.0)
	if err != nil {
		t.Fatalf("OdometerWheelStripOffsets failed: %v", err)
	}
	base := cloneFloat64s(targetOffsets)
	adjusted, err := OdometerSnapSettleWheelOffsets(pkg, 1.0, 0.0, previousOffsets, targetOffsets, base, 0.35)
	if err != nil {
		t.Fatalf("OdometerSnapSettleWheelOffsets failed: %v", err)
	}

	for index, offset := range adjusted {
		if offset < 0 {
			t.Fatalf("expected wheel %d settle offset to stay at or above zero, got %.2f", index, offset)
		}
	}
	if !almostEqual(adjusted[1], 0) {
		t.Fatalf("expected ones wheel at lower boundary to stay clamped to zero, got %.2f", adjusted[1])
	}
}

func TestOdometerSceneAppliesConfiguredDrumSlopToWheelPositions(t *testing.T) {
	pkg := loadOdometerScenePackage(t, "", false, []int{2, -1, 3})

	scene, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 12.3))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	wheels := wheelStripParts(scene)
	if len(wheels) != 3 {
		t.Fatalf("wheel parts = %d, want 3", len(wheels))
	}
	if wheels[0].Position[1] != 14 || wheels[1].Position[1] != 11 || wheels[2].Position[1] != 15 {
		t.Fatalf("drum slop positions = %v/%v/%v, want y 14/11/15", wheels[0].Position, wheels[1].Position, wheels[2].Position)
	}
	if !almostEqual(wheels[0].StripOffset, 20) || !almostEqual(wheels[1].StripOffset, 40) || !almostEqual(wheels[2].StripOffset, 60) {
		t.Fatalf("expected drum slop to keep wheel strip offsets unchanged, got %.2f/%.2f/%.2f", wheels[0].StripOffset, wheels[1].StripOffset, wheels[2].StripOffset)
	}
}

func TestOdometerSceneDefaultsToNoDrumSlop(t *testing.T) {
	pkg := loadOdometerScenePackage(t, "", false, nil)

	scene, err := OdometerScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("trip_distance", 12.3))
	if err != nil {
		t.Fatalf("OdometerScene returned error: %v", err)
	}

	wheels := wheelStripParts(scene)
	if len(wheels) != 3 {
		t.Fatalf("wheel parts = %d, want 3", len(wheels))
	}
	if wheels[0].Position[1] != 12 || wheels[1].Position[1] != 12 || wheels[2].Position[1] != 12 {
		t.Fatalf("default wheel positions = %v/%v/%v, want y=12 for all wheels", wheels[0].Position, wheels[1].Position, wheels[2].Position)
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

func TestSegmentedSceneUsesSparseThresholdsAndHysteresis(t *testing.T) {
	pkg := loadSegmentedScenePackage(t)

	below, _, err := SegmentedScene(pkg, Placement{Position: []int{10, 20}, Scale: 1.25}, okGaugeStateWithRange("rpm", -100, 0, 7000), nil)
	if err != nil {
		t.Fatalf("SegmentedScene returned error: %v", err)
	}
	if got := partLayerSequence(below); got != "layer:panel,layer:glass" {
		t.Fatalf("below-threshold sequence = %q", got)
	}

	first, previous, err := SegmentedScene(pkg, Placement{Position: []int{10, 20}, Scale: 1.25}, okGaugeStateWithRange("rpm", 3500, 0, 7000), nil)
	if err != nil {
		t.Fatalf("SegmentedScene returned error: %v", err)
	}
	if got := partLayerSequence(first); got != "layer:panel,layer:segments,layer:glass" {
		t.Fatalf("first threshold sequence = %q", got)
	}
	selected := firstLayerPart(first, "segments")
	if selected.Layer != "segments" || !strings.HasSuffix(selected.AssetPath, "rpm_050.png") {
		t.Fatalf("selected segment = %#v", selected)
	}

	highest, _, err := SegmentedScene(pkg, Placement{Position: []int{10, 20}, Scale: 1.25}, okGaugeStateWithRange("rpm", 7000, 0, 7000), nil)
	if err != nil {
		t.Fatalf("SegmentedScene returned error: %v", err)
	}
	if got := firstLayerPart(highest, "segments"); got.Layer != "segments" || !strings.HasSuffix(got.AssetPath, "rpm_100.png") {
		t.Fatalf("highest segment = %#v", got)
	}

	aboveMax, _, err := SegmentedScene(pkg, Placement{Position: []int{10, 20}, Scale: 1.25}, okGaugeStateWithRange("rpm", 9000, 0, 7000), nil)
	if err != nil {
		t.Fatalf("SegmentedScene returned error: %v", err)
	}
	if got := firstLayerPart(aboveMax, "segments"); got.Layer != "segments" || !strings.HasSuffix(got.AssetPath, "rpm_100.png") {
		t.Fatalf("above-max segment = %#v", got)
	}

	percent, _, err := SegmentedScene(pkg, Placement{Position: []int{10, 20}, Scale: 1.25}, sensors.SensorState{ID: "rpm", Value: 60, Min: 0, Max: 0, Status: sensors.StatusOK}, nil)
	if err != nil {
		t.Fatalf("SegmentedScene returned error: %v", err)
	}
	if got := firstLayerPart(percent, "segments"); got.Layer != "segments" || !strings.HasSuffix(got.AssetPath, "rpm_050.png") {
		t.Fatalf("percent segment = %#v", got)
	}

	held, previous, err := SegmentedScene(pkg, Placement{Position: []int{10, 20}, Scale: 1.25}, okGaugeStateWithRange("rpm", 3080, 0, 7000), previous)
	if err != nil {
		t.Fatalf("SegmentedScene returned error: %v", err)
	}
	if got := firstLayerPart(held, "segments"); got.AssetPath == "" || !strings.HasSuffix(got.AssetPath, "rpm_050.png") {
		t.Fatalf("held segment = %#v", got)
	}

	dropped, _, err := SegmentedScene(pkg, Placement{Position: []int{10, 20}, Scale: 1.25}, okGaugeStateWithRange("rpm", 3000, 0, 7000), previous)
	if err != nil {
		t.Fatalf("SegmentedScene returned error: %v", err)
	}
	if got := firstLayerPart(dropped, "segments"); got.AssetPath == "" || !strings.HasSuffix(got.AssetPath, "rpm_025.png") {
		t.Fatalf("dropped segment = %#v", got)
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
	return loadRadialScenePackageWithRealism(t, nil, nil, nil)
}

func loadRadialScenePackageWithNeedleShadow(t *testing.T, offset []int, alpha *float64) Package {
	return loadRadialScenePackageWithRealism(t, offset, alpha, nil)
}

func loadRadialScenePackageWithCalibrationOffset(t *testing.T, calibrationOffset *float64) Package {
	return loadRadialScenePackageWithRealism(t, nil, nil, calibrationOffset)
}

func loadRadialScenePackageWithRealism(t *testing.T, offset []int, alpha *float64, calibrationOffset *float64) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "simple_rpm")
	writeGaugeYAML(t, packageDir, radialGaugeYAML(offset, alpha, calibrationOffset))
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func loadOdometerScenePackage(t *testing.T, movement string, wraparound bool, drumSlop []int) Package {
	return loadOdometerScenePackageWithRealism(t, movement, wraparound, false, false, drumSlop)
}

func loadOdometerScenePackageWithRealism(t *testing.T, movement string, wraparound bool, carryDrag bool, snapSettle bool, drumSlop []int) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "trip")
	writeGaugeYAML(t, packageDir, odometerGaugeYAML(movement, wraparound, carryDrag, snapSettle, drumSlop))
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func loadFourWheelOdometerScenePackage(t *testing.T) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "trip_four_wheel")
	writeGaugeYAML(t, packageDir, fourWheelOdometerGaugeYAML())
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

func loadSegmentedScenePackage(t *testing.T) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "segmented", "rpm")
	files := []string{
		"panel.png",
		"glass.png",
		"levels/rpm_025.png",
		"levels/rpm_050.png",
		"levels/rpm_100.png",
		"levels/rpm_150.png",
	}
	for _, path := range files {
		fullPath := filepath.Join(packageDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(path), 0o600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}
	writeGaugeYAML(t, packageDir, `id: test_rpm_segmented
type: segmented
sensor: rpm
size:
  width: 120
  height: 120
layers:
  panel: panel.png
  segments: levels/rpm_{percent:03}.png
  glass: glass.png
`)
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

func radialGaugeYAML(offset []int, alpha *float64, calibrationOffset *float64) string {
	realismLines := []string{}
	if len(offset) == 2 {
		realismLines = append(realismLines, "  needle_shadow:")
		realismLines = append(realismLines, fmt.Sprintf("    offset: [%d, %d]", offset[0], offset[1]))
		if alpha != nil {
			realismLines = append(realismLines, fmt.Sprintf("    alpha: %.3f", *alpha))
		}
	}
	if calibrationOffset != nil {
		realismLines = append(realismLines, fmt.Sprintf("  calibration_offset: %.3f", *calibrationOffset))
	}
	realismBlock := ""
	if len(realismLines) > 0 {
		realismBlock = "realism:\n" + strings.Join(realismLines, "\n") + "\n"
	}
	return `id: simple_radial_rpm
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

func odometerGaugeYAML(movement string, wraparound bool, carryDrag bool, snapSettle bool, drumSlop []int) string {
	movementLine := ""
	if movement != "" {
		movementLine = fmt.Sprintf("  movement: %s\n", movement)
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
	if drumSlop != nil {
		values := make([]string, len(drumSlop))
		for i, slop := range drumSlop {
			values[i] = strconv.Itoa(slop)
		}
		realismLines = append(realismLines, "  drum_slop: ["+strings.Join(values, ", ")+"]")
	}
	realismBlock := ""
	if len(realismLines) > 0 {
		realismBlock = "realism:\n" + strings.Join(realismLines, "\n") + "\n"
	}
	return fmt.Sprintf(`id: test_trip_odometer
type: odometer
sensor: trip_distance
%ssize:
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
`, realismBlock, movementLine)
}

func fourWheelOdometerGaugeYAML() string {
	return `id: test_trip_odometer_four_wheel
type: odometer
sensor: trip_distance
size:
  width: 168
  height: 60
layers:
  panel: panel.png
  glass: glass.png
odometer:
  movement: linear
  wheels:
    - strip: digits.png
      position: [10, 12]
      window: { width: 12, height: 20 }
    - strip: digits.png
      position: [24, 12]
      window: { width: 12, height: 20 }
    - strip: digits.png
      position: [38, 12]
      window: { width: 12, height: 20 }
    - strip: red_digits.png
      position: [56, 12]
      window: { width: 12, height: 20 }
      offset: [2, 4]
      role: sub_unit
`
}

func float64SlicesAlmostEqual(left []float64, right []float64) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if !almostEqual(left[index], right[index]) {
			return false
		}
	}
	return true
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

func okGaugeStateWithRange(id string, value, min, max float64) sensors.SensorState {
	return sensors.SensorState{ID: id, Value: value, Min: min, Max: max, Status: sensors.StatusOK}
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
		case ScenePartKindNeedleShadow:
			parts = append(parts, fmt.Sprintf("needle_shadow:%.0f", part.Angle))
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

func wheelSliceDigits(slices []WheelSlice) []int {
	digits := make([]int, len(slices))
	for index, slice := range slices {
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

func firstLayerPart(scene Scene, layer string) ScenePart {
	for _, part := range scene.Parts {
		if part.Kind == ScenePartKindLayer && part.Layer == layer {
			return part
		}
	}
	return ScenePart{}
}
