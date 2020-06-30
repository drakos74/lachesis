package network

import (
	"github.com/drakos74/lachesis/store"
)

type CmdType int

const (
	Get CmdType = iota + 1
	Put
)

type Command interface {
	Type() CmdType
	Element() store.Element
}

type Response struct {
	store.Element
	Err error
}

type PutCommand struct {
	element store.Element
}

func NewPut(element store.Element) Command {
	return PutCommand{element: element}
}

func (p PutCommand) Type() CmdType {
	return Put
}

func (p PutCommand) Element() store.Element {
	return p.element
}

type GetCommand struct {
	key store.Key
}

func NewGet(key store.Key) Command {
	return GetCommand{key: key}
}

func (p GetCommand) Type() CmdType {
	return Get
}

func (p GetCommand) Element() store.Element {
	return store.NewElement(p.key, store.NilBytes)
}
