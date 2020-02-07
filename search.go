package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

var (
	errCheckmate checkmateError
	errStalemate = errors.New("stalemate")
	errFiftyMove = errors.New("fifty-move rule")
	// TODO: threefold repetition

	// openWindow encompasses all possible Rels.
	openWindow = Window{
		Rel{err: errCheckmate},        // worst case: mated
		Rel{err: errCheckmate.Prev()}, // best case: mate in 1 ply
	}
)

// Search contains parameters relevant to a position search.
type Search struct {
	// allowCutoff describes whether to allow alpha-beta pruning.
	allowCutoff bool

	// counters tracks the number of nodes searched at each depth.
	counters []int
}

// SearchPosition searches a Position to the specified depth via iterative deepening and returns the search results.
func SearchPosition(ctx context.Context, pos Position, depth int) Results {
	var rs Results
	s := Search{
		allowCutoff: true,
		counters:    make([]int, depth+1),
	}
	for d := 1; d <= depth; d++ {
		_, rs = s.negamax(ctx, pos, rs, openWindow, d)
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
) (bestScore Rel, results Results) {
	s.counters[len(s.counters)-1-depth]++

	err := checkTerminal(pos)
	if err != nil && (s.allowCutoff || err != errInsufficient) {
		// Do not cut off during perft in the case of insufficient material
		return Rel{err: err}, nil
	}
	if len(rs) == 0 {
		rs = legalResults(pos)
	}
	// Invariant: len(rs) > 0 and rs contains only legal moves
	if w.beta.err == errCheckmate {
		// An alternative move from this position's parent delivers mate; no need to search this one.
		return rs[0].score.Rel(pos.ToMove), rs
	}
	if depth == 0 {
		score := Eval(pos)
		return score.Rel(pos.ToMove), rs
	}
	if s.allowCutoff && deepEnough(rs, depth) {
		return rs[0].score.Rel(pos.ToMove), rs
	}
	if w.alpha.err == errCheckmate {
		// There is at least one legal move, so the worst case is not checkmate.
		w = Window{Rel{err: errCheckmate.Prev().Prev()}, w.beta}
	}

	for _, r := range rs {
		if r.score.err != nil && s.allowCutoff {
			// The move is already known not to avoid a game-ending state; no need to search it further.
			continue
		}

		score, cont := s.negamax(ctx, Make(pos, r.move), r.cont, w.Next(), depth-1)
		score = score.Prev()

		rs.Update(Result{move: r.move, score: score.Abs(pos.ToMove), depth: depth - 1, cont: cont})

		if depth >= 3 && ctx.Err() != nil {
			break
		}

		if !s.allowCutoff {
			continue
		}
		var ok bool
		w, ok = w.Constrain(score)
		if !ok {
			// beta cutoff
			break
		}
	}
	rs.SortFor(pos.ToMove)
	return w.alpha, rs
}

// checkTerminal returns an error describing the type of terminal position represented by pos, or nil if pos is not terminal.
func checkTerminal(pos Position) error {
	if IsMate(pos) {
		if IsCheck(pos) {
			return errCheckmate
		}
		return errStalemate
	}
	if pos.HalfMove >= 100 {
		// fifty-move rule
		return errFiftyMove
	}
	return nil
}

// legalResults returns a Results containing all legal moves in pos.
func legalResults(pos Position) Results {
	moves := LegalMoves(pos)
	rs := make(Results, 0, len(moves))
	for _, m := range moves {
		rs = append(rs, Result{move: m})
	}
	return rs
}

// deepEnough reports whether rs stores the results of a position search to at least the specified depth.
func deepEnough(rs Results, depth int) bool {
	// rs is already deep enough if all non-terminal elements have been searched to at least depth-1,
	// as long as they have been searched to non-zero depth.
	for _, r := range rs {
		if (r.depth == 0 || r.depth < depth-1) && r.score.err == nil {
			return false
		}
	}
	return true
}

// IsPseudoLegal returns whether a Move is pseudo-legal in a Position.
// A move is pseudo-legal if the square to be moved from contains the specified piece
// and the piece is capable of moving to the target square if doing so would not put the king in check.
func IsPseudoLegal(pos Position, move Move) bool {
	for _, m := range PseudoLegalMoves(pos) {
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

// IsMate returns whether or not a Position is checkmate or stalemate.
// A position is checkmate or stalemate if the side to move has no legal moves.
func IsMate(pos Position) bool { return len(LegalMoves(pos)) == 0 }

// Result holds a searched Move along with its evaluated score,
// the search depth, and the continuation Results of the search.
type Result struct {
	move  Move
	score Abs
	depth int
	cont  Results
}

// String returns a string representation of r, including its principal variation.
func (r Result) String() string {
	s := fmt.Sprintf("%v (%v) %v", r.score, r.depth, r.PV())
	if r.score.err != nil {
		return s + " " + r.score.err.Error()
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

// Results contains the results of a search. Functions returning a Results should sort it before returning.
type Results []Result

// SortFor sorts rs beginning with the best move for c.
// The sort is not guaranteed to be stable.
func (rs Results) SortFor(c Color) {
	// Sort first by terminal condition, then by depth decreasing,
	// then by Score decreasing/increasing for White/Black,
	// and then by origin and destination Square increasing.
	sort.Slice(rs, func(i, j int) bool {
		less := Less(Score(rs[i].score), Score(rs[j].score))
		if ierr, jerr := rs[i].score.err, rs[j].score.err; ierr != nil || jerr != nil {
			_, ich := ierr.(checkmateError)
			_, jch := jerr.(checkmateError)
			if ich || jch {
				return !less
			}
			return !less == (c == White)
		}
		if rs[i].depth != rs[j].depth {
			return rs[i].depth > rs[j].depth
		}
		if rs[i].score != rs[j].score {
			return !less == (c == White)
		}
		return rs.squareSort(i, j)
	})
}

// SortBySquares sorts rs first by the Result moves' origin Squares and then by their destination Squares,
// and finally by promotion piece in the case of pawn promotions.
func (rs Results) SortBySquares() { sort.Slice(rs, rs.squareSort) }

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

// checkmateError indicates the number of plies until a forced checkmate can be delivered.
// Odd values are wins for the player with the next move, and even values for the player with the previous move.
// The zero value of type checkmateError indicates that the current position is checkmate.
type checkmateError int

func (n checkmateError) Error() string {
	if n&1 == 0 {
		return fmt.Sprintf("-#%d", n/2)
	}
	return fmt.Sprintf("#%d", (n+1)/2)
}

// Prev returns the checkmateError corresponding to n's previous ply.
func (n checkmateError) Prev() checkmateError { return n + 1 }

// Next returns the checkmateError corresponding to n's following ply.
// It panics if n is not positive.
func (n checkmateError) Next() checkmateError {
	if n < 1 {
		panic(fmt.Sprintf("non-positive checkmateError %v", n))
	}
	return n - 1
}

// Prev returns the Rel corresponding to s's previous ply.
func (s Rel) Prev() Rel {
	if err, ok := s.err.(checkmateError); ok {
		return Rel{-s.n, err.Prev()}
	}
	return Rel{-s.n, s.err}
}

// Next returns the Rel corresponding to s's following ply.
func (s Rel) Next() Rel {
	if err, ok := s.err.(checkmateError); ok {
		return Rel{-s.n, err.Next()}
	}
	return Rel{-s.n, s.err}
}

// Window represents the bounds of a position's evaluation.
type Window struct{ alpha, beta Rel }

// Constrain updates the lower bound of w, if applicable, and returns the updated Window
// and a boolean value reporting whether the returned Window remains valid.
// Constrain employs fail-hard beta cutoff. Invariant: alpha <= beta in the returned Window.
func (w Window) Constrain(s Rel) (c Window, ok bool) {
	switch {
	case Less(Score(s), Score(w.alpha)):
		return w, true
	case Less(Score(w.beta), Score(s)):
		return Window{w.beta, w.beta}, false
	default:
		return Window{s, w.beta}, true
	}
}

// Next returns the Window corresponding to w's following ply.
func (w Window) Next() Window { return Window{w.beta.Next(), w.alpha.Next()} }

// String returns a string representation of w.
func (w Window) String() string { return fmt.Sprintf("(%v, %v)", w.alpha, w.beta) }
