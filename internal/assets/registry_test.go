package assets

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
)

func TestLoadMinimalAssetRegistry(t *testing.T) {
	root := t.TempDir()
	writePNG(t, root, "assets/panel.png")
	writePNG(t, root, "assets/digit_back.png")
	writePNG(t, root, "assets/0.png")
	writePNG(t, root, "assets/1.png")
	writePNG(t, root, "assets/minus.png")
	writePNG(t, root, "assets/dp.png")
	writePNG(t, root, "assets/off.png")
	writePNG(t, root, "assets/on.png")
	writePNG(t, root, "assets/unknown.png")

	registry, err := Load(v3config.AssetConfig{
		ImageSets: map[string]v3config.ImageSetConfig{
			"panel": {Image: "assets/panel.png"},
		},
		DigitSets: map[string]v3config.DigitSetConfig{
			"digits": {
				Background: "assets/digit_back.png",
				Characters: map[string]string{
					"0": "assets/0.png",
					"1": "assets/1.png",
					"-": "assets/minus.png",
				},
				DecimalPoint: "assets/dp.png",
				Spacing:      4,
			},
		},
		IndicatorSets: map[string]v3config.IndicatorSetConfig{
			"warning": {
				States: map[string]string{
					IndicatorStateOff:     "assets/off.png",
					IndicatorStateOn:      "assets/on.png",
					IndicatorStateUnknown: "assets/unknown.png",
				},
			},
		},
	}, root)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if registry.RepoRoot() != filepath.Clean(root) {
		t.Fatalf("unexpected repo root: %q", registry.RepoRoot())
	}
	panel, ok := registry.ImageSet("panel")
	if !ok || panel.Image == nil || panel.Image.Image == nil {
		t.Fatalf("expected decoded panel image asset")
	}
	digits, ok := registry.DigitSet("digits")
	if !ok {
		t.Fatalf("expected digit set")
	}
	if digits.Spacing != 4 || digits.Background == nil || digits.DecimalPoint == nil {
		t.Fatalf("expected digit metadata and optional layers")
	}
	if _, ok := digits.Characters["-"]; !ok {
		t.Fatalf("expected minus character to be loaded")
	}
	indicator, ok := registry.IndicatorSet("warning")
	if !ok {
		t.Fatalf("expected indicator set")
	}
	for _, state := range []string{IndicatorStateOff, IndicatorStateOn, IndicatorStateUnknown} {
		if indicator.States[state].Image == nil {
			t.Fatalf("expected decoded indicator state %q", state)
		}
	}
}

func TestLoadReportsMissingAssetClearly(t *testing.T) {
	_, err := Load(v3config.AssetConfig{
		ImageSets: map[string]v3config.ImageSetConfig{
			"panel": {Image: "assets/missing.png"},
		},
	}, t.TempDir())
	if err == nil {
		t.Fatalf("expected missing asset to fail")
	}
	assertContains(t, err.Error(), "assets.image_sets.panel.image")
	assertContains(t, err.Error(), "missing.png")
}

func TestLoadRejectsNonRootRelativeAssetPath(t *testing.T) {
	for _, badPath := range []string{"../escape.png", "/tmp/escape.png", "https://example.invalid/image.png"} {
		t.Run(badPath, func(t *testing.T) {
			_, err := Load(v3config.AssetConfig{
				ImageSets: map[string]v3config.ImageSetConfig{
					"panel": {Image: badPath},
				},
			}, t.TempDir())
			if err == nil {
				t.Fatalf("expected bad path to fail")
			}
			assertContains(t, err.Error(), "repository-root relative")
		})
	}
}

func TestLoadReportsMissingRequiredIndicatorState(t *testing.T) {
	root := t.TempDir()
	writePNG(t, root, "assets/off.png")
	writePNG(t, root, "assets/on.png")

	_, err := Load(v3config.AssetConfig{
		IndicatorSets: map[string]v3config.IndicatorSetConfig{
			"warning": {
				States: map[string]string{
					IndicatorStateOff: "assets/off.png",
					IndicatorStateOn:  "assets/on.png",
				},
			},
		},
	}, root)
	if err == nil {
		t.Fatalf("expected missing unknown state to fail")
	}
	assertContains(t, err.Error(), "assets.indicator_sets.warning.states.unknown")
}

func writePNG(t *testing.T, root, repoPath string) {
	t.Helper()
	fullPath := filepath.Join(root, filepath.FromSlash(repoPath))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	file, err := os.Create(fullPath)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer file.Close()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(file, img); err != nil {
		t.Fatalf("png.Encode failed: %v", err)
	}
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("expected %q to contain %q", got, want)
	}
}
