package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	lachesisbench "github.com/drakos74/lachesis/store/benchmark"

	"gioui.org/layout"
	"github.com/drakos74/oremi"
	"github.com/drakos74/oremi/bench"
)

func main() {

	file := flag.String("file", "store/benchmark/results/b.txt", "bench output file")

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
	for label, benchmark := range graphs {
		collections["latency"][label] = benchmark.Extract(
			bench.Value,
			bench.Latency,
			bench.Label(lachesisbench.ValueSizeScenario),
		).Color(colors.Get(label))
	}
	return collections
}
