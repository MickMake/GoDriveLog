package gauges

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/sensors"
)

const (
	ScenePartKindLayer        = "layer"
	ScenePartKindBackground   = "background"
	ScenePartKindCharacter    = "character"
	ScenePartKindDecimalPoint = "decimal_point"
	ScenePartKindForeground   = "foreground"
	ScenePartKindNeedle       = "needle"
)

type Placement struct {
	Position []int
	Scale    float64
}

type Scene struct {
	PackageID      string
	PackagePath    string
	Type           string
	SensorID       string
	Position       []int
	Scale          float64
	Size           Size
	Status         string
	Error          string
	Text           string
	DigitPositions [][]int
	FacePivot      Point
	NeedlePivot    Point
	Angle          float64
	Parts          []ScenePart
}

type ScenePart struct {
	Kind        string
	Layer       string
	AssetPath   string
	Slot        int
	Character   string
	Position    []int
	Angle       float64
	FacePivot   Point
	NeedlePivot Point
}

func NumericScene(pkg Package, placement Placement, state sensors.SensorState) (Scene, error) {
	if pkg.Type != TypeNumeric {
		return Scene{}, fmt.Errorf("gauge package %q type %q is not numeric", pkg.ID, pkg.Type)
	}
	if placement.Scale <= 0 {
		return Scene{}, fmt.Errorf("gauge package %q placement scale must be greater than zero", pkg.ID)
	}
	if pkg.Digits.Count <= 0 {
		return Scene{}, fmt.Errorf("gauge package %q digits count must be greater than zero", pkg.ID)
	}
	if len(pkg.Digits.Positions) != pkg.Digits.Count {
		return Scene{}, fmt.Errorf("gauge package %q must define %d digit positions", pkg.ID, pkg.Digits.Count)
	}

	state = stateForPackage(pkg.Sensor, state)
	scene := Scene{
		PackageID:      pkg.ID,
		PackagePath:    pkg.Path,
		Type:           pkg.Type,
		SensorID:       pkg.Sensor,
		Position:       cloneInts(placement.Position),
		Scale:          placement.Scale,
		Size:           pkg.Size,
		Status:         state.Status,
		Error:          state.Error,
		DigitPositions: cloneIntSlices(pkg.Digits.Positions),
		Parts:          underlayLayerParts(pkg.Layers),
	}

	if state.Status != sensors.StatusOK {
		scene.Parts = append(scene.Parts, overlayLayerParts(pkg.Layers)...)
		return scene, nil
	}

	text := formatValue(pkg.Format, state.Value)
	characters, decimalSlots, err := splitTextIntoSlots(text, pkg.Digits.Count)
	if err != nil {
		return Scene{}, fmt.Errorf("gauge package %q: %w", pkg.ID, err)
	}
	decimalBySlot := map[int]bool{}
	for _, slot := range decimalSlots {
		decimalBySlot[slot] = true
	}

	scene.Text = text
	for slot, ch := range characters {
		position := digitPosition(pkg, slot)
		if pkg.DigitSet.Background != "" {
			scene.Parts = append(scene.Parts, ScenePart{Kind: ScenePartKindBackground, AssetPath: pkg.DigitSet.Background, Slot: slot, Position: position})
		}
		if ch != " " {
			assetPath, ok := pkg.DigitSet.Characters[ch]
			if !ok {
				return Scene{}, fmt.Errorf("gauge package %q digit set has no character asset for %q", pkg.ID, ch)
			}
			scene.Parts = append(scene.Parts, ScenePart{Kind: ScenePartKindCharacter, AssetPath: assetPath, Slot: slot, Character: ch, Position: position})
		}
		if decimalBySlot[slot] {
			if pkg.DigitSet.DecimalPoint == "" {
				return Scene{}, fmt.Errorf("gauge package %q formatted output requires digit_set decimal_point", pkg.ID)
			}
			scene.Parts = append(scene.Parts, ScenePart{Kind: ScenePartKindDecimalPoint, AssetPath: pkg.DigitSet.DecimalPoint, Slot: slot, Position: position})
		}
		if pkg.DigitSet.Foreground != "" {
			scene.Parts = append(scene.Parts, ScenePart{Kind: ScenePartKindForeground, AssetPath: pkg.DigitSet.Foreground, Slot: slot, Position: position})
		}
	}
	scene.Parts = append(scene.Parts, overlayLayerParts(pkg.Layers)...)
	return scene, nil
}

func RadialScene(pkg Package, placement Placement, state sensors.SensorState) (Scene, error) {
	if pkg.Type != TypeRadial {
		return Scene{}, fmt.Errorf("gauge package %q type %q is not radial", pkg.ID, pkg.Type)
	}
	if placement.Scale <= 0 {
		return Scene{}, fmt.Errorf("gauge package %q placement scale must be greater than zero", pkg.ID)
	}
	needlePath := strings.TrimSpace(pkg.Layers["needle"])
	if needlePath == "" {
		return Scene{}, fmt.Errorf("gauge package %q radial layer needle must not be empty", pkg.ID)
	}
	if pkg.ValueMap.Max <= pkg.ValueMap.Min {
		return Scene{}, fmt.Errorf("gauge package %q value_map max must be greater than min", pkg.ID)
	}

	state = stateForPackage(pkg.Sensor, state)
	scene := Scene{
		PackageID:   pkg.ID,
		PackagePath: pkg.Path,
		Type:        pkg.Type,
		SensorID:    pkg.Sensor,
		Position:    cloneInts(placement.Position),
		Scale:       placement.Scale,
		Size:        pkg.Size,
		Status:      state.Status,
		Error:       state.Error,
		FacePivot:   pkg.Pivot.Face,
		NeedlePivot: pkg.Pivot.Needle,
		Parts:       radialUnderlayLayerParts(pkg.Layers),
	}

	if state.Status == sensors.StatusOK {
		angle, err := radialAngle(pkg.ValueMap, state.Value)
		if err != nil {
			return Scene{}, fmt.Errorf("gauge package %q: %w", pkg.ID, err)
		}
		scene.Angle = angle
		scene.Parts = append(scene.Parts, ScenePart{
			Kind:        ScenePartKindNeedle,
			Layer:       "needle",
			AssetPath:   needlePath,
			Angle:       angle,
			FacePivot:   pkg.Pivot.Face,
			NeedlePivot: pkg.Pivot.Needle,
		})
	}

	scene.Parts = append(scene.Parts, radialOverlayLayerParts(pkg.Layers)...)
	return scene, nil
}

func radialAngle(valueMap ValueMap, value float64) (float64, error) {
	if valueMap.Max <= valueMap.Min {
		return 0, fmt.Errorf("value_map max must be greater than min")
	}
	mappedValue := value
	if valueMap.Clamp {
		if mappedValue < valueMap.Min {
			mappedValue = valueMap.Min
		}
		if mappedValue > valueMap.Max {
			mappedValue = valueMap.Max
		}
	}
	ratio := (mappedValue - valueMap.Min) / (valueMap.Max - valueMap.Min)
	return valueMap.StartAngle + ratio*(valueMap.EndAngle-valueMap.StartAngle), nil
}

func (s Scene) Signature() string {
	var b strings.Builder
	b.WriteString(s.PackageID)
	b.WriteString("|")
	b.WriteString(s.PackagePath)
	b.WriteString("|")
	b.WriteString(s.Type)
	b.WriteString("|")
	b.WriteString(s.SensorID)
	b.WriteString("|")
	b.WriteString(formatIntSlice(s.Position))
	b.WriteString("|")
	b.WriteString(strconv.FormatFloat(s.Scale, 'f', -1, 64))
	b.WriteString("|")
	b.WriteString(strconv.Itoa(s.Size.Width))
	b.WriteString("x")
	b.WriteString(strconv.Itoa(s.Size.Height))
	b.WriteString("|")
	b.WriteString(s.Status)
	b.WriteString("|")
	b.WriteString(s.Error)
	b.WriteString("|")
	b.WriteString(s.Text)
	b.WriteString("|positions=")
	b.WriteString(formatIntSlices(s.DigitPositions))
	b.WriteString("|face=")
	b.WriteString(formatPoint(s.FacePivot))
	b.WriteString("|needle=")
	b.WriteString(formatPoint(s.NeedlePivot))
	b.WriteString("|angle=")
	b.WriteString(strconv.FormatFloat(s.Angle, 'f', -1, 64))
	b.WriteString("|")
	for _, part := range s.Parts {
		b.WriteString(part.Kind)
		b.WriteString("@")
		b.WriteString(part.Layer)
		b.WriteString("#")
		b.WriteString(strconv.Itoa(part.Slot))
		b.WriteString("=")
		b.WriteString(part.AssetPath)
		b.WriteString("#")
		b.WriteString(part.Character)
		b.WriteString("#")
		b.WriteString(formatIntSlice(part.Position))
		b.WriteString("#")
		b.WriteString(strconv.FormatFloat(part.Angle, 'f', -1, 64))
		b.WriteString("#")
		b.WriteString(formatPoint(part.FacePivot))
		b.WriteString("#")
		b.WriteString(formatPoint(part.NeedlePivot))
		b.WriteString(";")
	}
	return b.String()
}

func stateForPackage(sensorID string, state sensors.SensorState) sensors.SensorState {
	if state.ID == "" {
		state.ID = sensorID
	}
	if state.Status == "" {
		state.Status = sensors.StatusUnknown
	}
	return state
}

func underlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"background", "panel", "bezel", "face", "ticks"})
}

func overlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"glass", "overlay", "foreground"})
}

func radialUnderlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"background", "panel", "bezel", "face", "ticks"})
}

func radialOverlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"glass", "overlay", "foreground"})
}

func namedLayerParts(layers map[string]string, orderedNames []string) []ScenePart {
	parts := []ScenePart{}
	for _, name := range orderedNames {
		assetPath := strings.TrimSpace(layers[name])
		if assetPath == "" {
			continue
		}
		parts = append(parts, ScenePart{Kind: ScenePartKindLayer, Layer: name, AssetPath: assetPath})
	}
	return parts
}

func digitPosition(pkg Package, slot int) []int {
	if slot < 0 || slot >= len(pkg.Digits.Positions) {
		return nil
	}
	return cloneInts(pkg.Digits.Positions[slot])
}

func formatValue(format string, value float64) string {
	if strings.TrimSpace(format) == "" {
		return fmt.Sprintf("%.0f", value)
	}
	return fmt.Sprintf(format, value)
}

func splitTextIntoSlots(text string, slots int) ([]string, []int, error) {
	characters := make([]string, 0, len(text))
	decimalSlots := []int{}
	lastSlot := -1
	for _, r := range text {
		ch := string(r)
		if ch == "." {
			if lastSlot < 0 {
				return nil, nil, fmt.Errorf("decimal separator has no preceding digit slot")
			}
			decimalSlots = append(decimalSlots, lastSlot)
			continue
		}
		characters = append(characters, ch)
		lastSlot = len(characters) - 1
	}
	if len(characters) > slots {
		return nil, nil, fmt.Errorf("formatted output %q needs %d character slots, gauge package allows %d", text, len(characters), slots)
	}
	padded := make([]string, slots)
	padding := slots - len(characters)
	for i := 0; i < padding; i++ {
		padded[i] = " "
	}
	copy(padded[padding:], characters)
	if padding > 0 {
		for i := range decimalSlots {
			decimalSlots[i] += padding
		}
	}
	return padded, decimalSlots, nil
}

func cloneInts(values []int) []int {
	if values == nil {
		return nil
	}
	return append([]int(nil), values...)
}

func cloneIntSlices(values [][]int) [][]int {
	if values == nil {
		return nil
	}
	cloned := make([][]int, len(values))
	for i, value := range values {
		cloned[i] = cloneInts(value)
	}
	return cloned
}

func formatIntSlice(values []int) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = strconv.Itoa(value)
	}
	return strings.Join(parts, ",")
}

func formatIntSlices(values [][]int) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = formatIntSlice(value)
	}
	return strings.Join(parts, ";")
}

func formatPoint(point Point) string {
	if point == (Point{}) {
		return ""
	}
	return strconv.FormatFloat(point.X, 'f', -1, 64) + "," + strconv.FormatFloat(point.Y, 'f', -1, 64)
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.000001
}
