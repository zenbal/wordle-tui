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
	assign   map[int]int          // green
	veto     map[int]map[int]bool // yellow
	include  map[int]bool         // yellow & grey
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

	veto := make(map[int]map[int]bool, 5)
	for i := 0; i < 5; i++ {
		veto[i] = make(map[int]bool, 26)
	}

	wordle := &Wordle{
		board:    board,
		attempt:  0,
		solution: solution,
		trie:     trie,
		assign:   make(map[int]int),  // idx -> char_idx
		include:  make(map[int]bool), // char_idx -> bool
		veto:     veto,               // idx -> char_idx -> bool
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
		char_idx := alphabet_idx(char.value)
		if char.value == w.solution[i] {
			char.feedback = GREEN
			w.assign[i] = char_idx
			num_correct++
		} else if w.solutionContains(char.value) {
			char.feedback = YELLOW
			w.include[char_idx] = true
			w.veto[i][char_idx] = true
		} else {
			char.feedback = GREY
			w.include[char_idx] = false
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

func (w *Wordle) validate(guess Guess) bool {
    found_included := make(map[int]bool)
	for i, char := range guess {
        char_idx := alphabet_idx(char.value)

        if assigned, ok := w.assign[i]; ok {
            if assigned != char_idx {
                fmt.Printf("Tip: '%s' is at index %d of the solution\n", string(ALPHABET[assigned]), i)
                return false
            }
        }

        if included, ok := w.include[char_idx]; ok {
            if !included {
                fmt.Printf("Tip: '%s' is not part of the solution\n", string(char.value))
                return false
            } 
            found_included[char_idx] = true
        }
        
        if veto, ok := w.veto[i]; ok {
            if _, isVetoed := veto[char_idx]; isVetoed {
                fmt.Printf("Tip: '%s' can't be at index %d of the solution\n", string(char.value), i)
                return false
            }
        }
	}
    
    for char_idx, include := range w.include {
        if include && !found_included[char_idx] {
            fmt.Printf("Tip: '%s' is part of the solution\n", string(ALPHABET[char_idx]))
            return false
        }
    }

	return true
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
		// fmt.Println("Assign:", w.assign)
		// fmt.Println("Include:", w.include)
		// fmt.Println("Veto:", w.veto)
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
