package lb

import (
	"math/big"

	"github.com/drakos74/lachesis/network"
)

// ReplicaPartition creates a switch with partition strategy that implements replicas
func ReplicaPartition() network.Switch {
	return &ReplicaSwitch{}
}

// ReplicaSwitch emulates a network switch which replicates the request to several nodes
type ReplicaSwitch struct {
	network.Cluster
}

// DeRegister removes a node from the cluster
func (c *ReplicaSwitch) DeRegister(id int) {
	// do nothing
}

// Route returns the appropriate node to which the request should be routed based on the given key
func (r ReplicaSwitch) Route(key network.Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{mod(hash, len(r.Members())), mod(hash+1, len(r.Members())), mod(hash+2, len(r.Members()))}, nil
}
