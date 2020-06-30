package raft

import "github.com/drakos74/lachesis/network"

type Signal int

const (
	Append Signal = iota + 1
	Commit
)

func (s Signal) String() string {
	switch s {
	case Append:
		return "append"
	case Commit:
		return "commit"
	}
	return ""
}

type Epoch struct {
	term     int
	leaderID uint32
}

type Log struct {
	prevLogIndex int
	prevLogTerm  int
	logIndex     int
}

type HeartBeat struct {
	Epoch
	Log
}

var NoBeat = HeartBeat{}

type AppendRPC struct {
	HeartBeat
	command network.Command
}

type ResponseRPC struct {
	Signal
	HeartBeat
	response network.Response
}
