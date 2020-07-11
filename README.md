# lachesis
simple key-value storage experiment
## 
Lachesis is the second of the 3 Fates from ancient mythology (https://en.wikipedia.org/wiki/Moirai).
She is responsible for maintaining and storing the thread of life.

## Experiment tooling

- Storage interface / factory abstraction
- Test Suites (consistency , concurrency , failure-rate ... )
- Benchmarking Suites
- Benchmark Visualisation tool

### Interface

```go
// Storage is the low level interface for interacting with the underlying implementation in bytes
type Storage interface {
	Put(element Element) error
	Get(key Key) (Element, error)
}

// Key identifies the byte arrays used as keys of the storage
type Key []byte

// Value identifies the byte arrays used for the values of the storage
type Value []byte

// Element is a concrete implementation of the Element interface
type Element struct {
	Key
	Value
}

// StorageFactory generates a storage object
type StorageFactory func() Storage
```

### Storage Metadata

In many cases we want also to assert that our assumptions regarding the internals of the store 
are satisfied as well.

In that sense we want to extend the interface

```go
// Metadata stores internal statistics specific to the underlying storage implementation
type Metadata struct {
	Size        uint64
	KeysBytes   uint64
	ValuesBytes uint64
	Errors      errors
}
```

```go
// Storage is the low level interface for interacting with the underlying implementation in bytes
type Storage interface {
	Put(element Element) error
	Get(key Key) (Element, error)
	Metadata() Metadata
	Close() error
}
```

### Tests

We would ideally want to run the same test packages on different implementations

```go
func TestCache_KeyValueImplementation(t *testing.T) {
	new(test.ConsistencyWithMeta).Run(t, CacheFactory)
}

func testCacheSyncImplementation(t *testing.T) {
	new(test.Concurrency).Run(t, CacheFactory)
}
```

#### Suites

Note : To run the tests for the file storage the following folders must be present :

```
# storage specific test files
store/file/data/*
store/file/sync-data/*
store/file/bash/database
store/bolt/data/*
store/badger/data/*
```

### Benchmarks
Note : To run the benchmarks for the file storage the following folders must be present :

```
# project specific test files
store/benchmark/testdata/file/*
store/benchmark/testdata/scratchpad/*
store/benchmark/testdata/sync-scratchpad/*
store/benchmark/testdata/treepad/*
store/benchmark/testdata/sync-treepad/*
store/benchmark/testdata/badger/*
store/benchmark/testdata/bolt/*
```

### Visualisation