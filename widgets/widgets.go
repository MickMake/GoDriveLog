package widgets

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MickMake/GoDriveLog/widgets/bar"
	"github.com/MickMake/GoDriveLog/widgets/model"
	"github.com/MickMake/GoDriveLog/widgets/radial"
	"github.com/MickMake/GoDriveLog/widgets/ramped"
	"github.com/MickMake/GoDriveLog/widgets/speedhud"
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
	case "half_top1":
		return radial.NewHalfTop1(cfg), nil
	case "half_bottom1":
		return radial.NewHalfBottom1(cfg), nil
	case "quarter_tl1":
		return radial.NewQuarterTL1(cfg), nil
	case "quarter_tr1":
		return radial.NewQuarterTR1(cfg), nil
	case "quarter_bl1":
		return radial.NewQuarterBL1(cfg), nil
	case "quarter_br1":
		return radial.NewQuarterBR1(cfg), nil
	case "ramped1":
		return ramped.NewRamped1(cfg), nil
	case "ramped2":
		return ramped.NewRamped2(cfg), nil
	case "ramped3":
		return ramped.NewRamped3(cfg), nil
	case "speedhud1":
		return speedhud.NewSpeedHUD1(cfg), nil
	case "speedhud2":
		return speedhud.NewSpeedHUD2(cfg), nil
	case "speedhud3":
		return speedhud.NewSpeedHUD3(cfg), nil
	case "bar1":
		return model.NewNumericWidget("bar1", cfg), nil
	case "bar2":
		return bar.NewBar2(cfg), nil
	case "bar3":
		return bar.NewBar3(cfg), nil
	case "graph1":
		return model.NewNumericWidget("graph1", cfg), nil
	case "led1":
		return model.NewNumericWidget("led1", cfg), nil
	default:
		return nil, fmt.Errorf("unknown widget style %q", style)
	}
}

func Styles() []string {
	styles := []string{"bar1", "bar2", "bar3", "graph1", "half_bottom1", "half_top1", "led1", "quarter_bl1", "quarter_br1", "quarter_tl1", "quarter_tr1", "radial1", "radial2", "radial3", "ramped1", "ramped2", "ramped3", "speedhud1", "speedhud2", "speedhud3"}
	sort.Strings(styles)
	return styles
}
