package g_test

import (
	"context"
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

func TestMapWhile(t *testing.T) {
	// Test map while condition is true
	result := ToSlice(MapWhile(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) (int, bool) {
		if x < 4 {
			return x * 2, true
		}
		return 0, false
	}))
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapWhile() = %v, want %v", result, expected)
	}

	// Test map while with immediate false
	result2 := ToSlice(MapWhile(FromSlice([]int{1, 2, 3}), func(int) (int, bool) {
		return 0, false
	}))
	if len(result2) != 0 {
		t.Errorf("MapWhile(immediate false) = %v, want empty slice", result2)
	}
}

func TestPosition(t *testing.T) {
	// Test finding position
	pos, found := Position(FromSlice([]int{10, 20, 30, 40}), func(x int) bool {
		return x > 25
	})
	if !found || pos != 2 {
		t.Errorf("Position() = %d, %t, want 2, true", pos, found)
	}

	// Test not found
	pos2, found2 := Position(FromSlice([]int{10, 20, 30}), func(x int) bool {
		return x > 100
	})
	if found2 {
		t.Errorf("Position(not found) = %d, %t, want _, false", pos2, found2)
	}
}

func TestRPosition(t *testing.T) {
	// Test finding last position
	pos, found := RPosition(FromSlice([]int{10, 30, 20, 30}), func(x int) bool {
		return x == 30
	})
	if !found || pos != 3 {
		t.Errorf("RPosition() = %d, %t, want 3, true", pos, found)
	}

	// Test not found
	pos2, found2 := RPosition(FromSlice([]int{10, 20, 30}), func(x int) bool {
		return x > 100
	})
	if found2 {
		t.Errorf("RPosition(not found) = %d, %t, want _, false", pos2, found2)
	}
}

func TestCmp(t *testing.T) {
	// Test equal sequences
	cmp1 := Cmp(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp1 != 0 {
		t.Errorf("Cmp(equal) = %d, want 0", cmp1)
	}

	// Test first less than second
	cmp2 := Cmp(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp2 != -1 {
		t.Errorf("Cmp(less) = %d, want -1", cmp2)
	}

	// Test first shorter than second
	cmp3 := Cmp(FromSlice([]int{1, 2}), FromSlice([]int{1, 2, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp3 != -1 {
		t.Errorf("Cmp(shorter) = %d, want -1", cmp3)
	}
}

func TestEqual(t *testing.T) {
	// Test equal sequences
	if !Equal(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2, 3})) {
		t.Error("Equal(same) should be true")
	}

	// Test different sequences
	if Equal(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2, 4})) {
		t.Error("Equal(different) should be false")
	}

	// Test different lengths
	if Equal(FromSlice([]int{1, 2}), FromSlice([]int{1, 2, 3})) {
		t.Error("Equal(different lengths) should be false")
	}
}

func TestLt(t *testing.T) {
	less := func(a, b int) bool { return a < b }

	// Test less than
	if !Lt(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), less) {
		t.Error("Lt([1,2], [1,3]) should be true")
	}

	// Test not less than
	if Lt(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), less) {
		t.Error("Lt([1,3], [1,2]) should be false")
	}

	// Test equal sequences
	if Lt(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), less) {
		t.Error("Lt(equal) should be false")
	}
}

func TestOnce(t *testing.T) {
	result := ToSlice(Once(42))
	expected := []int{42}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Once(42) = %v, want %v", result, expected)
	}
}

func TestOnceWith(t *testing.T) {
	called := false
	result := ToSlice(OnceWith(func() int {
		called = true
		return 100
	}))
	expected := []int{100}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OnceWith() = %v, want %v", result, expected)
	}
	if !called {
		t.Error("OnceWith function should have been called")
	}
}

func TestEmpty(t *testing.T) {
	result := ToSlice(Empty[int]())
	if len(result) != 0 {
		t.Errorf("Empty() = %v, want empty slice", result)
	}

	// Additional test to ensure the yield function is actually called (it shouldn't be)
	called := false
	Empty[int]()(func(int) bool {
		called = true
		return true
	})

	if called {
		t.Errorf("Empty sequence should not call yield function")
	}
}

func TestRepeat(t *testing.T) {
	result := ToSlice(Take(Repeat(7), 3))
	expected := []int{7, 7, 7}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Take(Repeat(7), 3) = %v, want %v", result, expected)
	}
}

func TestRepeatWith(t *testing.T) {
	counter := 0
	result := ToSlice(Take(RepeatWith(func() int {
		counter++
		return counter
	}), 3))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Take(RepeatWith(counter), 3) = %v, want %v", result, expected)
	}
}

func TestIsPartitioned(t *testing.T) {
	isEven := func(x int) bool { return x%2 == 0 }

	// Test properly partitioned (even numbers first)
	if !IsPartitioned(FromSlice([]int{2, 4, 6, 1, 3, 5}), isEven) {
		t.Error("IsPartitioned([2,4,6,1,3,5], isEven) should be true")
	}

	// Test not partitioned
	if IsPartitioned(FromSlice([]int{2, 1, 4, 3}), isEven) {
		t.Error("IsPartitioned([2,1,4,3], isEven) should be false")
	}

	// Test all match predicate
	if !IsPartitioned(FromSlice([]int{2, 4, 6}), isEven) {
		t.Error("IsPartitioned([2,4,6], isEven) should be true")
	}

	// Test none match predicate
	if !IsPartitioned(FromSlice([]int{1, 3, 5}), isEven) {
		t.Error("IsPartitioned([1,3,5], isEven) should be true")
	}

	// Test empty sequence
	if !IsPartitioned(FromSlice([]int{}), isEven) {
		t.Error("IsPartitioned([], isEven) should be true")
	}
}

// ========== Tests for basic operations ==========

func TestForEach(t *testing.T) {
	var result []int
	ForEach(FromSlice([]int{1, 2, 3}), func(x int) {
		result = append(result, x*2)
	})
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ForEach() = %v, want %v", result, expected)
	}
}

func TestCount(t *testing.T) {
	count := Count(FromSlice([]int{1, 2, 3, 4, 5}))
	if count != 5 {
		t.Errorf("Count() = %d, want 5", count)
	}

	// Test empty sequence
	count2 := Count(FromSlice([]int{}))
	if count2 != 0 {
		t.Errorf("Count(empty) = %d, want 0", count2)
	}
}

func TestRange(t *testing.T) {
	var result []int
	Range(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		result = append(result, x)
		return x < 3
	})
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Range() = %v, want %v", result, expected)
	}
}

func TestTake(t *testing.T) {
	result := ToSlice(Take(FromSlice([]int{1, 2, 3, 4, 5}), 3))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Take() = %v, want %v", result, expected)
	}

	// Test take more than available
	result2 := ToSlice(Take(FromSlice([]int{1, 2}), 5))
	expected2 := []int{1, 2}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Take(more than available) = %v, want %v", result2, expected2)
	}

	// Test take zero
	result3 := ToSlice(Take(FromSlice([]int{1, 2, 3}), 0))
	if len(result3) != 0 {
		t.Errorf("Take(0) = %v, want empty", result3)
	}
}

func TestSkip(t *testing.T) {
	result := ToSlice(Skip(FromSlice([]int{1, 2, 3, 4, 5}), 2))
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip() = %v, want %v", result, expected)
	}

	// Test skip more than available
	result2 := ToSlice(Skip(FromSlice([]int{1, 2}), 5))
	if len(result2) != 0 {
		t.Errorf("Skip(more than available) = %v, want empty", result2)
	}
}

func TestStepBy(t *testing.T) {
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3, 4, 5, 6}), 2))
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StepBy(2) = %v, want %v", result, expected)
	}

	// Test step by 1
	result2 := ToSlice(StepBy(FromSlice([]int{1, 2, 3}), 1))
	expected2 := []int{1, 2, 3}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("StepBy(1) = %v, want %v", result2, expected2)
	}
}

func TestTakeWhile(t *testing.T) {
	result := ToSlice(TakeWhile(FromSlice([]int{1, 2, 3, 4, 1}), func(x int) bool {
		return x < 4
	}))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TakeWhile() = %v, want %v", result, expected)
	}
}

func TestSkipWhile(t *testing.T) {
	result := ToSlice(SkipWhile(FromSlice([]int{1, 2, 3, 4, 1}), func(x int) bool {
		return x < 3
	}))
	expected := []int{3, 4, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SkipWhile() = %v, want %v", result, expected)
	}
}

func TestChain(t *testing.T) {
	s1 := FromSlice([]int{1, 2})
	s2 := FromSlice([]int{3, 4})
	s3 := FromSlice([]int{5})
	result := ToSlice(Chain(s1, s2, s3))
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chain() = %v, want %v", result, expected)
	}
}

func TestMap(t *testing.T) {
	result := ToSlice(Map(FromSlice([]int{1, 2, 3}), func(x int) int {
		return x * 2
	}))
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map() = %v, want %v", result, expected)
	}
}

func TestMapTo(t *testing.T) {
	result := ToSlice(MapTo(FromSlice([]int{1, 2, 3}), func(x int) string {
		return fmt.Sprintf("num_%d", x)
	}))
	expected := []string{"num_1", "num_2", "num_3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapTo() = %v, want %v", result, expected)
	}
}

func TestInspect(t *testing.T) {
	var inspected []int
	result := ToSlice(Inspect(FromSlice([]int{1, 2, 3}), func(x int) {
		inspected = append(inspected, x)
	}))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Inspect() = %v, want %v", result, expected)
	}
	if !reflect.DeepEqual(inspected, expected) {
		t.Errorf("Inspect side effect = %v, want %v", inspected, expected)
	}
}

func TestFilter(t *testing.T) {
	result := ToSlice(Filter(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		return x%2 == 0
	}))
	expected := []int{2, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter() = %v, want %v", result, expected)
	}
}

func TestExclude(t *testing.T) {
	result := ToSlice(Exclude(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		return x%2 == 0
	}))
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Exclude() = %v, want %v", result, expected)
	}
}

func TestFilterMap(t *testing.T) {
	result := ToSlice(FilterMap(FromSlice([]int{1, 2, 3, 4}), func(x int) (int, bool) {
		if x%2 == 0 {
			return x * 2, true
		}
		return 0, false
	}))
	expected := []int{4, 8}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FilterMap() = %v, want %v", result, expected)
	}
}

func TestFind(t *testing.T) {
	val, found := Find(FromSlice([]int{1, 2, 3, 4}), func(x int) bool {
		return x > 2
	})
	if !found || val != 3 {
		t.Errorf("Find() = %d, %t, want 3, true", val, found)
	}

	// Test not found
	_, found2 := Find(FromSlice([]int{1, 2, 3}), func(x int) bool {
		return x > 10
	})
	if found2 {
		t.Error("Find(not found) should return false")
	}
}

func TestNth(t *testing.T) {
	val, found := Nth(FromSlice([]int{10, 20, 30}), 1)
	if !found || val != 20 {
		t.Errorf("Nth(1) = %d, %t, want 20, true", val, found)
	}

	// Test out of bounds
	_, found2 := Nth(FromSlice([]int{1, 2}), 5)
	if found2 {
		t.Error("Nth(out of bounds) should return false")
	}

	// Test negative index
	_, found3 := Nth(FromSlice([]int{1, 2, 3}), -1)
	if found3 {
		t.Error("Nth(negative) should return false")
	}
}

func TestAny(t *testing.T) {
	if !Any(FromSlice([]int{1, 2, 3}), func(x int) bool { return x > 2 }) {
		t.Error("Any() should return true when predicate matches")
	}

	if Any(FromSlice([]int{1, 2, 3}), func(x int) bool { return x > 10 }) {
		t.Error("Any() should return false when predicate never matches")
	}

	// Test empty sequence
	if Any(FromSlice([]int{}), func(int) bool { return true }) {
		t.Error("Any(empty) should return false")
	}
}

func TestAll(t *testing.T) {
	if !All(FromSlice([]int{2, 4, 6}), func(x int) bool { return x%2 == 0 }) {
		t.Error("All() should return true when all elements match")
	}

	if All(FromSlice([]int{1, 2, 3}), func(x int) bool { return x%2 == 0 }) {
		t.Error("All() should return false when some elements don't match")
	}

	// Test empty sequence
	if !All(FromSlice([]int{}), func(int) bool { return false }) {
		t.Error("All(empty) should return true")
	}
}

func TestFold(t *testing.T) {
	sum := Fold(FromSlice([]int{1, 2, 3}), 0, func(acc, x int) int {
		return acc + x
	})
	if sum != 6 {
		t.Errorf("Fold() = %d, want 6", sum)
	}

	// Test with different types
	concat := Fold(FromSlice([]string{"a", "b", "c"}), "", func(acc, x string) string {
		return acc + x
	})
	if concat != "abc" {
		t.Errorf("Fold(strings) = %s, want abc", concat)
	}
}

func TestReduce(t *testing.T) {
	sum, ok := Reduce(FromSlice([]int{1, 2, 3}), func(a, b int) int {
		return a + b
	})
	if !ok || sum != 6 {
		t.Errorf("Reduce() = %d, %t, want 6, true", sum, ok)
	}

	// Test empty sequence
	_, ok2 := Reduce(FromSlice([]int{}), func(a, b int) int {
		return a + b
	})
	if ok2 {
		t.Error("Reduce(empty) should return false")
	}
}

func TestEnumerate(t *testing.T) {
	result := ToPairs(Enumerate(FromSlice([]string{"a", "b", "c"}), 0))
	expected := []Pair[int, string]{{0, "a"}, {1, "b"}, {2, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Enumerate() = %v, want %v", result, expected)
	}

	// Test with different start
	result2 := ToPairs(Enumerate(FromSlice([]string{"x", "y"}), 10))
	expected2 := []Pair[int, string]{{10, "x"}, {11, "y"}}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Enumerate(start=10) = %v, want %v", result2, expected2)
	}
}

func TestZip(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]string{"a", "b"})
	result := ToPairs(Zip(s1, s2))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip() = %v, want %v", result, expected)
	}
}

func TestZipWith(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{10, 20, 30})
	result := ToSlice(ZipWith(s1, s2, func(a, b int) int {
		return a + b
	}))
	expected := []int{11, 22, 33}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ZipWith() = %v, want %v", result, expected)
	}
}

func TestInterleave(t *testing.T) {
	s1 := FromSlice([]int{1, 3, 5})
	s2 := FromSlice([]int{2, 4})
	result := ToSlice(Interleave(s1, s2))
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Interleave() = %v, want %v", result, expected)
	}
}

// ========== Tests for generators ==========

func TestIota(t *testing.T) {
	// Test basic range
	result := ToSlice(Iota(0, 5))
	expected := []int{0, 1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Iota(0, 5) = %v, want %v", result, expected)
	}

	// Test with step
	result2 := ToSlice(Iota(0, 10, 2))
	expected2 := []int{0, 2, 4, 6, 8}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Iota(0, 10, 2) = %v, want %v", result2, expected2)
	}

	// Test negative step
	result3 := ToSlice(Iota(10, 0, -2))
	expected3 := []int{10, 8, 6, 4, 2}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Iota(10, 0, -2) = %v, want %v", result3, expected3)
	}

	// Test zero step (should return empty)
	result4 := ToSlice(Iota(0, 5, 0))
	if len(result4) != 0 {
		t.Errorf("Iota(0, 5, 0) = %v, want empty", result4)
	}
}

func TestIotaInclusive(t *testing.T) {
	// Test inclusive range
	result := ToSlice(IotaInclusive(0, 3))
	expected := []int{0, 1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("IotaInclusive(0, 3) = %v, want %v", result, expected)
	}

	// Test inclusive with negative step
	result2 := ToSlice(IotaInclusive(5, 1, -1))
	expected2 := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("IotaInclusive(5, 1, -1) = %v, want %v", result2, expected2)
	}
}

// ========== Tests for source functions ==========

func TestFromSlice(t *testing.T) {
	result := ToSlice(FromSlice([]int{1, 2, 3}))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromSlice() = %v, want %v", result, expected)
	}

	// Test empty slice
	result2 := ToSlice(FromSlice([]int{}))
	if len(result2) != 0 {
		t.Errorf("FromSlice(empty) = %v, want empty", result2)
	}
}

func TestFromChan(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	result := ToSlice(FromChan(ch))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromChan() = %v, want %v", result, expected)
	}
}

func TestFromMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := ToMap(FromMap(m))

	if len(result) != 2 {
		t.Errorf("FromMap() length = %d, want 2", len(result))
	}
	if result["a"] != 1 || result["b"] != 2 {
		t.Errorf("FromMap() values incorrect: %v", result)
	}
}

// ========== Tests for advanced operations ==========

func TestWindows(t *testing.T) {
	result := ToSlice(Windows(FromSlice([]int{1, 2, 3, 4, 5}), 3))
	expected := [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Windows() = %v, want %v", result, expected)
	}

	// Test window size larger than sequence
	result2 := ToSlice(Windows(FromSlice([]int{1, 2}), 5))
	if len(result2) != 0 {
		t.Errorf("Windows(large size) = %v, want empty", result2)
	}

	// Test invalid window size
	result3 := ToSlice(Windows(FromSlice([]int{1, 2, 3}), 0))
	if len(result3) != 0 {
		t.Errorf("Windows(0) = %v, want empty", result3)
	}
}

func TestChunks(t *testing.T) {
	result := ToSlice(Chunks(FromSlice([]int{1, 2, 3, 4, 5}), 2))
	expected := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chunks() = %v, want %v", result, expected)
	}

	// Test exact division
	result2 := ToSlice(Chunks(FromSlice([]int{1, 2, 3, 4}), 2))
	expected2 := [][]int{{1, 2}, {3, 4}}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Chunks(exact) = %v, want %v", result2, expected2)
	}
}

func TestUnique(t *testing.T) {
	result := ToSlice(Unique(FromSlice([]int{1, 2, 2, 3, 1, 4})))
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unique() = %v, want %v", result, expected)
	}
}

func TestUniqueBy(t *testing.T) {
	result := ToSlice(UniqueBy(FromSlice([]int{1, -1, 2, -2, 3}), func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("UniqueBy() = %v, want %v", result, expected)
	}
}

func TestSortBy(t *testing.T) {
	result := ToSlice(SortBy(FromSlice([]int{3, 1, 4, 1, 5}), func(a, b int) bool {
		return a < b
	}))
	expected := []int{1, 1, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortBy() = %v, want %v", result, expected)
	}
}

func TestMinBy(t *testing.T) {
	min, found := MinBy(FromSlice([]int{3, 1, 4, 1, 5}), func(a, b int) bool {
		return a < b
	})
	if !found || min != 1 {
		t.Errorf("MinBy() = %d, %t, want 1, true", min, found)
	}

	// Test empty sequence
	_, found2 := MinBy(FromSlice([]int{}), func(a, b int) bool {
		return a < b
	})
	if found2 {
		t.Error("MinBy(empty) should return false")
	}
}

func TestMaxBy(t *testing.T) {
	max, found := MaxBy(FromSlice([]int{3, 1, 4, 1, 5}), func(a, b int) bool {
		return a < b
	})
	if !found || max != 5 {
		t.Errorf("MaxBy() = %d, %t, want 5, true", max, found)
	}
}

func TestToSlice(t *testing.T) {
	result := ToSlice(FromSlice([]int{1, 2, 3}))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToSlice() = %v, want %v", result, expected)
	}
}

func TestCountBy(t *testing.T) {
	result := CountBy(FromSlice([]int{1, 2, 2, 3, 1}), func(x int) int {
		return x
	})
	expected := map[int]int{1: 2, 2: 2, 3: 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("CountBy() = %v, want %v", result, expected)
	}
}

func TestPartition(t *testing.T) {
	left, right := Partition(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		return x%2 == 0
	})
	expectedLeft := []int{2, 4}
	expectedRight := []int{1, 3, 5}
	if !reflect.DeepEqual(left, expectedLeft) {
		t.Errorf("Partition left = %v, want %v", left, expectedLeft)
	}
	if !reflect.DeepEqual(right, expectedRight) {
		t.Errorf("Partition right = %v, want %v", right, expectedRight)
	}
}

func TestGroupByAdjacent(t *testing.T) {
	result := ToSlice(GroupByAdjacent(FromSlice([]int{1, 1, 2, 2, 2, 3, 1}), func(a, b int) bool {
		return a == b
	}))
	expected := [][]int{{1, 1}, {2, 2, 2}, {3}, {1}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GroupByAdjacent() = %v, want %v", result, expected)
	}
}

func TestCycle(t *testing.T) {
	result := ToSlice(Take(Cycle(FromSlice([]int{1, 2})), 5))
	expected := []int{1, 2, 1, 2, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Cycle() = %v, want %v", result, expected)
	}
}

func TestDedup(t *testing.T) {
	result := ToSlice(Dedup(FromSlice([]int{1, 1, 2, 2, 3, 1, 1})))
	expected := []int{1, 2, 3, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Dedup() = %v, want %v", result, expected)
	}
}

func TestDedupBy(t *testing.T) {
	result := ToSlice(DedupBy(FromSlice([]int{1, -1, 2, 2, -2, 3}), func(a, b int) bool {
		return (a < 0) == (b < 0)
	}))
	expected := []int{1, -1, 2, -2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DedupBy() = %v, want %v", result, expected)
	}
}

func TestIntersperse(t *testing.T) {
	result := ToSlice(Intersperse(FromSlice([]int{1, 2, 3}), 0))
	expected := []int{1, 0, 2, 0, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersperse() = %v, want %v", result, expected)
	}

	// Test empty sequence
	result2 := ToSlice(Intersperse(FromSlice([]int{}), 0))
	if len(result2) != 0 {
		t.Errorf("Intersperse(empty) = %v, want empty", result2)
	}
}

func TestFlatten(t *testing.T) {
	result := ToSlice(Flatten(FromSlice([][]int{{1, 2}, {3, 4}, {5}})))
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Flatten() = %v, want %v", result, expected)
	}
}

func TestFlattenSeq(t *testing.T) {
	seqs := []Seq[int]{
		FromSlice([]int{1, 2}),
		FromSlice([]int{3, 4}),
		FromSlice([]int{5}),
	}
	result := ToSlice(FlattenSeq(FromSlice(seqs)))
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlattenSeq() = %v, want %v", result, expected)
	}
}

func TestCombinations(t *testing.T) {
	result := ToSlice(Combinations(FromSlice([]int{1, 2, 3, 4}), 2))
	expected := [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Combinations() = %v, want %v", result, expected)
	}

	// Test k > n
	result2 := ToSlice(Combinations(FromSlice([]int{1, 2}), 5))
	if len(result2) != 0 {
		t.Errorf("Combinations(k>n) = %v, want empty", result2)
	}
}

func TestPermutations(t *testing.T) {
	result := ToSlice(Permutations(FromSlice([]int{1, 2, 3})))
	expected := [][]int{
		{1, 2, 3}, {2, 1, 3}, {3, 1, 2}, {1, 3, 2}, {2, 3, 1}, {3, 2, 1},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Permutations() = %v, want %v", result, expected)
	}
}

func TestReverse(t *testing.T) {
	result := ToSlice(FromSliceReverse([]int{1, 2, 3, 4, 5}))
	expected := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Reverse() = %v, want %v", result, expected)
	}
}

func TestCounter(t *testing.T) {
	result := Counter(FromSlice([]int{1, 2, 2, 3, 1}))
	expected := map[any]int{1: 2, 2: 2, 3: 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Counter() = %v, want %v", result, expected)
	}
}

func TestContains(t *testing.T) {
	if !Contains(FromSlice([]int{1, 2, 3}), 2) {
		t.Error("Contains() should return true when element exists")
	}

	if Contains(FromSlice([]int{1, 2, 3}), 5) {
		t.Error("Contains() should return false when element doesn't exist")
	}
}

// ========== Tests for Seq2 operations ==========

func TestForEach2(t *testing.T) {
	var result []string
	ForEach2(FromMap(map[string]int{"a": 1, "b": 2}), func(k string, v int) {
		result = append(result, fmt.Sprintf("%s:%d", k, v))
	})
	if len(result) != 2 {
		t.Errorf("ForEach2() processed %d items, want 2", len(result))
	}
}

func TestCount2(t *testing.T) {
	count := Count2(FromMap(map[string]int{"a": 1, "b": 2, "c": 3}))
	if count != 3 {
		t.Errorf("Count2() = %d, want 3", count)
	}
}

func TestRange2(t *testing.T) {
	var result []string
	Range2(Enumerate(FromSlice([]string{"x", "y", "z"}), 0), func(i int, s string) bool {
		result = append(result, fmt.Sprintf("%d:%s", i, s))
		return i < 1 // stop after first 2 elements
	})
	if len(result) != 2 {
		t.Errorf("Range2() processed %d items, want 2", len(result))
	}
}

func TestMap2(t *testing.T) {
	result := ToPairs(Map2(Enumerate(FromSlice([]string{"a", "b"}), 0), func(i int, s string) (string, int) {
		return s, i * 10
	}))
	expected := []Pair[string, int]{{"a", 0}, {"b", 10}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map2() = %v, want %v", result, expected)
	}
}

func TestFilter2(t *testing.T) {
	result := ToPairs(Filter2(Enumerate(FromSlice([]int{1, 2, 3, 4}), 0), func(_, v int) bool {
		return v%2 == 0
	}))
	expected := []Pair[int, int]{{1, 2}, {3, 4}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter2() = %v, want %v", result, expected)
	}
}

func TestExclude2(t *testing.T) {
	result := ToPairs(Exclude2(Enumerate(FromSlice([]int{1, 2, 3, 4}), 0), func(_, v int) bool {
		return v%2 == 0
	}))
	expected := []Pair[int, int]{{0, 1}, {2, 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Exclude2() = %v, want %v", result, expected)
	}
}

func TestFind2(t *testing.T) {
	k, v, found := Find2(Enumerate(FromSlice([]int{10, 20, 30}), 0), func(_, val int) bool {
		return val > 15
	})
	if !found || k != 1 || v != 20 {
		t.Errorf("Find2() = %d, %d, %t, want 1, 20, true", k, v, found)
	}

	// Test not found
	_, _, found2 := Find2(Enumerate(FromSlice([]int{1, 2, 3}), 0), func(_, val int) bool {
		return val > 10
	})
	if found2 {
		t.Error("Find2(not found) should return false")
	}
}

func TestKeys(t *testing.T) {
	result := ToSlice(Keys(Enumerate(FromSlice([]string{"a", "b", "c"}), 10)))
	expected := []int{10, 11, 12}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Keys() = %v, want %v", result, expected)
	}
}

func TestValues(t *testing.T) {
	result := ToSlice(Values(Enumerate(FromSlice([]string{"a", "b", "c"}), 0)))
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Values() = %v, want %v", result, expected)
	}
}

func TestOrderByKey(t *testing.T) {
	pairs := []Pair[int, string]{{3, "c"}, {1, "a"}, {2, "b"}}
	result := ToPairs(OrderByKey(FromPairs(pairs), func(a, b int) bool { return a < b }))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderByKey() = %v, want %v", result, expected)
	}
}

func TestOrderByValue(t *testing.T) {
	pairs := []Pair[string, int]{{"c", 3}, {"a", 1}, {"b", 2}}
	result := ToPairs(OrderByValue(FromPairs(pairs), func(a, b int) bool { return a < b }))
	expected := []Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderByValue() = %v, want %v", result, expected)
	}
}

func TestFromPairs(t *testing.T) {
	pairs := []Pair[string, int]{{"a", 1}, {"b", 2}}
	result := ToMap(FromPairs(pairs))
	expected := map[string]int{"a": 1, "b": 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromPairs() = %v, want %v", result, expected)
	}
}

func TestTake2(t *testing.T) {
	result := ToPairs(Take2(Enumerate(FromSlice([]int{1, 2, 3, 4, 5}), 0), 3))
	expected := []Pair[int, int]{{0, 1}, {1, 2}, {2, 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Take2() = %v, want %v", result, expected)
	}

	// Test take zero
	result2 := ToPairs(Take2(Enumerate(FromSlice([]int{1, 2, 3}), 0), 0))
	if len(result2) != 0 {
		t.Errorf("Take2(0) = %v, want empty", result2)
	}
}

func TestSkip2(t *testing.T) {
	result := ToPairs(Skip2(Enumerate(FromSlice([]int{1, 2, 3, 4, 5}), 0), 2))
	expected := []Pair[int, int]{{2, 3}, {3, 4}, {4, 5}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip2() = %v, want %v", result, expected)
	}

	// Test skip zero
	result2 := ToPairs(Skip2(Enumerate(FromSlice([]int{1, 2}), 0), 0))
	expected2 := []Pair[int, int]{{0, 1}, {1, 2}}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Skip2(0) = %v, want %v", result2, expected2)
	}
}

func TestStepBy2(t *testing.T) {
	result := ToPairs(StepBy2(Enumerate(FromSlice([]int{1, 2, 3, 4, 5, 6}), 0), 2))
	expected := []Pair[int, int]{{0, 1}, {2, 3}, {4, 5}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StepBy2() = %v, want %v", result, expected)
	}

	// Test invalid step
	result2 := ToPairs(StepBy2(Enumerate(FromSlice([]int{1, 2, 3}), 0), 0))
	if len(result2) != 0 {
		t.Errorf("StepBy2(0) = %v, want empty", result2)
	}
}

func TestChain2(t *testing.T) {
	s1 := Enumerate(FromSlice([]string{"a", "b"}), 0)
	s2 := Enumerate(FromSlice([]string{"c", "d"}), 10)
	result := ToPairs(Chain2(s1, s2))
	expected := []Pair[int, string]{{0, "a"}, {1, "b"}, {10, "c"}, {11, "d"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chain2() = %v, want %v", result, expected)
	}
}

func TestSortBy2(t *testing.T) {
	pairs := []Pair[string, int]{{"c", 3}, {"a", 1}, {"b", 2}}
	result := ToPairs(SortBy2(FromPairs(pairs), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}))
	expected := []Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortBy2() = %v, want %v", result, expected)
	}
}

func TestInspect2(t *testing.T) {
	var inspected []string
	result := ToPairs(Inspect2(Enumerate(FromSlice([]int{1, 2, 3}), 0), func(i, v int) {
		inspected = append(inspected, fmt.Sprintf("%d:%d", i, v))
	}))
	expected := []Pair[int, int]{{0, 1}, {1, 2}, {2, 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Inspect2() = %v, want %v", result, expected)
	}
	if len(inspected) != 3 {
		t.Errorf("Inspect2 side effect count = %d, want 3", len(inspected))
	}
}

func TestAny2(t *testing.T) {
	if !Any2(Enumerate(FromSlice([]int{1, 2, 3}), 0), func(_, v int) bool { return v > 2 }) {
		t.Error("Any2() should return true when predicate matches")
	}

	if Any2(Enumerate(FromSlice([]int{1, 2, 3}), 0), func(_, v int) bool { return v > 10 }) {
		t.Error("Any2() should return false when predicate never matches")
	}
}

func TestAll2(t *testing.T) {
	if !All2(Enumerate(FromSlice([]int{2, 4, 6}), 0), func(_, v int) bool { return v%2 == 0 }) {
		t.Error("All2() should return true when all elements match")
	}

	if All2(Enumerate(FromSlice([]int{1, 2, 3}), 0), func(_, v int) bool { return v%2 == 0 }) {
		t.Error("All2() should return false when some elements don't match")
	}
}

func TestFold2(t *testing.T) {
	sum := Fold2(Enumerate(FromSlice([]int{1, 2, 3}), 0), 0, func(acc, _, v int) int {
		return acc + v
	})
	if sum != 6 {
		t.Errorf("Fold2() = %d, want 6", sum)
	}
}

func TestReduce2(t *testing.T) {
	result, ok := Reduce2(Enumerate(FromSlice([]int{1, 2, 3}), 0), func(a, b Pair[int, int]) Pair[int, int] {
		return Pair[int, int]{a.Key, a.Value + b.Value}
	})
	if !ok || result.Value != 6 {
		t.Errorf("Reduce2() = %v, %t, want {_, 6}, true", result, ok)
	}

	// Test empty sequence
	_, ok2 := Reduce2(Enumerate(FromSlice([]int{}), 0), func(a, _ Pair[int, int]) Pair[int, int] {
		return a
	})
	if ok2 {
		t.Error("Reduce2(empty) should return false")
	}
}

func TestNth2(t *testing.T) {
	k, v, found := Nth2(Enumerate(FromSlice([]string{"a", "b", "c"}), 10), 1)
	if !found || k != 11 || v != "b" {
		t.Errorf("Nth2(1) = %d, %s, %t, want 11, b, true", k, v, found)
	}

	// Test out of bounds
	_, _, found2 := Nth2(Enumerate(FromSlice([]int{1, 2}), 0), 5)
	if found2 {
		t.Error("Nth2(out of bounds) should return false")
	}

	// Test negative index
	_, _, found3 := Nth2(Enumerate(FromSlice([]int{1, 2, 3}), 0), -1)
	if found3 {
		t.Error("Nth2(negative) should return false")
	}
}

// ========== Tests for missing functions ==========

func TestPull2(t *testing.T) {
	next, stop := Pull2(Enumerate(FromSlice([]int{1, 2, 3}), 10))
	defer stop()

	k, v, ok := next()
	if !ok || k != 10 || v != 1 {
		t.Errorf("Pull2() first = %d, %d, %t, want 10, 1, true", k, v, ok)
	}

	k, v, ok = next()
	if !ok || k != 11 || v != 2 {
		t.Errorf("Pull2() second = %d, %d, %t, want 11, 2, true", k, v, ok)
	}
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	result := ToSlice(Take(Context(FromSlice([]int{1, 2, 3, 4, 5}), ctx), 3))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Context() = %v, want %v", result, expected)
	}

	// Test with cancelled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	result2 := ToSlice(Context(FromSlice([]int{1, 2, 3}), cancelCtx))
	if len(result2) != 0 {
		t.Errorf("Context(cancelled) = %v, want empty", result2)
	}
}

func TestContext2(t *testing.T) {
	ctx := context.Background()
	result := ToPairs(Take2(Context2(Enumerate(FromSlice([]int{1, 2, 3}), 0), ctx), 2))
	expected := []Pair[int, int]{{0, 1}, {1, 2}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Context2() = %v, want %v", result, expected)
	}

	// Test with cancelled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	result2 := ToPairs(Context2(Enumerate(FromSlice([]int{1, 2, 3}), 0), cancelCtx))
	if len(result2) != 0 {
		t.Errorf("Context2(cancelled) = %v, want empty", result2)
	}
}

func TestToChan(t *testing.T) {
	ctx := context.Background()
	ch := ToChan(FromSlice([]int{1, 2, 3}), ctx)

	var result []int
	for v := range ch {
		result = append(result, v)
	}

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToChan() = %v, want %v", result, expected)
	}

	// Test with cancelled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	ch2 := ToChan(FromSlice([]int{1, 2, 3}), cancelCtx)

	var result2 []int
	for v := range ch2 {
		result2 = append(result2, v)
	}

	if len(result2) != 0 {
		t.Errorf("ToChan(cancelled) = %v, want empty", result2)
	}
}

func TestToChan2(t *testing.T) {
	ctx := context.Background()
	ch := ToChan2(Enumerate(FromSlice([]int{1, 2}), 10), ctx)

	var result []Pair[int, int]
	for kv := range ch {
		result = append(result, kv)
	}

	expected := []Pair[int, int]{{10, 1}, {11, 2}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToChan2() = %v, want %v", result, expected)
	}

	// Test with cancelled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	ch2 := ToChan2(Enumerate(FromSlice([]int{1, 2}), 0), cancelCtx)

	var result2 []Pair[int, int]
	for kv := range ch2 {
		result2 = append(result2, kv)
	}

	if len(result2) != 0 {
		t.Errorf("ToChan2(cancelled) = %v, want empty", result2)
	}
}

func TestLe(t *testing.T) {
	less := func(a, b int) bool { return a < b }

	// Test less than or equal
	if !Le(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), less) {
		t.Error("Le([1,2], [1,3]) should be true")
	}

	// Test equal sequences
	if !Le(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), less) {
		t.Error("Le(equal) should be true")
	}

	// Test greater than
	if Le(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), less) {
		t.Error("Le([1,3], [1,2]) should be false")
	}
}

func TestGt(t *testing.T) {
	less := func(a, b int) bool { return a < b }

	// Test greater than
	if !Gt(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), less) {
		t.Error("Gt([1,3], [1,2]) should be true")
	}

	// Test not greater than
	if Gt(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), less) {
		t.Error("Gt([1,2], [1,3]) should be false")
	}

	// Test equal sequences
	if Gt(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), less) {
		t.Error("Gt(equal) should be false")
	}
}

func TestGe(t *testing.T) {
	less := func(a, b int) bool { return a < b }

	// Test greater than or equal
	if !Ge(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), less) {
		t.Error("Ge([1,3], [1,2]) should be true")
	}

	// Test equal sequences
	if !Ge(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), less) {
		t.Error("Ge(equal) should be true")
	}

	// Test less than
	if Ge(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), less) {
		t.Error("Ge([1,2], [1,3]) should be false")
	}
}

// ========== Tests for additional coverage ==========

func TestSkipNegative(t *testing.T) {
	// Test negative skip (should return original sequence)
	result := ToSlice(Skip(FromSlice([]int{1, 2, 3}), -5))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip(negative) = %v, want %v", result, expected)
	}
}

func TestStepByNegative(t *testing.T) {
	// Test negative step (should return empty)
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3}), -1))
	if len(result) != 0 {
		t.Errorf("StepBy(negative) = %v, want empty", result)
	}
}

func TestChainEarlyStop(t *testing.T) {
	s1 := FromSlice([]int{1, 2})
	s2 := FromSlice([]int{3, 4})
	s3 := FromSlice([]int{5, 6})

	var result []int
	Range(Chain(s1, s2, s3), func(x int) bool {
		result = append(result, x)
		return x < 3 // Stop when we hit 3
	})

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chain early stop = %v, want %v", result, expected)
	}
}

func TestIotaNegativeStep(t *testing.T) {
	result := ToSlice(Iota(5, 1, -1))
	expected := []int{5, 4, 3, 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Iota(negative step) = %v, want %v", result, expected)
	}
}

func TestIotaInclusiveNegativeStep(t *testing.T) {
	result := ToSlice(IotaInclusive(5, 2, -2))
	expected := []int{5, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("IotaInclusive(negative step) = %v, want %v", result, expected)
	}
}

func TestFromChanEmpty(t *testing.T) {
	ch := make(chan int)
	close(ch) // Close empty channel

	result := ToSlice(FromChan(ch))
	if len(result) != 0 {
		t.Errorf("FromChan(empty) = %v, want empty", result)
	}
}

func TestEnumerateNonZeroStart(t *testing.T) {
	result := ToPairs(Enumerate(FromSlice([]string{"a", "b"}), 100))
	expected := []Pair[int, string]{{100, "a"}, {101, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Enumerate(start=100) = %v, want %v", result, expected)
	}
}

func TestZipEarlyStop(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3, 4, 5})
	s2 := FromSlice([]string{"a", "b"}) // shorter sequence

	result := ToPairs(Zip(s1, s2))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(different lengths) = %v, want %v", result, expected)
	}
}

func TestZipWithEarlyStop(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{10, 20}) // shorter sequence

	result := ToSlice(ZipWith(s1, s2, func(a, b int) int { return a + b }))
	expected := []int{11, 22}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ZipWith(different lengths) = %v, want %v", result, expected)
	}
}

func TestInterleaveExhausted(t *testing.T) {
	s1 := FromSlice([]int{1})          // short sequence
	s2 := FromSlice([]int{2, 3, 4, 5}) // longer sequence

	result := ToSlice(Interleave(s1, s2))
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Interleave(different lengths) = %v, want %v", result, expected)
	}
}

func TestChunksWithRemainder(t *testing.T) {
	result := ToSlice(Chunks(FromSlice([]int{1, 2, 3, 4, 5}), 3))
	expected := [][]int{{1, 2, 3}, {4, 5}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chunks with remainder = %v, want %v", result, expected)
	}
}

func TestIntersperseSingle(t *testing.T) {
	result := ToSlice(Intersperse(FromSlice([]int{42}), 0))
	expected := []int{42}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersperse(single) = %v, want %v", result, expected)
	}
}

func TestFlattenEmpty(t *testing.T) {
	result := ToSlice(Flatten(FromSlice([][]int{{}, {1, 2}, {}})))
	expected := []int{1, 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Flatten(with empty slices) = %v, want %v", result, expected)
	}
}

func TestCombinationsZeroK(t *testing.T) {
	result := ToSlice(Combinations(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("Combinations(k=0) = %v, want empty", result)
	}
}

func TestPermutationsEmpty(t *testing.T) {
	result := ToSlice(Permutations(FromSlice([]int{})))
	if len(result) != 0 {
		t.Errorf("Permutations(empty) = %v, want empty", result)
	}
}

func TestReverseEmpty(t *testing.T) {
	result := ToSlice(FromSliceReverse([]int{}))
	if len(result) != 0 {
		t.Errorf("Reverse(empty) = %v, want empty", result)
	}
}

func TestGroupByAdjacentEmpty(t *testing.T) {
	result := ToSlice(GroupByAdjacent(FromSlice([]int{}), func(a, b int) bool {
		return a == b
	}))
	if len(result) != 0 {
		t.Errorf("GroupByAdjacent(empty) = %v, want empty", result)
	}
}

func TestCmpDifferentLengths(t *testing.T) {
	cmp := Cmp(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp != 1 { // first is longer, so greater
		t.Errorf("Cmp(longer, shorter) = %d, want 1", cmp)
	}
}

func TestCmpSecondLonger(t *testing.T) {
	cmp := Cmp(FromSlice([]int{1, 2}), FromSlice([]int{1, 2, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp != -1 { // first is shorter, so less
		t.Errorf("Cmp(shorter, longer) = %d, want -1", cmp)
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Test cancellation during iteration
	count := 0
	seq := func(y func(int) bool) {
		for i := range 10 {
			if i == 3 {
				cancel() // Cancel after yielding 3 elements
			}
			if !y(i) {
				break
			}
			count++
		}
	}

	result := ToSlice(Context(seq, ctx))
	// Should get elements 0, 1, 2 before cancellation hits
	if len(result) > 4 {
		t.Errorf("Context cancellation should limit results, got %d elements", len(result))
	}
}

func TestContextAlreadyCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	count := 0
	seq := func(y func(int) bool) {
		for i := range 5 {
			count++
			if !y(i) {
				break
			}
		}
	}

	result := ToSlice(Context(seq, ctx))
	if len(result) != 0 {
		t.Errorf("Context with pre-cancelled context should yield empty result, got %v", result)
	}
	if count != 0 {
		t.Errorf("Context with pre-cancelled context should not execute sequence, but count = %d", count)
	}
}

func TestContext2Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	count := 0
	seq := func(y func(int, string) bool) {
		for i := range 10 {
			if i == 2 {
				cancel()
			}
			if !y(i, fmt.Sprintf("val%d", i)) {
				break
			}
			count++
		}
	}

	result := ToPairs(Context2(seq, ctx))
	if len(result) > 3 {
		t.Errorf("Context2 cancellation should limit results, got %d elements", len(result))
	}
}

func TestChainEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{4, 5, 6})
	s3 := FromSlice([]int{7, 8, 9})

	count := 0
	Range(Chain(s1, s2, s3), func(x int) bool {
		count++
		return x != 5 // Stop at element 5
	})

	if count != 5 {
		t.Errorf("Chain early termination should stop at element 5, got count %d", count)
	}
}

func TestChainWithEmptySequences(t *testing.T) {
	s1 := FromSlice([]int{1, 2})
	s2 := FromSlice([]int{}) // empty
	s3 := FromSlice([]int{3, 4})

	result := ToSlice(Chain(s1, s2, s3))
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chain with empty sequences = %v, want %v", result, expected)
	}
}

func TestChain2EarlyTermination(t *testing.T) {
	s1 := Enumerate(FromSlice([]int{1, 2, 3}), 0)
	s2 := Enumerate(FromSlice([]int{4, 5, 6}), 10)

	count := 0
	Range2(Chain2(s1, s2), func(_, _ int) bool {
		count++
		return count < 4 // Stop after 3 pairs
	})

	if count != 4 {
		t.Errorf("Chain2 early termination should stop after 4 elements, got count %d", count)
	}
}

func TestChain2WithEmptySequences(t *testing.T) {
	s1 := Enumerate(FromSlice([]int{1, 2}), 0)
	s2 := Enumerate(FromSlice([]int{}), 10) // empty
	s3 := Enumerate(FromSlice([]int{3}), 20)

	result := ToPairs(Chain2(s1, s2, s3))
	expected := []Pair[int, int]{{0, 1}, {1, 2}, {20, 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chain2 with empty sequences = %v, want %v", result, expected)
	}
}

func TestFromChanEarlyTermination(t *testing.T) {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	close(ch)

	count := 0
	Range(FromChan(ch), func(x int) bool {
		count++
		return x != 3 // Stop at element 3
	})

	if count != 3 {
		t.Errorf("FromChan early termination should stop at element 3, got count %d", count)
	}
}

func TestFromMapEarlyTermination(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	count := 0
	Range2(FromMap(m), func(_ string, _ int) bool {
		count++
		return count < 2 // Stop after 2 elements
	})

	if count != 2 {
		t.Errorf("FromMap early termination should stop after 2 elements, got count %d", count)
	}
}

func TestIotaEarlyTermination(t *testing.T) {
	count := 0
	Range(Iota(0, 10), func(x int) bool {
		count++
		return x != 5 // Stop at element 5
	})

	if count != 6 {
		t.Errorf("Iota early termination should stop after 6 elements, got count %d", count)
	}
}

func TestIotaNegativeStepEarlyTermination(t *testing.T) {
	count := 0
	Range(Iota(10, 0, -1), func(x int) bool {
		count++
		return x != 7 // Stop at element 7
	})

	if count != 4 {
		t.Errorf("Iota negative step early termination should stop after 4 elements, got count %d", count)
	}
}

func TestIotaInclusiveEarlyTermination(t *testing.T) {
	count := 0
	Range(IotaInclusive(0, 10), func(x int) bool {
		count++
		return x != 4 // Stop at element 4
	})

	if count != 5 {
		t.Errorf("IotaInclusive early termination should stop after 5 elements, got count %d", count)
	}
}

func TestIotaInclusiveNegativeStepEarlyTermination(t *testing.T) {
	count := 0
	Range(IotaInclusive(10, 0, -2), func(x int) bool {
		count++
		return x != 6 // Stop at element 6
	})

	if count != 3 {
		t.Errorf("IotaInclusive negative step early termination should stop after 3 elements, got count %d", count)
	}
}

func TestFromPairsEarlyTermination(t *testing.T) {
	pairs := []Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}, {"d", 4}}

	count := 0
	Range2(FromPairs(pairs), func(k string, _ int) bool {
		count++
		return k != "c" // Stop at key "c"
	})

	if count < 3 || count > 4 {
		t.Errorf("FromPairs early termination unexpected count %d", count)
	}
}

func TestInterleaveEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 3, 5, 7})
	s2 := FromSlice([]int{2, 4, 6, 8})

	count := 0
	Range(Interleave(s1, s2), func(x int) bool {
		count++
		return x != 4 // Stop at element 4
	})

	if count != 4 {
		t.Errorf("Interleave early termination should stop after 4 elements, got count %d", count)
	}
}

func TestToChanContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Test context already cancelled
	cancel()
	ch := ToChan(FromSlice([]int{1, 2, 3}), ctx)

	count := 0
	for range ch {
		count++
	}

	if count != 0 {
		t.Errorf("ToChan with cancelled context should yield no values, got %d", count)
	}
}

func TestToChan2ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Test context already cancelled
	cancel()
	ch := ToChan2(Enumerate(FromSlice([]int{1, 2, 3}), 0), ctx)

	count := 0
	for range ch {
		count++
	}

	if count != 0 {
		t.Errorf("ToChan2 with cancelled context should yield no values, got %d", count)
	}
}

func TestIotaZeroStep(t *testing.T) {
	result := ToSlice(Iota(1, 10, 0))
	if len(result) != 0 {
		t.Errorf("Iota with zero step should return empty, got %v", result)
	}
}

func TestIotaInclusiveZeroStep(t *testing.T) {
	result := ToSlice(IotaInclusive(1, 10, 0))
	if len(result) != 0 {
		t.Errorf("IotaInclusive with zero step should return empty, got %v", result)
	}
}

func TestZipEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3, 4})
	s2 := FromSlice([]string{"a", "b", "c", "d"})

	count := 0
	Range2(Zip(s1, s2), func(i int, _ string) bool {
		count++
		return i != 2 // Stop at element 2
	})

	if count != 2 {
		t.Errorf("Zip early termination should stop after 2 elements, got count %d", count)
	}
}

func TestZipWithEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3, 4})
	s2 := FromSlice([]int{10, 20, 30, 40})

	count := 0
	Range(ZipWith(s1, s2, func(a, b int) int { return a + b }), func(sum int) bool {
		count++
		return sum != 33 // Stop at sum 33 (3+30)
	})

	if count != 3 {
		t.Errorf("ZipWith early termination should stop after 3 elements, got count %d", count)
	}
}

func TestInterleaveRemainingElements(t *testing.T) {
	s1 := FromSlice([]int{1, 3})
	s2 := FromSlice([]int{2, 4, 6, 8}) // s2 longer than s1

	result := ToSlice(Interleave(s1, s2))
	expected := []int{1, 2, 3, 4, 6, 8}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Interleave with remaining elements = %v, want %v", result, expected)
	}
}

func TestIntersperseEarlyTermination(t *testing.T) {
	count := 0
	Range(Intersperse(FromSlice([]int{1, 2, 3, 4}), 0), func(int) bool {
		count++
		return count < 4 // Stop early
	})

	if count != 4 {
		t.Errorf("Intersperse early termination should stop after 4 elements, got count %d", count)
	}
}

func TestFlattenEarlyTermination(t *testing.T) {
	slices := [][]int{{1, 2}, {3, 4}, {5, 6}}

	count := 0
	Range(Flatten(FromSlice(slices)), func(x int) bool {
		count++
		return x != 4 // Stop at element 4
	})

	if count != 4 {
		t.Errorf("Flatten early termination should stop after 4 elements, got count %d", count)
	}
}

func TestChunksEarlyTermination(t *testing.T) {
	count := 0
	Range(Chunks(FromSlice([]int{1, 2, 3, 4, 5, 6}), 2), func(chunk []int) bool {
		count++
		return len(chunk) == 2 && chunk[1] != 4 // Stop when we see chunk [3,4]
	})

	if count < 2 {
		t.Errorf("Chunks early termination should have at least 2 chunks, got count %d", count)
	}
}

func TestGroupByAdjacentEarlyTermination(t *testing.T) {
	count := 0
	Range(
		GroupByAdjacent(FromSlice([]int{1, 1, 2, 2, 3, 3}), func(a, b int) bool { return a == b }),
		func([]int) bool {
			count++
			return count < 2 // Stop after 2 groups
		},
	)

	if count != 2 {
		t.Errorf("GroupByAdjacent early termination should stop after 2 groups, got count %d", count)
	}
}

func TestCombinationsEarlyTermination(t *testing.T) {
	count := 0
	Range(Combinations(FromSlice([]int{1, 2, 3, 4}), 2), func(combo []int) bool {
		count++
		return !reflect.DeepEqual(combo, []int{1, 3}) // Stop when we see [1,3]
	})

	if count < 2 {
		t.Errorf("Combinations early termination should have processed at least 2 combinations, got count %d", count)
	}
}

func TestPermutationsEarlyTermination(t *testing.T) {
	count := 0
	Range(Permutations(FromSlice([]int{1, 2, 3})), func(perm []int) bool {
		count++
		return !reflect.DeepEqual(perm, []int{1, 3, 2}) // Stop when we see [1,3,2]
	})

	if count < 2 {
		t.Errorf("Permutations early termination should have processed at least 2 permutations, got count %d", count)
	}
}

func TestChunksZeroSize(t *testing.T) {
	result := ToSlice(Chunks(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("Chunks with zero size should return empty, got %v", result)
	}
}

func TestWindowsZeroSize(t *testing.T) {
	result := ToSlice(Windows(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("Windows with zero size should return empty, got %v", result)
	}
}

func TestSortBy2EarlyTermination(t *testing.T) {
	pairs := []Pair[string, int]{{"b", 2}, {"a", 1}, {"c", 3}}

	count := 0
	Range2(SortBy2(FromPairs(pairs), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}), func(k string, _ int) bool {
		count++
		return k != "b" // Stop at "b"
	})

	if count < 2 {
		t.Errorf("SortBy2 early termination should process at least 2 elements, got count %d", count)
	}
}

func TestChainHeadEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{4, 5, 6})

	count := 0
	Range(Chain(s1, s2), func(x int) bool {
		count++
		return x != 2 // Stop at element 2 in head sequence
	})

	if count != 2 {
		t.Errorf("Chain head early termination should stop after 2 elements, got count %d", count)
	}
}

func TestChain2HeadEarlyTermination(t *testing.T) {
	s1 := Enumerate(FromSlice([]int{1, 2, 3}), 0)
	s2 := Enumerate(FromSlice([]int{4, 5, 6}), 10)

	count := 0
	Range2(Chain2(s1, s2), func(_, v int) bool {
		count++
		return v != 2 // Stop at value 2 in head sequence
	})

	if count != 2 {
		t.Errorf("Chain2 head early termination should stop after 2 elements, got count %d", count)
	}
}

func TestReverseEarlyTermination(t *testing.T) {
	count := 0
	Range(FromSliceReverse([]int{1, 2, 3, 4, 5}), func(x int) bool {
		count++
		return x != 3 // Stop at element 3 (in reverse: 5,4,3)
	})

	if count != 3 {
		t.Errorf("Reverse early termination should stop after 3 elements, got count %d", count)
	}
}

func TestOrderByKeyEarlyTermination(t *testing.T) {
	pairs := []Pair[string, int]{{"c", 3}, {"a", 1}, {"b", 2}}

	count := 0
	Range2(OrderByKey(FromPairs(pairs), func(a, b string) bool { return a < b }), func(k string, _ int) bool {
		count++
		return k != "b" // Stop at key "b"
	})

	if count < 2 {
		t.Errorf("OrderByKey early termination should process at least 2 elements, got count %d", count)
	}
}

func TestOrderByValueEarlyTermination(t *testing.T) {
	pairs := []Pair[string, int]{{"c", 3}, {"a", 1}, {"b", 2}}

	count := 0
	Range2(OrderByValue(FromPairs(pairs), func(a, b int) bool { return a < b }), func(_ string, v int) bool {
		count++
		return v != 2 // Stop at value 2
	})

	if count < 2 {
		t.Errorf("OrderByValue early termination should process at least 2 elements, got count %d", count)
	}
}

func TestStepByZero(t *testing.T) {
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("StepBy with zero step should return empty, got %v", result)
	}
}

func TestStepBy2Zero(t *testing.T) {
	result := ToPairs(StepBy2(Enumerate(FromSlice([]int{1, 2, 3}), 0), 0))
	if len(result) != 0 {
		t.Errorf("StepBy2 with zero step should return empty, got %v", result)
	}
}

func TestTakeZero(t *testing.T) {
	result := ToSlice(Take(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("Take with zero should return empty, got %v", result)
	}
}

func TestTake2Zero(t *testing.T) {
	result := ToPairs(Take2(Enumerate(FromSlice([]int{1, 2, 3}), 0), 0))
	if len(result) != 0 {
		t.Errorf("Take2 with zero should return empty, got %v", result)
	}
}

func TestSkipZero(t *testing.T) {
	result := ToSlice(Skip(FromSlice([]int{1, 2, 3}), 0))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip with zero = %v, want %v", result, expected)
	}
}

func TestSkip2Zero(t *testing.T) {
	result := ToPairs(Skip2(Enumerate(FromSlice([]int{1, 2, 3}), 0), 0))
	expected := []Pair[int, int]{{0, 1}, {1, 2}, {2, 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Skip2 with zero = %v, want %v", result, expected)
	}
}

func TestInterleaveFirstSequenceRemainderEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 3, 5, 7, 9}) // s1 longer than s2
	s2 := FromSlice([]int{2, 4})

	count := 0
	Range(Interleave(s1, s2), func(x int) bool {
		count++
		return x != 7 // Stop at element 7 (should be in s1 remainder)
	})

	if count < 5 {
		t.Errorf(
			"Interleave first sequence remainder early termination should process at least 5 elements, got count %d",
			count,
		)
	}
}

func TestInterleaveSecondSequenceRemainderEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 3})
	s2 := FromSlice([]int{2, 4, 6, 8, 10}) // s2 longer than s1

	count := 0
	Range(Interleave(s1, s2), func(x int) bool {
		count++
		return x != 8 // Stop at element 8 (should be in s2 remainder)
	})

	if count < 5 {
		t.Errorf(
			"Interleave second sequence remainder early termination should process at least 5 elements, got count %d",
			count,
		)
	}
}

func TestZipFirstExhaustedEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 2}) // s1 shorter
	s2 := FromSlice([]string{"a", "b", "c", "d"})

	count := 0
	Range2(Zip(s1, s2), func(_ int, s string) bool {
		count++
		return s != "b" // Stop at "b"
	})

	if count != 2 {
		t.Errorf("Zip first exhausted early termination should stop after 2 elements, got count %d", count)
	}
}

func TestToChanContextCancellationDuringIteration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a sequence that will trigger cancellation during iteration
	seq := func(y func(int) bool) {
		for i := range 10 {
			if i == 2 {
				cancel() // Cancel after yielding a couple elements
			}
			if !y(i) {
				return
			}
		}
	}

	ch := ToChan(seq, ctx)

	count := 0
	for range ch {
		count++
	}

	// Should get at least some elements before cancellation
	if count < 2 {
		t.Errorf("ToChan context cancellation during iteration should yield some elements, got %d", count)
	}
}

func TestToChan2ContextCancellationDuringIteration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a sequence that will trigger cancellation during iteration
	seq := func(y func(int, string) bool) {
		for i := range 10 {
			if i == 2 {
				cancel() // Cancel after yielding a couple elements
			}
			if !y(i, fmt.Sprintf("val%d", i)) {
				return
			}
		}
	}

	ch := ToChan2(seq, ctx)

	count := 0
	for range ch {
		count++
	}

	// Should get at least some elements before cancellation
	if count < 2 {
		t.Errorf("ToChan2 context cancellation during iteration should yield some elements, got %d", count)
	}
}

func TestInterleaveFirstSequenceEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 3, 5})
	s2 := FromSlice([]int{2, 4, 6})

	count := 0
	Range(Interleave(s1, s2), func(x int) bool {
		count++
		return x != 1 // Stop immediately at first element
	})

	if count != 1 {
		t.Errorf("Interleave first sequence early termination should stop after 1 element, got count %d", count)
	}
}

func TestInterleaveSecondSequenceEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 3, 5})
	s2 := FromSlice([]int{2, 4, 6})

	count := 0
	Range(Interleave(s1, s2), func(x int) bool {
		count++
		return x != 2 // Stop at second element (from s2)
	})

	if count != 2 {
		t.Errorf("Interleave second sequence early termination should stop after 2 elements, got count %d", count)
	}
}

func TestOrderByKeyReversedComparison(t *testing.T) {
	// Test case where less(b.K, a.K) returns true to cover the return 1 branch
	pairs := []Pair[int, string]{{3, "c"}, {1, "a"}, {2, "b"}}

	// Use reverse comparison to trigger the less(b.K, a.K) branch
	result := ToPairs(OrderByKey(FromPairs(pairs), func(a, b int) bool { return a > b }))

	// Should be sorted in descending order: 3, 2, 1
	if len(result) != 3 {
		t.Errorf("OrderByKey reverse comparison result length mismatch")
	}
	if result[0].Key != 3 || result[1].Key != 2 || result[2].Key != 1 {
		t.Errorf("OrderByKey reverse comparison not working as expected: %v", result)
	}
}

func TestOrderByValueReversedComparison(t *testing.T) {
	// Test case where less(b.V, a.V) returns true to cover the return 1 branch
	pairs := []Pair[int, string]{{1, "c"}, {2, "a"}, {3, "b"}}

	// Use reverse comparison to trigger the less(b.V, a.V) branch
	result := ToPairs(OrderByValue(FromPairs(pairs), func(a, b string) bool { return a > b }))

	// Should be sorted by values in descending order: "c", "b", "a"
	if len(result) != 3 {
		t.Errorf("OrderByValue reverse comparison result length mismatch")
	}
	if result[0].Value != "c" || result[1].Value != "b" || result[2].Value != "a" {
		t.Errorf("OrderByValue reverse comparison not working as expected: %v", result)
	}
}

func TestSortBy2ReversedComparison(t *testing.T) {
	// Test case where less(b, a) returns true to cover the return 1 branch
	pairs := []Pair[int, string]{{1, "c"}, {2, "a"}, {3, "b"}}

	// Use reverse comparison to trigger the less(b, a) branch
	result := ToPairs(SortBy2(FromPairs(pairs), func(a, b Pair[int, string]) bool {
		return a.Value > b.Value // Reverse comparison on values
	}))

	// Should be sorted by values in descending order: "c", "b", "a"
	if len(result) != 3 {
		t.Errorf("SortBy2 reverse comparison result length mismatch")
	}
	if result[0].Value != "c" || result[1].Value != "b" || result[2].Value != "a" {
		t.Errorf("SortBy2 reverse comparison not working as expected: %v", result)
	}
}

func TestZipSecondExhaustedFirst(t *testing.T) {
	// Test case where second sequence is exhausted first
	s1 := FromSlice([]int{1, 2, 3, 4, 5}) // longer sequence
	s2 := FromSlice([]string{"a", "b"})   // shorter sequence, exhausted first

	result := ToPairs(Zip(s1, s2))

	if len(result) != 2 {
		t.Errorf("Zip with second sequence exhausted first should yield 2 elements, got %d", len(result))
	}
}

func TestZipBothEmptySequences(t *testing.T) {
	// Test case with both empty sequences
	s1 := FromSlice([]int{})
	s2 := FromSlice([]string{})

	result := ToPairs(Zip(s1, s2))

	if len(result) != 0 {
		t.Errorf("Zip with both empty sequences should yield no elements, got %d", len(result))
	}
}

func TestOrderByKeyEqualElements(t *testing.T) {
	// Test case where keys are equal to cover the default case (return 0)
	pairs := []Pair[int, string]{{1, "a"}, {1, "b"}, {1, "c"}}

	// Equal elements should maintain their relative order
	result := ToPairs(OrderByKey(FromPairs(pairs), func(a, b int) bool { return a < b }))

	if len(result) != 3 {
		t.Errorf("OrderByKey with equal keys should preserve all elements")
	}
	// All keys should be 1
	for _, pair := range result {
		if pair.Key != 1 {
			t.Errorf("OrderByKey with equal keys unexpected key: %v", pair.Key)
		}
	}
}

func TestOrderByValueEqualElements(t *testing.T) {
	// Test case where values are equal to cover the default case (return 0)
	pairs := []Pair[int, string]{{1, "a"}, {2, "a"}, {3, "a"}}

	// Equal elements should maintain their relative order
	result := ToPairs(OrderByValue(FromPairs(pairs), func(a, b string) bool { return a < b }))

	if len(result) != 3 {
		t.Errorf("OrderByValue with equal values should preserve all elements")
	}
	// All values should be "a"
	for _, pair := range result {
		if pair.Value != "a" {
			t.Errorf("OrderByValue with equal values unexpected value: %v", pair.Value)
		}
	}
}

func TestSortBy2EqualElements(t *testing.T) {
	// Test case where elements are equal to cover the return 0 case
	pairs := []Pair[int, string]{{1, "a"}, {1, "a"}, {1, "a"}}

	// Equal elements should maintain their relative order
	result := ToPairs(SortBy2(FromPairs(pairs), func(a, b Pair[int, string]) bool {
		if a.Key != b.Key {
			return a.Key < b.Key
		}
		return a.Value < b.Value
	}))

	if len(result) != 3 {
		t.Errorf("SortBy2 with equal elements should preserve all elements")
	}
	// All elements should be {1, "a"}
	for _, pair := range result {
		if pair.Key != 1 || pair.Value != "a" {
			t.Errorf("SortBy2 with equal elements unexpected element: %v", pair)
		}
	}
}

// Additional tests to achieve 100% coverage

func TestCycleWithEmptySequence(t *testing.T) {
	// Test Cycle with empty sequence to cover the early return path
	result := ToSlice(Take(Cycle(Empty[int]()), 3))
	if len(result) != 0 {
		t.Errorf("Cycle with empty sequence should yield no elements, got %d", len(result))
	}
}

func TestFlattenSeqEarlyTermination(t *testing.T) {
	// Test FlattenSeq with early termination
	s1 := FromSlice([]int{1, 2})
	s2 := FromSlice([]int{3, 4})
	seqs := FromSlice([]Seq[int]{s1, s2})

	count := 0
	Range(FlattenSeq(seqs), func(int) bool {
		count++
		return count < 3 // Stop after 3 elements
	})

	if count != 3 {
		t.Errorf("FlattenSeq early termination should stop after 3 elements, got count %d", count)
	}
}

func TestStepByZeroStepEdgeCase(t *testing.T) {
	// Test StepBy with step 0 to cover the early return path
	result := ToSlice(StepBy(FromSlice([]int{1, 2, 3, 4}), 0))
	if len(result) != 0 {
		t.Errorf("StepBy with step 0 should yield no elements, got %d", len(result))
	}
}

func TestStepByEarlyTermination(t *testing.T) {
	// Test StepBy with early termination to cover the yield false path
	count := 0
	s := StepBy(FromSlice([]int{1, 2, 3, 4, 5, 6}), 2)
	s(func(int) bool {
		count++
		return count < 2 // Stop after first element
	})
	if count != 2 {
		t.Errorf("StepBy early termination should stop after 2 calls, got %d", count)
	}
}

func TestStepBy2ZeroStepEdgeCase(t *testing.T) {
	// Test StepBy2 with step 0 to cover the early return path
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToPairs(StepBy2(FromPairs(pairs), 0))
	if len(result) != 0 {
		t.Errorf("StepBy2 with step 0 should yield no elements, got %d", len(result))
	}
}

func TestStepBy2EarlyTermination(t *testing.T) {
	// Test StepBy2 with early termination to cover the yield false path
	count := 0
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {6, "f"}}
	s := StepBy2(FromPairs(pairs), 2)
	s(func(int, string) bool {
		count++
		return count < 2 // Stop after first element
	})
	if count != 2 {
		t.Errorf("StepBy2 early termination should stop after 2 calls, got %d", count)
	}
}

func TestWindowsZeroSizeEdgeCase(t *testing.T) {
	// Test Windows with size 0 to cover the early return path
	result := ToSlice(Windows(FromSlice([]int{1, 2, 3}), 0))
	if len(result) != 0 {
		t.Errorf("Windows with size 0 should yield no windows, got %d", len(result))
	}
}

func TestWindowsEarlyTermination(t *testing.T) {
	// Test Windows with early termination to cover the yield false path
	count := 0
	s := Windows(FromSlice([]int{1, 2, 3, 4, 5, 6}), 3)
	s(func([]int) bool {
		count++
		return count < 2 // Stop after first window
	})
	if count != 2 {
		t.Errorf("Windows early termination should stop after 2 calls, got %d", count)
	}
}

func TestLtWithEqualElements(t *testing.T) {
	// Test Lt with sequences where a is shorter
	s1 := FromSlice([]int{1, 2})    // shorter
	s2 := FromSlice([]int{1, 2, 3}) // longer

	result := Lt(s1, s2, func(a, b int) bool { return a < b })
	if !result {
		t.Errorf("Lt should return true when first sequence is shorter")
	}
}

func TestLeWithEqualElements(t *testing.T) {
	// Test Le with sequences where b is exhausted first
	s1 := FromSlice([]int{1, 2, 3}) // longer
	s2 := FromSlice([]int{1, 2})    // shorter

	result := Le(s1, s2, func(a, b int) bool { return a < b })
	if result {
		t.Errorf("Le should return false when second sequence is shorter")
	}
}

func TestGtWithEqualElements(t *testing.T) {
	// Test Gt with sequences where b is shorter
	s1 := FromSlice([]int{1, 2, 3}) // longer
	s2 := FromSlice([]int{1, 2})    // shorter

	result := Gt(s1, s2, func(a, b int) bool { return a < b })
	if !result {
		t.Errorf("Gt should return true when first sequence is longer")
	}
}

func TestGeWithEqualElements(t *testing.T) {
	// Test Ge with sequences where a is exhausted first
	s1 := FromSlice([]int{1, 2})    // shorter
	s2 := FromSlice([]int{1, 2, 3}) // longer

	result := Ge(s1, s2, func(a, b int) bool { return a < b })
	if result {
		t.Errorf("Ge should return false when first sequence is shorter")
	}
}
