package g_test

import (
	"reflect"
	"strings"
	"testing"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/iter"
)

func TestPermutations(t *testing.T) {
	empty := g.NewSlice[int]()
	permsEmpty := empty.Permutations()
	expectedPermsEmpty := []g.Slice[int]{empty}

	if !reflect.DeepEqual(permsEmpty, expectedPermsEmpty) {
		t.Errorf("expected %v, but got %v", expectedPermsEmpty, permsEmpty)
	}

	slice1 := g.SliceOf(1)
	perms1 := slice1.Permutations()
	expectedPerms1 := []g.Slice[int]{slice1}

	if !reflect.DeepEqual(perms1, expectedPerms1) {
		t.Errorf("expected %v, but got %v", expectedPerms1, perms1)
	}

	slice2 := g.SliceOf("a", "b")
	perms2 := slice2.Permutations()
	expectedPerms2 := []g.Slice[string]{
		{"a", "b"},
		{"b", "a"},
	}

	if !reflect.DeepEqual(perms2, expectedPerms2) {
		t.Errorf("expected %v, but got %v", expectedPerms2, perms2)
	}

	slice3 := g.SliceOf(1.0, 2.0, 3.0)
	perms3 := slice3.Permutations()

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
}

func TestSliceInsert(t *testing.T) {
	empty := g.NewSlice[int]()
	empty = empty.Insert(0, 1, 2)

	expectedEmpty := g.NewSlice[int]().Append(1, 2)
	if !reflect.DeepEqual(empty, expectedEmpty) {
		t.Errorf("expected %v, but got %v", expectedEmpty, empty)
	}

	slice1 := g.NewSlice[int]().Append(3, 4)
	slice1 = slice1.Insert(0, 1, 2)

	expected1 := g.NewSlice[int]().Append(1, 2, 3, 4)
	if !reflect.DeepEqual(slice1, expected1) {
		t.Errorf("expected %v, but got %v", expected1, slice1)
	}

	slice2 := g.NewSlice[string]().Append("foo", "bar", "baz")
	slice2 = slice2.Insert(1, "qux", "quux")

	expected2 := g.NewSlice[string]().Append("foo", "qux", "quux", "bar", "baz")
	if !reflect.DeepEqual(slice2, expected2) {
		t.Errorf("expected %v, but got %v", expected2, slice2)
	}

	slice3 := g.NewSlice[float64]().Append(1.23, 4.56)
	slice3 = slice3.Insert(slice3.Len(), 7.89)

	expected3 := g.NewSlice[float64]().Append(1.23, 4.56, 7.89)
	if !reflect.DeepEqual(slice3, expected3) {
		t.Errorf("expected %v, but got %v", expected3, slice3)
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

func TestSliceHMapHashedHnt(t *testing.T) {
	sl := g.Slice[g.Int]{1, 2, 3, 4, 5}
	m := sl.ToMapHashed()

	if m.Len() != sl.Len() {
		t.Errorf("Expected %d, got %d", sl.Len(), m.Len())
	}

	for _, v := range sl {
		if !m.Contains(v.Hash().MD5()) {
			t.Errorf("Expected %v, got %v", v, m[v.Hash().MD5()])
		}
	}
}

func TestSliceHMapHashedStrings(t *testing.T) {
	sl := g.Slice[g.String]{"1", "2", "3", "4", "5"}
	m := sl.ToMapHashed()

	if m.Len() != sl.Len() {
		t.Errorf("Expected %d, got %d", sl.Len(), m.Len())
	}

	for _, v := range sl {
		if !m.Contains(v.Hash().MD5()) {
			t.Errorf("Expected %v, got %v", v, m[v.Hash().MD5()])
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

func TestSliceChunks(t *testing.T) {
	tests := []struct {
		name     string
		input    g.Slice[int]
		expected []g.Slice[int]
		size     int
	}{
		{
			name:     "empty slice",
			input:    g.NewSlice[int](),
			expected: []g.Slice[int]{},
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
			result := tt.input.Chunks(tt.size)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d chunks, but got %d", len(tt.expected), len(result))
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

func TestSliceReverse(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	sl = sl.Reverse()

	if !reflect.DeepEqual(sl, g.Slice[int]{5, 4, 3, 2, 1}) {
		t.Errorf("Expected %v, got %v", g.Slice[int]{5, 4, 3, 2, 1}, sl)
	}
}

func TestSliceAll(t *testing.T) {
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
			got := tc.sl.All(tc.f)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSliceAny(t *testing.T) {
	sl1 := g.NewSlice[int]()
	f1 := func(x int) bool { return x > 0 }

	if sl1.Any(f1) {
		t.Errorf("Expected false for empty slice, got true")
	}

	sl2 := g.NewSlice[int]().Append(1, 2, 3)
	f2 := func(x int) bool { return x < 1 }

	if sl2.Any(f2) {
		t.Errorf("Expected false for slice with no matching elements, got true")
	}

	sl3 := g.NewSlice[string]().Append("foo", "bar")
	f3 := func(x string) bool { return x == "bar" }

	if !sl3.Any(f3) {
		t.Errorf("Expected true for slice with one matching element, got false")
	}

	sl4 := g.NewSlice[int]().Append(1, 2, 3, 4, 5)
	f4 := func(x int) bool { return x%2 == 0 }

	if !sl4.Any(f4) {
		t.Errorf("Expected true for slice with multiple matching elements, got false")
	}
}

func TestSliceReduce(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	sum := sl.Reduce(func(index, value int) int { return index + value }, 0)

	if sum != 15 {
		t.Errorf("Expected %d, got %d", 15, sum)
	}
}

func TestSliceFilter(t *testing.T) {
	var sl g.Slice[int]

	sl = sl.Append(1, 2, 3, 4, 5)
	result := sl.Filter(func(v int) bool { return v%2 == 0 })

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

func TestSliceMap(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	result := sl.Map(func(i int) int { return i * 2 })

	if result.Len() != sl.Len() {
		t.Errorf("Expected %d, got %d", sl.Len(), result.Len())
	}

	for i := 0; i < result.Len(); i++ {
		if result[i] != sl[i]*2 {
			t.Errorf("Expected %d, got %d", sl[i]*2, result[i])
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
}

func TestSliceSortBy(t *testing.T) {
	sl1 := g.NewSlice[int]().Append(3, 1, 4, 1, 5)
	expected1 := g.NewSlice[int]().Append(1, 1, 3, 4, 5)

	actual1 := sl1.SortBy(func(i, j int) bool { return sl1[i] < sl1[j] })

	if !actual1.Eq(expected1) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected1, actual1)
	}

	sl2 := g.NewSlice[string]().Append("foo", "bar", "baz")
	expected2 := g.NewSlice[string]().Append("foo", "baz", "bar")

	actual2 := sl2.SortBy(func(i, j int) bool { return sl2[i] > sl2[j] })

	if !actual2.Eq(expected2) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected2, actual2)
	}

	sl3 := g.NewSlice[int]()
	expected3 := g.NewSlice[int]()

	actual3 := sl3.SortBy(func(i, j int) bool { return sl3[i] < sl3[j] })

	if !actual3.Eq(expected3) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected3, actual3)
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

func TestSliceLast(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	if sl.Last() != 5 {
		t.Error("Last() failed")
	}
}

func TestSliceLastIndex(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5}
	if sl.LastIndex() != 4 {
		t.Error("LastIndex() failed")
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
}

func TestSliceRandom(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < 10; i++ {
		if sl.Random() < 1 || sl.Random() > 10 {
			t.Error("Random() should return a number between 1 and 10")
		}
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

func TestSliceFilterZeroValues(t *testing.T) {
	sl := g.Slice[int]{1, 2, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10}
	sl = sl.FilterZeroValues()

	if sl.Len() != 10 {
		t.Errorf("Expected 10, got %d", sl.Len())
	}

	for i := range iter.N(sl.Len()) {
		if sl[i] == 0 {
			t.Errorf("Expected non-zero value, got %d", sl[i])
		}
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

func TestSliceForEach(t *testing.T) {
	sl1 := g.NewSlice[int]().Append(1, 2, 3, 4, 5)
	sl2 := g.NewSlice[string]().Append("foo", "bar", "baz")
	sl3 := g.NewSlice[float64]().Append(1.1, 2.2, 3.3, 4.4)

	var result1 []int

	sl1.ForEach(func(i int) { result1 = append(result1, i) })

	if !reflect.DeepEqual(result1, []int{1, 2, 3, 4, 5}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			sl1,
			[]int{1, 2, 3, 4, 5},
			result1,
		)
	}

	var result2 []string

	sl2.ForEach(func(s string) { result2 = append(result2, s) })

	if !reflect.DeepEqual(result2, []string{"foo", "bar", "baz"}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			sl2,
			[]string{"foo", "bar", "baz"},
			result2,
		)
	}

	var result3 []float64

	sl3.ForEach(func(f float64) { result3 = append(result3, f) })

	if !reflect.DeepEqual(result3, []float64{1.1, 2.2, 3.3, 4.4}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			sl3,
			[]float64{1.1, 2.2, 3.3, 4.4},
			result3,
		)
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

func TestSliceMapParallel(t *testing.T) {
	sl := g.NewSlice[int](10).Fill(1)
	result := sl.MapParallel(func(x int) int { return x * 2 })
	expected := g.NewSlice[int](10).Fill(2)

	if !result.Eq(expected) {
		t.Errorf("Unexpected result: got %v, expected %v", result, expected)
	}

	sl = g.NewSlice[int](10000).Fill(1)
	result = sl.MapParallel(func(x int) int { return x * 2 })
	expected = g.NewSlice[int](10000).Fill(2)

	if !result.Eq(expected) {
		t.Errorf("Unexpected result: got %v, expected %v", result, expected)
	}
}

func TestSliceFilterParallel(t *testing.T) {
	sl := g.SliceOf(1, 2, 3, 4, 5)
	expected := g.SliceOf(2, 4)
	actual := sl.FilterParallel(func(x int) bool { return x%2 == 0 })

	if !actual.Eq(expected) {
		t.Errorf("FilterParallel failed. Expected %v, but got %v", expected, actual)
	}

	sl = g.SliceOf(2, 4, 6, 8, 10)
	expected = sl.Clone()
	actual = sl.FilterParallel(func(x int) bool { return x%2 == 0 })

	if !actual.Eq(expected) {
		t.Errorf("FilterParallel failed. Expected %v, but got %v", expected, actual)
	}

	sl = g.SliceOf(1, 3, 5, 7, 9)
	expected = g.NewSlice[int]()
	actual = sl.FilterParallel(func(x int) bool { return x%2 == 0 })

	if !actual.Eq(expected) {
		t.Errorf("FilterParallel failed. Expected %v, but got %v", expected, actual)
	}
}

func TestSliceReduceParallel(t *testing.T) {
	sl := g.NewSlice[int](10).Fill(1)
	result := sl.ReduceParallel(func(a, b int) int { return a + b }, 0)
	expected := sl.Reduce(func(a, b int) int { return a + b }, 0)

	if result != expected {
		t.Errorf("Unexpected result: got %d, expected %d", result, expected)
	}

	sl = g.NewSlice[int](10000).Fill(1)
	result = sl.ReduceParallel(func(a, b int) int { return a + b }, 0)
	expected = sl.Reduce(func(a, b int) int { return a + b }, 0)

	if result != expected {
		t.Errorf("Unexpected result: got %d, expected %d", result, expected)
	}
}

func TestSliceZip(t *testing.T) {
	s1 := g.SliceOf(1, 2, 3, 4)
	s2 := g.SliceOf(5, 6, 7, 8)
	expected := []g.Slice[int]{{1, 5}, {2, 6}, {3, 7}, {4, 8}}
	result := s1.Zip(s2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s1, s2, result, expected)
	}

	s3 := g.SliceOf(1, 2, 3)
	s4 := g.SliceOf(4, 5)
	expected = []g.Slice[int]{{1, 4}, {2, 5}}
	result = s3.Zip(s4)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s3, s4, result, expected)
	}

	s5 := g.SliceOf(1, 2, 3)
	s6 := g.SliceOf(4, 5, 6)
	s7 := g.SliceOf(7, 8, 9)
	expected = []g.Slice[int]{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}}
	result = s5.Zip(s6, s7)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v, %v) = %v, expected %v", s5, s6, s7, result, expected)
	}

	s8 := g.SliceOf(1, 2, 3)
	s9 := g.SliceOf(4, 5)
	s10 := g.SliceOf(6)
	expected = []g.Slice[int]{{1, 4, 6}}
	result = s8.Zip(s9, s10)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v, %v) = %v, expected %v", s8, s9, s10, result, expected)
	}
}

func TestSliceFlatten(t *testing.T) {
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
				g.Slice[int]{2, 3},
				"abc",
				g.Slice[string]{"def", "ghi"},
				g.Slice[float64]{4.5, 6.7},
			},
			expected: g.Slice[any]{1, 2, 3, "abc", "def", "ghi", 4.5, 6.7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Flatten()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Flatten() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSliceCounter(t *testing.T) {
	sl1 := g.Slice[int]{1, 2, 3, 2, 1, 4, 5, 4, 4}
	sl2 := g.Slice[string]{"apple", "banana", "orange", "apple", "apple", "orange", "grape"}

	expected1 := g.NewMap[any, int]()
	expected1.Set(3, 1)
	expected1.Set(5, 1)
	expected1.Set(1, 2)
	expected1.Set(2, 2)
	expected1.Set(4, 3)

	result1 := sl1.Counter()
	if !result1.Eq(expected1) {
		t.Errorf("Counter() returned %v, expected %v", result1, expected1)
	}

	// Test with string values
	expected2 := g.NewMap[any, int]()
	expected2.Set("banana", 1)
	expected2.Set("grape", 1)
	expected2.Set("orange", 2)
	expected2.Set("apple", 3)

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
		if result := tc.sl.ContainsAny(tc.other); result != tc.expect {
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
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1, 2, 3, 4, 5}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{6}, false},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1}, true},
		{g.Slice[int]{1, 2, 3, 4, 5}, g.Slice[int]{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 5, 5}, true},
	}

	for _, tc := range testCases {
		if result := tc.sl.ContainsAll(tc.other); result != tc.expect {
			t.Errorf("ContainsAll(%v, %v) = %v; want %v", tc.sl, tc.other, result, tc.expect)
		}
	}
}

func TestSliceUnique(t *testing.T) {
	testCases := []struct {
		input  g.Slice[int]
		output g.Slice[int]
	}{
		{
			input:  g.NewSlice[int]().Append(1, 2, 3, 4, 5),
			output: g.NewSlice[int]().Append(1, 2, 3, 4, 5),
		},
		{
			input:  g.NewSlice[int]().Append(1, 2, 3, 4, 5, 5, 4, 3, 2, 1),
			output: g.NewSlice[int]().Append(1, 2, 3, 4, 5),
		},
		{
			input:  g.NewSlice[int]().Append(1, 1, 1, 1, 1),
			output: g.NewSlice[int]().Append(1),
		},
		{
			input:  g.NewSlice[int](),
			output: g.NewSlice[int](),
		},
	}

	for _, tc := range testCases {
		actual := tc.input.Unique()
		if !reflect.DeepEqual(actual, tc.output) {
			t.Errorf("Unique(%v) returned %v, expected %v", tc.input, actual, tc.output)
		}
	}
}

func TestSliceRange(t *testing.T) {
	// Test with an empty slice
	emptySlice := g.Slice[int]{}
	emptySubSlice := emptySlice.Range(0)
	if !emptySubSlice.Empty() {
		t.Errorf("Expected empty slice for empty source slice, but got: %v", emptySubSlice)
	}

	// Test with a non-empty slice
	slice := g.Slice[int]{1, 2, 3, 4, 5}

	// Test a valid range within bounds
	subSlice := slice.Range(1, 4)
	expected := g.Slice[int]{2, 3, 4}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with a single negative index
	subSlice = slice.Range(-2)
	expected = g.Slice[int]{4, 5}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with a negative start and end index
	subSlice = slice.Range(-3, -1)
	expected = g.Slice[int]{3, 4}
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with a negative start index beyond slice length
	subSlice = slice.Range(-6)
	if !subSlice.Eq(slice) {
		t.Errorf("Expected empty slice for negative start index beyond slice length, but got: %v", subSlice)
	}

	// Test with a negative end index beyond slice length
	subSlice = slice.Range(0, -6)
	if !subSlice.Empty() {
		t.Errorf("Expected empty slice for negative end index beyond slice length, but got: %v", subSlice)
	}

	// Test with a start index beyond slice length
	subSlice = slice.Range(6)
	if !subSlice.Empty() {
		t.Errorf("Expected empty slice for start index beyond slice length, but got: %v", subSlice)
	}

	// Test with an end index beyond slice length
	subSlice = slice.Range(0, 6)
	expected = slice
	if !subSlice.Eq(expected) {
		t.Errorf("Expected subSlice: %v, but got: %v", expected, subSlice)
	}

	// Test with start index greater than end index
	subSlice = slice.Range(3, 2)
	if !subSlice.Empty() {
		t.Errorf("Expected empty slice for start index greater than end index, but got: %v", subSlice)
	}
}
