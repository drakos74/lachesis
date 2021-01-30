package raft

import (
	"fmt"
	"reflect"

	"github.com/drakos74/lachesis/benchmarks/network"
	"github.com/drakos74/lachesis/store/app/storage"
)

// Protocol implements the internal cluster communication requirements,
// e.g. the leader and followers interaction logic
func Protocol() network.ProtocolFactory {

	processor := network.ProcessorFactory(func(state *network.State, node *network.StorageNode, element storage.Element) (rpc interface{}, wait bool) {
		stMachine, err := retrieveStatMachine(state)
		if err != nil {
			return nil, false
		}
		return AppendRPC{
			HeartBeat: HeartBeat{
				Epoch: Epoch{
					leaderID: node.Cluster().ID,
				},
				Log: Log{
					prevLogIndex: stMachine.commitIndex - 1,
					logIndex:     stMachine.commitIndex,
				},
			},
			command: network.NewPut(element),
		}, true
	})

	// follower phase 1 processing logic
	processor.Propose(func(state *network.State, storage storage.Storage, msg interface{}) (interface{}, error) {
		stMachine, err := retrieveStatMachine(state)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve state machine '%w'", err)
		}
		// we expect a proposal message
		if cmd, ok := msg.(AppendRPC); ok {

			err := appendRPC(stMachine, cmd)

			if err != nil {
				return nil, err
			}

			updateStateMachine(state, stMachine)

			return ResponseRPC{
				HeartBeat: cmd.HeartBeat,
				response:  network.Response{},
			}, nil
		}
		return nil, fmt.Errorf("unexpected message received for proposal confirmation '%v'", reflect.TypeOf(msg))
	})

	// leader phase1 processing logic
	processor.Promise(func(state *network.State, storage storage.Storage, msg interface{}) (interface{}, error) {
		// we are doing the same work as the follower in the previous step
		// e.g. appending to our log the same as the follower did, so that we are fully aligned!
		// create our stateMachine if not already created
		if _, ok := state.Log[""]; !ok {
			state.Log[""] = newStateMachine()
		}

		// we expect a proposal message
		if cmd, ok := msg.(AppendRPC); ok {
			wal := state.Log[""]
			if stateMachine, ok := wal.(*stateMachine); ok {

				err := appendRPC(stateMachine, cmd)

				if err != nil {
					return nil, err
				}

				state.Log[""] = stateMachine

				return ResponseRPC{
					HeartBeat: cmd.HeartBeat,
					response:  network.Response{},
				}, nil
			}
			return nil, fmt.Errorf("could not retrieve state machine '%v'", reflect.TypeOf(wal))
		}
		return nil, fmt.Errorf("unexpected message received for proposal confirmation '%v'", reflect.TypeOf(msg))
	})

	// follower phase 2 processing logic
	processor.Commit(func(state *network.State, store storage.Storage, msg interface{}) (interface{}, error) {
		// for ease of use, we skip another verify at this stage
		// we assume the leader is 'sane' in the sense that it will adhere to the protocol
		if heartbeat, ok := msg.(ResponseRPC); ok {
			wal := state.Log[""]

			if stateMachine, ok := wal.(*stateMachine); ok {
				// get all the pending entries from the state machine and commit them to the store
				for i := state.Index; i <= heartbeat.logIndex; i++ {
					resp := network.Execute(store, stateMachine.states[i].cmd)
					if resp.Err == nil {
						stateMachine.states[i].committed = true
					}
				}

				stateMachine.commitIndex = heartbeat.logIndex + 1

				state.Log[""] = stateMachine

				return ResponseRPC{
					response: network.Response{
						Element: storage.Nil,
					},
				}, nil
			}
			return nil, fmt.Errorf("could not retrieve state machine '%v'", reflect.TypeOf(wal))
		}
		return nil, fmt.Errorf("unexpected message received for commit action '%v'", reflect.TypeOf(msg))
	})

	processor.Confirm(func(state *network.State, store storage.Storage, msg interface{}) (interface{}, error) {
		if heartbeat, ok := msg.(ResponseRPC); ok {
			wal := state.Log[""]

			if stateMachine, ok := wal.(*stateMachine); ok {
				// get all the pending entries from the state machine and commit them to the store
				for i := state.Index; i <= heartbeat.logIndex; i++ {
					resp := network.Execute(store, stateMachine.states[i].cmd)
					if resp.Err == nil {
						stateMachine.states[i].committed = true
					}
				}

				stateMachine.commitIndex = heartbeat.logIndex + 1

				state.Log[""] = stateMachine

				return ResponseRPC{
					response: network.Response{
						Element: storage.Nil,
					},
				}, nil
			}
			return nil, fmt.Errorf("could not retrieve state machine '%v'", reflect.TypeOf(wal))
		}
		return nil, fmt.Errorf("unexpected message received for commit confirmation '%v'", reflect.TypeOf(msg))
	})

	return network.ConsensusProtocol(*processor)
}

func appendRPC(stateMachine *stateMachine, cmd AppendRPC) error {
	err := stateMachine.verify(cmd.HeartBeat)

	if err != nil {
		return fmt.Errorf("inconsistent node state: %w", err)
	}

	// add the new state
	stateMachine.append(cmd)

	return nil
}

func retrieveStatMachine(state *network.State) (*stateMachine, error) {
	if _, ok := state.Log[""]; !ok {
		state.Log[""] = newStateMachine()
	}

	wal := state.Log[""]
	if stMachine, ok := wal.(*stateMachine); ok {
		return stMachine, nil
	}
	return nil, fmt.Errorf("could not verify state machine: %v", reflect.TypeOf(wal))
}

func updateStateMachine(state *network.State, stMachine *stateMachine) {
	state.Log[""] = stMachine
}
