package lb

import (
	"math/big"

	"github.com/drakos74/lachesis/network"
)

// ShardedPartition creates a switch with partition strategy that implements sharding
func ShardedPartition() network.Switch {
	return &ShardedSwitch{}
}

// ShardedSwitch emulates a network switch which segments the traffic to different nodes based on the elements key
type ShardedSwitch struct {
	network.Cluster
}

// Route returns the appropriate node to which the request should be routed based on the given key
func (s ShardedSwitch) Route(key network.Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{int(hash % uint64(len(s.Members())))}, nil
}
