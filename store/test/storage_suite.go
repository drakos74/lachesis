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

type KeyValue struct {
	Suite
}

func (s *KeyValue) TestVoidReadOperation() {
	storage := s.newStorage()
	VoidReadOperation(s.t, storage)
}

func (s *KeyValue) TestPutOperation() {
	storage := s.newStorage()
	ReadWriteOperation(s.t, storage, Random(10, 20))
}

func (s *KeyValue) TestReadOverwriteOperation() {
	storage := s.newStorage()
	ReadOverwriteOperation(s.t, storage, RandomValue(10, 20))
}

func (s *KeyValue) TestMultiReadWriteOperations() {
	storage := s.newStorage()
	MultiReadWriteOperations(s.t, storage, Random(10, 20))
}

func (s *KeyValue) Run(t *testing.T, factory func() store.Storage) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}

type Concurrent struct {
	Suite
}

func (s *Concurrent) TestMultiConcurrentReadWriteOperations() {
	storage := s.newStorage()
	MultiConcurrentReadWriteOperations(s.t, storage, Random(10, 20))
}

func (s *Concurrent) Run(t *testing.T, factory func() store.Storage) {
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}
