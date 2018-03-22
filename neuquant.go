package goquant

import (
	"image"
	"math"
	"image/color"
)

type Pixel struct {
	R, G, B, A float64
}

type Neuron struct {
	weights    Pixel
	bias, freq float64
}

type SOMNetwork struct {
	network     []Neuron
	input       []Pixel
	beta, gamma float64
	alpha       AlphaFunc
	radius      RadiusFunc
}

type AlphaFunc func(int) float64

var AlphaDefault = func(cycle int) float64 {
	return math.Exp(-0.03 * float64(cycle))
}

type RadiusFunc func(int) int

var RadiusDefault = func(cycle int) int {
	return int(math.Round(32 * math.Exp(-0.0325*float64(cycle))))
}

func ExtractPixels(m image.Image) []Pixel {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y
	pixels := make([]Pixel, 0, w*h)
	for y := m.Bounds().Min.Y; y < h; y++ {
		for x := m.Bounds().Min.X; x < w; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			pixels = append(pixels, Pixel{
				R: float64(r >> 8),
				G: float64(g >> 8),
				B: float64(b >> 8),
				A: float64(a >> 8),
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
				freq:    initialFreq,
				bias:    0.0,
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
		beta:    1.0 / 1024.0,
		gamma:   1024.0,
		alpha:   AlphaDefault,
		radius:  RadiusDefault,
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

func (som *SOMNetwork) FindBMU(p Pixel) int {
	l := len(som.network)
	bestDist := math.MaxFloat64
	bestBiasDist := bestDist
	bestPos := 0
	bestBiasPos := bestPos

	for i := 0; i < l; i++ {
		n := som.network[i].weights
		dist := math.Abs(n.R-p.R) + math.Abs(n.G-p.G) + math.Abs(n.B-p.B) + math.Abs(n.A-p.A)
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

	return bestBiasPos
}

func (som *SOMNetwork) FindIndex(p Pixel) int {
	l := len(som.network)
	bestDist := math.MaxFloat64
	pos := 0
	for i := 0; i < l; i++ {
		n := som.network[i].weights
		dist := math.Abs(n.R-p.R) + math.Abs(n.G-p.G) + math.Abs(n.B-p.B) + math.Abs(n.A-p.A)
		if dist < bestDist {
			bestDist = dist
			pos = i
		}
	}

	return pos
}

func (som *SOMNetwork) Learn(cycles, samplingFactor int) {
	getNextPoint := som.nextPoint("default")
	l := len(som.input)
	pixelsPerCycle := l / (cycles * samplingFactor)

	c := 0
	for c < cycles {
		i := 0
		alpha := som.alpha(c)
		radius := som.radius(c)
		for i < pixelsPerCycle {
			p := som.input[getNextPoint()]
			som.updateWeights(som.FindBMU(p), radius, p, alpha)
			i++
		}
		c++
	}
}

func (som *SOMNetwork) GetPalette() []color.Color {
	palette := make([]color.Color, 0, len(som.network))
	for i := 0; i< len(som.network); i++ {
		c := som.network[i].weights
		palette = append(palette, color.RGBA{
			R: uint8(int(c.R)),
			G: uint8(int(c.G)),
			B: uint8(int(c.B)),
			A: uint8(int(c.A)),
		})
	}
	return palette
}

func (som *SOMNetwork) updateWeights(bmuIndex, radius int, point Pixel, alpha float64) {
	min := int(math.Max(0, float64(bmuIndex-radius)))
	max := int(math.Min(float64(len(som.network)-1), float64(bmuIndex+radius)))
	i := min
	for i <= max {
		a := alpha*float64(1 - (i-bmuIndex)*(i-bmuIndex)/(radius*radius))
		som.network[i].weights.A = a*point.A + (1-a)*som.network[i].weights.A
		som.network[i].weights.B = a*point.B + (1-a)*som.network[i].weights.B
		som.network[i].weights.G = a*point.G + (1-a)*som.network[i].weights.G
		som.network[i].weights.R = a*point.R + (1-a)*som.network[i].weights.R
		i++
	}
}