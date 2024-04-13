package g_test

import (
	"testing"

	"github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func genM() g.Map[g.String, int] {
	mo := g.NewMap[g.String, int](10000)
	for i := range 10000 {
		mo.Set(g.NewInt(i).ToString(), i)
	}

	return mo
}

func BenchmarkMEq(b *testing.B) {
	m := genM()
	m2 := m.Clone()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = m.Eq(m2)
	}
}
