package lb

import (
	"testing"

	"github.com/drakos74/lachesis/network"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/test"
)

// consistent hashing network

func newConsistentNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(ConsistentPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Create()
}

// fixed failure condition with consistent hashing network
// TODO : improve distribution metric
func TestConsistentNetwork_Failure_Resilience(t *testing.T) {
	new(test.Consistency).Run(t, newConsistentNetwork(network.NewNodeDownEvent(5, 30)))
}

func TestConsistentNetwork_NodeDownEventFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newConsistentNetwork(network.NewNodeDownEvent(0, 30)), test.Limit{})
}
