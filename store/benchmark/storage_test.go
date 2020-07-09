package benchmark

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"
	"github.com/drakos74/lachesis/store/test"
)

// in-memory

func BenchmarkCache(b *testing.B) {
	executeBenchmarks(b, mem.CacheFactory)
}

func BenchmarkSyncCache(b *testing.B) {
	executeBenchmarks(b, mem.SyncCacheFactory)
}

func BenchmarkTrie(b *testing.B) {
	executeBenchmarks(b, mem.TrieFactory)
}

func BenchmarkSyncTrie(b *testing.B) {
	executeBenchmarks(b, mem.SyncTrieFactory)
}

func BenchmarkBTree(b *testing.B) {
	executeBenchmarks(b, mem.BTreeFactory)
}

func BenchmarkSyncBTree(b *testing.B) {
	executeBenchmarks(b, mem.SyncBTreeFactory)
}

func executeBenchmarks(b *testing.B, storageFactory func() store.Storage) {

	// reduce logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	scenarios := []Scenario{
		Benchmark(Evolution().
			add(limit(6)).
			add(num(pow(10))).
			create(),
			10, 10, 20),
		Benchmark(Evolution().
			add(limit(10)).
			add(key(add(4))).
			create(),
			1000, 4, 100),
		Benchmark(Evolution().
			add(limit(8)).
			add(value(pow(2))).
			create(),
			1000, 16, 2),
	}

	storage := storageFactory()
	for _, scenario := range scenarios {
		executeBenchmark(b, storage, scenario, put, get)
	}

}

func executeBenchmark(b *testing.B, storage store.Storage, scenario Scenario, execution ...benchmarkExecution) {

	for scenario.next() {

		currentScenario := scenario.get()
		elements := test.Elements(currentScenario.Num, test.Random(currentScenario.KeySize, currentScenario.ValueSize))

		for _, exec := range execution {
			b.Run(fmt.Sprintf("%s|%s|num-objects:%d|size-key:%d|size-value:%d|", reflect.TypeOf(storage).String(), getFuncName(exec), scenario.Num, scenario.KeySize, scenario.ValueSize), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					exec(storage, elements)
				}
			})
		}

	}

	err := storage.Close()
	if err != nil {
		log.Fatal()
	}
}

// benchmark utilities

type benchmarkExecution func(storage store.Storage, elements []store.Element)

var result store.Element

func put(storage store.Storage, elements []store.Element) {
	for _, element := range elements {
		err := storage.Put(element)
		if err != nil {
			log.Fatalf("error : %v", err)
		}
	}
}

func get(storage store.Storage, elements []store.Element) {
	for _, element := range elements {
		result, err := storage.Get(element.Key)
		if err != nil {
			log.Fatalf("error : %v", err)
		}
		consume(result)
	}
}

// TODO : add combination of puts and gets scenario execution

// TODO : add ordered get scenario

func getFuncName(exec benchmarkExecution) string {
	execName := runtime.FuncForPC(reflect.ValueOf(exec).Pointer()).Name()
	idx := strings.LastIndex(execName, ".")
	return execName[idx+1:]
}

func consume(element store.Element) {
	//... store in a environment variable to avoid optimizations https://stackoverflow.com/questions/36966947/do-go-testing-b-benchmarks-prevent-unwanted-optimizations
	result = element
}
