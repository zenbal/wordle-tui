package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var ALPHABET = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func inAlphabet(char byte) bool {
	index := sort.Search(len(ALPHABET), func(i int) bool {
		return ALPHABET[i] >= char
	})
	return index < len(ALPHABET) && ALPHABET[index] == char
}

type Wordle struct {
	board    []Guess
	attempt  int
	solution string
	status   GameStatus
	trie     Trie
}

type GameStatus int

const (
	ONGOING GameStatus = iota
	WIN
	LOOSE
)

type Feedback int

const (
	TBD Feedback = iota
	GREY
	YELLOW
	GREEN
)

type Guess []*GuessChar

type GuessChar struct {
	value    byte
	feedback Feedback
}

func NewGuess(word string) (Guess, error) {
	guess := make([]*GuessChar, 5)
	if len(word) != 5 {
		return guess, fmt.Errorf("Error: Guess has to be 5 characters long")
	}
	for i, char := range word {
		if inAlphabet(byte(char)) == false {
			return guess, fmt.Errorf("Error: Invalid character")
		}
		guess[i] = &GuessChar{value: byte(char), feedback: TBD}
	}
	return guess, nil
}

func NewWordle(solution string) *Wordle {
	if len(solution) != 5 {
		return nil
	}
	board := make([]Guess, 6)

	trie := NewTrie()
	if err := trie.insertWordleData(); err != nil {
		return nil
	}

	wordle := &Wordle{
		board:    board,
		attempt:  0,
		solution: solution,
		trie:     trie,
	}
	return wordle
}

func (w *Wordle) guess(word string) error {
	new_guess, err := NewGuess(word)
	if err != nil {
		return err
	}
    if valid := w.trie.findWord(word); valid == false {
        return fmt.Errorf("Error: Invalid word")
    }
	w.board[w.attempt] = new_guess
	num_correct := 0
	for i, char := range w.board[w.attempt] {
		if char.value == w.solution[i] {
			char.feedback = GREEN
			num_correct++
		} else if w.solutionContains(char.value) {
			char.feedback = YELLOW
		} else {
			char.feedback = GREY
		}
	}

	if num_correct == 5 {
		w.status = WIN
	} else if w.attempt == 5 {
		w.status = LOOSE
	} else {
		w.status = ONGOING
	}
	w.attempt++
	return nil
}

func (w *Wordle) solutionContains(char byte) bool {
	for _, c := range w.solution {
		if byte(c) == char {
			return true
		}
	}
	return false
}

// INFO: temporary
func (w *Wordle) play() {
	if w.status != ONGOING {
		fmt.Print(w.toString())
		if w.status == WIN {
			fmt.Println("You win!")
		} else {
			fmt.Println("You loose!")
		}
		return
	}
    reader := bufio.NewReader(os.Stdin)
    for true {
        fmt.Print(w.toString())
        fmt.Print("Enter guess: ")
        text, _ := reader.ReadString('\n')
        text = text[:len(text)-1]
        if err := w.guess(text); err != nil {
            fmt.Println(err)
            continue
        }
        break
    }
	w.play()
}

func (w *Wordle) toString() string {
	s := "Board:\n"
	for row := range w.board {
		if w.board[row] == nil {
			continue
		}
		for col := range w.board[row] {
			s += fmt.Sprintf("  %s(%d)  ", string(w.board[row][col].value), w.board[row][col].feedback)
		}
		s += "\n"
	}
	return s
}

func main() {
	wordle := NewWordle("earth")
	wordle.play()
}
