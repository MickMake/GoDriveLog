package harness

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	v3assets "github.com/MickMake/GoDriveLog/internal/assets"
	"github.com/MickMake/GoDriveLog/internal/config/v3config"
	"github.com/MickMake/GoDriveLog/internal/dashboard/v3dashboard"
)

func TestNormalizePatternRejectsUnknown(t *testing.T) {
	if _, err := NormalizePattern("wobble"); err == nil {
		t.Fatal("NormalizePattern accepted unknown pattern")
	}
}

func TestSweepPatternRisesHoldsThenFalls(t *testing.T) {
	checks := []struct {
		elapsed time.Duration
		want    float64
	}{
		{elapsed: 0, want: 0},
		{elapsed: 2500 * time.Millisecond, want: 50},
		{elapsed: 5 * time.Second, want: 100},
		{elapsed: 5500 * time.Millisecond, want: 100},
		{elapsed: 6 * time.Second, want: 100},
		{elapsed: 8500 * time.Millisecond, want: 50},
		{elapsed: 11 * time.Second, want: 0},
	}

	for _, check := range checks {
		got, err := ValueForPattern(PatternSweep, 0, 100, check.elapsed)
		if err != nil {
			t.Fatalf("ValueForPattern returned error: %v", err)
		}
		if got != check.want {
			t.Fatalf("ValueForPattern(%s) = %v, want %v", check.elapsed, got, check.want)
		}
	}
}

func TestHeartbeatPatternUsesTwoPeaksAndNegativeDip(t *testing.T) {
	baseline, err := ValueForPattern(PatternHeartbeat, 0, 100, 0)
	if err != nil {
		t.Fatal(err)
	}
	firstPeak, err := ValueForPattern(PatternHeartbeat, 0, 100, 450*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	negative, err := ValueForPattern(PatternHeartbeat, 0, 100, 950*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	secondPeak, err := ValueForPattern(PatternHeartbeat, 0, 100, 1250*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	cycleEnd, err := ValueForPattern(PatternHeartbeat, 0, 100, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if baseline <= 0 || baseline >= 20 {
		t.Fatalf("baseline = %v, want slightly above min", baseline)
	}
	if negative >= baseline {
		t.Fatalf("negative dip = %v, want below baseline %v", negative, baseline)
	}
	if firstPeak <= baseline || firstPeak >= secondPeak {
		t.Fatalf("first peak = %v, want between baseline %v and second peak %v", firstPeak, baseline, secondPeak)
	}
	if secondPeak != 100 {
		t.Fatalf("second peak = %v, want max", secondPeak)
	}
	if cycleEnd != baseline {
		t.Fatalf("cycle end = %v, want baseline %v", cycleEnd, baseline)
	}
}

func TestGaugeAwareSweepIncrementalProfile(t *testing.T) {
	source := sensorSource{ID: "speed", Unit: "km/h", Min: 0, Max: 180, SweepProfile: sweepProfileIncremental}
	checks := []struct {
		elapsed time.Duration
		want    float64
	}{
		{elapsed: 0, want: -20},
		{elapsed: 2500 * time.Millisecond, want: 0},
		{elapsed: 5 * time.Second, want: 20},
		{elapsed: 7500 * time.Millisecond, want: 25},
		{elapsed: 10 * time.Second, want: -20},
	}

	for _, check := range checks {
		got, typed, err := valueForSourcePattern(PatternSweep, source, check.elapsed)
		if err != nil {
			t.Fatalf("valueForSourcePattern(%s) returned error: %v", check.elapsed, err)
		}
		if got != check.want {
			t.Fatalf("valueForSourcePattern(%s) = %v, want %v", check.elapsed, got, check.want)
		}
		if typed.Kind != "numeric" {
			t.Fatalf("typed kind = %q, want numeric", typed.Kind)
		}
	}
}

func TestGaugeAwareSweepIndicatorProfile(t *testing.T) {
	source := sensorSource{ID: "check_engine", SweepProfile: sweepProfileIndicator}
	checks := []struct {
		elapsed time.Duration
		want    bool
	}{
		{elapsed: 0, want: true},
		{elapsed: 999 * time.Millisecond, want: true},
		{elapsed: 1 * time.Second, want: false},
		{elapsed: 5 * time.Second, want: true},
		{elapsed: 5250 * time.Millisecond, want: false},
	}

	for _, check := range checks {
		got, typed, err := valueForSourcePattern(PatternSweep, source, check.elapsed)
		if err != nil {
			t.Fatalf("valueForSourcePattern(%s) returned error: %v", check.elapsed, err)
		}
		if got != 0 && got != 1 {
			t.Fatalf("valueForSourcePattern(%s) value = %v, want 0 or 1", check.elapsed, got)
		}
		if typed.Kind != "bool" || typed.Bool == nil || *typed.Bool != check.want {
			t.Fatalf("typed bool at %s = %#v, want %v", check.elapsed, typed, check.want)
		}
	}
}

func TestGaugeAwareSweepHeartbeatProfile(t *testing.T) {
	source := sensorSource{ID: "fuel_level", Min: 0, Max: 100, SweepProfile: sweepProfileHeartbeat}
	checks := []struct {
		elapsed time.Duration
		want    float64
	}{
		{elapsed: 0, want: 12},
		{elapsed: 80 * time.Millisecond, want: 92},
		{elapsed: 160 * time.Millisecond, want: 38},
		{elapsed: 280 * time.Millisecond, want: 12},
		{elapsed: defaultBarPulseCycle, want: 12},
	}

	for _, check := range checks {
		got, typed, err := valueForSourcePattern(PatternSweep, source, check.elapsed)
		if err != nil {
			t.Fatalf("valueForSourcePattern(%s) returned error: %v", check.elapsed, err)
		}
		if got != check.want {
			t.Fatalf("valueForSourcePattern(%s) = %v, want %v", check.elapsed, got, check.want)
		}
		if typed.Kind != "numeric" {
			t.Fatalf("typed kind = %q, want numeric", typed.Kind)
		}
	}
}

func TestRunFeedsFakeEventsThroughDashboardScenePath(t *testing.T) {
	configPath := writeHarnessConfig(t)

	var sceneUpdates int
	summary, err := Run(context.Background(), Options{
		ConfigPath: configPath,
		VehicleID:  "demo",
		Pattern:    PatternHeartbeat,
		MaxEvents:  2,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
		Sink: func(scenes []v3dashboard.Scene) error {
			sceneUpdates++
			if len(scenes) != 1 {
				t.Fatalf("scene count = %d, want 1", len(scenes))
			}
			if scenes[0].DashboardID != "primary" {
				t.Fatalf("DashboardID = %q, want primary", scenes[0].DashboardID)
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if summary.Pattern != PatternHeartbeat {
		t.Fatalf("summary.Pattern = %q, want %q", summary.Pattern, PatternHeartbeat)
	}
	if summary.Events != 2 {
		t.Fatalf("summary.Events = %d, want 2", summary.Events)
	}
	if sceneUpdates != 1 {
		t.Fatalf("scene updates = %d, want one sink call for one harness tick", sceneUpdates)
	}
}

func TestBaselineDashboardConfigRunsHarnessPatterns(t *testing.T) {
	configPath := setupBaselineHarnessEnvironment(t)
	patterns := []string{PatternFixed, PatternSweep, PatternHeartbeat}

	for _, pattern := range patterns {
		pattern := pattern
		t.Run(pattern, func(t *testing.T) {
			var sceneUpdates int
			summary, err := Run(context.Background(), Options{
				ConfigPath: configPath,
				VehicleID:  "vw_caddy",
				Pattern:    pattern,
				MaxEvents:  3,
				Now: func() time.Time {
					return time.Unix(0, 0)
				},
				Sink: func(scenes []v3dashboard.Scene) error {
					sceneUpdates++
					if len(scenes) != 1 {
						t.Fatalf("scene count = %d, want 1", len(scenes))
					}
					if scenes[0].DashboardID != "baseline" {
						t.Fatalf("DashboardID = %q, want baseline", scenes[0].DashboardID)
					}
					if len(scenes[0].Widgets) != 4 {
						t.Fatalf("widget count = %d, want 4", len(scenes[0].Widgets))
					}
					widgetIDs := map[string]bool{}
					for _, widget := range scenes[0].Widgets {
						widgetIDs[widget.ID] = true
					}
					for _, id := range []string{"temp_3_digit", "speed_3_digit", "rpm_4_digit", "radial_rpm"} {
						if !widgetIDs[id] {
							t.Fatalf("missing widget %q in baseline scene", id)
						}
					}
					return nil
				},
			})
			if err != nil {
				t.Fatalf("Run returned error: %v", err)
			}
			if summary.VehicleID != "vw_caddy" {
				t.Fatalf("summary.VehicleID = %q, want vw_caddy", summary.VehicleID)
			}
			if summary.Pattern != pattern {
				t.Fatalf("summary.Pattern = %q, want %q", summary.Pattern, pattern)
			}
			if summary.SensorCount != 3 {
				t.Fatalf("summary.SensorCount = %d, want 3", summary.SensorCount)
			}
			if summary.DashboardCount != 1 {
				t.Fatalf("summary.DashboardCount = %d, want 1", summary.DashboardCount)
			}
			if summary.Events != 3 {
				t.Fatalf("summary.Events = %d, want 3", summary.Events)
			}
			if sceneUpdates != 1 {
				t.Fatalf("scene updates = %d, want 1", sceneUpdates)
			}
		})
	}
}

func TestFrameworkSmokeDashboardConfigRunsHarness(t *testing.T) {
	configPath := setupExampleHarnessEnvironment(t, filepath.Join("examples", "framework-smoke", "dashboard.yaml"), "demo")

	var sceneUpdates int
	summary, err := Run(context.Background(), Options{
		ConfigPath: configPath,
		VehicleID:  "demo",
		Pattern:    PatternSweep,
		MaxEvents:  3,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
		Sink: func(scenes []v3dashboard.Scene) error {
			sceneUpdates++
			if len(scenes) != 1 {
				t.Fatalf("scene count = %d, want 1", len(scenes))
			}
			if scenes[0].DashboardID != "framework_smoke" {
				t.Fatalf("DashboardID = %q, want framework_smoke", scenes[0].DashboardID)
			}
			if len(scenes[0].Widgets) != 4 {
				t.Fatalf("widget count = %d, want 4", len(scenes[0].Widgets))
			}
			widgetIDs := map[string]bool{}
			for _, widget := range scenes[0].Widgets {
				widgetIDs[widget.ID] = true
			}
			for _, id := range []string{"panel_backplate", "speed_digits", "rpm_digits", "engine_warning"} {
				if !widgetIDs[id] {
					t.Fatalf("missing widget %q in framework smoke scene", id)
				}
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if summary.VehicleID != "demo" {
		t.Fatalf("summary.VehicleID = %q, want demo", summary.VehicleID)
	}
	if summary.Pattern != PatternSweep {
		t.Fatalf("summary.Pattern = %q, want %q", summary.Pattern, PatternSweep)
	}
	if summary.SensorCount != 3 {
		t.Fatalf("summary.SensorCount = %d, want 3", summary.SensorCount)
	}
	if summary.DashboardCount != 1 {
		t.Fatalf("summary.DashboardCount = %d, want 1", summary.DashboardCount)
	}
	if summary.Events != 3 {
		t.Fatalf("summary.Events = %d, want 3", summary.Events)
	}
	if sceneUpdates != 1 {
		t.Fatalf("scene updates = %d, want 1", sceneUpdates)
	}
}

func TestOrnateTimberDashboardConfigRunsHarness(t *testing.T) {
	configPath := setupExampleHarnessEnvironment(t, filepath.Join("examples", "ornate-timber", "dashboard.yaml"), "demo")

	var sceneUpdates int
	summary, err := Run(context.Background(), Options{
		ConfigPath: configPath,
		VehicleID:  "demo",
		Pattern:    PatternSweep,
		MaxEvents:  3,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
		Sink: func(scenes []v3dashboard.Scene) error {
			sceneUpdates++
			if len(scenes) != 1 {
				t.Fatalf("scene count = %d, want 1", len(scenes))
			}
			if scenes[0].DashboardID != "ornate_timber" {
				t.Fatalf("DashboardID = %q, want ornate_timber", scenes[0].DashboardID)
			}
			if len(scenes[0].Widgets) != 7 {
				t.Fatalf("widget count = %d, want 7", len(scenes[0].Widgets))
			}
			widgetIDs := map[string]bool{}
			for _, widget := range scenes[0].Widgets {
				widgetIDs[widget.ID] = true
			}
			for _, id := range []string{
				"panel_backplate",
				"speed_numeric",
				"rpm_radial",
				"trip_odometer",
				"fuel_bar",
				"rpm_segmented",
				"check_engine_indicator",
			} {
				if !widgetIDs[id] {
					t.Fatalf("missing widget %q in ornate timber scene", id)
				}
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if summary.VehicleID != "demo" {
		t.Fatalf("summary.VehicleID = %q, want demo", summary.VehicleID)
	}
	if summary.Pattern != PatternSweep {
		t.Fatalf("summary.Pattern = %q, want %q", summary.Pattern, PatternSweep)
	}
	if summary.SensorCount != 6 {
		t.Fatalf("summary.SensorCount = %d, want 6", summary.SensorCount)
	}
	if summary.DashboardCount != 1 {
		t.Fatalf("summary.DashboardCount = %d, want 1", summary.DashboardCount)
	}
	if summary.Events != 3 {
		t.Fatalf("summary.Events = %d, want 3", summary.Events)
	}
	if sceneUpdates != 1 {
		t.Fatalf("scene updates = %d, want 1", sceneUpdates)
	}
}

func TestNeonGridDashboardConfigRunsHarness(t *testing.T) {
	configPath := setupExampleHarnessEnvironment(t, filepath.Join("examples", "neon-grid", "dashboard.yaml"), "demo")

	var sceneUpdates int
	summary, err := Run(context.Background(), Options{
		ConfigPath: configPath,
		VehicleID:  "demo",
		Pattern:    PatternSweep,
		MaxEvents:  3,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
		Sink: func(scenes []v3dashboard.Scene) error {
			sceneUpdates++
			if len(scenes) != 1 {
				t.Fatalf("scene count = %d, want 1", len(scenes))
			}
			if scenes[0].DashboardID != "neon_grid" {
				t.Fatalf("DashboardID = %q, want neon_grid", scenes[0].DashboardID)
			}
			if len(scenes[0].Widgets) != 7 {
				t.Fatalf("widget count = %d, want 7", len(scenes[0].Widgets))
			}
			widgetIDs := map[string]bool{}
			for _, widget := range scenes[0].Widgets {
				widgetIDs[widget.ID] = true
			}
			for _, id := range []string{
				"panel_backplate",
				"speed_numeric",
				"rpm_radial",
				"trip_odometer",
				"coolant_bar",
				"rpm_segmented",
				"check_engine_indicator",
			} {
				if !widgetIDs[id] {
					t.Fatalf("missing widget %q in neon-grid scene", id)
				}
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if summary.VehicleID != "demo" {
		t.Fatalf("summary.VehicleID = %q, want demo", summary.VehicleID)
	}
	if summary.Pattern != PatternSweep {
		t.Fatalf("summary.Pattern = %q, want %q", summary.Pattern, PatternSweep)
	}
	if summary.SensorCount != 6 {
		t.Fatalf("summary.SensorCount = %d, want 6", summary.SensorCount)
	}
	if summary.DashboardCount != 1 {
		t.Fatalf("summary.DashboardCount = %d, want 1", summary.DashboardCount)
	}
	if summary.Events != 3 {
		t.Fatalf("summary.Events = %d, want 3", summary.Events)
	}
	if sceneUpdates != 1 {
		t.Fatalf("scene updates = %d, want 1", sceneUpdates)
	}
}

func TestSteamScrapDashboardConfigRunsHarness(t *testing.T) {
	configPath := setupExampleHarnessEnvironment(t, filepath.Join("examples", "steam-scrap", "dashboard.yaml"), "demo")

	var sceneUpdates int
	summary, err := Run(context.Background(), Options{
		ConfigPath: configPath,
		VehicleID:  "demo",
		Pattern:    PatternSweep,
		MaxEvents:  3,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
		Sink: func(scenes []v3dashboard.Scene) error {
			sceneUpdates++
			if len(scenes) != 1 {
				t.Fatalf("scene count = %d, want 1", len(scenes))
			}
			if scenes[0].DashboardID != "steam_scrap" {
				t.Fatalf("DashboardID = %q, want steam_scrap", scenes[0].DashboardID)
			}
			if len(scenes[0].Widgets) != 7 {
				t.Fatalf("widget count = %d, want 7", len(scenes[0].Widgets))
			}
			widgetIDs := map[string]bool{}
			for _, widget := range scenes[0].Widgets {
				widgetIDs[widget.ID] = true
			}
			for _, id := range []string{
				"panel_backplate",
				"speed_numeric",
				"rpm_radial",
				"trip_odometer",
				"boiler_pressure_bar",
				"rpm_segmented",
				"boiler_warning_indicator",
			} {
				if !widgetIDs[id] {
					t.Fatalf("missing widget %q in steam-scrap scene", id)
				}
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if summary.VehicleID != "demo" {
		t.Fatalf("summary.VehicleID = %q, want demo", summary.VehicleID)
	}
	if summary.Pattern != PatternSweep {
		t.Fatalf("summary.Pattern = %q, want %q", summary.Pattern, PatternSweep)
	}
	if summary.SensorCount != 5 {
		t.Fatalf("summary.SensorCount = %d, want 5", summary.SensorCount)
	}
	if summary.DashboardCount != 1 {
		t.Fatalf("summary.DashboardCount = %d, want 1", summary.DashboardCount)
	}
	if summary.Events != 3 {
		t.Fatalf("summary.Events = %d, want 3", summary.Events)
	}
	if sceneUpdates != 1 {
		t.Fatalf("scene updates = %d, want 1", sceneUpdates)
	}
}

func TestSensorSourcesUseGaugeAwareProfilesForThemedDashboards(t *testing.T) {
	configPath := setupExampleHarnessEnvironment(t, filepath.Join("examples", "ornate-timber", "dashboard.yaml"), "demo")
	cfg, err := v3config.LoadFile(configPath)
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	plan, err := v3config.Resolve(cfg, "demo")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	searchPaths, err := v3assets.DefaultSearchPaths(configPath, plan.VehicleID)
	if err != nil {
		t.Fatalf("DefaultSearchPaths returned error: %v", err)
	}
	sources, err := sensorSources(plan, searchPaths)
	if err != nil {
		t.Fatalf("sensorSources returned error: %v", err)
	}
	profiles := sensorProfilesByID(sources)

	if profiles["speed"] != sweepProfileIncremental {
		t.Fatalf("speed profile = %q, want %q", profiles["speed"], sweepProfileIncremental)
	}
	if profiles["trip_distance"] != sweepProfileIncremental {
		t.Fatalf("trip_distance profile = %q, want %q", profiles["trip_distance"], sweepProfileIncremental)
	}
	if profiles["fuel_level"] != sweepProfileHeartbeat {
		t.Fatalf("fuel_level profile = %q, want %q", profiles["fuel_level"], sweepProfileHeartbeat)
	}
	if profiles["check_engine"] != sweepProfileIndicator {
		t.Fatalf("check_engine profile = %q, want %q", profiles["check_engine"], sweepProfileIndicator)
	}
	if profiles["rpm"] != sweepProfileRange {
		t.Fatalf("rpm profile = %q, want %q", profiles["rpm"], sweepProfileRange)
	}
}

func TestSensorSourcesPreferRangeSweepForSharedBaselineRPMSensor(t *testing.T) {
	configPath := setupBaselineHarnessEnvironment(t)
	cfg, err := v3config.LoadFile(configPath)
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	plan, err := v3config.Resolve(cfg, "vw_caddy")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	searchPaths, err := v3assets.DefaultSearchPaths(configPath, plan.VehicleID)
	if err != nil {
		t.Fatalf("DefaultSearchPaths returned error: %v", err)
	}
	sources, err := sensorSources(plan, searchPaths)
	if err != nil {
		t.Fatalf("sensorSources returned error: %v", err)
	}
	profiles := sensorProfilesByID(sources)

	if profiles["rpm"] != sweepProfileRange {
		t.Fatalf("rpm profile = %q, want %q", profiles["rpm"], sweepProfileRange)
	}
	if profiles["speed"] != sweepProfileIncremental {
		t.Fatalf("speed profile = %q, want %q", profiles["speed"], sweepProfileIncremental)
	}
}

func BenchmarkBaselineDashboardHarnessPatterns(b *testing.B) {
	configPath := setupBaselineHarnessEnvironment(b)
	patterns := []string{PatternFixed, PatternSweep, PatternHeartbeat}

	for _, pattern := range patterns {
		pattern := pattern
		b.Run(pattern, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Run(context.Background(), Options{
					ConfigPath: configPath,
					VehicleID:  "vw_caddy",
					Pattern:    pattern,
					MaxEvents:  3,
					Now: func() time.Time {
						return time.Unix(0, 0)
					},
					Sink: func([]v3dashboard.Scene) error {
						return nil
					},
				})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func setupBaselineHarnessEnvironment(tb testing.TB) string {
	return setupExampleHarnessEnvironment(tb, filepath.Join("examples", "baseline-dashboard.yaml"), "vw_caddy")
}

func setupExampleHarnessEnvironment(tb testing.TB, relativeConfigPath, vehicleID string) string {
	tb.Helper()
	repoRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		tb.Fatal(err)
	}
	previousWorkingDirectory, err := os.Getwd()
	if err != nil {
		tb.Fatal(err)
	}
	if err := os.Chdir(repoRoot); err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		if err := os.Chdir(previousWorkingDirectory); err != nil {
			tb.Fatalf("restore working directory: %v", err)
		}
	})

	configPath := filepath.Join(repoRoot, relativeConfigPath)
	previousArgs := append([]string(nil), os.Args...)
	os.Args = []string{"GoDriveLog.test", "--harness", "--config", configPath, "--vehicle", vehicleID}
	tb.Cleanup(func() {
		os.Args = previousArgs
	})

	return configPath
}

func writeHarnessConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"off", "on", "unknown"} {
		if err := writeTestPNG(filepath.Join(dir, "assets", name+".png")); err != nil {
			t.Fatal(err)
		}
	}
	configPath := filepath.Join(dir, "config.yaml")
	config := `vehicles:
  demo:
    name: Demo vehicle
    obd:
      address: serial:///dev/null
      timeout: 1000
    dashboards:
      - primary
sensors:
  pulse:
    type: obd
    pid: "010C"
    unit: rpm
    poll: 100
    min: 0
    max: 100
  temp:
    type: obd
    pid: "0105"
    unit: C
    poll: 100
    min: 0
    max: 100
assets:
  digit_sets: {}
  bar_sets: {}
  frame_sets: {}
  image_sets: {}
  indicator_sets:
    pulse_indicator:
      states:
        "off": assets/off.png
        "on": assets/on.png
        "unknown": assets/unknown.png
logs: {}
dashboards:
  primary:
    display: main
    size:
      width: 32
      height: 16
    widgets:
      - id: pulse_widget
        type: indicator
        sensor: pulse
        asset: pulse_indicator
        position: [0, 0]
      - id: temp_widget
        type: indicator
        sensor: temp
        asset: pulse_indicator
        position: [8, 0]
`
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	return configPath
}

func writeTestPNG(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	return png.Encode(file, img)
}

func sensorProfilesByID(sources []sensorSource) map[string]sweepProfile {
	profiles := map[string]sweepProfile{}
	for _, source := range sources {
		profiles[source.ID] = source.SweepProfile
	}
	return profiles
}
