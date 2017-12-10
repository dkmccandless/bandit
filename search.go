package main

const (
	evalInf = 50000
)

func SearchPosition(pos Position, depth int) (int, Move) {
	var score int
	var recommended Move
	for d := 1; d <= depth; d++ {
		score, recommended = negamax(pos, recommended, -evalInf, evalInf, d, true, negamax)
	}
	return score, recommended
}

type SearchFunc func(Position, Move, int, int, int, bool, SearchFunc) (int, Move)

func negamax(pos Position, recommended Move, alpha int, beta int, depth int, allowCutoff bool, search SearchFunc) (bestScore int, bestMove Move) {
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

		score, _ := search(newpos, Move{}, -beta, -alpha, depth-1, allowCutoff, search)
		score *= -1

		if score > alpha {
			alpha, bestMove = score, m
		}
		if alpha > beta && allowCutoff {
			break
		}
	}
	return alpha, bestMove
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
