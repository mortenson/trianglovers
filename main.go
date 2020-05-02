package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

const (
	width  = 800
	height = 600
)

type gameMode int

const (
	modeTitle gameMode = iota
	modeTutorial
	modeQuestion
)

type gameState struct {
	mode gameMode
}

type preferencePoints struct {
	a int
	b int
	c int
}

type trianglover struct {
	name   string
	points preferencePoints
}

var trianglovers []*trianglover

func drawPolygon(screen *ebiten.Image, clr color.Color, coordinates [][2]int) {
	image, _ := ebiten.NewImage(16, 16, ebiten.FilterDefault)
	image.Fill(clr)
	triopts := &ebiten.DrawTrianglesOptions{}
	vertices := make([]ebiten.Vertex, 0)
	indices := make([]uint16, 0)
	totalX := 0
	totalY := 0
	for i := 0; i < len(coordinates); i++ {
		totalX += coordinates[i][0]
		totalY += coordinates[i][1]
		vertices = append(vertices, ebiten.Vertex{
			DstX:   float32(coordinates[i][0]),
			DstY:   float32(coordinates[i][1]),
			SrcX:   0,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		})
		indices = append(indices, uint16(i))
		if i > 0 {
			indices = append(indices, uint16(len(coordinates)))
			indices = append(indices, uint16(i))
		}
	}
	indices = append(indices, 0)
	indices = append(indices, uint16(len(coordinates)))

	centerX := totalX / len(coordinates)
	centerY := totalY / len(coordinates)
	vertices = append(vertices, ebiten.Vertex{
		DstX:   float32(centerX),
		DstY:   float32(centerY),
		SrcX:   0,
		SrcY:   0,
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	})

	screen.DrawTriangles(vertices, indices, image, triopts)
}

func update(screen *ebiten.Image) error {
	if trianglovers == nil {
		trianglovers = make([]*trianglover, 0)
		for i := 0; i < 10; i++ {
			trianglovers = append(trianglovers, &trianglover{
				name: fmt.Sprintf("%d", i),
				points: preferencePoints{
					a: rand.Intn(33),
					b: rand.Intn(33) + 33,
					c: rand.Intn(33) + 66,
				},
			})
		}
	}

	hexPoints := [][2]int{
		{50, 0},
		{100, 25},
		{100, 75},
		{50, 100},
		{0, 75},
		{0, 25},
	}
	drawPolygon(screen, color.RGBA{255, 0, 0, 255}, hexPoints)
	a := 10
	b := 40
	c := 75
	points := make([][2]int, 0)
	for i := 0; i < len(hexPoints); i++ {
		nextIndex := i + 1
		if nextIndex >= len(hexPoints) {
			nextIndex = 0
		}
		xDelta := (hexPoints[nextIndex][0] - hexPoints[i][0]) / 17
		yDelta := (hexPoints[nextIndex][1] - hexPoints[i][1]) / 17
		for j := 0; j < 17; j++ {
			points = append(points, [2]int{
				hexPoints[i][0] + j*xDelta,
				hexPoints[i][1] + j*yDelta,
			})
		}
	}
	drawPolygon(screen, color.White, [][2]int{
		points[a],
		points[b],
		points[c],
	})

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := ebiten.Run(update, width, height, 2, "Trianglovers"); err != nil {
		panic(err)
	}
}
