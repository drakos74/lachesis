package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/io/mem"
	"github.com/drakos74/lachesis/store/test"
)

func TestTrie_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.TrieFactory)
}

func testTrieSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.TrieFactory)
}
