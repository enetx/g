package g_test

import (
	"sync"
	"testing"

	. "github.com/enetx/g"
)

// go test -bench=MapSafe -benchmem -count=4

func genMapSafe() *MapSafe[int, int] {
	ms := NewMapSafe[int, int]()
	for i := range 10000 {
		ms.Set(i, i)
	}

	return ms
}

func genSyncMap() *sync.Map {
	ms := &sync.Map{}
	for i := range 10000 {
		ms.Store(i, i)
	}

	return ms
}

func BenchmarkMapSafe(b *testing.B) {
	ms := genMapSafe()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ms.Iter().Range(func(key, _ int) bool {
			ms.Delete(key)
			return true
		})
	}
}

func BenchmarkMapSafeSyncMap(b *testing.B) {
	ms := genSyncMap()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ms.Range(func(key, _ any) bool {
			ms.Delete(key)
			return true
		})
	}
}
