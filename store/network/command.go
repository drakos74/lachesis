package network

import "github.com/drakos74/lachesis/store"

type CmdType int

const (
	Get CmdType = iota + 1
	Put
)

type Command interface {
	Type() CmdType
	Element() store.Element
	Exec() func(storage store.Storage) (store.Element, error)
}

type Response struct {
	store.Element
	Err error
}

type PutCommand struct {
	element store.Element
}

func (p PutCommand) Type() CmdType {
	return Put
}

func (p PutCommand) Element() store.Element {
	return p.element
}

func (p PutCommand) Exec() func(storage store.Storage) (store.Element, error) {
	return func(storage store.Storage) (element store.Element, e error) {
		err := storage.Put(p.Element())
		return store.Nil, err
	}
}

type GetCommand struct {
	key store.Key
}

func (p GetCommand) Type() CmdType {
	return Get
}

func (p GetCommand) Element() store.Element {
	return store.NewElement(p.key, store.NilBytes)
}

func (p GetCommand) Exec() func(storage store.Storage) (store.Element, error) {
	return func(storage store.Storage) (element store.Element, e error) {
		return storage.Get(p.key)
	}
}
