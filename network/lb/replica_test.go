package lb

import (
	"testing"

	"github.com/drakos74/lachesis/network"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/test"
)

// full replication network

func newReplicatedNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(ReplicaPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Create()
}

// fixed failure condition with replica network
func TestReplicaNetwork_Failure_Resilience(t *testing.T) {
	new(test.Consistency).Run(t, newReplicatedNetwork(network.NewNodeDownEvent(5, 30)))
}

func TestReplicaNetwork_NodeDownEventFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newReplicatedNetwork(network.NewNodeDownEvent(5, 30)), test.Limit{})
}

func TestReplicaNetwork_NodeLossEventFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newReplicatedNetwork(network.NewNodeDownEvent(5, 0)), test.Limit{})
}
