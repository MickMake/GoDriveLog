package gauges

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestSevenSegmentSceneUsesPackageOwnedFormatPositionsAndStaticLayers(t *testing.T) {
	pkg := loadSevenSegmentScenePackage(t, 4, "%04.0f")

	scene, err := SevenSegmentScene(pkg, Placement{Position: []int{780, 40}, Scale: 1.5}, okGaugeState("rpm", 12))
	if err != nil {
		t.Fatalf("SevenSegmentScene returned error: %v", err)
	}

	if scene.PackageID != "test_4_digit_rpm" || scene.SensorID != "rpm" || scene.Type != TypeSevenSegment {
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

func TestSevenSegmentSceneEmitsPanelUnderDigitsAndGlassOverDigits(t *testing.T) {
	pkg := loadSevenSegmentScenePackage(t, 4, "%04.0f")
	scene, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 12))
	if err != nil {
		t.Fatalf("SevenSegmentScene returned error: %v", err)
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

func TestSevenSegmentSceneSupportsTwoThreeFourAndFiveDigitShapes(t *testing.T) {
	for _, count := range []int{2, 3, 4, 5} {
		t.Run(fmt.Sprintf("%d_digits", count), func(t *testing.T) {
			format := fmt.Sprintf("%%0%d.0f", count)
			pkg := loadSevenSegmentScenePackage(t, count, format)
			scene, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 12))
			if err != nil {
				t.Fatalf("SevenSegmentScene returned error: %v", err)
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

func TestSevenSegmentSceneDoesNotRenderLiveDigitsForNonOKStates(t *testing.T) {
	pkg := loadSevenSegmentScenePackage(t, 4, "%04.0f")
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
			scene, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "rpm", Status: status, Error: "not live"})
			if err != nil {
				t.Fatalf("SevenSegmentScene returned error: %v", err)
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

func TestSevenSegmentSceneSignatureChangesWithFormattedOutput(t *testing.T) {
	pkg := loadSevenSegmentScenePackage(t, 4, "%04.0f")
	first, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 42.1))
	if err != nil {
		t.Fatalf("SevenSegmentScene returned error: %v", err)
	}
	unchanged, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 42.2))
	if err != nil {
		t.Fatalf("SevenSegmentScene returned error: %v", err)
	}
	changed, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 43))
	if err != nil {
		t.Fatalf("SevenSegmentScene returned error: %v", err)
	}

	if first.Signature() != unchanged.Signature() {
		t.Fatalf("expected unchanged rounded output to keep same signature")
	}
	if first.Signature() == changed.Signature() {
		t.Fatalf("expected changed formatted output to change signature")
	}
}

func TestSevenSegmentSceneSignatureIncludesDigitPositionsForNonOKState(t *testing.T) {
	pkg := loadSevenSegmentScenePackage(t, 4, "%04.0f")
	first, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout})
	if err != nil {
		t.Fatalf("SevenSegmentScene returned error: %v", err)
	}
	pkg.Digits.Positions[2] = []int{999, 12}
	changed, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, sensors.SensorState{ID: "rpm", Status: sensors.StatusTimeout})
	if err != nil {
		t.Fatalf("SevenSegmentScene returned error: %v", err)
	}
	if first.Signature() == changed.Signature() {
		t.Fatalf("expected non-ok signature to change when digit positions change")
	}
}

func TestSevenSegmentSceneRejectsMissingDecimalPointWhenFormatNeedsIt(t *testing.T) {
	pkg := loadSevenSegmentScenePackage(t, 4, "%.1f")
	pkg.DigitSet.DecimalPoint = ""

	_, err := SevenSegmentScene(pkg, Placement{Position: []int{0, 0}, Scale: 1}, okGaugeState("rpm", 12.3))
	if err == nil {
		t.Fatal("expected missing decimal point to fail")
	}
	assertErrorContains(t, err, "decimal_point")
}

func loadSevenSegmentScenePackage(t *testing.T, count int, format string) Package {
	t.Helper()
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "7Seg", "amber", fmt.Sprintf("%d_digit_rpm", count))
	writeGaugeYAML(t, packageDir, sevenSegmentGaugeYAML(count, format))
	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	return pkg
}

func sevenSegmentGaugeYAML(count int, format string) string {
	var positions strings.Builder
	for slot := 0; slot < count; slot++ {
		positions.WriteString(fmt.Sprintf("    - [%d, 12]\n", slot*10+2))
	}
	return fmt.Sprintf(`id: test_%d_digit_rpm
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
		default:
			parts = append(parts, part.Kind)
		}
	}
	return strings.Join(parts, ",")
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
