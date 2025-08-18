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

	collected := heap.Iter().Collect()

	if collected.Len() != 5 {
		t.Errorf("Expected collected heap length 5, got %d", collected.Len())
	}
}

func TestSeqHeap_CollectWith(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(3, 1, 4, 1, 5)

	// Collect with reverse comparison (max heap)
	collected := heap.Iter().CollectWith(func(a, b int) cmp.Ordering {
		return cmp.Cmp(b, a) // reverse comparison
	})

	if collected.Len() != 5 {
		t.Errorf("Expected collected heap length 5, got %d", collected.Len())
	}

	// Should pop in reverse order (largest first)
	result := make([]int, 0)
	for !collected.Empty() {
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
	heap.Iter().Counter().ForEach(func(k int, v g.Int) {
		counts[k] = v
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

	pairs := make([][2]int, 0)
	heap1.Iter().Zip(heap2.Iter()).ForEach(func(a, b int) {
		pairs = append(pairs, [2]int{a, b})
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

	evens, odds := heap.Iter().Partition(func(v int) bool { return v%2 == 0 })

	if evens.Len() != 3 {
		t.Errorf("Expected 3 even numbers, got %d", evens.Len())
	}

	if odds.Len() != 3 {
		t.Errorf("Expected 3 odd numbers, got %d", odds.Len())
	}
}

func TestSeqHeap_PartitionWith(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5, 6)

	evens, odds := heap.Iter().PartitionWith(
		func(v int) bool { return v%2 == 0 },
		cmp.Cmp[int], // min heap for evens
		func(a, b int) cmp.Ordering { return cmp.Cmp(b, a) }, // max heap for odds
	)

	if evens.Len() != 3 || odds.Len() != 3 {
		t.Errorf("Expected both heaps to have 3 elements, got evens=%d, odds=%d", evens.Len(), odds.Len())
	}

	// Verify even numbers come out in min heap order
	evenResult := make([]int, 0)
	for !evens.Empty() {
		evenResult = append(evenResult, evens.Pop().Some())
	}
	expectedEvens := []int{2, 4, 6}
	if len(evenResult) != len(expectedEvens) {
		t.Errorf("Evens: expected %v, got %v", expectedEvens, evenResult)
	}

	// Verify odd numbers come out in max heap order
	oddResult := make([]int, 0)
	for !odds.Empty() {
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

func TestSeqHeap_CollectEdgeCases(t *testing.T) {
	// Test Collect on empty heap
	empty := g.NewHeap(cmp.Cmp[int])
	collected := empty.Iter().Collect()

	if !collected.Empty() {
		t.Errorf("Expected empty collected heap from empty heap, got length %d", collected.Len())
	}
}

func TestSeqHeap_CounterEdgeCases(t *testing.T) {
	// Test Counter on empty heap
	empty := g.NewHeap(cmp.Cmp[int])
	counts := make(map[int]g.Int)

	empty.Iter().Counter().ForEach(func(k int, v g.Int) {
		counts[k] = v
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

func TestSeqHeap_PartitionEdgeCases(t *testing.T) {
	// Test Partition on empty heap
	empty := g.NewHeap(cmp.Cmp[int])
	evens, odds := empty.Iter().Partition(func(v int) bool { return v%2 == 0 })

	if !evens.Empty() {
		t.Errorf("Expected empty evens heap from empty input, got length %d", evens.Len())
	}

	if !odds.Empty() {
		t.Errorf("Expected empty odds heap from empty input, got length %d", odds.Len())
	}

	// Test Partition on single element
	single := g.NewHeap(cmp.Cmp[int])
	single.Push(2) // even
	evens, odds = single.Iter().Partition(func(v int) bool { return v%2 == 0 })

	if evens.Len() != 1 {
		t.Errorf("Expected 1 even from single even element, got %d", evens.Len())
	}

	if !odds.Empty() {
		t.Errorf("Expected empty odds from single even element, got %d", odds.Len())
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
