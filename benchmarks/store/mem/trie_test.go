package mem

import (
	"github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestTrie_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.TrieFactory)
}

func testTrieSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.TrieFactory)
}
