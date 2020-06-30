package lb

import (
	"testing"

	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/test"
)

func newLeaderNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(LeaderFollowerPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Node(network.SingleNode).
		Create()
}

func TestNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newLeaderNetwork())
}

func TestNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newLeaderNetwork())
}

func TestNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newLeaderNetwork())
}

func testNetwork_SyncFailingImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newLeaderNetwork(network.NewNodeDownEvent(0, 30)))
}

func TestNetwork_SimpleNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newLeaderNetwork(network.NewNodeDownEvent(3, 30)), test.Limit{})
}

func testNetwork_LeaderNodeDownEvent(t *testing.T) {
	// Note : we know from the routing strategy, that the leader is always the '0th' element
	new(test.FailureRate).Run(t, newLeaderNetwork(network.NewNodeDownEvent(0, 30)), test.Limit{})
}
