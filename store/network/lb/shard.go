package lb

import (
	"math/big"

	"github.com/drakos74/lachesis/store/network"
)

func ShardedPartition() network.Switch {
	return &ShardedSwitch{}
}

type ShardedSwitch struct {
	parallelism int
}

func (s *ShardedSwitch) Register(id int) {
	s.parallelism++
}

func (s *ShardedSwitch) DeRegister(id int) {
	// nothing to do
}

func (s ShardedSwitch) Route(key network.Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{int(hash % uint64(s.parallelism))}, nil
}
