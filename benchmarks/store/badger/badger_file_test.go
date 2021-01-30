package badger

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestBadgerFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, FileFactory("data"))
}

func TestBadgerFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, FileFactory("data"))
}
