package lb

import (
	"math/big"

	"github.com/drakos74/lachesis/store/network"
)

func ReplicaPartition() network.Switch {
	return &ReplicaSwitch{replicas: 3}
}

type ReplicaSwitch struct {
	replicas    int
	parallelism int
}

func (r *ReplicaSwitch) Register(id int) {
	r.parallelism++
}

func (r *ReplicaSwitch) DeRegister(id int) {
	// nothing to do
}

func (r ReplicaSwitch) Route(key network.Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{mod(hash, r.parallelism), mod(hash+1, r.parallelism), mod(hash+2, r.parallelism)}, nil
}
