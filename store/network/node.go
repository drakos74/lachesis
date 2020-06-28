package network

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/store"
	"github.com/google/uuid"
)

type Message struct {
	routingId uint32
	content   interface{}
}

type MsgProcessor func(in Message) Message

type Internal struct {
	ID      uint32
	in      chan Message
	out     chan Message
	Process MsgProcessor
}

func Protocol(process MsgProcessor) Internal {
	return Internal{
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

type StorageNode struct {
	Member
	store store.Storage
}

type Node func() *StorageNode

type NodeFactory func(newStorage store.StorageFactory, newCluster ProtocolFactory) *StorageNode

func SingleNode(newStorage store.StorageFactory, newCluster ProtocolFactory) *StorageNode {
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

func (n *StorageNode) start(ctx context.Context) error {

	// listen to internal cluster events
	go func() {
		for {
			select {
			case msg := <-n.Internal.in:
				n.Internal.out <- n.Internal.Process(msg)
			case <-ctx.Done():
				log.Debug().Msg("Closing member channel")
				return
			}
		}
	}()

	// listen to client events
	go func() {
		for {
			select {
			case cmd := <-n.Operation.in:
				response := n.Execute(cmd)
				n.Operation.out <- response
			case <-n.Meta.in:
				n.Meta.out <- n.Metadata()
			case <-ctx.Done():
				log.Debug().Msg("Closing storage channel")
				return
			}
		}
	}()

	return nil
}

// Execute will execute the command and produce the corresponding response
func (n *StorageNode) Execute(cmd Command) Response {
	element := store.Nil
	var err error
	element, err = cmd.Exec()(n)
	return Response{
		Element: element,
		Err:     err,
	}
}

// Storage interface

func (n *StorageNode) Put(element store.Element) error {
	println(fmt.Sprintf("storage put n = %v", n))
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
