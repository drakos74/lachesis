package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestSyncBTree_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, SyncBTreeFactory)
}

func TestSyncBTree_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, SyncBTreeFactory)
}
