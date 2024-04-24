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

func TestInsertThreeWords(t *testing.T) {
	trie := NewTrie()
	trie.insertWord("hello")
	trie.insertWord("world")
	trie.insertWord("ha")
	if trie.head.children[7].children[4].children[11].children[11].children[14].isWord != true {
		t.Errorf("Test failed: Word 'hello' not inserted correctly")
	}
	if trie.head.children[22].children[14].children[17].children[11].children[3].isWord != true {
		t.Errorf("Test failed: Word 'world' not inserted correctly")
	}
	if trie.head.children[7].children[0].isWord != true {
		t.Errorf("Test failed: Word 'ha' not inserted correctly")
	}
}

func TestFindWord(t *testing.T) {
	trie := NewTestTrie()
	if !trie.findWord("hello") {
		t.Errorf("Test failed: Expected findWord('hello') to return true but got false")
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

func TestInsertWordleData(t *testing.T) {
	trie := NewTrie()
	if err := trie.insertWordleData(); err != nil {
		t.Errorf("Test failed: Something went wrong")
	}
}
