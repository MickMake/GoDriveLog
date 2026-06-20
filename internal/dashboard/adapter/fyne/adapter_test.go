package fyne

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	fyneui "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func TestAdapterRendersScenePartsFromRepoRelativeAssets(t *testing.T) {
	dir := t.TempDir()
	if err := writeTestPNG(filepath.Join(dir, "assets", "test.png")); err != nil {
		t.Fatal(err)
	}

	adapter := newNoRefreshAdapter(t, dir)

	err := adapter.Update([]v3dashboard.Scene{{
		DashboardID: "primary",
		Size:        v3config.SizeConfig{Width: 24, Height: 16},
		Widgets: []v3dashboard.Widget{{
			ID:       "background",
			Type:     "image",
			Position: []int{3, 4},
			Parts: []v3dashboard.Part{{
				Kind:      v3dashboard.PartKindImage,
				AssetPath: "assets/test.png",
			}},
		}},
	}})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if adapter.RenderedObjectCount() != 1 {
		t.Fatalf("RenderedObjectCount = %d, want 1", adapter.RenderedObjectCount())
	}
	if adapter.CanvasObject().Size() != fyneui.NewSize(24, 16) {
		t.Fatalf("CanvasObject size = %v, want 24x16", adapter.CanvasObject().Size())
	}
}

func TestAdapterRendersScenePartsFromResolvedAssetPaths(t *testing.T) {
	dir := t.TempDir()
	assetPath := filepath.Join(dir, "assets", "resolved.png")
	if err := writeTestPNG(assetPath); err != nil {
		t.Fatal(err)
	}

	adapter := newNoRefreshAdapter(t, dir)

	err := adapter.Update([]v3dashboard.Scene{{
		DashboardID: "primary",
		Size:        v3config.SizeConfig{Width: 24, Height: 16},
		Widgets: []v3dashboard.Widget{{
			ID:       "resolved",
			Type:     "image",
			Position: []int{0, 0},
			Parts: []v3dashboard.Part{{
				Kind:      v3dashboard.PartKindImage,
				AssetPath: assetPath,
			}},
		}},
	}})
	if err != nil {
		t.Fatalf("Update returned error for resolved asset path: %v", err)
	}
	if adapter.RenderedObjectCount() != 1 {
		t.Fatalf("RenderedObjectCount = %d, want 1", adapter.RenderedObjectCount())
	}
}

func TestAdapterRejectsEscapingAssetPath(t *testing.T) {
	adapter := newNoRefreshAdapter(t, t.TempDir())

	err := adapter.Update([]v3dashboard.Scene{{
		DashboardID: "primary",
		Widgets: []v3dashboard.Widget{{
			ID: "bad",
			Parts: []v3dashboard.Part{{
				Kind:      v3dashboard.PartKindImage,
				AssetPath: "../outside.png",
			}},
		}},
	}})
	if err == nil {
		t.Fatal("Update succeeded for escaping asset path")
	}
}

func TestAdapterKeepsGlassOverlayLastAndReusesGaugeObjects(t *testing.T) {
	dir := t.TempDir()
	for _, asset := range []string{"panel.png", "digit0.png", "digit1.png", "glass.png"} {
		if err := writeTestPNG(filepath.Join(dir, "assets", asset)); err != nil {
			t.Fatal(err)
		}
	}
	adapter := newNoRefreshAdapter(t, dir)

	if err := adapter.Update([]v3dashboard.Scene{gaugeSceneWithDigit("assets/digit0.png")}); err != nil {
		t.Fatalf("first Update returned error: %v", err)
	}
	if got := adapter.RenderedObjectCount(); got != 3 {
		t.Fatalf("RenderedObjectCount = %d, want 3", got)
	}
	firstObjects := append([]fyneui.CanvasObject(nil), adapter.root.Objects...)
	assertLastResourceName(t, adapter, "assets/glass.png")

	if err := adapter.Update([]v3dashboard.Scene{gaugeSceneWithDigit("assets/digit1.png")}); err != nil {
		t.Fatalf("second Update returned error: %v", err)
	}
	if got := adapter.RenderedObjectCount(); got != 3 {
		t.Fatalf("RenderedObjectCount after digit change = %d, want 3", got)
	}
	for index, object := range adapter.root.Objects {
		if object != firstObjects[index] {
			t.Fatalf("object %d was rebuilt across a digit-only change", index)
		}
	}
	assertLastResourceName(t, adapter, "assets/glass.png")
}

func TestAdapterRendersRadialNeedleAroundPivots(t *testing.T) {
	dir := t.TempDir()
	for _, asset := range []struct {
		path   string
		width  int
		height int
	}{
		{path: "assets/background.png", width: 100, height: 100},
		{path: "assets/face.png", width: 100, height: 100},
		{path: "assets/ticks.png", width: 100, height: 100},
		{path: "assets/needle.png", width: 10, height: 20},
		{path: "assets/overlay.png", width: 100, height: 100},
	} {
		if err := writeSizedTestPNG(filepath.Join(dir, asset.path), asset.width, asset.height); err != nil {
			t.Fatal(err)
		}
	}
	adapter := newNoRefreshAdapter(t, dir)

	parts, err := adapter.renderWidgetParts("primary", radialWidgetWithNeedle(90), 0)
	if err != nil {
		t.Fatalf("renderWidgetParts returned error: %v", err)
	}
	if len(parts) != 5 {
		t.Fatalf("rendered radial part count = %d, want 5", len(parts))
	}
	needle := parts[3]
	if needle.size != fyneui.NewSize(40, 20) {
		t.Fatalf("rotated/scaled needle size = %v, want 40x20", needle.size)
	}
	if needle.x != 110 || needle.y != 110 {
		t.Fatalf("rotated needle position = (%v,%v), want (110,110)", needle.x, needle.y)
	}
}

func TestAdapterReusesRadialNeedleObjectAcrossAngleChanges(t *testing.T) {
	dir := t.TempDir()
	for _, asset := range []struct {
		path   string
		width  int
		height int
	}{
		{path: "assets/face.png", width: 100, height: 100},
		{path: "assets/needle.png", width: 10, height: 20},
		{path: "assets/overlay.png", width: 100, height: 100},
	} {
		if err := writeSizedTestPNG(filepath.Join(dir, asset.path), asset.width, asset.height); err != nil {
			t.Fatal(err)
		}
	}
	adapter := newNoRefreshAdapter(t, dir)

	if err := adapter.Update([]v3dashboard.Scene{radialSceneWithNeedle(0)}); err != nil {
		t.Fatalf("first radial Update returned error: %v", err)
	}
	if got := adapter.RenderedObjectCount(); got != 3 {
		t.Fatalf("RenderedObjectCount = %d, want 3", got)
	}
	firstNeedle := adapter.root.Objects[1]

	if err := adapter.Update([]v3dashboard.Scene{radialSceneWithNeedle(90)}); err != nil {
		t.Fatalf("second radial Update returned error: %v", err)
	}
	if got := adapter.RenderedObjectCount(); got != 3 {
		t.Fatalf("RenderedObjectCount after angle change = %d, want 3", got)
	}
	if adapter.root.Objects[1] != firstNeedle {
		t.Fatalf("radial needle object was rebuilt across an angle-only change")
	}
}

func newNoRefreshAdapter(t *testing.T, repoRoot string) *Adapter {
	t.Helper()
	adapter, err := New(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	disableRefresh(adapter)
	return adapter
}

func disableRefresh(adapter *Adapter) {
	adapter.refreshRoot = func(*fyneui.Container) {}
	adapter.refreshImage = func(*canvas.Image) {}
}

func gaugeSceneWithDigit(digitAsset string) v3dashboard.Scene {
	return v3dashboard.Scene{
		DashboardID: "primary",
		Size:        v3config.SizeConfig{Width: 24, Height: 16},
		Widgets: []v3dashboard.Widget{{
			ID:       "rpm",
			Type:     "gauge",
			Position: []int{3, 4},
			Scale:    1,
			Parts: []v3dashboard.Part{{
				Kind:      v3dashboard.PartKindLayer,
				Layer:     "panel",
				AssetPath: "assets/panel.png",
			}, {
				Kind:      v3dashboard.PartKindCharacter,
				AssetPath: digitAsset,
				Slot:      0,
				Position:  []int{2, 2},
			}, {
				Kind:      v3dashboard.PartKindLayer,
				Layer:     "glass",
				AssetPath: "assets/glass.png",
			}},
		}},
	}
}

func radialSceneWithNeedle(angle float64) v3dashboard.Scene {
	return v3dashboard.Scene{
		DashboardID: "primary",
		Size:        v3config.SizeConfig{Width: 240, Height: 240},
		Widgets:    []v3dashboard.Widget{radialWidgetWithNeedle(angle)},
	}
}

func radialWidgetWithNeedle(angle float64) v3dashboard.Widget {
	return v3dashboard.Widget{
		ID:       "rpm",
		Type:     "gauge",
		Position: []int{10, 20},
		Scale:    2,
		Parts: []v3dashboard.Part{{
			Kind:      v3dashboard.PartKindLayer,
			Layer:     "background",
			AssetPath: "assets/background.png",
		}, {
			Kind:      v3dashboard.PartKindLayer,
			Layer:     "face",
			AssetPath: "assets/face.png",
		}, {
			Kind:      v3dashboard.PartKindLayer,
			Layer:     "ticks",
			AssetPath: "assets/ticks.png",
		}, {
			Kind:        v3dashboard.PartKindNeedle,
			Layer:       "needle",
			AssetPath:   "assets/needle.png",
			Angle:       angle,
			FacePivot:   v3dashboard.GaugePoint{X: 0.5, Y: 0.5},
			NeedlePivot: v3dashboard.GaugePoint{X: 0.5, Y: 1},
		}, {
			Kind:      v3dashboard.PartKindLayer,
			Layer:     "overlay",
			AssetPath: "assets/overlay.png",
		}},
	}
}

func assertLastResourceName(t *testing.T, adapter *Adapter, want string) {
	t.Helper()
	if len(adapter.root.Objects) == 0 {
		t.Fatal("expected at least one rendered object")
	}
	image, ok := adapter.root.Objects[len(adapter.root.Objects)-1].(*canvas.Image)
	if !ok {
		t.Fatalf("last rendered object type = %T, want *canvas.Image", adapter.root.Objects[len(adapter.root.Objects)-1])
	}
	if image.Resource.Name() != want {
		t.Fatalf("last resource = %q, want glass overlay %q", image.Resource.Name(), want)
	}
}

func writeTestPNG(path string) error {
	return writeSizedTestPNG(path, 1, 1)
}

func writeSizedTestPNG(path string, width int, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}
