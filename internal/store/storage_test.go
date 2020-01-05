package store

import (
	"fmt"
	"lachesis/internal/model"
	"lachesis/internal/store/file"
	"lachesis/internal/store/mem"
	"lachesis/pkg"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
)

var testHash = strconv.Itoa(time.Now().Nanosecond())

var testKey = []byte(fmt.Sprintf("key_%s", testHash))
var testValue = []byte(fmt.Sprintf("value_%s", testHash))

// file storage

func getFileStorage(t *testing.T) Storage {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	db, err := file.NewSB("../../test/testdata/unit")
	assert.NoError(t, err)
	return db
}

func TestSB_PutGet(t *testing.T) {
	db := getFileStorage(t)
	testSingleRWOperation(t, db)
}

func TestSB_MultiPutGet(t *testing.T) {
	db := getFileStorage(t)
	testMultiRWOperation(t, db)
}

func TestSB_OverwritePutGet(t *testing.T) {
	db := getFileStorage(t)
	testROWOperation(t, db)
}

// naive file storage implementation is not thread-safe
func testSB_ConcurrentPutGet(t *testing.T) {
	db := getFileStorage(t)
	testConcurrentRWOperation(t, db)
}

// file storage

func getSyncFileStorage(t *testing.T) Storage {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	db, err := file.NewSyncSB("../../test/testdata/unit")
	assert.NoError(t, err)
	return db
}

func TestSyncSB_PutGet(t *testing.T) {
	db := getSyncFileStorage(t)
	testSingleRWOperation(t, db)
}

func TestSyncSB_MultiPutGet(t *testing.T) {
	db := getSyncFileStorage(t)
	testMultiRWOperation(t, db)
}

func TestSyncSB_OverwritePutGet(t *testing.T) {
	db := getSyncFileStorage(t)
	testROWOperation(t, db)
}

func TestSyncSB_ConcurrentPutGet(t *testing.T) {
	db := getSyncFileStorage(t)
	testConcurrentRWOperation(t, db)
}

// in-memory storage

func getMemStorage(t *testing.T) Storage {
	return mem.NewCache()
}

func TestMemory_PutGet(t *testing.T) {
	db := getMemStorage(t)
	testSingleRWOperation(t, db)
}

func TestMemory_MultiPutGet(t *testing.T) {
	db := getMemStorage(t)
	testMultiRWOperation(t, db)
}

func TestMemory_OverwritePutGet(t *testing.T) {
	db := getMemStorage(t)
	testROWOperation(t, db)
}

// naive in memory implementation is not thread-safe
func testMemory_ConcurrentPutGet(t *testing.T) {
	db := getMemStorage(t)
	testConcurrentRWOperation(t, db)
}

// in-memory sync storage

func getSyncMemStorage(t *testing.T) Storage {
	return mem.NewSyncCache()
}

func TestSyncMemory_PutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	testSingleRWOperation(t, db)
}

func TestSyncMemory_MultiPutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	testMultiRWOperation(t, db)
}

func TestSyncMemory_OverwritePutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	testROWOperation(t, db)
}

func TestSyncMemory_ConcurrentPutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	testConcurrentRWOperation(t, db)
}

// in-memory trie

func getTrieStorage(t *testing.T) Storage {
	return mem.NewTrie()
}

func TestTrie_PutGet(t *testing.T) {
	db := getTrieStorage(t)
	testSingleRWOperation(t, db)
}

func TestTrie_MultiPutGet(t *testing.T) {
	db := getTrieStorage(t)
	testMultiRWOperation(t, db)
}

func TestTrie_OverwritePutGet(t *testing.T) {
	db := getTrieStorage(t)
	testROWOperation(t, db)
}

// naive in memory trie implementation is not thread-safe
func testTrie_ConcurrentPutGet(t *testing.T) {
	db := getTrieStorage(t)
	testConcurrentRWOperation(t, db)
}

// in-memory sync trie

func getSyncTrieStorage(t *testing.T) Storage {
	return mem.NewSyncTrie()
}

func TestSyncTrie_PutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	testSingleRWOperation(t, db)
}

func TestSyncTrie_MultiPutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	testMultiRWOperation(t, db)
}

func TestSyncTrie_OverwritePutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	testROWOperation(t, db)
}

func TestSyncTrie_ConcurrentPutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	testConcurrentRWOperation(t, db)
}

func testSingleRWOperation(t *testing.T, storage Storage) {
	// write path
	element := model.NewObject(testKey, testValue)

	err := storage.Put(element)
	assert.NoError(t, err)

	// read path
	key := model.NewKey(testKey)

	value, err := storage.Get(key)
	assert.NoError(t, err)

	assert.Equal(t, testValue, value.Value())

	// flush path
	err = storage.Close()
	assert.NoError(t, err)
}

func testROWOperation(t *testing.T, storage Storage) {
	// write path
	element := model.NewObject(testKey, testValue)

	err := storage.Put(element)
	assert.NoError(t, err)

	// read path
	key := model.NewKey(testKey)

	value, err := storage.Get(key)
	assert.NoError(t, err)

	assert.Equal(t, testValue, value.Value())

	// do another insert
	s := []byte(strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int()))
	v, err := pkg.Concat(len(testValue)+len(s), testValue, s)
	assert.NoError(t, err)

	newElement := model.NewObject(testKey, v)
	err = storage.Put(newElement)
	assert.NoError(t, err)

	// get the new value
	newValue, err := storage.Get(key)
	assert.NoError(t, err)

	assert.Equal(t, newElement.Value(), newValue.Value())

	// flush path
	err = storage.Close()
	assert.NoError(t, err)
}

const num = 100000 // should be one million

func testMultiRWOperation(t *testing.T, storage Storage) {

	metadata := model.NewMetadata()

	// write path
	elements := make([]model.Element, 0)
	for i := 0; i < num; i++ {
		s := []byte(strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int()))
		k, err := pkg.Concat(len(testKey)+len(s), testKey, s)
		assert.NoError(t, err)
		v, err := pkg.Concat(len(testValue)+len(s), testValue, s)
		assert.NoError(t, err)
		obj := model.NewObject(k, v)
		metadata.Add(obj)
		elements = append(elements, obj)
	}

	for _, element := range elements {
		err := storage.Put(element)
		assert.NoError(t, err)
	}

	// read path

	for _, element := range elements {
		value, err := storage.Get(element)
		assert.NoError(t, err)
		assert.Equal(t, element.Value(), value.Value())
	}

	// flush path
	err := storage.Close()
	assert.NoError(t, err)

	assert.Equal(t, metadata.Size, storage.Metadata().Size)

}

// TODO : add a timeout to this test
func testConcurrentRWOperation(t *testing.T, storage Storage) {

	mutex := sync.WaitGroup{}

	var r int32
	var w int32

	for i := 0; i < num; i++ {

		mutex.Add(1)

		go func(storage Storage) {
			s := []byte(strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int()))
			k, err := pkg.Concat(len(testKey)+len(s), testKey, s)
			assert.NoError(t, err)
			v, err := pkg.Concat(len(testValue)+len(s), testValue, s)
			assert.NoError(t, err)
			obj := model.NewObject(k, v)

			go func(storage Storage, el model.Element) {
				err := storage.Put(el)
				if err != nil {
					t.Fail()
					return
				}
				atomic.AddInt32(&w, 1)
				go func(storage Storage, element model.Element) {
					key := model.NewKey(element.Key())
					result, err := storage.Get(key)
					if err != nil {
						t.Fail()
						return
					}
					atomic.AddInt32(&r, 1)
					assert.Equal(t, result.Value(), element.Value())
					mutex.Done()
				}(storage, el)

			}(storage, obj)

		}(storage)

	}

	mutex.Wait()

	// flush path
	err := storage.Close()
	assert.NoError(t, err)

	// NOTE : We will have key overlaps ... but the different stores will behave differently
	// so for now we just assert based on the read and write operations, and not the embedded metadata
	assert.Equal(t, w, r)

}
