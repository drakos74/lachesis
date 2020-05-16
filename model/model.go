package model

import "fmt"

// Element is the lowest level of representation of a key-value pair
type Element interface {
	Key() []byte
	Value() []byte
	Size() int
}

func String(e Element) string {
	return fmt.Sprintf("{%v,%v}", e.Key(), e.Value())
}

// Object is a concrete implementation of the Element interface
type Object struct {
	key   []byte
	value []byte
}

// Key returns the key associated with the object
func (o Object) Key() []byte {
	return o.key
}

// Value returns the key associated with the object
func (o Object) Value() []byte {
	return o.value
}

// Size returns the sum of the sizes of the key and the value
func (o Object) Size() int {
	return len(o.key) + len(o.value)
}

// NewObject creates a new Object
func NewObject(key, value []byte) Element {
	return Object{
		key:   key,
		value: value,
	}
}

// NewKey creates an object with only the key property
func NewKey(key []byte) Element {
	return Object{
		key: key,
	}
}

// Serializer converts an element into a consistent byte representation by merging the key and value
type Serializer func(element Element) ([]byte, error)

// Deserializer transform the byte array into an element object by splitting from the byte array keys and values
type Deserializer func(element Element, data []byte) (Element, error)

// Serdes combines the functionalities of the Serializer and deserializer into one single struct
type Serdes struct {
	Serializer
	Deserializer
}

type Metadata struct {
	Size        int
	KeysBytes   int
	ValuesBytes int
	Errors      errors
}

func NewMetadata() Metadata {
	return Metadata{
		Errors: make([]error, 0),
	}
}

// Merge combines 2 metadtaa instances into one
func (m *Metadata) Merge(metadata Metadata) {
	m.Size += metadata.Size
	m.KeysBytes += metadata.KeysBytes
	m.ValuesBytes += metadata.ValuesBytes
}

// Add increments the metadata state for an extra element
func (m *Metadata) Add(element Element) {
	m.Size++
	m.KeysBytes += len(element.Key())
	m.ValuesBytes += len(element.Value())
}

// Error adds the provided error to the metadata instance
func (m *Metadata) Error(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

type errors []error

func (err *errors) append(currentErr error) {
	*err = append(*err, currentErr)
}
