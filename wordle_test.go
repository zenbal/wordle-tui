package main

import (
	"testing"
)

func NewTestTrie() *Trie {
    trie := NewTrie()
    trie.insertWord("hello")
    trie.insertWord("world")
    return &trie
}

func TestInsertWord(t *testing.T) {
	trie := NewTrie()
	trie.insertWord("hello")
	if trie.head.children[7].children[4].children[11].children[11].children[14].isWord != true {
		t.Errorf("Test failed: Word 'hello' not inserted correctly")
	}
}

func TestInsertTwoWords(t *testing.T) {
	trie := NewTrie()
	trie.insertWord("hello")
	trie.insertWord("world")
	// Check if both words are inserted correctly
	if trie.head.children[7].children[4].children[11].children[11].children[14].isWord != true {
		t.Errorf("Test failed: Word 'hello' not inserted correctly")
	}
	if trie.head.children[22].children[14].children[17].children[11].children[3].isWord != true {
		t.Errorf("Test failed: Word 'world' not inserted correctly")
	}
}

func TestDeleteWord(t *testing.T) {
    trie := NewTestTrie()
    trie.deleteWord("hello")
    if trie.head.children[7] != nil {
        t.Errorf("Test failed: Word 'hello' not deleted correctly")
    }

    trie.deleteWord("world")
    if trie.head.children[22] != nil {
        t.Errorf("Test failed: Word 'world' not deleted correctly")
    }
}
