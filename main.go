package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"path/filepath"
	"time"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"golang.org/x/image/font"
)

const (
	width  = 800
	height = 600
)

type vertex [2]int

type trianglover struct {
	name           string
	headPoint      int
	points         [3]int
	guessPoints    [3]int
	questionsAsked []int
	answerIndex    int
}

type dragPoint struct {
	position vertex
	dragging bool
	hovering bool
}

type answer struct {
	ID     string
	ranges [][2]int
}

type question struct {
	ID      string
	answers []answer
}

type match struct {
	a     int
	b     int
	color color.Color
}

type gameModeType int

const (
	modeTitle gameModeType = iota
	modeIntro
	modeGuess
	modeMatch
	modeResult
)

type gameState struct {
	trianglovers      []*trianglover
	dragPoints        [3]*dragPoint
	dragTargets       []vertex
	currentLover      *trianglover
	currentLoverIndex int
	hoverQuestion     int
	currentQuestion   int
	gameMode          gameModeType
	matches           []match
	lastMatch         int
	lastMatchColor    color.Color
	eyeR              float64
	eyeDirection      float64
	loverOffset       float64
	loverDirection    float64
	maxQuestions      int
}

var defaultFont font.Face
var largeFont font.Face
var titleFont font.Face
var questions []question
var strings map[string][]string
var files map[string][]byte
var fontFiles map[string]*truetype.Font
var imageFiles map[string]*ebiten.Image
var audioFiles map[string]*audio.Player
var defaultColors map[string]color.Color
var audioToggle bool
var gState *gameState

func newGameState() *gameState {
	s := gameState{
		trianglovers:      []*trianglover{},
		dragTargets:       []vertex{},
		dragPoints:        [3]*dragPoint{{}, {}, {}},
		currentQuestion:   -1,
		hoverQuestion:     -1,
		currentLoverIndex: 0,
		gameMode:          modeTitle,
		lastMatch:         -1,
		matches:           []match{},
		eyeR:              0,
		eyeDirection:      1,
		loverOffset:       0,
		loverDirection:    1,
		maxQuestions:      6,
	}
	return &s
}

func (s *gameState) genTrianglovers(hardMode bool) {
	defaultNames := []string{
		"Digree",
		"Acutie",
		"Equilaten",
		"Hypotenate",
		"Obtussey",
		"Isos",
		"Petagorean",
		"Scaley",
		"Vert",
		"Anglea",
	}
	rand.Shuffle(len(defaultNames), func(i, j int) { defaultNames[i], defaultNames[j] = defaultNames[j], defaultNames[i] })
	var points [3]int
	pointHistory := [][3]int{}
	guessPoints := [3]int{
		0,
		34,
		68,
	}
	max := 5
	if hardMode {
		max = 3
	}
	for i := 0; i < max; i++ {
		// Generate random points that are unique enough to make the game fair.
		for {
			goodPoints := true
			points = [3]int{
				rand.Intn(34),
				rand.Intn(34) + 34,
				rand.Intn(34) + 68,
			}
			for _, h := range pointHistory {
				if math.Abs(float64(h[0]-points[0])) <= 8 && math.Abs(float64(h[1]-points[1])) <= 8 && math.Abs(float64(h[2]-points[2])) <= 8 {
					goodPoints = false
					break
				}
			}
			if goodPoints {
				pointHistory = append(pointHistory, points)
				break
			}
		}
		headPoint := rand.Intn(3)
		s.trianglovers = append(s.trianglovers, &trianglover{
			name:           defaultNames[i],
			points:         points,
			headPoint:      headPoint,
			guessPoints:    guessPoints,
			questionsAsked: make([]int, 0),
			answerIndex:    rand.Intn(2),
		})
		s.trianglovers = append(s.trianglovers, &trianglover{
			name:           defaultNames[i+max],
			points:         points,
			headPoint:      (headPoint + 1) % 3,
			guessPoints:    guessPoints,
			questionsAsked: make([]int, 0),
			answerIndex:    rand.Intn(2),
		})
		// In hard mode, two pairs of lovers have the same points in different
		// rotatations. This makes guessing based on shape alone impossible.
		if hardMode && i < 2 {
			offset := rand.Intn(34)
			for j := range points {
				points[j] += offset + (j * 34)
				if points[j] > 102 {
					points[j] -= 102
				}
			}
			pointHistory = append(pointHistory, points)
			headPoint := rand.Intn(3)
			s.trianglovers = append(s.trianglovers, &trianglover{
				name:           defaultNames[i+6],
				points:         points,
				headPoint:      headPoint,
				guessPoints:    guessPoints,
				questionsAsked: make([]int, 0),
			})
			s.trianglovers = append(s.trianglovers, &trianglover{
				name:           defaultNames[i+8],
				points:         points,
				headPoint:      (headPoint + 1) % 3,
				guessPoints:    guessPoints,
				questionsAsked: make([]int, 0),
			})
		}
	}
	rand.Shuffle(len(s.trianglovers), func(i, j int) { s.trianglovers[i], s.trianglovers[j] = s.trianglovers[j], s.trianglovers[i] })
	s.currentLover = s.trianglovers[0]
}

func init() {
	rand.Seed(time.Now().UnixNano())
	// Initialize game satte.
	gState = newGameState()
	// Load static assets.
	loadFiles()
	defaultFont = truetype.NewFace(fontFiles["Archivo-SemiBold.ttf"], &truetype.Options{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	largeFont = truetype.NewFace(fontFiles["LobsterTwo-Italic.ttf"], &truetype.Options{
		Size:    45,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	titleFont = truetype.NewFace(fontFiles["LobsterTwo-Italic.ttf"], &truetype.Options{
		Size:    60,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	// Helper function to ensure question ranges roll over.
	fixRanges := func(ranges [][2]int) [][2]int {
		for i := range ranges {
			for j := range ranges[i] {
				if ranges[i][j] < 0 {
					ranges[i][j] += 102
				} else if ranges[i][j] > 102 {
					ranges[i][j] -= 102
				}
			}
		}
		return ranges
	}
	// Generate question ranges programmatically. Hypothetically this format
	// allows for more complex ranges and questions, but for now "STRONG"
	// questions are within 4 units of a point, "NORMAL" is within 8, and
	// "DEFAULT" is everything else.
	questions = []question{}
	hexPoints := []struct {
		ID    string
		point int
	}{
		{"COMFORT", 0},
		{"WEALTH", 17},
		{"ADVENTURE", 34},
		{"EXCITEMENT", 51},
		{"ROMANCE", 68},
		{"FAMILY", 85},
	}
	for _, hexPoint := range hexPoints {
		questions = append(questions, question{
			ID: hexPoint.ID + "_A",
			answers: []answer{
				{
					ID: hexPoint.ID + "_A_STRONG",
					ranges: fixRanges([][2]int{
						{hexPoint.point - 4, hexPoint.point - 1},
						{hexPoint.point, hexPoint.point + 4},
					}),
				},
				{
					ID: hexPoint.ID + "_A_NORMAL",
					ranges: fixRanges([][2]int{
						{hexPoint.point - 8, hexPoint.point - 4},
						{hexPoint.point + 4, hexPoint.point + 8},
					}),
				},
			},
		})
	}
	strings = getStrings()
	defaultColors = map[string]color.Color{
		"darkPink": color.RGBA{245, 128, 193, 255},
		"pink":     color.RGBA{255, 187, 225, 255},
		"purple":   color.RGBA{175, 58, 141, 255},
		"white":    color.RGBA{255, 241, 241, 255},
	}
	audioToggle = true
}

// Helper function to load all assets from packr, processing them early if
// possible. Supports PNG images, WAV music, and TTF fonts.
func loadFiles() {
	files = make(map[string][]byte, 0)
	imageFiles = make(map[string]*ebiten.Image, 0)
	fontFiles = make(map[string]*truetype.Font, 0)
	audioFiles = make(map[string]*audio.Player, 0)
	packrBox := packr.New("assets", "./assets")
	audioContext, err := audio.NewContext(44100)
	if err != nil {
		panic(err)
	}
	for _, f := range packrBox.List() {
		b, err := packrBox.Find(f)
		if err != nil {
			panic(err)
		}
		s, err := packrBox.Open(f)
		if err != nil {
			panic(err)
		}
		switch filepath.Ext(f) {
		case ".png":
			img, _, err := image.Decode(bytes.NewReader(b))
			if err != nil {
				panic(err)
			}
			eimg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
			if err != nil {
				panic(err)
			}
			imageFiles[f] = eimg
		case ".ttf":
			ttf, err := truetype.Parse(b)
			if err != nil {
				panic(err)
			}
			fontFiles[f] = ttf
		case ".wav":
			d, err := wav.Decode(audioContext, s)
			if err != nil {
				panic(err)
			}
			audioPlayer, err := audio.NewPlayer(audioContext, d)
			if err != nil {
				panic(err)
			}
			audioFiles[f] = audioPlayer
		case ".mp3":
			d, err := mp3.Decode(audioContext, s)
			if err != nil {
				panic(err)
			}
			audioPlayer, err := audio.NewPlayer(audioContext, d)
			if err != nil {
				panic(err)
			}
			audioFiles[f] = audioPlayer
		default:
			files[f] = b
		}
	}
}

func update(screen *ebiten.Image) error {
	ebiten.SetScreenScale(1.5)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(imageFiles["background.png"], op)

	handleAudioToggle()
	drawAudioToggle(screen)

	switch gState.gameMode {
	case modeTitle:
		handleStart()
		drawTitle(screen)
	case modeIntro:
		handleIntro()
		drawIntro(screen)
	case modeGuess:
		handleGuess()
		drawGuess(screen)
	case modeMatch:
		handleMatch()
		drawMatchPage(screen)
	case modeResult:
		handleResult()
		drawResult(screen)
	}

	return nil
}

func main() {
	if err := ebiten.Run(update, width, height, 2, "Trianglovers"); err != nil {
		panic(err)
	}
}
