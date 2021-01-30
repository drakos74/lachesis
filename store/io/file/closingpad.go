package file

import (
	"fmt"
	"github.com/drakos74/lachesis"

	"github.com/drakos74/lachesis/io/bytes"
	"github.com/drakos74/lachesis/io/mem"
	"github.com/rs/zerolog/log"
)

type ClosingPad struct {
	ScratchPad
}

// TrieClosingPadFactory generates a file storage implementation
// with a trie as an index
func TrieClosingPadFactory(path string) store.StorageFactory {
	return func() store.Storage {
		pad, err := NewScratchPad(path, mem.SyncTrieFactory)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return &ClosingPad{*pad}
	}
}

// TreeClosingPadFactory generates a file storage implementation
// with a btree as an index
func TreeClosingPadFactory(path string) store.StorageFactory {
	return func() store.Storage {
		pad, err := NewScratchPad(path, mem.SyncBTreeFactory)
		if err != nil {
			panic(fmt.Sprintf("error during store creation: %v", err))
		}
		return &ClosingPad{*pad}
	}
}

// Put adds an element to the store
func (s *ClosingPad) Put(element store.Element) error {
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
	//  need to investigate the low level implications of this
	defer func() {
		syncErr := s.wrFile.Sync()
		if syncErr != nil {
			log.Err(syncErr)
		}
	}()
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
	return s.index.Put(store.NewElement(element.Key, index.Bytes()))
}
