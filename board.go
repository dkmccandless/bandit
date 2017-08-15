package main

import (
	"fmt"
	"strings"
)

// A Square represents a square on the board in little-endian rank-file order (a1 = 0, a2 = 1, h8 = 63).
// Behavior is undefined for values outside of [0, 64).
type Square byte

// Rank returns the rank of the input Square in the range [0, 7).
func (s Square) Rank() byte { return byte(s >> 3) }

// File returns the file of the input Square in the range [0, 7).
func (s Square) File() byte { return byte(s & 7) }

// Diagonal returns the number of the southwest-northeast diagonal on which the input Square lies, from 0 (h1) to 14 (a8).
func (s Square) Diagonal() byte { return 7 + s.Rank() - s.File() }

// AntiDiagonal returns the number of the northwest-southeast anti-diagonal on which the input Square lies, from 0 (a1) to 14 (h8).
func (s Square) AntiDiagonal() byte { return s.Rank() + s.File() }

// Board returns a Board in which only the bit corresponding to the input Square is set.
func (s Square) Board() Board { return 1 << s }

//go:generate stringer -type=Square
const (
	a1 Square = iota
	b1
	c1
	d1
	e1
	f1
	g1
	h1
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a8
	b8
	c8
	d8
	e8
	f8
	g8
	h8
)

// A Board represents a bitboard that describes some aspect of a chess board or position.
// Every bit corresponds to one square on the board.
type Board uint64

const (
	AFile Board = 0x0101010101010101 << iota
	BFile
	CFile
	DFile
	EFile
	FFile
	GFile
	HFile
)

const (
	Rank1 Board = 0xff << (8 * iota)
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

const (
	LongDiagonal     Board = 0x8040201008040201 // a1 to h8
	LongAntiDiagonal Board = 0x0102040810204080 // h1 to a8
)

var (
	Files       = []Board{AFile, BFile, CFile, DFile, EFile, FFile, GFile, HFile}
	Ranks       = []Board{Rank1, Rank2, Rank3, Rank4, Rank5, Rank6, Rank7, Rank8}
	fileLetters = []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	rankNumbers = []string{"1", "2", "3", "4", "5", "6", "7", "8"}

	// The squares that must be empty before castling
	QSCastleEmptySquares = []Board{
		(BFile | CFile | DFile) & Rank1,
		(BFile | CFile | DFile) & Rank8,
	}
	KSCastleEmptySquares = []Board{
		(FFile | GFile) & Rank1,
		(FFile | GFile) & Rank8,
	}
	// The squares that the king occupies during castling
	QSCastleKingSquares = []Board{
		(CFile | DFile | EFile) & Rank1,
		(CFile | DFile | EFile) & Rank8,
	}
	KSCastleKingSquares = []Board{
		(EFile | FFile | GFile) & Rank1,
		(EFile | FFile | GFile) & Rank8,
	}
)

// The two Colors of chess pieces are White and Black. A piece's color determines when it can move and whether it can move to an occupied Square.
type Color byte

//go:generate stringer -type=Color
const (
	White Color = iota
	Black
)

// A Piece is one of the six types of chess pieces, or an auxiliary value corresponding to none or all of them.
type Piece byte

//go:generate stringer -type=Piece
const (
	None Piece = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
	All
)

var pieceLetter = []string{"", "P", "N", "B", "R", "Q", "K"}

// A Position contains all information necessary to specify the current state of a game.
type Position struct {
	// Castling rights: a value of true indicates that the specified Color has the option of castling to the specified side, if it is legal to do so.
	QSCastle [2]bool
	KSCastle [2]bool

	// The unique square, if any, to which an en passant capture can be played by the side to move.
	// Valid values are in the ranges [16, 24) (the 3rd rank, with Black to move following a White pawn push)
	// and [40, 48) (the 6th rank, with White to move following a Black pawn push).
	// A value of 0 indicates that there is no en passant opportunity. Behavior is undefined for any other value.
	ep Square

	// The side to move in the current position, and the opposing side (which had the previous move).
	ToMove, Opp Color

	// The number of half-moves (plies) since the most recent capture or pawn move,
	// for use in determining eligibility for a draw under the fifty-move rule.
	HalfMove int

	// The number of the current move for the side to play.
	// This begins at 1 and increments after each Black move.
	FullMove int

	// The Square of each king.
	KingSquare [2]Square

	// Boards corresponding to the positions of all White and Black pieces,
	// indexed in the order specified by the Piece constants:
	// pawns, knights, bishops, rooks, queens, and the king,
	// with an additional Board representing their union.
	b [2][8]Board
}

var InitialPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// InitialPosition represents the position at the beginning of a game of chess.
var InitialPosition, _ = ParseFEN(InitialPositionFEN)

// PieceOn returns the Color and Piece type of the piece, if any, on the specified Square.
func (pos Position) PieceOn(s Square) (c Color, p Piece, ok bool) {
	b := s.Board()
	switch {
	case pos.b[White][All]&b != 0:
		c = White
	case pos.b[Black][All]&b != 0:
		c = Black
	default:
		return
	}
	for _, p = range []Piece{Pawn, Knight, Bishop, Rook, Queen, King} {
		if pos.b[c][p]&b != 0 {
			ok = true
			return
		}
	}
	// Error if neither return statement above is utilized
	panic(fmt.Sprintf("PieceOn: nonexistent %v piece on square %v in position %+v", c, s, pos))
	return
}

// deBruijn is a sequence of bits each consecutive six-bit subsequence of which corresponds to a different number in [0, 64).
// It is used with dbIndex to provide quick lookup of the position of a single set bit.
const deBruijn = 0x022c98fdaf386e95 // = 0b0000001000101100100110001111110110101111001110000110111010010101

// The positions of the consecutive six-bit numbers appearing as subsequences of the deBruijn sequence,
// in little-endian encoding but counted from the left.
// Left-shifting the deBruijn constant by i bits yields the number dbIndex[i] in the most significant six bits.
var dbIndex = []Square{
	0, 1, 2, 45, 3, 7, 46, 21,
	4, 14, 57, 8, 17, 47, 40, 22,
	62, 5, 55, 15, 60, 58, 9, 33,
	18, 11, 30, 48, 41, 51, 35, 23,
	63, 44, 6, 20, 13, 56, 16, 39,
	61, 54, 59, 32, 10, 29, 50, 34,
	43, 19, 12, 38, 53, 31, 28, 49,
	42, 37, 52, 27, 36, 26, 25, 24,
}

// LS1B returns a Board consisting of only the least significant 1 bit of the input Board.
func LS1B(b Board) Board {
	if b == 0 {
		panic("LS1B: Board is empty")
	}
	return b & -b
}

// LS1BIndex returns the position of the least significant 1 bit of the input Board.
// It left-shifts the deBruijn constant by multiplying it by the least significant 1 bit
// and then uses the most significant six bits to look up the corresponding value in dbIndex.
func LS1BIndex(b Board) Square {
	return dbIndex[(deBruijn*LS1B(b))>>58]
}

// ResetLS1B returns a Board consisting of all set bits of the input Board except for the least significant one.
func ResetLS1B(b Board) Board {
	if b == 0 {
		panic("ResetLS1B: Board is empty")
	}
	return b & (b - 1)
}

// PopCount returns the number of 1 bits in the input Board.
func PopCount(b Board) int {
	var n int
	for b != 0 {
		n++
		b = ResetLS1B(b)
	}
	return n
}

func pieceChar(c Color, p Piece) string {
	s := pieceLetter[p]
	if c == Black {
		s = strings.ToLower(s)
	}
	return s
}

func (pos Position) Display() string {
	s := "\n"
	for r := 7; r >= 0; r-- {
		for f := 0; f < 8; f++ {
			if c, p, ok := pos.PieceOn(Square(8*r + f)); ok {
				s += pieceChar(c, p)
			} else {
				s += "."
			}
			s += " "
		}
		s += "\n"
	}
	return s
}
