package btree

import (
	"github.com/drakos74/lachesis/store/store/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNode_New(t *testing.T) {

	n := new(node)

	newElement := test.Random(5, 10).ElementFactory

	for i := 0; i < 100; i++ {
		n.elements = append(n.elements, newElement())
	}

	assert.Equal(t, 100, len(n.elements))

}

func TestNode_Put(t *testing.T) {

	n := new(node)
	// cap at 2
	n.children = make(children, 0, 2)

	newElement := test.Random(5, 10).ElementFactory

	for i := 0; i < 5; i++ {
		n.insert(newElement(), 2)
	}

	assert.Equal(t, 5, len(n.elements))
	// note we will never end up adding more children, because we started with no children at all
	assert.Equal(t, 0, len(n.children))

}

func TestBTree_Put(t *testing.T) {

	n := New(2)

	newElement := test.Random(5, 10).ElementFactory

	for i := 0; i < 5; i++ {
		n.ReplaceOrInsert(newElement())
	}

	assert.Equal(t, 1, len(n.root.elements))
	assert.Equal(t, 2, len(n.root.children))

}

func TestBTree_PutGet(t *testing.T) {

	tree := New(2)

	elements := test.Elements(100, test.Random(5, 10))

	for _, element := range elements {
		tree.ReplaceOrInsert(element)
	}

	assert.Equal(t, 1, len(tree.root.elements))
	assert.Equal(t, 2, len(tree.root.children))

	for _, element := range elements {
		el := tree.Get(element)
		assert.Equal(t, element, el)
	}

}

func TestBTree_Count(t *testing.T) {

	tree := New(4)

	elements := test.Elements(100, test.Random(5, 10))

	for _, element := range elements {
		tree.ReplaceOrInsert(element)
	}

	assert.Equal(t, 2, len(tree.root.elements))
	assert.Equal(t, 3, len(tree.root.children))

	// elements
	var e uint64
	// node count
	var n uint64
	// depth count e.g. nodes with children
	var d uint64
	// key size
	var ks uint64
	var vs uint64

	stats(tree.root, &e, &n, &d, &ks, &vs)

	assert.Equal(t, 100, int(e))
	assert.Equal(t, 20, int(n))
	assert.Equal(t, 4, int(d))

	assert.Equal(t, 100*5, int(ks))
	assert.Equal(t, 100*10, int(vs))

}
