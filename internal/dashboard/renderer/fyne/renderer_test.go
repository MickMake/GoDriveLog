package fyne

import (
	"os"
	"path/filepath"
	"testing"

	fyneui "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

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

func TestRendererReusesUnchangedSpriteFrameObject(t *testing.T) {
	renderer := New(nil)
	sceneState := scene.Scene{Elements: []scene.Element{spriteFrameElement("rpm", 1, "frame-1.png")}}

	if err := renderer.Update(sceneState); err != nil {
		t.Fatalf("first update: %v", err)
	}
	if len(renderer.root.Objects) != 1 {
		t.Fatalf("root objects = %d, want 1", len(renderer.root.Objects))
	}
	firstObject := renderer.root.Objects[0]

	if err := renderer.Update(sceneState); err != nil {
		t.Fatalf("second update: %v", err)
	}
	if got := renderer.root.Objects[0]; got != firstObject {
		t.Fatalf("unchanged sprite frame object was rebuilt")
	}
}

func TestRendererUpdatesSpriteFrameResourceInPlace(t *testing.T) {
	renderer := New(nil)
	firstScene := scene.Scene{Elements: []scene.Element{spriteFrameElement("rpm", 1, "frame-1.png")}}
	secondScene := scene.Scene{Elements: []scene.Element{spriteFrameElement("rpm", 2, "frame-2.png")}}

	if err := renderer.Update(firstScene); err != nil {
		t.Fatalf("first update: %v", err)
	}
	imageObject, ok := renderer.root.Objects[0].(*canvas.Image)
	if !ok {
		t.Fatalf("root object type = %T, want *canvas.Image", renderer.root.Objects[0])
	}
	if imageObject.Resource == nil || imageObject.Resource.Name() != "frame-1.png" {
		t.Fatalf("first resource = %v, want frame-1.png", resourceName(imageObject.Resource))
	}

	if err := renderer.Update(secondScene); err != nil {
		t.Fatalf("second update: %v", err)
	}
	if got := renderer.root.Objects[0]; got != imageObject {
		t.Fatalf("changed sprite frame rebuilt image object")
	}
	if imageObject.Resource == nil || imageObject.Resource.Name() != "frame-2.png" {
		t.Fatalf("second resource = %v, want frame-2.png", resourceName(imageObject.Resource))
	}
}

func TestRendererReusesGroupAndUpdatesChildSpriteFrame(t *testing.T) {
	renderer := New(nil)
	firstScene := scene.Scene{Elements: []scene.Element{groupElement("panel", spriteFrameElement("rpm", 1, "frame-1.png"))}}
	secondScene := scene.Scene{Elements: []scene.Element{groupElement("panel", spriteFrameElement("rpm", 2, "frame-2.png"))}}

	if err := renderer.Update(firstScene); err != nil {
		t.Fatalf("first update: %v", err)
	}
	group, ok := renderer.root.Objects[0].(*fyneui.Container)
	if !ok {
		t.Fatalf("root object type = %T, want *fyne.Container", renderer.root.Objects[0])
	}
	childImage, ok := group.Objects[0].(*canvas.Image)
	if !ok {
		t.Fatalf("child object type = %T, want *canvas.Image", group.Objects[0])
	}

	if err := renderer.Update(secondScene); err != nil {
		t.Fatalf("second update: %v", err)
	}
	if got := renderer.root.Objects[0]; got != group {
		t.Fatalf("group object was rebuilt")
	}
	if got := group.Objects[0]; got != childImage {
		t.Fatalf("child image object was rebuilt")
	}
	if childImage.Resource == nil || childImage.Resource.Name() != "frame-2.png" {
		t.Fatalf("child resource = %v, want frame-2.png", resourceName(childImage.Resource))
	}
}

func TestRendererDoesNotReuseCanvasObjectForRepeatedElementID(t *testing.T) {
	renderer := New(nil)
	repeated := spriteFrameElement("rpm", 1, "frame-1.png")
	firstScene := scene.Scene{Elements: []scene.Element{
		groupElement("left", repeated),
		groupElement("right", repeated),
	}}
	secondScene := scene.Scene{Elements: []scene.Element{
		groupElement("left", spriteFrameElement("rpm", 2, "frame-2.png")),
		groupElement("right", spriteFrameElement("rpm", 2, "frame-2.png")),
	}}

	if err := renderer.Update(firstScene); err != nil {
		t.Fatalf("first update: %v", err)
	}
	leftGroup := renderer.root.Objects[0].(*fyneui.Container)
	rightGroup := renderer.root.Objects[1].(*fyneui.Container)
	leftImage := leftGroup.Objects[0].(*canvas.Image)
	rightImage := rightGroup.Objects[0].(*canvas.Image)
	if leftImage == rightImage {
		t.Fatalf("repeated element ID reused one canvas object across occurrences")
	}

	if err := renderer.Update(secondScene); err != nil {
		t.Fatalf("second update: %v", err)
	}
	if leftGroup.Objects[0] != leftImage {
		t.Fatalf("left repeated occurrence image was rebuilt")
	}
	if rightGroup.Objects[0] != rightImage {
		t.Fatalf("right repeated occurrence image was rebuilt")
	}
	if leftImage.Resource == nil || leftImage.Resource.Name() != "frame-2.png" {
		t.Fatalf("left resource = %v, want frame-2.png", resourceName(leftImage.Resource))
	}
	if rightImage.Resource == nil || rightImage.Resource.Name() != "frame-2.png" {
		t.Fatalf("right resource = %v, want frame-2.png", resourceName(rightImage.Resource))
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

func spriteFrameElement(id string, frameIndex int, path string) scene.Element {
	return scene.Element{
		ID:       id,
		Type:     config.DashboardBlockSpriteFrame,
		Visible:  true,
		Geometry: config.RectConfig{X: 1, Y: 2, Width: 30, Height: 40},
		Frame:    assets.Frame{Index: frameIndex, Path: path, Data: []byte{byte(frameIndex)}},
		HasFrame: true,
	}
}

func groupElement(id string, children ...scene.Element) scene.Element {
	return scene.Element{
		ID:       id,
		Type:     config.DashboardBlockGroup,
		Visible:  true,
		Geometry: config.RectConfig{X: 0, Y: 0, Width: 100, Height: 100},
		Children: children,
	}
}

func resourceName(resource fyneui.Resource) string {
	if resource == nil {
		return "<nil>"
	}
	return resource.Name()
}
