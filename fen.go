package main

import (
	"errors"
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
		err = fmt.Errorf("ParseFEN: %v fields (need 6)", len(fields))
	}
	switch wk, bk := strings.Count(fields[0], "K"), strings.Count(fields[0], "k"); {
	case wk == 0:
		err = errors.New("ParseFEN: No white king")
		return
	case wk > 1:
		err = fmt.Errorf("ParseFEN: %v white kings", wk)
		return
	case bk == 0:
		err = errors.New("ParseFEN: No black king")
		return
	case bk > 1:
		err = fmt.Errorf("ParseFEN: %v black kings", bk)
		return
	}
	rows := strings.Split(fields[0], "/")
	if len(rows) != 8 {
		err = fmt.Errorf("ParseFEN: %v rows (need 8)", len(rows))
	}
	for r, s := range rows {
		var n int
		for _, char := range s {
			switch {
			case isWhite(char) || isBlack(char):
				n++
			case isNumber(char):
				n += int(char - '1')
			default:
				err = fmt.Errorf("ParseFEN: Invalid character in row %v", s)
				return
			}
		}
		if n != 8 {
			err = fmt.Errorf("ParseFEN: %v squares in row %v (need 8)", n, s)
		}
		sq := Square(56 - 8*r)
		for _, char := range s {
			if isNumber(char) {
				sq += Square(char - '0')
				continue
			}
			c, ok := RuneToColor[char]
			if !ok {
				err = fmt.Errorf("ParseFEN: Invalid character in row %v", s)
				return
			}
			p, ok := RuneToPiece[char]
			if !ok {
				err = fmt.Errorf("ParseFEN: Invalid character in row %v", s)
				return
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
		pos.ToMove, pos.Opp = White, Black
	case "b":
		pos.ToMove, pos.Opp = Black, White
	default:
		err = fmt.Errorf("ParseFEN: Invalid active player field %v", fields[1])
		return
	}

	// castling
	if n := strings.IndexFunc(fields[2], func(char rune) bool {
		return char != '-' && RuneToPiece[char] != Queen && RuneToPiece[char] != King
	}); n != -1 {
		err = fmt.Errorf("ParseFEN: Invalid character in castling field %v", fields[2])
		return
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
	if n := strings.IndexFunc(fields[3], func(char rune) bool {
		return char != '-' && !('a' <= char && char <= 'h') && char != '3' && char != '6'
	}); n != -1 {
		err = fmt.Errorf("ParseFEN: Invalid character in en passant field %v", fields[2])
		return
	}
	if fields[3] != "-" {
		pos.ep += Square(8*(fields[3][1]-'1') + fields[3][0] - 'a')
	}

	// halfmove clock
	if pos.HalfMove, err = strconv.Atoi(fields[4]); err != nil {
		return
	}

	// fullmove number
	if pos.FullMove, err = strconv.Atoi(fields[5]); err != nil {
		return
	}

	return
}

// FEN converts a Position into an FEN record.
func FEN(pos Position) string {
	var s string
	for rank := 7; rank >= 0; rank-- {
		var gap int
		for sq := Square(8 * rank); sq < Square(8*(rank+1)); sq++ {
			c, p, ok := pos.PieceOn(sq)
			if ok {
				if gap != 0 {
					s += strconv.Itoa(gap)
					gap = 0
				}
				s += PieceChar(c, p)
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

func PieceChar(c Color, p Piece) string {
	var s string
	switch p {
	case Pawn:
		s = "P"
	case Knight:
		s = "N"
	case Bishop:
		s = "B"
	case Rook:
		s = "R"
	case Queen:
		s = "Q"
	case King:
		s = "K"
	}
	if c == Black {
		s = strings.ToLower(s)
	}
	return s
}

func isWhite(r rune) bool  { return RuneToColor[r] == White }
func isBlack(r rune) bool  { return RuneToColor[r] == Black }
func isNumber(r rune) bool { return '1' <= r && r <= '8' }
