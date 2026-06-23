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
