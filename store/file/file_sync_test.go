package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestSyncFileImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, SyncScratchPadFactory("sync-data"))
}

func TestSyncFile_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, SyncScratchPadFactory("sync-data"))
}
