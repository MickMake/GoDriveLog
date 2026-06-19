package v3config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDocsV3ConfigExample(t *testing.T) {
	cfg := loadTestConfig(t, "config.example.yaml")
	if len(cfg.Vehicles) != 1 {
		t.Fatalf("expected one vehicle, got %d", len(cfg.Vehicles))
	}
	if _, ok := cfg.Sensors["speed"]; !ok {
		t.Fatalf("expected global speed sensor")
	}
	if _, ok := cfg.Dashboards["simple_primary"]; !ok {
		t.Fatalf("expected simple_primary dashboard")
	}
}

func TestLoadDocsV3ConfigFull(t *testing.T) {
	cfg := loadTestConfig(t, "config.full.yaml")
	if len(cfg.Vehicles) < 2 {
		t.Fatalf("expected multiple vehicle profiles, got %d", len(cfg.Vehicles))
	}
	if err := ValidateSelectedVehicle(cfg, "bench_z31"); err != nil {
		t.Fatalf("expected selected bench_z31 vehicle to validate: %v", err)
	}
	if err := ValidateSelectedVehicle(cfg, ""); err == nil {
		t.Fatalf("expected explicit selected vehicle to be required for multiple vehicles")
	}
}

func TestLoadDocsV3StandaloneExamples(t *testing.T) {
	for _, path := range []string{
		"examples/simple_speed_warning.yaml",
		"examples/nissan_300zx_z31_inspired.yaml",
		"examples/honda_s2000_inspired.yaml",
	} {
		t.Run(path, func(t *testing.T) {
			loadTestConfig(t, path)
		})
	}
}

func TestRejectUnknownRootField(t *testing.T) {
	_, err := LoadBytes([]byte(validMinimalYAML() + "\nvehicle:\n  name: old shape\n"))
	if err == nil {
		t.Fatalf("expected unknown root field to fail")
	}
	assertErrorContains(t, err, "field vehicle not found")
}

func TestRejectUnknownNestedField(t *testing.T) {
	bad := strings.Replace(validMinimalYAML(), "timeout: 1000", "timeout: 1000\n      provider: elm327", 1)
	_, err := LoadBytes([]byte(bad))
	if err == nil {
		t.Fatalf("expected unknown nested field to fail")
	}
	assertErrorContains(t, err, "field provider not found")
}

func TestRejectBadVehicleReferences(t *testing.T) {
	bad := strings.Replace(validMinimalYAML(), "- jsonl", "- missing_log", 1)
	_, err := LoadBytes([]byte(bad))
	if err == nil {
		t.Fatalf("expected bad vehicle log reference to fail")
	}
	assertErrorContains(t, err, "missing_log")
}

func TestRejectBadLogSensorReference(t *testing.T) {
	bad := strings.Replace(validMinimalYAML(), "- speed", "- missing_sensor", 1)
	_, err := LoadBytes([]byte(bad))
	if err == nil {
		t.Fatalf("expected bad log sensor reference to fail")
	}
	assertErrorContains(t, err, "missing_sensor")
}

func TestRejectBadWidgetAssetFamily(t *testing.T) {
	bad := strings.Replace(validMinimalYAML(), "asset: digits", "asset: panel", 1)
	_, err := LoadBytes([]byte(bad))
	if err == nil {
		t.Fatalf("expected wrong widget asset family to fail")
	}
	assertErrorContains(t, err, "assets.digit_sets")
}

func TestAllowGaugeWidgetWithoutSensor(t *testing.T) {
	cfg, err := LoadBytes([]byte(validMinimalYAMLWithGaugeWidget("")))
	if err != nil {
		t.Fatalf("expected gauge widget config to load: %v", err)
	}

	widgets := cfg.Dashboards["simple_primary"].Widgets
	gauge := widgets[len(widgets)-1]
	if gauge.Type != WidgetTypeGauge || gauge.Sensor != "" || gauge.Asset != "" {
		t.Fatalf("gauge widget identity = %#v", gauge)
	}
	if gauge.Gauge != "assets/gauges/7Seg/amber/4_digit_rpm" || gauge.Scale != 1.0 {
		t.Fatalf("gauge widget placement = %#v", gauge)
	}
}

func TestRejectGaugeWidgetSensor(t *testing.T) {
	_, err := LoadBytes([]byte(validMinimalYAMLWithGaugeWidget("        sensor: speed\n")))
	if err == nil {
		t.Fatalf("expected gauge widget sensor to fail")
	}
	assertErrorContains(t, err, "sensor must be empty for gauge widgets")
}

func TestRejectInvalidIDs(t *testing.T) {
	bad := strings.Replace(validMinimalYAML(), "vw_caddy:", "VW-Caddy:", 1)
	_, err := LoadBytes([]byte(bad))
	if err == nil {
		t.Fatalf("expected invalid ID to fail")
	}
	assertErrorContains(t, err, "must match")
}

func TestRejectAssetPathUpwardEscapes(t *testing.T) {
	for _, badPath := range []string{
		"..",
		"assets/dashboard/..",
		"assets/dashboard/../x",
	} {
		t.Run(badPath, func(t *testing.T) {
			bad := strings.Replace(validMinimalYAML(), "assets/dashboard/simple/panel/background.png", badPath, 1)
			_, err := LoadBytes([]byte(bad))
			if err == nil {
				t.Fatalf("expected asset path %q to fail", badPath)
			}
			assertErrorContains(t, err, "repository-root relative")
		})
	}
}

func TestRejectDecimalFormatsWithoutDecimalPoint(t *testing.T) {
	for _, format := range []string{"%f", "%03f"} {
		t.Run(format, func(t *testing.T) {
			bad := removeDigitDecimalPoint(validMinimalYAML())
			bad = strings.Replace(bad, `format: "%03.0f"`, `format: "`+format+`"`, 1)
			_, err := LoadBytes([]byte(bad))
			if err == nil {
				t.Fatalf("expected format %q without decimal_point to fail", format)
			}
			assertErrorContains(t, err, "decimal_point")
		})
	}
}

func TestAllowZeroPrecisionFormatsWithoutDecimalPoint(t *testing.T) {
	for _, format := range []string{"%.0f", "%03.0f"} {
		t.Run(format, func(t *testing.T) {
			cfgText := removeDigitDecimalPoint(validMinimalYAML())
			cfgText = strings.Replace(cfgText, `format: "%03.0f"`, `format: "`+format+`"`, 1)
			if _, err := LoadBytes([]byte(cfgText)); err != nil {
				t.Fatalf("expected format %q without decimal_point to pass: %v", format, err)
			}
		})
	}
}

func loadTestConfig(t *testing.T, repoPath string) Config {
	t.Helper()
	path := filepath.Join("..", "..", "..", filepath.FromSlash(repoPath))
	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile(%s) failed: %v", repoPath, err)
	}
	return cfg
}

func assertErrorContains(t *testing.T, err error, want string) {
	t.Helper()
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected error to contain %q, got %q", want, err.Error())
	}
}

func removeDigitDecimalPoint(cfgText string) string {
	return strings.Replace(cfgText, "      decimal_point: assets/dashboard/simple/digits/dp.png\n", "", 1)
}

func validMinimalYAMLWithGaugeWidget(extraLines string) string {
	return validMinimalYAML() + `
      - id: rpm_gauge
        type: gauge
` + extraLines + `        gauge: assets/gauges/7Seg/amber/4_digit_rpm
        position: [780, 40]
        scale: 1.0`
}

func validMinimalYAML() string {
	return `vehicles:
  vw_caddy:
    name: VW Caddy
    obd:
      address: serial:///dev/ttyUSB0
      timeout: 1000
    logs:
      - jsonl
    dashboards:
      - simple_primary
sensors:
  speed:
    type: obd
    pid: "010D"
    unit: km/h
    poll: 250
    min: 0
    max: 220
assets:
  digit_sets:
    digits:
      characters:
        "0": assets/dashboard/simple/digits/0.png
        "1": assets/dashboard/simple/digits/1.png
        "2": assets/dashboard/simple/digits/2.png
        "3": assets/dashboard/simple/digits/3.png
        "4": assets/dashboard/simple/digits/4.png
        "5": assets/dashboard/simple/digits/5.png
        "6": assets/dashboard/simple/digits/6.png
        "7": assets/dashboard/simple/digits/7.png
        "8": assets/dashboard/simple/digits/8.png
        "9": assets/dashboard/simple/digits/9.png
      decimal_point: assets/dashboard/simple/digits/dp.png
  image_sets:
    panel:
      image: assets/dashboard/simple/panel/background.png
logs:
  jsonl:
    path: logs/godrivelog.jsonl
    sensors:
      - speed
dashboards:
  simple_primary:
    display: HDMI-1
    size:
      width: 800
      height: 480
    widgets:
      - id: panel_backplate
        type: image
        asset: panel
        position: [0, 0]
      - id: speed_digits
        type: digit_display
        sensor: speed
        asset: digits
        position: [40, 40]
        digits: 3
        format: "%03.0f"`
}
