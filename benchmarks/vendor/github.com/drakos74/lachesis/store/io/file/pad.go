package file

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/drakos74/lachesis/store/io/bytes"
	"github.com/drakos74/lachesis/store/io/mem"
	"github.com/rs/zerolog/log"
)

// ScratchPad is a single file wrapper for storing key value pairs
// it uses a Trie for storing the keys as a fileIndex for the file
type ScratchPad struct {
	wrFile *os.File
	rdFile *os.File
	// we store in the fileIndex a slice of bytes representing the stored object [Size,Size]
	index storage.Storage
	// we use this to encapsulate our concatenation logic
	concat   storage.ConcatOperator
	offset   int
	filename string
}

// NewScratchPad creates a new ScratchPad instance
func NewScratchPad(path string, index storage.StorageFactory) (*ScratchPad, error) {
	// generate a file name
	fileName := fmt.Sprintf("%s/%s.%s", path, strconv.FormatInt(time.Now().UnixNano(), 10), "lac")
	wrFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not create write file for ScratchPad %w", err)
	}
	rdFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not create read file for ScratchPad %w", err)
	}
	// TODO : we should build the fileIndex from the file ...
	//  to make the storage consistent
	//  but for that we need to store also the key
	log.Debug().
		Str("filename", wrFile.Name()).
		Msg("Open ScratchPad Storage")
	return &ScratchPad{wrFile: wrFile, rdFile: rdFile, concat: storage.IndexedConcat(), index: index(), filename: fileName}, nil
}

// TriePadFactory generates a file storage implementation
// with a trie as an index
func TriePadFactory(path string) storage.StorageFactory {
	return func() storage.Storage {
		pad, err := NewScratchPad(path, mem.SyncTrieFactory)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return pad
	}
}

// TreePadFactory generates a file storage implementation
// with a btree as an index
func TreePadFactory(path string) storage.StorageFactory {
	return func() storage.Storage {
		pad, err := NewScratchPad(path, mem.SyncBTreeFactory)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return pad
	}
}

// Put adds an element to the store
func (s *ScratchPad) Put(element storage.Element) error {
	bb, err := s.concat.Join(element)
	if err != nil {
		return fmt.Errorf("could not serialize element '%v' %w", element, err)
	}
	// Note : we leave the overwrites there ... just applying a new fileIndex !!!
	// We will silently remove them at the next 'compaction' operation
	n, err := s.wrFile.Write(bb)
	if err != nil {
		return fmt.Errorf("could not write element '%v' %w", element, err)
	}
	// TODO : seems we dont need to call 'sync' in order to flush to the file...
	//  while we write and read from the same process
	//  need to investigate the low level implications of this
	//  (could be because go uses an mmap under the curtains for file operations)
	//defer func(err errors) {
	//	syncErr := storage.wrFile.Sync()
	//	if syncErr != nil {
	//		err.append(syncErr)
	//	}
	//}(storage.err)
	if n != len(bb) {
		// TODO : handle the file corruption -> open new file
		return fmt.Errorf("write failed '%d' != %d", n, len(bb))
	}

	index, err := bytes.FileIndex(s.offset, n)
	if err != nil {
		return fmt.Errorf("could not create fileIndex '%v' %w", index, err)
	}
	s.offset += n

	log.Trace().
		Int64("offset", index.Offset()).
		Int("Size", index.Size()).
		Bytes("key", element.Key).
		Msg("Write_Index")
	// Note : we overwrite the element only in the key struct,
	// so the old value is not reachable from the outside world
	return s.index.Put(storage.NewElement(element.Key, index.Bytes()))
}

// Get retrieves the element corresponding to the provided key
// if a value is not found, it will return an error
func (s *ScratchPad) Get(key storage.Key) (storage.Element, error) {
	bb, err := s.index.Get(key)
	if err != nil {
		return storage.Element{}, fmt.Errorf(storage.NoValue, key)
	}

	index, err := bytes.ReadIndex(bb.Value)
	if err != nil {
		return storage.Element{}, fmt.Errorf("cannot read fileIndex '%v' %w", index, err)
	}

	log.Trace().
		Int64("offset", index.Offset()).
		Int("Size", index.Size()).
		Bytes("key", bb.Key).
		Msg("Read_Index")

	data := make([]byte, index.Size())
	n, err := s.rdFile.ReadAt(data, index.Offset())
	if err != nil {
		return storage.Element{}, fmt.Errorf("cannot read at '%d' bb '%d' found '%d' %w", index.Offset(), index.Size(), n, err)
	}
	if n != index.Size() {
		return storage.Element{}, fmt.Errorf("cannot read at '%d' bb '%d' found '%d'", index.Offset(), index.Size(), n)
	}
	result, err := s.concat.Split(key, data)
	if err != nil {
		return storage.Element{}, fmt.Errorf("cannot deserialize '%v' %w", data, err)
	}
	return result, nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (s *ScratchPad) Metadata() storage.Metadata {
	file, _ := os.Open(s.filename)
	fileScanner := bufio.NewScanner(file)
	l := 0
	var b uint64
	for fileScanner.Scan() {
		l++
		b += uint64(len(fileScanner.Bytes()))
	}
	keyMetadata := s.index.Metadata()
	return storage.Metadata{
		Size:        keyMetadata.Size,
		KeysBytes:   keyMetadata.ValuesBytes + keyMetadata.KeysBytes,
		ValuesBytes: b,
		Errors:      make([]error, 0),
	}
}

// Close closes the file and completes all clean-up operations needed
func (s *ScratchPad) Close() error {

	log.Debug().
		Str("filename", s.wrFile.Name()).
		Int("Offset", s.offset).
		Msg("Close ScratchPad Storage")

	wrErr := s.wrFile.Close()
	rdErr := s.rdFile.Close()

	if wrErr != nil || rdErr != nil {
		return fmt.Errorf("could not close ScratchPad [%v,%v]", wrErr, rdErr)
	}

	return nil
}
