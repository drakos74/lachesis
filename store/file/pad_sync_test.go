package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestSyncFileImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, SyncScratchPadFactory("sync-data"))
}

func TestSyncFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, SyncScratchPadFactory("sync-data"))
}

func TestSyncBTreeFileImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, SyncTreePadFactory("sync-data"))
}

func TestSyncBTreeFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, SyncTreePadFactory("sync-data"))
}
