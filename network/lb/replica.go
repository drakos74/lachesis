package lb

import (
	"math/big"

	"github.com/drakos74/lachesis/network"
)

func ReplicaPartition() network.Switch {
	return &ReplicaSwitch{replicas: 3}
}

type ReplicaSwitch struct {
	network.Cluster
	replicas int
}

func (r ReplicaSwitch) Route(key network.Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{mod(hash, len(r.Members())), mod(hash+1, len(r.Members())), mod(hash+2, len(r.Members()))}, nil
}
