package main

import (
	"testing"
)

func NewTestWordle() *Wordle {
	wordle := NewWordle()
	wordle.solution = "earth"
	return wordle
}

func TestGuess(t *testing.T) {
	wordle := NewTestWordle()
	guess := "adept"
	err := wordle.guess(guess)
	if err != nil {
		t.Errorf("Expected guess to be successful but got %s", err)
	}

	for i, char := range guess {
		if byte(char) != wordle.board[0][i].value {
			t.Errorf("Expected 'adept' to be in row 0 of the wordle board.")
		}
	}

	if wordle.board[0][0].feedback != YELLOW {
		t.Errorf("Expected feedback %d for '%s' at index %d", YELLOW, string(wordle.board[0][0].value), 0)
	}
	if wordle.board[0][1].feedback != GREY {
		t.Errorf("Expected feedback %d for '%s' at index %d", GREY, string(wordle.board[0][1].value), 1)
	}
	if wordle.board[0][2].feedback != YELLOW {
		t.Errorf("Expected feedback %d for '%s' at index %d", YELLOW, string(wordle.board[0][2].value), 2)
	}
	if wordle.board[0][3].feedback != GREY {
		t.Errorf("Expected feedback %d for '%s' at index %d", GREY, string(wordle.board[0][3].value), 3)
	}
	if wordle.board[0][4].feedback != YELLOW {
		t.Errorf("Expected feedback %d for '%s' at index %d", YELLOW, string(wordle.board[0][4].value), 4)
	}
}

func TestValidate(t *testing.T) {
	wordle := NewTestWordle()
	err := wordle.guess("adept")
	if err != nil {
		t.Errorf("Expected guess to be successful but got %s", err)
	}

	guess := make([]*GuessChar, 5)
	word := "taste"
	for i := range word {
		guess[i] = &GuessChar{value: byte(word[i]), feedback: TBD}
	}

	// valid
	if valid := wordle.validate(guess); !valid {
		t.Errorf("Expected 'taste' to be a valid guess following 'adept'")
	}

	// veto
	word = "adult"
	for i := range word {
		guess[i] = &GuessChar{value: byte(word[i]), feedback: TBD}
	}
	if valid := wordle.validate(guess); valid {
		t.Errorf("Expected 'adult' to be an invalid guess following 'adept'")
	}

	// include/exclude
	word = "drown"
	for i := range word {
		guess[i] = &GuessChar{value: byte(word[i]), feedback: TBD}
	}
	if valid := wordle.validate(guess); valid {
		t.Errorf("Expected 'drown' to be an invalid guess following 'adept'")
	}
}

func TestFindGuessBacktrack(t *testing.T) {
	wordle := NewTestWordle()
	err := wordle.guess("adept")
	if err != nil {
		t.Errorf("Expected guess to be successful but got %s", err)
	}

	guess := wordle.findGuessBacktrack()
	guess_str := ""
	for _, char := range guess {
		guess_str += string(char.value)
	}
	if guess_str != "baste" {
		t.Errorf("Expected suggested guess to be 'baste' but got '%s'", guess_str)
	}

	err = wordle.guess(guess_str)
	if err != nil {
		t.Errorf("Expected guess to be successful but got %s", err)
	}

	guess = wordle.findGuessBacktrack()
	guess_str = ""
	for _, char := range guess {
		guess_str += string(char.value)
	}
	if guess_str != "earth" {
		t.Errorf("Expected suggested guess to be 'earth' but got '%s'", guess_str)
	}

	err = wordle.guess(guess_str)
	if err != nil {
		t.Errorf("Expected guess to be successful but got %s", err)
	}

	if wordle.status != WIN {
		t.Errorf("Expected status to be 'WIN'")
	}
}
