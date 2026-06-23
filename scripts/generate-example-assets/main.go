package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"path/filepath"
	"strings"

	"github.com/MickMake/GoDriveLog/internal/assets/examplegen"
)

const frameworkSmokeTheme = "framework-smoke"
const ornateTimberTheme = "ornate-timber"

func main() {
	var (
		repoRoot = flag.String("repo-root", ".", "repository root")
		theme    = flag.String("theme", frameworkSmokeTheme, "theme to generate: framework-smoke, ornate-timber, or all")
	)
	flag.Parse()

	if err := run(*repoRoot, *theme); err != nil {
		log.Fatal(err)
	}
}

func run(repoRoot, theme string) error {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return fmt.Errorf("resolve repo root: %w", err)
	}

	switch strings.TrimSpace(theme) {
	case "", frameworkSmokeTheme:
		return generateFrameworkSmoke(root)
	case ornateTimberTheme:
		return generateOrnateTimber(root)
	case "all":
		if err := generateFrameworkSmoke(root); err != nil {
			return err
		}
		return generateOrnateTimber(root)
	default:
		return fmt.Errorf("unsupported theme %q; supported themes: %s, %s, all", theme, frameworkSmokeTheme, ornateTimberTheme)
	}
}

func generateFrameworkSmoke(repoRoot string) error {
	themeRoot := filepath.Join(repoRoot, "examples", "assets", "v3.4", frameworkSmokeTheme)

	if err := generatePanel(themeRoot); err != nil {
		return err
	}
	if err := generateDigits(themeRoot); err != nil {
		return err
	}
	if err := generateIndicator(themeRoot); err != nil {
		return err
	}

	fmt.Printf("generated %s assets under %s\n", frameworkSmokeTheme, themeRoot)
	return nil
}

func generatePanel(themeRoot string) error {
	background := examplegen.NewCanvas(1280, 720, color.NRGBA{R: 18, G: 24, B: 33, A: 255})
	background.DrawGrid(40, 1, color.NRGBA{R: 54, G: 74, B: 96, A: 22})
	background.FillRect(48, 48, 1184, 624, color.NRGBA{R: 25, G: 34, B: 46, A: 255})
	background.StrokeRect(48, 48, 1184, 624, 4, color.NRGBA{R: 91, G: 122, B: 153, A: 255})
	background.FillRect(72, 72, 1136, 84, color.NRGBA{R: 32, G: 44, B: 58, A: 255})
	background.FillRect(72, 182, 420, 204, color.NRGBA{R: 12, G: 17, B: 24, A: 255})
	background.FillRect(72, 412, 520, 204, color.NRGBA{R: 12, G: 17, B: 24, A: 255})
	background.FillRect(880, 182, 248, 248, color.NRGBA{R: 12, G: 17, B: 24, A: 255})
	background.FillRect(960, 488, 168, 112, color.NRGBA{R: 26, G: 16, B: 16, A: 255})
	background.StrokeRect(72, 182, 420, 204, 2, color.NRGBA{R: 118, G: 156, B: 193, A: 140})
	background.StrokeRect(72, 412, 520, 204, 2, color.NRGBA{R: 118, G: 156, B: 193, A: 140})
	background.StrokeRect(880, 182, 248, 248, 2, color.NRGBA{R: 118, G: 156, B: 193, A: 140})
	background.StrokeRect(960, 488, 168, 112, 2, color.NRGBA{R: 194, G: 118, B: 118, A: 140})
	background.FillCircle(1004, 306, 96, color.NRGBA{R: 18, G: 22, B: 30, A: 255})
	background.StrokeCircle(1004, 306, 96, 3, color.NRGBA{R: 118, G: 156, B: 193, A: 160})
	background.StrokeCircle(1004, 306, 72, 2, color.NRGBA{R: 73, G: 96, B: 122, A: 110})
	background.AddGrain(examplegen.HashSeed(frameworkSmokeTheme+":panel"), 10)

	foreground := examplegen.NewCanvas(1280, 720, color.NRGBA{})
	foreground.StrokeRect(58, 58, 1164, 604, 1, color.NRGBA{R: 210, G: 225, B: 242, A: 48})
	for _, point := range [][2]int{{72, 72}, {1208, 72}, {72, 648}, {1208, 648}} {
		foreground.FillCircle(point[0], point[1], 10, color.NRGBA{R: 42, G: 46, B: 54, A: 220})
		foreground.StrokeCircle(point[0], point[1], 10, 2, color.NRGBA{R: 180, G: 194, B: 210, A: 160})
	}
	foreground.FillRect(72, 72, 1136, 18, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	foreground.FillRect(72, 182, 420, 14, color.NRGBA{R: 255, G: 255, B: 255, A: 12})
	foreground.FillRect(72, 412, 520, 14, color.NRGBA{R: 255, G: 255, B: 255, A: 12})
	foreground.FillRect(880, 182, 248, 14, color.NRGBA{R: 255, G: 255, B: 255, A: 12})
	foreground.FillRect(960, 488, 168, 10, color.NRGBA{R: 255, G: 196, B: 196, A: 18})

	if err := background.WritePNG(filepath.Join(themeRoot, "panel", "background.png")); err != nil {
		return fmt.Errorf("write framework-smoke panel background: %w", err)
	}
	if err := foreground.WritePNG(filepath.Join(themeRoot, "panel", "foreground.png")); err != nil {
		return fmt.Errorf("write framework-smoke panel foreground: %w", err)
	}
	return nil
}

func generateDigits(themeRoot string) error {
	digitsRoot := filepath.Join(themeRoot, "digits")

	background := examplegen.NewCanvas(72, 108, color.NRGBA{R: 10, G: 13, B: 18, A: 255})
	background.FillRect(6, 6, 60, 96, color.NRGBA{R: 6, G: 8, B: 11, A: 255})
	background.StrokeRect(0, 0, 72, 108, 2, color.NRGBA{R: 72, G: 88, B: 105, A: 255})
	background.StrokeRect(6, 6, 60, 96, 1, color.NRGBA{R: 34, G: 42, B: 52, A: 255})
	cellWidth := background.Image.Bounds().Dx()
	cellHeight := background.Image.Bounds().Dy()

	glass := examplegen.NewCanvas(cellWidth, cellHeight, color.NRGBA{})
	glass.FillRect(4, 4, 64, 24, color.NRGBA{R: 255, G: 255, B: 255, A: 18})
	glass.StrokeRect(1, 1, 70, 106, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 26})

	if err := background.WritePNG(filepath.Join(digitsRoot, "digit_back.png")); err != nil {
		return fmt.Errorf("write digit background: %w", err)
	}
	if err := glass.WritePNG(filepath.Join(digitsRoot, "digit_glass.png")); err != nil {
		return fmt.Errorf("write digit glass: %w", err)
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
		canvas := examplegen.NewCanvas(cellWidth, cellHeight, color.NRGBA{})
		drawDigitSegments(canvas, segments, color.NRGBA{R: 255, G: 168, B: 64, A: 255})
		canvas.AddGrain(examplegen.HashSeed(frameworkSmokeTheme+":"+id), 4)
		filename := "digit_" + id + ".png"
		if id == "-" {
			filename = "digit_minus.png"
		}
		if err := canvas.WritePNG(filepath.Join(digitsRoot, filename)); err != nil {
			return fmt.Errorf("write digit %s: %w", id, err)
		}
	}

	decimalPoint := examplegen.NewCanvas(cellWidth, cellHeight, color.NRGBA{})
	decimalPoint.FillCircle(cellWidth-14, cellHeight-22, 8, color.NRGBA{R: 255, G: 168, B: 64, A: 255})
	decimalPoint.AddGrain(examplegen.HashSeed(frameworkSmokeTheme+":dp"), 4)
	if err := decimalPoint.WritePNG(filepath.Join(digitsRoot, "digit_dp.png")); err != nil {
		return fmt.Errorf("write decimal point: %w", err)
	}

	return nil
}

func drawDigitSegments(canvas *examplegen.Canvas, segments []string, fill color.NRGBA) {
	for _, segment := range segments {
		switch segment {
		case "a":
			canvas.FillRect(18, 10, 36, 8, fill)
		case "b":
			canvas.FillRect(50, 18, 8, 30, fill)
		case "c":
			canvas.FillRect(50, 58, 8, 30, fill)
		case "d":
			canvas.FillRect(18, 88, 36, 8, fill)
		case "e":
			canvas.FillRect(14, 58, 8, 30, fill)
		case "f":
			canvas.FillRect(14, 18, 8, 30, fill)
		case "g":
			canvas.FillRect(18, 49, 36, 8, fill)
		}
	}
}

func generateIndicator(themeRoot string) error {
	indicatorRoot := filepath.Join(themeRoot, "indicator")

	background := examplegen.NewCanvas(96, 96, color.NRGBA{})
	background.FillCircle(48, 48, 44, color.NRGBA{R: 14, G: 16, B: 20, A: 255})
	background.FillCircle(48, 48, 34, color.NRGBA{R: 8, G: 9, B: 11, A: 255})
	background.StrokeCircle(48, 48, 44, 3, color.NRGBA{R: 76, G: 89, B: 103, A: 255})

	if err := background.WritePNG(filepath.Join(indicatorRoot, "lamp_back.png")); err != nil {
		return fmt.Errorf("write indicator background: %w", err)
	}

	states := map[string]color.NRGBA{
		"off":     {R: 74, G: 26, B: 24, A: 255},
		"on":      {R: 255, G: 82, B: 50, A: 255},
		"unknown": {R: 255, G: 194, B: 70, A: 255},
	}
	for name, fill := range states {
		canvas := examplegen.NewCanvas(96, 96, color.NRGBA{})
		canvas.FillCircle(48, 48, 28, fill)
		canvas.StrokeCircle(48, 48, 28, 2, color.NRGBA{R: 255, G: 240, B: 220, A: 96})
		if name == "off" {
			canvas.FillCircle(48, 48, 16, color.NRGBA{R: 52, G: 18, B: 18, A: 220})
		}
		canvas.AddGrain(examplegen.HashSeed(frameworkSmokeTheme+":"+name), 5)
		if err := canvas.WritePNG(filepath.Join(indicatorRoot, "lamp_"+name+".png")); err != nil {
			return fmt.Errorf("write indicator state %s: %w", name, err)
		}
	}

	glass := examplegen.NewCanvas(96, 96, color.NRGBA{})
	glass.FillCircle(40, 32, 18, color.NRGBA{R: 255, G: 255, B: 255, A: 22})
	glass.StrokeCircle(48, 48, 44, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 28})
	if err := glass.WritePNG(filepath.Join(indicatorRoot, "lamp_glass.png")); err != nil {
		return fmt.Errorf("write indicator glass: %w", err)
	}

	return nil
}
