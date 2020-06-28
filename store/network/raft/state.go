package raft

import (
	"fmt"

	"github.com/drakos74/lachesis/store/network"
)

type State struct {
	term      int
	index     int
	cmd       network.Command
	committed bool
}

type stateMachine struct {
	commitIndex int
	states      []*State
}

func newStateMachine() *stateMachine {
	return &stateMachine{states: make([]*State, 0)}
}

func (sm *stateMachine) verify(heartBeat HeartBeat) error {

	if len(sm.states) == 0 {
		// accept everything
		return nil
	}
	state := sm.states[len(sm.states)-1]

	if state.term > heartBeat.term {
		return fmt.Errorf("stale term received '%v' compared to current '%v'", heartBeat.term, state.term)
	}

	var index int
	for i := len(sm.states) - 1; i >= 0; i-- {
		st := sm.states[i]
		// find previous log entry in our own state machine
		if st.index == heartBeat.prevLogIndex &&
			st.term == heartBeat.term {
			index = i
			break
		}
	}

	// delete any preceding entries
	sm.states = sm.states[0 : index+1]

	return nil
}

func (sm *stateMachine) append(cmd AppendRPC) {

	sm.states = append(sm.states, &State{
		term:      cmd.term,
		index:     cmd.logIndex,
		cmd:       cmd.command,
		committed: false,
	})

}
