package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

// Handles the behavior of matching lovers, which is complex.
func handleMatch() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	for i := range gState.trianglovers {
		x := ((i % 5) * 125) + 100
		y := (int(i/5) * 125) + 100
		if !isMouseColliding(x, y, 100, 100) {
			continue
		}
		// If a match already exists for this lover, remove it.
		for j, m := range gState.matches {
			if m.a == i || m.b == i {
				gState.matches = append(gState.matches[:j], gState.matches[j+1:]...)
				break
			}
		}
		if gState.lastMatch == i {
			// If a lover is clicked twice, cancel the match.
			gState.lastMatch = -1
		} else if gState.lastMatch != -1 {
			// Add a new match.
			gState.matches = append(gState.matches, match{
				a:     gState.lastMatch,
				b:     i,
				color: gState.lastMatchColor,
			})
			gState.lastMatch = -1
		} else {
			// Start a new match with a random color.
			gState.lastMatch = i
			gState.lastMatchColor = color.RGBA{uint8(55 + rand.Intn(200)), uint8(55 + rand.Intn(200)), uint8(55 + rand.Intn(200)), 255}
		}
	}
	// Returns players to the interview phase.
	if isButtonColliding("Go back", 380, 550) {
		gState.gameMode = modeGuess
	}
	// Submits matches, if all matches are made.
	if isButtonColliding("Submit matches!", 350, 400) {
		if len(gState.matches) == len(gState.trianglovers)/2 {
			gState.gameMode = modeResult
		}
	}
}

// Draws the match page.
func drawMatchPage(screen *ebiten.Image) {
	title := "Match the lovers"
	text.Draw(screen, title, largeFont, 400-(getTextWidth(title, largeFont)/2), 75, defaultColors["purple"])
	colormap := map[int]color.Color{}
	for _, m := range gState.matches {
		colormap[m.a] = m.color
		colormap[m.b] = m.color
	}
	for i, lover := range gState.trianglovers {
		x := ((i % 5) * 125) + 100
		y := (int(i/5) * 125) + 100
		clr, ok := colormap[i]
		if !ok {
			if isMouseColliding(x, y, 100, 100) {
				clr = defaultColors["pink"]
			} else {
				clr = defaultColors["darkPink"]
			}
		}
		if gState.lastMatch == i {
			clr = gState.lastMatchColor
		}
		drawMatchChart(screen, x, y, lover.guessPoints, false, clr)
		text.Draw(screen, lover.name, defaultFont, x+50-(getTextWidth(lover.name, defaultFont)/2), y+115, defaultColors["purple"])
	}
	drawButton(screen, "Go back", 380, 550)
	if len(gState.matches) == len(gState.trianglovers)/2 {
		drawButton(screen, "Submit matches!", 350, 400)
	}
}
