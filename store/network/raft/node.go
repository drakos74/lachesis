package raft

import (
	"fmt"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/network"
)

type Node struct {
	*stateMachine
	*network.StorageNode
}

func RaftNode(newStorage store.StorageFactory, newCluster network.ProtocolFactory) *network.StorageNode {
	return Node{
		stateMachine: newStateMachine(),
		StorageNode:  network.SingleNode(newStorage, newCluster),
	}.StorageNode
}

func NewNode(newNode network.Node) *Node {
	return &Node{
		stateMachine: newStateMachine(),
		StorageNode:  newNode(),
	}
}

func (n *Node) append(cmd AppendRPC) ResponseRPC {

	err := n.stateMachine.verify(cmd.HeartBeat)

	if err != nil {
		return ResponseRPC{
			response: network.Response{
				Element: store.Nil,
				Err:     fmt.Errorf("inconsistent node state: %w", err),
			},
		}
	}

	// add the new state
	n.stateMachine.append(cmd)

	return ResponseRPC{
		HeartBeat: HeartBeat{},
		response:  network.Response{},
	}

}

func (n *Node) commit(heartbeat HeartBeat) ResponseRPC {

	// for ease of use, we skip another verify at this stage
	// we assume the leader is 'sane' in the sense that it will adhere to the protocol

	// get all the pending entries from the state machine and commit them to the store
	for i := n.stateMachine.commitIndex; i <= heartbeat.logIndex; i++ {
		resp := n.StorageNode.Execute(n.stateMachine.states[i].cmd)
		if resp.Err == nil {
			n.stateMachine.states[i].committed = true
		}
	}

	n.stateMachine.commitIndex = heartbeat.logIndex + 1

	return ResponseRPC{
		response: network.Response{
			Element: store.Nil,
		},
	}
}

func (n *Node) Put(element store.Element) error {
	println(fmt.Sprintf("raft put n = %v", n))
	return n.StorageNode.Put(element)
}

func (n *Node) Get(key store.Key) (store.Element, error) {
	return n.StorageNode.Get(key)
}

func (n *Node) Metadata() store.Metadata {
	return n.StorageNode.Metadata()
}

func (n *Node) Close() error {
	return n.StorageNode.Close()
}
