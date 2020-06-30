package lb

import (
	"math/rand"

	"github.com/drakos74/lachesis/network"
)

func RandomPartition() network.Switch {
	return &RandomSwitch{}
}

type RandomSwitch struct {
	parallelism int
}

func (r *RandomSwitch) Register(id int) {
	r.parallelism++
}

func (r *RandomSwitch) DeRegister(id int) {
	// nothing to do
}

func (r RandomSwitch) Route(key network.Key) ([]int, error) {
	return []int{rand.Intn(r.parallelism)}, nil
}

func (r RandomSwitch) Members() []int {
	members := make([]int, r.parallelism)
	for i := 0; i < r.parallelism; i++ {
		members[i] = i
	}
	return members
}
