package bolt

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestBoltFile_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, BoltFileFactory("data"))
}

func TestBoltFile_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, BoltFileFactory("data"))
}
