package file

import (
	"fmt"
	"testing"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/test"
)

func newScratchPad() *ScratchPad {
	pad, err := NewScratchPad("data")
	if err != nil {
		panic(fmt.Sprintf("error during store creation: %v", err))
	}
	return pad
}

func TestFile_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, func() store.Storage {
		return newScratchPad()
	})
}

func testFile_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, func() store.Storage {
		return newScratchPad()
	})
}
