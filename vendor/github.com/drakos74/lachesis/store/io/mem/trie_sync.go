package mem

import (
	"fmt"
	"sync"

	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/drakos74/lachesis/store/datastruct/trie"
)

// SyncTrie is an in memory struct implementing the storage interface
// it s the most efficient one in terms of performance and is used for a baseline regarding tests
type SyncTrie struct {
	storage *trie.Trie
	sync.RWMutex
}

// NewSyncTrie creates a new Cache instance
func NewSyncTrie() *SyncTrie {
	return &SyncTrie{storage: trie.NewTrie(byte(' '))}
}

// SyncTrieFactory generates a SyncTrie storage implementation
func SyncTrieFactory() storage.Storage {
	return NewSyncTrie()
}

// Put adds an element to the trie
func (st *SyncTrie) Put(element storage.Element) error {
	st.Lock()
	defer st.Unlock()
	return st.storage.Commit(element.Key, element.Value)
}

// Get retrieves and element from the trie
func (st *SyncTrie) Get(key storage.Key) (storage.Element, error) {
	st.RLock()
	defer st.RUnlock()
	if data, ok := st.storage.Read(key); ok {
		return storage.NewElement(key, data), nil
	}
	return storage.Element{}, fmt.Errorf(storage.NoValue, key)
}

// Close will run any maintainance operations
func (st *SyncTrie) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (st *SyncTrie) Metadata() storage.Metadata {
	return trie.Metadata(st.storage)
}
