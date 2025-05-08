package g_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestParallelCollect(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	got := SliceOf(nums...).Iter().Parallel(3).Collect()

	if len(got) != len(nums) {
		t.Fatalf("Collect returned %d items, want %d", len(got), len(nums))
	}

	for _, v := range nums {
		if !got.Contains(v) {
			t.Errorf("Collect missing element %d in result %v", v, got)
		}
	}
}

func TestMapFilterP(t *testing.T) {
	nums := Slice[int]{1, 2, 3, 4, 5}

	res := SliceOf(nums...).Iter().Parallel(2).
		Map(func(v int) int { return v * 2 }).
		Filter(func(v int) bool { return v%4 == 0 }).
		Collect()

	res.SortBy(cmp.Cmp)
	expected := Slice[int]{4, 8}

	if res.Ne(expected) {
		t.Errorf("Map+Filter got %v, want %v", res, expected)
	}
}

func TestChain(t *testing.T) {
	a := SliceOf(1, 2)
	b := SliceOf(3, 4)

	res := a.Iter().Parallel(2).Chain(b.Iter().Parallel(2)).Collect()
	res.SortBy(cmp.Cmp)

	expected := Slice[int]{1, 2, 3, 4}

	if res.Ne(expected) {
		t.Errorf("Chain got %v, want %v", res, expected)
	}
}

func TestAllAnyCount(t *testing.T) {
	nums := []int{2, 4, 6, 8}
	seq := SliceOf(nums...).Iter().Parallel(4)

	if !seq.All(func(v int) bool { return v%2 == 0 }) {
		t.Error("All returned false for even slice")
	}

	if seq.Any(func(v int) bool { return v == 5 }) {
		t.Error("Any returned true for missing element")
	}

	if cnt := seq.Count().Std(); cnt != len(nums) {
		t.Errorf("Count returned %d, want %d", cnt, len(nums))
	}
}

func TestFindPartition(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	seq := SliceOf(nums...).Iter().Parallel(3)

	opt := seq.Find(func(v int) bool { return v > 3 })
	if !opt.IsSome() || opt.Some() < 4 {
		t.Errorf("Find got %v, want Some(4) or Some(5)", opt)
	}

	left, right := seq.Partition(func(v int) bool { return v%2 == 0 })
	left.SortBy(cmp.Cmp)
	right.SortBy(cmp.Cmp)

	expectedLeft := Slice[int]{2, 4}
	expectedRight := Slice[int]{1, 3, 5}

	if left.Ne(expectedLeft) {
		t.Errorf("Partition left got %v, want %v", left, expectedLeft)
	}

	if right.Ne(expectedRight) {
		t.Errorf("Partition right got %v, want %v", right, expectedRight)
	}
}

func TestFlatten(t *testing.T) {
	nested := SliceOf[any](
		[]int{1, 2},
		[]int{3},
		4,
		[]int{},
		[]int{5, 6},
	).Iter().Parallel(2).Flatten().Collect()

	expected := Slice[int]{1, 2, 3, 4, 5, 6}
	_transformed := TransformSlice(nested, func(t any) int { return t.(int) })
	_transformed.SortBy(cmp.Cmp)

	if _transformed.Ne(expected) {
		t.Errorf("Flatten got %v, want %v", nested, expected)
	}
}

func TestFoldForEach(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	sum := SliceOf(nums...).Iter().Parallel(2).Fold(0, func(acc, v int) int { return acc + v })
	if sum != 10 {
		t.Errorf("Fold got %d, want 10", sum)
	}

	var collected []int
	SliceOf(nums...).Iter().Parallel(3).ForEach(func(v int) { collected = append(collected, v) })

	if len(collected) != len(nums) {
		t.Errorf("ForEach appended %d items, want %d", len(collected), len(nums))
	}
}

func TestTimeoutLimit(t *testing.T) {
	n := 6
	limit := Int(2)
	start := time.Now()

	SliceOf(1, 2, 3, 4, 5, 6).Iter().Parallel(limit).
		Map(func(v int) int {
			time.Sleep(100 * time.Millisecond)
			return v
		}).
		Collect()

	elapsed := time.Since(start)
	expected := time.Duration((n+limit.Std()-1)/limit.Std()) * 100 * time.Millisecond

	if elapsed < expected {
		t.Errorf("Parallel with limit %d too fast: %v < %v", limit, elapsed, expected)
	}
}

func TestUnique(t *testing.T) {
	nums := []int{1, 2, 2, 3, 1, 4}
	res := SliceOf(nums...).Iter().Parallel(3).Unique().Collect()

	if len(res) != 4 {
		t.Fatalf("Unique returned %d items, want 4", len(res))
	}

	for _, v := range []int{1, 2, 3, 4} {
		if !res.Contains(v) {
			t.Errorf("Unique missing element %d in result %v", v, res)
		}
	}
}

func TestTake(t *testing.T) {
	nums := []int{10, 20, 30, 40, 50}
	res := SliceOf(nums...).Iter().Parallel(1).Take(3).Collect()
	expected := Slice[int]{10, 20, 30}

	if !res.Eq(expected) {
		t.Errorf("Take got %v, want %v", res, expected)
	}
}

func TestSkip(t *testing.T) {
	nums := []int{5, 6, 7, 8, 9}
	res := SliceOf(nums...).Iter().Parallel(1).Skip(2).Collect()
	expected := Slice[int]{7, 8, 9}

	if !res.Eq(expected) {
		t.Errorf("Skip got %v, want %v", res, expected)
	}
}

func TestExclude(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	res := SliceOf(nums...).Iter().Parallel(2).Exclude(func(v int) bool { return v%2 == 0 }).Collect()

	if len(res) != 3 {
		t.Fatalf("Exclude returned %d items, want 3", len(res))
	}

	for _, v := range []int{1, 3, 5} {
		if !res.Contains(v) {
			t.Errorf("Exclude missing element %d", v)
		}
	}
}

func TestInspect(t *testing.T) {
	nums := []int{100, 200, 300}
	var (
		seen []int
		mu   sync.Mutex
	)

	res := SliceOf(nums...).Iter().Parallel(2).
		Inspect(func(v int) {
			mu.Lock()
			seen = append(seen, v)
			mu.Unlock()
		}).
		Collect()

	if len(res) != len(nums) {
		t.Fatalf("Inspect changed output length: got %d, want %d", len(res), len(nums))
	}

	seenSlice := SliceOf(seen...)
	seenSlice.SortBy(cmp.Cmp)

	expected := Slice[int]{100, 200, 300}
	if !seenSlice.Eq(expected) {
		t.Errorf("Inspect saw %v, want %v", seenSlice, expected)
	}
}
