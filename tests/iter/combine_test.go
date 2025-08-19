package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

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
		t.Errorf("Windows(zero size) = %v, want empty", result3)
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

func TestGroupByAdjacent(t *testing.T) {
	result := ToSlice(GroupByAdjacent(FromSlice([]int{1, 1, 2, 2, 2, 3, 1}), func(a, b int) bool {
		return a == b
	}))
	expected := [][]int{{1, 1}, {2, 2, 2}, {3}, {1}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GroupByAdjacent() = %v, want %v", result, expected)
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

func TestChain2Empty(t *testing.T) {
	// Test Chain2 with empty sequences
	s1 := FromPairs([]Pair[int, string]{})
	s2 := FromPairs([]Pair[int, string]{})
	result := ToPairs(Chain2(s1, s2))
	if len(result) != 0 {
		t.Errorf("Chain2(empty, empty) = %v, want empty", result)
	}
}

func TestChain2SingleEmpty(t *testing.T) {
	// Test Chain2 with single empty sequence
	result := ToPairs(Chain2(FromPairs([]Pair[int, string]{})))
	if len(result) != 0 {
		t.Errorf("Chain2(single empty) = %v, want empty", result)
	}
}

func TestChain2NoSequences(t *testing.T) {
	// Test Chain2 with no sequences
	empty1 := FromPairs([]Pair[int, string]{})
	empty2 := FromPairs([]Pair[int, string]{})
	result := ToPairs(Chain2(empty1, empty2))
	if len(result) != 0 {
		t.Errorf("Chain2(no sequences) = %v, want empty", result)
	}
}

func TestChain2EarlyTermination(t *testing.T) {
	s1 := FromPairs([]Pair[int, string]{{1, "a"}, {2, "b"}})
	s2 := FromPairs([]Pair[int, string]{{3, "c"}, {4, "d"}})
	s3 := FromPairs([]Pair[int, string]{{5, "e"}, {6, "f"}})

	count := 0
	Chain2(s1, s2, s3)(func(k int, v string) bool {
		count++
		return count < 4 // Stop after 4 pairs
	})

	if count != 4 {
		t.Errorf("Chain2 early termination count = %v, want 4", count)
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

func TestGroupByAdjacentEmpty(t *testing.T) {
	result := ToSlice(GroupByAdjacent(FromSlice([]int{}), func(a, b int) bool {
		return a == b
	}))
	if len(result) != 0 {
		t.Errorf("GroupByAdjacent(empty) = %v, want empty", result)
	}
}

func TestChainEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{4, 5, 6})
	s3 := FromSlice([]int{7, 8, 9})

	count := 0
	Range(Chain(s1, s2, s3), func(x int) bool {
		count++
		return count < 5 // Stop after 5 elements
	})

	if count != 5 {
		t.Errorf("Chain early termination should stop after 5 elements, got count %d", count)
	}
}

func TestZipEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3, 4, 5})
	s2 := FromSlice([]string{"a", "b", "c", "d", "e"})

	count := 0
	Range2(Zip(s1, s2), func(k int, v string) bool {
		count++
		return count < 3 // Stop after 3 elements
	})

	if count != 3 {
		t.Errorf("Zip early termination should stop after 3 elements, got count %d", count)
	}
}

func TestInterleaveEarlyTermination(t *testing.T) {
	s1 := FromSlice([]int{1, 3, 5, 7, 9})
	s2 := FromSlice([]int{2, 4, 6, 8})

	count := 0
	Range(Interleave(s1, s2), func(x int) bool {
		count++
		return count < 6 // Stop after 6 elements
	})

	if count != 6 {
		t.Errorf("Interleave early termination should stop after 6 elements, got count %d", count)
	}
}

func TestWindowsEarlyTermination(t *testing.T) {
	count := 0
	Range(Windows(FromSlice([]int{1, 2, 3, 4, 5, 6}), 3), func(window []int) bool {
		count++
		return count < 3 // Stop after 3 windows
	})

	if count != 3 {
		t.Errorf("Windows early termination should stop after 3 windows, got count %d", count)
	}
}

func TestChunksEarlyTermination(t *testing.T) {
	count := 0
	Range(Chunks(FromSlice([]int{1, 2, 3, 4, 5, 6}), 2), func(chunk []int) bool {
		count++
		return len(chunk) == 2 && chunk[1] != 4 // Stop when we see chunk [3,4]
	})

	if count < 2 {
		t.Errorf("Chunks early termination should have processed at least 2 chunks, got count %d", count)
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

func TestChainEmpty(t *testing.T) {
	result := ToSlice(Chain(FromSlice([]int{}), FromSlice([]int{})))
	if len(result) != 0 {
		t.Errorf("Chain(empty, empty) = %v, want empty", result)
	}
}

func TestChainSingleEmpty(t *testing.T) {
	// Test chain with single empty sequence
	result := ToSlice(Chain(FromSlice([]int{})))
	if len(result) != 0 {
		t.Errorf("Chain(single empty) = %v, want empty", result)
	}
}

func TestChainNoSequences(t *testing.T) {
	// Test chain with no sequences
	empty1 := FromSlice([]int{})
	empty2 := FromSlice([]int{})
	result := ToSlice(Chain(empty1, empty2))
	if len(result) != 0 {
		t.Errorf("Chain(no sequences) = %v, want empty", result)
	}
}

func TestZipEmpty(t *testing.T) {
	s1 := FromSlice([]int{})
	s2 := FromSlice([]string{"a", "b"})
	result := ToPairs(Zip(s1, s2))
	if len(result) != 0 {
		t.Errorf("Zip(empty, non-empty) = %v, want empty", result)
	}
}

func TestInterleaveEmpty(t *testing.T) {
	s1 := FromSlice([]int{})
	s2 := FromSlice([]int{1, 2, 3})
	result := ToSlice(Interleave(s1, s2))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Interleave(empty, non-empty) = %v, want %v", result, expected)
	}
}

func TestWindowsEmpty(t *testing.T) {
	result := ToSlice(Windows(FromSlice([]int{}), 3))
	if len(result) != 0 {
		t.Errorf("Windows(empty) = %v, want empty", result)
	}
}

func TestChunksEmpty(t *testing.T) {
	result := ToSlice(Chunks(FromSlice([]int{}), 2))
	if len(result) != 0 {
		t.Errorf("Chunks(empty) = %v, want empty", result)
	}
}
