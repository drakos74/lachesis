package network

import "github.com/rs/zerolog/log"

// Key represents a generic byte array that is used for being able to partition traffic
type Key []byte

// Switch represents an entity much like a network switch or router
// it will keep track of all the cluster/network members
// and route traffic according to it's internal partitioning implementation
type Switch interface {
	Register(id int)
	DeRegister(id int)
	Route(key Key) ([]int, error)
	Members() []int
}

// PartitionStrategy is the factory type for creating a network switch abstraction
type PartitionStrategy func() Switch

// SinglePartition creates a switch that always delegates requests to a single node
func SinglePartition() Switch {
	return &UnarySwitch{}
}

// UnarySwitch is the single node Switch implementation
type UnarySwitch struct {
	Cluster
}

// Route returns the appropriate member for handling the given request based on the given key
func (u UnarySwitch) Route(key Key) ([]int, error) {
	return []int{0}, nil
}

// Members returns all the indexes of the current cluster members
func (u UnarySwitch) Members() []int {
	return []int{0}
}

// Cluster is the base implementation for the functionality of Register Deregister
// it keeps a slice of the indexes of the network members
// and appropriately removes or adds
type Cluster struct {
	members []int
}

// NewCluster creates a new cluster functionality
func NewCluster() Cluster {
	return Cluster{members: make([]int, 0)}
}

// Register registers a node to the cluster
func (c *Cluster) Register(id int) {
	log.Info().Int("Index", id).Msg("Register To Network")
	c.members = append(c.members, id)
}

// DeRegister removes a node from the cluster
func (c *Cluster) DeRegister(id int) {
	log.Info().Int("Index", id).Msg("De-Register From Network")
	copy(c.members[id:], c.members[id+1:]) // Shift a[i+1:] left one Index.
	c.members[len(c.members)-1] = 0        // Erase last element (write zero value).
	c.members = c.members[:len(c.members)-1]
}

// Members returns the cluster member indexes
func (c *Cluster) Members() []int {
	return c.members
}
