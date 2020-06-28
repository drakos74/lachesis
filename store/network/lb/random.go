package lb

import (
	"math/rand"

	"github.com/drakos74/lachesis/store/network"
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
