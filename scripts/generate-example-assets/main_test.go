package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFrameworkSmokeGeneratedDigitAssetsShareCellDimensions(t *testing.T) {
	root := t.TempDir()
	if err := generateFrameworkSmoke(root); err != nil {
		t.Fatalf("generateFrameworkSmoke: %v", err)
	}

	digitsRoot := filepath.Join(root, "examples", "assets", "v3.4", frameworkSmokeTheme, "digits")
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

	digitsRoot := filepath.Join(root, "examples", "assets", "v3.4", ornateTimberTheme, "gauges", "speed_numeric", "digits")
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

	themeRoot := filepath.Join(root, "examples", "assets", "v3.4", ornateTimberTheme)
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

	runtimeRoot := filepath.Join(root, "assets", "gauges", "v3.4", ornateTimberTheme)
	for _, relative := range []string{
		"speed_numeric/panel.png",
		"speed_numeric/digits/digit_8.png",
		"radial_rpm/needle.png",
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
