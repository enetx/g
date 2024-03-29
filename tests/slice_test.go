package g_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/enetx/g"
)

func TestSliceUnpack(t *testing.T) {
	tests := []struct {
		name     string
		slice    g.Slice[int]
		vars     []*int
		expected []int
	}{
		{
			name:     "Unpack with valid indices",
			slice:    g.Slice[int]{1, 2, 3, 4, 5},
			vars:     []*int{new(int), new(int), new(int)},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Unpack with invalid indices",
			slice:    g.Slice[int]{1, 2, 3},
			vars:     []*int{new(int), new(int), new(int), new(int)},
			expected: []int{1, 2, 3, 0}, // Expecting zero value for the fourth variable
		},
		{
			name:     "Unpack with empty slice",
			slice:    g.Slice[int]{},
			vars:     []*int{new(int)},
			expected: []int{0}, // Expecting zero value for the only variable
		},
		{
			name:     "Unpack with nil slice",
			slice:    nil,
			vars:     []*int{new(int)},
			expected: []int{0}, // Expecting zero value for the only variable
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.slice.Unpack(test.vars...)
			for i, v := range test.vars {
				if *v != test.expected[i] {
					t.Errorf("Expected %d but got %d", test.expected[i], *v)
				}
			}
		})
	}
}

func TestSubSliceWithStep(t *testing.T) {
	slice := g.SliceOf(1, 2, 3, 4, 5, 6, 7, 8, 9)

	testCases := []struct {
		start    int
		end      int
		step     int
		expected g.Slice[int]
	}{
		{1, 7, 2, g.SliceOf(2, 4, 6)},
		{0, 9, 3, g.SliceOf(1, 4, 7)},
		{2, 6, 1, g.SliceOf(3, 4, 5, 6)},
		{0, 9, 2, g.SliceOf(1, 3, 5, 7, 9)},
		{0, 9, 4, g.SliceOf(1, 5, 9)},
		{6, 1, -2, g.SliceOf(7, 5, 3)},
		{8, 1, -3, g.SliceOf(9, 6, 3)},
		{8, 0, -2, g.SliceOf(9, 7, 5, 3)},
		{8, 0, -1, g.SliceOf(9, 8, 7, 6, 5, 4, 3, 2)},
		{8, 0, -4, g.SliceOf(9, 5)},
		{-1, -6, -2, g.SliceOf(9, 7, 5)},
		{-2, -9, -3, g.SliceOf(8, 5, 2)},
		{-1, -8, -2, g.SliceOf(9, 7, 5, 3)},
		{-3, -10, -2, g.SliceOf(7, 5, 3, 1)},
		{-1, -10, -1, g.SliceOf(9, 8, 7, 6, 5, 4, 3, 2, 1)},
		{-5, -1, -1, g.Slice[int]{}},
		{-1, -5, -1, g.SliceOf(9, 8, 7, 6)},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("start:%d_end:%d_step:%d", tc.start, tc.end, tc.step), func(t *testing.T) {
			result := slice.SubSlice(tc.start, tc.end, tc.step)

			if !result.Eq(tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestSortInts(t *testing.T) {
	slice := g.Slice[int]{5, 2, 8, 1, 6}
	slice.Sort()

	expected := g.Slice[int]{1, 2, 5, 6, 8}

	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v but got %v", expected, slice)
	}
}

func TestSortStrings(t *testing.T) {
	slice := g.Slice[string]{"apple", "orange", "banana", "grape"}
	slice.Sort()

	expected := g.Slice[string]{"apple", "banana", "grape", "orange"}

	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v but got %v", expected, slice)
	}
}

func TestSliceSortFloats(t *testing.T) {
	slice := g.Slice[float64]{5.6, 2.3, 8.9, 1.2, 6.7}
	slice.Sort()

	expected := g.Slice[float64]{1.2, 2.3, 5.6, 6.7, 8.9}

	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v but got %v", expected, slice)
	}
}

func TestSliceInsert(t *testing.T) {
	// Test insertion in the middle
	slice := g.Slice[string]{"a", "b", "c", "d"}
	newSlice := slice.Insert(2, "e", "f")
	expected := g.Slice[string]{"a", "b", "e", "f", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion at the start
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Insert(0, "x", "y")
	expected = g.Slice[string]{"x", "y", "a", "b", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(0) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion at the end
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Insert(slice.Len(), "x", "y")
	expected = g.Slice[string]{"a", "b", "c", "d", "x", "y"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(end) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion with negative index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Insert(-2, "x", "y")
	expected = g.Slice[string]{"a", "b", "x", "y", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(-2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion at index 0 with an empty slice
	slice = g.Slice[string]{}
	newSlice = slice.Insert(0, "x", "y")
	expected = g.Slice[string]{"x", "y"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(0) with empty slice failed. Expected %v, but got %v", expected, newSlice)
	}
}

func TestSliceInsertInPlace(t *testing.T) {
	// Test insertion in the middle
	slice := g.Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(2, "e", "f")
	expected := g.Slice[string]{"a", "b", "e", "f", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion at the start
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(0, "x", "y")
	expected = g.Slice[string]{"x", "y", "a", "b", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(0) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion at the end
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(slice.Len(), "x", "y")
	expected = g.Slice[string]{"a", "b", "c", "d", "x", "y"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(end) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion with negative index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(-2, "x", "y")
	expected = g.Slice[string]{"a", "b", "x", "y", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(-2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion at index 0 with an empty slice
	slice = g.Slice[string]{}
	slice.InsertInPlace(0, "x", "y")
	expected = g.Slice[string]{"x", "y"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(0) with empty slice failed. Expected %v, but got %v", expected, slice)
	}
}

func TestSliceToSlice(t *testing.T) {
	sl := g.NewSlice[int]().Append(1, 2, 3, 4, 5)
	slice := sl.Std()

	if len(slice) != sl.Len() {
		t.Errorf("Expected length %d, but got %d", sl.Len(), len(slice))
	}

	for i, v := range sl {
		if v != slice[i] {
			t.Errorf("Expected value %d at index %d, but got %d", v, i, slice[i])
		}
	}
}

func TestSliceShuffle(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sl.Shuffle()

	if sl.Len() != 10 {
		t.Error("Expected length of 10, got ", sl.Len())
	}
}

func TestSliceReverse(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	sl.Reverse()

	if !reflect.DeepEqual(sl, g.Slice[int]{5, 4, 3, 2, 1}) {
		t.Errorf("Expected %v, got %v", g.Slice[int]{5, 4, 3, 2, 1}, sl)
	}
}

func TestSliceIndex(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}

	if sl.Index(1) != 0 {
		t.Error("Index of 1 should be 0")
	}

	if sl.Index(2) != 1 {
		t.Error("Index of 2 should be 1")
	}

	if sl.Index(3) != 2 {
		t.Error("Index of 3 should be 2")
	}

	if sl.Index(4) != 3 {
		t.Error("Index of 4 should be 3")
	}

	if sl.Index(5) != 4 {
		t.Error("Index of 5 should be 4")
	}

	if sl.Index(6) != -1 {
		t.Error("Index of 6 should be -1")
	}
}

func TestSliceRandomSample(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := sl.RandomSample(5)

	if result.Len() != 5 {
		t.Errorf("Expected result length to be 5, got %d", result.Len())
	}

	for _, item := range result {
		if !sl.Contains(item) {
			t.Errorf("Expected result to contain only items from the original slice, got %d", item)
		}
	}
}

func TestSliceAddUnique(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3}
	sl = sl.AddUnique(4, 5, 6)

	if !sl.Contains(4) {
		t.Error("AddUnique failed")
	}

	sl = sl.AddUnique(4, 5, 6)
	if sl.Len() != 6 {
		t.Error("AddUnique failed")
	}
}

func TestSliceCount(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5, 6, 7}

	if sl.Count(1) != 1 {
		t.Error("Expected 1, got ", sl.Count(1))
	}

	if sl.Count(2) != 1 {
		t.Error("Expected 1, got ", sl.Count(2))
	}

	if sl.Count(3) != 1 {
		t.Error("Expected 1, got ", sl.Count(3))
	}

	if sl.Count(4) != 1 {
		t.Error("Expected 1, got ", sl.Count(4))
	}

	if sl.Count(5) != 1 {
		t.Error("Expected 1, got ", sl.Count(5))
	}

	if sl.Count(6) != 1 {
		t.Error("Expected 1, got ", sl.Count(6))
	}

	if sl.Count(7) != 1 {
		t.Error("Expected 1, got ", sl.Count(7))
	}

	emptySlice := g.NewSlice[int]()
	if emptySlice.Count(7) != 0 {
		t.Error("Expected 0, got", emptySlice.Count(7))
	}
}

func TestSliceSortBy(t *testing.T) {
	sl1 := g.NewSlice[int]().Append(3, 1, 4, 1, 5)
	expected1 := g.NewSlice[int]().Append(1, 1, 3, 4, 5)

	sl1.SortBy(func(a, b int) bool { return a < b })

	if !sl1.Eq(expected1) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected1, sl1)
	}

	sl2 := g.NewSlice[string]().Append("foo", "bar", "baz")
	expected2 := g.NewSlice[string]().Append("foo", "baz", "bar")

	sl2.SortBy(func(a, b string) bool { return a > b })

	if !sl2.Eq(expected2) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected2, sl2)
	}

	sl3 := g.NewSlice[int]()
	expected3 := g.NewSlice[int]()

	sl3.SortBy(func(a, b int) bool { return a < b })

	if !sl3.Eq(expected3) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected3, sl3)
	}
}

func TestSliceJoin(t *testing.T) {
	sl := g.Slice[string]{"1", "2", "3", "4", "5"}
	str := sl.Join(",")

	if !strings.EqualFold("1,2,3,4,5", str.Std()) {
		t.Errorf("Expected 1,2,3,4,5, got %s", str.Std())
	}
}

func TestSliceToStringSlice(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3}
	result := sl.ToStringSlice()
	expected := []string{"1", "2", "3"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSliceAdd(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3}.Append(4, 5, 6)

	if !reflect.DeepEqual(sl, g.Slice[int]{1, 2, 3, 4, 5, 6}) {
		t.Error("Add failed")
	}
}

func TestSliceClone(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3}
	slClone := sl.Clone()

	if !sl.Eq(slClone) {
		t.Errorf("Clone() failed, expected %v, got %v", sl, slClone)
	}
}

func TestSliceCut(t *testing.T) {
	slice := g.Slice[int]{1, 2, 3, 4, 5}

	// Test normal range
	newSlice := slice.Cut(1, 3)
	expected := g.Slice[int]{1, 4, 5}
	if !newSlice.Eq(expected) {
		t.Errorf("Cut(1, 3) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test range with negative indices
	newSlice = slice.Cut(-3, -2)
	expected = g.Slice[int]{1, 2, 4, 5}
	if !newSlice.Eq(expected) {
		t.Errorf("Cut(-3, -2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test empty range
	newSlice = slice.Cut(0, 5)
	expected = g.Slice[int]{}
	if !newSlice.Eq(expected) {
		t.Errorf("Cut(3, 2) failed. Expected %v, but got %v", expected, newSlice)
	}
}

func TestSliceCutInPlace(t *testing.T) {
	slice := g.Slice[int]{1, 2, 3, 4, 5}

	// Test normal range
	slice.CutInPlace(1, 3)
	expected := g.Slice[int]{1, 4, 5}
	if !slice.Eq(expected) {
		t.Errorf("CutInPlace(1, 3) failed. Expected %v, but got %v", expected, slice)
	}

	// Test range with negative indices
	slice = g.Slice[int]{1, 2, 3, 4, 5} // Restore the original slice
	slice.CutInPlace(-3, -2)
	expected = g.Slice[int]{1, 2, 4, 5}
	if !slice.Eq(expected) {
		t.Errorf("CutInPlace(-3, -2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test empty range
	slice = g.Slice[int]{1, 2, 3, 4, 5} // Restore the original slice
	slice.CutInPlace(0, 5)
	expected = g.Slice[int]{}
	if !slice.Eq(expected) {
		t.Errorf("CutInPlace(0, 5) failed. Expected %v, but got %v", expected, slice)
	}
}

func TestSliceLast(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	if sl.Last() != 5 {
		t.Error("Last() failed")
	}
}

func TestSliceLen(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	if sl.Len() != 5 {
		t.Errorf("Expected 5, got %d", sl.Len())
	}
}

func TestSlicePop(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	last, sl := sl.Pop()

	if last != 5 {
		t.Errorf("Expected 5, got %v", last)
	}

	if sl.Len() != 4 {
		t.Errorf("Expected 4, got %v", sl.Len())
	}

	r := g.Slice[int]{1, 2, 3, 4}
	if sl.Ne(r) {
		t.Errorf("Expected %v, got %v", r, sl)
	}
}

func TestSliceMaxInt(t *testing.T) {
	sl := g.Slice[g.Int]{1, 2, 3, 4, 5}
	if max := sl.Max(); max != 5 {
		t.Errorf("Max() = %d, want: %d.", max, 5)
	}
}

func TestSliceMaxFloats(t *testing.T) {
	sl := g.Slice[g.Float]{2.2, 2.8, 2.1, 2.7}
	if max := sl.Max(); max != 2.8 {
		t.Errorf("Max() = %f, want: %f.", max, 2.8)
	}
}

func TestSliceMinFloat(t *testing.T) {
	sl := g.Slice[g.Float]{2.2, 2.8, 2.1, 2.7}
	if min := sl.Min(); min != 2.1 {
		t.Errorf("Min() = %f; want: %f", min, 2.1)
	}
}

func TestSliceMinInt(t *testing.T) {
	sl := g.Slice[g.Int]{1, 2, 3, 4, 5}
	if min := sl.Min(); min != 1 {
		t.Errorf("Min() = %d; want: %d", min, 1)
	}
}

func TestSliceDelete(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	sl = sl.Delete(2)

	if !reflect.DeepEqual(sl, g.Slice[int]{1, 2, 4, 5}) {
		t.Errorf("Delete(2) = %v, want %v", sl, g.Slice[int]{1, 2, 4, 5})
	}

	sl = g.Slice[int]{1, 2, 3, 4, 5}
	sl = sl.Delete(-2)

	if !reflect.DeepEqual(sl, g.Slice[int]{1, 2, 3, 5}) {
		t.Errorf("Delete(2) = %v, want %v", sl, g.Slice[int]{1, 2, 3, 5})
	}
}

func TestSliceSFill(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	sl.Fill(0)

	for _, v := range sl {
		if v != 0 {
			t.Errorf("Expected all elements to be 0, but found %d", v)
		}
	}
}

func TestSliceSet(t *testing.T) {
	sl := g.NewSlice[int](5)

	sl.Set(0, 1)
	sl.Set(0, 1)
	sl.Set(2, 2)
	sl.Set(4, 3)

	if !reflect.DeepEqual(sl, g.Slice[int]{1, 0, 2, 0, 3}) {
		t.Errorf("Set() = %v, want %v", sl, g.Slice[int]{1, 0, 2, 0, 3})
	}
}

func TestSliceCounter(t *testing.T) {
	sl1 := g.Slice[int]{1, 2, 3, 2, 1, 4, 5, 4, 4}
	sl2 := g.Slice[string]{"apple", "banana", "orange", "apple", "apple", "orange", "grape"}

	expected1 := g.NewMapOrd[int, uint]()
	expected1.
		Set(3, 1).
		Set(5, 1).
		Set(1, 2).
		Set(2, 2).
		Set(4, 3)

	result1 := sl1.Counter()
	if !result1.Eq(expected1) {
		t.Errorf("Counter() returned %v, expected %v", result1, expected1)
	}

	// Test with string values
	expected2 := g.NewMapOrd[string, uint]()
	expected2.
		Set("banana", 1).
		Set("grape", 1).
		Set("orange", 2).
		Set("apple", 3)

	result2 := sl2.Counter()
	if !result2.Eq(expected2) {
		t.Errorf("Counter() returned %v, expected %v", result2, expected2)
	}
}

func TestSliceReplace(t *testing.T) {
	tests := []struct {
		name     string
		input    g.Slice[string]
		i, j     int
		values   []string
		expected g.Slice[string]
	}{
		{
			name:     "basic test",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e", "f"},
			expected: g.Slice[string]{"a", "e", "f", "d"},
		},
		{
			name:     "replace at start",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        2,
			values:   []string{"e", "f"},
			expected: g.Slice[string]{"e", "f", "c", "d"},
		},
		{
			name:     "replace at end",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        2,
			j:        4,
			values:   []string{"e", "f"},
			expected: g.Slice[string]{"a", "b", "e", "f"},
		},
		{
			name:     "replace with more values",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        2,
			values:   []string{"e", "f", "g", "h"},
			expected: g.Slice[string]{"a", "e", "f", "g", "h", "c", "d"},
		},
		{
			name:     "replace with fewer values",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e"},
			expected: g.Slice[string]{"a", "e", "d"},
		},
		{
			name:     "replace entire slice",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        4,
			values:   []string{"e", "f", "g", "h"},
			expected: g.Slice[string]{"e", "f", "g", "h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Replace(tt.i, tt.j, tt.values...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSliceReplaceM(t *testing.T) {
	// Test replacement in the middle
	slice := g.Slice[string]{"a", "b", "c", "d"}
	newSlice := slice.Replace(1, 1, "zz", "xx")
	expected := g.Slice[string]{"a", "zz", "xx", "b", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(1, 1) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement with same start and end indices
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(2, 2, "zz", "xx")
	expected = g.Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(2, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from i to the end
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(2, slice.Len(), "zz", "xx")
	expected = g.Slice[string]{"a", "b", "zz", "xx"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(2, end) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from the start to j
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(0, 2, "zz", "xx")
	expected = g.Slice[string]{"zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(start, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test empty replacement
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(2, 2) // No replacement, should remain unchanged
	expected = g.Slice[string]{"a", "b", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(2, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from negative index to positive index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(-2, 2, "zz", "xx")
	expected = g.Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(-2, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from positive index to negative index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(1, -1, "zz", "xx")
	expected = g.Slice[string]{"a", "zz", "xx", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(1, -1) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from negative index to negative index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(-3, -2, "zz", "xx")
	expected = g.Slice[string]{"a", "zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(-3, -2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from negative index to positive index including negative values
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(-3, 3, "zz", "xx")
	expected = g.Slice[string]{"a", "zz", "xx", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(-3, 3) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement with empty slice
	slice = g.Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(1, 3)
	expected = g.Slice[string]{"a", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(1, 3) with empty slice failed. Expected %v, but got %v", expected, newSlice)
	}
}

func TestSliceReplaceInPlace(t *testing.T) {
	tests := []struct {
		name     string
		input    g.Slice[string]
		i, j     int
		values   []string
		expected g.Slice[string]
	}{
		{
			name:     "basic test",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e", "f"},
			expected: g.Slice[string]{"a", "e", "f", "d"},
		},
		{
			name:     "replace at start",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        2,
			values:   []string{"e", "f"},
			expected: g.Slice[string]{"e", "f", "c", "d"},
		},
		{
			name:     "replace at end",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        2,
			j:        4,
			values:   []string{"e", "f"},
			expected: g.Slice[string]{"a", "b", "e", "f"},
		},
		{
			name:     "replace with more values",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        2,
			values:   []string{"e", "f", "g", "h"},
			expected: g.Slice[string]{"a", "e", "f", "g", "h", "c", "d"},
		},
		{
			name:     "replace with fewer values",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e"},
			expected: g.Slice[string]{"a", "e", "d"},
		},
		{
			name:     "replace entire slice",
			input:    g.Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        4,
			values:   []string{"e", "f", "g", "h"},
			expected: g.Slice[string]{"e", "f", "g", "h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := &tt.input
			sl.ReplaceInPlace(tt.i, tt.j, tt.values...)
			if !reflect.DeepEqual(*sl, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, *sl)
			}
		})
	}
}

func TestSliceReplaceInPlaceM(t *testing.T) {
	// Test replacement in the middle
	slice := g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(1, 1, "zz", "xx")
	expected := g.Slice[string]{"a", "zz", "xx", "b", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(1, 1) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement with same start and end indices
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(2, 2, "zz", "xx")
	expected = g.Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(2, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from i to the end
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(2, slice.Len(), "zz", "xx")
	expected = g.Slice[string]{"a", "b", "zz", "xx"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(2, end) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from the start to j
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(0, 2, "zz", "xx")
	expected = g.Slice[string]{"zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(start, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test empty replacement
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(2, 2) // No replacement, should remain unchanged
	expected = g.Slice[string]{"a", "b", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(2, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from negative index to positive index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(-2, 2, "zz", "xx")
	expected = g.Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(-2, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from positive index to negative index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(1, -1, "zz", "xx")
	expected = g.Slice[string]{"a", "zz", "xx", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(1, -1) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from negative index to negative index
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(-3, -2, "zz", "xx")
	expected = g.Slice[string]{"a", "zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(-3, -2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from negative index to positive index including negative values
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(-3, 3, "zz", "xx")
	expected = g.Slice[string]{"a", "zz", "xx", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(-3, 3) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement with empty slice
	slice = g.Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(1, 3)
	expected = g.Slice[string]{"a", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(1, 3) with empty slice failed. Expected %v, but got %v", expected, slice)
	}
}

func TestSliceContainsAny(t *testing.T) {
	testCases := []struct {
		sl     g.Slice[int]
		other  g.Slice[int]
		expect bool
	}{
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{6, 7, 8, 9, 10}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{5, 6, 7, 8, 9}, true},
		{g.Slice[int]{}, g.Slice[int]{1, 2, 3, 4, 5}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1, 2, 3, 4, 5}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{6}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{6, 7, 8, 9, 0, 3}, true},
	}

	for _, tc := range testCases {
		if result := tc.sl.ContainsAny(tc.other...); result != tc.expect {
			t.Errorf("ContainsAny(%v, %v) = %v; want %v", tc.sl, tc.other, result, tc.expect)
		}
	}
}

func TestSliceContainsAll(t *testing.T) {
	testCases := []struct {
		sl     g.Slice[int]
		other  g.Slice[int]
		expect bool
	}{
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1, 2, 3}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1, 2, 3, 6}, false},
		{g.Slice[int]{}, g.Slice[int]{1, 2, 3, 4, 5}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1, 2, 3, 4, 5}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{6}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 5, 5}, true},
	}

	for _, tc := range testCases {
		if result := tc.sl.ContainsAll(tc.other...); result != tc.expect {
			t.Errorf("ContainsAll(%v, %v) = %v; want %v", tc.sl, tc.other, result, tc.expect)
		}
	}
}

func TestSliceSubSlice(t *testing.T) {
	// Test with an empty slice
	emptySlice := g.Slice[int]{}
	emptySubSlice := emptySlice.SubSlice(0, 0)
	if !emptySubSlice.Empty() {
		t.Errorf("Expected empty slice for empty source slice, but got: %v", emptySubSlice)
	}

	// Test with a non-empty slice
	slice := g.Slice[int]{1, 2, 3, 4, 5}

	// Test a valid range within bounds
	subSlice := slice.SubSlice(1, 4)
	expected := g.Slice[int]{2, 3, 4}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with a single negative index
	subSlice = slice.SubSlice(-2, slice.Len())
	expected = g.Slice[int]{4, 5}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with a negative start and end index
	subSlice = slice.SubSlice(-3, -1)
	expected = g.Slice[int]{3, 4}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}
}

func TestSubSliceOutOfBoundsStartIndex(t *testing.T) {
	slice := g.Slice[int]{1, 2, 3, 4, 5}

	// Test with start index beyond slice length (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for start index beyond slice length, but no panic occurred")
		}
	}()
	_ = slice.SubSlice(10, slice.Len())
}

func TestSubSliceOutOfBoundsNegativeStartIndex(t *testing.T) {
	slice := g.Slice[int]{1, 2, 3, 4, 5}

	// Test with a negative start index beyond slice length (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for negative start index beyond slice length, but no panic occurred")
		}
	}()
	_ = slice.SubSlice(-10, slice.Len())
}

func TestSubSliceOutOfBoundsEndIndex(t *testing.T) {
	slice := g.Slice[int]{1, 2, 3, 4, 5}

	// Test with an end index beyond slice length (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for end index beyond slice length, but no panic occurred")
		}
	}()
	_ = slice.SubSlice(2, 10)
}

func TestGrowSlice(t *testing.T) {
	// Initialize a slice with some elements.
	initialSlice := g.SliceOf(1, 2, 3)

	// Check the initial capacity of the slice.
	initialCapacity := initialSlice.Cap()

	// Grow the slice to accommodate more elements.
	newCapacity := initialCapacity + 5
	grownSlice := initialSlice.Grow(newCapacity - initialCapacity)

	// Check if the capacity of the grown slice is as expected.
	if grownSlice.Cap() != newCapacity {
		t.Errorf("Grow method failed: Expected capacity %d, got %d", newCapacity, grownSlice.Cap())
	}

	// Append new elements to the grown slice.
	for i := 0; i < 5; i++ {
		grownSlice = grownSlice.Append(i + 4)
	}

	// Check if the length of the grown slice is correct.
	if grownSlice.Len() != newCapacity {
		t.Errorf("Grow method failed: Expected length %d, got %d", newCapacity, grownSlice.Len())
	}
}

func TestSliceNotEmpty(t *testing.T) {
	// Test case 1: Slice with elements
	sl1 := g.SliceOf(1, 2, 3)
	if !sl1.NotEmpty() {
		t.Errorf("Test case 1 failed: Expected slice to be not empty")
	}

	// Test case 2: Empty slice
	sl2 := g.NewSlice[g.Int]()
	if sl2.NotEmpty() {
		t.Errorf("Test case 2 failed: Expected slice to be empty")
	}
}

func TestSliceAppendInPlace(t *testing.T) {
	// Create a slice with initial elements
	initialSlice := g.Slice[int]{1, 2, 3}

	// Append additional elements using AppendInPlace
	initialSlice.AppendInPlace(4, 5, 6)

	// Verify that the slice has the expected elements
	expected := g.Slice[int]{1, 2, 3, 4, 5, 6}
	if !initialSlice.Eq(expected) {
		t.Errorf("AppendInPlace failed. Expected: %v, Got: %v", expected, initialSlice)
	}
}

func TestSliceString(t *testing.T) {
	// Create a slice with some elements
	sl := g.SliceOf(1, 2, 3, 4, 5)

	// Define the expected string representation
	expected := "Slice[1, 2, 3, 4, 5]"

	// Get the string representation using the String method
	result := sl.String()

	// Compare the result with the expected value
	if result != expected {
		t.Errorf("Slice String method failed. Expected: %s, Got: %s", expected, result)
	}
}

func TestSliceEq(t *testing.T) {
	// Create two slices with the same elements
	sl1 := g.SliceOf(1, 2, 3, 4, 5)
	sl2 := g.SliceOf(1, 2, 3, 4, 5)

	// Check if the slices are equal
	if !sl1.Eq(sl2) {
		t.Error("Slice Eq method failed. Expected slices to be equal.")
	}

	// Create another slice with different elements
	sl3 := g.SliceOf(1, 2, 3)

	// Check if the slices are not equal
	if sl1.Eq(sl3) {
		t.Error("Slice Eq method failed. Expected slices to be not equal.")
	}

	// Create another slice with different elements
	sl4 := g.SliceOf(1, 2, 3, 4, 6)

	// Check if the slices are not equal
	if sl1.Eq(sl4) {
		t.Error("Slice Eq method failed. Expected slices to be not equal.")
	}
}

func TestSliceLastIndex(t *testing.T) {
	// Create a slice with some elements
	sl := g.SliceOf(1, 2, 3, 4, 5)

	// Get the last index of the slice
	lastIndex := sl.LastIndex()

	// Check if the last index is correct
	expectedLastIndex := sl.Len() - 1
	if lastIndex != expectedLastIndex {
		t.Errorf("Slice LastIndex method failed. Expected: %d, Got: %d", expectedLastIndex, lastIndex)
	}

	// Create an empty slice
	emptySlice := g.NewSlice[int]()

	// Get the last index of the empty slice
	emptyLastIndex := emptySlice.LastIndex()

	// Check if the last index of an empty slice is 0
	if emptyLastIndex != 0 {
		t.Errorf("Slice LastIndex method failed for empty slice. Expected: 0, Got: %d", emptyLastIndex)
	}
}

func TestSliceRandom(t *testing.T) {
	// Create a slice with some elements
	sl := g.SliceOf(1, 2, 3, 4, 5)

	// Get a random element from the slice
	randomElement := sl.Random()

	// Check if the random element is within the slice
	found := false
	sl.Iter().ForEach(func(v int) {
		if v == randomElement {
			found = true
		}
	})

	if !found {
		t.Errorf("Slice Random method failed. Random element %d not found in the slice", randomElement)
	}

	// Test for an empty slice
	emptySlice := g.NewSlice[int]()

	// Get a random element from the empty slice
	emptyRandomElement := emptySlice.Random()

	// Check if the random element from an empty slice is zero value
	if emptyRandomElement != 0 {
		t.Errorf("Slice Random method failed for empty slice. Expected: 0, Got: %d", emptyRandomElement)
	}
}

func TestSliceMaxMin(t *testing.T) {
	// Test cases for Int
	testMaxMin(t, g.SliceOf[g.Int](3, 1, 4, 1, 5), g.Int(5), g.Int(1))
	testMaxMin(t, g.SliceOf(3, 1, 4, 1, 5), 5, 1)

	// Test cases for Float
	testMaxMin(t, g.SliceOf[g.Float](3.14, 1.23, 4.56, 1.01, 5.67), g.Float(5.67), g.Float(1.01))
	testMaxMin(t, g.SliceOf(3.14, 1.23, 4.56, 1.01, 5.67), 5.67, 1.01)

	// Test cases for String
	testMaxMin(t, g.SliceOf[g.String]("apple", "banana", "orange", "grape"), g.String("orange"), g.String("apple"))
	testMaxMin(t, g.SliceOf("apple", "banana", "orange", "grape"), "orange", "apple")

	// Add more test cases for other types as needed
}

func testMaxMin[T comparable](t *testing.T, sl g.Slice[T], expectedMax, expectedMin T) {
	// Test Max method
	maxElement := sl.Max()
	if maxElement != expectedMax {
		t.Errorf("Slice Max method failed for type %T. Expected: %v, Got: %v", sl[0], expectedMax, maxElement)
	}

	// Test Min method
	minElement := sl.Min()
	if minElement != expectedMin {
		t.Errorf("Slice Min method failed for type %T. Expected: %v, Got: %v", sl[0], expectedMin, minElement)
	}

	// Test for an empty slice
	emptySlice := g.NewSlice[int]()
	emptyMaxElement := emptySlice.Max()
	if emptyMaxElement != 0 {
		t.Errorf("Slice Max method failed for empty slice. Expected: 0, Got: %v", emptyMaxElement)
	}
	emptyMinElement := emptySlice.Min()
	if emptyMinElement != 0 {
		t.Errorf("Slice Min method failed for empty slice. Expected: 0, Got: %v", emptyMinElement)
	}
}

func TestSliceAddUniqueInPlace(t *testing.T) {
	// Test cases for Int
	testAddUniqueInPlace(t, g.SliceOf(1, 2, 3, 4, 5), []int{3, 4, 5, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7})

	// Test cases for Float
	testAddUniqueInPlace(
		t,
		g.SliceOf(1.1, 2.2, 3.3, 4.4, 5.5),
		[]float64{3.3, 4.4, 5.5, 6.6, 7.7},
		[]float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7},
	)

	// Test cases for String
	testAddUniqueInPlace(
		t,
		g.SliceOf("apple", "banana", "orange", "grape"),
		[]string{"orange", "grape", "kiwi"},
		[]string{"apple", "banana", "orange", "grape", "kiwi"},
	)

	// Add more test cases for other types as needed
}

func testAddUniqueInPlace[T comparable](t *testing.T, sl g.Slice[T], elems, expected []T) {
	sl.AddUniqueInPlace(elems...)
	if !sl.Eq(g.SliceOf(expected...)) {
		t.Errorf("Slice AddUniqueInPlace method failed for type %T. Expected: %v, Got: %v", sl[0], expected, sl)
	}
}

func TestSliceAsAny(t *testing.T) {
	// Test cases for Int
	testSliceAsAny(t, g.SliceOf(1, 2, 3), []any{1, 2, 3})

	// Test cases for Float
	testSliceAsAny(t, g.SliceOf(1.1, 2.2, 3.3), []any{1.1, 2.2, 3.3})

	// Test cases for String
	testSliceAsAny(t, g.SliceOf("apple", "banana", "orange"), []any{"apple", "banana", "orange"})

	// Add more test cases for other types as needed
}

func testSliceAsAny[T any](t *testing.T, sl g.Slice[T], expected []any) {
	result := sl.AsAny()
	if !result.Eq(g.SliceOf(expected...)) {
		t.Errorf("Slice AsAny method failed for type %T. Expected: %v, Got: %v", sl[0], expected, result)
	}
}

func TestSliceLess(t *testing.T) {
	// Test 1: Comparing two Int elements
	slice1 := g.Slice[g.Int]{1, 2, 3}
	result1 := slice1.Less(0, 1)

	expectedResult1 := true
	if result1 != expectedResult1 {
		t.Errorf("Test 1: Expected %t, got %t", expectedResult1, result1)
	}

	// Test 2: Comparing two string elements
	slice2 := g.Slice[g.String]{"apple", "banana", "orange"}
	result2 := slice2.Less(0, 1)

	expectedResult2 := true
	if result2 != expectedResult2 {
		t.Errorf("Test 2: Expected %t, got %t", expectedResult2, result2)
	}

	// Test 3: Comparing Int and int elements
	slice3 := g.Slice[any]{1, 2.5, 3}
	result3 := slice3.Less(0, 1)

	expectedResult3 := false
	if result3 != expectedResult3 {
		t.Errorf("Test 3: Expected %t, got %t", expectedResult3, result3)
	}

	// Test 4: Comparing two float elements
	slice4 := g.Slice[g.Float]{1.2, 2.5, 3.1}
	result4 := slice4.Less(0, 1)

	expectedResult4 := true
	if result4 != expectedResult4 {
		t.Errorf("Test 4: Expected %t, got %t", expectedResult4, result4)
	}
}
