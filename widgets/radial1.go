package widgets

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
)

const (
	radial1ProgressSegments = 72
	radial1TickCount        = 50
	radial1PeakFade         = 600 * time.Millisecond
)

// Radial1 is a reusable dark-dashboard radial gauge.
// It is configured entirely by GaugeConfig and updated by SetValue.
type Radial1 struct {
	widget.BaseWidget
	config GaugeConfig

	value     float64
	lastValue float64

	peakValue float64
	peakAt    time.Time
}

func NewRadial1(cfg GaugeConfig) Widget {
	n := &Radial1{config: cfg.Normalize()}
	n.value = n.config.Min
	n.lastValue = n.value
	n.ExtendBaseWidget(n)
	return n
}

func (g *Radial1) Style() string { return "radial1" }

func (g *Radial1) Config() GaugeConfig { return g.config }

func (g *Radial1) Value() float64 { return g.value }

func (g *Radial1) SetValue(value float64) {
	value = clamp(value, g.config.Min, g.config.Max)

	if g.config.ShowPeak {
		span := g.config.Max - g.config.Min
		if span <= 0 {
			span = 1
		}
		drop := span * 0.03 // pragmatic: "noticeable" drop triggers the peak marker
		if g.lastValue-value >= drop {
			g.peakValue = g.lastValue
			g.peakAt = time.Now()
		}
	}
	g.value = value
	g.lastValue = value

	g.Refresh()
}

func (g *Radial1) Snapshot() Snapshot {
	value := g.Value()
	return Snapshot{
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

func (g *Radial1) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(parseHex(g.config.Theme.Background, color.NRGBA{R: 5, G: 7, B: 10, A: 255}))

	dial := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	dial.StrokeColor = parseHex(g.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255})
	dial.StrokeWidth = 4

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

	r := &radial1Renderer{
		gauge:     g,
		bg:        bg,
		dial:      dial,
		ticks:     ticks,
		rangeArc:  rangeArc,
		valueArc:  valueArc,
		labels:    labels,
		labelText: labelText,
		valueText: valueText,
		unitText:  unitText,
		minText:   minText,
		maxText:   maxText,
		peakLine:  peakLine,
		needle:    needle,
		cap:       cap,
	}
	r.Refresh()
	return r
}

type radial1Renderer struct {
	gauge *Radial1

	bg   *canvas.Rectangle
	dial *canvas.Circle

	ticks    []*canvas.Line
	rangeArc []*canvas.Line
	valueArc []*canvas.Line

	labels []*canvas.Text

	labelText *canvas.Text
	valueText *canvas.Text
	unitText  *canvas.Text
	minText   *canvas.Text
	maxText   *canvas.Text

	peakLine *canvas.Line
	needle   []*canvas.Line
	cap      *canvas.Circle
}

func (r *radial1Renderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.dial.Resize(size)

	cx, cy := size.Width/2, size.Height/2
	radius := math.Min(float64(cx), float64(cy)) - 20

	startAngle := 0.75 * math.Pi
	endAngle := 2.25 * math.Pi
	angleRange := endAngle - startAngle

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
		label.TextSize = float32(math.Max(10, radius/24))
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

	// Peak marker
	if !cfg.ShowPeak || r.gauge.peakAt.IsZero() {
		r.peakLine.Hide()
	} else {
		fade := 1.0 - float64(time.Since(r.gauge.peakAt))/float64(radial1PeakFade)
		if fade <= 0 {
			r.peakLine.Hide()
		} else {
			peakPct := clamp((r.gauge.peakValue-cfg.Min)/span, 0, 1)
			peakAngle := startAngle + peakPct*angleRange
			r.peakLine.Position1 = fyne.NewPos(cx, cy)
			r.peakLine.Position2 = fyne.NewPos(
				cx+float32((radius-30)*math.Cos(peakAngle)),
				cy+float32((radius-30)*math.Sin(peakAngle)),
			)
			base := parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
			r.peakLine.StrokeColor = withAlpha(base, uint8(255*fade))
			r.peakLine.Show()
		}
	}

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

func (r *radial1Renderer) layoutNeedle(cx, cy float32, radius, angle float64) {
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

func (r *radial1Renderer) MinSize() fyne.Size { return fyne.NewSize(480, 480) }

func (r *radial1Renderer) Refresh() {
	r.Layout(r.gauge.Size())
	canvas.Refresh(r.bg)
	canvas.Refresh(r.dial)
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
	canvas.Refresh(r.labelText)
	canvas.Refresh(r.valueText)
	canvas.Refresh(r.unitText)
	canvas.Refresh(r.minText)
	canvas.Refresh(r.maxText)
	canvas.Refresh(r.peakLine)
	for _, line := range r.needle {
		canvas.Refresh(line)
	}
	canvas.Refresh(r.cap)
}

func (r *radial1Renderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		r.bg,
		r.dial,
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
	objects = append(objects,
		r.peakLine,
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

func (r *radial1Renderer) Destroy() {}

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
