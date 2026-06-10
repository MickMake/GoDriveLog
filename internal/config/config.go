package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultLogRotate            = "daily"
	DefaultLogDirectory         = "./log"
	DefaultOBDProvider          = OBDProviderOBD
	DefaultOBDAddress           = "serial:///dev/ttyUSB0"
	DefaultDashboardRefreshMS   = 100
	DefaultDashboardRenderMinMS = 0

	OBDProviderOBD      = "obd"
	OBDProviderMock     = "mock"
	OBDProviderRaceDemo = "race-demo"
)

type Config struct {
	OBD       OBDConfig               `yaml:"obd"`
	Log       LogConfig               `yaml:"log"`
	Vehicle   VehicleConfig           `yaml:"vehicle"`
	Sensors   map[string]SensorConfig `yaml:"sensors"`
	Dashboard DashboardConfig         `yaml:"dashboard"`
}

type OBDConfig struct {
	Provider string `yaml:"provider"`
	MockMode bool   `yaml:"mock_mode"`
	Address  string `yaml:"address"`
	Debug    bool   `yaml:"debug"`
}

type LogConfig struct {
	Rotate    string `yaml:"rotate"`
	Directory string `yaml:"directory"`
}

type VehicleConfig struct {
	Name string `yaml:"name"`
}

type SensorConfig struct {
	Type    string  `yaml:"type"`
	PID     string  `yaml:"pid"`
	Unit    string  `yaml:"unit"`
	Refresh int     `yaml:"refresh"`
	Min     float64 `yaml:"min"`
	Max     float64 `yaml:"max"`
	Log     bool    `yaml:"log"`
	Display bool    `yaml:"display"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	applyDefaults(&cfg)
	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func NormalizeOBDProvider(provider string, mockMode bool) string {
	provider = strings.TrimSpace(strings.ToLower(provider))
	switch provider {
	case "":
		if mockMode {
			return OBDProviderMock
		}
		return DefaultOBDProvider
	case OBDProviderOBD, OBDProviderMock, OBDProviderRaceDemo:
		return provider
	case "race_demo", "racedemo":
		return OBDProviderRaceDemo
	default:
		return provider
	}
}

func applyDefaults(cfg *Config) {
	cfg.OBD.Provider = NormalizeOBDProvider(cfg.OBD.Provider, cfg.OBD.MockMode)

	if cfg.Log.Rotate == "" {
		cfg.Log.Rotate = DefaultLogRotate
	}
	if cfg.Log.Directory == "" {
		cfg.Log.Directory = DefaultLogDirectory
	}
	if cfg.OBD.Address == "" {
		cfg.OBD.Address = DefaultOBDAddress
	}
	if cfg.Dashboard.RefreshMS == 0 {
		cfg.Dashboard.RefreshMS = DefaultDashboardRefreshMS
	}
	if cfg.Dashboard.RenderMinMS == 0 {
		cfg.Dashboard.RenderMinMS = DefaultDashboardRenderMinMS
	}
}

func validate(cfg Config) error {
	if cfg.Vehicle.Name == "" {
		return fmt.Errorf("vehicle.name must not be empty")
	}
	if len(cfg.Sensors) == 0 {
		return fmt.Errorf("sensors must not be empty")
	}
	if cfg.Log.Rotate != DefaultLogRotate {
		return fmt.Errorf("log.rotate must be %q", DefaultLogRotate)
	}
	if cfg.Log.Directory == "" {
		return fmt.Errorf("log.directory must not be empty")
	}
	switch cfg.OBD.Provider {
	case OBDProviderOBD, OBDProviderMock, OBDProviderRaceDemo:
	default:
		return fmt.Errorf("obd.provider must be %q, %q, or %q", OBDProviderOBD, OBDProviderMock, OBDProviderRaceDemo)
	}
	if cfg.OBD.Address == "" {
		return fmt.Errorf("obd.address must not be empty")
	}

	for key, sensor := range cfg.Sensors {
		if sensor.Type != "obd" && sensor.Type != "virtual" {
			return fmt.Errorf("sensors.%s.type must be obd or virtual", key)
		}
		if sensor.Type == "obd" && sensor.PID == "" {
			return fmt.Errorf("sensors.%s.pid must not be empty for obd sensors", key)
		}
		if sensor.Unit == "" {
			return fmt.Errorf("sensors.%s.unit must not be empty", key)
		}
		if sensor.Refresh <= 0 {
			return fmt.Errorf("sensors.%s.refresh must be positive", key)
		}
		if sensor.Max <= sensor.Min {
			return fmt.Errorf("sensors.%s.max must be greater than min", key)
		}
	}

	if err := validateDashboard(cfg); err != nil {
		return err
	}

	return nil
}
