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
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/vector"
	"golang.org/x/image/font"
)

const (
	width  = 800
	height = 600
)

type vertex [2]int

type trianglover struct {
	name        string
	headPoint   int
	points      [3]int
	guessPoints [3]int
}

type dragPoint struct {
	position vertex
	dragging bool
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
var dragPoints [3]*dragPoint
var dragTargets []vertex
var currentLover *trianglover

type answer struct {
	ID     string
	ranges [][2]int
}

type question struct {
	ID      string
	answers []answer
}

var questions []question

func drawPolygon(screen *ebiten.Image, clr color.Color, coordinates []vertex) {
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

func distance(p1, p2 vertex) float64 {
	first := math.Pow(float64(p2[0]-p1[0]), 2)
	second := math.Pow(float64(p2[1]-p1[1]), 2)
	return math.Sqrt(first + second)
}

func angle(p1, p2, p3 vertex) float64 {
	radians := math.Atan2(float64(p3[1]-p1[1]), float64(p3[0]-p1[0])) - math.Atan2(float64(p2[1]-p1[1]), float64(p2[0]-p1[0]))
	degrees := radians * 180 / math.Pi
	if degrees > 0 {
		return degrees
	}
	return 360 + degrees
}

func getHexPoints(x, y int) []vertex {
	return []vertex{
		{50 + x, 0 + y},
		{100 + x, 25 + y},
		{100 + x, 75 + y},
		{50 + x, 100 + y},
		{0 + x, 75 + y},
		{0 + x, 25 + y},
	}
}

func getHexBoundaryPoints(hexPoints []vertex) []vertex {
	// Segment hexagon boundary into 102 points.
	points := make([]vertex, 0)
	for i := 0; i < len(hexPoints); i++ {
		nextIndex := i + 1
		if nextIndex >= len(hexPoints) {
			nextIndex = 0
		}
		xDelta := float64(hexPoints[nextIndex][0]-hexPoints[i][0]) / 17
		yDelta := float64(hexPoints[nextIndex][1]-hexPoints[i][1]) / 17
		for j := 0.0; j < 17; j++ {
			points = append(points, vertex{
				hexPoints[i][0] + int(j*xDelta),
				hexPoints[i][1] + int(j*yDelta),
			})
		}
	}
	return points
}

func drawMatchChart(screen *ebiten.Image, x, y int, prefPoints [3]int) {
	// Draw hexagon.
	hexPoints := getHexPoints(x, y)
	drawPolygon(screen, color.RGBA{255, 0, 0, 255}, hexPoints)
	points := getHexBoundaryPoints(hexPoints)
	// Draw the triangle.
	drawPolygon(screen, color.White, []vertex{
		points[prefPoints[0]],
		points[prefPoints[1]],
		points[prefPoints[2]],
	})
	// Register drag points (targets)
	dragTargets = points
	// Register drag points.
	dragPoints[0].position = points[prefPoints[0]]
	dragPoints[1].position = points[prefPoints[1]]
	dragPoints[2].position = points[prefPoints[2]]
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
	// Add angles (@todo).
	a := points[prefPoints[0]]
	b := points[prefPoints[1]]
	c := points[prefPoints[2]]
	A := angle(a, b, c)
	B := angle(b, c, a)
	C := angle(c, a, b)
	text.Draw(screen, fmt.Sprintf("%+v %+v %+v", points[prefPoints[0]], points[prefPoints[1]], points[prefPoints[2]]), defaultFont, 12, 12, color.White)
	text.Draw(screen, fmt.Sprintf("A: %f B: %f C: %f", A, B, C), defaultFont, 12, 24, color.White)
}

func handleDrag() {
	mouseX, mouseY := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, point := range dragPoints {
			point.dragging = false
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, point := range dragPoints {
			if distance(point.position, vertex{mouseX, mouseY}) < 10 {
				point.dragging = true
				break
			}
		}
	}
	for pointID, point := range dragPoints {
		if !point.dragging {
			continue
		}
		min := 0
		max := 34
		switch pointID {
		case 1:
			min = 34
			max = 68
		case 2:
			min = 68
			max = 102
		}
		closestID := -1
		closestDistance := 100.0
		for i := min; i < max; i++ {
			d := distance(dragTargets[i], vertex{mouseX, mouseY})
			if d < closestDistance {
				closestID = i
				closestDistance = d
			}
		}
		if closestID == -1 {
			continue
		}
		currentLover.points[pointID] = closestID
	}
}

func drawTrianglover(screen *ebiten.Image, lover *trianglover) {
	points := getHexBoundaryPoints(getHexPoints(100, 300))
	vertices := []vertex{points[lover.points[0]], points[lover.points[1]], points[lover.points[2]]}
	// Rotate triangle so that the head is facing up.
	for i := range vertices {
		if i == lover.headPoint {
			vertices[i] = points[0]
			continue
		}
		diff := lover.points[i] - lover.points[lover.headPoint]
		if diff < 0 {
			diff += len(points)
		}
		vertices[i] = points[diff]
	}
	// Scale triangle up.
	scaleX := vertices[0][0]
	scaleY := vertices[0][1]
	scale := 2
	for i := range vertices {
		vertices[i][0] = scaleX + ((vertices[i][0] - scaleX) * scale)
		vertices[i][1] = scaleY + ((vertices[i][1] - scaleY) * scale)
	}
	drawPolygon(screen, color.White, vertices)
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
	dragTargets = []vertex{}
	dragPoints = [3]*dragPoint{{}, {}, {}}
	hexPoints := map[string]int{
		"COMFORT":    0,
		"WEALTH":     17,
		"ADVENTURE":  34,
		"EXCITEMENT": 51,
		"ROMANCE":    68,
		"FAMILY":     85,
	}
	normalizePoint := func(p int) int {
		if p < 0 {
			return p + 102
		}
		return p
	}
	questions = []question{}
	for label, hexPoint := range hexPoints {
		questions = append(questions, question{
			ID: label + "_A",
			answers: []answer{
				{
					ID: label + "_A_STRONG",
					ranges: [][2]int{
						{normalizePoint(hexPoint - 8), hexPoint},
						{hexPoint, hexPoint + 8},
					},
				},
				{
					ID: label + "_A_NORMAL",
					ranges: [][2]int{
						{normalizePoint(hexPoint - 17), normalizePoint(hexPoint - 8)},
						{hexPoint + 8, hexPoint + 17},
					},
				},
				{
					ID: label + "_A_AGAINST",
					ranges: [][2]int{
						{hexPoint + 34, hexPoint + 68},
					},
				},
			},
		})
	}
}

func update(screen *ebiten.Image) error {
	ebiten.SetScreenScale(1.5)
	if trianglovers == nil {
		trianglovers = make([]*trianglover, 0)
		for i := 0; i < 10; i++ {
			trianglovers = append(trianglovers, &trianglover{
				name: fmt.Sprintf("%d", i),
				points: [3]int{
					rand.Intn(34),
					rand.Intn(34) + 34,
					rand.Intn(34) + 68,
				},
				headPoint: rand.Intn(3),
			})
		}
		currentLover = trianglovers[0]
	}

	handleDrag()
	drawMatchChart(screen, width-180, height-120, currentLover.points)
	drawTrianglover(screen, currentLover)

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := ebiten.Run(update, width, height, 2, "Trianglovers"); err != nil {
		panic(err)
	}
}
