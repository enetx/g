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
