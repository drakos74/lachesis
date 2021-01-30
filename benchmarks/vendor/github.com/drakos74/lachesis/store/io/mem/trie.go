package mem

import (
	"fmt"

	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/drakos74/lachesis/store/datastruct/trie"
)

// Trie is an in memory struct implementing the storage interface
type Trie struct {
	storage *trie.Trie
}

// NewTrie creates a new Cache instance
func NewTrie() *Trie {
	return &Trie{storage: trie.NewTrie(byte(' '))}
}

// TrieFactory generates a Trie storage implementation
func TrieFactory() storage.Storage {
	return NewTrie()
}

// Put adds an element to the trie
func (t *Trie) Put(element storage.Element) error {
	return t.storage.Commit(element.Key, element.Value)
}

// Get retrieves and element from the trie
func (t *Trie) Get(key storage.Key) (storage.Element, error) {
	if data, ok := t.storage.Read(key); ok {
		return storage.NewElement(key, data), nil
	}
	return storage.Nil, fmt.Errorf(storage.NoValue, key)
}

// Close will run any maintenance operations
func (t *Trie) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (t *Trie) Metadata() storage.Metadata {
	return trie.Metadata(t.storage)
}
