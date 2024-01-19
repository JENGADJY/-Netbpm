package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Pixel struct {
	R, G, B uint8
}

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

type Point struct {
	X, Y int
}

// ReadPPM lit une image PPM à partir d'un fichier et renvoie un objet PPM.

func ReadPPM(fileName string) (*PPM, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, errors.New("Error reading magic number")
	}
	magicNumber := scanner.Text()
	if magicNumber != "P3" && magicNumber != "P6" {
		return nil, errors.New("Unsupported PPM format")
	}

	if !scanner.Scan() {
		return nil, errors.New("Error reading dimensions")
	}
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, errors.New("Invalid dimensions in PPM file")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("Error converting width: %v", err)
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("Error converting height: %v", err)
	}

	if !scanner.Scan() {
		return nil, errors.New("Error reading max value")
	}
	maxVal, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("Error converting max value: %v", err)
	}

	data := make([][]Pixel, height)
	for i := range data {
		data[i] = make([]Pixel, width)
		for j := range data[i] {
			if !scanner.Scan() {
				return nil, errors.New("Error reading pixel values")
			}
			values := strings.Fields(scanner.Text())
			if len(values) != 3 {
				return nil, fmt.Errorf("Invalid number of pixel values in PPM file. Got: %v, Expected: 3", len(values))
			}

			fmt.Printf("Line: %v, Values: %v\n", i*height+j, values)

			r, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, fmt.Errorf("Error converting red value: %v", err)
			}
			g, err := strconv.Atoi(values[1])
			if err != nil {
				return nil, fmt.Errorf("Error converting green value: %v", err)
			}
			b, err := strconv.Atoi(values[2])
			if err != nil {
				return nil, fmt.Errorf("Error converting blue value: %v", err)
			}

			data[i][j] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
		}
	}

	return &PPM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         maxVal,
	}, nil
}

// Size renvoie la largeur et la hauteur de l'image PPM.

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At renvoie la couleur du pixel aux coordonnées spécifiées.

func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set définit la couleur du pixel aux coordonnées spécifiées avec la valeur de couleur donnée.

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save enregistre l'image PPM dans un fichier.

func (ppm *PPM) Save(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			fmt.Fprintf(writer, "%d %d %d\n", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
		}
	}

	return writer.Flush()
}

func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width/2; j++ {
			ppm.data[i][j], ppm.data[i][ppm.width-j-1] = ppm.data[i][ppm.width-j-1], ppm.data[i][j]
		}
	}
}

func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ {
		ppm.data[i], ppm.data[ppm.height-i-1] = ppm.data[ppm.height-i-1], ppm.data[i]
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = int(maxValue)
}

func (ppm *PPM) Rotate90CW() {
	newData := make([][]Pixel, ppm.width)
	for i := range newData {
		newData[i] = make([]Pixel, ppm.height)
		for j := range newData[i] {
			newData[i][j] = ppm.data[ppm.height-j-1][i]
		}
	}
	ppm.data = newData
	ppm.width, ppm.height = ppm.height, ppm.width
}

type PGM struct {
	data        [][]uint8
	width       int
	height      int
	magicNumber string
	max         int
}

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func (ppm *PPM) ToPGM() *PGM {
	pgmData := make([][]uint8, ppm.height)
	for i := 0; i < ppm.height; i++ {
		pgmData[i] = make([]uint8, ppm.width)
		for j := 0; j < ppm.width; j++ {
			grayValue := uint8((uint32(ppm.data[i][j].R) + uint32(ppm.data[i][j].G) + uint32(ppm.data[i][j].B)) / 3)
			pgmData[i][j] = grayValue
		}
	}

	maxValue := 255

	return &PGM{
		data:        pgmData,
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         maxValue,
	}
}

func (ppm *PPM) ToPBM() *PBM {
	threshold := uint8(ppm.max) / 2

	pbmData := make([][]bool, ppm.height)
	for i := 0; i < ppm.height; i++ {
		pbmData[i] = make([]bool, ppm.width)
		for j := 0; j < ppm.width; j++ {
			averageIntensity := (uint32(ppm.data[i][j].R) + uint32(ppm.data[i][j].G) + uint32(ppm.data[i][j].B)) / 3
			pbmData[i][j] = averageIntensity > uint32(threshold)
		}
	}

	return &PBM{
		data:        pbmData,
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}
}

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	x0, y0 := p1.X, p1.Y
	x1, y1 := p2.X, p2.Y

	dx := abs(x1 - x0)
	sx := x0
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	dy := -abs(y1 - y0)
	sy := y0
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx + dy
	e2 := 0

	for {
		ppm.Set(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 = 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}

}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	ppm.DrawLine(p1, Point{X: p1.X + width, Y: p1.Y}, color)
	ppm.DrawLine(Point{X: p1.X, Y: p1.Y + height}, Point{X: p1.X + width, Y: p1.Y + height}, color)
	ppm.DrawLine(p1, Point{X: p1.X, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y}, Point{X: p1.X + width, Y: p1.Y + height}, color)
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ppm.Set(p1.X+x, p1.Y+y, color)
		}
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	x := radius
	y := 0
	err := 0

	for x >= y {
		ppm.Set(center.X+x, center.Y+y, color)
		ppm.Set(center.X+y, center.Y+x, color)
		ppm.Set(center.X-y, center.Y+x, color)
		ppm.Set(center.X-x, center.Y+y, color)
		ppm.Set(center.X-x, center.Y-y, color)
		ppm.Set(center.X-y, center.Y-x, color)
		ppm.Set(center.X+y, center.Y-x, color)
		ppm.Set(center.X+x, center.Y-y, color)

		y += 1
		if err <= 0 {
			err += 2*y + 1
		} else {
			x -= 1
			err += 2*(y-x) + 1
		}
	}
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	x := radius
	y := 0
	err := 0

	for x >= y {
		ppm.fillScanline(center.X-y, center.X+y, center.Y+x, color)
		ppm.fillScanline(center.X-y, center.X+y, center.Y-x, color)
		ppm.fillScanline(center.X-x, center.X+x, center.Y+y, color)
		ppm.fillScanline(center.X-x, center.X+x, center.Y-y, color)

		y += 1
		if err <= 0 {
			err += 2*y + 1
		} else {
			x -= 1
			err += 2*(y-x) + 1
		}
	}
}

func (ppm *PPM) fillScanline(x1, x2, y int, color Pixel) {
	if y < 0 || y >= ppm.height {
		return
	}

	if x1 > x2 {
		x1, x2 = x2, x1
	}

	for x := x1; x <= x2; x++ {
		if x >= 0 && x < ppm.width {
			ppm.Set(x, y, color)
		}
	}
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	vertices := []Point{p1, p2, p3}
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].Y < vertices[j].Y
	})

	slope1 := float64(vertices[1].X-vertices[0].X) / float64(vertices[1].Y-vertices[0].Y)
	slope2 := float64(vertices[2].X-vertices[0].X) / float64(vertices[2].Y-vertices[0].Y)
	slope3 := float64(vertices[2].X-vertices[1].X) / float64(vertices[2].Y-vertices[1].Y)

	x1 := float64(vertices[0].X)
	x2 := float64(vertices[0].X)

	for y := vertices[0].Y; y <= vertices[1].Y; y++ {
		ppm.fillScanline(int(x1), int(x2), y, color)
		x1 += slope1
		x2 += slope2
	}

	x1 = float64(vertices[1].X)
	x2 = float64(vertices[0].X)

	for y := vertices[1].Y + 1; y <= vertices[2].Y; y++ {
		ppm.fillScanline(int(x1), int(x2), y, color)
		x1 += slope3
		x2 += slope2
	}
}

func (ppm *PPM) calculateSlope(p1, p2 Point) float64 {
	if p2.Y-p1.Y == 0 {
		return 0
	}
	return float64(p2.X-p1.X) / float64(p2.Y-p1.Y)
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	vertices := []Point{p1, p2, p3}
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].Y < vertices[j].Y
	})

	slope1 := ppm.calculateSlope(vertices[0], vertices[1])
	slope2 := ppm.calculateSlope(vertices[0], vertices[2])
	slope3 := ppm.calculateSlope(vertices[1], vertices[2])

	x1 := float64(vertices[0].X)
	x2 := float64(vertices[0].X)

	for y := vertices[0].Y; y <= vertices[1].Y; y++ {
		ppm.fillScanline(int(x1), int(x2), y, color)
		x1 += slope1
		x2 += slope2
	}

	x1 = float64(vertices[1].X)
	x2 = float64(vertices[0].X)

	for y := vertices[1].Y + 1; y <= vertices[2].Y; y++ {
		ppm.fillScanline(int(x1), int(x2), y, color)
		x1 += slope3
		x2 += slope2
	}
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	numPoints := len(points)

	for i := 0; i < numPoints-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}

	ppm.DrawLine(points[numPoints-1], points[0], color)
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	numPoints := len(points)

	minY, maxY := points[0].Y, points[0].Y
	for i := 1; i < numPoints; i++ {
		if points[i].Y < minY {
			minY = points[i].Y
		} else if points[i].Y > maxY {
			maxY = points[i].Y
		}
	}

	for y := minY; y <= maxY; y++ {
		intersectPoints := []Point{}
		for i := 0; i < numPoints; i++ {
			if (points[i].Y <= y && points[(i+1)%numPoints].Y > y) ||
				(points[i].Y > y && points[(i+1)%numPoints].Y <= y) {
				x := int(float64(points[i].X) + float64(y-points[i].Y)/float64(points[(i+1)%numPoints].Y-points[i].Y)*float64(points[(i+1)%numPoints].X-points[i].X))
				intersectPoints = append(intersectPoints, Point{X: x, Y: y})
			}
		}

		sort.Slice(intersectPoints, func(i, j int) bool {
			return intersectPoints[i].X < intersectPoints[j].X
		})

		for i := 0; i < len(intersectPoints)-1; i += 2 {
			ppm.DrawLine(intersectPoints[i], intersectPoints[i+1], color)
		}
	}
}

func (ppm *PPM) DrawKochSnowflake(center Point, radius, depth int, color Pixel) {

	p1 := Point{center.X, center.Y - radius}
	cos30 := math.Cos(math.Pi / 6)
	sin30 := math.Sin(math.Pi / 6)
	p2 := Point{center.X + int(float64(radius)*cos30), center.Y + int(float64(radius)*sin30)}
	p3 := Point{center.X - int(float64(radius)*cos30), center.Y + int(float64(radius)*sin30)}

	ppm.drawKochSegment(p1, p2, depth, color)
	ppm.drawKochSegment(p2, p3, depth, color)
	ppm.drawKochSegment(p3, p1, depth, color)
}

func (ppm *PPM) drawKochSegment(p1, p2 Point, depth int, color Pixel) {
	if depth == 0 {
		ppm.DrawLine(p1, p2, color)
	} else {

		p3 := calculateKochIntermediate(p1, p2, 1.0/3.0)
		p4 := calculateKochIntermediate(p1, p2, 2.0/3.0)

		tip := calculateKochTip(p1, p2)

		ppm.drawKochSegment(p1, p3, depth-1, color)
		ppm.drawKochSegment(p3, tip, depth-1, color)
		ppm.drawKochSegment(tip, p4, depth-1, color)
		ppm.drawKochSegment(p4, p2, depth-1, color)
	}
}

func calculateKochIntermediate(p1, p2 Point, ratio float64) Point {
	return Point{
		X: p1.X + int(float64(p2.X-p1.X)*ratio),
		Y: p1.Y + int(float64(p2.Y-p1.Y)*ratio),
	}
}

func calculateKochTip(p1, p2 Point) Point {
	angle := math.Pi / 3.0
	x := int(float64(p1.X) + float64(p2.X-p1.X)*math.Cos(angle) - float64(p2.Y-p1.Y)*math.Sin(angle))
	y := int(float64(p1.Y) + float64(p2.X-p1.X)*math.Sin(angle) + float64(p2.Y-p1.Y)*math.Cos(angle))
	return Point{X: x, Y: y}
}

func cos(angle float64) {
	panic("unimplemented")
}

func sin(angle float64) {
	panic("unimplemented")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {

	ppm, err := ReadPPM("duck.ppm")
	if err != nil {
		fmt.Println("Error reading PPM file:", err)
		return
	}

	width, height := ppm.Size()
	fmt.Println("Image size:", width, height)
	ppm.Invert()
	err = ppm.Save("duck2.ppm")
	if err != nil {
		fmt.Println("Error saving PPM file:", err)
	}
}
