package network

import (
	"github.com/drakos74/lachesis/store"
	"github.com/google/uuid"
)

// Void is an empty message
var Void = Message{}

// Message represents an internal communication object
// messages are used for the cluster nodes to communicate with each other
type Message struct {
	ID        uint32
	Source    uint32
	RoutingID uint32
	Type      MsgType
	Content   interface{}
	Err       error
}

// MsgRouter bears instructions on how to process a message
type MsgRouter func(members int, node Storage, in Message) Message

// Member represents a network member node
type Member struct {
	Operation
	Meta
	*Internal
}

// MsgType represents the type of a message
type MsgType int

// MsgProcessor bears the logic of processing a message
type MsgProcessor func(state *State, storage store.Storage, msg interface{}) (interface{}, error)

// State represents the wal of the storage node
type State struct {
	Index int64
	// TODO : make this also a storage interface to be able to inject different implementations for the WAL
	// only caveat is that the raw implementation is not enough, as we want to store an 'interface' and not raw bytes
	Log map[string]interface{}
}

// NewStateLog creates a new state log for a node
func NewStateLog() State {
	return State{
		Log: make(map[string]interface{}),
	}
}

// Processor encapsulates all logic related to the internal cluster communication protocol
type Processor struct {
	Store    store.Storage
	State    State
	initiate func(state *State, node *StorageNode, element store.Element) (rpc interface{}, wait bool)
	handle   map[MsgType]MsgProcessor
}

// ProcessorFactory creates a new processor
func ProcessorFactory(initPut func(state *State, node *StorageNode, element store.Element) (rpc interface{}, wait bool)) *Processor {
	return &Processor{
		initiate: initPut,
		handle:   make(map[MsgType]MsgProcessor),
	}
}

// Storage adds a storage implementation to the processor
func (p *Processor) Storage(storage store.Storage) *Processor {
	p.Store = storage
	return p
}

// Create creates a processor
func (p *Processor) Create() *Processor {
	// TODO : make checks
	return p
}

// Peer encapsulates functionality regarding peer to peer communication between cluster members
type Peer struct {
	msgID     uint32
	processor Processor
}

func (p *Peer) initPut(node *StorageNode, element store.Element) (rpc interface{}, wait bool) {
	return p.processor.initiate(&p.processor.State, node, element)
}

func (p *Peer) getProcessor(msgType MsgType) MsgProcessor {
	if process, ok := p.processor.handle[msgType]; ok {
		return func(state *State, storage store.Storage, msg interface{}) (interface{}, error) {
			return process(&p.processor.State, p.processor.Store, msg)
		}
	}
	return func(state *State, storage store.Storage, msg interface{}) (interface{}, error) {
		return nil, nil
	}
}

// Storage represents a network member node with key-value Storage capabilities
type Storage interface {
	store.Storage
	Cluster() Member
}

// StorageNode is the implementation for a node that can act as member of a network
// and also act as a key-value Storage
type StorageNode struct {
	Member
	*Peer
}

// NodeFactory is the factory type for a StorageNode
type NodeFactory func(newStorage store.StorageFactory, newCluster ProtocolFactory) Storage

// Node is the factory for creating a cluster node with the given properties
func Node(newStorage store.StorageFactory, newCluster ProtocolFactory) Storage {
	id := uuid.New().ID()
	protocol, peer := newCluster(id)
	storage := newStorage()
	peer.processor.Storage(storage)
	peer.processor.State = NewStateLog()
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
			Internal: protocol,
		},
		Peer: peer,
	}
}

// StorageNode internals

// Cluster exposes the cluster specific properties of the StorageNode
func (n *StorageNode) Cluster() Member {
	return n.Member
}

// Execute will execute the given command and produce the corresponding response
func Execute(storage store.Storage, cmd Command) Response {
	element := store.Nil
	var err error
	switch cmd.Type() {
	case Put:
		err = storage.Put(cmd.Element())
	case Get:
		element, err = storage.Get(cmd.Element().Key)
	}
	return Response{
		Element: element,
		Err:     err,
	}
}

// Storage interface

// Put writes an element to the Storage
func (n *StorageNode) Put(element store.Element) error {
	return n.processor.Store.Put(element)
}

// Get retrieves an element from the Storage
func (n *StorageNode) Get(key store.Key) (store.Element, error) {
	return n.processor.Store.Get(key)
}

// Metadata returns the metadata for the StorageNode
func (n *StorageNode) Metadata() store.Metadata {
	return n.processor.Store.Metadata()
}

// Close shuts down the internal Storage implementation of the StorageNode
func (n *StorageNode) Close() error {
	return n.processor.Store.Close()
}
