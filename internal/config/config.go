package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LogDir               string         `json:"log_dir"`
	EngineStartPID       string         `json:"engine_start_pid"`
	EngineStartThreshold float64        `json:"engine_start_threshold"`
	MockMode             bool           `json:"mock_mode"`
	Sensors              []SensorConfig `json:"sensors"`
}

type SensorConfig struct {
	PID       string        `json:"pid"`
	Name      string        `json:"name"`
	RefreshMS int           `json:"refresh_ms"`
	Style     string        `json:"style"`
	Min       float64       `json:"min"`
	Max       float64       `json:"max"`
	Display   DisplayConfig `json:"display"`
}

type DisplayConfig struct {
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.LogDir == "" {
		cfg.LogDir = "./logs"
	}
	if cfg.EngineStartPID == "" {
		cfg.EngineStartPID = "010C"
	}
	if cfg.EngineStartThreshold <= 0 {
		cfg.EngineStartThreshold = 50
	}
	if len(cfg.Sensors) == 0 {
		return Config{}, fmt.Errorf("config has no sensors")
	}
	for i := range cfg.Sensors {
		if cfg.Sensors[i].PID == "" || cfg.Sensors[i].Name == "" {
			return Config{}, fmt.Errorf("sensor %d requires pid and name", i)
		}
		if cfg.Sensors[i].RefreshMS <= 0 {
			cfg.Sensors[i].RefreshMS = 1000
		}
		if cfg.Sensors[i].Max <= cfg.Sensors[i].Min {
			cfg.Sensors[i].Max = cfg.Sensors[i].Min + 100
		}
		if cfg.Sensors[i].Display.Width <= 0 {
			cfg.Sensors[i].Display.Width = 300
		}
		if cfg.Sensors[i].Display.Height <= 0 {
			cfg.Sensors[i].Display.Height = 90
		}
	}

	return cfg, nil
}
