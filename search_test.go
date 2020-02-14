package main

import (
	"context"
	"testing"
)

func TestCheckTerminal(t *testing.T) {
	for _, test := range []struct {
		fen string
		err error
	}{
		{"6qk/8/5B2/8/8/8/8/K6R b - - 0 1", errCheckmate},
		{"R6k/6pp/8/8/8/8/8/K7 b - - 0 1", errCheckmate},
		{"6rk/5Npp/8/8/8/8/8/7K b - - 0 1", errCheckmate},
		{"kQK5/8/8/8/8/8/8/8 b - - 100 100", errCheckmate}, // checkmate supercedes fifty-move rule
		{"7k/6q1/6pK/7b/7b/8/8/Q7 w - - 0 1", nil},         // forced mate in 1
		{"3k4/3P4/3K4/8/8/8/8/8 b - - 0 1", errStalemate},
		{"K6k/2q5/8/8/8/8/8/8 w - - 0 1", errStalemate},
		{"3k4/7b/1p3p2/8/1p1K1p2/8/b7/8 w - - 0 1", errStalemate},
		{"3k4/7b/1p3p2/8/1p1K1p2/8/b7/8 w - - 100 100", errStalemate}, // stalemate supercedes fifty-move rule
		{"K1k5/q7/3Q4/8/8/8/8/8 w - - 0 1", nil},                      // forced stalemate in 1
		{"6k1/5ppp/8/8/8/8/5PPP/6K1 w - - 99 99", nil},
		{"6k1/5ppp/8/8/8/8/5PPP/6K1 w - - 100 100", errFiftyMove},
		{"6k1/5ppp/8/8/8/8/5PPP/6K1 w - - 101 101", errFiftyMove},
		{InitialPositionFEN, nil},
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := checkTerminal(pos); got != test.err {
			t.Errorf("TestCheckTerminal(%v): got %v, want %v", test.fen, got, test.err)
		}
	}
}

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

func TestIs(t *testing.T) {
	for _, test := range []struct {
		fen     string
		l, c, t bool
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", true, false, false},  // initial position
		{"r1b2rkB/1pp1ppbp/2n3p1/8/PpP3nP/8/3Kpp2/1N3BNR w - - 0 13", true, false, false}, // legal, not check

		{"8/8/8/8/8/8/q7/K6k w - - 0 1", true, true, false},                 // check with forced capture
		{"8/8/8/8/8/q7/8/KB5k w - - 0 1", true, true, false},                // check with forced block
		{"8/8/8/8/8/q7/8/K6k w - - 0 1", true, true, false},                 // check with forced evasion
		{"8/8/8/8/8/8/2b5/k2R3K b - - 0 1", true, true, false},              // check with capture, block, or evasion
		{"7k/8/5B2/8/8/8/8/K6R b - - 0 1", true, true, false},               // double check
		{"6qk/8/5B2/8/8/8/8/K6R b - - 0 1", true, true, true},               // checkmate by double check, either can be blocked
		{"R6k/6pp/8/8/8/8/8/K7 b - - 0 1", true, true, true},                // checkmate, back rank
		{"6rk/5Npp/8/8/8/8/8/7K b - - 0 1", true, true, true},               // checkmate, smothered
		{"8/8/8/8/8/8/q7/K6k b - - 0 1", false, false, false},               // illegal, side to move not in check but opponent in check
		{"K7/2n5/8/8/8/8/5N2/7k b - - 0 1", false, true, false},             // illegal, both sides in check
		{"4R2k/6pp/8/8/8/8/8/q6K b - - 0 1", false, true, true},             // illegal, last move delivered checkmate by putting own king in check
		{"k6R/pp6/8/8/8/8/6PP/r6K w - - 0 1", false, true, true},            // illegal, both sides in checkmate
		{"8/8/7p/6pP/5pP1/3kpPp1/4P1PN/6NK w - - 0 1", false, false, false}, // illegal, only pseudo-legal move is to capture the king
		{"K6k/2q5/8/8/8/8/8/8 w - - 0 1", true, false, true},                // stalemate
		{"K6k/2q5/8/8/8/8/8/8 b - - 0 1", true, false, false},               // would be stalemate if it were the opponent's turn
		{"8/8/6K1/8/8/8/2k5/8 b - - 0 1", true, false, false},               // legal, lone kings
		{"8/8/8/3k4/3K4/8/8/8 w - - 0 1", false, true, false},               // illegal, lone kings
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := IsLegal(pos); got != test.l {
			t.Errorf("TestIsLegal(%v): got %v, want %v", test.fen, got, test.l)
		}
		if got := IsCheck(pos); got != test.c {
			t.Errorf("TestIsCheck(%v): got %v, want %v", test.fen, got, test.c)
		}
		if got := IsMate(pos); got != test.t {
			t.Errorf("TestIsMate(%v): got %v, want %v", test.fen, got, test.t)
		}
	}
}

// isDraw returns whether s represents a drawing terminal condition.
func isDraw(s Abs) bool {
	if s.err == nil {
		return false
	}
	_, ok := s.err.(checkmateError)
	return !ok
}

func TestSortFor(t *testing.T) {
	// rs is an unsorted Results; want contains the same elements
	// sorted in order of desirability for White and Black, respectively,
	// so that want[c][i] should sort before want[c][j] precisely when i < j,
	// except that the sort order of two drawing errors is undefined.
	var (
		rs = Results{
			Result{depth: 5, score: Abs{err: checkmateError(5)}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: c8, CapturePiece: Bishop, CaptureSquare: c8, PromotePiece: Rook}},
			Result{depth: 0, score: Abs{err: checkmateError(0)}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 9, score: Abs{err: checkmateError(9)}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: f6, To: f7}},
			Result{depth: 1, score: Abs{err: checkmateError(1)}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 4, score: Abs{n: 1}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 3, score: Abs{n: 150}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 3, score: Abs{n: 200}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Queen}},
			Result{depth: 4, score: Abs{n: -1}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 5, score: Abs{err: errInsufficient}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 2, score: Abs{err: errStalemate}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 8, score: Abs{err: checkmateError(8)}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Knight}},
			Result{depth: 4, score: Abs{n: -100}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 8, score: Abs{err: errFiftyMove}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 4, score: Abs{n: 100}, move: Move{Piece: Rook, From: a1, To: a8}},
			Result{depth: 4, score: Abs{err: checkmateError(4)}, move: Move{Piece: Rook, From: a1, To: a8}},
		}
		want = []Results{
			{
				Result{depth: 1, score: Abs{err: checkmateError(1)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 5, score: Abs{err: checkmateError(5)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 9, score: Abs{err: checkmateError(9)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{n: 100}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: f6, To: f7}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: c8, CapturePiece: Bishop, CaptureSquare: c8, PromotePiece: Rook}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Knight}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Queen}},
				Result{depth: 4, score: Abs{n: 1}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 8, score: Abs{err: errFiftyMove}, move: Move{Piece: Rook, From: a1, To: a8}},    //
				Result{depth: 2, score: Abs{err: errStalemate}, move: Move{Piece: Rook, From: a1, To: a8}},    // relative order undefined
				Result{depth: 5, score: Abs{err: errInsufficient}, move: Move{Piece: Rook, From: a1, To: a8}}, //
				Result{depth: 4, score: Abs{n: -1}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{n: -100}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 3, score: Abs{n: 200}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 3, score: Abs{n: 150}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 8, score: Abs{err: checkmateError(8)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{err: checkmateError(4)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 0, score: Abs{err: checkmateError(0)}, move: Move{Piece: Rook, From: a1, To: a8}},
			},
			{
				Result{depth: 1, score: Abs{err: checkmateError(1)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 5, score: Abs{err: checkmateError(5)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 9, score: Abs{err: checkmateError(9)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{n: -100}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{n: -1}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 2, score: Abs{err: errStalemate}, move: Move{Piece: Rook, From: a1, To: a8}},    //
				Result{depth: 5, score: Abs{err: errInsufficient}, move: Move{Piece: Rook, From: a1, To: a8}}, // relative order undefined
				Result{depth: 8, score: Abs{err: errFiftyMove}, move: Move{Piece: Rook, From: a1, To: a8}},    //
				Result{depth: 4, score: Abs{n: 1}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: f6, To: f7}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: c8, CapturePiece: Bishop, CaptureSquare: c8, PromotePiece: Rook}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Knight}},
				Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Queen}},
				Result{depth: 4, score: Abs{n: 100}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 3, score: Abs{n: 150}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 3, score: Abs{n: 200}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 8, score: Abs{err: checkmateError(8)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 4, score: Abs{err: checkmateError(4)}, move: Move{Piece: Rook, From: a1, To: a8}},
				Result{depth: 0, score: Abs{err: checkmateError(0)}, move: Move{Piece: Rook, From: a1, To: a8}},
			},
		}
	)
	for _, c := range []Color{White, Black} {
		got := append(Results{}, rs...)
		got.SortFor(c)
		for i := range got {
			sg, sw := got[i].score, want[c][i].score
			if sg != sw && !(isDraw(sg) && isDraw(sw)) {
				t.Errorf("TestSortFor(%v): got %v, want %v", c, got, want[c])
				break
			}
		}
	}
}

func materel(n int) Rel { return Rel{err: checkmateError(n)} }

var relTests = []struct{ r, p, n Rel }{
	{Rel{n: -100}, Rel{n: 100}, Rel{n: 100}},
	{Rel{n: 0}, Rel{n: 0}, Rel{n: 0}},
	{Rel{n: 100}, Rel{n: -100}, Rel{n: -100}},
	{Rel{err: errStalemate}, Rel{err: errStalemate}, Rel{err: errStalemate}},
	{Rel{err: errInsufficient}, Rel{err: errInsufficient}, Rel{err: errInsufficient}},
	{Rel{err: errFiftyMove}, Rel{err: errFiftyMove}, Rel{err: errFiftyMove}},
	{materel(4), materel(5), materel(3)},
	{materel(3), materel(4), materel(2)},
	{materel(1), materel(2), Rel{err: errCheckmate}},
}

func TestRelPrev(t *testing.T) {
	for _, test := range relTests {
		if got := test.r.Prev(); got != test.p {
			t.Errorf("TestPrev(%v): got %v, want %v", test.r, got, test.p)
		}
	}
}

func TestRelNext(t *testing.T) {
	for _, test := range relTests {
		if got := test.r.Next(); got != test.n {
			t.Errorf("TestNext(%v): got %v, want %v", test.r, got, test.n)
		}
	}
}

func centwindow(a, b int) Window { return Window{Rel{n: a}, Rel{n: b}} }

func TestWindowNext(t *testing.T) {
	for _, test := range []struct {
		w, n Window
	}{
		{centwindow(-50, -30), centwindow(30, 50)},
		{centwindow(-20, 10), centwindow(-10, 20)},
		{centwindow(20, 100), centwindow(-100, -20)},
		{Window{materel(4), Rel{n: 100}}, Window{Rel{n: -100}, materel(3)}},
		{Window{Rel{n: 1}, materel(5)}, Window{materel(4), Rel{n: -1}}},
		{Window{materel(2), materel(9)}, Window{materel(8), materel(1)}},
	} {
		if got := test.w.Next(); got != test.n {
			t.Errorf("TestNext(%v): got %v, want %v", test.w, got, test.n)
		}
	}
}

func TestConstrain(t *testing.T) {
	for _, test := range []struct {
		w  Window
		n  Rel
		c  Window
		ok bool
	}{
		{centwindow(-50, -30), materel(6), centwindow(-50, -30), true},
		{centwindow(-50, -30), Rel{n: -100}, centwindow(-50, -30), true},
		{centwindow(-50, -30), Rel{n: -35}, centwindow(-35, -30), true},
		{centwindow(-50, -30), Rel{n: -30}, centwindow(-30, -30), true},
		{centwindow(-50, -30), Rel{n: 0}, centwindow(-30, -30), false},
		{centwindow(-50, -30), Rel{err: errStalemate}, centwindow(-30, -30), false},
		{centwindow(-50, -30), materel(5), centwindow(-30, -30), false},

		{centwindow(-20, 10), materel(6), centwindow(-20, 10), true},
		{centwindow(-20, 10), Rel{n: -30}, centwindow(-20, 10), true},
		{centwindow(-20, 10), Rel{n: 0}, centwindow(0, 10), true},
		{centwindow(-20, 10), Rel{err: errStalemate}, Window{Rel{err: errStalemate}, Rel{n: 10}}, true},
		{centwindow(-20, 10), Rel{n: 10}, centwindow(10, 10), true},
		{centwindow(-20, 10), Rel{n: 30}, centwindow(10, 10), false},
		{centwindow(-20, 10), materel(5), centwindow(10, 10), false},

		{centwindow(20, 100), materel(6), centwindow(20, 100), true},
		{centwindow(20, 100), Rel{err: errStalemate}, centwindow(20, 100), true},
		{centwindow(20, 100), Rel{n: 0}, centwindow(20, 100), true},
		{centwindow(20, 100), Rel{n: 60}, centwindow(60, 100), true},
		{centwindow(20, 100), Rel{n: 100}, centwindow(100, 100), true},
		{centwindow(20, 100), Rel{n: 120}, centwindow(100, 100), false},
		{centwindow(20, 100), materel(5), centwindow(100, 100), false},

		{Window{materel(4), Rel{n: 5}}, materel(2), Window{materel(4), Rel{n: 5}}, true},
		{Window{materel(4), Rel{n: 5}}, materel(6), Window{materel(6), Rel{n: 5}}, true},
		{Window{materel(4), Rel{n: 5}}, Rel{n: -800}, centwindow(-800, 5), true},
		{Window{materel(4), Rel{n: 5}}, Rel{err: errStalemate}, Window{Rel{err: errStalemate}, Rel{n: 5}}, true},
		{Window{materel(4), Rel{n: 5}}, Rel{n: 5}, centwindow(5, 5), true},
		{Window{materel(4), Rel{n: 5}}, Rel{n: 800}, centwindow(5, 5), false},
		{Window{materel(4), Rel{n: 5}}, materel(5), centwindow(5, 5), false},

		{Window{Rel{n: 5}, materel(5)}, materel(8), Window{Rel{n: 5}, materel(5)}, true},
		{Window{Rel{n: 5}, materel(5)}, Rel{n: -800}, Window{Rel{n: 5}, materel(5)}, true},
		{Window{Rel{n: 5}, materel(5)}, Rel{err: errStalemate}, Window{Rel{n: 5}, materel(5)}, true},
		{Window{Rel{n: 5}, materel(5)}, Rel{n: 15}, Window{Rel{n: 15}, materel(5)}, true},
		{Window{Rel{n: 5}, materel(5)}, materel(5), Window{materel(5), materel(5)}, true},
		{Window{Rel{n: 5}, materel(5)}, materel(1), Window{materel(5), materel(5)}, false},

		{Window{materel(8), materel(5)}, materel(2), Window{materel(8), materel(5)}, true},
		{Window{materel(8), materel(5)}, materel(24), Window{materel(24), materel(5)}, true},
		{Window{materel(8), materel(5)}, Rel{n: 50}, Window{Rel{n: 50}, materel(5)}, true},
		{Window{materel(8), materel(5)}, materel(25), Window{materel(25), materel(5)}, true},
		{Window{materel(8), materel(5)}, materel(5), Window{materel(5), materel(5)}, true},
		{Window{materel(8), materel(5)}, materel(3), Window{materel(5), materel(5)}, false},
	} {
		if gotc, gotok := test.w.Constrain(test.n); gotc != test.c || gotok != test.ok {
			t.Errorf("TestConstrain(%v, %v): got %v, %v; want %v, %v", test.w, test.n, gotc, gotok, test.c, test.ok)
		}
	}
}

func BenchmarkIsMate(b *testing.B) {
	for _, benchmark := range []struct{ name, fen string }{
		{"initial", InitialPositionFEN},
		{"check", "4k3/4q3/8/8/8/8/8/4K3 w - - 0 1"},
		{"forced", "k7/8/8/8/8/8/8/K1q5 w - - 0 1"},
		{"checkmate", "8/8/8/8/8/8/8/Kqk5 w - - 0 1"},
		{"stalemate", "8/8/8/8/8/8/2q5/K1k5 w - - 0 1"},
	} {
		b.Run(benchmark.name, func(b *testing.B) {
			pos, err := ParseFEN(benchmark.fen)
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < b.N; i++ {
				IsMate(pos)
			}
		})
	}
}

func BenchmarkSearchPosition(b *testing.B) {
	pos, err := ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		SearchPosition(ctx, pos, 2)
	}
}

func BenchmarkSortFor(b *testing.B) {
	var unsorted = Results{
		Result{depth: 5, score: Abs{err: checkmateError(5)}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: c8, CapturePiece: Bishop, CaptureSquare: c8, PromotePiece: Rook}},
		Result{depth: 0, score: Abs{err: checkmateError(0)}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 9, score: Abs{err: checkmateError(9)}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: f6, To: f7}},
		Result{depth: 1, score: Abs{err: checkmateError(1)}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 4, score: Abs{n: 1}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 3, score: Abs{n: 150}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 3, score: Abs{n: 200}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Queen}},
		Result{depth: 4, score: Abs{n: -1}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 5, score: Abs{err: errInsufficient}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 2, score: Abs{err: errStalemate}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 8, score: Abs{err: checkmateError(8)}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 4, score: Abs{n: 50}, move: Move{Piece: Pawn, From: d7, To: d8, PromotePiece: Knight}},
		Result{depth: 4, score: Abs{n: -100}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 8, score: Abs{err: errFiftyMove}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 4, score: Abs{n: 100}, move: Move{Piece: Rook, From: a1, To: a8}},
		Result{depth: 4, score: Abs{err: checkmateError(4)}, move: Move{Piece: Rook, From: a1, To: a8}},
	}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rs := append(Results{}, unsorted...)
		b.StartTimer()
		rs.SortFor(White)
	}
}
