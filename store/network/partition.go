package network

import (
	"math/big"
	"math/rand"
)

type Switch interface {
	Register(id int)
	Route(cmd Command) ([]int, error)
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

func (u UnarySwitch) Route(cmd Command) ([]int, error) {
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

func (r RandomSwitch) Route(cmd Command) ([]int, error) {
	return []int{rand.Intn(r.parallelism)}, nil
}

func ShardedPartition() Switch {
	return &ShardedSwitch{}
}

type ShardedSwitch struct {
	parallelism int
}

func (r *ShardedSwitch) Register(id int) {
	r.parallelism++
}

func (r ShardedSwitch) Route(cmd Command) ([]int, error) {

	key := cmd.Element().Key

	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()

	return []int{int(hash % uint64(r.parallelism))}, nil
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

func (r ReplicaSwitch) Route(cmd Command) ([]int, error) {

	key := cmd.Element().Key

	var i big.Int
	// convert to int
	hash := i.SetBytes(key).Uint64()

	return []int{mod(hash, r.parallelism), mod(hash+1, r.parallelism), mod(hash+2, r.parallelism)}, nil
}

func mod(hash uint64, parallelism int) int {
	return int(hash % uint64(parallelism))
}
