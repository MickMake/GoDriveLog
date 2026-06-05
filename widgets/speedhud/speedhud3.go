package speedhud

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/widgets/model"
)

// SpeedHUD3 is SpeedHUD1 plus a daily peak display ("trip" = daily).
type SpeedHUD3 struct {
	widget.BaseWidget
	config model.GaugeConfig
	value  float64

	smoothBuf   []float64
	smoothNext  int
	smoothCount int

	peakValue float64
	peakDay   string
}

func NewSpeedHUD3(cfg model.GaugeConfig) model.Widget {
	cfg = cfg.Normalize()
	w := cfg.SmoothingWindow
	if w <= 1 {
		w = 1
	}
	g := &SpeedHUD3{config: cfg, value: cfg.Min, smoothBuf: make([]float64, w)}
	g.smoothBuf[0] = g.value
	g.smoothCount = 1
	g.peakValue = g.value
	g.peakDay = time.Now().Format("2006-01-02")
	g.ExtendBaseWidget(g)
	return g
}

func (g *SpeedHUD3) Style() string { return "speedhud3" }

func (g *SpeedHUD3) Config() model.GaugeConfig { return g.config }

func (g *SpeedHUD3) Value() float64 { return g.value }

func (g *SpeedHUD3) SetValue(v float64) {
	v = clamp(v, g.config.Min, g.config.Max)
	g.value = g.smooth(v)
	g.updatePeakDaily(g.value)
	g.Refresh()
}

func (g *SpeedHUD3) smooth(v float64) float64 {
	if len(g.smoothBuf) <= 1 {
		return v
	}
	g.smoothBuf[g.smoothNext] = v
	g.smoothNext = (g.smoothNext + 1) % len(g.smoothBuf)
	if g.smoothCount < len(g.smoothBuf) {
		g.smoothCount++
	}
	var sum float64
	for i := 0; i < g.smoothCount; i++ {
		sum += g.smoothBuf[i]
	}
	return sum / float64(g.smoothCount)
}

func (g *SpeedHUD3) updatePeakDaily(v float64) {
	today := time.Now().Format("2006-01-02")
	if today != g.peakDay {
		g.peakDay = today
		g.peakValue = v
		return
	}
	if v > g.peakValue {
		g.peakValue = v
	}
}

func (g *SpeedHUD3) Snapshot() model.Snapshot {
	v := g.Value()
	return model.Snapshot{
		Style:      g.Style(),
		Label:      g.config.Label,
		Unit:       g.config.Unit,
		Min:        g.config.Min,
		Max:        g.config.Max,
		Value:      v,
		Normalised: normalise(v, g.config.Min, g.config.Max),
		Warning:    inRange(v, g.config.WarningRange),
		Danger:     inRange(v, g.config.DangerRange),
	}
}

func (g *SpeedHUD3) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(parseHex(g.config.Theme.Background, color.NRGBA{R: 5, G: 7, B: 10, A: 255}))

	arcBg := make([]*canvas.Line, 0, hudSegments)
	arcVal := make([]*canvas.Line, 0, hudSegments)
	for i := 0; i < hudSegments; i++ {
		b := canvas.NewLine(withAlpha(parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}), 45))
		b.StrokeWidth = 8
		arcBg = append(arcBg, b)

		v := canvas.NewLine(parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
		v.StrokeWidth = 8
		v.Hide()
		arcVal = append(arcVal, v)
	}

	label := canvas.NewText(g.config.Label, parseHex(g.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Bold: true}

	value := canvas.NewText("", parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
	value.Alignment = fyne.TextAlignCenter
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	unit := canvas.NewText(g.config.Unit, parseHex(g.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	unit.Alignment = fyne.TextAlignCenter
	unit.TextStyle = fyne.TextStyle{Bold: true}

	peak := canvas.NewText("", withAlpha(parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}), 180))
	peak.Alignment = fyne.TextAlignCenter
	peak.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	r := &speedHUD3Renderer{g: g, bg: bg, arcBg: arcBg, arcVal: arcVal, label: label, value: value, unit: unit, peak: peak}
	r.Refresh()
	return r
}

type speedHUD3Renderer struct {
	g *SpeedHUD3

	bg     *canvas.Rectangle
	arcBg  []*canvas.Line
	arcVal []*canvas.Line

	label *canvas.Text
	value *canvas.Text
	unit  *canvas.Text
	peak  *canvas.Text
}

func (r *speedHUD3Renderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	cfg := r.g.config.Normalize()
	span := cfg.Max - cfg.Min
	if span <= 0 {
		span = 1
	}
	pct := clamp((r.g.value-cfg.Min)/span, 0, 1)

	cx := size.Width / 2
	cy := size.Height * 0.78
	radius := math.Min(float64(size.Width), float64(size.Height)) * 0.40
	start := 1.10 * math.Pi
	end := 1.90 * math.Pi
	arc := end - start

	for i := 0; i < hudSegments; i++ {
		p0 := float64(i) / float64(hudSegments)
		p1 := float64(i+1) / float64(hudSegments)
		a0 := start + p0*arc
		a1 := start + p1*arc
		pos1 := fyne.NewPos(cx+float32(radius*math.Cos(a0)), cy+float32(radius*math.Sin(a0)))
		pos2 := fyne.NewPos(cx+float32(radius*math.Cos(a1)), cy+float32(radius*math.Sin(a1)))

		b := r.arcBg[i]
		b.Position1, b.Position2 = pos1, pos2

		v := r.arcVal[i]
		if p0 >= pct {
			v.Hide()
			continue
		}
		if p1 > pct {
			p1 = pct
			a1 = start + p1*arc
			pos2 = fyne.NewPos(cx+float32(radius*math.Cos(a1)), cy+float32(radius*math.Sin(a1)))
		}
		midValue := cfg.Min + ((p0+p1)/2)*span
		v.Position1, v.Position2 = pos1, pos2
		if inRange(midValue, cfg.DangerRange) {
			v.StrokeColor = parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
		} else if inRange(midValue, cfg.WarningRange) {
			v.StrokeColor = parseHex(cfg.Theme.Warning, color.NRGBA{R: 255, G: 176, B: 0, A: 255})
		} else {
			v.StrokeColor = parseHex(cfg.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255})
		}
		v.Show()
	}

	// Text stack
	r.label.Text = cfg.Label
	r.label.TextSize = float32(math.Max(12, float64(size.Height)*0.06))
	r.label.Refresh()
	r.label.Resize(fyne.NewSize(size.Width, r.label.MinSize().Height))
	r.label.Move(fyne.NewPos(0, size.Height*0.06))

	r.value.Text = formatValue(r.g.value, span)
	r.value.TextSize = float32(math.Max(36, float64(size.Height)*0.32))
	r.value.Refresh()
	r.value.Resize(fyne.NewSize(size.Width, r.value.MinSize().Height))
	r.value.Move(fyne.NewPos(0, size.Height*0.18))

	r.unit.Text = cfg.Unit
	r.unit.TextSize = float32(math.Max(12, float64(size.Height)*0.07))
	r.unit.Refresh()
	r.unit.Resize(fyne.NewSize(size.Width, r.unit.MinSize().Height))
	r.unit.Move(fyne.NewPos(0, size.Height*0.18+r.value.MinSize().Height-6))

	r.peak.Text = fmt.Sprintf("PEAK %s", formatValue(r.g.peakValue, span))
	r.peak.TextSize = float32(math.Max(10, float64(size.Height)*0.055))
	r.peak.Refresh()
	r.peak.Resize(fyne.NewSize(size.Width, r.peak.MinSize().Height))
	r.peak.Move(fyne.NewPos(0, size.Height*0.78-r.peak.MinSize().Height))
}

func (r *speedHUD3Renderer) MinSize() fyne.Size { return fyne.NewSize(360, 240) }

func (r *speedHUD3Renderer) Refresh() {
	r.Layout(r.g.Size())
	canvas.Refresh(r.bg)
	for _, l := range r.arcBg {
		canvas.Refresh(l)
	}
	for _, l := range r.arcVal {
		canvas.Refresh(l)
	}
	canvas.Refresh(r.label)
	canvas.Refresh(r.value)
	canvas.Refresh(r.unit)
	canvas.Refresh(r.peak)
}

func (r *speedHUD3Renderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.bg}
	for _, l := range r.arcBg {
		objs = append(objs, l)
	}
	for _, l := range r.arcVal {
		objs = append(objs, l)
	}
	objs = append(objs, r.label, r.value, r.unit, r.peak)
	return objs
}

func (r *speedHUD3Renderer) Destroy() {}

func formatValue(value, span float64) string {
	if span >= 50 {
		return fmt.Sprintf("%.0f", value)
	}
	return fmt.Sprintf("%.1f", value)
}

func parseHex(s string, fallback color.NRGBA) color.NRGBA {
	s = strings.TrimSpace(strings.TrimPrefix(s, "#"))
	if len(s) != 6 {
		return fallback
	}
	r, err1 := strconv.ParseUint(s[0:2], 16, 8)
	g, err2 := strconv.ParseUint(s[2:4], 16, 8)
	b, err3 := strconv.ParseUint(s[4:6], 16, 8)
	if err1 != nil || err2 != nil || err3 != nil {
		return fallback
	}
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

func withAlpha(c color.NRGBA, a uint8) color.NRGBA { c.A = a; return c }

func clamp(v, min, max float64) float64 { return math.Max(min, math.Min(max, v)) }

func normalise(value, min, max float64) float64 {
	if max == min {
		return 0
	}
	return clamp((value-min)/(max-min), 0, 1)
}

func inRange(value float64, r *model.Range) bool {
	if r == nil {
		return false
	}
	min, max := r.Min, r.Max
	if max < min {
		min, max = max, min
	}
	return value >= min && value <= max
}
