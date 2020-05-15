package main

import (
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

func handleAudioToggle() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && isMouseColliding(width-60, 0, 60, 14) {
		audioToggle = !audioToggle
	}

	// This fixes an odd bug where the first time the mp3 is played, a horrible
	// crackling noise starts.
	bgm := audioFiles["bgm.mp3"]
	if time.Now().After(audioBufferUntil) {
		bgm.SetVolume(1)
	} else {
		bgm.SetVolume(0)
	}

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
	clr := defaultColors["purple"]
	if isMouseColliding(width-60, 0, 60, 14) {
		clr = defaultColors["darkPink"]
	}
	text.Draw(screen, prompt, defaultFont, width-60, 14, clr)
}
