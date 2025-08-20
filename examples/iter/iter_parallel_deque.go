package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	dq := NewDeque[int]()
	for i := range 1000 {
		dq.PushBack(i)
	}

	start := time.Now()
	dq.Iter().Parallel(1000).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed ForEach: {}", time.Since(start))

	start = time.Now()
	dq.Iter().Parallel(1000).
		Take(500).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Take: {}", time.Since(start))

	start = time.Now()
	dq.Iter().Parallel(1000).
		Skip(500).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Skip: {}", time.Since(start))

	start = time.Now()
	dq.Iter().Parallel(1000).
		Range(func(_ int) bool {
			time.Sleep(100 * time.Millisecond)
			return true
		})

	Println("Elapsed Range: {}", time.Since(start))

	start = time.Now()
	_ = dq.Iter().Parallel(1000).
		Map(func(v int) int {
			time.Sleep(100 * time.Millisecond)
			return v * 2
		}).
		Filter(f.Ne(4)).
		Collect()

	Println("Elapsed Map+Filter: {}", time.Since(start))

	start = time.Now()
	_, _ = dq.Iter().Parallel(1000).
		Partition(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i%2 == 0
		})

	Println("Elapsed Partition: {}", time.Since(start))

	start = time.Now()
	_ = dq.Iter().Parallel(1000).
		Filter(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i > 100
		}).
		Collect()

	Println("Elapsed Filter: {}", time.Since(start))

	start = time.Now()
	_ = dq.Iter().Parallel(1000).
		Find(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i == 500
		})

	Println("Elapsed Find: {}", time.Since(start))

	start = time.Now()
	_ = dq.Iter().Parallel(1000).
		Inspect(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		}).
		Fold(0, func(acc, v int) int {
			return acc + v
		})

	Println("Elapsed Fold: {}", time.Since(start))

	start = time.Now()
	_ = DequeOf(1, 2, 3).Iter().Parallel(4).
		Chain(
			DequeOf(4, 5, 6).Iter().Parallel(4),
			DequeOf(7, 8, 9).Iter().Parallel(4),
		).
		Collect()

	Println("Elapsed Chain: {}", time.Since(start))

	// Example with nested data for Flatten
	start = time.Now()
	nestedData := []any{
		[]int{1, 2, 3},
		[]int{4, 5, 6},
		[]int{7, 8, 9},
	}
	nestedDeque := NewDeque[any]()
	for _, item := range nestedData {
		nestedDeque.PushBack(item)
	}

	_ = nestedDeque.Iter().Parallel(100).
		Inspect(func(any) {
			time.Sleep(100 * time.Millisecond)
		}).
		Flatten().
		Collect()

	Println("Elapsed Flatten: {}", time.Since(start))

	// Elapsed ForEach: 101.415541ms
	// Elapsed Take: 101.433208ms
	// Elapsed Skip: 101.4415ms
	// Elapsed Range: 101.661584ms
	// Elapsed Map+Filter: 103.20925ms
	// Elapsed Partition: 102.1155ms
	// Elapsed Filter: 102.110542ms
	// Elapsed Find: 100.547291ms
	// Elapsed Fold: 102.799875ms
	// Elapsed Chain: 147.25µs
	// Elapsed Flatten: 101.524791ms
}
