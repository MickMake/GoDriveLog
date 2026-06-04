package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultLogRotate   = "daily"
	DefaultLogDirectory = "./log"
	DefaultOBDAddress  = "serial:///dev/ttyUSB0"
)

type Config struct {
	MockMode   bool          `yaml:"mock_mode"`
	OBDAddress string        `yaml:"obd_address"`
	OBDDebug   bool          `yaml:"obd_debug"`
	Log        LogConfig     `yaml:"log"`
	Vehicle    VehicleConfig `yaml:"vehicle"`
}

type LogConfig struct {
	Rotate    string `yaml:"rotate"`
	Directory string `yaml:"directory"`
}

type VehicleConfig struct {
	Name string               `yaml:"name"`
	PIDs map[string]PIDConfig `yaml:"pids"`
}

type PIDConfig struct {
	Type    string        `yaml:"type"`
	PID     string        `yaml:"pid"`
	Unit    string        `yaml:"unit"`
	Refresh int           `yaml:"refresh"`
	Min     float64       `yaml:"min"`
	Max     float64       `yaml:"max"`
	Log     bool          `yaml:"log"`
	Display DisplayConfig `yaml:"display"`
}

type DisplayConfig struct {
	Enabled  bool           `yaml:"enabled"`
	Style    string         `yaml:"style"`
	Position PositionConfig `yaml:"position"`
}

type PositionConfig struct {
	X      float32 `yaml:"x"`
	Y      float32 `yaml:"y"`
	Width  float32 `yaml:"width"`
	Height float32 `yaml:"height"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	applyDefaults(&cfg)
	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Log.Rotate == "" {
		cfg.Log.Rotate = DefaultLogRotate
	}
	if cfg.Log.Directory == "" {
		cfg.Log.Directory = DefaultLogDirectory
	}
	if cfg.OBDAddress == "" {
		cfg.OBDAddress = DefaultOBDAddress
	}
}

func validate(cfg Config) error {
	if cfg.Vehicle.Name == "" {
		return fmt.Errorf("vehicle.name must not be empty")
	}
	if len(cfg.Vehicle.PIDs) == 0 {
		return fmt.Errorf("vehicle.pids must not be empty")
	}
	if cfg.Log.Rotate != DefaultLogRotate {
		return fmt.Errorf("log.rotate must be %q", DefaultLogRotate)
	}
	if cfg.Log.Directory == "" {
		return fmt.Errorf("log.directory must not be empty")
	}

	for key, pid := range cfg.Vehicle.PIDs {
		if pid.Type != "obd" && pid.Type != "virtual" {
			return fmt.Errorf("vehicle.pids.%s.type must be obd or virtual", key)
		}
		if pid.Type == "obd" && pid.PID == "" {
			return fmt.Errorf("vehicle.pids.%s.pid must not be empty for obd PIDs", key)
		}
		if (pid.Log || pid.Display.Enabled) && pid.Refresh <= 0 {
			return fmt.Errorf("vehicle.pids.%s.refresh must be positive for active PIDs", key)
		}
		if pid.Max <= pid.Min {
			return fmt.Errorf("vehicle.pids.%s.max must be greater than min", key)
		}
		if pid.Display.Enabled {
			if pid.Display.Style == "" {
				return fmt.Errorf("vehicle.pids.%s.display.style must not be empty when display is enabled", key)
			}
			if !validDisplayStyle(pid.Display.Style) {
				return fmt.Errorf("vehicle.pids.%s.display.style must be gauge, bar, or graph", key)
			}
			if pid.Display.Position.Width <= 0 || pid.Display.Position.Height <= 0 {
				return fmt.Errorf("vehicle.pids.%s.display.position width and height must be positive", key)
			}
		}
	}

	return nil
}

func validDisplayStyle(style string) bool {
	switch style {
	case "gauge", "bar", "graph":
		return true
	default:
		return false
	}
}
