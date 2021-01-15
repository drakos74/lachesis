package bolt

import (
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestBoltFile_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, FileFactory("data"))
}

func TestBoltFile_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, FileFactory("data"))
}
