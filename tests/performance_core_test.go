package g_test

import (
	"sync"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestPerformanceCoreParity(t *testing.T) {
	t.Run("bytes join", func(t *testing.T) {
		tests := []struct {
			values Slice[Bytes]
			sep    Bytes
			want   String
		}{
			{nil, nil, ""},
			{nil, Bytes("::"), ""},
			{SliceOf(Bytes("a")), Bytes(","), "a"},
			{SliceOf(Bytes(""), Bytes("b"), Bytes("")), Bytes("::"), "::b::"},
			{SliceOf(Bytes("привет"), Bytes("мир")), Bytes("/"), "привет/мир"},
		}
		for _, tt := range tests {
			if got := tt.values.Join(tt.sep); got != tt.want {
				t.Errorf("Join(%q) = %q, want %q", tt.sep, got, tt.want)
			}
		}
	})

	t.Run("justification", func(t *testing.T) {
		if got := String("go").LeftJustify(7, "аб"); got != "goабаба" {
			t.Fatalf("LeftJustify = %q", got)
		}
		if got := String("go").RightJustify(7, "аб"); got != "абабаgo" {
			t.Fatalf("RightJustify = %q", got)
		}
		if got := String("go").Center(7, "аб"); got != "абgoаба" {
			t.Fatalf("Center = %q", got)
		}
	})

	t.Run("binary", func(t *testing.T) {
		for value, want := range map[Int]String{0: "00000000", 5: "00000101", -5: "-0000101", 255: "11111111"} {
			if got := value.Binary(); got != want {
				t.Errorf("%d.Binary() = %q, want %q", value, got, want)
			}
		}
	})

	t.Run("heap collect", func(t *testing.T) {
		seq := SeqHeap[int](func(yield func(int) bool) {
			for _, value := range []int{3, 1, 2, 1} {
				if !yield(value) {
					return
				}
			}
		})
		heap := seq.Collect(cmp.Cmp[int])
		for _, want := range []int{1, 1, 2, 3} {
			if got := heap.Pop(); got.IsNone() || got.Some() != want {
				t.Fatalf("expected %d, got %v", want, got)
			}
		}

		empty := SeqHeap[int](func(func(int) bool) {}).Collect(cmp.Cmp[int])
		if !empty.IsEmpty() {
			t.Fatal("empty sequence must collect into an empty heap")
		}

		partial := SeqHeap[int](func(yield func(int) bool) {
			yield(4)
			yield(2)
		}).Collect(cmp.Cmp[int])
		if partial.Pop().Some() != 2 || partial.Pop().Some() != 4 {
			t.Fatal("partial sequence values were not heapified")
		}

		maxHeap := seq.Collect(func(a, b int) cmp.Ordering { return cmp.Cmp(b, a) })
		for _, want := range []int{3, 2, 1, 1} {
			if got := maxHeap.Pop(); got.IsNone() || got.Some() != want {
				t.Fatalf("max heap: expected %d, got %v", want, got)
			}
		}
	})
}

func TestMapSafeClearStructuralConsistency(t *testing.T) {
	for iteration := range 100 {
		m := NewMapSafe[int, int]()
		for key := range 64 {
			m.Insert(key, key)
		}

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			<-start
			m.Clear()
		}()
		go func() {
			defer wg.Done()
			<-start
			for key := range 64 {
				m.Insert(key, iteration)
			}
		}()
		go func() {
			defer wg.Done()
			<-start
			for key := range 64 {
				m.Entry(key).OrInsert(iteration)
				m.Remove(key)
			}
		}()
		close(start)
		wg.Wait()

		count := m.Iter().Count()
		if m.Len() != count {
			t.Fatalf("iteration %d: Len=%d Iter.Count=%d", iteration, m.Len(), count)
		}
	}
}

func BenchmarkPerformanceCore(b *testing.B) {
	b.Run("format-to-reset", func(b *testing.B) {
		var builder Builder
		for b.Loop() {
			builder.Reset()
			FormatTo(&builder, "{}={:04d}", "answer", 42)
			_ = builder.String()
		}
	})
	b.Run("bytes-join", func(b *testing.B) {
		values := SliceOf(Bytes("alpha"), Bytes("beta"), Bytes("gamma"), Bytes("delta"))
		for b.Loop() {
			_ = values.Join(Bytes("::"))
		}
	})
	b.Run("left-justify", func(b *testing.B) {
		for b.Loop() {
			_ = String("value").LeftJustify(128, "аб")
		}
	})
	b.Run("binary", func(b *testing.B) {
		for b.Loop() {
			_ = Int(123456).Binary()
		}
	})
	b.Run("heap-collect", func(b *testing.B) {
		seq := SeqHeap[int](func(yield func(int) bool) {
			for value := 1024; value > 0; value-- {
				if !yield(value) {
					return
				}
			}
		})
		for b.Loop() {
			_ = seq.Collect(cmp.Cmp[int])
		}
	})
}

func BenchmarkMapSafeStructural(b *testing.B) {
	b.Run("read", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		m.Insert(1, 1)
		for b.Loop() {
			_ = m.Get(1)
		}
	})
	b.Run("insert-replace", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		m.Insert(1, 1)
		for b.Loop() {
			m.Insert(1, 2)
		}
	})
	b.Run("insert-remove", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		for b.Loop() {
			m.Insert(1, 1)
			m.Remove(1)
		}
	})
	b.Run("clear", func(b *testing.B) {
		m := NewMapSafe[int, int]()
		for b.Loop() {
			m.Insert(1, 1)
			m.Clear()
		}
	})
}
