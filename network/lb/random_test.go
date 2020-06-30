package lb

import (
	"testing"

	"github.com/drakos74/lachesis/network"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/test"
)

// Note : All the faulty tests should fail
// This network is problematic because ...
// we are using random partitioninng (routing)
// but have no replication
// so you put something on one node, but try to retrieve it from another
func newFaultyNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(RandomPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Create()
}

func testFaultyNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newFaultyNetwork())
}

func testFaultyNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newFaultyNetwork())
}

func testFaultyNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newFaultyNetwork())
}

func TestFaultyNetwork_SimpleFailureRate(t *testing.T) {
	// {"level":"info","write":"0.00","read":"90.60","time":"2020-06-28T11:45:07+02:00","message":"Error Rate"}
	new(test.FailureRate).Run(t, newFaultyNetwork(), test.Limit{
		Read:  0.0,
		Write: 95.0,
	})
}
