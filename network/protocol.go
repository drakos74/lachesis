package network

import (
	"github.com/drakos74/lachesis/store"
)

// Internal defines the communication channels for inter-node communication within the network
// it acts as a central switch that routes traffic to the nodes, based on the message metadata
type Internal struct {
	ID      uint32
	in      chan Message
	out     chan Message
	Process MsgRouter
}

// ProtocolFactory allows the creation of Internal
type ProtocolFactory func(id uint32) (*Internal, *Peer)

// NoProtocol represents a void protocol
// this means there is no inter-node communication in the network
func NoProtocol(id uint32) (*Internal, *Peer) {
	return &Internal{ID: id}, &Peer{processor: *ProcessorFactory(func(state *State, node *StorageNode, element store.Element) (rpc interface{}, wait bool) {
		return nil, false
	})}
}
