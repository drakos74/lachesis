package lb

import (
	"github.com/drakos74/lachesis/benchmarks/network"
)

// LeaderFollowerPartition creates new leader follower switch
func LeaderFollowerPartition() network.Switch {
	return &LeaderFollower{Cluster: network.NewCluster()}
}

// LeaderFollower emulates a leader-follower cluster network switch
type LeaderFollower struct {
	network.Cluster
}

// Register registers a new node to the cluster
func (c *LeaderFollower) Register(id int) {
	c.Cluster.Register(id)
}

// Route returns always the leader of the cluster to route requests to
func (c *LeaderFollower) Route(key network.Key) ([]int, error) {
	// leader will always be the 'first' node
	return c.Members()[:1], nil
}
