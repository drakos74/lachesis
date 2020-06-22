package file

import (
	"testing"

	"github.com/drakos74/lachesis/store/test"
)

func TestFile_SimpleImplementation(t *testing.T) {
	new(test.Consistency).Run(t, FileStorageFactory("data"))
}
