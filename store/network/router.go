package network

type Key []byte

type Switch interface {
	Register(id int)
	DeRegister(id int)
	Route(key Key) ([]int, error)
}

type PartitionStrategy func() Switch

func SinglePartition() Switch {
	return &UnarySwitch{}
}

type UnarySwitch struct {
}

func (u *UnarySwitch) Register(id int) {
	// nothing to do
}

func (u *UnarySwitch) DeRegister(id int) {
	// nothing to do
}

func (u UnarySwitch) Route(key Key) ([]int, error) {
	return []int{0}, nil
}
