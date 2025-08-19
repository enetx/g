package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

func TestKeys(t *testing.T) {
	// Test keys operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToSlice(Keys(FromPairs(pairs)))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Keys() = %v, want %v", result, expected)
	}
}

func TestValues(t *testing.T) {
	// Test values operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToSlice(Values(FromPairs(pairs)))
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Values() = %v, want %v", result, expected)
	}
}

func TestFromPairs(t *testing.T) {
	// Test fromPairs operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := ToPairs(FromPairs(pairs))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromPairs() = %v, want %v", result, expected)
	}
}

func TestFromPairsEarlyTermination(t *testing.T) {
	// Test fromPairs with early termination
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}

	count := 0
	FromPairs(pairs)(func(k int, v string) bool {
		count++
		return count < 3
	})

	if count != 3 {
		t.Errorf("FromPairsEarlyTermination count = %v, want 3", count)
	}
}

func TestOrderByKey(t *testing.T) {
	// Test orderByKey operation
	pairs := []Pair[int, string]{{3, "c"}, {1, "a"}, {2, "b"}}
	result := ToPairs(OrderByKey(FromPairs(pairs), func(a, b int) bool { return a < b }))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderByKey() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := ToPairs(OrderByKey(FromPairs([]Pair[int, string]{}), func(a, b int) bool { return a < b }))
	if len(result2) != 0 {
		t.Errorf("OrderByKey(empty) = %v, want empty", result2)
	}
}

func TestOrderByValue(t *testing.T) {
	// Test orderByValue operation
	pairs := []Pair[int, string]{{1, "c"}, {2, "a"}, {3, "b"}}
	result := ToPairs(OrderByValue(FromPairs(pairs), func(a, b string) bool { return a < b }))
	expected := []Pair[int, string]{{2, "a"}, {3, "b"}, {1, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderByValue() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := ToPairs(OrderByValue(FromPairs([]Pair[int, string]{}), func(a, b string) bool { return a < b }))
	if len(result2) != 0 {
		t.Errorf("OrderByValue(empty) = %v, want empty", result2)
	}
}

func TestSortBy2(t *testing.T) {
	// Test sortBy2 operation
	pairs := []Pair[string, int]{{"c", 3}, {"a", 1}, {"b", 2}}
	result := ToPairs(SortBy2(FromPairs(pairs), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}))
	expected := []Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortBy2() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := ToPairs(SortBy2(FromPairs([]Pair[string, int]{}), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}))
	if len(result2) != 0 {
		t.Errorf("SortBy2(empty) = %v, want empty", result2)
	}

	// Test reverse sort by value
	result3 := ToPairs(SortBy2(FromPairs(pairs), func(a, b Pair[string, int]) bool {
		return a.Value > b.Value
	}))
	expected3 := []Pair[string, int]{{"c", 3}, {"b", 2}, {"a", 1}}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("SortBy2(reverse by value) = %v, want %v", result3, expected3)
	}

	// Test single element
	single := []Pair[string, int]{{"a", 1}}
	result4 := ToPairs(SortBy2(FromPairs(single), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}))
	if !reflect.DeepEqual(result4, single) {
		t.Errorf("SortBy2(single) = %v, want %v", result4, single)
	}

	// Test equal elements
	equal := []Pair[string, int]{{"a", 1}, {"a", 2}, {"a", 3}}
	result5 := ToPairs(SortBy2(FromPairs(equal), func(a, b Pair[string, int]) bool {
		return a.Value < b.Value
	}))
	expected5 := []Pair[string, int]{{"a", 1}, {"a", 2}, {"a", 3}}
	if !reflect.DeepEqual(result5, expected5) {
		t.Errorf("SortBy2(equal keys, sort by value) = %v, want %v", result5, expected5)
	}
}

func TestOrderByKeyAdvanced(t *testing.T) {
	// Test OrderByKey with single element
	single := []Pair[int, string]{{1, "a"}}
	result := ToPairs(OrderByKey(FromPairs(single), func(a, b int) bool { return a < b }))
	if !reflect.DeepEqual(result, single) {
		t.Errorf("OrderByKey(single) = %v, want %v", result, single)
	}

	// Test OrderByKey with already sorted data
	sorted := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result2 := ToPairs(OrderByKey(FromPairs(sorted), func(a, b int) bool { return a < b }))
	if !reflect.DeepEqual(result2, sorted) {
		t.Errorf("OrderByKey(already sorted) = %v, want %v", result2, sorted)
	}

	// Test OrderByKey with duplicate keys
	duplicates := []Pair[int, string]{{2, "b"}, {1, "a"}, {2, "c"}, {1, "d"}}
	result3 := ToPairs(OrderByKey(FromPairs(duplicates), func(a, b int) bool { return a < b }))
	// Should have keys in order 1, 1, 2, 2
	if len(result3) != 4 || result3[0].Key != 1 || result3[1].Key != 1 || result3[2].Key != 2 || result3[3].Key != 2 {
		t.Errorf("OrderByKey(duplicates) = %v, want keys ordered as [1,1,2,2]", result3)
	}
}

func TestOrderByValueAdvanced(t *testing.T) {
	// Test OrderByValue with single element
	single := []Pair[int, string]{{1, "a"}}
	result := ToPairs(OrderByValue(FromPairs(single), func(a, b string) bool { return a < b }))
	if !reflect.DeepEqual(result, single) {
		t.Errorf("OrderByValue(single) = %v, want %v", result, single)
	}

	// Test OrderByValue with already sorted data
	sorted := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result2 := ToPairs(OrderByValue(FromPairs(sorted), func(a, b string) bool { return a < b }))
	if !reflect.DeepEqual(result2, sorted) {
		t.Errorf("OrderByValue(already sorted) = %v, want %v", result2, sorted)
	}

	// Test OrderByValue with duplicate values
	duplicates := []Pair[int, string]{{2, "b"}, {1, "a"}, {3, "b"}, {4, "a"}}
	result3 := ToPairs(OrderByValue(FromPairs(duplicates), func(a, b string) bool { return a < b }))
	// Should have values in order a, a, b, b
	if len(result3) != 4 || result3[0].Value != "a" || result3[1].Value != "a" || result3[2].Value != "b" || result3[3].Value != "b" {
		t.Errorf("OrderByValue(duplicates) = %v, want values ordered as [a,a,b,b]", result3)
	}
}
