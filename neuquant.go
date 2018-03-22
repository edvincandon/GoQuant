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
	cycles int
	samplingFactor int
	initialRadius int
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
	// sets up the neurons with the initial weights
	// default method will use the diagonal (0,0,0,0) - (255,255,255,255)
	// add more later
	neurons := make([]Neuron, size)
	initialFreq :=  1.0 / float64(size)

	switch method {
		case "default":
			for i := 0; i < size; i++ {
				val := float64(i)
				neurons[i] = Neuron{
					weight: Pixel{val, val, val, val},
					freq: initialFreq,
					bias: 0.0}
			}
	}

	return neurons
}

func NewSOMNetwork(size int, input[]Pixel) SOMNetwork {
	neurons := initNeurons(size, "default")
	network := SOMNetwork{
		network: neurons,
		input: input,
		cycles: 100,
		samplingFactor: 3,
		initialRadius: 32}

	return network
}
