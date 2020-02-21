package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	// readMove returns errGo in response to the "go" command.
	// errGo instructs Human's Play method to return the engine's preferred move.
	errGo = errors.New("go")

	// readMove returns errResign in response to the "resign" command.
	// errResign instructs Human's Play method to resign the game.
	errResign = errors.New("resign")
)

// Player is the interface that wraps the Play method.
//
// Play analyzes a Position and returns an evaluation score and a Move in that Position.
// Returning the zero value Move{} indicates resignation.
type Player interface {
	Play(Position) (Abs, Move)
}

// Computer can Play without user input.
type Computer struct {
	moveTime time.Duration
	depth    int
}

// Play searches pos and returns an evaluation score and a preferred Move.
func (c Computer) Play(pos Position) (Abs, Move) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.moveTime)
	defer cancel()

	results := SearchPosition(ctx, pos, c.depth)
	fmt.Println(results)
	return results[0].score, results[0].move
}

// Human can Play via user input.
type Human struct{ s *bufio.Scanner }

// Play reads from standard input and returns the engine's evaluation of Pos and the input Move.
// It accepts moves as the concatenation of the origin and destination squares,
// followed by the promoted piece in the case of pawn promotion (e2e4, b1c3, e1g1, d7d8q).
// It also accepts the following commands:
// 	go 		immediately play the engine's preferred move
// 	resign 	resign the game
func (h Human) Play(pos Position) (Abs, Move) {
	ch := make(chan Results)
	// Wait for SearchPosition to return.
	defer func() { <-ch }()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		results := SearchPosition(ctx, pos, 100)
		ch <- results
		close(ch)
	}()

	for {
		m, err := h.readMove(pos)
		switch err {
		case errGo:
			cancel()
			results := <-ch
			return results[0].score, results[0].move
		case errResign:
			cancel()
			return Abs{}, Move{}
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Get the engine's opinion of Human's move.
		cancel()
		results := <-ch
		for _, r := range results {
			if r.move == m {
				return r.score, m
			}
		}
		panic("unreached")
	}
}

// readMove reads from standard input and returns the input Move.
// It accepts moves as the concatenation of origin and destination squares,
// followed by the promoted piece in the case of pawn promotion (e2e4, b1c3, e1g1, d7d8q).
// It also accepts the following commands:
// 	go 		immediately play the engine's preferred move
// 	resign 	resign the game
// It returns an error if the input is invalid or represents an illegal move.
// It also returns ErrGo as the result of the "go" command.
func (h Human) readMove(pos Position) (Move, error) {
	fmt.Printf("> ")
	if ok := h.s.Scan(); !ok {
		return Move{}, h.s.Err()
	}
	return h.parseInput(pos, h.s.Text())
}

// parseInput parses s as a user input command or move in pos.
func (h Human) parseInput(pos Position, s string) (Move, error) {
	var m Move
	var promote Piece
	switch {
	case s == "resign":
		return m, errResign
	case s == "go":
		return m, errGo
	case len(s) == 5:
		switch s[4:] {
		case "q":
			promote = Queen
			s = s[:4]
		case "r":
			promote = Rook
			s = s[:4]
		case "b":
			promote = Bishop
			s = s[:4]
		case "n":
			promote = Knight
			s = s[:4]
		}
	}
	from, to, err := ParseTwoSquares(s)
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
	if promote != None && (p != Pawn || (pos.ToMove == White && to.Rank() != 7) || (pos.ToMove == Black && to.Rank() != 0)) {
		return m, fmt.Errorf("illegal promotion")
	}
	var ep bool
	if pos.ep != 0 && to == pos.ep && cp == Pawn {
		ep = true
	}
	m = Move{From: from, To: to, Piece: p, CapturePiece: cp, EP: ep, PromotePiece: promote}
	if !IsPseudoLegal(pos, m) || !IsLegal(Make(pos, m)) {
		return m, fmt.Errorf("illegal move")
	}
	return m, nil
}

// ParseTwoSquares parses s as the concatenation of two Squares, e.g. "e2e4",
// and returns the corresponding Squares.
func ParseTwoSquares(s string) (from, to Square, err error) {
	if len(s) != 4 {
		return 0, 0, fmt.Errorf("length %v input (want 4)", len(s))
	}
	from, err = ParseSquare(s[:2])
	if err != nil {
		return
	}
	to, err = ParseSquare(s[2:])
	return
}
