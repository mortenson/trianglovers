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

func drawPolygon(screen *ebiten.Image, coordinates ...int) {
	image, _ := ebiten.NewImage(1, 1, ebiten.FilterDefault)
	image.Fill(color.White)
	triopts := &ebiten.DrawTrianglesOptions{}
	vertexes := make([]ebiten.Vertex, 0)
	indicies := make([]uint16, 0)
	for i := 0; i < len(coordinates)-1; i += 2 {
		vertexes = append(vertexes, ebiten.Vertex{
			DstX:   float32(coordinates[i]),
			DstY:   float32(coordinates[i+1]),
			SrcX:   0,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		})
		indicies = append(indicies, uint16(i/2))
	}
	screen.DrawTriangles(vertexes, indicies, image, triopts)
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

	drawPolygon(screen, 100, 100, 150, 150, 50, 150)

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := ebiten.Run(update, width, height, 2, "Trianglovers"); err != nil {
		panic(err)
	}
}
