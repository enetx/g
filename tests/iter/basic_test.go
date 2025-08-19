package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

func TestForEach(t *testing.T) {
	// Test forEach operation
	var sum int
	ForEach(FromSlice([]int{1, 2, 3}), func(x int) { sum += x })
	if sum != 6 {
		t.Errorf("ForEach() sum = %v, want 6", sum)
	}
}

func TestCount(t *testing.T) {
	// Test count operation
	result := Count(FromSlice([]int{1, 2, 3, 4, 5}))
	if result != 5 {
		t.Errorf("Count() = %v, want 5", result)
	}

	// Test count with empty sequence
	result2 := Count(FromSlice([]int{}))
	if result2 != 0 {
		t.Errorf("Count(empty) = %v, want 0", result2)
	}
}

func TestRange(t *testing.T) {
	// Test range operation
	var sum int
	Range(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		sum += x
		return x < 3
	})
	if sum != 6 { // 1 + 2 + 3
		t.Errorf("Range() sum = %v, want 6", sum)
	}
}

func TestTake(t *testing.T) {
	// Test take operation
	result := ToSlice(Take(FromSlice([]int{1, 2, 3, 4, 5}), 3))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Take() = %v, want %v", result, expected)
	}

	// Test take with more than available
	result2 := ToSlice(Take(FromSlice([]int{1, 2}), 5))
	expected2 := []int{1, 2}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Take(more than available) = %v, want %v", result2, expected2)
	}
}

func TestSkip(t *testing.T) {
	// Test skip operation
	result := ToSlice(Skip(FromSlice([]int{1, 2, 3, 4, 5}), 2))
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip() = %v, want %v", result, expected)
	}

	// Test skip with more than available
	result2 := ToSlice(Skip(FromSlice([]int{1, 2}), 5))
	if len(result2) != 0 {
		t.Errorf("Skip(more than available) = %v, want empty slice", result2)
	}
}

func TestStepBy(t *testing.T) {
	// Test stepBy operation
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8}), 3))
	expected := []int{1, 4, 7}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StepBy() = %v, want %v", result, expected)
	}
}

func TestTakeWhile(t *testing.T) {
	// Test takeWhile operation
	result := ToSlice(TakeWhile(FromSlice([]int{1, 2, 3, 4, 2, 1}), func(x int) bool { return x < 4 }))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TakeWhile() = %v, want %v", result, expected)
	}
}

func TestSkipWhile(t *testing.T) {
	// Test skipWhile operation
	result := ToSlice(SkipWhile(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool { return x < 3 }))
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SkipWhile() = %v, want %v", result, expected)
	}
}

func TestNth(t *testing.T) {
	// Test nth operation
	result, found := Nth(FromSlice([]int{10, 20, 30, 40}), 2)
	if !found || result != 30 {
		t.Errorf("Nth() = %v, %v, want 30, true", result, found)
	}

	// Test nth out of bounds
	_, found2 := Nth(FromSlice([]int{1, 2}), 5)
	if found2 {
		t.Errorf("Nth(out of bounds) found = %v, want false", found2)
	}
}

func TestForEach2(t *testing.T) {
	// Test forEach2 operation
	var sum int
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	ForEach2(FromPairs(pairs), func(k int, v string) { sum += k })
	if sum != 3 {
		t.Errorf("ForEach2() sum = %v, want 3", sum)
	}
}

func TestCount2(t *testing.T) {
	// Test count2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := Count2(FromPairs(pairs))
	if result != 2 {
		t.Errorf("Count2() = %v, want 2", result)
	}
}

func TestRange2(t *testing.T) {
	// Test range2 operation
	var sum int
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	Range2(FromPairs(pairs), func(k int, v string) bool {
		sum += k
		return k < 2
	})
	if sum != 3 { // 1 + 2
		t.Errorf("Range2() sum = %v, want 3", sum)
	}
}

func TestTake2(t *testing.T) {
	// Test take2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToPairs(Take2(FromPairs(pairs), 2))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Take2() = %v, want %v", result, expected)
	}
}

func TestSkip2(t *testing.T) {
	// Test skip2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToPairs(Skip2(FromPairs(pairs), 1))
	expected := []Pair[int, string]{{2, "b"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip2() = %v, want %v", result, expected)
	}
}

func TestStepBy2(t *testing.T) {
	// Test stepBy2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	result := ToPairs(StepBy2(FromPairs(pairs), 2))
	expected := []Pair[int, string]{{1, "a"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StepBy2() = %v, want %v", result, expected)
	}
}

func TestNth2(t *testing.T) {
	// Test nth2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	k, v, found := Nth2(FromPairs(pairs), 1)
	if !found || k != 2 || v != "b" {
		t.Errorf("Nth2() = %v, %v, %v, want 2, b, true", k, v, found)
	}
}

// Edge cases tests

func TestSkipNegative(t *testing.T) {
	// Test skip with negative value
	result := ToSlice(Skip(FromSlice([]int{1, 2, 3}), -1))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip(negative) = %v, want %v", result, expected)
	}
}

func TestStepByNegative(t *testing.T) {
	// Test stepBy with invalid step
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("StepBy(0) = %v, want empty slice", result)
	}
}

func TestTakeZero(t *testing.T) {
	// Test take with zero
	result := ToSlice(Take(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("Take(0) = %v, want empty slice", result)
	}
}

func TestTake2Zero(t *testing.T) {
	// Test take2 with zero
	pairs := []Pair[int, string]{{1, "a"}}
	result := ToPairs(Take2(FromPairs(pairs), 0))
	if len(result) != 0 {
		t.Errorf("Take2(0) = %v, want empty slice", result)
	}
}

func TestSkipZero(t *testing.T) {
	// Test skip with zero
	result := ToSlice(Skip(FromSlice([]int{1, 2, 3}), 0))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip(0) = %v, want %v", result, expected)
	}
}

func TestSkip2Zero(t *testing.T) {
	// Test skip2 with zero
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := ToPairs(Skip2(FromPairs(pairs), 0))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip2(0) = %v, want %v", result, expected)
	}
}

func TestStepByZero(t *testing.T) {
	// Test stepBy with zero step
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("StepBy(0) = %v, want empty slice", result)
	}
}

func TestStepBy2Zero(t *testing.T) {
	// Test stepBy2 with zero step
	pairs := []Pair[int, string]{{1, "a"}}
	result := ToPairs(StepBy2(FromPairs(pairs), 0))
	if len(result) != 0 {
		t.Errorf("StepBy2(0) = %v, want empty slice", result)
	}
}

func TestStepByZeroStepEdgeCase(t *testing.T) {
	// Test stepBy with zero step edge case
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3, 4, 5}), 0))
	if len(result) != 0 {
		t.Errorf("StepBy zero step = %v, want empty slice", result)
	}
}

func TestStepBy2ZeroStepEdgeCase(t *testing.T) {
	// Test stepBy2 with zero step edge case
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := ToPairs(StepBy2(FromPairs(pairs), 0))
	if len(result) != 0 {
		t.Errorf("StepBy2 zero step = %v, want empty slice", result)
	}
}

func TestStepByEarlyTermination(t *testing.T) {
	// Test stepBy with early termination
	count := 0
	StepBy(FromSlice([]int{1, 2, 3, 4, 5, 6}), 2)(func(x int) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("StepByEarlyTermination count = %v, want 2", count)
	}
}

func TestStepBy2EarlyTermination(t *testing.T) {
	// Test stepBy2 with early termination
	count := 0
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	StepBy2(FromPairs(pairs), 2)(func(k int, v string) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("StepBy2EarlyTermination count = %v, want 2", count)
	}
}
