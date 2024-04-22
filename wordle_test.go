package main

import (
	"fmt"
	"testing"
)

func NewTestWordle() *Wordle {
    return NewWordle("earth")
}

func TestGuess(t *testing.T) {
    wordle := NewTestWordle()
    wordle.guess("drake")
    fmt.Print(wordle.toString())
}
