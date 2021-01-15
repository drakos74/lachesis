package mem

import (
	"github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestSyncBTree_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.SyncBTreeFactory)
}

func TestSyncBTree_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.SyncBTreeFactory)
}
