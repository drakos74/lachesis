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

func TestFileImplementation(t *testing.T) {
	test.Execute(t, func() store.Storage {
		return newScratchPad()
	})
}
