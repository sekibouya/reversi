package main

type Location struct {
	x int
	y int
}

type Record struct {
	color         int
	placedPiece   Location
	flippedPieces []Location
}

func (r *Record) flip(piece Location) {
	r.flippedPieces = append(r.flippedPieces, piece)
}

type Board struct {
	board        [8][8]int
	counts       [2]int
	records      []*Record
	currentColor int
}

var DIRECTIONS = [8]Location{{x: -1, y: 0}, {x: -1, y: 1}, {x: 0, y: 1}, {x: 1, y: 1}, {x: 1, y: 0}, {x: 1, y: -1}, {x: 0, y: -1}, {x: -1, y: -1}}

func Flip(color int) int {
	switch color {
	case 0:
		return 1
	case 1:
		return 0
	default:
		return -1
	}
}

func InitBoard() *Board {
	var board [8][8]int
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			board[x][y] = -1
		}
	}
	board[4][3] = 0
	board[3][4] = 0
	board[3][3] = 1
	board[4][4] = 1
	var counts [2]int
	counts[0] = 2
	counts[1] = 2
	var records []*Record
	b := &Board{
		board:        board,
		counts:       counts,
		records:      records,
		currentColor: 0,
	}
	return b
}

func (b *Board) isLegal(x int, y int, direction int) bool {
	d := DIRECTIONS[direction]
	for i := 1; i < 8; i++ {
		x := x + d.x*i
		y := y + d.y*i
		if x < 0 || x >= 8 || y < 0 || y >= 8 {
			return false
		}
		c := b.board[x][y]
		if c == -1 {
			return false
		} else if c == b.currentColor {
			return i > 1
		}
	}
	return false
}

func (b *Board) IsLegal(x int, y int) bool {
	if x < 0 || x >= 8 || y < 0 || y >= 8 || b.board[x][y] != -1 {
		return false
	}
	for i := 0; i < 8; i++ {
		if b.isLegal(x, y, i) {
			return true
		}
	}
	return false
}

func (b *Board) IsLegalAll() bool {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if b.IsLegal(x, y) {
				return true
			}
		}
	}
	return false
}

func (b *Board) enumerateLegalLocations() []Location {
	var locs []Location
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			l := Location{x: x, y: y}
			if b.IsLegal(l.x, l.y) {
				locs = append(locs, l)
			}
		}
	}
	return locs
}

func (b *Board) Put(x int, y int) {
	var legalFlags [8]bool
	for i := 0; i < 8; i++ {
		legalFlags[i] = b.isLegal(x, y, i)
	}
	b.board[x][y] = b.currentColor
	b.counts[b.currentColor]++
	var locs []Location
	rec := &Record{color: b.currentColor, placedPiece: Location{x: x, y: y}, flippedPieces: locs}

	b.records = append(b.records, rec)
	opp := Flip(b.currentColor)
	for i := 0; i < 8; i++ {
		if legalFlags[i] {
			d := DIRECTIONS[i]
			for j := 1; j < 8; j++ {
				x := x + d.x*j
				y := y + d.y*j
				if b.board[x][y] == b.currentColor {
					break
				}
				b.board[x][y] = b.currentColor
				b.counts[b.currentColor]++
				b.counts[opp]--
				rec.flip(Location{x: x, y: y})
			}
		}
	}
	b.currentColor = Flip(b.currentColor)
}

func (b *Board) Pass() {
	b.records = append(b.records, &Record{})
	b.currentColor = Flip(b.currentColor)
}

func (b *Board) Undo() {
	rec := b.records[len(b.records)-1]
	b.records = b.records[:len(b.records)-1]
	if rec.flippedPieces != nil {
		b.board[rec.placedPiece.x][rec.placedPiece.y] = -1
		opp := (rec.color + 1) % 2
		flippedPieceCount := len(rec.flippedPieces)
		for i := 0; i < flippedPieceCount; i++ {
			p := rec.flippedPieces[i]
			b.board[p.x][p.y] = opp
		}
		b.counts[rec.color] -= 1 + flippedPieceCount
		b.counts[opp] += flippedPieceCount
	}
	b.currentColor = Flip(b.currentColor)
}
