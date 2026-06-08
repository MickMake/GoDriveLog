package assets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func TestLoadRegistryLoadsAndCachesAssets(t *testing.T) {
	root := makeAssetFixtures(t)
	cfg := config.DashboardConfig{
		AssetRoot: "assets",
		Assets: []config.DashboardAssetConfig{
			{ID: "background", Type: config.DashboardAssetImage, Path: "background.png"},
			{ID: "throttle", Type: config.DashboardAssetFrameSet, Pattern: "frames/frame_{index:03}.png", FrameCount: 3},
			{ID: "digits", Type: config.DashboardAssetCharset, Glyphs: map[string]string{"0": "digits/0.png", "1": "digits/1.png"}},
		},
	}

	registry, err := Load(cfg, filepath.Join(root, "dashboard.yaml"))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	background, err := registry.MustGet("background")
	if err != nil {
		t.Fatalf("MustGet(background): %v", err)
	}
	if background.Type != TypeImage || string(background.Data) != "background" {
		t.Fatalf("background = %#v, want cached image bytes", background)
	}

	throttle, err := registry.MustGet("throttle")
	if err != nil {
		t.Fatalf("MustGet(throttle): %v", err)
	}
	if len(throttle.Frames) != 3 {
		t.Fatalf("len(throttle.Frames) = %d, want 3", len(throttle.Frames))
	}
	if throttle.Frames[2].Index != 2 || string(throttle.Frames[2].Data) != "frame-2" {
		t.Fatalf("throttle frame 2 = %#v, want cached frame-2", throttle.Frames[2])
	}

	digits, err := registry.MustGet("digits")
	if err != nil {
		t.Fatalf("MustGet(digits): %v", err)
	}
	if len(digits.Glyphs) != 2 || string(digits.Glyphs["1"].Data) != "digit-1" {
		t.Fatalf("digits glyphs = %#v, want cached glyphs", digits.Glyphs)
	}

	if err := os.WriteFile(filepath.Join(root, "assets", "background.png"), []byte("changed"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	backgroundAgain, err := registry.MustGet("background")
	if err != nil {
		t.Fatalf("MustGet(background): %v", err)
	}
	if string(backgroundAgain.Data) != "background" {
		t.Fatalf("background cache = %q, want original cached bytes", string(backgroundAgain.Data))
	}
}

func TestLoadRegistrySupportsExplicitFrames(t *testing.T) {
	root := makeAssetFixtures(t)
	cfg := config.DashboardConfig{
		AssetRoot: "assets",
		Assets: []config.DashboardAssetConfig{
			{ID: "throttle", Type: config.DashboardAssetFrameSet, FrameCount: 2, Frames: []string{"frames/frame_000.png", "frames/frame_001.png"}},
		},
	}

	registry, err := Load(cfg, filepath.Join(root, "dashboard.yaml"))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	asset, err := registry.MustGet("throttle")
	if err != nil {
		t.Fatalf("MustGet(throttle): %v", err)
	}
	if len(asset.Frames) != 2 {
		t.Fatalf("len(asset.Frames) = %d, want 2", len(asset.Frames))
	}
}

func TestLoadRegistryRejectsInvalidAssets(t *testing.T) {
	root := makeAssetFixtures(t)
	tests := []struct {
		name  string
		asset config.DashboardAssetConfig
	}{
		{name: "missing image", asset: config.DashboardAssetConfig{ID: "missing", Type: config.DashboardAssetImage, Path: "missing.png"}},
		{name: "remote path", asset: config.DashboardAssetConfig{ID: "remote", Type: config.DashboardAssetImage, Path: "https://example.com/image.png"}},
		{name: "generated frame missing", asset: config.DashboardAssetConfig{ID: "frames", Type: config.DashboardAssetFrameSet, Pattern: "frames/missing_{index:03}.png", FrameCount: 2}},
		{name: "bad generated pattern", asset: config.DashboardAssetConfig{ID: "frames", Type: config.DashboardAssetFrameSet, Pattern: "frames/frame.png", FrameCount: 2}},
		{name: "frame count mismatch", asset: config.DashboardAssetConfig{ID: "frames", Type: config.DashboardAssetFrameSet, FrameCount: 3, Frames: []string{"frames/frame_000.png"}}},
		{name: "missing glyph", asset: config.DashboardAssetConfig{ID: "digits", Type: config.DashboardAssetCharset, Glyphs: map[string]string{"0": "digits/missing.png"}}},
		{name: "empty glyph key", asset: config.DashboardAssetConfig{ID: "digits", Type: config.DashboardAssetCharset, Glyphs: map[string]string{"": "digits/0.png"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(config.DashboardConfig{AssetRoot: "assets", Assets: []config.DashboardAssetConfig{tt.asset}}, filepath.Join(root, "dashboard.yaml"))
			if err == nil {
				t.Fatal("Load returned nil error, want error")
			}
		})
	}
}

func TestMustGetRejectsUnknownAsset(t *testing.T) {
	registry := &Registry{assets: map[string]Asset{}}
	if _, err := registry.MustGet("missing"); err == nil {
		t.Fatal("MustGet returned nil error, want error")
	}
}

func TestExpandFramePath(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		index   int
		want    string
	}{
		{name: "plain index", pattern: "frame_{index}.png", index: 7, want: "frame_7.png"},
		{name: "zero padded index", pattern: "frame_{index:03}.png", index: 7, want: "frame_007.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandFramePath(tt.pattern, tt.index)
			if err != nil {
				t.Fatalf("expandFramePath returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expandFramePath = %q, want %q", got, tt.want)
			}
		})
	}
}

func makeAssetFixtures(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"assets/background.png":       "background",
		"assets/frames/frame_000.png": "frame-0",
		"assets/frames/frame_001.png": "frame-1",
		"assets/frames/frame_002.png": "frame-2",
		"assets/digits/0.png":         "digit-0",
		"assets/digits/1.png":         "digit-1",
	}
	for path, content := range files {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}
	return root
}
