package retro1

import (
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/widgets/model"
)

type Ramp struct {
	widget.BaseWidget
	config model.GaugeConfig
	style  string
	level  int
	value  float64
}

func newRamp(cfg model.GaugeConfig, style string, level int) model.Widget {
	cfg = cfg.Normalize()
	r := &Ramp{config: cfg, style: style, level: level, value: cfg.Min}
	r.ExtendBaseWidget(r)
	return r
}

func (r *Ramp) Style() string { return r.style }
func (r *Ramp) Config() model.GaugeConfig { return r.config }
func (r *Ramp) Value() float64 { return r.value }
func (r *Ramp) SetValue(v float64) { r.value = clamp(v, r.config.Min, r.config.Max); r.Refresh() }
func (r *Ramp) Snapshot() model.Snapshot {
	v := r.Value()
	return model.Snapshot{Style: r.Style(), Label: r.config.Label, Unit: r.config.Unit, Min: r.config.Min, Max: r.config.Max, Value: v, Normalised: normalise(v, r.config.Min, r.config.Max), Warning: inRange(v, r.config.WarningRange), Danger: inRange(v, r.config.DangerRange)}
}

func (r *Ramp) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(panelBG)
	frame := canvas.NewRectangle(panelFrame)
	face := canvas.NewRectangle(panelInset)
	climb := canvas.NewLine(greenOn)
	top := canvas.NewLine(greenOn)
	climbGhost := canvas.NewLine(ghostGreen)
	topGhost := canvas.NewLine(ghostGreen)
	value := canvas.NewText("", textAmber)
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	label := canvas.NewText(r.config.Label, labelGreen)
	label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	rr := &rampRenderer{r: r, bg: bg, frame: frame, face: face, climb: climb, top: top, climbGhost: climbGhost, topGhost: topGhost, value: value, label: label}
	rr.Refresh()
	return rr
}

type rampRenderer struct {
	r *Ramp
	bg, frame, face *canvas.Rectangle
	climb, top, climbGhost, topGhost *canvas.Line
	value, label *canvas.Text
}

func (rr *rampRenderer) Layout(size fyne.Size) {
	rr.bg.Resize(size)
	pad := float32(math.Max(10, float64(size.Height)*0.08))
	rr.frame.Move(fyne.NewPos(pad/2, pad/2))
	rr.frame.Resize(fyne.NewSize(size.Width-pad, size.Height-pad))
	rr.face.Move(fyne.NewPos(pad, pad))
	rr.face.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))
	cfg := rr.r.config.Normalize()
	span := cfg.Max - cfg.Min
	if span <= 0 { span = 1 }
	pct := clamp((rr.r.value-cfg.Min)/span, 0, 1)
	start := fyne.NewPos(pad*1.8, size.Height-pad*3)
	corner := fyne.NewPos(size.Width*0.42, pad*2.2)
	end := fyne.NewPos(size.Width-pad*1.8, pad*2.2)
	stroke := float32(math.Max(22, float64(size.Height)*0.2))
	rr.climbGhost.Position1, rr.climbGhost.Position2 = start, corner
	rr.topGhost.Position1, rr.topGhost.Position2 = corner, end
	rr.climbGhost.StrokeWidth, rr.topGhost.StrokeWidth = stroke, stroke
	if rr.r.level < 2 { rr.climbGhost.StrokeColor, rr.topGhost.StrokeColor = panelInset, panelInset }
	rr.climb.Position1 = start
	rr.climb.Position2 = rampPoint(math.Min(pct*2, 1), start, corner, end)
	rr.top.Position1 = corner
	rr.top.Position2 = rampPoint(pct, start, corner, end)
	rr.climb.StrokeWidth, rr.top.StrokeWidth = stroke, stroke
	col := greenOn
	if inRange(rr.r.value, cfg.WarningRange) { col = amberOn }
	if inRange(rr.r.value, cfg.DangerRange) { col = redOn }
	rr.climb.StrokeColor, rr.top.StrokeColor = col, col
	if pct <= 0.5 { rr.top.StrokeColor = panelInset }
	rr.value.Text = fmt.Sprintf("%s %s", formatValue(rr.r.value, span), cfg.Unit)
	rr.value.TextSize = float32(math.Max(24, float64(size.Height)*0.19))
	rr.value.Refresh()
	rr.value.Resize(fyne.NewSize(size.Width-pad*2, rr.value.MinSize().Height))
	rr.value.Move(fyne.NewPos(pad, size.Height-pad*1.2-rr.value.MinSize().Height))
	rr.label.Text = cfg.Label
	rr.label.TextSize = float32(math.Max(13, float64(size.Height)*0.09))
	rr.label.Refresh()
	rr.label.Move(fyne.NewPos(pad*1.4, pad*1.15))
}

func rampPoint(t float64, start, corner, end fyne.Position) fyne.Position {
	if t <= 0.5 { p := t / 0.5; return fyne.NewPos(start.X+float32(p)*(corner.X-start.X), start.Y+float32(p)*(corner.Y-start.Y)) }
	p := (t - 0.5) / 0.5
	return fyne.NewPos(corner.X+float32(p)*(end.X-corner.X), corner.Y)
}

func (rr *rampRenderer) MinSize() fyne.Size { return fyne.NewSize(520, 170) }
func (rr *rampRenderer) Refresh() { rr.Layout(rr.r.Size()); canvas.Refresh(rr.bg) }
func (rr *rampRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{rr.bg, rr.frame, rr.face, rr.climbGhost, rr.topGhost, rr.climb, rr.top, rr.label, rr.value} }
func (rr *rampRenderer) Destroy() {}

var _ = container.NewWithoutLayout
