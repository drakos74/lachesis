package bytes

import (
	"encoding/binary"
	"fmt"
)

// Concat merges 2 arrays of bytes into one
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
func FileIndex(offset int, size int) (fileIndex, error) {

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

func ReadIndex(b []byte) (fileIndex, error) {
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

func (i fileIndex) Bytes() []byte {
	return i.bytes
}

func (i fileIndex) Offset() int64 {
	return int64(i.offset)
}

func (i fileIndex) Size() int {
	return int(i.size)
}
