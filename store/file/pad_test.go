package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, ScratchPadFactory("data"))
}

func testFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, ScratchPadFactory("data"))
}

func TestBTreeFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, TreePadFactory("data"))
}

func testBTreeFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, TreePadFactory("data"))
}
