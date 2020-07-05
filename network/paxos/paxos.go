package paxos

import (
	"fmt"
	"reflect"
	"time"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/network"
)

// Protocol implements the internal cluster communication requirements,
// e.g. the proposers and acceptors interaction logic
func Protocol() network.ProtocolFactory {

	processor := network.ProcessorFactory(func(state *network.State, node *network.StorageNode, element store.Element) (rpc interface{}, wait bool) {
		index := time.Now().UnixNano()
		return Proposal{
			index:   index,
			command: network.NewPut(element),
		}, true
	})

	// follower phase 1 processing logic
	processor.Propose(func(state *network.State, storage store.Storage, msg interface{}) (interface{}, error) {
		// we expect a proposal message
		if proposal, ok := msg.(Proposal); ok {
			if currentState, ok := state.Log[string(proposal.command.Element().Key)]; ok {
				entry := currentState.(Proposal)
				if entry.index > proposal.index {
					return nil, fmt.Errorf("a higher index '%v' already exists for key '%v'",
						entry.index, proposal.command.Element().Key)
				}
			}
			state.Log[string(proposal.command.Element().Key)] = proposal
			response := Promise{}
			return response, nil
		}
		return nil, fmt.Errorf("unexpected message received for proposal confirmation '%v'", reflect.TypeOf(msg))
	})

	// leader phase1 processing logic
	processor.Promise(func(state *network.State, storage store.Storage, msg interface{}) (interface{}, error) {
		// we are doing the same work as the follower in the previous step
		// e.g. appending to our log the same as the follower did, so that we are fully aligned!
		if proposal, ok := msg.(Proposal); ok {
			if currentState, ok := state.Log[string(proposal.command.Element().Key)]; ok {
				entry := currentState.(Proposal)
				if entry.index > proposal.index {
					return nil, fmt.Errorf("a higher index '%v' already exists for key '%v'",
						entry.index, proposal.command.Element().Key)
				}
			}
			state.Log[string(proposal.command.Element().Key)] = proposal
			response := Commit{key: proposal.command.Element().Key}
			return response, nil
		}
		return nil, fmt.Errorf("unexpected message to create commit '%v'", reflect.TypeOf(msg))
	})

	// follower phase 2 processing logic
	processor.Commit(func(state *network.State, storage store.Storage, msg interface{}) (interface{}, error) {
		if commit, ok := msg.(Commit); ok {
			currentState := state.Log[string(commit.key)]
			entry := currentState.(Proposal)
			network.Execute(storage, entry.command)
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected message received for commit confirmation '%v'", reflect.TypeOf(msg))
	})

	processor.Confirm(func(state *network.State, storage store.Storage, msg interface{}) (interface{}, error) {
		if commit, ok := msg.(Commit); ok {
			currentState := state.Log[string(commit.key)]
			entry := currentState.(Proposal)
			network.Execute(storage, entry.command)
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected message received for commit confirmation '%v'", reflect.TypeOf(msg))
	})

	return network.ConsensusProtocol(*processor)
}
