package mem

import (
	"testing"

	"github.com/drakos74/lachesis/store"

	"github.com/drakos74/lachesis/store/test"
)

func TestCacheImplementation(t *testing.T) {
	test.Execute(t, func() store.Storage {
		return NewCache()
	})
}
