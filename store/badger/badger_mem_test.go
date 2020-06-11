package badger

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestBadgerInMem_KeyValueImplementation(t *testing.T) {
	new(test.KeyValue).Run(t, BadgerMemoryFactory)
}

func TestBadgerInMem_SyncImplementation(t *testing.T) {
	new(test.Concurrent).Run(t, BadgerMemoryFactory)
}
