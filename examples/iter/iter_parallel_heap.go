package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

func main() {
	heap := NewHeap(cmp.Cmp[int])
	for i := range 1000 {
		heap.Push(i)
	}

	start := time.Now()
	heap.Iter().Parallel(1000).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed ForEach: {}", time.Since(start))

	start = time.Now()
	heap.Iter().Parallel(1000).
		Take(500).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Take: {}", time.Since(start))

	start = time.Now()
	heap.Iter().Parallel(1000).
		Skip(500).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Skip: {}", time.Since(start))

	start = time.Now()
	heap.Iter().Parallel(1000).
		Range(func(_ int) bool {
			time.Sleep(100 * time.Millisecond)
			return true
		})

	Println("Elapsed Range: {}", time.Since(start))

	start = time.Now()
	_ = heap.Iter().Parallel(1000).
		Map(func(v int) int {
			time.Sleep(100 * time.Millisecond)
			return v * 2
		}).
		Filter(f.Ne(4)).
		Collect()

	Println("Elapsed Map+Filter: {}", time.Since(start))

	start = time.Now()
	_, _ = heap.Iter().Parallel(1000).
		Partition(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i%2 == 0
		})

	Println("Elapsed Partition: {}", time.Since(start))

	start = time.Now()
	_, _ = heap.Iter().Parallel(1000).
		PartitionWith(
			func(i int) bool {
				time.Sleep(100 * time.Millisecond)
				return i > 500
			},
			cmp.Cmp[int],     // min-heap for left partition
			cmp.Reverse[int], // max-heap for right partition
		)

	Println("Elapsed PartitionWith: {}", time.Since(start))

	start = time.Now()
	_ = heap.Iter().Parallel(1000).
		Filter(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i > 100
		}).
		CollectWith(cmp.Reverse[int]) // collect into max-heap

	Println("Elapsed Filter+CollectWith: {}", time.Since(start))

	start = time.Now()
	_ = heap.Iter().Parallel(1000).
		Find(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i == 500
		})

	Println("Elapsed Find: {}", time.Since(start))

	start = time.Now()
	_ = heap.Iter().Parallel(1000).
		Inspect(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		}).
		Fold(0, func(acc, v int) int {
			return acc + v
		})

	Println("Elapsed Fold: {}", time.Since(start))

	start = time.Now()
	chained := NewHeap(cmp.Cmp[int])
	chained.Push(1, 2, 3)
	heap2 := NewHeap(cmp.Cmp[int])
	heap2.Push(4, 5, 6)
	heap3 := NewHeap(cmp.Cmp[int])
	heap3.Push(7, 8, 9)

	_ = chained.Iter().Parallel(4).
		Chain(
			heap2.Iter().Parallel(4),
			heap3.Iter().Parallel(4),
		).
		Collect()

	Println("Elapsed Chain: {}", time.Since(start))

	// Example with nested data for Flatten
	start = time.Now()
	nestedHeap := NewHeap(func(a, b any) cmp.Ordering {
		return cmp.Cmp(a.(Slice[int])[0], b.(Slice[int])[0])
	})

	nestedHeap.Push(
		SliceOf(10, 20, 30),
		SliceOf(40, 50, 60),
		SliceOf(70, 80, 90),
	)

	_ = nestedHeap.Iter().Parallel(1000).
		Inspect(func(any) {
			time.Sleep(100 * time.Millisecond)
		}).
		Flatten().
		Collect()

	Println("Elapsed Flatten: {}", time.Since(start))

	// Test All and Any operations
	start = time.Now()
	_ = heap.Iter().Parallel(1000).
		All(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i >= 0
		})

	Println("Elapsed All: {}", time.Since(start))

	start = time.Now()
	_ = heap.Iter().Parallel(1000).
		Any(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i > 500
		})

	Println("Elapsed Any: {}", time.Since(start))

	// Elapsed ForEach: 101.306084ms
	// Elapsed Take: 101.649917ms
	// Elapsed Skip: 101.895875ms
	// Elapsed Range: 101.152ms
	// Elapsed Map+Filter: 102.952417ms
	// Elapsed Partition: 101.640834ms
	// Elapsed PartitionWith: 102.014041ms
	// Elapsed Filter+CollectWith: 102.154125ms
	// Found: 500, Elapsed Find: 100.768125ms
	// Sum: 466054, Elapsed Fold: 102.659875ms
	// Chained: Heap[1, 2, 7, 8, 9, 3, 4, 5, 6], Elapsed Chain: 117.5µs
	// Flattened and mapped: Heap[180, 40, 160, 80, 60, 120, 20, 140, 100], Elapsed Flatten: 101.757ms
	// All positive: true, Elapsed All: 101.760542ms
	// Has large number: true, Elapsed Any: 102.917ms
}
