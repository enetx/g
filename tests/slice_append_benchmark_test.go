package g_test

import (
	"testing"

	"github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func BenchmarkAppendInPlace(b *testing.B) {
	b.ResetTimer()

	slice := g.NewSlice[g.String]()

	for i := range 10000000 {
		slice = slice.Append(g.NewInt(i).ToString())
	}
}

func BenchmarkAppend(b *testing.B) {
	b.ResetTimer()

	slice := g.NewSlice[g.String]()

	for i := range 10000000 {
		slice.AppendInPlace(g.NewInt(i).ToString())
	}
}
