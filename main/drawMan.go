package main

import (
	"github.com/fogleman/gg"
	"math"
)

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

func flipPositions(toPoints [][]float64, mapBounds [][]int) [][]float64 {
	const offset = 200
	mapHeight := mapBounds[1][0]
	//mapWidth := mapBounds[1][1]
	var retPoints [][]float64
	for _, toPoint := range toPoints {
		//flippedPoint := []float64{float64(mapWidth) - toPoint[0], float64(mapHeight) - toPoint[1]}
		flippedPoint := []float64{toPoint[0], float64(mapHeight) - toPoint[1] - offset}
		retPoints = append(retPoints, flippedPoint)
	}
	return retPoints
}

func scalePositions(toPoints [][]float64) [][]float64 {
	const scale = 1
	var retPoints [][]float64
	for _, toPoint := range toPoints {
		scaledPoint := []float64{scale * toPoint[0], scale * toPoint[1]}
		retPoints = append(retPoints, scaledPoint)
	}
	return retPoints
}

func drawPoints(toPoints [][]float64, drawingContext *gg.Context) {

	for _, toPoint := range toPoints {
		drawingContext.DrawCircle(toPoint[0], toPoint[1], 60)
		drawingContext.SetRGB(0, 1, 0)
		drawingContext.Fill()
	}
}

func drawPoint(position []float64, drawingContext *gg.Context) {
	drawingContext.DrawCircle(position[0], position[1], 60)
	drawingContext.SetRGB(1, 0, 1)
	drawingContext.Fill()
}

func chartMap(imageName string, spawnPosition []float64, exitPositions [][]float64, toPoints [][]float64, mapBounds [][]int, numberOfStashes int) {
	imgFile, loadError := gg.LoadPNG("./images/" + imageName)
	if loadError != nil {
		panic(loadError)
	}
	toPoints = flipPositions(toPoints, mapBounds)

	dc := gg.NewContextForImage(imgFile)
	drawPoint(spawnPosition, dc)

	// draw toPoints for debug
	drawPoints(toPoints, dc)

	// set initial point
	closestPoint, closestPointIndex := getClosestPoint(spawnPosition, toPoints)
	toPoints = remove(toPoints, closestPointIndex)

	// draw initial Line
	dc.Push()
	drawLineBetweenPoints(spawnPosition, closestPoint, dc)

	for len(toPoints) != 0 {
		if numberOfStashes == 0 {
			break
		}
		currClosestPoint, currClosestPointIndex := getClosestPoint(closestPoint, toPoints) // get point closest to last closest point
		toPoints = remove(toPoints, currClosestPointIndex)                                 // remove that point from the slice
		drawLineBetweenPoints(closestPoint, currClosestPoint, dc)                          // draw a line between them
		closestPoint = currClosestPoint                                                    // assign the closest point as the prev closest point
		numberOfStashes--

	}

	closestExit, _ := getClosestPoint(closestPoint, exitPositions)
	drawLineBetweenPoints(closestPoint, closestExit, dc)
	drawPoint(closestExit, dc)

	dc.Pop()
	drawnImageName := "DRAWN_" + imageName
	dc.SavePNG(drawnImageName)
}
