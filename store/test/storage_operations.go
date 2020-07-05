package test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/store"
	"github.com/stretchr/testify/assert"
)

// VoidReadOperation performs a read on a non-existing key
// expecting the results to be an error and an empty element
func VoidReadOperation(t *testing.T, storage store.Storage, checkMeta bool) {

	// read path
	key := RandomBytes(10)
	testElement, err := storage.Get(key)
	assert.Error(t, err)
	assert.Equal(t, store.Element{}, testElement)

	if checkMeta {
		// check if store is empty
		assertMeta(t, 0, 0, 0, storage.Metadata())
	} else {
		// print just the metadata
		log.Info().Msg(fmt.Sprintf("metadata = %v", storage.Metadata()))
	}

	// wrap up
	err = storage.Close()
	assert.NoError(t, err)
}

// IntermediateReadOperation performs a read on a given key
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
func ReadWriteOperation(t *testing.T, storage store.Storage, generator RandomFactory, checkMeta bool) {

	element := generator.ElementFactory()

	if checkMeta {
		// check if store is empty
		assertMeta(t, 0, 0, 0, storage.Metadata())
	}

	// write path
	err := storage.Put(element)
	assert.NoError(t, err)

	// read path
	savedElement := IntermediateReadOperation(t, storage, element.Key, element.Value)
	assert.Equal(t, element, savedElement)

	if checkMeta {
		// check the metadata
		assertMeta(t, 1, uint64(generator.KeySize), uint64(generator.ValueSize), storage.Metadata())
	} else {
		// print just the metadata
		log.Info().Msg(fmt.Sprintf("metadata = %v", storage.Metadata()))
	}

	// wrap up
	err = storage.Close()
	assert.NoError(t, err)
}

// ReadOverwriteOperation performs a write, read and write operation again on the same key
// to check the overwrite feature of the given storage
func ReadOverwriteOperation(t *testing.T, storage store.Storage, generator RandomFactory, checkMeta bool) {

	element1 := generator.ElementFactory()
	element2 := generator.ElementFactory()
	assert.NotEqual(t, element1, element2)

	if checkMeta {
		// check if store is empty
		assertMeta(t, 0, 0, 0, storage.Metadata())
	}

	// write path
	err := storage.Put(element1)
	assert.NoError(t, err)

	// overwrite path
	err = storage.Put(element2)
	assert.NoError(t, err)

	// read path
	savedElement := IntermediateReadOperation(t, storage, element1.Key, element2.Value)
	assert.Equal(t, element2, savedElement)

	if checkMeta {
		// check the metadata
		assert.Equal(t, uint64(1), storage.Metadata().Size)
	} else {
		// print just the metadata
		log.Info().Msg(fmt.Sprintf("metadata = %v", storage.Metadata()))
	}

	// wrap up
	err = storage.Close()
	assert.NoError(t, err)
}

const num = 1000

// MultiReadWriteOperations executes multiple read and write operations
func MultiReadWriteOperations(t *testing.T, storage store.Storage, generator RandomFactory, checkMeta bool) {

	metadata := store.NewMetadata()

	// generator the elements
	elements := make([]store.Element, 0)

	for i := 0; i < num; i++ {
		element := generator.ElementFactory()
		metadata.Add(element)
		elements = append(elements, element)
	}

	//  write path
	for _, element := range elements {
		err := storage.Put(element)
		assert.NoError(t, err)
	}

	// read path
	for _, element := range elements {
		value, err := storage.Get(element.Key)
		assert.NoError(t, err)
		assert.Equal(t, element.Value, value.Value)
	}

	if checkMeta {
		// assert internal stats
		assert.Equal(t, metadata.Size, storage.Metadata().Size)
	} else {
		// print just the metadata
		log.Info().Msg(fmt.Sprintf("metadata = %v", storage.Metadata()))
	}

	// wrap up
	err := storage.Close()
	assert.NoError(t, err)
}

// MultiConcurrentReadWriteOperations executes multiple concurrent read and write operations
func MultiConcurrentReadWriteOperations(t *testing.T, storage store.Storage, generator RandomFactory) {

	wg := sync.WaitGroup{}

	var r int32
	var w int32

	for i := 0; i < num; i++ {

		wg.Add(1)

		// TODO : try to make this linear
		// each element cycle is done in a different routine to generator more contention
		go func(storage store.Storage) {
			element := generator.ElementFactory()

			// put
			err := storage.Put(element)
			if err != nil {
				t.Fail()
			}
			atomic.AddInt32(&w, 1)

			// make sure we call read after the write finished
			go func() {
				// read
				key := element.Key
				result, err := storage.Get(key)
				if err != nil {
					panic(fmt.Errorf("error on read: %w", err))
				}
				atomic.AddInt32(&r, 1)
				wg.Done()
				assert.Equal(t, element.Value, result.Value)
			}()

		}(storage)

	}

	wg.Wait()

	// flush path
	err := storage.Close()
	assert.NoError(t, err)

	// NOTE : We might have key overlaps ... but the different stores will behave differently
	// so for now we just assert based on the read and write operations, and not the embedded metadata
	assert.Equal(t, w, r)

}

// Errors counts the number of errors during read and write operations
type Errors struct {
	write int32
	read  int32
}

// MultiConcurrentFailureRateOperations executes multiple concurrent operations and track the amount of errors encountered
func MultiConcurrentFailureRateOperations(t *testing.T, storage store.Storage, generator RandomFactory) (readError, writeError float64) {

	wg := sync.WaitGroup{}

	var r int32
	var w int32

	var errRead int32
	var errWrite int32

	for i := 0; i < num; i++ {

		wg.Add(1)

		// TODO : try to make this linear
		// each element cycle is done in a different routine to generate more contention
		go func(storage store.Storage) {
			element := generator.ElementFactory()

			// put
			err := storage.Put(element)

			if err != nil {
				atomic.AddInt32(&errWrite, 1)
			}
			atomic.AddInt32(&w, 1)

			// make sure we call read after the write finished
			go func() {
				// read
				key := element.Key
				result, err := storage.Get(key)
				if err != nil {
					atomic.AddInt32(&errRead, 1)
				} else {
					// we did not account for read error ...
					// but we still might have retrieved the wrong entry
					if !assert.ObjectsAreEqual(element.Value, result.Value) {
						atomic.AddInt32(&errRead, 1)
					}
				}
				atomic.AddInt32(&r, 1)
				wg.Done()
			}()

		}(storage)

	}

	wg.Wait()

	// flush path
	err := storage.Close()
	assert.NoError(t, err)

	assert.Equal(t, w, r)
	assert.Equal(t, int(w), num)
	assert.Equal(t, int(r), num)

	readError = 100 * float64(errWrite) / float64(w)
	writeError = 100 * float64(errRead) / float64(r)

	log.Info().
		Str("write", fmt.Sprintf("%.2f", writeError)).
		Str("read", fmt.Sprintf("%.2f", readError)).
		Msg("Error Rate")

	return

}

func assertMeta(t *testing.T, size, keysSize, vaLuesSize uint64, meta store.Metadata) {
	assert.Equal(t, size, meta.Size)
	// TODO : assert on the volume of the store
	//assert.Equal(t, keysSize, meta.KeysBytes)
	//assert.Equal(t, vaLuesSize, meta.ValuesBytes)
	assert.Equal(t, 0, len(meta.Errors))
}
