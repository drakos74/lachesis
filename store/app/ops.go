package app

import (
	"fmt"
	"github.com/drakos74/lachesis"
	"github.com/drakos74/lachesis/io/bytes"
)

// TODO :  make it simpler , if there is no other use fot it
// Serializer converts an element into a consistent byte representation by merging the key and value
type Join func(element store.Element) ([]byte, error)

// Deserializer transform the byte array into an element object by splitting from the byte array keys and values
type Split func(key store.Key, data []byte) (store.Element, error)

// ConcatOperator combines the functionalities of the Join and Split methods into one single struct
type ConcatOperator struct {
	Join
	Split
}

// IndexedConcat handles the concatenation logic
func IndexedConcat() ConcatOperator {
	nl := []byte{byte('\n')}
	return ConcatOperator{
		Join: func(element store.Element) ([]byte, error) {
			b, err := bytes.Concat(len(element.Value)+len(nl), element.Value, nl)
			if err != nil {
				return nil, fmt.Errorf("could not serialize value %w", err)
			}
			return b, nil
		},
		Split: func(key store.Key, data []byte) (store.Element, error) {
			n := len(data) - len(nl)
			return store.NewElement(key, data[0:n]), nil
		},
	}
}

// RawConcat handles the concatenation logic
func RawConcat() ConcatOperator {
	nl := []byte{byte('\n')}
	return ConcatOperator{
		Join: func(element store.Element) ([]byte, error) {
			b, err := bytes.Concat(len(element.Value)+len(nl), element.Value, nl)
			if err != nil {
				return nil, fmt.Errorf("could not serialize value %w", err)
			}
			return b, nil
		},
		Split: func(key store.Key, data []byte) (store.Element, error) {
			n := len(data) - len(nl)
			return store.NewElement(key, data[0:n+1]), nil
		},
	}
}
