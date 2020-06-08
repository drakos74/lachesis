package test

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/drakos74/lachesis/store"
)

// Factory generates a store.Element
type Factory func() store.Element

// RandomBytes generates an array of random bytes of the given size
func RandomBytes(size int) []byte {
	bb := make([]byte, size)
	rand.Read(bb)
	return bb
}

// Random returns a factory for generating a elements with a random key and value
// key and value sizes are provided as input arguments
func Random(keySize, valueSize int) Factory {
	return func() store.Element {
		key := RandomBytes(keySize)
		value := RandomBytes(valueSize)
		return store.NewElement(key, value)
	}
}

// RandomValue returns a factory for generating elements with random values
// but always with the same 'random' key
// key and value sizes are provided as input arguments
func RandomValue(keySize, valueSize int) Factory {
	key := RandomBytes(keySize)
	return func() store.Element {
		return store.NewElement(key, RandomBytes(valueSize))
	}
}

func Equals(expected, actual []byte) error {
	res := bytes.Compare(expected, actual)
	if res == 0 {
		return nil
	}
	return fmt.Errorf("expected: %v\nactual: %v", expected, actual)
}
