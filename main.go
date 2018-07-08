package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	ErrGo = errors.New("go")
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

	for !IsTerminal(pos) {
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
	// Returning the zero value Move{} indicates resignation.
	Play(Position) (int, Move)
}

type Computer struct {
	moveTime time.Duration
	depth    int
}

func (c Computer) Play(pos Position) (int, Move) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.moveTime)
	defer cancel()

	_, results := SearchPosition(ctx, pos, c.depth)
	return results[0].score, results[0].move
}

type Human struct{ s *bufio.Scanner }

func (h Human) Play(pos Position) (int, Move) {
	ch := make(chan Results)
	// Wait for SearchPosition to return
	defer func() { <-ch }()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		_, results := SearchPosition(ctx, pos, 100)
		ch <- results
		close(ch)
	}()

	for {
		m, err := h.readMove(pos)
		if err == ErrGo {
			cancel()
			results := <-ch
			return results[0].score, results[0].move
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		// get engine's opinion of Human's move
		cancel()
		results := <-ch
		for _, r := range results {
			if r.move == m {
				return r.score, m
			}
		}
		return 0, m
	}
}

func (h Human) readMove(pos Position) (m Move, err error) {
	fmt.Printf("> ")
	if ok := h.s.Scan(); !ok {
		return m, h.s.Err()
	}
	if h.s.Text() == "resign" {
		return
	}
	if h.s.Text() == "go" {
		return m, ErrGo
	}
	from, to, err := ParseUserMove(h.s.Text())
	if err != nil {
		return m, err
	}
	c, p := pos.PieceOn(from)
	if p == None {
		return m, fmt.Errorf("No piece on square %v", from)
	}
	if c != pos.ToMove {
		return m, fmt.Errorf("%v piece on square %v", c, from)
	}
	cc, cp := pos.PieceOn(to)
	if cp != None && cc == c {
		return m, fmt.Errorf("%v piece on square %v", cc, to)
	}
	var cs Square
	if cp != None {
		cs = to
	} // TODO: en passant, promotion
	m = Move{From: from, To: to, Piece: p, CapturePiece: cp, CaptureSquare: cs}
	if !IsPseudoLegal(pos, m) || !IsLegal(Make(pos, m)) {
		return m, fmt.Errorf("%v is illegal", m)
	}
	return m, nil
}
