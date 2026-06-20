package fyne

import (
	"path/filepath"
	"testing"

	fyneui "fyne.io/fyne/v2"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func TestAdapterPositionsGaugePartsFromPackageCoordinatesAndScale(t *testing.T) {
	dir := t.TempDir()
	if err := writeTestPNG(filepath.Join(dir, "assets", "digit.png")); err != nil {
		t.Fatal(err)
	}
	adapter, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	parts, err := adapter.renderWidgetParts("primary", v3dashboard.Widget{
		ID:       "rpm",
		Type:     v3dashboard.PartKindLayer,
		Position: []int{10, 20},
		Scale:    2,
		Parts: []v3dashboard.Part{{
			Kind:      v3dashboard.PartKindCharacter,
			AssetPath: "assets/digit.png",
			Position:  []int{3, 4},
		}},
	}, 5)
	if err != nil {
		t.Fatalf("renderWidgetParts returned error: %v", err)
	}
	if len(parts) != 1 {
		t.Fatalf("expected one rendered part, got %d", len(parts))
	}
	part := parts[0]
	if part.x != 16 || part.y != 33 {
		t.Fatalf("rendered position = (%v,%v), want (16,33)", part.x, part.y)
	}
	if part.size != fyneui.NewSize(2, 2) {
		t.Fatalf("rendered size = %v, want 2x2", part.size)
	}
}
