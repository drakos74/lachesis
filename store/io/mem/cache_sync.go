package mem

import (
	"fmt"
	"github.com/drakos74/lachesis"
	"sync"
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
func SyncCacheFactory() lachesis.Storage {
	return NewSyncCache()
}

// Put adds an element to the cache
func (sc *SyncCache) Put(element lachesis.Element) error {
	sc.storage.Store(string(element.Key), element.Value)
	return nil
}

// Get retrieves and element from the cache
func (sc *SyncCache) Get(key lachesis.Key) (lachesis.Element, error) {
	if result, ok := sc.storage.Load(string(key)); ok {
		return lachesis.NewElement(key, result.(lachesis.Value)), nil
	}
	return lachesis.Element{}, fmt.Errorf(lachesis.NoValue, key)
}

// Close will run any maintenance operations
func (sc *SyncCache) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (sc *SyncCache) Metadata() lachesis.Metadata {
	meta := lachesis.NewMetadata()
	sc.storage.Range(func(key, value interface{}) bool {
		meta.Merge(lachesis.Metadata{
			Size:        1,
			KeysBytes:   uint64(len([]byte(key.(string)))),
			ValuesBytes: uint64(len(value.(lachesis.Value))),
			Errors:      nil,
		})
		return true
	})
	return meta
}
