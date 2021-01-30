package paxos

import (
	"github.com/drakos74/lachesis/benchmarks/network"
	"github.com/drakos74/lachesis/store/app/storage"
)

// Proposal represents the paxos protocol proposal request
type Proposal struct {
	index   int64
	command network.Command
}

// Promise represents the paxos protocol promise response
type Promise struct {
	err error
}

// Commit represents the paxos protocol commit request
type Commit struct {
	index int64
	key   storage.Key
}
