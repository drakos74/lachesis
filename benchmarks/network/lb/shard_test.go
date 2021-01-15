package lb

import (
	mem2 "github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/network"
	"github.com/drakos74/lachesis/benchmarks/store/test"
	"github.com/drakos74/lachesis/internal/app/store"
)

// simple Sharded network

func newShardedNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(ShardedPartition).
		Storage(mem2.CacheFactory).
		Nodes(10).
		Create()
}

// Fixing the issue from above, by using a sharding strategy
// e.g. we route commands based on the key to a different node
// it will still fail for NodeDown events, but only a smaller subset of keys will be affected
func TestShardedNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newShardedNetwork())
}

func TestShardedNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newShardedNetwork())
}

func TestShardedNetwork_SimpleFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newShardedNetwork(), test.Limit{})
}

// Note : this will pass event for the non-concurrent-safe stores
// because of our inherent synchronization at network level
func TestShardedNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newShardedNetwork())
}

// failure for sharded network in case of node down event
func TestShardedNetwork_Failure(t *testing.T) {
	new(test.Consistency).Run(t, newShardedNetwork(network.NewNodeDownEvent(5, 30)))
}

// Note : we have intermittent failures if we choose limit{Read:0.0,Write:0.0}
// this is because of the randomisation of sharding and node down event
// but this should be enough to signify that a sharded network fails in cases of node outages
func TestShardedNetwork_NodeDownEventFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newShardedNetwork(network.NewNodeDownEvent(5, 30)), test.Limit{Write: 2, Read: 1})
}
