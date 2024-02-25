package main

import (
	"bufio"
	"fmt"
	imageColor "image/color"
	"os"
	"strconv"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	SCREEN_SIZE          = 700
	LINE_WEIGHT          = 2
	PIECE_OUTLINE_WEIGHT = 2
	TEXT_SIZE            = 16
	BOARD_OFFSET         = 0.05 * float32(SCREEN_SIZE)
	BOARD_SIZE           = 0.9 * float32(SCREEN_SIZE)
	PIECE_RADIUS         = BOARD_SIZE / 20
)

var BOARD_COLOR = [3]imageColor.RGBA{{R: 0, G: 161, B: 0, A: 255}, {R: 65, G: 105, B: 255, A: 255}, {R: 138, G: 81, B: 255, A: 255}}

type game struct {
	board                  *Board
	mouseClickedFlag       bool
	arrowLeftReleasedFlag  bool
	arrowRightReleasedFlag bool
	mouseX                 int
	mouseY                 int
	gameOverFlag           bool
	resultToShowFlag       bool
	analyzeFlag            bool
	analyzeLocations       []Location
	analyzeIndex           int
	analyzeBranchPoint     int
}

func transformBoardLocation(x int, y int) (float32, float32) {
	cx := BOARD_OFFSET + BOARD_SIZE*(float32(x)+0.5)/8
	cy := BOARD_OFFSET + BOARD_SIZE*(float32(y)+0.5)/8
	return cx, cy
}

func transformScreenPosition(cx float32, cy float32) (int, int) {
	x := 8 * (cx - BOARD_OFFSET) / BOARD_SIZE
	y := 8 * (cy - BOARD_OFFSET) / BOARD_SIZE
	return int(x), int(y)
}

func NewGame() *game {
	var g game
	g.reset()
	return &g
}

func drawPiece(screen *ebiten.Image, x int, y int, color int, alpha float32, outline float32) {
	c := uint8(255 * color)
	cx, cy := transformBoardLocation(x, y)
	a := imageColor.NRGBA{c, c, c, uint8(255 * alpha)}
	vector.DrawFilledCircle(screen, cx, cy, PIECE_RADIUS, a, true)
	if outline >= 0 {
		vector.StrokeCircle(screen, cx, cy, PIECE_RADIUS, PIECE_OUTLINE_WEIGHT*outline, imageColor.Gray{uint8(255 * (1 - color))}, true)
	}
}

func drawBoard(screen *ebiten.Image, b Board, mode int) {
	screen.Fill(imageColor.Gray{64})
	vector.DrawFilledRect(screen, BOARD_OFFSET, BOARD_OFFSET, BOARD_SIZE, BOARD_SIZE, BOARD_COLOR[mode], true)
	vector.StrokeRect(screen, BOARD_OFFSET, BOARD_OFFSET, BOARD_SIZE, BOARD_SIZE, LINE_WEIGHT, imageColor.Black, true)
	for i := 0; i <= 8; i++ {
		var p float32 = BOARD_OFFSET + BOARD_SIZE*float32(i)/8
		vector.StrokeLine(screen, BOARD_OFFSET, p, BOARD_OFFSET+BOARD_SIZE, p, LINE_WEIGHT, imageColor.Black, true)
		vector.StrokeLine(screen, p, BOARD_OFFSET, p, BOARD_OFFSET+BOARD_SIZE, LINE_WEIGHT, imageColor.Black, true)
	}
	var modeText string
	if mode == 0 {
		modeText = "Game Mode"
	} else if mode == 1 {
		modeText = "Analyze Mode"
	} else {
		modeText = "Analyze Mode (under consideration)"
	}
	text.Draw(screen, modeText, basicfont.Face7x13, SCREEN_SIZE*0.05, SCREEN_SIZE*0.03, imageColor.White)
	bc := b.counts[0]
	wc := b.counts[1]
	strBc := strconv.Itoa(bc)
	strWc := strconv.Itoa(wc)
	text.Draw(screen, "Black "+strBc+"-"+strWc+" White", basicfont.Face7x13, SCREEN_SIZE*0.4, SCREEN_SIZE*0.98, imageColor.White)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			c := b.board[x][y]
			if c != -1 {
				drawPiece(screen, x, y, c, 1, -1)
			}
		}
	}
	recs := b.records
	if len(recs) > 0 {
		lastRec := recs[len(recs)-1]
		if lastRec.flippedPieces != nil {
			color := lastRec.color
			placed := lastRec.placedPiece
			drawPiece(screen, placed.x, placed.y, color, 1, 1.5)
			for i := 0; i < len(lastRec.flippedPieces); i++ {
				flipped := lastRec.flippedPieces[i]
				drawPiece(screen, flipped.x, flipped.y, color, 1, 1)
			}
		}
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	switch {
	case g.analyzeFlag:
		x, y := transformScreenPosition(float32(g.mouseX), float32(g.mouseY))
		if g.analyzeBranchPoint == -10 {
			drawBoard(screen, *g.board, 1)
			locs := g.board.enumerateLegalLocations()
			for _, l := range locs {
				drawPiece(screen, l.x, l.y, g.board.currentColor, 0.3, -1)
			}
			var loc Location
			if g.arrowRightReleasedFlag {
				if g.analyzeIndex < len(g.analyzeLocations)-1 {
					g.analyzeIndex++
					loc = g.analyzeLocations[g.analyzeIndex]
					if g.board.IsLegal(loc.x, loc.y) {
						g.board.Put(loc.x, loc.y)
						if !g.board.IsLegalAll() {
							g.board.Pass()
						}
					} else {
						fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", "invalid value error")
						g.analyzeIndex--
					}
				}

			}
			if g.arrowLeftReleasedFlag {
				if g.analyzeIndex > -1 {
					g.analyzeIndex--
					g.board.Undo()
					if !g.board.IsLegalAll() {
						g.board.Undo()
					}
				}
			}

			if g.board.IsLegal(x, y) && g.mouseClickedFlag {
				g.board.Put(x, y)
				g.analyzeBranchPoint = g.analyzeIndex
				g.analyzeIndex++
				if !g.board.IsLegalAll() {
					g.board.Pass()
				}
			}
		} else {
			drawBoard(screen, *g.board, 2)
			locs := g.board.enumerateLegalLocations()
			for _, l := range locs {
				drawPiece(screen, l.x, l.y, g.board.currentColor, 0.3, -1)
			}
			if g.arrowLeftReleasedFlag {
				if g.analyzeIndex > -1 {
					g.analyzeIndex--
					if g.analyzeIndex == g.analyzeBranchPoint {
						g.analyzeBranchPoint = -10
					}
					g.board.Undo()
					if !g.board.IsLegalAll() {
						g.board.Undo()
					}
				}

			}
			if g.board.IsLegal(x, y) && g.mouseClickedFlag {
				g.board.Put(x, y)
				g.analyzeIndex++
				if !g.board.IsLegalAll() {
					g.board.Pass()
				}
			}
		}
	case g.resultToShowFlag:
		drawBoard(screen, *g.board, 0)
		bc := g.board.counts[0]
		wc := g.board.counts[1]
		fmt.Printf("\x1b[33m%s %d-%d %s\x1b[0m\n", "Black", bc, wc, "White")
		if bc > wc {
			fmt.Printf("\x1b[33m%s\x1b[0m\n", "**BLACK WINS**")
		} else if bc < wc {
			fmt.Printf("\x1b[33m%s\x1b[0m\n", "**WHITE WINS**")
		} else {
			fmt.Printf("\x1b[33m%s\x1b[0m\n", "**DRAW**")
		}
		fmt.Println("Press 'S' to export game records")
		fmt.Println("Press 'A' to analyze the game")
		g.resultToShowFlag = false
	case !g.gameOverFlag:
		nextTurnFlag := false
		x, y := transformScreenPosition(float32(g.mouseX), float32(g.mouseY))
		if g.board.IsLegal(x, y) && g.mouseClickedFlag {
			g.board.Put(x, y)
			nextTurnFlag = true
		} else {
			drawBoard(screen, *g.board, 0)
			locs := g.board.enumerateLegalLocations()
			for _, l := range locs {
				drawPiece(screen, l.x, l.y, g.board.currentColor, 0.3, -1)
			}
		}
		if nextTurnFlag {
			if !g.board.IsLegalAll() {
				g.board.Pass()
				if !g.board.IsLegalAll() {
					g.gameOverFlag = true
					g.resultToShowFlag = true
				}
			}
			drawBoard(screen, *g.board, 0)
		}
	default:
		drawBoard(screen, *g.board, 0)
	}
	g.mouseClickedFlag = false
	g.arrowRightReleasedFlag = false
	g.arrowLeftReleasedFlag = false
}

func (g *game) reset() {
	g.board = InitBoard()
	g.mouseClickedFlag = false
	g.gameOverFlag = false
	g.resultToShowFlag = false
	g.analyzeFlag = false
	g.analyzeIndex = -1
	g.analyzeLocations = []Location{}
	fmt.Printf("\x1b[33m%s\x1b[0m\n", "GAME START!")
	fmt.Println("Press 'I' to import game records and analyze")
	fmt.Println("Press 'R' to restart")
}

func (g *game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.mouseClickedFlag = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		g.arrowLeftReleasedFlag = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		g.arrowRightReleasedFlag = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		g.reset()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyI) {
		g.analyzeLocations = []Location{}
		var file string
		fmt.Printf("Input the path of the import file>>")
		fmt.Scan(&file)
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", err)
			f.Close()
			return nil
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			text := s.Text()
			var i int
			var x int
			var y int
			var str string
			for text != "" {
				str = text[0:1]
				text = text[1:]
				if i%2 == 0 {
					switch str {
					case "a":
						x = 0
					case "b":
						x = 1
					case "c":
						x = 2
					case "d":
						x = 3
					case "e":
						x = 4
					case "f":
						x = 5
					case "g":
						x = 6
					case "h":
						x = 7
					default:
						fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", "File format is invalid")
						return nil
					}
				} else {
					y, err = strconv.Atoi(str)
					if err != nil {
						fmt.Fprintf(os.Stderr, "\x1b[31m%s: %s\x1b[0m\n", "File format is invalid", err)
						return nil
					}
					y--
					g.analyzeLocations = append(g.analyzeLocations, Location{x, y})
				}
				i++
			}
		}
		if err := s.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", err)
			return nil
		}
		g.board = InitBoard()
		g.analyzeFlag = true
		g.analyzeIndex = -1
		g.gameOverFlag = false
		g.analyzeBranchPoint = -10
		fmt.Println("Imported:", file)
		fmt.Println("ANALYZE MODE!\nUse '←','→'")
	}
	if g.gameOverFlag && inpututil.IsKeyJustReleased(ebiten.KeyS) {
		var file string
		fmt.Printf("Input a name for the export file>>")
		fmt.Scan(&file)
		f, err := os.Create(file)
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", err)
			}
		}()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", err)
			return nil
		}
		var str string
		for _, v := range g.board.records {
			if v.flippedPieces != nil {
				l := v.placedPiece
				switch l.x {
				case 0:
					str += "a"
				case 1:
					str += "b"
				case 2:
					str += "c"
				case 3:
					str += "d"
				case 4:
					str += "e"
				case 5:
					str += "f"
				case 6:
					str += "g"
				case 7:
					str += "h"
				default:
					fmt.Fprintf(os.Stderr, "\x1b[31m%s: %s\x1b[0m\n", "invalid value error", err)
					return nil
				}
				str += strconv.Itoa(l.y + 1)
			}
		}

		if _, err := fmt.Fprintln(f, str); err != nil {
			fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", err)
			return nil
		}
		fmt.Println("Exported:", file)
	}
	if g.gameOverFlag && inpututil.IsKeyJustReleased(ebiten.KeyA) {
		g.analyzeLocations = []Location{}
		for _, v := range g.board.records {
			if v.flippedPieces != nil {
				g.analyzeLocations = append(g.analyzeLocations, v.placedPiece)
			}
		}
		g.board = InitBoard()
		g.analyzeFlag = true
		g.analyzeIndex = -1
		g.gameOverFlag = false
		g.analyzeBranchPoint = -10
		fmt.Println("ANALYZE MODE!\nUse '←','→'")
	}

	return nil
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_SIZE, SCREEN_SIZE
}
