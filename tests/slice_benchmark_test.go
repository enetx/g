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

func BenchmarkForEach(b *testing.B) {
	slice := genSlice().Iter()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.ForEach(func(s g.String) {
			s.Comp().Flate()
		})
	}
}

func BenchmarkMap(b *testing.B) {
	slice := genSlice().Iter()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Map(func(s g.String) g.String { return s.Comp().Flate() }).Collect()
	}
}

func BenchmarkFilter(b *testing.B) {
	slice := genSlice().Iter()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Filter(func(s g.String) bool { return s.Comp().Flate().Len()%2 == 0 }).Collect()
	}
}

func BenchmarkUnique(b *testing.B) {
	slice := genSlice().Iter()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Unique()
	}
}

func BenchmarkDedup(b *testing.B) {
	slice := genSlice().Iter()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Dedup().Collect()
	}
}

func BenchmarkDedup2(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Compact()
	}
}
