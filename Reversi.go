package main

import (
	"fmt"
	imageColor "image/color"
	"strconv"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	SCREEN_SIZE          = 700
	BOARD_COLOR_R        = 0
	BOARD_COLOR_G        = 161
	BOARD_COLOR_B        = 0
	BOARD_COLOR_A        = 255
	LINE_WEIGHT          = 2
	PIECE_OUTLINE_WEIGHT = 2
	TEXT_SIZE            = 16
	BOARD_OFFSET         = 0.05 * float32(SCREEN_SIZE)
	BOARD_SIZE           = 0.9 * float32(SCREEN_SIZE)
	PIECE_RADIUS         = BOARD_SIZE / 20
)

type game struct {
	board            *Board
	mouseClickedFlag bool
	mouseX           int
	mouseY           int
	gameOverFlag     bool
	resultToShowFlag bool
}

// func getColorName(color int) string {
// 	if color == 0 {
// 		return "黒"
// 	} else if color == 1 {
// 		return "白"
// 	} else {
// 		return strconv.Itoa(color)
// 	}
// }

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

func newGame() *game {
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

func drawBoard(screen *ebiten.Image, b Board) {
	screen.Fill(imageColor.Gray{64})
	vector.DrawFilledRect(screen, BOARD_OFFSET, BOARD_OFFSET, BOARD_SIZE, BOARD_SIZE, imageColor.RGBA{BOARD_COLOR_R, BOARD_COLOR_G, BOARD_COLOR_B, BOARD_COLOR_A}, true)
	vector.StrokeRect(screen, BOARD_OFFSET, BOARD_OFFSET, BOARD_SIZE, BOARD_SIZE, LINE_WEIGHT, imageColor.Black, true)
	for i := 0; i <= 8; i++ {
		var p float32 = BOARD_OFFSET + BOARD_SIZE*float32(i)/8
		vector.StrokeLine(screen, BOARD_OFFSET, p, BOARD_OFFSET+BOARD_SIZE, p, LINE_WEIGHT, imageColor.Black, true)
		vector.StrokeLine(screen, p, BOARD_OFFSET, p, BOARD_OFFSET+BOARD_SIZE, LINE_WEIGHT, imageColor.Black, true)
	}
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
	case g.resultToShowFlag:
		drawBoard(screen, *g.board)
		bc := g.board.counts[0]
		wc := g.board.counts[1]
		fmt.Printf("Black %d-%d White\n", bc, wc)
		if bc > wc {
			fmt.Println("**BLACK WINS**")
		} else if bc < wc {
			fmt.Println("**WHITE WINS**")
		} else {
			fmt.Println("**DRAW**")
		}
		fmt.Println("Press 'R' to restart")
		g.resultToShowFlag = false
	case !g.gameOverFlag:
		nextTurnFlag := false
		x, y := transformScreenPosition(float32(g.mouseX), float32(g.mouseY))
		if g.board.IsLegal(x, y) && g.mouseClickedFlag {
			g.board.Put(x, y)
			nextTurnFlag = true
		} else {
			drawBoard(screen, *g.board)
			locs := g.board.enumerateLegalLocations()
			for _, l := range locs {
				drawPiece(screen, l.x, l.y, g.board.currentColor, 0.3, -1)
			}
		}
		// if x < 0 || x >= 8 || y < 0 || y >= 8 || !g.board.IsLegal(x, y) {
		// 	drawBoard(screen, *g.board)
		// } else if g.mouseClickedFlag {
		// 	g.board.Put(x, y)
		// 	nextTurnFlag = true
		// } else {
		// 	drawBoard(screen, *g.board)
		// 	drawPiece(screen, x, y, g.board.currentColor, 0.5, -1)
		// }
		if nextTurnFlag {
			if !g.board.IsLegalAll() {
				g.board.Pass()
				if !g.board.IsLegalAll() {
					g.board.Undo()
					g.gameOverFlag = true
					g.resultToShowFlag = true
				}
			}
			drawBoard(screen, *g.board)
		}
	default:
		drawBoard(screen, *g.board)
	}
	g.mouseClickedFlag = false
}

func (g *game) reset() {
	g.board = InitBoard()
	g.mouseClickedFlag = false
	g.gameOverFlag = false
	g.resultToShowFlag = false
}

func (g *game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.mouseClickedFlag = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.reset()
	}
	return nil
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_SIZE, SCREEN_SIZE
}
