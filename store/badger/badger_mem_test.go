package badger

import (
	"fmt"
	"testing"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/test"
)

func newMemStore() *Store {
	s, err := NewMemoryStore()
	if err != nil {
		panic(fmt.Sprintf("error during store creation: %v", err))
	}
	return s
}

func TestBadgerInMemImplementation(t *testing.T) {
	test.Execute(t, func() store.Storage {
		return newMemStore()
	})
}
