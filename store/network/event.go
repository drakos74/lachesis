package network

import (
	"fmt"
)

type EventRotation struct {
	index  int
	events []Event
}

type Event interface {
	Switch
	Wrap(router Switch) Switch
	Reset() Switch
}

type NodeDown struct {
	index      int
	iterations int
	duration   int
	Switch
}

func NewNodeDownEvent(index, duration int) *NodeDown {
	return &NodeDown{
		index:    index,
		duration: duration,
	}
}

func (u *NodeDown) Wrap(router Switch) Switch {
	u.Switch = router
	return u
}

func (u *NodeDown) Reset() Switch {
	return u.Switch
}

func (u *NodeDown) Register(id int) {
	u.Switch.Register(id)
}

func (u *NodeDown) Route(cmd Command) ([]int, error) {
	ids, err := u.Switch.Route(cmd)
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
