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
	face := canvas.NewRectangle(color.NRGBA{R: 5, G: 10, B: 4, A: 255})
	windows := make([]*canvas.Rectangle, 4)
	ghost := make([]*canvas.Rectangle, 28)
	glow := make([]*canvas.Rectangle, 28)
	core := make([]*canvas.Rectangle, 28)
	for i := range windows {
		windows[i] = canvas.NewRectangle(color.NRGBA{R: 0, G: 7, B: 2, A: 255})
		windows[i].CornerRadius = 2
	}
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
	r := &sevenSegRenderer{s: s, bg: bg, frame: frame, face: face, windows: windows, ghost: ghost, glow: glow, core: core, label: label, unit: unit}
	r.Refresh()
	return r
}

type sevenSegRenderer struct {
	s *SevenSeg
	bg, frame, face *canvas.Rectangle
	windows []*canvas.Rectangle
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

	labelSize := float32(math.Max(12, float64(size.Height)*0.075))
	r.label.Text = cfg.Label
	r.label.TextSize = labelSize
	r.label.Refresh()
	r.label.Move(fyne.NewPos(pad*1.35, pad*1.05))
	r.unit.Text = cfg.Unit
	r.unit.TextSize = labelSize
	r.unit.Refresh()
	r.unit.Move(fyne.NewPos(size.Width-pad*5.0, pad*1.05))

	slotW := (size.Width - pad*4.4) / 4
	digitW := slotW * 0.70
	digitH := size.Height - pad*5.1
	if digitH < 70 { digitH = 70 }
	x := pad * 1.95
	y := pad * 2.7
	idx := 0
	for i, ch := range text {
		winPad := pad * 0.18
		r.windows[i].Move(fyne.NewPos(x-winPad, y-winPad))
		r.windows[i].Resize(fyne.NewSize(digitW+winPad*2, digitH+winPad*2))
		idx = r.layoutDigit(idx, x, y, digitW, digitH, sevenMask(ch), lit)
		x += slotW
	}
}

func (r *sevenSegRenderer) layoutDigit(idx int, x, y, w, h float32, mask [7]bool, lit color.NRGBA) int {
	th := float32(math.Max(7, math.Min(float64(w)*0.12, float64(h)*0.13)))
	gap := th * 0.52
	vh := (h - th*3 - gap*2) / 2
	if vh < th*1.5 { vh = th * 1.5 }
	p := [7][4]float32{{x+gap,y,w-gap*2,th},{x+w-th,y+th+gap,th,vh},{x+w-th,y+th+vh+gap*2,th,vh},{x+gap,y+h-th,w-gap*2,th},{x,y+th+vh+gap*2,th,vh},{x,y+th+gap,th,vh},{x+gap,y+h/2-th/2,w-gap*2,th}}
	for i := range p {
		g1 := r.ghost[idx]
		g2 := r.glow[idx]
		g3 := r.core[idx]
		idx++
		setSeg(g1, p[i], color.NRGBA{})
		setSeg(g2, inflateSeg(p[i], th*0.5), color.NRGBA{})
		setSeg(g3, p[i], color.NRGBA{})
		if r.s.level >= 2 { g1.FillColor = color.NRGBA{R: 30, G: 95, B: 22, A: 85} }
		if mask[i] {
			g3.FillColor = lit
			if r.s.level >= 3 { g2.FillColor = withAlpha(lit, 80) }
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
	for _, x := range r.windows { objs = append(objs, x) }
	for _, x := range r.ghost { objs = append(objs, x) }
	for _, x := range r.glow { objs = append(objs, x) }
	for _, x := range r.core { objs = append(objs, x) }
	return append(objs, r.label,r.unit)
}
func (r *sevenSegRenderer) Destroy() {}
