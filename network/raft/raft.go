package raft

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/network"
)

const consensusThreshold = 0.51

type confirmation struct {
	count   map[uint32]int
	trigger map[uint32]bool
}

// RaftProtocol implements the internal cluster communication requirements,
// e.g. the leader keeping all the state machines up to date,
// so that anyone can take over if needed
func RaftProtocol(group chan Signal) network.ProtocolFactory {
	return func(id uint32) network.Internal {
		// keep some local cache for the leader to count the responses
		tmpCache := confirmation{
			count:   make(map[uint32]int),
			trigger: make(map[uint32]bool),
		}

		return network.Protocol(id, func(members int, node network.Storage, msg network.Message) network.Message {

			member, _ := node.(*Node)

			// case it s an AppendRPC message
			// this means we are a follower
			// and need to send the message back to 'where it came from'
			if cmd, ok := msg.Content.(AppendRPC); ok {
				log.Debug().
					Str("type", "follower").
					Str("rpc", "append").
					Int("index", cmd.HeartBeat.logIndex).
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received rpc")
				return network.Message{
					ID:        msg.ID,
					Source:    member.Cluster().ID,
					RoutingId: msg.Source,
					Content:   member.append(cmd),
				}
			}

			if heartbeat, ok := msg.Content.(HeartBeat); ok {
				log.Debug().
					Str("type", "follower").
					Str("rpc", "commit").
					Int("index", heartbeat.logIndex).
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received rpc")
				return network.Message{
					ID:        msg.ID,
					Source:    member.Cluster().ID,
					RoutingId: msg.Source,
					Content:   member.commit(heartbeat),
				}
			}

			// case it s an ResponseRPC message
			// this means we are the leader now
			if rsp, ok := msg.Content.(ResponseRPC); ok {
				log.Debug().
					Str("type", "leader").
					Str("error", fmt.Sprintf("%v ", rsp.response.Err)).
					Uint32("id", msg.ID).
					Uint32("node", member.Cluster().ID).
					Uint32("from", msg.Source).
					Msg("received response")
				if _, ok := tmpCache.count[msg.ID]; !ok {
					tmpCache.count[msg.ID] = members - 2
					tmpCache.trigger[msg.ID] = false
				}

				if rsp.response.Err == nil {
					count := tmpCache.count[msg.ID]
					count--
					tmpCache.count[msg.ID] = count
					if count < 0 && !tmpCache.trigger[msg.ID] {
						log.Debug().
							Str("type", "leader").
							Uint32("id", msg.ID).
							Str("signal", rsp.Signal.String()).
							Uint32("node", member.Cluster().ID).
							Msg("received signal")
						group <- rsp.Signal
						tmpCache.trigger[msg.ID] = true
					}
				}
			}
			// dont trigger any other events
			return network.Void
		})
	}
}
