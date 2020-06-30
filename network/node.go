package network

import (
	"github.com/drakos74/lachesis/store"
	"github.com/google/uuid"
)

var Void = Message{}

type Message struct {
	ID        uint32
	Source    uint32
	RoutingId uint32
	Content   interface{}
}

type MsgProcessor func(members int, node Storage, in Message) Message

type Internal struct {
	ID      uint32
	in      chan Message
	out     chan Message
	Process MsgProcessor
}

func (i Internal) Send(msg Message) {
	i.out <- msg
	// wait for the responses

}

func Protocol(id uint32, process MsgProcessor) Internal {
	return Internal{
		ID:      id,
		in:      make(chan Message),
		out:     make(chan Message),
		Process: process,
	}
}

type ProtocolFactory func(id uint32) Internal

func NoProtocol(id uint32) Internal {
	return Internal{ID: id}
}

type Member struct {
	Operation
	Meta
	Internal
}

type Storage interface {
	store.Storage
	Cluster() Member
}

type StorageNode struct {
	Member
	store store.Storage
}

type Node func() Storage

type NodeFactory func(newStorage store.StorageFactory, newCluster ProtocolFactory) Storage

func SingleNode(newStorage store.StorageFactory, newCluster ProtocolFactory) Storage {
	id := uuid.New().ID()
	return &StorageNode{
		Member: Member{
			Operation: Operation{
				in:  make(chan Command),
				out: make(chan Response),
			},
			Meta: Meta{
				out: make(chan store.Metadata),
				in:  make(chan struct{}),
			},
			Internal: newCluster(id),
		},
		store: newStorage(),
	}
}

// StorageNode internals

func (n *StorageNode) Cluster() Member {
	return n.Member
}

// Execute will execute the command and produce the corresponding response
func Execute(node Storage, cmd Command) Response {
	element := store.Nil
	var err error
	switch cmd.Type() {
	case Put:
		err = node.Put(cmd.Element())
	case Get:
		element, err = node.Get(cmd.Element().Key)
	}
	return Response{
		Element: element,
		Err:     err,
	}
}

// Storage interface

func (n *StorageNode) Put(element store.Element) error {
	return n.store.Put(element)
}

func (n *StorageNode) Get(key store.Key) (store.Element, error) {
	return n.store.Get(key)
}

func (n *StorageNode) Metadata() store.Metadata {
	return n.store.Metadata()
}

func (n *StorageNode) Close() error {
	return n.store.Close()
}
