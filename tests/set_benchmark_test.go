package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func genSet() Set[String] {
	slice := NewSlice[String](0, 10000)
	for i := range 10000 {
		slice.Push(Int(i).String())
	}

	return SetOf(slice...)
}

func BenchmarkSymmetricDifference(b *testing.B) {
	set1 := genSet()
	set2 := genSet()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		set1.SymmetricDifference(set2).Collect()
	}
}
