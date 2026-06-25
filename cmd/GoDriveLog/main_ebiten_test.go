//go:build !fyne_legacy

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	v3harness "github.com/MickMake/GoDriveLog/internal/dashboard/harness"
)

func TestDashboardRunDiscoversSingleVehicleConfig(t *testing.T) {
	root := t.TempDir()
	writeTestConfig(t, filepath.Join(root, "dashboard.yaml"), singleVehicleConfigYAML("demo"))
	restoreWD := changeWorkingDirectory(t, root)
	defer restoreWD()

	var gotConfig string
	var gotVehicle string
	var gotDuration time.Duration
	previousRun := dashboardRunCommand
	dashboardRunCommand = func(configPath, vehicleID string, duration time.Duration) error {
		gotConfig = configPath
		gotVehicle = vehicleID
		gotDuration = duration
		return nil
	}
	defer func() {
		dashboardRunCommand = previousRun
	}()

	if err := runCLI([]string{"dashboard", "run"}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatalf("runCLI returned error: %v", err)
	}
	if filepath.Base(gotConfig) != "dashboard.yaml" {
		t.Fatalf("config path = %q, want dashboard.yaml", gotConfig)
	}
	if gotVehicle != "demo" {
		t.Fatalf("vehicle id = %q, want demo", gotVehicle)
	}
	if gotDuration != 0 {
		t.Fatalf("duration = %s, want 0", gotDuration)
	}
}

func TestDashboardRunStopsAtFirstValidMultiVehicleConfigWithoutVehicle(t *testing.T) {
	root := t.TempDir()
	writeTestConfig(t, filepath.Join(root, "config.yaml"), multiVehicleConfigYAML())
	writeTestConfig(t, filepath.Join(root, "godrivelog.yaml"), singleVehicleConfigYAML("later"))
	restoreWD := changeWorkingDirectory(t, root)
	defer restoreWD()

	err := runCLI([]string{"dashboard", "run"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected discovery to stop on the first valid multi-vehicle config")
	}
	if !strings.Contains(err.Error(), "config.yaml") || !strings.Contains(err.Error(), "bench_z31") || !strings.Contains(err.Error(), "vw_caddy") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDashboardRunSearchesConfigsForRequestedVehicle(t *testing.T) {
	root := t.TempDir()
	writeTestConfig(t, filepath.Join(root, "config.yaml"), singleVehicleConfigYAML("demo"))
	writeTestConfig(t, filepath.Join(root, "godrivelog.yaml"), multiVehicleConfigYAML())
	restoreWD := changeWorkingDirectory(t, root)
	defer restoreWD()

	var gotConfig string
	var gotVehicle string
	previousRun := dashboardRunCommand
	dashboardRunCommand = func(configPath, vehicleID string, duration time.Duration) error {
		gotConfig = configPath
		gotVehicle = vehicleID
		return nil
	}
	defer func() {
		dashboardRunCommand = previousRun
	}()

	if err := runCLI([]string{"dashboard", "run", "bench_z31"}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatalf("runCLI returned error: %v", err)
	}
	if filepath.Base(gotConfig) != "godrivelog.yaml" {
		t.Fatalf("config path = %q, want godrivelog.yaml", gotConfig)
	}
	if gotVehicle != "bench_z31" {
		t.Fatalf("vehicle id = %q, want bench_z31", gotVehicle)
	}
}

func TestDashboardHarnessUsesPositionalVehicleAndDefaults(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "dashboard.yaml")
	writeTestConfig(t, configPath, multiVehicleConfigYAML())

	var gotConfig string
	var gotVehicle string
	var gotPattern string
	var gotInterval time.Duration
	var gotDuration time.Duration
	previousHarness := dashboardHarnessCommand
	dashboardHarnessCommand = func(configPath, vehicleID, pattern string, interval, duration time.Duration) error {
		gotConfig = configPath
		gotVehicle = vehicleID
		gotPattern = pattern
		gotInterval = interval
		gotDuration = duration
		return nil
	}
	defer func() {
		dashboardHarnessCommand = previousHarness
	}()

	if err := runCLI([]string{"dashboard", "harness", "bench_z31", "--config", configPath}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatalf("runCLI returned error: %v", err)
	}
	if gotConfig != configPath {
		t.Fatalf("config path = %q, want %q", gotConfig, configPath)
	}
	if gotVehicle != "bench_z31" {
		t.Fatalf("vehicle id = %q, want bench_z31", gotVehicle)
	}
	if gotPattern != v3harness.PatternSweep {
		t.Fatalf("pattern = %q, want %q", gotPattern, v3harness.PatternSweep)
	}
	if gotInterval != 100*time.Millisecond {
		t.Fatalf("interval = %s, want 100ms", gotInterval)
	}
	if gotDuration != 0 {
		t.Fatalf("duration = %s, want 0", gotDuration)
	}
}

func TestDashboardValidateRejectsPositionalAndFlagConfig(t *testing.T) {
	err := runCLI([]string{"dashboard", "validate", "./one.yaml", "--config", "./two.yaml"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected conflicting config inputs to fail")
	}
	if !strings.Contains(err.Error(), "either a positional config file or --config") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDashboardExamplesExportsBuiltInTheme(t *testing.T) {
	repoRoot := repoRootFromTestFile(t)
	outputDir := filepath.Join(t.TempDir(), "framework-smoke")
	restoreWD := changeWorkingDirectory(t, repoRoot)
	defer restoreWD()

	stdout := &bytes.Buffer{}
	if err := runCLI([]string{"dashboard", "examples", "--theme", frameworkSmokeTheme, "--output", outputDir}, stdout, &bytes.Buffer{}); err != nil {
		t.Fatalf("runCLI returned error: %v", err)
	}
	for _, relative := range []string{
		"dashboard.yaml",
		filepath.Join("assets", "panel", "background.png"),
		filepath.Join("assets", "indicator", "lamp_on.png"),
	} {
		path := filepath.Join(outputDir, relative)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
	if !strings.Contains(stdout.String(), "exported dashboard example") {
		t.Fatalf("stdout = %q, want export confirmation", stdout.String())
	}
}

func TestDashboardExamplesFailsForNonEmptyOutputWithoutForce(t *testing.T) {
	repoRoot := repoRootFromTestFile(t)
	outputDir := filepath.Join(t.TempDir(), "framework-smoke")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "keep.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	restoreWD := changeWorkingDirectory(t, repoRoot)
	defer restoreWD()

	err := runCLI([]string{"dashboard", "examples", "--theme", frameworkSmokeTheme, "--output", outputDir}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected non-empty output directory to fail without --force")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDashboardHelpOutputsIncludeNewCommandTree(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "dashboard",
			args: []string{"dashboard", "--help"},
			want: []string{"GoDriveLog dashboard", "run", "harness", "examples", "validate"},
		},
		{
			name: "run",
			args: []string{"dashboard", "run", "--help"},
			want: []string{"GoDriveLog dashboard run [vehicle-id]", "--config", "--renderer", "--duration"},
		},
		{
			name: "harness",
			args: []string{"dashboard", "harness", "--help"},
			want: []string{"GoDriveLog dashboard harness [vehicle-id]", "--pattern", "--interval", "--duration"},
		},
		{
			name: "examples",
			args: []string{"dashboard", "examples", "--help"},
			want: []string{"GoDriveLog dashboard examples --output <directory>", "--output", "--theme", "--force"},
		},
		{
			name: "validate",
			args: []string{"dashboard", "validate", "--help"},
			want: []string{"GoDriveLog dashboard validate [config-file]", "--config"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			if err := runCLI(test.args, stdout, &bytes.Buffer{}); err != nil {
				t.Fatalf("runCLI returned error: %v", err)
			}
			for _, want := range test.want {
				if !strings.Contains(stdout.String(), want) {
					t.Fatalf("help output missing %q\n%s", want, stdout.String())
				}
			}
		})
	}
}

func repoRootFromTestFile(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func changeWorkingDirectory(t *testing.T, path string) func() {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(path); err != nil {
		t.Fatalf("Chdir(%s): %v", path, err)
	}
	return func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore Chdir(%s): %v", previous, err)
		}
	}
}

func writeTestConfig(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func singleVehicleConfigYAML(vehicleID string) string {
	return `vehicles:
  ` + vehicleID + `:
    name: Demo
    obd:
      address: serial:///dev/ttyUSB0
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - primary
sensors:
  speed:
    type: obd
    pid: "010D"
    unit: km/h
    poll: 250
assets:
  image_sets:
    panel:
      image: assets/panel.png
logs:
  jsonl:
    path: logs/demo.jsonl
    sensors:
      - speed
dashboards:
  primary:
    size:
      width: 800
      height: 480
    widgets:
      - id: panel_backplate
        type: image
        asset: panel
        position: [0, 0]
`
}

func multiVehicleConfigYAML() string {
	return `vehicles:
  bench_z31:
    name: Bench Z31
    obd:
      address: tcp://127.0.0.1:35000
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - primary
  vw_caddy:
    name: VW Caddy
    obd:
      address: serial:///dev/ttyUSB0
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - primary
sensors:
  speed:
    type: obd
    pid: "010D"
    unit: km/h
    poll: 250
assets:
  image_sets:
    panel:
      image: assets/panel.png
logs:
  jsonl:
    path: logs/demo.jsonl
    sensors:
      - speed
dashboards:
  primary:
    size:
      width: 800
      height: 480
    widgets:
      - id: panel_backplate
        type: image
        asset: panel
        position: [0, 0]
`
}
