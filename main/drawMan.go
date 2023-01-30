package main

import (
	"bufio"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func drawPoint() {

}

func calculateDistance(p1 []float64, p2 []float64) float64 {
	return math.Sqrt(math.Pow(p1[0]-p2[0], 2) + math.Pow(p1[1]-p2[1], 2))
}
func getClosestPoint(fromPoint []float64, toPoints [][]float64) ([]float64, int) {
	var closestPoint []float64
	var smallestDistance float64 = 100000
	var closestPointIndex int
	for i, toPoint := range toPoints {
		distance := calculateDistance(fromPoint, toPoint)
		if distance < smallestDistance {
			smallestDistance = distance
			closestPoint = toPoint
			closestPointIndex = i
		}
	}
	return closestPoint, closestPointIndex
}

func drawLineBetweenPoints(p1 []float64, p2 []float64, drawingContext *gg.Context) {
	drawingContext.SetRGBA(1, 0, 0, 1)
	drawingContext.SetLineWidth(40)
	drawingContext.DrawLine(p1[0], p1[1], p2[0], p2[1])
	drawingContext.Stroke()
}

func remove(slice [][]float64, s int) [][]float64 {
	return append(slice[:s], slice[s+1:]...)
}

func chartMap(imageName string, spawnPosition []float64, toPoints [][]float64) {
	imgFile, loadError := gg.LoadPNG("./images/" + imageName)
	if loadError != nil {
		panic(loadError)
	}
	// set initial point
	closestPoint, closestPointIndex := getClosestPoint(spawnPosition, toPoints)
	fmt.Printf("len = %v \n", len(toPoints))
	toPoints = remove(toPoints, closestPointIndex)

	// draw initial Point
	dc := gg.NewContextForImage(imgFile)
	dc.Push()
	drawLineBetweenPoints(spawnPosition, closestPoint, dc)

	for len(toPoints) != 0 {
		fmt.Printf("LEN = %v \n", len(toPoints))
		currClosestPoint, currClosestPointIndex := getClosestPoint(closestPoint, toPoints) // get point closest to last closest point
		toPoints = remove(toPoints, currClosestPointIndex)                                 // remove that point from the slice
		drawLineBetweenPoints(closestPoint, currClosestPoint, dc)                          // draw a line between them
		closestPoint = currClosestPoint                                                    // assign the closest point as the prev closest point
	}

	dc.Pop()
	dc.SavePNG("woods_drawn.png")

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
