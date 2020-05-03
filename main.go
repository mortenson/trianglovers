package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/vector"
	"golang.org/x/image/font"
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
	name        string
	points      preferencePoints
	guessPoints preferencePoints
}

var trianglovers []*trianglover
var defaultFont font.Face
var hexLabels = []string{
	"Comfort",
	"Wealth",
	"Adventure",
	"Excitement",
	"Romance",
	"Family",
}

func drawPolygon(screen *ebiten.Image, clr color.Color, coordinates [][2]int) {
	path := vector.Path{}
	path.MoveTo(float32(coordinates[0][0]), float32(coordinates[0][1]))
	for i := 1; i < len(coordinates); i++ {
		path.LineTo(float32(coordinates[i][0]), float32(coordinates[i][1]))
	}
	path.MoveTo(float32(coordinates[0][0]), float32(coordinates[0][1]))
	path.Fill(screen, &vector.FillOptions{
		Color: clr,
	})
}

func distance(p1, p2 [2]int) float64 {
	first := math.Pow(float64(p2[0]-p1[0]), 2)
	second := math.Pow(float64(p2[1]-p1[1]), 2)
	return math.Sqrt(first + second)
}

func angle(p1, p2, p3 [2]int) float64 {
	radians := math.Atan2(float64(p3[1]-p1[1]), float64(p3[0]-p1[0])) - math.Atan2(float64(p2[1]-p1[1]), float64(p2[0]-p1[0]))
	degrees := radians * 180 / math.Pi
	if degrees > 0 {
		return degrees
	} else {
		return 360 + degrees
	}
}

func drawMatchChart(screen *ebiten.Image, x, y int, prefPoints preferencePoints) {
	// Draw hexagon.
	hexPoints := [][2]int{
		{50 + x, 0 + y},
		{100 + x, 25 + y},
		{100 + x, 75 + y},
		{50 + x, 100 + y},
		{0 + x, 75 + y},
		{0 + x, 25 + y},
	}
	drawPolygon(screen, color.RGBA{255, 0, 0, 255}, hexPoints)
	// Segment hexagon boundary into 102 points.
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
	// Draw the triangle.
	drawPolygon(screen, color.White, [][2]int{
		points[prefPoints.a],
		points[prefPoints.b],
		points[prefPoints.c],
	})
	// Add labels.
	for i, hexLabel := range hexLabels {
		x := hexPoints[i][0]
		y := hexPoints[i][1]
		if i == 0 || i == 1 || i == 5 {
			y -= 7
		} else {
			y += 12
		}
		if i == 1 || i == 2 {
			x += 5
		} else if i == 4 || i == 5 {
			x -= 7*len(hexLabel) + 5
		} else {
			x -= (7 * len(hexLabel)) / 2
		}
		text.Draw(screen, hexLabel, defaultFont, x, y, color.White)
	}
	// Add angles.
	a := points[prefPoints.a]
	b := points[prefPoints.b]
	c := points[prefPoints.c]
	A := angle(a, b, c)
	B := angle(b, c, a)
	C := angle(c, a, b)
	text.Draw(screen, fmt.Sprintf("%+v %+v %+v", points[prefPoints.a], points[prefPoints.b], points[prefPoints.c]), defaultFont, 12, 12, color.White)
	text.Draw(screen, fmt.Sprintf("A: %f B: %f C: %f", A, B, C), defaultFont, 12, 24, color.White)
}

func drawTrianglover(screen *ebiten.Image, lover *trianglover) {

}

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	defaultFont = truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func update(screen *ebiten.Image) error {
	ebiten.SetScreenScale(1.5)
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

	drawMatchChart(screen, width-180, height-120, trianglovers[0].points)
	drawTrianglover(screen, trianglovers[0])

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := ebiten.Run(update, width, height, 2, "Trianglovers"); err != nil {
		panic(err)
	}
}
