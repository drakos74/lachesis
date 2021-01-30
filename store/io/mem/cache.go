package mem

import (
	"fmt"
	"github.com/drakos74/lachesis"
)

// Cache is an in memory struct implementing the storage interface
// it s the most efficient one in terms of performance and is used for a baseline regarding tests
type Cache struct {
	storage map[string]lachesis.Value
}

// NewCache creates a new Cache instance
func NewCache() *Cache {
	return &Cache{storage: make(map[string]lachesis.Value)}
}

// CacheFactory generates a Cache storage implementation
func CacheFactory() lachesis.Storage {
	return NewCache()
}

// Put adds an element to the cache
func (c *Cache) Put(element lachesis.Element) error {
	c.storage[string(element.Key)] = element.Value
	return nil
}

// Get retrieves and element from the cache
func (c *Cache) Get(key lachesis.Key) (lachesis.Element, error) {
	if result, ok := c.storage[string(key)]; ok {
		element := lachesis.NewElement(key, result)
		return element, nil
	}
	return lachesis.Nil, fmt.Errorf(lachesis.NoValue, key)
}

// Close will run any maintenance operations for the store
func (c *Cache) Close() error {
	return nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (c *Cache) Metadata() lachesis.Metadata {
	var keyBytes uint64
	var valueBytes uint64
	for k, v := range c.storage {
		keyBytes += uint64(len(k))
		valueBytes += uint64(len(v))
	}
	return lachesis.Metadata{
		Size:        uint64(len(c.storage)),
		KeysBytes:   keyBytes,
		ValuesBytes: valueBytes,
		Errors:      make([]error, 0),
	}
}
