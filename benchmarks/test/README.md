## Benchmarks

### Scenarios

Scenarios are the basic benchmark conditions.
As part of each benchmark, we are measuring latency and-/or memory allocations.

This abstraction allows us to create measurements for increasing parameters 
`key-size` , `value-size` , `number of elements`

```go
        scenarios := []Scenario{
		Benchmark(Evolution().
			add(limit(6)).     // make 6 iterations
			add(num(pow(10))). // increase the number of elements by a power of 10
			create(),
			10, 10, 20), // initial values
		Benchmark(Evolution().
			add(limit(10)).   // make 10 iterations
			add(key(add(4))). // add 4 to the key-size
			create(),
			1000, 4, 100), // initial values
		Benchmark(Evolution().
			add(limit(8)).      // 8 iterations
			add(value(pow(2))). // increase the value-size by a power of 2
			create(),
			1000, 16, 2), // initial values
	}
```

Scenarios bear evolution logic for the test parameters

```go

	for _, scenario := range scenarios {
		storage := storageFactory()
		executeBenchmark(b, storage, scenario, put, get)
	}
```

We let every scenario  evolve

```go
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
```

This will produce several benchmarks. The format of printing the benchmarks is important here

```go

```