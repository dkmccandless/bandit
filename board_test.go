package main

import "testing"

func TestSquareRank(t *testing.T) {
	ranks := []byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3,
		4, 4, 4, 4, 4, 4, 4, 4,
		5, 5, 5, 5, 5, 5, 5, 5,
		6, 6, 6, 6, 6, 6, 6, 6,
		7, 7, 7, 7, 7, 7, 7, 7,
	}
	for _, v := range bs {
		if got := v.s.Rank(); got != ranks[v.s] {
			t.Errorf("Square(%v).Rank: got %v, want %v", v.s, got, ranks[v.s])
		}
	}
}

func TestSquareFile(t *testing.T) {
	files := []byte{
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
	}
	for _, v := range bs {
		if got := v.s.File(); got != files[v.s] {
			t.Errorf("Square(%v).File: got %v, want %v", v.s, got, files[v.s])
		}
	}
}

func TestSquareDiagonal(t *testing.T) {
	diagonals := []byte{
		7, 6, 5, 4, 3, 2, 1, 0,
		8, 7, 6, 5, 4, 3, 2, 1,
		9, 8, 7, 6, 5, 4, 3, 2,
		10, 9, 8, 7, 6, 5, 4, 3,
		11, 10, 9, 8, 7, 6, 5, 4,
		12, 11, 10, 9, 8, 7, 6, 5,
		13, 12, 11, 10, 9, 8, 7, 6,
		14, 13, 12, 11, 10, 9, 8, 7,
	}
	for _, v := range bs {
		if got := v.s.Diagonal(); got != diagonals[v.s] {
			t.Errorf("Square(%v).Diagonal: got %v, want %v", v.s, got, diagonals[v.s])
		}
	}
}

func TestSquareAntiDiagonal(t *testing.T) {
	antiDiagonals := []byte{
		0, 1, 2, 3, 4, 5, 6, 7,
		1, 2, 3, 4, 5, 6, 7, 8,
		2, 3, 4, 5, 6, 7, 8, 9,
		3, 4, 5, 6, 7, 8, 9, 10,
		4, 5, 6, 7, 8, 9, 10, 11,
		5, 6, 7, 8, 9, 10, 11, 12,
		6, 7, 8, 9, 10, 11, 12, 13,
		7, 8, 9, 10, 11, 12, 13, 14,
	}
	for _, v := range bs {
		if got := v.s.AntiDiagonal(); got != antiDiagonals[v.s] {
			t.Errorf("Square(%v).AntiDiagonal: got %v, want %v", v.s, got, antiDiagonals[v.s])
		}
	}
}

func TestSquareToBoard(t *testing.T) {
	for _, v := range bs {
		if got := v.s.Board(); got != v.b {
			t.Errorf("Square(%v).Board: got %016x, want %016x", v.s, got, v.b)
		}
	}
}

func TestLS1B(t *testing.T) {
	for _, test := range ls1bBoards {
		if got := LS1B(test.input); got != test.ls1b {
			t.Errorf("LS1B(%016x): got %016x, want %016x", test.input, got, test.ls1b)
		}
	}
}

func TestResetLS1B(t *testing.T) {
	for _, test := range ls1bBoards {
		if got := ResetLS1B(test.input); got != test.input-test.ls1b {
			t.Errorf("ResetLS1B(%016x): got %016x, want %016x", test.input, got, test.input-test.ls1b)
		}
	}
}

func TestLS1BIndex(t *testing.T) {
	for _, v := range bs {
		if got := LS1BIndex(v.b); got != v.s {
			t.Errorf("LS1BIndex(%016x): got %v, want %v", v.b, got, v.s)
		}
	}
}

func TestPieceOn(t *testing.T) {
	for _, test := range []struct {
		s  Square
		c  Color
		p  Piece
		ok bool
	}{
		{a1, White, Rook, true},
		{b1, White, Knight, true},
		{c1, White, Bishop, true},
		{d1, White, Queen, true},
		{e1, White, King, true},
		{f1, White, Bishop, true},
		{g1, White, Knight, true},
		{h1, White, Rook, true},
		{a2, White, Pawn, true},
		{b2, White, Pawn, true},
		{c2, White, Pawn, true},
		{d2, White, Pawn, true},
		{e2, White, Pawn, true},
		{f2, White, Pawn, true},
		{g2, White, Pawn, true},
		{h2, White, Pawn, true},
		{a3, 0, 0, false},
		{b3, 0, 0, false},
		{c3, 0, 0, false},
		{d3, 0, 0, false},
		{e3, 0, 0, false},
		{f3, 0, 0, false},
		{g3, 0, 0, false},
		{h3, 0, 0, false},
		{a4, 0, 0, false},
		{b4, 0, 0, false},
		{c4, 0, 0, false},
		{d4, 0, 0, false},
		{e4, 0, 0, false},
		{f4, 0, 0, false},
		{g4, 0, 0, false},
		{h4, 0, 0, false},
		{a5, 0, 0, false},
		{b5, 0, 0, false},
		{c5, 0, 0, false},
		{d5, 0, 0, false},
		{e5, 0, 0, false},
		{f5, 0, 0, false},
		{g5, 0, 0, false},
		{h5, 0, 0, false},
		{a6, 0, 0, false},
		{b6, 0, 0, false},
		{c6, 0, 0, false},
		{d6, 0, 0, false},
		{e6, 0, 0, false},
		{f6, 0, 0, false},
		{g6, 0, 0, false},
		{h6, 0, 0, false},
		{a7, Black, Pawn, true},
		{b7, Black, Pawn, true},
		{c7, Black, Pawn, true},
		{d7, Black, Pawn, true},
		{e7, Black, Pawn, true},
		{f7, Black, Pawn, true},
		{g7, Black, Pawn, true},
		{h7, Black, Pawn, true},
		{a8, Black, Rook, true},
		{b8, Black, Knight, true},
		{c8, Black, Bishop, true},
		{d8, Black, Queen, true},
		{e8, Black, King, true},
		{f8, Black, Bishop, true},
		{g8, Black, Knight, true},
		{h8, Black, Rook, true},
	} {
		gotColor, gotPiece, gotOK := InitialPosition.PieceOn(test.s)
		if gotColor != test.c {
			t.Errorf("PieceOn(%v) in Position %+v: got Color %v, want %v", test.s, InitialPosition, gotColor, test.c)
		}
		if gotPiece != test.p {
			t.Errorf("PieceOn(%v) in Position %+v: got Piece %v, want %v", test.s, InitialPosition, gotPiece, test.p)
		}
		if gotOK != test.ok {
			t.Errorf("PieceOn(%v) in Position %+v: got bool %v, want %v", test.s, InitialPosition, gotOK, test.ok)
		}
	}
}

var bs = []struct {
	b Board
	s Square
}{
	{1 << 0, a1},
	{1 << 1, b1},
	{1 << 2, c1},
	{1 << 3, d1},
	{1 << 4, e1},
	{1 << 5, f1},
	{1 << 6, g1},
	{1 << 7, h1},
	{1 << 8, a2},
	{1 << 9, b2},
	{1 << 10, c2},
	{1 << 11, d2},
	{1 << 12, e2},
	{1 << 13, f2},
	{1 << 14, g2},
	{1 << 15, h2},
	{1 << 16, a3},
	{1 << 17, b3},
	{1 << 18, c3},
	{1 << 19, d3},
	{1 << 20, e3},
	{1 << 21, f3},
	{1 << 22, g3},
	{1 << 23, h3},
	{1 << 24, a4},
	{1 << 25, b4},
	{1 << 26, c4},
	{1 << 27, d4},
	{1 << 28, e4},
	{1 << 29, f4},
	{1 << 30, g4},
	{1 << 31, h4},
	{1 << 32, a5},
	{1 << 33, b5},
	{1 << 34, c5},
	{1 << 35, d5},
	{1 << 36, e5},
	{1 << 37, f5},
	{1 << 38, g5},
	{1 << 39, h5},
	{1 << 40, a6},
	{1 << 41, b6},
	{1 << 42, c6},
	{1 << 43, d6},
	{1 << 44, e6},
	{1 << 45, f6},
	{1 << 46, g6},
	{1 << 47, h6},
	{1 << 48, a7},
	{1 << 49, b7},
	{1 << 50, c7},
	{1 << 51, d7},
	{1 << 52, e7},
	{1 << 53, f7},
	{1 << 54, g7},
	{1 << 55, h7},
	{1 << 56, a8},
	{1 << 57, b8},
	{1 << 58, c8},
	{1 << 59, d8},
	{1 << 60, e8},
	{1 << 61, f8},
	{1 << 62, g8},
	{1 << 63, h8},
}

var ls1bBoards = []struct {
	input Board
	ls1b  Board
}{
	{0x10, 0x10},
	{0x12, 0x2},
	{0x48c289e000, 0x2000},
	{0xac << 30, 1 << 32},
	{InitialPosition.b[White][Bishop], 0x4},
	{InitialPosition.b[White][Pawn], 0x100},
	{InitialPosition.b[Black][All], 1 << 48},
	{InitialPosition.b[Black][Knight], 1 << 57},
}
