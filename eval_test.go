package main

import "testing"

func TestInsufficientMaterial(t *testing.T) {
	for _, fen := range []string{
		"K6k/8/8/8/8/8/8/8 w - - 0 1",
		"KN5k/8/8/8/8/8/8/8 w - - 0 1",
		"KB5k/8/8/8/8/8/8/8 w - - 0 1",
		"KB1b3k/8/8/8/8/8/8/8 w - - 0 1",
		"K6k/8/8/8/8/8/8/1B1B1B1B w - - 0 1",
		"K1b1b1bk/8/8/8/8/8/8/1B1B1B1B w - - 0 1",
	} {
		pos, err := ParseFEN(fen)
		if err != nil {
			panic(err)
		}
		if got := Eval(pos); got != 0 {
			t.Errorf("Insufficient Material: Eval(%v)==%v, want 0", fen, got)
		}
	}
}