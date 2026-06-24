package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/MickMake/GoDriveLog/internal/assets/examplegen"
)

var steamPalette = struct {
	iron      color.NRGBA
	darkIron  color.NRGBA
	plate     color.NRGBA
	brass     color.NRGBA
	copper    color.NRGBA
	agedCream color.NRGBA
	soot      color.NRGBA
	amber     color.NRGBA
	lamp      color.NRGBA
	warning   color.NRGBA
	wire      color.NRGBA
	verdigris color.NRGBA
}{
	iron:      color.NRGBA{R: 58, G: 57, B: 55, A: 255},
	darkIron:  color.NRGBA{R: 29, G: 31, B: 32, A: 255},
	plate:     color.NRGBA{R: 94, G: 82, B: 67, A: 255},
	brass:     color.NRGBA{R: 170, G: 134, B: 72, A: 255},
	copper:    color.NRGBA{R: 171, G: 96, B: 62, A: 255},
	agedCream: color.NRGBA{R: 202, G: 186, B: 150, A: 255},
	soot:      color.NRGBA{R: 18, G: 17, B: 16, A: 255},
	amber:     color.NRGBA{R: 255, G: 154, B: 78, A: 255},
	lamp:      color.NRGBA{R: 255, G: 219, B: 140, A: 255},
	warning:   color.NRGBA{R: 218, G: 93, B: 58, A: 255},
	wire:      color.NRGBA{R: 116, G: 92, B: 46, A: 255},
	verdigris: color.NRGBA{R: 88, G: 132, B: 118, A: 255},
}

var steamRuntimeGaugePackages = map[string]string{
	"boiler_pressure_bar": `id: steam_scrap_boiler_pressure_bar
type: bar
sensor: boiler_pressure

size:
  width: 154
  height: 300

layers:
  panel: panel.png
  level: level.png
  glass: glass.png

value_map:
  min: 0
  max: 240
  clamp: true

bar:
  mode: level
  axis: vertical
  origin: bottom
  bounds: [56, 34, 42, 224]
`,
	"boiler_warning_indicator": `id: steam_scrap_boiler_warning_indicator
type: indicator
sensor: boiler_warning

size:
  width: 112
  height: 112

layers:
  bezel: bezel.png
  face: face.png
  off: off.png
  on: on.png
  glass: glass.png
`,
	"radial_rpm": `id: steam_scrap_radial_rpm
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
	"rpm_segmented": `id: steam_scrap_rpm_segmented
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
	"speed_numeric": `id: steam_scrap_speed_numeric
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
	"trip_odometer": `id: steam_scrap_trip_odometer
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

func generateSteamScrap(repoRoot string) error {
	exampleThemeRoot := filepath.Join(repoRoot, "examples", steamScrapTheme, "assets")
	exampleGaugeRoot := filepath.Join(exampleThemeRoot, "gauges")

	if err := generateSteamScrapPanel(exampleThemeRoot); err != nil {
		return err
	}

	for _, generate := range []func(string) error{
		generateSteamScrapNumeric,
		generateSteamScrapRadial,
		generateSteamScrapOdometer,
		generateSteamScrapIndicator,
		generateSteamScrapBar,
		generateSteamScrapSegmented,
	} {
		if err := generate(exampleGaugeRoot); err != nil {
			return err
		}
	}
	if err := writeSteamScrapRuntimeGaugeYAML(exampleGaugeRoot); err != nil {
		return err
	}

	fmt.Printf("generated %s assets under %s\n", steamScrapTheme, exampleThemeRoot)
	return nil
}

func writeSteamScrapRuntimeGaugeYAML(runtimeGaugeRoot string) error {
	for packageDir, yaml := range steamRuntimeGaugePackages {
		path := filepath.Join(runtimeGaugeRoot, packageDir, "gauge.yaml")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create steam runtime gauge package dir %s: %w", packageDir, err)
		}
		if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
			return fmt.Errorf("write steam runtime gauge package %s: %w", packageDir, err)
		}
	}
	return nil
}

func generateSteamScrapPanel(themeRoot string) error {
	const width = 1280
	const height = 720

	background := examplegen.NewCanvas(width, height, steamPalette.darkIron)
	drawSteamPlate(background, 0, 0, width, height, steamPalette.iron, examplegen.HashSeed(steamScrapTheme+":panel:outer"))
	drawSteamPlate(background, 26, 26, 1228, 668, steamPalette.plate, examplegen.HashSeed(steamScrapTheme+":panel:mid"))
	drawSteamPlate(background, 58, 58, 1164, 604, shiftColor(steamPalette.plate, -10), examplegen.HashSeed(steamScrapTheme+":panel:inner"))
	drawSteamWindow(background, 72, 92, 412, 196)
	drawSteamRoundWindow(background, 466, 82, 374, 374)
	drawSteamWindow(background, 892, 104, 304, 172)
	drawSteamWindow(background, 1060, 194, 148, 318)
	drawSteamWindow(background, 94, 526, 420, 126)
	drawSteamWindow(background, 856, 520, 178, 126)

	drawSteamPipe(background, 0, 134, 210, 32, steamPalette.brass)
	drawSteamPipe(background, 0, 566, 246, 28, steamPalette.copper)
	drawSteamPipe(background, 1012, 0, 30, 236, steamPalette.iron)
	drawSteamPipe(background, 1176, 0, 28, 222, steamPalette.brass)
	drawSteamPipe(background, 1198, 346, 82, 28, steamPalette.copper)
	drawSteamPipe(background, 934, 640, 346, 26, steamPalette.iron)
	drawSteamCable(background, 830, 318, 1060, 534, steamPalette.wire)
	drawSteamCable(background, 96, 470, 320, 660, steamPalette.wire)
	drawSteamLamp(background, 1168, 118, 22)
	drawSteamLamp(background, 1188, 154, 16)
	background.AddGrain(examplegen.HashSeed(steamScrapTheme+":panel"), 9)

	foreground := examplegen.NewCanvas(width, height, color.NRGBA{})
	foreground.StrokeRect(48, 48, 1184, 624, 1, color.NRGBA{R: 255, G: 244, B: 214, A: 22})
	for _, rivet := range [][2]int{
		{74, 74}, {1206, 74}, {74, 646}, {1206, 646},
		{286, 110}, {930, 120}, {1110, 194}, {1012, 548},
	} {
		drawSteamRivet(foreground, rivet[0], rivet[1], 8, steamPalette.brass)
	}
	foreground.FillRect(84, 104, 390, 10, color.NRGBA{R: 255, G: 240, B: 198, A: 18})
	foreground.FillRect(904, 112, 280, 10, color.NRGBA{R: 255, G: 238, B: 192, A: 16})
	foreground.FillRect(1080, 202, 114, 12, color.NRGBA{R: 255, G: 236, B: 194, A: 14})
	foreground.FillCircle(602, 148, 42, color.NRGBA{R: 255, G: 255, B: 255, A: 12})

	if err := background.WritePNG(filepath.Join(themeRoot, "panel", "background.png")); err != nil {
		return fmt.Errorf("write steam-scrap panel background: %w", err)
	}
	if err := foreground.WritePNG(filepath.Join(themeRoot, "panel", "foreground.png")); err != nil {
		return fmt.Errorf("write steam-scrap panel foreground: %w", err)
	}
	return nil
}

func generateSteamScrapNumeric(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "speed_numeric")

	panel := examplegen.NewCanvas(380, 160, color.NRGBA{})
	drawSteamPlate(panel, 0, 0, 380, 160, steamPalette.plate, examplegen.HashSeed(steamScrapTheme+":numeric:panel"))
	panel.StrokeRect(0, 0, 380, 160, 2, shiftColor(steamPalette.darkIron, 30))
	panel.StrokeRect(4, 4, 372, 152, 1, color.NRGBA{R: 220, G: 196, B: 146, A: 84})
	panel.FillRect(28, 20, 324, 120, steamPalette.soot)
	panel.StrokeRect(28, 20, 324, 120, 2, steamPalette.brass)
	for _, x := range []int{40, 144, 248} {
		panel.FillRect(x, 30, 92, 96, color.NRGBA{R: 20, G: 18, B: 15, A: 255})
		panel.StrokeRect(x, 30, 92, 96, 1, steamPalette.copper)
	}
	for _, rivet := range [][2]int{{18, 18}, {362, 18}, {18, 142}, {362, 142}} {
		drawSteamRivet(panel, rivet[0], rivet[1], 7, steamPalette.brass)
	}

	glass := examplegen.NewCanvas(380, 160, color.NRGBA{})
	glass.FillRect(30, 22, 320, 24, color.NRGBA{R: 255, G: 246, B: 210, A: 20})
	glass.FillRect(30, 22, 28, 116, color.NRGBA{R: 255, G: 255, B: 255, A: 8})
	glass.StrokeRect(29, 21, 322, 118, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 22})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write steam numeric panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write steam numeric glass: %w", err)
	}

	digitsRoot := filepath.Join(gaugeRoot, "digits")
	background := examplegen.NewCanvas(92, 96, color.NRGBA{})
	background.FillRect(0, 0, 92, 96, color.NRGBA{R: 18, G: 15, B: 12, A: 255})
	background.StrokeRect(0, 0, 92, 96, 1, steamPalette.brass)
	background.FillRect(4, 4, 84, 88, color.NRGBA{R: 10, G: 8, B: 7, A: 255})

	digitGlass := examplegen.NewCanvas(92, 96, color.NRGBA{})
	digitGlass.FillRect(4, 4, 84, 18, color.NRGBA{R: 255, G: 244, B: 220, A: 16})
	digitGlass.StrokeRect(1, 1, 90, 94, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 20})

	if err := background.WritePNG(filepath.Join(digitsRoot, "digit_back.png")); err != nil {
		return fmt.Errorf("write steam numeric digit background: %w", err)
	}
	if err := digitGlass.WritePNG(filepath.Join(digitsRoot, "digit_glass.png")); err != nil {
		return fmt.Errorf("write steam numeric digit glass: %w", err)
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
		drawScaledSegments(canvas, 9, 8, 74, 80, segments, alphaColor(steamPalette.amber, 48))
		drawScaledSegments(canvas, 10, 8, 72, 78, segments, steamPalette.amber)
		drawScaledSegments(canvas, 10, 9, 72, 78, segments, color.NRGBA{R: 255, G: 234, B: 188, A: 120})
		canvas.AddGrain(examplegen.HashSeed(steamScrapTheme+":numeric:"+id), 3)
		filename := "digit_" + id + ".png"
		if id == "-" {
			filename = "digit_minus.png"
		}
		if err := canvas.WritePNG(filepath.Join(digitsRoot, filename)); err != nil {
			return err
		}
	}

	decimalPoint := examplegen.NewCanvas(92, 96, color.NRGBA{})
	decimalPoint.FillCircle(74, 74, 10, alphaColor(steamPalette.amber, 40))
	decimalPoint.FillCircle(74, 74, 7, steamPalette.amber)
	decimalPoint.FillCircle(74, 74, 3, color.NRGBA{R: 255, G: 240, B: 198, A: 220})
	if err := decimalPoint.WritePNG(filepath.Join(digitsRoot, "digit_dp.png")); err != nil {
		return fmt.Errorf("write steam numeric decimal point: %w", err)
	}
	return nil
}

func generateSteamScrapRadial(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "radial_rpm")
	const size = 360
	const center = size / 2

	background := examplegen.NewCanvas(size, size, color.NRGBA{})
	background.FillCircle(center, center, 170, steamPalette.iron)
	background.FillCircle(center, center, 152, shiftColor(steamPalette.plate, -16))
	background.StrokeCircle(center, center, 170, 3, shiftColor(steamPalette.darkIron, 42))
	background.StrokeCircle(center, center, 152, 2, steamPalette.brass)

	face := examplegen.NewCanvas(size, size, color.NRGBA{})
	face.FillCircle(center, center, 136, steamPalette.agedCream)
	face.FillCircle(center, center, 104, color.NRGBA{R: 40, G: 32, B: 24, A: 255})
	face.StrokeCircle(center, center, 136, 2, steamPalette.brass)
	face.StrokeCircle(center, center, 104, 1, color.NRGBA{R: 210, G: 186, B: 138, A: 128})
	face.FillRect(124, 244, 112, 10, steamPalette.copper)

	ticks := examplegen.NewCanvas(size, size, color.NRGBA{})
	drawArcDots(ticks, center, center, 118, -135, 135, 13, steamPalette.darkIron, steamPalette.copper)
	drawArcDots(ticks, center, center, 90, -135, 135, 7, alphaColor(steamPalette.brass, 120), alphaColor(steamPalette.brass, 180))

	needle := examplegen.NewCanvas(size, size, color.NRGBA{})
	needle.FillRect(center-5, 84, 10, 136, steamPalette.copper)
	needle.FillRect(center-11, 116, 22, 18, steamPalette.brass)
	needle.FillRect(center-2, 58, 4, 28, steamPalette.warning)
	needle.FillCircle(center, center, 20, steamPalette.brass)
	needle.FillCircle(center, center, 8, steamPalette.darkIron)
	needle.AddGrain(examplegen.HashSeed(steamScrapTheme+":radial:needle"), 4)

	overlay := examplegen.NewCanvas(size, size, color.NRGBA{})
	overlay.FillCircle(132, 114, 42, color.NRGBA{R: 255, G: 255, B: 255, A: 14})
	overlay.StrokeCircle(center, center, 170, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	for _, rivet := range [][2]int{{82, 82}, {278, 82}, {82, 278}, {278, 278}} {
		drawSteamRivet(overlay, rivet[0], rivet[1], 8, steamPalette.brass)
	}

	for name, canvas := range map[string]*examplegen.Canvas{
		"background.png": background,
		"face.png":       face,
		"ticks.png":      ticks,
		"needle.png":     needle,
		"overlay.png":    overlay,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write steam radial asset %s: %w", name, err)
		}
	}
	return nil
}

func generateSteamScrapOdometer(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "trip_odometer")

	panel := examplegen.NewCanvas(332, 118, color.NRGBA{})
	drawSteamPlate(panel, 0, 0, 332, 118, shiftColor(steamPalette.plate, -6), examplegen.HashSeed(steamScrapTheme+":odometer:panel"))
	panel.StrokeRect(0, 0, 332, 118, 2, shiftColor(steamPalette.darkIron, 32))
	for slot := 0; slot < 5; slot++ {
		x := 24 + slot*58
		panel.FillRect(x, 28, 40, 56, color.NRGBA{R: 18, G: 16, B: 14, A: 255})
		panel.StrokeRect(x, 28, 40, 56, 1, steamPalette.brass)
	}
	for _, rivet := range [][2]int{{18, 18}, {314, 18}, {18, 100}, {314, 100}} {
		drawSteamRivet(panel, rivet[0], rivet[1], 7, steamPalette.brass)
	}
	panel.FillRect(18, 90, 296, 7, steamPalette.copper)

	glass := examplegen.NewCanvas(332, 118, color.NRGBA{})
	glass.FillRect(18, 18, 296, 16, color.NRGBA{R: 255, G: 246, B: 214, A: 18})
	glass.StrokeRect(18, 28, 296, 56, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write steam odometer panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write steam odometer glass: %w", err)
	}

	if err := writeSteamOdometerStrip(filepath.Join(gaugeRoot, "digits.png"), steamPalette.agedCream, color.NRGBA{R: 20, G: 18, B: 16, A: 255}); err != nil {
		return err
	}
	if err := writeSteamOdometerStrip(filepath.Join(gaugeRoot, "tenths.png"), steamPalette.amber, color.NRGBA{R: 20, G: 18, B: 16, A: 255}); err != nil {
		return err
	}
	return nil
}

func generateSteamScrapIndicator(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "boiler_warning_indicator")

	bezel := examplegen.NewCanvas(112, 112, color.NRGBA{})
	drawSteamPlate(bezel, 0, 0, 112, 112, steamPalette.iron, examplegen.HashSeed(steamScrapTheme+":indicator:bezel"))
	bezel.FillCircle(56, 56, 40, shiftColor(steamPalette.brass, -12))
	bezel.FillCircle(56, 56, 34, steamPalette.soot)
	bezel.StrokeCircle(56, 56, 40, 2, steamPalette.brass)
	for _, rivet := range [][2]int{{16, 16}, {96, 16}, {16, 96}, {96, 96}} {
		drawSteamRivet(bezel, rivet[0], rivet[1], 6, steamPalette.copper)
	}

	face := examplegen.NewCanvas(112, 112, color.NRGBA{})
	face.FillCircle(56, 56, 24, color.NRGBA{R: 46, G: 30, B: 22, A: 255})
	face.StrokeCircle(56, 56, 24, 1, steamPalette.copper)

	off := examplegen.NewCanvas(112, 112, color.NRGBA{})
	off.FillCircle(56, 56, 20, color.NRGBA{R: 84, G: 46, B: 32, A: 255})
	off.FillCircle(56, 56, 8, color.NRGBA{R: 102, G: 58, B: 40, A: 160})

	on := examplegen.NewCanvas(112, 112, color.NRGBA{})
	drawGlowCircle(on, 56, 56, 22, alphaColor(steamPalette.warning, 24), 7)
	on.FillCircle(56, 56, 20, steamPalette.warning)
	on.FillCircle(56, 56, 8, color.NRGBA{R: 255, G: 232, B: 204, A: 220})

	glass := examplegen.NewCanvas(112, 112, color.NRGBA{})
	glass.FillCircle(44, 38, 18, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.StrokeCircle(56, 56, 40, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 24})

	for name, canvas := range map[string]*examplegen.Canvas{
		"bezel.png": bezel,
		"face.png":  face,
		"off.png":   off,
		"on.png":    on,
		"glass.png": glass,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write steam indicator asset %s: %w", name, err)
		}
	}
	return nil
}

func generateSteamScrapBar(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "boiler_pressure_bar")
	const width = 154
	const height = 300
	const slotX = 56
	const slotY = 34
	const slotW = 42
	const slotH = 224

	panel := examplegen.NewCanvas(width, height, color.NRGBA{})
	drawSteamPlate(panel, 0, 0, width, height, shiftColor(steamPalette.plate, -4), examplegen.HashSeed(steamScrapTheme+":bar:panel"))
	panel.StrokeRect(0, 0, width, height, 2, shiftColor(steamPalette.darkIron, 32))
	panel.FillRect(slotX, slotY, slotW, slotH, color.NRGBA{R: 20, G: 18, B: 16, A: 255})
	panel.StrokeRect(slotX, slotY, slotW, slotH, 2, steamPalette.brass)
	for offset := 0; offset <= slotH; offset += 44 {
		panel.FillRect(32, slotY+offset, 16, 3, steamPalette.agedCream)
	}
	for _, rivet := range [][2]int{{18, 18}, {136, 18}, {18, 282}, {136, 282}} {
		drawSteamRivet(panel, rivet[0], rivet[1], 7, steamPalette.copper)
	}

	level := examplegen.NewCanvas(width, height, color.NRGBA{})
	fillVerticalGradientRect(level, slotX, slotY, slotW, slotH, steamPalette.verdigris, steamPalette.amber)
	for y := slotY + 8; y < slotY+slotH; y += 24 {
		level.FillRect(slotX+6, y, slotW-12, 3, color.NRGBA{R: 255, G: 242, B: 206, A: 58})
	}

	glass := examplegen.NewCanvas(width, height, color.NRGBA{})
	glass.FillRect(slotX+2, slotY+2, slotW-4, 22, color.NRGBA{R: 255, G: 255, B: 255, A: 16})
	glass.StrokeRect(slotX+1, slotY+1, slotW-2, slotH-2, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})

	for name, canvas := range map[string]*examplegen.Canvas{
		"panel.png": panel,
		"level.png": level,
		"glass.png": glass,
	} {
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, name)); err != nil {
			return fmt.Errorf("write steam bar asset %s: %w", name, err)
		}
	}
	return nil
}

func generateSteamScrapSegmented(gaugesRoot string) error {
	gaugeRoot := filepath.Join(gaugesRoot, "rpm_segmented")
	const width = 240
	const height = 144
	slotXs := []int{30, 80, 130, 180}

	panel := examplegen.NewCanvas(width, height, color.NRGBA{})
	drawSteamPlate(panel, 0, 0, width, height, shiftColor(steamPalette.plate, -8), examplegen.HashSeed(steamScrapTheme+":segmented:panel"))
	panel.StrokeRect(0, 0, width, height, 2, shiftColor(steamPalette.darkIron, 30))
	panel.FillRect(20, 34, 200, 66, color.NRGBA{R: 26, G: 22, B: 18, A: 255})
	panel.StrokeRect(20, 34, 200, 66, 2, steamPalette.brass)
	for _, x := range slotXs {
		panel.FillRect(x, 46, 30, 42, steamPalette.soot)
		panel.StrokeRect(x, 46, 30, 42, 1, steamPalette.copper)
	}
	for _, rivet := range [][2]int{{18, 18}, {222, 18}, {18, 126}, {222, 126}} {
		drawSteamRivet(panel, rivet[0], rivet[1], 7, steamPalette.brass)
	}

	glass := examplegen.NewCanvas(width, height, color.NRGBA{})
	glass.FillRect(22, 36, 196, 14, color.NRGBA{R: 255, G: 248, B: 214, A: 18})
	glass.StrokeRect(22, 36, 196, 62, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 18})

	if err := panel.WritePNG(filepath.Join(gaugeRoot, "panel.png")); err != nil {
		return fmt.Errorf("write steam segmented panel: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(gaugeRoot, "glass.png")); err != nil {
		return fmt.Errorf("write steam segmented glass: %w", err)
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
			fill := steamPalette.copper
			glow := alphaColor(steamPalette.copper, 18)
			if slot == threshold.count-1 {
				fill = steamPalette.amber
				glow = alphaColor(steamPalette.amber, 24)
			}
			if threshold.value == 100 && slot == threshold.count-1 {
				fill = steamPalette.warning
				glow = alphaColor(steamPalette.warning, 30)
			}
			drawGlowRect(canvas, x+2, 48, 26, 38, glow, 4)
			canvas.FillRect(x+2, 48, 26, 38, fill)
			canvas.FillRect(x+6, 52, 18, 8, color.NRGBA{R: 255, G: 240, B: 198, A: 68})
		}
		name := fmt.Sprintf("rpm_%03d.png", threshold.value)
		if err := canvas.WritePNG(filepath.Join(gaugeRoot, "levels", name)); err != nil {
			return fmt.Errorf("write steam segmented level %s: %w", name, err)
		}
	}
	return nil
}

func drawSteamPlate(canvas *examplegen.Canvas, x, y, width, height int, base color.NRGBA, seed uint64) {
	canvas.FillRect(x, y, width, height, base)
	for row := 0; row < height; row++ {
		delta := int(ornateNoise(seed+31, row, width)%15) - 7
		line := shiftColor(base, delta)
		line.A = 20
		canvas.FillRect(x, y+row, width, 1, line)
		if row%17 == 0 {
			scar := shiftColor(base, -24)
			scar.A = 10
			canvas.FillRect(x+8, y+row, maxInt(0, width-16), 1, scar)
		}
	}
	for column := x + 16; column < x+width-16; column += maxInt(28, width/5) {
		oxidize := alphaColor(steamPalette.verdigris, 18)
		canvas.FillCircle(column, y+12+int(ornateNoise(seed+47, column, height)%uint64(maxInt(1, height-24))), 3, oxidize)
	}
}

func drawSteamWindow(canvas *examplegen.Canvas, x, y, width, height int) {
	canvas.FillRect(x, y, width, height, shiftColor(steamPalette.darkIron, 8))
	canvas.FillRect(x+10, y+10, width-20, height-20, steamPalette.soot)
	canvas.StrokeRect(x, y, width, height, 2, steamPalette.brass)
	canvas.StrokeRect(x+8, y+8, width-16, height-16, 1, alphaColor(steamPalette.agedCream, 48))
	for _, rivet := range [][2]int{
		{x + 14, y + 14},
		{x + width - 14, y + 14},
		{x + 14, y + height - 14},
		{x + width - 14, y + height - 14},
	} {
		drawSteamRivet(canvas, rivet[0], rivet[1], 6, steamPalette.copper)
	}
}

func drawSteamRoundWindow(canvas *examplegen.Canvas, x, y, width, height int) {
	cx := x + width/2
	cy := y + height/2
	radius := minInt(width, height)/2 - 4
	canvas.FillRect(x, y, width, height, shiftColor(steamPalette.darkIron, 8))
	canvas.FillCircle(cx, cy, radius, shiftColor(steamPalette.plate, -18))
	canvas.FillCircle(cx, cy, radius-18, steamPalette.soot)
	canvas.StrokeCircle(cx, cy, radius, 3, steamPalette.brass)
	canvas.StrokeCircle(cx, cy, radius-18, 1, alphaColor(steamPalette.agedCream, 52))
	for index := 0; index < 8; index++ {
		angle := float64(index) * 45
		rx := cx + int(cosDegrees(angle)*float64(radius-10))
		ry := cy + int(sinDegrees(angle)*float64(radius-10))
		drawSteamRivet(canvas, rx, ry, 6, steamPalette.copper)
	}
}

func drawSteamPipe(canvas *examplegen.Canvas, x, y, width, height int, fill color.NRGBA) {
	canvas.FillRect(x, y, width, height, fill)
	canvas.FillRect(x, y+height/3, width, maxInt(2, height/4), alphaColor(steamPalette.agedCream, 18))
	canvas.StrokeRect(x, y, width, height, 1, shiftColor(fill, -24))
	if width > height {
		drawSteamRivet(canvas, x+18, y+height/2, maxInt(4, height/3), steamPalette.brass)
		drawSteamRivet(canvas, x+width-18, y+height/2, maxInt(4, height/3), steamPalette.brass)
	} else {
		drawSteamRivet(canvas, x+width/2, y+18, maxInt(4, width/3), steamPalette.brass)
		drawSteamRivet(canvas, x+width/2, y+height-18, maxInt(4, width/3), steamPalette.brass)
	}
}

func drawSteamCable(canvas *examplegen.Canvas, x0, y0, x1, y1 int, fill color.NRGBA) {
	steps := maxInt(absInt(x1-x0), absInt(y1-y0))
	if steps == 0 {
		return
	}
	for step := 0; step <= steps; step++ {
		t := float64(step) / float64(steps)
		x := x0 + int(float64(x1-x0)*t)
		y := y0 + int(float64(y1-y0)*t)
		canvas.FillCircle(x, y, 3, fill)
		canvas.FillCircle(x, y, 1, alphaColor(steamPalette.agedCream, 90))
	}
}

func drawSteamLamp(canvas *examplegen.Canvas, cx, cy, radius int) {
	drawGlowCircle(canvas, cx, cy, radius+2, alphaColor(steamPalette.lamp, 18), 6)
	canvas.FillCircle(cx, cy, radius, steamPalette.lamp)
	canvas.FillCircle(cx, cy, maxInt(3, radius/3), color.NRGBA{R: 255, G: 248, B: 222, A: 220})
	canvas.StrokeCircle(cx, cy, radius, 1, alphaColor(steamPalette.copper, 144))
}

func drawSteamRivet(canvas *examplegen.Canvas, cx, cy, radius int, fill color.NRGBA) {
	canvas.FillCircle(cx, cy, radius, fill)
	canvas.FillCircle(cx-radius/4, cy-radius/4, maxInt(1, radius/3), alphaColor(steamPalette.lamp, 96))
	canvas.StrokeCircle(cx, cy, radius, 1, alphaColor(shiftColor(fill, 34), 160))
}

func writeSteamOdometerStrip(path string, fill, background color.NRGBA) error {
	const cellWidth = 40
	const cellHeight = 56
	canvas := examplegen.NewCanvas(cellWidth, cellHeight*10, color.NRGBA{})
	for digit := 0; digit < 10; digit++ {
		y := digit * cellHeight
		canvas.FillRect(0, y, cellWidth, cellHeight, background)
		canvas.StrokeRect(0, y, cellWidth, cellHeight, 1, alphaColor(steamPalette.brass, 108))
		drawScaledSegments(canvas, 4, y+4, 32, 46, segmentCharacters(fmt.Sprintf("%d", digit)[0]), fill)
		canvas.FillRect(4, y+4, 32, 8, color.NRGBA{R: 255, G: 244, B: 214, A: 14})
	}
	return canvas.WritePNG(path)
}

func cosDegrees(value float64) float64 {
	return cosSinDegrees(value, true)
}

func sinDegrees(value float64) float64 {
	return cosSinDegrees(value, false)
}

func cosSinDegrees(value float64, cosine bool) float64 {
	radians := value * 3.141592653589793 / 180
	if cosine {
		return cosApprox(radians)
	}
	return sinApprox(radians)
}

func sinApprox(x float64) float64 {
	for x > 3.141592653589793 {
		x -= 2 * 3.141592653589793
	}
	for x < -3.141592653589793 {
		x += 2 * 3.141592653589793
	}
	x2 := x * x
	return x * (1 - x2/6 + x2*x2/120)
}

func cosApprox(x float64) float64 {
	for x > 3.141592653589793 {
		x -= 2 * 3.141592653589793
	}
	for x < -3.141592653589793 {
		x += 2 * 3.141592653589793
	}
	x2 := x * x
	return 1 - x2/2 + x2*x2/24
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
