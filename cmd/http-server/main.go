package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"

	"github.com/edvincandon/GoQuant/neuquant"
)

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleForm(w, r)
		case http.MethodPost:
			handleUpload(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<h1>Upload image</h1>
		<form enctype="multipart/form-data" action="/upload" method="POST">
			<input name="file" type="file">
			<input type="submit" value="Upload">
		</form>
	`)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Please upload a file : %s", err.Error()), http.StatusBadRequest)
		return
	}
	img, f, err := image.Decode(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot decode image : %s", err.Error()), http.StatusBadRequest)
		return
	}

	imgOut := Quantize(img)

	imgInBase64, sizeIn, ok := encodeImg(w, img, f)
	if !ok {
		return
	}
	imgOutBase64, sizeOut, ok := encodeImg(w, imgOut, f)
	if !ok {
		return
	}

	tplStr := `	<!DOCTYPE html><html lang="en"><head></head><body>
		<h2>Original Image</h2>
		<img src="data:image/jpg;base64,{{ .ImgIn }}" height="400">
		<p>Size : {{ .SizeIn }} bytes</p>
		<h2>Compressed Image</h2>
		<img src="data:image/jpg;base64,{{ .ImgOut }}" height="400">
		<p>Size : {{ .SizeOut }} bytes</p>
		</body>`

	if tpl, err := template.New("image").Parse(tplStr); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println("unable to parse image template.")
		return
	} else {
		data := map[string]interface{}{"ImgIn": imgInBase64, "ImgOut": imgOutBase64, "SizeIn": sizeIn, "SizeOut": sizeOut}
		if err = tpl.Execute(w, data); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println("unable to execute template.")
			return
		}
	}
}

func encodeImg(w http.ResponseWriter, img image.Image, f string) (string, int, bool) {
	var buf bytes.Buffer
	var err error

	switch f {
	case "png":
		err = png.Encode(&buf, img)
	case "gif":
		err = gif.Encode(&buf, img, nil)
	case "jpeg":
		err = jpeg.Encode(&buf, img, nil)
	default:
		return "", 0, false
	}

	if err != nil {
		http.Error(w, "unable to encode image.", http.StatusInternalServerError)
		return "", 0, false
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), buf.Len(), true
}

func CloneToRGBA(src image.Image) *image.NRGBA {
	b := src.Bounds()
	dst := image.NewNRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}

func Quantize(imgOriginal image.Image) image.Image {
	img := CloneToRGBA(imgOriginal)
	_, pal := neuquant.Quantize(img)
	imgOut := image.NewPaletted(img.Bounds(), pal)

	draw.FloydSteinberg.Draw(imgOut, img.Bounds(), img, image.Point{0,0})
/*
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			i := som.FindClosestNeuronIndex(&neuquant.Pixel{
				R: float64(r >> 8),
				G: float64(g >> 8),
				B: float64(b >> 8),
				A: float64(a >> 8),
			})
			imgOut.SetColorIndex(x, y, uint8(i))

			dither(img, x, y, pal[i], img.At(x, y))
		}
	}
*/
	return imgOut
}

func dither(img *image.NRGBA, x, y int, actual, target color.Color) {
	if x == 0 || x >= img.Bounds().Max.X || y == 0 || y >= img.Bounds().Max.Y {
		return
	}

	r1, g1, b1, a1 := actual.RGBA()
	r2, g2, b2, a2 := target.RGBA()
	e := [4]float64{
		float64(r2>>8) - float64(r1>>8),
		float64(g2>>8) - float64(g1>>8),
		float64(b2>>8) - float64(b1>>8),
		float64(a2>>8) - float64(a1>>8),
	}

	ditherPixel(img, x+1, y, e, 7.0/16.0)
	ditherPixel(img, x-1, y+1, e, 3.0/16.0)
	ditherPixel(img, x, y+1, e, 5.0/16.0)
	ditherPixel(img, x+1, y+1, e, 1.0/16.0)

}

func ditherPixel(img *image.NRGBA, x, y int, e [4]float64, percent float64) {
	r, g, b, a := img.At(x, y).RGBA()
	r1 := float64(r>>8) + e[0]*percent/2
	g1 := float64(g>>8) + e[1]*percent/2
	b1 := float64(b>>8) + e[2]*percent/2
	a1 := float64(a>>8) + e[3]*percent

	img.Set(x, y, color.RGBA{uint8(r1), uint8(g1), uint8(b1), uint8(a1)})
}
