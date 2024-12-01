package g_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

func TestSliceUnpack(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[int]
		vars     []*int
		expected []int
	}{
		{
			name:     "Unpack with valid indices",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			vars:     []*int{new(int), new(int), new(int)},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Unpack with invalid indices",
			slice:    Slice[int]{1, 2, 3},
			vars:     []*int{new(int), new(int), new(int), new(int)},
			expected: []int{1, 2, 3, 0}, // Expecting zero value for the fourth variable
		},
		{
			name:     "Unpack with empty slice",
			slice:    Slice[int]{},
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
	slice := SliceOf(1, 2, 3, 4, 5, 6, 7, 8, 9)

	testCases := []struct {
		start    Int
		end      Int
		step     Int
		expected Slice[int]
	}{
		{1, 7, 2, SliceOf(2, 4, 6)},
		{0, 9, 3, SliceOf(1, 4, 7)},
		{2, 6, 1, SliceOf(3, 4, 5, 6)},
		{0, 9, 2, SliceOf(1, 3, 5, 7, 9)},
		{0, 9, 4, SliceOf(1, 5, 9)},
		{6, 1, -2, SliceOf(7, 5, 3)},
		{8, 1, -3, SliceOf(9, 6, 3)},
		{8, 0, -2, SliceOf(9, 7, 5, 3)},
		{8, 0, -1, SliceOf(9, 8, 7, 6, 5, 4, 3, 2)},
		{8, 0, -4, SliceOf(9, 5)},
		{-1, -6, -2, SliceOf(9, 7, 5)},
		{-2, -9, -3, SliceOf(8, 5, 2)},
		{-1, -8, -2, SliceOf(9, 7, 5, 3)},
		{-3, -10, -2, SliceOf(7, 5, 3, 1)},
		{-1, -10, -1, SliceOf(9, 8, 7, 6, 5, 4, 3, 2, 1)},
		{-5, -1, -1, Slice[int]{}},
		{-1, -5, -1, SliceOf(9, 8, 7, 6)},
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

// func TestSortInts(t *testing.T) {
// 	slice := Slice[int]{5, 2, 8, 1, 6}
// 	slice.Sort()

// 	expected := Slice[int]{1, 2, 5, 6, 8}

// 	if !reflect.DeepEqual(slice, expected) {
// 		t.Errorf("Expected %v but got %v", expected, slice)
// 	}
// }

// func TestSortStrings(t *testing.T) {
// 	slice := Slice[string]{"apple", "orange", "banana", "grape"}
// 	slice.Sort()

// 	expected := Slice[string]{"apple", "banana", "grape", "orange"}

// 	if !reflect.DeepEqual(slice, expected) {
// 		t.Errorf("Expected %v but got %v", expected, slice)
// 	}
// }

// func TestSliceSortFloats(t *testing.T) {
// 	slice := Slice[float64]{5.6, 2.3, 8.9, 1.2, 6.7}
// 	slice.Sort()

// 	expected := Slice[float64]{1.2, 2.3, 5.6, 6.7, 8.9}

// 	if !reflect.DeepEqual(slice, expected) {
// 		t.Errorf("Expected %v but got %v", expected, slice)
// 	}
// }

func TestSliceInsert(t *testing.T) {
	// Test insertion in the middle
	slice := Slice[string]{"a", "b", "c", "d"}
	newSlice := slice.Insert(2, "e", "f")
	expected := Slice[string]{"a", "b", "e", "f", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion at the start
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Insert(0, "x", "y")
	expected = Slice[string]{"x", "y", "a", "b", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(0) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion at the end
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Insert(slice.Len(), "x", "y")
	expected = Slice[string]{"a", "b", "c", "d", "x", "y"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(end) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion with negative index
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Insert(-2, "x", "y")
	expected = Slice[string]{"a", "b", "x", "y", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(-2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test insertion at index 0 with an empty slice
	slice = Slice[string]{}
	newSlice = slice.Insert(0, "x", "y")
	expected = Slice[string]{"x", "y"}
	if !newSlice.Eq(expected) {
		t.Errorf("Insert(0) with empty slice failed. Expected %v, but got %v", expected, newSlice)
	}
}

func TestSliceInsertInPlace(t *testing.T) {
	// Test insertion in the middle
	slice := Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(2, "e", "f")
	expected := Slice[string]{"a", "b", "e", "f", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion at the start
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(0, "x", "y")
	expected = Slice[string]{"x", "y", "a", "b", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(0) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion at the end
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(slice.Len(), "x", "y")
	expected = Slice[string]{"a", "b", "c", "d", "x", "y"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(end) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion with negative index
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.InsertInPlace(-2, "x", "y")
	expected = Slice[string]{"a", "b", "x", "y", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(-2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test insertion at index 0 with an empty slice
	slice = Slice[string]{}
	slice.InsertInPlace(0, "x", "y")
	expected = Slice[string]{"x", "y"}
	if !slice.Eq(expected) {
		t.Errorf("InsertInPlace(0) with empty slice failed. Expected %v, but got %v", expected, slice)
	}
}

func TestSliceToSlice(t *testing.T) {
	sl := NewSlice[int]().Append(1, 2, 3, 4, 5)
	slice := sl.Std()

	if len(slice) != sl.Len().Std() {
		t.Errorf("Expected length %d, but got %d", sl.Len(), len(slice))
	}

	for i, v := range sl {
		if v != slice[i] {
			t.Errorf("Expected value %d at index %d, but got %d", v, i, slice[i])
		}
	}
}

func TestSliceShuffle(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sl.Shuffle()

	if sl.Len() != 10 {
		t.Error("Expected length of 10, got ", sl.Len())
	}
}

func TestSliceReverse(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	sl.Reverse()

	if !reflect.DeepEqual(sl, Slice[int]{5, 4, 3, 2, 1}) {
		t.Errorf("Expected %v, got %v", Slice[int]{5, 4, 3, 2, 1}, sl)
	}
}

func TestSliceIndex(t *testing.T) {
	// Test case: Function returns an index for known types (Int)
	slInt := Slice[Int]{1, 2, 3, 4, 5}
	index := slInt.Index(3)
	if index != 2 {
		t.Errorf("Expected index 2, got %d", index)
	}

	// Test case: Function returns -1 for unknown types (String)
	slString := Slice[String]{"a", "b", "c"}
	index = slString.Index("d")
	if index != -1 {
		t.Errorf("Expected index -1, got %d", index)
	}

	// Test case: Function returns an index for known types (Float)
	slFloat := Slice[Float]{1.1, 2.2, 3.3}
	index = slFloat.Index(2.2)
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns -1 for empty slice (Bool)
	emptySliceBool := Slice[bool]{}
	index = emptySliceBool.Index(true)
	if index != -1 {
		t.Errorf("Expected index -1 for empty slice, got %d", index)
	}

	// Test case: Function returns an index for known types (Byte)
	slByte := Slice[byte]{byte('a'), byte('b'), byte('c')}
	index = slByte.Index(byte('b'))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (String)
	slString2 := Slice[string]{"apple", "banana", "cherry"}
	index = slString2.Index("banana")
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int)
	slInt2 := Slice[int]{10, 20, 30}
	index = slInt2.Index(20)
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int8)
	slInt8 := Slice[int8]{1, 2, 3}
	index = slInt8.Index(int8(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int16)
	slInt16 := Slice[int16]{1, 2, 3}
	index = slInt16.Index(int16(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int32)
	slInt32 := Slice[int32]{1, 2, 3}
	index = slInt32.Index(int32(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Int64)
	slInt64 := Slice[int64]{1, 2, 3}
	index = slInt64.Index(int64(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint)
	slUint := Slice[uint]{1, 2, 3, 4, 5}
	index = slUint.Index(uint(3))
	if index != 2 {
		t.Errorf("Expected index 2, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint8)
	slUint8 := Slice[uint8]{1, 2, 3}
	index = slUint8.Index(uint8(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint16)
	slUint16 := Slice[uint16]{1, 2, 3}
	index = slUint16.Index(uint16(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint32)
	slUint32 := Slice[uint32]{1, 2, 3}
	index = slUint32.Index(uint32(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Uint64)
	slUint64 := Slice[uint64]{1, 2, 3}
	index = slUint64.Index(uint64(2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Float32)
	slFloat32 := Slice[float32]{1.1, 2.2, 3.3}
	index = slFloat32.Index(float32(2.2))
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}

	// Test case: Function returns an index for known types (Float64)
	slFloat64 := Slice[float64]{1.1, 2.2, 3.3}
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
	slCustom := Slice[customType]{{Value: 1}, {Value: 2}, {Value: 3}}

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
	sl := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
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
	sl := Slice[int]{1, 2, 3}
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
	sl1 := NewSlice[int]().Append(3, 1, 4, 1, 5)
	expected1 := NewSlice[int]().Append(1, 1, 3, 4, 5)

	sl1.SortBy(func(a, b int) cmp.Ordering { return cmp.Cmp(a, b) })

	if !sl1.Eq(expected1) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected1, sl1)
	}

	sl2 := NewSlice[string]().Append("foo", "bar", "baz")
	expected2 := NewSlice[string]().Append("foo", "baz", "bar")

	sl2.SortBy(func(a, b string) cmp.Ordering { return cmp.Cmp(b, a) })

	if !sl2.Eq(expected2) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected2, sl2)
	}

	sl3 := NewSlice[int]()
	expected3 := NewSlice[int]()

	sl3.SortBy(func(a, b int) cmp.Ordering { return cmp.Cmp(a, b) })

	if !sl3.Eq(expected3) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected3, sl3)
	}
}

func TestSliceJoin(t *testing.T) {
	sl := Slice[string]{"1", "2", "3", "4", "5"}
	str := sl.Join(",")

	if !strings.EqualFold("1,2,3,4,5", str.Std()) {
		t.Errorf("Expected 1,2,3,4,5, got %s", str.Std())
	}
}

func TestSliceToStringSlice(t *testing.T) {
	sl := Slice[int]{1, 2, 3}
	result := sl.ToStringSlice()
	expected := []string{"1", "2", "3"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSliceAdd(t *testing.T) {
	sl := Slice[int]{1, 2, 3}.Append(4, 5, 6)

	if !reflect.DeepEqual(sl, Slice[int]{1, 2, 3, 4, 5, 6}) {
		t.Error("Add failed")
	}
}

func TestSliceClone(t *testing.T) {
	sl := Slice[int]{1, 2, 3}
	slClone := sl.Clone()

	if !sl.Eq(slClone) {
		t.Errorf("Clone() failed, expected %v, got %v", sl, slClone)
	}
}

func TestSliceCut(t *testing.T) {
	slice := Slice[int]{1, 2, 3, 4, 5}

	// Test normal range
	newSlice := slice.Cut(1, 3)
	expected := Slice[int]{1, 4, 5}
	if !newSlice.Eq(expected) {
		t.Errorf("Cut(1, 3) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test range with negative indices
	newSlice = slice.Cut(-3, -2)
	expected = Slice[int]{1, 2, 4, 5}
	if !newSlice.Eq(expected) {
		t.Errorf("Cut(-3, -2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test empty range
	newSlice = slice.Cut(0, 5)
	expected = Slice[int]{}
	if !newSlice.Eq(expected) {
		t.Errorf("Cut(3, 2) failed. Expected %v, but got %v", expected, newSlice)
	}
}

func TestSliceCutInPlace(t *testing.T) {
	slice := Slice[int]{1, 2, 3, 4, 5}

	// Test normal range
	slice.CutInPlace(1, 3)
	expected := Slice[int]{1, 4, 5}
	if !slice.Eq(expected) {
		t.Errorf("CutInPlace(1, 3) failed. Expected %v, but got %v", expected, slice)
	}

	// Test range with negative indices
	slice = Slice[int]{1, 2, 3, 4, 5} // Restore the original slice
	slice.CutInPlace(-3, -2)
	expected = Slice[int]{1, 2, 4, 5}
	if !slice.Eq(expected) {
		t.Errorf("CutInPlace(-3, -2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test empty range
	slice = Slice[int]{1, 2, 3, 4, 5} // Restore the original slice
	slice.CutInPlace(0, 5)
	expected = Slice[int]{}
	if !slice.Eq(expected) {
		t.Errorf("CutInPlace(0, 5) failed. Expected %v, but got %v", expected, slice)
	}
}

func TestSliceLast(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	if sl.Last() != 5 {
		t.Error("Last() failed")
	}
}

func TestSliceLen(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	if sl.Len() != 5 {
		t.Errorf("Expected 5, got %d", sl.Len())
	}
}

func TestSlicePop(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	last, sl := sl.Pop()

	if last != 5 {
		t.Errorf("Expected 5, got %v", last)
	}

	if sl.Len() != 4 {
		t.Errorf("Expected 4, got %v", sl.Len())
	}

	r := Slice[int]{1, 2, 3, 4}
	if sl.Ne(r) {
		t.Errorf("Expected %v, got %v", r, sl)
	}
}

func TestSliceMaxInt(t *testing.T) {
	sl := Slice[Int]{1, 2, 3, 4, 5}
	if max := cmp.Max(sl...); max != 5 {
		t.Errorf("Max() = %d, want: %d.", max, 5)
	}
}

func TestSliceMaxFloats(t *testing.T) {
	sl := Slice[Float]{2.2, 2.8, 2.1, 2.7}
	if max := cmp.Max(sl...); max != 2.8 {
		t.Errorf("Max() = %f, want: %f.", max, 2.8)
	}
}

func TestSliceMinFloat(t *testing.T) {
	sl := Slice[Float]{2.2, 2.8, 2.1, 2.7}
	if min := cmp.Min(sl...); min != 2.1 {
		t.Errorf("Min() = %f; want: %f", min, 2.1)
	}
}

func TestSliceMinInt(t *testing.T) {
	sl := Slice[Int]{1, 2, 3, 4, 5}
	if min := cmp.Min(sl...); min != 1 {
		t.Errorf("Min() = %d; want: %d", min, 1)
	}
}

func TestSliceDelete(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	sl = sl.Delete(2)

	if !reflect.DeepEqual(sl, Slice[int]{1, 2, 4, 5}) {
		t.Errorf("Delete(2) = %v, want %v", sl, Slice[int]{1, 2, 4, 5})
	}

	sl = Slice[int]{1, 2, 3, 4, 5}
	sl = sl.Delete(-2)

	if !reflect.DeepEqual(sl, Slice[int]{1, 2, 3, 5}) {
		t.Errorf("Delete(2) = %v, want %v", sl, Slice[int]{1, 2, 3, 5})
	}
}

func TestSliceSFill(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5}
	sl.Fill(0)

	for _, v := range sl {
		if v != 0 {
			t.Errorf("Expected all elements to be 0, but found %d", v)
		}
	}
}

func TestSliceSet(t *testing.T) {
	sl := NewSlice[int](5)

	sl.Set(0, 1)
	sl.Set(0, 1)
	sl.Set(2, 2)
	sl.Set(4, 3)

	if !reflect.DeepEqual(sl, Slice[int]{1, 0, 2, 0, 3}) {
		t.Errorf("Set() = %v, want %v", sl, Slice[int]{1, 0, 2, 0, 3})
	}
}

func TestSliceReplace(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice[string]
		i, j     Int
		values   []string
		expected Slice[string]
	}{
		{
			name:     "basic test",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e", "f"},
			expected: Slice[string]{"a", "e", "f", "d"},
		},
		{
			name:     "replace at start",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        2,
			values:   []string{"e", "f"},
			expected: Slice[string]{"e", "f", "c", "d"},
		},
		{
			name:     "replace at end",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        2,
			j:        4,
			values:   []string{"e", "f"},
			expected: Slice[string]{"a", "b", "e", "f"},
		},
		{
			name:     "replace with more values",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        2,
			values:   []string{"e", "f", "g", "h"},
			expected: Slice[string]{"a", "e", "f", "g", "h", "c", "d"},
		},
		{
			name:     "replace with fewer values",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e"},
			expected: Slice[string]{"a", "e", "d"},
		},
		{
			name:     "replace entire slice",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        4,
			values:   []string{"e", "f", "g", "h"},
			expected: Slice[string]{"e", "f", "g", "h"},
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
	slice := Slice[string]{"a", "b", "c", "d"}
	newSlice := slice.Replace(1, 1, "zz", "xx")
	expected := Slice[string]{"a", "zz", "xx", "b", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(1, 1) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement with same start and end indices
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(2, 2, "zz", "xx")
	expected = Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(2, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from i to the end
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(2, slice.Len(), "zz", "xx")
	expected = Slice[string]{"a", "b", "zz", "xx"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(2, end) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from the start to j
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(0, 2, "zz", "xx")
	expected = Slice[string]{"zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(start, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test empty replacement
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(2, 2) // No replacement, should remain unchanged
	expected = Slice[string]{"a", "b", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(2, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from negative index to positive index
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(-2, 2, "zz", "xx")
	expected = Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(-2, 2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from positive index to negative index
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(1, -1, "zz", "xx")
	expected = Slice[string]{"a", "zz", "xx", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(1, -1) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from negative index to negative index
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(-3, -2, "zz", "xx")
	expected = Slice[string]{"a", "zz", "xx", "c", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(-3, -2) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement from negative index to positive index including negative values
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(-3, 3, "zz", "xx")
	expected = Slice[string]{"a", "zz", "xx", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(-3, 3) failed. Expected %v, but got %v", expected, newSlice)
	}

	// Test replacement with empty slice
	slice = Slice[string]{"a", "b", "c", "d"}
	newSlice = slice.Replace(1, 3)
	expected = Slice[string]{"a", "d"}
	if !newSlice.Eq(expected) {
		t.Errorf("Replace(1, 3) with empty slice failed. Expected %v, but got %v", expected, newSlice)
	}
}

func TestSliceReplaceInPlace(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice[string]
		i, j     Int
		values   []string
		expected Slice[string]
	}{
		{
			name:     "basic test",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e", "f"},
			expected: Slice[string]{"a", "e", "f", "d"},
		},
		{
			name:     "replace at start",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        2,
			values:   []string{"e", "f"},
			expected: Slice[string]{"e", "f", "c", "d"},
		},
		{
			name:     "replace at end",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        2,
			j:        4,
			values:   []string{"e", "f"},
			expected: Slice[string]{"a", "b", "e", "f"},
		},
		{
			name:     "replace with more values",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        2,
			values:   []string{"e", "f", "g", "h"},
			expected: Slice[string]{"a", "e", "f", "g", "h", "c", "d"},
		},
		{
			name:     "replace with fewer values",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e"},
			expected: Slice[string]{"a", "e", "d"},
		},
		{
			name:     "replace entire slice",
			input:    Slice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        4,
			values:   []string{"e", "f", "g", "h"},
			expected: Slice[string]{"e", "f", "g", "h"},
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
	slice := Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(1, 1, "zz", "xx")
	expected := Slice[string]{"a", "zz", "xx", "b", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(1, 1) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement with same start and end indices
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(2, 2, "zz", "xx")
	expected = Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(2, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from i to the end
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(2, slice.Len(), "zz", "xx")
	expected = Slice[string]{"a", "b", "zz", "xx"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(2, end) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from the start to j
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(0, 2, "zz", "xx")
	expected = Slice[string]{"zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(start, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test empty replacement
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(2, 2) // No replacement, should remain unchanged
	expected = Slice[string]{"a", "b", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(2, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from negative index to positive index
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(-2, 2, "zz", "xx")
	expected = Slice[string]{"a", "b", "zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(-2, 2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from positive index to negative index
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(1, -1, "zz", "xx")
	expected = Slice[string]{"a", "zz", "xx", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(1, -1) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from negative index to negative index
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(-3, -2, "zz", "xx")
	expected = Slice[string]{"a", "zz", "xx", "c", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(-3, -2) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement from negative index to positive index including negative values
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(-3, 3, "zz", "xx")
	expected = Slice[string]{"a", "zz", "xx", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(-3, 3) failed. Expected %v, but got %v", expected, slice)
	}

	// Test replacement with empty slice
	slice = Slice[string]{"a", "b", "c", "d"}
	slice.ReplaceInPlace(1, 3)
	expected = Slice[string]{"a", "d"}
	if !slice.Eq(expected) {
		t.Errorf("ReplaceInPlace(1, 3) with empty slice failed. Expected %v, but got %v", expected, slice)
	}
}

func TestSliceContainsAny(t *testing.T) {
	testCases := []struct {
		sl     Slice[int]
		other  Slice[int]
		expect bool
	}{
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{6, 7, 8, 9, 10}, false},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{5, 6, 7, 8, 9}, true},
		{Slice[int]{}, Slice[int]{1, 2, 3, 4, 5}, false},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{}, false},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{1, 2, 3, 4, 5}, true},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{6}, false},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{1}, true},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{6, 7, 8, 9, 0, 3}, true},
	}

	for _, tc := range testCases {
		if result := tc.sl.ContainsAny(tc.other...); result != tc.expect {
			t.Errorf("ContainsAny(%v, %v) = %v; want %v", tc.sl, tc.other, result, tc.expect)
		}
	}
}

func TestSliceContainsAll(t *testing.T) {
	testCases := []struct {
		sl     Slice[int]
		other  Slice[int]
		expect bool
	}{
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{1, 2, 3}, true},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{1, 2, 3, 6}, false},
		{Slice[int]{}, Slice[int]{1, 2, 3, 4, 5}, false},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{}, false},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{1, 2, 3, 4, 5}, true},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{6}, false},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{1}, true},
		{Slice[int]{1, 2, 3, 4, 5}, Slice[int]{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 5, 5}, true},
	}

	for _, tc := range testCases {
		if result := tc.sl.ContainsAll(tc.other...); result != tc.expect {
			t.Errorf("ContainsAll(%v, %v) = %v; want %v", tc.sl, tc.other, result, tc.expect)
		}
	}
}

func TestSliceSubSlice(t *testing.T) {
	// Test with an empty slice
	emptySlice := Slice[int]{}
	emptySubSlice := emptySlice.SubSlice(0, 0)
	if !emptySubSlice.Empty() {
		t.Errorf("Expected empty slice for empty source slice, but got: %v", emptySubSlice)
	}

	// Test with a non-empty slice
	slice := Slice[int]{1, 2, 3, 4, 5}

	// Test a valid range within bounds
	subSlice := slice.SubSlice(1, 4)
	expected := Slice[int]{2, 3, 4}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with a single negative index
	subSlice = slice.SubSlice(-2, slice.Len())
	expected = Slice[int]{4, 5}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with a negative start and end index
	subSlice = slice.SubSlice(-3, -1)
	expected = Slice[int]{3, 4}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}
}

func TestSubSliceOutOfBoundsStartIndex(t *testing.T) {
	slice := Slice[int]{1, 2, 3, 4, 5}

	// Test with start index beyond slice length (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for start index beyond slice length, but no panic occurred")
		}
	}()
	_ = slice.SubSlice(10, slice.Len())
}

func TestSubSliceOutOfBoundsNegativeStartIndex(t *testing.T) {
	slice := Slice[int]{1, 2, 3, 4, 5}

	// Test with a negative start index beyond slice length (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for negative start index beyond slice length, but no panic occurred")
		}
	}()
	_ = slice.SubSlice(-10, slice.Len())
}

func TestSubSliceOutOfBoundsEndIndex(t *testing.T) {
	slice := Slice[int]{1, 2, 3, 4, 5}

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
	initialSlice := SliceOf(1, 2, 3)

	// Check the initial capacity of the slice.
	var initialCapacity Int = initialSlice.Cap()

	// Grow the slice to accommodate more elements.
	var newCapacity Int = initialCapacity + 5
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
	sl1 := SliceOf(1, 2, 3)
	if !sl1.NotEmpty() {
		t.Errorf("Test case 1 failed: Expected slice to be not empty")
	}

	// Test case 2: Empty slice
	sl2 := NewSlice[Int]()
	if sl2.NotEmpty() {
		t.Errorf("Test case 2 failed: Expected slice to be empty")
	}
}

func TestSliceAppendInPlace(t *testing.T) {
	// Create a slice with initial elements
	initialSlice := Slice[int]{1, 2, 3}

	// Append additional elements using AppendInPlace
	initialSlice.AppendInPlace(4, 5, 6)

	// Verify that the slice has the expected elements
	expected := Slice[int]{1, 2, 3, 4, 5, 6}
	if !initialSlice.Eq(expected) {
		t.Errorf("AppendInPlace failed. Expected: %v, Got: %v", expected, initialSlice)
	}
}

func TestSliceString(t *testing.T) {
	// Create a slice with some elements
	sl := SliceOf(1, 2, 3, 4, 5)

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
	slInt1 := Slice[Int]{1, 2, 3}
	slInt2 := Slice[Int]{1, 2, 3}
	if !slInt1.Eq(slInt2) {
		t.Errorf("Test 1: Expected slices to be equal")
	}

	// Test case: Function returns false for unequal slices of known types (String)
	slString1 := Slice[String]{"a", "b", "c"}
	slString2 := Slice[String]{"a", "x", "c"}
	if slString1.Eq(slString2) {
		t.Errorf("Test 2: Expected slices to be unequal")
	}

	// Test case: Function returns true for empty slices
	emptySlice1 := Slice[Float]{}
	emptySlice2 := Slice[Float]{}
	if !emptySlice1.Eq(emptySlice2) {
		t.Errorf("Test 3: Expected empty slices to be equal")
	}

	// Test case: Function returns false for slices of different lengths
	slFloat1 := Slice[Float]{1.1, 2.2, 3.3}
	slFloat2 := Slice[Float]{1.1, 2.2}
	if slFloat1.Eq(slFloat2) {
		t.Errorf("Test 4: Expected slices of different lengths to be unequal")
	}

	// Test case: Function returns true for equal slices of string type
	slString3 := Slice[string]{"apple", "banana", "cherry"}
	slString4 := Slice[string]{"apple", "banana", "cherry"}
	if !slString3.Eq(slString4) {
		t.Errorf("Test 5: Expected slices to be equal")
	}

	// Test case: Function returns false for unequal slices of int type
	slInt3 := Slice[int]{10, 20, 30}
	slInt4 := Slice[int]{10, 20, 40}
	if slInt3.Eq(slInt4) {
		t.Errorf("Test 6: Expected slices to be unequal")
	}

	// Test case: Function returns true for equal slices of float64 type
	slFloat64_1 := Slice[float64]{1.1, 2.2, 3.3}
	slFloat64_2 := Slice[float64]{1.1, 2.2, 3.3}
	if !slFloat64_1.Eq(slFloat64_2) {
		t.Errorf("Test 7: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of bool type
	slBool1 := Slice[bool]{true, false, true}
	slBool2 := Slice[bool]{true, false, true}
	if !slBool1.Eq(slBool2) {
		t.Errorf("Test 8: Expected slices to be equal")
	}

	// Test case: Function returns false for unequal slices of byte type
	slByte1 := Slice[byte]{1, 2, 3}
	slByte2 := Slice[byte]{1, 2, 4}
	if slByte1.Eq(slByte2) {
		t.Errorf("Test 9: Expected slices to be unequal")
	}

	// Test case: Function returns true for equal slices of int8 type
	slInt81 := Slice[int8]{1, 2, 3}
	slInt82 := Slice[int8]{1, 2, 3}
	if !slInt81.Eq(slInt82) {
		t.Errorf("Test 10: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of int16 type
	slInt161 := Slice[int16]{1, 2, 3}
	slInt162 := Slice[int16]{1, 2, 3}
	if !slInt161.Eq(slInt162) {
		t.Errorf("Test 11: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of int32 type
	slInt321 := Slice[int32]{1, 2, 3}
	slInt322 := Slice[int32]{1, 2, 3}
	if !slInt321.Eq(slInt322) {
		t.Errorf("Test 12: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of int64 type
	slInt641 := Slice[int64]{1, 2, 3}
	slInt642 := Slice[int64]{1, 2, 3}
	if !slInt641.Eq(slInt642) {
		t.Errorf("Test 13: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint type
	slUint1 := Slice[uint]{1, 2, 3}
	slUint2 := Slice[uint]{1, 2, 3}
	if !slUint1.Eq(slUint2) {
		t.Errorf("Test 14: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint8 type
	slUint81 := Slice[uint8]{1, 2, 3}
	slUint82 := Slice[uint8]{1, 2, 3}
	if !slUint81.Eq(slUint82) {
		t.Errorf("Test 15: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint16 type
	slUint161 := Slice[uint16]{1, 2, 3}
	slUint162 := Slice[uint16]{1, 2, 3}
	if !slUint161.Eq(slUint162) {
		t.Errorf("Test 16: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint32 type
	slUint321 := Slice[uint32]{1, 2, 3}
	slUint322 := Slice[uint32]{1, 2, 3}
	if !slUint321.Eq(slUint322) {
		t.Errorf("Test 17: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of uint64 type
	slUint641 := Slice[uint64]{1, 2, 3}
	slUint642 := Slice[uint64]{1, 2, 3}
	if !slUint641.Eq(slUint642) {
		t.Errorf("Test 18: Expected slices to be equal")
	}

	// Test case: Function returns true for equal slices of float32 type
	slFloat321 := Slice[float32]{1.1, 2.2, 3.3}
	slFloat322 := Slice[float32]{1.1, 2.2, 3.3}
	if !slFloat321.Eq(slFloat322) {
		t.Errorf("Test 19: Expected slices to be equal")
	}
}

func TestSliceLastIndex(t *testing.T) {
	// Create a slice with some elements
	sl := SliceOf(1, 2, 3, 4, 5)

	// Get the last index of the slice
	lastIndex := sl.LastIndex()

	// Check if the last index is correct
	expectedLastIndex := sl.Len() - 1
	if lastIndex != expectedLastIndex {
		t.Errorf("Slice LastIndex method failed. Expected: %d, Got: %d", expectedLastIndex, lastIndex)
	}

	// Create an empty slice
	emptySlice := NewSlice[int]()

	// Get the last index of the empty slice
	emptyLastIndex := emptySlice.LastIndex()

	// Check if the last index of an empty slice is 0
	if emptyLastIndex != 0 {
		t.Errorf("Slice LastIndex method failed for empty slice. Expected: 0, Got: %d", emptyLastIndex)
	}
}

func TestSliceRandom(t *testing.T) {
	// Create a slice with some elements
	sl := SliceOf(1, 2, 3, 4, 5)

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
	emptySlice := NewSlice[int]()

	// Get a random element from the empty slice
	emptyRandomElement := emptySlice.Random()

	// Check if the random element from an empty slice is zero value
	if emptyRandomElement != 0 {
		t.Errorf("Slice Random method failed for empty slice. Expected: 0, Got: %d", emptyRandomElement)
	}
}

func TestSliceAddUniqueInPlace(t *testing.T) {
	// Test cases for Int
	testAddUniqueInPlace(t, SliceOf(1, 2, 3, 4, 5), []int{3, 4, 5, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7})

	// Test cases for Float
	testAddUniqueInPlace(
		t,
		SliceOf(1.1, 2.2, 3.3, 4.4, 5.5),
		[]float64{3.3, 4.4, 5.5, 6.6, 7.7},
		[]float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7},
	)

	// Test cases for String
	testAddUniqueInPlace(
		t,
		SliceOf("apple", "banana", "orange", "grape"),
		[]string{"orange", "grape", "kiwi"},
		[]string{"apple", "banana", "orange", "grape", "kiwi"},
	)

	// Add more test cases for other types as needed
}

func testAddUniqueInPlace[T comparable](t *testing.T, sl Slice[T], elems, expected []T) {
	sl.AddUniqueInPlace(elems...)
	if !sl.Eq(SliceOf(expected...)) {
		t.Errorf("Slice AddUniqueInPlace method failed for type %T. Expected: %v, Got: %v", sl[0], expected, sl)
	}
}

func TestSliceAsAny(t *testing.T) {
	// Test cases for Int
	testSliceAsAny(t, SliceOf(1, 2, 3), []any{1, 2, 3})

	// Test cases for Float
	testSliceAsAny(t, SliceOf(1.1, 2.2, 3.3), []any{1.1, 2.2, 3.3})

	// Test cases for String
	testSliceAsAny(t, SliceOf("apple", "banana", "orange"), []any{"apple", "banana", "orange"})

	// Add more test cases for other types as needed
}

func testSliceAsAny[T any](t *testing.T, sl Slice[T], expected []any) {
	result := sl.AsAny()
	if !result.Eq(SliceOf(expected...)) {
		t.Errorf("Slice AsAny method failed for type %T. Expected: %v, Got: %v", sl[0], expected, result)
	}
}

func TestSliceContainsBy(t *testing.T) {
	// Test case 1: Slice contains the element that satisfies the provided function
	sl1 := Slice[Int]{1, 2, 3, 4, 5}
	contains1 := sl1.ContainsBy(f.Eq(Int(3)))

	if !contains1 {
		t.Errorf("Test 1: Expected true, got false")
	}

	// Test case 2: Slice does not contain the element that satisfies the provided function
	sl2 := Slice[String]{"apple", "banana", "cherry"}
	contains2 := sl2.ContainsBy(f.Eq(String("orange")))

	if contains2 {
		t.Errorf("Test 2: Expected false, got true")
	}

	// Test case 3: Slice contains the element that satisfies the provided function (using custom struct)
	type Person struct {
		Name string
		Age  int
	}

	sl3 := Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}, {Name: "Charlie", Age: 35}}

	contains3 := sl3.ContainsBy(func(x Person) bool { return x.Name == "Bob" && x.Age == 25 })
	if !contains3 {
		t.Errorf("Test 3: Expected true, got false")
	}
}

func TestSliceEqNeBy(t *testing.T) {
	// Test case 1: Slices are equal using the equality function
	sl1 := Slice[Int]{1, 2, 3}
	sl2 := Slice[Int]{1, 2, 3}

	eq1 := sl1.EqBy(sl2, func(x, y Int) bool { return x.Eq(y) })

	if !eq1 {
		t.Errorf("Test 1: Expected true, got false")
	}

	// Test case 2: Slices are not equal using the equality function
	sl3 := Slice[String]{"apple", "banana", "cherry"}
	sl4 := Slice[String]{"apple", "orange", "cherry"}

	eq2 := sl3.EqBy(sl4, func(x, y String) bool { return x.Eq(y) })

	if eq2 {
		t.Errorf("Test 2: Expected false, got true")
	}

	// Test case 3: Slices are equal using the equality function (using custom struct)
	type Person struct {
		Name string
		Age  int
	}

	sl5 := Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	sl6 := Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}

	eq3 := sl5.EqBy(sl6, func(x, y Person) bool {
		return x.Name == y.Name && x.Age == y.Age
	})

	if !eq3 {
		t.Errorf("Test 3: Expected true, got false")
	}

	// Additional tests for NeBy

	// Test case 4: Slices are not equal using the equality function
	sl7 := Slice[Int]{1, 2, 3}
	sl8 := Slice[Int]{1, 2, 4}

	ne1 := sl7.NeBy(sl8, func(x, y Int) bool { return x.Eq(y) })

	if !ne1 {
		t.Errorf("Test 4: Expected true, got false")
	}

	// Test case 5: Slices are equal using the equality function
	sl9 := Slice[String]{"apple", "banana", "cherry"}
	sl10 := Slice[String]{"apple", "banana", "cherry"}

	ne2 := sl9.NeBy(sl10, func(x, y String) bool { return x.Eq(y) })

	if ne2 {
		t.Errorf("Test 5: Expected false, got true")
	}
}

func TestSliceIndexBy(t *testing.T) {
	// Test case 1: Element satisfying the custom comparison function exists in the slice
	sl1 := Slice[int]{1, 2, 3, 4, 5}
	index1 := sl1.IndexBy(f.Eq(3))

	expectedIndex1 := Int(2)
	if index1 != expectedIndex1 {
		t.Errorf("Test 1: Expected index %d, got %d", expectedIndex1, index1)
	}

	// Test case 2: Element satisfying the custom comparison function doesn't exist in the slice
	sl2 := Slice[string]{"apple", "banana", "cherry"}
	index2 := sl2.IndexBy(f.Eq("orange"))

	expectedIndex2 := Int(-1)
	if index2 != expectedIndex2 {
		t.Errorf("Test 2: Expected index %d, got %d", expectedIndex2, index2)
	}

	// Test case 3: Element satisfying the custom comparison function exists in the slice (using custom struct)
	type Person struct {
		Name string
		Age  int
	}

	sl3 := Slice[Person]{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	index3 := sl3.IndexBy(func(x Person) bool { return x.Name == "Bob" && x.Age == 25 })

	expectedIndex3 := Int(1)
	if index3 != expectedIndex3 {
		t.Errorf("Test 3: Expected index %d, got %d", expectedIndex3, index3)
	}
}

func TestSliceSwap(t *testing.T) {
	// Define a slice to test
	s := Slice[int]{1, 2, 3, 4, 5}

	// Test swapping two elements
	s.Swap(1, 3)
	expected := Slice[int]{1, 4, 3, 2, 5}
	if !s.Eq(expected) {
		t.Errorf("Swap failed: got %v, want %v", s, expected)
	}

	// Test swapping two elements
	s.Swap(-1, 0)
	expected = Slice[int]{5, 4, 3, 2, 1}
	if !s.Eq(expected) {
		t.Errorf("Swap failed: got %v, want %v", s, expected)
	}
}

func TestSliceMaxBy(t *testing.T) {
	// Test case 1: Maximum integer
	s := Slice[int]{3, 1, 4, 2, 5}
	maxInt := s.MaxBy(cmp.Cmp)
	expectedMaxInt := 5
	if maxInt != expectedMaxInt {
		t.Errorf("s.MaxBy(IntCompare) = %d; want %d", maxInt, expectedMaxInt)
	}
}

func TestSliceMinBy(t *testing.T) {
	// Test case 1: Minimum integer
	s := Slice[int]{3, 1, 4, 2, 5}
	minInt := s.MinBy(cmp.Cmp)
	expectedMinInt := 1
	if minInt != expectedMinInt {
		t.Errorf("s.MinBy(IntCompare) = %d; want %d", minInt, expectedMinInt)
	}
}

func TestSliceTransform(t *testing.T) {
	original := Slice[int]{1, 2, 3}

	addElement := func(sl Slice[int]) Slice[int] {
		return append(sl, 4)
	}

	expected := Slice[int]{1, 2, 3, 4}
	result := original.Transform(addElement)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Transform failed: expected %v, got %v", expected, result)
	}

	removeLast := func(sl Slice[int]) Slice[int] {
		return sl[:len(sl)-1]
	}

	expectedAfterRemoval := Slice[int]{1, 2, 3}
	resultAfterRemoval := result.Transform(removeLast)

	if !reflect.DeepEqual(resultAfterRemoval, expectedAfterRemoval) {
		t.Errorf("Transform with removal failed: expected %v, got %v", expectedAfterRemoval, resultAfterRemoval)
	}
}
