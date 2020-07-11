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
	"github.com/drakos74/lachesis/store/badger"
	"github.com/drakos74/lachesis/store/bolt"
	"github.com/drakos74/lachesis/store/file"
	"github.com/drakos74/lachesis/store/test"
)

// in-memory

//func BenchmarkCache(b *testing.B) {
//	executeBenchmarks(b, mem.CacheFactory)
//}
//
//func BenchmarkSyncCache(b *testing.B) {
//	executeBenchmarks(b, mem.SyncCacheFactory)
//}
//
//func BenchmarkTrie(b *testing.B) {
//	executeBenchmarks(b, mem.TrieFactory)
//}
//
//func BenchmarkSyncTrie(b *testing.B) {
//	executeBenchmarks(b, mem.SyncTrieFactory)
//}
//
//func BenchmarkBTree(b *testing.B) {
//	executeBenchmarks(b, mem.BTreeFactory)
//}
//
//func BenchmarkSyncBTree(b *testing.B) {
//	executeBenchmarks(b, mem.SyncBTreeFactory)
//}

// file

// BenchmarkFileStorage executes the benchmarks for the file storage
//func BenchmarkFileStorage(b *testing.B) {
//	executeBenchmarks(b, file.StorageFactory("testdata/file"))
//}

// BenchmarkTriePad executes the benchmarks for the file storage
// with a trie for key indexing
func BenchmarkTriePad(b *testing.B) {
	executeBenchmarks(b, file.TriePadFactory("testdata/trie-pad"))
}

// BenchmarkClosingTriePad executes the benchmarks for the file storage
// with a trie for key indexing
func BenchmarkClosingTriePad(b *testing.B) {
	executeBenchmarks(b, file.TrieClosingPadFactory("testdata/trie-pad"))
}

// BenchmarkSyncTriePad executes the benchmarks for the thread-safe file storage
// with a trie for key indexing
func BenchmarkSyncTriePad(b *testing.B) {
	executeBenchmarks(b, file.SyncScratchPadFactory("testdata/sync-trie-pad"))
}

// BenchmarkTreePad executes the benchmarks for the file storage
// with a b-tree for key indexing
func BenchmarkTreePad(b *testing.B) {
	executeBenchmarks(b, file.TreePadFactory("testdata/btree-pad"))
}

// BenchmarkClosingTreePad executes the benchmarks for the file storage
// with a b-tree for key indexing
func BenchmarkClosingTreePad(b *testing.B) {
	executeBenchmarks(b, file.TreeClosingPadFactory("testdata/btree-pad"))
}

// BenchmarkSyncTreePad executes the benchmarks for the thread-safe file storage
// with a b-tree for key indexing
func BenchmarkSyncTreePad(b *testing.B) {
	executeBenchmarks(b, file.SyncTreePadFactory("testdata/sync-btree-pad"))
}

//BenchmarkMemBadger executes the benchmarks for badger in-memory store
func BenchmarkMemBadger(b *testing.B) {
	executeBenchmarks(b, badger.MemoryFactory)
}

// BenchmarkFileBadger executes the benchmarks for badger file store
func BenchmarkFileBadger(b *testing.B) {
	executeBenchmarks(b, badger.FileFactory("testdata/badger"))
}

// BenchmarkFileBolt executes the benchmarks for badger file store
func BenchmarkFileBolt(b *testing.B) {
	executeBenchmarks(b, bolt.FileFactory("testdata/bolt"))
}

func executeBenchmarks(b *testing.B, storageFactory func() store.Storage) {

	// reduce logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	scenarios := []Scenario{
		Benchmark(Evolution().
			add(limit(5)).
			add(num(pow(10))).
			create(),
			10, 4, 10). // initial values
			Name("num-objects"),
		Benchmark(Evolution().
			add(limit(5)).
			add(key(pow(2))).
			create(),
			10, 4, 10). // initial values
			Name("increasing-key-size"),
		Benchmark(Evolution().
			add(limit(5)).
			add(key(pow(2))).
			create(),
			10, 4, 10). // initial values
			Name("increasing-value-size"),
	}

	for _, scenario := range scenarios {
		storage := storageFactory()
		executeBenchmark(b, storage, scenario, put, get)
	}

}

func executeBenchmark(b *testing.B, storage store.Storage, scenario Scenario, execution ...benchmarkExecution) {

	for scenario.next() {
		currentScenario := scenario.get()
		elements := test.Elements(currentScenario.Num, test.Random(currentScenario.KeySize, currentScenario.ValueSize))
		for _, exec := range execution {
			b.Run(fmt.Sprintf("%s|%s|%s|num-objects:%d|size-key:%d|size-value:%d|", reflect.TypeOf(storage).String(), getFuncName(exec), scenario.name, scenario.Num, scenario.KeySize, scenario.ValueSize), func(b *testing.B) {
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
