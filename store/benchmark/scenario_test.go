package benchmark

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnceOffScenario(t *testing.T) {

	scenario := Benchmark(limit(1), 2, 10, 50)

	benchmarks := make([]Config, 0)
	scenario.execute(func(scenario Config) {
		benchmarks = append(benchmarks, scenario)
	})

	assert.Equal(t, 1, len(benchmarks))

	assertBenchmarkScenario(t, 2, 10, 50, benchmarks[0])

}

func TestSimpleEvolutionScenario(t *testing.T) {

	scenario := Benchmark(Evolution().
		add(limit(10)).
		add(num(add(1))).
		create(),
		2, 10, 50)

	benchmarks := make([]Config, 0)
	scenario.execute(func(scenario Config) {
		benchmarks = append(benchmarks, scenario)
	})

	assert.Equal(t, 10, len(benchmarks))

	assertBenchmarkScenario(t, 2, 10, 50, benchmarks[0])
	assertBenchmarkScenario(t, 11, 10, 50, benchmarks[9])

}

func TestSinglePropertyEvolutionScenario(t *testing.T) {

	scenario := Benchmark(Evolution().
		add(limit(5)).
		add(num(pow(5))).
		create(),
		2, 10, 50)

	benchmarks := make([]Config, 0)
	scenario.execute(func(scenario Config) {
		benchmarks = append(benchmarks, scenario)
	})

	assert.Equal(t, 5, len(benchmarks))

	assertBenchmarkScenario(t, 2, 10, 50, benchmarks[0])
	assertBenchmarkScenario(t, 2*int(math.Pow(5, 4)), 10, 50, benchmarks[4])

}

func TestComplexEvolutionScenario(t *testing.T) {

	scenario := Benchmark(Evolution().
		add(limit(5)).
		add(num(pow(10))).
		add(key(add(5))).
		add(value(pow(2))).
		create(),
		2, 2, 2)

	benchmarks := make([]Config, 0)
	scenario.execute(func(scenario Config) {
		benchmarks = append(benchmarks, scenario)
	})

	assert.Equal(t, 5, len(benchmarks))

	assertBenchmarkScenario(t, 2, 2, 2, benchmarks[0])
	assertBenchmarkScenario(t, 2*int(math.Pow(10, 4)), 2+(5*4), int(math.Pow(2, 5)), benchmarks[4])

}

func assertBenchmarkScenario(t *testing.T, num, keySize, valueSize int, singleBenchmark Config) {
	assert.Equal(t, num, singleBenchmark.Num)
	assert.Equal(t, keySize, singleBenchmark.KeySize)
	assert.Equal(t, valueSize, singleBenchmark.ValueSize)
}
