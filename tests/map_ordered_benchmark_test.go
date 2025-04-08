package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func genMO() MapOrd[String, int] {
	mo := NewMapOrd[String, int](10000)
	for i := range 10000 {
		mo.Set(Int(i).String(), i)
	}

	return mo
}

func BenchmarkMoContains(b *testing.B) {
	mo := genMO()

	for b.Loop() {
		_ = mo.Contains("9999")
	}
}

func BenchmarkMoEq(b *testing.B) {
	mo := genMO()
	mo2 := mo.Clone()

	b.ResetTimer()

	for b.Loop() {
		_ = mo.Eq(mo2)
	}
}

func BenchmarkMoGet(b *testing.B) {
	mo := genMO()

	for b.Loop() {
		_ = mo.Get("9999")
	}
}
