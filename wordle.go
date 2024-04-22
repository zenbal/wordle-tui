package main

// import "fmt"

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
    siblings := append(n.parent.children[:node_idx], n.parent.children[node_idx + 1:]...)
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
            new_node.parent = curr 
			curr.children[char_idx] = new_node
		}
		curr = new_node
	}
	curr.isWord = true
}

func (t *Trie) deleteWord(word string) {
	curr := t.head
	for i := 0; i < len(word); i++ {
		char := word[i]
		char_idx := alphabet_idx(char)
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

