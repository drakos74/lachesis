package network

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/drakos74/lachesis/store"
)

type Conn struct {
	ports   []chan Command
	reports []chan Response
}

type Comm struct {
	meta     []chan store.Metadata
	metadata []chan struct{}
}

type Network struct {
	Switch
	Conn
	Comm
	eventPool *EventRotation
	cnl       func()
	cycles    int
}

type NetworkFactory struct {
	router  PartitionStrategy
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

func (f *NetworkFactory) Router(router PartitionStrategy) *NetworkFactory {
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

func networkFactory(parallelism int, router PartitionStrategy, newStore store.StorageFactory, events ...Event) store.StorageFactory {

	return func() store.Storage {

		ctx, cnl := context.WithCancel(context.Background())

		partition := router()

		ports := make([]chan Command, parallelism)
		reports := make([]chan Response, parallelism)

		meta := make([]chan store.Metadata, parallelism)
		metadata := make([]chan struct{}, parallelism)

		for i := 0; i < parallelism; i++ {
			port := make(chan Command)
			report := make(chan Response)

			n := make(chan struct{})
			m := make(chan store.Metadata)

			// start up node
			go func() {

				storage := newStore()

				for {
					select {
					case c := <-port:
						element, err := c.Exec()(storage)
						report <- Response{
							Element: element,
							Err:     err,
						}
					case <-n:
						m <- storage.Metadata()
					case <-ctx.Done():
						return
					}
				}

			}()

			ports[i] = port
			reports[i] = report

			meta[i] = m
			metadata[i] = n

			partition.Register(i)

		}

		return &Network{
			Switch: partition,
			Comm: Comm{
				meta:     meta,
				metadata: metadata,
			},
			Conn: Conn{
				ports:   ports,
				reports: reports,
			},
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
			println(fmt.Sprintf("apply new event at %d - %d = %v", n.cycles, idx, n.eventPool.events[idx]))
			event := n.eventPool.events[idx]
			n.trigger(event)
			n.eventPool.index++
		}
	}
}

func (n *Network) Put(element store.Element) error {

	cmd := PutCommand{element: element}

	// emulate a network retry mechanism
	ids, err := retry(10, n.Route, cmd)
	if err != nil {
		return fmt.Errorf("error during put action: %w", err)
	}

	var response Response
	for _, id := range ids {
		// we emulate for now blocking communication
		n.ports[id] <- cmd
		response = <-n.reports[id]
	}

	n.cycle()

	return response.Err

}

func (n *Network) Get(key store.Key) (store.Element, error) {
	cmd := GetCommand{key: key}

	// emulate a network retry mechanism
	ids, err := retry(10, n.Route, cmd)
	if err != nil {
		return store.Nil, fmt.Errorf("error during get action: %w", err)
	}

	var response Response
	for _, id := range ids {
		// we emulate for now blocking communication
		n.ports[id] <- cmd
		response = <-n.reports[id]
		if response.Err == nil {
			break
		}
	}

	n.cycle()

	return response.Element, response.Err
}

func retry(iterations int, apply func(cmd Command) ([]int, error), cmd Command) ([]int, error) {
	ids := make([]int, 0)
	err := errors.New("")
	for i := 0; i < iterations; i++ {
		ids, err = apply(cmd)
		if err == nil {
			break
		}
	}
	return ids, err
}

func (n *Network) Metadata() store.Metadata {

	metadata := store.Metadata{}

	// keep track of our distribution factor

	counts := make([]float64, len(n.metadata))

	for i, m := range n.metadata {
		m <- struct{}{}
		meta := <-n.meta[i]
		// TODO : expose differently
		println(fmt.Sprintf("%v meta = %v", i, meta))
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
	fmt.Println("std : ", sd)
	fmt.Println("mtd : ", md)

	// TODO : formalize distribution metric
	println(fmt.Sprintf("distribution-metric = %.2f", md/(sd+0.0001)))

}

func (n *Network) Close() error {
	n.cnl()
	return nil
}
