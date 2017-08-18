package main

import "testing"

var positionTests = []struct {
	fen string
	pos Position
}{
	{InitialPositionFEN, InitialPosition},
	{"8/8/8/8/8/8/8/K1k5 b - - 13 71", Position{
		ToMove:     Black,
		HalfMove:   13,
		FullMove:   71,
		KingSquare: [2]Square{a1, c1},
		b: [2][8]Board{
			{0, 0, 0, 0, 0, 0, a1.Board(), a1.Board()},
			{0, 0, 0, 0, 0, 0, c1.Board(), c1.Board()},
		},
		z: pieceZobrist[White][King][a1] ^ pieceZobrist[Black][King][c1] ^ blackToMoveZobrist,
	}},
	{"r3k2r/8/8/3p4/3P4/8/8/R3K2R w Qkq d6 0 1", Position{
		QSCastle:   [2]bool{true, true},
		KSCastle:   [2]bool{false, true},
		ep:         d6,
		ToMove:     White,
		HalfMove:   0,
		FullMove:   1,
		KingSquare: [2]Square{e1, e8},
		b: [2][8]Board{
			{0, d4.Board(), 0, 0, a1.Board() ^ h1.Board(), 0, e1.Board(), a1.Board() ^ e1.Board() ^ h1.Board() ^ d4.Board()},
			{0, d5.Board(), 0, 0, a8.Board() ^ h8.Board(), 0, e8.Board(), d5.Board() ^ a8.Board() ^ e8.Board() ^ h8.Board()},
		},
		z: pieceZobrist[White][Rook][a1] ^ pieceZobrist[White][King][e1] ^ pieceZobrist[White][Rook][h1] ^ pieceZobrist[White][Pawn][d4] ^
			pieceZobrist[Black][Pawn][d5] ^ pieceZobrist[Black][Rook][a8] ^ pieceZobrist[Black][King][e8] ^ pieceZobrist[Black][Rook][h8] ^
			qsCastleZobrist[White] ^ qsCastleZobrist[Black] ^ ksCastleZobrist[Black], // no pawn can capture en passant
	}},
}

func TestParseFENInvalid(t *testing.T) {
	for _, test := range []string{
		"",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w",              // not enough fields
		"rnbqkbnr/pppppppp/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",     // not enough rows
		"rnbqkbnr/pppppppp/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", // too many rows
		"rnbqkbnr/pppppppp/8/8/8/7/PPPPPPPP/RNBQKBNR w KQkq - 0 1",   // not enough squares in row (empty)
		"rnbqkbnr/pppppppp/8/8/8/9/PPPPPPPP/RNBQKBNR w KQkq - 0 1",   // too many squares in row (empty)
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP/RNBQKBNR w KQkq - 0 1",    // not enough squares in row (pieces)
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPPP/RNBQKBNR w KQkq - 0 1",  // too many squares in row (pieces)
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQQBNR w KQkq - 0 1",   // not enough white kings
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKKBNR w KQkq - 0 1",   // too many white kings
		"rnbqqbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",   // not enough black kings
		"rnbkkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",   // too many black kings
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR g KQkq - 0 1",   // invalid character in active player field
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w AQKJ - 0 1",   // invalid character in castling field
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e 0 1",   // invalid en passant square (no rank)
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq 6 0 1",   // invalid en passant square (no file)
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e0 0 1",  // invalid en passant square (rank out of range)
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq j6 0 1",  // invalid en passant square (file out of range)
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e5 0 1",  // invalid en passant square (wrong rank)
	} {
		if pos, err := ParseFEN(test); err == nil {
			t.Errorf("ParseFEN(%v): got %v, %v; want error", test, pos, err)
		}
	}
}

func TestParseFENValid(t *testing.T) {
	for _, test := range positionTests {
		if got, err := ParseFEN(test.fen); err != nil {
			t.Errorf("ParseFEN(%v): got %+v, %v; want %+v, nil", test.fen, got, err, test.pos)
		}
	}
}

func TestParseSquareInvalid(t *testing.T) {
	for _, test := range []string{
		"",
		"-",
		"c",
		"7",
		"c0",
		"j7",
		"7c",
		"c7 ",
	} {
		if sq, err := ParseSquare(test); err == nil {
			t.Errorf("ParseSquare(%v): got %v, %v; want error", test, sq, err)
		}
	}
}

func TestParseSquareValid(t *testing.T) {
	for _, test := range []struct {
		s    string
		want Square
	}{
		{"a1", a1},
		{"b1", b1},
		{"c1", c1},
		{"d1", d1},
		{"e1", e1},
		{"f1", f1},
		{"g1", g1},
		{"h1", h1},
		{"a2", a2},
		{"b2", b2},
		{"c2", c2},
		{"d2", d2},
		{"e2", e2},
		{"f2", f2},
		{"g2", g2},
		{"h2", h2},
		{"a3", a3},
		{"b3", b3},
		{"c3", c3},
		{"d3", d3},
		{"e3", e3},
		{"f3", f3},
		{"g3", g3},
		{"h3", h3},
		{"a4", a4},
		{"b4", b4},
		{"c4", c4},
		{"d4", d4},
		{"e4", e4},
		{"f4", f4},
		{"g4", g4},
		{"h4", h4},
		{"a5", a5},
		{"b5", b5},
		{"c5", c5},
		{"d5", d5},
		{"e5", e5},
		{"f5", f5},
		{"g5", g5},
		{"h5", h5},
		{"a6", a6},
		{"b6", b6},
		{"c6", c6},
		{"d6", d6},
		{"e6", e6},
		{"f6", f6},
		{"g6", g6},
		{"h6", h6},
		{"a7", a7},
		{"b7", b7},
		{"c7", c7},
		{"d7", d7},
		{"e7", e7},
		{"f7", f7},
		{"g7", g7},
		{"h7", h7},
		{"a8", a8},
		{"b8", b8},
		{"c8", c8},
		{"d8", d8},
		{"e8", e8},
		{"f8", f8},
		{"g8", g8},
		{"h8", h8},
	} {
		if got, err := ParseSquare(test.s); err != nil {
			t.Errorf("ParseSquare(%v): got %v, %v; want %v, nil", test.s, got, err, test.want.String())
		}
	}
}

func TestFEN(t *testing.T) {
	for _, test := range positionTests {
		if got := FEN(test.pos); got != test.fen {
			t.Errorf("FEN(%+v): got %v, want %v", test.pos, got, test.fen)
		}
	}
}
