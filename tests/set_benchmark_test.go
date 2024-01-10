package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/iter"
)

// go test -bench=. -benchmem -count=4

func genSet() g.Set[g.String] {
	slice := g.NewSlice[g.String](0, 10000)
	for i := range iter.N(10000) {
		slice = slice.Append(g.NewInt(i).ToString())
	}

	return g.SetOf(slice...)
}

func BenchmarkSymmetricDifference(b *testing.B) {
	set1 := genSet()
	set2 := genSet()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		set1.SymmetricDifference(set2).Collect()
	}
}

func BenchmarkSymmetricDifference2(b *testing.B) {
	set1 := genSet()
	set2 := genSet()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		set1.SymmetricDifference(set2).Collect()
	}
}
