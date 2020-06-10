package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/test"
)

func TestSyncCache_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, func() store.Storage {
		return NewSyncCache()
	})
}

func TestSyncCache_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, func() store.Storage {
		return NewSyncCache()
	})
}
