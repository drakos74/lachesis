package lb

import (
	mem2 "github.com/drakos74/lachesis/internal/infra/mem"
	"testing"

	"github.com/drakos74/lachesis/benchmarks/network"

	"github.com/drakos74/lachesis/benchmarks/store/test"
	"github.com/drakos74/lachesis/internal/app/store"
)

// consistent hashing network

func newConsistentNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(ConsistentPartition).
		Storage(mem2.CacheFactory).
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
