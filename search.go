package main

const (
	evalInf = 50000
)

func SearchPosition(pos Position, depth int) (int, Move) {
	return negamax(pos, -evalInf, evalInf, depth, true, negamax)
}

type SearchFunc func(Position, int, int, int, bool, SearchFunc) (int, Move)

func negamax(pos Position, alpha int, beta int, depth int, allowCutoff bool, search SearchFunc) (bestScore int, bestMove Move) {
	if IsTerminal(pos) { // checkmate or stalemate
		if IsCheck(pos) {
			bestScore = -evalInf
		}
		return
	}
	if depth == 0 {
		bestScore = Eval(pos) * evalMult(pos.ToMove)
		return
	}

	moves := Candidates(pos) // pseudo-legal
	bestScore = alpha
	// Initialize bestMove with a legal move
	for i, m := range moves {
		if IsLegal(Make(pos, m)) {
			bestMove = m
			moves = moves[i:]
			break
		}
	}

	for _, m := range moves {
		newpos := Make(pos, m)
		if !IsLegal(newpos) {
			continue
		}

		score, _ := search(newpos, -beta, -alpha, depth-1, allowCutoff, search)
		score *= -1

		if score > bestScore {
			bestScore, bestMove = score, m
		}
		if score > alpha {
			alpha = score
		}
		if alpha > beta && allowCutoff {
			break
		}
	}
	return
}

// IsCheck returns whether the king of the side to move is in check.
func IsCheck(pos Position) bool {
	return IsAttacked(pos, pos.KingSquare[pos.ToMove], pos.Opp())
}

// IsLegal returns whether a Position results from a legal move.
// A position is illegal if the king of the side that just moved is in check.
func IsLegal(pos Position) bool {
	return !IsAttacked(pos, pos.KingSquare[pos.Opp()], pos.ToMove)
}

// IsTerminal returns whether or not a Position is checkmate or stalemate.
// A position is checkmate or stalemate if the side to move has no legal moves.
func IsTerminal(pos Position) bool {
	for _, m := range Candidates(pos) {
		if IsLegal(Make(pos, m)) {
			return false
		}
	}
	return true
}
