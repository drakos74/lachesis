package file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/drakos74/lachesis/store"
)

type FileStorage struct {
	wrFile *os.File
	rdFile *os.File
	// we use this to encapsulate our concatenation logic
	concat   concat
	filename string
}

func (f FileStorage) Put(element store.Element) error {
	bytes, err := f.concat.join(element)
	if err != nil {
		return fmt.Errorf("could not serialize element '%v' %w", element, err)
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

func (f FileStorage) Get(key store.Key) (store.Element, error) {
	scanner := bufio.NewScanner(f.rdFile)
	for scanner.Scan() {
		result, err := f.concat.split(key, scanner.Bytes())
		if err != nil {
			return store.Nil, fmt.Errorf("error during deserialisation: %w", err)
		}
		if bytes.Compare(result.Key, key) == 0 {
			// Note : overwrite will fail, as we are returning the first match
			return result, nil
		}
	}
	return store.Nil, fmt.Errorf(store.NoValue, key)
}

func (f FileStorage) Metadata() store.Metadata {

	var size uint64

	scanner := bufio.NewScanner(f.rdFile)
	for scanner.Scan() {
		atomic.AddUint64(&size, 1)
	}

	return store.Metadata{
		Size: size,
	}
}

func (f FileStorage) Close() error {
	wrErr := f.wrFile.Close()
	rdErr := f.rdFile.Close()

	if wrErr != nil || rdErr != nil {
		return fmt.Errorf("could not close ScratchPad [%v,%v]", wrErr, rdErr)
	}

	return nil
}

// NewFile creates a new FileStorage storage instance
func NewFileStorage(path string) (*FileStorage, error) {
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
	return &FileStorage{wrFile: wrFile, rdFile: rdFile, concat: newRawConcat(), filename: fileName}, nil
}

// FileStorageFactory generates a file storage implementation
func FileStorageFactory(path string) store.StorageFactory {
	return func() store.Storage {
		pad, err := NewFileStorage(path)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return pad
	}
}
