package network

import (
	"github.com/drakos74/lachesis/store"
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
	Element() store.Element
}

// Response is the result of the Command
type Response struct {
	store.Element
	Err error
}

// PutCommand represents a put action
type PutCommand struct {
	element store.Element
}

// NewPut creates a new PutCommand
func NewPut(element store.Element) Command {
	return PutCommand{element: element}
}

// Type returns the the type for the PutCommand
func (p PutCommand) Type() CmdType {
	return Put
}

// Element retrieves the element that needs to be written to the Storage
func (p PutCommand) Element() store.Element {
	return p.element
}

// GetCommand represents a get action
type GetCommand struct {
	key store.Key
}

// NewGet creates a new GetCommand
func NewGet(key store.Key) Command {
	return GetCommand{key: key}
}

// Type returns the the type for the PutCommand
func (p GetCommand) Type() CmdType {
	return Get
}

// Element retrieves the element which key needs to be looked up and returned from the Storage
func (p GetCommand) Element() store.Element {
	return store.NewElement(p.key, store.NilBytes)
}
