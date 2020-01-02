package store

import (
	"fmt"
	"lachesis/internal/model"
	"lachesis/internal/store/file"
	"lachesis/internal/store/mem"
	"log"
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

var result model.Element

func randomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func createElement(k, v int) model.Element {
	return model.NewObject(randomBytes(k), randomBytes(v))
}

func generate(n, k, v int) []model.Element {
	elements := make([]model.Element, 0)
	for i := 0; i < n; i++ {
		element := createElement(k, v)
		elements = append(elements, element)
	}
	return elements
}

func put(storage Storage, elements []model.Element) {
	for _, element := range elements {
		err := storage.Put(element)
		if err != nil {
			log.Fatal()
		}
	}
}

func read(storage Storage, elements []model.Element) {
	for _, element := range elements {
		result, err := storage.Get(element)
		if err != nil {
			log.Fatal()
		}
		consume(result)
	}
}

func consume(element model.Element) {
	//... store in a environment variable to avoid optimizations https://stackoverflow.com/questions/36966947/do-go-testing-b-benchmarks-prevent-unwanted-optimizations
	result = element
}

func BenchmarkSB(b *testing.B) {
	executeBenchmark(b, func() Storage {
		db, err := file.New("../../test/testdata/bench")
		if err != nil {
			log.Fatal()
		}
		return db
	})
}

func BenchmarkMemory(b *testing.B) {
	executeBenchmark(b, func() Storage {
		return mem.NewCache()
	})
}

func BenchmarkTrie(b *testing.B) {
	executeBenchmark(b, func() Storage {
		return mem.NewTrie()
	})
}

func executeBenchmark(b *testing.B, store func() Storage) {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	for k := 0; k <= 3; k++ {

		for l := 10; l <= 10; l++ {

			storage := store()

			num := int(math.Pow10(k))
			key := 10 * l
			value := key * 10

			elements := generate(num, key, value)

			b.Run(fmt.Sprintf("%s:%s/num:%d,size-key:%d,size-value:%d", reflect.TypeOf(storage).String(), "put", num, key, value), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					put(storage, elements)
				}
			})

			b.Run(fmt.Sprintf("%s:%s/num:%d,size-key:%d,size-value:%d", reflect.TypeOf(storage).String(), "get", num, key, value), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					read(storage, elements)
				}
			})

			err := storage.Close()
			if err != nil {
				log.Fatal()
			}

		}

	}

}
