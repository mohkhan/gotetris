/*
Package main contains a console-based implementation of Tetris.
See the README for more details.

I don't have any tests or that much in the way of documentation. It's
just a simple video game ;)
*/

package main

import (
	"fmt"
	//// 	"math/rand"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

// Colors
const backgroundColor = termbox.ColorBlue
const boardColor = termbox.ColorBlack
const instructionsColor = termbox.ColorWhite

var pieceColors = []termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorRed,
	termbox.ColorGreen,
	termbox.ColorYellow,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorCyan,
	termbox.ColorWhite,
}

// Layout
const defaultMarginWidth = 2
const defaultMarginHeight = 1
const titleStartX = defaultMarginWidth
const titleStartY = defaultMarginHeight
const titleHeight = 1
const titleEndY = titleStartY + titleHeight
const boardStartX = defaultMarginWidth
const boardStartY = titleEndY + defaultMarginHeight
const boardWidth = 10
const boardHeight = 16
const boardEndX = boardStartX + boardWidth
const boardEndY = boardStartY + boardHeight
const instructionsStartX = boardEndX + defaultMarginWidth
const instructionsStartY = boardStartY

// Text in the UI
const title = "Tetris Written in Go"

var instructions = []string{
	"Goal: Fill in 5 lines!",
	"",
	"\u2190      Left",
	"\u2192      Right",
	"\u2191      Rotate",
	"\u2193      Drop faster",
	"s      Start",
	"p      Pause",
	"esc    Exit",
	"",
	"Level: %v",
	"Lines: %v",
}

// Game play
const numSquares = 4
const numTypes = 7
const defaultLevel = 1
const maxLevel = 10
const rowsPerLevel = 5
const slowestSpeed = 700
const fastestSpeed = 60

// Keystroke processing
const initialDelay = 200
const repeatDelay = 20

type Game struct {
	curLevel    int
	curX        int
	curY        int
	curPiece    int
	skyline     int
	gamePaused  bool
	gameStarted bool
	timer       *time.Timer
	numLines    int
	speed       int
	board       [][]int // [y][x]
	xToErase    []int
	yToErase    []int
	dx          []int
	dy          []int
	dxPrime     []int
	dyPrime     []int
	dxBank      [][]int
	dyBank      [][]int

	// Keystroke processing
	isActiveLeft  bool
	isActiveRight bool
	isActiveUp    bool
	isActiveDown  bool
	timerLeft     *time.Timer
	timerRight    *time.Timer
	timerDown     *time.Timer
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func draw() {
	termbox.Clear(backgroundColor, backgroundColor)
	tbprint(titleStartX, titleStartY, instructionsColor, backgroundColor, title)
	for y := boardStartY; y < boardEndY; y++ {
		for x := boardStartX; x < boardEndX; x++ {
			termbox.SetCell(x, y, ' ', boardColor, boardColor)
		}
	}
	for i, instruction := range instructions {
		if strings.HasPrefix(instruction, "Level:") {
			instruction = fmt.Sprintf(instruction, 0)
		} else if strings.HasPrefix(instruction, "Lines:") {
			instruction = fmt.Sprintf(instruction, 0)
		}
		tbprint(instructionsStartX, instructionsStartY+i, instructionsColor, backgroundColor, instruction)
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	NewGame()

	draw()

loop:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}

			//// 	document.on.keyDown.add(g.onKeyDown)
			//// 	document.on.keyUp.add(onKeyUp)
			////
			//// 	SelectElement levelSelect = query("#level-select")
			//// 	levelSelect.on.change.add((Event e) {
			//// 		g.onLevelSelectChange()
			//// 		levelSelect.blur()
			//// 	})
			////
			//// 	InputElement startButton = query("#start-button")
			//// 	InputElement pauseButton = query("#pause-button")
			//// 	startButton.on.click.add((event) => g.start())
			//// 	pauseButton.on.click.add((event) => g.pause())

		default:
			draw()
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func NewGame() (g *Game) {
	g = new(Game)

	g.curLevel = defaultLevel
	g.curX = 1
	g.curY = 1
	g.skyline = boardHeight - 1
	g.gamePaused = false
	g.gameStarted = false
	g.numLines = 0
	g.speed = slowestSpeed - fastestSpeed*defaultLevel

	// Keystroke processing
	g.isActiveLeft = false
	g.isActiveRight = false
	g.isActiveUp = false
	g.isActiveDown = false

	g.board = make([][]int, boardHeight)
	for i := 0; i < boardHeight; i++ {
		g.board[i] = make([]int, boardWidth)
		for j := 0; j < boardWidth; j++ {
			g.board[i][j] = 0
		}
	}

	g.xToErase = []int{0, 0, 0, 0}
	g.yToErase = []int{0, 0, 0, 0}
	g.dx = []int{0, 0, 0, 0}
	g.dy = []int{0, 0, 0, 0}
	g.dxPrime = []int{0, 0, 0, 0}
	g.dyPrime = []int{0, 0, 0, 0}

	g.dxBank = [][]int{
		{},
		{0, 1, -1, 0},
		{0, 1, -1, -1},
		{0, 1, -1, 1},
		{0, -1, 1, 0},
		{0, 1, -1, 0},
		{0, 1, -1, -2},
		{0, 1, 1, 0},
	}

	g.dyBank = [][]int{
		{},
		{0, 0, 0, 1},
		{0, 0, 0, 1},
		{0, 0, 0, 1},
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 1, 1},
	}

	return
}

//// func (g *Game) run() {
//// 	g.drawBoard()
//// 	g.resetGame()
//// }
////
//// func (g *Game) start() {
//// 	if g.gameStarted {
//// 		if g.gamePaused {
//// 			g.resume()
//// 		}
//// 		return
//// 	}
//// 	g.getPiece()
//// 	g.drawPiece()
//// 	g.gameStarted = true
//// 	g.gamePaused = false
//// 	InputElement g.numLinesField = query("#num-lines")
//// 	g.numLinesField.value = g.numLines.toString()
//// 	g.timer = new(Timer(g.speed, (timer) => g.play()))
//// }
////
//// func (g *Game) drawBoard() {
//// 	DivElement g.boardDiv = query("#g.board-div")
//// 	pre := new(Element.tag("pre"))
//// 	g.boardDiv.nodes.add(pre)
//// 	pre.classes.add("g.board")
//// 	for i := 0; i < g.boardHeight; i++ {
//// 		div := new(Element.tag("div"))
//// 		pre.nodes.add(div)
//// 		for j := 0; j < g.boardWidth; j++ {
//// 			img := new(Element.tag("img"))
//// 			div.nodes.add(img)
//// 			img.id = "s-$i-$j"
//// 			img.src = "images/s${g.board[i][j].abs()}.png"
//// 		}
//// 		rightMargin := new(Element.tag("img"))
//// 		div.nodes.add(rightMargin)
//// 		rightMargin.src = "images/g.png"
//// 		rightMargin.width = 1
//// 	}
//// 	trailingDiv = new(Element.tag("div"))
//// 	pre.nodes.add(trailingDiv)
//// 	trailingImg = new(Element.tag("img"))
//// 	trailingDiv.nodes.add(trailingImg)
//// 	trailingImg.src = "images/g.png"
//// 	trailingImg.id = "g.board-trailing-img"
//// 	trailingImg.width = g.boardWidth * 16 + 1
//// 	trailingImg.height = 1
//// }
////
//// func (g *Game) resetGame() {
//// 	for i := 0; i < g.boardHeight; i++ {
//// 		for j := 0; j < g.boardWidth; j++ {
//// 			g.board[i][j] = 0
//// 			ImageElement img = query("#s-$i-$j")
//// 			img.src = "images/s0.png"
//// 		}
//// 	}
//// 	g.gameStarted = false
//// 	g.gamePaused = false
//// 	g.numLines = 0
//// 	g.curLevel = 1
//// 	g.skyline = g.boardHeight - 1
//// 	InputElement g.numLinesField = query("#num-lines")
//// 	g.numLinesField.value = g.numLines.toString()
//// 	SelectElement levelSelect = query("#level-select")
//// 	levelSelect.selectedIndex = 0
////
//// 	// I shouldn"t have to call this manually, but I do.
//// 	// See: http://code.google.com/p/dart/issues/detail?id=2325&thanks=2325&ts=1332879888
//// 	g.onLevelSelectChange()
//// }
////
//// func (g *Game) play() {
//// 	if g.moveDown() {
//// 		g.timer = new(Timer(g.speed, (timer) => g.play()))
//// 	} else {
//// 		g.fillMatrix()
//// 		g.removeLines()
//// 		if g.skyline > 0 && g.getPiece() {
//// 			g.timer = new(Timer(g.speed, (timer) => g.play()))
//// 		} else {
//// 			g.isActiveLeft = false
//// 			g.isActiveUp = false
//// 			g.isActiveRight = false
//// 			g.isActiveDown = false
//// 			window.alert("Game over!")
//// 			g.resetGame()
//// 		}
//// 	}
//// }
////
//// func (g *Game) pause() {
//// 	if g.gameStarted {
//// 		if g.gamePaused {
//// 			g.resume()
//// 			return
//// 		}
//// 		g.timer.cancel()
//// 		g.gamePaused = true
//// 	}
//// }
////
//// func (g *Game) onKeyDown(event KeyboardEvent) {
//// 	// I"m positive there are more modern ways to do keyboard event handling.
//// 	String keyNN, keyIE
//// 	keyNN = " ${event.keyCode} "
//// 	keyIE = " ${event.keyCode} "
////
//// 	// Only preventDefault if we can actually handle the keyDown event.  If we
//// 	// capture all keyDown events, we break things like using ^r to reload the page.
////
//// 	if !g.gameStarted || g.gamePaused {
//// 		return
//// 	}
////
//// 	if leftNN.indexOf(keyNN) != -1 || leftIE.indexOf(keyIE) != -1 {
//// 		if !g.isActiveLeft {
//// 			g.isActiveLeft = true
//// 			g.isActiveRight = false
//// 			g.moveLeft()
//// 			g.timerLeft = new(Timer(initialDelay, (timer) => g.slideLeft()))
//// 		}
//// 	} else if rightNN.indexOf(keyNN) != -1 || rightIE.indexOf(keyIE) != -1 {
//// 		if !g.isActiveRight {
//// 			g.isActiveRight = true
//// 			g.isActiveLeft = false
//// 			g.moveRight()
//// 			g.timerRight = new(Timer(initialDelay, (timer) => g.slideRight()))
//// 		}
//// 	} else if upNN.indexOf(keyNN) != -1 || upIE.indexOf(keyIE) != -1 {
//// 		if !g.isActiveUp {
//// 			g.isActiveUp = true
//// 			g.isActiveDown = false
//// 			g.rotate()
//// 		}
//// 	} else if spaceNN.indexOf(keyNN) != -1 || spaceIE.indexOf(keyIE) != -1 {
//// 		if !g.isActiveSpace {
//// 			g.isActiveSpace = true
//// 			g.isActiveDown = false
//// 			g.fall()
//// 		}
//// 	} else if downNN.indexOf(keyNN) != -1 || downIE.indexOf(keyIE) != -1 {
//// 		if !g.isActiveDown {
//// 			g.isActiveDown = true
//// 			g.isActiveUp = false
//// 			g.moveDown()
//// 			g.timerDown = new(Timer(initialDelay, (timer) => g.slideDown()))
//// 		}
//// 	} else {
//// 		return
//// 	}
//// 	event.preventDefault()
//// }
////
//// func (g *Game) onKeyUp(event KeyboardEvent) {
//// 	// See comments in g.onKeyDown.
//// 	keyNN := " ${event.keyCode} "
//// 	keyIE := " ${event.keyCode} "
//// 	if leftNN.indexOf(keyNN) != -1 || leftIE.indexOf(keyIE) != -1 {
//// 		g.isActiveLeft = false
//// 		g.timerLeft.cancel()
//// 	} else if rightNN.indexOf(keyNN) != -1 || rightIE.indexOf(keyIE) != -1 {
//// 		g.isActiveRight = false
//// 		g.timerRight.cancel()
//// 	} else if upNN.indexOf(keyNN) != -1 || upIE.indexOf(keyIE) != -1 {
//// 		g.isActiveUp = false
//// 	} else if downNN.indexOf(keyNN) != -1 || downIE.indexOf(keyIE) != -1 {
//// 		g.isActiveDown = false
//// 		g.timerDown.cancel()
//// 	} else if spaceNN.indexOf(keyNN) != -1 || spaceIE.indexOf(keyIE) != -1 {
//// 		g.isActiveSpace = false
//// 	} else {
//// 		return
//// 	}
//// 	event.preventDefault()
//// }
////
//// func (g *Game) fillMatrix() {
//// 	for k := 0; k < numSquares; k++ {
//// 		x := g.curX + g.dx[k]
//// 		y := g.curY + g.dy[k]
//// 		if 0 <= y && y < g.boardHeight && 0 <= x && x < g.boardWidth {
//// 			g.board[y][x] = g.curPiece
//// 			if y < g.skyline {
//// 				g.skyline = y
//// 			}
//// 		}
//// 	}
//// }
////
//// func (g *Game) removeLines() {
//// 	for i := 0; i < g.boardHeight; i++ {
//// 		gapFound := false
//// 		for j := 0; j < g.boardWidth; j++ {
//// 			if g.board[i][j] == 0 {
//// 				gapFound = true
//// 				break
//// 			}
//// 		}
//// 		if !gapFound {
//// 			for k := i; k >= g.skyline; k-- {
//// 				for j = 0; j < g.boardWidth; j++ {
//// 					g.board[k][j] = g.board[k - 1][j]
//// 					ImageElement img = query("#s-$k-$j")
//// 					img.src = pieceColors[g.board[k][j]].src
//// 				}
//// 			}
//// 			for j := 0; j < g.boardWidth; j++ {
//// 				g.board[0][j] = 0
//// 				ImageElement img = query("#s-0-$j")
//// 				img.src = pieceColors[0].src
//// 			}
//// 			g.numLines++
//// 			g.skyline++
//// 			InputElement g.numLinesField = query("#num-lines")
//// 			g.numLinesField.value = g.numLines.toString()
//// 			if g.numLines % rowsPerLevel == 0 && g.curLevel < maxLevel {
//// 				g.curLevel++
//// 			}
//// 			g.speed = slowestSpeed - fastestSpeed * g.curLevel
//// 			SelectElement levelSelect = query("#level-select")
//// 			levelSelect.selectedIndex = g.curLevel - 1
//// 		}
//// 	}
//// }
////
//// func (g *Game) pieceFits(x, y) bool {
//// 	for k := 0; k < numSquares; k++ {
//// 		theX := x + g.dxPrime[k]
//// 		theY := y + g.dyPrime[k]
//// 		if theX < 0 || theX >= g.boardWidth || theY >= g.boardHeight {
//// 			return false
//// 		}
//// 		if theY > -1 && g.board[theY][theX] > 0 {
//// 			return false
//// 		}
//// 	}
//// 	return true
//// }
////
//// func (g *Game) erasePiece() {
//// 	for k := 0; k < numSquares; k++ {
//// 		x := g.curX + g.dx[k]
//// 		y := g.curY + g.dy[k]
//// 		if 0 <= y && y < g.boardHeight && 0 <= x && x < g.boardWidth {
//// 			g.xToErase[k] = x
//// 			g.yToErase[k] = y
//// 			g.board[y][x] = 0
//// 		}
//// 	}
//// }
////
//// func (g *Game) drawPiece() {
//// 	for k := 0; k < numSquares; k++ {
//// 		x = g.curX + g.dx[k]
//// 		y = g.curY + g.dy[k]
//// 		if 0 <= y && y < g.boardHeight && 0 <= x && x < g.boardWidth && g.board[y][x] != -g.curPiece {
//// 			ImageElement img = query("#s-$y-$x")
//// 			img.src = pieceColors[g.curPiece].src
//// 			g.board[y][x] = -g.curPiece
//// 		}
//// 		x := g.xToErase[k]
//// 		y := g.yToErase[k]
//// 		if g.board[y][x] == 0 {
//// 			ImageElement img = query("#s-$y-$x")
//// 			img.src = pieceColors[0].src
//// 		}
//// 	}
//// }
////
//// func (g *Game) moveDown() {
//// 	for k := 0; k < numSquares; k++ {
//// 		g.dxPrime[k] = g.dx[k]
//// 		g.dyPrime[k] = g.dy[k]
//// 	}
//// 	if g.pieceFits(g.curX, g.curY + 1) {
//// 		g.erasePiece()
//// 		g.curY++
//// 		g.drawPiece()
//// 		return true
//// 	}
//// 	return false
//// }
////
//// func (g *Game) moveLeft() {
//// 	for k := 0; k < numSquares; k++ {
//// 		g.dxPrime[k] = g.dx[k]
//// 		g.dyPrime[k] = g.dy[k]
//// 	}
//// 	if g.pieceFits(g.curX - 1, g.curY) {
//// 		g.erasePiece()
//// 		g.curX--
//// 		g.drawPiece()
//// 	}
//// }
////
//// func (g *Game) moveRight() {
//// 	for k := 0; k < numSquares; k++ {
//// 		g.dxPrime[k] = g.dx[k]
//// 		g.dyPrime[k] = g.dy[k]
//// 	}
//// 	if g.pieceFits(g.curX + 1, g.curY) {
//// 		g.erasePiece()
//// 		g.curX++
//// 		g.drawPiece()
//// 	}
//// }
////
//// func (g *Game) getPiece() bool {
//// 	g.curPiece = 1 + rand.Int()%numTypes
//// 	g.curX = 5
//// 	g.curY = 0
//// 	for k := 0; k < numSquares; k++ {
//// 		g.dx[k] = g.dxBank[g.curPiece][k]
//// 		g.dy[k] = g.dyBank[g.curPiece][k]
//// 	}
//// 	for k = 0; k < numSquares; k++ {
//// 		g.dxPrime[k] = g.dx[k]
//// 		g.dyPrime[k] = g.dy[k]
//// 	}
//// 	if g.pieceFits(g.curX, g.curY) {
//// 		g.drawPiece()
//// 		return true
//// 	}
//// 	return false
//// }
////
//// func (g *Game) resume() {
//// 	if g.gameStarted && g.gamePaused {
//// 		g.play()
//// 		g.gamePaused = false
//// 	}
//// }
////
//// func (g *Game) onLevelSelectChange() {
//// 	SelectElement levelSelect = query("#level-select")
//// 	OptionElement selectedOption = levelSelect.options[levelSelect.selectedIndex]
//// 	g.curLevel = int.parse(selectedOption.value)
//// 	g.speed = slowestSpeed - fastestSpeed * g.curLevel
//// }
////
//// func (g *Game) rotate() {
//// 	for k := 0; k < numSquares; k++ {
//// 		g.dxPrime[k] = g.dy[k]
//// 		g.dyPrime[k] = -g.dx[k]
//// 	}
//// 	if g.pieceFits(g.curX, g.curY) {
//// 		g.erasePiece()
//// 		for k = 0; k < numSquares; k++ {
//// 			g.dx[k] = g.dxPrime[k]
//// 			g.dy[k] = g.dyPrime[k]
//// 		}
//// 		g.drawPiece()
//// 	}
//// }
////
//// func (g *Game) fall() {
//// 	for k := 0; k < numSquares; k++ {
//// 		g.dxPrime[k] = g.dx[k]
//// 		g.dyPrime[k] = g.dy[k]
//// 	}
//// 	if !g.pieceFits(g.curX, g.curY + 1) {
//// 		return
//// 	}
//// 	g.timer.cancel()
//// 	g.erasePiece()
//// 	while g.pieceFits(g.curX, g.curY + 1) {
//// 		g.curY++
//// 	}
//// 	g.drawPiece()
//// 	g.timer = new(Timer(g.speed, (timer) => g.play()))
//// }
////
//// func (g *Game) slideLeft() {
//// 	if g.isActiveLeft {
//// 		g.moveLeft()
//// 		g.timerLeft = new(Timer(repeatDelay, (timer) => g.slideLeft()))
//// 	}
//// }
////
//// func (g *Game) slideRight() {
//// 	if g.isActiveRight {
//// 		g.moveRight()
//// 		g.timerRight = new(Timer(repeatDelay, (timer) => g.slideRight()))
//// 	}
//// }
////
//// func (g *Game) slideDown() {
//// 	if g.isActiveDown {
//// 		g.moveDown()
//// 		g.timerDown = new(Timer(repeatDelay, (timer) => g.slideDown()))
//// 	}
//// }
////
//// func (g *Game) write(String message) {
//// 	p = new(Element.tag("p"))
//// 	p.text = message
//// 	document.body.nodes.add(p)
//// }
////
//// void main() {
//// 	new(Game).run()
//// }////
