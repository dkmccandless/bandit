package main

import (
	"log"
)

// A Move contains the information needed to transition from one Position to another.
// Behavior is undefined when the Move is not legal in the given Position.
type Move struct {
	From          Square
	To            Square
	Piece         Piece  // the moving Piece; invariant: not None
	CapturePiece  Piece  // the Piece being captured, or else None
	CaptureSquare Square // differs from To in the case of en passant captures
	PromotePiece  Piece  // the Piece being promoted to, or else None
}

// IsCapture returns whether the Move is a capture.
func (m Move) IsCapture() bool { return m.CapturePiece != None }

// IsPromotion returns whether the Move promotes a pawn.
func (m Move) IsPromotion() bool { return m.PromotePiece != None }

// IsDouble returns whether the Move is an initial pawn push.
func (m Move) IsDouble() bool {
	return m.Piece == Pawn && ((m.From.Rank() == 1 && m.To.Rank() == 3) || m.From.Rank() == 6 && m.To.Rank() == 4)
}

// IsQSCastle returns whether the Move castles queenside.
func (m Move) IsQSCastle() bool { return m.Piece == King && m.From-m.To == 2 }

// IsQSCastle returns whether the Move castles kingside.
func (m Move) IsKSCastle() bool { return m.Piece == King && m.To-m.From == 2 }

// Make applies a Move to a Position and returns the resulting Position.
func Make(pos Position, m Move) Position {
	// Remove en passant capturing rights from the Zobrist bitstring.
	// In the event of an en passant capture, this must be done before the pawn bitboard is changed.
	if pos.ep != 0 {
		for _, sq := range eligibleEPCapturers(pos) {
			pos.z.xor(canEPCaptureZobrist[sq.File()])
		}
	}

	// Move the piece
	fromTo := m.From.Board() ^ m.To.Board()
	pos.b[pos.ToMove][m.Piece] ^= fromTo
	pos.b[pos.ToMove][All] ^= fromTo
	pos.z.xorPiece(pos.ToMove, m.Piece, m.From)
	pos.z.xorPiece(pos.ToMove, m.Piece, m.To)

	switch {
	case m.IsCapture():
		// Remove the captured piece from CaptureSquare, not To
		pos.b[pos.Opp()][m.CapturePiece] ^= m.CaptureSquare.Board()
		pos.b[pos.Opp()][All] ^= m.CaptureSquare.Board()
		pos.z.xorPiece(pos.Opp(), m.CapturePiece, m.CaptureSquare)
		// Lose the relevant castling right
		switch {
		case (pos.ToMove == Black && m.CaptureSquare == a1) || (pos.ToMove == White && m.CaptureSquare == a8):
			if pos.QSCastle[pos.Opp()] {
				pos.z.xor(qsCastleZobrist[pos.Opp()])
			}
			pos.QSCastle[pos.Opp()] = false
		case (pos.ToMove == Black && m.CaptureSquare == h1) || (pos.ToMove == White && m.CaptureSquare == h8):
			if pos.KSCastle[pos.Opp()] {
				pos.z.xor(ksCastleZobrist[pos.Opp()])
			}
			pos.KSCastle[pos.Opp()] = false
		}
	case m.IsQSCastle():
		// Move the castling rook
		rookFromSquare, rookToSquare := m.From-4, m.From-1
		rookFromTo := rookFromSquare.Board() ^ rookToSquare.Board()
		pos.b[pos.ToMove][Rook] ^= rookFromTo
		pos.b[pos.ToMove][All] ^= rookFromTo
		pos.z.xorPiece(pos.ToMove, Rook, rookFromSquare)
		pos.z.xorPiece(pos.ToMove, Rook, rookToSquare)
	case m.IsKSCastle():
		// Move the castling rook
		rookFromSquare, rookToSquare := m.From+3, m.From+1
		rookFromTo := rookFromSquare.Board() ^ rookToSquare.Board()
		pos.b[pos.ToMove][Rook] ^= rookFromTo
		pos.b[pos.ToMove][All] ^= rookFromTo
		pos.z.xorPiece(pos.ToMove, Rook, rookFromSquare)
		pos.z.xorPiece(pos.ToMove, Rook, rookToSquare)
	}

	if m.IsPromotion() {
		// Replace the Pawn with PromotePiece
		pos.b[pos.ToMove][Pawn] ^= m.To.Board()
		pos.b[pos.ToMove][m.PromotePiece] ^= m.To.Board()
		pos.z.xorPiece(pos.ToMove, Pawn, m.To)
		pos.z.xorPiece(pos.ToMove, m.PromotePiece, m.To)
	}

	switch m.Piece {
	case King:
		// Update KingSquare and forfeit all castling rights
		pos.KingSquare[pos.ToMove] = m.To
		if pos.QSCastle[pos.ToMove] {
			pos.z.xor(qsCastleZobrist[pos.ToMove])
		}
		pos.QSCastle[pos.ToMove] = false
		if pos.KSCastle[pos.ToMove] {
			pos.z.xor(ksCastleZobrist[pos.ToMove])
		}
		pos.KSCastle[pos.ToMove] = false
	case Rook:
		// Forfeit the relevant castling right
		switch {
		case (pos.ToMove == White && m.From == a1) || (pos.ToMove == Black && m.From == a8):
			if pos.QSCastle[pos.ToMove] {
				pos.z.xor(qsCastleZobrist[pos.ToMove])
			}
			pos.QSCastle[pos.ToMove] = false
		case (pos.ToMove == White && m.From == h1) || (pos.ToMove == Black && m.From == h8):
			if pos.KSCastle[pos.ToMove] {
				pos.z.xor(ksCastleZobrist[pos.ToMove])
			}
			pos.KSCastle[pos.ToMove] = false
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
		for _, sq := range eligibleEPCapturers(pos) {
			pos.z.xor(canEPCaptureZobrist[sq.File()])
		}
	} else {
		pos.ep = 0
	}

	return pos
}

// Candidates returns a slice of all pseudo-legal Moves in the current Position.
func Candidates(pos Position) []Move {
	can := []Move{}
	can = append(can, PawnMoves(pos)...)
	can = append(can, KnightMoves(pos)...)
	can = append(can, BishopMoves(pos)...)
	can = append(can, RookMoves(pos)...)
	can = append(can, QueenMoves(pos)...)
	can = append(can, KingMoves(pos)...)

	// Counting sort into the order winning captures, equal captures, losing captures, non-captures.
	// (This terminology anticipates that the captured piece is defended and the capturing piece is liable to be captured in exchange.)
	const (
		winning = iota
		equal
		losing
		noncapture
	)
	bins := make([]int, 4)
	for _, m := range can {
		var moveType int
		switch {
		case !m.IsCapture():
			moveType = noncapture
		case m.Piece == m.CapturePiece || (m.Piece == Bishop && m.CapturePiece == Knight) || (m.Piece == Knight && m.CapturePiece == Bishop):
			moveType = equal
		case m.Piece < m.CapturePiece:
			moveType = winning
		case m.Piece > m.CapturePiece:
			moveType = losing
		}
		bins[moveType]++
	}
	index := make([]int, len(bins))
	for i := range index {
		for j := 0; j < i; j++ {
			index[i] += bins[j]
		}
	}
	sorted := make([]Move, len(can))
	for _, m := range can {
		var moveType int
		switch {
		case !m.IsCapture():
			moveType = noncapture
		case m.Piece == m.CapturePiece || (m.Piece == Bishop && m.CapturePiece == Knight) || (m.Piece == Knight && m.CapturePiece == Bishop):
			moveType = equal
		case m.Piece < m.CapturePiece:
			moveType = winning
		case m.Piece > m.CapturePiece:
			moveType = losing
		}
		sorted[index[moveType]] = m
		index[moveType]++
	}
	return sorted
}

var (
	kingAttacks   = make([]Board, 64)
	knightAttacks = make([]Board, 64)
)

func init() {
	for s := Square(0); s < 64; s++ {
		b := s.Board()
		// +7 +8 +9
		// -1  K +1
		// -9 -8 -7
		kingAttacks[s] = southwest(b) | south(b) | southeast(b) | west(b) | east(b) | northwest(b) | north(b) | northeast(b)
		//    +15  +17
		//  +6        +10
		//        N
		// -10         -6
		//    -17  -15
		knightAttacks[s] = southwest(south(b)) | southeast(south(b)) | southwest(west(b)) | southeast(east(b)) | northwest(west(b)) | northeast(east(b)) | northwest(north(b)) | northeast(north(b))
	}
}

var PieceMoves = []func(Position) []Move{
	func(Position) []Move { return []Move{} },
	PawnMoves,
	KnightMoves,
	BishopMoves,
	RookMoves,
	QueenMoves,
	KingMoves,
}

// PawnMoves returns a slice of all pseudo-legal Moves that pawns can make in the current Position.
func PawnMoves(pos Position) (moves []Move) {
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	for pawns := pos.b[pos.ToMove][Pawn]; pawns != 0; pawns = ResetLS1B(pawns) {
		f, from := LS1B(pawns), LS1BIndex(pawns)

		switch pos.ToMove {
		case White:
			// Pawns have disjoint capture and non-capture movesets
			for dst := whitePawnAdvances(f, empty); dst != 0; dst = ResetLS1B(dst) {
				to := LS1BIndex(dst)
				m := Move{From: from, To: to, Piece: Pawn}
				if to.Rank() == 7 {
					m.PromotePiece = Queen
					moves = append(moves, m)
					m.PromotePiece = Rook
					moves = append(moves, m)
					m.PromotePiece = Bishop
					moves = append(moves, m)
					m.PromotePiece = Knight
				}
				moves = append(moves, m)
			}
			for dst := whitePawnAttacks(f, empty) & pos.b[pos.Opp()][All]; dst != 0; dst = ResetLS1B(dst) {
				to := LS1BIndex(dst)
				captureColor, capturePiece, ok := pos.PieceOn(to)
				if !ok {
					log.Fatalf("PawnMoves (White): attempted capture on empty Square %v", to)
				}
				if captureColor != pos.Opp() {
					log.Fatalf("PawnMoves (White): attempted capture of %v %v on %v", captureColor, capturePiece, to)
				}
				m := Move{From: from, To: to, Piece: Pawn, CapturePiece: capturePiece, CaptureSquare: to}
				if to.Rank() == 7 {
					m.PromotePiece = Queen
					moves = append(moves, m)
					m.PromotePiece = Rook
					moves = append(moves, m)
					m.PromotePiece = Bishop
					moves = append(moves, m)
					m.PromotePiece = Knight
				}
				moves = append(moves, m)
			}
		case Black:
			// Pawns have disjoint capture and non-capture movesets
			for dst := blackPawnAdvances(f, empty); dst != 0; dst = ResetLS1B(dst) {
				to := LS1BIndex(dst)
				m := Move{From: from, To: to, Piece: Pawn}
				if to.Rank() == 0 {
					m.PromotePiece = Queen
					moves = append(moves, m)
					m.PromotePiece = Rook
					moves = append(moves, m)
					m.PromotePiece = Bishop
					moves = append(moves, m)
					m.PromotePiece = Knight
				}
				moves = append(moves, m)
			}
			for dst := blackPawnAttacks(f, empty) & pos.b[pos.Opp()][All]; dst != 0; dst = ResetLS1B(dst) {
				to := LS1BIndex(dst)
				captureColor, capturePiece, ok := pos.PieceOn(to)
				if !ok {
					log.Fatalf("PawnMoves (Black): attempted capture on empty Square %v", to)
				}
				if captureColor != pos.Opp() {
					log.Fatalf("PawnMoves (Black): attempted capture of %v %v on %v", captureColor, capturePiece, to)
				}
				m := Move{From: from, To: to, Piece: Pawn, CapturePiece: capturePiece, CaptureSquare: to}
				if to.Rank() == 0 {
					m.PromotePiece = Queen
					moves = append(moves, m)
					m.PromotePiece = Rook
					moves = append(moves, m)
					m.PromotePiece = Bishop
					moves = append(moves, m)
					m.PromotePiece = Knight
				}
				moves = append(moves, m)
			}
		}
	}

	if pos.ep != 0 {
		// Double pawn push occurred on the previous move
		epCaptureSquare := pos.ep ^ 8
		epSources := west(epCaptureSquare.Board()) | east(epCaptureSquare.Board())
		for src := epSources & pos.b[pos.ToMove][Pawn]; src != 0; src = ResetLS1B(src) {
			from := LS1BIndex(src)
			moves = append(moves, Move{
				From:          from,
				To:            pos.ep,
				Piece:         Pawn,
				CapturePiece:  Pawn,
				CaptureSquare: epCaptureSquare,
			})
		}
	}
	return
}

// KnightMoves returns a slice of all pseudo-legal Moves that knights can make in the current Position.
func KnightMoves(pos Position) (moves []Move) {
	for knights := pos.b[pos.ToMove][Knight]; knights != 0; knights = ResetLS1B(knights) {
		from := LS1BIndex(knights)
		for dst := knightAttacks[from] &^ pos.b[pos.ToMove][All]; dst != 0; dst = ResetLS1B(dst) {
			t, to := LS1B(dst), LS1BIndex(dst)
			m := Move{From: from, To: to, Piece: Knight}
			if t&pos.b[pos.Opp()][All] != 0 {
				captureColor, capturePiece, ok := pos.PieceOn(to)
				if !ok {
					log.Fatalf("KnightMoves: attempted capture on empty Square %v", to)
				}
				if captureColor != pos.Opp() {
					log.Fatalf("KnightMoves: attempted capture of %v %v on %v", captureColor, capturePiece, to)
				}
				m.CapturePiece = capturePiece
				m.CaptureSquare = to
			}
			moves = append(moves, m)
		}
	}
	return
}

// BishopMoves returns a slice of all pseudo-legal Moves that bishops can make in the current Position.
func BishopMoves(pos Position) (moves []Move) {
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	for bishops := pos.b[pos.ToMove][Bishop]; bishops != 0; bishops = ResetLS1B(bishops) {
		f, from := LS1B(bishops), LS1BIndex(bishops)
		for dst := bishopAttacks(f, empty) &^ pos.b[pos.ToMove][All]; dst != 0; dst = ResetLS1B(dst) {
			t, to := LS1B(dst), LS1BIndex(dst)
			m := Move{From: from, To: to, Piece: Bishop}
			if t&pos.b[pos.Opp()][All] != 0 {
				captureColor, capturePiece, ok := pos.PieceOn(to)
				if !ok {
					log.Fatalf("BishopMoves: attempted capture on empty Square %v", to)
				}
				if captureColor != pos.Opp() {
					log.Fatalf("BishopMoves: attempted capture of %v %v on %v", captureColor, capturePiece, to)
				}
				m.CapturePiece = capturePiece
				m.CaptureSquare = to
			}
			moves = append(moves, m)
		}
	}
	return
}

// RookMoves returns a slice of all pseudo-legal Moves that rooks can make in the current Position.
func RookMoves(pos Position) (moves []Move) {
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	for rooks := pos.b[pos.ToMove][Rook]; rooks != 0; rooks = ResetLS1B(rooks) {
		f, from := LS1B(rooks), LS1BIndex(rooks)
		for dst := rookAttacks(f, empty) &^ pos.b[pos.ToMove][All]; dst != 0; dst = ResetLS1B(dst) {
			t, to := LS1B(dst), LS1BIndex(dst)
			m := Move{From: from, To: to, Piece: Rook}
			if t&pos.b[pos.Opp()][All] != 0 {
				captureColor, capturePiece, ok := pos.PieceOn(to)
				if !ok {
					log.Fatalf("RookMoves: attempted capture on empty Square %v", to)
				}
				if captureColor != pos.Opp() {
					log.Fatalf("RookMoves: attempted capture of %v %v on %v", captureColor, capturePiece, to)
				}
				m.CapturePiece = capturePiece
				m.CaptureSquare = to
			}
			moves = append(moves, m)
		}
	}
	return
}

// QueenMoves returns a slice of all pseudo-legal Moves that queens can make in the current Position.
func QueenMoves(pos Position) (moves []Move) {
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	for queens := pos.b[pos.ToMove][Queen]; queens != 0; queens = ResetLS1B(queens) {
		f, from := LS1B(queens), LS1BIndex(queens)
		for dst := queenAttacks(f, empty) &^ pos.b[pos.ToMove][All]; dst != 0; dst = ResetLS1B(dst) {
			t, to := LS1B(dst), LS1BIndex(dst)
			m := Move{From: from, To: to, Piece: Queen}
			if t&pos.b[pos.Opp()][All] != 0 {
				captureColor, capturePiece, ok := pos.PieceOn(to)
				if !ok {
					log.Fatalf("QueenMoves: attempted capture on empty Square %v", to)
				}
				if captureColor != pos.Opp() {
					log.Fatalf("QueenMoves: attempted capture of %v %v on %v", captureColor, capturePiece, to)
				}
				m.CapturePiece = capturePiece
				m.CaptureSquare = to
			}
			moves = append(moves, m)
		}
	}
	return
}

// KingMoves returns a slice of all pseudo-legal Moves that the king can make in the current Position.
func KingMoves(pos Position) (moves []Move) {
	from := LS1BIndex(pos.b[pos.ToMove][King])
	for dst := kingAttacks[from] &^ pos.b[pos.ToMove][All]; dst != 0; dst = ResetLS1B(dst) {
		t, to := LS1B(dst), LS1BIndex(dst)
		m := Move{From: from, To: to, Piece: King}
		if t&pos.b[pos.Opp()][All] != 0 {
			captureColor, capturePiece, ok := pos.PieceOn(to)
			if !ok {
				log.Fatalf("KingMoves: attempted capture on empty Square %v", to)
			}
			if captureColor != pos.Opp() {
				log.Fatalf("KingMoves: attempted capture of %v %v on %v", captureColor, capturePiece, to)
			}
			m.CapturePiece = capturePiece
			m.CaptureSquare = to
		}
		moves = append(moves, m)
	}

	if canQSCastle(pos) {
		moves = append(moves, Move{From: from, To: pos.KingSquare[pos.ToMove] - 2, Piece: King})
	}
	if canKSCastle(pos) {
		moves = append(moves, Move{From: from, To: pos.KingSquare[pos.ToMove] + 2, Piece: King})
	}
	return
}

// canQSCastle returns whether castling queenside is pseudo-legal in the current Position.
func canQSCastle(pos Position) bool {
	if !pos.QSCastle[pos.ToMove] {
		return false
	}
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	if QSCastleEmptySquares[pos.ToMove]&^empty != 0 {
		return false
	}
	for dst := QSCastleKingSquares[pos.ToMove]; dst != 0; dst = ResetLS1B(dst) {
		if IsAttacked(pos, LS1BIndex(dst), pos.Opp()) {
			return false
		}
	}
	return true
}

// canKSCastle returns whether castling kingside is pseudo-legal in the current Position.
func canKSCastle(pos Position) bool {
	if !pos.KSCastle[pos.ToMove] {
		return false
	}
	empty := ^pos.b[White][All] & ^pos.b[Black][All]
	if KSCastleEmptySquares[pos.ToMove]&^empty != 0 {
		return false
	}
	for dst := KSCastleKingSquares[pos.ToMove]; dst != 0; dst = ResetLS1B(dst) {
		if IsAttacked(pos, LS1BIndex(dst), pos.Opp()) {
			return false
		}
	}
	return true
}

// whitePawnAdvances returns a Board consisting of all squares to which the input white pawn can advance.
func whitePawnAdvances(pawn, empty Board) Board {
	return (north(north(pawn)&empty)&empty)<<32>>32 | (north(pawn) & empty)
}

// blackPawnAdvances returns a Board consisting of all squares to which the input black pawn can advance.
func blackPawnAdvances(pawn, empty Board) Board {
	return (south(south(pawn)&empty)&empty)>>32<<32 | (south(pawn) & empty)
}

// whitePawnAttacks returns a Board consisting of all squares attacked/defended by the input white pawn.
func whitePawnAttacks(pawn, empty Board) Board {
	return (northwest(pawn) | northeast(pawn)) &^ empty
}

// blackPawnAttacks returns a Board consisting of all squares attacked/defended by the input black pawn.
func blackPawnAttacks(pawn, empty Board) Board {
	return (southwest(pawn) | southeast(pawn)) &^ empty
}

// bishopAttacks returns a Board consisting of all squares attacked/defended by the input bishop.
func bishopAttacks(piece, empty Board) Board {
	return attackFill(piece, empty, southwest) | attackFill(piece, empty, southeast) | attackFill(piece, empty, northwest) | attackFill(piece, empty, northeast)
}

// rookAttacks returns a Board consisting of all squares attacked/defended by the input rook.
func rookAttacks(piece, empty Board) Board {
	return attackFill(piece, empty, south) | attackFill(piece, empty, west) | attackFill(piece, empty, east) | attackFill(piece, empty, north)
}

// queenAttacks returns a Board consisting of all squares attacked/defended by the input queen.
func queenAttacks(piece, empty Board) Board {
	return rookAttacks(piece, empty) | bishopAttacks(piece, empty)
}

// IsAttacked returns whether the given Square is attacked by any piece of the given Color in the given Position.
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
		knightAttacks[s]&pos.b[c][Knight] != 0 ||
		kingAttacks[s]&pos.b[c][King] != 0
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

// eligibleEPCapturers returns a slice of the Squares of all pawns that may pseudo-legally capture en passant in a Position.
func eligibleEPCapturers(pos Position) []Square {
	var s []Square
	if pos.ep != 0 && pos.ep.File() != 0 {
		westCaptureSquare := pos.ep ^ 8 - 1
		if c, p, _ := pos.PieceOn(westCaptureSquare); c == pos.ToMove && p == Pawn {
			s = append(s, westCaptureSquare)
		}
	}
	if pos.ep != 0 && pos.ep.File() != 7 {
		eastCaptureSquare := pos.ep ^ 8 + 1
		if c, p, _ := pos.PieceOn(eastCaptureSquare); c == pos.ToMove && p == Pawn {
			s = append(s, eastCaptureSquare)
		}
	}
	return s
}

// algebraic returns the description of a Move in algebraic notation.
func algebraic(pos Position, m Move) string {
	var s string
	switch m.Piece {
	case Pawn:
		if m.IsCapture() {
			s = fileLetters[m.From.File()]
		}
	default:
		s = pieceLetter[m.Piece]
		var can Board
		for _, mm := range PieceMoves[m.Piece](pos) {
			if mm.To == m.To {
				can ^= mm.From.Board()
			}
		}
		switch {
		case PopCount(can) == 1:
			// don't need to specify
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
	if newpos := Make(pos, m); IsCheck(newpos) {
		if IsTerminal(newpos) {
			s += "#"
		} else {
			s += "+"
		}
	}
	return s
}
