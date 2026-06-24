package v3dashboard

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	v3gauges "github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

const (
	PartKindBackground   = "background"
	PartKindImage        = "image"
	PartKindCharacter    = "character"
	PartKindDecimalPoint = "decimal_point"
	PartKindState        = "state"
	PartKindCell         = "cell"
	PartKindFrame        = "frame"
	PartKindForeground   = "foreground"
	PartKindLayer        = "layer"
	PartKindNeedle       = "needle"
	PartKindBar          = "bar"
	PartKindWheelStrip   = "wheel_strip"
)

// Runtime owns selected v3 dashboard scene state. It consumes sensor state/events
// produced by the sensor runtime; it never reads endpoint or OBD code directly.
type Runtime struct {
	dashboards []Dashboard
	states     map[string]sensors.SensorState
	signatures map[string]string
	segments   map[string]v3gauges.SegmentedSelection
}

type Dashboard struct {
	ID     string
	Config v3config.DashboardConfig
	Assets *v3assets.Registry
}

type Scene struct {
	DashboardID string
	Display     string
	Size        v3config.SizeConfig
	Widgets     []Widget
}

type Widget struct {
	ID                  string
	Type                string
	SensorID            string
	AssetID             string
	GaugeID             string
	GaugePath           string
	GaugeDigitPositions [][]int
	GaugeFacePivot      v3gauges.Point
	GaugeNeedlePivot    v3gauges.Point
	GaugeAngle          float64
	GaugeMovement       string
	GaugeBarMode        string
	GaugeBarAxis        string
	GaugeBarOrigin      string
	GaugeBarBounds      []int
	Scale               float64
	Position            []int
	Status              string
	Text                string
	Parts               []Part
	Error               string
}

type Part struct {
	Kind        string
	Layer       string
	AssetPath   string
	Slot        int
	Character   string
	Position    []int
	State       string
	Cell        string
	Frame       int
	Angle       float64
	FacePivot   v3gauges.Point
	NeedlePivot v3gauges.Point
	Source      []int
	Window      v3gauges.Size
	StripOffset float64
	Role        string
}

// NewRuntime builds the selected-dashboard runtime from an already resolved
// RuntimePlan. Unselected dashboards are not present in the plan and stay inert.
func NewRuntime(plan v3config.RuntimePlan, registry *v3assets.Registry) (*Runtime, error) {
	if registry == nil {
		return nil, fmt.Errorf("v3 dashboard runtime requires an asset registry")
	}
	dashboards := make([]Dashboard, 0, len(plan.Dashboards))
	for _, selected := range plan.Dashboards {
		dashboard := Dashboard{ID: selected.ID, Config: selected.Config, Assets: registry}
		if err := dashboard.Validate(); err != nil {
			return nil, err
		}
		dashboards = append(dashboards, dashboard)
	}
	return &Runtime{
		dashboards: dashboards,
		states:     map[string]sensors.SensorState{},
		signatures: map[string]string{},
		segments:   map[string]v3gauges.SegmentedSelection{},
	}, nil
}

func (r *Runtime) DashboardCount() int {
	if r == nil {
		return 0
	}
	return len(r.dashboards)
}

func (r *Runtime) SetState(state sensors.SensorState) {
	if r == nil || state.ID == "" {
		return
	}
	r.states[state.ID] = state
}

func (r *Runtime) ApplyEvent(event sensors.SensorEvent) ([]Scene, bool, error) {
	if r == nil {
		return nil, false, fmt.Errorf("v3 dashboard runtime is nil")
	}
	if event.SensorID != "" {
		state := event.State
		if state.ID == "" {
			state.ID = event.SensorID
		}
		r.states[event.SensorID] = state
	}

	scenes, err := r.Snapshot()
	if err != nil {
		return nil, false, err
	}
	changed := false
	for _, scene := range scenes {
		signature := sceneSignature(scene)
		if r.signatures[scene.DashboardID] != signature {
			r.signatures[scene.DashboardID] = signature
			changed = true
		}
	}
	if !changed {
		return nil, false, nil
	}
	return scenes, true, nil
}

func (r *Runtime) Snapshot() ([]Scene, error) {
	if r == nil {
		return nil, fmt.Errorf("v3 dashboard runtime is nil")
	}
	scenes := make([]Scene, 0, len(r.dashboards))
	for _, dashboard := range r.dashboards {
		scene, err := dashboard.render(r.states, r.segments)
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, scene)
	}
	return scenes, nil
}

func (d Dashboard) Validate() error {
	if d.Assets == nil {
		return fmt.Errorf("dashboard %q requires an asset registry", d.ID)
	}
	for _, widget := range d.Config.Widgets {
		if strings.TrimSpace(widget.ID) == "" {
			return fmt.Errorf("dashboard %q has widget with empty id", d.ID)
		}
		if _, err := d.renderWidget(widget, map[string]sensors.SensorState{}, nil); err != nil && !isMissingSensorOnly(err) {
			return err
		}
	}
	return nil
}

func (d Dashboard) Render(states map[string]sensors.SensorState) (Scene, error) {
	return d.render(states, nil)
}

func (d Dashboard) render(states map[string]sensors.SensorState, segments map[string]v3gauges.SegmentedSelection) (Scene, error) {
	scene := Scene{DashboardID: d.ID, Display: d.Config.Display, Size: d.Config.Size}
	for _, configWidget := range d.Config.Widgets {
		widget, err := d.renderWidget(configWidget, states, segments)
		if err != nil {
			return Scene{}, err
		}
		scene.Widgets = append(scene.Widgets, widget)
	}
	return scene, nil
}

func (d Dashboard) renderWidget(configWidget v3config.WidgetConfig, states map[string]sensors.SensorState, segments map[string]v3gauges.SegmentedSelection) (Widget, error) {
	widget := Widget{
		ID:       configWidget.ID,
		Type:     configWidget.Type,
		SensorID: configWidget.Sensor,
		AssetID:  configWidget.Asset,
		Position: append([]int(nil), configWidget.Position...),
	}

	switch configWidget.Type {
	case v3config.WidgetTypeImage:
		set, ok := d.Assets.ImageSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.image_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		widget.Status = sensors.StatusOK
		widget.Parts = appendImageSetParts(widget.Parts, set)
	case v3config.WidgetTypeDigitDisplay:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.DigitSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.digit_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		if state.Status != sensors.StatusOK {
			widget.Parts = appendDigitLayerParts(widget.Parts, set, configWidget.Digits)
			return widget, nil
		}
		text := formatValue(configWidget.Format, state.Value)
		parts, err := digitParts(set, text, configWidget.Digits, d.ID, configWidget.ID, configWidget.Asset)
		if err != nil {
			return Widget{}, err
		}
		widget.Text = text
		widget.Parts = parts
	case v3config.WidgetTypeBarDisplay:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.BarSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.bar_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		parts, err := barParts(set, configWidget, state, d.ID)
		if err != nil {
			return Widget{}, err
		}
		widget.Parts = parts
	case v3config.WidgetTypeFrameGauge:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.FrameSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.frame_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		parts, err := frameGaugeParts(set, configWidget, state, d.ID)
		if err != nil {
			return Widget{}, err
		}
		widget.Parts = parts
	case v3config.WidgetTypeIndicator:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.IndicatorSet(configWidget.Asset)
		if !ok {
			return Widget{}, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.indicator_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		indicatorState := indicatorStateFor(state)
		widget.Parts = appendIndicatorParts(widget.Parts, set, indicatorState)
	case v3config.WidgetTypeGauge:
  pkg, err := v3gauges.LoadPackageWithSearchPaths(d.Assets.SearchPaths(), configWidget.Gauge)
		if err != nil {
			return Widget{}, fmt.Errorf("dashboard %q widget %q gauge %q could not load package: %w", d.ID, configWidget.ID, configWidget.Gauge, err)
		}
		state := stateForSensor(pkg.Sensor, states)
		var gaugeScene v3gauges.Scene
		switch pkg.Type {
		case v3gauges.TypeNumeric:
			gaugeScene, err = v3gauges.NumericScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
		case v3gauges.TypeRadial:
			gaugeScene, err = v3gauges.RadialScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
		case v3gauges.TypeOdometer:
			gaugeScene, err = v3gauges.OdometerScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
		case v3gauges.TypeIndicator:
			gaugeScene, err = v3gauges.IndicatorScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
		case v3gauges.TypeBar:
			gaugeScene, err = v3gauges.BarScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
		case v3gauges.TypeSegmented:
			var previous *v3gauges.SegmentedSelection
			selectionKey := segmentSelectionKey(d.ID, configWidget.ID, pkg.Path)
			if segments != nil {
				if selection, ok := segments[selectionKey]; ok {
					previous = &selection
				}
			}
			var nextSelection *v3gauges.SegmentedSelection
			gaugeScene, nextSelection, err = v3gauges.SegmentedScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state, previous)
			if err == nil && segments != nil {
				if nextSelection == nil {
					delete(segments, selectionKey)
				} else {
					segments[selectionKey] = *nextSelection
				}
			}
		default:
			return Widget{}, fmt.Errorf("dashboard %q widget %q gauge package type %q is not supported by dashboard scene runtime", d.ID, configWidget.ID, pkg.Type)
		}
		if err != nil {
			return Widget{}, fmt.Errorf("dashboard %q widget %q: %w", d.ID, configWidget.ID, err)
		}
		applyGaugeScene(&widget, gaugeScene)
	default:
		return Widget{}, fmt.Errorf("dashboard %q widget %q type %q is not supported", d.ID, configWidget.ID, configWidget.Type)
	}

	return widget, nil
}

func applyGaugeScene(widget *Widget, scene v3gauges.Scene) {
	widget.SensorID = scene.SensorID
	widget.GaugeID = scene.PackageID
	widget.GaugePath = scene.PackagePath
	widget.GaugeDigitPositions = cloneIntSlices(scene.DigitPositions)
	widget.GaugeFacePivot = scene.FacePivot
	widget.GaugeNeedlePivot = scene.NeedlePivot
	widget.GaugeAngle = scene.Angle
	widget.GaugeMovement = scene.Movement
	widget.GaugeBarMode = scene.BarMode
	widget.GaugeBarAxis = scene.BarAxis
	widget.GaugeBarOrigin = scene.BarOrigin
	widget.GaugeBarBounds = append([]int(nil), scene.BarBounds...)
	widget.Scale = scene.Scale
	widget.Status = scene.Status
	widget.Error = scene.Error
	widget.Text = scene.Text
	widget.Parts = gaugeSceneParts(scene)
}

func segmentSelectionKey(dashboardID string, widgetID string, packagePath string) string {
	return dashboardID + "|" + widgetID + "|" + packagePath
}

func stateForWidget(widget v3config.WidgetConfig, states map[string]sensors.SensorState) sensors.SensorState {
	return stateForSensor(widget.Sensor, states)
}

func stateForSensor(sensorID string, states map[string]sensors.SensorState) sensors.SensorState {
	state, ok := states[sensorID]
	if !ok {
		return sensors.SensorState{ID: sensorID, Status: sensors.StatusMissingUnsupported}
	}
	if state.ID == "" {
		state.ID = sensorID
	}
	if state.Status == "" {
		state.Status = sensors.StatusUnknown
	}
	return state
}

func formatValue(format string, value float64) string {
	if strings.TrimSpace(format) == "" {
		return fmt.Sprintf("%.0f", value)
	}
	return fmt.Sprintf(format, value)
}

func digitParts(set v3assets.DigitSet, text string, slots int, dashboardID, widgetID, assetID string) ([]Part, error) {
	if slots <= 0 {
		return nil, fmt.Errorf("dashboard %q widget %q digits must be greater than zero", dashboardID, widgetID)
	}
	characters, decimalSlots, err := splitDigitText(text, slots)
	if err != nil {
		return nil, fmt.Errorf("dashboard %q widget %q: %w", dashboardID, widgetID, err)
	}
	decimalBySlot := map[int]bool{}
	for _, slot := range decimalSlots {
		decimalBySlot[slot] = true
	}

	parts := make([]Part, 0, len(characters)*4)
	for slot, ch := range characters {
		if set.Background != nil {
			parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path, Slot: slot})
		}
		if ch != " " {
			asset, ok := set.Characters[ch]
			if !ok {
				return nil, fmt.Errorf("dashboard %q widget %q digit set %q has no character asset for %q", dashboardID, widgetID, assetID, ch)
			}
			parts = append(parts, Part{Kind: PartKindCharacter, AssetPath: asset.Path, Slot: slot, Character: ch})
		}
		if decimalBySlot[slot] {
			if set.DecimalPoint == nil {
				return nil, fmt.Errorf("dashboard %q widget %q digit set %q requires decimal_point", dashboardID, widgetID, assetID)
			}
			parts = append(parts, Part{Kind: PartKindDecimalPoint, AssetPath: set.DecimalPoint.Path, Slot: slot})
		}
		if set.Foreground != nil {
			parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path, Slot: slot})
		}
	}
	return parts, nil
}

func splitDigitText(text string, slots int) ([]string, []int, error) {
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
		return nil, nil, fmt.Errorf("formatted output %q needs %d character slots, dashboard config allows %d", text, len(characters), slots)
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

func appendDigitLayerParts(parts []Part, set v3assets.DigitSet, slots int) []Part {
	if slots <= 0 {
		return parts
	}
	for slot := 0; slot < slots; slot++ {
		if set.Background != nil {
			parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path, Slot: slot})
		}
		if set.Foreground != nil {
			parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path, Slot: slot})
		}
	}
	return parts
}

func gaugeSceneParts(scene v3gauges.Scene) []Part {
	parts := make([]Part, 0, len(scene.Parts))
	for _, scenePart := range scene.Parts {
		parts = append(parts, Part{
			Kind:        scenePart.Kind,
			Layer:       scenePart.Layer,
			AssetPath:   scenePart.AssetPath,
			Slot:        scenePart.Slot,
			Character:   scenePart.Character,
			Position:    append([]int(nil), scenePart.Position...),
			Angle:       scenePart.Angle,
			FacePivot:   scenePart.FacePivot,
			NeedlePivot: scenePart.NeedlePivot,
			Source:      append([]int(nil), scenePart.Source...),
			Window:      scenePart.Window,
			StripOffset: scenePart.StripOffset,
			Role:        scenePart.Role,
		})
	}
	return parts
}

func barParts(set v3assets.BarSet, widget v3config.WidgetConfig, state sensors.SensorState, dashboardID string) ([]Part, error) {
	if widget.Cells <= 0 {
		return nil, fmt.Errorf("dashboard %q widget %q cells must be greater than zero", dashboardID, widget.ID)
	}

	parts := []Part{}
	if set.Background != nil {
		parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path})
	}
	if state.Status != sensors.StatusOK {
		if set.Foreground != nil {
			parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path})
		}
		return parts, nil
	}

	if _, ok := set.Cells[v3assets.IndicatorStateOff]; !ok {
		return nil, fmt.Errorf("dashboard %q widget %q bar set %q requires off cell", dashboardID, widget.ID, widget.Asset)
	}
	filled, err := filledBarCells(widget, state.Value)
	if err != nil {
		return nil, fmt.Errorf("dashboard %q widget %q: %w", dashboardID, widget.ID, err)
	}
	cellName := v3assets.IndicatorStateOff
	if filled > 0 {
		cellName, err = barCellNameForValue(widget, set, state.Value)
		if err != nil {
			return nil, fmt.Errorf("dashboard %q widget %q: %w", dashboardID, widget.ID, err)
		}
	}

	for slot := 0; slot < widget.Cells; slot++ {
		name := v3assets.IndicatorStateOff
		if isFilledBarSlot(slot, filled, widget.Cells, widget.Reverse) {
			name = cellName
		}
		asset, ok := set.Cells[name]
		if !ok {
			return nil, fmt.Errorf("dashboard %q widget %q bar set %q has no cell asset for %q", dashboardID, widget.ID, widget.Asset, name)
		}
		parts = append(parts, Part{Kind: PartKindCell, AssetPath: asset.Path, Slot: slot, Cell: name})
	}
	if set.Foreground != nil {
		parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path})
	}
	return parts, nil
}

func filledBarCells(widget v3config.WidgetConfig, value float64) (int, error) {
	min, max := widgetRange(widget, float64(widget.Cells))
	if max <= min {
		return 0, fmt.Errorf("min must be less than max")
	}
	if value <= min {
		return 0, nil
	}
	if value >= max {
		return widget.Cells, nil
	}
	filled := int(math.Ceil(((value - min) / (max - min)) * float64(widget.Cells)))
	if filled < 0 {
		return 0, nil
	}
	if filled > widget.Cells {
		return widget.Cells, nil
	}
	return filled, nil
}

func widgetRange(widget v3config.WidgetConfig, defaultMax float64) (float64, float64) {
	min := 0.0
	max := defaultMax
	if widget.Min != nil {
		min = *widget.Min
	}
	if widget.Max != nil {
		max = *widget.Max
	}
	return min, max
}

func isFilledBarSlot(slot, filled, cells int, reverse bool) bool {
	if filled <= 0 {
		return false
	}
	if reverse {
		return slot >= cells-filled
	}
	return slot < filled
}

func barCellNameForValue(widget v3config.WidgetConfig, set v3assets.BarSet, value float64) (string, error) {
	if len(widget.Zones) == 0 {
		if _, ok := set.Cells[v3assets.IndicatorStateOn]; !ok {
			return "", fmt.Errorf("bar set %q requires on cell when zones are omitted", widget.Asset)
		}
		return v3assets.IndicatorStateOn, nil
	}

	last := widget.Zones[len(widget.Zones)-1].Cell
	for _, zone := range widget.Zones {
		if value <= zone.UpTo {
			last = zone.Cell
			break
		}
	}
	if _, ok := set.Cells[last]; !ok {
		return "", fmt.Errorf("bar set %q has no cell asset for zone cell %q", widget.Asset, last)
	}
	return last, nil
}

func frameGaugeParts(set v3assets.FrameSet, widget v3config.WidgetConfig, state sensors.SensorState, dashboardID string) ([]Part, error) {
	parts := []Part{}
	if set.Background != nil {
		parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path})
	}
	if state.Status == sensors.StatusOK {
		frame, err := frameForValue(set, widget, state.Value)
		if err != nil {
			return nil, fmt.Errorf("dashboard %q widget %q: %w", dashboardID, widget.ID, err)
		}
		asset, ok := set.Frames[frame]
		if !ok {
			return nil, fmt.Errorf("dashboard %q widget %q frame set %q has no frame %d", dashboardID, widget.ID, widget.Asset, frame)
		}
		parts = append(parts, Part{Kind: PartKindFrame, AssetPath: asset.Path, Frame: frame})
	}
	if set.Foreground != nil {
		parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path})
	}
	return parts, nil
}

func frameForValue(set v3assets.FrameSet, widget v3config.WidgetConfig, value float64) (int, error) {
	if set.Last < set.First {
		return 0, fmt.Errorf("frame set %q has invalid range", widget.Asset)
	}
	if set.Last == set.First {
		return set.First, nil
	}
	min, max := widgetRange(widget, 1)
	if max <= min {
		return 0, fmt.Errorf("min must be less than max")
	}
	if value <= min {
		return set.First, nil
	}
	if value >= max {
		return set.Last, nil
	}
	span := float64(set.Last - set.First)
	offset := int(math.Round(((value - min) / (max - min)) * span))
	frame := set.First + offset
	if frame < set.First {
		return set.First, nil
	}
	if frame > set.Last {
		return set.Last, nil
	}
	return frame, nil
}

func indicatorStateFor(state sensors.SensorState) string {
	if state.Status != sensors.StatusOK {
		return v3assets.IndicatorStateUnknown
	}
	if state.Value != 0 {
		return v3assets.IndicatorStateOn
	}
	return v3assets.IndicatorStateOff
}

func appendImageSetParts(parts []Part, set v3assets.ImageSet) []Part {
	if set.Background != nil {
		parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path})
	}
	if set.Image != nil {
		parts = append(parts, Part{Kind: PartKindImage, AssetPath: set.Image.Path})
	}
	if set.Foreground != nil {
		parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path})
	}
	return parts
}

func appendIndicatorParts(parts []Part, set v3assets.IndicatorSet, state string) []Part {
	if set.Background != nil {
		parts = append(parts, Part{Kind: PartKindBackground, AssetPath: set.Background.Path})
	}
	asset := set.States[state]
	parts = append(parts, Part{Kind: PartKindState, AssetPath: asset.Path, State: state})
	if set.Foreground != nil {
		parts = append(parts, Part{Kind: PartKindForeground, AssetPath: set.Foreground.Path})
	}
	return parts
}

func isMissingSensorOnly(err error) bool {
	return err != nil && strings.Contains(err.Error(), string(sensors.StatusMissingUnsupported))
}

func sceneSignature(scene Scene) string {
	var b strings.Builder
	b.WriteString(scene.DashboardID)
	b.WriteString("|")
	for _, widget := range scene.Widgets {
		b.WriteString(widget.ID)
		b.WriteString(":")
		b.WriteString(widget.Type)
		b.WriteString(":")
		b.WriteString(widget.GaugeID)
		b.WriteString(":")
		b.WriteString(widget.GaugePath)
		b.WriteString(":")
		b.WriteString(formatPartPositions(widget.GaugeDigitPositions))
		b.WriteString(":")
		b.WriteString(formatGaugePoint(widget.GaugeFacePivot))
		b.WriteString(":")
		b.WriteString(formatGaugePoint(widget.GaugeNeedlePivot))
		b.WriteString(":")
		b.WriteString(strconv.FormatFloat(widget.GaugeAngle, 'f', -1, 64))
		b.WriteString(":")
		b.WriteString(widget.GaugeMovement)
		b.WriteString(":")
		b.WriteString(widget.GaugeBarMode)
		b.WriteString(":")
		b.WriteString(widget.GaugeBarAxis)
		b.WriteString(":")
		b.WriteString(widget.GaugeBarOrigin)
		b.WriteString(":")
		b.WriteString(formatPartPosition(widget.GaugeBarBounds))
		b.WriteString(":")
		b.WriteString(strconv.FormatFloat(widget.Scale, 'f', -1, 64))
		b.WriteString(":")
		b.WriteString(widget.Status)
		b.WriteString(":")
		b.WriteString(widget.Text)
		b.WriteString(":")
		for _, part := range widget.Parts {
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
			b.WriteString(part.State)
			b.WriteString("#")
			b.WriteString(part.Cell)
			b.WriteString("#")
			b.WriteString(strconv.Itoa(part.Frame))
			b.WriteString("#")
			b.WriteString(formatPartPosition(part.Position))
			b.WriteString("#")
			b.WriteString(strconv.FormatFloat(part.Angle, 'f', -1, 64))
			b.WriteString("#")
			b.WriteString(formatGaugePoint(part.FacePivot))
			b.WriteString("#")
			b.WriteString(formatGaugePoint(part.NeedlePivot))
			b.WriteString("#")
			b.WriteString(formatPartPosition(part.Source))
			b.WriteString("#")
			b.WriteString(strconv.Itoa(part.Window.Width))
			b.WriteString("x")
			b.WriteString(strconv.Itoa(part.Window.Height))
			b.WriteString("#")
			b.WriteString(strconv.FormatFloat(part.StripOffset, 'f', -1, 64))
			b.WriteString("#")
			b.WriteString(part.Role)
			b.WriteString(";")
		}
		b.WriteString("|")
	}
	return b.String()
}

func formatPartPosition(position []int) string {
	if len(position) == 0 {
		return ""
	}
	parts := make([]string, len(position))
	for index, value := range position {
		parts[index] = strconv.Itoa(value)
	}
	return strings.Join(parts, ",")
}

func formatPartPositions(positions [][]int) string {
	if len(positions) == 0 {
		return ""
	}
	parts := make([]string, len(positions))
	for index, position := range positions {
		parts[index] = formatPartPosition(position)
	}
	return strings.Join(parts, ";")
}

func formatGaugePoint(point v3gauges.Point) string {
	if point == (v3gauges.Point{}) {
		return ""
	}
	return strconv.FormatFloat(point.X, 'f', -1, 64) + "," + strconv.FormatFloat(point.Y, 'f', -1, 64)
}

func cloneIntSlices(values [][]int) [][]int {
	if values == nil {
		return nil
	}
	cloned := make([][]int, len(values))
	for index, value := range values {
		cloned[index] = append([]int(nil), value...)
	}
	return cloned
}

func sortedWidgetIDs(widgets []Widget) []string {
	ids := make([]string, 0, len(widgets))
	for _, widget := range widgets {
		ids = append(ids, widget.ID)
	}
	sort.Strings(ids)
	return ids
}
