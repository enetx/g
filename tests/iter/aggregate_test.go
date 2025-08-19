package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

func TestFind(t *testing.T) {
	// Test find operation
	result, found := Find(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x > 3 })
	if !found || result != 4 {
		t.Errorf("Find() = %v, %v, want 4, true", result, found)
	}

	// Test find not found
	_, found2 := Find(FromSlice([]int{1, 2, 3}), func(x int) bool { return x > 10 })
	if found2 {
		t.Errorf("Find(not found) found = %v, want false", found2)
	}
}

func TestAny(t *testing.T) {
	// Test any operation
	result := Any(FromSlice([]int{1, 2, 3, 4}), func(x int) bool { return x > 3 })
	if !result {
		t.Errorf("Any() = %v, want true", result)
	}

	// Test any false
	result2 := Any(FromSlice([]int{1, 2, 3}), func(x int) bool { return x > 10 })
	if result2 {
		t.Errorf("Any(false) = %v, want false", result2)
	}
}

func TestAll(t *testing.T) {
	// Test all operation
	result := All(FromSlice([]int{1, 2, 3, 4}), func(x int) bool { return x > 0 })
	if !result {
		t.Errorf("All() = %v, want true", result)
	}

	// Test all false
	result2 := All(FromSlice([]int{1, 2, 3}), func(x int) bool { return x > 2 })
	if result2 {
		t.Errorf("All(false) = %v, want false", result2)
	}
}

func TestFold(t *testing.T) {
	// Test fold operation
	result := Fold(FromSlice([]int{1, 2, 3, 4}), 0, func(acc, x int) int { return acc + x })
	if result != 10 {
		t.Errorf("Fold() = %v, want 10", result)
	}

	// Test fold with empty sequence
	result2 := Fold(FromSlice([]int{}), 5, func(acc, x int) int { return acc + x })
	if result2 != 5 {
		t.Errorf("Fold(empty) = %v, want 5", result2)
	}
}

func TestReduce(t *testing.T) {
	// Test reduce operation
	result, found := Reduce(FromSlice([]int{1, 2, 3, 4}), func(a, b int) int { return a + b })
	if !found || result != 10 {
		t.Errorf("Reduce() = %v, %v, want 10, true", result, found)
	}

	// Test reduce with empty sequence
	_, found2 := Reduce(FromSlice([]int{}), func(a, b int) int { return a + b })
	if found2 {
		t.Errorf("Reduce(empty) found = %v, want false", found2)
	}
}

func TestContains(t *testing.T) {
	// Test contains operation
	result := Contains(FromSlice([]int{1, 2, 3, 4}), 3)
	if !result {
		t.Errorf("Contains() = %v, want true", result)
	}

	// Test contains false
	result2 := Contains(FromSlice([]int{1, 2, 3}), 5)
	if result2 {
		t.Errorf("Contains(false) = %v, want false", result2)
	}
}

func TestFind2(t *testing.T) {
	// Test find2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "bb"}, {3, "c"}}
	k, v, found := Find2(FromPairs(pairs), func(key int, val string) bool { return len(val) > 1 })
	if !found || k != 2 || v != "bb" {
		t.Errorf("Find2() = %v, %v, %v, want 2, bb, true", k, v, found)
	}
}

func TestAny2(t *testing.T) {
	// Test any2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "bb"}}
	result := Any2(FromPairs(pairs), func(k int, v string) bool { return len(v) > 1 })
	if !result {
		t.Errorf("Any2() = %v, want true", result)
	}
}

func TestAll2(t *testing.T) {
	// Test all2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := All2(FromPairs(pairs), func(k int, v string) bool { return len(v) == 1 })
	if !result {
		t.Errorf("All2() = %v, want true", result)
	}

	// Test all2 with early false
	pairs2 := []Pair[int, string]{{1, "a"}, {2, "bb"}, {3, "c"}}
	result2 := All2(FromPairs(pairs2), func(k int, v string) bool { return len(v) == 1 })
	if result2 {
		t.Errorf("All2(early false) = %v, want false", result2)
	}

	// Test all2 with empty sequence
	result3 := All2(FromPairs([]Pair[int, string]{}), func(k int, v string) bool { return len(v) == 1 })
	if !result3 {
		t.Errorf("All2(empty) = %v, want true", result3)
	}

	// Test all2 early termination
	count := 0
	pairs3 := []Pair[int, string]{{1, "aa"}, {2, "bb"}, {3, "cc"}}
	All2(FromPairs(pairs3), func(k int, v string) bool {
		count++
		return false // First check returns false, should terminate
	})
	if count != 1 {
		t.Errorf("All2 early termination count = %v, want 1", count)
	}
}

func TestFold2(t *testing.T) {
	// Test fold2 operation
	pairs := []Pair[int, int]{{1, 10}, {2, 20}}
	result := Fold2(FromPairs(pairs), 0, func(acc, k, v int) int { return acc + k + v })
	if result != 33 { // 0 + 1 + 10 + 2 + 20
		t.Errorf("Fold2() = %v, want 33", result)
	}
}

func TestReduce2(t *testing.T) {
	// Test reduce2 operation
	pairs := []Pair[int, int]{{1, 10}, {2, 20}}
	result, found := Reduce2(FromPairs(pairs), func(a, b Pair[int, int]) Pair[int, int] {
		return Pair[int, int]{a.Key + b.Key, a.Value + b.Value}
	})
	if !found || result.Key != 3 || result.Value != 30 {
		t.Errorf("Reduce2() = %v, %v, want {3, 30}, true", result, found)
	}
}

func TestMinBy(t *testing.T) {
	// Test basic minBy operation
	result, found := MinBy(FromSlice([]int{3, 1, 4, 1, 5}), func(a, b int) bool { return a < b })
	if !found || result != 1 {
		t.Errorf("MinBy() = %v, %v, want 1, true", result, found)
	}

	// Test empty sequence
	_, found2 := MinBy(FromSlice([]int{}), func(a, b int) bool { return a < b })
	if found2 {
		t.Errorf("MinBy(empty) found = %v, want false", found2)
	}
}

func TestMaxBy(t *testing.T) {
	// Test basic maxBy operation
	result, found := MaxBy(FromSlice([]int{3, 1, 4, 1, 5}), func(a, b int) bool { return a < b })
	if !found || result != 5 {
		t.Errorf("MaxBy() = %v, %v, want 5, true", result, found)
	}

	// Test empty sequence
	_, found2 := MaxBy(FromSlice([]int{}), func(a, b int) bool { return a < b })
	if found2 {
		t.Errorf("MaxBy(empty) found = %v, want false", found2)
	}
}

func TestCountBy(t *testing.T) {
	// Test countBy operation
	result := CountBy(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x%2 == 0 })
	expected := map[bool]int{true: 2, false: 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("CountBy() = %v, want %v", result, expected)
	}

	// Test empty sequence
	result2 := CountBy(FromSlice([]int{}), func(x int) bool { return x%2 == 0 })
	if len(result2) != 0 {
		t.Errorf("CountBy(empty) = %v, want empty map", result2)
	}
}

func TestPartition(t *testing.T) {
	// Test partition operation
	left, right := Partition(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x%2 == 0 })

	expectedLeft := []int{2, 4}
	expectedRight := []int{1, 3, 5}

	if !reflect.DeepEqual(left, expectedLeft) {
		t.Errorf("Partition() left = %v, want %v", left, expectedLeft)
	}
	if !reflect.DeepEqual(right, expectedRight) {
		t.Errorf("Partition() right = %v, want %v", right, expectedRight)
	}
}

func TestSortBy(t *testing.T) {
	// Test sortBy operation
	result := ToSlice(SortBy(FromSlice([]int{3, 1, 4, 1, 5}), func(a, b int) bool { return a < b }))
	expected := []int{1, 1, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortBy() = %v, want %v", result, expected)
	}

	// Test reverse sort
	result2 := ToSlice(SortBy(FromSlice([]int{3, 1, 4, 1, 5}), func(a, b int) bool { return a > b }))
	expected2 := []int{5, 4, 3, 1, 1}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("SortBy(reverse) = %v, want %v", result2, expected2)
	}
}

func TestPosition(t *testing.T) {
	// Test position operation
	result, found := Position(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x == 3 })
	if !found || result != 2 {
		t.Errorf("Position() = %v, %v, want 2, true", result, found)
	}

	// Test not found
	_, found2 := Position(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x == 10 })
	if found2 {
		t.Errorf("Position(not found) found = %v, want false", found2)
	}
}

func TestRPosition(t *testing.T) {
	// Test rposition operation
	result, found := RPosition(FromSlice([]int{1, 2, 3, 2, 5}), func(x int) bool { return x == 2 })
	if !found || result != 3 {
		t.Errorf("RPosition() = %v, %v, want 3, true", result, found)
	}

	// Test not found
	_, found2 := RPosition(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x == 10 })
	if found2 {
		t.Errorf("RPosition(not found) found = %v, want false", found2)
	}
}

func TestIsPartitioned(t *testing.T) {
	// Test is_partitioned operation
	result := IsPartitioned(FromSlice([]int{2, 4, 1, 3, 5}), func(x int) bool { return x%2 == 0 })
	if !result {
		t.Errorf("IsPartitioned() = %v, want true", result)
	}

	// Test not partitioned
	result2 := IsPartitioned(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x%2 == 0 })
	if result2 {
		t.Errorf("IsPartitioned(not partitioned) = %v, want false", result2)
	}

	// Test empty sequence
	result3 := IsPartitioned(FromSlice([]int{}), func(x int) bool { return x%2 == 0 })
	if !result3 {
		t.Errorf("IsPartitioned(empty) = %v, want true", result3)
	}
}
