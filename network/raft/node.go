package raft

import (
	"fmt"
	"time"

	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/store"
)

// Node represents a raft node implementation
type Node struct {
	msgID uint32
	*stateMachine
	storage network.Storage
	signal  chan Signal
}

// Cluster exposes the member properties of a nide
func (n *Node) Cluster() network.Member {
	return n.storage.Cluster()
}

// NeNode creates a new raft node
func NewNode(newCluster func(signal chan Signal) network.ProtocolFactory) network.NodeFactory {
	signal := make(chan Signal)
	return func(newStorage store.StorageFactory, clusterFactory network.ProtocolFactory) network.Storage {
		return node(signal, func(signal chan Signal) network.Storage {
			return network.Node(newStorage, newCluster(signal))
		})
	}
}

func node(signal chan Signal, n func(signal chan Signal) network.Storage) *Node {
	node := &Node{
		signal:       signal,
		stateMachine: newStateMachine(),
		storage:      n(signal),
	}
	return node
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
		Signal:    Append,
		HeartBeat: cmd.HeartBeat,
		response:  network.Response{},
	}

}

func (n *Node) commit(heartbeat HeartBeat) ResponseRPC {

	// for ease of use, we skip another verify at this stage
	// we assume the leader is 'sane' in the sense that it will adhere to the protocol

	// get all the pending entries from the state machine and commit them to the store
	for i := n.stateMachine.commitIndex; i <= heartbeat.logIndex; i++ {
		resp := network.Execute(n.storage, n.stateMachine.states[i].cmd)
		if resp.Err == nil {
			n.stateMachine.states[i].committed = true
		}
	}

	n.stateMachine.commitIndex = heartbeat.logIndex + 1

	return ResponseRPC{
		Signal: Commit,
		response: network.Response{
			Element: store.Nil,
		},
	}
}

// coordination logic

func (n *Node) Put(element store.Element) error {
	cmd := AppendRPC{
		HeartBeat: HeartBeat{
			Epoch: Epoch{
				leaderID: n.Cluster().ID,
			},
			Log: Log{
				prevLogIndex: n.commitIndex - 1,
				logIndex:     n.commitIndex,
			},
		},
		command: network.NewPut(element),
	}

	n.msgID++
	n.Cluster().Internal.Send(network.Message{ID: n.msgID, Source: n.Cluster().ID, Content: cmd})

	err := n.processSignal(Append, func() error {
		result := n.append(cmd)
		if result.response.Err != nil {
			return result.response.Err
		}
		n.msgID++
		n.Cluster().Internal.Send(network.Message{ID: n.msgID, Source: n.Cluster().ID, Content: cmd.HeartBeat})
		return nil
	})

	if err != nil {
		return fmt.Errorf("could not append log: %w", err)
	}

	return n.processSignal(Commit, func() error {
		//return nil
		rsp := n.commit(cmd.HeartBeat)
		return rsp.response.Err
	})
}

func (n *Node) processSignal(signal Signal, process func() error) error {
	select {
	case s := <-n.signal:
		if s == signal {
			return process()
		} else {
			return fmt.Errorf("unexpected signal received: '%v' instead of '%v'", s, signal)
		}
	case <-time.Tick(5 * time.Second):
		return fmt.Errorf("could not get consensus from cluster")
	}
}

func (n *Node) Get(key store.Key) (store.Element, error) {
	return n.storage.Get(key)
}

func (n *Node) Metadata() store.Metadata {
	return n.storage.Metadata()
}

func (n *Node) Close() error {
	return n.storage.Close()
}
