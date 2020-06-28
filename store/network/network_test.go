package network

import (
	"testing"

	"github.com/drakos74/lachesis/internal/partition"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/mem"

	"github.com/drakos74/lachesis/store/test"
)

// Note : single node network
// This network should be consistent in terms of operations
// in the same way as the individual network implementations
// but should fail in case of external cluster events i.e. node-down
func newNetwork(event ...Event) store.StorageFactory {
	return Factory(event...).
		Router(partition.SinglePartition).
		Storage(mem.CacheFactory).
		Nodes(1).
		Create()
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

func TestNetwork_SimpleNodeDownEvent(t *testing.T) {
	new(test.FailureRate).Run(t, newNetwork(), test.Limit{})
}

// Note : All the faulty tests should fail
// This network is problematic because ...
// we are using random partitioninng (routing)
// but have no replication
// so you put something on one node, but try to retrieve it from another
func newFaultyNetwork(event ...Event) store.StorageFactory {
	return Factory(event...).
		Router(partition.RandomPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Create()
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

func TestFaultyNetwork_SimpleFailureRate(t *testing.T) {
	// {"level":"info","write":"0.00","read":"90.60","time":"2020-06-28T11:45:07+02:00","message":"Error Rate"}
	new(test.FailureRate).Run(t, newFaultyNetwork(), test.Limit{
		Read:  0.0,
		Write: 95.0,
	})
}

// simple Sharded network

func newShardedNetwork(event ...Event) store.StorageFactory {
	return Factory(event...).
		Router(partition.ShardedPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Create()
}

// Fixing the issue from above, by using a sharding strategy
// e.g. we route commands based on the key to a different node
// it will still fail for NodeDown events, but only a smaller subset of keys will be affected
func TestShardedNetwork_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, newShardedNetwork())
}

func TestShardedNetwork_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, newShardedNetwork())
}

func TestShardedNetwork_SimpleFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newShardedNetwork(), test.Limit{})
}

// Note : this will pass event for the non-concurrent-safe stores
// because of our inherent synchronization at network level
func TestShardedNetwork_SyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, newShardedNetwork())
}

// failure for sharded network in case of node down event
func TestShardedNetwork_Failure(t *testing.T) {
	new(test.Consistency).Run(t, newShardedNetwork(NewNodeDownEvent(5, 30)))
}

// Note : we have intermittent failures if we choose limit{Read:0.0,Write:0.0}
// this is because of the randomisation of sharding and node down event
// but this should be enough to signify that a sharded network fails in cases of node outages
func TestShardedNetwork_NodeDownEventFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newShardedNetwork(NewNodeDownEvent(5, 30)), test.Limit{Write: 0.3, Read: 0.3})
}

// full replication network

func newReplicatedNetwork(event ...Event) store.StorageFactory {
	return Factory(event...).
		Router(partition.ReplicaPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Create()
}

// fixed failure condition with replica network
func TestReplicaNetwork_Failure_Resilience(t *testing.T) {
	new(test.Consistency).Run(t, newReplicatedNetwork(NewNodeDownEvent(5, 30)))
}

func TestReplicaNetwork_NodeDownEventFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newReplicatedNetwork(NewNodeDownEvent(5, 30)), test.Limit{})
}

// consistent hashing network

func newConsistentNetwork(event ...Event) store.StorageFactory {
	return Factory(event...).
		Router(partition.ConsistentPartition).
		Storage(mem.CacheFactory).
		Nodes(10).
		Create()
}

// fixed failure condition with consistent hashing network
// TODO : improve distribution metric
func TestConsistentNetwork_Failure_Resilience(t *testing.T) {
	new(test.Consistency).Run(t, newConsistentNetwork(NewNodeDownEvent(5, 30)))
}

func TestConsistentNetwork_NodeDownEventFailureRate(t *testing.T) {
	new(test.FailureRate).Run(t, newConsistentNetwork(NewNodeDownEvent(5, 30)), test.Limit{})
}

// TODO : implement also with LeaderFollowerPartitionStrategy
