package raft

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/drakos74/lachesis/store/test"

	"github.com/drakos74/lachesis/store/mem"

	"github.com/drakos74/lachesis/store/network"
)

func singleNode() *network.StorageNode {
	return network.SingleNode(mem.CacheFactory, network.NoProtocol)
}

func TestRaftLeaderFollower_Append(t *testing.T) {

	leader := NewNode(singleNode)

	followers := make([]*Node, 10)

	for i := 0; i < 10; i++ {
		followers[i] = NewNode(singleNode)
	}

	for index := 0; index < 10; index++ {
		var count int
		cmd := newAppendCommand(leader, test.Random(10, 100).ElementFactory)
		// start appending to the log
		for _, follower := range followers {
			resp := follower.append(cmd)
			if resp.response.Err == nil {
				count++
			}
		}
		if count > 10/2 {
			leader.append(cmd)
			for _, follower := range followers {
				_ = follower.commit(cmd.HeartBeat)
			}
			leader.commit(cmd.HeartBeat)
		}
	}

	// check that followers state  machine is committed
	for _, follower := range followers {
		for _, state := range follower.stateMachine.states {
			assert.True(t, state.committed)
		}
	}

}

func newAppendCommand(leader *Node, newElement test.ElementFactory) AppendRPC {
	return AppendRPC{
		HeartBeat: HeartBeat{
			Epoch: Epoch{
				leaderID: leader.ID,
			},
			Log: Log{
				prevLogIndex: leader.commitIndex - 1,
				logIndex:     leader.commitIndex,
			},
		},
		command: network.NewPut(newElement()),
	}
}
