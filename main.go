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

	stdin := bufio.NewScanner(os.Stdin)
	players := []Player{Computer{*moveTime, *depth}, Computer{*moveTime, *depth}}
	if *humanWhite {
		players[White] = Human{stdin}
	}
	if *humanBlack {
		players[Black] = Human{stdin}
	}

	startTime := time.Now()
	var movesText string
	posZobrists := make(map[Zobrist]int)

game:
	for {
		moveTime := time.Now()

		score, move := players[pos.ToMove].Play(pos)
		if move == (Move{}) {
			// player resigns
			switch pos.ToMove {
			case White:
				movesText += "0-1"
			case Black:
				movesText += "1-0"
			}
			break
		}

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
		if s, ok := score.err.(checkmateError); ok {
			score = Abs{err: s.Next()}
		}

		fmt.Printf("%v%v %v %v\n", movenum, alg, score, time.Since(moveTime).Truncate(time.Millisecond))
		fmt.Println(pos)

		// Check for end-of-game conditions
		switch score.err {
		case errCheckmate:
			if pos.ToMove == White {
				movesText += "0-1"
			} else {
				movesText += "1-0"
			}
			break game
		case errStalemate, errInsufficient, errFiftyMove:
			movesText += "1/2-1/2"
			break game
		}
		if posZobrists[pos.z]++; posZobrists[pos.z] == 3 {
			// threefold repetition
			movesText += "1/2-1/2"
			break
		}
	}

	fmt.Println(movesText)
	fmt.Println(time.Since(startTime).Truncate(time.Millisecond))
}
