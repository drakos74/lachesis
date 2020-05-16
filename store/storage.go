package store

import "github.com/drakos74/lachesis/model"

// Storage is the low level interface for interacting with the underlying storage in bytes
type Storage interface {
	Put(element model.Element) error
	Get(element model.Element) (model.Element, error)
	Metadata() model.Metadata
	Close() error
}

// Encoder is the implementation transforming the incoming objects to the byte slice
type Encoder func(v interface{}) (model.Element, error)

// Decoder is the implementation for transforming the byte slice to the appropriate struct
type Decoder func(element model.Element) (interface{}, error)

// Repository allows us to transform the raw byte representation into the model objects used by the application
type Repository struct {
	Encoder
	Decoder
	Storage
}
