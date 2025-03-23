package g_test

import (
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

// go test -bench=. -benchmem -count=4

func genSlice() Slice[String] {
	slice := NewSlice[String](0, 10000)
	for i := range 10000 {
		slice.Push(Int(i).String())
	}

	return slice
}

func BenchmarkContains(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Find(f.Eq(String("1000"))).IsSome()
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
		slice.Iter().ForEach(func(s String) { _ = s })
	}
}

func BenchmarkMap(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Map(func(s String) String { return s }).Collect()
	}
}

func BenchmarkFilter(b *testing.B) {
	slice := genSlice()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Iter().Filter(func(_ String) bool { return true }).Collect()
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
