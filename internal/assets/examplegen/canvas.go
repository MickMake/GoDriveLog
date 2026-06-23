package examplegen

import (
	"hash/fnv"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// Canvas wraps a mutable RGBA image with small procedural-drawing helpers.
type Canvas struct {
	Image *image.NRGBA
}

func NewCanvas(width, height int, fill color.NRGBA) *Canvas {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	canvas := &Canvas{Image: img}
	canvas.FillRect(0, 0, width, height, fill)
	return canvas
}

func (c *Canvas) FillRect(x, y, width, height int, fill color.NRGBA) {
	if c == nil || c.Image == nil || width <= 0 || height <= 0 {
		return
	}
	bounds := c.Image.Bounds()
	startX, endX := clampSpan(x, width, bounds.Min.X, bounds.Max.X)
	startY, endY := clampSpan(y, height, bounds.Min.Y, bounds.Max.Y)
	for py := startY; py < endY; py++ {
		for px := startX; px < endX; px++ {
			blendPixel(c.Image, px, py, fill)
		}
	}
}

func (c *Canvas) StrokeRect(x, y, width, height, thickness int, stroke color.NRGBA) {
	if thickness <= 0 || width <= 0 || height <= 0 {
		return
	}
	c.FillRect(x, y, width, thickness, stroke)
	c.FillRect(x, y+height-thickness, width, thickness, stroke)
	c.FillRect(x, y, thickness, height, stroke)
	c.FillRect(x+width-thickness, y, thickness, height, stroke)
}

func (c *Canvas) FillCircle(cx, cy, radius int, fill color.NRGBA) {
	if c == nil || c.Image == nil || radius <= 0 {
		return
	}
	radiusSquared := radius * radius
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= radiusSquared {
				blendPixel(c.Image, x, y, fill)
			}
		}
	}
}

func (c *Canvas) StrokeCircle(cx, cy, radius, thickness int, stroke color.NRGBA) {
	if c == nil || c.Image == nil || radius <= 0 || thickness <= 0 {
		return
	}
	outerSquared := radius * radius
	inner := radius - thickness
	if inner < 0 {
		inner = 0
	}
	innerSquared := inner * inner
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx := x - cx
			dy := y - cy
			distanceSquared := dx*dx + dy*dy
			if distanceSquared <= outerSquared && distanceSquared >= innerSquared {
				blendPixel(c.Image, x, y, stroke)
			}
		}
	}
}

func (c *Canvas) DrawGrid(spacing, thickness int, stroke color.NRGBA) {
	if c == nil || c.Image == nil || spacing <= 0 || thickness <= 0 {
		return
	}
	bounds := c.Image.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x += spacing {
		c.FillRect(x, bounds.Min.Y, thickness, bounds.Dy(), stroke)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y += spacing {
		c.FillRect(bounds.Min.X, y, bounds.Dx(), thickness, stroke)
	}
}

func (c *Canvas) AddGrain(seed uint64, strength uint8) {
	if c == nil || c.Image == nil || strength == 0 {
		return
	}
	bounds := c.Image.Bounds()
	half := int(strength) / 2
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			index := c.Image.PixOffset(x, y)
			alpha := c.Image.Pix[index+3]
			if alpha == 0 {
				continue
			}
			noise := int(noise64(seed, x, y)%uint64(strength+1)) - half
			c.Image.Pix[index+0] = clampChannel(int(c.Image.Pix[index+0]) + noise)
			c.Image.Pix[index+1] = clampChannel(int(c.Image.Pix[index+1]) + noise)
			c.Image.Pix[index+2] = clampChannel(int(c.Image.Pix[index+2]) + noise)
		}
	}
}

func (c *Canvas) WritePNG(path string) error {
	if c == nil || c.Image == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, c.Image)
}

func HashSeed(text string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(text))
	return h.Sum64()
}

func clampSpan(start, width, min, max int) (int, int) {
	if start < min {
		start = min
	}
	end := start + width
	if end > max {
		end = max
	}
	if start > max {
		start = max
	}
	if end < min {
		end = min
	}
	return start, end
}

func blendPixel(img *image.NRGBA, x, y int, src color.NRGBA) {
	if !image.Pt(x, y).In(img.Bounds()) {
		return
	}
	index := img.PixOffset(x, y)
	dst := color.NRGBA{
		R: img.Pix[index+0],
		G: img.Pix[index+1],
		B: img.Pix[index+2],
		A: img.Pix[index+3],
	}
	out := blend(dst, src)
	img.Pix[index+0] = out.R
	img.Pix[index+1] = out.G
	img.Pix[index+2] = out.B
	img.Pix[index+3] = out.A
}

func blend(dst, src color.NRGBA) color.NRGBA {
	srcA := int(src.A)
	dstA := int(dst.A)
	outA := srcA + (dstA*(255-srcA)+127)/255
	if outA == 0 {
		return color.NRGBA{}
	}

	blendChannel := func(dstC, srcC uint8) uint8 {
		dstPart := int(dstC) * dstA * (255 - srcA)
		srcPart := int(srcC) * srcA * 255
		value := (srcPart + dstPart + outA*127) / (outA * 255)
		return clampChannel(value)
	}

	return color.NRGBA{
		R: blendChannel(dst.R, src.R),
		G: blendChannel(dst.G, src.G),
		B: blendChannel(dst.B, src.B),
		A: uint8(outA),
	}
}

func noise64(seed uint64, x, y int) uint64 {
	value := seed ^ (uint64(uint32(x)) * 0x9e3779b97f4a7c15) ^ (uint64(uint32(y)) * 0xc2b2ae3d27d4eb4f)
	value ^= value >> 30
	value *= 0xbf58476d1ce4e5b9
	value ^= value >> 27
	value *= 0x94d049bb133111eb
	value ^= value >> 31
	return value
}

func clampChannel(value int) uint8 {
	switch {
	case value < 0:
		return 0
	case value > 255:
		return 255
	default:
		return uint8(value)
	}
}
