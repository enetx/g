package g_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

type concurrentCounter struct {
	inFlight    int64
	maxInFlight int64
	sleep       time.Duration
}

func (cc *concurrentCounter) Fn(int) {
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

func (cc *concurrentCounter) Max() int64 { return atomic.LoadInt64(&cc.maxInFlight) }

// TestParallelCollect verifies Collect correctness and that multiple workers run concurrently.
func TestParallelCollect(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	workers := Int(3)
	cc := &concurrentCounter{sleep: 50 * time.Millisecond}

	got := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Collect()

	if len(got) != len(nums) {
		t.Fatalf("Collect returned %d items, want %d", len(got), len(nums))
	}

	for _, v := range nums {
		if !got.Contains(v) {
			t.Errorf("Collect missing element %d in result %v", v, got)
		}
	}

	if cc.Max() < 2 {
		t.Errorf("expected at least 2 concurrent tasks, got max %d", cc.Max())
	}
}

// TestMapFilterParallel verifies Map+Filter correctness and that Map runs in parallel.
func TestMapFilterParallel(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	workers := Int(2)
	cc := &concurrentCounter{sleep: 30 * time.Millisecond}

	res := SliceOf(input...).
		Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Map(func(v int) int { return v * 2 }).
		Filter(func(v int) bool { return v%4 == 0 }).
		Collect()

	expected := Slice[int]{4, 8}
	res.SortBy(cmp.Cmp)

	if res.Ne(expected) {
		t.Errorf("Map+Filter got %v, want %v", res, expected)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Map, got max concurrency %d", cc.Max())
	}
}

// TestChainParallel verifies Chain correctness and parallel execution across both sequences.
func TestChainParallel(t *testing.T) {
	a := []int{1, 2}
	b := []int{3, 4}
	workers := Int(2)
	cc := &concurrentCounter{sleep: 20 * time.Millisecond}

	res := SliceOf(a...).
		Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Chain(
			SliceOf(b...).Iter().Parallel(workers).Inspect(cc.Fn),
		).
		Collect()

	expected := Slice[int]{1, 2, 3, 4}
	res.SortBy(cmp.Cmp)
	if res.Ne(expected) {
		t.Errorf("Chain got %v, want %v", res, expected)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Chain, got max concurrency %d", cc.Max())
	}
}

// TestAllAnyCountParallel verifies All, Any, Count and that All and Count run in parallel.
func TestAllAnyCountParallel(t *testing.T) {
	nums := []int{2, 4, 6, 8}
	workers := Int(4)
	cc := &concurrentCounter{sleep: 10 * time.Millisecond}

	seq := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(cc.Fn)

	if !seq.All(func(v int) bool { return v%2 == 0 }) {
		t.Error("All returned false for even slice")
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel All, got max concurrency %d", cc.Max())
	}

	if seq.Any(func(v int) bool { return v == 5 }) {
		t.Error("Any returned true for missing element")
	}

	cc2 := &concurrentCounter{sleep: 10 * time.Millisecond}
	cnt := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(cc2.Fn).
		Count().Std()

	if cnt != len(nums) {
		t.Errorf("Count returned %d, want %d", cnt, len(nums))
	}

	if cc2.Max() < 2 {
		t.Errorf("expected parallel Count, got max concurrency %d", cc2.Max())
	}
}

// TestFindPartitionParallel verifies Find and Partition correctness and parallelism.
func TestFindPartitionParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	workers := Int(3)

	ccFind := &concurrentCounter{sleep: 15 * time.Millisecond}
	opt := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(ccFind.Fn).
		Find(func(v int) bool { return v > 3 })

	if !opt.IsSome() || opt.Some() < 4 {
		t.Errorf("Find got %v, want Some(4) or Some(5)", opt)
	}

	if ccFind.Max() < 2 {
		t.Errorf("expected parallel Find, got max %d", ccFind.Max())
	}

	ccPart := &concurrentCounter{sleep: 15 * time.Millisecond}
	left, right := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(ccPart.Fn).
		Partition(func(v int) bool { return v%2 == 0 })

	left.SortBy(cmp.Cmp)
	right.SortBy(cmp.Cmp)

	if !left.Eq(Slice[int]{2, 4}) {
		t.Errorf("Partition left got %v, want %v", left, []int{2, 4})
	}

	if !right.Eq(Slice[int]{1, 3, 5}) {
		t.Errorf("Partition right got %v, want %v", right, []int{1, 3, 5})
	}

	if ccPart.Max() < 2 {
		t.Errorf("expected parallel Partition, got max %d", ccPart.Max())
	}
}

// TestFlattenParallel verifies Flatten correctness and that it runs in parallel.
func TestFlattenParallel(t *testing.T) {
	nested := []any{
		[]int{1, 2},
		[]int{3},
		4,
		[]int{},
		[]int{5, 6},
	}

	workers := Int(2)
	cc := &concurrentCounter{sleep: 20 * time.Millisecond}

	nestedSeq := SliceOf(nested...).
		Iter().
		Parallel(workers).
		Flatten().
		Inspect(func(v any) {
			cc.Fn(v.(int))
		}).
		Collect()

	trans := TransformSlice(nestedSeq, func(x any) int { return x.(int) })
	trans.SortBy(cmp.Cmp)

	expected := Slice[int]{1, 2, 3, 4, 5, 6}

	if !trans.Eq(expected) {
		t.Errorf("Flatten got %v, want %v", trans, expected)
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Flatten, got max %d", cc.Max())
	}
}

// TestFoldForEachParallel verifies Fold, ForEach and their parallel behavior.
func TestFoldForEachParallel(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	workers := Int(2)

	ccFold := &concurrentCounter{sleep: 10 * time.Millisecond}
	sum := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(ccFold.Fn).
		Fold(0, func(acc, v int) int { return acc + v })

	if sum != 10 {
		t.Errorf("Fold got %d, want 10", sum)
	}

	if ccFold.Max() < 2 {
		t.Errorf("expected parallel Fold, got max %d", ccFold.Max())
	}

	var collected []int
	var mu sync.Mutex
	ccFor := &concurrentCounter{sleep: 10 * time.Millisecond}

	SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(ccFor.Fn).
		ForEach(func(v int) {
			mu.Lock()
			collected = append(collected, v)
			mu.Unlock()
		})

	if len(collected) != len(nums) {
		t.Errorf("ForEach appended %d items, want %d", len(collected), len(nums))
	}

	if ccFor.Max() < 2 {
		t.Errorf("expected parallel ForEach, got max %d", ccFor.Max())
	}
}

// TestTimeoutLimitParallel verifies that limiting workers affects timing as expected.
func TestTimeoutLimitParallel(t *testing.T) {
	n := 6
	limit := Int(2)
	sleep := 50 * time.Millisecond
	cc := &concurrentCounter{sleep: sleep}

	start := time.Now()

	SliceOf(1, 2, 3, 4, 5, 6).
		Iter().
		Parallel(limit).
		Inspect(cc.Fn).
		Collect()

	elapsed := time.Since(start)

	batches := (n + limit.Std() - 1) / limit.Std()
	min := time.Duration(batches)*sleep - 10*time.Millisecond
	max := time.Duration(batches)*sleep + 50*time.Millisecond

	if elapsed < min || elapsed > max {
		t.Errorf("timing off with limit %d: elapsed %v, want ~%v", limit, elapsed, time.Duration(batches)*sleep)
	}

	if cc.Max() != int64(limit.Std()) {
		t.Errorf("expected max concurrency %d, got %d", limit, cc.Max())
	}
}

// TestUniqueParallel verifies Unique correctness and parallel execution.
func TestUniqueParallel(t *testing.T) {
	nums := []int{1, 2, 2, 3, 1, 4}
	workers := Int(3)
	cc := &concurrentCounter{sleep: 10 * time.Millisecond}

	res := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Unique().
		Collect()

	if len(res) != 4 {
		t.Fatalf("Unique returned %d items, want 4", len(res))
	}

	for _, v := range []int{1, 2, 3, 4} {
		if !res.Contains(v) {
			t.Errorf("Unique missing element %d in result %v", v, res)
		}
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Unique, got max %d", cc.Max())
	}
}

// TestTakeSkipExcludeInspectParallel verifies Take, Skip, Exclude, and Inspect parallelism.
func TestTakeSkipExcludeInspectParallel(t *testing.T) {
	nums1 := []int{10, 20, 30, 40, 50}
	nums2 := []int{5, 6, 7, 8, 9}
	nums3 := []int{1, 2, 3, 4, 5}
	workers := Int(2)
	sleep := 10 * time.Millisecond

	ccTake := &concurrentCounter{sleep: sleep}
	resTake := SliceOf(nums1...).
		Iter().
		Parallel(workers).
		Inspect(ccTake.Fn).
		Take(3).
		Collect()

	if len(resTake) != 3 {
		t.Fatalf("Take: expected 3 items, got %d", len(resTake))
	}

	if ccTake.Max() < 2 {
		t.Errorf("expected parallel Take, got max %d", ccTake.Max())
	}

	ccSkip := &concurrentCounter{sleep: sleep}
	resSkip := SliceOf(nums2...).
		Iter().
		Parallel(workers).
		Inspect(ccSkip.Fn).
		Skip(2).
		Collect()

	resSkip.SortBy(cmp.Cmp)

	if !resSkip.Eq(Slice[int]{7, 8, 9}) {
		t.Errorf("Skip got %v, want %v", resSkip, []int{7, 8, 9})
	}

	if ccSkip.Max() < 2 {
		t.Errorf("expected parallel Skip, got max %d", ccSkip.Max())
	}

	ccEx := &concurrentCounter{sleep: sleep}
	resEx := SliceOf(nums3...).
		Iter().
		Parallel(workers).
		Inspect(ccEx.Fn).
		Exclude(func(v int) bool { return v%2 == 0 }).
		Collect()

	if len(resEx) != 3 {
		t.Fatalf("Exclude returned %d items, want 3", len(resEx))
	}

	for _, v := range []int{1, 3, 5} {
		if !resEx.Contains(v) {
			t.Errorf("Exclude missing element %d", v)
		}
	}

	if ccEx.Max() < 2 {
		t.Errorf("expected parallel Exclude, got max %d", ccEx.Max())
	}

	var seen []int
	var mu sync.Mutex

	ccIns := &concurrentCounter{sleep: sleep}
	resIns := SliceOf([]int{100, 200, 300}...).
		Iter().
		Parallel(workers).
		Inspect(func(v int) {
			mu.Lock()
			seen = append(seen, v)
			mu.Unlock()
			ccIns.Fn(v)
		}).
		Collect()

	if len(resIns) != 3 {
		t.Fatalf("Inspect changed output length: got %d, want %d", len(resIns), 3)
	}

	seenSlice := SliceOf(seen...)
	seenSlice.SortBy(cmp.Cmp)

	if !seenSlice.Eq(Slice[int]{100, 200, 300}) {
		t.Errorf("Inspect saw %v, want %v", seenSlice, []int{100, 200, 300})
	}

	if ccIns.Max() < 2 {
		t.Errorf("expected parallel Inspect, got max %d", ccIns.Max())
	}
}

// Enhanced concurrent counter with more detailed tracking
type enhancedCounter struct {
	current       int64
	max           int64
	total         int64
	sleepDuration time.Duration
	sequenceID    string
}

func (cc *enhancedCounter) Fn(v int) {
	current := atomic.AddInt64(&cc.current, 1)
	atomic.AddInt64(&cc.total, 1)

	// Update max concurrency
	for {
		currentMax := atomic.LoadInt64(&cc.max)
		if current <= currentMax || atomic.CompareAndSwapInt64(&cc.max, currentMax, current) {
			break
		}
	}

	// Simulate work
	time.Sleep(cc.sleepDuration)
	atomic.AddInt64(&cc.current, -1)
}

func (cc *enhancedCounter) Max() int64 {
	return atomic.LoadInt64(&cc.max)
}

func (cc *enhancedCounter) Total() int64 {
	return atomic.LoadInt64(&cc.total)
}

// TestChainComprehensive tests Chain with multiple aspects
func TestChainComprehensive(t *testing.T) {
	t.Run("BasicParallelism", func(t *testing.T) {
		// Test basic parallel execution
		seq1 := []int{1, 2, 3, 4, 5}
		seq2 := []int{10, 20, 30, 40, 50}
		workers := Int(3)

		cc := &enhancedCounter{sleepDuration: 50 * time.Millisecond}

		start := time.Now()
		res := SliceOf(seq1...).
			Iter().
			Parallel(workers).
			Inspect(cc.Fn).
			Chain(
				SliceOf(seq2...).Iter().Parallel(workers).Inspect(cc.Fn),
			).
			Collect()
		duration := time.Since(start)

		// Verify results
		expected := append(seq1, seq2...)
		res.SortBy(cmp.Cmp)
		expectedSlice := SliceOf(expected...)
		expectedSlice.SortBy(cmp.Cmp)

		if res.Ne(expectedSlice) {
			t.Errorf("Chain got %v, want %v", res, expectedSlice)
		}

		// Verify parallelism (should be at least 2, ideally close to 6)
		if cc.Max() < 2 {
			t.Errorf("expected parallel execution, got max concurrency %d", cc.Max())
		}

		// Verify timing - with parallelism should be much faster than sequential
		expectedSequentialTime := time.Duration(len(seq1)+len(seq2)) * 50 * time.Millisecond
		if duration > expectedSequentialTime/2 {
			t.Errorf("execution too slow, might not be parallel: %v", duration)
		}

		t.Logf("Max concurrency: %d, Total processed: %d, Duration: %v",
			cc.Max(), cc.Total(), duration)
	})

	t.Run("MultipleSequencesChain", func(t *testing.T) {
		// Test chaining multiple sequences
		seq1 := []int{1, 2}
		seq2 := []int{10, 20}
		seq3 := []int{100, 200}
		seq4 := []int{1000, 2000}

		workers := Int(4)
		cc := &enhancedCounter{sleepDuration: 30 * time.Millisecond}

		res := SliceOf(seq1...).
			Iter().
			Parallel(workers).
			Inspect(cc.Fn).
			Chain(
				SliceOf(seq2...).Iter().Parallel(workers).Inspect(cc.Fn),
				SliceOf(seq3...).Iter().Parallel(workers).Inspect(cc.Fn),
				SliceOf(seq4...).Iter().Parallel(workers).Inspect(cc.Fn),
			).
			Collect()

		// All elements should be present
		if res.Len() != 8 {
			t.Errorf("expected 8 elements, got %d", res.Len())
		}

		// Should achieve high concurrency with 4 sequences
		if cc.Max() < 4 {
			t.Errorf("expected high concurrency with 4 sequences, got %d", cc.Max())
		}

		t.Logf("Multiple sequences - Max concurrency: %d", cc.Max())
	})

	t.Run("HeavyTransformationsParallel", func(t *testing.T) {
		// Test that heavy transformations in each sequence run in parallel
		seq1 := []int{1, 2, 3, 4, 5, 6, 7, 8}
		seq2 := []int{10, 20, 30, 40, 50, 60, 70, 80}

		workers1 := Int(3)
		workers2 := Int(5)

		cc1 := &enhancedCounter{sleepDuration: 40 * time.Millisecond, sequenceID: "seq1"}
		cc2 := &enhancedCounter{sleepDuration: 40 * time.Millisecond, sequenceID: "seq2"}

		heavyTransform := func(x int) int {
			time.Sleep(20 * time.Millisecond) // Additional heavy work
			return x * 2
		}

		start := time.Now()
		res := SliceOf(seq1...).
			Iter().
			Parallel(workers1).
			Map(heavyTransform).
			Inspect(cc1.Fn).
			Chain(
				SliceOf(seq2...).
					Iter().
					Parallel(workers2).
					Map(heavyTransform).
					Inspect(cc2.Fn),
			).
			Collect()
		duration := time.Since(start)

		// Verify both sequences achieved parallelism
		if cc1.Max() < 2 {
			t.Errorf("seq1 not parallel enough, max concurrency: %d", cc1.Max())
		}
		if cc2.Max() < 2 {
			t.Errorf("seq2 not parallel enough, max concurrency: %d", cc2.Max())
		}

		// Total concurrency should be sum of both sequences
		totalExpectedConcurrency := cc1.Max() + cc2.Max()
		if totalExpectedConcurrency < 4 {
			t.Errorf("total parallelism too low: seq1=%d, seq2=%d", cc1.Max(), cc2.Max())
		}

		// Verify results are transformed correctly
		expectedLen := len(seq1) + len(seq2)
		if len(res) != expectedLen {
			t.Errorf("expected %d elements, got %d", expectedLen, res.Len())
		}

		t.Logf("Heavy transforms - Seq1 concurrency: %d, Seq2 concurrency: %d, Duration: %v",
			cc1.Max(), cc2.Max(), duration)
	})

	t.Run("EarlyTermination", func(t *testing.T) {
		// Test early termination works correctly
		largeSeq1 := make([]int, 1000)
		largeSeq2 := make([]int, 1000)
		for i := range largeSeq1 {
			largeSeq1[i] = i
			largeSeq2[i] = i + 1000
		}

		workers := Int(4)
		var processedCount atomic.Int64

		start := time.Now()
		res := SliceOf(largeSeq1...).
			Iter().
			Parallel(workers).
			Inspect(func(v int) {
				processedCount.Add(1)
				time.Sleep(1 * time.Millisecond)
			}).
			Chain(
				SliceOf(largeSeq2...).
					Iter().
					Parallel(workers).
					Inspect(func(v int) {
						processedCount.Add(1)
						time.Sleep(1 * time.Millisecond)
					}),
			).
			Take(10). // Should stop early
			Collect()
		duration := time.Since(start)

		// Should only get 10 elements
		if res.Len() != 10 {
			t.Errorf("expected 10 elements with Take(10), got %d", res.Len())
		}

		// Should process significantly fewer than 2000 elements
		processed := processedCount.Load()
		if processed > 100 {
			t.Logf("Warning: processed %d elements, early termination might not be working optimally", processed)
		}

		// Should complete much faster than processing all elements
		maxExpectedDuration := 200 * time.Millisecond
		if duration > maxExpectedDuration {
			t.Errorf("early termination too slow: %v", duration)
		}

		t.Logf("Early termination - Processed: %d elements, Duration: %v", processed, duration)
	})

	t.Run("DifferentWorkerCounts", func(t *testing.T) {
		// Test sequences with different worker counts
		seq1 := make([]int, 20)
		seq2 := make([]int, 20)
		for i := range seq1 {
			seq1[i] = i
			seq2[i] = i + 100
		}

		workers1 := Int(2)
		workers2 := Int(8)

		cc1 := &enhancedCounter{sleepDuration: 25 * time.Millisecond}
		cc2 := &enhancedCounter{sleepDuration: 25 * time.Millisecond}

		res := SliceOf(seq1...).
			Iter().
			Parallel(workers1).
			Inspect(cc1.Fn).
			Chain(
				SliceOf(seq2...).
					Iter().
					Parallel(workers2).
					Inspect(cc2.Fn),
			).
			Collect()

		// Verify different concurrency levels
		if cc1.Max() > int64(workers1)+1 { // +1 for some tolerance
			t.Errorf("seq1 exceeded expected concurrency: got %d, expected ~%d", cc1.Max(), workers1)
		}
		if cc2.Max() > int64(workers2)+1 {
			t.Errorf("seq2 exceeded expected concurrency: got %d, expected ~%d", cc2.Max(), workers2)
		}

		// Both should achieve some level of parallelism
		if cc1.Max() < 1 || cc2.Max() < 2 {
			t.Errorf("sequences didn't achieve expected parallelism: seq1=%d, seq2=%d", cc1.Max(), cc2.Max())
		}

		// All elements should be present
		if res.Len() != 40 {
			t.Errorf("expected 40 elements, got %d", res.Len())
		}

		t.Logf("Different workers - Seq1 (%d workers): %d concurrency, Seq2 (%d workers): %d concurrency",
			workers1, cc1.Max(), workers2, cc2.Max())
	})
}

func TestMapWithFilteredInput(t *testing.T) {
	// Test Map when previous process function returns ok=false
	// This tests the "else" branch in Map function (lines 329-330)
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	workers := Int(2)

	// Filter out even numbers, then map the remaining odd numbers
	result := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Filter(func(v int) bool { return v%2 == 1 }). // This creates a process that can return ok=false
		Map(func(v int) int { return v * 10 }).       // This should handle the ok=false case
		Collect()

	// Should contain only odd numbers multiplied by 10
	expected := Slice[int]{10, 30, 50, 70, 90}
	result.SortBy(cmp.Cmp)

	if result.Ne(expected) {
		t.Errorf("Map with filtered input got %v, want %v", result, expected)
	}
}

func TestAnyParallelEarlyExit(t *testing.T) {
	// Test Any with early exit - should stop as soon as true is found
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	workers := Int(3)
	cc := &concurrentCounter{sleep: 20 * time.Millisecond}

	// Look for number 3 - should exit early and not process all elements
	found := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(cc.Fn).
		Any(func(v int) bool { return v == 3 })

	if !found {
		t.Errorf("Any should have found element 3")
	}

	// Should achieve some parallelism but may not process all elements due to early exit
	if cc.Max() < 1 {
		t.Errorf("Expected some parallel processing, got max concurrency %d", cc.Max())
	}

	// Test case where nothing is found - should process all elements
	cc2 := &concurrentCounter{sleep: 10 * time.Millisecond}
	notFound := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Inspect(cc2.Fn).
		Any(func(v int) bool { return v == 99 }) // Element that doesn't exist

	if notFound {
		t.Errorf("Any should not have found non-existent element 99")
	}

	if cc2.Max() < 2 {
		t.Errorf("Expected parallel processing when scanning all elements, got max concurrency %d", cc2.Max())
	}
}

func TestInspectWithFilteredInput(t *testing.T) {
	// Test Inspect when previous process function returns ok=false
	// This tests the "else" branch in Inspect function
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	workers := Int(2)

	var inspectedValues []int
	var mu sync.Mutex

	// Filter out even numbers, then inspect the remaining odd numbers
	result := SliceOf(nums...).
		Iter().
		Parallel(workers).
		Filter(func(v int) bool { return v%2 == 1 }). // This creates a process that can return ok=false
		Inspect(func(v int) {                         // This should handle the ok=false case
			mu.Lock()
			inspectedValues = append(inspectedValues, v)
			mu.Unlock()
		}).
		Collect()

	// Should contain only odd numbers
	expected := Slice[int]{1, 3, 5, 7, 9}
	result.SortBy(cmp.Cmp)

	if result.Ne(expected) {
		t.Errorf("Inspect with filtered input got %v, want %v", result, expected)
	}

	// Check that inspected values match the collected values
	inspectedSlice := SliceOf(inspectedValues...)
	inspectedSlice.SortBy(cmp.Cmp)

	if inspectedSlice.Ne(expected) {
		t.Errorf("Inspected values got %v, want %v", inspectedSlice, expected)
	}
}
