package main

const (
	evalInf = 50000
)

func SearchPosition(pos Position, depth int) (int, Move) {
	return negamax(pos, -evalInf, evalInf, depth, true, negamax)
}

type SearchFunc func(Position, int, int, int, bool, SearchFunc) (int, Move)

func negamax(pos Position, alpha int, beta int, depth int, allowCutoff bool, search SearchFunc) (bestScore int, bestMove Move) {
	if depth == 0 {
		return Eval(pos) * evalMult(pos.ToMove), bestMove
	}
	moves := Candidates(pos) // pseudo-legal

	bestScore = alpha
	var haveLegalMove bool

	for _, m := range moves {
		newpos := Make(pos, m)
		if !isLegal(newpos) {
			continue
		}
		haveLegalMove = true

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

	if !haveLegalMove { // checkmate or stalemate
		if isCheck(pos) {
			bestScore = -evalInf
		} else {
			bestScore = 0
		}
	}

	return
}

// isCheck returns whether the side to move's king is in check.
func isCheck(pos Position) bool {
	return IsAttacked(pos, pos.KingSquare[pos.ToMove], pos.Opp)
}

// isLegal returns whether the given Position results from a legal move.
// A move is illegal if it leaves one's own king in check.
func isLegal(pos Position) bool {
	return !IsAttacked(pos, pos.KingSquare[pos.Opp], pos.ToMove)
}
