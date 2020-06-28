package lb

import (
	"math/big"
	"sort"
	"strconv"

	"github.com/drakos74/lachesis/store/network"
)

func ConsistentPartition() network.Switch {
	return &ConsistentSwitch{replicas: 3, nodes: make([]int, 0), hashMap: make(map[int]int)}
}

const unit = 360

type ConsistentSwitch struct {
	replicas int
	hashMap  map[int]int
	nodes    []int
}

func (c *ConsistentSwitch) Register(id int) {
	for i := 0; i < c.replicas; i++ {
		hash := mod(byteHash([]byte(strconv.Itoa(i)+" "+string(id))), unit)
		c.nodes = append(c.nodes, hash)
		c.hashMap[hash] = id
	}
	sort.Slice(c.nodes, func(i, j int) bool {
		return c.nodes[i] < c.nodes[j]
	})
}

func (c *ConsistentSwitch) DeRegister(id int) {
	// nothing to do for now ...
}

func (c ConsistentSwitch) Route(key network.Key) ([]int, error) {
	// convert to int
	hash := mod(byteHash(key), unit)
	idx := sort.Search(len(c.nodes), func(i int) bool { return c.nodes[i] >= hash })
	if idx == len(c.nodes) {
		idx = 0
	}

	nodes := make([]int, 0)

	for i := range c.nodes {
		if c.nodes[i] >= hash {
			nodes = append(nodes, c.hashMap[i])
			if len(nodes) == c.replicas {
				break
			}
		}
	}

	if len(nodes) != c.replicas {
		for i := 0; i < len(nodes)-c.replicas; i++ {
			nodes = append(nodes, c.hashMap[i])
		}
	}

	return nodes, nil
}

// very simple hashing function
func byteHash(bytes []byte) uint64 {
	var i big.Int
	hash := i.SetBytes(bytes).Uint64()
	return hash
}

// mod operation for a uint64
func mod(hash uint64, m int) int {
	return int(hash % uint64(m))
}
