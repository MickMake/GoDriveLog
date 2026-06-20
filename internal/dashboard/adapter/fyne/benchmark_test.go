package fyne

import (
	"fmt"
	"path/filepath"
	"testing"

	fynetest "fyne.io/fyne/v2/test"

	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func BenchmarkSevenSegmentAdapterUpdateNoRefresh(b *testing.B) {
	benchmarkSevenSegmentAdapterUpdate(b, false)
}

func BenchmarkSevenSegmentAdapterUpdateWithFyneRefresh(b *testing.B) {
	app := fynetest.NewApp()
	defer app.Quit()
	benchmarkSevenSegmentAdapterUpdate(b, true)
}

func benchmarkSevenSegmentAdapterUpdate(b *testing.B, withRefresh bool) {
	dir := b.TempDir()
	for _, asset := range []string{"panel.png", "glass.png"} {
		if err := writeTestPNG(filepath.Join(dir, "assets", asset)); err != nil {
			b.Fatal(err)
		}
	}
	for digit := 0; digit <= 9; digit++ {
		if err := writeTestPNG(filepath.Join(dir, "assets", fmt.Sprintf("digit%d.png", digit))); err != nil {
			b.Fatal(err)
		}
	}
	adapter, err := New(dir)
	if err != nil {
		b.Fatal(err)
	}
	if !withRefresh {
		disableRefresh(adapter)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		digitAsset := fmt.Sprintf("assets/digit%d.png", i%10)
		if err := adapter.Update([]v3dashboard.Scene{gaugeSceneWithDigit(digitAsset)}); err != nil {
			b.Fatal(err)
		}
	}
}
