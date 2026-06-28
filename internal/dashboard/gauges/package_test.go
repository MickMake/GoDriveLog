package gauges

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPackageLoadsNumericGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "7Seg", "amber", "4_digit_rpm")
	writeGaugeYAML(t, packageDir, `id: amber_4_digit_rpm
type: numeric
sensor: rpm
format: "%04.0f"
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
  decimal_point: ../7SegDP.png
  spacing: 4
digits:
  count: 4
  positions:
    - [35, 35]
    - [117, 35]
    - [199, 35]
    - [281, 35]
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}

	if pkg.ID != "amber_4_digit_rpm" || pkg.Type != TypeNumeric || pkg.Sensor != "rpm" || pkg.Format != "%04.0f" {
		t.Fatalf("package identity = %#v", pkg)
	}
	if pkg.Size.Width != 398 || pkg.Size.Height != 150 {
		t.Fatalf("size = %#v, want 398x150", pkg.Size)
	}
	assertPath(t, pkg.Path, packageDir)
	assertPath(t, pkg.YAMLPath, filepath.Join(packageDir, "gauge.yaml"))
	assertPath(t, pkg.AssetRoot, filepath.Join(root, "assets"))
	assertPath(t, pkg.Layers["panel"], filepath.Join(root, "assets", "gauges", "7Seg", "7Seg4Digits.png"))
	assertPath(t, pkg.Layers["glass"], filepath.Join(root, "assets", "gauges", "7Seg", "Glass.png"))
	assertPath(t, pkg.DigitSet.Background, filepath.Join(root, "assets", "gauges", "7Seg", "7SegBack.png"))
	assertPath(t, pkg.DigitSet.Characters["0"], filepath.Join(root, "assets", "gauges", "7Seg", "amber", "7Seg0.png"))
	assertPath(t, pkg.DigitSet.DecimalPoint, filepath.Join(root, "assets", "gauges", "7Seg", "amber", "7SegDP.png"))
	if pkg.Digits.Count != 4 || len(pkg.Digits.Positions) != 4 || pkg.Digits.Positions[2][0] != 199 {
		t.Fatalf("digits = %#v, want four resolved positions", pkg.Digits)
	}
}

func TestDefaultGaugeSearchPathsIncludesDashboardConfigEnvPath(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "dashboard.yaml")
	if err := os.Setenv(dashboardConfigEnvVar, configPath); err != nil {
		t.Fatalf("Setenv: %v", err)
	}
	defer func() {
		_ = os.Unsetenv(dashboardConfigEnvVar)
	}()

	paths := defaultGaugeSearchPaths()
	want := filepath.Dir(configPath)
	if !containsPath(paths, want) {
		t.Fatalf("defaultGaugeSearchPaths() = %v, want %q", paths, want)
	}
}

func containsPath(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}

func TestLoadPackageLoadsRadialGaugeFromArbitraryDirectory(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "random", "not_semantic", "rpm_round")
	writeGaugeYAML(t, packageDir, `id: simple_radial_rpm
type: radial
sensor: rpm
size:
  width: 512
  height: 512
layers:
  background: ../../../shared/radial/simple_rpm/bezel.png
  face: ../../../shared/radial/simple_rpm/face.png
  ticks: ../../../shared/radial/simple_rpm/ticks.png
  needle: ../../../shared/radial/simple_rpm/needle.png
  overlay: ../../../shared/radial/simple_rpm/glass.png
pivot:
  face: { x: 0.5, y: 0.55 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 7000
  start_angle: -135
  end_angle: 135
  clamp: true
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}

	if pkg.Type != TypeRadial || pkg.Sensor != "rpm" {
		t.Fatalf("package = %#v", pkg)
	}
	assertPath(t, pkg.Layers["needle"], filepath.Join(root, "assets", "gauges", "shared", "radial", "simple_rpm", "needle.png"))
	if pkg.Pivot.Face.X != 0.5 || pkg.Pivot.Needle.Y != 0.9 {
		t.Fatalf("pivot = %#v", pkg.Pivot)
	}
	if pkg.ValueMap.Max != 7000 || pkg.ValueMap.StartAngle != -135 || !pkg.ValueMap.Clamp {
		t.Fatalf("value_map = %#v", pkg.ValueMap)
	}
}

func TestLoadPackageLoadsOdometerGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "trip")
	writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
size:
  width: 240
  height: 80
layers:
  panel: panel.png
  glass: glass.png
odometer:
  wheels:
    - strip: digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
    - strip: red_digits.png
      position: [40, 12]
      window: { width: 24, height: 36 }
      offset: [2, 4]
      role: sub_unit
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}

	if pkg.ID != "trip_odometer" || pkg.Type != TypeOdometer || pkg.Sensor != "trip_distance" {
		t.Fatalf("package identity = %#v", pkg)
	}
	if pkg.Odometer.Movement != MovementInstant {
		t.Fatalf("movement = %q, want default instant", pkg.Odometer.Movement)
	}
	if pkg.Realism.MovementPolicy != MovementPolicyImmediate {
		t.Fatalf("movement policy = %q, want default immediate", pkg.Realism.MovementPolicy)
	}
	if pkg.Realism.DrumSlopSet {
		t.Fatalf("expected omitted drum_slop to remain absent, got %#v", pkg.Realism)
	}
	if len(pkg.Odometer.Wheels) != 2 {
		t.Fatalf("wheels = %#v, want 2", pkg.Odometer.Wheels)
	}
	assertPath(t, pkg.Odometer.Wheels[0].Strip, filepath.Join(root, "assets", "gauges", "odometer", "trip", "digits.png"))
	assertPath(t, pkg.Odometer.Wheels[1].Strip, filepath.Join(root, "assets", "gauges", "odometer", "trip", "red_digits.png"))
	if pkg.Odometer.Wheels[1].Role != WheelRoleSubUnit || pkg.Odometer.Wheels[1].Offset[0] != 2 {
		t.Fatalf("sub-unit wheel = %#v", pkg.Odometer.Wheels[1])
	}
}

func TestLoadPackageAcceptsImplementedOdometerMovementValues(t *testing.T) {
	movements := []string{
		MovementInstant,
		MovementLinear,
		MovementEaseOut,
		MovementBell,
	}

	for _, movement := range movements {
		t.Run(movement, func(t *testing.T) {
			root := makeGaugeFixtures(t)
			packageDir := filepath.Join(root, "assets", "gauges", "odometer", movement)
			writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
size:
  width: 240
  height: 80
odometer:
  movement: `+movement+`
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
`)

			pkg, err := LoadPackage(packageDir)
			if err != nil {
				t.Fatalf("LoadPackage returned error: %v", err)
			}
			if pkg.Odometer.Movement != movement {
				t.Fatalf("movement = %q, want %q", pkg.Odometer.Movement, movement)
			}
		})
	}
}

func TestLoadPackageWarnsAndFallsBackForRecognizedOdometerMovementValues(t *testing.T) {
	tests := []string{MovementSmooth, MovementClick}

	for _, movement := range tests {
		t.Run(movement, func(t *testing.T) {
			root := makeGaugeFixtures(t)
			packageDir := filepath.Join(root, "assets", "gauges", "odometer", movement)
			writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
size:
  width: 240
  height: 80
odometer:
  movement: `+movement+`
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
`)

			var pkg Package
			logOutput := capturePackageLogs(t, func() {
				var err error
				pkg, err = LoadPackage(packageDir)
				if err != nil {
					t.Fatalf("LoadPackage returned error: %v", err)
				}
			})
			if pkg.Odometer.Movement != MovementInstant {
				t.Fatalf("movement = %q, want fallback %q", pkg.Odometer.Movement, MovementInstant)
			}
			if !strings.Contains(logOutput, `odometer movement "`+movement+`" is recognised but not implemented`) {
				t.Fatalf("warning log = %q, want mention of %q fallback", logOutput, movement)
			}
		})
	}
}

func TestLoadPackageRejectsUnknownOdometerMovementValue(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "unknown_movement")
	writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
size:
  width: 240
  height: 80
odometer:
  movement: wobble
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, `odometer movement "wobble" is not supported`)
}

func TestLoadPackageLoadsOdometerWraparoundRealism(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "wrap")
	writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
realism:
  wraparound: true
size:
  width: 240
  height: 80
odometer:
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	if pkg.Realism.Wraparound == nil || !*pkg.Realism.Wraparound {
		t.Fatalf("wraparound realism = %#v, want true", pkg.Realism)
	}
}

func TestLoadPackageLoadsOdometerCarryDragRealism(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "carry_drag")
	writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
realism:
  carry_drag: true
size:
  width: 240
  height: 80
odometer:
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	if pkg.Realism.CarryDrag == nil || !*pkg.Realism.CarryDrag {
		t.Fatalf("carry_drag realism = %#v, want true", pkg.Realism)
	}
}

func TestLoadPackageLoadsOdometerDrumSlopRealism(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "slop")
	writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
realism:
  drum_slop: [1, -2]
size:
  width: 240
  height: 80
odometer:
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
    - strip: ../trip/red_digits.png
      position: [40, 12]
      window: { width: 24, height: 36 }
      role: sub_unit
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}
	if !pkg.Realism.DrumSlopSet {
		t.Fatalf("expected drum_slop to be marked present, got %#v", pkg.Realism)
	}
	if len(pkg.Realism.DrumSlop) != 2 || pkg.Realism.DrumSlop[0] != 1 || pkg.Realism.DrumSlop[1] != -2 {
		t.Fatalf("drum slop realism = %#v, want [1 -2]", pkg.Realism.DrumSlop)
	}
}

func TestLoadPackageLoadsIndicatorGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "indicator", "check_engine")
	writeGaugeYAML(t, packageDir, `id: check_engine_indicator
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
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}

	if pkg.ID != "check_engine_indicator" || pkg.Type != TypeIndicator || pkg.Sensor != "check_engine" {
		t.Fatalf("package identity = %#v", pkg)
	}
	assertPath(t, pkg.Layers["off"], filepath.Join(root, "assets", "gauges", "indicator", "check_engine", "off.png"))
	assertPath(t, pkg.Layers["on"], filepath.Join(root, "assets", "gauges", "indicator", "check_engine", "on.png"))
}

func TestLoadPackageLoadsIndicatorGaugeWithOnlyOnLayer(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "indicator", "check_engine")
	writeGaugeYAML(t, packageDir, `id: check_engine_indicator
type: indicator
sensor: check_engine
size:
  width: 48
  height: 48
layers:
  bezel: bezel.png
  on: on.png
  glass: glass.png
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}

	if pkg.Type != TypeIndicator || pkg.Sensor != "check_engine" {
		t.Fatalf("package identity = %#v", pkg)
	}
	assertPath(t, pkg.Layers["on"], filepath.Join(root, "assets", "gauges", "indicator", "check_engine", "on.png"))
	if _, ok := pkg.Layers["off"]; ok {
		t.Fatalf("off layer should be absent, got %#v", pkg.Layers["off"])
	}
}

func TestLoadPackageLoadsBarGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "bar", "coolant")
	writeGaugeYAML(t, packageDir, `id: coolant_bar
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
`)

	pkg, err := LoadPackage(packageDir)
	if err != nil {
		t.Fatalf("LoadPackage returned error: %v", err)
	}

	if pkg.ID != "coolant_bar" || pkg.Type != TypeBar || pkg.Sensor != "coolant_temperature" {
		t.Fatalf("package identity = %#v", pkg)
	}
	if pkg.Bar.Mode != "level" || pkg.Bar.Axis != "vertical" || pkg.Bar.Origin != "bottom" {
		t.Fatalf("bar config = %#v", pkg.Bar)
	}
	if pkg.ValueMap.Min != 40 || pkg.ValueMap.Max != 120 || !pkg.ValueMap.Clamp {
		t.Fatalf("value_map = %#v", pkg.ValueMap)
	}
	if len(pkg.Bar.Bounds) != 4 || pkg.Bar.Bounds[0] != 40 || pkg.Bar.Bounds[3] != 180 {
		t.Fatalf("bar bounds = %#v", pkg.Bar.Bounds)
	}
	assertPath(t, pkg.Layers["level"], filepath.Join(root, "assets", "gauges", "bar", "coolant", "level.png"))
}

func TestLoadPackageLoadsSegmentedGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "segmented", "rpm")
	files := []string{
		"levels/rpm_000.png",
		"levels/rpm_025.png",
		"levels/rpm_050.png",
		"levels/rpm_150.png",
		"panel.png",
		"glass.png",
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
	writeGaugeYAML(t, packageDir, `id: rpm_segmented
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

	if pkg.ID != "rpm_segmented" || pkg.Type != TypeSegmented || pkg.Sensor != "rpm" {
		t.Fatalf("package identity = %#v", pkg)
	}
	if pkg.Segmented.Hysteresis == nil || *pkg.Segmented.Hysteresis != 25 {
		t.Fatalf("hysteresis = %#v, want default 25", pkg.Segmented.Hysteresis)
	}
	if len(pkg.Segmented.Images) != 3 {
		t.Fatalf("segmented images = %#v, want 3 valid thresholds", pkg.Segmented.Images)
	}
	if pkg.Segmented.Images[0].Threshold != 0 || pkg.Segmented.Images[1].Threshold != 25 || pkg.Segmented.Images[2].Threshold != 50 {
		t.Fatalf("segmented thresholds = %#v", pkg.Segmented.Images)
	}
	assertPath(t, pkg.Segmented.Images[0].Path, filepath.Join(packageDir, "levels", "rpm_000.png"))
	assertPath(t, pkg.Segmented.Images[2].Path, filepath.Join(packageDir, "levels", "rpm_050.png"))
}

func TestLoadPackageRejectsIndicatorMissingStateLayer(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "indicator", "bad")
	writeGaugeYAML(t, packageDir, `id: bad_indicator
type: indicator
sensor: check_engine
size:
  width: 48
  height: 48
layers:
  off: off.png
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "indicator layer on")
}

func TestLoadPackageRejectsBadOdometerMovement(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "bad")
	writeGaugeYAML(t, packageDir, `id: bad_odometer
type: odometer
sensor: trip_distance
size:
  width: 100
  height: 50
odometer:
  movement: elastic
  wheels:
    - strip: digits.png
      position: [0, 0]
      window: { width: 10, height: 20 }
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "movement")
}

func TestLoadPackageRejectsWraparoundOnNonOdometerGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "bad_wraparound")
	writeGaugeYAML(t, packageDir, `id: bad_radial
type: radial
sensor: rpm
realism:
  wraparound: true
size:
  width: 100
  height: 100
layers:
  needle: ../../shared/radial/simple_rpm/needle.png
pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 100
  start_angle: -90
  end_angle: 90
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "wraparound")
}

func TestLoadPackageAcceptsSharedMovementPolicies(t *testing.T) {
	policies := []string{
		MovementPolicyImmediate,
		MovementPolicyLinear,
		MovementPolicyEaseOut,
	}

	for _, policy := range policies {
		t.Run(policy, func(t *testing.T) {
			root := makeGaugeFixtures(t)
			packageDir := filepath.Join(root, "assets", "gauges", "radial", policy)
			writeGaugeYAML(t, packageDir, `id: policy_radial
type: radial
sensor: rpm
realism:
  movement_policy: `+policy+`
size:
  width: 100
  height: 100
layers:
  needle: ../../shared/radial/simple_rpm/needle.png
pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 100
  start_angle: -90
  end_angle: 90
`)

			pkg, err := LoadPackage(packageDir)
			if err != nil {
				t.Fatalf("LoadPackage returned error: %v", err)
			}
			if pkg.Realism.MovementPolicy != policy {
				t.Fatalf("movement policy = %q, want %q", pkg.Realism.MovementPolicy, policy)
			}
		})
	}
}

func TestLoadPackageRejectsInvalidSharedMovementPolicy(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "bad_policy")
	writeGaugeYAML(t, packageDir, `id: bad_radial
type: radial
sensor: rpm
realism:
  movement_policy: elastic
size:
  width: 100
  height: 100
layers:
  needle: ../../shared/radial/simple_rpm/needle.png
pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 100
  start_angle: -90
  end_angle: 90
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "movement_policy")
}

func TestLoadPackageRejectsMisspelledSharedMovementPolicyKey(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "bad_policy_key")
	writeGaugeYAML(t, packageDir, `id: bad_radial
type: radial
sensor: rpm
realism:
  movement_polciy: linear
size:
  width: 100
  height: 100
layers:
  needle: ../../shared/radial/simple_rpm/needle.png
pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 100
  start_angle: -90
  end_angle: 90
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "movement_polciy")
}

func TestLoadPackageRejectsExplicitEmptyDrumSlopOnNonOdometerGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "empty_drum_slop")
	writeGaugeYAML(t, packageDir, `id: bad_radial
type: radial
sensor: rpm
realism:
  drum_slop: []
size:
  width: 100
  height: 100
layers:
  needle: ../../shared/radial/simple_rpm/needle.png
pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 100
  start_angle: -90
  end_angle: 90
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "only supported for odometer")
}

func TestLoadPackageRejectsCarryDragOnNonOdometerGauge(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "radial", "carry_drag")
	writeGaugeYAML(t, packageDir, `id: bad_radial
type: radial
sensor: rpm
realism:
  carry_drag: true
size:
  width: 100
  height: 100
layers:
  needle: ../../shared/radial/simple_rpm/needle.png
pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 100
  start_angle: -90
  end_angle: 90
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "carry_drag")
}

func TestLoadPackageRejectsExplicitEmptyOdometerDrumSlop(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "odometer", "empty_drum_slop")
	writeGaugeYAML(t, packageDir, `id: trip_odometer
type: odometer
sensor: trip_distance
realism:
  drum_slop: []
size:
  width: 240
  height: 80
odometer:
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
    - strip: ../trip/red_digits.png
      position: [40, 12]
      window: { width: 24, height: 36 }
      role: sub_unit
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "exactly one offset per odometer wheel")
}

func TestLoadPackageRejectsInvalidOdometerDrumSlop(t *testing.T) {
	tests := []struct {
		name        string
		realismYAML string
		want        string
	}{
		{
			name: "wrong_wheel_count",
			realismYAML: `realism:
  drum_slop: [1]
`,
			want: "exactly one offset per odometer wheel",
		},
		{
			name: "too_large",
			realismYAML: `realism:
  drum_slop: [10, 0]
`,
			want: "exceeds",
		},
		{
			name: "wrong_type",
			realismYAML: `realism:
  drum_slop: [1]
`,
			want: "only supported for odometer",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := makeGaugeFixtures(t)
			packageDir := filepath.Join(root, "assets", "gauges", "odometer", test.name)
			yamlText := `id: trip_odometer
type: odometer
sensor: trip_distance
size:
  width: 240
  height: 80
odometer:
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
    - strip: ../trip/red_digits.png
      position: [40, 12]
      window: { width: 24, height: 36 }
      role: sub_unit
`
			if test.name == "wrong_type" {
				packageDir = filepath.Join(root, "assets", "gauges", "radial", test.name)
				yamlText = `id: bad_radial
type: radial
sensor: rpm
` + test.realismYAML + `size:
  width: 100
  height: 100
layers:
  needle: ../../shared/radial/simple_rpm/needle.png
pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.9 }
value_map:
  min: 0
  max: 100
  start_angle: -90
  end_angle: 90
`
			} else {
				yamlText = `id: trip_odometer
type: odometer
sensor: trip_distance
` + test.realismYAML + `size:
  width: 240
  height: 80
odometer:
  wheels:
    - strip: ../trip/digits.png
      position: [10, 12]
      window: { width: 24, height: 36 }
    - strip: ../trip/red_digits.png
      position: [40, 12]
      window: { width: 24, height: 36 }
      role: sub_unit
`
			}
			writeGaugeYAML(t, packageDir, yamlText)

			_, err := LoadPackage(packageDir)
			if err == nil {
				t.Fatal("LoadPackage returned nil error, want error")
			}
			assertErrorContains(t, err, test.want)
		})
	}
}

func TestLoadPackageRejectsMissingBarLevelLayer(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "bar", "bad")
	writeGaugeYAML(t, packageDir, `id: bad_bar
type: bar
sensor: coolant_temperature
size:
  width: 100
  height: 100
layers:
  panel: panel.png
value_map:
  min: 40
  max: 120
  clamp: true
bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [10, 10, 20, 60]
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "bar layer level")
}

func TestLoadPackageRejectsInvalidBarValueMap(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "bar", "bad_value_map")
	writeGaugeYAML(t, packageDir, `id: bad_bar
type: bar
sensor: coolant_temperature
size:
  width: 100
  height: 100
layers:
  level: level.png
value_map:
  min: 120
  max: 40
  clamp: true
bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [10, 10, 20, 60]
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "bar value_map max must be greater than min")
}

func TestLoadPackageRejectsSegmentedLayerWithoutPercentPlaceholder(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "segmented", "bad")
	writeGaugeYAML(t, packageDir, `id: bad_segmented
type: segmented
sensor: rpm
size:
  width: 100
  height: 100
layers:
  segments: levels/rpm.png
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "{percent}")
}

func TestLoadPackageRejectsMissingBarValueMap(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "bar", "missing_value_map")
	writeGaugeYAML(t, packageDir, `id: bad_bar
type: bar
sensor: coolant_temperature
size:
  width: 100
  height: 100
layers:
  level: level.png
bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [10, 10, 20, 60]
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "bar value_map max must be greater than min")
}

func TestLoadPackageRejectsInvalidBarModeAxisOriginAndBounds(t *testing.T) {
	tests := []struct {
		name    string
		barYAML string
		want    string
	}{
		{
			name: "mode",
			barYAML: `bar:
  mode: fill
  axis: vertical
  origin: bottom
  bounds: [10, 10, 20, 60]
`,
			want: `bar mode`,
		},
		{
			name: "axis",
			barYAML: `bar:
  mode: level
  axis: horizontal
  origin: bottom
  bounds: [10, 10, 20, 60]
`,
			want: `bar axis`,
		},
		{
			name: "origin",
			barYAML: `bar:
  mode: level
  axis: vertical
  origin: top
  bounds: [10, 10, 20, 60]
`,
			want: `bar origin`,
		},
		{
			name: "bounds_length",
			barYAML: `bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [10, 10, 20]
`,
			want: `bar bounds must contain x, y, width, and height`,
		},
		{
			name: "bounds_values",
			barYAML: `bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [-1, 10, 20, 60]
`,
			want: `bar bounds x and y must be non-negative`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := makeGaugeFixtures(t)
			packageDir := filepath.Join(root, "assets", "gauges", "bar", test.name)
			writeGaugeYAML(t, packageDir, `id: bad_bar
type: bar
sensor: coolant_temperature
size:
  width: 100
  height: 100
layers:
  level: level.png
value_map:
  min: 40
  max: 120
  clamp: true
`+test.barYAML)

			_, err := LoadPackage(packageDir)
			if err == nil {
				t.Fatal("LoadPackage returned nil error, want error")
			}
			assertErrorContains(t, err, test.want)
		})
	}
}

func TestLoadPackageRejectsMissingGaugeYAML(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "missing_yaml")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "gauge.yaml")
}

func TestLoadPackageRejectsUnsupportedType(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "bad_type")
	writeGaugeYAML(t, packageDir, `id: bad
type: steam_whistle
sensor: rpm
size:
  width: 100
  height: 100
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "not supported")
}

func TestLoadPackageRejectsPathsEscapingAssetTree(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "gauges", "evil")
	writeGaugeYAML(t, packageDir, `id: escape
type: numeric
sensor: rpm
size:
  width: 100
  height: 100
layers:
  panel: ../../../outside.png
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "escapes asset tree")
}

func TestLoadPackageRejectsPackagesOutsideAssetsDirectory(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "dashboard", "not_a_gauge")
	writeGaugeYAML(t, packageDir, `id: bad
type: radial
sensor: rpm
size:
  width: 100
  height: 100
`)

	_, err := LoadPackage(packageDir)
	if err == nil {
		t.Fatal("LoadPackage returned nil error, want error")
	}
	assertErrorContains(t, err, "assets directory")
}

func makeGaugeFixtures(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"assets/gauges/7Seg/7Seg4Digits.png",
		"assets/gauges/7Seg/Glass.png",
		"assets/gauges/7Seg/7SegBack.png",
		"assets/gauges/7Seg/amber/7Seg0.png",
		"assets/gauges/7Seg/amber/7Seg1.png",
		"assets/gauges/7Seg/amber/7SegDP.png",
		"assets/gauges/bar/coolant/panel.png",
		"assets/gauges/bar/coolant/level.png",
		"assets/gauges/bar/coolant/glass.png",
		"assets/gauges/shared/radial/simple_rpm/bezel.png",
		"assets/gauges/shared/radial/simple_rpm/face.png",
		"assets/gauges/shared/radial/simple_rpm/ticks.png",
		"assets/gauges/shared/radial/simple_rpm/needle.png",
		"assets/gauges/shared/radial/simple_rpm/glass.png",
		"assets/gauges/odometer/trip/panel.png",
		"assets/gauges/odometer/trip/glass.png",
		"assets/gauges/odometer/trip/digits.png",
		"assets/gauges/odometer/trip/red_digits.png",
		"assets/gauges/odometer/bad/digits.png",
		"assets/gauges/indicator/check_engine/bezel.png",
		"assets/gauges/indicator/check_engine/face.png",
		"assets/gauges/indicator/check_engine/off.png",
		"assets/gauges/indicator/check_engine/on.png",
		"assets/gauges/indicator/check_engine/glass.png",
		"assets/gauges/indicator/bad/off.png",
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
	return root
}

func writeGaugeYAML(t *testing.T, packageDir string, text string) {
	t.Helper()
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "gauge.yaml"), []byte(text), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func assertPath(t *testing.T, got string, want string) {
	t.Helper()
	got = filepath.Clean(got)
	want = filepath.Clean(want)
	if got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func assertErrorContains(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("error is nil, want %q", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error = %q, want substring %q", err.Error(), want)
	}
}

func capturePackageLogs(t *testing.T, fn func()) string {
	t.Helper()
	var buffer bytes.Buffer
	previousWriter := log.Writer()
	previousFlags := log.Flags()
	log.SetOutput(&buffer)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(previousWriter)
		log.SetFlags(previousFlags)
	}()
	fn()
	return buffer.String()
}
