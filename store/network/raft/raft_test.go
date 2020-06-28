package raft

import (
	"testing"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/network"
	"github.com/drakos74/lachesis/store/test"
)

func newRaftNetwork(event ...network.Event) store.StorageFactory {
	return network.Factory(event...).
		Router(LeaderFollowerPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Node(RaftNode).
		Create()
}

func TestNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newRaftNetwork())
}

func TestNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newRaftNetwork())
}

func TestNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newRaftNetwork())
}

func TestNetwork_SimpleNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newRaftNetwork(network.NewNodeDownEvent(3, 30)), test.Limit{})
}

func TestNetwork_LeaderNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newRaftNetwork(network.NewNodeDownEvent(9, 30)), test.Limit{})
}
