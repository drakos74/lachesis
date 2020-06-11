package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestTrie_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, TrieFactory)
}

func testTrie_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, TrieFactory)
}
