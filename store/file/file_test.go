package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestFile_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, ScratchPadFactory("data"))
}

func testFile_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, ScratchPadFactory("data"))
}
