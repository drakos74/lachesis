package mem

import (
	"fmt"
	"sync"

	"github.com/drakos74/lachesis/internal/datastruct/trie"
	"github.com/drakos74/lachesis/store"
)

// SyncTrie is an in memory struct implementing the storage interface
// it s the most efficient one in terms of performance and is used for a baseline regarding tests
type SyncTrie struct {
	storage *trie.Trie
	sync.RWMutex
}

// Put adds an element to the trie
func (st *SyncTrie) Put(element store.Element) error {
	st.Lock()
	defer st.Unlock()
	return st.storage.Commit(element.Key, element.Value)
}

// Get retrieves and element from the trie
func (st *SyncTrie) Get(key store.Key) (store.Element, error) {
	st.RLock()
	defer st.RUnlock()
	if data, ok := st.storage.Read(key); ok {
		return store.NewElement(key, data), nil
	}
	return store.Element{}, fmt.Errorf(store.NoValue, key)
}

// Close will run any maintainance operations
func (st *SyncTrie) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (st *SyncTrie) Metadata() store.Metadata {
	return trie.Metadata(st.storage)
}

// NewTrie creates a new Cache instance
func NewSyncTrie() *SyncTrie {
	return &SyncTrie{storage: trie.NewTrie(byte(' '))}
}
