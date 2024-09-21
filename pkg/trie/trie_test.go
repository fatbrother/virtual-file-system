package trie

import (
	"reflect"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()

	// Test Insert and Search
	trie.Insert("hello", 1)
	trie.Insert("world", 2)
	trie.Insert("hi", 3)

	tests := []struct {
		key      string
		wantVal  interface{}
		wantBool bool
	}{
		{"hello", 1, true},              // key exists
		{"world", 2, true},			     // key exists
		{"hi", 3, true},		     	 // key exists
		{"hel", nil, false},		 	 // key does not exist
		{"hello world", nil, false},	 // key does not exist
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			gotVal, gotBool := trie.Search(tt.key)
			if gotVal != tt.wantVal || gotBool != tt.wantBool {
				t.Errorf("Search() = (%v, %v), want (%v, %v)", gotVal, gotBool, tt.wantVal, tt.wantBool)
			}
		})
	}

	// Test Delete
	trie.Delete("hello")
	if val, found := trie.Search("hello"); found {
		t.Errorf("Delete() failed, key 'hello' still exists with value %v", val)
	}

	// Test PrefixSearch
	trie.Insert("help", 4)
	trie.Insert("helm", 5)

	prefixResults := trie.PrefixSearch("he")
	expectedResults := map[string]interface{}{
		"help": 4,
		"helm": 5,
	}

	if !reflect.DeepEqual(prefixResults, expectedResults) {
		t.Errorf("PrefixSearch() = %v, want %v", prefixResults, expectedResults)
	}
}