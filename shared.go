package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Draws a block color polygon.
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

// Draws an outlined polygon.
// Note that to account for the function not being able to properly stroke
// asymmetrical polygons, if width is less than 1 the stroke will proportional
// to the shape. So a width of ".8" would the stroke is 20% wide. This was
// added last minute and would ideally be removed so triangles could be evenly
// stroked.
func drawPolygonLine(screen *ebiten.Image, width float64, borderColor color.Color, fillColor color.Color, coordinates []vertex) {
	drawPolygon(screen, borderColor, coordinates)
	centerX := 0
	centerY := 0
	for i := range coordinates {
		centerX += coordinates[i][0]
		centerY += coordinates[i][1]
	}
	centerX = centerX / len(coordinates)
	centerY = centerY / len(coordinates)
	for i := range coordinates {
		if width < 1 {
			coordinates[i][0] = centerX + int(float64(coordinates[i][0]-centerX)*width)
			coordinates[i][1] = centerY + int(float64(coordinates[i][1]-centerY)*width)
			continue
		}
		offsetX := width
		offsetY := width
		if coordinates[i][0] > centerX {
			offsetX = offsetX * -1
		}
		if coordinates[i][1] > centerY {
			offsetY = offsetY * -1
		}
		coordinates[i][0] = centerX + int(float64(coordinates[i][0]-centerX)+offsetX)
		coordinates[i][1] = centerY + int(float64(coordinates[i][1]-centerY)+offsetY)
	}
	drawPolygon(screen, fillColor, coordinates)
}

// Calculates the distance between two vertexes.
func distance(p1, p2 vertex) float64 {
	first := math.Pow(float64(p2[0]-p1[0]), 2)
	second := math.Pow(float64(p2[1]-p1[1]), 2)
	return math.Sqrt(first + second)
}

// Used to calculate the match chart's points.
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

// Gets the boundary points of the match chart. These are used to draw the
// triangle on the match chart and handle drag and drop.
func getHexBoundaryPoints(hexPoints []vertex) []vertex {
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

// Checks if the mouse is colliding with a collision box.
func isMouseColliding(x, y, width, height int) bool {
	mouseX, mouseY := ebiten.CursorPosition()
	return mouseX >= x && mouseX <= x+width && mouseY >= y && mouseY <= y+height
}

// Gets the vertexes used to draw a button.
func getButtonBounds(buttonText string, x, y int) []vertex {
	width := getTextWidth(buttonText, defaultFont)
	return []vertex{{x, y}, {x + width + 20, y}, {x + width + 20, y + 40}, {x, y + 40}}
}

// Checks if the mouse is colliding with a button.
func isButtonColliding(buttonText string, x, y int) bool {
	bounds := getButtonBounds(buttonText, x, y)
	return isMouseColliding(x, y, bounds[2][0]-x, bounds[2][1]-y)
}

// Draws a button with the given text.
func drawButton(screen *ebiten.Image, buttonText string, x, y int) {
	drawPolygonLine(screen, 2, defaultColors["darkPink"], defaultColors["white"], getButtonBounds(buttonText, x, y))
	clr := defaultColors["purple"]
	if isButtonColliding(buttonText, x, y) {
		clr = defaultColors["darkPink"]
	}
	text.Draw(screen, buttonText, defaultFont, x+10, y+25, clr)
}

// Calculates text width - very useful for centering text.
func getTextWidth(text string, face font.Face) int {
	width := fixed.I(0)
	prevR := rune(-1)
	largestWidth := fixed.I(0)
	for _, r := range []rune(text) {
		if r == '\n' {
			width = 0
			prevR = rune(-1)
			continue
		}
		if prevR >= 0 {
			width += face.Kern(prevR, r)
		}
		a, ok := face.GlyphAdvance(r)
		if !ok {
			panic("Unable to determine glyph width")
		}
		width += a
		prevR = r
		if width > largestWidth {
			largestWidth = width
		}
	}
	return largestWidth.Round()
}
