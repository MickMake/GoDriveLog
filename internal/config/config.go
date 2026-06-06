package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultLogRotate    = "daily"
	DefaultLogDirectory = "./log"
	DefaultOBDAddress   = "serial:///dev/ttyUSB0"
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
	Enabled  bool               `yaml:"enabled"`
	Widget   string             `yaml:"widget"`
	Style    DisplayStyleConfig `yaml:"style"`
	Position PositionConfig     `yaml:"position"`
}

type DisplayStyleConfig struct {
	SmoothingWindow int    `yaml:"smoothing_window"`
	DialRotation    int    `yaml:"dial_rotation"`
	ViewRotation    int    `yaml:"view_rotation"`
	ScaleDirection  string `yaml:"scale_direction"`
}

type PositionConfig struct {
	X      float32 `yaml:"x"`
	Y      float32 `yaml:"y"`
	Width  float32 `yaml:"width"`
	Height float32 `yaml:"height"`
	Z      float32 `yaml:"z"`
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
		if pid.Unit == "" {
			return fmt.Errorf("vehicle.pids.%s.unit must not be empty", key)
		}
		if pid.Refresh <= 0 {
			return fmt.Errorf("vehicle.pids.%s.refresh must be positive", key)
		}
		if pid.Max <= pid.Min {
			return fmt.Errorf("vehicle.pids.%s.max must be greater than min", key)
		}
		if pid.Display.Enabled {
			if pid.Display.Widget == "" {
				return fmt.Errorf("vehicle.pids.%s.display.widget must not be empty when display is enabled", key)
			}
			if !validDisplayWidget(pid.Display.Widget) {
				return fmt.Errorf("vehicle.pids.%s.display.widget must be a supported widget style", key)
			}
			if pid.Display.Style.SmoothingWindow < 0 {
				return fmt.Errorf("vehicle.pids.%s.display.style.smoothing_window must be >= 0", key)
			}
			if pid.Display.Position.Width <= 0 || pid.Display.Position.Height <= 0 {
				return fmt.Errorf("vehicle.pids.%s.display.position width and height must be positive", key)
			}
			if !validRotation(pid.Display.Style.DialRotation) {
				return fmt.Errorf("vehicle.pids.%s.display.style.dial_rotation must be one of 0, 90, 180, 270", key)
			}
			if !validRotation(pid.Display.Style.ViewRotation) {
				return fmt.Errorf("vehicle.pids.%s.display.style.view_rotation must be one of 0, 90, 180, 270", key)
			}
			if !validScaleDirection(pid.Display.Style.ScaleDirection) {
				return fmt.Errorf("vehicle.pids.%s.display.style.scale_direction must be forward or reverse", key)
			}
		}
	}

	return nil
}

func validRotation(deg int) bool {
	switch deg {
	case 0, 90, 180, 270:
		return true
	default:
		return false
	}
}

func validScaleDirection(dir string) bool {
	dir = strings.ToLower(strings.TrimSpace(dir))
	if dir == "" {
		return true // default handled elsewhere
	}
	switch dir {
	case "forward", "reverse":
		return true
	default:
		return false
	}
}

func validDisplayWidget(widget string) bool {
	switch strings.ToLower(strings.TrimSpace(widget)) {
	case "radial1", "radial2", "radial3",
		"half_top1", "half_bottom1", "quarter_tl1", "quarter_tr1", "quarter_bl1", "quarter_br1",
		"sweep1", "sweep2", "sweep3",
		"speedhud1", "speedhud2", "speedhud3",
		"bar1", "bar2", "bar3", "graph1", "led1":
		return true
	default:
		return false
	}
}
