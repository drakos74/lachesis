package raft

import "github.com/drakos74/lachesis/store/network"

// RaftProtocol implements the internal cluster communication requirements,
// e.g. the leader keeping all the state machines up to date,
// so that anyone can take over if needed
func RaftProtocol() network.Internal {

	return network.Protocol(func(in network.Message) network.Message {
		return in
	})

}
