package ramped

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/widgets/model"
)

// Ramped3 is sweep1 with a daily peak marker line along the arc.
type Ramped3 struct {
	widget.BaseWidget
	config model.GaugeConfig
	value  float64

	smoothBuf   []float64
	smoothNext  int
	smoothCount int

	peakValue float64
	peakDay   string

	pulse model.PulseTracker
}

func NewRamped3(cfg model.GaugeConfig) model.Widget {
	cfg = cfg.Normalize()
	w := cfg.SmoothingWindow
	if w <= 1 {
		w = 1
	}
	g := &Ramped3{config: cfg, value: cfg.Min, smoothBuf: make([]float64, w), pulse: model.NewPulseTracker()}
	g.smoothBuf[0] = g.value
	g.smoothCount = 1
	g.peakValue = g.value
	g.peakDay = time.Now().Format("2006-01-02")
	g.ExtendBaseWidget(g)
	return g
}

func (g *Ramped3) Style() string { return "sweep3" }

func (g *Ramped3) Config() model.GaugeConfig { return g.config }

func (g *Ramped3) Value() float64 { return g.value }

func (g *Ramped3) SetValue(v float64) {
	cfg := g.config.Normalize()
	v = clamp(v, cfg.Min, cfg.Max)
	g.value = g.smooth(v)
	g.updatePeakDaily(g.value)
	g.pulse.Update(g.value, g.config.WarningRange, g.config.DangerRange)
	g.Refresh()
}

func (g *Ramped3) smooth(v float64) float64 {
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

func (g *Ramped3) updatePeakDaily(v float64) {
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

func (g *Ramped3) Snapshot() model.Snapshot {
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

func (g *Ramped3) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(parseHex(g.config.Theme.Background, color.NRGBA{R: 5, G: 7, B: 10, A: 255}))
	pulse := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pulse.Hide()

	arcBg := make([]*canvas.Line, 0, rampedSegments)
	arcVal := make([]*canvas.Line, 0, rampedSegments)
	for i := 0; i < rampedSegments; i++ {
		b := canvas.NewLine(withAlpha(parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}), 55))
		b.StrokeWidth = 10
		arcBg = append(arcBg, b)

		v := canvas.NewLine(parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
		v.StrokeWidth = 10
		v.Hide()
		arcVal = append(arcVal, v)
	}

	ticks := make([]*canvas.Line, 0, rampedTicks)
	for i := 0; i < rampedTicks; i++ {
		l := canvas.NewLine(parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}))
		l.StrokeWidth = 3
		ticks = append(ticks, l)
	}

	label := canvas.NewText(g.config.Label, parseHex(g.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	label.Alignment = fyne.TextAlignLeading
	label.TextStyle = fyne.TextStyle{Bold: true}

	value := canvas.NewText("", parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
	value.Alignment = fyne.TextAlignTrailing
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	unit := canvas.NewText(g.config.Unit, parseHex(g.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	unit.Alignment = fyne.TextAlignTrailing
	unit.TextStyle = fyne.TextStyle{Bold: true}

	needle := canvas.NewLine(parseHex(g.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255}))
	needle.StrokeWidth = 4

	peakLine := canvas.NewLine(parseHex(g.config.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255}))
	peakLine.StrokeWidth = 6
	peakLine.Hide()

	r := &ramped3Renderer{g: g, bg: bg, pulse: pulse, arcBg: arcBg, arcVal: arcVal, ticks: ticks, label: label, value: value, unit: unit, needle: needle, peakLine: peakLine}
	r.Refresh()
	return r
}

type ramped3Renderer struct {
	g *Ramped3

	bg    *canvas.Rectangle
	pulse *canvas.Rectangle

	arcBg  []*canvas.Line
	arcVal []*canvas.Line
	ticks  []*canvas.Line

	label  *canvas.Text
	value  *canvas.Text
	unit   *canvas.Text
	needle *canvas.Line

	peakLine *canvas.Line
}

func (r *ramped3Renderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.pulse.Resize(size)

	cfg := r.g.config.Normalize()
	span := cfg.Max - cfg.Min
	if span <= 0 {
		span = 1
	}
	pct := clamp((r.g.value-cfg.Min)/span, 0, 1)
	peakPct := clamp((r.g.peakValue-cfg.Min)/span, 0, 1)

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
	cy := size.Height * 0.72
	radius := math.Min(float64(size.Width), float64(size.Height)) * 0.45
	start := 1.15 * math.Pi
	end := 1.85 * math.Pi
	arc := end - start

	for i := 0; i < rampedSegments; i++ {
		p0 := float64(i) / float64(rampedSegments)
		p1 := float64(i+1) / float64(rampedSegments)
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

	for i := 0; i < rampedTicks; i++ {
		p := float64(i) / float64(rampedTicks-1)
		a := start + p*arc
		outer := radius + 8
		inner := radius - 18
		l := r.ticks[i]
		l.Position1 = fyne.NewPos(cx+float32(inner*math.Cos(a)), cy+float32(inner*math.Sin(a)))
		l.Position2 = fyne.NewPos(cx+float32(outer*math.Cos(a)), cy+float32(outer*math.Sin(a)))
	}

	a := start + pct*arc
	needleLen := radius - 24
	r.needle.Position1 = fyne.NewPos(cx, cy)
	r.needle.Position2 = fyne.NewPos(cx+float32(needleLen*math.Cos(a)), cy+float32(needleLen*math.Sin(a)))

	// Peak arc tick
	segHalf := (arc / float64(rampedSegments)) * 0.9
	pa := start + peakPct*arc
	r.peakLine.Position1 = fyne.NewPos(cx+float32(radius*math.Cos(pa-segHalf)), cy+float32(radius*math.Sin(pa-segHalf)))
	r.peakLine.Position2 = fyne.NewPos(cx+float32(radius*math.Cos(pa+segHalf)), cy+float32(radius*math.Sin(pa+segHalf)))
	r.peakLine.Show()

	// Text
	r.label.Text = cfg.Label
	r.label.TextSize = float32(math.Max(12, float64(size.Height)*0.07))
	r.label.Refresh()
	r.label.Move(fyne.NewPos(14, 10))
	r.label.Resize(fyne.NewSize(size.Width*0.6, r.label.MinSize().Height))

	r.value.Text = formatValue(r.g.value, span)
	r.value.TextSize = float32(math.Max(20, float64(size.Height)*0.16))
	r.value.Refresh()
	r.value.Resize(fyne.NewSize(size.Width*0.5, r.value.MinSize().Height))
	r.value.Move(fyne.NewPos(size.Width-r.value.Size().Width-14, size.Height*0.08))

	r.unit.Text = cfg.Unit
	r.unit.TextSize = float32(math.Max(10, float64(size.Height)*0.06))
	r.unit.Refresh()
	r.unit.Resize(fyne.NewSize(size.Width*0.25, r.unit.MinSize().Height))
	r.unit.Move(fyne.NewPos(size.Width-r.unit.Size().Width-14, size.Height*0.08+r.value.MinSize().Height-4))
}

func (r *ramped3Renderer) MinSize() fyne.Size { return fyne.NewSize(480, 240) }

func (r *ramped3Renderer) Refresh() {
	r.Layout(r.g.Size())
	canvas.Refresh(r.bg)
	canvas.Refresh(r.pulse)
	for _, l := range r.arcBg {
		canvas.Refresh(l)
	}
	for _, l := range r.arcVal {
		canvas.Refresh(l)
	}
	for _, l := range r.ticks {
		canvas.Refresh(l)
	}
	canvas.Refresh(r.peakLine)
	canvas.Refresh(r.needle)
	canvas.Refresh(r.label)
	canvas.Refresh(r.value)
	canvas.Refresh(r.unit)
}

func (r *ramped3Renderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.bg, r.pulse}
	for _, l := range r.arcBg {
		objs = append(objs, l)
	}
	for _, l := range r.arcVal {
		objs = append(objs, l)
	}
	for _, l := range r.ticks {
		objs = append(objs, l)
	}
	objs = append(objs, r.peakLine, r.needle, r.label, r.value, r.unit)
	return objs
}

func (r *ramped3Renderer) Destroy() {}
