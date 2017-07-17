package main

import "testing"

const aroundD4 = (CFile|DFile|EFile)&(Rank3|Rank4|Rank5) ^ (DFile & Rank4)

func TestSouthAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), 0},
		{h1.Board(), ^h1.Board(), 0},
		{a8.Board(), ^a8.Board(), AFile ^ a8.Board()},
		{h8.Board(), ^h8.Board(), HFile ^ h8.Board()},
		{d4.Board(), 0, d3.Board()},
		{d4.Board(), aroundD4, d3.Board() | d2.Board()},
	} {
		if got := attackFill(v.input, v.empty, south); got != v.want {
			t.Errorf("attackFill(%016x, %016x, south): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestWestAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), 0},
		{h1.Board(), ^h1.Board(), Rank1 ^ h1.Board()},
		{a8.Board(), ^a8.Board(), 0},
		{h8.Board(), ^h8.Board(), Rank8 ^ h8.Board()},
		{d4.Board(), 0, c4.Board()},
		{d4.Board(), aroundD4, c4.Board() | b4.Board()},
	} {
		if got := attackFill(v.input, v.empty, west); got != v.want {
			t.Errorf("attackFill(%016x, %016x, west): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestEastAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), Rank1 ^ a1.Board()},
		{h1.Board(), ^h1.Board(), 0},
		{a8.Board(), ^a8.Board(), Rank8 ^ a8.Board()},
		{h8.Board(), ^h8.Board(), 0},
		{d4.Board(), 0, e4.Board()},
		{d4.Board(), aroundD4, e4.Board() | f4.Board()},
	} {
		if got := attackFill(v.input, v.empty, east); got != v.want {
			t.Errorf("attackFill(%016x, %016x, east): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestNorthAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), AFile ^ a1.Board()},
		{h1.Board(), ^h1.Board(), HFile ^ h1.Board()},
		{a8.Board(), ^a8.Board(), 0},
		{h8.Board(), ^h8.Board(), 0},
		{d4.Board(), 0, d5.Board()},
		{d4.Board(), aroundD4, d5.Board() | d6.Board()},
	} {
		if got := attackFill(v.input, v.empty, north); got != v.want {
			t.Errorf("attackFill(%016x, %016x, north): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestSouthwestAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), 0},
		{h1.Board(), ^h1.Board(), 0},
		{a8.Board(), ^a8.Board(), 0},
		{h8.Board(), ^h8.Board(), LongDiagonal ^ h8.Board()},
		{d4.Board(), 0, c3.Board()},
		{d4.Board(), aroundD4, c3.Board() | b2.Board()},
	} {
		if got := attackFill(v.input, v.empty, southwest); got != v.want {
			t.Errorf("attackFill(%016x, %016x, southwest): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestSoutheastAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), 0},
		{h1.Board(), ^h1.Board(), 0},
		{a8.Board(), ^a8.Board(), LongAntiDiagonal ^ a8.Board()},
		{h8.Board(), ^h8.Board(), 0},
		{d4.Board(), 0, e3.Board()},
		{d4.Board(), aroundD4, e3.Board() | f2.Board()},
	} {
		if got := attackFill(v.input, v.empty, southeast); got != v.want {
			t.Errorf("attackFill(%016x, %016x, southeast): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestNorthwestAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), 0},
		{h1.Board(), ^h1.Board(), LongAntiDiagonal ^ h1.Board()},
		{a8.Board(), ^a8.Board(), 0},
		{h8.Board(), ^h8.Board(), 0},
		{d4.Board(), 0, c5.Board()},
		{d4.Board(), aroundD4, c5.Board() | b6.Board()},
	} {
		if got := attackFill(v.input, v.empty, northwest); got != v.want {
			t.Errorf("attackFill(%016x, %016x, northwest): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestNortheastAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), LongDiagonal ^ a1.Board()},
		{h1.Board(), ^h1.Board(), 0},
		{a8.Board(), ^a8.Board(), 0},
		{h8.Board(), ^h8.Board(), 0},
		{d4.Board(), 0, e5.Board()},
		{d4.Board(), aroundD4, e5.Board() | f6.Board()},
	} {
		if got := attackFill(v.input, v.empty, northeast); got != v.want {
			t.Errorf("attackFill(%016x, %016x, northeast): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestWhitePawnAdvances(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a2.Board(), ^a2.Board(), a3.Board() | a4.Board()},
		{a2.Board(), ^(a2.Board() | a4.Board()), a3.Board()},
		{a2.Board(), ^(a2.Board() | a3.Board()), 0},
		{a5.Board(), ^a5.Board(), a6.Board()},
		{a5.Board(), ^(a5.Board() | a6.Board()), 0},
		{a7.Board(), ^a7.Board(), a8.Board()},
		{a7.Board(), ^(a7.Board() | a8.Board()), 0},

		{e2.Board(), ^e2.Board(), e3.Board() | e4.Board()},
		{e2.Board(), ^(e2.Board() | e4.Board()), e3.Board()},
		{e2.Board(), ^(e2.Board() | e3.Board()), 0},
		{e5.Board(), ^e5.Board(), e6.Board()},
		{e5.Board(), ^(e5.Board() | e6.Board()), 0},
		{e7.Board(), ^e7.Board(), e8.Board()},
		{e7.Board(), ^(e7.Board() | e8.Board()), 0},

		{h2.Board(), ^h2.Board(), h3.Board() | h4.Board()},
		{h2.Board(), ^(h2.Board() | h4.Board()), h3.Board()},
		{h2.Board(), ^(h2.Board() | h3.Board()), 0},
		{h5.Board(), ^h5.Board(), h6.Board()},
		{h5.Board(), ^(h5.Board() | h6.Board()), 0},
		{h7.Board(), ^h7.Board(), h8.Board()},
		{h7.Board(), ^(h7.Board() | h8.Board()), 0},
	} {
		if got := whitePawnAdvances(v.input, v.empty); got != v.want {
			t.Errorf("whitePawnAdvances(%016x, %016x): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestBlackPawnAdvances(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a7.Board(), ^a7.Board(), a6.Board() | a5.Board()},
		{a7.Board(), ^(a7.Board() | a5.Board()), a6.Board()},
		{a7.Board(), ^(a7.Board() | a6.Board()), 0},
		{a4.Board(), ^a4.Board(), a3.Board()},
		{a4.Board(), ^(a4.Board() | a3.Board()), 0},
		{a2.Board(), ^a2.Board(), a1.Board()},
		{a2.Board(), ^(a2.Board() | a1.Board()), 0},

		{e7.Board(), ^e7.Board(), e6.Board() | e5.Board()},
		{e7.Board(), ^(e7.Board() | e5.Board()), e6.Board()},
		{e7.Board(), ^(e7.Board() | e6.Board()), 0},
		{e4.Board(), ^e4.Board(), e3.Board()},
		{e4.Board(), ^(e4.Board() | e3.Board()), 0},
		{e2.Board(), ^e2.Board(), e1.Board()},
		{e2.Board(), ^(e2.Board() | e1.Board()), 0},

		{h7.Board(), ^h7.Board(), h6.Board() | h5.Board()},
		{h7.Board(), ^(h7.Board() | h5.Board()), h6.Board()},
		{h7.Board(), ^(h7.Board() | h6.Board()), 0},
		{h4.Board(), ^h4.Board(), h3.Board()},
		{h4.Board(), ^(h4.Board() | h3.Board()), 0},
		{h2.Board(), ^h2.Board(), h1.Board()},
		{h2.Board(), ^(h2.Board() | h1.Board()), 0},
	} {
		if got := blackPawnAdvances(v.input, v.empty); got != v.want {
			t.Errorf("blackPawnAdvances(%016x, %016x): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestWhitePawnAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a2.Board(), 0, b3.Board()},
		{a2.Board(), BFile, 0},
		{a5.Board(), 0, b6.Board()},
		{a5.Board(), BFile, 0},
		{a7.Board(), 0, b8.Board()},
		{a7.Board(), BFile, 0},

		{e2.Board(), 0, d3.Board() | f3.Board()},
		{e2.Board(), DFile, f3.Board()},
		{e2.Board(), FFile, d3.Board()},
		{e2.Board(), DFile | FFile, 0},
		{e5.Board(), 0, d6.Board() | f6.Board()},
		{e5.Board(), DFile, f6.Board()},
		{e5.Board(), FFile, d6.Board()},
		{e5.Board(), DFile | FFile, 0},
		{e7.Board(), 0, d8.Board() | f8.Board()},
		{e7.Board(), DFile, f8.Board()},
		{e7.Board(), FFile, d8.Board()},
		{e7.Board(), DFile | FFile, 0},

		{h2.Board(), 0, g3.Board()},
		{h2.Board(), GFile, 0},
		{h5.Board(), 0, g6.Board()},
		{h5.Board(), GFile, 0},
		{h7.Board(), 0, g8.Board()},
		{h7.Board(), GFile, 0},
	} {
		if got := whitePawnAttacks(v.input, v.empty); got != v.want {
			t.Errorf("whitePawnAttacks(%016x, %016x): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestBlackPawnAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a7.Board(), 0, b6.Board()},
		{a7.Board(), BFile, 0},
		{a4.Board(), 0, b3.Board()},
		{a4.Board(), BFile, 0},
		{a2.Board(), 0, b1.Board()},
		{a2.Board(), BFile, 0},

		{e7.Board(), 0, d6.Board() | f6.Board()},
		{e7.Board(), DFile, f6.Board()},
		{e7.Board(), FFile, d6.Board()},
		{e7.Board(), DFile | FFile, 0},
		{e4.Board(), 0, d3.Board() | f3.Board()},
		{e4.Board(), DFile, f3.Board()},
		{e4.Board(), FFile, d3.Board()},
		{e4.Board(), DFile | FFile, 0},
		{e2.Board(), 0, d1.Board() | f1.Board()},
		{e2.Board(), DFile, f1.Board()},
		{e2.Board(), FFile, d1.Board()},
		{e2.Board(), DFile | FFile, 0},

		{h7.Board(), 0, g6.Board()},
		{h7.Board(), GFile, 0},
		{h4.Board(), 0, g3.Board()},
		{h4.Board(), GFile, 0},
		{h2.Board(), 0, g1.Board()},
		{h2.Board(), GFile, 0},
	} {
		if got := blackPawnAttacks(v.input, v.empty); got != v.want {
			t.Errorf("blackPawnAttacks(%016x, %016x): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestBishopAttacks(t *testing.T) {
	var ld, la = LongDiagonal, LongAntiDiagonal
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), LongDiagonal ^ a1.Board()},
		{e7.Board(), ^e7.Board(), (ld<<16 | la<<24) ^ e7.Board()},
		{d4.Board(), 0, c3.Board() | e3.Board() | c5.Board() | e5.Board()},
		{d4.Board(), aroundD4, b2.Board() | f2.Board() | c3.Board() | e3.Board() | c5.Board() | e5.Board() | b6.Board() | f6.Board()},
	} {
		if got := bishopAttacks(v.input, v.empty); got != v.want {
			t.Errorf("bishopAttacks(%016x, %016x): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestRookAttacks(t *testing.T) {
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), (AFile | Rank1) ^ a1.Board()},
		{e7.Board(), ^e7.Board(), (EFile | Rank7) ^ e7.Board()},
		{d4.Board(), 0, d3.Board() | c4.Board() | e4.Board() | d5.Board()},
		{d4.Board(), aroundD4, d2.Board() | d3.Board() | b4.Board() | c4.Board() | e4.Board() | f4.Board() | d5.Board() | d6.Board()},
	} {
		if got := rookAttacks(v.input, v.empty); got != v.want {
			t.Errorf("rookAttacks(%016x, %016x): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestQueenAttacks(t *testing.T) {
	var ld, la = LongDiagonal, LongAntiDiagonal
	for _, v := range []struct {
		input Board
		empty Board
		want  Board
	}{
		{a1.Board(), ^a1.Board(), (AFile | Rank1 | LongDiagonal) ^ a1.Board()},
		{e7.Board(), ^e7.Board(), (EFile | Rank7 | ld<<16 | la<<24) ^ e7.Board()},
		{d4.Board(), 0, aroundD4},
		{d4.Board(), aroundD4, aroundD4 | b2.Board() | d2.Board() | f2.Board() | b4.Board() | f4.Board() | b6.Board() | d6.Board() | f6.Board()},
	} {
		if got := queenAttacks(v.input, v.empty); got != v.want {
			t.Errorf("queenAttacks(%016x, %016x): got %016x, want %016x", v.input, v.empty, got, v.want)
		}
	}
}

func TestIsAttacked(t *testing.T) {
	for _, v := range []struct {
		pos  Position
		s    Square
		c    Color
		want bool
	}{
		{InitialPosition, a1, White, false},
		{InitialPosition, b1, White, true},
		{InitialPosition, b1, Black, false},
		{InitialPosition, f3, White, true},
		{InitialPosition, f4, White, false},
		{InitialPosition, f6, White, false},
		{InitialPosition, f6, Black, true},
		{InitialPosition, f8, Black, true},
		{InitialPosition, f8, White, false},
		{InitialPosition, h8, Black, false},
	} {
		if got := IsAttacked(v.pos, v.s, v.c); got != v.want {
			t.Errorf("IsAttacked(%v, %v, %v): got %v, want %v", v.pos, v.s, v.c, got, v.want)
		}
	}
}
