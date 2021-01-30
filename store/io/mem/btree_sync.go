package mem

import (
	"fmt"
	"github.com/drakos74/lachesis"
	"sync/atomic"

	"github.com/google/btree"
)

// SyncBTree implements a storage based on a Btree data struct
type SyncBTree struct {
	*btree.BTree
}

// SyncBTreeFactory generates a concurrently safe BTree storage implementation
func SyncBTreeFactory() lachesis.Storage {
	return &SyncBTree{btree.New(10)}
}

type item struct {
	lachesis.Element
}

// Less compares 2 items in terms of natural order
func (i item) Less(than btree.Item) bool {
	return lachesis.IsLess(i.Element, than.(item).Element)
}

// Put stores an element in the storage for the given key
func (s *SyncBTree) Put(element lachesis.Element) error {
	s.BTree.ReplaceOrInsert(item{element})
	return nil
}

// Get returns an element based on the given key
func (s *SyncBTree) Get(key lachesis.Key) (lachesis.Element, error) {
	e := s.BTree.Get(item{lachesis.NewElement(key, []byte{})})
	if e == nil {
		return lachesis.Nil, fmt.Errorf(lachesis.NoValue, key)
	}
	return e.(item).Element, nil
}

// Metadata returns the metadata of the given storage
func (s *SyncBTree) Metadata() lachesis.Metadata {
	var count uint64
	var keySize uint64
	var valueSize uint64
	s.BTree.Ascend(func(i btree.Item) bool {
		if i != nil {
			e := i.(item).Element
			if !lachesis.IsNil(e) {
				atomic.AddUint64(&count, 1)
				atomic.AddUint64(&keySize, uint64(len(e.Key)))
				atomic.AddUint64(&valueSize, uint64(len(e.Value)))
				return true
			}
		}
		return false
	})
	return lachesis.Metadata{
		Size:        count,
		KeysBytes:   keySize,
		ValuesBytes: valueSize,
	}
}

// Close shuts down the storage and performs any needed cleanup
func (s *SyncBTree) Close() error {
	// no need to close anything
	return nil
}
