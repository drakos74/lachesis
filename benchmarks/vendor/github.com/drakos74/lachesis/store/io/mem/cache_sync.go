package mem

import (
	"fmt"
	"sync"

	"github.com/drakos74/lachesis/store/app/storage"
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

// SyncCacheFactory generates a SyncCache storage implementation
func SyncCacheFactory() storage.Storage {
	return NewSyncCache()
}

// Put adds an element to the cache
func (sc *SyncCache) Put(element storage.Element) error {
	sc.storage.Store(string(element.Key), element.Value)
	return nil
}

// Get retrieves and element from the cache
func (sc *SyncCache) Get(key storage.Key) (storage.Element, error) {
	if result, ok := sc.storage.Load(string(key)); ok {
		return storage.NewElement(key, result.(storage.Value)), nil
	}
	return storage.Element{}, fmt.Errorf(storage.NoValue, key)
}

// Close will run any maintenance operations
func (sc *SyncCache) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (sc *SyncCache) Metadata() storage.Metadata {
	meta := storage.NewMetadata()
	sc.storage.Range(func(key, value interface{}) bool {
		meta.Merge(storage.Metadata{
			Size:        1,
			KeysBytes:   uint64(len([]byte(key.(string)))),
			ValuesBytes: uint64(len(value.(storage.Value))),
			Errors:      nil,
		})
		return true
	})
	return meta
}
