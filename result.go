package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

// Handles the game being replayed.
func handleResult() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	if isButtonColliding("Replay", 385, 550) {
		rand.Seed(time.Now().UnixNano())
		gState = newGameState()
		gState.gameMode = modeTitle
	}
}

// Draws the result page.
func drawResult(screen *ebiten.Image) {
	score := 0
	// Map lover indexes to whether or not the match succeeded.
	successMap := map[int]bool{}
	for _, m := range gState.matches {
		if gState.trianglovers[m.a].points == gState.trianglovers[m.b].points {
			score++
			successMap[m.a] = true
			successMap[m.b] = successMap[m.a]
		} else {
			successMap[m.a] = false
			successMap[m.b] = successMap[m.a]
		}
	}
	// Draw grid of match charts indicating success.
	for i, lover := range gState.trianglovers {
		x := ((i % 5) * 125) + 100
		y := (int(i/5) * 125) + 25
		success, ok := successMap[i]
		var clr color.Color
		if ok && success {
			clr = color.RGBA{104, 211, 116, 255}
		} else {
			clr = color.RGBA{223, 90, 117, 255}
		}
		drawMatchChart(screen, x, y, lover.points, false, clr)
		hexPoints := getHexPoints(x, y)
		points := getHexBoundaryPoints(hexPoints)
		// Overlay an opaque triangle representing the user's guess.
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
	drawButton(screen, "Replay", 385, 550)
}
