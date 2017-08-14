package main

import (
	"fmt"
	"time"
)

func main() {
	pos := InitialPosition
	moves := []Move{}
	var movesText string

	for !IsTerminal(pos) {
		moveTime := time.Now()

		score, move := SearchPosition(pos, 4)

		if pos.ToMove == White {
			fmt.Printf("\n%v.", pos.FullMove)
			movesText += fmt.Sprintf("%v.", pos.FullMove)
		} else {
			fmt.Printf("\n%v...", pos.FullMove)
		}
		fmt.Printf("%v %.2f %v", algebraic(pos, move), float64(score)/100, time.Since(moveTime))
		movesText += algebraic(pos, move) + " "

		moves = append(moves, move)
		pos = Make(pos, move)

		fmt.Print(pos.Display())
	}

	fmt.Printf("\n%v", movesText)
}
