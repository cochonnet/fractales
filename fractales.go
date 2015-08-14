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
	maxIter := 51
	sizeX := 1000
	sizeY := 1000
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
	maxSquAbs := 4.
	bounds := img.Bounds()
	stepX := math.Abs(minCx-2.) / float64(bounds.Dx()) / zoom
	stepY := math.Abs(minCy-2.) / float64(bounds.Dy()) / zoom
	//var wg sync.WaitGroup
	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		cy := minCy + float64(py)*stepY
		// Une goroutine par ligne.
		//wg.Add(1)
		//go func(cy float64) {
		//defer wg.Done()
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			cx := minCx + float64(px)*stepX
			itr := pointIteration(cx, cy, pointX, pointY, maxSquAbs, julia, maxIter)
			img.Set(px, py, chooseColor(itr, maxIter))
		}
		//}(cy)
	}
	//wg.Wait()
}

func pointIteration(cx, cy, pointX, pointY, maxSquAbs float64, julia bool, maxIter int) int {
	x := 0.
	y := 0.
	if julia {
		x = cx
		y = cy
		cx = pointX
		cy = pointY
	}
	iter := 0
	for squAbs := 0.; squAbs <= maxSquAbs && iter < maxIter; iter++ {
		xt := (x * x) - (y * y) + cx
		yt := (2. * x * y) + cy
		x = xt
		y = yt
		squAbs = (x * x) + (y * y)
	}
	return iter
}

var pal = palette.Rainbow(256, 0, 1, 1, 1, 1).Colors()

func chooseColor(iterValue int, maxIter int) color.Color {
	if iterValue != maxIter {
		return pal[uint8(iterValue*255/maxIter)]
		// Vert.
		//return color.NRGBA{0, val * uint8(255/maxIter), 0, 255}
	}
	return color.NRGBA{0, 0, 0, 255}
}
