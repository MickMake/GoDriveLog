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

	radial1TrailLen        = 28
	radial1TrailFade       = 900 * time.Millisecond
	radial1SmoothingWindow = 3
)

type trailSample struct {
	value float64
	at    time.Time
}

// Radial1 is a reusable dark-dashboard radial gauge.
// It is configured entirely by model.GaugeConfig and updated by SetValue.
type Radial1 struct {
	widget.BaseWidget
	config model.GaugeConfig

	value float64 // displayed (smoothed) value

	smoothBuf   [radial1SmoothingWindow]float64
	smoothNext  int
	smoothCount int

	trail     [radial1TrailLen]trailSample
	trailNext int
}

func NewRadial1(cfg model.GaugeConfig) model.Widget {
	n := &Radial1{config: cfg.Normalize()}
	n.value = n.config.Min
	// Prime smoothing buffer so the first few samples don't jitter.
	n.smoothBuf[0] = n.value
	n.smoothCount = 1
	n.ExtendBaseWidget(n)
	return n
}

func (g *Radial1) Style() string { return "radial1" }

func (g *Radial1) Config() model.GaugeConfig { return g.config }

func (g *Radial1) Value() float64 { return g.value }

func (g *Radial1) SetValue(value float64) {
	value = clamp(value, g.config.Min, g.config.Max)

	smoothed := g.smooth(value)
	g.value = smoothed
	g.pushTrail(smoothed)

	g.Refresh()
}

func (g *Radial1) smooth(value float64) float64 {
	g.smoothBuf[g.smoothNext] = value
	g.smoothNext = (g.smoothNext + 1) % radial1SmoothingWindow
	if g.smoothCount < radial1SmoothingWindow {
		g.smoothCount++
	}

	var sum float64
	for i := 0; i < g.smoothCount; i++ {
		sum += g.smoothBuf[i]
	}
	return sum / float64(g.smoothCount)
}

func (g *Radial1) pushTrail(value float64) {
	// Backwards compatible: config.ShowPeak now means "show trail".
	if !g.config.ShowPeak {
		return
	}
	g.trail[g.trailNext] = trailSample{value: value, at: time.Now()}
	g.trailNext = (g.trailNext + 1) % radial1TrailLen
}

func (g *Radial1) Snapshot() model.Snapshot {
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

	trailLines := make([]*canvas.Line, 0, radial1TrailLen)
	for i := 0; i < radial1TrailLen; i++ {
		line := canvas.NewLine(parseHex(g.config.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255}))
		line.StrokeWidth = 2
		line.Hide()
		trailLines = append(trailLines, line)
	}

	needle := make([]*canvas.Line, 3)
	for i := range needle {
		needle[i] = canvas.NewLine(parseHex(g.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255}))
		needle[i].StrokeWidth = 3
	}

	cap := canvas.NewCircle(parseHex(g.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255}))
	cap.StrokeColor = parseHex(g.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255})
	cap.StrokeWidth = 2

	r := &radial1Renderer{
		gauge:      g,
		bg:         bg,
		dial:       dial,
		ticks:      ticks,
		rangeArc:   rangeArc,
		valueArc:   valueArc,
		labels:     labels,
		labelText:  labelText,
		valueText:  valueText,
		unitText:   unitText,
		minText:    minText,
		maxText:    maxText,
		trailLines: trailLines,
		needle:     needle,
		cap:        cap,
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

	trailLines []*canvas.Line
	needle     []*canvas.Line
	cap        *canvas.Circle
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
		// Slightly larger than before.
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

	// Needle trail
	r.layoutTrail(cfg, cx, cy, radius, startAngle, angleRange, span)

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

func (r *radial1Renderer) layoutTrail(cfg model.GaugeConfig, cx, cy float32, radius float64, startAngle, angleRange, span float64) {
	if !cfg.ShowPeak {
		for _, line := range r.trailLines {
			line.Hide()
		}
		return
	}

	now := time.Now()
	base := parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
	n := len(r.trailLines)
	for i := 0; i < n; i++ {
		idx := (r.gauge.trailNext - 1 - i) % radial1TrailLen
		if idx < 0 {
			idx += radial1TrailLen
		}
		sample := r.gauge.trail[idx]
		line := r.trailLines[i]

		if sample.at.IsZero() {
			line.Hide()
			continue
		}

		age := now.Sub(sample.at)
		if age >= radial1TrailFade {
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

		fade := 1.0 - float64(age)/float64(radial1TrailFade)
		line.StrokeColor = withAlpha(base, uint8(255*fade))
		line.Show()
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
	for _, line := range r.trailLines {
		canvas.Refresh(line)
	}
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

func clamp(value, min, max float64) float64 { return math.Max(min, math.Min(max, value)) }

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
