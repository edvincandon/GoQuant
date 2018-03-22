package goquant

import (
	"image"
)

type Pixel struct {
	r, g, b, a uint32
}

func ExtractPixels(m image.Image) []Pixel {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y
	pixels := make([]Pixel, 0, w*h)
	for y := m.Bounds().Min.Y; y < h; y++ {
		for x := m.Bounds().Min.X; x < w; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			pixels = append(pixels, Pixel{
				r: r >> 8,
				g: g >> 8,
				b: b >> 8,
				a: a >>	 8,
			})
		}
	}
	return pixels
}

type SOMNetwork struct {
	network []Neuron
	input []Pixel
}

type Neuron struct {
	weight Pixel
	bias , freq float64
}

func NewSOMNetwork(size int, inut[]Pixel) SOMNetwork {
	
}