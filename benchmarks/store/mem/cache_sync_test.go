package mem

import (
	"github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestSyncCache_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.SyncCacheFactory)
}

func TestSyncCache_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.SyncCacheFactory)
}
