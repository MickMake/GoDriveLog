package model

import (
	"math"
	"strings"
	"time"
)

// Range describes an inclusive numeric range for warnings, danger zones, or scale limits.
type Range struct {
	Min float64
	Max float64
}

// Theme groups dashboard colours by purpose. Colour strings are intentionally plain
// so this package can stay renderer-agnostic (Fyne widgets translate them as needed).
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
	ShowPeak        bool

	// SmoothingWindow applies a simple moving average over N samples.
	// 1 means no smoothing.
	SmoothingWindow int

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
		ShowPeak:        false,
		SmoothingWindow: 1,
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
	if c.SmoothingWindow <= 0 {
		c.SmoothingWindow = 1
	}
	if c.Theme == (Theme{}) {
		c.Theme = DefaultTheme()
	}
	return c
}

// AlertState is the current threshold state derived from WarningRange/DangerRange.
type AlertState int

const (
	AlertNormal AlertState = iota
	AlertWarning
	AlertDanger
)

// PulseTracker implements a small state machine for edge-triggered warning/danger pulses.
//
// - It only triggers when entering Warning or Danger.
// - It re-triggers only after leaving that state and re-entering ("re-arm on exit").
// - If Warning/Danger ranges are not set, it stays idle.
type PulseTracker struct {
	state       AlertState
	pulseState  AlertState
	pulseStart  time.Time
	pulseLength time.Duration
}

func NewPulseTracker() PulseTracker {
	return PulseTracker{state: AlertNormal, pulseState: AlertNormal, pulseLength: 1500 * time.Millisecond}
}

// Update recalculates the current state and arms/triggers pulses on state transitions.
func (p *PulseTracker) Update(value float64, warn, danger *Range) AlertState {
	newState := AlertNormal
	if inRange(value, danger) {
		newState = AlertDanger
	} else if inRange(value, warn) {
		newState = AlertWarning
	}

	if newState != p.state {
		// Trigger on entry to warning/danger.
		switch newState {
		case AlertWarning:
			p.pulseState = AlertWarning
			p.pulseStart = time.Now()
		case AlertDanger:
			p.pulseState = AlertDanger
			p.pulseStart = time.Now()
		}
		p.state = newState
	}

	return p.state
}

// Pulse returns the currently active pulse type and intensity (0..1), if any.
func (p *PulseTracker) Pulse(now time.Time) (AlertState, float64) {
	if p.pulseState == AlertNormal || p.pulseStart.IsZero() {
		return AlertNormal, 0
	}
	age := now.Sub(p.pulseStart)
	if age >= p.pulseLength {
		p.pulseState = AlertNormal
		return AlertNormal, 0
	}
	intensity := 1.0 - float64(age)/float64(p.pulseLength)
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}
	return p.pulseState, intensity
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

// NumericWidget is a reusable state-only numeric widget stub.
// Renderers can wrap this state object later without creating GoDriveLog dependencies.
type NumericWidget struct {
	style  string
	config GaugeConfig
	value  float64
}

func NewNumericWidget(style string, cfg GaugeConfig) *NumericWidget {
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
