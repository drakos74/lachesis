package badger

import (
	"fmt"
	"testing"
	"time"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/test"
)

func newFileStore() *Store {
	// use nano, in order to create a new store each time (we want the tests to remain independent at this stage)
	s, err := NewFileStore(fmt.Sprintf("data/%v", time.Now().UnixNano()))
	if err != nil {
		panic(fmt.Sprintf("error during store creation: %v", err))
	}
	return s
}

func TestBadgerFileImplementation(t *testing.T) {
	test.Execute(t, func() store.Storage {
		return newFileStore()
	})
}
