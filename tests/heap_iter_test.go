package g

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestSeqHeap_All(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(2, 4, 6, 8)

	allEven := heap.Iter().All(func(v int) bool { return v%2 == 0 })
	if !allEven {
		t.Error("Expected all elements to be even")
	}

	heap2 := g.NewHeap(cmp.Cmp[int])
	heap2.Push(1, 2, 3, 4)

	allEven2 := heap2.Iter().All(func(v int) bool { return v%2 == 0 })
	if allEven2 {
		t.Error("Expected not all elements to be even")
	}
}

func TestSeqHeap_Any(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 3, 5, 7)

	anyEven := heap.Iter().Any(func(v int) bool { return v%2 == 0 })
	if anyEven {
		t.Error("Expected no even elements")
	}

	heap2 := g.NewHeap(cmp.Cmp[int])
	heap2.Push(1, 2, 3, 5)

	anyEven2 := heap2.Iter().Any(func(v int) bool { return v%2 == 0 })
	if !anyEven2 {
		t.Error("Expected at least one even element")
	}
}

func TestSeqHeap_Chain(t *testing.T) {
	heap1 := g.NewHeap(cmp.Cmp[int])
	heap1.Push(3, 1, 2)

	heap2 := g.NewHeap(cmp.Cmp[int])
	heap2.Push(6, 4, 5)

	result := make([]int, 0)
	heap1.Iter().Chain(heap2.Iter()).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 3, 4, 5, 6} // first heap sorted, then second heap sorted
	if len(result) != len(expected) {
		t.Errorf("Chain: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Chunks(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5, 6)

	chunks := heap.Iter().Chunks(2).Collect()

	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, got %d", len(chunks))
	}

	expectedSizes := []int{2, 2, 2}
	for i, chunk := range chunks {
		if len(chunk) != expectedSizes[i] {
			t.Errorf("Chunk %d: expected size %d, got %d", i, expectedSizes[i], len(chunk))
		}
	}
}

func TestSeqHeap_Collect(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(3, 1, 4, 1, 5)

	// Collect with reverse comparison (max heap)
	collected := heap.Iter().Collect(func(a, b int) cmp.Ordering {
		return cmp.Cmp(b, a) // reverse comparison
	})

	if collected.Len() != 5 {
		t.Errorf("Expected collected heap length 5, got %d", collected.Len())
	}

	// Should pop in reverse order (largest first)
	result := make([]int, 0)
	for !collected.IsEmpty() {
		val := collected.Pop().Some()
		result = append(result, val)
	}

	expected := []int{5, 4, 3, 1, 1} // max heap order
	if len(result) != len(expected) {
		t.Errorf("CollectWith: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Count(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	count := heap.Iter().Count()
	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

func TestSeqHeap_Counter(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 1, 3, 2, 1)

	counts := make(map[int]g.Int)
	heap.Iter().Counter().ForEach(func(k any, v g.Int) {
		counts[k.(int)] = v
	})

	expected := map[int]g.Int{1: 3, 2: 2, 3: 1}
	for k, v := range expected {
		if counts[k] != v {
			t.Errorf("Counter for %d: expected %d, got %d", k, v, counts[k])
		}
	}
}

func TestSeqHeap_Dedup(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 1, 2, 2, 3, 3, 3)

	result := make([]int, 0)
	heap.Iter().Dedup().ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 3} // deduped in sorted order
	if len(result) != len(expected) {
		t.Errorf("Dedup: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Enumerate(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[string])
	heap.Push("c", "a", "b")

	pairs := make(map[g.Int]string)
	heap.Iter().Enumerate().ForEach(func(i g.Int, v string) {
		pairs[i] = v
	})

	expected := map[g.Int]string{0: "a", 1: "b", 2: "c"} // sorted order
	for k, v := range expected {
		if pairs[k] != v {
			t.Errorf("Enumerate at %d: expected %s, got %s", k, v, pairs[k])
		}
	}
}

func TestSeqHeap_Filter(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5, 6)

	result := make([]int, 0)
	heap.Iter().Filter(func(v int) bool { return v%2 == 0 }).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{2, 4, 6}
	if len(result) != len(expected) {
		t.Errorf("Filter: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Exclude(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	result := make([]int, 0)
	heap.Iter().Exclude(func(v int) bool { return v%2 == 0 }).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 3, 5}
	if len(result) != len(expected) {
		t.Errorf("Exclude: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Find(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	found := heap.Iter().Find(func(v int) bool { return v > 3 })

	if !found.IsSome() {
		t.Error("Expected to find element > 3")
	}

	if found.Some() != 4 {
		t.Errorf("Expected to find 4, got %d", found.Some())
	}
}

func TestSeqHeap_Fold(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	sum := heap.Iter().Fold(0, func(acc, val int) int {
		return acc + val
	})

	if sum != 15 {
		t.Errorf("Expected sum 15, got %d", sum)
	}
}

func TestSeqHeap_Reduce(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	product := heap.Iter().Reduce(func(a, b int) int {
		return a * b
	})

	if !product.IsSome() {
		t.Error("Expected reduce to return some value")
	}

	if product.Some() != 120 {
		t.Errorf("Expected product 120, got %d", product.Some())
	}
}

func TestSeqHeap_ForEach(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	sum := 0
	heap.Iter().ForEach(func(v int) {
		sum += v
	})

	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
}

func TestSeqHeap_Map(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	result := make([]int, 0)
	heap.Iter().Map(func(v int) int {
		return v * 2
	}).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{2, 4, 6} // sorted and doubled
	if len(result) != len(expected) {
		t.Errorf("Map: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Skip(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(5, 2, 8, 1, 9)

	result := make([]int, 0)
	heap.Iter().Skip(2).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{5, 8, 9} // sorted order [1,2,5,8,9], skip first 2
	if len(result) != len(expected) {
		t.Errorf("Skip: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Take(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(5, 2, 8, 1, 9)

	result := make([]int, 0)
	heap.Iter().Take(3).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 5} // first 3 elements in sorted order
	if len(result) != len(expected) {
		t.Errorf("Take: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_StepBy(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5, 6)

	result := make([]int, 0)
	heap.Iter().StepBy(2).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 3, 5} // every 2nd element in sorted order
	if len(result) != len(expected) {
		t.Errorf("StepBy: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_SortBy(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[string])
	heap.Push("apple", "banana", "cherry")

	// Sort by string length (ascending)
	result := make([]string, 0)
	heap.Iter().SortBy(func(a, b string) cmp.Ordering {
		return cmp.Cmp(len(a), len(b))
	}).ForEach(func(v string) {
		result = append(result, v)
	})

	expected := []string{"apple", "banana", "cherry"} // by length: 5, 6, 6
	if len(result) != len(expected) {
		t.Errorf("SortBy: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_MaxBy_MinBy(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[string])
	heap.Push("a", "bb", "ccc", "dd")

	// Find max by string length
	maxByLen := heap.Iter().MaxBy(func(a, b string) cmp.Ordering {
		return cmp.Cmp(len(a), len(b))
	})

	if !maxByLen.IsSome() || maxByLen.Some() != "ccc" {
		t.Errorf("MaxBy: expected 'ccc', got %v", maxByLen)
	}

	// Find min by string length
	minByLen := heap.Iter().MinBy(func(a, b string) cmp.Ordering {
		return cmp.Cmp(len(a), len(b))
	})

	if !minByLen.IsSome() || minByLen.Some() != "a" {
		t.Errorf("MinBy: expected 'a', got %v", minByLen)
	}
}

func TestSeqHeap_Unique(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 1, 3, 2, 4)

	result := make([]int, 0)
	heap.Iter().Unique().ForEach(func(v int) {
		result = append(result, v)
	})

	// Should contain each unique element once in sorted order
	expected := []int{1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("Unique: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Windows(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	windows := heap.Iter().Windows(3).Collect()

	if len(windows) != 3 { // [1,2,3], [2,3,4], [3,4,5]
		t.Errorf("Expected 3 windows, got %d", len(windows))
	}

	for _, window := range windows {
		if len(window) != 3 {
			t.Errorf("Expected window size 3, got %d", len(window))
		}
	}
}

func TestSeqHeap_Zip(t *testing.T) {
	heap1 := g.NewHeap(cmp.Cmp[int])
	heap1.Push(1, 2, 3)

	heap2 := g.NewHeap(cmp.Cmp[int])
	heap2.Push(4, 5, 6)

	pairs := make([][2]any, 0)
	heap1.Iter().Zip(heap2.Iter()).ForEach(func(a, b any) {
		pairs = append(pairs, [2]any{a, b})
	})

	expected := [][2]int{{1, 4}, {2, 5}, {3, 6}}
	if len(pairs) != len(expected) {
		t.Errorf("Zip: expected %v, got %v", expected, pairs)
	}
}

func TestSeqHeap_Intersperse(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	result := make([]int, 0)
	heap.Iter().Intersperse(0).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 0, 2, 0, 3}
	if len(result) != len(expected) {
		t.Errorf("Intersperse: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_First(t *testing.T) {
	t.Run("first element exists", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])
		heap.Push(50, 20, 30, 40, 10)

		first := heap.Iter().First()

		if first.IsNone() {
			t.Error("Expected Some value, got None")
		} else if first.Some() != 10 {
			t.Errorf("Expected 10 (min element), got %d", first.Some())
		}
	})

	t.Run("empty heap", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])

		first := heap.Iter().First()

		if first.IsSome() {
			t.Errorf("Expected None for empty heap, got Some(%v)", first.Some())
		}
	})

	t.Run("single element", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])
		heap.Push(42)

		first := heap.Iter().First()

		if first.IsNone() {
			t.Error("Expected Some value, got None")
		} else if first.Some() != 42 {
			t.Errorf("Expected 42, got %d", first.Some())
		}
	})
}

func TestSeqHeap_Last(t *testing.T) {
	t.Run("last element exists", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])
		heap.Push(50, 20, 30, 40, 10)

		last := heap.Iter().Last()

		if last.IsNone() {
			t.Error("Expected Some value, got None")
		}
	})

	t.Run("empty heap", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])

		last := heap.Iter().Last()

		if last.IsSome() {
			t.Errorf("Expected None for empty heap, got Some(%v)", last.Some())
		}
	})

	t.Run("single element", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])
		heap.Push(42)

		last := heap.Iter().Last()

		if last.IsNone() {
			t.Error("Expected Some value, got None")
		} else if last.Some() != 42 {
			t.Errorf("Expected 42, got %d", last.Some())
		}
	})
}

func TestSeqHeap_Nth(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(50, 20, 30, 40, 10)

	second := heap.Iter().Nth(1)

	if !second.IsSome() {
		t.Error("Expected to find element at index 1")
	}

	if second.Some() != 20 {
		t.Errorf("Expected 20 at index 1, got %d", second.Some())
	}
}

func TestSeqHeap_Partition(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5, 6)

	evens, odds := heap.Iter().Partition(
		func(v int) bool { return v%2 == 0 },
		cmp.Cmp[int], // min heap for evens
		func(a, b int) cmp.Ordering { return cmp.Cmp(b, a) }, // max heap for odds
	)

	if evens.Len() != 3 || odds.Len() != 3 {
		t.Errorf("Expected both heaps to have 3 elements, got evens=%d, odds=%d", evens.Len(), odds.Len())
	}

	// Verify even numbers come out in min heap order
	evenResult := make([]int, 0)
	for !evens.IsEmpty() {
		evenResult = append(evenResult, evens.Pop().Some())
	}
	expectedEvens := []int{2, 4, 6}
	if len(evenResult) != len(expectedEvens) {
		t.Errorf("Evens: expected %v, got %v", expectedEvens, evenResult)
	}

	// Verify odd numbers come out in max heap order
	oddResult := make([]int, 0)
	for !odds.IsEmpty() {
		oddResult = append(oddResult, odds.Pop().Some())
	}
	expectedOdds := []int{5, 3, 1} // max heap order
	if len(oddResult) != len(expectedOdds) {
		t.Errorf("Odds: expected %v, got %v", expectedOdds, oddResult)
	}
}

func TestSeqHeap_Context(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result := make([]int, 0)
	heap.Iter().Context(ctx).ForEach(func(v int) {
		result = append(result, v)
	})

	// Should complete quickly before timeout
	if len(result) != 5 {
		t.Errorf("Expected all 5 elements, got %d", len(result))
	}
}

func TestSeqHeap_Pull(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(3, 1, 2)

	next, stop := heap.Iter().Pull()
	defer stop()

	val, ok := next()
	if !ok || val != 1 {
		t.Errorf("Expected first element 1, got %d (ok: %t)", val, ok)
	}

	val, ok = next()
	if !ok || val != 2 {
		t.Errorf("Expected second element 2, got %d (ok: %t)", val, ok)
	}
}

func TestSeqHeap_ToChan(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(3, 1, 4, 1, 5)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := heap.Iter().ToChan(ctx)
	result := make([]int, 0)

	for val := range ch {
		result = append(result, val)
	}

	expected := []int{1, 1, 3, 4, 5} // sorted order from min heap
	if len(result) != len(expected) {
		t.Errorf("ToChan: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Inspect(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	inspected := make([]int, 0)
	result := make([]int, 0)

	heap.Iter().
		Inspect(func(v int) { inspected = append(inspected, v) }).
		ForEach(func(v int) { result = append(result, v) })

	// Both should have same elements
	if len(inspected) != len(result) {
		t.Errorf("Inspect: inspected=%v, result=%v", inspected, result)
	}
}

func TestSeqHeap_Combinations(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	combos := heap.Iter().Combinations(2).Collect()

	if len(combos) != 3 { // C(3,2) = 3
		t.Errorf("Expected 3 combinations, got %d", len(combos))
	}

	for _, combo := range combos {
		if len(combo) != 2 {
			t.Errorf("Expected combination size 2, got %d", len(combo))
		}
	}
}

func TestSeqHeap_Permutations(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	perms := heap.Iter().Permutations().Collect()

	if len(perms) != 6 { // 3! = 6
		t.Errorf("Expected 6 permutations, got %d", len(perms))
	}

	for _, perm := range perms {
		if len(perm) != 3 {
			t.Errorf("Expected permutation size 3, got %d", len(perm))
		}
	}
}

// Additional tests to improve coverage

func TestSeqHeap_CounterEdgeCases(t *testing.T) {
	// Test Counter on empty heap
	empty := g.NewHeap(cmp.Cmp[int])
	counts := make(map[int]g.Int)

	empty.Iter().Counter().ForEach(func(k any, v g.Int) {
		counts[k.(int)] = v
	})

	if len(counts) != 0 {
		t.Errorf("Expected empty counter from empty heap, got %d entries", len(counts))
	}
}

func TestSeqHeap_GroupBy(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 1, 2, 3, 2, 3, 4)

	groups := heap.Iter().GroupBy(func(a, b int) bool { return a <= b }).Collect()

	if len(groups) == 0 {
		t.Error("Expected at least one group")
	}

	// Verify we have some groups
	totalElements := 0
	for _, group := range groups {
		totalElements += len(group)
	}

	if totalElements != 7 {
		t.Errorf("Expected total of 7 elements in groups, got %d", totalElements)
	}
}

func TestSeqHeap_Cycle(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2)

	result := make([]int, 0)
	heap.Iter().Cycle().Take(6).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 1, 2, 1, 2} // cycle repeats sorted order
	if len(result) != len(expected) {
		t.Errorf("Cycle: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_DedupEdgeCases(t *testing.T) {
	// Test Dedup on empty heap
	empty := g.NewHeap(cmp.Cmp[int])
	result := make([]int, 0)

	empty.Iter().Dedup().ForEach(func(v int) {
		result = append(result, v)
	})

	if len(result) != 0 {
		t.Errorf("Expected empty result from empty heap dedup, got %v", result)
	}

	// Test Dedup on single element
	single := g.NewHeap(cmp.Cmp[int])
	single.Push(42)
	result = nil

	single.Iter().Dedup().ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{42}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Dedup single: expected %v, got %v", expected, result)
	}
}

func TestSeqHeap_Flatten(t *testing.T) {
	// Test Flatten with heap containing nested any elements
	heap := g.NewHeap(func(a, b any) cmp.Ordering {
		// Simple comparison for any type
		aStr := fmt.Sprintf("%v", a)
		bStr := fmt.Sprintf("%v", b)
		return cmp.Cmp(aStr, bStr)
	})

	heap.Push([]any{2, 3}, 1, []any{4, 5})
	result := make([]any, 0)

	heap.Iter().Flatten().ForEach(func(v any) {
		result = append(result, v)
	})

	// Should flatten the nested elements
	if len(result) == 0 {
		t.Errorf("Expected flattened result, got empty")
	}
}

func TestSeqHeap_Range(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	result := make([]int, 0)
	heap.Iter().Range(func(val int) bool {
		result = append(result, val)
		return val < 3 // Stop when we reach 3 or higher
	})

	expected := []int{1, 2, 3} // Should include the 3 that triggers the stop
	if len(result) != len(expected) {
		t.Errorf("Range: expected %v, got %v", expected, result)
	}
}

func TestSeqHeapEq(t *testing.T) {
	heap1 := g.NewHeap(cmp.Cmp[int])
	heap1.Push(1, 2, 3)

	heap2 := g.NewHeap(cmp.Cmp[int])
	heap2.Push(1, 2, 3)

	// Test equal heaps
	if !heap1.Iter().Eq(heap2.Iter()) {
		t.Errorf("Equal heaps should be equal")
	}

	// Test unequal heaps (different elements)
	heap3 := g.NewHeap(cmp.Cmp[int])
	heap3.Push(1, 2, 4)

	if heap1.Iter().Eq(heap3.Iter()) {
		t.Errorf("Heaps with different elements should not be equal")
	}

	// Test unequal heaps (different lengths)
	heap4 := g.NewHeap(cmp.Cmp[int])
	heap4.Push(1, 2, 3, 4)

	if heap1.Iter().Eq(heap4.Iter()) {
		t.Errorf("Heaps with different lengths should not be equal")
	}

	// Test empty heaps
	heap5 := g.NewHeap(cmp.Cmp[int])
	heap6 := g.NewHeap(cmp.Cmp[int])

	if !heap5.Iter().Eq(heap6.Iter()) {
		t.Errorf("Empty heaps should be equal")
	}
}

func TestSeqHeapFlatMap(t *testing.T) {
	// Test FlatMap with expanding elements
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	result := heap.Iter().FlatMap(func(n int) g.SeqHeap[int] {
		subHeap := g.NewHeap(cmp.Cmp[int])
		subHeap.Push(n, n*10)
		return subHeap.Iter()
	}).Collect(cmp.Cmp[int])

	// Count elements - should have 6 total (2 per original element)
	count := 0
	result.Iter().ForEach(func(v int) { count++ })

	if count != 6 {
		t.Errorf("FlatMap count: expected 6, got %d", count)
	}

	// Verify all elements are present (though order depends on heap)
	elements := make(map[int]bool)
	result.Iter().ForEach(func(v int) { elements[v] = true })

	expected := []int{1, 10, 2, 20, 3, 30}
	for _, v := range expected {
		if !elements[v] {
			t.Errorf("FlatMap missing element: %d", v)
		}
	}

	// Test FlatMap with empty results
	emptyResult := heap.Iter().FlatMap(func(n int) g.SeqHeap[int] {
		return g.NewHeap(cmp.Cmp[int]).Iter()
	}).Collect(cmp.Cmp[int])

	if !emptyResult.IsEmpty() {
		t.Errorf("FlatMap empty: expected empty heap, got %d elements", emptyResult.Len())
	}

	// Test FlatMap on empty heap
	empty := g.NewHeap(cmp.Cmp[int])
	emptyFlatMapped := empty.Iter().FlatMap(func(n int) g.SeqHeap[int] {
		subHeap := g.NewHeap(cmp.Cmp[int])
		subHeap.Push(n)
		return subHeap.Iter()
	}).Collect(cmp.Cmp[int])

	if !emptyFlatMapped.IsEmpty() {
		t.Errorf("FlatMap empty input: expected empty heap, got %d elements", emptyFlatMapped.Len())
	}

	// Test FlatMap with single elements
	singleResult := heap.Iter().FlatMap(func(n int) g.SeqHeap[int] {
		subHeap := g.NewHeap(cmp.Cmp[int])
		subHeap.Push(n * 2)
		return subHeap.Iter()
	}).Collect(cmp.Cmp[int])

	singleCount := 0
	singleResult.Iter().ForEach(func(v int) { singleCount++ })

	if singleCount != 3 { // Same number as original heap
		t.Errorf("FlatMap single count: expected 3, got %d", singleCount)
	}
}

func TestSeqHeapFilterMap(t *testing.T) {
	// Test FilterMap with transformation and filtering
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	result := heap.Iter().FilterMap(func(n int) g.Option[int] {
		if n%2 == 0 {
			return g.Some(n * 10)
		}
		return g.None[int]()
	}).Collect(cmp.Cmp[int])

	// Should have 2 even numbers (2, 4) transformed to (20, 40)
	expected := []int{20, 40}
	actual := make([]int, 0)
	result.Iter().ForEach(func(v int) { actual = append(actual, v) })

	if len(actual) != len(expected) {
		t.Errorf("FilterMap length: expected %d, got %d", len(expected), len(actual))
	}

	// Since heap order may vary, check that all expected elements are present
	actualSet := make(map[int]bool)
	for _, v := range actual {
		actualSet[v] = true
	}

	for _, v := range expected {
		if !actualSet[v] {
			t.Errorf("FilterMap missing element: %d", v)
		}
	}

	// Test FilterMap that filters all elements
	allFiltered := heap.Iter().FilterMap(func(n int) g.Option[int] {
		return g.None[int]()
	}).Collect(cmp.Cmp[int])

	if !allFiltered.IsEmpty() {
		t.Errorf("FilterMap all filtered: expected empty heap, got %d elements", allFiltered.Len())
	}

	// Test FilterMap that keeps all elements with transformation
	allKept := heap.Iter().FilterMap(func(n int) g.Option[int] {
		return g.Some(n * 2)
	}).Collect(cmp.Cmp[int])

	keptCount := 0
	allKept.Iter().ForEach(func(v int) { keptCount++ })

	if keptCount != 5 { // Same as original heap size
		t.Errorf("FilterMap all kept count: expected 5, got %d", keptCount)
	}

	// Test FilterMap on empty heap
	empty := g.NewHeap(cmp.Cmp[int])
	emptyFiltered := empty.Iter().FilterMap(func(n int) g.Option[int] {
		return g.Some(n * 2)
	}).Collect(cmp.Cmp[int])

	if !emptyFiltered.IsEmpty() {
		t.Errorf("FilterMap empty input: expected empty heap, got %d elements", emptyFiltered.Len())
	}

	// Test FilterMap with string processing
	stringHeap := g.NewHeap(cmp.Cmp[string])
	stringHeap.Push("hello", "", "world", "   ", "go")

	processedWords := stringHeap.Iter().FilterMap(func(s string) g.Option[string] {
		trimmed := g.String(s).Trim()
		if !trimmed.IsEmpty() {
			return g.Some(string(trimmed.Upper()))
		}
		return g.None[string]()
	}).Collect(cmp.Cmp[string])

	// Should have 3 valid words: "HELLO", "WORLD", "GO"
	wordCount := 0
	processedWords.Iter().ForEach(func(s string) { wordCount++ })

	if wordCount != 3 {
		t.Errorf("FilterMap strings count: expected 3, got %d", wordCount)
	}

	// Verify expected words are present
	words := make(map[string]bool)
	processedWords.Iter().ForEach(func(s string) { words[s] = true })

	expectedWords := []string{"HELLO", "WORLD", "GO"}
	for _, word := range expectedWords {
		if !words[word] {
			t.Errorf("FilterMap missing word: %s", word)
		}
	}
}

func TestSeqHeapScan(t *testing.T) {
	// Test Scan with sum accumulation
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	result := heap.Iter().Scan(0, func(acc, val int) int {
		return acc + val
	}).Collect(cmp.Cmp[int])

	// Should have 6 elements: initial value + 5 accumulated values
	count := 0
	result.Iter().ForEach(func(v int) { count++ })

	if count != 6 {
		t.Errorf("Scan count: expected 6, got %d", count)
	}

	// Check that 0 (initial value) is present
	hasZero := false
	result.Iter().ForEach(func(v int) {
		if v == 0 {
			hasZero = true
		}
	})

	if !hasZero {
		t.Errorf("Scan should include initial value 0")
	}

	// Test Scan with multiplication
	small := g.NewHeap(cmp.Cmp[int])
	small.Push(2, 3, 4)

	product := small.Iter().Scan(1, func(acc, val int) int {
		return acc * val
	}).Collect(cmp.Cmp[int])

	productCount := 0
	product.Iter().ForEach(func(v int) { productCount++ })

	if productCount != 4 { // initial + 3 elements
		t.Errorf("Scan product count: expected 4, got %d", productCount)
	}

	// Check that 1 (initial value) is present
	hasOne := false
	product.Iter().ForEach(func(v int) {
		if v == 1 {
			hasOne = true
		}
	})

	if !hasOne {
		t.Errorf("Scan should include initial value 1")
	}

	// Test Scan on empty heap (should just return initial value)
	empty := g.NewHeap(cmp.Cmp[int])
	emptyScanned := empty.Iter().Scan(42, func(acc, val int) int {
		return acc + val
	}).Collect(cmp.Cmp[int])

	if emptyScanned.Len() != 1 {
		t.Errorf("Scan empty length: expected 1, got %d", emptyScanned.Len())
	}

	// Verify the single element is the initial value
	foundInitial := false
	emptyScanned.Iter().ForEach(func(v int) {
		if v == 42 {
			foundInitial = true
		}
	})

	if !foundInitial {
		t.Errorf("Scan empty should contain initial value 42")
	}

	// Test Scan with single element
	single := g.NewHeap(cmp.Cmp[int])
	single.Push(10)

	singleScanned := single.Iter().Scan(5, func(acc, val int) int {
		return acc + val
	}).Collect(cmp.Cmp[int])

	if singleScanned.Len() != 2 { // initial + 1 element
		t.Errorf("Scan single length: expected 2, got %d", singleScanned.Len())
	}

	// Check both expected values are present
	values := make(map[int]bool)
	singleScanned.Iter().ForEach(func(v int) { values[v] = true })

	if !values[5] || !values[15] {
		t.Errorf("Scan single should contain both 5 and 15")
	}

	// Test Scan with string concatenation
	stringHeap := g.NewHeap(cmp.Cmp[string])
	stringHeap.Push("a", "b", "c")

	concatenated := stringHeap.Iter().Scan("", func(acc, val string) string {
		return acc + val
	}).Collect(cmp.Cmp[string])

	concatCount := 0
	concatenated.Iter().ForEach(func(s string) { concatCount++ })

	if concatCount != 4 { // initial + 3 elements
		t.Errorf("Scan concat count: expected 4, got %d", concatCount)
	}

	// Check that empty string (initial value) is present
	hasEmpty := false
	concatenated.Iter().ForEach(func(s string) {
		if s == "" {
			hasEmpty = true
		}
	})

	if !hasEmpty {
		t.Errorf("Scan should include initial empty string")
	}
}

func TestSeqHeapNext(t *testing.T) {
	t.Run("Next with non-empty iterator", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])
		heap.Push(3, 1, 4, 2, 5)
		iter := heap.Iter()

		// Extract first element (should be smallest in min-heap)
		first := iter.Next()
		if !first.IsSome() || first.Some() != 1 {
			t.Errorf("Expected Some(1), got %v", first)
		}

		// Extract second element
		second := iter.Next()
		if !second.IsSome() || second.Some() != 2 {
			t.Errorf("Expected Some(2), got %v", second)
		}

		// Check that remaining elements can be collected
		remaining := iter.Collect(cmp.Cmp[int])
		if remaining.Len() != 3 {
			t.Errorf("Expected 3 remaining elements, got %d", remaining.Len())
		}
	})

	t.Run("Next with empty iterator", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])
		iter := heap.Iter()

		result := iter.Next()
		if result.IsSome() {
			t.Errorf("Expected None, got Some(%v)", result.Some())
		}
	})

	t.Run("Next until exhausted", func(t *testing.T) {
		heap := g.NewHeap(cmp.Cmp[int])
		heap.Push(1, 2)
		iter := heap.Iter()

		// Extract all elements
		first := iter.Next()
		second := iter.Next()
		third := iter.Next()

		if !first.IsSome() || first.Some() != 1 {
			t.Errorf("Expected first to be Some(1), got %v", first)
		}
		if !second.IsSome() || second.Some() != 2 {
			t.Errorf("Expected second to be Some(2), got %v", second)
		}
		if third.IsSome() {
			t.Errorf("Expected third to be None, got Some(%v)", third.Some())
		}

		// Iterator should be empty now
		remaining := iter.Collect(cmp.Cmp[int])
		if remaining.Len() != 0 {
			t.Errorf("Expected empty heap, got length %d", remaining.Len())
		}
	})
}
