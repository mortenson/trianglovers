package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

func handleAudioToggle() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && isMouseColliding(width-70, 0, 60, 14) {
		audioToggle = !audioToggle
	}

	bgm := audioFiles["bgm.mp3"]
	if !audioToggle && bgm.IsPlaying() {
		bgm.Pause()
	} else if audioToggle && !bgm.IsPlaying() {
		bgm.Rewind()
		bgm.Play()
	}
}

func drawAudioToggle(screen *ebiten.Image) {
	prompt := "audio on"
	if !audioToggle {
		prompt = "audio off"
	}
	text.Draw(screen, prompt, defaultFont, width-60, 14, defaultColors["purple"])
}
