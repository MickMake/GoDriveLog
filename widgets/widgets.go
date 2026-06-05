package widgets

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MickMake/GoDriveLog/widgets/model"
	"github.com/MickMake/GoDriveLog/widgets/radial"
)

// Re-export shared model types from widgets/model so subpackages can depend on model
// without creating import cycles.
type Range = model.Range
type Theme = model.Theme
type GaugeConfig = model.GaugeConfig
type Widget = model.Widget
type Snapshot = model.Snapshot

func DefaultTheme() Theme { return model.DefaultTheme() }
func DefaultGaugeConfig() GaugeConfig { return model.DefaultGaugeConfig() }

// New returns a widget by config style name, such as radial1, bar1, graph1, or led1.
func New(style string, cfg GaugeConfig) (Widget, error) {
	switch strings.ToLower(strings.TrimSpace(style)) {
	case "radial1":
		return radial.NewRadial1(cfg), nil
	case "radial2":
		return radial.NewRadial2(cfg), nil
	case "radial3":
		return radial.NewRadial3(cfg), nil
	case "bar1":
		return model.NewNumericWidget("bar1", cfg), nil
	case "graph1":
		return model.NewNumericWidget("graph1", cfg), nil
	case "led1":
		return model.NewNumericWidget("led1", cfg), nil
	default:
		return nil, fmt.Errorf("unknown widget style %q", style)
	}
}

func Styles() []string {
	styles := []string{"bar1", "graph1", "led1", "radial1", "radial2", "radial3"}
	sort.Strings(styles)
	return styles
}
