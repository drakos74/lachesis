package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestCache_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, CacheFactory)
}

func testCache_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, CacheFactory)
}
