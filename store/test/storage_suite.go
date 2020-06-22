package test

import (
	"testing"

	"github.com/drakos74/lachesis/store"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	t          *testing.T
	newStorage func() store.Storage
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
