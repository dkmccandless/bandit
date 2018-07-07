package main

const (
	opening = iota
	endgame

	pawnPhase   = 1
	knightPhase = 2
	bishopPhase = 2
	rookPhase   = 4
	queenPhase  = 8

	totalPhase = 2 * (queenPhase + 2*rookPhase + 2*bishopPhase + 2*knightPhase + 8*pawnPhase)
)

// The static evaluation of each type of Piece.
var pieceEval = [6][2]int{
	// opening, endgame
	{},
	{100, 100},
	{320, 300},
	{315, 335},
	{500, 525},
	{900, 900},
}

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
		// Black pawn
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
		// Black knight
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
		// Black bishop
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
		// Black rook
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
		// Black queen
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

var kingps = [2][2]pieceSquare{
	{},
	{
		// Black opening
		{
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-20, -30, -30, -40, -40, -30, -30, -20,
			-10, -20, -20, -20, -20, -20, -20, -10,
			20, 20, 0, 0, 0, 0, 20, 20,
			20, 30, 10, 0, 0, 10, 30, 20,
		},
		// Black endgame
		{
			-50, -40, -30, -20, -20, -30, -40, -50,
			-30, -20, -10, 0, 0, -10, -20, -30,
			-30, -10, 20, 30, 30, 20, -10, -30,
			-30, -10, 30, 40, 40, 30, -10, -30,
			-30, -10, 30, 40, 40, 30, -10, -30,
			-30, -10, 20, 30, 30, 20, -10, -30,
			-30, -30, 0, 0, 0, 0, -30, -30,
			-50, -30, -30, -30, -30, -30, -30, -50,
		},
	},
}

func init() {
	for piece := range ps[Black] {
		ps[White][piece] = ps[Black][piece].flip()
	}
	for phase := range kingps[Black] {
		kingps[White][phase] = kingps[Black][phase].flip()
	}
}

// Eval returns a Position's evaluation score in centipawns relative to White.
func Eval(pos Position) int {
	wp := PopCount(pos.b[White][Pawn])
	wn := PopCount(pos.b[White][Knight])
	wb := PopCount(pos.b[White][Bishop])
	wr := PopCount(pos.b[White][Rook])
	wq := PopCount(pos.b[White][Queen])

	bp := PopCount(pos.b[Black][Pawn])
	bn := PopCount(pos.b[Black][Knight])
	bb := PopCount(pos.b[Black][Bishop])
	br := PopCount(pos.b[Black][Rook])
	bq := PopCount(pos.b[Black][Queen])

	npawns, nknights, nbishops, nrooks, nqueens := wp+bp, wn+bn, wb+bb, wr+br, wq+bq

	// Check for draw due to insufficient material
	if npawns == 0 && nrooks == 0 && nqueens == 0 {
		switch {
		case nknights+nbishops <= 1:
			// KvK, KNvK, KBvK
			return 0
		case nknights == 0:
			if bishops := (pos.b[White][Bishop] | pos.b[Black][Bishop]); bishops&DarkSquares == 0 || bishops&LightSquares == 0 {
				// kings and any number of same color bishops
				return 0
			}
		}
	}

	phase := queenPhase*nqueens + rookPhase*nrooks + bishopPhase*nbishops + knightPhase*nknights + pawnPhase*npawns

	var eval int
	for sq := a1; sq <= h8; sq++ {
		switch c, p := pos.PieceOn(sq); {
		case p == None:
			continue
		case p == King:
			eval += evalMult(c) * taper(kingps[c][opening][sq], kingps[c][endgame][sq], phase)
		default:
			eval += evalMult(c) * (taper(pieceEval[p][opening], pieceEval[p][endgame], phase) + ps[c][p][sq])
		}
	}
	return eval
}

func taper(open, end, phase int) int {
	return (open*phase + end*(totalPhase-phase)) / totalPhase
}

// evalMult returns an evaluation sign multiplier of 1 for White and -1 for Black.
func evalMult(c Color) int {
	return 1 - 2*int(c)
}
