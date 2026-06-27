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
	ScenePartKindBar          = "bar"
	ScenePartKindWheelStrip   = "wheel_strip"
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
	Movement       string
	BarMode        string
	BarAxis        string
	BarOrigin      string
	BarBounds      []int
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
	Source      []int
	Window      Size
	StripOffset float64
	Wraparound  bool
	Role        string
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

func OdometerScene(pkg Package, placement Placement, state sensors.SensorState) (Scene, error) {
	if pkg.Type != TypeOdometer {
		return Scene{}, fmt.Errorf("gauge package %q type %q is not odometer", pkg.ID, pkg.Type)
	}
	if placement.Scale <= 0 {
		return Scene{}, fmt.Errorf("gauge package %q placement scale must be greater than zero", pkg.ID)
	}
	if len(pkg.Odometer.Wheels) == 0 {
		return Scene{}, fmt.Errorf("gauge package %q odometer wheels must not be empty", pkg.ID)
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
		Movement:    pkg.Odometer.Movement,
		Parts:       odometerUnderlayLayerParts(pkg.Layers),
	}

	if state.Status == sensors.StatusOK {
		scene.Text = formatValue("%.1f", state.Value)
		digitPlaces := odometerDigitPlaces(pkg.Odometer.Wheels)
		wraparound := odometerWraparoundEnabled(pkg)
		for index, wheel := range pkg.Odometer.Wheels {
			offset, err := odometerWheelOffset(pkg.Odometer.Movement, wraparound, wheel, digitPlaces[index], state.Value)
			if err != nil {
				return Scene{}, fmt.Errorf("gauge package %q: %w", pkg.ID, err)
			}
			sourceX, sourceY := odometerWheelSource(wheel, offset)
			scene.Parts = append(scene.Parts, ScenePart{
				Kind:        ScenePartKindWheelStrip,
				AssetPath:   wheel.Strip,
				Slot:        index,
				Position:    cloneInts(wheel.Position),
				Source:      []int{sourceX, sourceY},
				Window:      wheel.Window,
				StripOffset: offset,
				Wraparound:  wraparound,
				Role:        odometerWheelRole(wheel),
			})
		}
	}

	scene.Parts = append(scene.Parts, odometerOverlayLayerParts(pkg.Layers)...)
	return scene, nil
}

func IndicatorScene(pkg Package, placement Placement, state sensors.SensorState) (Scene, error) {
	if pkg.Type != TypeIndicator {
		return Scene{}, fmt.Errorf("gauge package %q type %q is not indicator", pkg.ID, pkg.Type)
	}
	if placement.Scale <= 0 {
		return Scene{}, fmt.Errorf("gauge package %q placement scale must be greater than zero", pkg.ID)
	}
	offPath := strings.TrimSpace(pkg.Layers["off"])
	onPath := strings.TrimSpace(pkg.Layers["on"])
	if onPath == "" {
		return Scene{}, fmt.Errorf("gauge package %q indicator layer on must not be empty", pkg.ID)
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
		Parts:       indicatorUnderlayLayerParts(pkg.Layers),
	}

	if indicatorStateOn(state) {
		scene.Parts = append(scene.Parts, ScenePart{Kind: ScenePartKindLayer, Layer: "on", AssetPath: onPath})
	} else if offPath != "" {
		scene.Parts = append(scene.Parts, ScenePart{Kind: ScenePartKindLayer, Layer: "off", AssetPath: offPath})
	}
	scene.Parts = append(scene.Parts, indicatorOverlayLayerParts(pkg.Layers)...)
	return scene, nil
}

func BarScene(pkg Package, placement Placement, state sensors.SensorState) (Scene, error) {
	if pkg.Type != TypeBar {
		return Scene{}, fmt.Errorf("gauge package %q type %q is not bar", pkg.ID, pkg.Type)
	}
	if placement.Scale <= 0 {
		return Scene{}, fmt.Errorf("gauge package %q placement scale must be greater than zero", pkg.ID)
	}
	levelPath := strings.TrimSpace(pkg.Layers["level"])
	if levelPath == "" {
		return Scene{}, fmt.Errorf("gauge package %q bar layer level must not be empty", pkg.ID)
	}
	if len(pkg.Bar.Bounds) != 4 {
		return Scene{}, fmt.Errorf("gauge package %q bar bounds must contain x, y, width, and height", pkg.ID)
	}
	if pkg.ValueMap.Max <= pkg.ValueMap.Min {
		return Scene{}, fmt.Errorf("gauge package %q bar value_map max must be greater than min", pkg.ID)
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
		BarMode:     pkg.Bar.Mode,
		BarAxis:     pkg.Bar.Axis,
		BarOrigin:   pkg.Bar.Origin,
		BarBounds:   cloneInts(pkg.Bar.Bounds),
		Parts:       barUnderlayLayerParts(pkg.Layers),
	}

	if state.Status != sensors.StatusOK {
		scene.Parts = append(scene.Parts, barOverlayLayerParts(pkg.Layers)...)
		return scene, nil
	}

	normalizedPercent, err := barNormalizedPercent(pkg.ValueMap, state.Value)
	if err != nil {
		return Scene{}, fmt.Errorf("gauge package %q: %w", pkg.ID, err)
	}
	windowHeight := pkg.Bar.Bounds[3]
	revealPercent := normalizedPercent
	if revealPercent < 0 {
		revealPercent = 0
	}
	if revealPercent > 100 {
		revealPercent = 100
	}
	revealHeight := int(math.Round((revealPercent / 100) * float64(windowHeight)))
	if revealHeight > 0 {
		boundsX := pkg.Bar.Bounds[0]
		boundsY := pkg.Bar.Bounds[1]
		boundsWidth := pkg.Bar.Bounds[2]
		sourceY := boundsY + (windowHeight - revealHeight)
		scene.Parts = append(scene.Parts, ScenePart{
			Kind:      ScenePartKindBar,
			Layer:     "level",
			AssetPath: levelPath,
			Position:  []int{boundsX, sourceY},
			Source:    []int{boundsX, sourceY},
			Window:    Size{Width: boundsWidth, Height: revealHeight},
		})
	}
	scene.Parts = append(scene.Parts, barOverlayLayerParts(pkg.Layers)...)
	return scene, nil
}

func SegmentedScene(pkg Package, placement Placement, state sensors.SensorState, previous *SegmentedSelection) (Scene, *SegmentedSelection, error) {
	if pkg.Type != TypeSegmented {
		return Scene{}, nil, fmt.Errorf("gauge package %q type %q is not segmented", pkg.ID, pkg.Type)
	}
	if placement.Scale <= 0 {
		return Scene{}, nil, fmt.Errorf("gauge package %q placement scale must be greater than zero", pkg.ID)
	}
	segmentsPath := strings.TrimSpace(pkg.Layers["segments"])
	if segmentsPath == "" {
		return Scene{}, nil, fmt.Errorf("gauge package %q segmented layer segments must not be empty", pkg.ID)
	}
	if len(pkg.Segmented.Images) == 0 {
		return Scene{}, nil, fmt.Errorf("gauge package %q segmented layer segments has no discovered images", pkg.ID)
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
		Parts:       underlayLayerParts(pkg.Layers),
	}

	if state.Status != sensors.StatusOK {
		scene.Parts = append(scene.Parts, overlayLayerParts(pkg.Layers)...)
		return scene, nil, nil
	}

	percent := segmentedNormalizedPercent(state)
	index := segmentedSelectionIndex(pkg.Segmented.Images, percent, previous, segmentedHysteresis(pkg.Segmented))
	var nextSelection *SegmentedSelection
	if index >= 0 {
		selected := pkg.Segmented.Images[index]
		scene.Parts = append(scene.Parts, ScenePart{
			Kind:      ScenePartKindLayer,
			Layer:     "segments",
			AssetPath: selected.Path,
		})
		nextSelection = &SegmentedSelection{Threshold: selected.Threshold, Path: selected.Path}
	}

	scene.Parts = append(scene.Parts, overlayLayerParts(pkg.Layers)...)
	return scene, nextSelection, nil
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

func barNormalizedPercent(valueMap ValueMap, value float64) (float64, error) {
	if valueMap.Max <= valueMap.Min {
		return 0, fmt.Errorf("bar value_map max must be greater than min")
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
	return ((mappedValue - valueMap.Min) / (valueMap.Max - valueMap.Min)) * 100, nil
}

func segmentedNormalizedPercent(state sensors.SensorState) float64 {
	percent := state.Value
	if state.Max > state.Min {
		percent = ((state.Value - state.Min) / (state.Max - state.Min)) * 100
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	return percent
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
	b.WriteString("|bar=")
	b.WriteString(s.BarMode)
	b.WriteString(",")
	b.WriteString(s.BarAxis)
	b.WriteString(",")
	b.WriteString(s.BarOrigin)
	b.WriteString(",")
	b.WriteString(formatIntSlice(s.BarBounds))
	b.WriteString("|movement=")
	b.WriteString(s.Movement)
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
		b.WriteString("#")
		b.WriteString(formatIntSlice(part.Source))
		b.WriteString("#")
		b.WriteString(strconv.Itoa(part.Window.Width))
		b.WriteString("x")
		b.WriteString(strconv.Itoa(part.Window.Height))
		b.WriteString("#")
		b.WriteString(strconv.FormatFloat(part.StripOffset, 'f', -1, 64))
		b.WriteString("#")
		b.WriteString(strconv.FormatBool(part.Wraparound))
		b.WriteString("#")
		b.WriteString(part.Role)
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

func odometerUnderlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"background", "panel", "bezel", "face", "ticks"})
}

func odometerOverlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"glass", "overlay", "foreground"})
}

func indicatorUnderlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"background", "panel", "bezel", "face"})
}

func indicatorOverlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"glass", "overlay", "foreground"})
}

func barUnderlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"background", "panel", "bezel", "face", "ticks"})
}

func barOverlayLayerParts(layers map[string]string) []ScenePart {
	return namedLayerParts(layers, []string{"glass", "overlay", "foreground"})
}

func indicatorStateOn(state sensors.SensorState) bool {
	if state.Status != sensors.StatusOK {
		return false
	}
	if state.TypedValue.Kind == sensors.ValueKindBool && state.TypedValue.Bool != nil {
		return *state.TypedValue.Bool
	}
	return state.Value != 0
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

func odometerDigitPlaces(wheels []OdometerWheel) []int {
	places := make([]int, len(wheels))
	place := 0
	for index := len(wheels) - 1; index >= 0; index-- {
		if odometerWheelRole(wheels[index]) == WheelRoleSubUnit {
			places[index] = -1
			continue
		}
		places[index] = place
		place++
	}
	return places
}

func odometerWheelOffset(movement string, wraparound bool, wheel OdometerWheel, place int, value float64) (float64, error) {
	position, err := odometerWheelPosition(wheel, place, value, wraparound)
	if err != nil {
		return 0, err
	}
	if movement == MovementClick {
		position = math.Floor(position)
	}
	if position < 0 {
		position = 0
	}
	return position * float64(wheel.Window.Height), nil
}

func odometerWheelPosition(wheel OdometerWheel, place int, value float64, wraparound bool) (float64, error) {
	if wheel.Window.Height <= 0 {
		return 0, fmt.Errorf("odometer wheel window height must be positive")
	}
	absolute := math.Abs(value)
	if odometerWheelRole(wheel) == WheelRoleSubUnit {
		if wraparound {
			return absolute * 10, nil
		}
		return math.Mod(absolute*10, 10), nil
	}
	if place < 0 {
		return 0, fmt.Errorf("odometer wheel place must not be negative")
	}
	divisor := math.Pow10(place)
	if wraparound {
		return absolute / divisor, nil
	}
	return math.Mod(absolute/divisor, 10), nil
}

func odometerWraparoundEnabled(pkg Package) bool {
	return pkg.Realism.Wraparound != nil && *pkg.Realism.Wraparound
}

func odometerWheelSource(wheel OdometerWheel, stripOffset float64) (int, int) {
	sourceX := 0
	sourceY := int(math.Floor(stripOffset))
	if len(wheel.Offset) >= 2 {
		sourceX += wheel.Offset[0]
		sourceY += wheel.Offset[1]
	}
	return sourceX, sourceY
}

func odometerWheelRole(wheel OdometerWheel) string {
	if strings.TrimSpace(wheel.Role) == "" {
		return WheelRoleDigit
	}
	return wheel.Role
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

func segmentedHysteresis(segmented Segmented) float64 {
	if segmented.Hysteresis == nil {
		return 25
	}
	return *segmented.Hysteresis
}

func segmentedSelectionIndex(images []SegmentedImage, value float64, previous *SegmentedSelection, hysteresis float64) int {
	if len(images) == 0 {
		return -1
	}
	candidate := segmentedHighestIndex(images, value)
	if previous == nil {
		return candidate
	}

	previousIndex := segmentedImageIndex(images, previous.Threshold)
	if previousIndex < 0 {
		return candidate
	}
	if candidate >= previousIndex {
		return candidate
	}

	current := previousIndex
	for current > 0 {
		lower := current - 1
		release := segmentedReleaseThreshold(images[current].Threshold, images[lower].Threshold, hysteresis)
		if value >= release {
			return current
		}
		current = lower
	}

	if value < float64(images[0].Threshold) {
		return -1
	}
	return 0
}

func segmentedHighestIndex(images []SegmentedImage, value float64) int {
	index := -1
	for i, image := range images {
		if value >= float64(image.Threshold) {
			index = i
			continue
		}
		break
	}
	return index
}

func segmentedImageIndex(images []SegmentedImage, threshold int) int {
	for index, image := range images {
		if image.Threshold == threshold {
			return index
		}
	}
	return -1
}

func segmentedReleaseThreshold(upper int, lower int, hysteresis float64) float64 {
	if upper <= lower {
		return float64(upper)
	}
	return float64(upper) - (float64(upper-lower) * hysteresis / 100)
}
