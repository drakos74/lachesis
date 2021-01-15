package mem

import (
	"github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestSyncTrie_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.SyncTrieFactory)
}

func TestSyncTrie_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.SyncTrieFactory)
}
