package lb

import (
	"testing"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/network"
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

func TestNetwork_SimpleNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newLeaderNetwork(network.NewNodeDownEvent(3, 30)), test.Limit{})
}

func TestNetwork_LeaderNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newLeaderNetwork(network.NewNodeDownEvent(9, 30)), test.Limit{})
}
