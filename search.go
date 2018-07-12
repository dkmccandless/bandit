package main

import (
	"context"
	"fmt"
	"sort"
)

const (
	evalInf = 50000
)

// SearchPosition searches a Position to the specified depth via iterative deepening
// and returns the search results.
func SearchPosition(ctx context.Context, pos Position, depth int) Results {
	var results Results
	for d := 1; d <= depth; d++ {
		_, results = negamax(ctx, pos, results, NewWindow(-evalInf, evalInf), d, true, make([]int, d+1))
		if ctx.Err() != nil {
			break
		}
	}
	return results
}

// negamax recursively searches a Position to the specified depth and returns the evaluation score
// relative to the side to move and the search results. It employs alpha-beta pruning outside of
// the specified Window. If recommended is zero length, negamax will generate and search all legal
// moves; if recommended moves are provided, they must all be legal, and only they will be searched.
func negamax(ctx context.Context, pos Position, recommended Results, w Window, depth int, allowCutoff bool, counters []int) (bestScore int, results Results) {
	counters[0]++

	if len(recommended) == 0 {
		moves := Candidates(pos) // pseudo-legal
		recommended = make(Results, 0, len(moves))
		for _, m := range moves {
			if !IsLegal(Make(pos, m)) {
				continue
			}
			recommended = append(recommended, Result{move: m})
		}
		if len(recommended) == 0 { // checkmate or stalemate
			if IsCheck(pos) {
				bestScore = -evalInf
			}
			return
		}
		if depth == 0 {
			bestScore = Eval(pos) * evalMult(pos.ToMove)
			return
		}
	}
	// Invariant: len(recommended) > 0 and recommended contains only legal moves

	results = make(Results, 0, len(recommended))
	defer func() { results.SortFor(pos.ToMove) }()

	for _, r := range recommended {
		score, cont := negamax(ctx, Make(pos, r.move), Results{}, w.Neg(), depth-1, allowCutoff, counters[1:])
		score *= -1

		// Store the score in results relative to White
		results = results.Update(Result{move: r.move, score: score * evalMult(pos.ToMove), depth: depth - 1, cont: cont})

		if depth >= 3 && ctx.Err() != nil {
			break
		}

		var constrained, ok bool
		w, constrained, ok = w.Constrain(score)
		if constrained {
			// improved lower bound
		}
		if !ok && allowCutoff {
			break
		}
	}

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

// A Result holds the score of a searched Move (relative to White),
// the search depth, and the continuation Results of the search.
type Result struct {
	move  Move
	score int
	depth int
	cont  Results
}

// String returns a string representation of r, including its principal variation.
func (r Result) String() string {
	return fmt.Sprintf("%v (%v) %v", float64(r.score)/100, r.depth, r.PV())
}

// PV returns a string representation of r's principal variation.
func (r Result) PV() string {
	if r.depth == 0 || len(r.cont) == 0 {
		return LongAlgebraic(r.move)
	}
	return LongAlgebraic(r.move) + " " + r.cont[0].PV()
}

// A Results contains the results of a search. Results satisfies sort.Interface.
// A Results should be sorted by the function generating it before it is returned.
type Results []Result

func (rs Results) Len() int      { return len(rs) }
func (rs Results) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }

// Less reports whether the Result with index i should sort before the Result with index j.
// It sorts first by depth and then by score.
func (rs Results) Less(i, j int) bool {
	return rs[i].depth < rs[j].depth || rs[i].depth == rs[j].depth && rs[i].score < rs[j].score
}

// SortFor sorts r beginning with the best move for c.
func (rs Results) SortFor(c Color) {
	if c == White {
		// highest score first
		sort.Sort(sort.Reverse(rs))
	} else {
		sort.Sort(rs)
	}
}

func (rs Results) Update(r Result) Results {
	for i := range rs {
		if rs[i].move == r.move {
			rs[i] = r
			return rs
		}
	}
	return append(rs, r)
}

// String returns a string representation of all Result values in r.
func (rs Results) String() string {
	var s string
	for _, r := range rs {
		s += fmt.Sprintf("%v\n", r.String())
	}
	return s
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
// whether the lower bound was changed, and whether the returned Window remains valid.
func (w Window) Constrain(n int) (c Window, constrained bool, ok bool) {
	if n <= w.alpha {
		return w, false, true
	}
	return Window{n, w.beta}, true, n <= w.beta
}

// Neg returns the additive inverse of w.
func (w Window) Neg() Window { return Window{-w.beta, -w.alpha} }
