package g_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

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
	iter := g.FromChan(ch)

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
	slice1 := g.Slice[int]{1, 2, 3, 4, 5}
	isEven := func(val int) bool {
		return val%2 == 0
	}

	evens, odds := slice1.Iter().Partition(isEven)
	expectedEvens := g.Slice[int]{2, 4}
	expectedOdds := g.Slice[int]{1, 3, 5}

	if !reflect.DeepEqual(evens, expectedEvens) {
		t.Errorf("Expected evens %v, but got %v", expectedEvens, evens)
	}

	if !reflect.DeepEqual(odds, expectedOdds) {
		t.Errorf("Expected odds %v, but got %v", expectedOdds, odds)
	}

	// Test case 2: Partitioning with strings
	slice2 := g.Slice[string]{"apple", "banana", "cherry", "date"}
	hasA := func(val string) bool {
		return strings.Contains(val, "a")
	}

	withA, withoutA := slice2.Iter().Partition(hasA)
	expectedWithA := g.Slice[string]{"apple", "banana", "date"}
	expectedWithoutA := g.Slice[string]{"cherry"}

	if !reflect.DeepEqual(withA, expectedWithA) {
		t.Errorf("Expected withA %v, but got %v", expectedWithA, withA)
	}

	if !reflect.DeepEqual(withoutA, expectedWithoutA) {
		t.Errorf("Expected withoutA %v, but got %v", expectedWithoutA, withoutA)
	}

	// Test case 3: Partitioning an empty slice
	emptySlice := g.Slice[int]{}
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
	slice1 := g.Slice[int]{0, 1, 2, 3}
	combs1 := slice1.Iter().Combinations(3).Collect()
	expectedCombs1 := []g.Slice[int]{
		{0, 1, 2},
		{0, 1, 3},
		{0, 2, 3},
		{1, 2, 3},
	}

	if !reflect.DeepEqual(combs1, expectedCombs1) {
		t.Errorf("Test case 1 failed: expected %v, but got %v", expectedCombs1, combs1)
	}

	// Test case 2: Combinations of strings
	p1 := g.SliceOf[g.String]("a", "b")
	p2 := g.SliceOf[g.String]("c", "d")
	combs2 := p1.Iter().Chain(p2.Iter()).Map(g.String.Upper).Combinations(2).Collect()
	expectedCombs2 := []g.Slice[g.String]{
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
	p3 := g.SliceOf[any]("x", "y")
	p4 := g.SliceOf[any](1, 2)
	combs3 := p3.Iter().Chain(p4.Iter()).Combinations(2).Collect()
	expectedCombs3 := []g.Slice[any]{
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
	emptySlice := g.Slice[int]{}
	combs4 := emptySlice.Iter().Combinations(2).Collect()
	expectedCombs4 := []g.Slice[int]{}

	if !reflect.DeepEqual(combs4, expectedCombs4) {
		t.Errorf("Test case 4 failed: expected %v, but got %v", expectedCombs4, combs4)
	}

	// Test case 5: Combinations with k greater than slice length
	slice5 := g.Slice[int]{1, 2, 3}
	combs5 := slice5.Iter().Combinations(4).Collect()
	expectedCombs5 := []g.Slice[int]{}

	if !reflect.DeepEqual(combs5, expectedCombs5) {
		t.Errorf("Test case 5 failed: expected %v, but got %v", expectedCombs5, combs5)
	}
}

func TestSliceIterSortBy(t *testing.T) {
	sl1 := g.NewSlice[int]().Append(3, 1, 4, 1, 5)
	expected1 := g.NewSlice[int]().Append(1, 1, 3, 4, 5)

	actual1 := sl1.Iter().SortBy(cmp.Cmp).Collect()

	if !actual1.Eq(expected1) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected1, actual1)
	}

	sl2 := g.NewSlice[string]().Append("foo", "bar", "baz")
	expected2 := g.NewSlice[string]().Append("foo", "baz", "bar")

	actual2 := sl2.Iter().SortBy(func(a, b string) cmp.Ordering { return cmp.Cmp(b, a) }).Collect()

	if !actual2.Eq(expected2) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected2, actual2)
	}

	sl3 := g.NewSlice[int]()
	expected3 := g.NewSlice[int]()

	actual3 := sl3.Iter().SortBy(cmp.Cmp).Collect()

	if !actual3.Eq(expected3) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected3, actual3)
	}
}

func TestSliceIterDedup(t *testing.T) {
	// Test case 1: Dedup with consecutive duplicate elements for int
	sliceInt := g.Slice[int]{1, 2, 2, 3, 4, 4, 4, 5}
	expectedResultInt := g.Slice[int]{1, 2, 3, 4, 5}

	iterInt := sliceInt.Iter().Dedup()
	resultInt := iterInt.Collect()

	if !reflect.DeepEqual(resultInt, expectedResultInt) {
		t.Errorf("Dedup failed for int. Expected %v, got %v", expectedResultInt, resultInt)
	}

	// Test case 2: Dedup with consecutive duplicate elements for string
	sliceString := g.Slice[string]{"apple", "orange", "orange", "banana", "banana", "grape"}
	expectedResultString := g.Slice[string]{"apple", "orange", "banana", "grape"}

	iterString := sliceString.Iter().Dedup()
	resultString := iterString.Collect()

	if !reflect.DeepEqual(resultString, expectedResultString) {
		t.Errorf("Dedup failed for string. Expected %v, got %v", expectedResultString, resultString)
	}

	// Test case 3: Dedup with consecutive duplicate elements for float64
	sliceFloat64 := g.Slice[float64]{1.2, 2.3, 2.3, 3.4, 4.5, 4.5, 4.5, 5.6}
	expectedResultFloat64 := g.Slice[float64]{1.2, 2.3, 3.4, 4.5, 5.6}

	iterFloat64 := sliceFloat64.Iter().Dedup()
	resultFloat64 := iterFloat64.Collect()

	if !reflect.DeepEqual(resultFloat64, expectedResultFloat64) {
		t.Errorf("Dedup failed for float64. Expected %v, got %v", expectedResultFloat64, resultFloat64)
	}

	// Test case 4: Dedup with consecutive duplicate elements for custom non-comparable struct
	type myStruct struct {
		val []int
	}

	sliceStruct := g.Slice[myStruct]{
		{val: []int{1}},
		{val: []int{2}},
		{val: []int{2}},
		{val: []int{3}},
		{val: []int{3}},
		{val: []int{4}},
	}

	expectedResultStruct := g.Slice[myStruct]{{val: []int{1}}, {val: []int{2}}, {val: []int{3}}, {val: []int{4}}}

	iterStruct := sliceStruct.Iter().Dedup()
	resultStruct := iterStruct.Collect()

	if !reflect.DeepEqual(resultStruct, expectedResultStruct) {
		t.Errorf("Dedup failed for custom struct. Expected %v, got %v", expectedResultStruct, resultStruct)
	}
}

func TestSliceIterStepBy(t *testing.T) {
	// Test case 1: StepBy with a step size of 3
	slice := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedResult := g.Slice[int]{1, 4, 7, 10}

	iter := slice.Iter().StepBy(3)
	result := iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 2: StepBy with a step size of 2
	slice = g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedResult = g.Slice[int]{1, 3, 5, 7, 9}

	iter = slice.Iter().StepBy(2)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 3: StepBy with a step size larger than the slice length
	slice = g.Slice[int]{1, 2, 3, 4, 5}
	expectedResult = g.Slice[int]{1}

	iter = slice.Iter().StepBy(10)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 4: StepBy with a step size of 1
	slice = g.Slice[int]{1, 2, 3, 4, 5}
	expectedResult = g.Slice[int]{1, 2, 3, 4, 5}

	iter = slice.Iter().StepBy(1)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}
}

func TestSliceIterPermutations(t *testing.T) {
	// Test case 1: Single element slice
	slice1 := g.SliceOf(1)
	perms1 := slice1.Iter().Permutations().Collect()
	expectedPerms1 := []g.Slice[int]{slice1}

	if !reflect.DeepEqual(perms1, expectedPerms1) {
		t.Errorf("expected %v, but got %v", expectedPerms1, perms1)
	}

	// Test case 2: Two-element string slice
	slice2 := g.SliceOf("a", "b")
	perms2 := slice2.Iter().Permutations().Collect()
	expectedPerms2 := []g.Slice[string]{
		{"a", "b"},
		{"b", "a"},
	}

	if !reflect.DeepEqual(perms2, expectedPerms2) {
		t.Errorf("expected %v, but got %v", expectedPerms2, perms2)
	}

	// Test case 3: Three-element float64 slice
	slice3 := g.SliceOf(1.0, 2.0, 3.0)
	perms3 := slice3.Iter().Permutations().Collect()
	expectedPerms3 := []g.Slice[float64]{
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
	slice4 := g.Slice[any]{}
	perms4 := slice4.Iter().Permutations().Collect()
	expectedPerms4 := []g.Slice[any]{slice4}

	if !reflect.DeepEqual(perms4, expectedPerms4) {
		t.Errorf("expected %v, but got %v", expectedPerms4, perms4)
	}

	// Additional Test case 5: Four-element mixed-type slice
	slice5 := g.SliceOf[any]("a", 1, 2.5, true)
	perms5 := slice5.Iter().Permutations().Collect()
	expectedPerms5 := []g.Slice[any]{
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
		input    g.Slice[int]
		expected []g.Slice[int]
		size     g.Int
	}{
		{
			name:     "empty slice",
			input:    g.NewSlice[int](),
			expected: []g.Slice[int]{g.NewSlice[int]()},
			size:     2,
		},
		{
			name:     "single chunk",
			input:    g.NewSlice[int]().Append(1, 2, 3),
			expected: []g.Slice[int]{g.NewSlice[int]().Append(1, 2, 3)},
			size:     3,
		},
		{
			name:  "multiple chunks",
			input: g.NewSlice[int]().Append(1, 2, 3, 4, 5, 6),
			expected: []g.Slice[int]{
				g.NewSlice[int]().Append(1, 2),
				g.NewSlice[int]().Append(3, 4),
				g.NewSlice[int]().Append(5, 6),
			},
			size: 2,
		},
		{
			name:  "last chunk is smaller",
			input: g.NewSlice[int]().Append(1, 2, 3, 4, 5),
			expected: []g.Slice[int]{
				g.NewSlice[int]().Append(1, 2),
				g.NewSlice[int]().Append(3, 4),
				g.NewSlice[int]().Append(5),
			},
			size: 2,
		},
		{
			name:     "chunk size bigger than slice length",
			input:    g.NewSlice[int]().Append(1, 2, 3, 4),
			expected: []g.Slice[int]{g.NewSlice[int]().Append(1, 2, 3, 4)},
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
	sl1 := g.NewSlice[int]()
	sl2 := g.NewSlice[int]().Append(1, 2, 3)
	sl3 := g.NewSlice[int]().Append(2, 4, 6)

	testCases := []struct {
		f    func(int) bool
		name string
		sl   g.Slice[int]
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
	sl1 := g.NewSlice[int]()
	f1 := func(x int) bool { return x > 0 }

	if sl1.Iter().Any(f1) {
		t.Errorf("Expected false for empty slice, got true")
	}

	sl2 := g.NewSlice[int]().Append(1, 2, 3)
	f2 := func(x int) bool { return x < 1 }

	if sl2.Iter().Any(f2) {
		t.Errorf("Expected false for slice with no matching elements, got true")
	}

	sl3 := g.NewSlice[string]().Append("foo", "bar")
	f3 := func(x string) bool { return x == "bar" }

	if !sl3.Iter().Any(f3) {
		t.Errorf("Expected true for slice with one matching element, got false")
	}

	sl4 := g.NewSlice[int]().Append(1, 2, 3, 4, 5)
	f4 := func(x int) bool { return x%2 == 0 }

	if !sl4.Iter().Any(f4) {
		t.Errorf("Expected true for slice with multiple matching elements, got false")
	}
}

func TestSliceIterFold(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	sum := sl.Iter().Fold(0, func(index, value int) int { return index + value })

	if sum != 15 {
		t.Errorf("Expected %d, got %d", 15, sum)
	}
}

func TestSliceIterFilter(t *testing.T) {
	var sl g.Slice[int]

	sl = sl.Append(1, 2, 3, 4, 5)
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
	sl := g.Slice[int]{1, 2, 3, 4, 5}
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
	sl := g.Slice[int]{1, 2, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10}
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
	sl1 := g.NewSlice[int]().Append(1, 2, 3, 4, 5)
	sl2 := g.NewSlice[string]().Append("foo", "bar", "baz")
	sl3 := g.NewSlice[float64]().Append(1.1, 2.2, 3.3, 4.4)

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
	s1 := g.SliceOf(1, 2, 3, 4)
	s2 := g.SliceOf(5, 6, 7, 8)
	expected := g.MapOrd[int, int]{{1, 5}, {2, 6}, {3, 7}, {4, 8}}
	result := s1.Iter().Zip(s2.Iter()).Collect()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s1, s2, result, expected)
	}

	s3 := g.SliceOf(1, 2, 3)
	s4 := g.SliceOf(4, 5)
	expected = g.MapOrd[int, int]{{1, 4}, {2, 5}}
	result = s3.Iter().Zip(s4.Iter()).Collect()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s3, s4, result, expected)
	}
}

func TestSliceIterFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    g.Slice[any]
		expected g.Slice[any]
	}{
		{
			name:     "Empty slice",
			input:    g.Slice[any]{},
			expected: g.Slice[any]{},
		},
		{
			name:     "Flat slice",
			input:    g.Slice[any]{1, "abc", 3.14},
			expected: g.Slice[any]{1, "abc", 3.14},
		},
		{
			name: "Nested slice",
			input: g.Slice[any]{
				1,
				g.SliceOf(2, 3),
				"abc",
				g.SliceOf("def", "ghi"),
				g.SliceOf(4.5, 6.7),
			},
			expected: g.Slice[any]{1, 2, 3, "abc", "def", "ghi", 4.5, 6.7},
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
		slice := g.Slice[int]{1, 2, 3, 4, 5}
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
		slice := g.Slice[int]{1, 2, 3, 4, 5}
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
		emptySlice := g.Slice[int]{}
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
	seq := g.Slice[int]{1, 2, 3, 4, 5}
	count := seq.Iter().Count()
	if count != 5 {
		t.Errorf("Expected count to be %d, got %d", 5, count)
	}

	// Test case 2: Empty sequence
	emptySeq := g.Slice[int]{}
	emptyCount := emptySeq.Iter().Count()
	if emptyCount != 0 {
		t.Errorf("Expected count of an empty sequence to be %d, got %d", 0, emptyCount)
	}
}

func TestSliceIterCycle(t *testing.T) {
	// Test case 1: Cyclic behavior
	seq := g.Slice[int]{1, 2, 3}
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
	seq := g.Slice[string]{"bbb", "ddd", "xxx", "aaa", "ccc"}
	enumerated := seq.Iter().Enumerate().Collect()

	expected := g.NewMapOrd[g.Int, string]()
	expected.Set(0, "bbb").Set(1, "ddd").Set(2, "xxx").Set(3, "aaa").Set(4, "ccc")

	for i, v := range enumerated {
		if expected[i] != v {
			t.Errorf("Expected element at index %d to be %v, got %v", i, expected[i], v)
		}
	}
}

func TestSliceIterSkip(t *testing.T) {
	// Test case 1: Skip elements
	seq := g.Slice[int]{1, 2, 3, 4, 5, 6}
	skipped := seq.Iter().Skip(3).Collect()
	expected := g.Slice[int]{4, 5, 6}
	if len(skipped) != len(expected) {
		t.Errorf("Expected skipped slice to have length %d, got %d", len(expected), len(skipped))
	}
	for i := range expected {
		if skipped[i] != expected[i] {
			t.Errorf("Expected element at index %d to be %d, got %d", i, expected[i], skipped[i])
		}
	}

	// Test case 2: Skip all elements
	seq2 := g.Slice[string]{"a", "b", "c"}
	skipped2 := seq2.Iter().Skip(3).Collect()
	if len(skipped2) != 0 {
		t.Errorf("Expected skipped slice of all elements to be empty, got length %d", len(skipped2))
	}
}

func TestSliceIterUnique(t *testing.T) {
	// Test case 1: Unique elements
	seq := g.Slice[int]{1, 2, 3, 2, 4, 5, 3}
	unique := seq.Iter().Unique().Collect()

	expected := g.Slice[int]{1, 2, 3, 4, 5}
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
	seq := g.Slice[int]{1, 2, 3, 4, 5}
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
	seq := g.Slice[int]{1, 2, 3, 4, 5, 6}
	windows := seq.Iter().Windows(3).Collect()

	expected := []g.Slice[int]{
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
	seq := g.Slice[int]{1, 2, 3}

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
	s := g.Slice[int]{1, 2, 3}

	// Define a slice to store the inspected elements
	var inspectedElements g.Slice[int]

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
	sl1 := g.Slice[int]{1, 2, 3, 2, 1, 4, 5, 4, 4}
	sl2 := g.Slice[string]{"apple", "banana", "orange", "apple", "apple", "orange", "grape"}

	expected1 := g.NewMapOrd[int, g.Int]()
	expected1.
		Set(3, 1).
		Set(5, 1).
		Set(1, 2).
		Set(2, 2).
		Set(4, 3).
		SortByKey(cmp.Cmp)

	result1 := sl1.Iter().Counter().Collect().SortByKey(cmp.Cmp)
	if !result1.Eq(expected1) {
		t.Errorf("Counter() returned %v, expected %v", result1, expected1)
	}

	// Test with string values
	expected2 := g.NewMapOrd[string, g.Int]()
	expected2.
		Set("banana", 1).
		Set("grape", 1).
		Set("orange", 2).
		Set("apple", 3).
		SortByKey(cmp.Cmp)

	result2 := sl2.Iter().Counter().Collect().SortByKey(cmp.Cmp)
	if !result2.Eq(expected2) {
		t.Errorf("Counter() returned %v, expected %v", result2, expected2)
	}
}

func TestSliceIntersperse(t *testing.T) {
	// Test case 1: Intersperse strings with a comma
	testSlice := g.Slice[string]{"apple", "banana", "orange"}
	expected := g.Slice[string]{"apple", ", ", "banana", ", ", "orange"}
	interspersed := testSlice.Iter().Intersperse(", ").Collect()

	if interspersed.Ne(expected) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expected, interspersed)
	}

	// Test case 2: Intersperse strings with a dash
	testSlice = g.Slice[string]{"apple", "banana", "orange"}
	expected = g.Slice[string]{"apple", "-", "banana", "-", "orange"}
	interspersed = testSlice.Iter().Intersperse("-").Collect()

	if interspersed.Ne(expected) {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expected, interspersed)
	}

	// Test case 3: Intersperse empty slice
	emptySlice := g.Slice[string]{}    // Create an empty slice of strings
	expectedEmpty := g.Slice[string]{} // Expected empty slice
	interspersedEmpty := emptySlice.Iter().Intersperse(", ").Collect()

	if interspersedEmpty.Ne(expectedEmpty) {
		t.Errorf("Test case 3 failed. Expected: %v, Got: %v", expectedEmpty, interspersedEmpty)
	}
}
