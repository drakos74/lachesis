package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/io/file"
	"github.com/drakos74/lachesis/store/test"
)

func TestClosingFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, file.TrieClosingPadFactory("data"))
}

func testClosingFileSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, file.TrieClosingPadFactory("data"))
}

func TestBTreeClosingFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, file.TreeClosingPadFactory("data"))
}

func testBTreeClosingFileSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, file.TreeClosingPadFactory("data"))
}
