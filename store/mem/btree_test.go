package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestBTree_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, BTreeFactory)
}

func testBTree_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, BTreeFactory)
}
