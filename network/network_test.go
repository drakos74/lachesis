package network

import (
	"testing"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/mem"

	"github.com/drakos74/lachesis/store/test"
)

// Note : single node network
// This network should be consistent in terms of operations
// in the same way as the individual network implementations
// but should fail in case of external cluster events i.e. node-down
func newNetwork(event ...Event) store.StorageFactory {
	return Factory(event...).
		Router(SinglePartition).
		Storage(mem.CacheFactory).
		Nodes(1).
		Create()
}

func TestNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newNetwork())
}

func TestNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newNetwork())
}

// Note : this will pass event for the
func TestNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newNetwork())
}

func TestNetwork_SimpleNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newNetwork(), test.Limit{})
}
