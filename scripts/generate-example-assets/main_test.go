package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFrameworkSmokeGeneratedDigitAssetsShareCellDimensions(t *testing.T) {
	root := t.TempDir()
	if err := generateFrameworkSmoke(root); err != nil {
		t.Fatalf("generateFrameworkSmoke: %v", err)
	}

	digitsRoot := filepath.Join(root, "examples", frameworkSmokeTheme, "assets", "digits")
	want := readPNGSize(t, filepath.Join(digitsRoot, "digit_back.png"))

	for _, name := range []string{
		"digit_back.png",
		"digit_glass.png",
		"digit_0.png",
		"digit_1.png",
		"digit_2.png",
		"digit_3.png",
		"digit_4.png",
		"digit_5.png",
		"digit_6.png",
		"digit_7.png",
		"digit_8.png",
		"digit_9.png",
		"digit_minus.png",
		"digit_dp.png",
	} {
		got := readPNGSize(t, filepath.Join(digitsRoot, name))
		if got != want {
			t.Fatalf("%s size = %v, want %v for the framework-smoke digit cell", name, got, want)
		}
	}
}

func TestOrnateTimberGeneratedDigitAssetsShareCellDimensions(t *testing.T) {
	root := t.TempDir()
	if err := generateOrnateTimber(root); err != nil {
		t.Fatalf("generateOrnateTimber: %v", err)
	}

	digitsRoot := filepath.Join(root, "examples", ornateTimberTheme, "assets", "gauges", "speed_numeric", "digits")
	want := readPNGSize(t, filepath.Join(digitsRoot, "digit_back.png"))

	for _, name := range []string{
		"digit_back.png",
		"digit_glass.png",
		"digit_0.png",
		"digit_1.png",
		"digit_2.png",
		"digit_3.png",
		"digit_4.png",
		"digit_5.png",
		"digit_6.png",
		"digit_7.png",
		"digit_8.png",
		"digit_9.png",
		"digit_minus.png",
		"digit_dp.png",
	} {
		got := readPNGSize(t, filepath.Join(digitsRoot, name))
		if got != want {
			t.Fatalf("%s size = %v, want %v for the ornate-timber digit cell", name, got, want)
		}
	}
}

func TestOrnateTimberGeneratedAssetsIncludeExpectedPaths(t *testing.T) {
	root := t.TempDir()
	if err := generateOrnateTimber(root); err != nil {
		t.Fatalf("generateOrnateTimber: %v", err)
	}

	themeRoot := filepath.Join(root, "examples", ornateTimberTheme, "assets")
	runtimeRoot := filepath.Join(themeRoot, "gauges")
	for _, relative := range []string{
		"panel/background.png",
		"panel/foreground.png",
		"gauges/speed_numeric/panel.png",
		"gauges/speed_numeric/digits/digit_8.png",
		"gauges/radial_rpm/needle.png",
		"gauges/trip_odometer/digits.png",
		"gauges/trip_odometer/tenths.png",
		"gauges/check_engine_indicator/on.png",
		"gauges/fuel_bar/level.png",
		"gauges/rpm_segmented/levels/rpm_100.png",
	} {
		path := filepath.Join(themeRoot, relative)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	for _, relative := range []string{
		"check_engine_indicator/gauge.yaml",
		"fuel_bar/gauge.yaml",
		"radial_rpm/gauge.yaml",
		"rpm_segmented/gauge.yaml",
		"speed_numeric/panel.png",
		"speed_numeric/gauge.yaml",
		"speed_numeric/digits/digit_8.png",
		"radial_rpm/needle.png",
		"trip_odometer/gauge.yaml",
		"trip_odometer/tenths.png",
		"check_engine_indicator/on.png",
		"fuel_bar/level.png",
		"rpm_segmented/levels/rpm_100.png",
	} {
		path := filepath.Join(runtimeRoot, relative)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	segmentedGaugePath := filepath.Join(themeRoot, "gauges", "rpm_segmented", "gauge.yaml")
	segmentedGauge, err := os.ReadFile(segmentedGaugePath)
	if err != nil {
		t.Fatalf("read %s: %v", segmentedGaugePath, err)
	}
	if !strings.Contains(string(segmentedGauge), "sensor: rpm") {
		t.Fatalf("%s does not contain sensor: rpm", segmentedGaugePath)
	}

	for _, relative := range []string{
		"examples/dashboards/framework-smoke.yaml",
		"examples/dashboards/ornate-timber.yaml",
		"examples/assets/v3.4/framework-smoke",
		"examples/assets/v3.4/ornate-timber",
		"assets/gauges/v3.4/ornate-timber",
	} {
		path := filepath.Join(root, relative)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed, got err=%v", path, err)
		}
	}
}

func TestCommittedExampleDashboardsUseSelfContainedLayout(t *testing.T) {
	repoRoot := testRepoRoot(t)

	for _, relative := range []string{
		"examples/framework-smoke/dashboard.yaml",
		"examples/framework-smoke/assets/panel/background.png",
		"examples/framework-smoke/assets/panel/foreground.png",
		"examples/framework-smoke/assets/digits/digit_back.png",
		"examples/framework-smoke/assets/indicator/lamp_on.png",
		"examples/ornate-timber/dashboard.yaml",
		"examples/ornate-timber/assets/panel/background.png",
		"examples/ornate-timber/assets/panel/foreground.png",
		"examples/ornate-timber/assets/gauges/speed_numeric/gauge.yaml",
		"examples/ornate-timber/assets/gauges/radial_rpm/gauge.yaml",
		"examples/ornate-timber/assets/gauges/trip_odometer/gauge.yaml",
		"examples/ornate-timber/assets/gauges/check_engine_indicator/gauge.yaml",
		"examples/ornate-timber/assets/gauges/fuel_bar/gauge.yaml",
		"examples/ornate-timber/assets/gauges/rpm_segmented/gauge.yaml",
	} {
		path := filepath.Join(repoRoot, relative)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func TestAmberSevenSegReferenceAssetKeepsSourceDimensions(t *testing.T) {
	repoRoot := testRepoRoot(t)
	path := filepath.Join(repoRoot, "examples", "assets", "gauges", "7Seg", "amber", "7Seg8.png")
	got := readPNGSize(t, path)
	want := image.Pt(230, 318)
	if got != want {
		t.Fatalf("%s size = %v, want %v", path, got, want)
	}
}

func testRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func readPNGSize(t *testing.T, path string) image.Point {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()

	config, err := png.DecodeConfig(file)
	if err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return image.Pt(config.Width, config.Height)
}
