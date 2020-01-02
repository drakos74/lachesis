package store

import (
	"fmt"
	"lachesis/internal/model"
	"lachesis/internal/store/file"
	"lachesis/internal/store/mem"
	"lachesis/pkg"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testHash = strconv.Itoa(time.Now().Nanosecond())

var testKey = []byte(fmt.Sprintf("key_%s", testHash))
var testValue = []byte(fmt.Sprintf("value_%s", testHash))

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

func testMultiRWOperation(t *testing.T, storage Storage) {
	// write path
	elements := make([]model.Element, 0)
	for i := 0; i < 100; i++ {
		s := []byte(strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int()))
		k, err := pkg.Concat(len(testKey)+len(s), testKey, s)
		assert.NoError(t, err)
		v, err := pkg.Concat(len(testValue)+len(s), testValue, s)
		assert.NoError(t, err)
		obj := model.NewObject(k, v)
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
}

func getFileStorage(t *testing.T) Storage {
	db, err := file.New("../../test/testdata/unit")
	assert.NoError(t, err)
	return db
}
