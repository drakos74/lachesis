package test

import (
	"math/rand"

	"github.com/drakos74/lachesis/store"
)

type Factory func() store.Element

func RandomBytes(size int) []byte {
	bb := make([]byte, size)
	rand.Read(bb)
	return bb
}

func Random(keySize, valueSize int) Factory {
	key := RandomBytes(keySize)
	value := RandomBytes(valueSize)
	return func() store.Element {
		return store.NewElement(key, value)
	}
}

func RandomValue(key []byte, valueSize int) Factory {
	return func() store.Element {
		return store.NewElement(key, RandomBytes(valueSize))
	}
}
