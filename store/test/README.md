## Test Suite

Test Suites provide the functionality to build generic tests on the storage interfaces.

```go
// Consistency is the storage consistency test suite
type Consistency struct {
	Suite
}

// ConsistencyWithMeta is the testing suite that checks consistency of the given storage implementation
// and additionally asserts the internal metadata of the given storage
type ConsistencyWithMeta struct {
	Suite
}

// Concurrency is the concurrent test suite
type Concurrency struct {
	Suite
}

// FailureRate tracks the number of failures
// it is used for complex and intermittent test scenarios
// where no exact assumptions can be made, due to randomness of initial conditions
type FailureRate struct {
	Suite
}

... 

```

For example ... the consistency suite consists of 3 test scenarios

- Get operation on an empty store
```go
// TestVoidReadOperation tests the given storage on a read operation
func (s *Consistency) TestVoidReadOperation() {
	storage := s.newStorage()
	VoidReadOperation(s.t, storage, false)
}
```

- Put and Get Operations for the same key
```go
// TestPutOperation tests the given storage on a single put operation
func (s *Consistency) TestPutOperation() {
	storage := s.newStorage()
	ReadWriteOperation(s.t, storage, Random(10, 20), false)
}
```

- Multiple Put and Get Operations for different keys
```go
// TestMultiReadWriteOperations tests the given storage on multiple read and write operations
func (s *Consistency) TestMultiReadWriteOperations() {
	storage := s.newStorage()
	MultiReadWriteOperations(s.t, storage, Random(10, 20), false)
}
```