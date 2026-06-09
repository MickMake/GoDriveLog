package config

import "fmt"

const (
	DashboardAssetImage    = "image"
	DashboardAssetFrameSet = "frame_set"
	DashboardAssetCharset  = "charset"

	DashboardDecoderNormalize    = "normalize"
	DashboardDecoderThreshold    = "threshold"
	DashboardDecoderFrameIndex   = "frame_index"
	DashboardDecoderFormatNumber = "format_number"
	DashboardDecoderDigits       = "digits"
	DashboardDecoderBoolean      = "boolean"

	DashboardBlockImage       = "image"
	DashboardBlockSpriteFrame = "sprite_frame"
	DashboardBlockSpriteText  = "sprite_text"
	DashboardBlockGroup       = "group"
	DashboardBlockText        = "text"

	DashboardBlockSevenSegmentNumber  = "seven_segment_number"
	DashboardBlockPercentFrameBar     = "percent_frame_bar"
	DashboardBlockStateLamp           = "state_lamp"
	DashboardBlockGlowingNumberBox    = "glowing_number_box"
	DashboardBlockLabelledSensorValue = "labelled_sensor_value"
	DashboardBlockWarningOverlay      = "warning_overlay"
	DashboardBlockStaleOverlay        = "stale_overlay"
	DashboardBlockStaticPanel         = "static_panel"
)

type DashboardConfig struct {
	RefreshMS   int                      `yaml:"refresh_ms"`
	RenderMinMS int                      `yaml:"render_min_ms"`
	Canvas      CanvasConfig             `yaml:"canvas"`
	AssetRoot   string                   `yaml:"asset_root"`
	Assets      []DashboardAssetConfig   `yaml:"assets"`
	Decoders    []DashboardDecoderConfig `yaml:"decoders"`
	Blocks      []DashboardBlockConfig   `yaml:"blocks"`
	Layers      []DashboardLayerConfig   `yaml:"layers"`
}

type CanvasConfig struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type DashboardAssetConfig struct {
	ID         string            `yaml:"id"`
	Type       string            `yaml:"type"`
	Path       string            `yaml:"path"`
	Pattern    string            `yaml:"pattern"`
	FrameCount int               `yaml:"frame_count"`
	Frames     []string          `yaml:"frames"`
	Glyphs     map[string]string `yaml:"glyphs"`
}

type DashboardDecoderConfig struct {
	ID         string            `yaml:"id"`
	Type       string            `yaml:"type"`
	Sensor     string            `yaml:"sensor"`
	Input      string            `yaml:"input"`
	Asset      string            `yaml:"asset"`
	Format     string            `yaml:"format"`
	FrameCount int               `yaml:"frame_count"`
	Thresholds []ThresholdConfig `yaml:"thresholds"`
}

type ThresholdConfig struct {
	At    float64 `yaml:"at"`
	Value string  `yaml:"value"`
}

type DashboardConditionConfig struct {
	Sensor    string   `yaml:"sensor"`
	Decoder   string   `yaml:"decoder"`
	Status    string   `yaml:"status"`
	Equals    string   `yaml:"equals"`
	NotEquals string   `yaml:"not_equals"`
	Min       *float64 `yaml:"min"`
	Max       *float64 `yaml:"max"`
}

type DashboardBlockConfig struct {
	ID        string                   `yaml:"id"`
	Type      string                   `yaml:"type"`
	Asset     string                   `yaml:"asset"`
	Decoder   string                   `yaml:"decoder"`
	Blocks    []string                 `yaml:"blocks"`
	Condition DashboardConditionConfig `yaml:"condition"`
	Geometry  RectConfig               `yaml:"geometry"`
}

type RectConfig struct {
	X      float64 `yaml:"x"`
	Y      float64 `yaml:"y"`
	Width  float64 `yaml:"width"`
	Height float64 `yaml:"height"`
}

type DashboardLayerConfig struct {
	ID     string   `yaml:"id"`
	Z      int      `yaml:"z"`
	Blocks []string `yaml:"blocks"`
}

func validateDashboard(cfg Config) error {
	if cfg.Dashboard.RefreshMS <= 0 {
		return fmt.Errorf("dashboard.refresh_ms must be positive")
	}
	if cfg.Dashboard.RenderMinMS < 0 {
		return fmt.Errorf("dashboard.render_min_ms must not be negative")
	}
	if cfg.Dashboard.Canvas.Width <= 0 {
		return fmt.Errorf("dashboard.canvas.width must be positive")
	}
	if cfg.Dashboard.Canvas.Height <= 0 {
		return fmt.Errorf("dashboard.canvas.height must be positive")
	}

	assetIDs, err := validateAssets(cfg.Dashboard.Assets)
	if err != nil {
		return err
	}
	decoderIDs, err := validateDecoders(cfg.Dashboard.Decoders, cfg.Sensors, assetIDs)
	if err != nil {
		return err
	}
	blockIDs, err := validateBlocks(cfg.Dashboard.Blocks, assetIDs, decoderIDs, cfg.Sensors)
	if err != nil {
		return err
	}
	if err := validateLayers(cfg.Dashboard.Layers, blockIDs); err != nil {
		return err
	}

	return nil
}

func validateAssets(assets []DashboardAssetConfig) (map[string]bool, error) {
	ids := map[string]bool{}
	for i, asset := range assets {
		path := fmt.Sprintf("dashboard.assets[%d]", i)
		if asset.ID == "" {
			return nil, fmt.Errorf("%s.id must not be empty", path)
		}
		if ids[asset.ID] {
			return nil, fmt.Errorf("dashboard.assets id %q must be unique", asset.ID)
		}
		ids[asset.ID] = true

		switch asset.Type {
		case DashboardAssetImage:
			if asset.Path == "" {
				return nil, fmt.Errorf("%s.path must not be empty for image assets", path)
			}
		case DashboardAssetFrameSet:
			hasExplicitFrames := len(asset.Frames) > 0
			hasGeneratedFrames := asset.Pattern != "" || asset.FrameCount > 0
			if !hasExplicitFrames && !hasGeneratedFrames {
				return nil, fmt.Errorf("%s.frames or pattern/frame_count must not be empty for frame_set assets", path)
			}
			if hasExplicitFrames && asset.Pattern != "" {
				return nil, fmt.Errorf("%s must not define both frames and pattern", path)
			}
			if asset.Pattern != "" && asset.FrameCount <= 0 {
				return nil, fmt.Errorf("%s.frame_count must be positive for generated frame_set assets", path)
			}
			if asset.Pattern == "" && asset.FrameCount > 0 && len(asset.Frames) != asset.FrameCount {
				return nil, fmt.Errorf("%s.frame_count must match frames length", path)
			}
		case DashboardAssetCharset:
			if len(asset.Glyphs) == 0 {
				return nil, fmt.Errorf("%s.glyphs must not be empty for charset assets", path)
			}
		default:
			return nil, fmt.Errorf("%s.type must be image, frame_set, or charset", path)
		}
	}
	return ids, nil
}

func validateDecoders(decoders []DashboardDecoderConfig, sensors map[string]SensorConfig, assets map[string]bool) (map[string]bool, error) {
	ids := map[string]bool{}
	for i, decoder := range decoders {
		path := fmt.Sprintf("dashboard.decoders[%d]", i)
		if decoder.ID == "" {
			return nil, fmt.Errorf("%s.id must not be empty", path)
		}
		if ids[decoder.ID] {
			return nil, fmt.Errorf("dashboard.decoders id %q must be unique", decoder.ID)
		}

		switch decoder.Type {
		case DashboardDecoderNormalize, DashboardDecoderThreshold, DashboardDecoderFrameIndex, DashboardDecoderFormatNumber, DashboardDecoderDigits, DashboardDecoderBoolean:
		default:
			return nil, fmt.Errorf("%s.type must be a supported decoder type", path)
		}

		if decoder.Sensor == "" && decoder.Input == "" {
			return nil, fmt.Errorf("%s must define sensor or input", path)
		}
		if decoder.Sensor != "" && decoder.Input != "" {
			return nil, fmt.Errorf("%s must not define both sensor and input", path)
		}
		if decoder.Sensor != "" {
			if _, ok := sensors[decoder.Sensor]; !ok {
				return nil, fmt.Errorf("%s.sensor %q must reference a configured sensor", path, decoder.Sensor)
			}
		}
		if decoder.Input != "" && !ids[decoder.Input] {
			return nil, fmt.Errorf("%s.input %q must reference an earlier decoder", path, decoder.Input)
		}
		if decoder.Asset != "" && !assets[decoder.Asset] {
			return nil, fmt.Errorf("%s.asset %q must reference a configured asset", path, decoder.Asset)
		}
		if decoder.Type == DashboardDecoderFrameIndex && decoder.FrameCount <= 0 {
			return nil, fmt.Errorf("%s.frame_count must be positive for frame_index decoders", path)
		}
		if decoder.Type == DashboardDecoderThreshold && len(decoder.Thresholds) == 0 {
			return nil, fmt.Errorf("%s.thresholds must not be empty for threshold decoders", path)
		}

		ids[decoder.ID] = true
	}
	return ids, nil
}

func validateBlocks(blocks []DashboardBlockConfig, assets map[string]bool, decoders map[string]bool, sensors map[string]SensorConfig) (map[string]bool, error) {
	ids := map[string]bool{}
	for i, block := range blocks {
		path := fmt.Sprintf("dashboard.blocks[%d]", i)
		if block.ID == "" {
			return nil, fmt.Errorf("%s.id must not be empty", path)
		}
		if ids[block.ID] {
			return nil, fmt.Errorf("dashboard.blocks id %q must be unique", block.ID)
		}
		ids[block.ID] = true

		if !isSupportedDashboardBlockType(block.Type) {
			return nil, fmt.Errorf("%s.type must be a supported block type", path)
		}

		if block.Asset != "" && !assets[block.Asset] {
			return nil, fmt.Errorf("%s.asset %q must reference a configured asset", path, block.Asset)
		}
		if block.Decoder != "" && !decoders[block.Decoder] {
			return nil, fmt.Errorf("%s.decoder %q must reference a configured decoder", path, block.Decoder)
		}
		if err := validateCondition(path, block.Condition, sensors, decoders); err != nil {
			return nil, err
		}
		if isDashboardGroupBlock(block) && len(block.Blocks) == 0 {
			return nil, fmt.Errorf("%s.blocks must not be empty for group blocks", path)
		}
		if !isDashboardGroupBlock(block) && block.Geometry.Width <= 0 {
			return nil, fmt.Errorf("%s.geometry.width must be positive", path)
		}
		if !isDashboardGroupBlock(block) && block.Geometry.Height <= 0 {
			return nil, fmt.Errorf("%s.geometry.height must be positive", path)
		}
	}

	for i, block := range blocks {
		if !isDashboardGroupBlock(block) {
			continue
		}
		path := fmt.Sprintf("dashboard.blocks[%d]", i)
		for _, childID := range block.Blocks {
			if !ids[childID] {
				return nil, fmt.Errorf("%s.blocks %q must reference a configured block", path, childID)
			}
		}
	}

	return ids, nil
}

func isSupportedDashboardBlockType(blockType string) bool {
	switch blockType {
	case DashboardBlockImage, DashboardBlockSpriteFrame, DashboardBlockSpriteText, DashboardBlockGroup, DashboardBlockText,
		DashboardBlockSevenSegmentNumber, DashboardBlockPercentFrameBar, DashboardBlockStateLamp, DashboardBlockGlowingNumberBox,
		DashboardBlockLabelledSensorValue, DashboardBlockWarningOverlay, DashboardBlockStaleOverlay, DashboardBlockStaticPanel:
		return true
	default:
		return false
	}
}

func isDashboardGroupBlock(block DashboardBlockConfig) bool {
	return block.Type == DashboardBlockGroup || (block.Type == DashboardBlockGlowingNumberBox && len(block.Blocks) > 0)
}

func validateCondition(path string, condition DashboardConditionConfig, sensors map[string]SensorConfig, decoders map[string]bool) error {
	if isEmptyDashboardCondition(condition) {
		return nil
	}
	if condition.Sensor != "" && condition.Decoder != "" {
		return fmt.Errorf("%s.condition must not define both sensor and decoder", path)
	}
	if condition.Sensor == "" && condition.Decoder == "" {
		return fmt.Errorf("%s.condition must define sensor or decoder", path)
	}
	if condition.Sensor != "" {
		if _, ok := sensors[condition.Sensor]; !ok {
			return fmt.Errorf("%s.condition.sensor %q must reference a configured sensor", path, condition.Sensor)
		}
	}
	if condition.Decoder != "" && !decoders[condition.Decoder] {
		return fmt.Errorf("%s.condition.decoder %q must reference a configured decoder", path)
	}
	return nil
}

func isEmptyDashboardCondition(condition DashboardConditionConfig) bool {
	return condition.Sensor == "" && condition.Decoder == "" && condition.Status == "" && condition.Equals == "" && condition.NotEquals == "" && condition.Min == nil && condition.Max == nil
}

func validateLayers(layers []DashboardLayerConfig, blocks map[string]bool) error {
	if len(layers) == 0 {
		return fmt.Errorf("dashboard.layers must not be empty")
	}
	ids := map[string]bool{}
	for i, layer := range layers {
		path := fmt.Sprintf("dashboard.layers[%d]", i)
		if layer.ID == "" {
			return fmt.Errorf("%s.id must not be empty", path)
		}
		if ids[layer.ID] {
			return fmt.Errorf("dashboard.layers id %q must be unique")
		}
		ids[layer.ID] = true
		if len(layer.Blocks) == 0 {
			return fmt.Errorf("%s.blocks must not be empty", path)
		}
		for _, blockID := range layer.Blocks {
			if !blocks[blockID] {
				return fmt.Errorf("%s.blocks %q must reference a configured block", path, blockID)
			}
		}
	}
	return nil
}
