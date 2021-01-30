package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/io/mem"
	"github.com/drakos74/lachesis/store/test"
)

func TestBTree_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.BTreeFactory)
}

func testBTreeSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.BTreeFactory)
}
