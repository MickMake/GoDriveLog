package bar

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
	bar2Segments = 24
)

// Bar2 is a segmented horizontal "thermometer" bar with a value marker.
// Good for load/boost/temps.
type Bar2 struct {
	widget.BaseWidget
	config model.GaugeConfig
	value  float64

	smoothBuf   []float64
	smoothNext  int
	smoothCount int

	pulse model.PulseTracker
}

func NewBar2(cfg model.GaugeConfig) model.Widget {
	cfg = cfg.Normalize()
	w := cfg.SmoothingWindow
	if w <= 1 {
		w = 1
	}
	b := &Bar2{config: cfg, value: cfg.Min, smoothBuf: make([]float64, w), pulse: model.NewPulseTracker()}
	b.smoothBuf[0] = b.value
	b.smoothCount = 1
	b.ExtendBaseWidget(b)
	return b
}

func (b *Bar2) Style() string { return "bar2" }

func (b *Bar2) Config() model.GaugeConfig { return b.config }

func (b *Bar2) Value() float64 { return b.value }

func (b *Bar2) SetValue(v float64) {
	v = clamp(v, b.config.Min, b.config.Max)
	b.value = b.smooth(v)
	b.pulse.Update(b.value, b.config.WarningRange, b.config.DangerRange)
	b.Refresh()
}

func (b *Bar2) smooth(v float64) float64 {
	if len(b.smoothBuf) <= 1 {
		return v
	}
	b.smoothBuf[b.smoothNext] = v
	b.smoothNext = (b.smoothNext + 1) % len(b.smoothBuf)
	if b.smoothCount < len(b.smoothBuf) {
		b.smoothCount++
	}
	var sum float64
	for i := 0; i < b.smoothCount; i++ {
		sum += b.smoothBuf[i]
	}
	return sum / float64(b.smoothCount)
}

func (b *Bar2) Snapshot() model.Snapshot {
	v := b.Value()
	return model.Snapshot{
		Style:      b.Style(),
		Label:      b.config.Label,
		Unit:       b.config.Unit,
		Min:        b.config.Min,
		Max:        b.config.Max,
		Value:      v,
		Normalised: normalise(v, b.config.Min, b.config.Max),
		Warning:    inRange(v, b.config.WarningRange),
		Danger:     inRange(v, b.config.DangerRange),
	}
}

func (b *Bar2) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(parseHex(b.config.Theme.Background, color.NRGBA{R: 5, G: 7, B: 10, A: 255}))
	pulse := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pulse.Hide()

	segments := make([]*canvas.Rectangle, 0, bar2Segments)
	for i := 0; i < bar2Segments; i++ {
		r := canvas.NewRectangle(withAlpha(parseHex(b.config.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}), 55))
		r.CornerRadius = 3
		segments = append(segments, r)
	}

	marker := canvas.NewRectangle(parseHex(b.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255}))
	marker.CornerRadius = 2

	label := canvas.NewText(b.config.Label, parseHex(b.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	label.TextStyle = fyne.TextStyle{Bold: true}

	value := canvas.NewText("", parseHex(b.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	value.Alignment = fyne.TextAlignTrailing

	r := &bar2Renderer{b: b, bg: bg, pulse: pulse, segments: segments, marker: marker, label: label, value: value}
	r.Refresh()
	return r
}

type bar2Renderer struct {
	b *Bar2

	bg       *canvas.Rectangle
	pulse    *canvas.Rectangle
	segments []*canvas.Rectangle
	marker   *canvas.Rectangle
	label    *canvas.Text
	value    *canvas.Text
}

func (r *bar2Renderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.pulse.Resize(size)

	cfg := r.b.config.Normalize()
	span := cfg.Max - cfg.Min
	if span <= 0 {
		span = 1
	}
	pct := clamp((r.b.value-cfg.Min)/span, 0, 1)

	pState, p := r.b.pulse.Pulse(time.Now())
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

	pad := float32(14)
	top := float32(10)
	barY := size.Height * 0.55
	barH := float32(math.Max(18, float64(size.Height)*0.22))
	barW := size.Width - pad*2

	// Text
	r.label.Text = cfg.Label
	r.label.TextSize = float32(math.Max(12, float64(size.Height)*0.09))
	r.label.Refresh()
	r.label.Move(fyne.NewPos(pad, top))

	r.value.Text = fmt.Sprintf("%s %s", formatValue(r.b.value, span), cfg.Unit)
	r.value.TextSize = r.label.TextSize
	r.value.Refresh()
	r.value.Resize(fyne.NewSize(size.Width-pad*2, r.value.MinSize().Height))
	r.value.Move(fyne.NewPos(pad, top))

	segGap := float32(4)
	segW := (barW - segGap*float32(len(r.segments)-1)) / float32(len(r.segments))

	filled := int(math.Round(pct * float64(len(r.segments))))
	if filled < 0 {
		filled = 0
	}
	if filled > len(r.segments) {
		filled = len(r.segments)
	}

	for i, s := range r.segments {
		x := pad + float32(i)*(segW+segGap)
		s.Resize(fyne.NewSize(segW, barH))
		s.Move(fyne.NewPos(x, barY))

		midValue := cfg.Min + (float64(i)+0.5)/float64(len(r.segments))*span
		col := withAlpha(parseHex(cfg.Theme.Tick, color.NRGBA{R: 127, G: 138, B: 153, A: 255}), 55)
		if i < filled {
			if inRange(midValue, cfg.DangerRange) {
				col = parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
			} else if inRange(midValue, cfg.WarningRange) {
				col = parseHex(cfg.Theme.Warning, color.NRGBA{R: 255, G: 176, B: 0, A: 255})
			} else {
				col = parseHex(cfg.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255})
			}
		}
		s.FillColor = col
		s.Refresh()
	}

	// Marker line
	mx := pad + pct*barW
	r.marker.Resize(fyne.NewSize(4, barH+10))
	r.marker.Move(fyne.NewPos(mx-2, barY-5))
}

func (r *bar2Renderer) MinSize() fyne.Size { return fyne.NewSize(420, 120) }

func (r *bar2Renderer) Refresh() {
	r.Layout(r.b.Size())
	canvas.Refresh(r.bg)
	canvas.Refresh(r.pulse)
	for _, s := range r.segments {
		canvas.Refresh(s)
	}
	canvas.Refresh(r.marker)
	canvas.Refresh(r.label)
	canvas.Refresh(r.value)
}

func (r *bar2Renderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.bg, r.pulse, r.label, r.value}
	for _, s := range r.segments {
		objs = append(objs, s)
	}
	objs = append(objs, r.marker)
	return objs
}

func (r *bar2Renderer) Destroy() {}

func formatValue(value, span float64) string {
	if span >= 1000 {
		return fmt.Sprintf("%.0f", value)
	}
	if span >= 100 {
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
