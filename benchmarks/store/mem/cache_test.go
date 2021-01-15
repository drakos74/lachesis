package mem

import (
	"github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/store/test"
)

func TestCache_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, mem.CacheFactory)
}

func testCacheSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, mem.CacheFactory)
}
