package g

import (
	"context"
	"testing"
	"time"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestSeqDeque_All(t *testing.T) {
	deque := g.DequeOf(2, 4, 6, 8)
	allEven := deque.Iter().All(func(v int) bool { return v%2 == 0 })
	if !allEven {
		t.Error("Expected all elements to be even")
	}

	deque2 := g.DequeOf(1, 2, 3, 4)
	allEven2 := deque2.Iter().All(func(v int) bool { return v%2 == 0 })
	if allEven2 {
		t.Error("Expected not all elements to be even")
	}
}

func TestSeqDeque_Any(t *testing.T) {
	deque := g.DequeOf(1, 3, 5, 7)
	anyEven := deque.Iter().Any(func(v int) bool { return v%2 == 0 })
	if anyEven {
		t.Error("Expected no even elements")
	}

	deque2 := g.DequeOf(1, 2, 3, 5)
	anyEven2 := deque2.Iter().Any(func(v int) bool { return v%2 == 0 })
	if !anyEven2 {
		t.Error("Expected at least one even element")
	}
}

func TestSeqDeque_Chain(t *testing.T) {
	deque1 := g.DequeOf(1, 2, 3)
	deque2 := g.DequeOf(4, 5, 6)

	result := make([]int, 0)
	deque1.Iter().Chain(deque2.Iter()).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 3, 4, 5, 6}
	if len(result) != len(expected) {
		t.Errorf("Chain: expected %v, got %v", expected, result)
	}

	for i, v := range expected {
		if i < len(result) && result[i] != v {
			t.Errorf("Chain at index %d: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestSeqDeque_Chunks(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5, 6)
	chunks := deque.Iter().Chunks(2).Collect()

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

func TestSeqDeque_Collect(t *testing.T) {
	deque := g.DequeOf(3, 1, 4, 1, 5)
	collected := deque.Iter().Collect()

	if collected.Len() != 5 {
		t.Errorf("Expected collected deque length 5, got %d", collected.Len())
	}

	// Verify elements using iterator
	expected := []int{3, 1, 4, 1, 5}
	actual := make([]int, 0)
	deque.Iter().ForEach(func(v int) {
		actual = append(actual, v)
	})

	for i, v := range expected {
		if i < len(actual) && actual[i] != v {
			t.Errorf("Collect at index %d: expected %d, got %d", i, v, actual[i])
		}
	}
}

func TestSeqDeque_Count(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	count := deque.Iter().Count()
	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}

func TestSeqDeque_Counter(t *testing.T) {
	deque := g.DequeOf(1, 2, 1, 3, 2, 1)
	counts := make(map[int]g.Int)

	deque.Iter().Counter().ForEach(func(k any, v g.Int) {
		counts[k.(int)] = v
	})

	expected := map[int]g.Int{1: 3, 2: 2, 3: 1}
	for k, v := range expected {
		if counts[k] != v {
			t.Errorf("Counter for %d: expected %d, got %d", k, v, counts[k])
		}
	}
}

func TestSeqDeque_Dedup(t *testing.T) {
	deque := g.DequeOf(1, 1, 2, 2, 3, 3, 3)
	result := make([]int, 0)

	deque.Iter().Dedup().ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("Dedup: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Enumerate(t *testing.T) {
	deque := g.DequeOf("a", "b", "c")
	pairs := make(map[g.Int]string)

	deque.Iter().Enumerate().ForEach(func(i g.Int, v string) {
		pairs[i] = v
	})

	expected := map[g.Int]string{0: "a", 1: "b", 2: "c"}
	for k, v := range expected {
		if pairs[k] != v {
			t.Errorf("Enumerate at %d: expected %s, got %s", k, v, pairs[k])
		}
	}
}

func TestSeqDeque_Filter(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5, 6)
	result := make([]int, 0)

	deque.Iter().Filter(func(v int) bool { return v%2 == 0 }).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{2, 4, 6}
	if len(result) != len(expected) {
		t.Errorf("Filter: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Exclude(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	result := make([]int, 0)

	deque.Iter().Exclude(func(v int) bool { return v%2 == 0 }).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 3, 5}
	if len(result) != len(expected) {
		t.Errorf("Exclude: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Find(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	found := deque.Iter().Find(func(v int) bool { return v > 3 })

	if !found.IsSome() {
		t.Error("Expected to find element > 3")
	}

	if found.Some() != 4 {
		t.Errorf("Expected to find 4, got %d", found.Some())
	}
}

func TestSeqDeque_Fold(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	sum := deque.Iter().Fold(0, func(acc, val int) int {
		return acc + val
	})

	if sum != 15 {
		t.Errorf("Expected sum 15, got %d", sum)
	}
}

func TestSeqDeque_Reduce(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	product := deque.Iter().Reduce(func(a, b int) int {
		return a * b
	})

	if !product.IsSome() {
		t.Error("Expected reduce to return some value")
	}

	if product.Some() != 120 {
		t.Errorf("Expected product 120, got %d", product.Some())
	}
}

func TestSeqDeque_ForEach(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	sum := 0

	deque.Iter().ForEach(func(v int) {
		sum += v
	})

	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
}

func TestSeqDeque_Map(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	result := make([]int, 0)

	deque.Iter().Map(func(v int) int {
		return v * 2
	}).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{2, 4, 6}
	if len(result) != len(expected) {
		t.Errorf("Map: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Skip(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	result := make([]int, 0)

	deque.Iter().Skip(2).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Skip: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Take(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	result := make([]int, 0)

	deque.Iter().Take(3).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("Take: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_StepBy(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5, 6)
	result := make([]int, 0)

	deque.Iter().StepBy(2).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 3, 5}
	if len(result) != len(expected) {
		t.Errorf("StepBy: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Unique(t *testing.T) {
	deque := g.DequeOf(1, 2, 1, 3, 2, 4)
	result := make([]int, 0)

	deque.Iter().Unique().ForEach(func(v int) {
		result = append(result, v)
	})

	// Should contain each unique element once
	expected := []int{1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("Unique: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Windows(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	windows := deque.Iter().Windows(3).Collect()

	if len(windows) != 3 { // [1,2,3], [2,3,4], [3,4,5]
		t.Errorf("Expected 3 windows, got %d", len(windows))
	}

	for _, window := range windows {
		if len(window) != 3 {
			t.Errorf("Expected window size 3, got %d", len(window))
		}
	}
}

func TestSeqDeque_Zip(t *testing.T) {
	deque1 := g.DequeOf(1, 2, 3)
	deque2 := g.DequeOf(4, 5, 6)

	pairs := make([][2]any, 0)
	deque1.Iter().Zip(deque2.Iter()).ForEach(func(a, b any) {
		pairs = append(pairs, [2]any{a, b})
	})

	expected := [][2]int{{1, 4}, {2, 5}, {3, 6}}
	if len(pairs) != len(expected) {
		t.Errorf("Zip: expected %v, got %v", expected, pairs)
	}
}

func TestSeqDeque_Intersperse(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	result := make([]int, 0)

	deque.Iter().Intersperse(0).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 0, 2, 0, 3}
	if len(result) != len(expected) {
		t.Errorf("Intersperse: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_First(t *testing.T) {
	t.Run("first element exists", func(t *testing.T) {
		deque := g.DequeOf(10, 20, 30, 40, 50)

		first := deque.Iter().First()

		if first.IsNone() {
			t.Error("Expected Some value, got None")
		} else if first.Some() != 10 {
			t.Errorf("Expected 10, got %d", first.Some())
		}
	})

	t.Run("empty deque", func(t *testing.T) {
		deque := g.DequeOf[int]()

		first := deque.Iter().First()

		if first.IsSome() {
			t.Errorf("Expected None for empty deque, got Some(%v)", first.Some())
		}
	})

	t.Run("single element", func(t *testing.T) {
		deque := g.DequeOf(42)

		first := deque.Iter().First()

		if first.IsNone() {
			t.Error("Expected Some value, got None")
		} else if first.Some() != 42 {
			t.Errorf("Expected 42, got %d", first.Some())
		}
	})
}

func TestSeqDeque_Last(t *testing.T) {
	t.Run("last element exists", func(t *testing.T) {
		deque := g.DequeOf(10, 20, 30, 40, 50)

		last := deque.Iter().Last()

		if last.IsNone() {
			t.Error("Expected Some value, got None")
		} else if last.Some() != 50 {
			t.Errorf("Expected 50, got %d", last.Some())
		}
	})

	t.Run("empty deque", func(t *testing.T) {
		deque := g.DequeOf[int]()

		last := deque.Iter().Last()

		if last.IsSome() {
			t.Errorf("Expected None for empty deque, got Some(%v)", last.Some())
		}
	})

	t.Run("single element", func(t *testing.T) {
		deque := g.DequeOf(42)

		last := deque.Iter().Last()

		if last.IsNone() {
			t.Error("Expected Some value, got None")
		} else if last.Some() != 42 {
			t.Errorf("Expected 42, got %d", last.Some())
		}
	})
}

func TestSeqDeque_Nth(t *testing.T) {
	deque := g.DequeOf(10, 20, 30, 40, 50)
	second := deque.Iter().Nth(1)

	if !second.IsSome() {
		t.Error("Expected to find element at index 1")
	}

	if second.Some() != 20 {
		t.Errorf("Expected 20 at index 1, got %d", second.Some())
	}
}

func TestSeqDeque_Partition(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5, 6)
	evens, odds := deque.Iter().Partition(func(v int) bool { return v%2 == 0 })

	if evens.Len() != 3 {
		t.Errorf("Expected 3 even numbers, got %d", evens.Len())
	}

	if odds.Len() != 3 {
		t.Errorf("Expected 3 odd numbers, got %d", odds.Len())
	}
}

func TestSeqDeque_Context(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result := make([]int, 0)
	deque.Iter().Context(ctx).ForEach(func(v int) {
		result = append(result, v)
	})

	// Should complete quickly before timeout
	if len(result) != 5 {
		t.Errorf("Expected all 5 elements, got %d", len(result))
	}
}

func TestSeqDeque_Pull(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	next, stop := deque.Iter().Pull()
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

func TestSeqDeque_ToChan(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := deque.Iter().Chan(ctx)
	result := make([]int, 0)

	for val := range ch {
		result = append(result, val)
	}

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("ToChan: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Combinations(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	combos := deque.Iter().Combinations(2).Collect()

	if len(combos) != 3 { // C(3,2) = 3
		t.Errorf("Expected 3 combinations, got %d", len(combos))
	}

	for _, combo := range combos {
		if len(combo) != 2 {
			t.Errorf("Expected combination size 2, got %d", len(combo))
		}
	}
}

func TestSeqDeque_Cycle(t *testing.T) {
	deque := g.DequeOf(1, 2)
	result := make([]int, 0)

	deque.Iter().Cycle().Take(6).ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 1, 2, 1, 2} // cycle repeats
	if len(result) != len(expected) {
		t.Errorf("Cycle: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_Flatten(t *testing.T) {
	// Test with deque containing nested any elements that can be flattened
	deque := g.DequeOf[any](1, []any{2, 3}, 4, []any{5, 6})
	result := make([]any, 0)

	deque.Iter().Flatten().ForEach(func(v any) {
		result = append(result, v)
	})

	expected := []any{1, 2, 3, 4, 5, 6}
	if len(result) != len(expected) {
		t.Errorf("Flatten: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_GroupBy(t *testing.T) {
	deque := g.DequeOf(1, 1, 2, 3, 2, 3, 4)
	groups := deque.Iter().GroupBy(func(a, b int) bool { return a <= b }).Collect()

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

func TestSeqDeque_Inspect(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	inspected := make([]int, 0)
	result := make([]int, 0)

	deque.Iter().
		Inspect(func(v int) { inspected = append(inspected, v) }).
		ForEach(func(v int) { result = append(result, v) })

	// Both should have same elements
	if len(inspected) != len(result) {
		t.Errorf("Inspect: inspected=%v, result=%v", inspected, result)
	}
}

func TestSeqDeque_MaxBy(t *testing.T) {
	deque := g.DequeOf("a", "bb", "ccc", "dd")

	// Find max by string length
	maxByLen := deque.Iter().MaxBy(func(a, b string) cmp.Ordering {
		return cmp.Cmp(len(a), len(b))
	})

	if !maxByLen.IsSome() || maxByLen.Some() != "ccc" {
		t.Errorf("MaxBy: expected 'ccc', got %v", maxByLen)
	}
}

func TestSeqDeque_MinBy(t *testing.T) {
	deque := g.DequeOf("a", "bb", "ccc", "dd")

	// Find min by string length
	minByLen := deque.Iter().MinBy(func(a, b string) cmp.Ordering {
		return cmp.Cmp(len(a), len(b))
	})

	if !minByLen.IsSome() || minByLen.Some() != "a" {
		t.Errorf("MinBy: expected 'a', got %v", minByLen)
	}
}

func TestSeqDeque_Permutations(t *testing.T) {
	deque := g.DequeOf(1, 2, 3)
	perms := deque.Iter().Permutations().Collect()

	if len(perms) != 6 { // 3! = 6
		t.Errorf("Expected 6 permutations, got %d", len(perms))
	}

	for _, perm := range perms {
		if len(perm) != 3 {
			t.Errorf("Expected permutation size 3, got %d", len(perm))
		}
	}
}

func TestSeqDeque_Range(t *testing.T) {
	deque := g.DequeOf(1, 2, 3, 4, 5)
	result := make([]int, 0)

	deque.Iter().Range(func(val int) bool {
		result = append(result, val)
		return val < 3 // Stop when we reach 3 or higher
	})

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("Range: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_SortBy(t *testing.T) {
	deque := g.DequeOf("apple", "banana", "cherry")
	result := make([]string, 0)

	// Sort by string length (ascending)
	deque.Iter().SortBy(func(a, b string) cmp.Ordering {
		return cmp.Cmp(len(a), len(b))
	}).ForEach(func(v string) {
		result = append(result, v)
	})

	expected := []string{"apple", "banana", "cherry"} // by length: 5, 6, 6
	if len(result) != len(expected) {
		t.Errorf("SortBy: expected %v, got %v", expected, result)
	}
}

// Additional tests to improve coverage

func TestSeqDeque_CounterEdgeCases(t *testing.T) {
	// Test Counter on empty deque
	empty := g.NewDeque[int]()
	counts := make(map[int]g.Int)

	empty.Iter().Counter().ForEach(func(k any, v g.Int) {
		counts[k.(int)] = v
	})

	if len(counts) != 0 {
		t.Errorf("Expected empty counter from empty deque, got %d entries", len(counts))
	}
}

func TestSeqDeque_DedupEdgeCases(t *testing.T) {
	// Test Dedup on empty deque
	empty := g.NewDeque[int]()
	result := make([]int, 0)

	empty.Iter().Dedup().ForEach(func(v int) {
		result = append(result, v)
	})

	if len(result) != 0 {
		t.Errorf("Expected empty result from empty deque dedup, got %v", result)
	}

	// Test Dedup on single element
	single := g.DequeOf(42)
	result = nil

	single.Iter().Dedup().ForEach(func(v int) {
		result = append(result, v)
	})

	expected := []int{42}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Dedup single: expected %v, got %v", expected, result)
	}
}

func TestSeqDeque_FlattenEdgeCases(t *testing.T) {
	// Test Flatten on empty deque
	empty := g.NewDeque[any]()
	result := make([]any, 0)

	empty.Iter().Flatten().ForEach(func(v any) {
		result = append(result, v)
	})

	if len(result) != 0 {
		t.Errorf("Expected empty result from empty deque flatten, got %v", result)
	}

	// Test Flatten with non-slice elements
	simple := g.DequeOf[any](1, 2, 3)
	result = nil

	simple.Iter().Flatten().ForEach(func(v any) {
		result = append(result, v)
	})

	expected := []any{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("Flatten simple: expected %v, got %v", expected, result)
	}
}

func TestSeqDequeFlatMap(t *testing.T) {
	// Test FlatMap with expanding elements
	deque := g.DequeOf(1, 2, 3)
	result := deque.Iter().FlatMap(func(n int) g.SeqDeque[int] {
		return g.DequeOf(n, n*10).Iter()
	}).Collect()

	expected := []int{1, 10, 2, 20, 3, 30}
	actual := make([]int, 0)
	result.Iter().ForEach(func(v int) { actual = append(actual, v) })

	if len(actual) != len(expected) {
		t.Errorf("FlatMap length: expected %d, got %d", len(expected), len(actual))
	}

	for i, v := range expected {
		if i >= len(actual) || actual[i] != v {
			t.Errorf("FlatMap result[%d]: expected %d, got %v", i, v, actual)
			break
		}
	}

	// Test FlatMap with empty results
	emptyResult := deque.Iter().FlatMap(func(n int) g.SeqDeque[int] {
		return g.NewDeque[int]().Iter()
	}).Collect()

	if emptyResult.Len() != 0 {
		t.Errorf("FlatMap empty: expected 0 elements, got %d", emptyResult.Len())
	}

	// Test FlatMap with single elements
	singleResult := deque.Iter().FlatMap(func(n int) g.SeqDeque[int] {
		return g.DequeOf(n * 2).Iter()
	}).Collect()

	expectedSingle := []int{2, 4, 6}
	actualSingle := make([]int, 0)
	singleResult.Iter().ForEach(func(v int) { actualSingle = append(actualSingle, v) })

	for i, v := range expectedSingle {
		if i >= len(actualSingle) || actualSingle[i] != v {
			t.Errorf("FlatMap single[%d]: expected %d, got %v", i, v, actualSingle)
		}
	}

	// Test FlatMap on empty deque
	empty := g.NewDeque[int]()
	emptyFlatMapped := empty.Iter().FlatMap(func(n int) g.SeqDeque[int] {
		return g.DequeOf(n).Iter()
	}).Collect()

	if emptyFlatMapped.Len() != 0 {
		t.Errorf("FlatMap empty input: expected 0 elements, got %d", emptyFlatMapped.Len())
	}
}

func TestSeqDequeFilterMap(t *testing.T) {
	// Test FilterMap with transformation and filtering
	deque := g.DequeOf(1, 2, 3, 4, 5)
	result := deque.Iter().FilterMap(func(n int) g.Option[int] {
		if n%2 == 0 {
			return g.Some(n * 10)
		}
		return g.None[int]()
	}).Collect()

	expected := []int{20, 40}
	actual := make([]int, 0)
	result.Iter().ForEach(func(v int) { actual = append(actual, v) })

	if len(actual) != len(expected) {
		t.Errorf("FilterMap length: expected %d, got %d", len(expected), len(actual))
	}

	for i, v := range expected {
		if i >= len(actual) || actual[i] != v {
			t.Errorf("FilterMap result[%d]: expected %d, got %v", i, v, actual)
			break
		}
	}

	// Test FilterMap that filters all elements
	allFiltered := deque.Iter().FilterMap(func(n int) g.Option[int] {
		return g.None[int]()
	}).Collect()

	if allFiltered.Len() != 0 {
		t.Errorf("FilterMap all filtered: expected 0 elements, got %d", allFiltered.Len())
	}

	// Test FilterMap that keeps all elements with transformation
	allKept := deque.Iter().FilterMap(func(n int) g.Option[int] {
		return g.Some(n * 2)
	}).Collect()

	expectedAll := []int{2, 4, 6, 8, 10}
	actualAll := make([]int, 0)
	allKept.Iter().ForEach(func(v int) { actualAll = append(actualAll, v) })

	if len(actualAll) != len(expectedAll) {
		t.Errorf("FilterMap all kept length: expected %d, got %d", len(expectedAll), len(actualAll))
	}

	// Test FilterMap on empty deque
	empty := g.NewDeque[int]()
	emptyFiltered := empty.Iter().FilterMap(func(n int) g.Option[int] {
		return g.Some(n * 2)
	}).Collect()

	if emptyFiltered.Len() != 0 {
		t.Errorf("FilterMap empty input: expected 0 elements, got %d", emptyFiltered.Len())
	}

	// Test FilterMap with string processing
	words := g.DequeOf("hello", "", "world", "   ", "go")
	processedWords := words.Iter().FilterMap(func(s string) g.Option[string] {
		trimmed := g.String(s).Trim()
		if !trimmed.IsEmpty() {
			return g.Some(string(trimmed.Upper()))
		}
		return g.None[string]()
	}).Collect()

	expectedWords := []string{"HELLO", "WORLD", "GO"}
	actualWords := make([]string, 0)
	processedWords.Iter().ForEach(func(s string) { actualWords = append(actualWords, s) })

	if len(actualWords) != len(expectedWords) {
		t.Errorf("FilterMap strings length: expected %d, got %d", len(expectedWords), len(actualWords))
	}

	for i, expected := range expectedWords {
		if i >= len(actualWords) || actualWords[i] != expected {
			t.Errorf("FilterMap strings[%d]: expected %s, got %v", i, expected, actualWords)
		}
	}
}

func TestSeqDequeScan(t *testing.T) {
	// Test Scan with sum accumulation
	deque := g.DequeOf(1, 2, 3, 4, 5)
	result := deque.Iter().Scan(0, func(acc, val int) int {
		return acc + val
	}).Collect()

	expected := []int{0, 1, 3, 6, 10, 15}
	actual := make([]int, 0)
	result.Iter().ForEach(func(v int) { actual = append(actual, v) })

	if len(actual) != len(expected) {
		t.Errorf("Scan length: expected %d, got %d", len(expected), len(actual))
	}

	for i, v := range expected {
		if i >= len(actual) || actual[i] != v {
			t.Errorf("Scan result[%d]: expected %d, got %v", i, v, actual)
			break
		}
	}

	// Test Scan with multiplication
	small := g.DequeOf(2, 3, 4)
	product := small.Iter().Scan(1, func(acc, val int) int {
		return acc * val
	}).Collect()

	expectedProduct := []int{1, 2, 6, 24}
	actualProduct := make([]int, 0)
	product.Iter().ForEach(func(v int) { actualProduct = append(actualProduct, v) })

	if len(actualProduct) != len(expectedProduct) {
		t.Errorf("Scan product length: expected %d, got %d", len(expectedProduct), len(actualProduct))
	}

	for i, v := range expectedProduct {
		if i >= len(actualProduct) || actualProduct[i] != v {
			t.Errorf("Scan product[%d]: expected %d, got %v", i, v, actualProduct)
		}
	}

	// Test Scan on empty deque (should just return initial value)
	empty := g.NewDeque[int]()
	emptyScanned := empty.Iter().Scan(42, func(acc, val int) int {
		return acc + val
	}).Collect()

	if emptyScanned.Len() != 1 {
		t.Errorf("Scan empty length: expected 1, got %d", emptyScanned.Len())
	}

	if emptyScanned.Front().UnwrapOr(0) != 42 {
		t.Errorf("Scan empty value: expected 42, got %v", emptyScanned.Front())
	}

	// Test Scan with single element
	single := g.DequeOf(10)
	singleScanned := single.Iter().Scan(5, func(acc, val int) int {
		return acc + val
	}).Collect()

	expectedSingle := []int{5, 15}
	actualSingle := make([]int, 0)
	singleScanned.Iter().ForEach(func(v int) { actualSingle = append(actualSingle, v) })

	if len(actualSingle) != len(expectedSingle) {
		t.Errorf("Scan single length: expected %d, got %d", len(expectedSingle), len(actualSingle))
	}

	for i, v := range expectedSingle {
		if i >= len(actualSingle) || actualSingle[i] != v {
			t.Errorf("Scan single[%d]: expected %d, got %v", i, v, actualSingle)
		}
	}

	// Test Scan with string concatenation
	stringDeque := g.DequeOf("a", "b", "c")
	concatenated := stringDeque.Iter().Scan("", func(acc, val string) string {
		return acc + val
	}).Collect()

	expectedConcat := []string{"", "a", "ab", "abc"}
	actualConcat := make([]string, 0)
	concatenated.Iter().ForEach(func(s string) { actualConcat = append(actualConcat, s) })

	if len(actualConcat) != len(expectedConcat) {
		t.Errorf("Scan concat length: expected %d, got %d", len(expectedConcat), len(actualConcat))
	}

	for i, v := range expectedConcat {
		if i >= len(actualConcat) || actualConcat[i] != v {
			t.Errorf("Scan concat[%d]: expected %s, got %v", i, v, actualConcat)
		}
	}
}

func TestSeqDequeNext(t *testing.T) {
	t.Run("Next with non-empty iterator", func(t *testing.T) {
		iter := g.DequeOf(1, 2, 3, 4, 5).Iter()

		// First element
		first := iter.Next()
		if !first.IsSome() || first.Some() != 1 {
			t.Errorf("Expected Some(1), got %v", first)
		}

		// Second element
		second := iter.Next()
		if !second.IsSome() || second.Some() != 2 {
			t.Errorf("Expected Some(2), got %v", second)
		}

		// Remaining elements
		remaining := iter.Collect()
		expected := g.DequeOf(3, 4, 5)
		if !remaining.Eq(expected) {
			t.Errorf("Expected remaining %v, got %v", expected, remaining)
		}
	})

	t.Run("Next with empty iterator", func(t *testing.T) {
		deque := g.Deque[int]{}
		iter := deque.Iter()

		result := iter.Next()
		if result.IsSome() {
			t.Errorf("Expected None, got Some(%v)", result.Some())
		}
	})

	t.Run("Next until exhausted", func(t *testing.T) {
		iter := g.DequeOf(1, 2).Iter()

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
		remaining := iter.Collect()
		if remaining.Len() != 0 {
			t.Errorf("Expected empty deque, got %v", remaining)
		}
	})
}
