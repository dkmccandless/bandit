package main

import (
	"fmt"
	"testing"
)

var algebraicTests = []struct {
	fen  string
	move Move
	alg  string
	long string
}{
	{InitialPositionFEN, Move{From: e2, To: e4, Piece: Pawn}, "e4", "e2-e4"},
	{InitialPositionFEN, Move{From: g2, To: g3, Piece: Pawn}, "g3", "g2-g3"},
	{InitialPositionFEN, Move{From: g1, To: f3, Piece: Knight}, "Nf3", "Ng1-f3"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e1, To: f2, Piece: King}, "Kf2", "Ke1-f2"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: g2, To: g4, Piece: Pawn}, "g4", "g2-g4"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: d5, To: e6, Piece: Pawn, CapturePiece: Pawn, EP: true}, "dxe6", "d5xe6"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: a1, To: a8, Piece: Rook, CapturePiece: Rook}, "Rxa8+", "Ra1xa8"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: b7, To: b8, Piece: Pawn, PromotePiece: Queen}, "b8Q+", "b7-b8Q"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: b7, To: a8, Piece: Pawn, CapturePiece: Rook, PromotePiece: Queen}, "bxa8Q+", "b7xa8Q"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e1, To: c1, Piece: King}, "O-O-O", "O-O-O"},
	{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e1, To: g1, Piece: King}, "O-O+", "O-O"},
	{"r4k2/1P6/3P2Q1/4p3/8/8/6P1/R3K2R w KQ - 0 1", Move{From: e1, To: g1, Piece: King}, "O-O#", "O-O"},
	{"7k/8/8/8/8/8/8/R4RK1 w - - 0 1", Move{From: a1, To: d1, Piece: Rook}, "Rad1", "Ra1-d1"},
	{"7k/R7/8/8/8/8/8/R5K1 w - - 0 1", Move{From: a7, To: a6, Piece: Rook}, "R7a6", "Ra7-a6"},
	{"7k/R7/8/8/8/8/8/R5K1 w - - 0 1", Move{From: a7, To: a8, Piece: Rook}, "Ra8+", "Ra7-a8"},
	{"8/B7/8/8/8/6k1/6P1/B5K1 w - - 0 1", Move{From: a7, To: d4, Piece: Bishop}, "B7d4", "Ba7-d4"},
	{"8/B5B1/8/8/8/6k1/6P1/B5K1 w - - 0 1", Move{From: a1, To: d4, Piece: Bishop}, "B1d4", "Ba1-d4"},
	{"8/B5B1/8/8/8/6k1/6P1/B5K1 w - - 0 1", Move{From: g7, To: d4, Piece: Bishop}, "Bgd4", "Bg7-d4"},
	{"8/B5B1/8/8/8/6k1/6P1/B5K1 w - - 0 1", Move{From: a7, To: d4, Piece: Bishop}, "Ba7d4", "Ba7-d4"},

	// Cases with pinned pieces that should not be disambiguated
	{"7k/8/8/8/8/8/R7/KR5r w - - 0 1", Move{From: a2, To: b2, Piece: Rook}, "Rb2", "Ra2-b2"},
	{"r6k/8/8/8/8/8/R7/KR6 w - - 0 1", Move{From: b1, To: b2, Piece: Rook}, "Rb2", "Rb1-b2"},
	{"7k/8/8/8/8/1R6/R7/KR5r w - - 0 1", Move{From: a2, To: b2, Piece: Rook}, "Rab2", "Ra2-b2"},
	{"7k/8/8/8/8/1R6/R7/KR5r w - - 0 1", Move{From: b3, To: b2, Piece: Rook}, "Rbb2", "Rb3-b2"},
	{"r6k/8/8/8/8/1R6/R7/KR6 w - - 0 1", Move{From: b1, To: b2, Piece: Rook}, "R1b2", "Rb1-b2"},
	{"r6k/8/8/8/8/1R6/R7/KR6 w - - 0 1", Move{From: b3, To: b2, Piece: Rook}, "R3b2", "Rb3-b2"},
	{"5b1k/4Q3/2Q5/2K2Q1r/1Q6/b3Q3/2Q5/2r3b1 w - - 0 1", Move{From: c6, To: e4, Piece: Queen}, "Qe4", "Qc6-e4"},
	{"2r2b1k/4Q3/2Q5/2K2Q1r/1Q6/b3Q3/2Q5/2r5 w - - 0 1", Move{From: e3, To: e4, Piece: Queen}, "Qe4", "Qe3-e4"},
	{"5b1k/4Q3/2Q5/2K2Q1r/1Q6/b3Q3/2Q5/2r5 w - - 0 1", Move{From: e3, To: e4, Piece: Queen}, "Qee4", "Qe3-e4"},
	{"7k/4Q3/2Q5/2K2Q1r/1Q6/b3Q3/2Q5/2r5 w - - 0 1", Move{From: c6, To: e4, Piece: Queen}, "Qce4", "Qc6-e4"},
	{"7k/4Q3/2Q5/2K2Q1r/1Q6/b3Q3/2Q5/2r5 w - - 0 1", Move{From: e3, To: e4, Piece: Queen}, "Q3e4", "Qe3-e4"},
	{"7k/4Q3/2Q5/2K2Q1r/1Q6/b3Q3/2Q5/2r5 w - - 0 1", Move{From: e7, To: e4, Piece: Queen}, "Q7e4", "Qe7-e4"},
}

func TestLongAlgebraic(t *testing.T) {
	for _, test := range algebraicTests {
		if got := LongAlgebraic(test.move); got != test.long {
			t.Errorf("LongAlgebraic(%v): got %v, want %v", test.move, got, test.long)
		}
	}
}

func TestAlgebraic(t *testing.T) {
	for _, test := range algebraicTests {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := Algebraic(pos, test.move); got != test.alg {
			t.Errorf("Algebraic(%v, %+v): got %v, want %v", test.fen, test.move, got, test.alg)
		}
	}
}

var moveAlgebraicTests = []struct {
	fen    string
	num    string
	move   Move
	numalg string
}{
	{InitialPositionFEN, "1.", Move{From: g1, To: f3, Piece: Knight}, "1.Nf3"},
	{"r1bqkbnr/1ppp1ppp/p1n5/4p3/B3P3/5N2/PPPP1PPP/RNBQK2R b - - 0 4", "4...", Move{From: b7, To: b5, Piece: Pawn}, "4...b5"},
	{"r1b2rkB/1pp1ppbp/2n3p1/8/PpP3nP/8/3Kpp2/1N3BNR w - - 0 13", "13.", Move{From: c4, To: c5, Piece: Pawn}, "13.c5"},
	{"r3k2K/8/8/5nn1/8/8/8/8 b q - 13 54", "54...", Move{From: e8, To: c8, Piece: King}, "54...O-O-O#"},
}

func TestMoveNumber(t *testing.T) {
	for _, test := range moveAlgebraicTests {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := moveNumber(pos); got != test.num {
			t.Errorf("moveNumber(%v): got %v, want %v", test.fen, got, test.num)
		}
	}
}

func TestNumberedAlgebraic(t *testing.T) {
	for _, test := range moveAlgebraicTests {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := numberedAlgebraic(pos, test.move); got != test.numalg {
			t.Errorf("numberedAlgebraic(%v, %v): got %v, want %v", test.fen, test.move, got, test.num)
		}
	}
}

func TestText(t *testing.T) {
	for _, test := range []struct {
		fen   string
		moves []Move
		want  string
	}{
		{InitialPositionFEN, []Move{}, ""},
		{InitialPositionFEN, []Move{
			Move{From: e2, To: e4, Piece: Pawn},
		}, "1.e4"},
		{InitialPositionFEN, []Move{
			Move{From: e2, To: e4, Piece: Pawn},
			Move{From: c7, To: c5, Piece: Pawn},
		}, "1.e4 c5"},
		{InitialPositionFEN, []Move{
			Move{From: e2, To: e4, Piece: Pawn},
			Move{From: c7, To: c5, Piece: Pawn},
			Move{From: g1, To: f3, Piece: Knight},
		}, "1.e4 c5 2.Nf3"},
		{"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b - - 0 1", []Move{
			Move{From: g8, To: f6, Piece: Knight},
			Move{From: c2, To: c4, Piece: Pawn},
			Move{From: g7, To: g6, Piece: Pawn},
		}, "1...Nf6 2.c4 g6"},
		{"r3k2K/8/8/5nn1/8/8/8/8 b q - 13 54", []Move{
			Move{From: e8, To: c8, Piece: King},
		}, "54...O-O-O#"},
		{InitialPositionFEN, []Move{
			Move{From: c2, To: c4, Piece: Pawn},
			Move{From: g8, To: f6, Piece: Knight},
			Move{From: g1, To: f3, Piece: Knight},
			Move{From: g7, To: g6, Piece: Pawn},
			Move{From: f3, To: g1, Piece: Knight},
			Move{From: f8, To: g7, Piece: Bishop},
			Move{From: d1, To: a4, Piece: Queen},
			Move{From: e8, To: g8, Piece: King},
			Move{From: a4, To: d7, Piece: Queen, CapturePiece: Pawn},
			Move{From: d8, To: d7, Piece: Queen, CapturePiece: Queen},
			Move{From: g2, To: g4, Piece: Pawn},
			Move{From: d7, To: d2, Piece: Queen, CapturePiece: Pawn},
			Move{From: e1, To: d2, Piece: King, CapturePiece: Queen},
			Move{From: f6, To: g4, Piece: Knight, CapturePiece: Pawn},
			Move{From: b2, To: b4, Piece: Pawn},
			Move{From: a7, To: a5, Piece: Pawn},
			Move{From: a2, To: a4, Piece: Pawn},
			Move{From: g7, To: a1, Piece: Bishop, CapturePiece: Rook},
			Move{From: c1, To: b2, Piece: Bishop},
			Move{From: b8, To: c6, Piece: Knight},
			Move{From: b2, To: h8, Piece: Bishop},
			Move{From: a1, To: g7, Piece: Bishop},
			Move{From: h2, To: h4, Piece: Pawn},
			Move{From: a5, To: b4, Piece: Pawn, CapturePiece: Pawn},
		}, "1.c4 Nf6 2.Nf3 g6 3.Ng1 Bg7 4.Qa4 O-O 5.Qxd7 Qxd7 6.g4 Qxd2+ 7.Kxd2 Nxg4 8.b4 a5 9.a4 Bxa1 10.Bb2 Nc6 11.Bh8 Bg7 12.h4 axb4"},
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := Text(pos, test.moves); got != test.want {
			var s string
			for _, m := range test.moves {
				s += fmt.Sprintf("%v ", m)
			}
			t.Errorf("Text(%v, %v): got %v, want %v", test.fen, s, got, test.want)
		}
	}
}
