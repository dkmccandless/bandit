package main

import (
	"fmt"
	"sort"
)

const (
	evalInf = 50000
)

func SearchPosition(pos Position, depth int) (int, Results) {
	var score int
	var results Results
	for d := 1; d <= depth; d++ {
		score, results = negamax(pos, results, NewWindow(-evalInf, evalInf), d, true, negamax)
	}
	return score, results
}

type SearchFunc func(Position, Results, Window, int, bool, SearchFunc) (int, Results)

func negamax(pos Position, recommended Results, w Window, depth int, allowCutoff bool, search SearchFunc) (bestScore int, results Results) {
	if len(recommended) == 0 {
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

		recommended = make(Results, 0, len(moves))
		for _, m := range moves {
			recommended = append(recommended, Result{move: m})
		}
	}

	results = make(Results, 0, len(recommended))
	for _, r := range recommended {
		newpos := Make(pos, r.move)
		if !IsLegal(newpos) {
			continue
		}

		score, _ := search(newpos, Results{}, w.Neg(), depth-1, allowCutoff, search)
		score *= -1

		results = append(results, Result{move: r.move, score: score, depth: depth - 1})

		var constrained, ok bool
		w, constrained, ok = w.Constrain(score)
		if constrained {
			// improved lower bound
		}
		if !ok && allowCutoff {
			break
		}
	}

	sort.Sort(results)
	return w.alpha, results
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

// A Result holds the score of a searched Move and the search depth.
type Result struct {
	move  Move
	score int
	depth int
}

// A Results contains the results of a search. Results satisfies sort.Interface.
// A Results should be sorted by the function generating it before it is returned.
type Results []Result

func (r Results) Len() int      { return len(r) }
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Less sorts two Result elements first by depth and then by score.
func (r Results) Less(i, j int) bool {
	return r[i].depth > r[j].depth || r[i].depth == r[j].depth && r[i].score > r[j].score
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
