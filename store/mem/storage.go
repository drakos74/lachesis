package mem

import (
	"fmt"
	"sync"

	"github.com/drakos74/lachesis/model"
	"github.com/drakos74/lachesis/store/trie"
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

// Metadata returns implementation statistics
func (c *Cache) Metadata() model.Metadata {
	keyBytes := 0
	valueBytes := 0
	for k, v := range c.storage {
		keyBytes += len(k)
		valueBytes += len(v)
	}
	return model.Metadata{
		Size:        len(c.storage),
		KeysBytes:   keyBytes,
		ValuesBytes: valueBytes,
		Errors:      make([]error, 0),
	}
}

// NewCache creates a new Cache instance
func NewCache() *Cache {
	return &Cache{storage: make(map[string]payload)}
}

// Trie implementation

// Trie is an in memory struct implementing the storage interface
// it s the most efficient one in terms of performance and is used for a baseline regarding tests
type Trie struct {
	storage *trie.Trie
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

// Close will run any maintenance operations
func (t *Trie) Close() error {
	return nil
}

// Metadata returns implementation statistics
func (t *Trie) Metadata() model.Metadata {
	return trie.Metadata(t.storage)
}

// NewTrie creates a new Cache instance
func NewTrie() *Trie {
	return &Trie{storage: trie.NewTrie(byte(' '))}
}

// sync Cache implementation

// SyncCache is an in memory struct implementing the storage interface backed by a sync.Map
type SyncCache struct {
	storage sync.Map
}

// Put adds an element to the cache
func (sc *SyncCache) Put(element model.Element) error {
	sc.storage.Store(string(element.Key()), element.Value())
	return nil
}

// Get retrieves and element from the cache
func (sc *SyncCache) Get(element model.Element) (model.Element, error) {
	if result, ok := sc.storage.Load(string(element.Key())); ok {
		return model.NewObject(element.Key(), result.([]byte)), nil
	}
	return nil, fmt.Errorf("could not find element for key %v", element.Key())
}

// Close will run any maintainance operations
func (sc *SyncCache) Close() error {
	return nil
}

// Metadata returns implementation statistics
func (sc *SyncCache) Metadata() model.Metadata {
	meta := model.NewMetadata()
	sc.storage.Range(func(key, value interface{}) bool {
		meta.Merge(model.Metadata{
			Size:        1,
			KeysBytes:   len(key.(string)), // TODO : not exact count, but should do for now
			ValuesBytes: len(value.([]byte)),
			Errors:      nil,
		})
		return true
	})
	return meta
}

// NewCache creates a new Cache instance
func NewSyncCache() *SyncCache {
	return &SyncCache{}
}

// Sync Trie implementation

// SyncTrie is an in memory struct implementing the storage interface
// it s the most efficient one in terms of performance and is used for a baseline regarding tests
type SyncTrie struct {
	storage *trie.Trie
	sync.RWMutex
}

// Put adds an element to the trie
func (st *SyncTrie) Put(element model.Element) error {
	st.Lock()
	defer st.Unlock()
	return st.storage.Commit(element.Key(), element.Value())
}

// Get retrieves and element from the trie
func (st *SyncTrie) Get(element model.Element) (model.Element, error) {
	st.RLock()
	defer st.RUnlock()
	if data, ok := st.storage.Read(element.Key()); ok {
		return model.NewObject(element.Key(), data), nil
	}
	return nil, fmt.Errorf("could not find element for key %v", element.Key())
}

// Close will run any maintainance operations
func (st *SyncTrie) Close() error {
	return nil
}

// Metadata returns implementation statistics
func (st *SyncTrie) Metadata() model.Metadata {
	return trie.Metadata(st.storage)
}

// NewTrie creates a new Cache instance
func NewSyncTrie() *SyncTrie {
	return &SyncTrie{storage: trie.NewTrie(byte(' '))}
}
