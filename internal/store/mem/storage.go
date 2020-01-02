package mem

import (
	"fmt"
	"lachesis/internal/model"
	"lachesis/internal/store/trie"
)

type payload []byte

// Map implementation

// Cache is an in memory struct implementing the storage interface
// it s the most efficient one in terms of performance and is used for a baseline regarding tests
type Cache struct {
	storage map[string]payload
}

// Put adds an element to the cache
func (c *Cache) Put(element model.Element) error {
	c.storage[string(element.Key())] = element.Value()
	return nil
}

// Get retrieves and element from the cache
func (c *Cache) Get(element model.Element) (model.Element, error) {
	if result, ok := c.storage[string(element.Key())]; ok {
		return model.NewObject(element.Key(), result), nil
	}
	return nil, fmt.Errorf("could not find element for key %v", element.Key())
}

// Close will run any maintainance operations
func (c *Cache) Close() error {
	return nil
}

// NewCache creates a new Cache instance
func NewCache() *Cache {
	return &Cache{storage: make(map[string]payload)}
}

// Trie implementation

// Cache is an in memory struct implementing the storage interface
// it s the most efficient one in terms of performance and is used for a baseline regarding tests
type Trie struct {
	storage trie.Trie
}

// Put adds an element to the trie
func (t *Trie) Put(element model.Element) error {
	return t.storage.Commit(element.Key(), element.Value())
}

// Get retrieves and element from the trie
func (t *Trie) Get(element model.Element) (model.Element, error) {
	if data, ok := t.storage.Read(element.Key()); ok {
		return model.NewObject(element.Key(), data), nil
	}
	return nil, fmt.Errorf("could not find element for key %v", element.Key())
}

// Close will run any maintainance operations
func (t *Trie) Close() error {
	return nil
}

// NewTrie creates a new Cache instance
func NewTrie() *Trie {
	return &Trie{storage: trie.NewTrie(byte(' '))}
}
