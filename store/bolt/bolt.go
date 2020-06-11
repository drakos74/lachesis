package bolt

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/store"

	"github.com/boltdb/bolt"
)

const bucket = "bucket"

type Store struct {
	db *bolt.DB
}

func (s Store) Put(element store.Element) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(element.Key, element.Value)
		return err
	})
}

func (s Store) Get(key store.Key) (store.Element, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			// no bucket, no value
			return fmt.Errorf(store.NoValue, key)
		}
		value = b.Get(key)
		if value == nil {
			return fmt.Errorf(store.NoValue, key)
		}
		return nil
	})
	if err != nil {
		return store.Element{}, err
	}
	return store.NewElement(key, value), nil
}

func (s Store) Metadata() store.Metadata {
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
		panic(fmt.Errorf(store.InternalError, "metadata", bucket, err))
	}

	return store.Metadata{
		Size:        count,
		KeysBytes:   keyBytes,
		ValuesBytes: valueBytes,
	}
}

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

func NewFileStore(f string) (*Store, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	log.Info().Str("bolt-dir", dir)
	return newBolt(bolt.Open(fmt.Sprintf("%s.db", f), 0600, &bolt.Options{Timeout: 1 * time.Second}))
}

// BoltFileFactory generates a bolt storage implementation
func BoltFileFactory(path string) store.StorageFactory {
	return func() store.Storage {
		// use nano, in order to create a new store each time (we want the tests to remain independent at this stage)
		s, err := NewFileStore(fmt.Sprintf("%s/%v", path, time.Now().UnixNano()))
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return s
	}
}
