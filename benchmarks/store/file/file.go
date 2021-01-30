package file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/rs/zerolog/log"
)

// Storage is the file storage implementation for the Storage interface
type Storage struct {
	wrFile *os.File
	rdFile *os.File
	// we use this to encapsulate our concatenation logic
	concat   storage.ConcatOperator
	filename string
}

// Put adds an element to the file storage
func (f Storage) Put(element storage.Element) error {
	bytes, err := f.concat.Join(element)
	if err != nil {
		return fmt.Errorf("could not noncat element '%v' %w", element, err)
	}
	n, err := f.wrFile.Write(bytes)
	if err != nil {
		return fmt.Errorf("could not write element '%v' %w", element, err)
	}
	if n != len(bytes) {
		return fmt.Errorf("write failed '%d' != %d", n, len(bytes))
	}
	return nil
}

// Get retrieves a value from the file storage based on the given key
func (f Storage) Get(key storage.Key) (storage.Element, error) {
	scanner := bufio.NewScanner(f.rdFile)
	for scanner.Scan() {
		result, err := f.concat.Split(key, scanner.Bytes())
		if err != nil {
			return storage.Nil, fmt.Errorf("error during deserialisation: %w", err)
		}
		if bytes.Compare(result.Key, key) == 0 {
			// Note : overwrite will fail, as we are returning the first match
			return result, nil
		}
	}
	return storage.Nil, fmt.Errorf(storage.NoValue, key)
}

// Metadata returns the internal stats of the file storage implementation
func (f Storage) Metadata() storage.Metadata {

	var size uint64

	scanner := bufio.NewScanner(f.rdFile)
	for scanner.Scan() {
		atomic.AddUint64(&size, 1)
	}

	return storage.Metadata{
		Size: size,
	}
}

// Close closes the files related to the storage implementation
func (f Storage) Close() error {
	wrErr := f.wrFile.Close()
	rdErr := f.rdFile.Close()

	if wrErr != nil || rdErr != nil {
		return fmt.Errorf("could not close ScratchPad [%v,%v]", wrErr, rdErr)
	}

	return nil
}

// NewFileStorage creates a new Storage storage instance
func NewFileStorage(path string) (*Storage, error) {
	// generate a file name
	fileName := fmt.Sprintf("%s/%s.%s", path, strconv.FormatInt(time.Now().UnixNano(), 10), "lac")
	wrFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not create write file for storage %w", err)
	}
	rdFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not create read file for storage %w", err)
	}
	log.Debug().
		Str("filename", wrFile.Name()).
		Msg("Open ScratchPad Storage")
	return &Storage{wrFile: wrFile, rdFile: rdFile, concat: storage.RawConcat(), filename: fileName}, nil
}

// StorageFactory generates a file storage implementation
func StorageFactory(path string) storage.StorageFactory {
	return func() storage.Storage {
		pad, err := NewFileStorage(path)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return pad
	}
}
