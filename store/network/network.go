package network

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/internal/partition"

	"github.com/drakos74/lachesis/store"
)

const (
	eventInterval = 100
)

type Operation struct {
	in  chan Command
	out chan Response
}

type Meta struct {
	out chan store.Metadata
	in  chan struct{}
}

type Operations []Operation

type Metadata []Meta

type WorldClock struct {
	tick      chan struct{}
	tock      chan Event
	eventPool *EventRotation
	cycles    int
}

func (wc WorldClock) startTicking() {
	for range wc.tick {
		wc.cycles++
		// TODO : fix the abstraction
		// leave some time to warm up, and use the same amount to move to the next events
		if wc.cycles > wc.eventPool.warmUp {
			idx := wc.eventPool.index
			if idx < len(wc.eventPool.events) {
				// TODO : track differently
				log.Info().
					Str("Type", "EVENT").
					Msg(fmt.Sprintf("apply new event at %d - %d = %v", wc.cycles, idx, wc.eventPool.events[idx]))
				event := wc.eventPool.events[idx]
				wc.tock <- event
				wc.eventPool.index++
			}
			wc.cycles = 0
		}
	}
}

type Network struct {
	partition.Switch
	WorldClock
	nodes []*StorageNode
	cnl   func()
}

type NetworkFactory struct {
	router      partition.PartitionStrategy
	storage     store.StorageFactory
	node        NodeFactory
	protocol    Protocol
	parallelism int
	events      []Event
}

func Factory(events ...Event) *NetworkFactory {
	return &NetworkFactory{
		events:   events,
		protocol: NoProtocol,
		node:     SingleNode,
	}
}

func (f *NetworkFactory) Storage(storage store.StorageFactory) *NetworkFactory {
	f.storage = storage
	return f
}

func (f *NetworkFactory) Nodes(parallelism int) *NetworkFactory {
	f.parallelism = parallelism
	return f
}

func (f *NetworkFactory) Router(router partition.PartitionStrategy) *NetworkFactory {
	f.router = router
	return f
}

func (f *NetworkFactory) Protocol(protocol Protocol) *NetworkFactory {
	f.protocol = protocol
	return f
}

func (f *NetworkFactory) validate() {
	if f.parallelism == 0 {
		panic("cannot create network without amount of parallelism")
	}

	if f.storage == nil {
		panic("cannot create network without a storage implementation")
	}

	if f.router == nil {
		panic("cannot create network without a routing implementation")
	}

	if f.protocol == nil {
		panic("cannot create network without a cluster protocol")
	}

	if f.node == nil {
		panic("cannot create network without a node implementation")
	}
}

func (f *NetworkFactory) Create() store.StorageFactory {
	f.validate()

	return func() store.Storage {

		ctx, cnl := context.WithCancel(context.Background())

		route := f.router()

		nodes := make([]*StorageNode, 0)

		for i := 0; i < f.parallelism; i++ {
			node := f.node(f.storage, f.protocol)
			err := node.start(ctx)
			if err == nil {
				// TODO : register with nodeId
				// register node to the network interface
				route.Register(len(nodes))
				nodes = append(nodes, node)
				// emulate the node internal protocol communication layer
				go func() {
					for msg := range node.Internal.out {
						for _, n := range nodes {
							if n.ID == msg.routingId {
								n.Internal.in <- msg
							}
						}
					}
				}()
			} else {
				// TODO : what can we do here ???
			}
		}

		net := &Network{
			Switch: route,
			WorldClock: WorldClock{
				tick: make(chan struct{}),
				tock: make(chan Event),
				eventPool: &EventRotation{
					warmUp: eventInterval,
					events: f.events,
				},
			},
			nodes: nodes,
			cnl:   cnl,
		}

		// listen to the world clock for external events
		go func() {
			for ev := range net.tock {
				net.trigger(ev)
			}
		}()

		// start the world clock
		go net.WorldClock.startTicking()

		return net

	}

}

func (n *Network) trigger(event Event) {
	if ev, ok := n.Switch.(Event); ok {
		// get back the initial router implementation
		n.Switch = ev.Reset()
	}
	// wrap with the new one
	n.Switch = event.Wrap(n.Switch)
}

func (n *Network) Put(element store.Element) error {

	cmd := PutCommand{element: element}
	// emulate a network retry mechanism
	// i.e. to capture cases where node is down
	ids, err := retry(10, n.Route, cmd.Element().Key)
	if err != nil {
		return fmt.Errorf("error during put action: %w", err)
	}

	var response Response
	for _, id := range ids {
		// we emulate for now blocking communication
		n.nodes[id].Operation.in <- cmd
		nodeResponse := <-n.nodes[id].Operation.out
		if nodeResponse.Err == nil {
			// pick the non-failing response to send to the client
			// TODO : investigate also the fail-fast approach by DeRegistering parallelism
			response = nodeResponse
		} else {
			// TODO : track these events differently
			log.Info().Str("Type", "ERROR").Msg(fmt.Sprintf("node %d returned an error = %v", id, err))
		}
	}

	n.WorldClock.tick <- struct{}{}

	return response.Err

}

func (n *Network) Get(key store.Key) (store.Element, error) {
	cmd := GetCommand{key: key}
	// emulate a network retry mechanism
	ids, err := retry(10, n.Route, cmd.Element().Key)
	if err != nil {
		return store.Nil, fmt.Errorf("error during get action: %w", err)
	}

	var response Response
	for _, id := range ids {
		// we emulate for now blocking communication
		n.nodes[id].Operation.in <- cmd
		response = <-n.nodes[id].Operation.out

		// stop at the first successful response
		if response.Err == nil {
			break
		}
	}

	n.WorldClock.tick <- struct{}{}

	return response.Element, response.Err
}

func retry(iterations int, apply func(key partition.Key) ([]int, error), key []byte) ([]int, error) {
	ids := make([]int, 0)
	err := errors.New("")
	for i := 0; i < iterations; i++ {
		ids, err = apply(key)
		if err == nil {
			break
		}
	}
	return ids, err
}

func (n *Network) Metadata() store.Metadata {

	metadata := store.Metadata{}

	// keep track of our distribution factor

	counts := make([]float64, len(n.nodes))

	for i, m := range n.nodes {
		m.Meta.in <- struct{}{}
		meta := <-m.Meta.out
		// TODO : expose differently
		log.Info().Str("Type", "META").Msg(fmt.Sprintf("%v meta = %v", i, meta))
		metadata.Merge(meta)
		counts[i] = float64(meta.Size)
	}

	std(counts)

	return metadata
}

func std(num []float64) {
	size := float64(len(num))
	var sum, mean, sd, md, dv float64
	for i := range num {
		sum += num[i]
	}
	mean = sum / size

	for j := 0; j < len(num); j++ {
		// The use of Pow math function func Pow(x, y float64) float64
		sd += math.Pow(num[j]-mean, 2)
		dv += math.Abs(num[j] - mean)
	}
	// The use of Sqrt math function func Sqrt(x float64) float64
	sd = math.Sqrt(sd / size)
	md = dv / size

	// TODO : formalize distribution metric
	log.Info().
		Str("Type", "META").
		Float64("std", sd).
		Float64("mtd", md).
		Msg(fmt.Sprintf("distribution-metric = %.2f", md/(sd+0.0001)))
}

func (n *Network) Close() error {
	n.cnl()
	return nil
}
