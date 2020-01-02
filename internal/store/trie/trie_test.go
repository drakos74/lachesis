package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func new() Trie {
	k := byte('A')
	trie := NewTrie(k)
	return trie
}

func TestTrie_add(t *testing.T) {

	trie := new()

	b := byte('s')

	err := trie.add([]byte{b}, []byte{1})

	assert.NoError(t, err)

	assertRead(t, trie, []byte{b}, []byte{1})
}

func TestTrie_addOverwrite(t *testing.T) {

	trie := new()

	b := byte('s')

	err := trie.add([]byte{b}, []byte{1})

	assert.NoError(t, err)

	assertRead(t, trie, []byte{b}, []byte{1})

	err = trie.add([]byte{b}, []byte{2})
	assert.NoError(t, err)

	assertRead(t, trie, []byte{b}, []byte{2})

}

func TestTrie_Commit(t *testing.T) {

	trie := new()

	b := []byte("demo")

	err := trie.Commit(b, []byte{1})
	assert.NoError(t, err)

	assertRead(t, trie, b, []byte{1})
}

func TestTrie_MultiCommit(t *testing.T) {

	trie := new()

	b1 := []byte("demo1")
	b2 := []byte("demo2")
	b3 := []byte("demo3")

	err := trie.Commit(b1, []byte{1})
	assert.NoError(t, err)

	err = trie.Commit(b2, []byte{2})
	assert.NoError(t, err)

	err = trie.Commit(b3, []byte{3})
	assert.NoError(t, err)

	assertRead(t, trie, b1, []byte{1})
	assertRead(t, trie, b2, []byte{2})
	assertRead(t, trie, b3, []byte{3})

}

func TestTrie_CommitOvewrite(t *testing.T) {

	trie := new()

	b := []byte("demo")

	err := trie.Commit(b, []byte{1})
	assert.Nil(t, err)

	err = trie.Commit(b, []byte{2})
	assert.NoError(t, err)

	assertRead(t, trie, b, []byte{2})

}

func TestTrie_ReadNotFoundExtended(t *testing.T) {
	trie := new()

	b := []byte("demo")

	err := trie.Commit(b, []byte{1})
	assert.Nil(t, err)

	assertRead(t, trie, b, []byte{1})

	v, ok := trie.Read([]byte("demo1"))

	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestTrie_ReadNotFoundWithin(t *testing.T) {
	trie := new()

	b := []byte("demo")

	err := trie.Commit(b, []byte{1})
	assert.Nil(t, err)

	assertRead(t, trie, b, []byte{1})

	v, ok := trie.Read([]byte("dem"))

	assert.False(t, ok)
	assert.Nil(t, v)
}

func assertRead(t *testing.T, trie Trie, key []byte, value []byte) {
	v, ok := trie.Read(key)
	assert.True(t, ok)
	assert.Equal(t, value, v)
}
