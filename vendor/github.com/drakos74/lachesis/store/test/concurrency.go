package test

import (
	"bytes"
	"runtime"
	"strconv"
)

// CurrentRoutineID returns the unique identifier of the go routine within which it gets executed
func CurrentRoutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
