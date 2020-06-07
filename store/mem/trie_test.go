package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/test"
)

func TestTrieImplementation(t *testing.T) {
	test.Execute(t, func() store.Storage {
		return NewTrie()
	})
}
