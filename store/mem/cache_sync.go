package mem

import (
	"fmt"
	"sync"

	"github.com/drakos74/lachesis/store"
)

// SyncCache is an in memory struct implementing the storage interface
// this implementation is thread-safe
type SyncCache struct {
	storage sync.Map
}

// NewSyncCache creates a new Cache instance
func NewSyncCache() *SyncCache {
	return &SyncCache{}
}

// Put adds an element to the cache
func (sc *SyncCache) Put(element store.Element) error {
	sc.storage.Store(string(element.Key), element.Value)
	return nil
}

// Get retrieves and element from the cache
func (sc *SyncCache) Get(key store.Key) (store.Element, error) {
	if result, ok := sc.storage.Load(string(key)); ok {
		return store.NewElement(key, result.(store.Value)), nil
	}
	return store.Element{}, fmt.Errorf("could not find element for key %v", key)
}

// Close will run any maintenance operations
func (sc *SyncCache) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (sc *SyncCache) Metadata() store.Metadata {
	meta := store.NewMetadata()
	sc.storage.Range(func(key, value interface{}) bool {
		meta.Merge(store.Metadata{
			Size:        1,
			KeysBytes:   len(key.(string)), // TODO : not exact count, but should do for now
			ValuesBytes: len(value.(store.Value)),
			Errors:      nil,
		})
		return true
	})
	return meta
}
