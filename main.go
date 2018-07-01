package main

import (
	"flag"
	"fmt"
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

	players := []Player{Computer{*depth}, Computer{*depth}}

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
