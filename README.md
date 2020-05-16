# lachesis
simple key-value storage experiment
## 
Lachesis is the second of the 3 Fates from ancient mythology (https://en.wikipedia.org/wiki/Moirai).
She is responsible for maintaining and storing the thread of life.

## Benchmarks
Note : To run the tests and benchmarks for the file storage the following folders must be present :

- test/testdata/bench
- test/testdata/unit
```
$ go test ./... -gcflags=-N -run=xxx -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: lachesis/internal/store
BenchmarkSB/*file.SB:put/num:1,size-key:100,size-value:1000-16                             83684             13760 ns/op            1080 B/op          3 allocs/op
BenchmarkSB/*file.SB:get/num:1,size-key:100,size-value:1000-16                            174919              6414 ns/op            1072 B/op          2 allocs/op
BenchmarkSB/*file.SB:put/num:10,size-key:100,size-value:1000-16                            10000            136816 ns/op           10800 B/op         30 allocs/op
BenchmarkSB/*file.SB:get/num:10,size-key:100,size-value:1000-16                            17689            113908 ns/op           10720 B/op         20 allocs/op
BenchmarkSB/*file.SB:put/num:100,size-key:100,size-value:1000-16                             391           2699301 ns/op          108002 B/op        300 allocs/op
BenchmarkSB/*file.SB:get/num:100,size-key:100,size-value:1000-16                            1328           1149386 ns/op          107200 B/op        200 allocs/op
BenchmarkSB/*file.SB:put/num:1000,size-key:100,size-value:1000-16                             38          28929329 ns/op         1080002 B/op       3000 allocs/op
BenchmarkSB/*file.SB:get/num:1000,size-key:100,size-value:1000-16                            100          16263561 ns/op         1072000 B/op       2000 allocs/op
BenchmarkMemory/*mem.Cache:put/num:1,size-key:100,size-value:1000-16                    12915903               129 ns/op             112 B/op          1 allocs/op
BenchmarkMemory/*mem.Cache:get/num:1,size-key:100,size-value:1000-16                     8707650               127 ns/op              48 B/op          1 allocs/op
BenchmarkMemory/*mem.Cache:put/num:10,size-key:100,size-value:1000-16                     785540              1370 ns/op            1120 B/op         10 allocs/op
BenchmarkMemory/*mem.Cache:get/num:10,size-key:100,size-value:1000-16                     923350              1290 ns/op             480 B/op         10 allocs/op
BenchmarkMemory/*mem.Cache:put/num:100,size-key:100,size-value:1000-16                    121375             11535 ns/op           11200 B/op        100 allocs/op
BenchmarkMemory/*mem.Cache:get/num:100,size-key:100,size-value:1000-16                     71890             15675 ns/op            4800 B/op        100 allocs/op
BenchmarkMemory/*mem.Cache:put/num:1000,size-key:100,size-value:1000-16                     9344            195839 ns/op          112000 B/op       1000 allocs/op
BenchmarkMemory/*mem.Cache:get/num:1000,size-key:100,size-value:1000-16                     9018            254647 ns/op           48000 B/op       1000 allocs/op
BenchmarkTrie/*mem.Trie:put/num:1,size-key:100,size-value:1000-16                         235987              5279 ns/op              48 B/op          1 allocs/op
BenchmarkTrie/*mem.Trie:get/num:1,size-key:100,size-value:1000-16                         274467              6750 ns/op              48 B/op          1 allocs/op
BenchmarkTrie/*mem.Trie:put/num:10,size-key:100,size-value:1000-16                         17958             56874 ns/op             480 B/op         10 allocs/op
BenchmarkTrie/*mem.Trie:get/num:10,size-key:100,size-value:1000-16                         26592             53625 ns/op             480 B/op         10 allocs/op
BenchmarkTrie/*mem.Trie:put/num:100,size-key:100,size-value:1000-16                         2504            427245 ns/op            4800 B/op        100 allocs/op
BenchmarkTrie/*mem.Trie:get/num:100,size-key:100,size-value:1000-16                         3224            377288 ns/op            4800 B/op        100 allocs/op
BenchmarkTrie/*mem.Trie:put/num:1000,size-key:100,size-value:1000-16                         186           6209379 ns/op           48000 B/op       1000 allocs/op
BenchmarkTrie/*mem.Trie:get/num:1000,size-key:100,size-value:1000-16                         183           6580642 ns/op           48000 B/op       1000 allocs/op
PASS
ok      lachesis/internal/store 43.969s
```

### Cache

- simplest but harder to scale/improve

### Trie

- lower size than the keys and values added together

### File

- need to do clean up tasks (for updates etc ...)