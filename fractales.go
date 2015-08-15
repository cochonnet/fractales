// Copyright 2015 Cochonnet.

package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"net/http"

	"github.com/biogo/graphics/palette"
)

func main() {
	port := flag.String("port", "localhost:8080", "Port to listen to")
	flag.Parse()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/img.png", imgHandler)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	http.Redirect(w, r, "/img.png", 301)
}

func imgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	//w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
	pointX := -2. // X coordinate of starting point of Mandelbrot or fix point for Julia (range: 2.0 to 2.0)
	pointY := -2. // Y coordinate of starting point of Mandelbrot or fix point for Julia (range: 2.0 to 2.0)
	zoom := 1.    // Zoom level (only working properly for Mandelbrot)
	julia := false
	maxIter := 61
	sizeX := 4000
	sizeY := 4000
	img := image.NewNRGBA(image.Rectangle{Max: image.Point{sizeX, sizeY}})
	calculateImg(img, pointX, pointY, zoom, julia, maxIter)
	png.Encode(w, img)
}

func calculateImg(img *image.NRGBA, pointX, pointY, zoom float64, julia bool, maxIter int) {
	minCx := -2.
	minCy := -2.
	if !julia {
		minCx = pointX
		minCy = pointY
	}
	bounds := img.Bounds()
	stepX := math.Abs(minCx-2.) / float64(bounds.Dx()) / zoom
	stepY := math.Abs(minCy-2.) / float64(bounds.Dy()) / zoom
	pal := paletteToNRGBA(palette.Rainbow(maxIter, 0, 1, 1, 1, 1))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		cy := minCy + float64(y)*stepY
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.SetNRGBA(x, y, pal[pointIteration(minCx+float64(x)*stepX, cy, pointX, pointY, julia, maxIter)])
		}
	}
}

// pointIteration returns the number of iterations at cx, cy before it diverges.
func pointIteration(cx, cy, pointX, pointY float64, julia bool, maxIter int) int {
	x := 0.
	y := 0.
	if julia {
		x = cx
		y = cy
		cx = pointX
		cy = pointY
	}
	iter := 0
	for squAbs := 0.; squAbs <= 4. && iter < maxIter; iter++ {
		xt := (x * x) - (y * y) + cx
		yt := (2. * x * y) + cy
		x = xt
		y = yt
		squAbs = (x * x) + (y * y)
	}
	return iter
}

// paletteToNRGBA renders the colors as NRGBA and append black.
func paletteToNRGBA(pal palette.Palette) []color.NRGBA {
	colors := pal.Colors()
	out := make([]color.NRGBA, len(colors)+1)
	for i := range colors {
		r, g, b, a := colors[i].RGBA()
		out[i] = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	out[len(colors)] = color.NRGBA{0, 0, 0, 255}
	return out
}
