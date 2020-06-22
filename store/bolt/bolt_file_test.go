package bolt

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestBoltFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, BoltFileFactory("data"))
}

func TestBoltFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, BoltFileFactory("data"))
}
