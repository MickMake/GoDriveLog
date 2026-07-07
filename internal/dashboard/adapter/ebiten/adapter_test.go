package ebiten

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	v3gauges "github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func TestRenderWidgetPartsTreatsNeedleMinAsRotatingNeedlePart(t *testing.T) {
	rendered := renderSingleNeedleLikePart(t, v3dashboard.PartKindNeedleMin)

	if !rendered.needle {
		t.Fatal("expected needle_min to use rotating needle path")
	}
	if rendered.shadow {
		t.Fatal("expected needle_min not to be treated as a shadow part")
	}
	assertRenderedNeedleGeometry(t, rendered)
}

func TestRenderWidgetPartsTreatsNeedleMaxAsRotatingNeedlePart(t *testing.T) {
	rendered := renderSingleNeedleLikePart(t, v3dashboard.PartKindNeedleMax)

	if !rendered.needle {
		t.Fatal("expected needle_max to use rotating needle path")
	}
	if rendered.shadow {
		t.Fatal("expected needle_max not to be treated as a shadow part")
	}
	assertRenderedNeedleGeometry(t, rendered)
}

func TestRenderWidgetPartsKeepsLiveNeedleAndShadowBehavior(t *testing.T) {
	needle := renderSingleNeedleLikePart(t, v3dashboard.PartKindNeedle)
	if !needle.needle {
		t.Fatal("expected live needle to use rotating needle path")
	}
	if needle.shadow {
		t.Fatal("expected live needle not to be treated as a shadow part")
	}
	assertRenderedNeedleGeometry(t, needle)

	shadow := renderSingleNeedleLikePart(t, v3dashboard.PartKindNeedleShadow)
	if !shadow.needle {
		t.Fatal("expected needle shadow to use rotating needle path")
	}
	if !shadow.shadow {
		t.Fatal("expected needle shadow to stay marked as a shadow part")
	}
	assertRenderedNeedleGeometry(t, shadow)
}

func renderSingleNeedleLikePart(t *testing.T, kind string) renderedPart {
	t.Helper()

	root := t.TempDir()
	writeAdapterPNG(t, filepath.Join(root, "face.png"), 100, 80)
	writeAdapterPNG(t, filepath.Join(root, "needle.png"), 20, 40)

	adapter, err := New(root, 0, 0)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	widget := v3dashboard.Widget{
		ID:       "rpm",
		Position: []int{10, 20},
		Scale:    2,
		Parts: []v3dashboard.Part{
			{Kind: v3dashboard.PartKindLayer, Layer: "face", AssetPath: "face.png"},
			{
				Kind:        kind,
				Layer:       kind,
				AssetPath:   "needle.png",
				Angle:       27.5,
				FacePivot:   v3gauges.Point{X: 0.25, Y: 0.75},
				NeedlePivot: v3gauges.Point{X: 0.4, Y: 0.6},
				Position:    []int{3, 4},
				Alpha:       0.35,
			},
		},
	}

	parts, err := adapter.renderWidgetParts("primary", widget, 0)
	if err != nil {
		t.Fatalf("renderWidgetParts returned error: %v", err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 rendered parts, got %d", len(parts))
	}
	return parts[1]
}

func assertRenderedNeedleGeometry(t *testing.T, rendered renderedPart) {
	t.Helper()

	if rendered.angle != 27.5 {
		t.Fatalf("angle = %v, want 27.5", rendered.angle)
	}
	if rendered.x != 66 {
		t.Fatalf("x = %v, want 66", rendered.x)
	}
	if rendered.y != 148 {
		t.Fatalf("y = %v, want 148", rendered.y)
	}
	if rendered.pivotX != 8 {
		t.Fatalf("pivotX = %v, want 8", rendered.pivotX)
	}
	if rendered.pivotY != 24 {
		t.Fatalf("pivotY = %v, want 24", rendered.pivotY)
	}
	if rendered.scale != 2 {
		t.Fatalf("scale = %v, want 2", rendered.scale)
	}
}

func writeAdapterPNG(t *testing.T, path string, width int, height int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	defer file.Close()

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	if err := png.Encode(file, img); err != nil {
		t.Fatalf("Encode: %v", err)
	}
}
