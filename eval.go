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

// Score represents the engine's evaluation of a Position in centipawns.
// When the evaluation indicates a game-ending condition, err describes the condition.
type Score struct {
	n   int
	err error
}

// String returns a string representation of s.
func (s Score) String() string {
	if s.err != nil {
		return s.err.Error()
	}
	return fmt.Sprintf("%.2f", float64(s.n)/100)
}

// Less reports whether a is lower than b.
// Mates in an even number of plies are lower than any non-mate score,
// mates in an odd number of plies are higher than any non-mate score,
// and mates in fewer plies are more extremal than mates in more plies.
func Less(a, b Score) bool {
	aply, ach := a.err.(checkmateError)
	bply, bch := b.err.(checkmateError)
	switch {
	case ach && bch:
		switch awin, bwin := aply&1 != 0, bply&1 != 0; {
		case awin && bwin:
			return aply > bply
		case !awin && !bwin:
			return aply < bply
		default:
			return bwin
		}
	case ach:
		return aply&1 == 0
	case bch:
		return bply&1 != 0
	}
	return a.n < b.n
}

// Abs represents the engine's evaluation relative to White.
type Abs Score

// Rel returns the Rel of s with respect to c.
// This is equal to s for White and -s for Black.
func (s Abs) Rel(c Color) Rel {
	if c == White {
		return Rel(s)
	}
	return Rel{-s.n, s.err}
}

// String returns a string representation of s.
func (s Abs) String() string { return Score(s).String() }

// Rel represents the engine's evaluation relative to the side to move.
type Rel Score

// Abs returns the Abs of s relative to White, where s is with respect to c.
// This is equal to s for White and -s for Black.
func (s Rel) Abs(c Color) Abs {
	if c == White {
		return Abs(s)
	}
	return Abs{-s.n, s.err}
}

// String returns a string representation of s.
func (s Rel) String() string { return Score(s).String() }

// relcent represents the numerical component of an evaluation in relative centipawns.
type relcent int

// pieceEval represents the static evaluation of each type of Piece.
var pieceEval = [6][2]relcent{
	// opening, endgame
	{},
	{100, 100},
	{320, 300},
	{315, 335},
	{500, 525},
	{900, 900},
}

// pieceSquare is a slice of relcents modifying the evaluation of a Piece depending on its Square.
type pieceSquare [64]relcent

// ps provides piece-square tables for each Piece of each Color.
// Values are based on Tomasz Michniewski's "Unified Evaluation" test tournament tables.
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

// flip returns the vertical transposition of ps.
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

// Eval returns a Position's Abs evaluation score.
// It returns errInsufficient in the case of insufficient material.
func Eval(pos Position) Abs {
	if IsInsufficient(pos) {
		return Abs{err: errInsufficient}
	}
	var (
		npawns   = PopCount(pos.b[White][Pawn] | pos.b[Black][Pawn])
		nknights = PopCount(pos.b[White][Knight] | pos.b[Black][Knight])
		nbishops = PopCount(pos.b[White][Bishop] | pos.b[Black][Bishop])
		nrooks   = PopCount(pos.b[White][Rook] | pos.b[Black][Rook])
		nqueens  = PopCount(pos.b[White][Queen] | pos.b[Black][Queen])
		phase    = queenPhase*nqueens + rookPhase*nrooks + bishopPhase*nbishops + knightPhase*nknights + pawnPhase*npawns
		eval     int
	)
	for sq := a1; sq <= h8; sq++ {
		c, p := pos.PieceOn(sq)
		if p == None {
			continue
		}
		var r relcent
		switch p {
		case King:
			r = taper(kingps[c][opening][sq], kingps[c][endgame][sq], phase)
		default:
			r = taper(pieceEval[p][opening], pieceEval[p][endgame], phase) + ps[c][p][sq]
		}
		if c == White {
			eval += int(r)
		} else {
			eval -= int(r)
		}
	}
	return Abs{n: eval}
}

// IsInsufficient reports whether pos contains insufficient material to deliver checkmate.
// This is a subset of the condition of impossibility of checkmate, which results in an automatic draw.
func IsInsufficient(pos Position) bool {
	if pos.b[White][Pawn] != 0 || pos.b[Black][Pawn] != 0 ||
		pos.b[White][Rook] != 0 || pos.b[Black][Rook] != 0 ||
		pos.b[White][Queen] != 0 || pos.b[Black][Queen] != 0 {
		return false
	}
	var (
		nknights = PopCount(pos.b[White][Knight] | pos.b[Black][Knight])
		nbishops = PopCount(pos.b[White][Bishop] | pos.b[Black][Bishop])
	)
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
func taper(open, end relcent, phase int) relcent {
	return (open*relcent(phase) + end*(totalPhase-relcent(phase))) / totalPhase
}
