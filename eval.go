package main

const (
	PawnEval   = 100
	KnightEval = 320
	BishopEval = 333
	RookEval   = 500
	QueenEval  = 900
)

// evalMult returns an evaluation sign multiplier of 1 for White and -1 for Black.
func evalMult(c Color) int {
	if c == White {
		return 1
	} else {
		return -1
	}
}

// Eval returns a Position's evaluation score in centipawns relative to White.
func Eval(pos Position) int {
	whiteMaterial :=
		PopCount(pos.b[White][Pawn])*PawnEval +
			PopCount(pos.b[White][Knight])*KnightEval +
			PopCount(pos.b[White][Bishop])*BishopEval +
			PopCount(pos.b[White][Rook])*RookEval +
			PopCount(pos.b[White][Queen])*QueenEval
	blackMaterial :=
		PopCount(pos.b[Black][Pawn])*PawnEval +
			PopCount(pos.b[Black][Knight])*KnightEval +
			PopCount(pos.b[Black][Bishop])*BishopEval +
			PopCount(pos.b[Black][Rook])*RookEval +
			PopCount(pos.b[Black][Queen])*QueenEval
	return whiteMaterial - blackMaterial
}
