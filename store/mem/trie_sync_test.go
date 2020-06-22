package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestSyncTrie_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, SyncTrieFactory)
}

func TestSyncTrie_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, SyncTrieFactory)
}
