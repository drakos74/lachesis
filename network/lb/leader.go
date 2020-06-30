package lb

import (
	"github.com/drakos74/lachesis/network"
)

func LeaderFollowerPartition() network.Switch {
	return &LeaderFollower{Cluster: network.NewCluster()}
}

type LeaderFollower struct {
	network.Cluster
}

func (c *LeaderFollower) Register(id int) {
	c.Cluster.Register(id)
}

func (c *LeaderFollower) Route(key network.Key) ([]int, error) {
	// leader will always be the 'first' node
	return c.Members()[:1], nil
}
