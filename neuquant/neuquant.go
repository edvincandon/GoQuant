package neuquant

import (
	"github.com/edvincandon/GoQuant/kohonen"
	"image"
	"image/color"
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

	return baseDist
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

func Quantize(img image.Image) (kohonen.SOM, []color.Color) {
	som := kohonen.NewSOM(
		256,
		func(i int) kohonen.Neuron {
			var node Pixel
			switch(i) {
			case 0:
				node = Pixel{255.0, 255.0, 255.0, 255.0}
			case 1:
				node = Pixel{255.0, 255.0, 255.0, 0.0}
			case 2:
				node = Pixel{0.0, 0.0, 0.0, 255.0}
			case 3:
				node = Pixel{0.0, 0.0, 0.0, 0.0}
			default:
				node = Pixel{float64(i), float64(i), float64(i), 255.0}
			}
			return kohonen.Neuron{
				Node: &node,
				Freq: 1.0 / 256.0,
				Bias: 0.0,
			}
		},
		kohonen.SOMConfig{
			NCycle:   100,
			Sampling: 1,
			Beta:     1.0 / 1024.0,
			Gamma:    1024.0,
			Alpha:    kohonen.AlphaDefault,
			Radius:   kohonen.RadiusDefault,
			Input:    kohonen.InputDefault,
		},
	)

	pixels := ExtractPixels(img)
	nodes := som.Learn(pixels)

	colors := make([]color.Color, 0, 256)
	for _, n := range nodes {
		p, ok := n.(*Pixel)
		if !ok {
			panic("expected pixel")
		}
		colors = append(colors, color.RGBA{
			R: uint8(p.R / p.A),
			G: uint8(p.G / p.A),
			B: uint8(p.B / p.A),
			A: uint8(p.A),
		})
	}

	return som, colors
}

func ExtractPixels(m image.Image) []kohonen.Node {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y
	pixels := make([]kohonen.Node, 0, w*h)
	for y := m.Bounds().Min.Y; y < h; y++ {
		for x := m.Bounds().Min.X; x < w; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			pixels = append(pixels, &Pixel{
				R: float64(r >> 8) * float64(a >> 8),
				G: float64(g >> 8) * float64(a >> 8),
				B: float64(b >> 8) * float64(a >> 8),
				A: float64(a >> 8),
			})
		}
	}
	return pixels
}
