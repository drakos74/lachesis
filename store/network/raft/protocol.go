package raft

import "github.com/drakos74/lachesis/store/network"

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

type AppendRPC struct {
	HeartBeat
	command network.Command
}

type ResponseRPC struct {
	HeartBeat
	response network.Response
}
