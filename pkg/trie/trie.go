package trie

import (
	"strings"
	"sync"
)

// Node represents a node in the Trie
type Node struct {
	children map[rune]*Node
	isEnd    bool
	value    interface{}
}

// Trie represents a trie data structure
type Trie struct {
	root *Node
	mu   sync.RWMutex
}

// NewTrie creates a new Trie
func NewTrie() *Trie {
	return &Trie{
		root: &Node{children: make(map[rune]*Node)},
	}
}

// Insert adds a key-value pair to the trie
func (t *Trie) Insert(key string, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	for _, ch := range strings.ToLower(key) {
		if node.children[ch] == nil {
			node.children[ch] = &Node{children: make(map[rune]*Node)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
	node.value = value
}

// Search looks for a key in the trie and returns its value
func (t *Trie) Search(key string) (interface{}, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, ch := range strings.ToLower(key) {
		if node.children[ch] == nil {
			return nil, false
		}
		node = node.children[ch]
	}
	return node.value, node.isEnd
}

// Delete removes a key from the trie
func (t *Trie) Delete(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.delete(t.root, strings.ToLower(key), 0)
}

func (t *Trie) delete(node *Node, key string, depth int) bool {
	if node == nil {
		return false
	}

	if depth == len(key) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		node.value = nil
		return len(node.children) == 0
	}

	ch := rune(key[depth])
	if t.delete(node.children[ch], key, depth+1) {
		delete(node.children, ch)
		return !node.isEnd && len(node.children) == 0
	}

	return false
}

// PrefixSearch returns all key-value pairs with the given prefix
func (t *Trie) PrefixSearch(prefix string) map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, ch := range strings.ToLower(prefix) {
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}

	results := make(map[string]interface{})
	t.collect(node, prefix, results)
	return results
}

func (t *Trie) collect(node *Node, prefix string, results map[string]interface{}) {
	if node.isEnd {
		results[prefix] = node.value
	}
	for ch, child := range node.children {
		t.collect(child, prefix+string(ch), results)
	}
}