package main

import "testing"

func TestLess(t *testing.T) {
	// ss is a slice of Rels sorted such that ss[i] is lower than ss[j] precisely when i < j,
	// except that no drawing evaluation is lower than another.
	var ss = []Rel{
		Rel{err: checkmateError(0)},
		Rel{err: checkmateError(4)},
		Rel{err: checkmateError(14)},
		Rel{n: -100},
		Rel{n: -1},
		Rel{err: errFiftyMove},
		Rel{err: errInsufficient},
		Rel{err: errStalemate},
		Rel{n: 1},
		Rel{n: 100},
		Rel{err: checkmateError(15)},
		Rel{err: checkmateError(5)},
		Rel{err: checkmateError(1)},
	}
	for i, a := range ss {
		for j, b := range ss {
			_, ach := a.err.(checkmateError)
			_, bch := b.err.(checkmateError)
			adraw := a.err != nil && !ach
			bdraw := b.err != nil && !bch
			want := i < j && !(adraw && bdraw)
			if got := Less(Score(a), Score(b)); got != want {
				t.Errorf("TestLess(%v, %v): got %v, want %v", a, b, got, want)
			}
		}
	}
}

func TestIsInsufficient(t *testing.T) {
	for _, test := range []struct {
		fen  string
		want bool
	}{
		{"K6k/8/8/8/8/8/8/8 w - - 0 1", true},
		{"KN5k/8/8/8/8/8/8/8 w - - 0 1", true},
		{"KB5k/8/8/8/8/8/8/8 w - - 0 1", true},
		{"KB1b3k/8/8/8/8/8/8/8 w - - 0 1", true},
		{"K6k/8/8/8/8/8/8/1B1B1B1B w - - 0 1", true},
		{"K1b1b1bk/8/8/8/8/8/8/1B1B1B1B w - - 0 1", true},
		{"KN4nk/8/8/8/8/8/8/8 w - - 0 1", false},
		{"KNN4k/8/8/8/8/8/8/8 w - - 0 1", false},
		{"KNB4k/8/8/8/8/8/8/8 w - - 0 1", false},
		{"KNb4k/8/8/8/8/8/8/8 w - - 0 1", false},
		{"K6k/P7/8/8/8/8/8/8 w - - 0 1", false},
		{"K6k/R7/8/8/8/8/8/8 w - - 0 1", false},
		{"K6k/Q7/8/8/8/8/8/8 w - - 0 1", false},
		{InitialPositionFEN, false},
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := IsInsufficient(pos); got != test.want {
			t.Errorf("IsInsufficient(%v): got %v, want %v", test.fen, got, test.want)
		}
	}
}

func BenchmarkEval(b *testing.B) {
	pos, err := ParseFEN("r1b2rkB/1pp1ppbp/2n3p1/8/PpP3nP/8/3KPP2/1N3BNR w - - 0 13")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		Eval(pos)
	}
}
