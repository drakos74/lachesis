package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/io/file"
	"github.com/drakos74/lachesis/store/test"
)

func TestFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, file.TriePadFactory("data"))
}

func testFileSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, file.TriePadFactory("data"))
}

func TestBTreeFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, file.TreePadFactory("data"))
}

func testBTreeFileSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, file.TreePadFactory("data"))
}
