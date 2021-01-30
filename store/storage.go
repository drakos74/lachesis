package store

import (
	"bytes"
	"fmt"
)

const (
	// NoValue represents an error message string in the case where there is no value for a given key
	NoValue = "could not find element for key %v"
	// NoIndex represents the error message string in the case where no index was found for a given key
	NoIndex = "could not find index for key %v"
	// InternalError represents the error message in the case where there was an error internal to the storage implementation
	InternalError = "could not complete operation %v for %v: %w"
)

// Storage is the low level interface for interacting with the underlying implementation in bytes
type Storage interface {
	Put(element Element) error
	Get(key Key) (Element, error)
	Metadata() Metadata
	Close() error
}

// Key identifies the byte arrays used as keys of the storage
type Key []byte

// Value identifies the byte arrays used for the values of the storage
type Value []byte

// Element is a concrete implementation of the Element interface
type Element struct {
	Key
	Value
}

// StorageFactory generates a storage object
type StorageFactory func() Storage

// String returns a readable representation of an Element
func String(e Element) string {
	return fmt.Sprintf("{%v,%v}", e.Key, e.Value)
}

// Size returns the sum of the sizes of the key and the value
func (o Element) Size() int {
	return len(o.Key) + len(o.Value)
}

// NewElement creates a new Element
func NewElement(key, value []byte) Element {
	return Element{
		key,
		value,
	}
}

// Metadata stores internal statistics specific to the underlying storage implementation
type Metadata struct {
	Size        uint64
	KeysBytes   uint64
	ValuesBytes uint64
	Errors      errors
}

// NewMetadata create a new metadata struct
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
	m.KeysBytes += uint64(len(element.Key))
	m.ValuesBytes += uint64(len(element.Value))
}

// Error adds the provided error to the metadata instance
func (m *Metadata) Error(err error) {
	if err != nil {
		m.Errors.append(err)
	}
}

type errors []error

// TODO : test
func (err *errors) append(currentErr error) {
	*err = append(*err, currentErr)
}

// handle nil

// NilBytes represents an empty byte array
var NilBytes = make([]byte, 0)

// Nil represents an element that has not been initialised with ay values
var Nil = Element{}

// IsNil checks if an element has not been initialised with any properties
func IsNil(e Element) bool {
	return len(e.Key) == 0 && len(e.Value) == 0
}

// equal

// BytesEqual compares to byte arrays
func BytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// IsEqual tests the equality of 2 elements based on their keys and values
func IsEqual(a, b Element) bool {
	return BytesEqual(a.Key, b.Key) && BytesEqual(a.Value, b.Value)
}

// ordering

// IsLess compares to elements based on the natural ordering of their key bytes
func IsLess(a, b Element) bool {
	return bytes.Compare(a.Key, b.Key) < 0
}
