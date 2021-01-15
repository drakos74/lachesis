package network

import (
	"fmt"
	"time"

	"github.com/drakos74/lachesis/internal/app/store"

	"github.com/rs/zerolog/log"
)

// Internal defines the communication channels for inter-node communication within the network
// it acts as a central switch that routes traffic to the nodes, based on the message metadata
type Internal struct {
	ID      uint32
	in      chan Message
	out     chan Message
	signal  chan Signal
	Process MsgRouter
}

// Send propagates a message to the appropriate communication channel for inter-node communication
func (i *Internal) Send(msg Message) {
	i.out <- msg
}

func (i *Internal) processSignal(signal Signal, process func() error) error {
	select {
	case s := <-i.signal:
		if s == signal {
			return process()
		}
		return fmt.Errorf("unexpected signal received: '%v' instead of '%v'", s, signal)
	case <-time.Tick(5 * time.Second):
		return fmt.Errorf("could not get consensus from cluster")
	}
}

// buffer defines the buffering capabilities of the network for inter-node communication
const buffer = 0

// Protocol creates the internal communication pipeline by providing the processing instructions
// for how to handle incoming messages
func Protocol(id uint32, signal chan Signal, process MsgRouter, processor Processor) (*Internal, *Peer) {
	return &Internal{
			ID:      id,
			in:      make(chan Message, buffer),
			out:     make(chan Message, buffer),
			signal:  signal,
			Process: process,
		}, &Peer{
			processor: processor,
		}
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

const consensusThreshold = 1

type confirmation struct {
	msgType map[uint32]MsgType
	count   map[uint32]int
	trigger map[uint32]bool
}

func (c *confirmation) reached(id uint32) bool {
	if c.count[id] < 0 && !c.trigger[id] {
		c.trigger[id] = true
		return true
	}
	return false
}

// ConsensusProtocol implements the internal cluster communication requirements
func ConsensusProtocol(processor Processor) ProtocolFactory {

	return func(id uint32) (*Internal, *Peer) {
		// keep some local cache for the leader to count the responses
		consensus := confirmation{
			msgType: make(map[uint32]MsgType),
			count:   make(map[uint32]int),
			trigger: make(map[uint32]bool),
		}

		signal := make(chan Signal)

		return Protocol(id, signal, func(members int, node Storage, msg Message) Message {

			member, _ := node.(*StorageNode)

			switch msg.Type {
			// follower logic
			case Propose:
				fallthrough
			case Commit:
				log.Debug().
					Str("type", "follower").
					Str("rpc", msg.Type.String()).
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received rpc")
				response, err := member.getProcessor(msg.Type)(&processor.State, processor.Store, msg.Content)
				return Message{
					ID:        msg.ID,
					Source:    member.Cluster().ID,
					RoutingID: msg.Source,
					Content:   response,
					Type:      msg.Type.Next(),
					Err:       err,
				}
			// initiator logic
			case Promise:
				fallthrough
			case Confirm:
				log.Debug().
					Str("type", "leader").
					Str("rpc", msg.Type.String()).
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received rpc response")
				if _, ok := consensus.count[msg.ID]; !ok {
					consensus.count[msg.ID] = int(float64(members-2) * consensusThreshold)
					consensus.trigger[msg.ID] = false
					consensus.msgType[msg.ID] = msg.Type
				}
				if msg.Err == nil {
					consensus.count[msg.ID]--
					if consensus.reached(msg.ID) {
						log.Debug().
							Str("type", "leader").
							Str("rpc", msg.Type.String()).
							Uint32("id", msg.ID).
							Uint32("node", member.Cluster().ID).
							Msg("received signal")
						signal <- msg.Type.Phase()
						consensus.trigger[msg.ID] = true
					}
				} else {
					log.Err(msg.Err).Msg("cannot read consensus response")
					// TODO : return the error
				}
			}

			// dont trigger any other events
			return Void
		}, processor)
	}
}
