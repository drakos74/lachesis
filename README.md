# lachesis
simple key-value storage experiment
## 
Lachesis is the second of the 3 Fates from ancient mythology (https://en.wikipedia.org/wiki/Moirai).
She is responsible for maintaining and storing the thread of life.

## Benchmarks
Note : To run the tests and benchmarks for the file storage the following folders must be present :

```
# project specific test files
store/benchmark/testdata/file/*
store/benchmark/testdata/scratchpad/*
store/benchmark/testdata/sync-scratchpad/*
store/benchmark/testdata/treepad/*
store/benchmark/testdata/sync-treepad/*
store/benchmark/testdata/badger/*
store/benchmark/testdata/bolt/*

# storage specific test files
store/file/data/*
store/file/sync-data/*
store/file/bash/database
store/bolt/data/*
store/badger/data/*
```