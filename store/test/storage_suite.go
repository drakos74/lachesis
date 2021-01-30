package test

import (
	"github.com/drakos74/lachesis"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

// Suite is the testing suite base struct
type Suite struct {
	suite.Suite
	t          *testing.T
	limit      Limit
	newStorage lachesis.StorageFactory
}

// Limit is used for asserting counts on read and write operations
type Limit struct {
	Read  float64
	Write float64
}

// SetupTest initiates a new test
// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (s *Suite) SetupTest() {
	// TODO : remove if nothing else todo here
}

// Consistency is the storage consistency test suite
type Consistency struct {
	Suite
}

// TestVoidReadOperation tests the given storage on a read operation
func (s *Consistency) TestVoidReadOperation() {
	storage := s.newStorage()
	VoidReadOperation(s.t, storage, false)
}

// TestPutOperation tests the given storage on a single put operation
func (s *Consistency) TestPutOperation() {
	storage := s.newStorage()
	ReadWriteOperation(s.t, storage, Random(10, 20), false)
}

// TestMultiReadWriteOperations tests the given storage on multiple read and write operations
func (s *Consistency) TestMultiReadWriteOperations() {
	storage := s.newStorage()
	MultiReadWriteOperations(s.t, storage, Random(10, 20), false)
}

// Run executes the Consistency test suite
func (s *Consistency) Run(t *testing.T, factory lachesis.StorageFactory) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}

// ConsistencyWithMeta is the testing suite that checks consistency of the given storage implementation
// and additionally asserts the internal metadata of the given storage
type ConsistencyWithMeta struct {
	Suite
}

// TestVoidReadOperation executes a get operation on an empty store
func (s *ConsistencyWithMeta) TestVoidReadOperation() {
	storage := s.newStorage()
	VoidReadOperation(s.t, storage, true)
}

// TestPutOperation executes a test for the put operation
func (s *ConsistencyWithMeta) TestPutOperation() {
	storage := s.newStorage()
	ReadWriteOperation(s.t, storage, Random(10, 20), true)
}

// TestReadOverwriteOperation executes a test on a write and read operation
func (s *ConsistencyWithMeta) TestReadOverwriteOperation() {
	storage := s.newStorage()
	ReadOverwriteOperation(s.t, storage, RandomValue(10, 20), true)
}

// TestMultiReadWriteOperations executes tests on several read and write operations
func (s *ConsistencyWithMeta) TestMultiReadWriteOperations() {
	storage := s.newStorage()
	MultiReadWriteOperations(s.t, storage, Random(10, 20), true)
}

// Run executes the ConsistencyWithMeta test suite
func (s *ConsistencyWithMeta) Run(t *testing.T, factory lachesis.StorageFactory) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}

// Concurrency is tge concurrent test suite
type Concurrency struct {
	Suite
}

// TestMultiConcurrentReadWriteOperations executes multiple concurrent read and write operations
func (s *Concurrency) TestMultiConcurrentReadWriteOperations() {
	storage := s.newStorage()
	MultiConcurrentReadWriteOperations(s.t, storage, Random(10, 20))
}

// Run executes the concurrency test suite
func (s *Concurrency) Run(t *testing.T, factory func() lachesis.Storage) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}

// FailureRate tracks the number of failures
// it is used for complex and intermittent test scenarios
// where no exact assumptions can be made, due to randomness of initial conditions
type FailureRate struct {
	Suite
}

// TestMultiConcurrentFailureRateOperations tests the given storage implementation on multiple concurrent read and write calls
func (s *FailureRate) TestMultiConcurrentFailureRateOperations() {
	storage := s.newStorage()
	r, w := MultiConcurrentFailureRateOperations(s.t, storage, Random(10, 20))

	// check the limits ...
	// we just need to be careful in terms of logical buffers
	s.True(r <= s.limit.Read, "error limit of %v breached for read %v", s.limit.Read, r)
	s.True(w <= s.limit.Write, "error limit of %v breached for write %v", s.limit.Write, w)

	// make sure we bound this also from the bottom
	// e.g. if we spacified an error limit, we should at least have encountered some errors
	if s.limit.Read > 0 {
		s.True(r > 0, "we should have encountered at lest some read error")
	}
	if s.limit.Write > 0 {
		s.True(w > 0, "we should have encountered at lest some write error")
	}

}

// Run executes the test suite
func (s *FailureRate) Run(t *testing.T, factory func() lachesis.Storage, limit Limit) {
	// reduce logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	s.t = t
	s.newStorage = factory
	s.limit = limit
	suite.Run(t, s)
}
