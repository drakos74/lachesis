package bash

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/drakos74/lachesis/internal/app/store"
)

// DB is the storage implementation backed by a simple bash script
type DB struct {
	path string
}

// DBFactory create a DB storage implementation
func DBFactory(path string) store.StorageFactory {
	return func() store.Storage {
		println(fmt.Sprintf("path = %v", path))
		return DB{path: path}
	}
}

// Put adds an element to the Bash store
func (b DB) Put(element store.Element) error {
	println(fmt.Sprintf("element = %v", element))
	cmd := exec.Command("bash", fmt.Sprintf("%s%s.sh", b.path, "db_set"), string(element.Key), string(element.Value))
	var out bytes.Buffer
	cmd.Stdout = &out
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	return cmd.Run()
}

// Get performs a value retrieval from th Bash store based on the given key
func (b DB) Get(key store.Key) (store.Element, error) {
	cmd := exec.Command("bash", fmt.Sprintf("%s%s.sh", b.path, "db_get"), string(key))
	var out bytes.Buffer
	cmd.Stdout = &out
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	value := out.String()
	if len(value) > 0 {
		return store.NewElement(key, []byte(value[0:len(value)-1])), nil
	}
	return store.Nil, err
}

// Metadata returns the internal metadata of the given storage implementation
func (b DB) Metadata() store.Metadata {
	return store.Metadata{}
}

// Close closes the DB storage and performs any needed cleanup
func (b DB) Close() error {
	// nothing to do
	return nil
}
