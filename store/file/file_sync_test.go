package file

import (
	"fmt"
	"testing"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/test"
)

func newSyncScratchPad() *SyncScratchPad {
	pad, err := NewSyncScratchPad("data")
	if err != nil {
		panic(fmt.Sprintf("error during store creation: %v", err))
	}
	return pad
}

func TestSyncFileImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, func() store.Storage {
		return newSyncScratchPad()
	})
}

func TestSyncFile_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, func() store.Storage {
		return newSyncScratchPad()
	})
}
