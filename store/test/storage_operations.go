package test

import (
	"testing"

	"github.com/drakos74/lachesis/store"
	"github.com/stretchr/testify/assert"
)

// VoidReadOperation performs a read on a non-existing key
// expecting the results to be an error and an empty element
func VoidReadOperation(t *testing.T, storage store.Storage) {

	// read path
	key := RandomBytes(10)
	testElement, err := storage.Get(key)
	assert.Error(t, err)
	assert.Equal(t, store.Element{}, testElement)

	// check if store is empty
	assert.Equal(t, 0, storage.Metadata().Size)

	// wrap up
	err = storage.Close()
	assert.NoError(t, err)
}

// ReadOperation performs a read on a given key
// expecting the results to be the expected value
func IntermediateReadOperation(t *testing.T, storage store.Storage, key store.Key, expectedValue store.Value) store.Element {

	testElement, err := storage.Get(key)
	assert.NoError(t, err)

	// main assertion
	assert.Equal(t, expectedValue, testElement.Value)
	return testElement
}

// ReadWriteOperation performs a write and a following read on the storage
// it asserts that we got back the value that was put into the store
func ReadWriteOperation(t *testing.T, storage store.Storage, generate Factory) {

	element := generate()

	// check if store is empty
	assert.Equal(t, 0, storage.Metadata().Size)

	// write path
	err := storage.Put(element)
	assert.NoError(t, err)

	// read path
	savedElement := IntermediateReadOperation(t, storage, element.Key, element.Value)
	assert.Equal(t, element, savedElement)

	// check the metadata
	assert.Equal(t, 1, storage.Metadata().Size)

	// wrap up
	err = storage.Close()
	assert.NoError(t, err)
}

func ReadOverwriteOperation(t *testing.T, storage store.Storage, generate Factory) {

	element1 := generate()
	element2 := generate()
	assert.NotEqual(t, element1, element2)

	// check if store is empty
	assert.Equal(t, 0, storage.Metadata().Size)

	// write path
	err := storage.Put(element1)
	assert.NoError(t, err)

	// overwrite path
	err = storage.Put(element2)
	assert.NoError(t, err)

	// read path
	savedElement := IntermediateReadOperation(t, storage, element1.Key, element2.Value)
	assert.Equal(t, element2, savedElement)

	// check the metadata
	assert.Equal(t, 1, storage.Metadata().Size)

	// wrap up
	err = storage.Close()
	assert.NoError(t, err)
}
