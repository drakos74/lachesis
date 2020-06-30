package network

import "github.com/rs/zerolog/log"

type Key []byte

type Switch interface {
	Register(id int)
	DeRegister(id int)
	Route(key Key) ([]int, error)
	Members() []int
}

type PartitionStrategy func() Switch

func SinglePartition() Switch {
	return &UnarySwitch{}
}

type UnarySwitch struct {
	Cluster
}

func (u UnarySwitch) Route(key Key) ([]int, error) {
	return []int{0}, nil
}

func (u UnarySwitch) Members() []int {
	return []int{0}
}

// base implementation for Register Deregister
type Cluster struct {
	members []int
}

func NewCluster() Cluster {
	return Cluster{members: make([]int, 0)}
}

func (c *Cluster) Register(id int) {
	log.Info().Int("index", id).Msg("Register To Network")
	c.members = append(c.members, id)
}

func (c *Cluster) DeRegister(id int) {
	log.Info().Int("index", id).Msg("De-Register From Network")
	copy(c.members[id:], c.members[id+1:]) // Shift a[i+1:] left one index.
	c.members[len(c.members)-1] = 0        // Erase last element (write zero value).
	c.members = c.members[:len(c.members)-1]
}

func (c *Cluster) Members() []int {
	return c.members
}
