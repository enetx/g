package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	Println("=== Parallel FlatMap Examples ===\n")

	// Example 1: Parallel FlatMap with Slices - expanding numbers
	Println("1. Parallel FlatMap with Slices:")
	start := time.Now()

	result1 := SliceOf(1, 2, 3, 4, 5).
		Iter().
		Parallel(3).
		FlatMap(func(n int) SeqSlice[int] {
			// Simulate some work
			time.Sleep(50 * time.Millisecond)
			// Expand each number to a range
			return Range(n*10, n*10+3).Collect().Iter()
		}).
		Collect()
	result1.SortBy(cmp.Cmp)

	duration1 := time.Since(start)
	Println("Result: {}", result1)
	Println("Duration: {} (should be ~50ms with parallelism)\n", duration1)

	// Example 2: Parallel FlatMap with Deques - word processing
	Println("2. Parallel FlatMap with Deques:")
	deque := NewDeque[String]()
	deque.PushBack("hello world")
	deque.PushBack("go programming")
	deque.PushBack("parallel computing")

	start = time.Now()
	result2 := deque.Iter().
		Parallel(2).
		FlatMap(func(sentence String) SeqDeque[String] {
			// Simulate processing time
			time.Sleep(30 * time.Millisecond)
			// Split sentence into words and create deque
			words := NewDeque[String]()
			sentence.Fields().ForEach(func(word String) {
				words.PushBack(word.Upper())
			})
			return words.Iter()
		}).
		Collect()

	duration2 := time.Since(start)
	Println("Result: {}", result2.String())
	Println("Duration: {} (should be ~60ms with parallelism)\n", duration2)

	// Example 3: Parallel FlatMap with Heaps - matrix expansion
	Println("3. Parallel FlatMap with Heaps:")
	heap := NewHeap(cmp.Cmp[int])
	heap.Push(2, 3, 4)

	start = time.Now()
	result3 := heap.Iter().
		Parallel(3).
		FlatMap(func(size int) SeqHeap[int] {
			// Simulate computation
			time.Sleep(40 * time.Millisecond)
			// Create a heap with multiplication table
			h := NewHeap(cmp.Cmp[int])
			for i := range Int(size) {
				h.Push(size * int(i+1))
			}
			return h.Iter()
		}).
		Collect()

	duration3 := time.Since(start)
	Println("Result: {} elements", result3.Len())
	allElements := make([]int, 0)
	for !result3.Empty() {
		allElements = append(allElements, result3.Pop().Some())
	}
	sortedElements := SliceOf(allElements...)
	sortedElements.SortBy(cmp.Cmp)
	sortedElements.Println()
	Println("Duration: {} (should be ~40ms with parallelism)\n", duration3)

	// Example 4: Complex nested data processing
	Println("4. Complex nested data processing:")
	data := SliceOf[String]("user1,user2", "admin1,admin2,admin3", "guest1")

	start = time.Now()
	result4 := data.Iter().
		Parallel(3).
		FlatMap(func(userGroup String) SeqSlice[String] {
			// Simulate user processing
			time.Sleep(25 * time.Millisecond)
			// Split the group and process each user
			return userGroup.Split(",").Map(func(user String) String {
				return user.Append("_processed")
			})
		}).
		Collect()

	duration4 := time.Since(start)
	Println("Result: {}", result4)
	Println("Duration: {} (should be ~25ms with parallelism)\n", duration4)

	// Example 5: Performance comparison - Sequential vs Parallel
	Println("5. Performance Comparison:")
	largeData := Range(1, 21).Collect() // 20 items

	// Sequential processing
	start = time.Now()
	seqResult := largeData.Iter().
		FlatMap(func(n int) SeqSlice[int] {
			time.Sleep(10 * time.Millisecond)
			return SliceOf(n, n*2).Iter()
		}).
		Collect()
	seqDuration := time.Since(start)

	// Parallel processing
	start = time.Now()
	parResult := largeData.Iter().
		Parallel(10).
		FlatMap(func(n int) SeqSlice[int] {
			time.Sleep(10 * time.Millisecond)
			return SliceOf(n, n*2).Iter()
		}).
		Collect()
	parDuration := time.Since(start)

	Println("Sequential duration: {} (should be ~200ms)", seqDuration)
	Println("Parallel duration: {} (should be ~20-30ms)", parDuration)
	Println("Speedup: {}x", Float(seqDuration.Milliseconds())/Float(parDuration.Milliseconds()))
	seqResult.SortBy(cmp.Cmp)
	parResult.SortBy(cmp.Cmp)
	Println("Results equal: {}", seqResult.Eq(parResult))

	/* Expected Output:
	=== Parallel FlatMap Examples ===

	1. Parallel FlatMap with Slices:
	Result: Slice[10, 11, 12, 20, 21, 22, 30, 31, 32, 40, 41, 42, 50, 51, 52]
	Duration: 52ms (should be ~50ms with parallelism)

	2. Parallel FlatMap with Deques:
	Result: Deque[HELLO, WORLD, GO, PROGRAMMING, PARALLEL, COMPUTING]
	Duration: 62ms (should be ~60ms with parallelism)

	3. Parallel FlatMap with Heaps:
	Result: 9 elements
	Slice[2, 3, 4, 6, 8, 9, 12, 16]
	Duration: 41ms (should be ~40ms with parallelism)

	4. Complex nested data processing:
	Result: Slice[user1_processed, user2_processed, admin1_processed, admin2_processed, admin3_processed, guest1_processed]
	Duration: 26ms (should be ~25ms with parallelism)

	5. Performance Comparison:
	Sequential duration: 201ms (should be ~200ms)
	Parallel duration: 23ms (should be ~20-30ms)
	Speedup: 8.7x
	Results equal: true
	*/
}
