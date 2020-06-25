package network

import (
	"testing"

	"github.com/drakos74/lachesis/internal/partition"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/mem"

	"github.com/drakos74/lachesis/store/test"
)

func newNetwork() store.StorageFactory {
	return networkFactory(1, partition.SinglePartition, mem.CacheFactory)
}

func TestNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newNetwork())
}

func TestNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newNetwork())
}

// Note : this will pass event for the
func TestNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newNetwork())
}

// Note : All the faulty tests should fail
func newFaultyNetwork() store.StorageFactory {
	return networkFactory(10, partition.RandomPartition, mem.CacheFactory)
}

func testFaultyNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newFaultyNetwork())
}

func testFaultyNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newFaultyNetwork())
}

func testFaultyNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newFaultyNetwork())
}

// simple Sharded network

func TestShardedNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, networkFactory(10, partition.ShardedPartition, mem.CacheFactory))
}

func TestShardedNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, networkFactory(10, partition.ShardedPartition, mem.CacheFactory))
}

// Note : this will pass event for the non-concurrent-safe stores
// because of our inherent synchronization at network level
func TestShardedNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, networkFactory(10, partition.ShardedPartition, mem.CacheFactory))
}

// failure for sharded network in case of node down event
func TestShardedNetwork_Failure(t *testing.T) {
	new(test.Consistency).Run(t, networkFactory(10, partition.ShardedPartition, mem.CacheFactory, NewNodeDownEvent(5, 30)))
}

// full replication network

// fixed failure condition with replica network
func TestReplicaNetwork_Failure_Resilience(t *testing.T) {
	new(test.Consistency).Run(t, networkFactory(10, partition.ReplicaPartition, mem.CacheFactory, NewNodeDownEvent(5, 30)))
}

// consistent hashing network

// fixed failure condition with consistent hashing network
// TODO : improve distribution metric
func TestConsistentNetwork_Failure_Resilience(t *testing.T) {
	new(test.Consistency).Run(t, networkFactory(10, partition.ConsistentPartition, mem.CacheFactory, NewNodeDownEvent(5, 30)))
}
