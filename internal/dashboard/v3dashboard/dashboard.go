package v3dashboard

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	v3gauges "github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

var gaugePackageLoader = v3gauges.LoadPackageWithSearchPaths

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
	dashboards      []Dashboard
	states          map[string]sensors.SensorState
	signatures      map[string]string
	segments        map[string]v3gauges.SegmentedSelection
	movements       map[string]widgetMovementState
	movementPlanner movementPlanner
	clock           func() time.Time
}

const defaultOdometerMovementDuration = 200 * time.Millisecond

type movementPhase string

const (
	movementPhaseStatic      movementPhase = "static"
	movementPhaseValueChange movementPhase = "value_changed"
	movementPhaseMoving      movementPhase = "moving"
	movementPhaseSettled     movementPhase = "settled"
)

type widgetMovementState struct {
	Phase                movementPhase
	Policy               string
	Mode                 string
	PreviousDisplayValue float64
	DisplayValue         float64
	TargetValue          float64
	Duration             time.Duration
	StartedAt            time.Time
	HasValue             bool
	PreviousWheelOffsets []float64
	WheelOffsets         []float64
	TargetWheelOffsets   []float64
}

type movementContext struct {
	DashboardID string
	WidgetID    string
	SensorID    string
	GaugeType   string
	GaugeMode   string
}

type movementPlanner func(movementContext, sensors.SensorState, widgetMovementState) time.Duration

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
	Wraparound  bool
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
		dashboards:      dashboards,
		states:          map[string]sensors.SensorState{},
		signatures:      map[string]string{},
		segments:        map[string]v3gauges.SegmentedSelection{},
		movements:       map[string]widgetMovementState{},
		movementPlanner: defaultMovementPlanner,
		clock:           time.Now,
	}, nil
}

func defaultMovementPlanner(context movementContext, state sensors.SensorState, current widgetMovementState) time.Duration {
	if context.GaugeType == v3gauges.TypeOdometer {
		switch context.GaugeMode {
		case v3gauges.MovementLinear, v3gauges.MovementEaseOut, v3gauges.MovementBell:
			return defaultOdometerMovementDuration
		default:
			return 0
		}
	}
	if current.Policy == "" || current.Policy == v3gauges.MovementPolicyImmediate {
		return 0
	}
	return 0
}

func WithGaugePackageLoader(loader func([]string, string) (v3gauges.Package, error), fn func() error) error {
	if loader == nil {
		return fn()
	}
	previous := gaugePackageLoader
	gaugePackageLoader = loader
	defer func() {
		gaugePackageLoader = previous
	}()
	return fn()
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
	now := eventTime(event, r.now())
	if event.SensorID != "" {
		state := event.State
		if state.ID == "" {
			state.ID = event.SensorID
		}
		r.states[event.SensorID] = state
	}

	scenes, movementChanged, err := r.snapshotAt(now)
	if err != nil {
		return nil, false, err
	}
	changed := movementChanged
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

func eventTime(event sensors.SensorEvent, fallback time.Time) time.Time {
	if !event.Timestamp.IsZero() {
		return event.Timestamp
	}
	if !event.ReadAt.IsZero() {
		return event.ReadAt
	}
	return fallback
}

func (r *Runtime) Snapshot() ([]Scene, error) {
	if r == nil {
		return nil, fmt.Errorf("v3 dashboard runtime is nil")
	}
	scenes, _, err := r.snapshotAt(r.now())
	return scenes, err
}

func (r *Runtime) Tick(now time.Time) ([]Scene, bool, error) {
	if r == nil {
		return nil, false, fmt.Errorf("v3 dashboard runtime is nil")
	}
	if now.IsZero() {
		now = r.now()
	}
	scenes, movementChanged, err := r.snapshotAt(now)
	if err != nil {
		return nil, false, err
	}
	changed := movementChanged
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

func (r *Runtime) HasActiveMovement() bool {
	if r == nil {
		return false
	}
	for _, movement := range r.movements {
		if movementActive(movement) {
			return true
		}
	}
	return false
}

func (r *Runtime) snapshotAt(now time.Time) ([]Scene, bool, error) {
	if r == nil {
		return nil, false, fmt.Errorf("v3 dashboard runtime is nil")
	}
	scenes := make([]Scene, 0, len(r.dashboards))
	movementChanged := false
	for _, dashboard := range r.dashboards {
		scene, dashboardMovementChanged, err := dashboard.render(r.states, r.segments, r.movements, r.movementPlanner, now)
		if err != nil {
			return nil, false, err
		}
		scenes = append(scenes, scene)
		movementChanged = movementChanged || dashboardMovementChanged
	}
	return scenes, movementChanged, nil
}

func (d Dashboard) Validate() error {
	if d.Assets == nil {
		return fmt.Errorf("dashboard %q requires an asset registry", d.ID)
	}
	for _, widget := range d.Config.Widgets {
		if strings.TrimSpace(widget.ID) == "" {
			return fmt.Errorf("dashboard %q has widget with empty id", d.ID)
		}
		if _, _, err := d.renderWidget(widget, map[string]sensors.SensorState{}, nil, nil, nil, time.Time{}); err != nil && !isMissingSensorOnly(err) {
			return err
		}
	}
	return nil
}

func (d Dashboard) Render(states map[string]sensors.SensorState) (Scene, error) {
	scene, _, err := d.render(states, nil, nil, nil, time.Time{})
	return scene, err
}

func (d Dashboard) render(states map[string]sensors.SensorState, segments map[string]v3gauges.SegmentedSelection, movements map[string]widgetMovementState, planner movementPlanner, now time.Time) (Scene, bool, error) {
	scene := Scene{DashboardID: d.ID, Display: d.Config.Display, Size: d.Config.Size}
	movementChanged := false
	for _, configWidget := range d.Config.Widgets {
		widget, widgetMovementChanged, err := d.renderWidget(configWidget, states, segments, movements, planner, now)
		if err != nil {
			return Scene{}, false, err
		}
		scene.Widgets = append(scene.Widgets, widget)
		movementChanged = movementChanged || widgetMovementChanged
	}
	return scene, movementChanged, nil
}

func (d Dashboard) renderWidget(configWidget v3config.WidgetConfig, states map[string]sensors.SensorState, segments map[string]v3gauges.SegmentedSelection, movements map[string]widgetMovementState, planner movementPlanner, now time.Time) (Widget, bool, error) {
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
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.image_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		widget.Status = sensors.StatusOK
		widget.Parts = appendImageSetParts(widget.Parts, set)
	case v3config.WidgetTypeDigitDisplay:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.DigitSet(configWidget.Asset)
		if !ok {
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.digit_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		if state.Status != sensors.StatusOK {
			widget.Parts = appendDigitLayerParts(widget.Parts, set, configWidget.Digits)
			return widget, false, nil
		}
		text := formatValue(configWidget.Format, state.Value)
		parts, err := digitParts(set, text, configWidget.Digits, d.ID, configWidget.ID, configWidget.Asset)
		if err != nil {
			return Widget{}, false, err
		}
		widget.Text = text
		widget.Parts = parts
	case v3config.WidgetTypeBarDisplay:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.BarSet(configWidget.Asset)
		if !ok {
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.bar_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		parts, err := barParts(set, configWidget, state, d.ID)
		if err != nil {
			return Widget{}, false, err
		}
		widget.Parts = parts
	case v3config.WidgetTypeFrameGauge:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.FrameSet(configWidget.Asset)
		if !ok {
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.frame_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		parts, err := frameGaugeParts(set, configWidget, state, d.ID)
		if err != nil {
			return Widget{}, false, err
		}
		widget.Parts = parts
	case v3config.WidgetTypeIndicator:
		state := stateForWidget(configWidget, states)
		widget.Status = state.Status
		widget.Error = state.Error
		set, ok := d.Assets.IndicatorSet(configWidget.Asset)
		if !ok {
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q asset %q must reference assets.indicator_sets", d.ID, configWidget.ID, configWidget.Asset)
		}
		indicatorState := indicatorStateFor(state)
		widget.Parts = appendIndicatorParts(widget.Parts, set, indicatorState)
	case v3config.WidgetTypeGauge:
		pkg, err := gaugePackageLoader(d.Assets.SearchPaths(), configWidget.Gauge)
		if err != nil {
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q gauge %q could not load package: %w", d.ID, configWidget.ID, configWidget.Gauge, err)
		}
		state := stateForSensor(pkg.Sensor, states)
		context := movementContext{
			DashboardID: d.ID,
			WidgetID:    configWidget.ID,
			SensorID:    pkg.Sensor,
			GaugeType:   pkg.Type,
			GaugeMode:   pkg.Odometer.Movement,
		}
		var movementChanged bool
		var movement widgetMovementState
		if pkg.Type == v3gauges.TypeOdometer {
			state, movement, movementChanged, err = resolveOdometerMovementState(movements, movementKey(d.ID, configWidget.ID), context, pkg, state, planner, now)
			if err != nil {
				return Widget{}, false, fmt.Errorf("dashboard %q widget %q: %w", d.ID, configWidget.ID, err)
			}
		} else {
			state, movementChanged = resolveMovementState(movements, movementKey(d.ID, configWidget.ID), context, state, pkg.Realism.MovementPolicy, planner, now)
		}
		var gaugeScene v3gauges.Scene
		switch pkg.Type {
		case v3gauges.TypeNumeric:
			gaugeScene, err = v3gauges.NumericScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
		case v3gauges.TypeRadial:
			gaugeScene, err = v3gauges.RadialScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
		case v3gauges.TypeOdometer:
			if movementActive(movement) && len(movement.WheelOffsets) == len(pkg.Odometer.Wheels) {
				gaugeScene, err = v3gauges.OdometerSceneWithWheelOffsets(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state, movement.WheelOffsets)
			} else {
				gaugeScene, err = v3gauges.OdometerScene(pkg, v3gauges.Placement{Position: configWidget.Position, Scale: configWidget.Scale}, state)
			}
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
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q gauge package type %q is not supported by dashboard scene runtime", d.ID, configWidget.ID, pkg.Type)
		}
		if err != nil {
			return Widget{}, false, fmt.Errorf("dashboard %q widget %q: %w", d.ID, configWidget.ID, err)
		}
		applyGaugeScene(&widget, gaugeScene)
		return widget, movementChanged, nil
	default:
		return Widget{}, false, fmt.Errorf("dashboard %q widget %q type %q is not supported", d.ID, configWidget.ID, configWidget.Type)
	}

	return widget, false, nil
}

func movementKey(dashboardID string, widgetID string) string {
	return dashboardID + "|" + widgetID
}

func resolveMovementState(movements map[string]widgetMovementState, key string, context movementContext, source sensors.SensorState, policy string, planner movementPlanner, now time.Time) (sensors.SensorState, bool) {
	if movements == nil {
		return source, false
	}
	if strings.TrimSpace(policy) == "" {
		policy = v3gauges.MovementPolicyImmediate
	}
	if source.Status != sensors.StatusOK {
		previous, hadMovement := movements[key]
		if hadMovement {
			delete(movements, key)
		}
		return source, hadMovement && movementActive(previous)
	}
	movement := movements[key]
	previous := movement
	if !movement.HasValue {
		movement = widgetMovementState{
			Phase:                movementPhaseStatic,
			Policy:               policy,
			Mode:                 context.GaugeMode,
			PreviousDisplayValue: source.Value,
			DisplayValue:         source.Value,
			TargetValue:          source.Value,
			HasValue:             true,
		}
	} else if source.Value != movement.TargetValue {
		if movementActive(movement) {
			movement = advanceMovementState(movement, now)
		}
		movement.Policy = policy
		movement.PreviousDisplayValue = movement.DisplayValue
		movement.TargetValue = source.Value
		duration := time.Duration(0)
		if planner != nil {
			duration = planner(context, source, movement)
		}
		if duration <= 0 || movement.DisplayValue == source.Value || movement.Policy == v3gauges.MovementPolicyImmediate {
			movement.DisplayValue = source.Value
			movement.Phase = movementPhaseStatic
			movement.Duration = 0
			movement.StartedAt = time.Time{}
		} else {
			movement.Duration = duration
			movement.StartedAt = now
			movement.Phase = movementPhaseValueChange
		}
	}
	movement = advanceMovementState(movement, now)
	movements[key] = movement
	source.Value = movement.DisplayValue
	source.TypedValue = sensors.NewNumericValue(source.Value, source.Unit)
	return source, movementActive(movement) != movementActive(previous)
}

func resolveOdometerMovementState(movements map[string]widgetMovementState, key string, context movementContext, pkg v3gauges.Package, source sensors.SensorState, planner movementPlanner, now time.Time) (sensors.SensorState, widgetMovementState, bool, error) {
	if movements == nil {
		return source, widgetMovementState{}, false, nil
	}
	if source.Status != sensors.StatusOK {
		previous, hadMovement := movements[key]
		if hadMovement {
			delete(movements, key)
		}
		return source, widgetMovementState{}, hadMovement && movementActive(previous), nil
	}

	targetOffsets, err := v3gauges.OdometerWheelStripOffsets(pkg, source.Value)
	if err != nil {
		return source, widgetMovementState{}, false, err
	}

	movement := movements[key]
	previous := movement
	if !movement.HasValue {
		movement = widgetMovementState{
			Phase:                movementPhaseStatic,
			Mode:                 pkg.Odometer.Movement,
			PreviousDisplayValue: source.Value,
			DisplayValue:         source.Value,
			TargetValue:          source.Value,
			HasValue:             true,
			PreviousWheelOffsets: cloneFloat64s(targetOffsets),
			WheelOffsets:         cloneFloat64s(targetOffsets),
			TargetWheelOffsets:   cloneFloat64s(targetOffsets),
		}
	} else if source.Value != movement.TargetValue {
		if movementActive(movement) {
			movement = advanceMovementState(movement, now)
		}
		movement.Mode = pkg.Odometer.Movement
		movement.PreviousDisplayValue = movement.DisplayValue
		movement.DisplayValue = source.Value
		movement.TargetValue = source.Value
		movement.PreviousWheelOffsets = cloneFloat64s(currentWheelOffsets(movement))
		movement.TargetWheelOffsets = cloneFloat64s(targetOffsets)
		duration := time.Duration(0)
		if planner != nil {
			duration = planner(context, source, movement)
		}
		if duration <= 0 || odometerWheelOffsetsEqual(movement.PreviousWheelOffsets, movement.TargetWheelOffsets) {
			movement.WheelOffsets = cloneFloat64s(targetOffsets)
			movement.Phase = movementPhaseStatic
			movement.Duration = 0
			movement.StartedAt = time.Time{}
		} else {
			movement.Duration = duration
			movement.StartedAt = now
			movement.Phase = movementPhaseValueChange
			movement.WheelOffsets = cloneFloat64s(movement.PreviousWheelOffsets)
		}
	}
	movement = advanceMovementState(movement, now)
	movements[key] = movement
	return source, movement, movementActive(movement) != movementActive(previous), nil
}

func advanceMovementState(movement widgetMovementState, now time.Time) widgetMovementState {
	switch movement.Phase {
	case movementPhaseValueChange:
		movement.Phase = movementPhaseMoving
	case movementPhaseMoving:
		if movement.Duration <= 0 || now.IsZero() || !movement.HasValue {
			movement.DisplayValue = movement.TargetValue
			if odometerMovementInFlight(movement) {
				movement.WheelOffsets = cloneFloat64s(movement.TargetWheelOffsets)
			}
			movement.Phase = movementPhaseSettled
			return movement
		}
		elapsed := now.Sub(movement.StartedAt)
		if elapsed < 0 {
			elapsed = 0
		}
		progress := float64(elapsed) / float64(movement.Duration)
		if progress >= 1 {
			movement.DisplayValue = movement.TargetValue
			if odometerMovementInFlight(movement) {
				movement.WheelOffsets = cloneFloat64s(movement.TargetWheelOffsets)
			}
			movement.Phase = movementPhaseSettled
			return movement
		}
		if progress < 0 {
			progress = 0
		}
		if odometerMovementInFlight(movement) {
			progress = applyOdometerMovementCurve(progress, movement.Mode)
			movement.WheelOffsets = interpolateWheelOffsets(movement.PreviousWheelOffsets, movement.TargetWheelOffsets, progress)
			movement.DisplayValue = movement.TargetValue
		} else {
			progress = applyMovementPolicy(progress, movement.Policy)
			movement.DisplayValue = movement.PreviousDisplayValue + ((movement.TargetValue - movement.PreviousDisplayValue) * progress)
		}
	case movementPhaseSettled:
		movement.Phase = movementPhaseStatic
		movement.DisplayValue = movement.TargetValue
		if odometerMovementInFlight(movement) {
			movement.WheelOffsets = cloneFloat64s(movement.TargetWheelOffsets)
		}
	}
	return movement
}

func applyMovementPolicy(progress float64, policy string) float64 {
	switch policy {
	case "", v3gauges.MovementPolicyImmediate:
		return 1
	case v3gauges.MovementPolicyEaseOut:
		return 1 - ((1 - progress) * (1 - progress))
	default:
		return progress
	}
}

func applyOdometerMovementCurve(progress float64, mode string) float64 {
	switch mode {
	case v3gauges.MovementEaseOut:
		return 1 - ((1 - progress) * (1 - progress))
	case v3gauges.MovementBell:
		return progress * progress * (3 - (2 * progress))
	default:
		return progress
	}
}

func odometerMovementInFlight(movement widgetMovementState) bool {
	return len(movement.PreviousWheelOffsets) > 0 && len(movement.PreviousWheelOffsets) == len(movement.TargetWheelOffsets)
}

func currentWheelOffsets(movement widgetMovementState) []float64 {
	if len(movement.WheelOffsets) > 0 {
		return movement.WheelOffsets
	}
	if len(movement.TargetWheelOffsets) > 0 {
		return movement.TargetWheelOffsets
	}
	return nil
}

func interpolateWheelOffsets(previous []float64, target []float64, progress float64) []float64 {
	interpolated := make([]float64, len(previous))
	for index := range previous {
		interpolated[index] = previous[index] + ((target[index] - previous[index]) * progress)
	}
	return interpolated
}

func odometerWheelOffsetsEqual(left []float64, right []float64) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if math.Abs(left[index]-right[index]) > 0.001 {
			return false
		}
	}
	return true
}

func cloneFloat64s(values []float64) []float64 {
	if values == nil {
		return nil
	}
	return append([]float64(nil), values...)
}

func movementActive(movement widgetMovementState) bool {
	switch movement.Phase {
	case movementPhaseValueChange, movementPhaseMoving, movementPhaseSettled:
		return true
	default:
		return false
	}
}

func (r *Runtime) now() time.Time {
	if r != nil && r.clock != nil {
		return r.clock()
	}
	return time.Now()
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
			Wraparound:  scenePart.Wraparound,
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
			b.WriteString(strconv.FormatBool(part.Wraparound))
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
