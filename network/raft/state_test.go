package raft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateVerifyAppend(t *testing.T) {

	machine := newStateMachine()

	term := 10
	for i := 0; i < 10; i++ {

		cmd := AppendRPC{
			HeartBeat: HeartBeat{
				Epoch: Epoch{
					term: term,
				},
				Log: Log{
					prevLogIndex: i - 1,
					prevLogTerm:  9,
					logIndex:     i,
				},
			},
		}

		err := machine.verify(cmd.HeartBeat)
		assert.NoError(t, err)

		machine.append(cmd)
	}

	assert.Equal(t, 10, len(machine.states))

	// now try to roll back some states
	err := machine.verify(HeartBeat{
		Epoch: Epoch{
			term: term,
		},
		Log: Log{
			prevLogIndex: 5,
			prevLogTerm:  9,
			logIndex:     5,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 6, len(machine.states))

	assert.Equal(t, machine.states[len(machine.states)-1], &State{
		term:  10,
		index: 5,
	})

}

func TestStateMachine_VerifyOverflow(t *testing.T) {

	machine := newStateMachine()

	preFillStates(machine, 2, 3, 4, 5, 6, 7)

	// Note : newIndex is irrelevant for now
	err := machine.verify(prevState(5))
	assert.NoError(t, err)

	assert.Equal(t, 4, len(machine.states))

}

func prevState(prevIndex int) HeartBeat {
	return HeartBeat{
		Epoch: Epoch{
			term: 10,
		},
		Log: Log{
			prevLogIndex: prevIndex,
			prevLogTerm:  10,
		},
	}
}

func newState(prevIndex, newIndex int) HeartBeat {
	return HeartBeat{
		Epoch: Epoch{
			term: 10,
		},
		Log: Log{
			prevLogIndex: prevIndex,
			prevLogTerm:  10,
			logIndex:     newIndex,
		},
	}
}

func preFillStates(machine *stateMachine, indexes ...int) {

	for _, index := range indexes {
		cmd := AppendRPC{
			HeartBeat: HeartBeat{
				Epoch: Epoch{
					term: 10,
				},
				Log: Log{
					logIndex: index,
				},
			},
		}
		machine.append(cmd)
	}

}
