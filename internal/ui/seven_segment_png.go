package ui

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const defaultSevenSegmentAssetDir = "assets/digits/normal"

type digitCell struct {
	Symbol  string
	Decimal bool
}

type pngDigitResources struct {
	dir      string
	bySymbol map[string]fyne.Resource
	dp       fyne.Resource
}

type PNGSevenSegmentDisplay struct {
	root      *fyne.Container
	resources *pngDigitResources
	cells     []*pngDigitCell
	digitSize fyne.Size
}

type pngDigitCell struct {
	root    *fyne.Container
	base    *canvas.Image
	decimal *canvas.Image
	symbol  string
	dp      bool
}

func NewPNGSevenSegmentDisplay(value string, digits int, x, y, width, height float32) *PNGSevenSegmentDisplay {
	display := &PNGSevenSegmentDisplay{
		resources: loadPNGDigitResources(defaultSevenSegmentAssetDir),
		digitSize: fyne.NewSize(width/float32(maxInt(digits, 1)), height),
	}
	display.root = container.NewWithoutLayout()
	display.root.Move(fyne.NewPos(x, y))
	display.root.Resize(fyne.NewSize(width, height))
	display.SetValue(value, digits)
	return display
}

func (d *PNGSevenSegmentDisplay) CanvasObject() fyne.CanvasObject { return d.root }

func (d *PNGSevenSegmentDisplay) SetValue(value string, digits int) {
	cells := formatDigitCells(value, digits)
	for len(d.cells) < len(cells) {
		cell := newPNGDigitCell(d.resources, d.digitSize)
		d.cells = append(d.cells, cell)
		d.root.Add(cell.root)
	}
	for i, parsed := range cells {
		cell := d.cells[i]
		cell.root.Move(fyne.NewPos(float32(i)*d.digitSize.Width, 0))
		cell.root.Resize(d.digitSize)
		cell.set(d.resources, parsed.Symbol, parsed.Decimal, d.digitSize)
		cell.root.Show()
	}
	for i := len(cells); i < len(d.cells); i++ {
		d.cells[i].root.Hide()
	}
	d.root.Refresh()
}

func newPNGDigitCell(resources *pngDigitResources, size fyne.Size) *pngDigitCell {
	base := canvas.NewImageFromResource(resources.resource("blank"))
	base.FillMode = canvas.ImageFillContain
	base.Resize(size)
	decimalResource := resources.dp
	if decimalResource == nil {
		decimalResource = resources.resource("blank")
	}
	decimal := canvas.NewImageFromResource(decimalResource)
	decimal.FillMode = canvas.ImageFillContain
	decimal.Resize(size)
	decimal.Hide()
	root := container.NewWithoutLayout(base, decimal)
	root.Resize(size)
	return &pngDigitCell{root: root, base: base, decimal: decimal, symbol: "blank"}
}

func (c *pngDigitCell) set(resources *pngDigitResources, symbol string, dp bool, size fyne.Size) {
	if symbol == "" {
		symbol = "blank"
	}
	if c.symbol != symbol {
		c.base.Resource = resources.resource(symbol)
		c.symbol = symbol
		c.base.Refresh()
	}
	c.base.Resize(size)
	c.decimal.Resize(size)
	if dp != c.dp {
		c.dp = dp
		if dp && resources.dp != nil {
			c.decimal.Show()
		} else {
			c.decimal.Hide()
		}
	}
}

func loadPNGDigitResources(dir string) *pngDigitResources {
	resources := &pngDigitResources{dir: dir, bySymbol: make(map[string]fyne.Resource, 12)}
	for _, symbol := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "dash", "blank"} {
		resources.bySymbol[symbol] = loadPNGResource(filepath.Join(dir, symbol+".png"))
	}
	resources.dp = loadPNGResource(filepath.Join(dir, "dp.png"))
	return resources
}

func (r *pngDigitResources) resource(symbol string) fyne.Resource {
	if res := r.bySymbol[symbol]; res != nil {
		return res
	}
	if res := r.bySymbol["blank"]; res != nil {
		return res
	}
	return fyne.NewStaticResource("blank.png", blankPNG())
}

func loadPNGResource(path string) fyne.Resource {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return fyne.NewStaticResource(filepath.Base(path), bytes)
}

func formatDigitCells(value string, width int) []digitCell {
	cells := make([]digitCell, 0, len(value))
	for _, r := range strings.TrimSpace(value) {
		switch {
		case r >= '0' && r <= '9':
			cells = append(cells, digitCell{Symbol: string(r)})
		case r == '-':
			cells = append(cells, digitCell{Symbol: "dash"})
		case r == ' ':
			cells = append(cells, digitCell{Symbol: "blank"})
		case r == '.':
			if len(cells) > 0 {
				cells[len(cells)-1].Decimal = true
			}
		default:
			// Keep noisy input from breaking the dashboard; it has enough goblins already.
		}
	}
	if width <= 0 {
		return cells
	}
	if len(cells) > width {
		return cells[len(cells)-width:]
	}
	for len(cells) < width {
		cells = append([]digitCell{{Symbol: "blank"}}, cells...)
	}
	return cells
}

func blankPNG() []byte {
	// 1x1 transparent PNG fallback. Real dashboard assets should live under assets/digits/normal.
	return []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 8, 6, 0, 0, 0, 31, 21, 196, 137, 0, 0, 0, 13, 73, 68, 65, 84, 120, 156, 99, 248, 15, 4, 0, 9, 251, 3, 253, 167, 182, 129, 129, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
