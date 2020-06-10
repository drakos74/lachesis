package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/test"
)

func TestCache_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, func() store.Storage {
		return NewCache()
	})
}

func testCache_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, func() store.Storage {
		return NewCache()
	})
}
