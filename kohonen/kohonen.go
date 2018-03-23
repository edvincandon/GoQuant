package kohonen

import (
	"math"
)

type Node interface {
	Distance(Node) float64
	Move(float64, Node)
}

type Neuron struct {
	Node       Node
	Bias, Freq float64
}

type SOM struct {
	network []Neuron
	size    int
	config  SOMConfig
}

type SOMConfig struct {
	NCycle   int
	Sampling int
	Beta     float64
	Gamma    float64
	Alpha    AlphaFunc
	Radius   RadiusFunc
	Input    InputFunc
}

type InitNeuronFunc func(int) Neuron

type AlphaFunc func(int) float64

type RadiusFunc func(int) int

type InputFunc func(int, int) int

var AlphaDefault = func(cycle int) float64 {
	return math.Exp(-0.03 * float64(cycle))
}

var RadiusDefault = func(cycle int) int {
	return int(math.Round(32 * math.Exp(-0.0325*float64(cycle))))
}

var InputDefault = func(i, l int) int {
	primes := []int{487, 491, 499, 503}

	for _, p := range primes {
		if l%p != 0 {
			continue
		}
		return i * p % l
	}

	return i
}

func NewSOM(size int, init InitNeuronFunc, config SOMConfig) SOM {
	network := make([]Neuron, 0, size)
	for i := 0; i < size; i++ {
		network = append(network, init(i))
	}
	return SOM{
		network: network,
		size:    size,
		config:  config,
	}
}

func (som *SOM) Learn(in []Node) []Node {
	nPerCycle := len(in) / (som.config.NCycle * som.config.Sampling)

	for c := 0; c < som.config.NCycle; c++ {
		alpha := som.config.Alpha(c)
		radius := som.config.Radius(c)
		for i := 0; i < nPerCycle; i++ {
			inputIndex := som.config.Input(c*nPerCycle+i, len(in))
			closestNeuronIndex := som.findClosestNeuronIndex(in[inputIndex], true)
			som.updateNetwork(closestNeuronIndex, radius, in[inputIndex], alpha)
		}
	}

	nodes := make([]Node, 0, som.size)
	for _, neuron := range som.network {
		nodes = append(nodes, neuron.Node)
	}

	return nodes
}

func (som *SOM) FindClosestNeuronIndex(n Node) int {
	return som.findClosestNeuronIndex(n, false)
}

func (som *SOM) findClosestNeuronIndex(n Node, updateNetwork bool) int {
	bestDistance := math.MaxFloat64
	bestBiasDistance := math.MaxFloat64
	bestPos := 0
	bestBiasPos := 0

	for i := 0; i < som.size; i++ {
		d := n.Distance(som.network[i].Node)
		if d < bestDistance {
			bestDistance = d
			bestPos = i
		}

		if !updateNetwork {
			continue
		}

		biasDistance := d - som.network[i].Bias
		if biasDistance < bestBiasDistance {
			bestBiasDistance = biasDistance
			bestBiasPos = i
		}

		som.network[i].Bias += som.config.Beta
		som.network[i].Freq -= som.config.Beta * som.network[i].Freq
	}

	if updateNetwork {
		som.network[bestPos].Bias -= som.config.Beta * som.config.Gamma
		som.network[bestPos].Freq += som.config.Beta

		return bestBiasPos
	}

	return bestPos
}

func (som *SOM) updateNetwork(closestNeuronIndex, radius int, target Node, alpha float64) {
	min := int(math.Max(0, float64(closestNeuronIndex-radius)))
	max := int(math.Min(float64(som.size-1), float64(closestNeuronIndex+radius)))

	for i := min; i <= max; i++ {
		q := float64(i - closestNeuronIndex)
		r2 := float64(radius * radius)
		a := alpha * (r2 - q*q) / r2
		som.network[i].Node.Move(a, target)
	}
}
