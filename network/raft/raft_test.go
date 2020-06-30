package raft

import (
	"testing"

	"github.com/drakos74/lachesis/network/lb"

	"github.com/stretchr/testify/assert"

	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/test"
)

func newRaftNetwork(event ...network.Event) store.StorageFactory {
	signal := make(chan Signal)
	return network.Factory(event...).
		Router(lb.LeaderFollowerPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Protocol(RaftProtocol(signal)).
		Node(RaftNode(RaftProtocol)).
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
