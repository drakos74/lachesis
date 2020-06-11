package badger

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestBadgerFile_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, BadgerFileFactory("data"))
}

func TestBadgerFile_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, BadgerFileFactory("data"))
}
