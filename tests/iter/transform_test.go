package g_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

func TestScan(t *testing.T) {
	// Test basic scan operation
	result := ToSlice(Scan(FromSlice([]int{1, 2, 3}), 0, func(acc, x int) int {
		return acc + x
	}))
	expected := []int{1, 3, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Scan() = %v, want %v", result, expected)
	}

	// Test scan with empty sequence
	result2 := ToSlice(Scan(FromSlice([]int{}), 10, func(acc, x int) int {
		return acc + x
	}))
	if len(result2) != 0 {
		t.Errorf("Scan(empty) = %v, want empty slice", result2)
	}
}

func TestMap(t *testing.T) {
	// Test basic map operation
	result := ToSlice(Map(FromSlice([]int{1, 2, 3}), func(x int) int { return x * 2 }))
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map() = %v, want %v", result, expected)
	}
}

func TestMapTo(t *testing.T) {
	// Test map to different type
	result := ToSlice(MapTo(FromSlice([]int{1, 2, 3}), func(x int) string { return fmt.Sprintf("%d", x) }))
	expected := []string{"1", "2", "3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapTo() = %v, want %v", result, expected)
	}
}

func TestInspect(t *testing.T) {
	// Test inspect operation
	var sum int
	ToSlice(Inspect(FromSlice([]int{1, 2, 3}), func(x int) { sum += x }))
	if sum != 6 {
		t.Errorf("Inspect() sum = %v, want 6", sum)
	}
}

func TestFilter(t *testing.T) {
	// Test basic filter operation
	result := ToSlice(Filter(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x%2 == 0 }))
	expected := []int{2, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter() = %v, want %v", result, expected)
	}
}

func TestExclude(t *testing.T) {
	// Test basic exclude operation
	result := ToSlice(Exclude(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x%2 == 0 }))
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Exclude() = %v, want %v", result, expected)
	}
}

func TestFilterMap(t *testing.T) {
	// Test filter map operation
	result := ToSlice(FilterMap(FromSlice([]string{"1", "2", "abc", "3"}), func(s string) (int, bool) {
		if len(s) == 1 {
			return len(s), true
		}
		return 0, false
	}))
	expected := []int{1, 1, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FilterMap() = %v, want %v", result, expected)
	}
}

func TestEnumerate(t *testing.T) {
	// Test enumerate operation
	result := ToPairs(Enumerate(FromSlice([]string{"a", "b", "c"}), 0))
	expected := []Pair[int, string]{{0, "a"}, {1, "b"}, {2, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Enumerate() = %v, want %v", result, expected)
	}
}

func TestUnique(t *testing.T) {
	// Test unique operation
	result := ToSlice(Unique(FromSlice([]int{1, 2, 2, 3, 3, 3, 4})))
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unique() = %v, want %v", result, expected)
	}
}

func TestUniqueBy(t *testing.T) {
	// Test unique by operation - consecutive elements with same key are removed
	result := ToSlice(UniqueBy(FromSlice([]string{"a", "b", "bb", "ccc", "dd", "e"}), func(s string) int { return len(s) }))
	expected := []string{"a", "bb", "ccc", "dd", "e"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("UniqueBy() = %v, want %v", result, expected)
	}
}

func TestMap2(t *testing.T) {
	// Test map2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := ToPairs(Map2(FromPairs(pairs), func(k int, v string) (int, string) {
		return k * 2, v + v
	}))
	expected := []Pair[int, string]{{2, "aa"}, {4, "bb"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map2() = %v, want %v", result, expected)
	}
}

func TestFilter2(t *testing.T) {
	// Test filter2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToPairs(Filter2(FromPairs(pairs), func(k int, v string) bool { return k%2 == 0 }))
	expected := []Pair[int, string]{{2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter2() = %v, want %v", result, expected)
	}
}

func TestExclude2(t *testing.T) {
	// Test exclude2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToPairs(Exclude2(FromPairs(pairs), func(k int, v string) bool { return k%2 == 0 }))
	expected := []Pair[int, string]{{1, "a"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Exclude2() = %v, want %v", result, expected)
	}
}

func TestInspect2(t *testing.T) {
	// Test inspect2 operation
	var sum int
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	ToPairs(Inspect2(FromPairs(pairs), func(k int, v string) { sum += k }))
	if sum != 3 {
		t.Errorf("Inspect2() sum = %v, want 3", sum)
	}
}

func TestEnumerateNonZeroStart(t *testing.T) {
	// Test enumerate with non-zero start
	result := ToPairs(Enumerate(FromSlice([]string{"a", "b"}), 5))
	expected := []Pair[int, string]{{5, "a"}, {6, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Enumerate(start=5) = %v, want %v", result, expected)
	}
}

func TestMapWhile(t *testing.T) {
	// Test mapWhile operation
	result := ToSlice(MapWhile(FromSlice([]string{"1", "2", "abc", "3"}), func(s string) (int, bool) {
		if len(s) == 1 {
			return len(s), true
		}
		return 0, false // Stop at "abc"
	}))
	expected := []int{1, 1} // Should stop at "abc"
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapWhile() = %v, want %v", result, expected)
	}

	// Test mapWhile with all matching elements
	result2 := ToSlice(MapWhile(FromSlice([]string{"a", "b", "c"}), func(s string) (int, bool) {
		return len(s), true
	}))
	expected2 := []int{1, 1, 1}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("MapWhile(all match) = %v, want %v", result2, expected2)
	}

	// Test mapWhile with empty sequence
	result3 := ToSlice(MapWhile(FromSlice([]string{}), func(s string) (int, bool) {
		return len(s), true
	}))
	if len(result3) != 0 {
		t.Errorf("MapWhile(empty) = %v, want empty", result3)
	}

	// Test mapWhile with no matching elements
	result4 := ToSlice(MapWhile(FromSlice([]string{"abc", "def"}), func(s string) (int, bool) {
		return len(s), len(s) == 1 // None match
	}))
	if len(result4) != 0 {
		t.Errorf("MapWhile(no match) = %v, want empty", result4)
	}
}
