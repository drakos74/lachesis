package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"gioui.org/layout"
	"github.com/drakos74/oremi"
	"github.com/drakos74/oremi/bench"
)

func main() {

	file := flag.String("file", "store/benchmark/results/benchmark_hashtable.txt", "bench output file")

	flag.Parse()

	println(fmt.Sprintf("parsing benchmark file = %v", *file))

	benchmarks, err := bench.New(*file)

	if err != nil {
		log.Fatalf("could not parse benchamrks from file '%s': %v", *file, err)
	}

	oremi.Draw("benchmarks", layout.Horizontal, 1400, 800, gatherBenchmarks(benchmarks))

}

func gatherBenchmarks(benchmarks bench.Benchmarks) map[string]map[string]oremi.Collection {

	graphs := make(map[string]bench.Benchmarks)
	colors := bench.Palette(10)

	for _, b := range benchmarks {
		label := b.Labels()[0]
		i := strings.Index(label, "/")
		l := label[0:i]
		if _, ok := graphs[l]; !ok {
			graphs[l] = make([]bench.Benchmark, 0)
		}
		graphs[l] = append(graphs[l], b)
	}

	collections := make(map[string]map[string]oremi.Collection)
	collections["latency"] = make(map[string]oremi.Collection)
	//collections["memory"] = make(map[string]oremi.Collection)
	for label, graph := range graphs {
		collections["latency"][label] = graph.Extract(
			bench.Key,
			bench.Latency,
			//bench.Include(map[string]float64{bench.Num: 1000}),
			//bench.Exclude(map[string]float64{bench.Key: 16}),
		).Color(colors.Get(label))
		//collections["latency"][label] = graph.Extract(bench.Value, bench.Latency, bench.Include(map[string]float64{bench.Num: 1000}), bench.Exclude(map[string]float64{bench.Value: 100})).
		//	Color(colors.Get(label))
		//collections["memory"][label] = graph.Extract(bench.Heap, bench.Throughput).
		//	Color(colors.Get(label))
	}
	return collections
}
