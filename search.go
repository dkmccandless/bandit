package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

const (
	evalInf = 50000
)

var (
	errCheckmate checkmateError
	errStalemate = errors.New("stalemate")
	errFiftyMove = errors.New("fifty-move rule")
	// TODO: threefold repetition
)

type Search struct {
	allowCutoff bool
	counters    []int
}

// SearchPosition searches a Position to the specified depth via iterative deepening
// and returns the search results.
func SearchPosition(ctx context.Context, pos Position, depth int) Results {
	var rs Results
	s := Search{
		allowCutoff: true,
		counters:    make([]int, depth+1),
	}
	for d := 1; d <= depth; d++ {
		_, rs, _ = s.negamax(ctx, pos, rs, Window{-evalInf, evalInf}, d)
		if ctx.Err() != nil {
			break
		}
	}
	return rs
}

// negamax recursively searches a Position to the specified depth and returns the evaluation score
// relative to the side to move and the search results. It employs alpha-beta pruning outside of
// the specified Window. If rs is zero length, negamax will generate and search all legal
// moves; if recommended moves are provided, they must all be legal, and only they will be searched.
func (s *Search) negamax(
	ctx context.Context,
	pos Position,
	rs Results,
	w Window,
	depth int,
) (bestScore RelScore, results Results, err error) {
	s.counters[len(s.counters)-1-depth]++

	score, rs, err := checkDone(pos, rs)
	if err != nil && (s.allowCutoff || err != errInsufficient) {
		// Do not cut off during perft in the case of insufficient material
		return score, rs, err
	}
	// Invariant: len(rs) > 0 and rs contains only legal moves
	if depth == 0 {
		score, err := Eval(pos)
		return score.Rel(pos.ToMove), rs, err
	}
	if s.allowCutoff && deepEnough(rs, depth) {
		return rs[0].score.Rel(pos.ToMove), rs, rs[0].err
	}

	for _, r := range rs {
		if r.err != nil && s.allowCutoff {
			// The move is already known not to avoid a game-ending state; no need to search it further
			continue
		}

		score, cont, err := s.negamax(ctx, Make(pos, r.move), r.cont, w.Neg(), depth-1)
		score *= -1
		if e, ok := err.(checkmateError); ok {
			err = e.Prev()
		}

		rs.Update(Result{move: r.move, score: score.Abs(pos.ToMove), depth: depth - 1, cont: cont, err: err})

		if depth >= 3 && ctx.Err() != nil {
			break
		}

		var constrained, ok bool
		w, constrained, ok = w.Constrain(score)
		if constrained {
			// improved lower bound
		}
		if !ok && s.allowCutoff {
			break
		}
	}
	rs.SortFor(pos.ToMove)
	return w.alpha, rs, rs[0].err
}

// checkDone reports whether pos represents a game-ending position.
// If so, the RelScore and error value indicate the type of ending.
// If not, the RelScore is 0 and the error is nil.
// In either case, checkDone returns a Results of all legal moves in pos.
// If rs is not provided, checkDone generates the legal moves first.
func checkDone(pos Position, rs Results) (RelScore, Results, error) {
	if len(rs) == 0 {
		moves := LegalMoves(pos)
		rs = make(Results, 0, len(moves))
		for _, m := range moves {
			rs = append(rs, Result{move: m})
		}
	}
	if len(rs) == 0 {
		// no legal moves
		if IsCheck(pos) {
			return -evalInf, rs, errCheckmate
		}
		return 0, rs, errStalemate
	}
	if pos.HalfMove == 100 {
		// fifty-move rule
		return 0, rs, errFiftyMove
	}
	return 0, rs, nil
}

// deepEnough reports whether rs stores the results of a position search to at least the specified depth.
func deepEnough(rs Results, depth int) bool {
	// rs is already deep enough if all non-terminal elements have been searched to at least depth-1,
	// as long as they have been searched to non-zero depth
	for _, r := range rs {
		if (r.depth == 0 || r.depth < depth-1) && r.err == nil {
			return false
		}
	}
	return true
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

// IsLegal returns whether a Position is legal.
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
func IsTerminal(pos Position) bool { return len(LegalMoves(pos)) == 0 }

// A Result holds a searched Move along with its evaluated Score,
// the search depth, and the continuation Results of the search.
type Result struct {
	move  Move
	score Score
	depth int
	cont  Results
	err   error
}

// String returns a string representation of r, including its principal variation.
func (r Result) String() string {
	s := fmt.Sprintf("%v (%v) %v", r.score, r.depth, r.PV())
	if r.err != nil {
		return s + " " + r.err.Error()
	}
	return s
}

// PV returns a string representation of r's principal variation.
func (r Result) PV() string {
	if r.depth == 0 || len(r.cont) == 0 {
		return LongAlgebraic(r.move)
	}
	return LongAlgebraic(r.move) + " " + r.cont[0].PV()
}

// A Results contains the results of a search.
// Functions returning a Results should sort it before returning.
type Results []Result

// SortFor sorts rs beginning with the best move for c.
// The sort is not guaranteed to be stable.
func (rs Results) SortFor(c Color) {
	// Sort first by mate condition, then by depth decreasing,
	// then by Score decreasing/increasing for White/Black,
	// and then by origin and destination Square increasing
	sort.Slice(rs, func(i, j int) bool {
		if less, ok := rs.mateSort(i, j); ok {
			return less
		}
		if rs[i].depth != rs[j].depth {
			return rs[i].depth > rs[j].depth
		}
		if rs[i].score != rs[j].score {
			return (rs[i].score > rs[j].score) == (c == White)
		}
		return rs.squareSort(i, j)
	})
}

// SortBySquares sorts rs first by the Result moves' origin Squares and then by their destination Squares,
// and finally by promotion piece in the case of pawn promotions.
func (rs Results) SortBySquares() { sort.Slice(rs, rs.squareSort) }

// mateSort reports how two Result elements should be sorted if at least one of them leads to checkmate.
// ok reports whether either Result contains an error of type checkmateError.
// If ok is true, less reports whether the Result with index i should sort before the Result with index j.
// If ok is false, the value of less is meaningless.
func (rs Results) mateSort(i, j int) (less bool, ok bool) {
	ich, iok := rs[i].err.(checkmateError)
	jch, jok := rs[j].err.(checkmateError)
	switch {
	case iok && jok:
		switch iwin, jwin := ich&1 != 0, jch&1 != 0; {
		case iwin && jwin:
			// faster winning mate first
			return ich < jch, true
		case !iwin && !jwin:
			// slower losing mate first
			return ich > jch, true
		default:
			return iwin, true
		}
	case iok:
		return ich&1 != 0, true
	case jok:
		return jch&1 == 0, true
	}
	return false, false
}

// squareSort reports how two Result elements should be sorted by their moves' origin and destination Squares,
// and then finally by promotion piece in the case of promoting pawns.
func (rs Results) squareSort(i, j int) bool {
	if rs[i].move.From != rs[j].move.From {
		return rs[i].move.From < rs[j].move.From
	}
	if rs[i].move.To != rs[j].move.To {
		return rs[i].move.To < rs[j].move.To
	}
	// From and To are non-unique in the case of pawn promotion
	return rs[i].move.PromotePiece < rs[j].move.PromotePiece
}

// Update finds the Result in rs with the same Move as r and replaces it with r.
// It panics if rs does not contain any elements with r's Move.
func (rs Results) Update(r Result) {
	for i := range rs {
		if rs[i].move == r.move {
			rs[i] = r
			return
		}
	}
	// rs should contain all legal moves
	panic("unreached")
}

// String returns a string representation of all Result values in r.
func (rs Results) String() string {
	var s string
	for _, r := range rs {
		s += fmt.Sprintf("%v\n", r.String())
	}
	return s
}

// A checkmateError indicates the number of plies until a forced checkmate can be delivered.
// Odd values are wins for the player with the next move, and even values for the player with the previous move.
// The zero value of type checkmateError indicates that the current position is checkmate.
type checkmateError int

func (n checkmateError) Error() string {
	if n&1 == 0 {
		return fmt.Sprintf("-#%d", n/2)
	}
	return fmt.Sprintf("#%d", (n+1)/2)
}

// Prev returns the checkmateError corresponding to the previous ply.
func (n checkmateError) Prev() checkmateError { return n + 1 }

// A Window represents the bounds of a position's evaluation.
type Window struct{ alpha, beta RelScore }

// Constrain updates the lower bound of w, if applicable, and returns the updated window,
// whether the lower bound was changed, and whether the returned Window remains valid.
func (w Window) Constrain(n RelScore) (c Window, constrained bool, ok bool) {
	if n <= w.alpha {
		return w, false, true
	}
	return Window{n, w.beta}, true, n <= w.beta
}

// Neg returns the additive inverse of w.
func (w Window) Neg() Window { return Window{-w.beta, -w.alpha} }
