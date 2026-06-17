package harness

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestRunFeedsFakeEventsThroughDashboardScenePath(t *testing.T) {
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
sensors:
  pulse:
    type: obd
    pid: "010C"
    unit: rpm
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
        off: assets/off.png
        on: assets/on.png
        unknown: assets/unknown.png
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
`
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	var sceneUpdates int
	summary, err := Run(context.Background(), Options{
		ConfigPath: configPath,
		RepoRoot:   dir,
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
	if sceneUpdates == 0 {
		t.Fatal("harness did not emit any dashboard scenes")
	}
}

func writeTestPNG(path string) error {
	const png1x1 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGJgYGAAAAAEAAGjChXjAAAAAElFTkSuQmCC"
	data, err := base64.StdEncoding.DecodeString(png1x1)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
