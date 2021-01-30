package bolt

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/rs/zerolog/log"
)

const bucket = "bucket"

// Store is the storage implementation backed by a bolt files store
type Store struct {
	db *bolt.DB
}

// Put writes an element to the bolt file storage
func (s Store) Put(element storage.Element) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(element.Key, element.Value)
		return err
	})
}

// Get retrieves a value from the bolt file storage based on the given key
func (s Store) Get(key storage.Key) (storage.Element, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			// no bucket, no value
			return fmt.Errorf(storage.NoValue, key)
		}
		value = b.Get(key)
		if value == nil {
			return fmt.Errorf(storage.NoValue, key)
		}
		return nil
	})
	if err != nil {
		return storage.Element{}, err
	}
	return storage.NewElement(key, value), nil
}

// Metadata returns internal stats regarding the bolt file storage implementation
func (s Store) Metadata() storage.Metadata {
	var count uint64
	var keyBytes uint64
	var valueBytes uint64

	err := s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			atomic.AddUint64(&count, 1)
			atomic.AddUint64(&keyBytes, uint64(len(k)))
			atomic.AddUint64(&valueBytes, uint64(len(v)))
		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf(storage.InternalError, "metadata", bucket, err))
	}

	return storage.Metadata{
		Size:        count,
		KeysBytes:   keyBytes,
		ValuesBytes: valueBytes,
	}
}

// Close closes the bolt file store
func (s Store) Close() error {
	return s.db.Close()
}

func newBolt(db *bolt.DB, err error) (*Store, error) {
	if err != nil {
		return nil, fmt.Errorf("could not create store: %w", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	//defer db.Close()
	return &Store{db: db}, err
}

// NewFileStore creates a new storage implementation backed by the bolt file store
func NewFileStore(f string) (*Store, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	log.Info().Str("bolt-dir", dir)
	return newBolt(bolt.Open(fmt.Sprintf("%s.db", f), 0600, &bolt.Options{Timeout: 1 * time.Second}))
}

// FileFactory generates a bolt storage implementation
func FileFactory(path string) storage.StorageFactory {
	return func() storage.Storage {
		// use nano, in order to create a new store each time (we want the tests to remain independent at this stage)
		s, err := NewFileStore(fmt.Sprintf("%s/%v", path, time.Now().UnixNano()))
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return s
	}
}
