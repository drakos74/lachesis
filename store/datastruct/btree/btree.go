package btree

import (
	"sort"
	"sync/atomic"

	"github.com/drakos74/lachesis"
)

// BTree is a b-tree implementation of the storage interface
type BTree struct {
	degree int
	length int
	root   *node
}

// New creates a new btree storage implementation
func New(degree int) *BTree {
	return &BTree{
		degree: degree,
	}
}

// maxElements returns the max number of elements to allow per syncNode.
func (t *BTree) maxElements() int {
	return t.degree*2 - 1
}

// minElements returns the min number of elements to allow per syncNode (ignored for the
// root syncNode).
func (t *BTree) minElements() int {
	return t.degree - 1
}

// ReplaceOrInsert adds the given item to the tree.  If an item in the tree
// already equals the given one, it is removed from the tree and returned.
func (t *BTree) ReplaceOrInsert(item lachesis.Element) lachesis.Element {
	if lachesis.IsNil(item) {
		panic("nil item being added to BTree")
	}
	if t.root == nil {
		t.root = new(node)
		t.root.elements = append(t.root.elements, item)
		t.length++
		return lachesis.Nil
	}
	if len(t.root.elements) >= t.maxElements() {
		item2, second := t.root.split(t.maxElements() / 2)
		oldroot := t.root
		t.root = new(node)
		t.root.elements = append(t.root.elements, item2)
		t.root.children = append(t.root.children, oldroot, second)
	}
	out := t.root.insert(item, t.maxElements())
	if lachesis.IsNil(out) {
		t.length++
	}
	return out
}

// Get looks for the key item in the tree, returning it.  It returns nil if
// unable to find that item.
func (t *BTree) Get(key lachesis.Element) lachesis.Element {
	if t.root == nil {
		return lachesis.Nil
	}
	return t.root.get(key)
}

// Stats returns the stats of the Btree
func (t *BTree) Stats() (count, keySize, valueSize uint64) {
	var nodeCount uint64
	var depth uint64
	stats(t.root, &count, &nodeCount, &depth, &keySize, &valueSize)
	return
}

func stats(n *node, count, nodes, depth, keySize, valueSize *uint64) {

	if n == nil {
		return
	}

	atomic.AddUint64(count, uint64(len(n.elements)))

	for _, e := range n.elements {
		atomic.AddUint64(keySize, uint64(len(e.Key)))
		atomic.AddUint64(valueSize, uint64(len(e.Value)))
	}

	atomic.AddUint64(nodes, uint64(len(n.children)))

	if len(n.children) > 0 {
		atomic.AddUint64(depth, 1)
	}

	for _, c := range n.children {
		stats(c, count, nodes, depth, keySize, valueSize)
	}

}

type node struct {
	elements elements
	children children
}

// insert inserts an item into the subtree rooted at this syncNode, making sure
// no nodes in the subtree exceed maxElements elements.  Should an equivalent item be
// be found/replaced by insert, it will be returned.
func (n *node) insert(item lachesis.Element, maxElements int) lachesis.Element {
	i, found := n.elements.find(item)
	if found {
		out := n.elements[i]
		// replace element
		n.elements[i] = item
		return out
	}
	if len(n.children) == 0 {
		// if there are no available children, add to the elements
		n.elements.insertAt(i, item)
		return lachesis.Nil
	}
	if n.maybeSplitChild(i, maxElements) {
		inTree := n.elements[i]
		switch {
		case lachesis.IsLess(item, inTree):
			// no change, we want first split syncNode
		case lachesis.IsLess(inTree, item):
			i++ // we want second split syncNode
		default:
			// is equal
			out := n.elements[i]
			n.elements[i] = item
			return out
		}
	}
	return n.children[i].insert(item, maxElements)
}

// get finds the given key in the subtree and returns it.
func (n *node) get(key lachesis.Element) lachesis.Element {
	i, found := n.elements.find(key)
	if found {
		return n.elements[i]
	} else if len(n.children) > 0 {
		return n.children[i].get(key)
	}
	return lachesis.Nil
}

// maybeSplitChild checks if a child should be split, and if so splits it.
// Returns whether or not a split occurred.
func (n *node) maybeSplitChild(i, maxElements int) bool {

	if len(n.children[i].elements) < maxElements {
		return false
	}
	first := n.children[i]
	item, second := first.split(maxElements / 2)
	n.elements.insertAt(i, item)
	n.children.insertAt(i+1, second)
	return true
}

// split splits the node at the given index
func (n *node) split(i int) (lachesis.Element, *node) {
	item := n.elements[i]
	next := new(node)
	// fill up the elements for the 'next' node
	next.elements = append(next.elements, n.elements[i+1:]...)
	// fix the elements on 'this' node
	n.elements.truncate(i)

	// do the same on the children
	if len(n.children) > 0 {
		next.children = append(next.children, n.children[i+1:]...)
		n.children.truncate(i + 1)
	}
	return item, next
}

// items stores items in a syncNode.
type elements []lachesis.Element

// insertAt inserts a value into the given index, pushing all subsequent values
// forward.
func (s *elements) insertAt(index int, item lachesis.Element) {
	*s = append(*s, lachesis.Nil)
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:])
	} // else ... let it break
	(*s)[index] = item
}

// removeAt removes a value at a given index, pulling all subsequent values
// back.
func (s *elements) removeAt(index int) lachesis.Element {
	item := (*s)[index]
	copy((*s)[index:], (*s)[index+1:])
	(*s)[len(*s)-1] = lachesis.Nil
	*s = (*s)[:len(*s)-1]
	return item
}

// pop removes and returns the last element in the list.
func (s *elements) pop() (out lachesis.Element) {
	index := len(*s) - 1
	out = (*s)[index]
	(*s)[index] = lachesis.Nil
	*s = (*s)[:index]
	return
}

// truncate truncates this instance at index so that it contains only the
// first index items. index must be less than or equal to length.
func (s *elements) truncate(index int) {
	var toClear elements
	*s, toClear = (*s)[:index], (*s)[index:]
	for len(toClear) > 0 {
		toClear = toClear[copy(toClear, make(elements, 16)):]
	}
}

// find returns the index where the given item should be inserted into this
// list.  'found' is true if the item already exists in the list at the given
// index.
func (s elements) find(item lachesis.Element) (index int, found bool) {
	i := sort.Search(len(s), func(i int) bool {
		return lachesis.IsLess(item, s[i])
	})
	// if we found an index , and that index is not less than the next
	// e.g. this corresponds to an equality operation
	if i > 0 && !lachesis.IsLess(s[i-1], item) {
		return i - 1, true
	}
	return i, false
}

// children stores child nodes in a node.
type children []*node

// insertAt inserts a value into the given index, pushing all subsequent values
// forward.
func (s *children) insertAt(index int, n *node) {
	*s = append(*s, nil)
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:])
	}
	(*s)[index] = n
}

// removeAt removes a value at a given index, pulling all subsequent values
// back.
func (s *children) removeAt(index int) *node {
	n := (*s)[index]
	copy((*s)[index:], (*s)[index+1:])
	(*s)[len(*s)-1] = nil
	*s = (*s)[:len(*s)-1]
	return n
}

// pop removes and returns the last element in the list.
func (s *children) pop() (out *node) {
	index := len(*s) - 1
	out = (*s)[index]
	(*s)[index] = nil
	*s = (*s)[:index]
	return
}

// truncate truncates this instance at index so that it contains only the
// first index children. index must be less than or equal to length.
func (s *children) truncate(index int) {
	var toClear children
	*s, toClear = (*s)[:index], (*s)[index:]
	for len(toClear) > 0 {
		toClear = toClear[copy(toClear, make(children, 16)):]
	}
}
