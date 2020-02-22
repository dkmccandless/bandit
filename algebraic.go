package main

import "fmt"

// Algebraic returns the description of a Move in standard algebraic notation.
func Algebraic(pos Position, m Move) string {
	var s string
	switch side, ok := m.IsCastle(); {
	case ok && side == QS:
		s = "O-O-O"
	case ok && side == KS:
		s = "O-O"
	default:
		if m.Piece == Pawn {
			if m.IsCapture() {
				s = fileLetters[m.From.File()]
			}
		} else {
			s = pieceLetter[m.Piece]
			var can Board
			for _, mm := range LegalMoves(pos) {
				if mm.To == m.To && mm.Piece == m.Piece {
					can ^= mm.From.Board()
				}
			}
			switch {
			case PopCount(can) == 1:
				// no need to specify
			case PopCount(can&Files[m.From.File()]) == 1:
				s += fileLetters[m.From.File()]
			case PopCount(can&Ranks[m.From.Rank()]) == 1:
				s += rankNumbers[m.From.Rank()]
			default:
				s += fileLetters[m.From.File()]
				s += rankNumbers[m.From.Rank()]
			}
		}
		if m.IsCapture() {
			s += "x"
		}
		s += m.To.String()
		if m.IsPromotion() {
			s += pieceLetter[m.PromotePiece]
		}
	}
	if newpos := Make(pos, m); IsCheck(newpos) {
		if IsMate(newpos) {
			s += "#"
		} else {
			s += "+"
		}
	}
	return s
}

// LongAlgebraic returns the description of a Move in long algebraic notation without check.
func LongAlgebraic(m Move) string {
	switch side, ok := m.IsCastle(); {
	case ok && side == QS:
		return "O-O-O"
	case ok && side == KS:
		return "O-O"
	}
	var s string
	if m.Piece != Pawn {
		s = pieceLetter[m.Piece]
	}
	s += m.From.String()
	if m.IsCapture() {
		s += "x"
	} else {
		s += "-"
	}
	s += m.To.String()
	if m.IsPromotion() {
		s += pieceLetter[m.PromotePiece]
	}
	return s
}

// moveNumber returns the number prefix of a move in algebraic notation.
func moveNumber(pos Position) string {
	s := fmt.Sprintf("%v.", pos.FullMove)
	if pos.ToMove == Black {
		s += ".."
	}
	return s
}

// numberedAlgebraic returns the description of m in standard algebraic notation, prefixed by the move number.
func numberedAlgebraic(pos Position, m Move) string { return moveNumber(pos) + Algebraic(pos, m) }

// Text returns the standard algebraic description of a sequence of Moves starting from pos.
// Each element of moves must be legal in the Position arising from the sequential application to pos of the preceding elements.
func Text(pos Position, moves []Move) string {
	var s string
	if pos.ToMove == Black {
		s = moveNumber(pos)
	}
	for i, m := range moves {
		if pos.ToMove == White {
			s += moveNumber(pos)
		}
		s += Algebraic(pos, m)
		if i < len(moves)-1 {
			s += " "
			pos = Make(pos, m)
		}
	}
	return s
}
