package main

import (
	"fmt"
	"sort"
)

var (
	ALPHABET        = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	ALPHABET_LENGTH = 26
	GUESS_LENGTH    = 5
	MAX_GUESSES     = 5 // zero indexed
)

func inAlphabet(char byte) bool {
	index := sort.Search(len(ALPHABET), func(i int) bool {
		return ALPHABET[i] >= char
	})
	return index < len(ALPHABET) && ALPHABET[index] == char
}

type Wordle struct {
	board     []Guess
	attempt   int
	solution  string
	status    GameStatus
	trie      Trie
	guessTrie Trie
	assign    map[int]int          // green
	veto      map[int]map[int]bool // yellow
	include   map[int]bool         // yellow & grey
	message   string
}

type GameStatus int

const (
	ONGOING GameStatus = iota
	WIN
	LOSE
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
	guess := make([]*GuessChar, GUESS_LENGTH)
	if len(word) != GUESS_LENGTH {
		return guess, fmt.Errorf("Error: Guess has to be %d characters long", GUESS_LENGTH)
	}
	for i, char := range word {
		if !inAlphabet(byte(char)) {
			return guess, fmt.Errorf("Error: Invalid character")
		}
		guess[i] = &GuessChar{value: byte(char), feedback: TBD}
	}
	return guess, nil
}

func NewWordle() *Wordle {
	board := make([]Guess, MAX_GUESSES+1)

	trie := NewTrie()
	if err := trie.insertWordleData(wordleSolutionsCSV); err != nil {
		return nil
	}

	guessTrie := NewTrie()
	if err := guessTrie.insertWordleData(wordleGuessesCSV); err != nil {
		return nil
	}
    if err := guessTrie.insertWordleData(wordleSolutionsCSV); err != nil {
        return nil
    }

	veto := make(map[int]map[int]bool, GUESS_LENGTH)
	for i := 0; i < GUESS_LENGTH; i++ {
		veto[i] = make(map[int]bool, ALPHABET_LENGTH)
	}

	wordle := &Wordle{
		board:     board,
		attempt:   0,
		trie:      trie,
		guessTrie: guessTrie,
		assign:    make(map[int]int),  // idx -> char_idx
		include:   make(map[int]bool), // char_idx -> bool
		veto:      veto,               // idx -> char_idx -> bool
	}
	wordle.solution = trie.randomWord()

	return wordle
}

func (w *Wordle) guess(word string) error {
	new_guess, err := NewGuess(word)
	if err != nil {
		return err
	}
	if valid := w.guessTrie.findWord(word); !valid {
		w.message = fmt.Sprintf("'%s' is not a valid word", word)
		return fmt.Errorf("Error: Invalid word")
	}
	w.board[w.attempt] = new_guess
	num_correct := 0
	for i, char := range w.board[w.attempt] {
		char_idx := alphabetIdx(char.value)
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

	if num_correct == GUESS_LENGTH {
		w.status = WIN
	} else if w.attempt == MAX_GUESSES {
		w.status = LOSE
	} else {
		w.status = ONGOING
	}
	w.attempt++
	return nil
}

func (w *Wordle) validate(guess Guess) bool {
	for i, char := range guess {
		char_idx := alphabetIdx(char.value)

		if assigned, ok := w.assign[i]; ok {
			if assigned != char_idx {
				w.message = fmt.Sprintf("'%s' is at index %d of the solution", string(ALPHABET[assigned]), i)
				return false
			}
		}

		if included, ok := w.include[char_idx]; ok {
			if !included {
				w.message = fmt.Sprintf("'%s' is not part of the solution", string(char.value))
				return false
			}
		}

		if veto, ok := w.veto[i]; ok {
			if _, isVetoed := veto[char_idx]; isVetoed {
				w.message = fmt.Sprintf("'%s' can't be at index %d of the solution", string(char.value), i)
				return false
			}
		}
	}

	w.message = ""
	return true
}

func (w *Wordle) validateFull(guess Guess) bool {
	if !w.validate(guess) {
		return false
	}

	found_included := make(map[int]bool)
	for _, char := range guess {
		char_idx := alphabetIdx(char.value)
		if included, ok := w.include[char_idx]; ok && included {
			found_included[char_idx] = true
		}
	}

	for char_idx, include := range w.include {
		if include && !found_included[char_idx] {
			w.message = fmt.Sprintf("'%s' is part of the solution", string(ALPHABET[char_idx]))
			return false
		}
	}

	w.message = ""
	return true
}

func (w *Wordle) suggestNextGuess() string {
	guess := w.findGuessBacktrack()
	word := ""
	for _, char := range guess {
		word += string(char.value)
	}
	return word
}

func (w *Wordle) findGuessBacktrack() Guess {
	if w.attempt == 0 {
		random_guess, err := NewGuess(w.trie.randomWord())
		if err != nil {
			return nil
		}
		return random_guess
	}

	guess := make(Guess, 0)
	return w.backtrack(guess, w.trie.head)
}

func (w *Wordle) backtrack(guess Guess, curr *Node) Guess {
	if len(guess) == GUESS_LENGTH && curr.isWord {
		return guess
	}

	for _, child := range curr.children {
		if child == nil {
			continue
		}
		temp_guess := append(guess, &GuessChar{value: child.value, feedback: TBD})
		if w.validate(temp_guess) {
			result := w.backtrack(temp_guess, child)
			if result != nil && w.validateFull(result) {
				return result
			}
		}
	}

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
