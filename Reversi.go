package main

import (
	imageColor "image/color"
	"strconv"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	//"golang.org/x/image/font/basicfont"
)

const (
	SCREEN_SIZE          = 500
	BOARD_COLOR          = 0xff00a000
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

func getColorName(color int) string {
	if color == 0 {
		return "黒"
	} else if color == 1 {
		return "白"
	} else {
		return strconv.Itoa(color)
	}
}

func transformBoardLocation(x int, y int) (int, int) {
	cx := BOARD_OFFSET + BOARD_SIZE*(float32(x)+0.5)/8
	cy := BOARD_OFFSET + BOARD_SIZE*(float32(y)+0.5)/8
	return int(cx), int(cy)
}

func transformScreenPosition(cx int, cy int) (int, int) {
	x := 8 * (float32(cx) - BOARD_OFFSET) / BOARD_SIZE
	y := 8 * (float32(cy) - BOARD_OFFSET) / BOARD_SIZE
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
	ebiten.DrawFilledCircle(screen, cx, cy, PIECE_RADIUS, a, true)
	if outline >= 0 {
		ebiten.StrokeCircle(screen, cx, cy, PIECE_RADIUS, PIECE_OUTLINE_WEIGHT*outline, 255*(1-color), true)
	}
}

func drawBoard(screen *ebiten.Image, b Board) {
	screen.Fill(imageColor.Gray{64})
	ebiten.DrawFilledRect(screen, BOARD_OFFSET, BOARD_OFFSET, BOARD_SIZE, BOARD_SIZE, BOARD_COLOR, true)
	ebiten.StrokeRect(screen, BOARD_OFFSET, BOARD_OFFSET, BOARD_SIZE, BOARD_SIZE, LINE_WEIGHT, 0, true)
	for i := 0; i <= 8; i++ {
		var p float32 = BOARD_OFFSET + BOARD_SIZE*float32(i)/8
		ebiten.StrokeLine(screen, BOARD_OFFSET, p, BOARD_OFFSET+BOARD_SIZE, p, LINE_WEIGHT, 0, true)
		ebiten.StrokeLine(screen, p, BOARD_OFFSET, p, BOARD_OFFSET+BOARD_SIZE, LINE_WEIGHT, 0, true)
	}
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
		bc := g.board.counts[0]
		wc := g.board.counts[1]
		str := "黒" + strconv.Itoa(bc) + "-" + strconv.Itoa(wc) + "白\n"
		if bc > wc {
			str += "黒の勝ち"
		} else if bc < wc {
			str += "白の勝ち"
		} else {
			str += "引き分け"
		}
		ebiten.text.Draw(screen, str, ebiten.basicfont, SCREEN_SIZE/2, SCREEN_SIZE/2, 0)
		g.resultToShowFlag = false
	case !g.gameOverFlag:
		nextTurnFlag := false
		x, y := transformScreenPosition(g.mouseX, g.mouseY)
		if x < 0 || x >= 8 || y < 0 || y >= 8 || !g.board.IsLegal(x, y) {
			drawBoard(screen, *g.board)
		} else if g.mouseClickedFlag {
			g.board.Put(x, y)
			nextTurnFlag = true
		} else {
			drawBoard(screen, *g.board)
			drawPiece(screen, x, y, g.board.currentColor, 0.5, -1)
		}
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
	}
	g.mouseClickedFlag = false
}

func (g *game) reset() {
	g.board = InitBoard()
	g.mouseClickedFlag = false
	g.gameOverFlag = false
	// g.resultToShowFlag=false
	//mouse
}

func (g *game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.mouseClickedFlag = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.board = InitBoard()
		g.gameOverFlag = false
	}
	return nil
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_SIZE, SCREEN_SIZE
}
