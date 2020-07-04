package paxos

import (
	"testing"

	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/network/lb"
	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/test"
	"github.com/stretchr/testify/assert"
)

func newPaxosNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(lb.RandomPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Protocol(PaxosProtocol()).
		Node(network.Node).
		Create()
}

func TestPaxosProtocol(t *testing.T) {
	net := newPaxosNetwork()()

	err := net.Put(test.Random(10, 100).ElementFactory())
	assert.NoError(t, err)

	net.Metadata()
}

func TestNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newPaxosNetwork())
}

func TestNetwork_SimpleNodeDownImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newPaxosNetwork(network.NewNodeDownEvent(0, 30)))
}

func testNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newPaxosNetwork())
}

func testNetwork_SyncNodeDownImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newPaxosNetwork(network.NewNodeDownEvent(0, 30)))
}

func testNetwork_SimpleNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newPaxosNetwork(network.NewNodeDownEvent(3, 30)), test.Limit{})
}

func testNetwork_LeaderNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newPaxosNetwork(network.NewNodeDownEvent(0, 30)), test.Limit{})
}
