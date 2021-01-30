package trie

import (
	"fmt"
	"github.com/drakos74/lachesis/store"
)

// Trie is a trie structure with keys and values made of byte arrays
// TODO : review, check implementation and flavours of Trie(s) based on the theory
type Trie struct {
	key   byte
	value []byte
	tries map[byte]Trie
}

// NewTrie creates a new Trie
func NewTrie(b byte) *Trie {
	return &Trie{
		key:   b,
		value: make([]byte, 0),
		tries: make(map[byte]Trie),
	}
}

// Commit adds the corresponding key-value pair to the Trie
func (t *Trie) Commit(key []byte, value []byte) error {

	if len(key) == 1 {
		return t.add(key[:1], value)
	}

	b := key[0]
	trie, ok := t.tries[b]

	if ok {
		// we already have this node ...
		return trie.Commit(key[1:], value)
	}
	// we dont have the rest of these nodes in the trie
	return t.add(key, value)
}

// Read reads the value for the corresponding key
func (t *Trie) Read(key []byte) ([]byte, bool) {

	b := key[0]
	trie, ok := t.tries[b]
	if ok {
		if len(key) > 1 {
			return trie.Read(key[1:])
		} else if len(trie.value) > 0 {
			return trie.value, true
		}
	}
	return nil, false
}

// with will directly override the value of the trie node.
func (t *Trie) with(value []byte) Trie {
	t.value = value
	return *t
}

// add will add the value to the first and only child of the trie
func (t *Trie) add(key []byte, value []byte) error {

	b := key[0]
	trie := NewTrie(b)
	if len(key) > 1 {
		t.tries[b] = *trie
		return trie.add(key[1:], value)
	}
	t.tries[b] = trie.with(value)

	return nil
}

// Metadata returns the internal stats for the Trie storage implementation
func Metadata(trie *Trie) store.Metadata {
	metadata := store.NewMetadata()
	for _, t := range trie.tries {
		metadata.Merge(Metadata(&t))
		metadata.KeysBytes++
	}
	if trie.value != nil && len(trie.value) > 0 {
		metadata.Size++
		metadata.ValuesBytes += uint64(len(trie.value))
	}
	metadata.KeysBytes++
	return metadata
}

// String prints the contents of the whole Tries
func (t *Trie) String() string {
	return fmt.Sprintf("key:%v,value:%v,tries:\n\t-> %v", t.key, t.value, t.tries)
}
