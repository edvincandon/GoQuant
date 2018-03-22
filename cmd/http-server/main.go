package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
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

func Quantize(img image.Image) image.Image {
	return img
}
