package main

import "testing"

func TestIsPseudoLegal(t *testing.T) {
	for _, test := range []struct {
		fen  string
		move Move
		want bool
	}{
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e1, To: f2, Piece: King}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: g2, To: g4, Piece: Pawn}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: d5, To: e6, Piece: Pawn, CapturePiece: Pawn, CaptureSquare: e5}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: a1, To: a8, Piece: Rook, CapturePiece: Rook, CaptureSquare: a8}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: b7, To: b8, Piece: Pawn, PromotePiece: Queen}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: b7, To: a8, Piece: Pawn, CapturePiece: Rook, CaptureSquare: a8, PromotePiece: Queen}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e1, To: c1, Piece: King}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e1, To: g1, Piece: King}, true},
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e2, To: e4, Piece: Pawn}, false},                       // no piece on From square
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: e5, To: e4, Piece: Pawn}, false},                       // piece of the wrong color
		{"r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1", Move{From: a1, To: h8, Piece: Rook}, false},                       // piece can't reach To square
		{"r4k2/1P6/8/3Pp3/8/8/6B1/R3K2R w KQ e6 0 1", Move{From: g2, To: h1, Piece: Bishop, CapturePiece: Rook}, false}, // capture own color piece
		{"r4k2/1P6/8/3Pp3/8/8/1b4P1/R3K2R w KQ e6 0 1", Move{From: e1, To: e2, Piece: King}, true},                      // puts king in check
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := IsPseudoLegal(pos, test.move); got != test.want {
			t.Errorf("TestIsPseudoLegal(%v, %+v): got %v, want %v", test.fen, test.move, got, test.want)
		}
	}
}

func TestIsLegal(t *testing.T) {
	for _, test := range []struct {
		fen  string
		want bool
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", true},  // initial position
		{"r1b2rkB/1pp1ppbp/2n3p1/8/PpP3nP/8/3Kpp2/1N3BNR w - - 0 13", true}, // legal, not check

		{"8/8/8/8/8/8/q7/K6k w - - 0 1", true},                // check with forced capture
		{"8/8/8/8/8/q7/8/KB5k w - - 0 1", true},               // check with forced block
		{"8/8/8/8/8/q7/8/K6k w - - 0 1", true},                // check with forced evasion
		{"8/8/8/8/8/8/2b5/k2R3K b - - 0 1", true},             // check with capture, block, or evasion
		{"7k/8/5B2/8/8/8/8/K6R b - - 0 1", true},              // double check
		{"6qk/8/5B2/8/8/8/8/K6R b - - 0 1", true},             // checkmate by double check, either can be blocked
		{"R6k/6pp/8/8/8/8/8/K7 b - - 0 1", true},              // checkmate, back rank
		{"6rk/5Npp/8/8/8/8/8/7K b - - 0 1", true},             // checkmate, smothered
		{"8/8/8/8/8/8/q7/K6k b - - 0 1", false},               // illegal, side to move not in check but opponent in check
		{"K7/2n5/8/8/8/8/5N2/7k b - - 0 1", false},            // illegal, both sides in check
		{"4R2k/6pp/8/8/8/8/8/q6K b - - 0 1", false},           // illegal, last move delivered checkmate by putting own king in check
		{"k6R/pp6/8/8/8/8/6PP/r6K w - - 0 1", false},          // illegal, both sides in checkmate
		{"8/8/7p/6pP/5pP1/3kpPp1/4P1PN/6NK w - - 0 1", false}, // illegal, only pseudo-legal move is to capture the king
		{"K6k/2q5/8/8/8/8/8/8 w - - 0 1", true},               // stalemate
		{"K6k/2q5/8/8/8/8/8/8 b - - 0 1", true},               // would be stalemate if it were the opponent's turn
		{"8/8/6K1/8/8/8/2k5/8 b - - 0 1", true},               // legal, lone kings
		{"8/8/8/3k4/3K4/8/8/8 w - - 0 1", false},              // illegal, lone kings
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := IsLegal(pos); got != test.want {
			t.Errorf("TestIsLegal(%v): got %v, want %v", test.fen, got, test.want)
		}
	}
}

func TestIsCheck(t *testing.T) {
	for _, test := range []struct {
		fen  string
		want bool
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},  // initial position
		{"r1b2rkB/1pp1ppbp/2n3p1/8/PpP3nP/8/3Kpp2/1N3BNR w - - 0 13", false}, // legal, not check

		{"8/8/8/8/8/8/q7/K6k w - - 0 1", true},                // check with forced capture
		{"8/8/8/8/8/q7/8/KB5k w - - 0 1", true},               // check with forced block
		{"8/8/8/8/8/q7/8/K6k w - - 0 1", true},                // check with forced evasion
		{"8/8/8/8/8/8/2b5/k2R3K b - - 0 1", true},             // check with capture, block, or evasion
		{"7k/8/5B2/8/8/8/8/K6R b - - 0 1", true},              // double check
		{"6qk/8/5B2/8/8/8/8/K6R b - - 0 1", true},             // checkmate by double check, either can be blocked
		{"R6k/6pp/8/8/8/8/8/K7 b - - 0 1", true},              // checkmate, back rank
		{"6rk/5Npp/8/8/8/8/8/7K b - - 0 1", true},             // checkmate, smothered
		{"8/8/8/8/8/8/q7/K6k b - - 0 1", false},               // illegal, side to move not in check but opponent in check
		{"K7/2n5/8/8/8/8/5N2/7k b - - 0 1", true},             // illegal, both sides in check
		{"4R2k/6pp/8/8/8/8/8/q6K b - - 0 1", true},            // illegal, last move delivered checkmate by putting own king in check
		{"k6R/pp6/8/8/8/8/6PP/r6K w - - 0 1", true},           // illegal, both sides in checkmate
		{"8/8/7p/6pP/5pP1/3kpPp1/4P1PN/6NK w - - 0 1", false}, // illegal, only pseudo-legal move is to capture the king
		{"K6k/2q5/8/8/8/8/8/8 w - - 0 1", false},              // stalemate
		{"K6k/2q5/8/8/8/8/8/8 b - - 0 1", false},              // would be stalemate if it were the opponent's turn
		{"8/8/6K1/8/8/8/2k5/8 b - - 0 1", false},              // legal, lone kings
		{"8/8/8/3k4/3K4/8/8/8 w - - 0 1", true},               // illegal, lone kings
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := IsCheck(pos); got != test.want {
			t.Errorf("TestIsCheck(%v): got %v, want %v", test.fen, got, test.want)
		}
	}
}

func TestIsTerminal(t *testing.T) {
	for _, test := range []struct {
		fen  string
		want bool
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},  // initial position
		{"r1b2rkB/1pp1ppbp/2n3p1/8/PpP3nP/8/3Kpp2/1N3BNR w - - 0 13", false}, // legal, not check

		{"8/8/8/8/8/8/q7/K6k w - - 0 1", false},               // check with forced capture
		{"8/8/8/8/8/q7/8/KB5k w - - 0 1", false},              // check with forced block
		{"8/8/8/8/8/q7/8/K6k w - - 0 1", false},               // check with forced evasion
		{"8/8/8/8/8/8/2b5/k2R3K b - - 0 1", false},            // check with capture, block, or evasion
		{"7k/8/5B2/8/8/8/8/K6R b - - 0 1", false},             // double check
		{"6qk/8/5B2/8/8/8/8/K6R b - - 0 1", true},             // checkmate by double check, either can be blocked
		{"R6k/6pp/8/8/8/8/8/K7 b - - 0 1", true},              // checkmate, back rank
		{"6rk/5Npp/8/8/8/8/8/7K b - - 0 1", true},             // checkmate, smothered
		{"8/8/8/8/8/8/q7/K6k b - - 0 1", false},               // illegal, side to move not in check but opponent in check
		{"K7/2n5/8/8/8/8/5N2/7k b - - 0 1", false},            // illegal, both sides in check
		{"4R2k/6pp/8/8/8/8/8/q6K b - - 0 1", true},            // illegal, last move delivered checkmate by putting own king in check
		{"k6R/pp6/8/8/8/8/6PP/r6K w - - 0 1", true},           // illegal, both sides in checkmate
		{"8/8/7p/6pP/5pP1/3kpPp1/4P1PN/6NK w - - 0 1", false}, // illegal, only pseudo-legal move is to capture the king
		{"K6k/2q5/8/8/8/8/8/8 w - - 0 1", true},               // stalemate
		{"K6k/2q5/8/8/8/8/8/8 b - - 0 1", false},              // would be stalemate if it were the opponent's turn
		{"8/8/6K1/8/8/8/2k5/8 b - - 0 1", false},              // legal, lone kings
		{"8/8/8/3k4/3K4/8/8/8 w - - 0 1", false},              // illegal, lone kings
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := IsTerminal(pos); got != test.want {
			t.Errorf("TestIsTerminal(%v): got %v, want %v", test.fen, got, test.want)
		}
	}
}

func BenchmarkSearchPosition(b *testing.B) {
	pos, err := ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = SearchPosition(pos, 2)
	}
}

var windowTests = []struct {
	w, neg          Window
	n               int
	c               Window
	constrained, ok bool
}{
	{Window{-50, -30}, Window{30, 50}, -100, Window{-50, -30}, false, true},
	{Window{-50, -30}, Window{30, 50}, -35, Window{-35, -30}, true, true},
	{Window{-50, -30}, Window{30, 50}, 0, Window{0, -30}, true, false},
	{Window{-20, 10}, Window{-10, 20}, -30, Window{-20, 10}, false, true},
	{Window{-20, 10}, Window{-10, 20}, 0, Window{0, 10}, true, true},
	{Window{-20, 10}, Window{-10, 20}, 30, Window{30, 10}, true, false},
	{Window{20, 100}, Window{-100, -20}, 0, Window{20, 100}, false, true},
	{Window{20, 100}, Window{-100, -20}, 60, Window{60, 100}, true, true},
	{Window{20, 100}, Window{-100, -20}, 120, Window{120, 100}, true, false},
}

func TestNeg(t *testing.T) {
	for _, test := range windowTests {
		if got := test.w.Neg(); got != test.neg {
			t.Errorf("TestNeg(%v): got %v, want %v", test.w, got, test.neg)
		}
	}
}

func TestConstrain(t *testing.T) {
	for _, test := range windowTests {
		if gotc, gotconstrained, gotok := test.w.Constrain(test.n); gotc != test.c || gotconstrained != test.constrained || gotok != test.ok {
			t.Errorf("TestConstrain(%v, %v): got %v, %v, %v; want %v, %v, %v", test.w, test.n, gotc, gotconstrained, gotok, test.c, test.constrained, test.ok)
		}
	}
}
