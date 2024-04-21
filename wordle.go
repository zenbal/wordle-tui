package main

// import "fmt"

func alphabet_idx(char byte) int {
	return int(char - 'a')
}

type Node struct {
	value    byte
	children []*Node
	isWord   bool
}

func NewNode(value byte) *Node {
	return &Node{
		value:    value,
		children: make([]*Node, 26),
	}
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
	for i := 0; i < len(word); i++ {
		char := word[i]
		char_idx := alphabet_idx(char)
		var new_node *Node = nil
		if next := curr.children[char_idx]; next == nil {
			new_node = NewNode(word[i])
			curr.children[char_idx] = new_node
		}
		curr = new_node
	}
	curr.isWord = true
}
