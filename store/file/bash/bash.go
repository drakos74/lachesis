package bash

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/drakos74/lachesis/store"
)

type BashDB struct {
}

func BashDBFactory() store.Storage {
	return BashDB{}
}

func (b BashDB) Put(element store.Element) error {
	cmd := exec.Command("bash", fmt.Sprintf("%s%s.sh", "", "db_set"), string(element.Key), string(element.Value))
	var out bytes.Buffer
	cmd.Stdout = &out
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	return cmd.Run()
}

func (b BashDB) Get(key store.Key) (store.Element, error) {
	cmd := exec.Command("bash", fmt.Sprintf("%s%s.sh", "", "db_get"), string(key))
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

func (b BashDB) Metadata() store.Metadata {
	panic("implement me")
}

func (b BashDB) Close() error {
	panic("implement me")
}
