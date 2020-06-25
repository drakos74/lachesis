package network

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/store"
	"github.com/google/uuid"
)

type MessageType int

const (
	Leader = iota + 1
	Confirm
)

type Message struct {
	msgType MessageType
}

type State struct {
	id     uint32
	leader uint32
}

func (s *State) isActive() bool {
	return s.leader == s.id
}

type Member struct {
	State
	Port
	Meta
	in  chan Message
	out chan Message
}

type Node struct {
	Member
	store store.Storage
}

func NewNode(factory store.StorageFactory) *Node {
	// TODO : for now start as a leader
	id := uuid.New().ID()

	return &Node{
		Member: Member{
			State: State{
				id:     id,
				leader: id,
			},
			Port: Port{
				in:  make(chan Command),
				out: make(chan Response),
			},
			Meta: Meta{
				out: make(chan store.Metadata),
				in:  make(chan struct{}),
			},
			in:  make(chan Message),
			out: make(chan Message),
		},
		store: factory(),
	}
}

// Node internals

func (n *Node) start(ctx context.Context) error {

	// listen to internal cluster events
	go func() {
		for {
			select {
			case msg := <-n.in:
				switch msg.msgType {
				case Leader:
				case Confirm:
				}
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
			case cmd := <-n.Port.in:
				// handle only if we are the leader
				// for whatever this might mean
				// Note : we can always have an implementation where the node is fixed as leader e.g. no leader
				element := store.Nil
				var err error
				if n.isActive() {
					element, err = cmd.Exec()(n)
				} else {
					log.Err(fmt.Errorf("node is not active: '%v'", n.State))
				}
				n.Port.out <- Response{
					Element: element,
					Err:     err,
				}
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

// Storage interface

func (n *Node) Put(element store.Element) error {
	return n.store.Put(element)
}

func (n *Node) Get(key store.Key) (store.Element, error) {
	return n.store.Get(key)
}

func (n *Node) Metadata() store.Metadata {
	return n.store.Metadata()
}

func (n *Node) Close() error {
	return n.store.Close()
}
