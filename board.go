package main

import (
	"fmt"
	"math/bits"
	"strings"
)

// Square represents a square on the board in little-endian rank-file order (a1 = 0, a2 = 1, h8 = 63).
// Behavior is undefined for values outside of [0, 64).
type Square byte

// Rank returns the rank of s in the range [0, 7).
func (s Square) Rank() byte { return byte(s >> 3) }

// File returns the file of s in the range [0, 7).
func (s Square) File() byte { return byte(s & 7) }

// Diagonal returns the southwest-northeast diagonal of s, from 0 (h1) to 14 (a8).
func (s Square) Diagonal() byte { return 7 + s.Rank() - s.File() }

// AntiDiagonal returns the northwest-southeast anti-diagonal of s, from 0 (a1) to 14 (h8).
func (s Square) AntiDiagonal() byte { return s.Rank() + s.File() }

// Board returns a Board in which only the bit corresponding to s is set.
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

// Board represents a bitboard that describes some aspect of a chess board or position.
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
	DarkSquares      Board = 0xaa55aa55aa55aa55
	LightSquares     Board = ^DarkSquares
	LongDiagonal     Board = 0x8040201008040201 // a1 to h8
	LongAntiDiagonal Board = 0x0102040810204080 // h1 to a8
)

var (
	Files       = []Board{AFile, BFile, CFile, DFile, EFile, FFile, GFile, HFile}
	Ranks       = []Board{Rank1, Rank2, Rank3, Rank4, Rank5, Rank6, Rank7, Rank8}
	fileLetters = []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	rankNumbers = []string{"1", "2", "3", "4", "5", "6", "7", "8"}

	// CastleEmptySquares describes the squares that must be empty when castling.
	CastleEmptySquares = [][]Board{
		{
			(BFile | CFile | DFile) & Rank1,
			(FFile | GFile) & Rank1,
		},
		{
			(BFile | CFile | DFile) & Rank8,
			(FFile | GFile) & Rank8,
		},
	}

	// CastleKingSquares describes the squares that the king occupies during castling.
	CastleKingSquares = [][]Board{
		{
			(CFile | DFile | EFile) & Rank1,
			(EFile | FFile | GFile) & Rank1,
		},
		{
			(CFile | DFile | EFile) & Rank8,
			(EFile | FFile | GFile) & Rank8,
		},
	}
)

// Color represents the color of a chess piece. A piece's color determines when it can move and whether it can move to an occupied Square.
type Color byte

//go:generate stringer -type=Color
const (
	White Color = iota
	Black
)

// Piece is one of the six types of chess pieces, or an auxiliary value corresponding to none or all of them.
type Piece byte

//go:generate stringer -type=Piece
const (
	// None is used in Move's CapturePiece and PromotePiece fields to denote that a move is not a capture or a promotion.
	None Piece = iota

	Pawn
	Knight
	Bishop
	Rook
	Queen
	King

	// All is updated in a Position's b elements to track each side's pieces regardless of type.
	// This is an optimization for the move generation functions.
	All
)

var pieceLetter = []string{"", "P", "N", "B", "R", "Q", "K"}

// Side represents the two sides of the board to which the king can castle.
type Side byte

const (
	QS Side = iota
	KS
)

// Position contains all information necessary to specify the current state of a game.
type Position struct {
	// Castle describes castling rights.
	// Castle[c][side] indicates whether the c retains the option of castling to side, if it is legal to do so.
	Castle [2][2]bool

	// ep is the unique square, if any, to which an en passant capture can be played by the side to move.
	// Valid values are in the ranges [16, 24) (the 3rd rank, with Black to move following a White pawn push)
	// and [40, 48) (the 6th rank, with White to move following a Black pawn push).
	// A value of 0 indicates that there is no en passant opportunity. Behavior is undefined for any other value.
	ep Square

	// ToMove is the side to move.
	ToMove Color

	// HalfMove is the number of half-moves (plies) since the most recent capture or pawn move,
	// for use in determining eligibility for a draw under the fifty-move rule.
	HalfMove int

	// FullMove is the number of the current move for the side to play.
	// This begins at 1 and increments after each Black move.
	FullMove int

	// KingSquare is the Square of the indexed Color's king.
	KingSquare [2]Square

	// b contains Boards describing the positions of all White and Black pieces.
	// The set bits of b[c][p] give the locations of all pieces of color c and type p.
	// b[c][All] gives the union of all pieces of color c.
	b [2][8]Board

	// z is the position's Zobrist bitstring.
	z Zobrist
}

// InitialPositionFEN is the FEN record of the initial position of the pieces.
var InitialPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// InitialPosition is the initial Position.
var InitialPosition, _ = ParseFEN(InitialPositionFEN)

// Opp returns the Color of the player who does not have the move.
func (pos Position) Opp() Color { return pos.ToMove ^ 1 }

// PieceOn returns the Color and Piece type of the piece, if any, on s.
// It returns a Piece of None if the Square is empty.
func (pos Position) PieceOn(s Square) (c Color, p Piece) {
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
			return
		}
	}
	// Error if neither return statement above is utilized
	panic(fmt.Sprintf("PieceOn: invalid piece on square %v in position %+v", s, pos))
	return
}

// LS1B returns a Board consisting of only the least significant 1 bit of b.
func LS1B(b Board) Board {
	if b == 0 {
		panic("LS1B: Board is empty")
	}
	return b & -b
}

// LS1BIndex returns the position of the least significant 1 bit of b.
func LS1BIndex(b Board) Square {
	return Square(bits.TrailingZeros64(uint64(b)))
}

// ResetLS1B returns a Board consisting of all set bits of b except for the least significant one.
func ResetLS1B(b Board) Board {
	if b == 0 {
		panic("ResetLS1B: Board is empty")
	}
	return b & (b - 1)
}

// PopCount returns the number of 1 bits in b.
func PopCount(b Board) int {
	return bits.OnesCount64(uint64(b))
}

func pieceChar(c Color, p Piece) string {
	s := pieceLetter[p]
	if c == Black {
		s = strings.ToLower(s)
	}
	return s
}

func (pos Position) String() string {
	var s string
	for r := 7; r >= 0; r-- {
		for f := 0; f < 8; f++ {
			if c, p := pos.PieceOn(Square(8*r + f)); p != None {
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
