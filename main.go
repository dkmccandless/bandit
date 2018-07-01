package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var (
		depth = flag.Int("depth", 4, "the search depth")
		fen   = flag.String("fen", InitialPositionFEN, "the FEN record of the starting position")
	)
	flag.Parse()
	if *depth < 1 {
		panic("invalid search depth")
	}

	pos, err := ParseFEN(*fen)
	if err != nil {
		panic(err)
	}

	stdin := bufio.NewScanner(os.Stdin)
	players := []Player{Human{stdin}, Computer{*depth}}

	startTime := time.Now()
	var movesText string
	posZobrists := make(map[Zobrist]int)

	for !IsTerminal(pos) {
		moveTime := time.Now()

		score, move := players[pos.ToMove].Play(pos)

		alg := Algebraic(pos, move)
		movenum := fmt.Sprintf("%v.", pos.FullMove)
		switch pos.ToMove {
		case White:
			movesText += movenum
		case Black:
			movenum += ".."
		}
		movesText += alg + " "

		pos = Make(pos, move)

		fmt.Printf("%v%v %.2f %v\n", movenum, alg, float64(score)/100, time.Since(moveTime).Truncate(time.Millisecond))
		fmt.Println(pos.String())

		// Check for threefold repetition
		posZobrists[pos.z]++
		if posZobrists[pos.z] == 3 {
			movesText += "1/2-1/2"
			break
		}
	}

	fmt.Println(movesText, time.Since(startTime).Truncate(time.Millisecond))
}

type Player interface {
	Play(Position) (int, Move)
}

type Computer struct{ depth int }

func (c Computer) Play(pos Position) (int, Move) {
	_, results := SearchPosition(pos, c.depth)
	return results[0].score, results[0].move
}

type Human struct{ s *bufio.Scanner }

func (h Human) Play(pos Position) (int, Move) {
	fmt.Printf("> ")
	if ok := h.s.Scan(); !ok {
		panic(h.s.Err())
	}
	from, to, err := ParseUserMove(h.s.Text())
	if err != nil {
		panic(err)
	}
	c, p := pos.PieceOn(from)
	if p == None {
		panic(fmt.Sprintf("No piece on square %v", from))
	}
	if c != pos.ToMove {
		panic(fmt.Sprintf("%v piece on square %v", c, from))
	}
	cc, cp := pos.PieceOn(to)
	if cp != None && cc == c {
		panic(fmt.Sprintf("%v piece on square %v", cc, to))
	}
	var cs Square
	if cp != None {
		cs = to
	} // TODO: en passant, promotion
	m := Move{From: from, To: to, Piece: p, CapturePiece: cp, CaptureSquare: cs}
	if !IsPseudoLegal(pos, m) || !IsLegal(Make(pos, m)) {
		panic(fmt.Sprintf("%v is illegal", m))
	}
	return 0, m
}
