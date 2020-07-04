package paxos

import (
	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/store"
)

type Proposal struct {
	index   int64
	command network.Command
}

type Promise struct {
	err error
}

type Commit struct {
	index int64
	key   store.Key
}
