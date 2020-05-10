package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"path/filepath"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/vector"
	"golang.org/x/image/font"
)

const (
	width  = 800
	height = 600
)

type vertex [2]int

type trianglover struct {
	name        string
	headPoint   int
	points      [3]int
	guessPoints [3]int
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
	modeGuess
	modeMatch
	modeResult
)

var trianglovers []*trianglover
var defaultFont font.Face
var largeFont font.Face
var titleFont font.Face
var hexLabels = []string{
	"Comfort",
	"Wealth",
	"Adventure",
	"Excitement",
	"Romance",
	"Family",
}
var dragPoints [3]*dragPoint
var dragTargets []vertex
var currentLover *trianglover
var currentLoverIndex int
var hoverQuestion int
var currentQuestion int
var questions []question
var gameMode gameModeType
var matches []match
var lastMatch int
var lastMatchColor color.Color
var strings map[string]string
var files map[string][]byte
var fontFiles map[string]*truetype.Font
var imageFiles map[string]*ebiten.Image
var defaultColors map[string]color.Color

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

func distance(p1, p2 vertex) float64 {
	first := math.Pow(float64(p2[0]-p1[0]), 2)
	second := math.Pow(float64(p2[1]-p1[1]), 2)
	return math.Sqrt(first + second)
}

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

func getHexBoundaryPoints(hexPoints []vertex) []vertex {
	// Segment hexagon boundary into 102 points.
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

func drawMatchChart(screen *ebiten.Image, x, y int, prefPoints [3]int, drawLabels bool, clr color.Color) {
	// Draw hexagon.
	hexPoints := getHexPoints(x, y)
	drawPolygon(screen, clr, hexPoints)
	points := getHexBoundaryPoints(hexPoints)
	// Draw the triangle.
	drawPolygon(screen, defaultColors["white"], []vertex{
		points[prefPoints[0]],
		points[prefPoints[1]],
		points[prefPoints[2]],
	})
	// Register drag points (targets)
	dragTargets = points
	// Register drag points.
	dragPoints[0].position = points[prefPoints[0]]
	dragPoints[1].position = points[prefPoints[1]]
	dragPoints[2].position = points[prefPoints[2]]
	// Draw hovered drag points.
	for _, point := range dragPoints {
		if point.hovering {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(point.position[0]-10), float64(point.position[1]-10))
			screen.DrawImage(imageFiles["dragcursor.png"], op)
		}
	}
	// Add labels.
	if !drawLabels {
		return
	}
	for i, hexLabel := range hexLabels {
		x := hexPoints[i][0]
		y := hexPoints[i][1]
		if i == 0 || i == 1 || i == 5 {
			y -= 7
		} else {
			y += 12
		}
		if i == 1 || i == 2 {
			x += 5
		} else if i == 4 || i == 5 {
			x -= 7*len(hexLabel) + 10
		} else {
			x -= (7 * len(hexLabel)) / 2
		}
		text.Draw(screen, hexLabel, defaultFont, x, y, defaultColors["purple"])
	}
}

func handleDrag() {
	mouseX, mouseY := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, point := range dragPoints {
			point.dragging = false
		}
	}
	for _, point := range dragPoints {
		if distance(point.position, vertex{mouseX, mouseY}) < 10 {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				point.dragging = true
			}
			point.hovering = true
		} else if !point.dragging {
			point.hovering = false
		}
	}
	for pointID, point := range dragPoints {
		if !point.dragging {
			continue
		}
		min := 0
		max := 34
		switch pointID {
		case 1:
			min = 34
			max = 68
		case 2:
			min = 68
			max = 102
		}
		closestID := -1
		closestDistance := 100.0
		for i := min; i < max; i++ {
			d := distance(dragTargets[i], vertex{mouseX, mouseY})
			if d < closestDistance {
				closestID = i
				closestDistance = d
			}
		}
		if closestID == -1 {
			continue
		}
		currentLover.guessPoints[pointID] = closestID
	}
}

var eyeR float64
var eyeDirection float64

func drawTrianglover(screen *ebiten.Image, lover *trianglover) {
	points := getHexBoundaryPoints(getHexPoints(100, 300))
	vertices := []vertex{points[lover.points[0]], points[lover.points[1]], points[lover.points[2]]}
	// Rotate triangle so that the head is facing up.
	for i := range vertices {
		if i == lover.headPoint {
			vertices[i] = points[0]
			continue
		}
		diff := lover.points[i] - lover.points[lover.headPoint]
		if diff < 0 {
			diff += len(points)
		}
		vertices[i] = points[diff]
	}
	// Scale triangle up.
	scaleX := vertices[0][0]
	scaleY := vertices[0][1]
	scale := 2
	lowestY := -1
	lowestX := -1
	for i := range vertices {
		vertices[i][0] = scaleX + ((vertices[i][0] - scaleX) * scale)
		vertices[i][1] = scaleY + ((vertices[i][1] - scaleY) * scale)
		if lowestY == -1 || vertices[i][1] < lowestY {
			lowestY = vertices[i][1]
		}
		if lowestX == -1 || vertices[i][0] < lowestX {
			lowestX = vertices[i][0]
		}
	}
	// Move triangles to consistent coordinates.
	yDiff := 150 - lowestY
	xDiff := 50 - lowestX
	highestX := -1
	for i := range vertices {
		vertices[i][1] += yDiff
		vertices[i][0] += xDiff
		if highestX == -1 || vertices[i][0] > highestX {
			highestX = vertices[i][0]
		}
	}
	drawPolygonLine(screen, .9, defaultColors["darkPink"], defaultColors["pink"], vertices)
	// Draw face.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25, 175)
	screen.DrawImage(imageFiles["eyeball.png"], op)
	if eyeR <= 0 || (eyeR < 1.5*math.Pi && rand.Intn(100) > 98) {
		eyeDirection = 1
	} else if eyeR >= 2*math.Pi || (eyeR >= 1.5*math.Pi && rand.Intn(100) > 98) {
		eyeDirection = -1
	}
	eyeR += eyeDirection * .02
	eyeWidth := 50
	eyeInnerWidth := 18
	middleOffset := float64((eyeWidth / 2) - (eyeInnerWidth / 2))
	xOffset := middleOffset + float64(eyeWidth/4)*math.Cos(eyeR)
	yOffset := middleOffset + float64(eyeWidth/4)*math.Sin(eyeR)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25+xOffset, 175+yOffset)
	screen.DrawImage(imageFiles["eyeinner.png"], op)
	op = &ebiten.DrawImageOptions{}
	eye2X := float64(highestX - 25)
	op.GeoM.Translate(eye2X, 175)
	screen.DrawImage(imageFiles["eyeball.png"], op)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(eye2X+xOffset, 175+yOffset)
	screen.DrawImage(imageFiles["eyeinner.png"], op)
	// Draw mouth.
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50+(float64(highestX-50)/2)-25, 250)
	screen.DrawImage(imageFiles["mouth.png"], op)
}

func drawQuestions(screen *ebiten.Image) {
	text.Draw(screen, "Ask a question", largeFont, 400, 70, defaultColors["purple"])
	drawPolygonLine(screen, 2, defaultColors["darkPink"], defaultColors["white"], []vertex{{400, 100}, {780, 100}, {780, 325}, {400, 325}})
	x := 410
	y := 122
	for i, q := range questions {
		var clr color.Color
		if currentQuestion == i {
			clr = defaultColors["darkPink"]
		} else if hoverQuestion == i {
			clr = defaultColors["darkPink"]
		} else {
			clr = defaultColors["purple"]
		}
		text.Draw(screen, strings[q.ID], defaultFont, x, y+(i*25), clr)
	}
}

func handleQuestions() {
	mouseX, mouseY := ebiten.CursorPosition()
	y := 122
	hoverQuestion = -1
	for i := range questions {
		qY := y + (i * 25)
		if mouseX >= 410 && mouseX <= 780 && mouseY <= qY && mouseY >= qY-12 {
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				currentQuestion = i
			} else {
				hoverQuestion = i
			}
			break
		}
	}
}

func handleNextPrevious() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	if isButtonColliding("Previous Lover", 400, 335) {
		if currentLoverIndex > 0 {
			currentLoverIndex--
			currentLover = trianglovers[currentLoverIndex]
			currentQuestion = -1
			hoverQuestion = -1
		}
	} else if isButtonColliding("Next Lover", 530, 335) {
		if currentLoverIndex < len(trianglovers)-1 {
			currentLoverIndex++
			currentLover = trianglovers[currentLoverIndex]
			currentQuestion = -1
			hoverQuestion = -1
		}
	} else if currentLoverIndex == len(trianglovers)-1 && isButtonColliding("Match!", 710, 335) {
		gameMode = modeMatch
	}
}

func getButtonBounds(buttonText string, x, y int) []vertex {
	width := getTextWidth(buttonText, defaultFont)
	return []vertex{{x, y}, {x + width + 20, y}, {x + width + 20, y + 40}, {x, y + 40}}
}

func isMouseColliding(x, y, width, height int) bool {
	mouseX, mouseY := ebiten.CursorPosition()
	return mouseX >= x && mouseX <= x+width && mouseY >= y && mouseY <= y+height
}

func isButtonColliding(buttonText string, x, y int) bool {
	bounds := getButtonBounds(buttonText, x, y)
	return isMouseColliding(x, y, bounds[2][0]-x, bounds[2][1]-y)
}

func drawButton(screen *ebiten.Image, buttonText string, x, y int) {
	drawPolygonLine(screen, 2, defaultColors["darkPink"], defaultColors["white"], getButtonBounds(buttonText, x, y))
	clr := defaultColors["purple"]
	if isButtonColliding(buttonText, x, y) {
		clr = defaultColors["darkPink"]
	}
	text.Draw(screen, buttonText, defaultFont, x+10, y+25, clr)
}

func drawNextPrevious(screen *ebiten.Image) {
	drawButton(screen, "Previous Lover", 400, 335)
	drawButton(screen, "Next Lover", 530, 335)
	if currentLoverIndex == len(trianglovers)-1 {
		drawButton(screen, "Match!", 710, 335)
	}
}

func drawAnswer(screen *ebiten.Image) {
	drawPolygonLine(screen, 2, defaultColors["darkPink"], defaultColors["white"], []vertex{{20, 500}, {500, 500}, {500, 580}, {20, 580}})
	if currentQuestion == -1 {
		return
	}
	chosenAnswer := -1
	for i, a := range questions[currentQuestion].answers {
		for _, r := range a.ranges {
			for _, p := range currentLover.points {
				if (chosenAnswer == -1 || i <= chosenAnswer) && p >= r[0] && p <= r[1] {
					chosenAnswer = i
				}
			}
		}
	}
	var answerID string
	if chosenAnswer != -1 {
		answerID = questions[currentQuestion].answers[chosenAnswer].ID
	} else {
		answerID = questions[currentQuestion].ID + "_DEFAULT"
	}
	text.Draw(screen, strings[answerID], defaultFont, 30, 510+12, defaultColors["purple"])
}

func handleStart() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	gameMode = modeGuess
}

func getTextWidth(text string, face font.Face) int {
	width := fixed.I(0)
	prevR := rune(-1)
	for _, r := range []rune(text) {
		if prevR >= 0 {
			width += face.Kern(prevR, r)
		}
		a, ok := face.GlyphAdvance(r)
		if !ok {
			panic("Unable to determine glyph width")
		}
		width += a
		prevR = r
	}
	return width.Round()
}

func drawTitle(screen *ebiten.Image) {
	title := "Trianglovers"
	text.Draw(screen, title, titleFont, (width/2)-(getTextWidth(title, titleFont)/2), (height/2)-45, defaultColors["purple"])
	button := "Click to start"
	text.Draw(screen, button, largeFont, (width/2)-(getTextWidth(button, largeFont)/2), (height/2)+60, defaultColors["purple"])
}

func handleMatch() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	for i := range trianglovers {
		x := ((i % 5) * 125) + 100
		y := (int(i/5) * 125) + 100
		if !isMouseColliding(x, y, 100, 100) {
			continue
		}
		for j, m := range matches {
			if m.a == i || m.b == i {
				matches = append(matches[:j], matches[j+1:]...)
				break
			}
		}
		if lastMatch == i {
			lastMatch = -1
		} else if lastMatch != -1 {
			matches = append(matches, match{
				a:     lastMatch,
				b:     i,
				color: lastMatchColor,
			})
			lastMatch = -1
		} else {
			lastMatch = i
			lastMatchColor = color.RGBA{uint8(55 + rand.Intn(200)), uint8(55 + rand.Intn(200)), uint8(55 + rand.Intn(200)), 255}
		}
	}
	if isButtonColliding("Submit matches!", 350, 400) {
		if len(matches) == len(trianglovers)/2 {
			gameMode = modeResult
		}
	}
}

func drawMatchPage(screen *ebiten.Image) {
	title := "Match the lovers"
	text.Draw(screen, title, largeFont, 400-(getTextWidth(title, largeFont)/2), 75, defaultColors["purple"])
	colormap := map[int]color.Color{}
	for _, m := range matches {
		colormap[m.a] = m.color
		colormap[m.b] = m.color
	}
	for i, lover := range trianglovers {
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
		if lastMatch == i {
			clr = lastMatchColor
		}
		drawMatchChart(screen, x, y, lover.guessPoints, false, clr)
		text.Draw(screen, lover.name, defaultFont, x+50-(getTextWidth(lover.name, defaultFont)/2), y+115, defaultColors["purple"])
	}
	if len(matches) == len(trianglovers)/2 {
		drawButton(screen, "Submit matches!", 350, 400)
	}
}

func drawResult(screen *ebiten.Image) {
	colors := []color.Color{
		color.RGBA{41, 237, 255, 255}, // blue
		color.RGBA{57, 255, 137, 255}, // green
		color.RGBA{255, 159, 44, 255}, // orange
		color.RGBA{255, 92, 113, 255}, // red
		color.RGBA{255, 92, 251, 255}, // purple
	}
	colormap := map[[3]int]color.Color{}
	for _, lover := range trianglovers {
		_, ok := colormap[lover.points]
		if ok {
			continue
		}
		colormap[lover.points] = colors[len(colormap)]
	}
	for i, lover := range trianglovers {
		x := ((i % 5) * 125) + 100
		y := (int(i/5) * 125) + 25
		clr, ok := colormap[lover.points]
		if !ok {
			clr = defaultColors["darkPink"]
		}
		drawMatchChart(screen, x, y, lover.points, false, clr)
		text.Draw(screen, lover.name, defaultFont, x+50-(getTextWidth(lover.name, defaultFont)/2), y+115, defaultColors["purple"])
	}
	score := 0
	for _, m := range matches {
		if trianglovers[m.a].points == trianglovers[m.b].points {
			score++
		}
	}
	title := fmt.Sprintf("Correct matches: %d/%d", score, len(matches))
	text.Draw(screen, title, largeFont, 400-(getTextWidth(title, largeFont)/2), 350, defaultColors["purple"])
	title = "Thanks for playing!"
	text.Draw(screen, title, titleFont, 400-(getTextWidth(title, titleFont)/2), 425, defaultColors["purple"])
}

func init() {
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
	dragTargets = []vertex{}
	dragPoints = [3]*dragPoint{{}, {}, {}}
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
						{hexPoint.point - 8, hexPoint.point - 1},
						{hexPoint.point, hexPoint.point + 8},
					}),
				},
				{
					ID: hexPoint.ID + "_A_NORMAL",
					ranges: fixRanges([][2]int{
						{hexPoint.point - 17, hexPoint.point - 8},
						{hexPoint.point + 8, hexPoint.point + 17},
					}),
				},
			},
		})
	}
	currentQuestion = -1
	hoverQuestion = -1
	trianglovers = make([]*trianglover, 0)
	guessPoints := [3]int{
		0,
		34,
		68,
	}
	for i := 0; i < 5; i++ {
		points := [3]int{
			rand.Intn(34),
			rand.Intn(34) + 34,
			rand.Intn(34) + 68,
		}
		headPoint := rand.Intn(3)
		trianglovers = append(trianglovers, &trianglover{
			name:        "Lover " + fmt.Sprintf("%d", i),
			points:      points,
			headPoint:   headPoint,
			guessPoints: guessPoints,
		})
		trianglovers = append(trianglovers, &trianglover{
			name:        "Lover " + fmt.Sprintf("%d", i),
			points:      points,
			headPoint:   (headPoint + 1) % 3,
			guessPoints: guessPoints,
		})
	}
	rand.Shuffle(len(trianglovers), func(i, j int) { trianglovers[i], trianglovers[j] = trianglovers[j], trianglovers[i] })
	currentLover = trianglovers[0]
	currentLoverIndex = 0
	gameMode = modeMatch
	lastMatch = -1
	matches = make([]match, 0)
	strings = getStrings()
	defaultColors = map[string]color.Color{
		"darkPink": color.RGBA{245, 128, 193, 255},
		"pink":     color.RGBA{255, 187, 225, 255},
		"purple":   color.RGBA{175, 58, 141, 255},
		"white":    color.RGBA{255, 241, 241, 255},
	}
	eyeR = 0
	eyeDirection = 1
}

func loadFiles() {
	files = make(map[string][]byte, 0)
	imageFiles = make(map[string]*ebiten.Image, 0)
	fontFiles = make(map[string]*truetype.Font, 0)
	packrBox := packr.NewBox("./assets")
	for _, f := range packrBox.List() {
		b, err := packrBox.Find(f)
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

	switch gameMode {
	case modeTitle:
		handleStart()
		drawTitle(screen)
	case modeGuess:
		handleDrag()
		handleQuestions()
		handleNextPrevious()
		drawMatchChart(screen, width-180, height-120, currentLover.guessPoints, true, defaultColors["darkPink"])
		drawTrianglover(screen, currentLover)
		drawQuestions(screen)
		drawAnswer(screen)
		drawNextPrevious(screen)
	case modeMatch:
		handleMatch()
		drawMatchPage(screen)
	case modeResult:
		drawResult(screen)
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := ebiten.Run(update, width, height, 2, "Trianglovers"); err != nil {
		panic(err)
	}
}
