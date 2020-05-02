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

func drawPolygon(screen *ebiten.Image, clr color.Color, coordinates ...int) {
	image, _ := ebiten.NewImage(16, 16, ebiten.FilterDefault)
	image.Fill(clr)
	triopts := &ebiten.DrawTrianglesOptions{}
	vertices := make([]ebiten.Vertex, 0)
	indices := make([]uint16, 0)
	for i := 0; i < len(coordinates)-1; i += 2 {
		vertices = append(vertices, ebiten.Vertex{
			DstX:   float32(coordinates[i]),
			DstY:   float32(coordinates[i+1]),
			SrcX:   0,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		})
		indices = append(indices, uint16(i/2))
		if i > 0 {
			indices = append(indices, uint16(len(coordinates)/2))
			indices = append(indices, uint16(i/2))
		}
	}
	indices = append(indices, 0)
	indices = append(indices, uint16(len(coordinates)/2))

	totalX := 0
	totalY := 0
	for i := 0; i < len(coordinates)-1; i += 2 {
		totalX += coordinates[i]
		totalY += coordinates[i+1]
	}
	centerX := totalX / (len(coordinates) / 2)
	centerY := totalY / (len(coordinates) / 2)
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

	drawPolygon(screen, color.RGBA{255, 0, 0, 255}, 50, 0, 100, 50, 100, 100, 50, 150, 0, 100, 0, 50)
	//drawPolygon(screen, color.RGBA{255, 0, 0, 255}, 100, 100, 150, 150, 50, 150)

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := ebiten.Run(update, width, height, 2, "Trianglovers"); err != nil {
		panic(err)
	}
}
