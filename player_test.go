package main

import (
	"testing"
)

func TestParseTwoSquaresInvalid(t *testing.T) {
	for _, test := range []string{
		"",
		"e",
		"e1",
		"e1g",
		"e1f1g1",
		"abcd",
		"a0a1",
		"i1h1",
		"a1a0",
		"a8a9",
	} {
		if from, to, err := ParseTwoSquares(test); err == nil {
			t.Errorf("ParseTwoSquares(%v): got %v, %v, nil; want error", test, from, to)
		}
	}
}

func TestParseTwoSquaresValid(t *testing.T) {
	for from := a1; from <= h8; from++ {
		for to := a1; to <= h8; to++ {
			s := from.String() + to.String()
			if gotfrom, gotto, err := ParseTwoSquares(s); err != nil || gotfrom != from || gotto != to {
				t.Errorf("ParseTwoSquares(%v): got %v, %v, %v, want %v, %v, nil", s, gotfrom, gotto, err, from, to)
			}
		}
	}
}

func TestParseInputInvalid(t *testing.T) {
	h := new(Human)
	fen := "4k3/7P/8/b7/8/8/3N4/R3K3 w Q - 0 1"
	pos, err := ParseFEN(fen)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []string{
		"",
		"r",
		"e",
		"e1",
		"e1c",
		"e1d1c1",
		"a1a9",  // no such square
		"d4d5",  // no piece
		"a5b4",  // opponent's piece
		"a1e1",  // capture own piece
		"a1h8",  // violates piece movement rule
		"a1a6",  // not pseudolegal in pos
		"d2c4",  // illegal in pos
		"h7h8z", // no such promotion piece
		"h7h8k", // prohibited promotion piece
		"h7h8",  // promotion piece not specified
	} {
		if m, err := h.parseInput(pos, test); err == nil {
			t.Errorf("parseInput(%v, %v): got %v, nil; want error", InitialPositionFEN, test, m)
		}
	}
}

func TestParseInputValid(t *testing.T) {
	h := new(Human)
	fen := "4k3/7P/8/b7/8/8/3N4/R3K3 w Q - 0 1"
	pos, err := ParseFEN(fen)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		s   string
		m   Move
		err error
	}{
		{"go", Move{}, errGo},
		{"resign", Move{}, errResign},
		{"e1e2", Move{From: e1, To: e2, Piece: King}, nil},
		{"e1c1", Move{From: e1, To: c1, Piece: King}, nil},
		{"a1a5", Move{From: a1, To: a5, Piece: Rook, CapturePiece: Bishop}, nil},
		{"h7h8q", Move{From: h7, To: h8, Piece: Pawn, PromotePiece: Queen}, nil},
		{"h7h8r", Move{From: h7, To: h8, Piece: Pawn, PromotePiece: Rook}, nil},
		{"h7h8b", Move{From: h7, To: h8, Piece: Pawn, PromotePiece: Bishop}, nil},
		{"h7h8n", Move{From: h7, To: h8, Piece: Pawn, PromotePiece: Knight}, nil},
	} {
		if m, err := h.parseInput(pos, test.s); m != test.m || err != test.err {
			t.Errorf("parseInput(%v, %v): got %v, %v; want %v, %v", fen, test.s, m, err, test.m, test.err)
		}
	}
}
