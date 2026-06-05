package widgets

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Range describes an inclusive numeric range for warnings, danger zones, or scale limits.
type Range struct {
	Min float64
	Max float64
}

// Theme groups dashboard colours by purpose. Colour strings are intentionally plain
// so this package can stay independent of Fyne internals until renderers need them.
type Theme struct {
	Background string
	Foreground string
	Tick       string
	Label      string
	Needle     string
	Value      string
	Warning    string
	Danger     string
	GraphLine  string
	Grid       string
}

// DefaultTheme returns dark automotive dashboard colours suitable for first-pass widgets.
func DefaultTheme() Theme {
	return Theme{
		Background: "#05070a",
		Foreground: "#d7dde8",
		Tick:       "#7f8a99",
		Label:      "#c7d0dd",
		Needle:     "#f0f4ff",
		Value:      "#7df9ff",
		Warning:    "#ffb000",
		Danger:     "#ff3030",
		GraphLine:  "#66d9ef",
		Grid:       "#263241",
	}
}

// GaugeConfig contains the shared options for numeric dashboard widgets.
type GaugeConfig struct {
	Label string
	Unit  string
	Min   float64
	Max   float64

	ShowLabel       bool
	ShowUnit        bool
	ShowMin         bool
	ShowMax         bool
	ShowValue       bool
	ShowTicks       bool
	ShowMajorLabels bool

	WarningRange *Range
	DangerRange  *Range
	Theme        Theme
}

// DefaultGaugeConfig returns percentage-mode defaults.
func DefaultGaugeConfig() GaugeConfig {
	return GaugeConfig{
		Label:           "Value",
		Unit:            "%",
		Min:             0,
		Max:             100,
		ShowLabel:       true,
		ShowUnit:        true,
		ShowMin:         true,
		ShowMax:         true,
		ShowValue:       true,
		ShowTicks:       true,
		ShowMajorLabels: true,
		Theme:           DefaultTheme(),
	}
}

// Normalize fills in safe defaults without hiding an explicitly supplied zero min or max.
func (c GaugeConfig) Normalize() GaugeConfig {
	if strings.TrimSpace(c.Label) == "" {
		c.Label = "Value"
	}
	if c.Min == c.Max {
		c.Min = 0
		c.Max = 100
	}
	if c.Max < c.Min {
		c.Min, c.Max = c.Max, c.Min
	}
	if c.Theme == (Theme{}) {
		c.Theme = DefaultTheme()
	}
	return c
}

// Widget is the small runtime contract GoDriveLog can call without knowing widget internals.
type Widget interface {
	Style() string
	Config() GaugeConfig
	Value() float64
	SetValue(float64)
	Snapshot() Snapshot
}

// Snapshot is a renderer-neutral view of a widget's current state.
type Snapshot struct {
	Style      string
	Label      string
	Unit       string
	Min        float64
	Max        float64
	Value      float64
	Normalised float64
	Warning    bool
	Danger     bool
}

// NumericWidget is a first-pass reusable numeric widget stub. Fyne renderers can wrap
// this state object later without creating GoDriveLog dependencies.
type NumericWidget struct {
	style  string
	config GaugeConfig
	value  float64
}

func newNumericWidget(style string, cfg GaugeConfig) *NumericWidget {
	cfg = cfg.Normalize()
	return &NumericWidget{style: style, config: cfg, value: cfg.Min}
}

func (w *NumericWidget) Style() string { return w.style }

func (w *NumericWidget) Config() GaugeConfig { return w.config }

func (w *NumericWidget) Value() float64 { return w.value }

func (w *NumericWidget) SetValue(value float64) { w.value = clamp(value, w.config.Min, w.config.Max) }

func (w *NumericWidget) Snapshot() Snapshot {
	value := w.Value()
	return Snapshot{
		Style:      w.style,
		Label:      w.config.Label,
		Unit:       w.config.Unit,
		Min:        w.config.Min,
		Max:        w.config.Max,
		Value:      value,
		Normalised: normalise(value, w.config.Min, w.config.Max),
		Warning:    inRange(value, w.config.WarningRange),
		Danger:     inRange(value, w.config.DangerRange),
	}
}

// New returns a widget by config style name, such as radial1, bar1, graph1, or led1.
func New(style string, cfg GaugeConfig) (Widget, error) {
	switch strings.ToLower(strings.TrimSpace(style)) {
	case "radial1":
		return NewRadial1(cfg), nil
	case "bar1":
		return NewBar1(cfg), nil
	case "graph1":
		return NewGraph1(cfg), nil
	case "led1":
		return NewLED1(cfg), nil
	default:
		return nil, fmt.Errorf("unknown widget style %q", style)
	}
}

func NewRadial1(cfg GaugeConfig) Widget { return newNumericWidget("radial1", cfg) }

func NewBar1(cfg GaugeConfig) Widget { return newNumericWidget("bar1", cfg) }

func NewGraph1(cfg GaugeConfig) Widget { return newNumericWidget("graph1", cfg) }

func NewLED1(cfg GaugeConfig) Widget { return newNumericWidget("led1", cfg) }

func Styles() []string {
	styles := []string{"bar1", "graph1", "led1", "radial1"}
	sort.Strings(styles)
	return styles
}

func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func normalise(value, min, max float64) float64 {
	if max == min {
		return 0
	}
	return clamp((value-min)/(max-min), 0, 1)
}

func inRange(value float64, r *Range) bool {
	if r == nil {
		return false
	}
	min, max := r.Min, r.Max
	if max < min {
		min, max = max, min
	}
	return value >= min && value <= max
}
