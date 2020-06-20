package file

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/drakos74/lachesis/store"
	"github.com/drakos74/lachesis/store/mem"

	"github.com/rs/zerolog/log"
)

// ScratchPad is a single file wrapper for storing key value pairs
// it uses a Trie for storing the keys as an fileIndex for the file
type ScratchPad struct {
	wrFile *os.File
	rdFile *os.File
	// we store in the fileIndex a slice of bytes representing the stored object [Size,Size]
	index *mem.SyncTrie
	// we use this to encapsulate our concatenation logic
	concat   concat
	offset   int
	filename string
}

// NewScratchPad creates a new ScratchPad instance
func NewScratchPad(path string) (*ScratchPad, error) {
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
	index := mem.NewSyncTrie()
	log.Debug().
		Str("filename", wrFile.Name()).
		Msg("Open ScratchPad Storage")
	return &ScratchPad{wrFile: wrFile, rdFile: rdFile, concat: newConcat(), index: index, filename: fileName}, nil
}

// ScratchPadFactory generates a file storage implementation
func ScratchPadFactory(path string) store.StorageFactory {
	return func() store.Storage {
		pad, err := NewScratchPad(path)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return pad
	}
}

// Put adds an element to the store
func (s *ScratchPad) Put(element store.Element) error {
	bytes, err := s.concat.join(element)
	if err != nil {
		return fmt.Errorf("could not serialize element '%v' %w", element, err)
	}
	// Note : we leave the overwrites there ... just applying a new fileIndex !!!
	// We will silently remove them at the next 'compaction' operation
	n, err := s.wrFile.Write(bytes)
	if err != nil {
		return fmt.Errorf("could not write element '%v' %w", element, err)
	}
	// TODO : seems we dont need to call 'sync' in order to flush to the file...
	//  need to investigate the low level implications of this
	//defer func(err errors) {
	//	syncErr := store.wrFile.Sync()
	//	if syncErr != nil {
	//		err.append(syncErr)
	//	}
	//}(store.err)
	if n != len(bytes) {
		// TODO : handle the file corruption -> open new file
		return fmt.Errorf("write failed '%d' != %d", n, len(bytes))
	}

	index, err := newFileIndex(s.offset, n)
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
	return s.index.Put(store.NewElement(element.Key, index.bytes))
}

// Get retrieves the element corresponding to the provided key
// if a value is not found, it will return an error
func (s *ScratchPad) Get(key store.Key) (store.Element, error) {
	keyBytes, err := s.index.Get(key)
	if err != nil {
		return store.Element{}, fmt.Errorf(store.NoValue, key)
	}

	index, err := readIndex(keyBytes.Value)
	if err != nil {
		return store.Element{}, fmt.Errorf("cannot read fileIndex '%v' %w", index, err)
	}

	log.Trace().
		Int64("offset", index.Offset()).
		Int("Size", index.Size()).
		Bytes("key", keyBytes.Key).
		Msg("Read_Index")

	data := make([]byte, index.Size())
	n, err := s.rdFile.ReadAt(data, index.Offset())
	if err != nil {
		return store.Element{}, fmt.Errorf("cannot read at '%d' keyBytes '%d' found '%d' %w", index.Offset(), index.Size(), n, err)
	}
	if n != index.Size() {
		return store.Element{}, fmt.Errorf("cannot read at '%d' keyBytes '%d' found '%d'", index.Offset(), index.Size(), n)
	}
	result, err := s.concat.split(key, data)
	if err != nil {
		return store.Element{}, fmt.Errorf("cannot deserialize '%v' %w", data, err)
	}
	return result, nil
}

// Metadata returns internal statistics about the storage
// It s not meant to serve anny functionality, but used only for testing
func (s *ScratchPad) Metadata() store.Metadata {
	file, _ := os.Open(s.filename)
	fileScanner := bufio.NewScanner(file)
	l := 0
	var b uint64
	for fileScanner.Scan() {
		l++
		b += uint64(len(fileScanner.Bytes()))
	}
	keyMetadata := s.index.Metadata()
	return store.Metadata{
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
