package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

// Handles the start button being cicked.
func handleIntro() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	if isButtonColliding("Start matching", 350, 450) {
		gState.gameMode = modeGuess
	}
}

// Draws intro text describing how the game works.
func drawIntro(screen *ebiten.Image) {
	title := "Welcome to your new job!"
	text.Draw(screen, title, titleFont, 400-(getTextWidth(title, titleFont)/2), 100, defaultColors["purple"])
	intro := `You've been hired as a matchmaker for 10 eligible Trianglovers!

Your task is to interview each Trianglover, changing their match
chart based on their answers. To change the match chart, drag its
points closer to topics the Lover is pasionate about.

A Lover's true match chart is the same shape as them, but may not
appear at the same rotation. Use the Lover's shape as a clue for
what their final match chart should look like.

When you've finished interviewing, it's time to make your matches.
Lovers are a match when their match charts are exactly the same.
If you've done your job well five happy couples should be matched.`
	text.Draw(screen, intro, defaultFont, 400-(getTextWidth(intro, defaultFont)/2), 150, defaultColors["purple"])
	title = "Good luck!"
	text.Draw(screen, title, largeFont, 400-(getTextWidth(title, largeFont)/2), 425, defaultColors["purple"])
	drawButton(screen, "Start matching", 350, 450)
}
