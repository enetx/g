package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	Println("=== Parallel StepBy Examples ===\n")

	// Example 1: Parallel StepBy with Slices - sampling data
	Println("1. Parallel StepBy with Slices - Data sampling:")
	start := time.Now()

	data := Range(1, 21).Collect() // 1 to 20
	result1 := data.Iter().
		Parallel(4).
		Inspect(func(int) {
			// Simulate processing time
			time.Sleep(10 * time.Millisecond)
		}).
		StepBy(3). // Take every 3rd element
		Collect()
	result1.SortBy(cmp.Cmp)

	duration1 := time.Since(start)
	Println("Original data: {}", data)
	Println("Every 3rd element: {}", result1)
	Println("Duration: {} (should be ~50ms with parallelism)\n", duration1)

	// Example 2: Parallel StepBy with Deques - log processing
	Println("2. Parallel StepBy with Deques - Log sampling:")
	deque := NewDeque[String]()
	logEntries := SliceOf[String](
		"INFO: Server started",
		"DEBUG: Connection established",
		"ERROR: Database timeout",
		"INFO: Processing request",
		"DEBUG: Query executed",
		"WARN: High memory usage",
		"INFO: Request completed",
		"DEBUG: Connection closed",
		"ERROR: Network error",
		"INFO: Server healthy",
	)
	logEntries.Iter().ForEach(func(entry String) {
		deque.PushBack(entry)
	})

	start = time.Now()
	result2 := deque.Iter().
		Parallel(3).
		Inspect(func(String) {
			// Simulate log parsing
			time.Sleep(15 * time.Millisecond)
		}).
		StepBy(2). // Sample every 2nd log entry
		Collect()

	duration2 := time.Since(start)
	Println("Sampled log entries:")
	result2.Iter().ForEach(func(entry String) {
		Println("  {}", entry)
	})
	Println("Duration: {} (should be ~75ms with parallelism)\n", duration2)

	// Example 3: Parallel StepBy with Heaps - priority sampling
	Println("3. Parallel StepBy with Heaps - Priority task sampling:")
	heap := NewHeap(cmp.Reverse[int]) // Max heap for priorities
	priorities := SliceOf(10, 25, 5, 30, 15, 8, 22, 18, 12, 28)
	priorities.Iter().ForEach(func(priority int) {
		heap.Push(priority)
	})

	start = time.Now()
	result3 := heap.Iter().
		Parallel(3).
		Inspect(func(int) {
			// Simulate task evaluation
			time.Sleep(12 * time.Millisecond)
		}).
		StepBy(3). // Sample every 3rd priority
		Collect(cmp.Cmp)

	duration3 := time.Since(start)
	Println("High-priority tasks (every 3rd): {} elements", result3.Len())
	sampled := make([]int, 0)
	for !result3.Empty() {
		sampled = append(sampled, result3.Pop().Some())
	}
	sorted := SliceOf(sampled...)
	sorted.SortBy(cmp.Reverse)
	sorted.Println()
	Println("Duration: {} (should be ~40ms with parallelism)\n", duration3)

	// Example 4: Statistical sampling - large dataset
	Println("4. Statistical sampling - Large dataset:")
	largeData := make([]int, 100)
	for i := range 100 {
		largeData[i] = i*2 + 10 // Simple linear function
	}

	start = time.Now()
	result4 := SliceOf(largeData...).
		Iter().
		Parallel(10).
		Inspect(func(int) {
			// Simulate data analysis
			time.Sleep(3 * time.Millisecond)
		}).
		StepBy(10). // Take every 10th sample
		Collect()

	duration4 := time.Since(start)
	Println("Sampled {} values from {} total", result4.Len(), len(largeData))
	firstFive := result4.Iter().Take(5).Collect()
	Println("First 5 samples: {}", firstFive)
	Println("Duration: {} (should be ~30ms with parallelism)\n", duration4)

	// Example 5: Different step sizes comparison
	Println("5. Different step sizes comparison:")
	testData := Range(1, 25).Collect() // 1 to 24

	stepSizes := SliceOf[uint](1, 2, 3, 4, 6)
	stepSizes.Iter().ForEach(func(step uint) {
		start = time.Now()
		result := testData.Iter().
			Parallel(3).
			Inspect(func(int) {
				time.Sleep(5 * time.Millisecond)
			}).
			StepBy(step).
			Collect()
		duration := time.Since(start)

		Println("Step {}: {} elements, Duration: {}", step, result.Len(), duration)
		if result.Len() <= 10 {
			result.SortBy(cmp.Cmp)
			Println("  Values: {}", result)
		} else {
			result.SortBy(cmp.Cmp)
			firstTen := result.Iter().Take(10).Collect()
			Println("  First 10: {}", firstTen)
		}
	})
	Println("")

	// Example 6: Edge cases and special behaviors
	Println("6. Edge cases and special behaviors:")

	// StepBy(0) should default to StepBy(1)
	start = time.Now()
	resultAll := SliceOf(1, 2, 3, 4, 5).Iter().
		Parallel(2).
		StepBy(0). // Should behave like StepBy(1)
		Collect()
	resultAll.SortBy(cmp.Cmp)
	Println("StepBy(0) result (should be all elements): {}", resultAll)

	// Empty collection
	emptyResult := NewSlice[int]().Iter().
		Parallel(2).
		StepBy(3).
		Collect()
	Println("Empty collection StepBy result: {}", emptyResult)

	// Single element
	singleResult := SliceOf(42).Iter().
		Parallel(2).
		StepBy(5).
		Collect()
	Println("Single element StepBy(5) result: {}", singleResult)
	Println("")

	// Example 7: Performance comparison - Sequential vs Parallel
	Println("7. Performance Comparison:")
	perfData := Range(1, 201).Collect() // 200 items

	// Sequential processing
	start = time.Now()
	seqResult := perfData.Iter().
		Inspect(func(int) {
			time.Sleep(1 * time.Millisecond)
		}).
		StepBy(5).
		Collect()
	seqDuration := time.Since(start)

	// Parallel processing
	start = time.Now()
	parResult := perfData.Iter().
		Parallel(10).
		Inspect(func(int) {
			time.Sleep(1 * time.Millisecond)
		}).
		StepBy(5).
		Collect()
	parDuration := time.Since(start)

	Println("Sequential duration: {} (should be ~200ms)", seqDuration)
	Println("Parallel duration: {} (should be ~20-30ms)", parDuration)
	Println("Speedup: {}x", Float(seqDuration.Milliseconds())/Float(parDuration.Milliseconds()))
	seqResult.SortBy(cmp.Cmp)
	parResult.SortBy(cmp.Cmp)
	Println("Results equal: {}", seqResult.Eq(parResult))
	Println("Sample count: {} from {} total", seqResult.Len(), perfData.Len())
}
