package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultLogRotate    = "daily"
	DefaultLogDirectory = "./log"
	DefaultOBDAddress   = "serial:///dev/ttyUSB0"
)

type Config struct {
	MockMode   bool                       `yaml:"mock_mode"`
	OBDAddress string                     `yaml:"obd_address"`
	OBDDebug   bool                       `yaml:"obd_debug"`
	Log        LogConfig                  `yaml:"log"`
	Vehicle    VehicleConfig              `yaml:"vehicle"`
	Sensors    map[string]SensorConfig    `yaml:"sensors"`
	Dashboard  DashboardConfig            `yaml:"dashboard"`
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
	if len(cfg.Sensors) == 0 {
		return fmt.Errorf("sensors must not be empty")
	}
	if cfg.Log.Rotate != DefaultLogRotate {
		return fmt.Errorf("log.rotate must be %q", DefaultLogRotate)
	}
	if cfg.Log.Directory == "" {
		return fmt.Errorf("log.directory must not be empty")
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
