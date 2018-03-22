package goquant

import (
	"image"
)

type Pixel struct {
	r, g, b, a float64
}

type Neuron struct {
	weight     Pixel
	bias, freq float64
}

type SOMNetwork struct {
	network []Neuron
	input   []Pixel
}

func (n *SOMNetwork) nextPoint(nextType string) func() int {
	const (
		prime1 = 499
		prime2 = 491
		prime3 = 487
		prime4 = 503
	)

	var fn func() int
	switch nextType {
	case "default":
		s := len(n.input)
		var p int
		switch {
		case s%prime1 != 0:
			p = prime1
		case s%prime2 != 0:
			p = prime2
		case s%prime3 != 0:
			p = prime3
		default:
			p = prime4
		}

		pos := 0
		fn = func() int {
			pos += p
			pos = pos % s
			return pos
		}
	}

	return fn
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
				a: float64(a >> 8),
			})
		}
	}
	return pixels
}

func initNeurons(size int, initType string) []Neuron {
	// sets up the neurons with the initial weights
	// default method will use the diagonal (0,0,0,0) - (255,255,255,255)
	// add more later
	neurons := make([]Neuron, size)
	initialFreq := 1.0 / float64(size)

	switch initType {
	case "default":
		for i := 0; i < size; i++ {
			val := float64(i)
			neurons[i] = Neuron{
				weight: Pixel{val, val, val, val},
				freq:   initialFreq,
				bias:   0.0,
			}
		}
	}

	return neurons
}

func NewSOMNetwork(size int, input []Pixel) SOMNetwork {
	neurons := initNeurons(size, "default")
	network := SOMNetwork{
		network: neurons,
		input:   input,
	}

	return network
}

func (som *SOMNetwork) Learn(cycles, samplingFactor int, initialRadius float64) {
	l := len(som.input)
	pixelsPerCycle := l / (cycles * samplingFactor)

	c := 0
	for c < cycles {
		i := 0
		for i < pixelsPerCycle {
			i++
		}
		c++
	}
}