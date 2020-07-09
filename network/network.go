package network

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/store"
)

const (
	eventInterval = 100
)

// Operation is the communication channel for outside clients with the network
type Operation struct {
	in  chan Command
	out chan Response
}

// Meta is the communication channel for metadata about the network
type Meta struct {
	out chan store.Metadata
	in  chan struct{}
}

// Operations is a collection of Operation
type Operations []Operation

// Metadata is a collection of Meta
type Metadata []Meta

// Network emulates a cluster of nodes and the corresponding operational components of it
// the network has a Storage functionality as a whole and can be viewed as a distributed Storage
type Network struct {
	Switch
	WorldClock
	nodes []Storage
	cnl   func()
}

// FactoryBuilder is the builder factory for a network
type FactoryBuilder struct {
	router      PartitionStrategy
	storage     store.StorageFactory
	nodeFactory NodeFactory
	protocol    ProtocolFactory
	parallelism int
	events      []Event
}

// Factory creates a new NodeFactory
func Factory(events ...Event) *FactoryBuilder {
	return &FactoryBuilder{
		events:      events,
		protocol:    NoProtocol,
		nodeFactory: Node,
	}
}

// Storage specifies the underlying Storage implementation for the network
func (f *FactoryBuilder) Storage(storage store.StorageFactory) *FactoryBuilder {
	f.storage = storage
	return f
}

// Nodes specifies the amount of nodes for the distributed network
func (f *FactoryBuilder) Nodes(parallelism int) *FactoryBuilder {
	f.parallelism = parallelism
	return f
}

// Router specifies the routing / switch implementation to route external client requests to individual nodes
func (f *FactoryBuilder) Router(router PartitionStrategy) *FactoryBuilder {
	f.router = router
	return f
}

// Protocol specifies the internal communication protocol for the network
// i.e. how it's members communicate with each other
func (f *FactoryBuilder) Protocol(protocol ProtocolFactory) *FactoryBuilder {
	f.protocol = protocol
	return f
}

// Node specifies the factory for creating new network members
func (f *FactoryBuilder) Node(nodeFactory NodeFactory) *FactoryBuilder {
	f.nodeFactory = nodeFactory
	return f
}

func (f *FactoryBuilder) validate() {
	if f.parallelism == 0 {
		panic("cannot create network without amount of parallelism")
	}

	if f.storage == nil {
		panic("cannot create network without a Storage implementation")
	}

	if f.router == nil {
		panic("cannot create network without a routing implementation")
	}

	if f.protocol == nil {
		panic("cannot create network without a cluster protocol")
	}

	if f.nodeFactory == nil {
		panic("cannot create network without a node implementation")
	}
}

// Create returns a functioning network implementation that can be used as a distributed Storage
func (f *FactoryBuilder) Create() store.StorageFactory {
	f.validate()

	return func() store.Storage {
		ctx, cnl := context.WithCancel(context.Background())
		// create the network router
		route := f.router()

		nodes := make([]Storage, 0)
		// apply properties to the members of the network
		for i := 0; i < f.parallelism; i++ {
			node := f.nodeFactory(f.storage, f.protocol)

			// listen to internal cluster events
			go func() {
				for {
					select {
					case msg := <-node.Cluster().Internal.in:
						node.Cluster().Internal.out <- node.Cluster().Internal.Process(len(route.Members()), node, msg)
					case <-ctx.Done():
						log.Debug().Msg("Closing member channel")
						return
					}
				}
			}()

			// listen to client events
			go func() {
				for {
					select {
					case cmd := <-node.Cluster().Operation.in:
						element := store.Nil
						var err error
						switch cmd.Type() {
						case Put:
							err = node.Put(cmd.Element())
						case Get:
							element, err = node.Get(cmd.Element().Key)
						}
						node.Cluster().Operation.out <- Response{
							Element: element,
							Err:     err,
						}
					case <-node.Cluster().Meta.in:
						node.Cluster().Meta.out <- node.Metadata()
					case <-ctx.Done():
						log.Debug().Msg("Closing Storage channel")
						return
					}
				}
			}()

			// register node to the network interface
			route.Register(len(nodes))
			nodes = append(nodes, node)

		}
		// create the network struct
		net := &Network{
			Switch: route,
			WorldClock: WorldClock{
				tick: make(chan struct{}),
				tock: make(chan Event),
				eventPool: &Events{
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

// trigger initiates an event for the network
func (n *Network) trigger(event Event) {
	if ev, ok := n.Switch.(Event); ok {
		// get back the initial router implementation
		n.Switch = ev.Reset()
	}
	// wrap with the new one
	n.Switch = event.Wrap(n.Switch)
	n.DeRegister(event.Index())
}

// Put allows clients to initPut a write command to the distributed Storage
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
		n.nodes[id].Cluster().Operation.in <- cmd
		nodeResponse := <-n.nodes[id].Cluster().Operation.out
		if nodeResponse.Err == nil {
			// pick the non-failing response to send to the client
			response = nodeResponse
		} else {
			log.Info().Str("Type", "ERROR").Msg(fmt.Sprintf("node %d returned an error = %v", id, nodeResponse.Err))
		}
	}

	n.WorldClock.tick <- struct{}{}

	return response.Err

}

// Get initiates a read requests to the distributed Storage
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
		n.nodes[id].Cluster().Operation.in <- cmd
		response = <-n.nodes[id].Cluster().Operation.out
		// stop at the first successful response
		if response.Err == nil {
			break
		}
	}

	n.WorldClock.tick <- struct{}{}

	return response.Element, response.Err
}

// retry emulates a retry mechanism, in case a node is down
func retry(iterations int, apply func(key Key) ([]int, error), key []byte) ([]int, error) {
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

// Metadata returns the network metadata
func (n *Network) Metadata() store.Metadata {

	metadata := store.Metadata{}

	// keep track of our distribution factor

	counts := make([]float64, len(n.nodes))

	for i, node := range n.nodes {
		node.Cluster().Meta.in <- struct{}{}
		meta := <-node.Cluster().Meta.out
		log.Info().Str("Type", "META").Msg(fmt.Sprintf("%v meta = %v", i, meta))
		metadata.Merge(meta)
		counts[i] = float64(meta.Size)
	}

	std(counts)

	return metadata
}

// std computes the standard deviation of a population of floats
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

// Close runs any cleanup / shutdown actions on the Storage
func (n *Network) Close() error {
	n.cnl()
	return nil
}
