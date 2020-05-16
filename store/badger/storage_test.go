package badger

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/drakos74/lachesis/store"
)

// file storage

func getFileStorage(t *testing.T) store.Storage {
	s, err := NewFile(fmt.Sprintf("/Users/drakos74/Projects/lachesis/test/testdata/badger/%v", time.Now().Unix()))
	assert.NoError(t, err)
	return s
}

func TestFile_PutGet(t *testing.T) {
	db := getFileStorage(t)
	store.ExecSingleRWOperation(t, db)
}

func TestFile_MultiPutGet(t *testing.T) {
	db := getFileStorage(t)
	store.ExecMultiRWOperation(t, db)
}

func TestFile_OverwritePutGet(t *testing.T) {
	db := getFileStorage(t)
	store.ExecROWOperation(t, db)
}

func TestFile_ConcurrentPutGet(t *testing.T) {
	db := getFileStorage(t)
	store.ExecConcurrentRWOperation(t, db)
}

// memory storage

func getMemStorage(t *testing.T) store.Storage {
	s, err := NewMem()
	assert.NoError(t, err)
	return s
}

func TestSyncMem_PutGet(t *testing.T) {
	db := getMemStorage(t)
	store.ExecSingleRWOperation(t, db)
}

func TestSyncMem_MultiPutGet(t *testing.T) {
	db := getMemStorage(t)
	store.ExecMultiRWOperation(t, db)
}

func TestSyncMem_OverwritePutGet(t *testing.T) {
	db := getMemStorage(t)
	store.ExecROWOperation(t, db)
}

func TestSyncMem_ConcurrentPutGet(t *testing.T) {
	db := getMemStorage(t)
	store.ExecConcurrentRWOperation(t, db)
}
