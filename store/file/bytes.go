package file

import (
	"encoding/binary"
	"fmt"

	"github.com/drakos74/lachesis/store"
)

// TODO :  make it simpler , if there is no other use fot it
// Serializer converts an element into a consistent byte representation by merging the key and value
type join func(element store.Element) ([]byte, error)

// Deserializer transform the byte array into an element object by splitting from the byte array keys and values
type split func(key store.Key, data []byte) (store.Element, error)

// concat combines the functionalities of the join and split methods into one single struct
type concat struct {
	join
	split
}

// Handle the concatenation logic
func newIndexedConcat() concat {
	nl := []byte{byte('\n')}
	return concat{
		join: func(element store.Element) ([]byte, error) {
			b, err := Concat(len(element.Value)+len(nl), element.Value, nl)
			if err != nil {
				return nil, fmt.Errorf("could not serialize value %w", err)
			}
			return b, nil
		},
		split: func(key store.Key, data []byte) (store.Element, error) {
			n := len(data) - len(nl)
			return store.NewElement(key, data[0:n]), nil
		},
	}
}

// Handle the concatenation logic
func newRawConcat() concat {
	nl := []byte{byte('\n')}
	return concat{
		join: func(element store.Element) ([]byte, error) {
			b, err := Concat(len(element.Value)+len(nl), element.Value, nl)
			if err != nil {
				return nil, fmt.Errorf("could not serialize value %w", err)
			}
			return b, nil
		},
		split: func(key store.Key, data []byte) (store.Element, error) {
			n := len(data) - len(nl)
			return store.NewElement(key, data[0:n+1]), nil
		},
	}
}

// TODO : optimise by using append
func Concat(size int, arrays ...[]byte) ([]byte, error) {
	arr := make([]byte, size)
	i := 0
	for _, array := range arrays {
		for _, b := range array {
			arr[i] = b
			i++
		}
	}
	if i != size {
		return nil, fmt.Errorf("size argument does not match %d vs %d", size, i)
	}
	return arr, nil
}

// Handle the ScratchPad fileIndex

const (
	maxKeySize   = 65535
	maxValueSize = 2147483647
)

type fileIndex struct {
	bytes  []byte
	offset uint64
	size   uint16
}

// TODO : consider using unsafe ... at least to test performance gain
func newFileIndex(offset int, size int) (fileIndex, error) {

	if size > maxKeySize {
		return fileIndex{}, fmt.Errorf("cannot store key of size bigger than %d. size was %d", maxKeySize, size)
	}
	ss := uint16(size)

	if offset > maxValueSize {
		return fileIndex{}, fmt.Errorf("cannot store value of size bigger than %d. size was %d", maxValueSize, offset)
	}
	oo := uint32(offset)

	s := make([]byte, 2)
	binary.LittleEndian.PutUint16(s, ss)
	o := make([]byte, 4)
	binary.LittleEndian.PutUint32(o, oo)
	b, err := Concat(6, o, s)
	if err != nil {
		return fileIndex{}, fmt.Errorf("could not create fileIndex for [%d,%d] %w", offset, size, err)
	}
	return fileIndex{
		bytes:  b,
		offset: uint64(oo),
		size:   ss,
	}, nil
}

func readIndex(b []byte) (fileIndex, error) {
	if len(b) != 6 {
		return fileIndex{}, fmt.Errorf("cannot read size from fileIndex %v", b)
	}
	o := binary.LittleEndian.Uint32(b[:4])
	s := binary.LittleEndian.Uint16(b[4:])
	return fileIndex{
		bytes:  b,
		offset: uint64(o),
		size:   s,
	}, nil
}

func (i fileIndex) Offset() int64 {
	return int64(i.offset)
}

func (i fileIndex) Size() int {
	return int(i.size)
}
