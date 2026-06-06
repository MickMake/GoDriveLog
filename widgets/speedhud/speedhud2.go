package speedhud

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/widgets/model"
)

// SpeedHUD2 is SpeedHUD1 plus a simple glow effect for the numeric value.
// (Fake glow: extra text layers behind the main value.)
type SpeedHUD2 struct {
	widget.BaseWidget
	config model.GaugeConfig
	value  float64

	smoothBuf   []float64
	smoothNext  int
	smoothCount int

	pulse model.PulseTracker
}

func NewSpeedHUD2(cfg model.GaugeConfig) model.Widget {
	cfg = cfg.Normalize()
	w := cfg.SmoothingWindow
	if w <= 1 {
		w = 1
	}
	g := &SpeedHUD2{config: cfg, value: cfg.Min, smoothBuf: make([]float64, w), pulse: model.NewPulseTracker()}
	g.smoothBuf[0] = g.value
	g.smoothCount = 1
	g.ExtendBaseWidget(g)
	return g
}

func (g *SpeedHUD2) Style() string { return "speedhud2" }

func (g *SpeedHUD2) Config() model.GaugeConfig { return g.config }

func (g *SpeedHUD2) Value() float64 { return g.value }

func (g *SpeedHUD2) SetValue(v float64) {
	v = clamp(v, g.config.Min, g.config.Max)
	g.value = g.smooth(v)
	g.pulse.Update(g.value, g.config.WarningRange, g.config.DangerRange)
	g.Refresh()
}

func (g *SpeedHUD2) smooth(v float64) float64 {
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

func (g *SpeedHUD2) Snapshot() model.Snapshot {
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

func (g *SpeedHUD2) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(parseHex(g.config.Theme.Background, color.NRGBA{R: 5, G: 7, B: 10, A: 255}))
	pulse := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pulse.Hide()

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

	// Glow layers (behind)
	glow1 := canvas.NewText("", withAlpha(parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}), 80))
	glow1.Alignment = fyne.TextAlignCenter
	glow1.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	glow2 := canvas.NewText("", withAlpha(parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}), 40))
	glow2.Alignment = fyne.TextAlignCenter
	glow2.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	value := canvas.NewText("", parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
	value.Alignment = fyne.TextAlignCenter
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	unit := canvas.NewText(g.config.Unit, parseHex(g.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	unit.Alignment = fyne.TextAlignCenter
	unit.TextStyle = fyne.TextStyle{Bold: true}

	r := &speedHUD2Renderer{g: g, bg: bg, pulse: pulse, arcBg: arcBg, arcVal: arcVal, label: label, glow1: glow1, glow2: glow2, value: value, unit: unit}
	r.Refresh()
	return r
}

type speedHUD2Renderer struct {
	g *SpeedHUD2

	bg    *canvas.Rectangle
	pulse *canvas.Rectangle
	arcBg  []*canvas.Line
	arcVal []*canvas.Line

	label *canvas.Text
	glow1 *canvas.Text
	glow2 *canvas.Text
	value *canvas.Text
	unit  *canvas.Text
}

func (r *speedHUD2Renderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.pulse.Resize(size)

	cfg := r.g.config.Normalize()
	span := cfg.Max - cfg.Min
	if span <= 0 {
		span = 1
	}
	pct := clamp((r.g.value-cfg.Min)/span, 0, 1)

	pState, p := r.g.pulse.Pulse(time.Now())
	if p > 0 {
		col := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		switch pState {
		case model.AlertWarning:
			col = parseHex(cfg.Theme.Warning, color.NRGBA{R: 255, G: 176, B: 0, A: 255})
		case model.AlertDanger:
			col = parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
		}
		r.pulse.FillColor = withAlpha(col, uint8(70*p))
		r.pulse.Show()
	} else {
		r.pulse.Hide()
	}

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

	valStr := formatValue(r.g.value, span)
	fontSize := float32(math.Max(36, float64(size.Height)*0.32))

	r.glow2.Text = valStr
	r.glow2.TextSize = fontSize + 8
	r.glow2.Refresh()
	r.glow2.Resize(fyne.NewSize(size.Width, r.glow2.MinSize().Height))
	r.glow2.Move(fyne.NewPos(0, size.Height*0.20-4))

	r.glow1.Text = valStr
	r.glow1.TextSize = fontSize + 4
	r.glow1.Refresh()
	r.glow1.Resize(fyne.NewSize(size.Width, r.glow1.MinSize().Height))
	r.glow1.Move(fyne.NewPos(0, size.Height*0.20-2))

	r.value.Text = valStr
	r.value.TextSize = fontSize
	r.value.Refresh()
	r.value.Resize(fyne.NewSize(size.Width, r.value.MinSize().Height))
	r.value.Move(fyne.NewPos(0, size.Height*0.20))

	r.unit.Text = cfg.Unit
	r.unit.TextSize = float32(math.Max(12, float64(size.Height)*0.07))
	r.unit.Refresh()
	r.unit.Resize(fyne.NewSize(size.Width, r.unit.MinSize().Height))
	r.unit.Move(fyne.NewPos(0, size.Height*0.20+r.value.MinSize().Height-6))
}

func (r *speedHUD2Renderer) MinSize() fyne.Size { return fyne.NewSize(360, 240) }

func (r *speedHUD2Renderer) Refresh() {
	r.Layout(r.g.Size())
	canvas.Refresh(r.bg)
	canvas.Refresh(r.pulse)
	for _, l := range r.arcBg {
		canvas.Refresh(l)
	}
	for _, l := range r.arcVal {
		canvas.Refresh(l)
	}
	canvas.Refresh(r.label)
	canvas.Refresh(r.glow2)
	canvas.Refresh(r.glow1)
	canvas.Refresh(r.value)
	canvas.Refresh(r.unit)
}

func (r *speedHUD2Renderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.bg, r.pulse}
	for _, l := range r.arcBg {
		objs = append(objs, l)
	}
	for _, l := range r.arcVal {
		objs = append(objs, l)
	}
	objs = append(objs, r.label, r.glow2, r.glow1, r.value, r.unit)
	return objs
}

func (r *speedHUD2Renderer) Destroy() {}
