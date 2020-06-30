package paxos

import (
	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/store"
)

type Signal int

const (
	Prepare Signal = iota + 1
	Accept
)

func (s Signal) String() string {
	switch s {
	case Prepare:
		return "prepare"
	case Accept:
		return "accept"
	}
	return ""
}

type Proposal struct {
	index   int64
	command network.Command
}

type Commit struct {
	index int64
	key   store.Key
}

type Promise struct {
	err error
}

type Acceptance struct {
	response network.Response
}
