package badger

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/drakos74/lachesis/model"

	"github.com/dgraph-io/badger/v2"
)

type Store struct {
	db *badger.DB
}

func (s *Store) Put(element model.Element) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(element.Key(), element.Value())
		return err
	})
}

func (s *Store) Get(element model.Element) (model.Element, error) {
	var key []byte
	var value []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(element.Key())
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
		return nil, err
	}
	return model.NewObject(key, value), nil
}

func (s *Store) Metadata() model.Metadata {
	var count int32
	err := s.db.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(badger.DefaultIteratorOptions)
		defer itr.Close()
		itr.Rewind()
		for itr.Item(); itr.Valid(); itr.Next() {
			atomic.AddInt32(&count, 1)
		}
		return nil
	})

	if err != nil {
		println(fmt.Sprintf("err = %v", err))
	}

	return model.Metadata{
		Size: int(count),
	}
}

func (s *Store) Close() error {
	return s.db.Close()
}

func NewFile(f string) (*Store, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	println(fmt.Sprintf("dir = %s", dir))
	return new(badger.Open(badger.DefaultOptions(f)))
}

func NewMem() (*Store, error) {
	return new(badger.Open(badger.DefaultOptions("").WithInMemory(true)))

}

func new(db *badger.DB, err error) (*Store, error) {
	if err != nil {
		return nil, fmt.Errorf("could not create store: %w", err)
	}
	//defer db.Close()
	return &Store{db: db}, nil
}
