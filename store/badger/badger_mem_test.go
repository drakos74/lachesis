package badger

import (
	"fmt"
	"testing"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/test"
)

func newMemStore() *Store {
	s, err := NewMemoryStore()
	if err != nil {
		panic(fmt.Sprintf("error during store creation: %v", err))
	}
	return s
}

func TestBadgerInMem_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, func() store.Storage {
		return newMemStore()
	})
}

func TestBadgerInMem_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, func() store.Storage {
		return newMemStore()
	})
}
