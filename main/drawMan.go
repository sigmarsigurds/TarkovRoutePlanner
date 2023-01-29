package main

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func drawPoint() {

}

func drawLine(imageName string) {
	imgFile, loadError := os.Open(imageName)
	if loadError != nil {
		panic(loadError)
	}
	defer imgFile.Close()
	reader := bufio.NewReader(imgFile)

	// Load a PNG image
	img, _ := png.Decode(reader)

	// Create a new RGBA image
	dst := image.NewRGBA(img.Bounds())
	// Draw the PNG image onto the new RGBA image
	draw.Draw(dst, dst.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Draw a line from (10,10) to (100,100)
	color := color.RGBA{255, 0, 0, 255}
	for x := 10; x <= 100; x++ {
		dst.Set(x, 10+(x-10)*(100-10)/(100-10), color)
		dst.Set(x, 20+(x-10)*(100-10)/(100-10), color)
	}

	// Save the RGBA image as a PNG file
	f, _ := os.Create("output.png")
	defer f.Close()
	encodeError := png.Encode(f, dst)
	if encodeError != nil {
		panic(encodeError)
	}
}
