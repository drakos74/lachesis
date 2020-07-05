package paxos

import (
	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/store"
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
	key   store.Key
}
