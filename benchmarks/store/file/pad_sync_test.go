package file

import (
	"github.com/drakos74/lachesis/internal/infra/file"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestSyncFileImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, file.SyncScratchPadFactory("sync-data"))
}

func TestSyncFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, file.SyncScratchPadFactory("sync-data"))
}

func TestSyncBTreeFileImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, file.SyncTreePadFactory("sync-data"))
}

func TestSyncBTreeFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, file.SyncTreePadFactory("sync-data"))
}
