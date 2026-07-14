package g_test

import (
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func BenchmarkHeapOperations(b *testing.B) {
	b.Run("push-pop", func(b *testing.B) {
		heap := NewHeap(cmp.Cmp[int])
		for b.Loop() {
			heap.Push(1)
			_ = heap.Pop()
		}
	})

	b.Run("bulk-small-existing", func(b *testing.B) {
		base := make(Slice[int], 4096)
		for i := range base {
			base[i] = i
		}
		b.ResetTimer()
		for b.Loop() {
			heap := HeapFromSlice(cmp.Cmp[int], base)
			heap.Push(-2, -1)
		}
	})

	b.Run("iterate-clone", func(b *testing.B) {
		heap := NewHeap(cmp.Cmp[int])
		for i := range 1024 {
			heap.Push(i)
		}
		b.ResetTimer()
		for b.Loop() {
			for range heap.Iter() {
			}
		}
	})
}

func BenchmarkDequeOperations(b *testing.B) {
	b.Run("push-pop-ends", func(b *testing.B) {
		deque := NewDeque[int](16)
		for b.Loop() {
			deque.PushFront(1)
			deque.PushBack(2)
			_ = deque.PopFront()
			_ = deque.PopBack()
		}
	})

	b.Run("extend", func(b *testing.B) {
		values := make([]int, 1024)
		for b.Loop() {
			deque := NewDeque[int]()
			deque.Extend(values...)
		}
	})

	b.Run("rotate-right", func(b *testing.B) {
		deque := NewDeque[int](1024)
		deque.Extend(make([]int, 1024)...)
		b.ResetTimer()
		for b.Loop() {
			deque.RotateRight(1)
		}
	})

	b.Run("binary-search-wrapped", func(b *testing.B) {
		deque := NewDeque[int](1024)
		for i := range 1024 {
			deque.PushBack(i)
		}
		for range 256 {
			value := deque.PopFront().Some()
			deque.PushBack(value + 1024)
		}
		b.ResetTimer()
		for b.Loop() {
			_, _ = deque.BinarySearch(900, cmp.Cmp[int])
		}
	})
}
