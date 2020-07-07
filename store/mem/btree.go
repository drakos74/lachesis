package mem

import (
	"fmt"

	"github.com/drakos74/lachesis/datastruct/btree"
	"github.com/drakos74/lachesis/store"
)

// Btree is a btree storage implementation
type Btree struct {
	*btree.BTree
}

// BTreeFactory generates a Cache storage implementation
func BTreeFactory() store.Storage {
	return &Btree{btree.New(10)}
}

// Put stores an element in the storage based on the given key
func (b *Btree) Put(element store.Element) error {
	b.BTree.ReplaceOrInsert(element)
	return nil
}

// Get retrieves an element based on the given key
func (b *Btree) Get(key store.Key) (store.Element, error) {
	e := b.BTree.Get(store.NewElement(key, []byte{}))
	var err error
	if store.IsNil(e) {
		err = fmt.Errorf(store.NoValue, key)
	}
	return e, err
}

// Metadata returns the metadata for the given storage
func (b *Btree) Metadata() store.Metadata {
	c, ks, vs := b.Stats()
	return store.Metadata{
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
