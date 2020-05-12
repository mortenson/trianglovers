package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

func handleStart() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	if isButtonColliding("Normal mode", 300, (height/2)+60) {
		gState.gameMode = modeIntro
	} else if isButtonColliding("Hard mode", 420, (height/2)+60) {
		gState.gameMode = modeIntro
		gState.maxQuestions = 4
	}
}

func drawTitle(screen *ebiten.Image) {
	title := "Trianglovers"
	text.Draw(screen, title, titleFont, (width/2)-(getTextWidth(title, titleFont)/2), (height/2)-45, defaultColors["purple"])
	drawButton(screen, "Normal mode", 300, (height/2)+60)
	drawButton(screen, "Hard mode", 420, (height/2)+60)
	credits := "by Sam and Mykal Mortenson"
	text.Draw(screen, credits, defaultFont, (width/2)-(getTextWidth(credits, defaultFont)/2), height-20, defaultColors["purple"])
}
