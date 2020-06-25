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

type Port struct {
	in  chan Command
	out chan Response
}

type Meta struct {
	out chan store.Metadata
	in  chan struct{}
}

type Conn struct {
	ports []Port
}

type Comm struct {
	metadata []Meta
}

type Network struct {
	partition.Switch
	nodes     []*Node
	eventPool *EventRotation
	cnl       func()
	cycles    int
}

type NetworkFactory struct {
	router  partition.PartitionStrategy
	storage store.StorageFactory
	nodes   int
	events  []Event
}

func Factory() *NetworkFactory {
	return &NetworkFactory{}
}

func (f *NetworkFactory) Storage(storage store.StorageFactory) *NetworkFactory {
	f.storage = storage
	return f
}

func (f *NetworkFactory) Nodes(nodes int) *NetworkFactory {
	f.nodes = nodes
	return f
}

func (f *NetworkFactory) Router(router partition.PartitionStrategy) *NetworkFactory {
	f.router = router
	return f
}

func (f *NetworkFactory) Events(events ...Event) *NetworkFactory {
	f.events = events
	return f
}

func (f *NetworkFactory) Create() store.StorageFactory {
	if f.nodes == 0 {
		panic("cannot create network without amount of nodes")
	}

	if f.storage == nil {
		panic("cannot create network without a storage implementation")
	}

	if f.router == nil {
		panic("cannot create network without a routing implementation")
	}

	return networkFactory(f.nodes, f.router, f.storage, f.events...)
}

func networkFactory(parallelism int, router partition.PartitionStrategy, newStore store.StorageFactory, events ...Event) store.StorageFactory {

	return func() store.Storage {

		ctx, cnl := context.WithCancel(context.Background())

		route := router()

		nodes := make([]*Node, 0)

		for i := 0; i < parallelism; i++ {
			node := NewNode(newStore)
			err := node.start(ctx)
			if err == nil {
				// TODO : register with nodeId
				route.Register(len(nodes))
				nodes = append(nodes, node)
			} else {
				// TODO : what can we do here ???
			}
		}

		return &Network{
			Switch: route,
			nodes:  nodes,
			eventPool: &EventRotation{
				events: events,
			},
			cnl: cnl,
		}

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

func (n *Network) cycle() {
	n.cycles++
	// leave some time to
	if n.cycles%100 == 0 {
		idx := n.eventPool.index
		if idx < len(n.eventPool.events) {
			// TODO : track differently
			log.Info().
				Str("Type", "EVENT").
				Msg(fmt.Sprintf("apply new event at %d - %d = %v", n.cycles, idx, n.eventPool.events[idx]))
			event := n.eventPool.events[idx]
			n.trigger(event)
			n.eventPool.index++
		}
	}
}

func (n *Network) Put(element store.Element) error {

	cmd := PutCommand{element: element}
	// emulate a network retry mechanism
	ids, err := retry(10, n.Route, cmd.Element().Key)
	if err != nil {
		return fmt.Errorf("error during put action: %w", err)
	}

	var response Response
	for _, id := range ids {
		// we emulate for now blocking communication
		n.nodes[id].Port.in <- cmd
		nodeResponse := <-n.nodes[id].Port.out
		if nodeResponse.Err == nil {
			// pick the non-failing response to send to the client
			// TODO : investigate also the fail-fast approach by DeRegistering nodes
			response = nodeResponse
		} else {
			// TODO : track these events differently
			log.Info().Str("Type", "ERROR").Msg(fmt.Sprintf("node %d returned an error = %v", id, err))
		}
	}

	n.cycle()

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
		n.nodes[id].Port.in <- cmd
		response = <-n.nodes[id].Port.out

		// stop at the first successful response
		if response.Err == nil {
			break
		}
	}

	n.cycle()

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
