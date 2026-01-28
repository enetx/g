package g_test

import (
	"sort"
	"sync"
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

	skipSlice := skipRes.Iter().Collect().Slice()
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

	takeSlice := takeRes.Iter().Collect().Slice()
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

	result := res.Iter().Collect().Slice()
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

// TestDequeParFlatMap tests the new FlatMap method for SeqDequePar
func TestDequeParFlatMap(t *testing.T) {
	t.Run("basic flat mapping", func(t *testing.T) {
		dq := NewDeque[int]()
		dq.PushBack(1)
		dq.PushBack(2)
		dq.PushBack(3)

		result := dq.Iter().
			Parallel(2).
			FlatMap(func(x int) SeqDeque[int] {
				innerDq := NewDeque[int]()
				innerDq.PushBack(x)
				innerDq.PushBack(x * 10)
				return innerDq.Iter()
			}).
			Collect()

		if result.Len() != 6 {
			t.Errorf("Expected 6 elements, got %d", result.Len())
		}

		// Check that we have the expected values (order may vary due to parallelism)
		valueCount := make(map[int]int)
		result.Iter().Range(func(v int) bool {
			valueCount[v]++
			return true
		})

		expected := map[int]int{1: 1, 10: 1, 2: 1, 20: 1, 3: 1, 30: 1}
		for k, v := range expected {
			if valueCount[k] != v {
				t.Errorf("Expected %d occurrences of %d, got %d", v, k, valueCount[k])
			}
		}
	})

	t.Run("empty input", func(t *testing.T) {
		dq := NewDeque[int]()
		result := dq.Iter().
			Parallel(2).
			FlatMap(func(x int) SeqDeque[int] {
				innerDq := NewDeque[int]()
				innerDq.PushBack(x * 2)
				return innerDq.Iter()
			}).
			Collect()

		if result.Len() != 0 {
			t.Errorf("Expected empty result, got %d elements", result.Len())
		}
	})

	t.Run("parallelism verification", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 0; i < 20; i++ {
			dq.PushBack(i)
		}

		cc := &concurrentCounterDeque{sleep: 30 * time.Millisecond}

		result := dq.Iter().
			Parallel(4).
			Inspect(cc.Fn).
			FlatMap(func(x int) SeqDeque[int] {
				innerDq := NewDeque[int]()
				innerDq.PushBack(x)
				return innerDq.Iter()
			}).
			Collect()

		if result.Len() != 20 {
			t.Errorf("Expected 20 elements, got %d", result.Len())
		}

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})
}

// TestDequeParFilterMap tests the new FilterMap method for SeqDequePar
func TestDequeParFilterMap(t *testing.T) {
	t.Run("filter and transform", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 6; i++ {
			dq.PushBack(i)
		}

		result := dq.Iter().
			Parallel(2).
			FilterMap(func(x int) Option[int] {
				if x%2 == 0 {
					return Some(x * 10)
				}
				return None[int]()
			}).
			Collect()

		if result.Len() != 3 {
			t.Errorf("Expected 3 elements, got %d", result.Len())
		}

		valueCount := make(map[int]int)
		result.Iter().Range(func(v int) bool {
			valueCount[v]++
			return true
		})

		expected := map[int]int{20: 1, 40: 1, 60: 1}
		for k, v := range expected {
			if valueCount[k] != v {
				t.Errorf("Expected %d occurrences of %d, got %d", v, k, valueCount[k])
			}
		}
	})

	t.Run("all filtered out", func(t *testing.T) {
		dq := NewDeque[int]()
		dq.PushBack(1)
		dq.PushBack(3)
		dq.PushBack(5)
		dq.PushBack(7)

		result := dq.Iter().
			Parallel(2).
			FilterMap(func(x int) Option[int] {
				if x%2 == 0 {
					return Some(x * 10)
				}
				return None[int]()
			}).
			Collect()

		if result.Len() != 0 {
			t.Errorf("Expected empty result, got %d elements", result.Len())
		}
	})

	t.Run("parallelism with filtering", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 0; i < 100; i++ {
			dq.PushBack(i)
		}

		cc := &concurrentCounterDeque{sleep: 10 * time.Millisecond}

		result := dq.Iter().
			Parallel(4).
			Inspect(cc.Fn).
			FilterMap(func(x int) Option[int] {
				if x%3 == 0 {
					return Some(x * 2)
				}
				return None[int]()
			}).
			Collect()

		// Should get numbers divisible by 3, transformed
		expectedCount := 100 / 3 // about 33
		if result.Len() < Int(expectedCount-1) || result.Len() > Int(expectedCount+1) {
			t.Errorf("Expected ~%d elements, got %d", expectedCount, result.Len())
		}

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})
}

// TestDequeParStepBy tests the new StepBy method for SeqDequePar
func TestDequeParStepBy(t *testing.T) {
	t.Run("step by 2", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 8; i++ {
			dq.PushBack(i)
		}

		result := dq.Iter().
			Parallel(2).
			StepBy(2).
			Collect()

		if result.Len() != 4 {
			t.Errorf("Expected 4 elements, got %d", result.Len())
		}
	})

	t.Run("step by 3", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 9; i++ {
			dq.PushBack(i)
		}

		result := dq.Iter().
			Parallel(2).
			StepBy(3).
			Collect()

		if result.Len() != 3 {
			t.Errorf("Expected 3 elements, got %d", result.Len())
		}
	})

	t.Run("step by 0 defaults to 1", func(t *testing.T) {
		dq := NewDeque[int]()
		dq.PushBack(1)
		dq.PushBack(2)
		dq.PushBack(3)

		result := dq.Iter().
			Parallel(2).
			StepBy(0).
			Collect()

		if result.Len() != 3 {
			t.Errorf("Expected 3 elements, got %d", result.Len())
		}
	})

	t.Run("parallel step counting", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 0; i < 50; i++ {
			dq.PushBack(i)
		}

		cc := &concurrentCounterDeque{sleep: 5 * time.Millisecond}

		result := dq.Iter().
			Parallel(4).
			Inspect(cc.Fn).
			StepBy(5).
			Collect()

		expectedCount := 10 // 50/5 = 10
		if result.Len() != Int(expectedCount) {
			t.Errorf("Expected %d elements, got %d", expectedCount, result.Len())
		}

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})
}

// TestDequeParMaxMinBy tests the new MaxBy/MinBy methods for SeqDequePar
func TestDequeParMaxMinBy(t *testing.T) {
	t.Run("find maximum", func(t *testing.T) {
		dq := NewDeque[int]()
		values := []int{3, 1, 4, 1, 5, 9, 2, 6}
		for _, v := range values {
			dq.PushBack(v)
		}

		result := dq.Iter().
			Parallel(2).
			MaxBy(func(a, b int) cmp.Ordering {
				return cmp.Cmp(a, b)
			})

		if result.IsNone() {
			t.Error("Expected Some value, got None")
		}

		if result.Some() != 9 {
			t.Errorf("Expected maximum 9, got %d", result.Some())
		}
	})

	t.Run("find minimum", func(t *testing.T) {
		dq := NewDeque[int]()
		values := []int{3, 1, 4, 1, 5, 9, 2, 6}
		for _, v := range values {
			dq.PushBack(v)
		}

		result := dq.Iter().
			Parallel(2).
			MinBy(func(a, b int) cmp.Ordering {
				return cmp.Cmp(a, b)
			})

		if result.IsNone() {
			t.Error("Expected Some value, got None")
		}

		if result.Some() != 1 {
			t.Errorf("Expected minimum 1, got %d", result.Some())
		}
	})

	t.Run("empty collection", func(t *testing.T) {
		dq := NewDeque[int]()

		maxResult := dq.Iter().
			Parallel(2).
			MaxBy(func(a, b int) cmp.Ordering {
				return cmp.Cmp(a, b)
			})

		minResult := dq.Iter().
			Parallel(2).
			MinBy(func(a, b int) cmp.Ordering {
				return cmp.Cmp(a, b)
			})

		if maxResult.IsSome() {
			t.Errorf("Expected None for max, got Some(%v)", maxResult.Some())
		}

		if minResult.IsSome() {
			t.Errorf("Expected None for min, got Some(%v)", minResult.Some())
		}
	})

	t.Run("custom comparison with parallelism", func(t *testing.T) {
		dq := NewDeque[string]()
		strings := []string{"a", "bb", "ccc", "d", "ee", "ffff"}
		for _, s := range strings {
			dq.PushBack(s)
		}

		cc := &concurrentCounterDeque{sleep: 20 * time.Millisecond}

		maxResult := dq.Iter().
			Parallel(3).
			Inspect(func(s string) { cc.Fn(len(s)) }).
			MaxBy(func(a, b string) cmp.Ordering {
				return cmp.Cmp(len(a), len(b))
			})

		if maxResult.IsNone() {
			t.Error("Expected Some value for max, got None")
		}

		if maxResult.Some() != "ffff" {
			t.Errorf("Expected longest string 'ffff', got %s", maxResult.Some())
		}

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})
}

// TestDequeParExclude tests the Exclude method for parallel deque iterators.
func TestDequeParExclude(t *testing.T) {
	t.Run("basic_exclude", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 10; i++ {
			dq.PushBack(i)
		}

		// Exclude even numbers
		result := dq.Iter().
			Parallel(3).
			Exclude(func(n int) bool {
				return n%2 == 0
			}).
			Collect()

		// Should have odd numbers: 1, 3, 5, 7, 9
		expected := []int{1, 3, 5, 7, 9}
		if result.Len() != Int(len(expected)) {
			t.Errorf("Expected %d elements, got %d", len(expected), result.Len())
		}

		resultSlice := result.Slice()
		resultSlice.SortBy(cmp.Cmp)
		for i, exp := range expected {
			if resultSlice[i] != exp {
				t.Errorf("Expected %d at index %d, got %d", exp, i, resultSlice[i])
			}
		}
	})

	t.Run("exclude_none", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 5; i++ {
			dq.PushBack(i)
		}

		// Exclude nothing (predicate always returns false)
		result := dq.Iter().
			Parallel(2).
			Exclude(func(n int) bool {
				return false
			}).
			Collect()

		if result.Len() != Int(5) {
			t.Errorf("Expected all 5 elements, got %d", result.Len())
		}
	})

	t.Run("exclude_all", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 5; i++ {
			dq.PushBack(i)
		}

		// Exclude everything (predicate always returns true)
		result := dq.Iter().
			Parallel(2).
			Exclude(func(n int) bool {
				return true
			}).
			Collect()

		if result.Len() != Int(0) {
			t.Errorf("Expected 0 elements, got %d", result.Len())
		}
	})

	t.Run("parallel_execution", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 20; i++ {
			dq.PushBack(i)
		}

		cc := &concurrentCounterDeque{sleep: 10 * time.Millisecond}

		result := dq.Iter().
			Parallel(4).
			Inspect(cc.Fn).
			Exclude(func(n int) bool {
				return n > 10 // Exclude numbers > 10
			}).
			Collect()

		// Should have numbers 1-10
		if result.Len() != Int(10) {
			t.Errorf("Expected 10 elements, got %d", result.Len())
		}

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})
}

// TestDequeParForEach tests the ForEach method for parallel deque iterators.
func TestDequeParForEach(t *testing.T) {
	t.Run("basic_foreach", func(t *testing.T) {
		dq := NewDeque[int]()
		for i := 1; i <= 10; i++ {
			dq.PushBack(i)
		}

		var processed []int
		var mu sync.Mutex

		dq.Iter().
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
		dq := NewDeque[int]()
		for i := 1; i <= 20; i++ {
			dq.PushBack(i)
		}

		cc := &concurrentCounterDeque{sleep: 10 * time.Millisecond}

		dq.Iter().
			Parallel(4).
			ForEach(func(n int) {
				cc.Fn(n)
			})

		if cc.Max() < 2 {
			t.Errorf("Expected parallel execution, got max concurrency %d", cc.Max())
		}
	})

	t.Run("empty_deque", func(t *testing.T) {
		dq := NewDeque[int]()
		executed := false

		dq.Iter().
			Parallel(2).
			ForEach(func(n int) {
				executed = true
			})

		if executed {
			t.Error("ForEach should not execute for empty deque")
		}
	})
}
