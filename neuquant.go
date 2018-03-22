package goquant

import (
	"image"
	"math"
)

type Pixel struct {
	r, g, b, a float64
}

type Neuron struct {
	weights     Pixel
	bias, freq float64
}

type SOMNetwork struct {
	network []Neuron
	input   []Pixel
	beta, gamma float64
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
	neurons := make([]Neuron, size)
	initialFreq := 1.0 / float64(size)

	switch initType {
	case "default":
		for i := 0; i < size; i++ {
			val := float64(i)
			neurons[i] = Neuron{
				weights: Pixel{val, val, val, val},
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
		beta: 1.0 / 1024.0,
		gamma : 1024.0,
	}

	return network
}

func (som *SOMNetwork) nextPoint(nextType string) func() int {
	var fn func() int
	l := len(som.input)

	const (
		prime1 = 499
		prime2 = 491
		prime3 = 487
		prime4 = 503
	)

	switch nextType {
	case "default":
		var p int
		switch {
		case l%prime1 != 0:
			p = prime1
		case l%prime2 != 0:
			p = prime2
		case l%prime3 != 0:
			p = prime3
		default:
			p = prime4
		}

		pos := 0
		fn = func() int {
			pos += p
			pos = pos % l
			return pos
		}
	}

	return fn
}

func (som *SOMNetwork) findBMU(p Pixel) int {
	l := len(som.network)
	bestDist := math.MaxFloat64
	bestBiasDist := bestDist
	bestPos := 0
	bestBiasPos := bestPos

	for i := 0; i < l; i++ {
			n := som.network[i].weights
			dist := math.Abs(n.r - p.r) + math.Abs(n.g - p.g) + math.Abs(n.b - p.b) + math.Abs(n.a - p.a)
			if dist < bestDist {
				bestDist = dist
				bestPos = i
			}

			biasDist := dist - som.network[i].bias
			if biasDist < bestBiasDist {
				bestBiasDist = biasDist
				bestBiasPos = i
			}

			som.network[i].freq -= som.beta * som.network[i].freq
			som.network[i].bias += som.beta * som.gamma * som.network[i].freq
		}

	som.network[bestPos].freq += som.beta
	som.network[bestPos].bias -= som.beta * som.gamma

	return bestPos
}

func (som *SOMNetwork) Learn(cycles, samplingFactor int, initialRadius float64) {
	getNextPoint := som.nextPoint("default")
	l := len(som.input)
	pixelsPerCycle := l / (cycles * samplingFactor)

	c := 0
	for c < cycles {
		i := 0
		for i < pixelsPerCycle {
			p := som.input[getNextPoint()]
			bmu := som.findBMU(p)
			i++
		}
		c++
	}
}
