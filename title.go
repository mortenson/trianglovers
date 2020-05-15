package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

// Handles the game starting.
func handleStart() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	if isButtonColliding("Normal mode", 300, (height/2)+60) {
		gState.gameMode = modeIntro
		gState.genTrianglovers(false)
	} else if isButtonColliding("Hard mode", 420, (height/2)+60) {
		gState.gameMode = modeIntro
		gState.maxQuestions = 4
		gState.genTrianglovers(true)
	}
}

// Draws the title screen.
func drawTitle(screen *ebiten.Image) {
	title := "Trianglovers"
	text.Draw(screen, title, titleFont, (width/2)-(getTextWidth(title, titleFont)/2), (height/2)-45, defaultColors["purple"])
	drawButton(screen, "Normal mode", 300, (height/2)+60)
	drawButton(screen, "Hard mode", 420, (height/2)+60)
	credits := "by Sam and Mykal Mortenson"
	text.Draw(screen, credits, defaultFont, (width/2)-(getTextWidth(credits, defaultFont)/2), height-20, defaultColors["purple"])
	help := "(f) fullscreen (s) window size"
	text.Draw(screen, help, defaultFont, width-getTextWidth(help, defaultFont)-5, height-5, defaultColors["purple"])
}
