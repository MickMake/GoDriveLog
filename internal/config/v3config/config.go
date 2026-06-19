package v3config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	WidgetTypeImage        = "image"
	WidgetTypeDigitDisplay = "digit_display"
	WidgetTypeBarDisplay   = "bar_display"
	WidgetTypeFrameGauge   = "frame_gauge"
	WidgetTypeIndicator    = "indicator"
	WidgetTypeGauge        = "gauge"

	ValueKindNumeric = "numeric"
	ValueKindBool    = "bool"
	ValueKindString  = "string"
)

type Config struct {
	Vehicles   map[string]VehicleConfig   `yaml:"vehicles"`
	Sensors    map[string]SensorConfig    `yaml:"sensors"`
	Assets     AssetConfig                `yaml:"assets"`
	Logs       map[string]LogConfig       `yaml:"logs"`
	Dashboards map[string]DashboardConfig `yaml:"dashboards"`
}

type VehicleConfig struct {
	Name       string    `yaml:"name"`
	OBD        OBDConfig `yaml:"obd"`
	Logs       []string  `yaml:"logs,omitempty"`
	Dashboards []string  `yaml:"dashboards,omitempty"`
}

type OBDConfig struct {
	Address string `yaml:"address"`
	Timeout int    `yaml:"timeout"`
}

type SensorConfig struct {
	Type      string   `yaml:"type"`
	PID       string   `yaml:"pid,omitempty"`
	ValueKind string   `yaml:"value_kind,omitempty"`
	Unit      string   `yaml:"unit"`
	Poll      int      `yaml:"poll"`
	Min       *float64 `yaml:"min,omitempty"`
	Max       *float64 `yaml:"max,omitempty"`
}

type AssetConfig struct {
	DigitSets     map[string]DigitSetConfig     `yaml:"digit_sets"`
	BarSets       map[string]BarSetConfig       `yaml:"bar_sets"`
	FrameSets     map[string]FrameSetConfig     `yaml:"frame_sets"`
	IndicatorSets map[string]IndicatorSetConfig `yaml:"indicator_sets"`
	ImageSets     map[string]ImageSetConfig     `yaml:"image_sets"`
}

type DigitSetConfig struct {
	Background   string            `yaml:"background,omitempty"`
	Characters   map[string]string `yaml:"characters"`
	DecimalPoint string            `yaml:"decimal_point,omitempty"`
	Foreground   string            `yaml:"foreground,omitempty"`
	Spacing      int               `yaml:"spacing,omitempty"`
}

type BarSetConfig struct {
	Background string            `yaml:"background,omitempty"`
	Cells      map[string]string `yaml:"cells"`
	Foreground string            `yaml:"foreground,omitempty"`
	Spacing    int               `yaml:"spacing,omitempty"`
}

type FrameSetConfig struct {
	Background string           `yaml:"background,omitempty"`
	Frames     FrameRangeConfig `yaml:"frames"`
	Foreground string           `yaml:"foreground,omitempty"`
}

type FrameRangeConfig struct {
	Path  string `yaml:"path"`
	First int    `yaml:"first"`
	Last  int    `yaml:"last"`
}

type IndicatorSetConfig struct {
	Background string            `yaml:"background,omitempty"`
	States     map[string]string `yaml:"states"`
	Foreground string            `yaml:"foreground,omitempty"`
}

type ImageSetConfig struct {
	Image      string `yaml:"image,omitempty"`
	Background string `yaml:"background,omitempty"`
	Foreground string `yaml:"foreground,omitempty"`
}

type LogConfig struct {
	Path    string   `yaml:"path"`
	Sensors []string `yaml:"sensors"`
}

type DashboardConfig struct {
	Display string         `yaml:"display"`
	Size    SizeConfig     `yaml:"size"`
	Widgets []WidgetConfig `yaml:"widgets"`
}

type SizeConfig struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type WidgetConfig struct {
	ID       string       `yaml:"id"`
	Type     string       `yaml:"type"`
	Sensor   string       `yaml:"sensor,omitempty"`
	Asset    string       `yaml:"asset,omitempty"`
	Gauge    string       `yaml:"gauge,omitempty"`
	Position []int        `yaml:"position"`
	Scale    float64      `yaml:"scale,omitempty"`
	Digits   int          `yaml:"digits,omitempty"`
	Format   string       `yaml:"format,omitempty"`
	Cells    int          `yaml:"cells,omitempty"`
	Min      *float64     `yaml:"min,omitempty"`
	Max      *float64     `yaml:"max,omitempty"`
	Reverse  bool         `yaml:"reverse,omitempty"`
	Zones    []ZoneConfig `yaml:"zones,omitempty"`
}

type ZoneConfig struct {
	UpTo float64 `yaml:"up_to"`
	Cell string  `yaml:"cell"`
}

func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return LoadBytes(data)
}

func LoadBytes(data []byte) (Config, error) {
	var cfg Config
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	if err := Validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func SensorDeclaredValueKind(sensor SensorConfig) string {
	return sensor.ValueKind
}

func SensorOutputValueKind(sensor SensorConfig) string {
	switch sensor.Type {
	case "obd":
		return ValueKindNumeric
	default:
		return ""
	}
}

func SensorEffectiveValueKind(sensor SensorConfig) string {
	if sensor.ValueKind != "" {
		return sensor.ValueKind
	}
	return SensorOutputValueKind(sensor)
}
