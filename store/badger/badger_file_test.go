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

func TestBadgerFile_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, func() store.Storage {
		return newFileStore()
	})
}

func TestBadgerFile_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, func() store.Storage {
		return newFileStore()
	})
}
