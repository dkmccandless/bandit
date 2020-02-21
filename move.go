package main

// Move contains the information needed to transition from one Position to another.
type Move struct {
	From         Square
	To           Square // invariant: not equal to From
	Piece        Piece  // the moving Piece; invariant: not None
	CapturePiece Piece  // the Piece being captured, or else None
	EP           bool   // whether the Move is an en passant capture
	PromotePiece Piece  // the Piece being promoted to, or else None
}

// IsCapture reports whether m is a capture.
func (m Move) IsCapture() bool { return m.CapturePiece != None }

// IsPromotion reports whether m promotes a pawn.
func (m Move) IsPromotion() bool { return m.PromotePiece != None }

// IsDouble reports whether m is an initial pawn push.
func (m Move) IsDouble() bool {
	return m.Piece == Pawn && ((m.From.Rank() == 1 && m.To.Rank() == 3) || m.From.Rank() == 6 && m.To.Rank() == 4)
}

// IsCastle reports whether m castles and to which side.
// If ok is true, side reports the Side to which m castles.
func (m Move) IsCastle() (side Side, ok bool) {
	switch {
	case m.Piece != King:
		return
	case m.From-m.To == 2:
		return QS, true
	case m.To-m.From == 2:
		return KS, true
	default:
		return
	}
}

// Make applies m to pos and reports the resulting Position.
// Behavior is undefined when m is illegal in pos.
func Make(pos Position, m Move) Position {
	// update changes pos.z and the Boards of c when a Piece of type p moves to or from s
	update := func(c Color, p Piece, s Square) {
		b := s.Board()
		pos.b[c][p] ^= b
		pos.b[c][All] ^= b
		pos.z.xorPiece(c, p, s)
	}

	// Remove en passant capturing rights from the Zobrist bitstring.
	// In the event of an en passant capture, this must be done before the pawn bitboard is changed.
	if pos.ep != 0 {
		if a, b := eligibleEPCapturers(pos); a != 0 {
			pos.z.xor(canEPCaptureZobrist[a.File()])
			if b != 0 {
				pos.z.xor(canEPCaptureZobrist[b.File()])
			}
		}
	}

	// Move the piece
	update(pos.ToMove, m.Piece, m.From)
	update(pos.ToMove, m.Piece, m.To)

	if m.IsCapture() {
		captureSquare := m.To
		if m.EP {
			// The captured pawn is one square centerward of To.
			captureSquare ^= 8
		}
		update(pos.Opp(), m.CapturePiece, captureSquare)
		// Lose the relevant castling right
		switch {
		case (pos.ToMove == Black && m.To == a1) || (pos.ToMove == White && m.To == a8):
			if pos.Castle[pos.Opp()][QS] {
				pos.z.xor(castleZobrist[pos.Opp()][QS])
			}
			pos.Castle[pos.Opp()][QS] = false
		case (pos.ToMove == Black && m.To == h1) || (pos.ToMove == White && m.To == h8):
			if pos.Castle[pos.Opp()][KS] {
				pos.z.xor(castleZobrist[pos.Opp()][KS])
			}
			pos.Castle[pos.Opp()][KS] = false
		}
	}

	switch side, ok := m.IsCastle(); {
	case ok && side == QS:
		// Move the castling rook
		rookFrom, rookTo := m.From-4, m.From-1
		update(pos.ToMove, Rook, rookFrom)
		update(pos.ToMove, Rook, rookTo)
	case ok && side == KS:
		// Move the castling rook
		rookFrom, rookTo := m.From+3, m.From+1
		update(pos.ToMove, Rook, rookFrom)
		update(pos.ToMove, Rook, rookTo)
	}

	if m.IsPromotion() {
		// Replace the Pawn with PromotePiece
		update(pos.ToMove, Pawn, m.To)
		update(pos.ToMove, m.PromotePiece, m.To)
	}

	switch m.Piece {
	case King:
		// Update KingSquare and forfeit all castling rights
		pos.KingSquare[pos.ToMove] = m.To
		if pos.Castle[pos.ToMove][QS] {
			pos.z.xor(castleZobrist[pos.ToMove][QS])
		}
		pos.Castle[pos.ToMove][QS] = false
		if pos.Castle[pos.ToMove][KS] {
			pos.z.xor(castleZobrist[pos.ToMove][KS])
		}
		pos.Castle[pos.ToMove][KS] = false
	case Rook:
		// Forfeit the relevant castling right
		switch {
		case (pos.ToMove == White && m.From == a1) || (pos.ToMove == Black && m.From == a8):
			if pos.Castle[pos.ToMove][QS] {
				pos.z.xor(castleZobrist[pos.ToMove][QS])
			}
			pos.Castle[pos.ToMove][QS] = false
		case (pos.ToMove == White && m.From == h1) || (pos.ToMove == Black && m.From == h8):
			if pos.Castle[pos.ToMove][KS] {
				pos.z.xor(castleZobrist[pos.ToMove][KS])
			}
			pos.Castle[pos.ToMove][KS] = false
		}
	}

	if m.Piece == Pawn || m.IsCapture() {
		pos.HalfMove = 0
	} else {
		pos.HalfMove++
	}
	if pos.ToMove == Black {
		pos.FullMove++
	}
	pos.ToMove = pos.Opp()
	pos.z.xor(blackToMoveZobrist)

	if m.IsDouble() {
		pos.ep = (m.From + m.To) / 2
		// Add en passant capturing rights to the Zobrist bitstring.
		// This must be done after the side to move is changed.
		if a, b := eligibleEPCapturers(pos); a != 0 {
			pos.z.xor(canEPCaptureZobrist[a.File()])
			if b != 0 {
				pos.z.xor(canEPCaptureZobrist[b.File()])
			}
		}
	} else {
		pos.ep = 0
	}

	return pos
}

// PseudoLegalMoves returns all pseudo-legal Moves in pos.
func PseudoLegalMoves(pos Position) []Move {
	moves := make([]Move, 0, 100)
	moves = append(moves, PawnMoves(pos)...)
	moves = append(moves, KnightMoves(pos)...)
	moves = append(moves, BishopMoves(pos)...)
	moves = append(moves, RookMoves(pos)...)
	moves = append(moves, QueenMoves(pos)...)
	moves = append(moves, KingMoves(pos)...)

	// Counting sort into the order winning captures, equal captures, losing captures, non-captures.
	// (This terminology anticipates that the captured piece is defended and the capturing piece is liable to be captured in exchange.)
	const (
		winning = iota
		equal
		losing
		noncapture
	)
	moveType := func(m Move) int {
		switch {
		case !m.IsCapture():
			return noncapture
		case m.Piece == m.CapturePiece ||
			(m.Piece == Bishop && m.CapturePiece == Knight) ||
			(m.Piece == Knight && m.CapturePiece == Bishop):
			return equal
		case m.Piece < m.CapturePiece:
			return winning
		default:
			return losing
		}
	}
	bins := make([]int, 4)
	for _, m := range moves {
		bins[moveType(m)]++
	}
	index := make([]int, len(bins))
	for i := range index {
		for j := 0; j < i; j++ {
			index[i] += bins[j]
		}
	}
	sorted := make([]Move, len(moves))
	for _, m := range moves {
		mt := moveType(m)
		sorted[index[mt]] = m
		index[mt]++
	}
	return sorted
}

// LegalMoves returns all legal Moves in pos.
func LegalMoves(pos Position) []Move {
	pl := PseudoLegalMoves(pos)
	legal := make([]Move, 0, len(pl))
	for _, m := range pl {
		if !IsLegal(Make(pos, m)) {
			continue
		}
		legal = append(legal, m)
	}
	return legal
}

// hasLegalMove reports whether pos has at least one legal move.
func hasLegalMove(pos Position) bool {
	// Generate the pseudolegal moves of each type of piece.
	// If one of them is legal, return immediately.
	for _, pl := range [][]Move{
		PawnMoves(pos),
		KnightMoves(pos),
		BishopMoves(pos),
		RookMoves(pos),
		QueenMoves(pos),
		KingMoves(pos),
	} {
		for _, m := range pl {
			if IsLegal(Make(pos, m)) {
				return true
			}
		}
	}
	return false
}

// rangeBits applies f sequentially to each set bit in board.
func rangeBits(board Board, f func(Board, Square)) {
	for bits := board; bits != 0; bits = ResetLS1B(bits) {
		b, s := LS1B(bits), LS1BIndex(bits)
		f(b, s)
	}
}

// PawnMoves returns a slice of all pseudo-legal Moves that pawns can make in pos.
func PawnMoves(pos Position) []Move {
	moves := make([]Move, 0, 8*2*4) // cap: all pawns are on the 7th rank and can promote via capture to either side
	empty := ^pos.b[White][All] & ^pos.b[Black][All]

	// Pawn movesets are asymmetrical and their capture and non-capture movesets are disjoint
	var pawnAdv, pawnAtk func(Board, Board) Board
	var promoteRank byte
	switch pos.ToMove {
	case White:
		pawnAdv, pawnAtk, promoteRank = whitePawnAdvances, whitePawnAttacks, 7
	case Black:
		pawnAdv, pawnAtk, promoteRank = blackPawnAdvances, blackPawnAttacks, 0
	}
	rangeBits(pos.b[pos.ToMove][Pawn], func(f Board, from Square) {
		rangeBits(pawnAdv(f, empty), func(_ Board, to Square) {
			if to.Rank() == promoteRank {
				for _, pp := range []Piece{Queen, Rook, Bishop, Knight} {
					moves = append(moves, Move{From: from, To: to, Piece: Pawn, PromotePiece: pp})
				}
				return
			}
			moves = append(moves, Move{From: from, To: to, Piece: Pawn})
		})
		rangeBits(pawnAtk(f, empty)&pos.b[pos.Opp()][All], func(_ Board, to Square) {
			_, cp := pos.PieceOn(to)
			if to.Rank() == promoteRank {
				for _, pp := range []Piece{Queen, Rook, Bishop, Knight} {
					moves = append(moves, Move{From: from, To: to, Piece: Pawn, CapturePiece: cp, PromotePiece: pp})
				}
				return
			}
			moves = append(moves, Move{From: from, To: to, Piece: Pawn, CapturePiece: cp})
		})
	})
	if pos.ep != 0 {
		// Double pawn push occurred on the previous move
		epcs := pos.ep ^ 8
		epSources := west(epcs.Board()) | east(epcs.Board())
		rangeBits(epSources&pos.b[pos.ToMove][Pawn], func(_ Board, s Square) {
			moves = append(moves, Move{From: s, To: pos.ep, Piece: Pawn, CapturePiece: Pawn, EP: true})
		})
	}
	return moves
}

// KnightMoves returns a slice of all pseudo-legal Moves that knights can make in pos.
func KnightMoves(pos Position) []Move { return pMoves(pos, Knight) }

// BishopMoves returns a slice of all pseudo-legal Moves that bishops can make in pos.
func BishopMoves(pos Position) []Move { return pMoves(pos, Bishop) }

// RookMoves returns a slice of all pseudo-legal Moves that rooks can make in pos.
func RookMoves(pos Position) []Move { return pMoves(pos, Rook) }

// QueenMoves returns a slice of all pseudo-legal Moves that queens can make in pos.
func QueenMoves(pos Position) []Move { return pMoves(pos, Queen) }

// KingMoves returns a slice of all pseudo-legal Moves that the king can make in pos.
func KingMoves(pos Position) []Move {
	moves := pMoves(pos, King)
	from := pos.KingSquare[pos.ToMove]
	if canCastle(pos, QS) {
		moves = append(moves, Move{From: from, To: from - 2, Piece: King})
	}
	if canCastle(pos, KS) {
		moves = append(moves, Move{From: from, To: from + 2, Piece: King})
	}
	return moves
}

// canCastle returns whether castling to side is legal in pos.
func canCastle(pos Position, side Side) bool {
	if !pos.Castle[pos.ToMove][side] {
		return false
	}
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	if CastleEmptySquares[pos.ToMove][side]&^empty != 0 {
		return false
	}
	var attacked bool
	rangeBits(CastleKingSquares[pos.ToMove][side], func(_ Board, s Square) {
		if attacked {
			return
		}
		attacked = IsAttacked(pos, s, pos.Opp())
	})
	return !attacked
}

// pMoves returns a slice of all pseudo-legal Moves that non-pawn pieces of type p can make in pos, excluding castling.
func pMoves(pos Position, p Piece) []Move {
	moves := make([]Move, 0, 28) // two bishops, two rooks, or one queen can have 28 moves
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	var pAttacks func(Board, Board) Board
	switch p {
	case Knight:
		pAttacks = knightAttacks
	case Bishop:
		pAttacks = bishopAttacks
	case Rook:
		pAttacks = rookAttacks
	case Queen:
		pAttacks = queenAttacks
	case King:
		pAttacks = kingAttacks
	}
	rangeBits(pos.b[pos.ToMove][p], func(f Board, from Square) {
		rangeBits(pAttacks(f, empty)&^pos.b[pos.ToMove][All], func(t Board, to Square) {
			m := Move{From: from, To: to, Piece: p}
			if t&pos.b[pos.Opp()][All] != 0 {
				_, capturePiece := pos.PieceOn(to)
				m.CapturePiece = capturePiece
			}
			moves = append(moves, m)
		})
	})
	return moves
}

var (
	kAttacks = make([]Board, 64)
	nAttacks = make([]Board, 64)
)

func init() {
	for s := a1; s <= h8; s++ {
		b := s.Board()
		// +7 +8 +9
		// -1  K +1
		// -9 -8 -7
		kAttacks[s] = southwest(b) | south(b) | southeast(b) | west(b) | east(b) | northwest(b) | north(b) | northeast(b)
		//    +15  +17
		//  +6        +10
		//        N
		// -10         -6
		//    -17  -15
		nAttacks[s] = southwest(south(b)) | southeast(south(b)) | southwest(west(b)) | southeast(east(b)) | northwest(west(b)) | northeast(east(b)) | northwest(north(b)) | northeast(north(b))
	}
}

// whitePawnAdvances returns a Board of all squares to which a white pawn at p can advance when there are no pieces at empty.
func whitePawnAdvances(p, empty Board) Board {
	return (north(north(p)&empty)&empty)<<32>>32 | (north(p) & empty)
}

// blackPawnAdvances returns a Board of all squares to which a black pawn at p can advance when there are no pieces at empty.
func blackPawnAdvances(p, empty Board) Board {
	return (south(south(p)&empty)&empty)>>32<<32 | (south(p) & empty)
}

// whitePawnAttacks returns a Board of all squares attacked by a white pawn at p.
func whitePawnAttacks(p, _ Board) Board { return northwest(p) | northeast(p) }

// blackPawnAttacks returns a Board of all squares attacked by a black pawn at p.
func blackPawnAttacks(p, _ Board) Board { return southwest(p) | southeast(p) }

// knightAttacks returns a Board of all squares attacked by a knight at p.
func knightAttacks(p, _ Board) Board { return nAttacks[LS1BIndex(p)] }

// bishopAttacks returns a Board of all squares attacked by a bishop at p when there are no pieces at empty.
func bishopAttacks(p, empty Board) Board {
	return attackFill(p, empty, southwest) | attackFill(p, empty, southeast) | attackFill(p, empty, northwest) | attackFill(p, empty, northeast)
}

// rookAttacks returns a Board of all squares attacked by a rook at p when there are no pieces at empty.
func rookAttacks(p, empty Board) Board {
	return attackFill(p, empty, south) | attackFill(p, empty, west) | attackFill(p, empty, east) | attackFill(p, empty, north)
}

// queenAttacks returns a Board of all squares attacked by a queen at p when there are no pieces at empty.
func queenAttacks(p, empty Board) Board { return rookAttacks(p, empty) | bishopAttacks(p, empty) }

// kingAttacks returns a Board of all squares attacked by a king at p.
func kingAttacks(p, _ Board) Board { return kAttacks[LS1BIndex(p)] }

// IsAttacked returns whether s is attacked by any piece of Color c in pos.
func IsAttacked(pos Position, s Square, c Color) bool {
	b := s.Board()
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	switch c {
	case White:
		if blackPawnAttacks(b, empty)&pos.b[c][Pawn] != 0 {
			return true
		}
	case Black:
		if whitePawnAttacks(b, empty)&pos.b[c][Pawn] != 0 {
			return true
		}
	}
	return rookAttacks(b, empty)&(pos.b[c][Rook]|pos.b[c][Queen]) != 0 ||
		bishopAttacks(b, empty)&(pos.b[c][Bishop]|pos.b[c][Queen]) != 0 ||
		knightAttacks(b, empty)&pos.b[c][Knight] != 0 ||
		kingAttacks(b, empty)&pos.b[c][King] != 0
}

// attackFill returns a Board showing all of the squares attacked by the input Board in the direction represented by shift.
func attackFill(piece, empty Board, shift func(Board) Board) Board {
	var fill Board
	for piece != 0 {
		fill, piece = fill|piece, shift(piece)&empty
	}
	return shift(fill) // Include the blocking piece and not the sliding piece
}

func south(b Board) Board { return b >> 8 }
func west(b Board) Board  { return b >> 1 &^ HFile }
func east(b Board) Board  { return b << 1 &^ AFile }
func north(b Board) Board { return b << 8 }

func southwest(b Board) Board { return west(south(b)) }
func southeast(b Board) Board { return east(south(b)) }
func northwest(b Board) Board { return west(north(b)) }
func northeast(b Board) Board { return east(north(b)) }

// eligibleEPCapturers returns the Squares of the pawns, if any, that may pseudo-legally capture en passant in a Position.
// If only one pawn can capture en passant, its Square is the first return value and the second is 0.
func eligibleEPCapturers(pos Position) (Square, Square) {
	var a, b Square
	if pos.ep == 0 {
		return 0, 0
	}
	if pos.ep.File() != 0 {
		westcs := pos.ep ^ 8 - 1
		if c, p := pos.PieceOn(westcs); c == pos.ToMove && p == Pawn {
			a = westcs
		}
	}
	if pos.ep.File() != 7 {
		eastcs := pos.ep ^ 8 + 1
		if c, p := pos.PieceOn(eastcs); c == pos.ToMove && p == Pawn {
			if a != 0 {
				b = eastcs
			} else {
				a = eastcs
			}
		}
	}
	return a, b
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
