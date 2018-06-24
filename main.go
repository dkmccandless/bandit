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

	tt := time.Now()
	var movesText string
	posZobrists := make(map[Zobrist]int)

	for !IsTerminal(pos) {
		moveTime := time.Now()

		score, results := SearchPosition(pos, *depth)
		move := results[0].move

		if pos.ToMove == White {
			fmt.Printf("\n%v.", pos.FullMove)
			movesText += fmt.Sprintf("%v.", pos.FullMove)
		} else {
			fmt.Printf("\n%v...", pos.FullMove)
		}
		fmt.Printf("%v %.2f %v\n", algebraic(pos, move), float64(score*evalMult(pos.ToMove))/100, time.Since(moveTime))
		movesText += algebraic(pos, move) + " "

		pos = Make(pos, move)
		if pos.z != pos.Zobrist() {
			panic(fmt.Sprintf("pos.z is %x, want %x", pos.z, pos.Zobrist()))
		}

		fmt.Print(pos.String())

		// Check for threefold repetition
		posZobrists[pos.z]++
		if posZobrists[pos.z] == 3 {
			movesText += "1/2-1/2"
			break
		}
	}

	fmt.Printf("\n%v %v\n", movesText, time.Since(tt))
}
