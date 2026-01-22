package g_test

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

type concurrentCounterHeap struct {
	inFlight    int64
	maxInFlight int64
	sleep       time.Duration
}

func (cc *concurrentCounterHeap) Fn(int) {
	cur := atomic.AddInt64(&cc.inFlight, 1)
	for {
		prev := atomic.LoadInt64(&cc.maxInFlight)
		if cur <= prev || atomic.CompareAndSwapInt64(&cc.maxInFlight, prev, cur) {
			break
		}
	}

	time.Sleep(cc.sleep)
	atomic.AddInt64(&cc.inFlight, -1)
}

func (cc *concurrentCounterHeap) Max() int64 { return atomic.LoadInt64(&cc.maxInFlight) }

// TestHeapParallelCollect verifies Collect correctness and that multiple workers run concurrently.
func TestHeapParallelCollect(t *testing.T) {
	nums := []int{6, 2, 8, 1, 4, 3}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(3)
	cc := &concurrentCounterHeap{sleep: 50 * time.Millisecond}

	got := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Collect(cmp.Cmp[int])

	if got.Len() != Int(len(nums)) {
		t.Fatalf("Collect returned %d items, want %d", got.Len(), len(nums))
	}

	// Verify all elements are present (heap might reorder)
	for _, v := range nums {
		found := false
		for !got.IsEmpty() {
			if got.Pop().Some() == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Collect missing element %d", v)
		}
		// Restore heap for next iteration
		for _, n := range nums {
			got.Push(n)
		}
	}

	if cc.Max() < 2 {
		t.Errorf("expected at least 2 concurrent tasks, got max %d", cc.Max())
	}
}

// TestHeapCollectWithParallel verifies CollectWith with custom comparison.
func TestHeapCollectWithParallel(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5, 9}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(2)
	cc := &concurrentCounterHeap{sleep: 30 * time.Millisecond}

	// Collect into a max-heap
	got := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Collect(cmp.Reverse[int])

	if got.Len() != Int(len(nums)) {
		t.Fatalf("CollectWith returned %d items, want %d", got.Len(), len(nums))
	}

	// Max heap should return largest first
	largest := got.Pop()
	if !largest.IsSome() || largest.Some() != 9 {
		t.Errorf("Max heap top got %v, want Some(9)", largest)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel CollectWith, got max concurrency %d", cc.Max())
	}
}

// TestHeapMapFilterParallel verifies Map+Filter correctness and that Map runs in parallel.
func TestHeapMapFilterParallel(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range input {
		heap.Push(n)
	}

	workers := Int(2)
	cc := &concurrentCounterHeap{sleep: 30 * time.Millisecond}

	res := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Map(func(v int) int { return v * 2 }).
		Filter(func(v int) bool { return v%4 == 0 }).
		Collect(cmp.Cmp[int])

	// Expect doubled values divisible by 4: [4, 8]
	result := make([]int, 0)
	for !res.IsEmpty() {
		result = append(result, res.Pop().Some())
	}

	SliceOf(result...).SortBy(cmp.Cmp)
	expected := []int{4, 8}

	if !SliceOf(result...).Eq(expected) {
		t.Errorf("Map+Filter got %v, want %v", result, expected)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Map, got max concurrency %d", cc.Max())
	}
}

// TestHeapChainParallel verifies Chain correctness and parallel execution across both sequences.
func TestHeapChainParallel(t *testing.T) {
	heapA := NewHeap(cmp.Cmp[int])
	heapA.Push(1, 2)
	heapB := NewHeap(cmp.Cmp[int])
	heapB.Push(3, 4)

	workers := Int(2)
	cc := &concurrentCounterHeap{sleep: 20 * time.Millisecond}

	res := heapA.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Chain(
			heapB.Iter().Parallel(workers).Inspect(cc.Fn),
		).
		Collect(cmp.Cmp[int])

	result := make([]int, 0)
	for !res.IsEmpty() {
		result = append(result, res.Pop().Some())
	}

	SliceOf(result...).SortBy(cmp.Cmp)
	expected := []int{1, 2, 3, 4}

	if !SliceOf(result...).Eq(expected) {
		t.Errorf("Chain got %v, want %v", result, expected)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Chain, got max concurrency %d", cc.Max())
	}
}

// TestHeapAllAnyCountParallel verifies All, Any, Count and that All and Count run in parallel.
func TestHeapAllAnyCountParallel(t *testing.T) {
	nums := []int{2, 4, 6, 8}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(4)
	cc := &concurrentCounterHeap{sleep: 10 * time.Millisecond}

	seq := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn)

	if !seq.All(func(v int) bool { return v%2 == 0 }) {
		t.Error("All returned false for even heap")
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel All, got max concurrency %d", cc.Max())
	}

	if seq.Any(func(v int) bool { return v == 5 }) {
		t.Error("Any returned true for missing element")
	}

	cc2 := &concurrentCounterHeap{sleep: 10 * time.Millisecond}
	cnt := heap.Iter().
		Parallel(workers).
		Inspect(cc2.Fn).
		Count()

	if cnt != Int(len(nums)) {
		t.Errorf("Count returned %d, want %d", cnt, len(nums))
	}

	if cc2.Max() < 2 {
		t.Errorf("expected parallel Count, got max concurrency %d", cc2.Max())
	}
}

// TestHeapPartitionParallel verifies PartitionWith with custom comparisons.
func TestHeapPartitionParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(3)
	cc := &concurrentCounterHeap{sleep: 15 * time.Millisecond}

	// Left: min-heap for even numbers, Right: max-heap for odd numbers
	left, right := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Partition(
			func(v int) bool { return v%2 == 0 },
			cmp.Cmp[int],     // min-heap
			cmp.Reverse[int], // max-heap
		)

	// Left should be min-heap with even numbers
	if left.IsEmpty() || left.Peek().Some() != 2 {
		t.Errorf("PartitionWith left min-heap peek got %v, want Some(2)", left.Peek())
	}

	// Right should be max-heap with odd numbers
	if right.IsEmpty() || right.Peek().Some() != 5 {
		t.Errorf("PartitionWith right max-heap peek got %v, want Some(5)", right.Peek())
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel PartitionWith, got max %d", cc.Max())
	}
}

// TestHeapSkipTakeParallel verifies Skip and Take with parallelism.
func TestHeapSkipTakeParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(3)

	ccSkip := &concurrentCounterHeap{sleep: 10 * time.Millisecond}
	skipRes := heap.Iter().
		Parallel(workers).
		Inspect(ccSkip.Fn).
		Skip(3).
		Collect(cmp.Cmp[int])

	skipResult := make([]int, 0)
	for !skipRes.IsEmpty() {
		skipResult = append(skipResult, skipRes.Pop().Some())
	}

	SliceOf(skipResult...).SortBy(cmp.Cmp)

	// Skip first 3 elements from sorted heap: [1,2,3] -> remaining [4,5,6,7,8]
	if !SliceOf(skipResult...).Eq(Slice[int]{4, 5, 6, 7, 8}) {
		t.Errorf("Skip got %v, want %v", skipResult, []int{4, 5, 6, 7, 8})
	}

	ccTake := &concurrentCounterHeap{sleep: 10 * time.Millisecond}
	takeRes := heap.Iter().
		Parallel(workers).
		Inspect(ccTake.Fn).
		Take(3).
		Collect(cmp.Cmp[int])

	takeResult := make([]int, 0)
	for !takeRes.IsEmpty() {
		takeResult = append(takeResult, takeRes.Pop().Some())
	}
	SliceOf(takeResult...).SortBy(cmp.Cmp)

	// Take first 3 elements from sorted heap: [1,2,3]
	if !SliceOf(takeResult...).Eq(Slice[int]{1, 2, 3}) {
		t.Errorf("Take got %v, want %v", takeResult, []int{1, 2, 3})
	}

	if ccSkip.Max() < 2 {
		t.Errorf("expected parallel Skip, got max %d", ccSkip.Max())
	}

	if ccTake.Max() < 2 {
		t.Errorf("expected parallel Take, got max %d", ccTake.Max())
	}
}

// TestHeapReduceFoldParallel verifies Reduce and Fold operations.
func TestHeapReduceFoldParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(2)

	sum := heap.Iter().
		Parallel(workers).
		Fold(0, func(acc, v int) int { return acc + v })

	if sum != 15 {
		t.Errorf("Fold got %d, want 15", sum)
	}

	product := heap.Iter().
		Parallel(workers).
		Reduce(func(a, b int) int { return a * b })

	if !product.IsSome() || product.Some() != 120 {
		t.Errorf("Reduce got %v, want Some(120)", product)
	}
}

// TestHeapUniqueParallel verifies Unique operation with parallelism.
func TestHeapUniqueParallel(t *testing.T) {
	nums := []int{1, 2, 3, 2, 4, 3, 5, 1}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(2)
	cc := &concurrentCounterHeap{sleep: 5 * time.Millisecond}

	res := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Unique().
		Collect(cmp.Cmp[int])

	result := make([]int, 0)
	for !res.IsEmpty() {
		result = append(result, res.Pop().Some())
	}

	SliceOf(result...).SortBy(cmp.Cmp)
	expected := []int{1, 2, 3, 4, 5}

	if !SliceOf(result...).Eq(expected) {
		t.Errorf("Unique got %v, want %v", result, expected)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Unique, got max %d", cc.Max())
	}
}

// TestHeapFlattenParallel verifies Flatten correctness and parallel execution for Heap.
func TestHeapFlattenParallel(t *testing.T) {
	// Create nested test data - using heap-specific comparable data
	nestedData := []any{
		SliceOf(1, 2, 3),
		SliceOf(4, 5, 6),
		SliceOf(7, 8, 9),
		SliceOf(10, 11),
		SliceOf(12, 13, 14, 15),
	}

	cmpFn := func(a, b any) cmp.Ordering {
		// Simple comparison for any type using string representation
		aStr := fmt.Sprintf("%v", a)
		bStr := fmt.Sprintf("%v", b)
		return cmp.Cmp(aStr, bStr)
	}

	// Create heap with nested data
	heap := NewHeap[any](cmpFn)

	for _, item := range nestedData {
		heap.Push(item)
	}

	workers := Int(3)
	cc := &concurrentCounterHeap{sleep: 50 * time.Millisecond}

	start := time.Now()
	result := heap.
		Iter().
		Parallel(workers).
		Inspect(func(any) { cc.Fn(0) }). // Use generic counter
		Flatten().
		Collect(cmpFn)
	duration := time.Since(start)

	// Expected flattened elements: 3+3+3+2+4 = 15
	expectedCount := 15
	if result.Len() != Int(expectedCount) {
		t.Errorf("Expected %d flattened items, got %d", expectedCount, result.Len())
	}

	// Should achieve parallel execution of nested structures
	if cc.Max() < 2 {
		t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
	}

	// Should be faster than sequential (5 * 50ms = 250ms)
	maxExpectedTime := 200 * time.Millisecond
	if duration > maxExpectedTime {
		t.Logf("Warning: execution might not be optimal: %v", duration)
	}

	// Verify all expected elements are present (order may vary due to heap nature)
	expectedElements := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	for _, expected := range expectedElements {
		found := false
		resultSlice := result.Iter().
			Collect(func(a1, a2 any) cmp.Ordering { return cmp.Cmp(a1.(int), a2.(int)) }).
			ToSlice()
		for _, actual := range resultSlice {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing expected element: %d", expected)
		}
	}

	t.Logf("Heap Flatten - Max concurrency: %d, Duration: %v, Items: %d",
		cc.Max(), duration, result.Len())
}

// TestHeapParFlatMap verifies FlatMap correctness and parallel execution.
func TestHeapParFlatMap(t *testing.T) {
	nums := []int{1, 2, 3}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(2)
	cc := &concurrentCounterHeap{sleep: 30 * time.Millisecond}

	start := time.Now()
	result := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		FlatMap(func(v int) SeqHeap[int] {
			h := NewHeap(cmp.Cmp[int])
			h.Push(v*10, v*10+1)
			return h.Iter()
		}).
		Collect(cmp.Cmp[int])
	duration := time.Since(start)

	// Expected: [10,11, 20,21, 30,31] = 6 elements
	if result.Len() != 6 {
		t.Errorf("FlatMap result length: got %d, want 6", result.Len())
	}

	// Check that all expected values are present
	resultSlice := make([]int, 0)
	for !result.IsEmpty() {
		resultSlice = append(resultSlice, result.Pop().Some())
	}
	SliceOf(resultSlice...).SortBy(cmp.Cmp)

	expected := []int{10, 11, 20, 21, 30, 31}
	if !SliceOf(resultSlice...).Eq(expected) {
		t.Errorf("FlatMap result: got %v, want %v", resultSlice, expected)
	}

	// Verify parallelism
	if cc.Max() < 2 {
		t.Errorf("Expected parallel FlatMap, got max concurrency %d", cc.Max())
	}

	// Should be faster than sequential (3 * 30ms = 90ms)
	if duration > 80*time.Millisecond {
		t.Logf("Warning: FlatMap might not be parallel: %v", duration)
	}

	t.Logf("FlatMap - Max concurrency: %d, Duration: %v", cc.Max(), duration)

	// Test with empty result
	emptyResult := heap.Iter().
		Parallel(workers).
		FlatMap(func(int) SeqHeap[int] {
			return NewHeap(cmp.Cmp[int]).Iter()
		}).
		Collect(cmp.Cmp[int])

	if emptyResult.Len() != 0 {
		t.Errorf("FlatMap empty result: got %d, want 0", emptyResult.Len())
	}
}

// TestHeapParFilterMap verifies FilterMap correctness and parallel execution.
func TestHeapParFilterMap(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(3)
	cc := &concurrentCounterHeap{sleep: 20 * time.Millisecond}

	start := time.Now()
	result := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		FilterMap(func(v int) Option[int] {
			if v%2 == 0 {
				return Some(v * 2)
			}
			return None[int]()
		}).
		Collect(cmp.Cmp[int])
	duration := time.Since(start)

	// Expected: even numbers * 2 = [4, 8, 12]
	resultSlice := make([]int, 0)
	for !result.IsEmpty() {
		resultSlice = append(resultSlice, result.Pop().Some())
	}
	SliceOf(resultSlice...).SortBy(cmp.Cmp)

	expected := []int{4, 8, 12}
	if !SliceOf(resultSlice...).Eq(expected) {
		t.Errorf("FilterMap result: got %v, want %v", resultSlice, expected)
	}

	// Verify parallelism
	if cc.Max() < 2 {
		t.Errorf("Expected parallel FilterMap, got max concurrency %d", cc.Max())
	}

	// Should be faster than sequential (6 * 20ms = 120ms)
	if duration > 100*time.Millisecond {
		t.Logf("Warning: FilterMap might not be parallel: %v", duration)
	}

	t.Logf("FilterMap - Max concurrency: %d, Duration: %v", cc.Max(), duration)

	// Test with all filtered out
	emptyResult := heap.Iter().
		Parallel(workers).
		FilterMap(func(v int) Option[int] {
			return None[int]()
		}).
		Collect(cmp.Cmp[int])

	if emptyResult.Len() != 0 {
		t.Errorf("FilterMap all filtered: got %d, want 0", emptyResult.Len())
	}
}

// TestHeapParStepBy verifies StepBy correctness and parallel execution.
func TestHeapParStepBy(t *testing.T) {
	nums := make([]int, 10)
	for i := 0; i < 10; i++ {
		nums[i] = i + 1
	}

	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(3)
	cc := &concurrentCounterHeap{sleep: 15 * time.Millisecond}

	start := time.Now()
	result := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		StepBy(3).
		Collect(cmp.Cmp[int])
	duration := time.Since(start)

	// StepBy(3) should return approximately 1/3 of elements (3-4 elements from 10 total)
	resultSlice := make([]int, 0)
	for !result.IsEmpty() {
		resultSlice = append(resultSlice, result.Pop().Some())
	}
	SliceOf(resultSlice...).SortBy(cmp.Cmp)

	// Should have about 3-4 elements (every 3rd from 10 elements)
	expectedCount := 4 // positions 0, 3, 6, 9
	if len(resultSlice) != expectedCount {
		t.Errorf("StepBy result count: got %d, want %d", len(resultSlice), expectedCount)
	}

	// All results should be valid numbers from 1-10
	for _, v := range resultSlice {
		if v < 1 || v > 10 {
			t.Errorf("StepBy invalid result: got %d, want 1-10", v)
		}
	}

	// Verify parallelism
	if cc.Max() < 2 {
		t.Errorf("Expected parallel StepBy, got max concurrency %d", cc.Max())
	}

	// Should be faster than sequential
	if duration > 100*time.Millisecond {
		t.Logf("Warning: StepBy might not be parallel: %v", duration)
	}

	t.Logf("StepBy - Max concurrency: %d, Duration: %v", cc.Max(), duration)

	// Test StepBy(0) should default to StepBy(1)
	allResult := heap.Iter().
		Parallel(workers).
		StepBy(0).
		Collect(cmp.Cmp[int])

	if allResult.Len() != Int(len(nums)) {
		t.Errorf("StepBy(0) result length: got %d, want %d", allResult.Len(), len(nums))
	}
}

// TestHeapParMaxMinBy verifies MaxBy and MinBy correctness and parallel execution.
func TestHeapParMaxMinBy(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(3)
	cc := &concurrentCounterHeap{sleep: 10 * time.Millisecond}

	start := time.Now()
	maxResult := heap.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})
	maxDuration := time.Since(start)

	if !maxResult.IsSome() || maxResult.Some() != 9 {
		t.Errorf("MaxBy result: got %v, want Some(9)", maxResult)
	}

	// Verify parallelism
	if cc.Max() < 2 {
		t.Errorf("Expected parallel MaxBy, got max concurrency %d", cc.Max())
	}

	t.Logf("MaxBy - Max concurrency: %d, Duration: %v", cc.Max(), maxDuration)

	// Test MinBy
	cc2 := &concurrentCounterHeap{sleep: 10 * time.Millisecond}
	minResult := heap.Iter().
		Parallel(workers).
		Inspect(cc2.Fn).
		MinBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	if !minResult.IsSome() || minResult.Some() != 1 {
		t.Errorf("MinBy result: got %v, want Some(1)", minResult)
	}

	if cc2.Max() < 2 {
		t.Errorf("Expected parallel MinBy, got max concurrency %d", cc2.Max())
	}

	// Test with empty heap
	emptyHeap := NewHeap(cmp.Cmp[int])
	emptyMax := emptyHeap.Iter().
		Parallel(workers).
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	if emptyMax.IsSome() {
		t.Errorf("MaxBy empty heap: got %v, want None", emptyMax)
	}

	emptyMin := emptyHeap.Iter().
		Parallel(workers).
		MinBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	if emptyMin.IsSome() {
		t.Errorf("MinBy empty heap: got %v, want None", emptyMin)
	}
}

// TestHeapParExclude tests the Exclude method for parallel heap iterators.
func TestHeapParExclude(t *testing.T) {
	t.Run("basic_exclude", func(t *testing.T) {
		heap := NewHeap(cmp.Cmp[int])
		for i := 1; i <= 10; i++ {
			heap.Push(i)
		}

		// Exclude even numbers
		result := heap.Iter().
			Parallel(3).
			Exclude(func(n int) bool {
				return n%2 == 0
			}).
			Collect(cmp.Cmp[int])

		// Should have odd numbers: 1, 3, 5, 7, 9
		expected := []int{1, 3, 5, 7, 9}
		if result.Len() != Int(len(expected)) {
			t.Errorf("Expected %d elements, got %d", len(expected), result.Len())
		}

		resultSlice := result.ToSlice()
		resultSlice.SortBy(cmp.Cmp)
		for i, exp := range expected {
			if resultSlice[i] != exp {
				t.Errorf("Expected %d at index %d, got %d", exp, i, resultSlice[i])
			}
		}
	})

	t.Run("exclude_none", func(t *testing.T) {
		heap := NewHeap(cmp.Cmp[int])
		for i := 1; i <= 5; i++ {
			heap.Push(i)
		}

		// Exclude nothing (predicate always returns false)
		result := heap.Iter().
			Parallel(2).
			Exclude(func(n int) bool {
				return false
			}).
			Collect(cmp.Cmp[int])

		if result.Len() != Int(5) {
			t.Errorf("Expected all 5 elements, got %d", result.Len())
		}
	})

	t.Run("exclude_all", func(t *testing.T) {
		heap := NewHeap(cmp.Cmp[int])
		for i := 1; i <= 5; i++ {
			heap.Push(i)
		}

		// Exclude everything (predicate always returns true)
		result := heap.Iter().
			Parallel(2).
			Exclude(func(n int) bool {
				return true
			}).
			Collect(cmp.Cmp[int])

		if result.Len() != Int(0) {
			t.Errorf("Expected 0 elements, got %d", result.Len())
		}
	})

	t.Run("parallel_execution", func(t *testing.T) {
		heap := NewHeap(cmp.Cmp[int])
		for i := 1; i <= 20; i++ {
			heap.Push(i)
		}

		cc := &concurrentCounterHeap{sleep: 10 * time.Millisecond}

		result := heap.Iter().
			Parallel(4).
			Inspect(cc.Fn).
			Exclude(func(n int) bool {
				return n > 10 // Exclude numbers > 10
			}).
			Collect(cmp.Cmp[int])

		// Should have numbers 1-10
		if result.Len() != Int(10) {
			t.Errorf("Expected 10 elements, got %d", result.Len())
		}

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})
}

// TestHeapParForEach tests the ForEach method for parallel heap iterators.
func TestHeapParForEach(t *testing.T) {
	t.Run("basic_foreach", func(t *testing.T) {
		heap := NewHeap(cmp.Cmp[int])
		for i := 1; i <= 10; i++ {
			heap.Push(i)
		}

		var processed []int
		var mu sync.Mutex

		heap.Iter().
			Parallel(3).
			ForEach(func(n int) {
				mu.Lock()
				processed = append(processed, n)
				mu.Unlock()
			})

		if len(processed) != 10 {
			t.Errorf("Expected 10 processed elements, got %d", len(processed))
		}

		// Check all numbers are present (order may vary due to parallelism)
		sort.Ints(processed)
		for i, expected := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
			if processed[i] != expected {
				t.Errorf("Expected %d at index %d, got %d", expected, i, processed[i])
			}
		}
	})

	t.Run("parallel_execution", func(t *testing.T) {
		heap := NewHeap(cmp.Cmp[int])
		for i := 1; i <= 20; i++ {
			heap.Push(i)
		}

		cc := &concurrentCounterHeap{sleep: 10 * time.Millisecond}

		heap.Iter().
			Parallel(4).
			ForEach(func(n int) {
				cc.Fn(n)
			})

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})

	t.Run("empty_heap", func(t *testing.T) {
		heap := NewHeap(cmp.Cmp[int])
		executed := false

		heap.Iter().
			Parallel(2).
			ForEach(func(n int) {
				executed = true
			})

		if executed {
			t.Error("ForEach should not execute for empty heap")
		}
	})
}

// TestHeapParFlattenComprehensive tests Flatten method edge cases for better coverage.
func TestHeapParFlattenComprehensive(t *testing.T) {
	t.Run("nested_slices", func(t *testing.T) {
		// Create heap with any type to hold nested structures
		heap := NewHeap(func(a, b any) cmp.Ordering {
			// Simple comparison based on string representation
			aStr := fmt.Sprintf("%v", a)
			bStr := fmt.Sprintf("%v", b)
			return cmp.Cmp(aStr, bStr)
		})

		heap.Push(SliceOf(1, 2))
		heap.Push(SliceOf(3, 4, 5))

		resultSlice := heap.Iter().
			Parallel(2).
			Flatten().
			Collect(func(a1, a2 any) cmp.Ordering { return cmp.Cmp(a1.(int), a2.(int)) })

		// Should get elements from the nested slices
		if resultSlice.Len() < Int(1) {
			t.Errorf("Expected at least 1 element, got %d", resultSlice.Len())
		}
	})

	t.Run("parallel_execution", func(t *testing.T) {
		heap := NewHeap(func(a, b any) cmp.Ordering {
			aStr := fmt.Sprintf("%v", a)
			bStr := fmt.Sprintf("%v", b)
			return cmp.Cmp(aStr, bStr)
		})

		for i := 0; i < 5; i++ {
			heap.Push(SliceOf(i*2, i*2+1))
		}

		cc := &concurrentCounterHeap{sleep: 5 * time.Millisecond}

		resultSlice := heap.Iter().
			Parallel(4).
			Inspect(func(v any) { cc.Fn(1) }).
			Flatten().
			Collect(func(a1, a2 any) cmp.Ordering { return cmp.Cmp(a1.(int), a2.(int)) })

		// Should have elements from the nested slices
		if resultSlice.Len() < Int(5) {
			t.Errorf("Expected at least 5 elements, got %d", resultSlice.Len())
		}

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})

	t.Run("early_termination", func(t *testing.T) {
		heap := NewHeap(func(a, b any) cmp.Ordering {
			aStr := fmt.Sprintf("%v", a)
			bStr := fmt.Sprintf("%v", b)
			return cmp.Cmp(aStr, bStr)
		})

		heap.Push(SliceOf(1, 2, 3, 4, 5))
		heap.Push(SliceOf(6, 7, 8, 9, 10))

		resultSlice := heap.Iter().
			Parallel(2).
			Flatten().
			Take(3). // Force early termination
			Collect(func(a1, a2 any) cmp.Ordering { return cmp.Cmp(a1.(int), a2.(int)) })

		if resultSlice.Len() != Int(3) {
			t.Errorf("Expected exactly 3 elements due to Take(3), got %d", resultSlice.Len())
		}
	})
}
