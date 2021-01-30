package bench

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/drakos74/oremi"

	"github.com/drakos74/oremi/internal/data/model"
)

const (
	Operations = "ops"
	Latency    = "ns/op"
	Throughput = "B/op"
	Heap       = "allocs/op"
	Key        = "size-key"
	Value      = "size-value"
	Num        = "num-objects"
)

// Benchmarks is a collection of Benchmark results
type Benchmarks []Benchmark

// Extract extracts Latency and operation information from the given benchmarks
// x value to be used for the x-axis
// y value to be used for the y-axis
func (b Benchmarks) Extract(x, y string, filters ...Filter) *oremi.Collection {
	series := model.NewSeries(x, y)
	for _, benchmark := range b {
		x, hasX := benchmark.read(x)
		y, hasY := benchmark.read(y)

		matchIn := true
		matchOut := true
		label := true
		for _, filter := range filters {
			switch filter.Type {
			case IN:
				filter.Apply(benchmark, &matchIn)
			case OUT:
				filter.Apply(benchmark, &matchOut)
			case LABEL:
				filter.Apply(benchmark, &label)
			}
		}

		if hasX && hasY && matchIn && matchOut && label {
			series.Add(model.NewVector(benchmark.labels, x, y))
		}
	}
	return oremi.New(series)
}

// New creates a collection of Benchmark items from a given becnhmark output file
func New(f string) (Benchmarks, error) {

	var benchmarks []Benchmark

	file, err := os.Open(f)
	if err != nil {
		return benchmarks, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		b, err := tryParseBenchmark(line)
		if err == nil {
			benchmarks = append(benchmarks, *b)
		} else {
			println(fmt.Sprintf("ignoring line %s : %v", line, err))
		}
	}

	if err := scanner.Err(); err != nil {
		return benchmarks, err
	}

	return benchmarks, nil

}

type Benchmark struct {
	labels []string
	// numeric labels
	numLabels map[string]float64
}

func (b Benchmark) hasLabel(label string) bool {
	for _, l := range b.labels {
		if l == label {
			return true
		}
	}
	return false
}

func (b Benchmark) read(numLabel string) (float64, bool) {
	if a, ok := b.numLabels[numLabel]; ok {
		return a, true
	} else {
		panic(fmt.Sprintf("num-label not found: %v", numLabel))
	}
	return 0, false
}

func (b Benchmark) Latency() (float64, bool) {
	return b.read(Latency)
}

func (b Benchmark) Operations() (float64, bool) {
	return b.read(Operations)
}

func (b Benchmark) Throughput() (float64, bool) {
	return b.read(Throughput)
}

func (b Benchmark) Heap() (float64, bool) {
	return b.read(Heap)
}

func (b Benchmark) Labels() []string {
	return b.labels
}

func tryParseBenchmark(line string) (*Benchmark, error) {
	parts := strings.Fields(line)

	if isBenchmarkOutput(parts) {
		return parseBenchmark(parts), nil
	}

	return nil, fmt.Errorf("could not find bench within %v", parts)

}

func parseBenchmark(parts []string) *Benchmark {

	ops, _ := strconv.Atoi(parts[1])
	lat, _ := strconv.ParseFloat(parts[2], 64)

	labels, numLabels := parseLabels(parts[0])

	numLabels[Latency] = lat
	numLabels[Operations] = float64(ops)

	if hasAtIndex(5, Throughput, parts) {
		numLabels[Throughput], _ = strconv.ParseFloat(parts[4], 64)
	}

	if hasAtIndex(7, Heap, parts) {
		numLabels[Heap], _ = strconv.ParseFloat(parts[6], 64)
	}

	return &Benchmark{
		labels:    labels,
		numLabels: numLabels,
	}
}

func parseLabels(str string) (labels []string, numLabels map[string]float64) {

	labels = make([]string, 0)
	numLabels = make(map[string]float64)

	parts := strings.Split(str, "|")

	for _, p := range parts {
		label := strings.Split(p, ":")
		if len(label) > 1 {
			num, _ := strconv.ParseFloat(label[1], 64)
			numLabels[label[0]] = num
		} else {
			labels = append(labels, label[0])
		}

	}

	return

}

func isBenchmarkOutput(parts []string) bool {
	return hasAtIndex(3, Latency, parts)
}

func hasAtIndex(index int, part string, parts []string) bool {
	if len(parts) > index {
		return parts[index] == part
	}
	return false
}
