package mem

import (
	"fmt"

	"github.com/drakos74/lachesis/internal/datastruct/trie"
	"github.com/drakos74/lachesis/store"
)

// Trie is an in memory struct implementing the storage interface
type Trie struct {
	storage *trie.Trie
}

// Put adds an element to the trie
func (t *Trie) Put(element store.Element) error {
	return t.storage.Commit(element.Key, element.Value)
}

// Get retrieves and element from the trie
func (t *Trie) Get(key store.Key) (store.Element, error) {
	if data, ok := t.storage.Read(key); ok {
		return store.NewElement(key, data), nil
	}
	return store.Element{}, fmt.Errorf("could not find element for key %v", key)
}

// Close will run any maintenance operations
func (t *Trie) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (t *Trie) Metadata() store.Metadata {
	return trie.Metadata(t.storage)
}

// NewTrie creates a new Cache instance
func NewTrie() *Trie {
	return &Trie{storage: trie.NewTrie(byte(' '))}
}
