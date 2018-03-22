package neuquant

import (
	"github.com/edvincandon/GoQuant/kohonen"
	"math"
)

type Pixel struct {
	R, G, B, A float64
}

func (p *Pixel) Distance(node kohonen.Node) float64 {
	n, ok := node.(*Pixel)
	if !ok {
		panic("cannot compare without a pixel")
	}
	baseDist := math.Abs(n.R-p.R) + math.Abs(n.G-p.G) + math.Abs(n.B-p.B) + math.Abs(n.A-p.A)

	deltaR := n.R - p.R
	deltaG := n.G - p.G
	deltaB := n.B - p.B
	deltaAlpha := n.A - p.A
	rgbDistanceSquared := (deltaR * deltaR + deltaG * deltaG + deltaB * deltaB) / 3.0

	return deltaAlpha * deltaAlpha / 2.0 + rgbDistanceSquared * n.A * p.A / (255.0 * baseDist)
}

func (p *Pixel) Move(a float64, node kohonen.Node) {
	n, ok := node.(*Pixel)
	if !ok {
		panic("cannot compare without a pixel")
	}

	p.A = a*n.A + (1-a)*p.A
	p.B = a*n.B + (1-a)*p.B
	p.G = a*n.G + (1-a)*p.G
	p.R = a*n.R + (1-a)*p.R
}
