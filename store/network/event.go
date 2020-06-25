package network

import (
	"fmt"

	"github.com/drakos74/lachesis/internal/partition"
)

type EventRotation struct {
	index  int
	events []Event
}

type Event interface {
	partition.Switch
	Wrap(router partition.Switch) partition.Switch
	Reset() partition.Switch
}

type NodeDown struct {
	index      int
	iterations int
	duration   int
	partition.Switch
}

func NewNodeDownEvent(index, duration int) *NodeDown {
	return &NodeDown{
		index:    index,
		duration: duration,
	}
}

func (u *NodeDown) Wrap(router partition.Switch) partition.Switch {
	u.Switch = router
	return u
}

func (u *NodeDown) Reset() partition.Switch {
	return u.Switch
}

func (u *NodeDown) Register(id int) {
	u.Switch.Register(id)
}

func (u *NodeDown) Route(key partition.Key) ([]int, error) {
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
