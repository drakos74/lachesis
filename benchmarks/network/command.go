package network

import (
	"github.com/drakos74/lachesis/store/app/storage"
)

// CmdType specifies the type of command
type CmdType int

const (
	// Get represents a 'get' request
	Get CmdType = iota + 1
	// Put represents a 'put' request
	Put
)

// Command is the interface for every command object
type Command interface {
	Type() CmdType
	Element() storage.Element
}

// Response is the result of the Command
type Response struct {
	storage.Element
	Err error
}

// PutCommand represents a put action
type PutCommand struct {
	element storage.Element
}

// NewPut creates a new PutCommand
func NewPut(element storage.Element) Command {
	return PutCommand{element: element}
}

// Type returns the the type for the PutCommand
func (p PutCommand) Type() CmdType {
	return Put
}

// Element retrieves the element that needs to be written to the Storage
func (p PutCommand) Element() storage.Element {
	return p.element
}

// GetCommand represents a get action
type GetCommand struct {
	key storage.Key
}

// NewGet creates a new GetCommand
func NewGet(key storage.Key) Command {
	return GetCommand{key: key}
}

// Type returns the the type for the PutCommand
func (p GetCommand) Type() CmdType {
	return Get
}

// Element retrieves the element which key needs to be looked up and returned from the Storage
func (p GetCommand) Element() storage.Element {
	return storage.NewElement(p.key, storage.NilBytes)
}
