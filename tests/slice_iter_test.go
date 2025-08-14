package g_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

// helper: compare groups with expected [][]int
func assertGroupsInt(t *testing.T, got []Slice[int], want [][]int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("groups length mismatch: got %d, want %d\n got=%v\nwant=%v", len(got), len(want), got, want)
	}
	for i := range want {
		if !reflect.DeepEqual([]int(got[i]), want[i]) {
			t.Fatalf("group %d mismatch: got %v, want %v", i, []int(got[i]), want[i])
		}
	}
}

func TestSliceIterGroupBy_EqualRuns(t *testing.T) {
	in := SliceOf(1, 1, 1, 3, 3, 2, 2, 2)
	got := in.Iter().GroupBy(func(a, b int) bool { return a == b }).Collect()

	want := [][]int{
		{1, 1, 1},
		{3, 3},
		{2, 2, 2},
	}
	assertGroupsInt(t, got, want)
}

func TestSliceIterGroupBy_LessEqRuns(t *testing.T) {
	in := SliceOf(1, 1, 2, 3, 2, 3, 2, 3, 4)
	got := in.Iter().GroupBy(func(a, b int) bool { return a <= b }).Collect()

	want := [][]int{
		{1, 1, 2, 3},
		{2, 3},
		{2, 3, 4},
	}
	assertGroupsInt(t, got, want)
}

func TestSliceIterGroupBy_Empty(t *testing.T) {
	in := Slice[int]{}
	got := in.Iter().GroupBy(func(a, b int) bool { return a == b }).Collect()
	if len(got) != 0 {
		t.Fatalf("expected 0 groups, got %d: %v", len(got), got)
	}
}

func TestSliceIterGroupBy_AlwaysTrue(t *testing.T) {
	in := SliceOf(7, 8, 9)
	got := in.Iter().GroupBy(func(a, b int) bool { return true }).Collect()
	want := [][]int{{7, 8, 9}}
	assertGroupsInt(t, got, want)
}

func TestSliceIterGroupBy_AlwaysFalse(t *testing.T) {
	in := SliceOf(7, 8, 9)
	got := in.Iter().GroupBy(func(a, b int) bool { return false }).Collect()
	want := [][]int{{7}, {8}, {9}}
	assertGroupsInt(t, got, want)
}

func TestSliceIterGroupBy_GroupsAreCopies(t *testing.T) {
	in := SliceOf(1, 1, 2, 2)
	orig := append([]int(nil), []int(in)...)

	groups := in.Iter().GroupBy(func(a, b int) bool { return a == b }).Collect()
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d: %v", len(groups), groups)
	}

	// mutate first group
	g := groups[0]
	g[0] = 99

	// source must be intact
	if !reflect.DeepEqual([]int(in), orig) {
		t.Fatalf("source mutated through group: src=%v, want=%v", []int(in), orig)
	}

	// and groups must not share backing with each other
	before := append([]int(nil), []int(groups[1])...)
	groups[0][0] = 100
	if !reflect.DeepEqual([]int(groups[1]), before) {
		t.Fatalf("groups share backing array unexpectedly: g1=%v, g2(before)=%v g2(after)=%v",
			[]int(groups[0]), before, []int(groups[1]))
	}
}

func TestSliceIterGroupBy_Strings(t *testing.T) {
	in := SliceOf("a", "a", "b", "bb", "bb", "a")
	got := in.Iter().GroupBy(func(a, b string) bool { return len(a) == len(b) }).Collect()

	want := [][]string{
		{"a", "a", "b"},
		{"bb", "bb"},
		{"a"},
	}

	if len(got) != len(want) {
		t.Fatalf("groups length mismatch: got %d, want %d\n got=%v\nwant=%v", len(got), len(want), got, want)
	}
	for i := range want {
		if !reflect.DeepEqual([]string(got[i]), want[i]) {
			t.Fatalf("group %d mismatch: got %v, want %v", i, []string(got[i]), want[i])
		}
	}
}

func TestSliceIntoIter(t *testing.T) {
	s := SliceOf(1, 2, 3, 4, 5)

	if len(s) != 5 {
		t.Fatalf("expected slice to have 5 elements, got %d", len(s))
	}

	iter := s.IntoIter()

	if len(s) != 0 {
		t.Fatalf("expected slice to be empty after IntoIter, got %d elements", len(s))
	}

	result := iter.Map(func(x int) int { return x * 10 }).Collect()

	expected := SliceOf(10, 20, 30, 40, 50)
	if !result.Eq(expected) {
		t.Errorf("expected result %v, got %v", expected, result)
	}
}

func TestSliceIterFromChan(t *testing.T) {
	// Create a channel and populate it with some test data
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 5; i++ {
			ch <- i
		}
	}()

	// Convert the channel into an iterator
	iter := FromChan(ch)

	// Create a slice to collect elements from the iterator
	var collected []int

	// Define a function to be used as a callback for iterator
	yield := func(v int) bool {
		if v == 3 {
			return false // Return false when element equals 3 to test premature exit
		}
		collected = append(collected, v)
		return true
	}

	// Iterate through the elements using the iterator and collect them
	iter(yield)

	// Define the expected result
	expected := []int{1, 2}

	// Compare the collected elements with the expected result
	if len(collected) != len(expected) {
		t.Errorf("Length mismatch: expected %d elements, got %d", len(expected), len(collected))
		return
	}

	for i, v := range collected {
		if v != expected[i] {
			t.Errorf("Element mismatch at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestSliceIterPartition(t *testing.T) {
	// Test case 1: Basic partitioning with integers
	slice1 := Slice[int]{1, 2, 3, 4, 5}
	isEven := func(val int) bool {
		return val%2 == 0
	}

	evens, odds := slice1.Iter().Partition(isEven)
	expectedEvens := Slice[int]{2, 4}
	expectedOdds := Slice[int]{1, 3, 5}

	if !reflect.DeepEqual(evens, expectedEvens) {
		t.Errorf("Expected evens %v, but got %v", expectedEvens, evens)
	}

	if !reflect.DeepEqual(odds, expectedOdds) {
		t.Errorf("Expected odds %v, but got %v", expectedOdds, odds)
	}

	// Test case 2: Partitioning with strings
	slice2 := Slice[string]{"apple", "banana", "cherry", "date"}
	hasA := func(val string) bool {
		return strings.Contains(val, "a")
	}

	withA, withoutA := slice2.Iter().Partition(hasA)
	expectedWithA := Slice[string]{"apple", "banana", "date"}
	expectedWithoutA := Slice[string]{"cherry"}

	if !reflect.DeepEqual(withA, expectedWithA) {
		t.Errorf("Expected withA %v, but got %v", expectedWithA, withA)
	}

	if !reflect.DeepEqual(withoutA, expectedWithoutA) {
		t.Errorf("Expected withoutA %v, but got %v", expectedWithoutA, withoutA)
	}

	// Test case 3: Partitioning an empty slice
	emptySlice := Slice[int]{}
	all, none := emptySlice.Iter().Partition(func(_ int) bool { return true })

	if len(all) != 0 {
		t.Errorf("Expected empty slice for 'all', but got %v", all)
	}

	if len(none) != 0 {
		t.Errorf("Expected empty slice for 'none', but got %v", none)
	}
}

func TestSliceIterCombinations(t *testing.T) {
	// Test case 1: Combinations of integers
	slice1 := Slice[int]{0, 1, 2, 3}
	combs1 := slice1.Iter().Combinations(3).Collect()
	expectedCombs1 := []Slice[int]{
		{0, 1, 2},
		{0, 1, 3},
		{0, 2, 3},
		{1, 2, 3},
	}

	if !reflect.DeepEqual(combs1, expectedCombs1) {
		t.Errorf("Test case 1 failed: expected %v, but got %v", expectedCombs1, combs1)
	}

	// Test case 2: Combinations of strings
	p1 := SliceOf[String]("a", "b")
	p2 := SliceOf[String]("c", "d")
	combs2 := p1.Iter().Chain(p2.Iter()).Map(String.Upper).Combinations(2).Collect()
	expectedCombs2 := []Slice[String]{
		{"A", "B"},
		{"A", "C"},
		{"A", "D"},
		{"B", "C"},
		{"B", "D"},
		{"C", "D"},
	}

	if !reflect.DeepEqual(combs2, expectedCombs2) {
		t.Errorf("Test case 2 failed: expected %v, but got %v", expectedCombs2, combs2)
	}

	// Test case 3: Combinations of mixed types
	p3 := SliceOf[any]("x", "y")
	p4 := SliceOf[any](1, 2)
	combs3 := p3.Iter().Chain(p4.Iter()).Combinations(2).Collect()
	expectedCombs3 := []Slice[any]{
		{"x", "y"},
		{"x", 1},
		{"x", 2},
		{"y", 1},
		{"y", 2},
		{1, 2},
	}

	if !reflect.DeepEqual(combs3, expectedCombs3) {
		t.Errorf("Test case 3 failed: expected %v, but got %v", expectedCombs3, combs3)
	}

	// Test case 4: Empty slice
	emptySlice := Slice[int]{}
	combs4 := emptySlice.Iter().Combinations(2).Collect()
	expectedCombs4 := []Slice[int]{}

	if !reflect.DeepEqual(combs4, expectedCombs4) {
		t.Errorf("Test case 4 failed: expected %v, but got %v", expectedCombs4, combs4)
	}

	// Test case 5: Combinations with k greater than slice length
	slice5 := Slice[int]{1, 2, 3}
	combs5 := slice5.Iter().Combinations(4).Collect()
	expectedCombs5 := []Slice[int]{}

	if !reflect.DeepEqual(combs5, expectedCombs5) {
		t.Errorf("Test case 5 failed: expected %v, but got %v", expectedCombs5, combs5)
	}
}

func TestSliceIterSortBy(t *testing.T) {
	sl1 := NewSlice[int]()
	sl1.Push(3, 1, 4, 1, 5)

	expected1 := NewSlice[int]()
	expected1.Push(1, 1, 3, 4, 5)

	actual1 := sl1.Iter().SortBy(cmp.Cmp).Collect()

	if !actual1.Eq(expected1) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected1, actual1)
	}

	sl2 := NewSlice[string]()
	sl2.Push("foo", "bar", "baz")
	expected2 := NewSlice[string]()
	expected2.Push("foo", "baz", "bar")

	actual2 := sl2.Iter().SortBy(func(a, b string) cmp.Ordering { return cmp.Cmp(b, a) }).Collect()

	if !actual2.Eq(expected2) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected2, actual2)
	}

	sl3 := NewSlice[int]()
	expected3 := NewSlice[int]()

	actual3 := sl3.Iter().SortBy(cmp.Cmp).Collect()

	if !actual3.Eq(expected3) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected3, actual3)
	}
}

func TestSliceIterDedup(t *testing.T) {
	// Test case 1: Dedup with consecutive duplicate elements for int
	sliceInt := Slice[int]{1, 2, 2, 3, 4, 4, 4, 5}
	expectedResultInt := Slice[int]{1, 2, 3, 4, 5}

	iterInt := sliceInt.Iter().Dedup()
	resultInt := iterInt.Collect()

	if !reflect.DeepEqual(resultInt, expectedResultInt) {
		t.Errorf("Dedup failed for int. Expected %v, got %v", expectedResultInt, resultInt)
	}

	// Test case 2: Dedup with consecutive duplicate elements for string
	sliceString := Slice[string]{"apple", "orange", "orange", "banana", "banana", "grape"}
	expectedResultString := Slice[string]{"apple", "orange", "banana", "grape"}

	iterString := sliceString.Iter().Dedup()
	resultString := iterString.Collect()

	if !reflect.DeepEqual(resultString, expectedResultString) {
		t.Errorf("Dedup failed for strin Expected %v, got %v", expectedResultString, resultString)
	}

	// Test case 3: Dedup with consecutive duplicate elements for float64
	sliceFloat64 := Slice[float64]{1.2, 2.3, 2.3, 3.4, 4.5, 4.5, 4.5, 5.6}
	expectedResultFloat64 := Slice[float64]{1.2, 2.3, 3.4, 4.5, 5.6}

	iterFloat64 := sliceFloat64.Iter().Dedup()
	resultFloat64 := iterFloat64.Collect()

	if !reflect.DeepEqual(resultFloat64, expectedResultFloat64) {
		t.Errorf("Dedup failed for float64. Expected %v, got %v", expectedResultFloat64, resultFloat64)
	}

	// Test case 4: Dedup with consecutive duplicate elements for custom non-comparable struct
	type myStruct struct {
		val []int
	}

	sliceStruct := Slice[myStruct]{
		{val: []int{1}},
		{val: []int{2}},
		{val: []int{2}},
		{val: []int{3}},
		{val: []int{3}},
		{val: []int{4}},
	}

	expectedResultStruct := Slice[myStruct]{{val: []int{1}}, {val: []int{2}}, {val: []int{3}}, {val: []int{4}}}

	iterStruct := sliceStruct.Iter().Dedup()
	resultStruct := iterStruct.Collect()

	if !reflect.DeepEqual(resultStruct, expectedResultStruct) {
		t.Errorf("Dedup failed for custom struct. Expected %v, got %v", expectedResultStruct, resultStruct)
	}
}

func TestSliceIterStepBy(t *testing.T) {
	// Test case 1: StepBy with a step size of 3
	slice := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedResult := Slice[int]{1, 4, 7, 10}

	iter := slice.Iter().StepBy(3)
	result := iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 2: StepBy with a step size of 2
	slice = Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedResult = Slice[int]{1, 3, 5, 7, 9}

	iter = slice.Iter().StepBy(2)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 3: StepBy with a step size larger than the slice length
	slice = Slice[int]{1, 2, 3, 4, 5}
	expectedResult = Slice[int]{1}

	iter = slice.Iter().StepBy(10)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 4: StepBy with a step size of 1
	slice = Slice[int]{1, 2, 3, 4, 5}
	expectedResult = Slice[int]{1, 2, 3, 4, 5}

	iter = slice.Iter().StepBy(1)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}
}

func TestSliceIterPermutations(t *testing.T) {
	// Test case 1: Single element slice
	slice1 := SliceOf(1)
	perms1 := slice1.Iter().Permutations().Collect()
	expectedPerms1 := []Slice[int]{slice1}

	if !reflect.DeepEqual(perms1, expectedPerms1) {
		t.Errorf("expected %v, but got %v", expectedPerms1, perms1)
	}

	// Test case 2: Two-element string slice
	slice2 := SliceOf("a", "b")
	perms2 := slice2.Iter().Permutations().Collect()
	expectedPerms2 := []Slice[string]{
		{"a", "b"},
		{"b", "a"},
	}

	if !reflect.DeepEqual(perms2, expectedPerms2) {
		t.Errorf("expected %v, but got %v", expectedPerms2, perms2)
	}

	// Test case 3: Three-element float64 slice
	slice3 := SliceOf(1.0, 2.0, 3.0)
	perms3 := slice3.Iter().Permutations().Collect()
	expectedPerms3 := []Slice[float64]{
		{1.0, 2.0, 3.0},
		{1.0, 3.0, 2.0},
		{2.0, 1.0, 3.0},
		{2.0, 3.0, 1.0},
		{3.0, 1.0, 2.0},
		{3.0, 2.0, 1.0},
	}

	if !reflect.DeepEqual(perms3, expectedPerms3) {
		t.Errorf("expected %v, but got %v", expectedPerms3, perms3)
	}

	// Additional Test case 4: Empty slice
	slice4 := Slice[any]{}
	perms4 := slice4.Iter().Permutations().Collect()
	expectedPerms4 := []Slice[any]{}

	if !reflect.DeepEqual(perms4, expectedPerms4) {
		t.Errorf("expected %v, but got %v", expectedPerms4, perms4)
	}

	// Additional Test case 5: Four-element mixed-type slice
	slice5 := SliceOf[any]("a", 1, 2.5, true)
	perms5 := slice5.Iter().Permutations().Collect()
	expectedPerms5 := []Slice[any]{
		{"a", 1, 2.5, true},
		{"a", 1, true, 2.5},
		{"a", 2.5, 1, true},
		{"a", 2.5, true, 1},
		{"a", true, 1, 2.5},
		{"a", true, 2.5, 1},
		{1, "a", 2.5, true},
		{1, "a", true, 2.5},
		{1, 2.5, "a", true},
		{1, 2.5, true, "a"},
		{1, true, "a", 2.5},
		{1, true, 2.5, "a"},
		{2.5, "a", 1, true},
		{2.5, "a", true, 1},
		{2.5, 1, "a", true},
		{2.5, 1, true, "a"},
		{2.5, true, "a", 1},
		{2.5, true, 1, "a"},
		{true, "a", 1, 2.5},
		{true, "a", 2.5, 1},
		{true, 1, "a", 2.5},
		{true, 1, 2.5, "a"},
		{true, 2.5, "a", 1},
		{true, 2.5, 1, "a"},
	}

	if !reflect.DeepEqual(perms5, expectedPerms5) {
		t.Errorf("expected %v, but got %v", expectedPerms5, perms5)
	}
}

func TestSliceIterChunks(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice[int]
		expected []Slice[int]
		size     Int
	}{
		{
			name:     "empty slice",
			input:    NewSlice[int](),
			expected: []Slice[int]{},
			size:     2,
		},
		{
			name:     "single chunk",
			input:    Slice[int]{1, 2, 3},
			expected: []Slice[int]{{1, 2, 3}},
			size:     3,
		},
		{
			name:  "multiple chunks",
			input: Slice[int]{1, 2, 3, 4, 5, 6},
			expected: []Slice[int]{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			size: 2,
		},
		{
			name:  "last chunk is smaller",
			input: Slice[int]{1, 2, 3, 4, 5},
			expected: []Slice[int]{
				{1, 2},
				{3, 4},
				{5},
			},
			size: 2,
		},
		{
			name:     "chunk size bigger than slice length",
			input:    Slice[int]{1, 2, 3, 4},
			expected: []Slice[int]{{1, 2, 3, 4}},
			size:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Iter().Chunks(tt.size).Collect()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %d, but got %d", tt.expected, result)
				return
			}

			for i, chunk := range result {
				if !chunk.Eq(tt.expected[i]) {
					t.Errorf("Chunk %d does not match the expected result", i)
				}
			}
		})
	}
}

func TestSliceIterAll(t *testing.T) {
	sl1 := NewSlice[int]()
	sl2 := NewSlice[int]()
	sl2.Push(1, 2, 3)
	sl3 := NewSlice[int]()
	sl3.Push(2, 4, 6)

	testCases := []struct {
		f    func(int) bool
		name string
		sl   Slice[int]
		want bool
	}{
		{
			name: "empty slice",
			f:    func(x int) bool { return x%2 == 0 },
			sl:   sl1,
			want: true,
		},
		{
			name: "all elements satisfy the condition",
			f:    func(x int) bool { return x%2 != 0 },
			sl:   sl2,
			want: false,
		},
		{
			name: "not all elements satisfy the condition",
			f:    func(x int) bool { return x%2 == 0 },
			sl:   sl3,
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.sl.Iter().All(tc.f)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSliceIterAny(t *testing.T) {
	sl1 := NewSlice[int]()
	f1 := func(x int) bool { return x > 0 }

	if sl1.Iter().Any(f1) {
		t.Errorf("Expected false for empty slice, got true")
	}

	sl2 := NewSlice[int]()
	sl2.Push(1, 2, 3)
	f2 := func(x int) bool { return x < 1 }

	if sl2.Iter().Any(f2) {
		t.Errorf("Expected false for slice with no matching elements, got true")
	}

	sl3 := NewSlice[string]()
	sl3.Push("foo", "bar")
	f3 := func(x string) bool { return x == "bar" }

	if !sl3.Iter().Any(f3) {
		t.Errorf("Expected true for slice with one matching element, got false")
	}

	sl4 := NewSlice[int]()
	sl4.Push(1, 2, 3, 4, 5)
	f4 := func(x int) bool { return x%2 == 0 }

	if !sl4.Iter().Any(f4) {
		t.Errorf("Expected true for slice with multiple matching elements, got false")
	}
}

func TestSliceIterFold(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	sum := sl.Iter().Fold(0, func(index, value int) int { return index + value })

	if sum != 15 {
		t.Errorf("Expected %d, got %d", 15, sum)
	}
}

func TestSeqSliceReduce_Sum(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}

	got := sl.Iter().Reduce(func(a, b int) int { return a + b })
	if got.IsNone() {
		t.Fatalf("Reduce should return Some for non-empty sequence")
	}

	if v := got.Some(); v != 15 {
		t.Errorf("expected 15, got %d", v)
	}
}

func TestSeqSliceReduce_Empty(t *testing.T) {
	var sl Slice[int]

	got := sl.Iter().Reduce(func(a, b int) int { return a + b })
	if got.IsSome() {
		t.Fatalf("Reduce should return None for empty sequence")
	}

	if v := got.UnwrapOr(-1); v != -1 {
		t.Errorf("expected UnwrapOr(-1) == -1, got %d", v)
	}
}

func TestSeqSliceReduce_Single(t *testing.T) {
	sl := Slice[int]{42}

	got := sl.Iter().Reduce(func(a, b int) int { return a + b })
	if got.IsNone() {
		t.Fatalf("Reduce should return Some for single-element sequence")
	}

	if v := got.Unwrap(); v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestSliceIterFilter(t *testing.T) {
	var sl Slice[int]

	sl.Push(1, 2, 3, 4, 5)
	result := sl.Iter().Filter(func(v int) bool { return v%2 == 0 }).Collect()

	if result.Len() != 2 {
		t.Errorf("Expected 2, got %d", result.Len())
	}

	if result[0] != 2 {
		t.Errorf("Expected 2, got %d", result[0])
	}

	if result[1] != 4 {
		t.Errorf("Expected 4, got %d", result[1])
	}
}

func TestSliceIterMap(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	result := sl.Iter().Map(func(i int) int { return i * 2 }).Collect()

	if result.Len() != sl.Len() {
		t.Errorf("Expected %d, got %d", sl.Len(), result.Len())
	}

	for i := range result.Len() {
		if result[i] != sl[i]*2 {
			t.Errorf("Expected %d, got %d", sl[i]*2, result[i])
		}
	}
}

func TestSliceIterExcludeZeroValues(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10}
	sl = sl.Iter().Exclude(f.IsZero).Collect()

	if sl.Len() != 10 {
		t.Errorf("Expected 10, got %d", sl.Len())
	}

	for i := range sl.Len() {
		if sl[i] == 0 {
			t.Errorf("Expected non-zero value, got %d", sl[i])
		}
	}
}

func TestSliceIterForEach(t *testing.T) {
	sl1 := NewSlice[int]()
	sl1.Push(1, 2, 3, 4, 5)
	sl2 := NewSlice[string]()
	sl2.Push("foo", "bar", "baz")
	sl3 := NewSlice[float64]()
	sl3.Push(1.1, 2.2, 3.3, 4.4)

	var result1 []int

	sl1.Iter().ForEach(func(i int) { result1 = append(result1, i) })

	if !reflect.DeepEqual(result1, []int{1, 2, 3, 4, 5}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			sl1,
			[]int{1, 2, 3, 4, 5},
			result1,
		)
	}

	var result2 []string

	sl2.Iter().ForEach(func(s string) { result2 = append(result2, s) })

	if !reflect.DeepEqual(result2, []string{"foo", "bar", "baz"}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			sl2,
			[]string{"foo", "bar", "baz"},
			result2,
		)
	}

	var result3 []float64

	sl3.Iter().ForEach(func(f float64) { result3 = append(result3, f) })

	if !reflect.DeepEqual(result3, []float64{1.1, 2.2, 3.3, 4.4}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			sl3,
			[]float64{1.1, 2.2, 3.3, 4.4},
			result3,
		)
	}
}

func TestSliceIterZip(t *testing.T) {
	s1 := SliceOf(1, 2, 3, 4)
	s2 := SliceOf(5, 6, 7, 8)
	expected := MapOrd[int, int]{{1, 5}, {2, 6}, {3, 7}, {4, 8}}
	result := s1.Iter().Zip(s2.Iter()).Collect()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s1, s2, result, expected)
	}

	s3 := SliceOf(1, 2, 3)
	s4 := SliceOf(4, 5)
	expected = MapOrd[int, int]{{1, 4}, {2, 5}}
	result = s3.Iter().Zip(s4.Iter()).Collect()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s3, s4, result, expected)
	}
}

func TestSliceIterFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice[any]
		expected Slice[any]
	}{
		{
			name:     "Empty slice",
			input:    Slice[any]{},
			expected: Slice[any]{},
		},
		{
			name:     "Flat slice",
			input:    Slice[any]{1, "abc", 3.14},
			expected: Slice[any]{1, "abc", 3.14},
		},
		{
			name: "Nested slice",
			input: Slice[any]{
				1,
				SliceOf(2, 3),
				"abc",
				SliceOf("def", "ghi"),
				SliceOf(4.5, 6.7),
			},
			expected: Slice[any]{1, 2, 3, "abc", "def", "ghi", 4.5, 6.7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Iter().Flatten().Collect()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Flatten() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSliceIterRange(t *testing.T) {
	// Test scenario: Function stops at a specific value
	t.Run("FunctionStopsAtThree", func(t *testing.T) {
		slice := Slice[int]{1, 2, 3, 4, 5}
		expected := []int{1, 2, 3}

		var result []int
		stopAtThree := func(val int) bool {
			result = append(result, val)
			return val != 3
		}

		slice.Iter().Range(stopAtThree)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})

	// Test scenario: Function always returns true
	t.Run("FunctionAlwaysTrue", func(t *testing.T) {
		slice := Slice[int]{1, 2, 3, 4, 5}
		expected := []int{1, 2, 3, 4, 5}

		var result []int
		alwaysTrue := func(val int) bool {
			result = append(result, val)
			return true
		}

		slice.Iter().Range(alwaysTrue)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})

	// Test scenario: Empty slice
	t.Run("EmptySlice", func(t *testing.T) {
		emptySlice := Slice[int]{}
		expected := []int{}

		result := []int{}
		anyFunc := func(val int) bool {
			result = append(result, val)
			return true
		}

		emptySlice.Iter().Range(anyFunc)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})
}

func TestSliceIterCount(t *testing.T) {
	// Test case 1: Count elements from the sequence
	seq := Slice[int]{1, 2, 3, 4, 5}
	count := seq.Iter().Count()
	if count != 5 {
		t.Errorf("Expected count to be %d, got %d", 5, count)
	}

	// Test case 2: Empty sequence
	emptySeq := Slice[int]{}
	emptyCount := emptySeq.Iter().Count()
	if emptyCount != 0 {
		t.Errorf("Expected count of an empty sequence to be %d, got %d", 0, emptyCount)
	}
}

func TestSliceIterCycle(t *testing.T) {
	// Test case 1: Cyclic behavior
	seq := Slice[int]{1, 2, 3}
	cycle := seq.Iter().Cycle().Take(9).Collect()

	expected := []int{1, 2, 3, 1, 2, 3, 1, 2, 3}
	for i := 0; i < len(expected); i++ {
		if cycle[i] != expected[i] {
			t.Errorf("Expected element at index %d to be %d, got %d", i, expected[i], cycle[i])
		}
	}
}

func TestSliceIterEnumerate(t *testing.T) {
	// Test case 1: Enumerate elements
	seq := Slice[string]{"bbb", "ddd", "xxx", "aaa", "ccc"}
	enumerated := seq.Iter().Enumerate().Collect()

	expected := NewMapOrd[Int, string]()
	expected.Set(0, "bbb")
	expected.Set(1, "ddd")
	expected.Set(2, "xxx")
	expected.Set(3, "aaa")
	expected.Set(4, "ccc")

	for i, v := range enumerated {
		if expected[i] != v {
			t.Errorf("Expected element at index %d to be %v, got %v", i, expected[i], v)
		}
	}
}

func TestSliceIterSkip(t *testing.T) {
	// Test case 1: Skip elements
	seq := Slice[int]{1, 2, 3, 4, 5, 6}
	skipped := seq.Iter().Skip(3).Collect()
	expected := Slice[int]{4, 5, 6}
	if len(skipped) != len(expected) {
		t.Errorf("Expected skipped slice to have length %d, got %d", len(expected), len(skipped))
	}
	for i := range expected {
		if skipped[i] != expected[i] {
			t.Errorf("Expected element at index %d to be %d, got %d", i, expected[i], skipped[i])
		}
	}

	// Test case 2: Skip all elements
	seq2 := Slice[string]{"a", "b", "c"}
	skipped2 := seq2.Iter().Skip(3).Collect()
	if len(skipped2) != 0 {
		t.Errorf("Expected skipped slice of all elements to be empty, got length %d", len(skipped2))
	}
}

func TestSliceIterUnique(t *testing.T) {
	// Test case 1: Unique elements
	seq := Slice[int]{1, 2, 3, 2, 4, 5, 3}
	unique := seq.Iter().Unique().Collect()

	expected := Slice[int]{1, 2, 3, 4, 5}
	if len(unique) != len(expected) {
		t.Errorf("Expected unique iterator length to be %d, got %d", len(expected), len(unique))
	}
	for i, v := range unique {
		if v != expected[i] {
			t.Errorf("Expected element at index %d to be %d, got %d", i, expected[i], v)
		}
	}
}

func TestSliceIterFind(t *testing.T) {
	// Test case 1: Element found
	seq := Slice[int]{1, 2, 3, 4, 5}
	found := seq.Iter().Find(func(i int) bool {
		return i == 2
	})
	if !found.IsSome() {
		t.Error("Expected found option to be Some")
	}
	if found.Some() != 2 {
		t.Errorf("Expected found element to be 2, got %d", found.Some())
	}

	// Test case 2: Element not found
	notFound := seq.Iter().Find(func(i int) bool {
		return i == 6
	})
	if notFound.IsSome() {
		t.Error("Expected not found option to be None")
	}
}

func TestSliceIterWindows(t *testing.T) {
	// Test case 1: Windows of correct size
	seq := Slice[int]{1, 2, 3, 4, 5, 6}
	windows := seq.Iter().Windows(3).Collect()

	expected := []Slice[int]{
		{1, 2, 3},
		{2, 3, 4},
		{3, 4, 5},
		{4, 5, 6},
	}

	if len(windows) != len(expected) {
		t.Errorf("Expected %d windows, got %d", len(expected), len(windows))
	}

	for i, window := range windows {
		if len(window) != len(expected[i]) {
			t.Errorf("Expected window %d length to be %d, got %d", i, len(expected[i]), len(window))
		}
		for j, v := range window {
			if v != expected[i][j] {
				t.Errorf("Expected window %d element at index %d to be %d, got %d", i, j, expected[i][j], v)
			}
		}
	}
}

func TestSliceIterToChannel(t *testing.T) {
	// Test case 1: Channel streaming without cancellation
	seq := Slice[int]{1, 2, 3}

	ch := seq.Iter().ToChan()
	var result []int
	for val := range ch {
		result = append(result, val)
	}

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected element at index %d to be %d, got %d", i, expected[i], v)
		}
	}

	// Test case 2: Channel streaming with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Ensure cancellation to avoid goroutine leaks.

	ch = seq.Iter().ToChan(ctx)
	var result2 []int
	for val := range ch {
		result2 = append(result2, val)
	}

	if len(result2) != 0 {
		t.Error("Expected no elements due to cancellation, got some elements")
	}
}

func TestSliceIterInspect(t *testing.T) {
	// Define a slice to iterate over
	s := Slice[int]{1, 2, 3}

	// Define a slice to store the inspected elements
	var inspectedElements Slice[int]

	// Create a new iterator with Inspect and collect the elements
	s.Iter().Inspect(func(v int) {
		inspectedElements = append(inspectedElements, v)
	}).Collect()

	if inspectedElements.Len() != s.Len() {
		t.Errorf("Expected %d inspected elements, got %d", s.Len(), inspectedElements.Len())
	}

	if inspectedElements.Ne(s) {
		t.Errorf("Expected %v, got %v", s, inspectedElements)
	}
}

func TestSliceIterCounter(t *testing.T) {
	sl1 := Slice[int]{1, 2, 3, 2, 1, 4, 5, 4, 4}
	sl2 := Slice[string]{"apple", "banana", "orange", "apple", "apple", "orange", "grape"}

	expected1 := NewMapOrd[int, Int]()
	expected1.Set(3, 1)
	expected1.Set(5, 1)
	expected1.Set(1, 2)
	expected1.Set(2, 2)
	expected1.Set(4, 3)
	expected1.SortByKey(cmp.Cmp)

	result1 := sl1.Iter().Counter().Collect()
	result1.SortByKey(cmp.Cmp)
	if !result1.Eq(expected1) {
		t.Errorf("Counter() returned %v, expected %v", result1, expected1)
	}

	// Test with string values
	expected2 := NewMapOrd[string, Int]()
	expected2.Set("banana", 1)
	expected2.Set("grape", 1)
	expected2.Set("orange", 2)
	expected2.Set("apple", 3)
	expected2.SortByKey(cmp.Cmp)

	result2 := sl2.Iter().Counter().Collect()
	result2.SortByKey(cmp.Cmp)
	if !result2.Eq(expected2) {
		t.Errorf("Counter() returned %v, expected %v", result2, expected2)
	}
}

func TestSliceIntersperse(t *testing.T) {
	// Test case 1: Intersperse strings with a comma
	testSlice := Slice[string]{"apple", "banana", "orange"}
	expected := Slice[string]{"apple", ", ", "banana", ", ", "orange"}
	interspersed := testSlice.Iter().Intersperse(", ").Collect()

	if interspersed.Ne(expected) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expected, interspersed)
	}

	// Test case 2: Intersperse strings with a dash
	testSlice = Slice[string]{"apple", "banana", "orange"}
	expected = Slice[string]{"apple", "-", "banana", "-", "orange"}
	interspersed = testSlice.Iter().Intersperse("-").Collect()

	if interspersed.Ne(expected) {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expected, interspersed)
	}

	// Test case 3: Intersperse empty slice
	emptySlice := Slice[string]{}    // Create an empty slice of strings
	expectedEmpty := Slice[string]{} // Expected empty slice
	interspersedEmpty := emptySlice.Iter().Intersperse(", ").Collect()

	if interspersedEmpty.Ne(expectedEmpty) {
		t.Errorf("Test case 3 failed. Expected: %v, Got: %v", expectedEmpty, interspersedEmpty)
	}
}

func TestSliceIterReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice[int]
		expected []int
	}{
		{
			name:     "empty slice",
			input:    Slice[int]{},
			expected: nil,
		},
		{
			name:     "single element",
			input:    Slice[int]{1},
			expected: []int{1},
		},
		{
			name:     "multiple elements",
			input:    Slice[int]{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iterator := tt.input.IterReverse()
			var result []int
			iterator.ForEach(func(element int) {
				result = append(result, element)
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IterReverse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name        string
		start, stop int
		step        []int
		want        Slice[int]
	}{
		{
			name:  "default step",
			start: 0, stop: 3,
			step: nil,
			want: Slice[int]{0, 1, 2},
		},
		{
			name:  "positive step 1",
			start: 0, stop: 5,
			step: []int{1},
			want: Slice[int]{0, 1, 2, 3, 4},
		},
		{
			name:  "custom positive step",
			start: 2, stop: 10,
			step: []int{2},
			want: Slice[int]{2, 4, 6, 8},
		},
		{
			name:  "negative step",
			start: 5, stop: 0,
			step: []int{-1},
			want: Slice[int]{5, 4, 3, 2, 1},
		},
		{
			name:  "zero step yields nothing",
			start: 0, stop: 5,
			step: []int{0},
			want: Slice[int]{},
		},
		{
			name:  "empty range when start == stop",
			start: 3, stop: 3,
			step: nil,
			want: Slice[int]{},
		},
		{
			name:  "step never approaches stop",
			start: 0, stop: 5,
			step: []int{-1},
			want: Slice[int]{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Range(tc.start, tc.stop, tc.step...).Collect()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Range(%d, %d, %v) = %v; want %v",
					tc.start, tc.stop, tc.step, got, tc.want)
			}
		})
	}
}

func TestRangeInclusive(t *testing.T) {
	tests := []struct {
		name        string
		start, stop int
		step        []int
		want        Slice[int]
	}{
		{
			name:  "default step",
			start: 0, stop: 3,
			step: nil,
			want: Slice[int]{0, 1, 2, 3},
		},
		{
			name:  "positive step 1",
			start: 0, stop: 5,
			step: []int{1},
			want: Slice[int]{0, 1, 2, 3, 4, 5},
		},
		{
			name:  "custom positive step",
			start: 2, stop: 10,
			step: []int{2},
			want: Slice[int]{2, 4, 6, 8, 10},
		},
		{
			name:  "negative step",
			start: 5, stop: 0,
			step: []int{-1},
			want: Slice[int]{5, 4, 3, 2, 1, 0},
		},
		{
			name:  "zero step yields nothing",
			start: 0, stop: 5,
			step: []int{0},
			want: Slice[int]{},
		},
		{
			name:  "inclusive when start == stop",
			start: 3, stop: 3,
			step: nil,
			want: Slice[int]{3},
		},
		{
			name:  "step never approaches stop",
			start: 0, stop: 5,
			step: []int{-1},
			want: Slice[int]{},
		},
		{
			name:  "step overshoots stop",
			start: 0, stop: 5,
			step: []int{6},
			want: Slice[int]{0},
		},
		{
			name:  "step exactly reaches stop",
			start: 0, stop: 6,
			step: []int{3},
			want: Slice[int]{0, 3, 6},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RangeInclusive(tc.start, tc.stop, tc.step...).Collect()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("RangeInclusive(%d, %d, %v) = %v; want %v",
					tc.start, tc.stop, tc.step, got, tc.want)
			}
		})
	}
}
