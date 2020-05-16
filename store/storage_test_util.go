package store

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/drakos74/lachesis/internal"

	"github.com/drakos74/lachesis/model"
	"github.com/drakos74/lachesis/store/file"
	"github.com/drakos74/lachesis/store/mem"

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
	ExecSingleRWOperation(t, db)
}

func TestSB_MultiPutGet(t *testing.T) {
	db := getFileStorage(t)
	ExecMultiRWOperation(t, db)
}

func TestSB_OverwritePutGet(t *testing.T) {
	db := getFileStorage(t)
	ExecROWOperation(t, db)
}

// naive file storage implementation is not thread-safe
func testSB_ConcurrentPutGet(t *testing.T) {
	db := getFileStorage(t)
	ExecConcurrentRWOperation(t, db)
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
	ExecSingleRWOperation(t, db)
}

func TestSyncSB_MultiPutGet(t *testing.T) {
	db := getSyncFileStorage(t)
	ExecMultiRWOperation(t, db)
}

func TestSyncSB_OverwritePutGet(t *testing.T) {
	db := getSyncFileStorage(t)
	ExecROWOperation(t, db)
}

func TestSyncSB_ConcurrentPutGet(t *testing.T) {
	db := getSyncFileStorage(t)
	ExecConcurrentRWOperation(t, db)
}

// in-memory storage

func getMemStorage(t *testing.T) Storage {
	return mem.NewCache()
}

func TestMemory_PutGet(t *testing.T) {
	db := getMemStorage(t)
	ExecSingleRWOperation(t, db)
}

func TestMemory_MultiPutGet(t *testing.T) {
	db := getMemStorage(t)
	ExecMultiRWOperation(t, db)
}

func TestMemory_OverwritePutGet(t *testing.T) {
	db := getMemStorage(t)
	ExecROWOperation(t, db)
}

// naive in memory implementation is not thread-safe
func testMemory_ConcurrentPutGet(t *testing.T) {
	db := getMemStorage(t)
	ExecConcurrentRWOperation(t, db)
}

// in-memory sync storage

func getSyncMemStorage(t *testing.T) Storage {
	return mem.NewSyncCache()
}

func TestSyncMemory_PutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	ExecSingleRWOperation(t, db)
}

func TestSyncMemory_MultiPutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	ExecMultiRWOperation(t, db)
}

func TestSyncMemory_OverwritePutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	ExecROWOperation(t, db)
}

func TestSyncMemory_ConcurrentPutGet(t *testing.T) {
	db := getSyncMemStorage(t)
	ExecConcurrentRWOperation(t, db)
}

// in-memory trie

func getTrieStorage(t *testing.T) Storage {
	return mem.NewTrie()
}

func TestTrie_PutGet(t *testing.T) {
	db := getTrieStorage(t)
	ExecSingleRWOperation(t, db)
}

func TestTrie_MultiPutGet(t *testing.T) {
	db := getTrieStorage(t)
	ExecMultiRWOperation(t, db)
}

func TestTrie_OverwritePutGet(t *testing.T) {
	db := getTrieStorage(t)
	ExecROWOperation(t, db)
}

// naive in memory trie implementation is not thread-safe
func testTrie_ConcurrentPutGet(t *testing.T) {
	db := getTrieStorage(t)
	ExecConcurrentRWOperation(t, db)
}

// in-memory sync trie

func getSyncTrieStorage(t *testing.T) Storage {
	return mem.NewSyncTrie()
}

func TestSyncTrie_PutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	ExecSingleRWOperation(t, db)
}

func TestSyncTrie_MultiPutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	ExecMultiRWOperation(t, db)
}

func TestSyncTrie_OverwritePutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	ExecROWOperation(t, db)
}

func TestSyncTrie_ConcurrentPutGet(t *testing.T) {
	db := getSyncTrieStorage(t)
	ExecConcurrentRWOperation(t, db)
}

func ExecSingleRWOperation(t *testing.T, storage Storage) {
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

func ExecROWOperation(t *testing.T, storage Storage) {
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
	v, err := internal.Concat(len(testValue)+len(s), testValue, s)
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

func ExecMultiRWOperation(t *testing.T, storage Storage) {

	metadata := model.NewMetadata()

	// write path
	elements := make([]model.Element, 0)

	for i := 0; i < num; i++ {
		s := []byte(strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int()))
		k, err := internal.Concat(len(testKey)+len(s), testKey, s)
		assert.NoError(t, err)
		v, err := internal.Concat(len(testValue)+len(s), testValue, s)
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

	assert.Equal(t, metadata.Size, storage.Metadata().Size)

	err := storage.Close()
	assert.NoError(t, err)

}

// TODO : add a timeout to this test
func ExecConcurrentRWOperation(t *testing.T, storage Storage) {

	mutex := sync.WaitGroup{}

	var r int32
	var w int32

	for i := 0; i < num; i++ {

		mutex.Add(1)

		go func(storage Storage) {
			s := []byte(strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int()))
			k, err := internal.Concat(len(testKey)+len(s), testKey, s)
			assert.NoError(t, err)
			v, err := internal.Concat(len(testValue)+len(s), testValue, s)
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
