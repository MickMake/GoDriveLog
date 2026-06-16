package fyne

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	fyneui "fyne.io/fyne/v2"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

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
		Size:        struct{ Width, Height int }{Width: 24, Height: 16},
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

func writeTestPNG(path string) error {
	const png1x1 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGJgYGAAAAAEAAGjChXjAAAAAElFTkSuQmCC"
	data, err := base64.StdEncoding.DecodeString(png1x1)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
