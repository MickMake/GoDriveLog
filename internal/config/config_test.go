package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDashboardV2Config(t *testing.T) {
	cfg := loadConfig(t, `
mock_mode: true
log:
  directory: ./test-log
vehicle:
  name: Test Van
dashboard:
  canvas:
    width: 800
    height: 480
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    refresh: 250
    min: 0
    max: 7000
    log: true
`)

	if cfg.Vehicle.Name != "Test Van" {
		t.Fatalf("Vehicle.Name = %q, want Test Van", cfg.Vehicle.Name)
	}
	if cfg.Dashboard.Canvas.Width != 800 || cfg.Dashboard.Canvas.Height != 480 {
		t.Fatalf("Dashboard.Canvas = %#v, want 800x480", cfg.Dashboard.Canvas)
	}
	if len(cfg.Sensors) != 1 {
		t.Fatalf("len(Sensors) = %d, want 1", len(cfg.Sensors))
	}
	if cfg.Log.Rotate != DefaultLogRotate {
		t.Fatalf("Log.Rotate = %q, want %q", cfg.Log.Rotate, DefaultLogRotate)
	}
	if cfg.OBDAddress != DefaultOBDAddress {
		t.Fatalf("OBDAddress = %q, want %q", cfg.OBDAddress, DefaultOBDAddress)
	}
}

func TestLoadRejectsOldPIDOnlyConfig(t *testing.T) {
	_, err := loadConfigFile(t, `
mock_mode: true
vehicle:
  name: Test Van
  pids:
    rpm:
      type: obd
      pid: "010C"
      unit: rpm
      refresh: 250
      min: 0
      max: 7000
      log: true
      display:
        enabled: true
        widget: radial1
`)
	if err == nil {
		t.Fatal("Load succeeded for old PID/display-only config, want error")
	}
}

func TestActiveSensorsUsesTopLevelSensors(t *testing.T) {
	cfg := Config{
		Sensors: map[string]SensorConfig{
			"rpm": {
				Type:    "obd",
				PID:     "010C",
				Unit:    "rpm",
				Refresh: 250,
				Min:     0,
				Max:     7000,
				Log:     true,
			},
			"speed": {
				Type:    "obd",
				PID:     "010D",
				Unit:    "km/h",
				Refresh: 500,
				Min:     0,
				Max:     160,
				Log:     false,
			},
			"moving": {
				Type:    "virtual",
				Unit:    "bool",
				Refresh: 1000,
				Min:     0,
				Max:     1,
				Log:     true,
			},
		},
	}

	active := ActiveSensors(cfg)
	if len(active) != 1 {
		t.Fatalf("len(active) = %d, want 1", len(active))
	}
	if active[0].Key != "rpm" || active[0].RawPID != "010C" {
		t.Fatalf("active[0] = %#v, want rpm/010C", active[0])
	}
}

func TestValidateRequiresDashboardCanvas(t *testing.T) {
	_, err := loadConfigFile(t, `
mock_mode: true
vehicle:
  name: Test Van
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    refresh: 250
    min: 0
    max: 7000
    log: true
`)
	if err == nil {
		t.Fatal("Load succeeded without dashboard canvas, want error")
	}
}

func loadConfig(t *testing.T, content string) Config {
	t.Helper()

	cfg, err := loadConfigFile(t, content)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	return cfg
}

func loadConfigFile(t *testing.T, content string) (Config, error) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return Load(path)
}
