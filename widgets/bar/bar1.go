package bar

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/widgets/model"
)

// StyleBar1 is the first reusable horizontal bar gauge style.
const StyleBar1 = "bar1"

// Bar1 is a simple continuous horizontal bar gauge.
type Bar1 struct {
	widget.BaseWidget
	config model.GaugeConfig
	value  float64

	smoothBuf   []float64
	smoothNext  int
	smoothCount int

	pulse model.PulseTracker
}

func NewBar1(cfg model.GaugeConfig) model.Widget {
	cfg = cfg.Normalize()
	w := cfg.SmoothingWindow
	if w <= 1 {
		w = 1
	}
	b := &Bar1{config: cfg, value: cfg.Min, smoothBuf: make([]float64, w), pulse: model.NewPulseTracker()}
	b.smoothBuf[0] = b.value
	b.smoothCount = 1
	b.ExtendBaseWidget(b)
	return b
}

func (b *Bar1) Style() string { return StyleBar1 }

func (b *Bar1) Config() model.GaugeConfig { return b.config }

func (b *Bar1) Value() float64 { return b.value }

func (b *Bar1) SetValue(v float64) {
	v = clamp(v, b.config.Min, b.config.Max)
	b.value = b.smooth(v)
	b.pulse.Update(b.value, b.config.WarningRange, b.config.DangerRange)
	b.Refresh()
}

func (b *Bar1) smooth(v float64) float64 {
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

func (b *Bar1) Snapshot() model.Snapshot {
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

func (b *Bar1) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(parseHex(b.config.Theme.Background, color.NRGBA{R: 5, G: 7, B: 10, A: 255}))
	pulse := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pulse.Hide()
	track := canvas.NewRectangle(withAlpha(parseHex(b.config.Theme.Grid, color.NRGBA{R: 38, G: 50, B: 65, A: 255}), 210))
	track.CornerRadius = 4
	fill := canvas.NewRectangle(parseHex(b.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
	fill.CornerRadius = 4
	marker := canvas.NewRectangle(parseHex(b.config.Theme.Needle, color.NRGBA{R: 240, G: 244, B: 255, A: 255}))
	marker.CornerRadius = 2
	label := canvas.NewText(b.config.Label, parseHex(b.config.Theme.Label, color.NRGBA{R: 199, G: 208, B: 221, A: 255}))
	label.TextStyle = fyne.TextStyle{Bold: true}
	value := canvas.NewText("", parseHex(b.config.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255}))
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	value.Alignment = fyne.TextAlignTrailing
	r := &bar1Renderer{b: b, bg: bg, pulse: pulse, track: track, fill: fill, marker: marker, label: label, value: value}
	r.Refresh()
	return r
}

type bar1Renderer struct {
	b *Bar1

	bg     *canvas.Rectangle
	pulse  *canvas.Rectangle
	track  *canvas.Rectangle
	fill   *canvas.Rectangle
	marker *canvas.Rectangle
	label  *canvas.Text
	value  *canvas.Text
}

func (r *bar1Renderer) Layout(size fyne.Size) {
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
		r.pulse.FillColor = withAlpha(col, uint8(65*p))
		r.pulse.Show()
	} else {
		r.pulse.Hide()
	}

	pad := float32(14)
	top := float32(10)
	barY := size.Height * 0.55
	barH := float32(math.Max(18, float64(size.Height)*0.24))
	barW := size.Width - pad*2

	r.label.Text = cfg.Label
	r.label.TextSize = float32(math.Max(12, float64(size.Height)*0.09))
	r.label.Refresh()
	r.label.Move(fyne.NewPos(pad, top))

	r.value.Text = fmt.Sprintf("%s %s", formatValue(r.b.value, span), cfg.Unit)
	r.value.TextSize = r.label.TextSize
	r.value.Refresh()
	r.value.Resize(fyne.NewSize(size.Width-pad*2, r.value.MinSize().Height))
	r.value.Move(fyne.NewPos(pad, top))

	r.track.Resize(fyne.NewSize(barW, barH))
	r.track.Move(fyne.NewPos(pad, barY))

	fillW := float32(pct) * barW
	if fillW < 1 && pct > 0 {
		fillW = 1
	}
	r.fill.Resize(fyne.NewSize(fillW, barH))
	r.fill.Move(fyne.NewPos(pad, barY))

	midValue := cfg.Min + pct*span
	if inRange(midValue, cfg.DangerRange) {
		r.fill.FillColor = parseHex(cfg.Theme.Danger, color.NRGBA{R: 255, G: 48, B: 48, A: 255})
	} else if inRange(midValue, cfg.WarningRange) {
		r.fill.FillColor = parseHex(cfg.Theme.Warning, color.NRGBA{R: 255, G: 176, B: 0, A: 255})
	} else {
		r.fill.FillColor = parseHex(cfg.Theme.Value, color.NRGBA{R: 125, G: 249, B: 255, A: 255})
	}

	mx := pad + float32(pct)*barW
	r.marker.Resize(fyne.NewSize(4, barH+10))
	r.marker.Move(fyne.NewPos(mx-2, barY-5))
}

func (r *bar1Renderer) MinSize() fyne.Size { return fyne.NewSize(420, 120) }

func (r *bar1Renderer) Refresh() {
	r.Layout(r.b.Size())
	canvas.Refresh(r.bg)
	canvas.Refresh(r.pulse)
	canvas.Refresh(r.track)
	canvas.Refresh(r.fill)
	canvas.Refresh(r.marker)
	canvas.Refresh(r.label)
	canvas.Refresh(r.value)
}

func (r *bar1Renderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.pulse, r.label, r.value, r.track, r.fill, r.marker}
}

func (r *bar1Renderer) Destroy() {}
