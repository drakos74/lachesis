package file

import (
	"encoding/binary"
	"fmt"
	"lachesis/internal/model"
	"lachesis/internal/store/trie"
	"lachesis/pkg"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

//type errors []error
//
//func (err *errors) append(currentErr error) {
//	*err = append(*err, currentErr)
//}

// SB is a single file wrapper for storing key value pairs
// it uses a Trie for storing the keys as an index for the file
type SB struct {
	wrFile *os.File
	rdFile *os.File
	//err    errors
	// we store in the index a slice of bytes representing the stored object [size,size]
	index trie.Trie
	// we use this to encapsulate our serialization / deserialization logic
	serdes   model.Serdes
	size     int
	filename string
}

// Close closes the file and completes all clean-up operations needed
func (sb *SB) Close() error {

	log.Debug().
		Str("filename", sb.wrFile.Name()).
		Int("size", sb.size).
		Msg("Close ScratchPad Storage")

	wrErr := sb.wrFile.Close()
	rdErr := sb.rdFile.Close()

	if wrErr != nil || rdErr != nil {
		return fmt.Errorf("could not close SB [%v,%v]", wrErr, rdErr)
	}

	return nil
}

// TODO : make a builder
// NewTrie creates a new SB instance
func New(path string) (*SB, error) {
	// TODO : make the randmness better and dont let it overflow
	fileName := fmt.Sprintf("%s/%s.%s", path, strconv.FormatInt(time.Now().UnixNano(), 10), "lac")
	wrFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not create write file for SB %w", err)
	}
	rdFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not create read file for SB %w", err)
	}
	// TODO : we should build the index from the file ... but for that we need to store also the key
	index := trie.NewTrie('0')
	log.Debug().
		Str("filename", wrFile.Name()).
		Msg("Open ScratchPad Storage")
	return &SB{wrFile: wrFile, rdFile: rdFile, serdes: createSerdes(), index: index, filename: fileName}, nil
}

// Put adds an element to the store
func (sb *SB) Put(element model.Element) error {
	bytes, err := sb.serdes.Serializer(element)
	if err != nil {
		return fmt.Errorf("could not serialize element '%v' %w", element, err)
	}
	// Note : we leave the overwrites there ... just applying a new index !!!
	n, err := sb.wrFile.Write(bytes)
	if err != nil {
		return fmt.Errorf("could not write element '%v' %w", element, err)
	}
	// TODO : seems we dont need to call 'sync' in order to flush to the file... need to investigate the low level implications of this
	//defer func(err errors) {
	//	syncErr := sb.wrFile.Sync()
	//	if syncErr != nil {
	//		err.append(syncErr)
	//	}
	//}(sb.err)
	if n != len(bytes) {
		// TODO : handle the file corruption -> open new file
		return fmt.Errorf("write failed '%d' != %d", n, len(bytes))
	}
	index, err := createIndex(sb.size, n)
	if err != nil {
		return fmt.Errorf("could not create index '%v' %w", index, err)
	}

	log.Trace().
		Int64("size", index.Offset()).
		Int("size", index.Size()).
		Bytes("key", element.Key()).
		Msg("Write_Index")

	sb.size += n

	log.Trace().
		Int("size", sb.size).
		Msg("ScratchPad")

	return sb.index.Commit(element.Key(), index.bytes)
}

// Get retrieves the element corresponding to the provided key
// if a value is not found, it will return an error
func (sb *SB) Get(element model.Element) (model.Element, error) {
	if bytes, ok := sb.index.Read(element.Key()); ok {

		index, err := readIndex(bytes)
		if err != nil {
			return nil, fmt.Errorf("cannot read index '%v' %w", index, err)
		}

		log.Trace().
			Int64("size", index.Offset()).
			Int("size", index.Size()).
			Bytes("key", element.Key()).
			Msg("Read_Index")

		data := make([]byte, index.Size())
		n, err := sb.rdFile.ReadAt(data, index.Offset())
		if err != nil {
			return nil, fmt.Errorf("cannot read at '%d' bytes '%d' found '%d' %w", index.Offset(), index.Size(), n, err)
		}
		if n != index.Size() {
			return nil, fmt.Errorf("cannot read at '%d' bytes '%d' found '%d'", index.Offset(), index.Size(), n)
		}
		result, err := sb.serdes.Deserializer(element, data)
		if err != nil {
			return nil, fmt.Errorf("cannot deserialize '%v' %w", data, err)
		}
		return result, nil
	}
	return nil, fmt.Errorf("cannot find index for element '%v'", element)
}

// Handle the SB index

const (
	maxKeySize   = 65535
	maxValueSize = 2147483647
)

type index struct {
	bytes  []byte
	offset uint64
	size   uint16
}

// TODO : consider using unsafe ... at least to test performance gain
func createIndex(offset int, size int) (index, error) {

	if size > maxKeySize {
		return index{}, fmt.Errorf("cannot store key of size bigger than %d. size was %d", maxKeySize, size)
	}
	ss := uint16(size)

	if offset > maxValueSize {
		return index{}, fmt.Errorf("cannot store value of size bigger than %d. size was %d", maxValueSize, offset)
	}
	oo := uint32(offset)

	s := make([]byte, 2)
	binary.LittleEndian.PutUint16(s, ss)
	o := make([]byte, 4)
	binary.LittleEndian.PutUint32(o, oo)
	b, err := pkg.Concat(6, o, s)
	if err != nil {
		return index{}, fmt.Errorf("could not create index for [%d,%d] %w", offset, size, err)
	}
	return index{
		bytes:  b,
		offset: uint64(oo),
		size:   ss,
	}, nil
}

func readIndex(b []byte) (index, error) {
	if len(b) != 6 {
		return index{}, fmt.Errorf("cannot read size from index %v", b)
	}
	o := binary.LittleEndian.Uint32(b[:4])
	s := binary.LittleEndian.Uint16(b[4:])
	return index{
		bytes:  b,
		offset: uint64(o),
		size:   s,
	}, nil
}

func (i index) Offset() int64 {
	return int64(i.offset)
}

func (i index) Size() int {
	return int(i.size)
}

// Handle the internal serialization

func createSerdes() model.Serdes {
	nl := []byte{byte('\n')}
	return model.Serdes{
		Serializer: func(element model.Element) ([]byte, error) {
			b, err := pkg.Concat(len(element.Value())+len(nl), element.Value(), nl)
			if err != nil {
				return nil, fmt.Errorf("could not serialize value %w", err)
			}
			return b, nil
		},
		Deserializer: func(element model.Element, data []byte) (model.Element, error) {
			n := len(data) - len(nl)
			return model.NewObject(element.Key(), data[0:n]), nil
		},
	}
}
