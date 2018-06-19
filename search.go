package main

import "fmt"

const (
	evalInf = 50000
)

func SearchPosition(pos Position, depth int) (int, Move) {
	var score int
	var recommended Move
	for d := 1; d <= depth; d++ {
		score, recommended = negamax(pos, recommended, NewWindow(-evalInf, evalInf), d, true, negamax)
	}
	return score, recommended
}

type SearchFunc func(Position, Move, Window, int, bool, SearchFunc) (int, Move)

func negamax(pos Position, recommended Move, w Window, depth int, allowCutoff bool, search SearchFunc) (bestScore int, bestMove Move) {
	moves := Candidates(pos) // pseudo-legal

	if !anyLegal(pos, moves) { // checkmate or stalemate
		if IsCheck(pos) {
			bestScore = -evalInf
		}
		return
	}
	if depth == 0 {
		bestScore = Eval(pos) * evalMult(pos.ToMove)
		return
	}

	// Initialize bestMove with a legal Move if no recommended Move is provided
	if recommended != (Move{}) {
		moves = reorder(moves, recommended)
		bestMove = recommended
	} else {
		for i, m := range moves {
			if IsLegal(Make(pos, m)) {
				bestMove = m
				moves = moves[i:]
				break
			}
		}
	}

	for _, m := range moves {
		newpos := Make(pos, m)
		if !IsLegal(newpos) {
			continue
		}

		score, _ := search(newpos, Move{}, w.Neg(), depth-1, allowCutoff, search)
		score *= -1

		var constrained, ok bool
		w, constrained, ok = w.Constrain(score)
		if constrained {
			bestMove = m
		}
		if !ok && allowCutoff {
			break
		}
	}
	return w.alpha, bestMove
}

// IsPseudoLegal returns whether a Move is pseudo-legal in a Position.
// A move is pseudo-legal if the square to be moved from contains the specified piece
// and the piece is capable of moving to the target square if doing so would not put the king in check.
func IsPseudoLegal(pos Position, move Move) bool {
	for _, m := range PieceMoves[move.Piece](pos) {
		if m == move {
			return true
		}
	}
	return false
}

// IsLegal returns whether a Position results from a legal move.
// A position is illegal if the king of the side that just moved is in check.
func IsLegal(pos Position) bool {
	return !IsAttacked(pos, pos.KingSquare[pos.Opp()], pos.ToMove)
}

// IsCheck returns whether the king of the side to move is in check.
func IsCheck(pos Position) bool {
	return IsAttacked(pos, pos.KingSquare[pos.ToMove], pos.Opp())
}

// IsTerminal returns whether or not a Position is checkmate or stalemate.
// A position is checkmate or stalemate if the side to move has no legal moves.
func IsTerminal(pos Position) bool { return !anyLegal(pos, Candidates(pos)) }

// anyLegal returns whether any of the given Moves are legal in the given Position.
func anyLegal(pos Position, moves []Move) bool {
	for _, m := range moves {
		if IsLegal(Make(pos, m)) {
			return true
		}
	}
	return false
}

// reorder returns a reordered slice of Moves with the specified Move first.
// The slice is not modified if it does not contain the specified Move.
func reorder(moves []Move, m Move) []Move {
	for _, n := range moves {
		if n == m {
			s := []Move{m}
			for _, n := range moves {
				if n != m {
					s = append(s, n)
				}
			}
			return s
		}
	}
	return moves
}

// A Window represents the bounds of a position's evaluation.
type Window struct{ alpha, beta int }

// NewWindow returns a Window with the given bounds. It panics if alpha > beta.
func NewWindow(alpha, beta int) Window {
	if alpha > beta {
		panic(fmt.Sprintf("invalid window bounds %v, %v", alpha, beta))
	}
	return Window{alpha, beta}
}

// Constrain updates the lower bound of w, if applicable, and returns the updated window,
// whether the lower bound was changed, and whether w remains a valid Window.
func (w Window) Constrain(n int) (c Window, constrained bool, ok bool) {
	if n <= w.alpha {
		return w, false, true
	}
	return Window{n, w.beta}, true, n <= w.beta
}

// Neg returns the additive inverse of w.
func (w Window) Neg() Window { return Window{-w.beta, -w.alpha} }
