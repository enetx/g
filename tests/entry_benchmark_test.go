package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func BenchmarkMapEntry(b *testing.B) {
	b.Run("occupied/or-insert", func(b *testing.B) {
		m := Map[int, int]{0: 1}
		for b.Loop() {
			_ = m.Entry(0).OrInsert(2)
		}
	})

	b.Run("vacant/or-insert", func(b *testing.B) {
		m := NewMap[int, int](1)
		for b.Loop() {
			delete(m, 0)
			_ = m.Entry(0).OrInsert(1)
		}
	})

	b.Run("and-modify", func(b *testing.B) {
		m := Map[int, int]{0: 0}
		for b.Loop() {
			m.Entry(0).AndModify(func(value *int) { *value++ })
		}
	})
}

func BenchmarkMapOrdEntry(b *testing.B) {
	mo := NewMapOrd[int, int](32)
	for i := range 32 {
		mo.Insert(i, i)
	}

	b.Run("occupied/or-insert", func(b *testing.B) {
		for b.Loop() {
			_ = mo.Entry(31).OrInsert(2)
		}
	})

	b.Run("and-modify", func(b *testing.B) {
		for b.Loop() {
			mo.Entry(31).AndModify(func(value *int) { *value++ })
		}
	})
}

func BenchmarkMapSafeEntry(b *testing.B) {
	b.Run("occupied/or-insert", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		m.Insert(0, 1)
		for b.Loop() {
			_ = m.Entry(0).OrInsert(2)
		}
	})

	b.Run("vacant/or-insert", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		for b.Loop() {
			m.Remove(0)
			_ = m.Entry(0).OrInsert(1)
		}
	})

	b.Run("and-modify", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		m.Insert(0, 0)
		for b.Loop() {
			m.Entry(0).AndModify(func(value *int) { *value++ })
		}
	})

	b.Run("and-modify/contention", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		m.Insert(0, 0)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.Entry(0).AndModify(func(value *int) { *value++ })
			}
		})
	})
}
