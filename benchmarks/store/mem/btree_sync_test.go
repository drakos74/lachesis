package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/io/mem"
	"github.com/drakos74/lachesis/store/test"
)

func TestSyncBTree_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.SyncBTreeFactory)
}

func TestSyncBTree_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.SyncBTreeFactory)
}
