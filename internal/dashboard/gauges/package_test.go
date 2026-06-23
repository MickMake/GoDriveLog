package gauges

import (
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
	if pkg.Odometer.Movement != MovementSmooth {
		t.Fatalf("movement = %q, want default smooth", pkg.Odometer.Movement)
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

func TestLoadPackageRejectsPackagesOutsideAssetsGauges(t *testing.T) {
	root := makeGaugeFixtures(t)
	packageDir := filepath.Join(root, "assets", "dashboard", "not_a_gauge")
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
	assertErrorContains(t, err, "assets/gauges")
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
