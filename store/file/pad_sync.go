package file

import (
	"fmt"
	"sync"

	"github.com/drakos74/lachesis/store/mem"

	"github.com/drakos74/lachesis/store"
)

// SyncScratchPad is a thread-safe implementation of  file store
type SyncScratchPad struct {
	store *ScratchPad
	mutex sync.RWMutex
}

// NewSyncScratchPad creates a new file store that is thread-safe
// using as an index a trie
func NewSyncScratchPad(path string) (*SyncScratchPad, error) {
	sb, err := NewScratchPad(path, mem.SyncTrieFactory)
	if err != nil {
		return nil, err
	}
	return &SyncScratchPad{
		store: sb,
	}, nil
}

// SyncScratchPadFactory generates a synced file storage implementation
func SyncScratchPadFactory(path string) store.StorageFactory {
	return func() store.Storage {
		pad, err := NewSyncScratchPad(path)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return pad
	}
}

// NewSyncTreePad creates a new file store that is thread-safe
// using as an index a Btree
func NewSyncTreePad(path string) (*SyncScratchPad, error) {
	sb, err := NewScratchPad(path, mem.SyncBTreeFactory)
	if err != nil {
		return nil, err
	}
	return &SyncScratchPad{
		store: sb,
	}, nil
}

// SyncScratchPadFactory generates a synced file storage implementation
func SyncTreePadFactory(path string) store.StorageFactory {
	return func() store.Storage {
		pad, err := NewSyncTreePad(path)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return pad
	}
}

// Put adds an element to the store while using a write lock
func (ss *SyncScratchPad) Put(element store.Element) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	return ss.store.Put(element)
}

// Get retrieves an element from the store while using a read lock
func (ss *SyncScratchPad) Get(key store.Key) (store.Element, error) {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	return ss.store.Get(key)
}

// Close does any clean up
func (ss *SyncScratchPad) Close() error {
	return ss.store.Close()
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (ss *SyncScratchPad) Metadata() store.Metadata {
	return ss.store.Metadata()
}
