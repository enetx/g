package g_test

import (
	"fmt"
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
		Collect()

	if got.Len() != Int(len(nums)) {
		t.Fatalf("Collect returned %d items, want %d", got.Len(), len(nums))
	}

	// Verify all elements are present (heap might reorder)
	for _, v := range nums {
		found := false
		for !got.Empty() {
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
		CollectWith(cmp.Reverse[int])

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
		Collect()

	// Expect doubled values divisible by 4: [4, 8]
	result := make([]int, 0)
	for !res.Empty() {
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
		Collect()

	result := make([]int, 0)
	for !res.Empty() {
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

// TestHeapFindPartitionParallel verifies Find and Partition correctness and parallelism.
func TestHeapFindPartitionParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	heap := NewHeap(cmp.Cmp[int])
	for _, n := range nums {
		heap.Push(n)
	}

	workers := Int(3)

	ccFind := &concurrentCounterHeap{sleep: 15 * time.Millisecond}
	opt := heap.Iter().
		Parallel(workers).
		Inspect(ccFind.Fn).
		Find(func(v int) bool { return v > 3 })

	if !opt.IsSome() || opt.Some() < 4 {
		t.Errorf("Find got %v, want Some(4) or Some(5)", opt)
	}

	if ccFind.Max() < 2 {
		t.Errorf("expected parallel Find, got max %d", ccFind.Max())
	}

	ccPart := &concurrentCounterHeap{sleep: 15 * time.Millisecond}
	left, right := heap.Iter().
		Parallel(workers).
		Inspect(ccPart.Fn).
		Partition(func(v int) bool { return v%2 == 0 })

	leftResult := make([]int, 0)
	for !left.Empty() {
		leftResult = append(leftResult, left.Pop().Some())
	}

	rightResult := make([]int, 0)
	for !right.Empty() {
		rightResult = append(rightResult, right.Pop().Some())
	}

	SliceOf(leftResult...).SortBy(cmp.Cmp)
	SliceOf(rightResult...).SortBy(cmp.Cmp)

	if !SliceOf(leftResult...).Eq(Slice[int]{2, 4}) {
		t.Errorf("Partition left got %v, want %v", leftResult, []int{2, 4})
	}

	if !SliceOf(rightResult...).Eq(Slice[int]{1, 3, 5}) {
		t.Errorf("Partition right got %v, want %v", rightResult, []int{1, 3, 5})
	}

	if ccPart.Max() < 2 {
		t.Errorf("expected parallel Partition, got max %d", ccPart.Max())
	}
}

// TestHeapPartitionWithParallel verifies PartitionWith with custom comparisons.
func TestHeapPartitionWithParallel(t *testing.T) {
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
		PartitionWith(
			func(v int) bool { return v%2 == 0 },
			cmp.Cmp[int],     // min-heap
			cmp.Reverse[int], // max-heap
		)

	// Left should be min-heap with even numbers
	if left.Empty() || left.Peek().Some() != 2 {
		t.Errorf("PartitionWith left min-heap peek got %v, want Some(2)", left.Peek())
	}

	// Right should be max-heap with odd numbers
	if right.Empty() || right.Peek().Some() != 5 {
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
		Collect()

	skipResult := make([]int, 0)
	for !skipRes.Empty() {
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
		Collect()

	takeResult := make([]int, 0)
	for !takeRes.Empty() {
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
		Collect()

	result := make([]int, 0)
	for !res.Empty() {
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

	// Create heap with nested data
	heap := NewHeap[any](func(a, b any) cmp.Ordering {
		// Simple comparison for any type using string representation
		aStr := fmt.Sprintf("%v", a)
		bStr := fmt.Sprintf("%v", b)
		return cmp.Cmp(aStr, bStr)
	})

	for _, item := range nestedData {
		heap.Push(item)
	}

	workers := Int(3)
	cc := &concurrentCounterHeap{sleep: 50 * time.Millisecond}

	start := time.Now()
	result := heap.
		Iter().
		Parallel(workers).
		Inspect(func(v any) { cc.Fn(0) }). // Use generic counter
		Flatten().
		Collect()
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
		resultSlice := result.Iter().Collect().ToSlice()
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
