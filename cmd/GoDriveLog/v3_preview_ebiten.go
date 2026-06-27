//go:build !fyne_legacy

package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
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
	"gopkg.in/yaml.v3"
)

type previewGaugeCandidate struct {
	DashboardID string
	Widget      v3config.WidgetConfig
}

type GaugePreviewSpec struct {
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
	spec, err := LoadGaugePreviewFile(options.ConfigPath, options.GaugeID)
	if err != nil {
		return err
	}
	if options.InitialValue != nil {
		spec.InitialValue = clampPreviewValue(*options.InitialValue, spec.Min, spec.Max)
	}
	spec.Step, spec.FineStep, spec.CoarseStep = previewStepSizes(spec.Min, spec.Max, options.Step, options.FineStep, options.CoarseStep)

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
	adapter.SetUpdateHook(controller.tick)

	runErr := v3dashboard.WithGaugePackageLoader(loadPreviewGaugePackageWithSearchPaths, func() error {
		if err := controller.renderValue(spec.InitialValue); err != nil {
			return err
		}
		return adapter.Run(ctx, spec.Title)
	})
	stop()
	return ignoreContextStop(runErr)
}

func LoadGaugePreviewFile(path string, gaugeID string) (GaugePreviewSpec, error) {
	configPath, cfg, err := loadGaugePreviewConfig(path)
	if err != nil {
		return GaugePreviewSpec{}, err
	}

	vehicleIDs := sortedVehicleIDs(cfg.Vehicles)
	if len(vehicleIDs) != 1 {
		return GaugePreviewSpec{}, fmt.Errorf("dashboard preview %q must define exactly one vehicle; found: %s", configPath, strings.Join(vehicleIDs, ", "))
	}
	vehicleID := vehicleIDs[0]
	selectedDashboards, err := previewSelectedDashboards(cfg, vehicleID)
	if err != nil {
		return GaugePreviewSpec{}, err
	}

	candidate, err := selectPreviewGauge(selectedDashboards, strings.TrimSpace(gaugeID))
	if err != nil {
		return GaugePreviewSpec{}, err
	}

	resolvedGaugePath, searchPaths, err := resolvePreviewGaugePath(configPath, candidate.Widget.Gauge)
	if err != nil {
		return GaugePreviewSpec{}, fmt.Errorf("resolve preview gauge %q: %w", previewGaugeLabel(candidate), err)
	}

	pkg, err := v3gauges.LoadPackageForPreview(resolvedGaugePath)
	if err != nil {
		return GaugePreviewSpec{}, fmt.Errorf("load preview gauge %q: %w", previewGaugeLabel(candidate), err)
	}
	sensorCfg, ok := cfg.Sensors[pkg.Sensor]
	if !ok {
		return GaugePreviewSpec{}, fmt.Errorf("preview gauge %q references sensor %q which is not defined in top-level sensors", previewGaugeLabel(candidate), pkg.Sensor)
	}

	minValue, maxValue, err := previewGaugeRange(sensorCfg, pkg)
	if err != nil {
		return GaugePreviewSpec{}, fmt.Errorf("preview gauge %q: %w", previewGaugeLabel(candidate), err)
	}

	registry, err := v3assets.LoadWithSearchPaths(cfg.Assets, searchPaths)
	if err != nil {
		return GaugePreviewSpec{}, fmt.Errorf("load preview assets: %w", err)
	}

	previewWidget := candidate.Widget
	previewWidget.Gauge = resolvedGaugePath
	if len(previewWidget.Position) != 2 {
		previewWidget.Position = []int{0, 0}
	}
	if previewWidget.Scale <= 0 {
		previewWidget.Scale = 1
	}

	dashboardConfig := v3config.DashboardConfig{
		Display: candidate.DashboardID + "_preview",
		Size: v3config.SizeConfig{
			Width:  previewDashboardWidth(candidate.Widget, pkg),
			Height: previewDashboardHeight(candidate.Widget, pkg),
		},
		Widgets: []v3config.WidgetConfig{previewWidget},
	}
	if dashboardConfig.Size.Width <= 0 {
		dashboardConfig.Size.Width = 800
	}
	if dashboardConfig.Size.Height <= 0 {
		dashboardConfig.Size.Height = 480
	}

	var runtime *v3dashboard.Runtime
	err = v3dashboard.WithGaugePackageLoader(loadPreviewGaugePackageWithSearchPaths, func() error {
		var runtimeErr error
		runtime, runtimeErr = v3dashboard.NewRuntime(v3config.RuntimePlan{
			VehicleID: vehicleID,
			Assets:    cfg.Assets,
			Dashboards: []v3config.ResolvedDashboard{{
				ID:     candidate.DashboardID,
				Config: dashboardConfig,
			}},
		}, registry)
		return runtimeErr
	})
	if err != nil {
		return GaugePreviewSpec{}, fmt.Errorf("create preview runtime: %w", err)
	}

	return GaugePreviewSpec{
		Runtime:       runtime,
		Title:         fmt.Sprintf("GoDriveLog preview - %s", previewGaugeLabel(candidate)),
		SensorID:      pkg.Sensor,
		Unit:          sensorCfg.Unit,
		Min:           minValue,
		Max:           maxValue,
		InitialValue:  midpoint(minValue, maxValue),
		Step:          0,
		FineStep:      0,
		CoarseStep:    0,
		SelectedGauge: previewGaugeLabel(candidate),
	}, nil
}

func buildGaugePreviewSpec(options dashboardPreviewOptions) (GaugePreviewSpec, error) {
	spec, err := LoadGaugePreviewFile(options.ConfigPath, options.GaugeID)
	if err != nil {
		return GaugePreviewSpec{}, err
	}
	if options.InitialValue != nil {
		spec.InitialValue = clampPreviewValue(*options.InitialValue, spec.Min, spec.Max)
	}
	spec.Step, spec.FineStep, spec.CoarseStep = previewStepSizes(spec.Min, spec.Max, options.Step, options.FineStep, options.CoarseStep)
	return spec, nil
}

func loadGaugePreviewConfig(path string) (string, v3config.Config, error) {
	configPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", v3config.Config{}, fmt.Errorf("resolve dashboard preview %q: %w", path, err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", v3config.Config{}, fmt.Errorf("read dashboard preview %q: %w", path, err)
	}

	var cfg v3config.Config
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return "", v3config.Config{}, fmt.Errorf("parse dashboard preview %q: %w", path, err)
	}
	return configPath, cfg, nil
}

func previewSelectedDashboards(cfg v3config.Config, vehicleID string) ([]v3config.ResolvedDashboard, error) {
	vehicle, ok := cfg.Vehicles[vehicleID]
	if !ok {
		return nil, fmt.Errorf("preview vehicle %q is not defined", vehicleID)
	}
	dashboardIDs := append([]string(nil), vehicle.Dashboards...)
	if len(dashboardIDs) == 0 {
		dashboardIDs = sortedDashboardIDs(cfg.Dashboards)
	}
	if len(dashboardIDs) == 0 {
		return nil, fmt.Errorf("dashboard preview file does not define any dashboards")
	}
	selected := make([]v3config.ResolvedDashboard, 0, len(dashboardIDs))
	for _, dashboardID := range dashboardIDs {
		dashboardCfg, ok := cfg.Dashboards[dashboardID]
		if !ok {
			return nil, fmt.Errorf("preview vehicle %q references dashboard %q which is not defined", vehicleID, dashboardID)
		}
		selected = append(selected, v3config.ResolvedDashboard{ID: dashboardID, Config: dashboardCfg})
	}
	return selected, nil
}

func sortedDashboardIDs(dashboards map[string]v3config.DashboardConfig) []string {
	ids := make([]string, 0, len(dashboards))
	for id := range dashboards {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func resolvePreviewGaugePath(configPath, gaugePath string) (string, []string, error) {
	cleanedGaugePath := strings.TrimSpace(gaugePath)
	if cleanedGaugePath == "" {
		return "", nil, fmt.Errorf("gauge path must not be empty")
	}
	previewDir := filepath.Dir(configPath)
	searchPaths := previewGaugeSearchPaths(previewDir)
	for _, root := range searchPaths {
		candidate := cleanedGaugePath
		if !filepath.IsAbs(candidate) {
			candidate = filepath.Join(root, cleanedGaugePath)
		}
		resolved, ok := previewGaugePackageDir(candidate)
		if ok {
			return resolved, searchPaths, nil
		}
	}
	if filepath.IsAbs(cleanedGaugePath) {
		return "", searchPaths, fmt.Errorf("gauge package %q was not found", cleanedGaugePath)
	}
	tried := make([]string, 0, len(searchPaths))
	for _, root := range searchPaths {
		tried = append(tried, filepath.Join(root, cleanedGaugePath))
	}
	return "", searchPaths, fmt.Errorf("gauge package %q was not found in preview search paths: %s", cleanedGaugePath, strings.Join(tried, ", "))
}

func previewGaugeSearchPaths(previewDir string) []string {
	paths := []string{}
	if previewDir != "" {
		paths = append(paths, previewDir)
	}
	if pwd, err := os.Getwd(); err == nil {
		paths = append(paths, pwd)
		if repoRoot := findRepoRoot(pwd); repoRoot != "" {
			paths = append(paths, repoRoot)
		}
	}
	cleaned := make([]string, 0, len(paths))
	seen := map[string]bool{}
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil || seen[abs] {
			continue
		}
		seen[abs] = true
		cleaned = append(cleaned, abs)
	}
	return cleaned
}

func findRepoRoot(start string) string {
	current := start
	for {
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func previewGaugePackageDir(candidate string) (string, bool) {
	if strings.TrimSpace(candidate) == "" {
		return "", false
	}
	resolved := filepath.Clean(candidate)
	info, err := os.Stat(resolved)
	if err == nil && !info.IsDir() && filepath.Base(resolved) == "gauge.yaml" {
		resolved = filepath.Dir(resolved)
	} else if err == nil && !info.IsDir() {
		return "", false
	}
	if _, err := os.Stat(filepath.Join(resolved, "gauge.yaml")); err == nil {
		abs, absErr := filepath.Abs(resolved)
		if absErr != nil {
			return "", false
		}
		return abs, true
	}
	return "", false
}

func previewDashboardWidth(widget v3config.WidgetConfig, pkg v3gauges.Package) int {
	return previewDimension(widget.Position, 0, widget.Scale, pkg.Size.Width)
}

func previewDashboardHeight(widget v3config.WidgetConfig, pkg v3gauges.Package) int {
	return previewDimension(widget.Position, 1, widget.Scale, pkg.Size.Height)
}

func previewDimension(position []int, index int, scale float64, size int) int {
	offset := 0
	if len(position) > index {
		offset = position[index]
	}
	if scale <= 0 {
		scale = 1
	}
	if size <= 0 {
		return 0
	}
	return offset + int(math.Ceil(float64(size)*scale))
}

func loadPreviewGaugePackageWithSearchPaths(searchPaths []string, packageDir string) (v3gauges.Package, error) {
	if filepath.IsAbs(packageDir) {
		return v3gauges.LoadPackageForPreview(packageDir)
	}
	for _, root := range searchPaths {
		resolved, ok := previewGaugePackageDir(filepath.Join(root, packageDir))
		if ok {
			return v3gauges.LoadPackageForPreview(resolved)
		}
	}
	return v3gauges.LoadPackageForPreview(packageDir)
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
