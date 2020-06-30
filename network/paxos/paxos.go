package paxos

import (
	"fmt"

	"github.com/drakos74/lachesis/network"
	"github.com/rs/zerolog/log"
)

const consensusThreshold = 1

type confirmation struct {
	count   map[uint32]int
	trigger map[uint32]bool
}

func (c confirmation) reached(id uint32) bool {
	return c.count[id] < 0 && !c.trigger[id]
}

// PaxosProtocol implements the internal cluster communication requirements,
// e.g. the proposers and acceptors communication
func PaxosProtocol(group chan Signal) network.ProtocolFactory {
	return func(id uint32) network.Internal {
		// keep some local cache for the leader to count the responses
		consensus := confirmation{
			count:   make(map[uint32]int),
			trigger: make(map[uint32]bool),
		}

		return network.Protocol(id, func(members int, node network.Storage, msg network.Message) network.Message {

			member, _ := node.(*Node)

			// acceptor reactions

			if proposal, ok := msg.Content.(Proposal); ok {
				log.Debug().
					Str("type", "acceptor").
					Str("rpc", "proposal").
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received rpc")
				return network.Message{
					ID:        msg.ID,
					Source:    member.Cluster().ID,
					RoutingId: msg.Source,
					Content:   member.promise(proposal),
				}
			}

			if commit, ok := msg.Content.(Commit); ok {
				log.Debug().
					Str("type", "acceptor").
					Str("rpc", "commit").
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received rpc")
				return network.Message{
					ID:        msg.ID,
					Source:    member.Cluster().ID,
					RoutingId: msg.Source,
					Content:   member.accept(commit),
				}
			}

			// proposer reactions

			if promise, ok := msg.Content.(Promise); ok {
				log.Debug().
					Str("type", "proposer").
					Str("error", fmt.Sprintf("%v ", promise.err)).
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received promise")
				if _, ok := consensus.count[msg.ID]; !ok {
					consensus.count[msg.ID] = int(float64(members-2) * consensusThreshold)
					consensus.trigger[msg.ID] = false
				}

				if promise.err == nil {
					consensus.count[msg.ID]--
					if consensus.reached(msg.ID) {
						log.Debug().
							Str("type", "leader").
							Uint32("id", msg.ID).
							Uint32("node", member.Cluster().ID).
							Msg("received signal")
						group <- Prepare
						consensus.trigger[msg.ID] = true
					}
				}
			}

			if response, ok := msg.Content.(network.Response); ok {
				log.Debug().
					Str("type", "proposer").
					Str("error", fmt.Sprintf("%v ", response.Err)).
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received promise")
				if _, ok := consensus.count[msg.ID]; !ok {
					consensus.count[msg.ID] = int(float64(members-2) * consensusThreshold)
					consensus.trigger[msg.ID] = false
				}

				if response.Err == nil {
					consensus.count[msg.ID]--
					if consensus.reached(msg.ID) {
						log.Debug().
							Str("type", "leader").
							Uint32("id", msg.ID).
							Uint32("node", member.Cluster().ID).
							Msg("received accept")
						group <- Accept
						consensus.trigger[msg.ID] = true
					}
				}
			}
			// dont trigger any other events
			return network.Void
		})
	}
}
