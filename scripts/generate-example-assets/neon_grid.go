package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/MickMake/GoDriveLog/internal/assets/examplegen"
)

var neonPalette = struct {
	space      color.NRGBA
	panel      color.NRGBA
	grid       color.NRGBA
	gridAccent color.NRGBA
	blue       color.NRGBA
	cyan       color.NRGBA
	white      color.NRGBA
	indigo     color.NRGBA
	warning    color.NRGBA
	critical   color.NRGBA
}{
	space:      color.NRGBA{R: 5, G: 10, B: 19, A: 255},
	panel:      color.NRGBA{R: 9, G: 18, B: 33, A: 255},
	grid:       color.NRGBA{R: 24, G: 84, B: 140, A: 48},
	gridAccent: color.NRGBA{R: 40, G: 170, B: 220, A: 90},
	blue:       color.NRGBA{R: 20, G: 150, B: 255, A: 255},
	cyan:       color.NRGBA{R: 74, G: 238, B: 255, A: 255},
	white:      color.NRGBA{R: 234, G: 252, B: 255, A: 255},
	indigo:     color.NRGBA{R: 24, G: 42, B: 100, A: 255},
	warning:    color.NRGBA{R: 255, G: 190, B: 72, A: 255},
	critical:   color.NRGBA{R: 255, G: 88, B: 116, A: 255},
}

var neonRuntimeGaugePackages = map[string]string{
	"check_engine_indicator": `id: neon_grid_check_engine_indicator
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
	"coolant_bar": `id: neon_grid_coolant_bar
type: bar
sensor: coolant_temperature

size:
  width: 144
  height: 300

layers:
  panel: panel.png
  level: level.png
  glass: glass.png

value_map:
  min: 40
  max: 120
  clamp: true

bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [54, 34, 36, 232]
`,
	"radial_rpm": `id: neon_grid_radial_rpm
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
	"rpm_segmented": `id: neon_grid_rpm_segmented
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
	"speed_numeric": `id: neon_grid_speed_numeric
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
	"trip_odometer": `id: neon_grid_trip_odometer
type: odometer
sensor: trip_distance

realism:
  drum_slop: [2, -1, 1, -2, 2]

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

func generateNeonGrid(repoRoot string) error {
	exampleThemeRoot := filepath.Join(repoRoot, "examples", neonGridTheme, "assets")
	exampleGaugeRoot := filepath.Join(exampleThemeRoot, "gauges")

	if err := generateNeonGridPanel(exampleThemeRoot); err != nil {
		return err
	}

	for _, generate := range []func(string) error{
		generateNeonGridNumeric,
		generateNeonGridRadial,
		generateNeonGridOdometer,
		generateNeonGridIndicator,
		generateNeonGridBar,
		generateNeonGridSegmented,
	} {
		if err := generate(exampleGaugeRoot); err != nil {
			return err
		}
	}
	if err := writeNeonGridRuntimeGaugeYAML(exampleGaugeRoot); err != nil {
		return err
	}

	fmt.Printf("generated %s assets under %s\n", neonGridTheme, exampleThemeRoot)
	return nil
}

func writeNeonGridRuntimeGaugeYAML(runtimeGaugeRoot string) error {
	for packageDir, yaml := range neonRuntimeGaugePackages {
		path := filepath.Join(runtimeGaugeRoot, packageDir, "gauge.yaml")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create neon runtime gauge package dir %s: %w", packageDir, err)
		}
		if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
			return fmt.Errorf("write neon runtime gauge package %s: %w", packageDir, err)
		}
	}
	return nil
}

func generateNeonGridPanel(themeRoot string) error {
	const width = 1280
	const height = 720

	background := examplegen.NewCanvas(width, height, neonPalette.space)
	background.DrawGrid(32, 1, neonPalette.grid)
	background.DrawGrid(128, 2, alphaColor(neonPalette.gridAccent, 34))
	background.FillRect(0, 0, width, 120, alphaColor(neonPalette.indigo, 28))
	background.FillRect(0, height-96, width, 96, alphaColor(neonPalette.indigo, 20))
	drawNeonTrace(background, 54, 94, 312, 20)
	drawNeonTrace(background, 932, 96, 236, 18)
	drawNeonTrace(background, 116, 618, 280, 16)
	drawNeonTrace(background, 884, 612, 160, 14)
	drawNeonInsetRect(background, 74, 92, 420, 192)
	drawNeonInsetCircle(background, 462, 78, 372, 372)
	drawNeonInsetRect(background, 900, 102, 292, 170)
	drawNeonInsetRect(background, 1060, 186, 144, 320)
	drawNeonInsetRect(background, 84, 528, 430, 128)
	drawNeonInsetRect(background, 850, 520, 180, 126)
	background.AddGrain(examplegen.HashSeed(neonGridTheme+":panel"), 6)

	foreground := examplegen.NewCanvas(width, height, color.NRGBA{})
	drawNeonFrame(foreground, 26, 26, 1228, 668, neonPalette.cyan)
	foreground.StrokeRect(48, 48, 1184, 624, 1, alphaColor(neonPalette.white, 28))
	for _, point := range [][2]int{{72, 72}, {1208, 72}, {72, 648}, {1208, 648}, {434, 126}, {878, 126}, {1034, 548}} {
		drawNeonNode(foreground, point[0], point[1], 9, neonPalette.cyan)
	}
	foreground.FillRect(80, 98, 408, 10, alphaColor(neonPalette.white, 18))
	foreground.FillRect(904, 108, 284, 10, alphaColor(neonPalette.white, 18))
	foreground.FillRect(1066, 194, 132, 12, alphaColor(neonPalette.white, 14))
	foreground.FillCircle(594, 136, 38, alphaColor(neonPalette.white, 12))

	if err := background.WritePNG(filepath.Join(themeRoot, "panel", "background.png")); err != nil {
		return fmt.Errorf("write neon-grid panel background: %w", err)
	}
	if err := foreground.WritePNG(filepath.Join(themeRoot, "panel", "foreground.png")); err != nil {
		return fmt.Errorf("write neon-grid panel foreground: %w", err)
	}
	return nil
}

func generateNeonGridNumeric(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "speed_numeric")

	panel := examplegen.NewCanvas(380, 160, color.NRGBA{})
	panel.FillRect(0, 0, 380, 160, neonPalette.panel)
	drawNeonFrame(panel, 0, 0, 380, 160, neonPalette.blue)
	panel.FillRect(28, 20, 324, 120, color.NRGBA{R: 5, G: 11, B: 20, A: 255})
	panel.StrokeRect(28, 20, 324, 120, 2, alphaColor(neonPalette.cyan, 96))
	for _, x := range []int{40, 144, 248} {
		drawGlowRect(panel, x, 30, 92, 96, alphaColor(neonPalette.blue, 28), 8)
		panel.FillRect(x, 30, 92, 96, color.NRGBA{R: 4, G: 9, B: 16, A: 255})
		panel.StrokeRect(x, 30, 92, 96, 1, alphaColor(neonPalette.cyan, 72))
	}
	panel.FillRect(34, 132, 312, 6, alphaColor(neonPalette.cyan, 80))

	glass := examplegen.NewCanvas(380, 160, color.NRGBA{})
	glass.FillRect(30, 22, 320, 24, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.FillRect(30, 22, 24, 116, color.NRGBA{R: 255, G: 255, B: 255, A: 8})
	glass.StrokeRect(29, 21, 322, 118, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 22})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write neon numeric panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write neon numeric glass: %w", err)
	}

	digitsRoot := filepath.Join(gaugeRoot, "digits")
	background := examplegen.NewCanvas(92, 96, color.NRGBA{})
	background.FillRect(0, 0, 92, 96, color.NRGBA{R: 3, G: 8, B: 15, A: 255})
	background.StrokeRect(0, 0, 92, 96, 1, alphaColor(neonPalette.cyan, 74))
	background.FillRect(4, 4, 84, 88, color.NRGBA{R: 6, G: 12, B: 20, A: 255})

	digitGlass := examplegen.NewCanvas(92, 96, color.NRGBA{})
	digitGlass.FillRect(4, 4, 84, 18, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	digitGlass.StrokeRect(1, 1, 90, 94, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 24})

	if err := background.WritePNG(filepath.Join(digitsRoot, "digit_back.png")); err != nil {
		return fmt.Errorf("write neon numeric digit background: %w", err)
	}
	if err := digitGlass.WritePNG(filepath.Join(digitsRoot, "digit_glass.png")); err != nil {
		return fmt.Errorf("write neon numeric digit glass: %w", err)
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
		drawScaledSegments(canvas, 9, 8, 74, 80, segments, alphaColor(neonPalette.blue, 36))
		drawScaledSegments(canvas, 10, 9, 72, 78, segments, alphaColor(neonPalette.cyan, 72))
		drawScaledSegments(canvas, 10, 8, 72, 78, segments, neonPalette.white)
		canvas.AddGrain(examplegen.HashSeed(neonGridTheme+":numeric:"+id), 3)
		filename := "digit_" + id + ".png"
		if id == "-" {
			filename = "digit_minus.png"
		}
		if err := canvas.WritePNG(filepath.Join(digitsRoot, filename)); err != nil {
			return err
		}
	}

	decimalPoint := examplegen.NewCanvas(92, 96, color.NRGBA{})
	decimalPoint.FillCircle(74, 74, 10, alphaColor(neonPalette.blue, 40))
	decimalPoint.FillCircle(74, 74, 8, neonPalette.cyan)
	decimalPoint.FillCircle(74, 74, 3, neonPalette.white)
	if err := decimalPoint.WritePNG(filepath.Join(digitsRoot, "digit_dp.png")); err != nil {
		return fmt.Errorf("write neon numeric decimal point: %w", err)
	}
	return nil
}

func generateNeonGridRadial(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "radial_rpm")
	const size = 360
	const center = size / 2

	background := examplegen.NewCanvas(size, size, color.NRGBA{})
	drawGlowCircle(background, center, center, 170, alphaColor(neonPalette.blue, 16), 20)
	background.FillCircle(center, center, 168, color.NRGBA{R: 7, G: 13, B: 24, A: 255})
	background.FillCircle(center, center, 152, color.NRGBA{R: 5, G: 10, B: 18, A: 255})
	background.StrokeCircle(center, center, 168, 3, alphaColor(neonPalette.cyan, 96))
	background.StrokeCircle(center, center, 152, 2, alphaColor(neonPalette.blue, 88))

	face := examplegen.NewCanvas(size, size, color.NRGBA{})
	drawGlowCircle(face, center, center, 132, alphaColor(neonPalette.cyan, 16), 10)
	face.FillCircle(center, center, 132, color.NRGBA{R: 4, G: 10, B: 18, A: 255})
	face.FillCircle(center, center, 104, color.NRGBA{R: 3, G: 7, B: 12, A: 255})
	face.StrokeCircle(center, center, 132, 2, alphaColor(neonPalette.cyan, 84))
	face.StrokeCircle(center, center, 104, 1, alphaColor(neonPalette.white, 42))
	face.FillRect(128, 246, 104, 10, alphaColor(neonPalette.blue, 72))

	ticks := examplegen.NewCanvas(size, size, color.NRGBA{})
	drawArcDots(ticks, center, center, 120, -135, 135, 13, alphaColor(neonPalette.blue, 72), alphaColor(neonPalette.cyan, 128))
	drawArcDots(ticks, center, center, 120, -135, 135, 13, neonPalette.blue, neonPalette.cyan)
	drawArcDots(ticks, center, center, 94, -135, 135, 7, alphaColor(neonPalette.blue, 48), alphaColor(neonPalette.cyan, 96))

	needle := examplegen.NewCanvas(size, size, color.NRGBA{})
	drawGlowRect(needle, center-5, 86, 10, 134, alphaColor(neonPalette.blue, 24), 5)
	needle.FillRect(center-5, 86, 10, 134, neonPalette.cyan)
	needle.FillRect(center-2, 56, 4, 36, neonPalette.white)
	needle.FillRect(center-12, 120, 24, 14, alphaColor(neonPalette.white, 112))
	drawGlowCircle(needle, center, center, 21, alphaColor(neonPalette.blue, 24), 8)
	needle.FillCircle(center, center, 19, neonPalette.cyan)
	needle.FillCircle(center, center, 8, neonPalette.white)

	overlay := examplegen.NewCanvas(size, size, color.NRGBA{})
	overlay.FillCircle(132, 114, 46, color.NRGBA{R: 255, G: 255, B: 255, A: 14})
	overlay.StrokeCircle(center, center, 168, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 22})
	for _, point := range [][2]int{{82, 82}, {278, 82}, {82, 278}, {278, 278}} {
		drawNeonNode(overlay, point[0], point[1], 8, neonPalette.cyan)
	}

	for name, canvas := range map[string]*examplegen.Canvas{
		"background.png": background,
		"face.png":       face,
		"ticks.png":      ticks,
		"needle.png":     needle,
		"overlay.png":    overlay,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write neon radial asset %s: %w", name, err)
		}
	}
	return nil
}

func generateNeonGridOdometer(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "trip_odometer")

	panel := examplegen.NewCanvas(332, 118, color.NRGBA{})
	panel.FillRect(0, 0, 332, 118, neonPalette.panel)
	drawNeonFrame(panel, 0, 0, 332, 118, neonPalette.blue)
	for slot := 0; slot < 5; slot++ {
		x := 24 + slot*58
		drawGlowRect(panel, x, 28, 40, 56, alphaColor(neonPalette.blue, 16), 4)
		panel.FillRect(x, 28, 40, 56, color.NRGBA{R: 4, G: 10, B: 16, A: 255})
		panel.StrokeRect(x, 28, 40, 56, 1, alphaColor(neonPalette.cyan, 66))
	}
	panel.FillRect(18, 90, 296, 6, alphaColor(neonPalette.cyan, 76))

	glass := examplegen.NewCanvas(332, 118, color.NRGBA{})
	glass.FillRect(18, 18, 296, 16, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.StrokeRect(18, 28, 296, 56, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write neon odometer panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write neon odometer glass: %w", err)
	}

	if err := writeOdometerStrip(filepath.Join(gaugeRoot, "digits.png"), neonPalette.cyan, color.NRGBA{R: 4, G: 10, B: 16, A: 255}); err != nil {
		return err
	}
	if err := writeOdometerStrip(filepath.Join(gaugeRoot, "tenths.png"), neonPalette.warning, color.NRGBA{R: 4, G: 10, B: 16, A: 255}); err != nil {
		return err
	}
	return nil
}

func generateNeonGridIndicator(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "check_engine_indicator")

	bezel := examplegen.NewCanvas(108, 108, color.NRGBA{})
	bezel.FillRect(0, 0, 108, 108, neonPalette.panel)
	drawNeonFrame(bezel, 0, 0, 108, 108, neonPalette.blue)
	bezel.FillCircle(54, 54, 38, color.NRGBA{R: 4, G: 10, B: 18, A: 255})
	bezel.StrokeCircle(54, 54, 38, 2, alphaColor(neonPalette.cyan, 86))

	face := examplegen.NewCanvas(108, 108, color.NRGBA{})
	face.FillCircle(54, 54, 24, color.NRGBA{R: 6, G: 16, B: 28, A: 255})
	face.StrokeCircle(54, 54, 24, 1, alphaColor(neonPalette.cyan, 96))

	off := examplegen.NewCanvas(108, 108, color.NRGBA{})
	drawGlowCircle(off, 54, 54, 18, alphaColor(neonPalette.blue, 18), 5)
	off.FillCircle(54, 54, 18, color.NRGBA{R: 16, G: 42, B: 72, A: 255})
	off.FillCircle(54, 54, 8, color.NRGBA{R: 34, G: 92, B: 126, A: 180})

	on := examplegen.NewCanvas(108, 108, color.NRGBA{})
	drawGlowCircle(on, 54, 54, 20, alphaColor(neonPalette.critical, 24), 8)
	on.FillCircle(54, 54, 18, neonPalette.critical)
	on.FillCircle(54, 54, 8, color.NRGBA{R: 255, G: 226, B: 234, A: 220})

	glass := examplegen.NewCanvas(108, 108, color.NRGBA{})
	glass.FillCircle(44, 38, 16, color.NRGBA{R: 255, G: 255, B: 255, A: 20})
	glass.StrokeCircle(54, 54, 38, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 24})

	for name, canvas := range map[string]*examplegen.Canvas{
		"bezel.png": bezel,
		"face.png":  face,
		"off.png":   off,
		"on.png":    on,
		"glass.png": glass,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write neon indicator asset %s: %w", name, err)
		}
	}
	return nil
}

func generateNeonGridBar(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "coolant_bar")
	const width = 144
	const height = 300
	const slotX = 54
	const slotY = 34
	const slotW = 36
	const slotH = 232

	panel := examplegen.NewCanvas(width, height, color.NRGBA{})
	panel.FillRect(0, 0, width, height, neonPalette.panel)
	drawNeonFrame(panel, 0, 0, width, height, neonPalette.blue)
	drawGlowRect(panel, slotX, slotY, slotW, slotH, alphaColor(neonPalette.blue, 14), 6)
	panel.FillRect(slotX, slotY, slotW, slotH, color.NRGBA{R: 4, G: 10, B: 16, A: 255})
	panel.StrokeRect(slotX, slotY, slotW, slotH, 2, alphaColor(neonPalette.cyan, 84))
	for offset := 0; offset <= slotH; offset += 46 {
		panel.FillRect(28, slotY+offset, 14, 3, alphaColor(neonPalette.white, 62))
	}

	level := examplegen.NewCanvas(width, height, color.NRGBA{})
	fillVerticalGradientRect(level, slotX, slotY, slotW, slotH, color.NRGBA{R: 30, G: 110, B: 224, A: 255}, neonPalette.cyan)
	for y := slotY + 8; y < slotY+slotH; y += 22 {
		level.FillRect(slotX+5, y, slotW-10, 2, color.NRGBA{R: 255, G: 255, B: 255, A: 64})
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
			return fmt.Errorf("write neon bar asset %s: %w", name, err)
		}
	}
	return nil
}

func generateNeonGridSegmented(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "rpm_segmented")
	const width = 240
	const height = 144
	slotXs := []int{30, 80, 130, 180}

	panel := examplegen.NewCanvas(width, height, color.NRGBA{})
	panel.FillRect(0, 0, width, height, neonPalette.panel)
	drawNeonFrame(panel, 0, 0, width, height, neonPalette.blue)
	panel.FillRect(20, 34, 200, 66, color.NRGBA{R: 4, G: 10, B: 16, A: 255})
	panel.StrokeRect(20, 34, 200, 66, 2, alphaColor(neonPalette.cyan, 92))
	for _, x := range slotXs {
		drawGlowRect(panel, x, 46, 30, 42, alphaColor(neonPalette.blue, 16), 4)
		panel.FillRect(x, 46, 30, 42, color.NRGBA{R: 8, G: 14, B: 22, A: 255})
		panel.StrokeRect(x, 46, 30, 42, 1, alphaColor(neonPalette.cyan, 72))
	}
	panel.FillRect(30, 108, 180, 6, alphaColor(neonPalette.cyan, 74))

	glass := examplegen.NewCanvas(width, height, color.NRGBA{})
	glass.FillRect(22, 36, 196, 14, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.StrokeRect(22, 36, 196, 62, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write neon segmented panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write neon segmented glass: %w", err)
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
			fill := neonPalette.blue
			glow := alphaColor(neonPalette.blue, 24)
			if slot == threshold.count-1 {
				fill = neonPalette.cyan
				glow = alphaColor(neonPalette.cyan, 28)
			}
			if threshold.value == 100 && slot == threshold.count-1 {
				fill = neonPalette.white
				glow = alphaColor(neonPalette.cyan, 34)
			}
			drawGlowRect(canvas, x+2, 48, 26, 38, glow, 4)
			canvas.FillRect(x+2, 48, 26, 38, fill)
			canvas.FillRect(x+6, 52, 18, 8, color.NRGBA{R: 255, G: 255, B: 255, A: 64})
		}
		name := fmt.Sprintf("rpm_%03d.png", threshold.value)
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, "levels", name)); err != nil {
			return fmt.Errorf("write neon segmented level %s: %w", name, err)
		}
	}
	return nil
}

func drawNeonFrame(canvas *examplegen.Canvas, x, y, width, height int, stroke color.NRGBA) {
	drawGlowRect(canvas, x, y, width, height, alphaColor(stroke, 12), 8)
	canvas.StrokeRect(x, y, width, height, 2, alphaColor(stroke, 120))
	canvas.StrokeRect(x+5, y+5, width-10, height-10, 1, alphaColor(neonPalette.white, 28))
}

func drawNeonInsetRect(canvas *examplegen.Canvas, x, y, width, height int) {
	drawGlowRect(canvas, x, y, width, height, alphaColor(neonPalette.blue, 12), 10)
	canvas.FillRect(x, y, width, height, color.NRGBA{R: 7, G: 13, B: 22, A: 255})
	canvas.FillRect(x+12, y+12, width-24, height-24, color.NRGBA{R: 4, G: 10, B: 16, A: 255})
	canvas.StrokeRect(x, y, width, height, 2, alphaColor(neonPalette.cyan, 82))
	canvas.StrokeRect(x+8, y+8, width-16, height-16, 1, alphaColor(neonPalette.white, 24))
}

func drawNeonInsetCircle(canvas *examplegen.Canvas, x, y, width, height int) {
	cx := x + width/2
	cy := y + height/2
	radius := minInt(width, height)/2 - 4
	drawGlowCircle(canvas, cx, cy, radius, alphaColor(neonPalette.blue, 12), 14)
	canvas.FillRect(x, y, width, height, color.NRGBA{R: 8, G: 14, B: 24, A: 255})
	canvas.FillCircle(cx, cy, radius, color.NRGBA{R: 6, G: 12, B: 20, A: 255})
	canvas.FillCircle(cx, cy, radius-16, color.NRGBA{R: 3, G: 8, B: 14, A: 255})
	canvas.StrokeCircle(cx, cy, radius, 2, alphaColor(neonPalette.cyan, 88))
	canvas.StrokeCircle(cx, cy, radius-16, 1, alphaColor(neonPalette.white, 26))
}

func drawNeonTrace(canvas *examplegen.Canvas, x, y, width, height int) {
	drawGlowRect(canvas, x, y, width, height, alphaColor(neonPalette.blue, 12), 4)
	canvas.FillRect(x, y, width, height, alphaColor(neonPalette.blue, 44))
	canvas.FillRect(x+width/2-4, y-height/2, 8, height*2, alphaColor(neonPalette.cyan, 40))
	drawNeonNode(canvas, x+12, y+height/2, 6, neonPalette.cyan)
	drawNeonNode(canvas, x+width-12, y+height/2, 6, neonPalette.cyan)
}

func drawNeonNode(canvas *examplegen.Canvas, cx, cy, radius int, fill color.NRGBA) {
	drawGlowCircle(canvas, cx, cy, radius+2, alphaColor(fill, 20), 5)
	canvas.FillCircle(cx, cy, radius, fill)
	canvas.FillCircle(cx, cy, maxInt(1, radius/3), neonPalette.white)
}

func drawGlowRect(canvas *examplegen.Canvas, x, y, width, height int, glow color.NRGBA, pad int) {
	for _, layer := range []struct {
		pad   int
		alpha uint8
	}{
		{pad: pad, alpha: glow.A / 3},
		{pad: maxInt(1, pad/2), alpha: glow.A / 2},
		{pad: maxInt(1, pad/4), alpha: glow.A},
	} {
		fill := glow
		fill.A = layer.alpha
		canvas.FillRect(x-layer.pad, y-layer.pad, width+layer.pad*2, height+layer.pad*2, fill)
	}
}

func drawGlowCircle(canvas *examplegen.Canvas, cx, cy, radius int, glow color.NRGBA, spread int) {
	for _, layer := range []struct {
		radius int
		alpha  uint8
	}{
		{radius: radius + spread, alpha: glow.A / 4},
		{radius: radius + maxInt(1, spread/2), alpha: glow.A / 2},
		{radius: radius + maxInt(1, spread/4), alpha: glow.A},
	} {
		fill := glow
		fill.A = layer.alpha
		canvas.FillCircle(cx, cy, layer.radius, fill)
	}
}

func alphaColor(base color.NRGBA, alpha uint8) color.NRGBA {
	base.A = alpha
	return base
}
