package lb

import (
	"math/big"

	"github.com/drakos74/lachesis/network"
)

func ShardedPartition() network.Switch {
	return &ShardedSwitch{}
}

type ShardedSwitch struct {
	network.Cluster
}

func (s ShardedSwitch) Route(key network.Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{int(hash % uint64(len(s.Members())))}, nil
}
