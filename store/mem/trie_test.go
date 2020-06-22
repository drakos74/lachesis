package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestTrie_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, TrieFactory)
}

func testTrie_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, TrieFactory)
}
