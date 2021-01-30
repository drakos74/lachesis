package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/io/mem"
	"github.com/drakos74/lachesis/store/test"
)

func TestSyncCache_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.SyncCacheFactory)
}

func TestSyncCache_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.SyncCacheFactory)
}
