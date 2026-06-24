package scene

import (
	"testing"

	"github.com/MickMake/GoDriveLog/internal/config"
)

func TestReusableBlockAliasesResolveToScenePrimitives(t *testing.T) {
	registry := makeRegistry(t)
	tests := []struct {
		name           string
		block          config.DashboardBlockConfig
		wantType       string
		wantText       string
		wantFrame      bool
		wantAssetID    string
		wantGlyphs     int
		wantChildCount int
	}{
		{
			name:       "seven segment number",
			block:      config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockSevenSegmentNumber, Asset: "yellow_digits", Decoder: "rpm_digits", Geometry: config.RectConfig{X: 10, Y: 10, Width: 200, Height: 80}},
			wantType:   config.DashboardBlockSpriteText,
			wantText:   "10",
			wantGlyphs: 2,
		},
		{
			name:       "labelled sensor value",
			block:      config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockLabelledSensorValue, Asset: "yellow_digits", Decoder: "rpm_digits", Geometry: config.RectConfig{X: 10, Y: 10, Width: 200, Height: 80}},
			wantType:   config.DashboardBlockSpriteText,
			wantText:   "10",
			wantGlyphs: 2,
		},
		{
			name:      "percent frame bar",
			block:     config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockPercentFrameBar, Asset: "throttle_frames", Decoder: "throttle_frame", Geometry: config.RectConfig{X: 10, Y: 10, Width: 200, Height: 30}},
			wantType:  config.DashboardBlockSpriteFrame,
			wantFrame: true,
		},
		{
			name:        "state lamp",
			block:       config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockStateLamp, Asset: "background", Geometry: config.RectConfig{X: 10, Y: 10, Width: 40, Height: 40}},
			wantType:    config.DashboardBlockImage,
			wantAssetID: "background",
		},
		{
			name:        "warning overlay",
			block:       config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockWarningOverlay, Asset: "background", Geometry: config.RectConfig{X: 0, Y: 0, Width: 800, Height: 480}},
			wantType:    config.DashboardBlockImage,
			wantAssetID: "background",
		},
		{
			name:        "stale overlay",
			block:       config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockStaleOverlay, Asset: "background", Geometry: config.RectConfig{X: 0, Y: 0, Width: 800, Height: 480}},
			wantType:    config.DashboardBlockImage,
			wantAssetID: "background",
		},
		{
			name:        "static panel",
			block:       config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockStaticPanel, Asset: "background", Geometry: config.RectConfig{X: 0, Y: 0, Width: 800, Height: 480}},
			wantType:    config.DashboardBlockImage,
			wantAssetID: "background",
		},
		{
			name:       "glowing number box as text",
			block:      config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockGlowingNumberBox, Asset: "yellow_digits", Decoder: "rpm_digits", Geometry: config.RectConfig{X: 10, Y: 10, Width: 200, Height: 80}},
			wantType:   config.DashboardBlockSpriteText,
			wantText:   "10",
			wantGlyphs: 2,
		},
		{
			name:           "glowing number box as group",
			block:          config.DashboardBlockConfig{ID: "subject", Type: config.DashboardBlockGlowingNumberBox, Blocks: []string{"panel", "digits"}},
			wantType:       config.DashboardBlockGroup,
			wantChildCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dashboard := baseDashboard()
			dashboard.Blocks = []config.DashboardBlockConfig{tt.block}
			if tt.wantChildCount > 0 {
				dashboard.Blocks = append(dashboard.Blocks,
					config.DashboardBlockConfig{ID: "panel", Type: config.DashboardBlockStaticPanel, Asset: "background", Geometry: config.RectConfig{X: 0, Y: 0, Width: 240, Height: 100}},
					config.DashboardBlockConfig{ID: "digits", Type: config.DashboardBlockSevenSegmentNumber, Asset: "yellow_digits", Decoder: "rpm_digits", Geometry: config.RectConfig{X: 20, Y: 10, Width: 200, Height: 80}},
				)
			}
			dashboard.Layers = []config.DashboardLayerConfig{{ID: "base", Z: 0, Blocks: []string{"subject"}}}

			scene, err := Evaluate(dashboard, registry, baseDecoderValues(), nil, Options{})
			if err != nil {
				t.Fatalf("Evaluate returned error: %v", err)
			}

			element := findElement(t, scene.Elements, "subject")
			if element.Type != tt.wantType {
				t.Fatalf("element.Type = %q, want %q", element.Type, tt.wantType)
			}
			if tt.wantText != "" && element.Text != tt.wantText {
				t.Fatalf("element.Text = %q, want %q", element.Text, tt.wantText)
			}
			if tt.wantGlyphs > 0 && len(element.Glyphs) != tt.wantGlyphs {
				t.Fatalf("len(element.Glyphs) = %d, want %d", len(element.Glyphs), tt.wantGlyphs)
			}
			if tt.wantFrame && !element.HasFrame {
				t.Fatal("element.HasFrame = false, want true")
			}
			if tt.wantAssetID != "" && element.AssetID != tt.wantAssetID {
				t.Fatalf("element.AssetID = %q, want %q", element.AssetID, tt.wantAssetID)
			}
			if tt.wantChildCount > 0 && len(element.Children) != tt.wantChildCount {
				t.Fatalf("len(element.Children) = %d, want %d", len(element.Children), tt.wantChildCount)
			}
		})
	}
}
