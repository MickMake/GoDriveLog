package harness

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	v3gauges "github.com/MickMake/GoDriveLog/internal/dashboard/gauges"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
	"github.com/MickMake/GoDriveLog/internal/sensors"
)

const (
	PatternFixed     = "fixed"
	PatternSweep     = "sweep"
	PatternHeartbeat = "heartbeat"

	defaultInterval          = 100 * time.Millisecond
	defaultSweepCycle        = 11 * time.Second
	defaultSweepRise         = 5 * time.Second
	defaultSweepHold         = time.Second
	defaultHeartbeatCycle    = 10 * time.Second
	defaultHeartbeatBaseline = 0.08
	defaultIncrementCycle    = 10 * time.Second
	defaultIncrementRise     = 5 * time.Second
	defaultIndicatorSlow     = time.Second
	defaultIndicatorFast     = 250 * time.Millisecond
	defaultIndicatorCycle    = 10 * time.Second
	defaultBarPulseCycle     = (2 * time.Second) / 3
)

// Scene is the harness dashboard scene boundary type.
type Scene = v3dashboard.Scene

// SceneSink is the harness output boundary. Production display adapters can
// consume this just like the normal v3 dashboard runtime boundary.
type SceneSink func([]Scene) error

// Options controls the v3 dashboard/gauge harness. The harness deliberately
// avoids endpoint access and feeds synthetic sensor events through the real v3
// dashboard event/state path.
type Options struct {
	ConfigPath       string
	VehicleID        string
	AssetSearchPaths []string
	Pattern          string
	Interval         time.Duration
	Sink             SceneSink
	Logger           *log.Logger
	Now              func() time.Time

	// MaxEvents is intended for focused tests. Zero means run until ctx is done.
	MaxEvents int
}

// Summary describes one harness run.
type Summary struct {
	VehicleID      string
	VehicleName    string
	SensorCount    int
	DashboardCount int
	Pattern        string
	Interval       time.Duration
	Events         int
}

type sweepProfile string

const (
	sweepProfileRange       sweepProfile = "range"
	sweepProfileIncremental sweepProfile = "incremental"
	sweepProfileIndicator   sweepProfile = "indicator"
	sweepProfileHeartbeat   sweepProfile = "heartbeat"
)

type sensorSource struct {
	ID           string
	Unit         string
	Min          float64
	Max          float64
	SweepProfile sweepProfile
	first        bool
}

// NormalizePattern validates and normalises a harness pattern name.
func NormalizePattern(pattern string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(pattern)) {
	case "", PatternSweep:
		return PatternSweep, nil
	case PatternHeartbeat:
		return PatternHeartbeat, nil
	case PatternFixed:
		return PatternFixed, nil
	default:
		return "", fmt.Errorf("unknown harness pattern %q; expected %s, %s, or %s", pattern, PatternSweep, PatternHeartbeat, PatternFixed)
	}
}

// Run loads one selected v3 dashboard path, then drives it with fake sensor
// events. It uses the same dashboard Runtime.ApplyEvent boundary as the normal
// v3 command path, but it never connects to OBD or starts production polling.
func Run(ctx context.Context, opts Options) (Summary, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(opts.ConfigPath) == "" {
		return Summary{}, fmt.Errorf("v3 dashboard harness requires a config path")
	}
	if opts.Sink == nil {
		return Summary{}, fmt.Errorf("v3 dashboard harness requires a scene sink")
	}
	pattern, err := NormalizePattern(opts.Pattern)
	if err != nil {
		return Summary{}, err
	}
	interval := opts.Interval
	if interval <= 0 {
		interval = defaultInterval
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}

	cfg, err := v3config.LoadFile(opts.ConfigPath)
	if err != nil {
		return Summary{}, fmt.Errorf("load v3 config: %w", err)
	}
	plan, err := v3config.Resolve(cfg, opts.VehicleID)
	if err != nil {
		return Summary{}, fmt.Errorf("resolve v3 runtime plan: %w", err)
	}
	if len(plan.Dashboards) == 0 {
		return Summary{}, fmt.Errorf("v3 dashboard harness requires at least one selected dashboard")
	}
	if len(plan.Sensors) == 0 {
		return Summary{}, fmt.Errorf("v3 dashboard harness requires at least one sensor")
	}

	searchPaths := opts.AssetSearchPaths
	if len(searchPaths) == 0 {
		searchPaths, err = v3assets.DefaultSearchPaths(opts.ConfigPath, plan.VehicleID)
		if err != nil {
			return Summary{}, err
		}
	}
	registry, err := v3assets.LoadWithSearchPaths(plan.Assets, searchPaths)
	if err != nil {
		return Summary{}, fmt.Errorf("load v3 dashboard assets: %w", err)
	}
	dashboardRuntime, err := v3dashboard.NewRuntime(plan, registry)
	if err != nil {
		return Summary{}, fmt.Errorf("create v3 dashboard runtime: %w", err)
	}

	sources, err := sensorSources(plan, searchPaths)
	if err != nil {
		return Summary{}, err
	}
	summary := Summary{
		VehicleID:      plan.VehicleID,
		VehicleName:    plan.Vehicle.Name,
		SensorCount:    len(sources),
		DashboardCount: len(plan.Dashboards),
		Pattern:        pattern,
		Interval:       interval,
	}
	if opts.Logger != nil {
		opts.Logger.Printf("v3 dashboard harness starting: vehicle=%s sensors=%d dashboards=%d pattern=%s interval=%s", summary.VehicleID, summary.SensorCount, summary.DashboardCount, pattern, interval)
	}

	started := now()
	emitAll := func(at time.Time) (bool, error) {
		elapsed := at.Sub(started)
		var latestScenes []Scene
		anyChanged := false
		for i := range sources {
			event := eventForSource(&sources[i], pattern, elapsed, at)
			scenes, changed, err := dashboardRuntime.ApplyEvent(event)
			if err != nil {
				return false, err
			}
			if changed {
				latestScenes = scenes
				anyChanged = true
			}
			summary.Events++
			if opts.MaxEvents > 0 && summary.Events >= opts.MaxEvents {
				if anyChanged {
					if err := opts.Sink(latestScenes); err != nil {
						return false, err
					}
				}
				return true, nil
			}
		}
		if anyChanged {
			if err := opts.Sink(latestScenes); err != nil {
				return false, err
			}
		}
		return false, nil
	}

	if done, err := emitAll(started); done || err != nil {
		return summary, err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			if opts.Logger != nil {
				opts.Logger.Printf("v3 dashboard harness stopped: vehicle=%s events=%d", summary.VehicleID, summary.Events)
			}
			return summary, nil
		case at := <-ticker.C:
			done, err := emitAll(at)
			if err != nil {
				return summary, err
			}
			if done {
				if opts.Logger != nil {
					opts.Logger.Printf("v3 dashboard harness stopped: vehicle=%s events=%d", summary.VehicleID, summary.Events)
				}
				return summary, nil
			}
		}
	}
}

// ValueForPattern returns the fake value for one sensor at elapsed time. It is
// exported so tests and later docs/examples can share the exact pattern rules.
func ValueForPattern(pattern string, minValue, maxValue float64, elapsed time.Duration) (float64, error) {
	pattern, err := NormalizePattern(pattern)
	if err != nil {
		return 0, err
	}
	minValue, maxValue = normalRange(minValue, maxValue)
	switch pattern {
	case PatternFixed:
		return minValue + ((maxValue - minValue) / 2), nil
	case PatternSweep:
		return sweepValue(minValue, maxValue, elapsed), nil
	case PatternHeartbeat:
		return heartbeatValue(minValue, maxValue, elapsed), nil
	default:
		return 0, fmt.Errorf("unknown harness pattern %q", pattern)
	}
}

func eventForSource(source *sensorSource, pattern string, elapsed time.Duration, at time.Time) sensors.SensorEvent {
	value, typedValue, err := valueForSourcePattern(pattern, *source, elapsed)
	if err != nil {
		value = source.Min
		typedValue = sensors.NewNumericValue(value, source.Unit)
	}
	kind := sensors.EventKindValueChange
	previousStatus := sensors.StatusOK
	if source.first {
		kind = sensors.EventKindFirstRead
		previousStatus = ""
		source.first = false
	}
	minValue, maxValue := normalRange(source.Min, source.Max)
	state := sensors.SensorState{
		ID:         source.ID,
		Value:      value,
		TypedValue: typedValue,
		Unit:       source.Unit,
		Min:        minValue,
		Max:        maxValue,
		Status:     sensors.StatusOK,
		UpdatedAt:  at,
		StaleAfter: sensors.StaleAfterForPoll(defaultInterval),
	}
	return sensors.SensorEvent{
		Kind:           kind,
		SensorID:       source.ID,
		State:          state,
		PreviousStatus: previousStatus,
		Timestamp:      at,
		ReadAt:         at,
	}
}

func sensorSources(plan v3config.RuntimePlan, searchPaths []string) ([]sensorSource, error) {
	profiles, err := sensorSweepProfiles(plan.Dashboards, searchPaths)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(plan.Sensors))
	for id := range plan.Sensors {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	sources := make([]sensorSource, 0, len(ids))
	for _, id := range ids {
		cfg := plan.Sensors[id]
		minValue := 0.0
		maxValue := 1.0
		if cfg.Min != nil {
			minValue = *cfg.Min
		}
		if cfg.Max != nil {
			maxValue = *cfg.Max
		}
		minValue, maxValue = normalRange(minValue, maxValue)
		profile := profiles[id]
		if profile == "" {
			profile = sweepProfileRange
		}
		sources = append(sources, sensorSource{
			ID:           id,
			Unit:         cfg.Unit,
			Min:          minValue,
			Max:          maxValue,
			SweepProfile: profile,
			first:        true,
		})
	}
	return sources, nil
}

func sensorSweepProfiles(dashboards []v3config.ResolvedDashboard, searchPaths []string) (map[string]sweepProfile, error) {
	profiles := map[string]sweepProfile{}
	for _, dashboard := range dashboards {
		for _, widget := range dashboard.Config.Widgets {
			sensorID, profile, ok, err := widgetSweepProfile(widget, searchPaths)
			if err != nil {
				return nil, fmt.Errorf("dashboard %q widget %q sweep profile: %w", dashboard.ID, widget.ID, err)
			}
			if !ok || strings.TrimSpace(sensorID) == "" {
				continue
			}
			profiles[sensorID] = mergeSweepProfile(profiles[sensorID], profile)
		}
	}
	return profiles, nil
}

func widgetSweepProfile(widget v3config.WidgetConfig, searchPaths []string) (string, sweepProfile, bool, error) {
	switch widget.Type {
	case v3config.WidgetTypeImage:
		return "", "", false, nil
	case v3config.WidgetTypeDigitDisplay:
		return widget.Sensor, sweepProfileIncremental, true, nil
	case v3config.WidgetTypeBarDisplay:
		return widget.Sensor, sweepProfileHeartbeat, true, nil
	case v3config.WidgetTypeFrameGauge:
		return widget.Sensor, sweepProfileRange, true, nil
	case v3config.WidgetTypeIndicator:
		return widget.Sensor, sweepProfileIndicator, true, nil
	case v3config.WidgetTypeGauge:
		pkg, err := v3gauges.LoadPackageWithSearchPaths(searchPaths, widget.Gauge)
		if err != nil {
			return "", "", false, err
		}
		return pkg.Sensor, sweepProfileForGaugeType(pkg.Type), true, nil
	default:
		if strings.TrimSpace(widget.Sensor) == "" {
			return "", "", false, nil
		}
		return widget.Sensor, sweepProfileRange, true, nil
	}
}

func sweepProfileForGaugeType(gaugeType string) sweepProfile {
	switch gaugeType {
	case v3gauges.TypeNumeric, v3gauges.TypeOdometer:
		return sweepProfileIncremental
	case v3gauges.TypeIndicator:
		return sweepProfileIndicator
	case v3gauges.TypeBar:
		return sweepProfileHeartbeat
	case v3gauges.TypeSegmented, v3gauges.TypeRadial:
		return sweepProfileRange
	default:
		return sweepProfileRange
	}
}

// Shared sensors can only emit one synthetic waveform, so mixed-use sensors
// resolve to the highest-impact profile instead of duplicating sensor IDs.
func mergeSweepProfile(current, next sweepProfile) sweepProfile {
	if sweepProfilePriority(next) > sweepProfilePriority(current) {
		return next
	}
	if current == "" {
		return next
	}
	return current
}

func sweepProfilePriority(profile sweepProfile) int {
	switch profile {
	case sweepProfileIndicator:
		return 4
	case sweepProfileHeartbeat:
		return 3
	case sweepProfileRange:
		return 2
	case sweepProfileIncremental:
		return 1
	default:
		return 0
	}
}

func valueForSourcePattern(pattern string, source sensorSource, elapsed time.Duration) (float64, sensors.Value, error) {
	pattern, err := NormalizePattern(pattern)
	if err != nil {
		return 0, sensors.Value{}, err
	}
	minValue, maxValue := normalRange(source.Min, source.Max)
	switch pattern {
	case PatternFixed:
		value := minValue + ((maxValue - minValue) / 2)
		return value, sensors.NewNumericValue(value, source.Unit), nil
	case PatternSweep:
		return sweepValueForSource(source, elapsed)
	case PatternHeartbeat:
		value := heartbeatValue(minValue, maxValue, elapsed)
		return value, sensors.NewNumericValue(value, source.Unit), nil
	default:
		return 0, sensors.Value{}, fmt.Errorf("unknown harness pattern %q", pattern)
	}
}

func sweepValueForSource(source sensorSource, elapsed time.Duration) (float64, sensors.Value, error) {
	switch source.SweepProfile {
	case sweepProfileIncremental:
		value := incrementalSweepValue(elapsed)
		return value, sensors.NewNumericValue(value, source.Unit), nil
	case sweepProfileIndicator:
		on := indicatorSweepOn(elapsed)
		if on {
			return 1, sensors.NewBoolValue(true), nil
		}
		return 0, sensors.NewBoolValue(false), nil
	case sweepProfileHeartbeat:
		value := barPulseValue(source.Min, source.Max, elapsed)
		return value, sensors.NewNumericValue(value, source.Unit), nil
	case sweepProfileRange, "":
		value := sweepValue(source.Min, source.Max, elapsed)
		return value, sensors.NewNumericValue(value, source.Unit), nil
	default:
		return 0, sensors.Value{}, fmt.Errorf("unknown harness sweep profile %q", source.SweepProfile)
	}
}

func normalRange(minValue, maxValue float64) (float64, float64) {
	if math.IsNaN(minValue) || math.IsInf(minValue, 0) {
		minValue = 0
	}
	if math.IsNaN(maxValue) || math.IsInf(maxValue, 0) || maxValue <= minValue {
		maxValue = minValue + 1
	}
	return minValue, maxValue
}

func incrementalSweepValue(elapsed time.Duration) float64 {
	cycle := positiveModulo(elapsed, defaultIncrementCycle)
	if cycle <= defaultIncrementRise {
		fraction := float64(cycle) / float64(defaultIncrementRise)
		return -20 + fraction*40
	}
	secondElapsed := cycle - defaultIncrementRise
	fraction := float64(secondElapsed) / float64(defaultIncrementRise)
	return 20 + fraction*10
}

func indicatorSweepOn(elapsed time.Duration) bool {
	cycle := positiveModulo(elapsed, defaultIndicatorCycle)
	if cycle < defaultIncrementRise {
		return blinkOn(cycle, defaultIndicatorSlow)
	}
	return blinkOn(cycle-defaultIncrementRise, defaultIndicatorFast)
}

func blinkOn(elapsed, interval time.Duration) bool {
	if interval <= 0 {
		return true
	}
	return int(elapsed/interval)%2 == 0
}

// sweepValue uses an 11 second cycle: 5 seconds min->max, 1 second pause at
// max, then 5 seconds max->min.
func sweepValue(minValue, maxValue float64, elapsed time.Duration) float64 {
	cycle := positiveModulo(elapsed, defaultSweepCycle)
	if cycle <= defaultSweepRise {
		fraction := float64(cycle) / float64(defaultSweepRise)
		return minValue + fraction*(maxValue-minValue)
	}
	if cycle <= defaultSweepRise+defaultSweepHold {
		return maxValue
	}
	fallElapsed := cycle - defaultSweepRise - defaultSweepHold
	fraction := float64(fallElapsed) / float64(defaultSweepRise)
	return maxValue - fraction*(maxValue-minValue)
}

func barPulseValue(minValue, maxValue float64, elapsed time.Duration) float64 {
	cycle := positiveModulo(elapsed, defaultBarPulseCycle)
	rangeValue := maxValue - minValue
	baseline := minValue + rangeValue*0.12
	peak := minValue + rangeValue*0.92
	recovery := minValue + rangeValue*0.38
	points := []wavePoint{
		{at: 0, value: baseline},
		{at: 80 * time.Millisecond, value: peak},
		{at: 160 * time.Millisecond, value: recovery},
		{at: 280 * time.Millisecond, value: baseline},
		{at: defaultBarPulseCycle, value: baseline},
	}
	return interpolate(points, cycle)
}

// heartbeatValue uses a 10 second cycle with a slightly-above-min baseline, a
// small first peak, a deeper negative dip, and a larger second peak.
func heartbeatValue(minValue, maxValue float64, elapsed time.Duration) float64 {
	cycle := positiveModulo(elapsed, defaultHeartbeatCycle)
	rangeValue := maxValue - minValue
	baseline := minValue + rangeValue*defaultHeartbeatBaseline
	negative := minValue + rangeValue*0.02
	firstPeak := minValue + rangeValue*0.62
	secondPeak := maxValue

	points := []wavePoint{
		{at: 0, value: baseline},
		{at: 200 * time.Millisecond, value: baseline},
		{at: 450 * time.Millisecond, value: firstPeak},
		{at: 700 * time.Millisecond, value: baseline},
		{at: 950 * time.Millisecond, value: negative},
		{at: 1250 * time.Millisecond, value: secondPeak},
		{at: 1700 * time.Millisecond, value: baseline},
		{at: defaultHeartbeatCycle, value: baseline},
	}
	return interpolate(points, cycle)
}

type wavePoint struct {
	at    time.Duration
	value float64
}

func interpolate(points []wavePoint, elapsed time.Duration) float64 {
	if len(points) == 0 {
		return 0
	}
	if elapsed <= points[0].at {
		return points[0].value
	}
	for i := 1; i < len(points); i++ {
		previous := points[i-1]
		current := points[i]
		if elapsed > current.at {
			continue
		}
		if current.at <= previous.at {
			return current.value
		}
		fraction := float64(elapsed-previous.at) / float64(current.at-previous.at)
		return previous.value + fraction*(current.value-previous.value)
	}
	return points[len(points)-1].value
}

func positiveModulo(value, period time.Duration) time.Duration {
	if period <= 0 {
		return 0
	}
	mod := value % period
	if mod < 0 {
		mod += period
	}
	return mod
}
