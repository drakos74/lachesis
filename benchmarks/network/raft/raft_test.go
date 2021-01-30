package raft

import (
	"testing"

	"github.com/drakos74/lachesis/benchmarks/network"
	"github.com/drakos74/lachesis/benchmarks/network/lb"
	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/drakos74/lachesis/store/io/mem"
	"github.com/drakos74/lachesis/store/test"
	"github.com/stretchr/testify/assert"
)

func newRaftNetwork(event ...network.Event) storage.StorageFactory {
	return network.Factory(event...).
		Router(lb.LeaderFollowerPartition).
		Storage(mem.SyncCacheFactory).
		Nodes(10).
		Protocol(Protocol()).
		Node(network.Node).
		Create()
}

func TestRaftProtocol(t *testing.T) {
	net := newRaftNetwork()()

	err := net.Put(test.Random(10, 100).ElementFactory())
	assert.NoError(t, err)

	net.Metadata()
}

func TestNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newRaftNetwork())
}

func TestNetwork_SimpleNodeDownImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newRaftNetwork(network.NewNodeDownEvent(0, 30)))
}

func TestNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newRaftNetwork())
}

func TestNetwork_SyncNodeDownImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newRaftNetwork(network.NewNodeDownEvent(0, 30)))
}

func TestNetwork_SimpleNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newRaftNetwork(network.NewNodeDownEvent(3, 30)), test.Limit{})
}

func TestNetwork_LeaderNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newRaftNetwork(network.NewNodeDownEvent(0, 30)), test.Limit{})
}
