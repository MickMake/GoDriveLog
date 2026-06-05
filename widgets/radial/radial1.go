package radial

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

const (
	radial1ProgressSegments = 72
	radial1TickCount        = 50

	radialTrailLen  = 28
	radialTrailFade = 900 * time.Millisecond
)

type radialMode int

const (
	radialModePlain radialMode = iota
	radialModeTrail
	radialModePeakDaily
)

type trailSample struct {
	value float64
	at    time.Time
}

// Radial is a reusable dark-dashboard radial gauge.
// It is configured entirely by model.GaugeConfig and updated by SetValue.
type Radial struct {
	widget.BaseWidget
	style  string
	mode   radialMode
	config model.GaugeConfig

	startAngle float64
	endAngle   float64

	value float64 // displayed (smoothed) value

	smoothBuf   []float64
	smoothNext  int
	smoothCount int

	trail     [radialTrailLen]trailSample
	trailNext int

	peakValue float64
	peakDay   string // YYYY-MM-DD

	pulse model.PulseTracker
}

func NewRadial(style string, mode radialMode, cfg model.GaugeConfig) model.Widget {
	// Default: classic 270-degree sweep.
	return NewRadialArc(style, mode, cfg, 0.75*math.Pi, 2.25*math.Pi)
}

func NewRadialArc(style string, mode radialMode, cfg model.GaugeConfig, startAngle, endAngle float64) model.Widget {
	cfg = cfg.Normalize()

	w := cfg.SmoothingWindow
	if w <= 1 {
		w = 1
	}

	n := &Radial{style: style, mode: mode, config: cfg, startAngle: startAngle, endAngle: endAngle, pulse: model.NewPulseTracker()}
	n.value = n.config.Min
	n.peakValue = n.value
	n.peakDay = time.Now().Format("2006-01-02")
	// Prime smoothing buffer so the first few samples don't jitter.
	n.smoothBuf = make([]float64, w)
	n.smoothBuf[0] = n.value
	n.smoothCount = 1

	n.ExtendBaseWidget(n)
	return n
}

// Classic full-size radial variants.
func NewRadial1(cfg model.GaugeConfig) model.Widget { return NewRadial("radial1", radialModePlain, cfg) }
func NewRadial2(cfg model.GaugeConfig) model.Widget { return NewRadial("radial2", radialModeTrail, cfg) }
func NewRadial3(cfg model.GaugeConfig) model.Widget { return NewRadial("radial3", radialModePeakDaily, cfg) }

// Half radials.
func NewHalfTop1(cfg model.GaugeConfig) model.Widget {
	// Top half: left->right across top.
	return NewRadialArc("half_top1", radialModePlain, cfg, math.Pi, 2*math.Pi)
}
func NewHalfBottom1(cfg model.GaugeConfig) model.Widget {
	return NewRadialArc("half_bottom1", radialModePlain, cfg, 0, math.Pi)
}

// Quarter radials (quadrants).
func NewQuarterTL1(cfg model.GaugeConfig) model.Widget {
	return NewRadialArc("quarter_tl1", radialModePlain, cfg, math.Pi, 1.5*math.Pi)
}
func NewQuarterTR1(cfg model.GaugeConfig) model.Widget {
	return NewRadialArc("quarter_tr1", radialModePlain, cfg, 1.5*math.Pi, 2*math.Pi)
}
func NewQuarterBL1(cfg model.GaugeConfig) model.Widget {
	return NewRadialArc("quarter_bl1", radialModePlain, cfg, 0.5*math.Pi, math.Pi)
}
func NewQuarterBR1(cfg model.GaugeConfig) model.Widget {
	return NewRadialArc("quarter_br1", radialModePlain, cfg, 0, 0.5*math.Pi)
}

func (g *Radial) Style() string { return g.style }

func (g *Radial) Config() model.GaugeConfig { return g.config }

func (g *Radial) Value() float64 { return g.value }

func (g *Radial) SetValue(value float64) {
	value = clamp(value, g.config.Min, g.config.Max)

	smoothed := g.smooth(value)
	g.value = smoothed

	// Built-in alert pulse (only does anything if ranges are set).
	g.pulse.Update(g.value, g.config.WarningRange, g.config.DangerRange)

	switch g.mode {
	case radialModeTrail:
		g.pushTrail(smoothed)
	case radialModePeakDaily:
		g.updatePeakDaily(smoothed)
	}

	g.Refresh()
}

func (g *Radial) smooth(value float64) float64 {
	if len(g.smoothBuf) <= 1 {
		return value
	}

	g.smoothBuf[g.smoothNext] = value
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

func (g *Radial) pushTrail(value float64) {
	g.trail[g.trailNext] = trailSample{value: value, at: time.Now()}
	g.trailNext = (g.trailNext + 1) % radialTrailLen
}

func (g *Radial) updatePeakDaily(value float64) {
	today := time.Now().Format("2006-01-02")
	if today != g.peakDay {
		g.peakDay = today
		g.peakValue = value
		return
	}
	if value > g.peakValue {
		g.peakValue = value
	}
}

func (g *Radial) Snapshot() model.Snapshot {
	value := g.Value()
	return model.Snapshot{
		Style:      g.Style(),
		Label:      g.config.Label,
		Unit:       g.config.Unit,
		Min:        g.config.Min,
		Max:        g.config.Max,
		Value:      value,
		Normalised: normalise(value, g.config.Min, g.config.Max),
		Warning:    inRange(value, g.config.WarningRange),
		Danger:     inRange(value, g.config.DangerRange),
	}
}

func (g *Radial) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(parseHex(g.config.Theme.Background, color.NRGBA{R: 5, G: 7, B: 10, A: 255}))

	dial := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	dial.StrokeColor = parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255})
	dial.StrokeWidth = 4

	pulseRing := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pulseRing.StrokeColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	pulseRing.StrokeWidth = 10
	pulseRing.Hide()

	ticks := make([]*canvas.Line, 0, radial1TickCount+1)
	for i := 0; i <= radial1TickCount; i++ {
		line := canvas.NewLine(parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}))
		switch {
		case i%10 == 0:
			line.StrokeWidth = 4
		case i%5 == 0:
			line.StrokeWidth = 3
		default:
			line.StrokeWidth = 2
		}
		ticks = append(ticks, line)
	}

	rangeArc := make([]*canvas.Line, 0, radial1ProgressSegments)
	valueArc := make([]*canvas.Line, 0, radial1ProgressSegments)
	for i := 0; i < radial1ProgressSegments; i++ {
		r := canvas.NewLine(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		r.StrokeWidth = 7
		r.Hide()
		rangeArc = append(rangeArc, r)

		v := canvas.NewLine(parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
		v.StrokeWidth = 7
		v.Hide()
		valueArc = append(valueArc, v)
	}

	labelText := canvas.NewText(g.config.Label, parseHex(g.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	labelText.Alignment = fyne.TextAlignCenter
	labelText.TextStyle = fyne.TextStyle{Bold: true}

	valueText := canvas.NewText("", parseHex(g.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
	valueText.Alignment = fyne.TextAlignCenter
	valueText.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	unitText := canvas.NewText(g.config.Unit, parseHex(g.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	unitText.Alignment = fyne.TextAlignCenter
	unitText.TextStyle = fyne.TextStyle{Bold: true}

	minText := canvas.NewText("", parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}))
	minText.Alignment = fyne.TextAlignLeading

	maxText := canvas.NewText("", parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}))
	maxText.Alignment = fyne.TextAlignTrailing

	labels := make([]*canvas.Text, 0, 6)
	for i := 0; i < 6; i++ {
		lt := canvas.NewText("", parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}))
		lt.Alignment = fyne.TextAlignCenter
		lt.TextStyle = fyne.TextStyle{Bold: true}
		labels = append(labels, lt)
	}

	trailLines := make([]*canvas.Line, 0, radialTrailLen)
	for i := 0; i < radialTrailLen; i++ {
		line := canvas.NewLine(parseHex(g.config.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255}))
		line.StrokeWidth = 2
		line.Hide()
		trailLines = append(trailLines, line)
	}

	peakLine := canvas.NewLine(parseHex(g.config.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255}))
	peakLine.StrokeWidth = 2
	peakLine.Hide()

	needle := make([]*canvas.Line, 3)
	for i := range needle {
		needle[i] = canvas.NewLine(parseHex(g.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255}))
		needle[i].StrokeWidth = 3
	}

	cap := canvas.NewCircle(parseHex(g.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255}))
	cap.StrokeColor = parseHex(g.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255})
	cap.StrokeWidth = 2

	r := &radialRenderer{
		gauge:     g,
		bg:        bg,
		dial:      dial,
		pulseRing: pulseRing,
		ticks:     ticks,
		rangeArc:  rangeArc,
		valueArc:  valueArc,
		labels:    labels,
		labelText: labelText,
		valueText: valueText,
		unitText:  unitText,
		minText:   minText,
		maxText:   maxText,
		trailLines: trailLines,
		peakLine:   peakLine,
		needle:     needle,
		cap:        cap,
	}
	r.Refresh()
	return r
}

type radialRenderer struct {
	gauge *Radial

	bg        *canvas.Rectangle
	dial      *canvas.Circle
	pulseRing *canvas.Circle

	ticks    []*canvas.Line
	rangeArc []*canvas.Line
	valueArc []*canvas.Line

	labels []*canvas.Text

	labelText *canvas.Text
	valueText *canvas.Text
	unitText  *canvas.Text
	minText   *canvas.Text
	maxText   *canvas.Text

	trailLines []*canvas.Line
	peakLine   *canvas.Line
	needle     []*canvas.Line
	cap        *canvas.Circle
}

func (r *radialRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.dial.Resize(size)

	// Center positioning: for partial arcs, move the center away from the arc midpoint
	// so the active sweep uses more of the available space.
	startAngle := r.gauge.startAngle
	endAngle := r.gauge.endAngle
	angleRange := endAngle - startAngle
	mid := startAngle + angleRange/2
	vx := math.Cos(mid)
	vy := math.Sin(mid)

	margin := float32(18)
	cx := size.Width/2 - float32(vx)*size.Width*0.18
	cy := size.Height/2 - float32(vy)*size.Height*0.18

	cx = clampF(cx, margin, size.Width-margin)
	cy = clampF(cy, margin, size.Height-margin)

	radius := math.Min(float64(size.Width), float64(size.Height))*0.45

	// Ring bounds
	r.dial.Resize(fyne.NewSize(float32(radius*2), float32(radius*2)))
	r.dial.Move(fyne.NewPos(cx-float32(radius), cy-float32(radius)))
	r.pulseRing.Resize(r.dial.Size())
	r.pulseRing.Move(r.dial.Position())

	// Pulse overlay
	pState, p := r.gauge.pulse.Pulse(time.Now())
	if p > 0 {
		col := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		switch pState {
		case model.AlertWarning:
			col = parseHex(r.gauge.config.Theme.Warning, color.NRGBA{R: 255, G: 176, B: 0, A: 255})
		case model.AlertDanger:
			col = parseHex(r.gauge.config.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
		}
		r.pulseRing.StrokeColor = withAlpha(col, uint8(220*p))
		r.pulseRing.Show()
	} else {
		r.pulseRing.Hide()
	}

	// Ticks
	for i, tick := range r.ticks {
		if !r.gauge.config.ShowTicks {
			tick.Hide()
			continue
		}
		pct := float64(i) / float64(len(r.ticks)-1)
		angle := startAngle + pct*angleRange

		inner := radius - 12
		if i%10 == 0 {
			inner = radius - 24
		} else if i%5 == 0 {
			inner = radius - 18
		}

		tick.Position1 = fyne.NewPos(
			cx+float32(inner*math.Cos(angle)),
			cy+float32(inner*math.Sin(angle)),
		)
		tick.Position2 = fyne.NewPos(
			cx+float32(radius*math.Cos(angle)),
			cy+float32(radius*math.Sin(angle)),
		)
		tick.Show()
	}

	// Range arc + value arc
	cfg := r.gauge.config.Normalize()
	span := cfg.Max - cfg.Min
	if span <= 0 {
		span = 1
	}
	currentPct := clamp((r.gauge.value-cfg.Min)/span, 0, 1)
	arcRadius := radius - 6

	warn := cfg.WarningRange
	danger := cfg.DangerRange

	for i := 0; i < radial1ProgressSegments; i++ {
		segStartPct := float64(i) / float64(radial1ProgressSegments)
		segEndPct := float64(i+1) / float64(radial1ProgressSegments)

		midValue := cfg.Min + ((segStartPct+segEndPct)/2)*span
		segStart := startAngle + segStartPct*angleRange
		segEnd := startAngle + segEndPct*angleRange

		pos1 := fyne.NewPos(
			cx+float32(arcRadius*math.Cos(segStart)),
			cy+float32(arcRadius*math.Sin(segStart)),
		)
		pos2 := fyne.NewPos(
			cx+float32(arcRadius*math.Cos(segEnd)),
			cy+float32(arcRadius*math.Sin(segEnd)),
		)

		// Background range arc
		bgSeg := r.rangeArc[i]
		bgSeg.Position1, bgSeg.Position2 = pos1, pos2
		if inRange(midValue, danger) {
			bgSeg.StrokeColor = withAlpha(parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255}), 90)
			bgSeg.Show()
		} else if inRange(midValue, warn) {
			bgSeg.StrokeColor = withAlpha(parseHex(cfg.Theme.Warning, color.NRGBA{R: 255, G: 176, B: 0, A: 255}), 90)
			bgSeg.Show()
		} else {
			bgSeg.Hide()
		}

		// Foreground value arc
		fgSeg := r.valueArc[i]
		if segStartPct >= currentPct {
			fgSeg.Hide()
			continue
		}
		if segEndPct > currentPct {
			segEndPct = currentPct
			segEnd = startAngle + segEndPct*angleRange
			pos2 = fyne.NewPos(
				cx+float32(arcRadius*math.Cos(segEnd)),
				cy+float32(arcRadius*math.Sin(segEnd)),
			)
		}
		fgSeg.Position1, fgSeg.Position2 = pos1, pos2
		if inRange(midValue, danger) {
			fgSeg.StrokeColor = parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
		} else if inRange(midValue, warn) {
			fgSeg.StrokeColor = parseHex(cfg.Theme.Warning, color.NRGBA{R: 255, G: 176, B: 0, A: 255})
		} else {
			fgSeg.StrokeColor = parseHex(cfg.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255})
		}
		fgSeg.Show()
	}

	// Major labels
	for i, label := range r.labels {
		if !cfg.ShowMajorLabels {
			label.Hide()
			continue
		}
		pct := float64(i) / float64(len(r.labels)-1)
		value := cfg.Min + pct*span
		label.Text = formatTick(value, span)
		label.TextSize = float32(math.Max(12, radius/22))
		label.Refresh()

		angle := startAngle + pct*angleRange
		labelRadius := radius - 48
		size := fyne.NewSize(60, 18)
		label.Resize(size)
		label.Move(fyne.NewPos(
			cx+float32(labelRadius*math.Cos(angle))-size.Width/2,
			cy+float32(labelRadius*math.Sin(angle))-size.Height/2,
		))
		label.Show()
	}

	// Trail / peak overlays
	r.layoutTrail(cfg, cx, cy, radius, startAngle, angleRange, span)
	r.layoutPeak(cfg, cx, cy, radius, startAngle, angleRange, span)

	// Needle
	currentAngle := startAngle + currentPct*angleRange
	r.layoutNeedle(cx, cy, radius, currentAngle)

	// Label/value/unit + min/max
	r.labelText.Text = cfg.Label
	r.labelText.TextSize = float32(math.Max(12, radius/18))
	r.labelText.Refresh()

	if cfg.ShowLabel {
		r.labelText.Resize(fyne.NewSize(size.Width, 24))
		r.labelText.Move(fyne.NewPos(0, float32(size.Height*0.06)))
		r.labelText.Show()
	} else {
		r.labelText.Hide()
	}

	if cfg.ShowValue {
		r.valueText.Text = formatValue(r.gauge.value, span)
		r.valueText.TextSize = float32(math.Max(18, radius/9))
		r.valueText.Refresh()

		valueY := float32(size.Height*0.66) - r.valueText.MinSize().Height/2
		r.valueText.Resize(fyne.NewSize(size.Width, r.valueText.MinSize().Height))
		r.valueText.Move(fyne.NewPos(0, valueY))
		r.valueText.Show()
	} else {
		r.valueText.Hide()
	}

	if cfg.ShowUnit && strings.TrimSpace(cfg.Unit) != "" {
		r.unitText.Text = cfg.Unit
		r.unitText.TextSize = float32(math.Max(10, radius/20))
		r.unitText.Refresh()

		y := float32(size.Height*0.66) + r.valueText.MinSize().Height/2 + 2
		r.unitText.Resize(fyne.NewSize(size.Width, r.unitText.MinSize().Height))
		r.unitText.Move(fyne.NewPos(0, y))
		r.unitText.Show()
	} else {
		r.unitText.Hide()
	}

	r.minText.TextSize = float32(math.Max(10, radius/26))
	r.maxText.TextSize = r.minText.TextSize
	r.minText.Text = formatTick(cfg.Min, span)
	r.maxText.Text = formatTick(cfg.Max, span)
	r.minText.Refresh()
	r.maxText.Refresh()

	if cfg.ShowMin {
		r.minText.Resize(fyne.NewSize(size.Width/2, r.minText.MinSize().Height))
		r.minText.Move(fyne.NewPos(12, size.Height-r.minText.MinSize().Height-10))
		r.minText.Show()
	} else {
		r.minText.Hide()
	}

	if cfg.ShowMax {
		r.maxText.Resize(fyne.NewSize(size.Width/2-12, r.maxText.MinSize().Height))
		r.maxText.Move(fyne.NewPos(size.Width/2, size.Height-r.maxText.MinSize().Height-10))
		r.maxText.Show()
	} else {
		r.maxText.Hide()
	}
}

func (r *radialRenderer) layoutTrail(cfg model.GaugeConfig, cx, cy float32, radius float64, startAngle, angleRange, span float64) {
	if r.gauge.mode != radialModeTrail {
		for _, line := range r.trailLines {
			line.Hide()
		}
		return
	}

	now := time.Now()
	base := parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
	n := len(r.trailLines)
	for i := 0; i < n; i++ {
		idx := (r.gauge.trailNext - 1 - i) % radialTrailLen
		if idx < 0 {
			idx += radialTrailLen
		}
		sample := r.gauge.trail[idx]
		line := r.trailLines[i]

		if sample.at.IsZero() {
			line.Hide()
			continue
		}

		age := now.Sub(sample.at)
		if age >= radialTrailFade {
			line.Hide()
			continue
		}

		pct := clamp((sample.value-cfg.Min)/span, 0, 1)
		angle := startAngle + pct*angleRange

		line.Position1 = fyne.NewPos(cx, cy)
		line.Position2 = fyne.NewPos(
			cx+float32((radius-30)*math.Cos(angle)),
			cy+float32((radius-30)*math.Sin(angle)),
		)

		fade := 1.0 - float64(age)/float64(radialTrailFade)
		line.StrokeColor = withAlpha(base, uint8(255*fade))
		line.Show()
	}
}

func (r *radialRenderer) layoutPeak(cfg model.GaugeConfig, cx, cy float32, radius float64, startAngle, angleRange, span float64) {
	if r.gauge.mode != radialModePeakDaily {
		r.peakLine.Hide()
		return
	}

	pct := clamp((r.gauge.peakValue-cfg.Min)/span, 0, 1)
	angle := startAngle + pct*angleRange

	r.peakLine.Position1 = fyne.NewPos(cx, cy)
	r.peakLine.Position2 = fyne.NewPos(
		cx+float32((radius-30)*math.Cos(angle)),
		cy+float32((radius-30)*math.Sin(angle)),
	)
	r.peakLine.Show()
}

func (r *radialRenderer) layoutNeedle(cx, cy float32, radius, angle float64) {
	tipRadius := radius - 30
	baseBack := 18.0
	baseHalfWidth := 7.0

	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	perpX := -sinA
	perpY := cosA

	tip := fyne.NewPos(
		cx+float32(tipRadius*cosA),
		cy+float32(tipRadius*sinA),
	)
	leftBase := fyne.NewPos(
		cx+float32(-baseBack*cosA+baseHalfWidth*perpX),
		cy+float32(-baseBack*sinA+baseHalfWidth*perpY),
	)
	rightBase := fyne.NewPos(
		cx+float32(-baseBack*cosA-baseHalfWidth*perpX),
		cy+float32(-baseBack*sinA-baseHalfWidth*perpY),
	)

	r.needle[0].Position1 = leftBase
	r.needle[0].Position2 = tip
	r.needle[1].Position1 = tip
	r.needle[1].Position2 = rightBase
	r.needle[2].Position1 = rightBase
	r.needle[2].Position2 = leftBase

	capSize := float32(18)
	r.cap.Resize(fyne.NewSize(capSize, capSize))
	r.cap.Move(fyne.NewPos(cx-capSize/2, cy-capSize/2))
}

func (r *radialRenderer) MinSize() fyne.Size { return fyne.NewSize(480, 480) }

func (r *radialRenderer) Refresh() {
	r.Layout(r.gauge.Size())
	canvas.Refresh(r.bg)
	canvas.Refresh(r.dial)
	canvas.Refresh(r.pulseRing)
	for _, tick := range r.ticks {
		canvas.Refresh(tick)
	}
	for _, seg := range r.rangeArc {
		canvas.Refresh(seg)
	}
	for _, seg := range r.valueArc {
		canvas.Refresh(seg)
	}
	for _, label := range r.labels {
		canvas.Refresh(label)
	}
	for _, line := range r.trailLines {
		canvas.Refresh(line)
	}
	canvas.Refresh(r.peakLine)
	canvas.Refresh(r.labelText)
	canvas.Refresh(r.valueText)
	canvas.Refresh(r.unitText)
	canvas.Refresh(r.minText)
	canvas.Refresh(r.maxText)
	for _, line := range r.needle {
		canvas.Refresh(line)
	}
	canvas.Refresh(r.cap)
}

func (r *radialRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		r.bg,
		r.dial,
		r.pulseRing,
	}
	for _, seg := range r.rangeArc {
		objects = append(objects, seg)
	}
	for _, seg := range r.valueArc {
		objects = append(objects, seg)
	}
	for _, tick := range r.ticks {
		objects = append(objects, tick)
	}
	for _, label := range r.labels {
		objects = append(objects, label)
	}
	objects = append(objects, r.peakLine)
	for _, line := range r.trailLines {
		objects = append(objects, line)
	}
	objects = append(objects,
		r.labelText,
		r.valueText,
		r.unitText,
		r.minText,
		r.maxText,
	)
	for _, line := range r.needle {
		objects = append(objects, line)
	}
	objects = append(objects, r.cap)
	return objects
}

func (r *radialRenderer) Destroy() {}

func formatValue(value, span float64) string {
	if span >= 1000 {
		return fmt.Sprintf("%.0f", value)
	}
	if span >= 100 {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.2f", value)
}

func formatTick(value, span float64) string {
	if span >= 1000 {
		return fmt.Sprintf("%.0f", value)
	}
	if span >= 100 {
		return fmt.Sprintf("%.0f", value)
	}
	if span >= 10 {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.2f", value)
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

func withAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = a
	return c
}

func clamp(value, min, max float64) float64 { return math.Max(min, math.Min(max, value)) }

func clampF(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

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
