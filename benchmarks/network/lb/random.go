package lb

import (
	"math/rand"

	"github.com/drakos74/lachesis/benchmarks/network"
)

// RandomPartition creates a switch with partition strategy that routes requests to random nodes
func RandomPartition() network.Switch {
	return &RandomSwitch{}
}

// RandomSwitch emulates a network switch logic based on randomness
type RandomSwitch struct {
	parallelism int
}

// Register registers a node to the network emulation switch
func (r *RandomSwitch) Register(id int) {
	r.parallelism++
}

// DeRegister removes a node from the list of active network nodes
func (r *RandomSwitch) DeRegister(id int) {
	// nothing to do
}

// Route returns the appropriate node to which the request should be routed based on the given key
func (r RandomSwitch) Route(key network.Key) ([]int, error) {
	return []int{rand.Intn(r.parallelism)}, nil
}

// Members returns the current active cluster members
func (r RandomSwitch) Members() []int {
	members := make([]int, r.parallelism)
	for i := 0; i < r.parallelism; i++ {
		members[i] = i
	}
	return members
}
