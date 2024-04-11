package g_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/enetx/g"
	"github.com/enetx/g/f"
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
	// Test case: Function returns an index for known types (Int)
	slInt := g.Slice[g.Int]{1, 2, 3, 4, 5}
	index := slInt.Index(3)
	if index != 2 {
		t.Errorf("Expected index 2, got %d", index)
	}

	// Test case: Function returns -1 for unknown types (String)
	slString := g.Slice[g.String]{"a", "b", "c"}
	index = slString.Index("d")
	if index != -1 {
		t.Errorf("Expected index -1, got %d", index)
	}

	// Test case: Function returns an index for known types (Float)
	slFloat := g.Slice[g.Float]{1.1, 2.2, 3.3}
	index = slFloat.Index(2.2)
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns -1 for empty slice (Bool)
	emptySliceBool := g.Slice[bool]{}
	index = emptySliceBool.Index(true)
	if index != -1 {
		t.Errorf("Expected index -1 for empty slice, got %d", index)
	}

	// Test case: Function returns an index for known types (Byte)
	slByte := g.Slice[byte]{byte('a'), byte('b'), byte('c')}
	index = slByte.Index(byte('b'))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (String)
	slString2 := g.Slice[string]{"apple", "banana", "cherry"}
	index = slString2.Index("banana")
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int)
	slInt2 := g.Slice[int]{10, 20, 30}
	index = slInt2.Index(20)
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int8)
	slInt8 := g.Slice[int8]{1, 2, 3}
	index = slInt8.Index(int8(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int16)
	slInt16 := g.Slice[int16]{1, 2, 3}
	index = slInt16.Index(int16(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int32)
	slInt32 := g.Slice[int32]{1, 2, 3}
	index = slInt32.Index(int32(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int64)
	slInt64 := g.Slice[int64]{1, 2, 3}
	index = slInt64.Index(int64(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint)
	slUint := g.Slice[uint]{1, 2, 3, 4, 5}
	index = slUint.Index(uint(3))
	if index != 2 {
		t.Errorf("Expected index 2, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint8)
	slUint8 := g.Slice[uint8]{1, 2, 3}
	index = slUint8.Index(uint8(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint16)
	slUint16 := g.Slice[uint16]{1, 2, 3}
	index = slUint16.Index(uint16(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint32)
	slUint32 := g.Slice[uint32]{1, 2, 3}
	index = slUint32.Index(uint32(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint64)
	slUint64 := g.Slice[uint64]{1, 2, 3}
	index = slUint64.Index(uint64(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Float32)
	slFloat32 := g.Slice[float32]{1.1, 2.2, 3.3}
	index = slFloat32.Index(float32(2.2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Float64)
	slFloat64 := g.Slice[float64]{1.1, 2.2, 3.3}
	index = slFloat64.Index(float64(2.2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}
}

func TestSliceIndexFunc(t *testing.T) {
	// Define a custom slice type
	type customType struct {
		Value int
	}

	// Create a slice with custom type
	slCustom := g.Slice[customType]{{Value: 1}, {Value: 2}, {Value: 3}}

	// Test case: Function returns an index for custom types using IndexFunc
	index := slCustom.Index(customType{Value: 2})
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns -1 for custom types not found using IndexFunc
	index = slCustom.Index(customType{Value: 4})
	if index != -1 {
		t.Errorf("Expected index -1, got %d", index)
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
	// Test case: Function returns true for equal slices of known types (Int)
	slInt1 := g.Slice[g.Int]{1, 2, 3}
	slInt2 := g.Slice[g.Int]{1, 2, 3}
	if !slInt1.Eq(slInt2) {
		t.Errorf("Test 1: Expected slices to be equal")
	}

	// Test case: Function returns false for unequal slices of known types (String)
	slString1 := g.Slice[g.String]{"a", "b", "c"}
	slString2 := g.Slice[g.String]{"a", "x", "c"}
	if slString1.Eq(slString2) {
		t.Errorf("Test 2: Expected slices to be unequal")
	}

	// Test case: Function returns true for empty slices
	emptySlice1 := g.Slice[g.Float]{}
	emptySlice2 := g.Slice[g.Float]{}
	if !emptySlice1.Eq(emptySlice2) {
		t.Errorf("Test 3: Expected empty slices to be equal")
	}

	// Test case: Function returns false for slices of different lengths
	slFloat1 := g.Slice[g.Float]{1.1, 2.2, 3.3}
	slFloat2 := g.Slice[g.Float]{1.1, 2.2}
	if slFloat1.Eq(slFloat2) {
		t.Errorf("Test 4: Expected slices of different lengths to be unequal")
	}

	// Test case: Function returns true for equal slices of string type
	slString3 := g.Slice[string]{"apple", "banana", "cherry"}
	slString4 := g.Slice[string]{"apple", "banana", "cherry"}
	if !slString3.Eq(slString4) {
		t.Errorf("Test 5: Expected slices to be equal")
	}

	// Test case: Function returns false for unequal slices of int type
	slInt3 := g.Slice[int]{10, 20, 30}
	slInt4 := g.Slice[int]{10, 20, 40}
	if slInt3.Eq(slInt4) {
		t.Errorf("Test 6: Expected slices to be unequal")
	}

	// Test case: Function returns true for equal slices of float64 type
	slFloat64_1 := g.Slice[float64]{1.1, 2.2, 3.3}
	slFloat64_2 := g.Slice[float64]{1.1, 2.2, 3.3}
	if !slFloat64_1.Eq(slFloat64_2) {
		t.Errorf("Test 7: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of bool type
	slBool1 := g.Slice[bool]{true, false, true}
	slBool2 := g.Slice[bool]{true, false, true}
	if !slBool1.Eq(slBool2) {
		t.Errorf("Test 8: Expected slices to be equal")
	}

	// Test case: Function returns false for unequal slices of byte type
	slByte1 := g.Slice[byte]{1, 2, 3}
	slByte2 := g.Slice[byte]{1, 2, 4}
	if slByte1.Eq(slByte2) {
		t.Errorf("Test 9: Expected slices to be unequal")
	}

	// Test case: Function returns true for equal slices of int8 type
	slInt81 := g.Slice[int8]{1, 2, 3}
	slInt82 := g.Slice[int8]{1, 2, 3}
	if !slInt81.Eq(slInt82) {
		t.Errorf("Test 10: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of int16 type
	slInt161 := g.Slice[int16]{1, 2, 3}
	slInt162 := g.Slice[int16]{1, 2, 3}
	if !slInt161.Eq(slInt162) {
		t.Errorf("Test 11: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of int32 type
	slInt321 := g.Slice[int32]{1, 2, 3}
	slInt322 := g.Slice[int32]{1, 2, 3}
	if !slInt321.Eq(slInt322) {
		t.Errorf("Test 12: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of int64 type
	slInt641 := g.Slice[int64]{1, 2, 3}
	slInt642 := g.Slice[int64]{1, 2, 3}
	if !slInt641.Eq(slInt642) {
		t.Errorf("Test 13: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint type
	slUint1 := g.Slice[uint]{1, 2, 3}
	slUint2 := g.Slice[uint]{1, 2, 3}
	if !slUint1.Eq(slUint2) {
		t.Errorf("Test 14: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint8 type
	slUint81 := g.Slice[uint8]{1, 2, 3}
	slUint82 := g.Slice[uint8]{1, 2, 3}
	if !slUint81.Eq(slUint82) {
		t.Errorf("Test 15: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint16 type
	slUint161 := g.Slice[uint16]{1, 2, 3}
	slUint162 := g.Slice[uint16]{1, 2, 3}
	if !slUint161.Eq(slUint162) {
		t.Errorf("Test 16: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint32 type
	slUint321 := g.Slice[uint32]{1, 2, 3}
	slUint322 := g.Slice[uint32]{1, 2, 3}
	if !slUint321.Eq(slUint322) {
		t.Errorf("Test 17: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint64 type
	slUint641 := g.Slice[uint64]{1, 2, 3}
	slUint642 := g.Slice[uint64]{1, 2, 3}
	if !slUint641.Eq(slUint642) {
		t.Errorf("Test 18: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of float32 type
	slFloat321 := g.Slice[float32]{1.1, 2.2, 3.3}
	slFloat322 := g.Slice[float32]{1.1, 2.2, 3.3}
	if !slFloat321.Eq(slFloat322) {
		t.Errorf("Test 19: Expected slices to be equal")
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
	// Test case 1: Comparing two Int elements
	sliceInt := g.Slice[g.Int]{1, 2, 3}
	resultInt := sliceInt.Less(0, 1)
	expectedResultInt := true
	if resultInt != expectedResultInt {
		t.Errorf("Test 1: Expected %t, got %t", expectedResultInt, resultInt)
	}

	// Test case 2: Comparing two String elements
	sliceString := g.Slice[g.String]{"apple", "banana", "orange"}
	resultString := sliceString.Less(0, 1)
	expectedResultString := true
	if resultString != expectedResultString {
		t.Errorf("Test 2: Expected %t, got %t", expectedResultString, resultString)
	}

	// Test case 3: Comparing Int and int elements
	sliceMixed := g.Slice[any]{1, 2.5, 3}
	resultMixed := sliceMixed.Less(0, 1)
	expectedResultMixed := false
	if resultMixed != expectedResultMixed {
		t.Errorf("Test 3: Expected %t, got %t", expectedResultMixed, resultMixed)
	}

	// Test case 4: Comparing two Float elements
	sliceFloat := g.Slice[g.Float]{1.2, 2.5, 3.1}
	resultFloat := sliceFloat.Less(0, 1)
	expectedResultFloat := true
	if resultFloat != expectedResultFloat {
		t.Errorf("Test 4: Expected %t, got %t", expectedResultFloat, resultFloat)
	}

	// Test case 5: Comparing two Bool elements
	sliceBool := g.Slice[bool]{true, false}
	resultBool := sliceBool.Less(0, 1)
	expectedResultBool := false
	if resultBool != expectedResultBool {
		t.Errorf("Test 5: Expected %t, got %t", expectedResultBool, resultBool)
	}

	// Test case 6: Comparing two Uint elements
	sliceUint := g.Slice[uint]{1, 2, 3}
	resultUint := sliceUint.Less(0, 1)
	expectedResultUint := true
	if resultUint != expectedResultUint {
		t.Errorf("Test 6: Expected %t, got %t", expectedResultUint, resultUint)
	}

	// Test case 7: Comparing two Uint8 elements
	sliceUint8 := g.Slice[uint8]{1, 2, 3}
	resultUint8 := sliceUint8.Less(0, 1)
	expectedResultUint8 := true
	if resultUint8 != expectedResultUint8 {
		t.Errorf("Test 7: Expected %t, got %t", expectedResultUint8, resultUint8)
	}

	// Test case 8: Comparing two Uint16 elements
	sliceUint16 := g.Slice[uint16]{1, 2, 3}
	resultUint16 := sliceUint16.Less(0, 1)
	expectedResultUint16 := true
	if resultUint16 != expectedResultUint16 {
		t.Errorf("Test 8: Expected %t, got %t", expectedResultUint16, resultUint16)
	}

	// Test case 9: Comparing two Uint32 elements
	sliceUint32 := g.Slice[uint32]{1, 2, 3}
	resultUint32 := sliceUint32.Less(0, 1)
	expectedResultUint32 := true
	if resultUint32 != expectedResultUint32 {
		t.Errorf("Test 9: Expected %t, got %t", expectedResultUint32, resultUint32)
	}

	// Test case 10: Comparing two Uint64 elements
	sliceUint64 := g.Slice[uint64]{1, 2, 3}
	resultUint64 := sliceUint64.Less(0, 1)
	expectedResultUint64 := true
	if resultUint64 != expectedResultUint64 {
		t.Errorf("Test 10: Expected %t, got %t", expectedResultUint64, resultUint64)
	}

	// Test case 11: Comparing two Int8 elements
	sliceInt8 := g.Slice[int8]{1, 2, 3}
	resultInt8 := sliceInt8.Less(0, 1)
	expectedResultInt8 := true
	if resultInt8 != expectedResultInt8 {
		t.Errorf("Test 11: Expected %t, got %t", expectedResultInt8, resultInt8)
	}

	// Test case 12: Comparing two Int16 elements
	sliceInt16 := g.Slice[int16]{1, 2, 3}
	resultInt16 := sliceInt16.Less(0, 1)
	expectedResultInt16 := true
	if resultInt16 != expectedResultInt16 {
		t.Errorf("Test 12: Expected %t, got %t", expectedResultInt16, resultInt16)
	}

	// Test case 13: Comparing two Int32 elements
	sliceInt32 := g.Slice[int32]{1, 2, 3}
	resultInt32 := sliceInt32.Less(0, 1)
	expectedResultInt32 := true
	if resultInt32 != expectedResultInt32 {
		t.Errorf("Test 13: Expected %t, got %t", expectedResultInt32, resultInt32)
	}

	// Test case 14: Comparing two Int64 elements
	sliceInt64 := g.Slice[int64]{1, 2, 3}
	resultInt64 := sliceInt64.Less(0, 1)
	expectedResultInt64 := true
	if resultInt64 != expectedResultInt64 {
		t.Errorf("Test 14: Expected %t, got %t", expectedResultInt64, resultInt64)
	}

	// Test case 15: Comparing two Float32 elements
	sliceFloat32 := g.Slice[float32]{1.2, 2.5, 3.1}
	resultFloat32 := sliceFloat32.Less(0, 1)
	expectedResultFloat32 := true
	if resultFloat32 != expectedResultFloat32 {
		t.Errorf("Test 15: Expected %t, got %t", expectedResultFloat32, resultFloat32)
	}

	// Test case 16: Comparing two Float64 elements
	sliceFloat64 := g.Slice[float64]{1.2, 2.5, 3.1}
	resultFloat64 := sliceFloat64.Less(0, 1)
	expectedResultFloat64 := true
	if resultFloat64 != expectedResultFloat64 {
		t.Errorf("Test 16: Expected %t, got %t", expectedResultFloat64, resultFloat64)
	}
}

func TestSliceContainsBy(t *testing.T) {
	// Test case 1: Slice contains the element that satisfies the provided function
	sl1 := g.Slice[g.Int]{1, 2, 3, 4, 5}
	contains1 := sl1.ContainsBy(f.Eq(g.Int(3)))

	if !contains1 {
		t.Errorf("Test 1: Expected true, got false")
	}

	// Test case 2: Slice does not contain the element that satisfies the provided function
	sl2 := g.Slice[g.String]{"apple", "banana", "cherry"}
	contains2 := sl2.ContainsBy(f.Eq(g.String("orange")))

	if contains2 {
		t.Errorf("Test 2: Expected false, got true")
	}

	// Test case 3: Slice contains the element that satisfies the provided function (using custom struct)
	type Person struct {
		Name string
		Age  int
	}

	sl3 := g.Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}, {Name: "Charlie", Age: 35}}

	contains3 := sl3.ContainsBy(func(x Person) bool { return x.Name == "Bob" && x.Age == 25 })
	if !contains3 {
		t.Errorf("Test 3: Expected true, got false")
	}
}

func TestSliceEqBy(t *testing.T) {
	// Test case 1: Slices are equal using the equality function
	sl1 := g.Slice[g.Int]{1, 2, 3}
	sl2 := g.Slice[g.Int]{1, 2, 3}

	eq1 := sl1.EqBy(sl2, func(x, y g.Int) bool { return x.Eq(y) })

	if !eq1 {
		t.Errorf("Test 1: Expected true, got false")
	}

	// Test case 2: Slices are not equal using the equality function
	sl3 := g.Slice[g.String]{"apple", "banana", "cherry"}
	sl4 := g.Slice[g.String]{"apple", "orange", "cherry"}

	eq2 := sl3.EqBy(sl4, func(x, y g.String) bool { return x.Eq(y) })

	if eq2 {
		t.Errorf("Test 2: Expected false, got true")
	}

	// Test case 3: Slices are equal using the equality function (using custom struct)
	type Person struct {
		Name string
		Age  int
	}

	sl5 := g.Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	sl6 := g.Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}

	eq3 := sl5.EqBy(sl6, func(x, y Person) bool {
		return x.Name == y.Name && x.Age == y.Age
	})

	if !eq3 {
		t.Errorf("Test 3: Expected true, got false")
	}
}

func TestSliceIndexBy(t *testing.T) {
	// Test case 1: Element satisfying the custom comparison function exists in the slice
	sl1 := g.Slice[int]{1, 2, 3, 4, 5}
	index1 := sl1.IndexBy(f.Eq(3))

	expectedIndex1 := 2
	if index1 != expectedIndex1 {
		t.Errorf("Test 1: Expected index %d, got %d", expectedIndex1, index1)
	}

	// Test case 2: Element satisfying the custom comparison function doesn't exist in the slice
	sl2 := g.Slice[string]{"apple", "banana", "cherry"}
	index2 := sl2.IndexBy(f.Eq("orange"))

	expectedIndex2 := -1
	if index2 != expectedIndex2 {
		t.Errorf("Test 2: Expected index %d, got %d", expectedIndex2, index2)
	}

	// Test case 3: Element satisfying the custom comparison function exists in the slice (using custom struct)
	type Person struct {
		Name string
		Age  int
	}

	sl3 := g.Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	index3 := sl3.IndexBy(func(x Person) bool { return x.Name == "Bob" && x.Age == 25 })

	expectedIndex3 := 1
	if index3 != expectedIndex3 {
		t.Errorf("Test 3: Expected index %d, got %d", expectedIndex3, index3)
	}
}
