package g_test

import (
	"testing"

	"github.com/enetx/g"
	"github.com/enetx/g/f"
)

// go test -bench=. -benchmem -count=4

func genSlice() g.Slice[g.String] {
	slice := g.NewSlice[g.String](0, 10000)
	for i := range 10000 {
		slice = slice.Append(g.NewInt(i).ToString())
	}

	return slice
}

func BenchmarkContains(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Find(f.Eq(g.String("1000"))).IsSome()
	}
}

func BenchmarkContains2(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Contains("10000")
	}
}

func BenchmarkForEach(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().ForEach(func(s g.String) { _ = s })
	}
}

func BenchmarkMap(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Map(func(s g.String) g.String { return s }).Collect()
	}
}

func BenchmarkFilter(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Filter(func(_ g.String) bool { return true }).Collect()
	}
}

func BenchmarkUnique(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Unique().Collect()
	}
}

func BenchmarkDedup(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Dedup().Collect()
	}
}
