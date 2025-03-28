package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func BenchmarkAppendInPlace(b *testing.B) {
	b.ResetTimer()

	slice := NewSlice[String]()

	for i := range 10000000 {
		slice = append(slice, Int(i).String())
	}
}

func BenchmarkAppend(b *testing.B) {
	b.ResetTimer()

	slice := NewSlice[String]()

	for i := range 10000000 {
		slice.Push(Int(i).String())
	}
}
