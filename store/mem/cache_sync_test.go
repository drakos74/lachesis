package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestSyncCache_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, SyncCacheFactory)
}

func TestSyncCache_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, SyncCacheFactory)
}
