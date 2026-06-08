package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

func TestNewDashboardWithConfigPathResolvesAssetsRelativeToConfigFile(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "configs")
	assetPath := filepath.Join(configDir, "assets", "background.png")
	if err := os.MkdirAll(filepath.Dir(assetPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(assetPath, []byte("image"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	dashboard, err := NewDashboardWithConfigPath(config.DashboardConfig{
		Canvas:    config.CanvasConfig{Width: 320, Height: 240},
		AssetRoot: "assets",
		Assets: []config.DashboardAssetConfig{
			{ID: "background", Type: config.DashboardAssetImage, Path: "background.png"},
		},
		Blocks: []config.DashboardBlockConfig{
			{ID: "background_panel", Type: config.DashboardBlockImage, Asset: "background", Geometry: config.RectConfig{Width: 320, Height: 240}},
		},
		Layers: []config.DashboardLayerConfig{
			{ID: "base", Z: 0, Blocks: []string{"background_panel"}},
		},
	}, filepath.Join(configDir, "dashboard.yaml"), sensors.NewStateStore(nil))
	if err != nil {
		t.Fatalf("NewDashboardWithConfigPath returned error: %v", err)
	}
	if dashboard.LastError() != nil {
		t.Fatalf("dashboard LastError = %v, want nil", dashboard.LastError())
	}
}
