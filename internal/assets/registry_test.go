package assets

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
)

func TestDefaultSearchPathsUsesExpectedOrder(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "config")
	pwdDir := filepath.Join(root, "work")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(pwdDir, 0o755); err != nil {
		t.Fatal(err)
	}

	oldPwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldPwd); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
	if err := os.Chdir(pwdDir); err != nil {
		t.Fatal(err)
	}
	currentPwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	configRoot, err := filepath.Abs(configDir)
	if err != nil {
		t.Fatal(err)
	}

	got, err := DefaultSearchPaths(filepath.Join(configDir, "config.v3.yaml"), "test_vehicle")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(configRoot, "test_vehicle"),
		filepath.Join(currentPwd, "test_vehicle"),
		configRoot,
		currentPwd,
		"/etc/godrivelog",
		"/usr/local/etc/godrivelog",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DefaultSearchPaths order = %#v, want %#v", got, want)
	}
}

func TestLoadWithSearchPathsUsesFirstMatchingAsset(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	writePNGSize(t, first, "assets/panel.png", 2, 2)
	writePNGSize(t, second, "assets/panel.png", 5, 5)

	registry, err := LoadWithSearchPaths(v3config.AssetConfig{
		ImageSets: map[string]v3config.ImageSetConfig{
			"panel": {Image: "assets/panel.png"},
		},
	}, []string{first, second})
	if err != nil {
		t.Fatalf("LoadWithSearchPaths failed: %v", err)
	}
	panel, ok := registry.ImageSet("panel")
	if !ok || panel.Image == nil {
		t.Fatalf("expected panel image asset")
	}
	if got := panel.Image.Bounds.Dx(); got != 2 {
		t.Fatalf("loaded asset width = %d, want first matching asset width 2", got)
	}
	if want := filepath.Join(first, "assets", "panel.png"); panel.Image.Path != want {
		t.Fatalf("asset path = %q, want %q", panel.Image.Path, want)
	}
}

func TestLoadMinimalAssetRegistry(t *testing.T) {
	root := t.TempDir()
	writePNG(t, root, "assets/panel.png")
	writePNG(t, root, "assets/digit_back.png")
	writePNG(t, root, "assets/minus.png")
	writePNG(t, root, "assets/dp.png")
	writePNG(t, root, "assets/off.png")
	writePNG(t, root, "assets/on.png")
	writePNG(t, root, "assets/unknown.png")
	for _, ch := range requiredDigitCharacters {
		writePNG(t, root, "assets/"+ch+".png")
	}

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
					"2": "assets/2.png",
					"3": "assets/3.png",
					"4": "assets/4.png",
					"5": "assets/5.png",
					"6": "assets/6.png",
					"7": "assets/7.png",
					"8": "assets/8.png",
					"9": "assets/9.png",
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
	if got := registry.SearchPaths(); len(got) != 1 || got[0] != filepath.Clean(root) {
		t.Fatalf("unexpected asset search paths: %#v", got)
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
	for _, ch := range requiredDigitCharacters {
		if _, ok := digits.Characters[ch]; !ok {
			t.Fatalf("expected required digit character %q to be loaded", ch)
		}
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

func TestLoadRicherAssetFamilies(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		"assets/bar/back.png",
		"assets/bar/off.png",
		"assets/bar/on.png",
		"assets/bar/warning.png",
		"assets/bar/front.png",
		"assets/frame/back.png",
		"assets/frame/frame_000.png",
		"assets/frame/frame_001.png",
		"assets/frame/frame_002.png",
		"assets/frame/front.png",
	} {
		writePNG(t, root, path)
	}

	registry, err := Load(v3config.AssetConfig{
		BarSets: map[string]v3config.BarSetConfig{
			"temperature_bar": {
				Background: "assets/bar/back.png",
				Cells: map[string]string{
					"off":     "assets/bar/off.png",
					"on":      "assets/bar/on.png",
					"warning": "assets/bar/warning.png",
				},
				Foreground: "assets/bar/front.png",
				Spacing:    2,
			},
		},
		FrameSets: map[string]v3config.FrameSetConfig{
			"throttle_frames": {
				Background: "assets/frame/back.png",
				Frames: v3config.FrameRangeConfig{
					Path:  "assets/frame/frame_%03d.png",
					First: 0,
					Last:  2,
				},
				Foreground: "assets/frame/front.png",
			},
		},
	}, root)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	bar, ok := registry.BarSet("temperature_bar")
	if !ok {
		t.Fatalf("expected bar set")
	}
	if bar.Spacing != 2 || bar.Background == nil || bar.Foreground == nil {
		t.Fatalf("expected bar metadata and optional layers")
	}
	for _, cell := range []string{"off", "on", "warning"} {
		if bar.Cells[cell].Image == nil {
			t.Fatalf("expected decoded bar cell %q", cell)
		}
	}

	frames, ok := registry.FrameSet("throttle_frames")
	if !ok {
		t.Fatalf("expected frame set")
	}
	if frames.First != 0 || frames.Last != 2 || len(frames.Frames) != 3 {
		t.Fatalf("expected decoded frame range 0..2, got %#v", frames)
	}
	if frames.Background == nil || frames.Foreground == nil {
		t.Fatalf("expected frame optional layers")
	}
	for frame := 0; frame <= 2; frame++ {
		if frames.Frames[frame].Image == nil {
			t.Fatalf("expected decoded frame %d", frame)
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

func TestLoadRejectsNonSearchPathRelativeAssetPath(t *testing.T) {
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
			assertContains(t, err.Error(), "search-path relative")
		})
	}
}

func TestLoadReportsMissingRequiredDigitCharacter(t *testing.T) {
	root := t.TempDir()
	for _, ch := range []string{"0", "1", "3", "4", "5", "6", "7", "8", "9"} {
		writePNG(t, root, "assets/"+ch+".png")
	}

	_, err := Load(v3config.AssetConfig{
		DigitSets: map[string]v3config.DigitSetConfig{
			"digits": {
				Characters: map[string]string{
					"0": "assets/0.png",
					"1": "assets/1.png",
					"3": "assets/3.png",
					"4": "assets/4.png",
					"5": "assets/5.png",
					"6": "assets/6.png",
					"7": "assets/7.png",
					"8": "assets/8.png",
					"9": "assets/9.png",
				},
			},
		},
	}, root)
	if err == nil {
		t.Fatalf("expected missing digit character to fail")
	}
	assertContains(t, err.Error(), "assets.digit_sets.digits.characters.2")
}

func TestLoadReportsMissingRequiredBarOffCell(t *testing.T) {
	root := t.TempDir()
	writePNG(t, root, "assets/on.png")

	_, err := Load(v3config.AssetConfig{
		BarSets: map[string]v3config.BarSetConfig{
			"temperature_bar": {Cells: map[string]string{"on": "assets/on.png"}},
		},
	}, root)
	if err == nil {
		t.Fatalf("expected missing off cell to fail")
	}
	assertContains(t, err.Error(), "assets.bar_sets.temperature_bar.cells.off")
}

func TestLoadReportsBarCellDimensionMismatch(t *testing.T) {
	root := t.TempDir()
	writePNGSize(t, root, "assets/off.png", 2, 2)
	writePNGSize(t, root, "assets/on.png", 3, 2)

	_, err := Load(v3config.AssetConfig{
		BarSets: map[string]v3config.BarSetConfig{
			"temperature_bar": {Cells: map[string]string{"off": "assets/off.png", "on": "assets/on.png"}},
		},
	}, root)
	if err == nil {
		t.Fatalf("expected mismatched cell dimensions to fail")
	}
	assertContains(t, err.Error(), "assets.bar_sets.temperature_bar.cells.on")
	assertContains(t, err.Error(), "dimensions")
}

func TestLoadReportsInvalidFrameRange(t *testing.T) {
	_, err := Load(v3config.AssetConfig{
		FrameSets: map[string]v3config.FrameSetConfig{
			"throttle_frames": {Frames: v3config.FrameRangeConfig{Path: "assets/frame_%03d.png", First: 3, Last: 1}},
		},
	}, t.TempDir())
	if err == nil {
		t.Fatalf("expected invalid frame range to fail")
	}
	assertContains(t, err.Error(), "assets.frame_sets.throttle_frames.frames.first")
}

func TestLoadRejectsMultiFrameLiteralPath(t *testing.T) {
	root := t.TempDir()
	writePNG(t, root, "assets/frame.png")

	_, err := Load(v3config.AssetConfig{
		FrameSets: map[string]v3config.FrameSetConfig{
			"throttle_frames": {Frames: v3config.FrameRangeConfig{Path: "assets/frame.png", First: 0, Last: 2}},
		},
	}, root)
	if err == nil {
		t.Fatalf("expected multi-frame literal path to fail")
	}
	assertContains(t, err.Error(), "assets.frame_sets.throttle_frames.frames.path")
	assertContains(t, err.Error(), "printf placeholder")
	assertContains(t, err.Error(), "multi-frame range")
}

func TestLoadAllowsSingleFrameLiteralPath(t *testing.T) {
	root := t.TempDir()
	writePNG(t, root, "assets/frame.png")

	registry, err := Load(v3config.AssetConfig{
		FrameSets: map[string]v3config.FrameSetConfig{
			"throttle_frames": {Frames: v3config.FrameRangeConfig{Path: "assets/frame.png", First: 0, Last: 0}},
		},
	}, root)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	frames, ok := registry.FrameSet("throttle_frames")
	if !ok || len(frames.Frames) != 1 || frames.Frames[0].Image == nil {
		t.Fatalf("expected single literal frame to load, got %#v", frames)
	}
}

func TestLoadReportsFrameDimensionMismatch(t *testing.T) {
	root := t.TempDir()
	writePNGSize(t, root, "assets/frame_000.png", 2, 2)
	writePNGSize(t, root, "assets/frame_001.png", 2, 3)

	_, err := Load(v3config.AssetConfig{
		FrameSets: map[string]v3config.FrameSetConfig{
			"throttle_frames": {Frames: v3config.FrameRangeConfig{Path: "assets/frame_%03d.png", First: 0, Last: 1}},
		},
	}, root)
	if err == nil {
		t.Fatalf("expected mismatched frame dimensions to fail")
	}
	assertContains(t, err.Error(), "assets.frame_sets.throttle_frames.frames.1")
	assertContains(t, err.Error(), "dimensions")
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
	writePNGSize(t, root, repoPath, 2, 2)
}

func writePNGSize(t *testing.T, root, repoPath string, width, height int) {
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
	img := image.NewRGBA(image.Rect(0, 0, width, height))
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
