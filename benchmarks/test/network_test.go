package test

import (
	"testing"

	"github.com/drakos74/lachesis/benchmarks/network"
	"github.com/drakos74/lachesis/benchmarks/network/lb"
	"github.com/drakos74/lachesis/internal/infra/file"
	"github.com/drakos74/lachesis/internal/infra/mem"
)

// No benefit from distribution
// because we only use memory

func BenchmarkCacheNetwork_SinglePartition(b *testing.B) {
	executeBenchmarks(b, network.Factory().
		Nodes(1).
		Storage(mem.CacheFactory).
		Router(lb.ShardedPartition).
		Create())
}

func BenchmarkCacheNetwork_MultiPartition(b *testing.B) {
	executeBenchmarks(b, network.Factory().
		Nodes(10).
		Storage(mem.CacheFactory).
		Router(lb.ShardedPartition).
		Create())
}

// Note difference for distributed case
// as we have also files involved

func BenchmarkPadNetwork_SinglePartition(b *testing.B) {
	executeBenchmarks(b, network.Factory().
		Nodes(1).
		Storage(file.SyncScratchPadFactory("testdata/scratchpad")).
		Router(lb.ShardedPartition).
		Create())
}

func BenchmarkPadNetwork_MultiPartition(b *testing.B) {
	executeBenchmarks(b, network.Factory().
		Nodes(10).
		Storage(file.SyncScratchPadFactory("testdata/scratchpad")).
		Router(lb.ShardedPartition).
		Create())
}
