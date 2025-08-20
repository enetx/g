package g_test

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

type concurrentCounterDeque struct {
	inFlight    int64
	maxInFlight int64
	sleep       time.Duration
}

func (cc *concurrentCounterDeque) Fn(int) {
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

func (cc *concurrentCounterDeque) Max() int64 { return atomic.LoadInt64(&cc.maxInFlight) }

// TestDequeParallelCollect verifies Collect correctness and that multiple workers run concurrently.
func TestDequeParallelCollect(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	dq := NewDeque[int]()
	for _, n := range nums {
		dq.PushBack(n)
	}

	workers := Int(3)
	cc := &concurrentCounterDeque{sleep: 50 * time.Millisecond}

	got := dq.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Collect()

	if got.Len() != Int(len(nums)) {
		t.Fatalf("Collect returned %d items, want %d", got.Len(), len(nums))
	}

	for _, v := range nums {
		found := false
		got.Iter().Range(func(item int) bool {
			if item == v {
				found = true
				return false
			}
			return true
		})
		if !found {
			t.Errorf("Collect missing element %d in result", v)
		}
	}

	if cc.Max() < 2 {
		t.Errorf("expected at least 2 concurrent tasks, got max %d", cc.Max())
	}
}

// TestDequeMapFilterParallel verifies Map+Filter correctness and that Map runs in parallel.
func TestDequeMapFilterParallel(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	dq := NewDeque[int]()
	for _, n := range input {
		dq.PushBack(n)
	}

	workers := Int(2)
	cc := &concurrentCounterDeque{sleep: 30 * time.Millisecond}

	res := dq.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Map(func(v int) int { return v * 2 }).
		Filter(func(v int) bool { return v%4 == 0 }).
		Collect()

	expected := []int{4, 8}
	result := res.Iter().Collect()

	// Check elements are present
	for _, v := range expected {
		found := false
		result.Iter().Range(func(item int) bool {
			if item == v {
				found = true
				return false
			}
			return true
		})
		if !found {
			t.Errorf("Map+Filter missing element %d in result", v)
		}
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Map, got max concurrency %d", cc.Max())
	}
}

// TestDequeChainParallel verifies Chain correctness and parallel execution across both sequences.
func TestDequeChainParallel(t *testing.T) {
	dqA := DequeOf(1, 2)
	dqB := DequeOf(3, 4)
	workers := Int(2)
	cc := &concurrentCounterDeque{sleep: 20 * time.Millisecond}

	res := dqA.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Chain(
			dqB.Iter().Parallel(workers).Inspect(cc.Fn),
		).
		Collect()

	expected := []int{1, 2, 3, 4}

	// Check all elements are present
	for _, v := range expected {
		found := false
		res.Iter().Range(func(item int) bool {
			if item == v {
				found = true
				return false
			}
			return true
		})
		if !found {
			t.Errorf("Chain missing element %d in result", v)
		}
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Chain, got max concurrency %d", cc.Max())
	}
}

// TestDequeAllAnyCountParallel verifies All, Any, Count and that All and Count run in parallel.
func TestDequeAllAnyCountParallel(t *testing.T) {
	nums := []int{2, 4, 6, 8}
	dq := NewDeque[int]()
	for _, n := range nums {
		dq.PushBack(n)
	}

	workers := Int(4)
	cc := &concurrentCounterDeque{sleep: 10 * time.Millisecond}

	seq := dq.Iter().
		Parallel(workers).
		Inspect(cc.Fn)

	if !seq.All(func(v int) bool { return v%2 == 0 }) {
		t.Error("All returned false for even deque")
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel All, got max concurrency %d", cc.Max())
	}

	if seq.Any(func(v int) bool { return v == 5 }) {
		t.Error("Any returned true for missing element")
	}

	cc2 := &concurrentCounterDeque{sleep: 10 * time.Millisecond}
	cnt := dq.Iter().
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

// TestDequeFindPartitionParallel verifies Find and Partition correctness and parallelism.
func TestDequeFindPartitionParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	dq := NewDeque[int]()
	for _, n := range nums {
		dq.PushBack(n)
	}

	workers := Int(3)

	ccFind := &concurrentCounterDeque{sleep: 15 * time.Millisecond}
	opt := dq.Iter().
		Parallel(workers).
		Inspect(ccFind.Fn).
		Find(func(v int) bool { return v > 3 })

	if !opt.IsSome() || opt.Some() < 4 {
		t.Errorf("Find got %v, want Some(4) or Some(5)", opt)
	}

	if ccFind.Max() < 2 {
		t.Errorf("expected parallel Find, got max %d", ccFind.Max())
	}

	ccPart := &concurrentCounterDeque{sleep: 15 * time.Millisecond}
	left, right := dq.Iter().
		Parallel(workers).
		Inspect(ccPart.Fn).
		Partition(func(v int) bool { return v%2 == 0 })

	expectedLeft := []int{2, 4}
	expectedRight := []int{1, 3, 5}

	// Check left partition
	for _, v := range expectedLeft {
		found := false
		left.Iter().Range(func(item int) bool {
			if item == v {
				found = true
				return false
			}
			return true
		})
		if !found {
			t.Errorf("Partition left missing element %d", v)
		}
	}

	// Check right partition
	for _, v := range expectedRight {
		found := false
		right.Iter().Range(func(item int) bool {
			if item == v {
				found = true
				return false
			}
			return true
		})
		if !found {
			t.Errorf("Partition right missing element %d", v)
		}
	}

	if ccPart.Max() < 2 {
		t.Errorf("expected parallel Partition, got max %d", ccPart.Max())
	}
}

// TestDequeSkipTakeParallel verifies Skip and Take with parallelism.
func TestDequeSkipTakeParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	dq := NewDeque[int]()
	for _, n := range nums {
		dq.PushBack(n)
	}

	workers := Int(3)

	ccSkip := &concurrentCounterDeque{sleep: 10 * time.Millisecond}
	skipRes := dq.Iter().
		Parallel(workers).
		Inspect(ccSkip.Fn).
		Skip(3).
		Collect()

	skipSlice := skipRes.Iter().Collect().ToSlice()
	skipSlice.SortBy(cmp.Cmp)

	if !skipSlice.Eq(Slice[int]{4, 5, 6, 7, 8}) {
		t.Errorf("Skip got %v, want %v", skipSlice, []int{4, 5, 6, 7, 8})
	}

	ccTake := &concurrentCounterDeque{sleep: 10 * time.Millisecond}
	takeRes := dq.Iter().
		Parallel(workers).
		Inspect(ccTake.Fn).
		Take(3).
		Collect()

	takeSlice := takeRes.Iter().Collect().ToSlice()
	takeSlice.SortBy(cmp.Cmp)

	if !takeSlice.Eq(Slice[int]{1, 2, 3}) {
		t.Errorf("Take got %v, want %v", takeSlice, []int{1, 2, 3})
	}

	if ccSkip.Max() < 2 {
		t.Errorf("expected parallel Skip, got max %d", ccSkip.Max())
	}

	if ccTake.Max() < 2 {
		t.Errorf("expected parallel Take, got max %d", ccTake.Max())
	}
}

// TestDequeReduceFoldParallel verifies Reduce and Fold operations.
func TestDequeReduceFoldParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	dq := NewDeque[int]()
	for _, n := range nums {
		dq.PushBack(n)
	}

	workers := Int(2)

	sum := dq.Iter().
		Parallel(workers).
		Fold(0, func(acc, v int) int { return acc + v })

	if sum != 15 {
		t.Errorf("Fold got %d, want 15", sum)
	}

	product := dq.Iter().
		Parallel(workers).
		Reduce(func(a, b int) int { return a * b })

	if !product.IsSome() || product.Some() != 120 {
		t.Errorf("Reduce got %v, want Some(120)", product)
	}
}

// TestDequeUniqueParallel verifies Unique operation with parallelism.
func TestDequeUniqueParallel(t *testing.T) {
	nums := []int{1, 2, 3, 2, 4, 3, 5, 1}
	dq := NewDeque[int]()
	for _, n := range nums {
		dq.PushBack(n)
	}

	workers := Int(2)
	cc := &concurrentCounterDeque{sleep: 5 * time.Millisecond}

	res := dq.Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Unique().
		Collect()

	result := res.Iter().Collect().ToSlice()
	result.SortBy(cmp.Cmp)
	expected := []int{1, 2, 3, 4, 5}

	if !result.Eq(expected) {
		t.Errorf("Unique got %v, want %v", result, expected)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Unique, got max %d", cc.Max())
	}
}

// TestDequeFlattenParallel verifies Flatten correctness and parallel execution for Deque.
func TestDequeFlattenParallel(t *testing.T) {
	// Create nested test data
	nestedData := []any{
		SliceOf(1, 2, 3),
		SliceOf(4, 5, 6),
		SliceOf("a", "b"),
		SliceOf(7, 8, 9, 10),
		SliceOf("c", "d", "e"),
	}

	// Create deque with nested data
	dq := NewDeque[any]()
	for _, item := range nestedData {
		dq.PushBack(item)
	}

	workers := Int(3)
	cc := &concurrentCounterDeque{sleep: 50 * time.Millisecond}

	start := time.Now()
	result := dq.
		Iter().
		Parallel(workers).
		Inspect(func(v any) { cc.Fn(0) }). // Use generic counter
		Flatten().
		Collect()
	duration := time.Since(start)

	// Expected flattened elements: 3+3+2+4+3 = 15
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

	// Verify all expected elements are present
	expectedElements := []any{1, 2, 3, 4, 5, 6, "a", "b", 7, 8, 9, 10, "c", "d", "e"}
	for _, expected := range expectedElements {
		found := false
		result.Iter().Range(func(actual any) bool {
			if actual == expected {
				found = true
				return false
			}
			return true
		})
		if !found {
			t.Errorf("Missing expected element: %v", expected)
		}
	}

	t.Logf("Deque Flatten - Max concurrency: %d, Duration: %v, Items: %d",
		cc.Max(), duration, result.Len())
}
