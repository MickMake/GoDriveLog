package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"

	"github.com/MickMake/GoDriveLog/internal/assets/examplegen"
)

var ornatePalette = struct {
	walnut   color.NRGBA
	jarrah   color.NRGBA
	maple    color.NRGBA
	oak      color.NRGBA
	charcoal color.NRGBA
	brass    color.NRGBA
	amber    color.NRGBA
	green    color.NRGBA
	cream    color.NRGBA
}{
	walnut:   color.NRGBA{R: 70, G: 43, B: 28, A: 255},
	jarrah:   color.NRGBA{R: 116, G: 61, B: 38, A: 255},
	maple:    color.NRGBA{R: 196, G: 162, B: 110, A: 255},
	oak:      color.NRGBA{R: 152, G: 118, B: 72, A: 255},
	charcoal: color.NRGBA{R: 28, G: 24, B: 23, A: 255},
	brass:    color.NRGBA{R: 181, G: 149, B: 88, A: 255},
	amber:    color.NRGBA{R: 242, G: 171, B: 74, A: 255},
	green:    color.NRGBA{R: 121, G: 212, B: 149, A: 255},
	cream:    color.NRGBA{R: 223, G: 211, B: 180, A: 255},
}

var ornateRuntimeGaugePackages = map[string]string{
	"check_engine_indicator": `id: ornate_timber_check_engine_indicator
type: indicator
sensor: check_engine

size:
  width: 108
  height: 108

layers:
  bezel: bezel.png
  face: face.png
  off: off.png
  on: on.png
  glass: glass.png
`,
	"fuel_bar": `id: ornate_timber_fuel_bar
type: bar
sensor: fuel_level

size:
  width: 152
  height: 300

layers:
  panel: panel.png
  level: level.png
  glass: glass.png

value_map:
  min: 0
  max: 100
  clamp: true

bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [56, 34, 40, 220]
`,
	"radial_rpm": `id: ornate_timber_radial_rpm
type: radial
sensor: rpm

size:
  width: 360
  height: 360

layers:
  background: background.png
  face: face.png
  ticks: ticks.png
  needle: needle.png
  overlay: overlay.png

pivot:
  face: { x: 0.5, y: 0.5 }
  needle: { x: 0.5, y: 0.5 }

value_map:
  min: 0
  max: 8000
  start_angle: -135
  end_angle: 135
  clamp: true
`,
	"rpm_segmented": `id: ornate_timber_rpm_segmented
type: segmented
sensor: rpm

size:
  width: 240
  height: 144

layers:
  panel: panel.png
  segments: levels/rpm_{percent:03}.png
  glass: glass.png

segmented:
  hysteresis: 25
`,
	"speed_numeric": `id: ornate_timber_speed_numeric
type: numeric
sensor: speed
format: "%03.0f"

size:
  width: 380
  height: 160

layers:
  panel: panel.png
  glass: glass.png

digit_set:
  background: digits/digit_back.png
  characters:
    "0": digits/digit_0.png
    "1": digits/digit_1.png
    "2": digits/digit_2.png
    "3": digits/digit_3.png
    "4": digits/digit_4.png
    "5": digits/digit_5.png
    "6": digits/digit_6.png
    "7": digits/digit_7.png
    "8": digits/digit_8.png
    "9": digits/digit_9.png
    "-": digits/digit_minus.png
  decimal_point: digits/digit_dp.png
  foreground: digits/digit_glass.png
  spacing: 12

digits:
  count: 3
  positions:
    - [40, 30]
    - [144, 30]
    - [248, 30]
`,
	"trip_odometer": `id: ornate_timber_trip_odometer
type: odometer
sensor: trip_distance

size:
  width: 332
  height: 118

layers:
  panel: panel.png
  glass: glass.png

odometer:
  movement: smooth
  wheels:
    - strip: digits.png
      position: [24, 28]
      window: { width: 40, height: 56 }
    - strip: digits.png
      position: [82, 28]
      window: { width: 40, height: 56 }
    - strip: digits.png
      position: [140, 28]
      window: { width: 40, height: 56 }
    - strip: digits.png
      position: [198, 28]
      window: { width: 40, height: 56 }
    - strip: tenths.png
      position: [256, 28]
      window: { width: 40, height: 56 }
      role: sub_unit
`,
}

func generateOrnateTimber(repoRoot string) error {
	exampleThemeRoot := filepath.Join(repoRoot, "examples", "assets", "v3.4", ornateTimberTheme)
	exampleGaugeRoot := filepath.Join(exampleThemeRoot, "gauges")
	runtimeGaugeRoot := filepath.Join(repoRoot, "assets", "gauges", "v3.4", ornateTimberTheme)

	if err := generateOrnateTimberPanel(exampleThemeRoot); err != nil {
		return err
	}

	for _, gaugesRoot := range []string{exampleGaugeRoot, runtimeGaugeRoot} {
		for _, generate := range []func(string) error{
			generateOrnateTimberNumeric,
			generateOrnateTimberRadial,
			generateOrnateTimberOdometer,
			generateOrnateTimberIndicator,
			generateOrnateTimberBar,
			generateOrnateTimberSegmented,
		} {
			if err := generate(gaugesRoot); err != nil {
				return err
			}
		}
	}
	if err := writeOrnateTimberRuntimeGaugeYAML(runtimeGaugeRoot); err != nil {
		return err
	}

	fmt.Printf("generated %s assets under %s and %s\n", ornateTimberTheme, exampleThemeRoot, runtimeGaugeRoot)
	return nil
}

func writeOrnateTimberRuntimeGaugeYAML(runtimeGaugeRoot string) error {
	for packageDir, yaml := range ornateRuntimeGaugePackages {
		path := filepath.Join(runtimeGaugeRoot, packageDir, "gauge.yaml")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create ornate runtime gauge package dir %s: %w", packageDir, err)
		}
		if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
			return fmt.Errorf("write ornate runtime gauge package %s: %w", packageDir, err)
		}
	}
	return nil
}

func generateOrnateTimberPanel(themeRoot string) error {
	const width = 1280
	const height = 720

	background := examplegen.NewCanvas(width, height, ornatePalette.walnut)
	drawWoodRect(background, 0, 0, width, height, ornatePalette.walnut, examplegen.HashSeed(ornateTimberTheme+":panel:outer"))
	drawWoodRect(background, 30, 30, 1220, 660, ornatePalette.jarrah, examplegen.HashSeed(ornateTimberTheme+":panel:mid"))
	drawWoodRect(background, 64, 64, 1152, 592, ornatePalette.oak, examplegen.HashSeed(ornateTimberTheme+":panel:inner"))
	background.StrokeRect(30, 30, 1220, 660, 3, shiftColor(ornatePalette.walnut, -28))
	background.StrokeRect(64, 64, 1152, 592, 2, shiftColor(ornatePalette.maple, 22))

	drawInlayBand(background, 92, 102, 280, 12, ornatePalette.brass)
	drawInlayBand(background, 910, 102, 278, 12, ornatePalette.brass)
	drawInsetRect(background, 84, 96, 404, 176)
	drawInsetCircle(background, 486, 98, 332, 332)
	drawInsetRect(background, 846, 96, 304, 176)
	drawInsetRect(background, 1038, 190, 134, 318)
	drawInsetRect(background, 90, 522, 414, 122)
	drawInsetRect(background, 852, 520, 156, 122)

	foreground := examplegen.NewCanvas(width, height, color.NRGBA{})
	foreground.StrokeRect(52, 52, 1176, 616, 1, color.NRGBA{R: 255, G: 246, B: 222, A: 34})
	for _, screw := range [][2]int{{72, 72}, {1208, 72}, {72, 648}, {1208, 648}, {392, 132}, {890, 132}, {1038, 548}} {
		drawScrewHead(foreground, screw[0], screw[1], 9)
	}
	foreground.FillRect(84, 96, 404, 10, color.NRGBA{R: 255, G: 240, B: 212, A: 18})
	foreground.FillRect(846, 96, 304, 10, color.NRGBA{R: 255, G: 240, B: 212, A: 18})
	foreground.FillCircle(590, 158, 54, color.NRGBA{R: 255, G: 255, B: 255, A: 14})
	foreground.FillRect(1048, 192, 12, 312, color.NRGBA{R: 255, G: 255, B: 255, A: 10})

	if err := background.WritePNG(filepath.Join(themeRoot, "panel", "background.png")); err != nil {
		return fmt.Errorf("write ornate timber panel background: %w", err)
	}
	if err := foreground.WritePNG(filepath.Join(themeRoot, "panel", "foreground.png")); err != nil {
		return fmt.Errorf("write ornate timber panel foreground: %w", err)
	}
	return nil
}

func generateOrnateTimberNumeric(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "speed_numeric")

	panel := examplegen.NewCanvas(380, 160, color.NRGBA{})
	drawWoodRect(panel, 0, 0, 380, 160, ornatePalette.jarrah, examplegen.HashSeed(ornateTimberTheme+":numeric:panel"))
	panel.StrokeRect(0, 0, 380, 160, 3, shiftColor(ornatePalette.walnut, -24))
	panel.StrokeRect(6, 6, 368, 148, 1, shiftColor(ornatePalette.maple, 24))
	panel.FillRect(28, 20, 324, 120, color.NRGBA{R: 18, G: 16, B: 16, A: 255})
	panel.StrokeRect(28, 20, 324, 120, 2, color.NRGBA{R: 122, G: 100, B: 64, A: 255})
	panel.FillRect(40, 30, 92, 96, color.NRGBA{R: 9, G: 11, B: 13, A: 255})
	panel.FillRect(144, 30, 92, 96, color.NRGBA{R: 9, G: 11, B: 13, A: 255})
	panel.FillRect(248, 30, 92, 96, color.NRGBA{R: 9, G: 11, B: 13, A: 255})
	drawInlayBand(panel, 34, 132, 312, 8, ornatePalette.brass)
	drawScrewHead(panel, 18, 18, 7)
	drawScrewHead(panel, 362, 18, 7)
	drawScrewHead(panel, 18, 142, 7)
	drawScrewHead(panel, 362, 142, 7)

	glass := examplegen.NewCanvas(380, 160, color.NRGBA{})
	glass.FillRect(30, 22, 320, 26, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.FillRect(30, 22, 30, 116, color.NRGBA{R: 255, G: 255, B: 255, A: 8})
	glass.StrokeRect(29, 21, 322, 118, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 30})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write ornate numeric panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write ornate numeric glass: %w", err)
	}

	digitsRoot := filepath.Join(gaugeRoot, "digits")
	background := examplegen.NewCanvas(92, 96, color.NRGBA{})
	background.FillRect(0, 0, 92, 96, color.NRGBA{R: 8, G: 9, B: 10, A: 255})
	background.StrokeRect(0, 0, 92, 96, 1, color.NRGBA{R: 122, G: 100, B: 64, A: 255})
	background.FillRect(4, 4, 84, 88, color.NRGBA{R: 12, G: 12, B: 13, A: 255})

	digitGlass := examplegen.NewCanvas(92, 96, color.NRGBA{})
	digitGlass.FillRect(4, 4, 84, 18, color.NRGBA{R: 255, G: 255, B: 255, A: 16})
	digitGlass.StrokeRect(1, 1, 90, 94, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 24})

	if err := background.WritePNG(filepath.Join(digitsRoot, "digit_back.png")); err != nil {
		return fmt.Errorf("write ornate numeric digit background: %w", err)
	}
	if err := digitGlass.WritePNG(filepath.Join(digitsRoot, "digit_glass.png")); err != nil {
		return fmt.Errorf("write ornate numeric digit glass: %w", err)
	}

	characters := map[string][]string{
		"0": {"a", "b", "c", "d", "e", "f"},
		"1": {"b", "c"},
		"2": {"a", "b", "g", "e", "d"},
		"3": {"a", "b", "c", "d", "g"},
		"4": {"f", "g", "b", "c"},
		"5": {"a", "f", "g", "c", "d"},
		"6": {"a", "f", "e", "d", "c", "g"},
		"7": {"a", "b", "c"},
		"8": {"a", "b", "c", "d", "e", "f", "g"},
		"9": {"a", "b", "c", "d", "f", "g"},
		"-": {"g"},
	}
	for id, segments := range characters {
		canvas := examplegen.NewCanvas(92, 96, color.NRGBA{})
		drawScaledSegments(canvas, 10, 8, 72, 78, segments, ornatePalette.amber)
		canvas.AddGrain(examplegen.HashSeed(ornateTimberTheme+":numeric:"+id), 5)
		filename := "digit_" + id + ".png"
		if id == "-" {
			filename = "digit_minus.png"
		}
		if err := canvas.WritePNG(filepath.Join(digitsRoot, filename)); err != nil {
			return err
		}
	}

	decimalPoint := examplegen.NewCanvas(92, 96, color.NRGBA{})
	decimalPoint.FillCircle(74, 74, 9, ornatePalette.green)
	decimalPoint.FillCircle(74, 74, 4, color.NRGBA{R: 234, G: 252, B: 236, A: 172})
	if err := decimalPoint.WritePNG(filepath.Join(digitsRoot, "digit_dp.png")); err != nil {
		return fmt.Errorf("write ornate numeric decimal point: %w", err)
	}
	return nil
}

func generateOrnateTimberRadial(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "radial_rpm")
	const size = 360
	const center = size / 2

	background := examplegen.NewCanvas(size, size, color.NRGBA{})
	background.FillCircle(center, center, 170, ornatePalette.walnut)
	background.FillCircle(center, center, 154, ornatePalette.jarrah)
	background.StrokeCircle(center, center, 170, 3, shiftColor(ornatePalette.walnut, -24))
	background.StrokeCircle(center, center, 154, 2, shiftColor(ornatePalette.maple, 24))

	face := examplegen.NewCanvas(size, size, color.NRGBA{})
	face.FillCircle(center, center, 136, color.NRGBA{R: 30, G: 28, B: 27, A: 255})
	face.FillCircle(center, center, 104, color.NRGBA{R: 18, G: 18, B: 18, A: 255})
	face.StrokeCircle(center, center, 136, 2, color.NRGBA{R: 102, G: 79, B: 52, A: 255})
	face.StrokeCircle(center, center, 104, 1, color.NRGBA{R: 124, G: 104, B: 72, A: 180})
	drawInlayBand(face, 126, 244, 108, 10, ornatePalette.brass)

	ticks := examplegen.NewCanvas(size, size, color.NRGBA{})
	drawArcDots(ticks, center, center, 118, -135, 135, 13, ornatePalette.maple, ornatePalette.brass)
	drawArcDots(ticks, center, center, 90, -135, 135, 7, color.NRGBA{R: 92, G: 72, B: 48, A: 180}, color.NRGBA{R: 92, G: 72, B: 48, A: 180})

	needle := examplegen.NewCanvas(size, size, color.NRGBA{})
	needle.FillRect(center-6, 84, 12, 136, ornatePalette.maple)
	needle.FillRect(center-11, 112, 22, 18, ornatePalette.oak)
	needle.FillRect(center-3, 62, 6, 26, color.NRGBA{R: 244, G: 226, B: 188, A: 255})
	needle.FillCircle(center, center, 20, ornatePalette.brass)
	needle.FillCircle(center, center, 8, shiftColor(ornatePalette.walnut, -30))
	needle.AddGrain(examplegen.HashSeed(ornateTimberTheme+":radial:needle"), 4)

	overlay := examplegen.NewCanvas(size, size, color.NRGBA{})
	overlay.FillCircle(136, 118, 42, color.NRGBA{R: 255, G: 255, B: 255, A: 16})
	overlay.StrokeCircle(center, center, 170, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 22})
	for _, screw := range [][2]int{{84, 84}, {276, 84}, {84, 276}, {276, 276}} {
		drawScrewHead(overlay, screw[0], screw[1], 8)
	}

	for name, canvas := range map[string]*examplegen.Canvas{
		"background.png": background,
		"face.png":       face,
		"ticks.png":      ticks,
		"needle.png":     needle,
		"overlay.png":    overlay,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write ornate radial asset %s: %w", name, err)
		}
	}
	return nil
}

func generateOrnateTimberOdometer(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "trip_odometer")

	panel := examplegen.NewCanvas(332, 118, color.NRGBA{})
	drawWoodRect(panel, 0, 0, 332, 118, ornatePalette.oak, examplegen.HashSeed(ornateTimberTheme+":odometer:panel"))
	panel.StrokeRect(0, 0, 332, 118, 3, shiftColor(ornatePalette.walnut, -24))
	panel.StrokeRect(5, 5, 322, 108, 1, shiftColor(ornatePalette.maple, 24))
	for slot := 0; slot < 5; slot++ {
		x := 24 + slot*58
		panel.FillRect(x, 28, 40, 56, color.NRGBA{R: 10, G: 10, B: 11, A: 255})
		panel.StrokeRect(x, 28, 40, 56, 1, color.NRGBA{R: 122, G: 100, B: 64, A: 255})
	}
	drawInlayBand(panel, 18, 90, 296, 8, ornatePalette.brass)
	drawScrewHead(panel, 18, 18, 7)
	drawScrewHead(panel, 314, 18, 7)
	drawScrewHead(panel, 18, 100, 7)
	drawScrewHead(panel, 314, 100, 7)

	glass := examplegen.NewCanvas(332, 118, color.NRGBA{})
	glass.FillRect(18, 18, 296, 16, color.NRGBA{R: 255, G: 255, B: 255, A: 16})
	glass.StrokeRect(18, 28, 296, 56, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write ornate odometer panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write ornate odometer glass: %w", err)
	}

	if err := writeOdometerStrip(filepath.Join(gaugeRoot, "digits.png"), ornatePalette.amber, ornatePalette.charcoal); err != nil {
		return err
	}
	if err := writeOdometerStrip(filepath.Join(gaugeRoot, "tenths.png"), ornatePalette.green, ornatePalette.charcoal); err != nil {
		return err
	}
	return nil
}

func generateOrnateTimberIndicator(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "check_engine_indicator")

	bezel := examplegen.NewCanvas(108, 108, color.NRGBA{})
	drawWoodRect(bezel, 0, 0, 108, 108, ornatePalette.jarrah, examplegen.HashSeed(ornateTimberTheme+":indicator:bezel"))
	bezel.FillCircle(54, 54, 40, ornatePalette.walnut)
	bezel.FillCircle(54, 54, 34, color.NRGBA{R: 22, G: 20, B: 20, A: 255})
	bezel.StrokeCircle(54, 54, 40, 2, color.NRGBA{R: 120, G: 94, B: 64, A: 255})
	bezel.StrokeRect(0, 0, 108, 108, 2, shiftColor(ornatePalette.walnut, -24))

	face := examplegen.NewCanvas(108, 108, color.NRGBA{})
	face.FillCircle(54, 54, 24, color.NRGBA{R: 42, G: 34, B: 22, A: 255})
	face.StrokeCircle(54, 54, 24, 1, color.NRGBA{R: 168, G: 134, B: 82, A: 220})

	off := examplegen.NewCanvas(108, 108, color.NRGBA{})
	off.FillCircle(54, 54, 20, color.NRGBA{R: 68, G: 24, B: 18, A: 255})
	off.FillCircle(54, 54, 9, color.NRGBA{R: 88, G: 32, B: 24, A: 160})

	on := examplegen.NewCanvas(108, 108, color.NRGBA{})
	on.FillCircle(54, 54, 20, ornatePalette.amber)
	on.FillCircle(54, 54, 9, color.NRGBA{R: 255, G: 238, B: 194, A: 200})

	glass := examplegen.NewCanvas(108, 108, color.NRGBA{})
	glass.FillCircle(44, 38, 16, color.NRGBA{R: 255, G: 255, B: 255, A: 20})
	glass.StrokeCircle(54, 54, 40, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 24})

	for name, canvas := range map[string]*examplegen.Canvas{
		"bezel.png": bezel,
		"face.png":  face,
		"off.png":   off,
		"on.png":    on,
		"glass.png": glass,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write ornate indicator asset %s: %w", name, err)
		}
	}
	return nil
}

func generateOrnateTimberBar(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "fuel_bar")
	const width = 152
	const height = 300
	const slotX = 56
	const slotY = 34
	const slotW = 40
	const slotH = 220

	panel := examplegen.NewCanvas(width, height, color.NRGBA{})
	drawWoodRect(panel, 0, 0, width, height, ornatePalette.oak, examplegen.HashSeed(ornateTimberTheme+":bar:panel"))
	panel.StrokeRect(0, 0, width, height, 3, shiftColor(ornatePalette.walnut, -24))
	panel.FillRect(slotX, slotY, slotW, slotH, color.NRGBA{R: 18, G: 18, B: 18, A: 255})
	panel.StrokeRect(slotX, slotY, slotW, slotH, 2, ornatePalette.brass)
	for offset := 0; offset <= slotH; offset += 44 {
		panel.FillRect(34, slotY+offset, 14, 3, ornatePalette.maple)
	}
	drawScrewHead(panel, 18, 18, 7)
	drawScrewHead(panel, 134, 18, 7)
	drawScrewHead(panel, 18, 282, 7)
	drawScrewHead(panel, 134, 282, 7)

	level := examplegen.NewCanvas(width, height, color.NRGBA{})
	fillVerticalGradientRect(level, slotX, slotY, slotW, slotH, color.NRGBA{R: 92, G: 178, B: 116, A: 255}, color.NRGBA{R: 238, G: 186, B: 82, A: 255})
	for y := slotY + 8; y < slotY+slotH; y += 24 {
		level.FillRect(slotX+6, y, slotW-12, 3, color.NRGBA{R: 255, G: 236, B: 198, A: 56})
	}

	glass := examplegen.NewCanvas(width, height, color.NRGBA{})
	glass.FillRect(slotX+2, slotY+2, slotW-4, 22, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.StrokeRect(slotX+1, slotY+1, slotW-2, slotH-2, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 20})

	for name, canvas := range map[string]*examplegen.Canvas{
		"panel.png": panel,
		"level.png": level,
		"glass.png": glass,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write ornate bar asset %s: %w", name, err)
		}
	}
	return nil
}

func generateOrnateTimberSegmented(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "rpm_segmented")
	const width = 240
	const height = 144
	slotXs := []int{30, 80, 130, 180}

	panel := examplegen.NewCanvas(width, height, color.NRGBA{})
	drawWoodRect(panel, 0, 0, width, height, ornatePalette.jarrah, examplegen.HashSeed(ornateTimberTheme+":segmented:panel"))
	panel.StrokeRect(0, 0, width, height, 3, shiftColor(ornatePalette.walnut, -24))
	panel.FillRect(20, 34, 200, 66, color.NRGBA{R: 18, G: 16, B: 16, A: 255})
	panel.StrokeRect(20, 34, 200, 66, 2, ornatePalette.brass)
	for _, x := range slotXs {
		panel.FillRect(x, 46, 30, 42, color.NRGBA{R: 30, G: 26, B: 24, A: 255})
		panel.StrokeRect(x, 46, 30, 42, 1, color.NRGBA{R: 102, G: 78, B: 50, A: 255})
	}
	drawInlayBand(panel, 30, 108, 180, 8, ornatePalette.maple)
	drawScrewHead(panel, 18, 18, 7)
	drawScrewHead(panel, 222, 18, 7)
	drawScrewHead(panel, 18, 126, 7)
	drawScrewHead(panel, 222, 126, 7)

	glass := examplegen.NewCanvas(width, height, color.NRGBA{})
	glass.FillRect(22, 36, 196, 14, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.StrokeRect(22, 36, 196, 62, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write ornate segmented panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write ornate segmented glass: %w", err)
	}

	thresholds := []struct {
		value int
		count int
	}{
		{value: 0, count: 0},
		{value: 25, count: 1},
		{value: 50, count: 2},
		{value: 75, count: 3},
		{value: 100, count: 4},
	}
	for _, threshold := range thresholds {
		canvas := examplegen.NewCanvas(width, height, color.NRGBA{})
		for slot := 0; slot < threshold.count; slot++ {
			x := slotXs[slot]
			fill := ornatePalette.amber
			if slot == threshold.count-1 && threshold.value == 100 {
				fill = ornatePalette.green
			}
			canvas.FillRect(x+2, 48, 26, 38, fill)
			canvas.FillRect(x+6, 52, 18, 8, color.NRGBA{R: 255, G: 236, B: 202, A: 64})
		}
		name := fmt.Sprintf("rpm_%03d.png", threshold.value)
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, "levels", name)); err != nil {
			return fmt.Errorf("write ornate segmented level %s: %w", name, err)
		}
	}
	return nil
}

func writeOdometerStrip(path string, fill, background color.NRGBA) error {
	const cellWidth = 40
	const cellHeight = 56
	canvas := examplegen.NewCanvas(cellWidth, cellHeight*10, color.NRGBA{})
	for digit := 0; digit < 10; digit++ {
		y := digit * cellHeight
		canvas.FillRect(0, y, cellWidth, cellHeight, background)
		canvas.StrokeRect(0, y, cellWidth, cellHeight, 1, color.NRGBA{R: 88, G: 72, B: 52, A: 255})
		drawScaledSegments(canvas, 4, y+4, 32, 46, segmentCharacters(fmt.Sprintf("%d", digit)[0]), fill)
		canvas.FillRect(4, y+4, 32, 8, color.NRGBA{R: 255, G: 255, B: 255, A: 14})
	}
	return canvas.WritePNG(path)
}

func drawScaledSegments(canvas *examplegen.Canvas, x, y, width, height int, segments []string, fill color.NRGBA) {
	if width <= 0 || height <= 0 {
		return
	}
	hThickness := maxInt(5, height/10)
	vThickness := maxInt(5, width/8)
	marginX := maxInt(6, width/5)
	marginY := maxInt(6, height/10)
	hWidth := width - marginX*2
	upperHeight := height/2 - marginY - hThickness
	lowerTop := y + height/2 + hThickness/2

	for _, segment := range segments {
		switch segment {
		case "a":
			canvas.FillRect(x+marginX, y+marginY, hWidth, hThickness, fill)
		case "b":
			canvas.FillRect(x+width-marginX-vThickness/2, y+marginY+hThickness, vThickness, upperHeight, fill)
		case "c":
			canvas.FillRect(x+width-marginX-vThickness/2, lowerTop, vThickness, upperHeight, fill)
		case "d":
			canvas.FillRect(x+marginX, y+height-marginY-hThickness, hWidth, hThickness, fill)
		case "e":
			canvas.FillRect(x+marginX-vThickness/2, lowerTop, vThickness, upperHeight, fill)
		case "f":
			canvas.FillRect(x+marginX-vThickness/2, y+marginY+hThickness, vThickness, upperHeight, fill)
		case "g":
			canvas.FillRect(x+marginX, y+height/2-hThickness/2, hWidth, hThickness, fill)
		}
	}
}

func segmentCharacters(ch byte) []string {
	switch ch {
	case '0':
		return []string{"a", "b", "c", "d", "e", "f"}
	case '1':
		return []string{"b", "c"}
	case '2':
		return []string{"a", "b", "g", "e", "d"}
	case '3':
		return []string{"a", "b", "c", "d", "g"}
	case '4':
		return []string{"f", "g", "b", "c"}
	case '5':
		return []string{"a", "f", "g", "c", "d"}
	case '6':
		return []string{"a", "f", "e", "d", "c", "g"}
	case '7':
		return []string{"a", "b", "c"}
	case '8':
		return []string{"a", "b", "c", "d", "e", "f", "g"}
	case '9':
		return []string{"a", "b", "c", "d", "f", "g"}
	default:
		return []string{"g"}
	}
}

func drawWoodRect(canvas *examplegen.Canvas, x, y, width, height int, base color.NRGBA, seed uint64) {
	canvas.FillRect(x, y, width, height, base)
	for row := 0; row < height; row++ {
		delta := int(ornateNoise(seed, row, width)%19) - 9
		line := shiftColor(base, delta)
		line.A = 30
		canvas.FillRect(x, y+row, width, 1, line)
		if row%11 == 0 {
			ridge := shiftColor(base, -12)
			ridge.A = 14
			canvas.FillRect(x, y+row, width, 1, ridge)
		}
	}
	for column := x + 18; column < x+width-18; column += maxInt(26, width/5) {
		if ornateNoise(seed+17, column, height)%3 == 0 {
			continue
		}
		knotY := y + 12 + int(ornateNoise(seed+31, column, height)%uint64(maxInt(1, height-24)))
		knotRadius := 4 + int(ornateNoise(seed+53, column, width)%4)
		knot := shiftColor(base, -22)
		knot.A = 18
		canvas.FillCircle(column, knotY, knotRadius, knot)
	}
}

func drawInsetRect(canvas *examplegen.Canvas, x, y, width, height int) {
	canvas.FillRect(x, y, width, height, color.NRGBA{R: 34, G: 28, B: 24, A: 255})
	canvas.FillRect(x+12, y+12, width-24, height-24, color.NRGBA{R: 22, G: 19, B: 18, A: 255})
	canvas.StrokeRect(x, y, width, height, 2, shiftColor(ornatePalette.walnut, -24))
	canvas.StrokeRect(x+8, y+8, width-16, height-16, 1, ornatePalette.brass)
}

func drawInsetCircle(canvas *examplegen.Canvas, x, y, width, height int) {
	cx := x + width/2
	cy := y + height/2
	radius := minInt(width, height)/2 - 4
	canvas.FillRect(x, y, width, height, color.NRGBA{R: 44, G: 36, B: 30, A: 255})
	canvas.FillCircle(cx, cy, radius, color.NRGBA{R: 32, G: 26, B: 24, A: 255})
	canvas.FillCircle(cx, cy, radius-16, color.NRGBA{R: 20, G: 20, B: 20, A: 255})
	canvas.StrokeCircle(cx, cy, radius, 2, ornatePalette.brass)
	canvas.StrokeCircle(cx, cy, radius-16, 1, color.NRGBA{R: 110, G: 90, B: 60, A: 180})
}

func drawScrewHead(canvas *examplegen.Canvas, cx, cy, radius int) {
	canvas.FillCircle(cx, cy, radius, color.NRGBA{R: 86, G: 70, B: 52, A: 255})
	canvas.StrokeCircle(cx, cy, radius, 1, color.NRGBA{R: 214, G: 190, B: 140, A: 150})
	canvas.FillRect(cx-radius/2, cy-1, radius, 2, color.NRGBA{R: 48, G: 42, B: 36, A: 255})
	canvas.FillRect(cx-1, cy-radius/2, 2, radius, color.NRGBA{R: 48, G: 42, B: 36, A: 255})
}

func drawInlayBand(canvas *examplegen.Canvas, x, y, width, height int, fill color.NRGBA) {
	canvas.FillRect(x, y, width, height, fill)
	canvas.StrokeRect(x, y, width, height, 1, shiftColor(fill, -26))
}

func drawArcDots(canvas *examplegen.Canvas, cx, cy, radius int, startAngle, endAngle float64, count int, regular, major color.NRGBA) {
	if count < 2 {
		return
	}
	step := (endAngle - startAngle) / float64(count-1)
	for index := 0; index < count; index++ {
		angle := (startAngle + step*float64(index)) * math.Pi / 180
		x := cx + int(math.Cos(angle)*float64(radius))
		y := cy + int(math.Sin(angle)*float64(radius))
		r := 4
		fill := regular
		if index == 0 || index == count-1 || index == count/2 || index%3 == 0 {
			r = 6
			fill = major
		}
		canvas.FillCircle(x, y, r, fill)
	}
}

func fillVerticalGradientRect(canvas *examplegen.Canvas, x, y, width, height int, top, bottom color.NRGBA) {
	for row := 0; row < height; row++ {
		t := float64(row) / float64(maxInt(1, height-1))
		canvas.FillRect(x, y+row, width, 1, lerpColor(top, bottom, t))
	}
}

func lerpColor(a, b color.NRGBA, t float64) color.NRGBA {
	channel := func(start, end uint8) uint8 {
		return uint8(float64(start) + (float64(end)-float64(start))*t)
	}
	return color.NRGBA{
		R: channel(a.R, b.R),
		G: channel(a.G, b.G),
		B: channel(a.B, b.B),
		A: channel(a.A, b.A),
	}
}

func shiftColor(base color.NRGBA, delta int) color.NRGBA {
	return color.NRGBA{
		R: shiftChannel(base.R, delta),
		G: shiftChannel(base.G, delta),
		B: shiftChannel(base.B, delta),
		A: base.A,
	}
}

func shiftChannel(channel uint8, delta int) uint8 {
	value := int(channel) + delta
	switch {
	case value < 0:
		return 0
	case value > 255:
		return 255
	default:
		return uint8(value)
	}
}

func ornateNoise(seed uint64, x, y int) uint64 {
	value := seed ^ (uint64(uint32(x)) * 0x9e3779b97f4a7c15) ^ (uint64(uint32(y)) * 0xc2b2ae3d27d4eb4f)
	value ^= value >> 30
	value *= 0xbf58476d1ce4e5b9
	value ^= value >> 27
	value *= 0x94d049bb133111eb
	value ^= value >> 31
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
