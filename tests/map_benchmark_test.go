package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func genM() Map[String, int] {
	mo := NewMap[String, int](10000)
	for i := range 10000 {
		mo.Set(Int(i).String(), i)
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
