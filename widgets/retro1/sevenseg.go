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
	ghost := make([]*canvas.Rectangle, 112)
	glow := make([]*canvas.Rectangle, 112)
	core := make([]*canvas.Rectangle, 112)
	for i := range ghost {
		ghost[i] = canvas.NewRectangle(color.NRGBA{})
		glow[i] = canvas.NewRectangle(color.NRGBA{})
		core[i] = canvas.NewRectangle(color.NRGBA{})
		ghost[i].CornerRadius = 4
		glow[i].CornerRadius = 7
		core[i].CornerRadius = 4
	}
	label := canvas.NewText(s.config.Label, labelGreen)
	label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	unit := canvas.NewText(s.config.Unit, labelGreen)
	unit.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	r := &sevenSegRenderer{s: s, bg: bg, frame: frame, face: face, ghost: ghost, glow: glow, core: core, label: label, unit: unit}
	r.Refresh()
	return r
}

type sevenSegRenderer struct {
	s *SevenSeg
	bg, frame, face *canvas.Rectangle
	ghost, glow, core []*canvas.Rectangle
	label, unit *canvas.Text
}

func (r *sevenSegRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	pad := float32(math.Max(10, float64(size.Height)*0.08))
	r.frame.Move(fyne.NewPos(pad/2, pad/2))
	r.frame.Resize(fyne.NewSize(size.Width-pad, size.Height-pad))
	r.face.Move(fyne.NewPos(pad, pad))
	r.face.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))
	cfg := r.s.config.Normalize()
	n := int(math.Round(clamp(r.s.value, 0, 9999)))
	text := fmt.Sprintf("%04d", n)
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
	dw := (size.Width - pad*4.0) / 4
	dh := size.Height - pad*4.3
	x := pad * 1.55
	y := pad * 2.55
	idx := 0
	for _, ch := range text {
		idx = r.layoutDigit(idx, x, y, dw*0.82, dh, sevenMask(ch), lit)
		x += dw
	}
}

func (r *sevenSegRenderer) layoutDigit(idx int, x, y, w, h float32, mask [7]bool, lit color.NRGBA) int {
	th := float32(math.Max(7, float64(w)*0.16))
	gap := th * 0.35
	vh := (h - th*3 - gap*2) / 2
	p := [7][4]float32{{x+gap,y,w-gap*2,th},{x+w-th,y+gap,th,vh},{x+w-th,y+th+vh+gap,th,vh},{x+gap,y+h-th,w-gap*2,th},{x,y+th+vh+gap,th,vh},{x,y+gap,th,vh},{x+gap,y+h/2-th/2,w-gap*2,th}}
	for i := range p {
		g1 := r.ghost[idx]
		g2 := r.glow[idx]
		g3 := r.core[idx]
		idx++
		setSeg(g1, p[i], color.NRGBA{})
		setSeg(g2, inflateSeg(p[i], 7), color.NRGBA{})
		setSeg(g3, p[i], color.NRGBA{})
		if r.s.level >= 2 { g1.FillColor = ghostAmber }
		if mask[i] {
			g3.FillColor = lit
			if r.s.level >= 3 { g2.FillColor = withAlpha(lit, 70) }
		}
	}
	return idx
}

func setSeg(r *canvas.Rectangle, p [4]float32, c color.NRGBA) {
	r.Move(fyne.NewPos(p[0], p[1]))
	r.Resize(fyne.NewSize(p[2], p[3]))
	r.FillColor = c
}

func inflateSeg(p [4]float32, d float32) [4]float32 {
	return [4]float32{p[0]-d, p[1]-d, p[2]+d*2, p[3]+d*2}
}

func sevenMask(ch rune) [7]bool {
	switch ch {
	case '0': return [7]bool{true,true,true,true,true,true,false}
	case '1': return [7]bool{false,true,true,false,false,false,false}
	case '2': return [7]bool{true,true,false,true,true,false,true}
	case '3': return [7]bool{true,true,true,true,false,false,true}
	case '4': return [7]bool{false,true,true,false,false,true,true}
	case '5': return [7]bool{true,false,true,true,false,true,true}
	case '6': return [7]bool{true,false,true,true,true,true,true}
	case '7': return [7]bool{true,true,true,false,false,false,false}
	case '8': return [7]bool{true,true,true,true,true,true,true}
	case '9': return [7]bool{true,true,true,true,false,true,true}
	}
	return [7]bool{}
}

func (r *sevenSegRenderer) MinSize() fyne.Size { return fyne.NewSize(520, 170) }
func (r *sevenSegRenderer) Refresh() { r.Layout(r.s.Size()); for _, obj := range r.Objects() { canvas.Refresh(obj) } }
func (r *sevenSegRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{r.bg,r.frame,r.face}
	for _, x := range r.ghost { objs = append(objs, x) }
	for _, x := range r.glow { objs = append(objs, x) }
	for _, x := range r.core { objs = append(objs, x) }
	return append(objs, r.label,r.unit)
}
func (r *sevenSegRenderer) Destroy() {}
