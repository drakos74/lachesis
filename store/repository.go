package store

import (
	"encoding/json"
	"fmt"
)

// KV is a generic key-value structure
type KV struct {
	Key   interface{}
	Value interface{}
}

// Encoder is the implementation transforming the incoming objects to the byte slice
type Encoder func(v KV) (Element, error)

// Decoder is the implementation for transforming the byte slice to the appropriate struct
type Decoder func(element Element) (KV, error)

// ValueConstructor re-creates the value object from the serialized form
type ValueConstructor func() interface{}

// Repository is the high level implementation  allowing to store arbitrary key value pairs
type Repository struct {
	encoder   Encoder
	decoder   Decoder
	generator ValueConstructor
	storage   Storage
}

// NewRepository creates a new repository
func NewRepository(constructor ValueConstructor, storage Storage) *Repository {
	return &Repository{
		encoder: encode,
		decoder: decode(constructor),
		storage: storage,
	}
}

// Put puts a key value pair into the repository
func (repo *Repository) Put(kv KV) error {
	element, err := repo.encoder(kv)
	if err != nil {
		return fmt.Errorf("could not put %v: %w", kv, err)
	}
	return repo.storage.Put(element)
}

// Get retrieves a value from the repository for the given key
func (repo *Repository) Get(k interface{}) (KV, error) {
	key, err := encode(KV{Key: k})
	if err != nil {
		return KV{}, fmt.Errorf("could not encode key %v: %w", k, err)
	}
	// TODO : fix this
	element, err := repo.storage.Get(key.Key)
	if err != nil {
		return KV{}, fmt.Errorf("could not retrieve element for key %v: %w", k, err)
	}
	return repo.decoder(element)
}

// Metadata returns the repository metadata
func (repo *Repository) Metadata() Metadata {
	return repo.storage.Metadata()
}

// Close closes the repository
func (repo *Repository) Close() error {
	return repo.storage.Close()
}

func encode(kv KV) (Element, error) {
	bk, err := json.Marshal(kv.Key)
	if err != nil {
		return Element{}, fmt.Errorf("could not marshall key %v: %v", kv.Key, err)
	}
	bv, err := json.Marshal(kv.Value)
	if err != nil {
		return Element{}, fmt.Errorf("could not marshall value %v: %v", kv.Value, err)
	}
	return NewElement(bk, bv), nil
}

func decode(constr ValueConstructor) func(element Element) (KV, error) {
	return func(element Element) (kv KV, e error) {
		response := constr()
		err := json.Unmarshal(element.Value, response)
		return KV{Value: response}, err
	}
}
