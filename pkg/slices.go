package pkg

import "fmt"

func Concat(size int, arrays ...[]byte) ([]byte, error) {
	arr := make([]byte, size)
	i := 0
	for _, array := range arrays {
		for _, b := range array {
			arr[i] = b
			i++
		}
	}
	if i != size {
		return nil, fmt.Errorf("size argument does not match %d vs %d", size, i)
	}
	return arr, nil
}
