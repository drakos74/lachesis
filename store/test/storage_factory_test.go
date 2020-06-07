package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {

	var bytes []byte
	for i := 0; i < 10000; i++ {
		newBytes := RandomBytes(10)
		assert.NotEqual(t, bytes, newBytes)
		bytes = newBytes
		assert.Equal(t, bytes, newBytes)
	}

}
