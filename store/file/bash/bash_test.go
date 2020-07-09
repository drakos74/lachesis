package bash

import (
	"fmt"
	"testing"
	"time"

	"github.com/drakos74/lachesis/store"

	"github.com/stretchr/testify/assert"
)

func TestBash_Put(t *testing.T) {
	bash := DB{}

	element := store.NewElement([]byte(fmt.Sprintf("key-%d", time.Now().Unix())), []byte(fmt.Sprintf("value-%d", time.Now().Unix())))

	err := bash.Put(element)
	assert.NoError(t, err)

	e, err := bash.Get(element.Key)
	assert.Equal(t, element, e)
	assert.NoError(t, err)

}
