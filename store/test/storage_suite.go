package test

import (
	"testing"

	"github.com/rs/zerolog"

	"github.com/drakos74/lachesis/store"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	t          *testing.T
	limit      Limit
	newStorage func() store.Storage
}

type Limit struct {
	Read  float64
	Write float64
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (s *Suite) SetupTest() {
	// TODO : remove if nothing else todo here
}

type Consistency struct {
	Suite
}

func (s *Consistency) TestVoidReadOperation() {
	storage := s.newStorage()
	VoidReadOperation(s.t, storage, false)
}

func (s *Consistency) TestPutOperation() {
	storage := s.newStorage()
	ReadWriteOperation(s.t, storage, Random(10, 20), false)
}

func (s *Consistency) TestMultiReadWriteOperations() {
	storage := s.newStorage()
	MultiReadWriteOperations(s.t, storage, Random(10, 20), false)
}

func (s *Consistency) Run(t *testing.T, factory store.StorageFactory) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}

type ConsistencyWithMeta struct {
	Suite
}

func (s *ConsistencyWithMeta) TestVoidReadOperation() {
	storage := s.newStorage()
	VoidReadOperation(s.t, storage, true)
}

func (s *ConsistencyWithMeta) TestPutOperation() {
	storage := s.newStorage()
	ReadWriteOperation(s.t, storage, Random(10, 20), true)
}

func (s *ConsistencyWithMeta) TestReadOverwriteOperation() {
	storage := s.newStorage()
	ReadOverwriteOperation(s.t, storage, RandomValue(10, 20), true)
}

func (s *ConsistencyWithMeta) TestMultiReadWriteOperations() {
	storage := s.newStorage()
	MultiReadWriteOperations(s.t, storage, Random(10, 20), true)
}

func (s *ConsistencyWithMeta) Run(t *testing.T, factory store.StorageFactory) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}

type Concurrency struct {
	Suite
}

func (s *Concurrency) TestMultiConcurrentReadWriteOperations() {
	storage := s.newStorage()
	MultiConcurrentReadWriteOperations(s.t, storage, Random(10, 20))
}

func (s *Concurrency) Run(t *testing.T, factory func() store.Storage) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}

type FailureRate struct {
	Suite
}

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

func (s *FailureRate) Run(t *testing.T, factory func() store.Storage, limit Limit) {
	// reduce logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	s.t = t
	s.newStorage = factory
	s.limit = limit
	suite.Run(t, s)
}
