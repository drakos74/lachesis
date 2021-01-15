package badger

import (
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestBadgerInMem_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, MemoryFactory)
}

func TestBadgerInMem_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, MemoryFactory)
}
