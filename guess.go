package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

// Draws the next/previous/match buttons that sit below the questions.
func drawNextPrevious(screen *ebiten.Image) {
	if gState.currentLoverIndex > 0 {
		drawButton(screen, "Previous Lover", 400, 335)
	}
	if gState.currentLoverIndex < len(gState.trianglovers)-1 {
		drawButton(screen, "Next Lover", 530, 335)
	}
	if gState.currentLoverIndex == len(gState.trianglovers)-1 {
		drawButton(screen, "Match!", 710, 335)
	}
	count := len(gState.trianglovers) - gState.currentLoverIndex - 1
	if count > 0 {
		text.Draw(screen, fmt.Sprintf("%d remaining", count), defaultFont, 630, 360, defaultColors["purple"])
	}
}

// Draws an answer based on the current question.
func drawAnswer(screen *ebiten.Image) {
	drawPolygonLine(screen, 2, defaultColors["darkPink"], defaultColors["white"], []vertex{{20, 500}, {500, 500}, {500, 580}, {20, 580}})
	if gState.currentQuestion == -1 {
		return
	}
	// Determine the correct answer by checking the answer ranges against the
	// triangle's points. The answer with the lowest weight is preferred.
	chosenAnswer := -1
	for i, a := range questions[gState.currentQuestion].answers {
		for _, r := range a.ranges {
			for _, p := range gState.currentLover.points {
				if (chosenAnswer == -1 || i <= chosenAnswer) && p >= r[0] && p <= r[1] {
					chosenAnswer = i
				}
			}
		}
	}
	// Get the answer text, or fall back to the default answer.
	var answerID string
	if chosenAnswer != -1 {
		answerID = questions[gState.currentQuestion].answers[chosenAnswer].ID
	} else {
		answerID = questions[gState.currentQuestion].ID + "_DEFAULT"
	}
	text.Draw(screen, strings[answerID][gState.currentLover.answerIndex], defaultFont, 30, 510+12, defaultColors["purple"])
}

// Draws an animated trianglover on the screen.
// Originally this was going to support drawing multiple lovers on one page,
// but time pressure pushed the code to be less abstract.
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
	centerX := 0
	centerY := 0
	highestX := -1
	for i := range vertices {
		vertices[i][1] += yDiff
		vertices[i][0] += xDiff
		centerX += vertices[i][0]
		centerY += vertices[i][1]
		if highestX == -1 || vertices[i][0] > highestX {
			highestX = vertices[i][0]
		}
	}
	centerX = centerX / 3
	centerY = centerY / 3
	// Animate points.
	if gState.loverOffset <= 0 {
		gState.loverDirection = 1
	} else if gState.loverOffset >= 10 {
		gState.loverDirection = -1
	}
	gState.loverOffset += gState.loverDirection * .2
	for i := range []int{0, 1} {
		if vertices[i][0] > centerX {
			vertices[i][0] += int(gState.loverOffset)
		} else {
			vertices[i][0] -= int(gState.loverOffset)
		}
		if vertices[i][1] > centerY {
			vertices[i][1] += int(gState.loverOffset)
		} else {
			vertices[i][1] -= int(gState.loverOffset)
		}
	}
	drawPolygonLine(screen, .9, defaultColors["darkPink"], defaultColors["pink"], vertices)
	// Draw name.
	text.Draw(screen, lover.name, titleFont, 25, 70, defaultColors["purple"])
	// Draw face.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25, 175)
	screen.DrawImage(imageFiles["eyeball.png"], op)
	// Calculate the current eyeball rotation - it basically spins ina circle
	// but randomly changes direction 2% of the time (checked every frame).
	// Time isn't used for animation, but should be to account for FPS changes.
	if gState.eyeR <= 0 || (gState.eyeR < 1.5*math.Pi && rand.Intn(100) > 98) {
		gState.eyeDirection = 1
	} else if gState.eyeR >= 2*math.Pi || (gState.eyeR >= 1.5*math.Pi && rand.Intn(100) > 98) {
		gState.eyeDirection = -1
	}
	gState.eyeR += gState.eyeDirection * .02
	eyeWidth := 50
	eyeInnerWidth := 18
	middleOffset := float64((eyeWidth / 2) - (eyeInnerWidth / 2))
	xOffset := middleOffset + float64(eyeWidth/4)*math.Cos(gState.eyeR)
	yOffset := middleOffset + float64(eyeWidth/4)*math.Sin(gState.eyeR)
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

// Draws the match chart. This function is also used for matching and results.
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
	// Register drag points (targets).
	// This shouldn't be in the draw method, it's just easy to put it here.
	gState.dragTargets = points
	// Register drag points.
	gState.dragPoints[0].position = points[prefPoints[0]]
	gState.dragPoints[1].position = points[prefPoints[1]]
	gState.dragPoints[2].position = points[prefPoints[2]]
	// Draw hovered drag points.
	for _, point := range gState.dragPoints {
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
	hexLabels := []string{
		"Comfort",
		"Wealth",
		"Adventure",
		"Excitement",
		"Romance",
		"Family",
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
		clr := defaultColors["purple"]
		if i < 3 && gState.dragPoints[0].dragging {
			clr = defaultColors["darkPink"]
		} else if i > 1 && i < 5 && gState.dragPoints[1].dragging {
			clr = defaultColors["darkPink"]
		} else if (i > 3 || i == 0) && gState.dragPoints[2].dragging {
			clr = defaultColors["darkPink"]
		}
		text.Draw(screen, hexLabel, defaultFont, x, y, clr)
	}
}

// Handles the match chart's drag behavior.
func handleDrag() {
	mouseX, mouseY := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, point := range gState.dragPoints {
			point.dragging = false
		}
	}
	for _, point := range gState.dragPoints {
		if distance(point.position, vertex{mouseX, mouseY}) < 10 {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				point.dragging = true
			}
			point.hovering = true
		} else if !point.dragging {
			point.hovering = false
		}
	}
	for pointID, point := range gState.dragPoints {
		if !point.dragging {
			continue
		}
		min := 0
		// Points are locked into specific ranges because not all triangles on
		// the chart are valid. This reflects how trianglovers are randomly
		// generated.
		max := 34
		switch pointID {
		case 1:
			min = 34
			max = 68
		case 2:
			min = 68
			max = 102
		}
		// Find the closest point to the mouse and update the guessPoints.
		closestID := -1
		closestDistance := 100.0
		for i := min; i < max; i++ {
			d := distance(gState.dragTargets[i], vertex{mouseX, mouseY})
			if d < closestDistance {
				closestID = i
				closestDistance = d
			}
		}
		if closestID == -1 {
			continue
		}
		gState.currentLover.guessPoints[pointID] = closestID
	}
}

// Checks if a question has already been asked by the current lover.
func hasBeenAsked(q int) bool {
	for _, id := range gState.currentLover.questionsAsked {
		if q == id {
			return true
		}
	}
	return false
}

// Draws questions onto the screen.
func drawQuestions(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("Ask a question (%d/%d)", len(gState.currentLover.questionsAsked), gState.maxQuestions), largeFont, 400, 85, defaultColors["purple"])
	drawPolygonLine(screen, 2, defaultColors["darkPink"], defaultColors["white"], []vertex{{400, 100}, {780, 100}, {780, 325}, {400, 325}})
	x := 410
	y := 122
	for i, q := range questions {
		var clr color.Color
		if len(gState.currentLover.questionsAsked) >= gState.maxQuestions && !hasBeenAsked(i) {
			clr = defaultColors["pink"]
		} else if gState.currentQuestion == i {
			clr = defaultColors["darkPink"]
		} else if gState.hoverQuestion == i {
			clr = defaultColors["darkPink"]
		} else {
			clr = defaultColors["purple"]
		}
		text.Draw(screen, strings[q.ID][0], defaultFont, x, y+(i*25), clr)
	}
}

// Handles when questions have been clicked. In hard mode, only four questions
// can be asked of each lover, so this keeps track of that.
func handleQuestions() {
	mouseX, mouseY := ebiten.CursorPosition()
	y := 122
	gState.hoverQuestion = -1
	for i := range questions {
		if len(gState.currentLover.questionsAsked) >= gState.maxQuestions && !hasBeenAsked(i) {
			continue
		}
		qY := y + (i * 25)
		if mouseX >= 410 && mouseX <= 780 && mouseY <= qY && mouseY >= qY-12 {
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				if !hasBeenAsked(i) {
					gState.currentLover.questionsAsked = append(gState.currentLover.questionsAsked, i)
				}
				gState.currentQuestion = i
			} else {
				gState.hoverQuestion = i
			}
			break
		}
	}
}

// Handles when the next/previous/match buttons have been clicked.
func handleNextPrevious() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	if isButtonColliding("Previous Lover", 400, 335) {
		if gState.currentLoverIndex > 0 {
			gState.currentLoverIndex--
			gState.currentLover = gState.trianglovers[gState.currentLoverIndex]
			gState.currentQuestion = -1
			gState.hoverQuestion = -1
		}
	} else if isButtonColliding("Next Lover", 530, 335) {
		if gState.currentLoverIndex < len(gState.trianglovers)-1 {
			gState.currentLoverIndex++
			gState.currentLover = gState.trianglovers[gState.currentLoverIndex]
			gState.currentQuestion = -1
			gState.hoverQuestion = -1
		}
	} else if gState.currentLoverIndex == len(gState.trianglovers)-1 && isButtonColliding("Match!", 710, 335) {
		gState.gameMode = modeMatch
	}
}

// Informs users that they can drag points.
func drawHelper(screen *ebiten.Image) {
	if gState.currentLoverIndex == 0 && len(gState.currentLover.questionsAsked) > 0 && gState.currentLover.guessPoints == [3]int{0, 34, 68} {
		text.Draw(screen, "Drag the chart points!", defaultFont, width-360, height-170, defaultColors["purple"])
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(width-210, height-180)
		screen.DrawImage(imageFiles["arrow.png"], op)
	}
}

// The main wrapper for all guess page behaviors.
func handleGuess() {
	handleDrag()
	handleQuestions()
	handleNextPrevious()
}

// The main wrapper for all drawn elements on the guess page.
func drawGuess(screen *ebiten.Image) {
	drawMatchChart(screen, width-180, height-120, gState.currentLover.guessPoints, true, defaultColors["darkPink"])
	drawTrianglover(screen, gState.currentLover)
	drawQuestions(screen)
	drawAnswer(screen)
	drawNextPrevious(screen)
	drawHelper(screen)
}
