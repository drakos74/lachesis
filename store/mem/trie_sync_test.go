package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/test"
)

func TestSyncTrieImplementation(t *testing.T) {
	test.Execute(t, func() store.Storage {
		return NewSyncTrie()
	})
}
