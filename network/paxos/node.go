package paxos

import (
	"fmt"
	"time"

	"github.com/drakos74/lachesis/network"
	"github.com/drakos74/lachesis/store"
)

type wal struct {
	index int64
	log   map[string]Proposal
}

type Node struct {
	wal
	msgID   uint32
	storage network.Storage
	signal  chan Signal
}

func (n *Node) Cluster() network.Member {
	return n.storage.Cluster()
}

func PaxosNode(newCluster func(signal chan Signal) network.ProtocolFactory) network.NodeFactory {
	signal := make(chan Signal)
	return func(newStorage store.StorageFactory, clusterFactory network.ProtocolFactory) network.Storage {
		return node(signal, func(signal chan Signal) network.Storage {
			return network.SingleNode(newStorage, newCluster(signal))
		})
	}
}

func node(signal chan Signal, n func(signal chan Signal) network.Storage) *Node {
	node := &Node{
		signal:  signal,
		storage: n(signal),
		wal: wal{
			log: make(map[string]Proposal),
		},
	}
	return node
}

func (n *Node) promise(proposal Proposal) Promise {
	if entry, ok := n.log[string(proposal.command.Element().Key)]; ok {
		if entry.index > proposal.index {
			return Promise{err: fmt.Errorf("a higher index '%v' already exists for key '%v'",
				entry.index, proposal.command.Element().Key)}
		}
	}
	n.log[string(proposal.command.Element().Key)] = proposal
	return Promise{}
}

func (n *Node) accept(commit Commit) network.Response {
	return network.Execute(n.storage, n.log[string(commit.key)].command)
}

// coordination logic

func (n *Node) Put(element store.Element) error {

	index := time.Now().UnixNano()
	proposal := Proposal{
		index:   index,
		command: network.NewPut(element),
	}

	n.msgID++
	n.Cluster().Internal.Send(network.Message{ID: n.msgID, Source: n.Cluster().ID, Content: proposal})

	err := n.processSignal(Prepare, func() error {
		n.msgID++
		n.Cluster().Internal.Send(network.Message{ID: n.msgID, Source: n.Cluster().ID, Content: Commit{index: proposal.index, key: element.Key}})
		return nil
	})

	if err != nil {
		return fmt.Errorf("could not prepare log: %w", err)
	}

	return n.processSignal(Accept, func() error {
		return n.storage.Put(element)
	})
}

func (n *Node) processSignal(signal Signal, process func() error) error {
	select {
	case s := <-n.signal:
		if s == signal {
			return process()
		} else {
			return fmt.Errorf("unexpected signal received: %v instead of %v", s, signal)
		}
	case <-time.Tick(15 * time.Second):
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
