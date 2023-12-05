package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/iter"
)

// go test -bench=. -benchmem -count=4

func genSlice() g.Slice[g.String] {
	slice := g.NewSlice[g.String](0, 10000)
	for i := range iter.N(10000) {
		slice = slice.Append(g.NewInt(i).ToString())
	}

	return slice
}

func BenchmarkMap(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Map(func(s g.String) g.String { return s.Comp().Flate() })
	}
}

func BenchmarkMapParallel(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.MapParallel(func(s g.String) g.String { return s.Comp().Flate() })
	}
}

func BenchmarkFilter(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Filter(func(s g.String) bool { return s.Comp().Flate().Len()%2 == 0 })
	}
}

func BenchmarkFilterParallel(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.FilterParallel(func(s g.String) bool { return s.Comp().Flate().Len()%2 == 0 })
	}
}
