package file

import (
	"github.com/drakos74/lachesis/internal/infra/file"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
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
