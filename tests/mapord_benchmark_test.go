package g_test

import (
	"testing"

	"github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func genMO() g.MapOrd[g.String, int] {
	mo := g.NewMapOrd[g.String, int](10000)
	for i := range 10000 {
		mo.Set(g.NewInt(i).String(), i)
	}

	return mo
}

func BenchmarkMoContains(b *testing.B) {
	mo := genMO()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = mo.Contains("9999")
	}
}

func BenchmarkMoEq(b *testing.B) {
	mo := genMO()
	mo2 := mo.Clone()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = mo.Eq(mo2)
	}
}

func BenchmarkMoGet(b *testing.B) {
	mo := genMO()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = mo.Get("9999")
	}
}
