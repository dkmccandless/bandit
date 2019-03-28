package main

import (
	"errors"
	"fmt"
)

const (
	opening = iota
	endgame

	pawnPhase = 1 << iota
	bishopPhase
	rookPhase
	queenPhase
	knightPhase = bishopPhase

	totalPhase = 2 * (queenPhase + 2*rookPhase + 2*bishopPhase + 2*knightPhase + 8*pawnPhase)
)

// errInsufficient is returned by Eval when neither player has sufficient material
// to deliver checkmate by any sequence of legal moves.
var errInsufficient = errors.New("insufficient material")

// A Score represents the engine's evaluation of a Position in centipawns relative to White.
type Score int

// String returns a string representation of s.
func (s Score) String() string { return fmt.Sprintf("%.2f", float64(s)/100) }

// Rel returns the RelScore of s with respect to c.
// This is equal to s for White and -s for Black.
func (s Score) Rel(c Color) RelScore {
	r := RelScore(s)
	if c == Black {
		return -r
	}
	return r
}

// A RelScore represents the engine's evaluation of a Position in centipawns relative to the side to move.
type RelScore int

// Abs returns the Score of r relative to White, where r is with respect to c.
// This is equal to r for White and -r for Black.
func (r RelScore) Abs(c Color) Score {
	s := Score(r)
	if c == Black {
		return -s
	}
	return s
}

// The static evaluation of each type of Piece.
var pieceEval = [6][2]RelScore{
	// opening, endgame
	{},
	{100, 100},
	{320, 300},
	{315, 335},
	{500, 525},
	{900, 900},
}

// Piece-square tables modifying the evaluation of a Piece depending on its Square.
type pieceSquare [64]RelScore

// Values based on Tomasz Michniewski's "Unified Evaluation" test tournament tables
var ps = [2][6]pieceSquare{
	// Generate White tables in init
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

func (ps pieceSquare) flip() pieceSquare {
	var f pieceSquare
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			f[8*rank+file] = ps[8*(7-rank)+file]
		}
	}
	return f
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
// It returns errInsufficient in the case of insufficient material.
func Eval(pos Position) (Score, error) {
	var (
		npawns   = PopCount(pos.b[White][Pawn] | pos.b[Black][Pawn])
		nknights = PopCount(pos.b[White][Knight] | pos.b[Black][Knight])
		nbishops = PopCount(pos.b[White][Bishop] | pos.b[Black][Bishop])
		nrooks   = PopCount(pos.b[White][Rook] | pos.b[Black][Rook])
		nqueens  = PopCount(pos.b[White][Queen] | pos.b[Black][Queen])
	)

	if IsInsufficient(pos, npawns, nknights, nbishops, nrooks, nqueens) {
		return 0, errInsufficient
	}

	phase := queenPhase*nqueens + rookPhase*nrooks + bishopPhase*nbishops + knightPhase*nknights + pawnPhase*npawns
	var eval Score
	for sq := a1; sq <= h8; sq++ {
		switch c, p := pos.PieceOn(sq); {
		case p == None:
			continue
		case p == King:
			eval += taper(kingps[c][opening][sq], kingps[c][endgame][sq], phase).Abs(c)
		default:
			eval += (taper(pieceEval[p][opening], pieceEval[p][endgame], phase) + ps[c][p][sq]).Abs(c)
		}
	}
	return eval, nil
}

// IsInsufficient reports whether a collection of pieces constitutes insufficient material
// to deliver checkmate. This is a subset of the condition of impossibility of checkmate,
// which results in an automatic draw.
func IsInsufficient(pos Position, npawns, nknights, nbishops, nrooks, nqueens int) bool {
	if npawns > 0 || nrooks > 0 || nqueens > 0 {
		return false
	}
	if nknights+nbishops <= 1 {
		// KvK, KNvK, KBvK
		return true
	}
	if nknights > 0 {
		// knight and at least one other minor piece; in KNNvK, KBvKN, and KNvKN,
		// mate can't be forced, although it can be given by a series of legal moves
		return false
	}
	if bishops := (pos.b[White][Bishop] | pos.b[Black][Bishop]); bishops&DarkSquares == 0 || bishops&LightSquares == 0 {
		// kings and any number of same color bishops only
		return true
	}
	// opposite color bishops; KBvKB can mate, although not by force
	return false
}

// taper returns the weighted sum of open and end according to the fraction phase/totalPhase.
// This mitigates evaluation discontinuity in the event of rapid loss of material.
func taper(open, end RelScore, phase int) RelScore {
	return (open*RelScore(phase) + end*(totalPhase-RelScore(phase))) / totalPhase
}
