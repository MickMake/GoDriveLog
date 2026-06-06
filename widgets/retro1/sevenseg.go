package retro1

import (
	"embed"
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/MickMake/GoDriveLog/widgets/model"
)

// Expected asset layout, relative to this file:
//
//	widgets/retro1/assets/7seg/green/digit_0_on.png
//	widgets/retro1/assets/7seg/green/digit_1_on.png
//	...
//	widgets/retro1/assets/7seg/green/digit_9_on.png
//	widgets/retro1/assets/7seg/green/digit_off.png
//	widgets/retro1/assets/7seg/green/dot_on.png
//	widgets/retro1/assets/7seg/green/dot_off.png
//
// Same structure may exist for yellow and red.
//
//go:embed assets/7seg/*/*.png
var sevenSegAssets embed.FS

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

func (s *SevenSeg) SetValue(v float64) {
	s.value = clamp(v, s.config.Min, s.config.Max)
	s.Refresh()
}

func (s *SevenSeg) Snapshot() model.Snapshot {
	v := s.Value()
	return model.Snapshot{
		Style:      s.Style(),
		Label:      s.config.Label,
		Unit:       s.config.Unit,
		Min:        s.config.Min,
		Max:        s.config.Max,
		Value:      v,
		Normalised: normalise(v, s.config.Min, s.config.Max),
		Warning:    inRange(v, s.config.WarningRange),
		Danger:     inRange(v, s.config.DangerRange),
	}
}

func (s *SevenSeg) CreateRenderer() fyne.WidgetRenderer {
	if s.level >= 3 {
		return newSevenSegSpriteRenderer(s)
	}
	return newSevenSegVectorRenderer(s)
}

// -----------------------------------------------------------------------------
// Sprite renderer for retro1_7seg3.
// -----------------------------------------------------------------------------

type sevenSegSpriteRenderer struct {
	s *SevenSeg

	bg    *canvas.Rectangle
	frame *canvas.Rectangle
	face  *canvas.Rectangle
	label *canvas.Text
	unit  *canvas.Text

	offDigits [4]*canvas.Image
	onDigits  [4]*canvas.Image
	dotOff    [4]*canvas.Image
	dotOn     [4]*canvas.Image

	objects []fyne.CanvasObject
}

func newSevenSegSpriteRenderer(s *SevenSeg) fyne.WidgetRenderer {
	r := &sevenSegSpriteRenderer{
		s:     s,
		bg:    canvas.NewRectangle(panelBG),
		frame: canvas.NewRectangle(panelFrame),
		face:  canvas.NewRectangle(color.NRGBA{R: 2, G: 4, B: 2, A: 255}),
		label: canvas.NewText(s.config.Label, labelGreen),
		unit:  canvas.NewText(s.config.Unit, labelGreen),
	}
	r.frame.CornerRadius = 5
	r.face.CornerRadius = 3
	r.label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	r.unit.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	blank := blankResource()
	r.objects = []fyne.CanvasObject{r.bg, r.frame, r.face}
	for i := 0; i < 4; i++ {
		r.offDigits[i] = canvas.NewImageFromResource(blank)
		r.onDigits[i] = canvas.NewImageFromResource(blank)
		r.dotOff[i] = canvas.NewImageFromResource(blank)
		r.dotOn[i] = canvas.NewImageFromResource(blank)

		for _, img := range []*canvas.Image{r.offDigits[i], r.onDigits[i], r.dotOff[i], r.dotOn[i]} {
			img.FillMode = canvas.ImageFillContain
			img.ScaleMode = canvas.ImageScaleSmooth
		}

		r.objects = append(r.objects, r.offDigits[i], r.onDigits[i], r.dotOff[i], r.dotOn[i])
	}
	r.objects = append(r.objects, r.label, r.unit)

	r.Refresh()
	return r
}

func (r *sevenSegSpriteRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	pad := float32(math.Max(10, float64(size.Height)*0.08))
	r.frame.Move(fyne.NewPos(pad/2, pad/2))
	r.frame.Resize(fyne.NewSize(size.Width-pad, size.Height-pad))

	r.face.Move(fyne.NewPos(pad, pad))
	r.face.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))

	cfg := r.s.config.Normalize()
	value := int(math.Round(clamp(r.s.value, 0, 9999)))
	text := fmt.Sprintf("%04d", value)
	colourName := spriteColourForValue(r.s.value, cfg)

	labelSize := float32(math.Max(12, float64(size.Height)*0.075))
	r.label.Text = cfg.Label
	r.label.TextSize = labelSize
	r.label.Refresh()
	r.label.Move(fyne.NewPos(pad*1.25, pad*0.95))

	r.unit.Text = cfg.Unit
	r.unit.TextSize = labelSize
	r.unit.Refresh()
	r.unit.Move(fyne.NewPos(size.Width-pad*5.0, pad*0.95))

	digitAreaTop := pad * 2.35
	digitAreaBottom := size.Height - pad*1.35
	digitH := digitAreaBottom - digitAreaTop
	slotW := (size.Width - pad*4.2) / 4
	digitW := slotW * 0.82
	x := pad * 1.6

	for i, ch := range text {
		digit := ch
		if digit < '0' || digit > '9' {
			digit = '0'
		}

		offPath := fmt.Sprintf("assets/7seg/%s/digit_off.png", colourName)
		onPath := fmt.Sprintf("assets/7seg/%s/digit_%c_on.png", colourName, digit)
		dotOffPath := fmt.Sprintf("assets/7seg/%s/dot_off.png", colourName)
		dotOnPath := fmt.Sprintf("assets/7seg/%s/dot_on.png", colourName)

		setImageResource(r.offDigits[i], loadSevenSegResource(offPath))
		setImageResource(r.onDigits[i], loadSevenSegResource(onPath))
		setImageResource(r.dotOff[i], loadSevenSegResource(dotOffPath))
		setImageResource(r.dotOn[i], loadSevenSegResource(dotOnPath))

		pos := fyne.NewPos(x, digitAreaTop)
		sz := fyne.NewSize(digitW, digitH)

		r.offDigits[i].Move(pos)
		r.offDigits[i].Resize(sz)
		r.onDigits[i].Move(pos)
		r.onDigits[i].Resize(sz)

		// Dot images are full digit-canvas placeholders. Draw them at the same
		// position as the digit. If they are transparent placeholders, no harm done.
		r.dotOff[i].Move(pos)
		r.dotOff[i].Resize(sz)
		r.dotOn[i].Move(pos)
		r.dotOn[i].Resize(sz)

		// Current gauge value is four integer digits, so decimal dots stay hidden.
		r.dotOn[i].Hide()
		r.dotOff[i].Hide()

		x += slotW
	}
}

func (r *sevenSegSpriteRenderer) MinSize() fyne.Size { return fyne.NewSize(520, 170) }

func (r *sevenSegSpriteRenderer) Refresh() {
	r.Layout(r.s.Size())
	for _, obj := range r.objects {
		canvas.Refresh(obj)
	}
}

func (r *sevenSegSpriteRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *sevenSegSpriteRenderer) Destroy() {}

func spriteColourForValue(value float64, cfg model.GaugeConfig) string {
	if inRange(value, cfg.DangerRange) {
		return "red"
	}
	if inRange(value, cfg.WarningRange) {
		return "yellow"
	}
	return "green"
}

func loadSevenSegResource(path string) fyne.Resource {
	data, err := sevenSegAssets.ReadFile(path)
	if err != nil {
		return blankResource()
	}
	return fyne.NewStaticResource(path, data)
}

func setImageResource(img *canvas.Image, res fyne.Resource) {
	img.Resource = res
	img.Refresh()
}

func blankResource() fyne.Resource {
	return fyne.NewStaticResource("blank.png", blankPNG)
}

// 1x1 transparent PNG.
var blankPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

// -----------------------------------------------------------------------------
// Simple vector fallback for retro1_7seg1 and retro1_7seg2.
// -----------------------------------------------------------------------------

type sevenSegVectorRenderer struct {
	s *SevenSeg

	bg    *canvas.Rectangle
	frame *canvas.Rectangle
	face  *canvas.Rectangle
	label *canvas.Text
	unit  *canvas.Text

	ghost []*canvas.Rectangle
	glow  []*canvas.Rectangle
	core  []*canvas.Rectangle

	objects []fyne.CanvasObject
}

func newSevenSegVectorRenderer(s *SevenSeg) fyne.WidgetRenderer {
	r := &sevenSegVectorRenderer{
		s:     s,
		bg:    canvas.NewRectangle(panelBG),
		frame: canvas.NewRectangle(panelFrame),
		face:  canvas.NewRectangle(color.NRGBA{R: 8, G: 2, B: 0, A: 255}),
		label: canvas.NewText(s.config.Label, labelGreen),
		unit:  canvas.NewText(s.config.Unit, labelGreen),
		ghost: make([]*canvas.Rectangle, 28),
		glow:  make([]*canvas.Rectangle, 28),
		core:  make([]*canvas.Rectangle, 28),
	}
	r.frame.CornerRadius = 5
	r.face.CornerRadius = 3
	r.label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	r.unit.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	r.objects = []fyne.CanvasObject{r.bg, r.frame, r.face}
	for i := range r.ghost {
		r.ghost[i] = canvas.NewRectangle(color.NRGBA{})
		r.glow[i] = canvas.NewRectangle(color.NRGBA{})
		r.core[i] = canvas.NewRectangle(color.NRGBA{})
		r.ghost[i].CornerRadius = 4
		r.glow[i].CornerRadius = 7
		r.core[i].CornerRadius = 4
		r.objects = append(r.objects, r.ghost[i], r.glow[i], r.core[i])
	}
	r.objects = append(r.objects, r.label, r.unit)

	r.Refresh()
	return r
}

func (r *sevenSegVectorRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	pad := float32(math.Max(10, float64(size.Height)*0.08))
	r.frame.Move(fyne.NewPos(pad/2, pad/2))
	r.frame.Resize(fyne.NewSize(size.Width-pad, size.Height-pad))

	r.face.Move(fyne.NewPos(pad, pad))
	r.face.Resize(fyne.NewSize(size.Width-pad*2, size.Height-pad*2))

	cfg := r.s.config.Normalize()
	value := int(math.Round(clamp(r.s.value, 0, 9999)))
	text := fmt.Sprintf("%04d", value)

	lit := textAmber
	if inRange(r.s.value, cfg.WarningRange) {
		lit = amberOn
	}
	if inRange(r.s.value, cfg.DangerRange) {
		lit = redOn
	}

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
	if digitH < 70 {
		digitH = 70
	}

	x := pad * 1.95
	y := pad * 2.7
	idx := 0

	for _, ch := range text {
		idx = r.layoutDigit(idx, x, y, digitW, digitH, sevenMask(ch), lit)
		x += slotW
	}
}

func (r *sevenSegVectorRenderer) layoutDigit(idx int, x, y, w, h float32, mask [7]bool, lit color.NRGBA) int {
	th := float32(math.Max(7, math.Min(float64(w)*0.12, float64(h)*0.13)))
	gap := th * 0.52
	vh := (h - th*3 - gap*2) / 2
	if vh < th*1.5 {
		vh = th * 1.5
	}

	p := [7][4]float32{
		{x + gap, y, w - gap*2, th},
		{x + w - th, y + th + gap, th, vh},
		{x + w - th, y + th + vh + gap*2, th, vh},
		{x + gap, y + h - th, w - gap*2, th},
		{x, y + th + vh + gap*2, th, vh},
		{x, y + th + gap, th, vh},
		{x + gap, y + h/2 - th/2, w - gap*2, th},
	}

	for i := range p {
		g1 := r.ghost[idx]
		g2 := r.glow[idx]
		g3 := r.core[idx]
		idx++

		setSeg(g1, p[i], color.NRGBA{})
		setSeg(g2, inflateSeg(p[i], th*0.5), color.NRGBA{})
		setSeg(g3, p[i], color.NRGBA{})

		if r.s.level >= 2 {
			g1.FillColor = color.NRGBA{R: 30, G: 95, B: 22, A: 85}
		}
		if mask[i] {
			g3.FillColor = lit
			if r.s.level >= 3 {
				g2.FillColor = withAlpha(lit, 80)
			}
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
	return [4]float32{p[0] - d, p[1] - d, p[2] + d*2, p[3] + d*2}
}

func sevenMask(ch rune) [7]bool {
	switch ch {
	case '0':
		return [7]bool{true, true, true, true, true, true, false}
	case '1':
		return [7]bool{false, true, true, false, false, false, false}
	case '2':
		return [7]bool{true, true, false, true, true, false, true}
	case '3':
		return [7]bool{true, true, true, true, false, false, true}
	case '4':
		return [7]bool{false, true, true, false, false, true, true}
	case '5':
		return [7]bool{true, false, true, true, false, true, true}
	case '6':
		return [7]bool{true, false, true, true, true, true, true}
	case '7':
		return [7]bool{true, true, true, false, false, false, false}
	case '8':
		return [7]bool{true, true, true, true, true, true, true}
	case '9':
		return [7]bool{true, true, true, true, false, true, true}
	default:
		return [7]bool{}
	}
}

func (r *sevenSegVectorRenderer) MinSize() fyne.Size { return fyne.NewSize(520, 170) }

func (r *sevenSegVectorRenderer) Refresh() {
	r.Layout(r.s.Size())
	for _, obj := range r.objects {
		canvas.Refresh(obj)
	}
}

func (r *sevenSegVectorRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *sevenSegVectorRenderer) Destroy() {}
