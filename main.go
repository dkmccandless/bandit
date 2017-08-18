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

	var movesText string

	for !IsTerminal(pos) {
		moveTime := time.Now()

		score, move := SearchPosition(pos, *depth)

		if pos.ToMove == White {
			fmt.Printf("\n%v.", pos.FullMove)
			movesText += fmt.Sprintf("%v.", pos.FullMove)
		} else {
			fmt.Printf("\n%v...", pos.FullMove)
		}
		fmt.Printf("%v %.2f %v", algebraic(pos, move), float64(score)/100, time.Since(moveTime))
		movesText += algebraic(pos, move) + " "

		pos = Make(pos, move)
		pos.z = pos.Zobrist()

		fmt.Print(pos.Display())
	}

	fmt.Printf("\n%v", movesText)
}
