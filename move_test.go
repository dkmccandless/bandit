package main

import "testing"

const aroundD4 = (CFile|DFile|EFile)&(Rank3|Rank4|Rank5) ^ (DFile & Rank4)

type boardsTest struct{ input, empty, want Board }

func TestAttackFill(t *testing.T) {
	for _, dir := range []struct {
		f     func(Board) Board
		name  string
		tests []boardsTest
	}{
		{
			south, "south", []boardsTest{
				{a1.Board(), ^a1.Board(), 0},
				{h1.Board(), ^h1.Board(), 0},
				{a8.Board(), ^a8.Board(), AFile ^ a8.Board()},
				{h8.Board(), ^h8.Board(), HFile ^ h8.Board()},
				{d4.Board(), 0, d3.Board()},
				{d4.Board(), aroundD4, d3.Board() | d2.Board()},
			},
		},
		{
			west, "west", []boardsTest{
				{a1.Board(), ^a1.Board(), 0},
				{h1.Board(), ^h1.Board(), Rank1 ^ h1.Board()},
				{a8.Board(), ^a8.Board(), 0},
				{h8.Board(), ^h8.Board(), Rank8 ^ h8.Board()},
				{d4.Board(), 0, c4.Board()},
				{d4.Board(), aroundD4, c4.Board() | b4.Board()},
			},
		},
		{
			east, "east", []boardsTest{
				{a1.Board(), ^a1.Board(), Rank1 ^ a1.Board()},
				{h1.Board(), ^h1.Board(), 0},
				{a8.Board(), ^a8.Board(), Rank8 ^ a8.Board()},
				{h8.Board(), ^h8.Board(), 0},
				{d4.Board(), 0, e4.Board()},
				{d4.Board(), aroundD4, e4.Board() | f4.Board()},
			},
		},
		{
			north, "north", []boardsTest{
				{a1.Board(), ^a1.Board(), AFile ^ a1.Board()},
				{h1.Board(), ^h1.Board(), HFile ^ h1.Board()},
				{a8.Board(), ^a8.Board(), 0},
				{h8.Board(), ^h8.Board(), 0},
				{d4.Board(), 0, d5.Board()},
				{d4.Board(), aroundD4, d5.Board() | d6.Board()},
			},
		},
		{
			southwest, "southwest", []boardsTest{
				{a1.Board(), ^a1.Board(), 0},
				{h1.Board(), ^h1.Board(), 0},
				{a8.Board(), ^a8.Board(), 0},
				{h8.Board(), ^h8.Board(), LongDiagonal ^ h8.Board()},
				{d4.Board(), 0, c3.Board()},
				{d4.Board(), aroundD4, c3.Board() | b2.Board()},
			},
		},
		{
			southeast, "southeast", []boardsTest{
				{a1.Board(), ^a1.Board(), 0},
				{h1.Board(), ^h1.Board(), 0},
				{a8.Board(), ^a8.Board(), LongAntiDiagonal ^ a8.Board()},
				{h8.Board(), ^h8.Board(), 0},
				{d4.Board(), 0, e3.Board()},
				{d4.Board(), aroundD4, e3.Board() | f2.Board()},
			},
		},
		{
			northwest, "northwest", []boardsTest{
				{a1.Board(), ^a1.Board(), 0},
				{h1.Board(), ^h1.Board(), LongAntiDiagonal ^ h1.Board()},
				{a8.Board(), ^a8.Board(), 0},
				{h8.Board(), ^h8.Board(), 0},
				{d4.Board(), 0, c5.Board()},
				{d4.Board(), aroundD4, c5.Board() | b6.Board()},
			},
		},
		{
			northeast, "northeast", []boardsTest{
				{a1.Board(), ^a1.Board(), LongDiagonal ^ a1.Board()},
				{h1.Board(), ^h1.Board(), 0},
				{a8.Board(), ^a8.Board(), 0},
				{h8.Board(), ^h8.Board(), 0},
				{d4.Board(), 0, e5.Board()},
				{d4.Board(), aroundD4, e5.Board() | f6.Board()},
			},
		},
	} {
		for _, test := range dir.tests {
			if got := attackFill(test.input, test.empty, dir.f); got != test.want {
				t.Errorf("attackFill(%016x, %016x, %v): got %016x, want %016x", test.input, test.empty, dir.name, got, test.want)
			}
		}
	}
}

func TestPieceAttacks(t *testing.T) {
	var ld, la = LongDiagonal, LongAntiDiagonal
	for _, set := range []struct {
		f     func(Board, Board) Board
		name  string
		tests []boardsTest
	}{
		{
			whitePawnAdvances, "whitePawnAdvances", []boardsTest{
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
			},
		},
		{
			blackPawnAdvances, "blackPawnAdvances", []boardsTest{
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
			},
		},
		{
			whitePawnAttacks, "whitePawnAttacks", []boardsTest{
				{a2.Board(), 0, b3.Board()},
				{a5.Board(), 0, b6.Board()},
				{a7.Board(), 0, b8.Board()},
				{e2.Board(), 0, d3.Board() | f3.Board()},
				{e5.Board(), 0, d6.Board() | f6.Board()},
				{e7.Board(), 0, d8.Board() | f8.Board()},
				{h2.Board(), 0, g3.Board()},
				{h5.Board(), 0, g6.Board()},
				{h7.Board(), 0, g8.Board()},
			},
		},
		{
			blackPawnAttacks, "blackPawnAttacks", []boardsTest{
				{a7.Board(), 0, b6.Board()},
				{a4.Board(), 0, b3.Board()},
				{a2.Board(), 0, b1.Board()},
				{e7.Board(), 0, d6.Board() | f6.Board()},
				{e4.Board(), 0, d3.Board() | f3.Board()},
				{e2.Board(), 0, d1.Board() | f1.Board()},
				{h7.Board(), 0, g6.Board()},
				{h4.Board(), 0, g3.Board()},
				{h2.Board(), 0, g1.Board()},
			},
		},
		{
			bishopAttacks, "bishopAttacks", []boardsTest{
				{a1.Board(), ^a1.Board(), LongDiagonal ^ a1.Board()},
				{e7.Board(), ^e7.Board(), north(north(ld)) | north(north(north(la))) ^ e7.Board()},
				{d4.Board(), 0, c3.Board() | e3.Board() | c5.Board() | e5.Board()},
				{d4.Board(), aroundD4, b2.Board() | f2.Board() | c3.Board() | e3.Board() | c5.Board() | e5.Board() | b6.Board() | f6.Board()},
			},
		},
		{
			rookAttacks, "rookAttacks", []boardsTest{
				{a1.Board(), ^a1.Board(), (AFile | Rank1) ^ a1.Board()},
				{e7.Board(), ^e7.Board(), (EFile | Rank7) ^ e7.Board()},
				{d4.Board(), 0, d3.Board() | c4.Board() | e4.Board() | d5.Board()},
				{d4.Board(), aroundD4, d2.Board() | d3.Board() | b4.Board() | c4.Board() | e4.Board() | f4.Board() | d5.Board() | d6.Board()},
			},
		},
		{
			queenAttacks, "queenAttacks", []boardsTest{
				{a1.Board(), ^a1.Board(), (AFile | Rank1 | LongDiagonal) ^ a1.Board()},
				{e7.Board(), ^e7.Board(), EFile | Rank7 | north(north(ld)) | north(north(north(la))) ^ e7.Board()},
				{d4.Board(), 0, aroundD4},
				{d4.Board(), aroundD4, aroundD4 | b2.Board() | d2.Board() | f2.Board() | b4.Board() | f4.Board() | b6.Board() | d6.Board() | f6.Board()},
			},
		},
	} {
		for _, test := range set.tests {
			if got := set.f(test.input, test.empty); got != test.want {
				t.Errorf("%v(%016x, %016x): got %016x, want %016x", set.name, test.input, test.empty, got, test.want)
			}
		}
	}
}

func TestIsAttacked(t *testing.T) {
	for _, test := range []struct {
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
		if got := IsAttacked(test.pos, test.s, test.c); got != test.want {
			t.Errorf("IsAttacked(%v, %v, %v): got %v, want %v", test.pos, test.s, test.c, got, test.want)
		}
	}
}

func TestHasLegalMove(t *testing.T) {
	for _, test := range []struct {
		fen  string
		want bool
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", true},  // initial position
		{"r1b2rkB/1pp1ppbp/2n3p1/8/PpP3nP/8/3Kpp2/1N3BNR w - - 0 13", true}, // legal, not check

		{"8/8/8/8/8/8/q7/K6k w - - 0 1", true},     // check with forced capture
		{"8/8/8/8/8/q7/8/KB5k w - - 0 1", true},    // check with forced block
		{"8/8/8/8/8/q7/8/K6k w - - 0 1", true},     // check with forced evasion
		{"8/8/8/8/8/8/2b5/k2R3K b - - 0 1", true},  // check with capture, block, or evasion
		{"7k/8/5B2/8/8/8/8/K6R b - - 0 1", true},   // double check
		{"6qk/8/5B2/8/8/8/8/K6R b - - 0 1", false}, // checkmate by double check, either can be blocked
		{"R6k/6pp/8/8/8/8/8/K7 b - - 0 1", false},  // checkmate, back rank
		{"6rk/5Npp/8/8/8/8/8/7K b - - 0 1", false}, // checkmate, smothered
		{"K6k/2q5/8/8/8/8/8/8 w - - 0 1", false},   // stalemate
		{"K6k/2q5/8/8/8/8/8/8 b - - 0 1", true},    // would be stalemate if it were the opponent's turn
		{"8/8/6K1/8/8/8/2k5/8 b - - 0 1", true},    // legal, lone kings
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := hasLegalMove(pos); got != test.want {
			t.Errorf("hasLegalMove(%v): got %v, want %v", test.fen, got, test.want)
		}
	}

}

func TestCanCastleQS(t *testing.T) {
	for _, test := range []struct {
		fen  string
		want bool
	}{
		{"4k3/8/8/8/8/8/8/4K2R w Q - 0 1", true},
		{"4k3/8/8/8/8/8/8/4K2R w - - 0 1", false},
		{"4k3/8/8/8/8/8/8/4K2R b Q - 0 1", false},
		{"4k2r/8/8/8/8/8/8/4K2R w Q - 0 1", true},
		{"1r2k3/8/8/8/8/8/8/4K2R w Q - 0 1", true},
		{"3rk3/8/8/8/8/8/8/4K2R w Q - 0 1", false},
		{"4k3/8/8/8/8/8/8/RN2K3 w Q - 0 1", false},
		{"4k3/8/8/8/8/8/8/R3K1N1 w Q - 0 1", true},

		{"4k2r/8/8/8/8/8/8/4K3 b q - 0 1", true},
		{"4k2r/8/8/8/8/8/8/4K3 b - - 0 1", false},
		{"4k2r/8/8/8/8/8/8/4K3 w q - 0 1", false},
		{"4k2r/8/8/8/8/8/8/4K2R b q - 0 1", true},
		{"4k2r/8/8/8/8/8/8/1R2K3 b q - 0 1", true},
		{"4k2r/8/8/8/8/8/8/3RK3 b q - 0 1", false},
		{"rn2k3/8/8/8/8/8/8/4K3 b q - 0 1", false},
		{"r3k1n1/8/8/8/8/8/8/4K3 b q - 0 1", true},
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := canCastle(pos, QS); got != test.want {
			t.Errorf("canCastleQS(%v): got %v, want %v", test.fen, got, test.want)
		}
	}
}

func TestCanCastleKS(t *testing.T) {
	for _, test := range []struct {
		fen  string
		want bool
	}{
		{"4k3/8/8/8/8/8/8/4K2R w K - 0 1", true},
		{"4k3/8/8/8/8/8/8/4K2R w - - 0 1", false},
		{"4k3/8/8/8/8/8/8/4K2R b K - 0 1", false},
		{"4k2r/8/8/8/8/8/8/4K2R w K - 0 1", true},
		{"4k1r1/8/8/8/8/8/8/4K2R w K - 0 1", false},
		{"4kr2/8/8/8/8/8/8/4K2R w K - 0 1", false},
		{"4k3/8/8/8/8/8/8/4K1NR w K - 0 1", false},
		{"4k3/8/8/8/8/8/8/RN2K2R w K - 0 1", true},

		{"4k2r/8/8/8/8/8/8/4K3 b k - 0 1", true},
		{"4k2r/8/8/8/8/8/8/4K3 b - - 0 1", false},
		{"4k2r/8/8/8/8/8/8/4K3 w k - 0 1", false},
		{"4k2r/8/8/8/8/8/8/4K2R b k - 0 1", true},
		{"4k2r/8/8/8/8/8/8/4K1R1 b k - 0 1", false},
		{"4k2r/8/8/8/8/8/8/4KR2 b k - 0 1", false},
		{"4k1nr/8/8/8/8/8/8/4K3 b k - 0 1", false},
		{"rn2k2r/8/8/8/8/8/8/4K3 b k - 0 1", true},
	} {
		pos, err := ParseFEN(test.fen)
		if err != nil {
			t.Fatal(err)
		}
		if got := canCastle(pos, KS); got != test.want {
			t.Errorf("canCastleKS(%v): got %v, want %v", test.fen, got, test.want)
		}
	}
}

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
			t.Errorf("algebraic(%v, %+v): got %v, want %v", test.fen, test.move, got, test.alg)
		}
	}
}

func BenchmarkConstructMove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Move{From: b7, To: a8, Piece: Pawn, CapturePiece: Rook, PromotePiece: Queen}
	}
}

func BenchmarkReadMove(b *testing.B) {
	m := Move{From: b7, To: a8, Piece: Pawn, CapturePiece: Rook, PromotePiece: Queen}
	for i := 0; i < b.N; i++ {
		_ = m.From
		_ = m.To
		_ = m.Piece
		_ = m.CapturePiece
		_ = m.EP
		_ = m.PromotePiece
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
		{"en passant", Move{From: d5, To: e6, Piece: Pawn, CapturePiece: Pawn, EP: true}},
		{"capture", Move{From: a1, To: a8, Piece: Rook, CapturePiece: Rook}},
		{"promotion", Move{From: b7, To: b8, Piece: Pawn, PromotePiece: Queen}},
		{"capture promotion", Move{From: b7, To: a8, Piece: Pawn, CapturePiece: Rook, PromotePiece: Queen}},
		{"castle queenside", Move{From: e1, To: c1, Piece: King}},
		{"castle kingside", Move{From: e1, To: g1, Piece: King}},
	} {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Make(pos, benchmark.move)
			}
		})
	}
}

func BenchmarkPseudoLegalMoves(b *testing.B) {
	pos, err := ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		PseudoLegalMoves(pos)
	}
}
