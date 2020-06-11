package test

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/drakos74/lachesis/store"
)

// Factory generates a store.Element
type Factory func() store.Element

type RandomFactory struct {
	Factory
	KeySize   int
	ValueSize int
}

// RandomBytes generates an array of random bytes of the given size
func RandomBytes(size int) []byte {
	bb := make([]byte, size)
	rand.Read(bb)
	return bb
}

// Random returns a factory for generating a elements with a random key and value
// key and value sizes are provided as input arguments
func Random(keySize, valueSize int) RandomFactory {
	return RandomFactory{
		Factory: func() store.Element {
			key := RandomBytes(keySize)
			value := RandomBytes(valueSize)
			return store.NewElement(key, value)
		},
		KeySize:   keySize,
		ValueSize: valueSize,
	}
}

// RandomValue returns a factory for generating elements with random values
// but always with the same 'random' key
// key and value sizes are provided as input arguments
func RandomValue(keySize, valueSize int) RandomFactory {
	key := RandomBytes(keySize)
	return RandomFactory{
		Factory: func() store.Element {
			return store.NewElement(key, RandomBytes(valueSize))
		},
		KeySize:   keySize,
		ValueSize: valueSize,
	}
}

// Elements will create the given number of elements with the provided factory
// it will return the elements in a slice
func Elements(n int, generate Factory) []store.Element {
	elements := make([]store.Element, n)
	for i := 0; i < n; i++ {
		elements[i] = generate()
	}
	return elements
}

func Equals(expected, actual []byte) error {
	res := bytes.Compare(expected, actual)
	if res == 0 {
		return nil
	}
	return fmt.Errorf("expected: %v\nactual: %v", expected, actual)
}
