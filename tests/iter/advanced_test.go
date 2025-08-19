package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

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

func TestFlattenSeqWithEmpty(t *testing.T) {
	seqs := []Seq[int]{
		FromSlice([]int{1, 2}),
		FromSlice([]int{}), // empty sequence
		FromSlice([]int{3, 4}),
	}
	result := ToSlice(FlattenSeq(FromSlice(seqs)))
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlattenSeq(with empty) = %v, want %v", result, expected)
	}
}

func TestFlattenSeqEarlyTermination(t *testing.T) {
	seqs := []Seq[int]{
		FromSlice([]int{1, 2, 3}),
		FromSlice([]int{4, 5, 6}),
		FromSlice([]int{7, 8, 9}),
	}

	count := 0
	Range(FlattenSeq(FromSlice(seqs)), func(x int) bool {
		count++
		return x != 5 // Stop when we reach 5
	})

	if count != 5 { // Should process 1, 2, 3, 4, 5
		t.Errorf("FlattenSeq early termination should stop after 5 elements, got count %d", count)
	}
}

func TestFlattenSeqAllEmpty(t *testing.T) {
	seqs := []Seq[int]{
		FromSlice([]int{}),
		FromSlice([]int{}),
		FromSlice([]int{}),
	}
	result := ToSlice(FlattenSeq(FromSlice(seqs)))
	if len(result) != 0 {
		t.Errorf("FlattenSeq(all empty) = %v, want empty", result)
	}
}

func TestFlattenSeqEmpty(t *testing.T) {
	// Empty sequence of sequences
	result := ToSlice(FlattenSeq(FromSlice([]Seq[int]{})))
	if len(result) != 0 {
		t.Errorf("FlattenSeq(empty seq of seqs) = %v, want empty", result)
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

func TestCycleEarlyTermination(t *testing.T) {
	count := 0
	Range(Cycle(FromSlice([]int{1, 2, 3})), func(x int) bool {
		count++
		return count < 7 // Stop after 7 elements
	})

	if count != 7 {
		t.Errorf("Cycle early termination should stop after 7 elements, got count %d", count)
	}
}

func TestCycleEmpty(t *testing.T) {
	// Cycle of empty sequence should return empty
	result := ToSlice(Take(Cycle(FromSlice([]int{})), 5))
	if len(result) != 0 {
		t.Errorf("Cycle(empty) = %v, want empty", result)
	}
}

func TestDedupEmpty(t *testing.T) {
	result := ToSlice(Dedup(FromSlice([]int{})))
	if len(result) != 0 {
		t.Errorf("Dedup(empty) = %v, want empty", result)
	}
}

func TestDedupByEmpty(t *testing.T) {
	result := ToSlice(DedupBy(FromSlice([]int{}), func(a, b int) bool { return a == b }))
	if len(result) != 0 {
		t.Errorf("DedupBy(empty) = %v, want empty", result)
	}
}
