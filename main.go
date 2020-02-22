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
		defaultTime = 3 * time.Second
		moveTime    = flag.Duration("time", 0, fmt.Sprintf("computer's time per move (default %v if no depth limit set)", defaultTime))
		depth       = flag.Int("depth", 0, "the search depth")
		fen         = flag.String("fen", InitialPositionFEN, "the FEN record of the starting position")
		humanWhite  = flag.Bool("w", false, "user plays White")
		humanBlack  = flag.Bool("b", false, "user plays Black")
	)
	flag.Parse()
	if *moveTime <= 0 {
		*moveTime = 86164091 * time.Millisecond
		if *depth <= 0 {
			*moveTime = defaultTime
		}
	}
	if *depth <= 0 {
		*depth = 100
	}

	pos, err := ParseFEN(*fen)
	if err != nil {
		panic(err)
	}
	startpos := pos

	stdin := bufio.NewScanner(os.Stdin)
	players := []Player{Computer{*moveTime, *depth}, Computer{*moveTime, *depth}}
	if *humanWhite {
		players[White] = Human{stdin}
	}
	if *humanBlack {
		players[Black] = Human{stdin}
	}

	fmt.Println(pos)
	startTime := time.Now()
	var moves []Move
	var resultText string
	posZobrists := make(map[Zobrist]int)

game:
	for {
		moveTime := time.Now()

		score, move := players[pos.ToMove].Play(pos)
		if move == (Move{}) {
			// player resigns
			resultText = []string{"1-0", "0-1"}[pos.Opp()]
			break
		}

		numalg := numberedAlgebraic(pos, move) // before Make
		moves = append(moves, move)
		pos = Make(pos, move)
		if s, ok := score.err.(checkmateError); ok {
			score = Abs{err: s.Next()}
		}

		fmt.Printf("%v %v %v\n", numalg, score, time.Since(moveTime).Truncate(time.Millisecond))
		fmt.Println(pos)

		// Check for end-of-game conditions
		switch score.err {
		case errCheckmate:
			resultText = []string{"1-0", "0-1"}[pos.Opp()]
			break game
		case errStalemate, errInsufficient, errFiftyMove:
			resultText = "1/2-1/2"
			break game
		}
		if posZobrists[pos.z]++; posZobrists[pos.z] == 3 {
			// threefold repetition
			resultText = "1/2-1/2"
			break
		}
	}

	fmt.Printf("%v %v\n", Text(startpos, moves), resultText)
	fmt.Println(time.Since(startTime).Truncate(time.Millisecond))
}
