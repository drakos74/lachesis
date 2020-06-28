package partition

import (
	"math/big"
	"math/rand"
	"sort"
	"strconv"
)

type Key []byte

type Switch interface {
	Register(id int)
	DeRegister(id int)
	Route(key Key) ([]int, error)
}

type PartitionStrategy func() Switch

func SinglePartition() Switch {
	return &UnarySwitch{}
}

type UnarySwitch struct {
}

func (u *UnarySwitch) Register(id int) {
	// nothing to do
}

func (u *UnarySwitch) DeRegister(id int) {
	// nothing to do
}

func (u UnarySwitch) Route(key Key) ([]int, error) {
	return []int{0}, nil
}

func RandomPartition() Switch {
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

func (r RandomSwitch) Route(key Key) ([]int, error) {
	return []int{rand.Intn(r.parallelism)}, nil
}

func ShardedPartition() Switch {
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

func (s ShardedSwitch) Route(key Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{int(hash % uint64(s.parallelism))}, nil
}

func ReplicaPartition() Switch {
	return &ReplicaSwitch{replicas: 3}
}

type ReplicaSwitch struct {
	replicas    int
	parallelism int
}

func (r *ReplicaSwitch) Register(id int) {
	r.parallelism++
}

func (r *ReplicaSwitch) DeRegister(id int) {
	// nothing to do
}

func (r ReplicaSwitch) Route(key Key) ([]int, error) {
	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()
	return []int{mod(hash, r.parallelism), mod(hash+1, r.parallelism), mod(hash+2, r.parallelism)}, nil
}

func ConsistentPartition() Switch {
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

func (c ConsistentSwitch) Route(key Key) ([]int, error) {
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

func LeaderFollowerPartition() Switch {
	return &CaptainCrewSwitch{replicas: 1}
}

type CaptainCrewSwitch struct {
	replicas    int
	parallelism int
}

func (r *CaptainCrewSwitch) Register(id int) {
	// this is always the captain
	// captain should know about the others
	r.parallelism = id
}

func (r *CaptainCrewSwitch) DeRegister(id int) {
	// nothing to do
}

func (r CaptainCrewSwitch) Route(key Key) ([]int, error) {
	return []int{r.parallelism}, nil
}