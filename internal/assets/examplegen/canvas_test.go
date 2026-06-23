package examplegen

import (
	"bytes"
	"image/color"
	"testing"
)

func TestAddGrainDeterministic(t *testing.T) {
	left := NewCanvas(12, 12, color.NRGBA{R: 20, G: 30, B: 40, A: 255})
	right := NewCanvas(12, 12, color.NRGBA{R: 20, G: 30, B: 40, A: 255})

	seed := HashSeed("framework-smoke")
	left.AddGrain(seed, 12)
	right.AddGrain(seed, 12)

	if !bytes.Equal(left.Image.Pix, right.Image.Pix) {
		t.Fatal("grain output differs for identical seed")
	}
}
