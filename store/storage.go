package store

import "fmt"

const (
	NoValue       = "could not find element for key %v"
	NoIndex       = "could not find index for key %v"
	InternalError = "could not complete operation %v for %v: %w"
)

// Storage is the low level interface for interacting with the underlying implementation in bytes
type Storage interface {
	Put(element Element) error
	Get(key Key) (Element, error)
	Metadata() Metadata
	Close() error
}

type Key []byte
type Value []byte

// Element is a concrete implementation of the Element interface
type Element struct {
	Key
	Value
}

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
