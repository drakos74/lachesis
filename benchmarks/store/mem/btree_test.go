package mem

import (
	"github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestBTree_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.BTreeFactory)
}

func testBTreeSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.BTreeFactory)
}
