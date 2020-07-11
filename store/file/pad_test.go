package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, TriePadFactory("data"))
}

func testFileSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, TriePadFactory("data"))
}

func TestBTreeFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, TreePadFactory("data"))
}

func testBTreeFileSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, TreePadFactory("data"))
}
