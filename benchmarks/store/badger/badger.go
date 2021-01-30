package badger

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/rs/zerolog/log"
)

// Store is the storage implementation backed by a badger store
type Store struct {
	db *badger.DB
}

// Put writes a key into the badger store
func (s *Store) Put(element storage.Element) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(element.Key, element.Value)
		return err
	})
}

// Get retrieves a value for the given key from the badger storage implementation
func (s *Store) Get(key storage.Key) (storage.Element, error) {
	var value []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf(storage.NoValue, key)
		}
		key = item.KeyCopy(nil)
		value, err = item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf(storage.InternalError, "get", key, err)
		}

		return nil
	})
	if err != nil {
		return storage.Element{}, err
	}
	return storage.NewElement(key, value), nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (s *Store) Metadata() storage.Metadata {
	var count uint64
	var keySize uint64
	var valueSize uint64
	err := s.db.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(badger.DefaultIteratorOptions)
		defer itr.Close()
		for itr.Rewind(); itr.Valid(); itr.Next() {
			atomic.AddUint64(&count, 1)
			item := itr.Item()
			keySize += uint64(item.KeySize())
			valueSize += uint64(item.ValueSize())
		}
		return nil
	})

	if err != nil {
		println(fmt.Sprintf("err = %v", err))
	}

	// TODO : add also sizes ...
	return storage.Metadata{
		Size:        count,
		KeysBytes:   keySize,
		ValuesBytes: valueSize,
	}
}

// Close closes the badger storage implementation
func (s *Store) Close() error {
	return s.db.Close()
}

func newBadger(db *badger.DB, err error) (*Store, error) {
	if err != nil {
		return nil, fmt.Errorf("could not create store: %w", err)
	}
	//defer db.Close()
	return &Store{db: db}, nil
}

// FileFactory generates a badger file storage implementation with the predefined path.
func FileStore(path string) storage.StorageFactory {
	// use nano, in order to create a new store each time (we want the tests to remain independent at this stage)
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("error during badger store creation: %v", err))
	}
	log.Info().Str("badger-dir", dir)
	return func() storage.Storage {
		s, err := newBadger(badger.Open(badger.DefaultOptions(fmt.Sprintf("%s", path))))
		if err != nil {
			panic(fmt.Sprintf("error during badger store creation: %v", err))
		}
		return s
	}
}

// FileFactory generates a badger file storage implementation with a unique path.
// This method should be used for local testing.
func FileFactory(path string) storage.StorageFactory {
	// use nano, in order to create a new store each time (we want the tests to remain independent at this stage)
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("error during badger store creation: %v", err))
	}
	log.Info().Str("badger-dir", dir)
	return func() storage.Storage {
		s, err := newBadger(badger.Open(badger.DefaultOptions(fmt.Sprintf("%s/%v", path, time.Now().UnixNano()))))
		if err != nil {
			panic(fmt.Sprintf("error during badger store creation: %v", err))
		}
		return s
	}
}

// NewMemoryStore creates a new badger memory store
func NewMemoryStore() (*Store, error) {
	return newBadger(badger.Open(badger.DefaultOptions("").WithInMemory(true)))
}

// MemoryFactory generates a badger in-memory storage implementation
func MemoryFactory() storage.Storage {
	s, err := NewMemoryStore()
	if err != nil {
		panic(fmt.Sprintf("error during store creation: %v", err))
	}
	return s
}
