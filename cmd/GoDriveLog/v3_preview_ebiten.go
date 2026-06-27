//go:build !fyne_legacy

package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	v3ebitenadapter "github.com/MickMake/GoDriveLog/internal/dashboard/adapter/ebiten"
	v3gauges "github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
	"github.com/MickMake/GoDriveLog/internal/sensors"
	ebitenui "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type previewGaugeCandidate struct {
	DashboardID string
	Widget      v3config.WidgetConfig
}

type previewGaugeSpec struct {
	Runtime       *v3dashboard.Runtime
	Title         string
	SensorID      string
	Unit          string
	Min           float64
	Max           float64
	InitialValue  float64
	Step          float64
	FineStep      float64
	CoarseStep    float64
	SelectedGauge string
}

type gaugePreviewController struct {
	runtime *v3dashboard.Runtime
	adapter *v3ebitenadapter.Adapter
	stop    func()

	sensorID   string
	unit       string
	min        float64
	max        float64
	midpoint   float64
	step       float64
	fineStep   float64
	coarseStep float64

	currentValue  float64
	lastFrom      float64
	lastTo        float64
	hasLast       bool
	replayPending bool
}

func runV3EbitenPreviewCommand(options dashboardPreviewOptions) error {
	spec, err := buildGaugePreviewSpec(options)
	if err != nil {
		return err
	}

	ctx, stop := newV3Context(0)
	defer stop()

	adapter, err := v3ebitenadapter.New(".", 800, 480)
	if err != nil {
		return err
	}
	controller := &gaugePreviewController{
		runtime:      spec.Runtime,
		adapter:      adapter,
		stop:         stop,
		sensorID:     spec.SensorID,
		unit:         spec.Unit,
		min:          spec.Min,
		max:          spec.Max,
		midpoint:     midpoint(spec.Min, spec.Max),
		step:         spec.Step,
		fineStep:     spec.FineStep,
		coarseStep:   spec.CoarseStep,
		currentValue: spec.InitialValue,
	}
	if err := controller.renderValue(spec.InitialValue); err != nil {
		return err
	}
	adapter.SetUpdateHook(controller.tick)

	runErr := adapter.Run(ctx, spec.Title)
	stop()
	return ignoreContextStop(runErr)
}

func buildGaugePreviewSpec(options dashboardPreviewOptions) (previewGaugeSpec, error) {
	cfg, err := v3config.LoadFile(options.ConfigPath)
	if err != nil {
		return previewGaugeSpec{}, fmt.Errorf("load dashboard preview %q: %w", options.ConfigPath, err)
	}

	vehicleIDs := sortedVehicleIDs(cfg.Vehicles)
	if len(vehicleIDs) != 1 {
		return previewGaugeSpec{}, fmt.Errorf("dashboard preview %q must define exactly one vehicle; found: %s", options.ConfigPath, strings.Join(vehicleIDs, ", "))
	}

	plan, err := v3config.Resolve(cfg, vehicleIDs[0])
	if err != nil {
		return previewGaugeSpec{}, fmt.Errorf("resolve dashboard preview %q: %w", options.ConfigPath, err)
	}

	searchPaths, err := v3assets.DefaultSearchPaths(options.ConfigPath, plan.VehicleID)
	if err != nil {
		return previewGaugeSpec{}, err
	}
	candidate, err := selectPreviewGauge(plan.Dashboards, strings.TrimSpace(options.GaugeID))
	if err != nil {
		return previewGaugeSpec{}, err
	}

	pkg, err := v3gauges.LoadPackageWithSearchPaths(searchPaths, candidate.Widget.Gauge)
	if err != nil {
		return previewGaugeSpec{}, fmt.Errorf("load preview gauge %q: %w", previewGaugeLabel(candidate), err)
	}
	sensorCfg, ok := cfg.Sensors[pkg.Sensor]
	if !ok {
		return previewGaugeSpec{}, fmt.Errorf("preview gauge %q references sensor %q which is not defined in top-level sensors", previewGaugeLabel(candidate), pkg.Sensor)
	}

	minValue, maxValue, err := previewGaugeRange(sensorCfg, pkg)
	if err != nil {
		return previewGaugeSpec{}, fmt.Errorf("preview gauge %q: %w", previewGaugeLabel(candidate), err)
	}
	step, fineStep, coarseStep := previewStepSizes(minValue, maxValue, options.Step, options.FineStep, options.CoarseStep)
	initialValue := midpoint(minValue, maxValue)
	if options.InitialValue != nil {
		initialValue = clampPreviewValue(*options.InitialValue, minValue, maxValue)
	}

	registry, err := v3assets.LoadWithSearchPaths(plan.Assets, searchPaths)
	if err != nil {
		return previewGaugeSpec{}, fmt.Errorf("load preview assets: %w", err)
	}

	previewWidget := candidate.Widget
	previewWidget.Position = []int{0, 0}
	previewPlan := plan
	previewPlan.Dashboards = []v3config.ResolvedDashboard{{
		ID: candidate.DashboardID,
		Config: v3config.DashboardConfig{
			Display: candidate.DashboardID + "_preview",
			Widgets: []v3config.WidgetConfig{previewWidget},
		},
	}}

	runtime, err := v3dashboard.NewRuntime(previewPlan, registry)
	if err != nil {
		return previewGaugeSpec{}, fmt.Errorf("create preview runtime: %w", err)
	}

	return previewGaugeSpec{
		Runtime:       runtime,
		Title:         fmt.Sprintf("GoDriveLog preview - %s", previewGaugeLabel(candidate)),
		SensorID:      pkg.Sensor,
		Unit:          sensorCfg.Unit,
		Min:           minValue,
		Max:           maxValue,
		InitialValue:  initialValue,
		Step:          step,
		FineStep:      fineStep,
		CoarseStep:    coarseStep,
		SelectedGauge: previewGaugeLabel(candidate),
	}, nil
}

func selectPreviewGauge(dashboards []v3config.ResolvedDashboard, requested string) (previewGaugeCandidate, error) {
	candidates := make([]previewGaugeCandidate, 0)
	for _, dashboard := range dashboards {
		for _, widget := range dashboard.Config.Widgets {
			if widget.Type != v3config.WidgetTypeGauge {
				continue
			}
			candidates = append(candidates, previewGaugeCandidate{
				DashboardID: dashboard.ID,
				Widget:      widget,
			})
		}
	}
	if len(candidates) == 0 {
		return previewGaugeCandidate{}, fmt.Errorf("dashboard preview file does not contain any gauge widgets")
	}

	requested = strings.TrimSpace(requested)
	if requested == "" {
		if len(candidates) == 1 {
			return candidates[0], nil
		}
		return previewGaugeCandidate{}, fmt.Errorf("dashboard preview file contains multiple gauges; supply --gauge with one of: %s", joinPreviewGaugeLabels(candidates))
	}

	var exactMatches []previewGaugeCandidate
	var shortMatches []previewGaugeCandidate
	for _, candidate := range candidates {
		if previewGaugeLabel(candidate) == requested {
			exactMatches = append(exactMatches, candidate)
		}
		if candidate.Widget.ID == requested {
			shortMatches = append(shortMatches, candidate)
		}
	}
	if len(exactMatches) == 1 {
		return exactMatches[0], nil
	}
	if len(shortMatches) == 1 {
		return shortMatches[0], nil
	}
	if len(shortMatches) > 1 {
		return previewGaugeCandidate{}, fmt.Errorf("preview gauge %q is ambiguous; use one of: %s", requested, joinPreviewGaugeLabels(shortMatches))
	}
	return previewGaugeCandidate{}, fmt.Errorf("preview gauge %q was not found; available gauges: %s", requested, joinPreviewGaugeLabels(candidates))
}

func joinPreviewGaugeLabels(candidates []previewGaugeCandidate) string {
	labels := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		labels = append(labels, previewGaugeLabel(candidate))
	}
	return strings.Join(labels, ", ")
}

func previewGaugeLabel(candidate previewGaugeCandidate) string {
	return candidate.DashboardID + "/" + candidate.Widget.ID
}

func previewGaugeRange(sensorCfg v3config.SensorConfig, pkg v3gauges.Package) (float64, float64, error) {
	if sensorCfg.Min != nil && sensorCfg.Max != nil && *sensorCfg.Max > *sensorCfg.Min {
		return *sensorCfg.Min, *sensorCfg.Max, nil
	}
	switch pkg.Type {
	case v3gauges.TypeRadial, v3gauges.TypeBar:
		if pkg.ValueMap.Max > pkg.ValueMap.Min {
			return pkg.ValueMap.Min, pkg.ValueMap.Max, nil
		}
	case v3gauges.TypeSegmented:
		return 0, 100, nil
	case v3gauges.TypeIndicator:
		return 0, 1, nil
	}
	return 0, 0, fmt.Errorf("could not infer a preview range; define sensor min/max for %q", pkg.Sensor)
}

func previewStepSizes(minValue, maxValue float64, step, fineStep, coarseStep *float64) (float64, float64, float64) {
	span := maxValue - minValue
	base := nicePreviewStep(span / 40)
	if base <= 0 {
		base = 1
	}
	fine := nicePreviewStep(base / 10)
	if fine <= 0 {
		fine = base
	}
	coarse := nicePreviewStep(base * 5)
	if coarse <= 0 {
		coarse = base
	}
	if step != nil {
		base = *step
	}
	if fineStep != nil {
		fine = *fineStep
	}
	if coarseStep != nil {
		coarse = *coarseStep
	}
	return base, fine, coarse
}

func nicePreviewStep(value float64) float64 {
	value = math.Abs(value)
	if value == 0 {
		return 0
	}
	exponent := math.Floor(math.Log10(value))
	scale := math.Pow(10, exponent)
	normalized := value / scale
	switch {
	case normalized <= 1:
		return 1 * scale
	case normalized <= 2:
		return 2 * scale
	case normalized <= 2.5:
		return 2.5 * scale
	case normalized <= 5:
		return 5 * scale
	default:
		return 10 * scale
	}
}

func midpoint(minValue, maxValue float64) float64 {
	return minValue + ((maxValue - minValue) / 2)
}

func clampPreviewValue(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func (c *gaugePreviewController) tick() error {
	if c == nil {
		return nil
	}
	if inpututil.IsKeyJustPressed(ebitenui.KeyEscape) || inpututil.IsKeyJustPressed(ebitenui.KeyQ) {
		c.stop()
		return nil
	}

	if c.replayPending {
		c.replayPending = false
		return c.renderValue(c.lastTo)
	}

	if inpututil.IsKeyJustPressed(ebitenui.KeySpace) && c.hasLast && c.lastFrom != c.lastTo {
		c.replayPending = true
		return c.renderValue(c.lastFrom)
	}
	if inpututil.IsKeyJustPressed(ebitenui.KeyR) {
		return c.applyValue(c.midpoint)
	}
	if inpututil.IsKeyJustPressed(ebitenui.KeyLeft) {
		return c.applyValue(c.min)
	}
	if inpututil.IsKeyJustPressed(ebitenui.KeyRight) {
		return c.applyValue(c.max)
	}
	if inpututil.IsKeyJustPressed(ebitenui.KeyUp) {
		return c.applyValue(c.currentValue + c.activeStep())
	}
	if inpututil.IsKeyJustPressed(ebitenui.KeyDown) {
		return c.applyValue(c.currentValue - c.activeStep())
	}

	_, wheelY := ebitenui.Wheel()
	if wheelY > 0 {
		return c.applyValue(c.currentValue + c.step)
	}
	if wheelY < 0 {
		return c.applyValue(c.currentValue - c.step)
	}
	return nil
}

func (c *gaugePreviewController) activeStep() float64 {
	if ebitenui.IsKeyPressed(ebitenui.KeyControl) || ebitenui.IsKeyPressed(ebitenui.KeyMeta) {
		return c.fineStep
	}
	if ebitenui.IsKeyPressed(ebitenui.KeyShift) {
		return c.coarseStep
	}
	return c.step
}

func (c *gaugePreviewController) applyValue(target float64) error {
	target = clampPreviewValue(target, c.min, c.max)
	if target == c.currentValue {
		return nil
	}
	c.lastFrom = c.currentValue
	c.lastTo = target
	c.hasLast = true
	return c.renderValue(target)
}

func (c *gaugePreviewController) renderValue(value float64) error {
	state := sensors.SensorState{
		ID:         c.sensorID,
		Value:      value,
		TypedValue: sensors.NewNumericValue(value, c.unit),
		Unit:       c.unit,
		Min:        c.min,
		Max:        c.max,
		Status:     sensors.StatusOK,
		UpdatedAt:  time.Now(),
	}
	c.runtime.SetState(state)
	scenes, err := c.runtime.Snapshot()
	if err != nil {
		return err
	}
	if err := c.adapter.UpdateScenes(scenes); err != nil {
		return err
	}
	c.currentValue = value
	return nil
}
