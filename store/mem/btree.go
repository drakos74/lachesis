package mem

import (
	"fmt"

	"github.com/drakos74/lachesis/internal/datastruct/btree"
	"github.com/drakos74/lachesis/store"
)

type Btree struct {
	*btree.BTree
}

// BTreeFactory generates a Cache storage implementation
func BTreeFactory() store.Storage {
	return &Btree{btree.New(10)}
}

func (b *Btree) Put(element store.Element) error {
	b.BTree.ReplaceOrInsert(element)
	return nil
}

func (b *Btree) Get(key store.Key) (store.Element, error) {
	e := b.BTree.Get(store.NewElement(key, []byte{}))
	var err error
	if store.IsNil(e) {
		err = fmt.Errorf(store.NoValue, key)
	}
	return e, err
}

func (b *Btree) Metadata() store.Metadata {
	c, ks, vs := b.Stats()
	return store.Metadata{
		Size:        c,
		KeysBytes:   ks,
		ValuesBytes: vs,
	}
}

func (b *Btree) Close() error {
	// nothing to do here
	return nil
}
