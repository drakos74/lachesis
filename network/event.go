package network

import (
	"fmt"
)

// Events is a collection of Events
type Events struct {
	warmUp int
	index  int
	events []Event
}

// Event represents a generic network event that gets applied at the SWitch level
type Event interface {
	Switch
	Wrap(router Switch) Switch
	Reset() Switch
	Index() int
}

// NodeDown emulates the case where a node is not responsive
type NodeDown struct {
	index      int
	iterations int
	duration   int
	Switch
}

// NewNodeDownEvent creates a new NodeDown Event
func NewNodeDownEvent(index, duration int) *NodeDown {
	return &NodeDown{
		index:    index,
		duration: duration,
	}
}

// Wrap will decorate the current Switch implementation with the events one
func (u *NodeDown) Wrap(router Switch) Switch {
	u.Switch = router
	return u
}

// Reset returns the previous network state
func (u *NodeDown) Reset() Switch {
	return u.Switch
}

// Register registers a new member to the network
func (u *NodeDown) Register(id int) {
	u.Switch.Register(id)
}

// Index returns the Index of the node that is supposed to be unresponsive
func (u *NodeDown) Index() int {
	return u.index
}

// Route returns the node responsible for serving the current request
func (u *NodeDown) Route(key Key) ([]int, error) {
	ids, err := u.Switch.Route(key)
	liveIds := make([]int, 0)
	for _, id := range ids {
		if u.index >= 0 && u.index == id && u.iterations < u.duration {
			// we need to ignore this one e.g. node is down
			u.iterations++
		} else {
			liveIds = append(liveIds, id)
		}
	}
	if len(liveIds) == 0 {
		return []int{}, fmt.Errorf("node %d is not responding", u.index)
	}
	return liveIds, err
}
