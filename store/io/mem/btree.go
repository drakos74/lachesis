package mem

import (
	"fmt"
	"github.com/drakos74/lachesis"

	"github.com/drakos74/lachesis/datastruct/btree"
)

// Btree is a btree storage implementation
type Btree struct {
	*btree.BTree
}

// BTreeFactory generates a Cache storage implementation
func BTreeFactory() lachesis.Storage {
	return &Btree{btree.New(10)}
}

// Put stores an element in the storage based on the given key
func (b *Btree) Put(element lachesis.Element) error {
	b.BTree.ReplaceOrInsert(element)
	return nil
}

// Get retrieves an element based on the given key
func (b *Btree) Get(key lachesis.Key) (lachesis.Element, error) {
	e := b.BTree.Get(lachesis.NewElement(key, []byte{}))
	var err error
	if lachesis.IsNil(e) {
		err = fmt.Errorf(lachesis.NoValue, key)
	}
	return e, err
}

// Metadata returns the metadata for the given storage
func (b *Btree) Metadata() lachesis.Metadata {
	c, ks, vs := b.Stats()
	return lachesis.Metadata{
		Size:        c,
		KeysBytes:   ks,
		ValuesBytes: vs,
	}
}

// Close shuts down the storage and performs any needed cleanup operations
func (b *Btree) Close() error {
	// nothing to do here
	return nil
}
