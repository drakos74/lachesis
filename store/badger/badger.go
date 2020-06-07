package badger

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/drakos74/lachesis/store"

	"github.com/dgraph-io/badger/v2"
)

type Store struct {
	db *badger.DB
}

func (s *Store) Put(element store.Element) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(element.Key, element.Value)
		return err
	})
}

func (s *Store) Get(key store.Key) (store.Element, error) {
	var value []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		key = item.KeyCopy(nil)
		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return store.Element{}, err
	}
	return store.NewElement(key, value), nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (s *Store) Metadata() store.Metadata {
	var count int32
	err := s.db.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(badger.DefaultIteratorOptions)
		defer itr.Close()
		for itr.Rewind(); itr.Valid(); itr.Next() {
			atomic.AddInt32(&count, 1)
		}
		return nil
	})

	if err != nil {
		println(fmt.Sprintf("err = %v", err))
	}

	return store.Metadata{
		Size: int(count),
	}
}

func (s *Store) Close() error {
	return s.db.Close()
}

func NewFileStore(f string) (*Store, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	println(fmt.Sprintf("dir = %s", dir))
	return newBadger(badger.Open(badger.DefaultOptions(f)))
}

func NewMemoryStore() (*Store, error) {
	return newBadger(badger.Open(badger.DefaultOptions("").WithInMemory(true)))

}

func newBadger(db *badger.DB, err error) (*Store, error) {
	if err != nil {
		return nil, fmt.Errorf("could not create store: %w", err)
	}
	//defer db.Close()
	return &Store{db: db}, nil
}
