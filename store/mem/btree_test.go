package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestBTree_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, BTreeFactory)
}

func testBTreeSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, BTreeFactory)
}
