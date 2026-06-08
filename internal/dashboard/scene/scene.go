package scene

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/assets"
	"github.com/MickMake/GoDriveLog/internal/dashboard/decoders"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

type Options struct {
	Conditions map[string]Condition
}

type Condition struct {
	Sensor    string
	Decoder   string
	Equals    string
	NotEquals string
	Min       *float64
	Max       *float64
}

type Scene struct {
	Elements []Element
}

type Element struct {
	ID       string
	Type     string
	LayerID  string
	Z        int
	AssetID  string
	Decoder  string
	Geometry config.RectConfig
	Visible  bool
	Children []Element
	Frame    assets.Frame
	HasFrame bool
	Text     string
	Glyphs   []assets.Glyph
}

func Evaluate(dashboard config.DashboardConfig, assetRegistry *assets.Registry, decoderValues map[string]decoders.Value, sensorStates map[string]sensors.SensorState, options Options) (Scene, error) {
	if assetRegistry == nil {
		return Scene{}, fmt.Errorf("asset registry must not be nil")
	}
	blocks := map[string]config.DashboardBlockConfig{}
	for _, block := range dashboard.Blocks {
		blocks[block.ID] = block
	}

	layers := append([]config.DashboardLayerConfig(nil), dashboard.Layers...)
	sort.SliceStable(layers, func(i, j int) bool {
		return layers[i].Z < layers[j].Z
	})

	scene := Scene{}
	for _, layer := range layers {
		for _, blockID := range layer.Blocks {
			block, ok := blocks[blockID]
			if !ok {
				return Scene{}, fmt.Errorf("layer %q block %q is not configured", layer.ID, blockID)
			}
			element, err := resolveBlock(block, layer.ID, layer.Z, blocks, assetRegistry, decoderValues, sensorStates, options, nil)
			if err != nil {
				return Scene{}, err
			}
			scene.Elements = append(scene.Elements, element)
		}
	}
	return scene, nil
}

func resolveBlock(block config.DashboardBlockConfig, layerID string, z int, blocks map[string]config.DashboardBlockConfig, assetRegistry *assets.Registry, decoderValues map[string]decoders.Value, sensorStates map[string]sensors.SensorState, options Options, stack []string) (Element, error) {
	activeStack, err := pushResolutionStack(stack, block.ID)
	if err != nil {
		return Element{}, err
	}

	visible, err := evaluateCondition(options.Conditions[block.ID], decoderValues, sensorStates)
	if err != nil {
		return Element{}, fmt.Errorf("block %q condition: %w", block.ID, err)
	}

	element := Element{ID: block.ID, Type: block.Type, LayerID: layerID, Z: z, AssetID: block.Asset, Decoder: block.Decoder, Geometry: block.Geometry, Visible: visible}
	if !visible {
		return element, nil
	}

	switch block.Type {
	case config.DashboardBlockImage:
		asset, err := requireAsset(assetRegistry, block.Asset, assets.TypeImage)
		if err != nil {
			return Element{}, fmt.Errorf("block %q: %w", block.ID, err)
		}
		element.AssetID = asset.ID
	case config.DashboardBlockSpriteFrame:
		asset, err := requireAsset(assetRegistry, block.Asset, assets.TypeFrameSet)
		if err != nil {
			return Element{}, fmt.Errorf("block %q: %w", block.ID, err)
		}
		frameIndex, err := frameIndexFor(block, decoderValues)
		if err != nil {
			return Element{}, err
		}
		if frameIndex < 0 || frameIndex >= len(asset.Frames) {
			return Element{}, fmt.Errorf("block %q frame index %d is outside asset %q frame range", block.ID, frameIndex, block.Asset)
		}
		element.Frame = asset.Frames[frameIndex]
		element.HasFrame = true
	case config.DashboardBlockSpriteText:
		asset, err := requireAsset(assetRegistry, block.Asset, assets.TypeCharset)
		if err != nil {
			return Element{}, fmt.Errorf("block %q: %w", block.ID, err)
		}
		text, err := textFor(block, decoderValues)
		if err != nil {
			return Element{}, err
		}
		element.Text = text
		for _, ch := range text {
			glyph, ok := asset.Glyphs[string(ch)]
			if !ok {
				return Element{}, fmt.Errorf("block %q charset asset %q has no glyph for %q", block.ID, block.Asset, string(ch))
			}
			element.Glyphs = append(element.Glyphs, glyph)
		}
	case config.DashboardBlockGroup:
		for _, childID := range block.Blocks {
			child, ok := blocks[childID]
			if !ok {
				return Element{}, fmt.Errorf("block %q child %q is not configured", block.ID, childID)
			}
			childElement, err := resolveBlock(child, layerID, z, blocks, assetRegistry, decoderValues, sensorStates, options, activeStack)
			if err != nil {
				return Element{}, err
			}
			element.Children = append(element.Children, childElement)
		}
	default:
		return Element{}, fmt.Errorf("block %q type %q is not a supported scene primitive", block.ID, block.Type)
	}

	return element, nil
}

func pushResolutionStack(stack []string, blockID string) ([]string, error) {
	for i, activeID := range stack {
		if activeID == blockID {
			cycle := append(append([]string(nil), stack[i:]...), blockID)
			return nil, fmt.Errorf("cyclic dashboard scene block reference detected: %s", strings.Join(cycle, " -> "))
		}
	}
	return append(append([]string(nil), stack...), blockID), nil
}

func requireAsset(registry *assets.Registry, id string, assetType string) (assets.Asset, error) {
	asset, err := registry.MustGet(id)
	if err != nil {
		return assets.Asset{}, err
	}
	if asset.Type != assetType {
		return assets.Asset{}, fmt.Errorf("asset %q type is %q, want %q", id, asset.Type, assetType)
	}
	return asset, nil
}

func frameIndexFor(block config.DashboardBlockConfig, values map[string]decoders.Value) (int, error) {
	value, ok := values[block.Decoder]
	if !ok {
		return 0, fmt.Errorf("block %q decoder %q is not available", block.ID, block.Decoder)
	}
	if value.Type == decoders.ValueTypeFrameIndex {
		return value.FrameIndex, nil
	}
	number, err := value.NumberValue()
	if err != nil {
		return 0, fmt.Errorf("block %q decoder %q is not a frame index: %w", block.ID, block.Decoder, err)
	}
	return int(number), nil
}

func textFor(block config.DashboardBlockConfig, values map[string]decoders.Value) (string, error) {
	value, ok := values[block.Decoder]
	if !ok {
		return "", fmt.Errorf("block %q decoder %q is not available", block.ID, block.Decoder)
	}
	switch value.Type {
	case decoders.ValueTypeDigits:
		text := ""
		for _, digit := range value.Digits {
			text += digit
		}
		return text, nil
	case decoders.ValueTypeText:
		return value.Text, nil
	case decoders.ValueTypeNumber, decoders.ValueTypeFrameIndex:
		return strconv.FormatFloat(value.Number, 'f', -1, 64), nil
	case decoders.ValueTypeBoolean:
		return strconv.FormatBool(value.Bool), nil
	default:
		return "", fmt.Errorf("block %q decoder %q has unsupported text value type %q", block.ID, block.Decoder, value.Type)
	}
}

func evaluateCondition(condition Condition, values map[string]decoders.Value, sensorStates map[string]sensors.SensorState) (bool, error) {
	if isEmptyCondition(condition) {
		return true, nil
	}
	value, err := conditionValue(condition, values, sensorStates)
	if err != nil {
		return false, err
	}
	if condition.Equals != "" && value.text != condition.Equals {
		return false, nil
	}
	if condition.NotEquals != "" && value.text == condition.NotEquals {
		return false, nil
	}
	if condition.Min != nil {
		if !value.hasNumber {
			return false, fmt.Errorf("min requires a numeric value")
		}
		if value.number < *condition.Min {
			return false, nil
		}
	}
	if condition.Max != nil {
		if !value.hasNumber {
			return false, fmt.Errorf("max requires a numeric value")
		}
		if value.number > *condition.Max {
			return false, nil
		}
	}
	return true, nil
}

type comparableValue struct {
	text      string
	number    float64
	hasNumber bool
}

func conditionValue(condition Condition, values map[string]decoders.Value, sensorStates map[string]sensors.SensorState) (comparableValue, error) {
	if condition.Sensor != "" && condition.Decoder != "" {
		return comparableValue{}, fmt.Errorf("condition must not define both sensor and decoder")
	}
	if condition.Sensor != "" {
		state, ok := sensorStates[condition.Sensor]
		if !ok {
			return comparableValue{}, fmt.Errorf("sensor %q is not available", condition.Sensor)
		}
		return comparableValue{text: strconv.FormatFloat(state.Value, 'f', -1, 64), number: state.Value, hasNumber: true}, nil
	}
	if condition.Decoder != "" {
		value, ok := values[condition.Decoder]
		if !ok {
			return comparableValue{}, fmt.Errorf("decoder %q is not available", condition.Decoder)
		}
		return comparableDecoderValue(value)
	}
	return comparableValue{}, fmt.Errorf("condition must define sensor or decoder")
}

func comparableDecoderValue(value decoders.Value) (comparableValue, error) {
	switch value.Type {
	case decoders.ValueTypeBoolean:
		return comparableValue{text: strconv.FormatBool(value.Bool), number: boolNumber(value.Bool), hasNumber: true}, nil
	case decoders.ValueTypeText:
		return comparableValue{text: value.Text}, nil
	case decoders.ValueTypeDigits:
		text := ""
		for _, digit := range value.Digits {
			text += digit
		}
		return comparableValue{text: text}, nil
	case decoders.ValueTypeNumber, decoders.ValueTypeFrameIndex:
		return comparableValue{text: strconv.FormatFloat(value.Number, 'f', -1, 64), number: value.Number, hasNumber: true}, nil
	default:
		return comparableValue{}, fmt.Errorf("decoder value type %q cannot be used in a condition", value.Type)
	}
}

func boolNumber(value bool) float64 {
	if value {
		return 1
	}
	return 0
}

func isEmptyCondition(condition Condition) bool {
	return condition.Sensor == "" && condition.Decoder == "" && condition.Equals == "" && condition.NotEquals == "" && condition.Min == nil && condition.Max == nil
}
