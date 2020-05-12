package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

func drawResult(screen *ebiten.Image) {
	score := 0
	colormap := map[int]color.Color{}
	for _, m := range gState.matches {
		if gState.trianglovers[m.a].points == gState.trianglovers[m.b].points {
			score++
			colormap[m.a] = color.RGBA{104, 211, 116, 255}
			colormap[m.b] = colormap[m.a]
		} else {
			colormap[m.a] = color.RGBA{223, 90, 117, 255}
			colormap[m.b] = colormap[m.a]
		}
	}
	for i, lover := range gState.trianglovers {
		x := ((i % 5) * 125) + 100
		y := (int(i/5) * 125) + 25
		clr, ok := colormap[i]
		if !ok {
			clr = defaultColors["darkPink"]
		}
		drawMatchChart(screen, x, y, lover.points, false, clr)
		hexPoints := getHexPoints(x, y)
		points := getHexBoundaryPoints(hexPoints)
		drawPolygon(screen, color.RGBA{255, 255, 255, 50}, []vertex{
			points[lover.guessPoints[0]],
			points[lover.guessPoints[1]],
			points[lover.guessPoints[2]],
		})
		text.Draw(screen, lover.name, defaultFont, x+50-(getTextWidth(lover.name, defaultFont)/2), y+115, defaultColors["purple"])
	}
	title := fmt.Sprintf("Correct matches: %d/%d", score, len(gState.matches))
	text.Draw(screen, title, largeFont, 400-(getTextWidth(title, largeFont)/2), 350, defaultColors["purple"])
	title = "Thanks for playing!"
	text.Draw(screen, title, titleFont, 400-(getTextWidth(title, titleFont)/2), 425, defaultColors["purple"])
}
