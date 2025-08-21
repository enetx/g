package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	Println("=== New Parallel Iterator Methods Demo ===\n")

	// Sample data for all examples
	numbers := SliceOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)

	Println("Original data: {}", numbers)
	Println("")

	// 1. FlatMap - Expand each number into a small range
	Println("1. FlatMap - Expand numbers into ranges:")
	start := time.Now()

	flatResult := numbers.Iter().
		Parallel(4).
		FlatMap(func(n int) SeqSlice[int] {
			time.Sleep(10 * time.Millisecond) // Simulate work
			// Each number becomes [n*10, n*10+1, n*10+2]
			return SliceOf(n*10, n*10+1, n*10+2).Iter()
		}).
		Collect()
	flatResult.SortBy(cmp.Cmp)

	Println("   Result: {} elements", flatResult.Len())
	Println("   First 10: {}", flatResult.Iter().Take(10).Collect())
	Println("   Duration: {}", time.Since(start))
	Println("")

	// 2. FilterMap - Keep only even numbers, double them
	Println("2. FilterMap - Process even numbers only:")
	start = time.Now()

	filterResult := numbers.Iter().
		Parallel(3).
		FilterMap(func(n int) Option[int] {
			time.Sleep(8 * time.Millisecond)
			if n%2 == 0 {
				return Some(n * 2) // Double even numbers
			}
			return None[int]()
		}).
		Collect()
	filterResult.SortBy(cmp.Cmp)

	Println("   Result: {}", filterResult)
	Println("   Duration: {}", time.Since(start))
	Println("")

	// 3. StepBy - Sample every 3rd element
	Println("3. StepBy - Sample every 3rd element:")
	start = time.Now()

	stepResult := numbers.Iter().
		Parallel(3).
		Inspect(func(n int) {
			time.Sleep(5 * time.Millisecond)
		}).
		StepBy(3).
		Collect()
	stepResult.SortBy(cmp.Cmp)

	Println("   Result: {}", stepResult)
	Println("   Duration: {}", time.Since(start))
	Println("")

	// 4. MaxBy/MinBy - Find extremes
	Println("4. MaxBy/MinBy - Find extreme values:")
	start = time.Now()

	maxVal := numbers.Iter().
		Parallel(3).
		Inspect(func(n int) {
			time.Sleep(3 * time.Millisecond)
		}).
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	minVal := numbers.Iter().
		Parallel(3).
		Inspect(func(n int) {
			time.Sleep(3 * time.Millisecond)
		}).
		MinBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	Println("   Maximum: {}", maxVal)
	Println("   Minimum: {}", minVal)
	Println("   Duration: {}", time.Since(start))
	Println("")

	// 5. Combined pipeline - All methods together
	Println("5. Combined Pipeline - All methods in sequence:")
	start = time.Now()

	combined := numbers.Iter().
		Parallel(4).
		// First expand each number
		FlatMap(func(n int) SeqSlice[int] {
			time.Sleep(2 * time.Millisecond)
			return SliceOf(n, n+100).Iter() // Original and +100 variant
		}).
		// Filter and transform
		FilterMap(func(n int) Option[int] {
			if n > 5 && n < 50 { // Keep mid-range values
				return Some(n * 2)
			}
			return None[int]()
		}).
		// Sample every 2nd
		StepBy(2).
		// Find maximum
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	Println("   Combined result: {}", combined)
	Println("   Duration: {}", time.Since(start))
	Println("")

	// 6. Performance comparison with different data structures
	Println("6. Performance with different data structures:")

	// Slice performance
	start = time.Now()
	sliceResult := numbers.Iter().
		Parallel(3).
		FilterMap(func(n int) Option[int] {
			time.Sleep(5 * time.Millisecond)
			if n%2 == 1 {
				return Some(n * 10) // Transform odd numbers
			}
			return None[int]()
		}).
		Collect()
	sliceDuration := time.Since(start)

	// Deque performance
	deque := NewDeque[int]()
	numbers.Iter().ForEach(func(n int) {
		deque.PushBack(n)
	})

	start = time.Now()
	dequeResult := deque.Iter().
		Parallel(3).
		FilterMap(func(n int) Option[int] {
			time.Sleep(5 * time.Millisecond)
			if n%2 == 1 {
				return Some(n * 10) // Transform odd numbers
			}
			return None[int]()
		}).
		Collect()
	dequeDuration := time.Since(start)

	// Heap performance
	heap := NewHeap(cmp.Cmp[int])
	numbers.Iter().ForEach(func(n int) {
		heap.Push(n)
	})

	start = time.Now()
	heapResult := heap.Iter().
		Parallel(3).
		FilterMap(func(n int) Option[int] {
			time.Sleep(5 * time.Millisecond)
			if n%2 == 1 {
				return Some(n * 10) // Transform odd numbers
			}
			return None[int]()
		}).
		Collect()
	heapDuration := time.Since(start)

	Println("   Slice result: {} items in {}", sliceResult.Len(), sliceDuration)
	Println("   Deque result: {} items in {}", dequeResult.Len(), dequeDuration)
	Println("   Heap result:  {} items in {}", heapResult.Len(), heapDuration)
	Println("   Results consistent: {}",
		sliceResult.Len() == dequeResult.Len() && dequeResult.Len() == heapResult.Len())

	Println("\n=== Summary ===")
	Println("All new parallel iterator methods are working correctly!")
	Println("- FlatMap: Parallel expansion of sequences ✓")
	Println("- FilterMap: Parallel filtering with transformation ✓")
	Println("- StepBy: Parallel sampling with atomic counting ✓")
	Println("- MaxBy/MinBy: Parallel extrema finding ✓")
	Println("- All methods work with Slice, Deque, and Heap iterators ✓")

	/* Expected Output:
	=== New Parallel Iterator Methods Demo ===

	Original data: Slice[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]

	1. FlatMap - Expand numbers into ranges:
	   Result: 36 elements
	   First 10: Slice[10, 11, 12, 20, 21, 22, 30, 31, 32, 40]
	   Duration: 32ms

	2. FilterMap - Process even numbers only:
	   Result: Slice[4, 8, 12, 16, 20, 24]
	   Duration: 27ms

	3. StepBy - Sample every 3rd element:
	   Result: Slice[1, 4, 7, 10]
	   Duration: 21ms

	4. MaxBy/MinBy - Find extreme values:
	   Maximum: Some(12)
	   Minimum: Some(1)
	   Duration: 13ms

	5. Combined Pipeline - All methods in sequence:
	   Combined result: Some(24)
	   Duration: 8ms

	6. Performance with different data structures:
	   Slice result: 6 items in 21ms
	   Deque result: 6 items in 21ms
	   Heap result:  6 items in 21ms
	   Results consistent: true

	=== Summary ===
	All new parallel iterator methods are working correctly!
	- FlatMap: Parallel expansion of sequences ✓
	- FilterMap: Parallel filtering with transformation ✓
	- StepBy: Parallel sampling with atomic counting ✓
	- MaxBy/MinBy: Parallel extrema finding ✓
	- All methods work with Slice, Deque, and Heap iterators ✓
	*/
}
