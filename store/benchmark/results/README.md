## Visualising Benchmarks

We are using [oremi](www.github.com/drakos74/oremi) for visualising the benchmarks

This includes  2 steps

- parse the file into benchmarks

```go
file := flag.String("file", "store/benchmark/results/benchmark_hashtable.txt", "bench output file")

	benchmarks, err := bench.New(*file)

	oremi.Draw("benchmarks", layout.Horizontal, 1400, 800, gatherBenchmarks(benchmarks))
```

- gather and filter the relevant values for the visualisation graph

```go
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
	for label, graph := range graphs {
		collections["latency"][label] = graph.
            Extract(bench.Key, 
                    bench.Latency, 
                    bench.Include(map[string]float64{bench.Num: 1000}), 
                    bench.Exclude(map[string]float64{bench.Key: 16})).
			Color(colors.Get(label))
	}
	return collections
```