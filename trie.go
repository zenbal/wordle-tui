package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"math/rand"
)

func alphabet_idx(char byte) int {
	return int(char - 'a')
}

type Node struct {
	value    byte
	children []*Node
	parent   *Node
	isWord   bool
}

func NewNode(value byte) *Node {
	return &Node{
		value:    value,
		children: make([]*Node, 26),
	}
}

func (n *Node) siblings() []*Node {
	node_idx := alphabet_idx(n.value)
	siblings := append(n.parent.children[:node_idx], n.parent.children[node_idx+1:]...)
	return siblings
}

func (n *Node) anySiblings() bool {
	siblings := n.siblings()
	for _, sibling := range siblings {
		if sibling != nil {
			return true
		}
	}
	return false
}

func (n *Node) getChildren() []*Node {
	result := make([]*Node, 0)
	for _, child := range n.children {
		if child != nil {
			result = append(result, child)
		}
	}
	return result
}

type Trie struct {
	head *Node
}

func NewTrie() Trie {
	head := NewNode(' ')
	return Trie{
		head: head,
	}
}

func (t *Trie) insertWord(word string) {
	curr := t.head
	for _, char := range word {
		char_idx := alphabet_idx(byte(char))
		var new_node *Node = nil
		if next := curr.children[char_idx]; next == nil {
			new_node = NewNode(byte(char))
			new_node.parent = curr
			curr.children[char_idx] = new_node
			curr = new_node
		} else {
			curr = next
		}
	}
	curr.isWord = true
}

func (t *Trie) findWord(word string) bool {
	curr := t.head
	for _, char := range word {
		char_idx := alphabet_idx(byte(char))
		if next := curr.children[char_idx]; next == nil {
			return false
		} else {
			curr = next
		}
	}
	return curr.isWord
}

func (t *Trie) deleteWord(word string) {
	curr := t.head
	for _, char := range word {
		char_idx := alphabet_idx(byte(char))
		if next := curr.children[char_idx]; next != nil {
			curr = next
		} else {
			return
		}
	}
	curr_idx := alphabet_idx(curr.value)
	for i := 0; i < len(word); i++ {
		if curr.anySiblings() {
			curr.parent.children[curr_idx] = nil
			return
		}
		curr = curr.parent
		curr_idx = alphabet_idx(curr.value)
	}
}

func (t *Trie) randomWord() string {
	curr := t.head
	word := ""
	for i := 0; i < 5; i++ {
		children := curr.getChildren()
		next := children[rand.Intn(len(children))]
		word += string(next.value)
		curr = next
	}
	return word
}

//go:embed valid_wordle_solutions.csv
var wordleDataCSV []byte

func (t *Trie) insertWordleData() error {
	reader := csv.NewReader(bytes.NewReader(wordleDataCSV))
	words, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, word := range words {
		t.insertWord(word[0])
	}
	return nil
}
