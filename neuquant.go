package goquant

import (
	"image"
)

type Pixel struct {
	r, g, b, a float64
}

type SOMNetwork struct {
	network []Neuron
	input []Pixel
}

type Neuron struct {
	weight Pixel
	bias, freq float64
}

func ExtractPixels(m image.Image) []Pixel {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y
	pixels := make([]Pixel, 0, w*h)
	for y := m.Bounds().Min.Y; y < h; y++ {
		for x := m.Bounds().Min.X; x < w; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			pixels = append(pixels, Pixel{
				r: float64(r >> 8),
				g: float64(g >> 8),
				b: float64(b >> 8),
				a: float64(a >>	8),
			})
		}
	}
	return pixels
}

func initNeurons(size int, method string) []Neuron {
	neurons := make([]Neuron, size)

	switch method {
		case "default":
			for i := 0; i < size; i++ {
				val := float64(i)
				neurons[i] = Neuron{ weight: Pixel{val, val, val, val}, bias: 0.0, freq: 0.0 }
			}
	}

	return neurons
}

func NewSOMNetwork(size int, input[]Pixel) SOMNetwork {
	neurons := initNeurons(size, "default")
	network := SOMNetwork{neurons, input}

	return network
}
