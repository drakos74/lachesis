package raft

import "github.com/drakos74/lachesis/benchmarks/network"

// Epoch represents the leader and its term according to the raft protocol definition
type Epoch struct {
	term     int
	leaderID uint32
}

// Log stores the current indexes for the wal
type Log struct {
	prevLogIndex int64
	prevLogTerm  int64
	logIndex     int64
}

// HeartBeat is used to track the epoch and log for a raft node
type HeartBeat struct {
	Epoch
	Log
}

// AppendRPC is the append RPC command for the raft protocol
type AppendRPC struct {
	HeartBeat
	command network.Command
}

// ResponseRPC is the response object for the raft protocol
type ResponseRPC struct {
	HeartBeat
	response network.Response
}
