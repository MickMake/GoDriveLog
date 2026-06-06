package retro1

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/widgets/model"
)

type SevenSeg struct {
	widget.BaseWidget
	config model.GaugeConfig
	style  string
	level  int
	value  float64
}

func newSevenSeg(cfg model.GaugeConfig, style string, level int) model.Widget {
	cfg = cfg.Normalize()
	s := &SevenSeg{config: cfg, style: style, level: level, value: cfg.Min}
	s.ExtendBaseWidget(s)
	return s
}

func (s *SevenSeg) Style() string { return s.style }
func (s *SevenSeg) Config() model.GaugeConfig { return s.config }
func (s *SevenSeg) Value() float64 { return s.value }
func (s *SevenSeg) SetValue(v float64) { s.value = clamp(v, s.config.Min, s.config.Max); s.Refresh() }
func (s *SevenSeg) Snapshot() model.Snapshot {
	v := s.Value()
	return model.Snapshot{Style: s.Style(), Label: s.config.Label, Unit: s.config.Unit, Min: s.config.Min, Max: s.config.Max, Value: v, Normalised: normalise(v, s.config.Min, s.config.Max), Warning: inRange(v, s.config.WarningRange), Danger: inRange(v, s.config.DangerRange)}
}

func (s *SevenSeg) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(panelBG)
	frame := canvas.NewRectangle(panelFrame)
	face := canvas.NewRectangle(color.NRGBA{R: 14, G: 6, B: 0, A: 255})
	ghost := canvas.NewText("8888", ghostAmber)
	ghost.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	glow1 := canvas.NewText("", color.NRGBA{})
	glow1.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	glow2 := canvas.NewText("", color.NRGBA{})
	glow2.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	value := canvas.NewText("", textAmber)
	value.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	label := canvas.NewText(s.config.Label, labelGreen)
	label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	unit := canvas.NewText(s.config.Unit, labelGreen)
	unit.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	r := &sevenSegRenderer{s: s, bg: bg, frame: frame, face: face, ghost: ghost, glow1: glow1, glow2: glow2, value: value, label: label, unit: unit}
	r.Refresh()
	return r
}

type sevenSegRenderer struct {
	s *SevenSeg
	bg, frame, face *canvas.Rectangle
	ghost, glow1, glow2, value, label, unit *canvas.Text
}

func (r *sevenSegRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	pad := float32(math.Max(10, float64(size.Height)*0.08))
	r.frame.Move(fyne.NewPos(pad/2, pad/2))
	r.frame.Resize(fyne.NewSize(size.Width-pad, size.Height-pad))
	r.face.Move(fyne.NewPos(pad, pad))
	r.face.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))
	cfg := r.s.config.Normalize()
	span := cfg.Max - cfg.Min
	if span <= 0 { span = 1 }
	n := int(math.Round(clamp(r.s.value, 0, 9999)))
	text := fmt.Sprintf("%04d", n)
	digitSize := float32(math.Max(46, float64(size.Height)*0.48))
	x := pad * 1.5
	y := size.Height*0.5 - digitSize*0.45
	lit := textAmber
	if inRange(r.s.value, cfg.WarningRange) { lit = amberOn }
	if inRange(r.s.value, cfg.DangerRange) { lit = redOn }

	r.label.Text = cfg.Label
	r.label.TextSize = float32(math.Max(13, float64(size.Height)*0.09))
	r.label.Refresh()
	r.label.Move(fyne.NewPos(pad*1.4, pad*1.15))
	r.unit.Text = cfg.Unit
	r.unit.TextSize = r.label.TextSize
	r.unit.Refresh()
	r.unit.Move(fyne.NewPos(size.Width-pad*6, pad*1.15))

	r.ghost.Text = "8888"
	r.ghost.TextSize = digitSize
	r.ghost.Color = color.NRGBA{}
	if r.s.level >= 2 { r.ghost.Color = ghostAmber }
	r.ghost.Refresh()
	r.ghost.Move(fyne.NewPos(x, y))

	for _, g := range []*canvas.Text{r.glow1, r.glow2} {
		g.Text = text
		g.TextSize = digitSize
		g.Color = color.NRGBA{}
		g.Refresh()
	}
	if r.s.level >= 3 {
		r.glow1.Color = withAlpha(lit, 45)
		r.glow2.Color = withAlpha(lit, 80)
	}
	r.glow1.Move(fyne.NewPos(x-5, y-4))
	r.glow2.Move(fyne.NewPos(x-2, y-2))

	r.value.Text = text
	r.value.TextSize = digitSize
	r.value.Color = lit
	r.value.Refresh()
	r.value.Move(fyne.NewPos(x, y))
	_ = span
}

func (r *sevenSegRenderer) MinSize() fyne.Size { return fyne.NewSize(520, 170) }
func (r *sevenSegRenderer) Refresh() { r.Layout(r.s.Size()); for _, obj := range r.Objects() { canvas.Refresh(obj) } }
func (r *sevenSegRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.bg, r.frame, r.face, r.ghost, r.glow1, r.glow2, r.value, r.label, r.unit} }
func (r *sevenSegRenderer) Destroy() {}
