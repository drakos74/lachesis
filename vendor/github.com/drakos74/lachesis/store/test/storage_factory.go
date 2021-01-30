package test

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/drakos74/lachesis/store/app/storage"
)

// ElementFactory generates a storage.Element
type ElementFactory func() storage.Element

// RandomFactory is an element factory that generates random byte arrays as entries for the elements of the storage
type RandomFactory struct {
	ElementFactory
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
		ElementFactory: func() storage.Element {
			key := RandomBytes(keySize)
			value := RandomBytes(valueSize)
			return storage.NewElement(key, value)
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
		ElementFactory: func() storage.Element {
			return storage.NewElement(key, RandomBytes(valueSize))
		},
		KeySize:   keySize,
		ValueSize: valueSize,
	}
}

// TODO : sequential key factory , to assert ordered key structs

// Elements will create the given number of elements with the provided factory
// it will return the elements in a slice
func Elements(n int, generator RandomFactory) []storage.Element {
	elements := make([]storage.Element, n)
	for i := 0; i < n; i++ {
		elements[i] = generator.ElementFactory()
	}
	return elements
}

// Equals is used to check equality of byte arrays
func Equals(expected, actual []byte) error {
	res := bytes.Compare(expected, actual)
	if res == 0 {
		return nil
	}
	return fmt.Errorf("expected: %v\nactual: %v", expected, actual)
}
