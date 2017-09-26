package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	RuneToColor = map[rune]Color{
		'P': White,
		'N': White,
		'B': White,
		'R': White,
		'Q': White,
		'K': White,
		'p': Black,
		'n': Black,
		'b': Black,
		'r': Black,
		'q': Black,
		'k': Black,
	}
	RuneToPiece = map[rune]Piece{
		'P': Pawn,
		'p': Pawn,
		'N': Knight,
		'n': Knight,
		'B': Bishop,
		'b': Bishop,
		'R': Rook,
		'r': Rook,
		'Q': Queen,
		'q': Queen,
		'K': King,
		'k': King,
	}
)

// ParseFEN converts an FEN record to a Position.
func ParseFEN(fen string) (pos Position, err error) {
	fields := strings.Fields(fen)
	if len(fields) != 6 {
		return pos, fmt.Errorf("ParseFEN: %v fields (need 6)", len(fields))
	}
	if wk := strings.Count(fields[0], "K"); wk != 1 {
		return pos, fmt.Errorf("ParseFEN: %v white kings", wk)
	}
	if bk := strings.Count(fields[0], "k"); bk != 1 {
		return pos, fmt.Errorf("ParseFEN: %v black kings", bk)
	}
	rows := strings.Split(fields[0], "/")
	if len(rows) != 8 {
		return pos, fmt.Errorf("ParseFEN: %v rows (need 8)", len(rows))
	}
	for r, s := range rows {
		var n int
		for _, char := range s {
			switch {
			case isWhite(char) || isBlack(char):
				n++
			case isNumber(char):
				n += int(char - '0')
			default:
				return pos, fmt.Errorf("ParseFEN: Invalid character in row %v", s)
			}
		}
		if n != 8 {
			return pos, fmt.Errorf("ParseFEN: %v squares in row %v (need 8)", n, s)
		}
		sq := Square(56 - 8*r)
		for _, char := range s {
			if isNumber(char) {
				sq += Square(char - '0')
				continue
			}
			c, ok := RuneToColor[char]
			if !ok {
				return pos, fmt.Errorf("ParseFEN: Invalid character in row %v", s)
			}
			p, ok := RuneToPiece[char]
			if !ok {
				return pos, fmt.Errorf("ParseFEN: Invalid character in row %v", s)
			}
			if p == Pawn && (sq.Rank() == 0 || sq.Rank() == 7) {
				return pos, fmt.Errorf("ParseFEN: Pawn on invalid square %v", sq.String())
			}
			pos.b[c][p] ^= sq.Board()
			pos.b[c][All] ^= sq.Board()
			if p == King {
				pos.KingSquare[c] = sq
			}
			sq++
		}
	}

	// side to move
	switch fields[1] {
	case "w":
		pos.ToMove = White
	case "b":
		pos.ToMove = Black
	default:
		return pos, fmt.Errorf("ParseFEN: Invalid active player field %v", fields[1])
	}

	// castling
	if n := strings.IndexFunc(fields[2], func(char rune) bool {
		return char != '-' && RuneToPiece[char] != Queen && RuneToPiece[char] != King
	}); n != -1 {
		return pos, fmt.Errorf("ParseFEN: Invalid character in castling field %v", fields[2])
	}
	for _, char := range fields[2] {
		switch RuneToPiece[char] {
		case Queen:
			pos.QSCastle[RuneToColor[char]] = true
		case King:
			pos.KSCastle[RuneToColor[char]] = true
		}
	}

	// en passant
	if fields[3] != "-" {
		if pos.ep, err = ParseSquare(fields[3]); err != nil {
			return
		}
		if pos.ep.Rank() != 2 && pos.ep.Rank() != 5 {
			return pos, fmt.Errorf("ParseFEN: Invalid en passant square %v", pos.ep.String())
		}
	}

	// halfmove clock
	if pos.HalfMove, err = strconv.Atoi(fields[4]); err != nil {
		return
	}

	// fullmove number
	if pos.FullMove, err = strconv.Atoi(fields[5]); err != nil {
		return
	}

	// Zobrist bitstring
	pos.z = pos.Zobrist()

	return
}

// ParseSquare returns the Square represented by the string.
// It returns an error if the string does not consist of one letter in [a, h] followed by one number in [1, 8].
func ParseSquare(s string) (sq Square, err error) {
	if len(s) != 2 || !isFile(rune(s[0])) || !isNumber(rune(s[1])) {
		return 0, fmt.Errorf("ParseSquare(%v): Invalid character", s)
	}
	return Square(8*(s[1]-'1') + s[0] - 'a'), nil
}

// FEN converts a Position into an FEN record.
func FEN(pos Position) string {
	var s string
	for rank := 7; rank >= 0; rank-- {
		var gap int
		for sq := Square(8 * rank); sq < Square(8*(rank+1)); sq++ {
			c, p := pos.PieceOn(sq)
			if p != None {
				if gap != 0 {
					s += strconv.Itoa(gap)
					gap = 0
				}
				s += pieceChar(c, p)
			} else {
				gap++
			}
		}
		if gap != 0 {
			s += strconv.Itoa(gap)
		}
		if rank != 0 {
			s += "/"
		}
	}
	switch pos.ToMove {
	case White:
		s += " w "
	case Black:
		s += " b "
	}
	if pos.KSCastle[White] {
		s += "K"
	}
	if pos.QSCastle[White] {
		s += "Q"
	}
	if pos.KSCastle[Black] {
		s += "k"
	}
	if pos.QSCastle[Black] {
		s += "q"
	}
	if !(pos.KSCastle[White] || pos.QSCastle[White] || pos.KSCastle[Black] || pos.QSCastle[Black]) {
		s += "-"
	}
	if pos.ep == 0 {
		s += " - "
	} else {
		s += " " + pos.ep.String() + " "
	}
	s += strconv.Itoa(pos.HalfMove) + " " + strconv.Itoa(pos.FullMove)
	return s
}

func isWhite(r rune) bool  { return r == 'P' || r == 'N' || r == 'B' || r == 'R' || r == 'Q' || r == 'K' }
func isBlack(r rune) bool  { return r == 'p' || r == 'n' || r == 'b' || r == 'r' || r == 'q' || r == 'k' }
func isNumber(r rune) bool { return '1' <= r && r <= '8' }
func isFile(r rune) bool   { return 'a' <= r && r <= 'h' }
