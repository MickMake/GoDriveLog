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

func TestLoadDashboardV21Config(t *testing.T) {
	cfg := loadConfig(t, validDashboardV21Config())

	if len(cfg.Dashboard.Assets) != 3 {
		t.Fatalf("len(Dashboard.Assets) = %d, want 3", len(cfg.Dashboard.Assets))
	}
	if len(cfg.Dashboard.Decoders) != 3 {
		t.Fatalf("len(Dashboard.Decoders) = %d, want 3", len(cfg.Dashboard.Decoders))
	}
	if len(cfg.Dashboard.Blocks) != 4 {
		t.Fatalf("len(Dashboard.Blocks) = %d, want 4", len(cfg.Dashboard.Blocks))
	}
	if len(cfg.Dashboard.Layers) != 1 {
		t.Fatalf("len(Dashboard.Layers) = %d, want 1", len(cfg.Dashboard.Layers))
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

func TestDashboardV21ValidationRejectsInvalidConfigs(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "duplicate asset id",
			content: replaceConfigText(validDashboardV21Config(), `
    path: assets/dashboard/background.png
`, `
    path: assets/dashboard/background.png
  - id: background
    type: image
    path: assets/dashboard/duplicate.png
`),
		},
		{
			name: "missing asset id",
			content: replaceConfigText(validDashboardV21Config(), `id: background`, `id: ""`),
		},
		{
			name: "unknown asset type",
			content: replaceConfigText(validDashboardV21Config(), `type: image`, `type: nonsense`),
		},
		{
			name: "unknown decoder type",
			content: replaceConfigText(validDashboardV21Config(), `type: digits`, `type: nonsense`),
		},
		{
			name: "decoder missing sensor reference",
			content: replaceConfigText(validDashboardV21Config(), `sensor: rpm`, `sensor: missing_sensor`),
		},
		{
			name: "decoder missing asset reference",
			content: replaceConfigText(validDashboardV21Config(), `asset: throttle_frames`, `asset: missing_asset`),
		},
		{
			name: "invalid frame count",
			content: replaceConfigText(validDashboardV21Config(), `frame_count: 10`, `frame_count: 0`),
		},
		{
			name: "unknown block type",
			content: replaceConfigText(validDashboardV21Config(), `type: sprite_text`, `type: nonsense`),
		},
		{
			name: "block missing asset reference",
			content: replaceConfigText(validDashboardV21Config(), `asset: background`, `asset: missing_asset`),
		},
		{
			name: "block missing decoder reference",
			content: replaceConfigText(validDashboardV21Config(), `decoder: throttle_frame`, `decoder: missing_decoder`),
		},
		{
			name: "invalid geometry",
			content: replaceConfigText(validDashboardV21Config(), `      width: 800`, `      width: 0`),
		},
		{
			name: "group missing child block reference",
			content: replaceConfigText(validDashboardV21Config(), `- throttle_bar`, `- missing_block`),
		},
		{
			name: "layer missing block reference",
			content: replaceConfigText(validDashboardV21Config(), `- main_cluster`, `- missing_block`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadConfigFile(t, tt.content)
			if err == nil {
				t.Fatal("Load succeeded, want validation error")
			}
		})
	}
}

func validDashboardV21Config() string {
	return `
mock_mode: true
vehicle:
  name: Test Van
dashboard:
  canvas:
    width: 800
    height: 480
  assets:
  - id: background
    type: image
    path: assets/dashboard/background.png
  - id: throttle_frames
    type: frame_set
    frames:
    - assets/dashboard/throttle/000.png
    - assets/dashboard/throttle/001.png
  - id: yellow_digits
    type: charset
    glyphs:
      "0": assets/dashboard/digits/yellow/0.png
      "1": assets/dashboard/digits/yellow/1.png
  decoders:
  - id: rpm_digits
    type: digits
    sensor: rpm
    asset: yellow_digits
  - id: throttle_frame
    type: frame_index
    sensor: throttle_position
    asset: throttle_frames
    frame_count: 10
  - id: rpm_warning
    type: threshold
    sensor: rpm
    thresholds:
    - at: 6000
      value: warning
  blocks:
  - id: background_panel
    type: image
    asset: background
    geometry:
      x: 0
      y: 0
      width: 800
      height: 480
  - id: rpm_display
    type: sprite_text
    asset: yellow_digits
    decoder: rpm_digits
    geometry:
      x: 100
      y: 60
      width: 240
      height: 80
  - id: throttle_bar
    type: sprite_frame
    asset: throttle_frames
    decoder: throttle_frame
    geometry:
      x: 100
      y: 170
      width: 300
      height: 40
  - id: main_cluster
    type: group
    blocks:
    - rpm_display
    - throttle_bar
  layers:
  - id: base
    z: 0
    blocks:
    - background_panel
    - main_cluster
sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    refresh: 250
    min: 0
    max: 7000
    log: true
  throttle_position:
    type: obd
    pid: "0111"
    unit: "%"
    refresh: 500
    min: 0
    max: 100
    log: true
`
}

func replaceConfigText(config string, old string, new string) string {
	for i := 0; i+len(old) <= len(config); i++ {
		if config[i:i+len(old)] == old {
			return config[:i] + new + config[i+len(old):]
		}
	}
	return config
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
