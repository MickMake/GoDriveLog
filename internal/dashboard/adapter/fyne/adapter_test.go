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
	fynetest "fyne.io/fyne/v2/test"

	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func TestMain(m *testing.M) {
	app := fynetest.NewApp()
	code := m.Run()
	app.Quit()
	os.Exit(code)
}

func TestAdapterRendersScenePartsFromRepoRelativeAssets(t *testing.T) {
	dir := t.TempDir()
	if err := writeTestPNG(filepath.Join(dir, "assets", "test.png")); err != nil {
		t.Fatal(err)
	}

	adapter, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = adapter.Update([]v3dashboard.Scene{{
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

	adapter, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = adapter.Update([]v3dashboard.Scene{{
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
	adapter, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	err = adapter.Update([]v3dashboard.Scene{{
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
	adapter, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

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
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 255, B: 255, A: 255})
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
