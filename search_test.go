package main

import "testing"

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

func BenchmarkMake(b *testing.B) {
	pos, err := ParseFEN("r4k2/1P6/8/3Pp3/8/8/6P1/R3K2R w KQ e6 0 1")
	if err != nil {
		b.Fatal(err)
	}
	for _, benchmark := range []struct {
		name string
		move Move
	}{
		{"quiet", Move{From: e1, To: f2, Piece: King}},
		{"double", Move{From: g2, To: g4, Piece: Pawn}},
		{"en passant", Move{From: d5, To: e6, Piece: Pawn, CapturePiece: Pawn, CaptureSquare: e5}},
		{"capture", Move{From: a1, To: a8, Piece: Rook, CapturePiece: Rook, CaptureSquare: a8}},
		{"promotion", Move{From: b7, To: b8, Piece: Pawn, PromotePiece: Queen}},
		{"capture promotion", Move{From: b7, To: a8, Piece: Pawn, CapturePiece: Rook, CaptureSquare: a8, PromotePiece: Queen}},
		{"castle queenside", Move{From: e1, To: c1, Piece: King}},
		{"castle kingside", Move{From: e1, To: g1, Piece: King}},
	} {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Make(pos, benchmark.move)
			}
		})
	}
}
