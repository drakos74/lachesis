package bash

import (
	"fmt"
	"testing"
	"time"

	"github.com/drakos74/lachesis/store/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestBash_Put(t *testing.T) {
	bash := DB{}

	element := storage.NewElement([]byte(fmt.Sprintf("key-%d", time.Now().Unix())), []byte(fmt.Sprintf("value-%d", time.Now().Unix())))

	err := bash.Put(element)
	assert.NoError(t, err)

	e, err := bash.Get(element.Key)
	assert.Equal(t, element, e)
	assert.NoError(t, err)

}

func TestBash_PutOverwrite(t *testing.T) {
	bash := DB{}

	key := fmt.Sprintf("key-%d", time.Now().Unix())

	value1 := fmt.Sprintf("value-%d", time.Now().Unix())

	value2 := fmt.Sprintf("value-2-%d", time.Now().Unix())

	element := storage.NewElement([]byte(key), []byte(value1))

	err := bash.Put(element)
	assert.NoError(t, err)

	e, err := bash.Get([]byte(key))
	assert.Equal(t, element, e)
	assert.NoError(t, err)

	element2 := storage.NewElement([]byte(key), []byte(value2))

	err2 := bash.Put(element2)
	assert.NoError(t, err2)

	e2, err := bash.Get([]byte(key))
	assert.Equal(t, element2, e2)
	assert.NoError(t, err)

}
