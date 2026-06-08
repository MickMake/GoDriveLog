package fyne

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/assets"
	"github.com/MickMake/GoDriveLog/internal/dashboard/scene"
)

func TestUpdateRendersVisibleElementsInSceneOrder(t *testing.T) {
	renderer := New(makeRegistry(t))

	err := renderer.Update(scene.Scene{Elements: []scene.Element{
		{ID: "back", Type: config.DashboardBlockImage, AssetID: "background", Visible: true, Geometry: config.RectConfig{Width: 100, Height: 50}},
		{ID: "hidden", Type: config.DashboardBlockImage, AssetID: "background", Visible: false, Geometry: config.RectConfig{Width: 100, Height: 50}},
		{ID: "front", Type: config.DashboardBlockSpriteFrame, Visible: true, HasFrame: true, Frame: assets.Frame{Index: 1, Data: []byte("frame")}, Geometry: config.RectConfig{Width: 10, Height: 10}},
	}})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if len(renderer.root.Objects) != 2 {
		t.Fatalf("len(root.Objects) = %d, want 2", len(renderer.root.Objects))
	}
}

func TestUpdateRendersSpriteTextGlyphs(t *testing.T) {
	renderer := New(makeRegistry(t))

	err := renderer.Update(scene.Scene{Elements: []scene.Element{
		{
			ID:       "digits",
			Type:     config.DashboardBlockSpriteText,
			Visible:  true,
			Geometry: config.RectConfig{X: 10, Y: 20, Width: 200, Height: 50},
			Glyphs: []assets.Glyph{
				{Char: "1", Data: []byte("glyph-1")},
				{Char: "2", Data: []byte("glyph-2")},
			},
		},
	}})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if len(renderer.root.Objects) != 1 {
		t.Fatalf("len(root.Objects) = %d, want 1", len(renderer.root.Objects))
	}
	if renderer.root.Objects[0].Position().X != 10 || renderer.root.Objects[0].Position().Y != 20 {
		t.Fatalf("sprite text position = %#v, want 10,20", renderer.root.Objects[0].Position())
	}
}

func TestUpdateRendersGroupChildren(t *testing.T) {
	renderer := New(makeRegistry(t))

	err := renderer.Update(scene.Scene{Elements: []scene.Element{
		{
			ID:      "group",
			Type:    config.DashboardBlockGroup,
			Visible: true,
			Children: []scene.Element{
				{ID: "child", Type: config.DashboardBlockImage, AssetID: "background", Visible: true, Geometry: config.RectConfig{Width: 100, Height: 50}},
			},
		},
	}})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if len(renderer.root.Objects) != 1 {
		t.Fatalf("len(root.Objects) = %d, want 1", len(renderer.root.Objects))
	}
}

func makeRegistry(t *testing.T) *assets.Registry {
	t.Helper()
	root := t.TempDir()
	assetPath := filepath.Join(root, "assets", "background.png")
	if err := os.MkdirAll(filepath.Dir(assetPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(assetPath, []byte("image"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	registry, err := assets.Load(config.DashboardConfig{
		AssetRoot: "assets",
		Assets: []config.DashboardAssetConfig{
			{ID: "background", Type: config.DashboardAssetImage, Path: "background.png"},
		},
	}, filepath.Join(root, "dashboard.yaml"))
	if err != nil {
		t.Fatalf("assets.Load returned error: %v", err)
	}
	return registry
}
