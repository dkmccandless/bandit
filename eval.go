package main

var pieceEval = [6]int{0, 100, 320, 333, 500, 900}

// Piece-square tables modifying the evaluation of a Piece depending on its Square.
type pieceSquare [64]int

func (ps pieceSquare) flip() pieceSquare {
	var f pieceSquare
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			f[8*rank+file] = ps[8*(7-rank)+file]
		}
	}
	return f
}

// Values based on Tomasz Michniewski's "Unified Evaluation" test tournament tables
var ps = [2][6]pieceSquare{
	{},
	{
		{},
		{
			0, 0, 0, 0, 0, 0, 0, 0,
			50, 50, 50, 50, 50, 50, 50, 50,
			10, 10, 20, 30, 30, 20, 10, 10,
			5, 5, 10, 25, 25, 10, 5, 5,
			0, 0, 0, 20, 20, 0, 0, 0,
			5, -5, -10, 0, 0, -10, -5, 5,
			5, 10, 10, -20, -20, 10, 10, 5,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		{
			-50, -40, -30, -30, -30, -30, -40, -50,
			-40, -20, 0, 5, 5, 0, -20, -40,
			-30, 0, 15, 20, 20, 15, 0, -30,
			-30, 5, 20, 25, 25, 20, 5, -30,
			-30, 0, 15, 20, 20, 15, 0, -30,
			-30, 0, 10, 15, 15, 10, 0, -30,
			-40, -20, 0, 5, 5, 0, -20, -40,
			-50, -40, -30, -30, -30, -30, -40, -50,
		},
		{
			-20, -10, -10, -10, -10, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 5, 5, 5, 0, -10,
			-10, 0, 5, 10, 10, 5, 0, -10,
			-10, 0, 5, 10, 10, 5, 0, -10,
			-10, 0, 5, 5, 5, 5, 0, -10,
			-10, 5, 0, 0, 0, 0, 5, -10,
			-20, -10, -10, -10, -10, -10, -10, -20,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0,
			5, 10, 10, 10, 10, 10, 10, 5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			0, 0, 0, 5, 5, 0, 0, 0,
		},
		{
			-20, -10, -10, -5, -5, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 5, 5, 5, 0, -10,
			-5, 0, 5, 5, 5, 5, 0, -5,
			0, 0, 5, 5, 5, 5, 0, -5,
			-10, 5, 5, 5, 5, 5, 0, -10,
			-10, 0, 5, 0, 0, 0, 0, -10,
			-20, -10, -10, -5, -5, -10, -10, -20,
		},
	},
}

func init() {
	for piece := range ps[White] {
		ps[White][piece] = ps[Black][piece].flip()
	}
}

// Eval returns a Position's evaluation score in centipawns relative to White.
func Eval(pos Position) int {
	var whiteMaterial, blackMaterial int
	for sq := a1; sq <= h8; sq++ {
		c, p, ok := pos.PieceOn(sq)
		switch {
		case !ok:
			continue
		case p == King:
			continue
		case c == White:
			whiteMaterial += pieceEval[p] + ps[White][p][sq]
		case c == Black:
			blackMaterial += pieceEval[p] + ps[Black][p][sq]
		}
	}
	return whiteMaterial - blackMaterial
}

// evalMult returns an evaluation sign multiplier of 1 for White and -1 for Black.
func evalMult(c Color) int {
	return 1 - 2*int(c)
}
