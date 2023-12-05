package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/iter"
)

// go test -bench=. -benchmem -count=4

func BenchmarkAppendInPlace(b *testing.B) {
	b.ResetTimer()

	slice := g.NewSlice[g.String]()

	for i := range iter.N(10000000) {
		slice = slice.Append(g.NewInt(i).ToString())
	}
}

func BenchmarkAppend(b *testing.B) {
	b.ResetTimer()

	slice := g.NewSlice[g.String]()

	for i := range iter.N(10000000) {
		slice.AppendInPlace(g.NewInt(i).ToString())
	}
}
