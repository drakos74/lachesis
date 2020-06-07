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
func (suite *Suite) SetupTest() {
	// TODO : remove if nothing else todo here
}

func (suite *Suite) TestVoidReadOperation() {
	storage := suite.newStorage()
	VoidReadOperation(suite.t, storage)
}

func (suite *Suite) TestPutOperation() {
	storage := suite.newStorage()
	ReadWriteOperation(suite.t, storage, Random(10, 20))
}

func (suite *Suite) TestReadOverwriteOperation() {
	storage := suite.newStorage()
	ReadOverwriteOperation(suite.t, storage, RandomValue(10, 20))
}

func (suite *Suite) TestMultiReadWriteOperations() {
	storage := suite.newStorage()
	MultiReadWriteOperations(suite.t, storage, Random(10, 20))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func Execute(t *testing.T, factory func() store.Storage) {
	s := new(Suite)
	s.t = t
	s.newStorage = factory
	suite.Run(t, s)
}
