package g_test

import (
	"context"
	"reflect"
	"slices"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestMapOrdIterReverse(t *testing.T) {
	tests := []struct {
		name     string
		actions  func(m MapOrd[int, int]) MapOrd[int, int]
		expected []Pair[int, int]
	}{
		{
			name: "empty map",
			actions: func(m MapOrd[int, int]) MapOrd[int, int] {
				return m
			},
			expected: nil,
		},
		{
			name: "single element",
			actions: func(m MapOrd[int, int]) MapOrd[int, int] {
				m.Set(1, 100)
				return m
			},
			expected: []Pair[int, int]{{1, 100}},
		},
		{
			name: "multiple elements",
			actions: func(m MapOrd[int, int]) MapOrd[int, int] {
				m.Set(1, 100)
				m.Set(2, 200)
				m.Set(3, 300)
				return m
			},
			expected: []Pair[int, int]{{3, 300}, {2, 200}, {1, 100}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapOrd := NewMapOrd[int, int]()
			mapOrd = tt.actions(mapOrd)
			iterator := mapOrd.IterReverse()
			var result []Pair[int, int]
			iterator.ForEach(func(k, v int) {
				result = append(result, Pair[int, int]{k, v})
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IterReverse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMapOrdIterSortBy(t *testing.T) {
	// Sample data
	data := NewMapOrd[int, string]()
	data.Set(1, "d")
	data.Set(3, "b")
	data.Set(2, "c")
	data.Set(5, "e")
	data.Set(4, "a")

	// Expected result
	expected := NewMapOrd[int, string]()
	expected.Set(1, "d")
	expected.Set(2, "c")
	expected.Set(3, "b")
	expected.Set(4, "a")
	expected.Set(5, "e")

	sortedItems := data.Iter().
		SortBy(func(a, b Pair[int, string]) cmp.Ordering { return cmp.Cmp(a.Key, b.Key) }).
		Collect()

	// Check if the result matches the expected output
	if !reflect.DeepEqual(sortedItems, expected) {
		t.Errorf("Expected %v, got %v", expected, sortedItems)
	}

	expected = NewMapOrd[int, string]()
	expected.Set(4, "a")
	expected.Set(3, "b")
	expected.Set(2, "c")
	expected.Set(1, "d")
	expected.Set(5, "e")

	sortedItems = data.Iter().
		SortBy(func(a, b Pair[int, string]) cmp.Ordering { return cmp.Cmp(a.Value, b.Value) }).
		Collect()

	// Check if the result matches the expected output
	if !reflect.DeepEqual(sortedItems, expected) {
		t.Errorf("Expected %v, got %v", expected, sortedItems)
	}
}

func TestMapOrdIterStepBy(t *testing.T) {
	// Test case 1: StepBy with a step size of 2
	mapData := NewMapOrd[string, int]()
	mapData.Set("one", 1)
	mapData.Set("two", 2)
	mapData.Set("three", 3)
	mapData.Set("four", 4)
	mapData.Set("five", 5)

	expectedResult := NewMapOrd[string, int]()
	expectedResult.Set("one", 1)
	expectedResult.Set("three", 3)
	expectedResult.Set("five", 5)

	iter := mapData.Iter().StepBy(2)
	result := iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 2: StepBy with a step size of 3
	mapData = NewMapOrd[string, int]()
	mapData.Set("one", 1)
	mapData.Set("two", 2)
	mapData.Set("three", 3)
	mapData.Set("four", 4)
	mapData.Set("five", 5)

	expectedResult = NewMapOrd[string, int]()
	expectedResult.Set("one", 1)
	expectedResult.Set("four", 4)

	iter = mapData.Iter().StepBy(3)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 3: StepBy with a step size larger than the map length

	mapData = NewMapOrd[string, int]()
	mapData.Set("one", 1)
	mapData.Set("two", 2)
	mapData.Set("three", 3)

	expectedResult = NewMapOrd[string, int]()
	expectedResult.Set("one", 1)

	iter = mapData.Iter().StepBy(5)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 4: StepBy with a step size of 1
	mapData = NewMapOrd[string, int]()
	mapData.Set("one", 1)
	mapData.Set("two", 2)
	mapData.Set("three", 3)

	expectedResult = NewMapOrd[string, int]()
	expectedResult.Set("one", 1)
	expectedResult.Set("two", 2)
	expectedResult.Set("three", 3)

	iter = mapData.Iter().StepBy(1)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}
}

func TestMapOrdIterRange(t *testing.T) {
	// Test scenario: Function stops at a specific key-value pair
	t.Run("FunctionStopsAtSpecificPair", func(t *testing.T) {
		orderedMap := MapOrd[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
		}
		expected := map[string]int{"a": 1, "b": 2}

		result := make(map[string]int)
		stopAtB := func(key string, val int) bool {
			result[key] = val
			return key != "b"
		}

		orderedMap.Iter().Range(stopAtB)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})

	// Test scenario: Function always returns true
	t.Run("FunctionAlwaysTrue", func(t *testing.T) {
		orderedMap := MapOrd[string, int]{
			{"a", 1},
			{"b", 2},
			{"c", 3},
		}

		expected := map[string]int{"a": 1, "b": 2, "c": 3}

		result := make(map[string]int)
		alwaysTrue := func(key string, val int) bool {
			result[key] = val
			return true
		}

		orderedMap.Iter().Range(alwaysTrue)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})

	// Test scenario: Empty ordered map
	t.Run("EmptyMap", func(t *testing.T) {
		emptyMap := MapOrd[string, int]{}
		expected := make(map[string]int)

		result := make(map[string]int)
		anyFunc := func(key string, val int) bool {
			result[key] = val
			return true
		}

		emptyMap.Iter().Range(anyFunc)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})
}

func TestMapOrdNe(t *testing.T) {
	// Test case 1: Maps are equal
	m1 := NewMapOrd[int, string]()
	m2 := NewMapOrd[int, string]()
	m1.Set(1, "a")
	m2.Set(1, "a")
	if m1.Ne(m2) {
		t.Errorf("Expected maps to be equal")
	}

	// Test case 2: Maps are not equal
	m2.Set(2, "b")
	if !m1.Ne(m2) {
		t.Errorf("Expected maps to be not equal")
	}
}

func TestMapOrdNotEmpty(t *testing.T) {
	// Test case 1: Map is empty
	m := NewMapOrd[int, string]()
	if m.NotEmpty() {
		t.Errorf("Expected map to be empty")
	}

	// Test case 2: Map is not empty
	m.Set(1, "a")
	if !m.NotEmpty() {
		t.Errorf("Expected map to be not empty")
	}
}

func TestMapOrdString(t *testing.T) {
	// Test case 1: Map with elements
	m := NewMapOrd[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")
	expected := "MapOrd{1:a, 2:b, 3:c}"
	if str := m.String(); str != expected {
		t.Errorf("Expected string representation to be %s, got %s", expected, str)
	}

	// Test case 2: Empty Map
	m2 := NewMapOrd[string, int]()
	expected2 := "MapOrd{}"
	if str := m2.String(); str != expected2 {
		t.Errorf("Expected string representation to be %s, got %s", expected2, str)
	}
}

func TestMapOrdClear(t *testing.T) {
	// Test case 1: Map with elements
	m := NewMapOrd[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Clear()
	if !m.Empty() {
		t.Errorf("Expected map to be empty after clearing")
	}

	// Test case 2: Empty Map
	m2 := NewMapOrd[string, int]()
	m2.Clear()
	if !m2.Empty() {
		t.Errorf("Expected empty map to remain empty after clearing")
	}
}

func TestMapOrdContains(t *testing.T) {
	// Test case 1: Map contains the key
	m := NewMapOrd[int, string]()
	m.Set(1, "a")
	if !m.Contains(1) {
		t.Errorf("Expected map to contain the key")
	}

	// Test case 2: Map doesn't contain the key
	if m.Contains(2) {
		t.Errorf("Expected map not to contain the key")
	}

	// Test case 3: Map contains the key
	m2 := NewMapOrd[[]int, []int]()
	m2.Set([]int{0}, []int{1})

	if !m2.Contains([]int{0}) {
		t.Errorf("Expected map to contain the key")
	}
}

func TestMapOrdValues(t *testing.T) {
	// Test case 1: Map with elements
	m := NewMapOrd[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")
	expected := Slice[string]{"a", "b", "c"}
	values := m.Values()
	if len(values) != len(expected) {
		t.Errorf("Expected values to have length %d, got %d", len(expected), len(values))
	}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected value at index %d to be %s, got %s", i, expected[i], v)
		}
	}

	// Test case 2: Empty Map
	m2 := NewMapOrd[string, int]()
	values2 := m2.Values()
	if len(values2) != 0 {
		t.Errorf("Expected values to be empty for an empty map")
	}
}

func TestMapOrdInvert(t *testing.T) {
	// Test case 1: Map with elements
	m := NewMapOrd[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")
	inverted := m.Invert()
	expected := NewMapOrd[string, int]()
	expected.Set("a", 1)
	expected.Set("b", 2)
	expected.Set("c", 3)
	if inverted.Len() != expected.Len() {
		t.Errorf("Expected inverted map to have length %d, got %d", expected.Len(), inverted.Len())
	}

	inverted.Iter().ForEach(func(k string, v int) {
		if !expected.Contains(k) {
			t.Errorf("Expected inverted map to contain key-value pair %s:%d", k, v)
		}
	})

	// Test case 2: Empty Map
	m2 := NewMapOrd[string, int]()
	inverted2 := m2.Invert()
	if inverted2.Len() != 0 {
		t.Errorf("Expected inverted map of an empty map to be empty")
	}
}

func TestMapOrdGetOrSet(t *testing.T) {
	// Test case 1: Key exists
	m := NewMapOrd[string, int]()
	m.Set("key1", 10)
	defaultValue := 20
	result := m.GetOrSet("key1", defaultValue)
	if result != 10 {
		t.Errorf("Expected value to be 10, got %d", result)
	}

	// Test case 2: Key doesn't exist
	result = m.GetOrSet("key2", defaultValue)
	if result != defaultValue {
		t.Errorf("Expected value to be %d, got %d", defaultValue, result)
	}
	if value := m.Get("key2"); value.Some() != defaultValue {
		t.Errorf("Expected key2 to be set with default value")
	}
}

func TestMapOrdClone(t *testing.T) {
	// Test case 1: Map with elements
	m := NewMapOrd[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")
	cloned := m.Clone()
	if cloned.Len() != m.Len() {
		t.Errorf("Expected cloned map to have length %d, got %d", m.Len(), cloned.Len())
	}
	cloned.Iter().ForEach(func(k int, v string) {
		if m.Get(k).Some() != v {
			t.Errorf("Expected cloned map to have key-value pair %d:%s", k, v)
		}
	})

	// Test case 2: Empty Map
	m2 := NewMapOrd[string, int]()
	cloned2 := m2.Clone()
	if cloned2.Len() != 0 {
		t.Errorf("Expected cloned map of an empty map to be empty")
	}
}

func TestMapOrdCopy(t *testing.T) {
	// Test case 1: Map with elements
	m := NewMapOrd[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")

	src := NewMapOrd[int, string]()
	src.Set(4, "d")
	src.Set(5, "e")

	m.Copy(src)
	if m.Len() != 5 {
		t.Errorf("Expected copied map to have length %d, got %d", 5, m.Len())
	}

	src.Iter().ForEach(func(k int, v string) {
		if m.Get(k).Some() != v {
			t.Errorf("Expected copied map to have key-value pair %d:%s", k, v)
		}
	})

	// Test case 2: Empty Source Map
	m2 := NewMapOrd[string, int]()
	src2 := NewMapOrd[string, int]()
	m2.Copy(src2)
	if m2.Len() != 0 {
		t.Errorf("Expected copied map of an empty source map to be empty")
	}
}

func TestMapOrdSortBy(t *testing.T) {
	// Test case 1: Sort by key
	m := MapOrd[string, int]{
		{"b", 2},
		{"c", 3},
		{"a", 1},
	}
	m.SortBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Key, b.Key) })
	expectedKeyOrder := []string{"a", "b", "c"}
	for i, p := range m {
		if p.Key != expectedKeyOrder[i] {
			t.Errorf("Expected key at index %d to be %s, got %s", i, expectedKeyOrder[i], p.Key)
		}
	}

	// Test case 2: Sort by value
	m2 := MapOrd[string, int]{
		{"a", 3},
		{"b", 1},
		{"c", 2},
	}
	m2.SortBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Value, b.Value) })
	expectedValueOrder := []int{1, 2, 3}
	for i, p := range m2 {
		if p.Value != expectedValueOrder[i] {
			t.Errorf("Expected value at index %d to be %d, got %d", i, expectedValueOrder[i], p.Value)
		}
	}
}

func TestSortByKey(t *testing.T) {
	// Create a sample MapOrd to test sorting by keys
	mo := MapOrd[string, int]{
		{"b", 2},
		{"a", 1},
		{"c", 3},
	}

	// Sort the MapOrd by keys using the custom comparison function
	mo.SortByKey(cmp.Cmp)

	// Expected sorted order by keys
	expected := MapOrd[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	// Check if the MapOrd is sorted as expected
	if !slices.Equal(mo, expected) {
		t.Errorf("SortByKey failed: expected %v, got %v", expected, mo)
	}
}

func TestSortByValue(t *testing.T) {
	// Create a sample MapOrd to test sorting by values
	mo := MapOrd[string, int]{
		{"a", 2},
		{"b", 1},
		{"c", 3},
	}

	// Define a custom comparison function for integers
	customIntCmp := func(a, b int) cmp.Ordering {
		if a < b {
			return cmp.Less
		} else if a > b {
			return cmp.Greater
		}
		return cmp.Equal
	}

	// Sort the MapOrd by values using the custom comparison function
	mo.SortByValue(customIntCmp)

	// Expected sorted order by values
	expected := MapOrd[string, int]{
		{"b", 1},
		{"a", 2},
		{"c", 3},
	}

	// Check if the MapOrd is sorted as expected
	if !slices.Equal(mo, expected) {
		t.Errorf("SortByValue failed: expected %v, got %v", expected, mo)
	}
}

func TestSortIterByKey(t *testing.T) {
	// Create a sample MapOrd to test sorting by keys
	mo := MapOrd[string, int]{
		{"b", 2},
		{"a", 1},
		{"c", 3},
	}

	// Sort the MapOrd by keys using the custom comparison function
	mo = mo.Iter().SortByKey(cmp.Cmp).Collect()

	// Expected sorted order by keys
	expected := MapOrd[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	// Check if the MapOrd is sorted as expected
	if !slices.Equal(mo, expected) {
		t.Errorf("SortByKey failed: expected %v, got %v", expected, mo)
	}
}

func TestSortIterByValue(t *testing.T) {
	// Create a sample MapOrd to test sorting by values
	mo := MapOrd[string, int]{
		{"a", 2},
		{"b", 1},
		{"c", 3},
	}

	// Define a custom comparison function for integers
	customIntCmp := func(a, b int) cmp.Ordering {
		if a < b {
			return cmp.Less
		} else if a > b {
			return cmp.Greater
		}
		return cmp.Equal
	}

	// Sort the MapOrd by values using the custom comparison function
	mo = mo.Iter().SortByValue(customIntCmp).Collect()

	// Expected sorted order by values
	expected := MapOrd[string, int]{
		{"b", 1},
		{"a", 2},
		{"c", 3},
	}

	// Check if the MapOrd is sorted as expected
	if !slices.Equal(mo, expected) {
		t.Errorf("SortByValue failed: expected %v, got %v", expected, mo)
	}
}

func TestMapOrdFromMap(t *testing.T) {
	// Test case 1: Map with elements
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	mapOrd := m.ToMapOrd()
	mapOrd.SortBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Value, b.Value) })

	expected := []Pair[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	for i, p := range mapOrd {
		if p != expected[i] {
			t.Errorf("Expected mapOrd[%d] to be %v, got %v", i, expected[i], p)
		}
	}

	// Test case 2: Empty Map
	m2 := NewMap[string, int]()
	mapOrd2 := m2.ToMapOrd()
	if len(mapOrd2) != 0 {
		t.Errorf("Expected mapOrd2 to be empty")
	}
}

func TestMapOrdFromStd(t *testing.T) {
	// Test case 1: Map with elements
	inputMap := map[string]int{"a": 1, "b": 2, "c": 3}
	orderedMap := MapOrdFromStd(inputMap)
	if len(orderedMap) != len(inputMap) {
		t.Errorf("Expected ordered map to have length %d, got %d", len(inputMap), orderedMap.Len())
	}
	for key, value := range inputMap {
		if orderedMap.Get(key).Some() != value {
			t.Errorf("Expected ordered map to have key-value pair %s:%d", key, value)
		}
	}

	// Test case 2: Empty Map
	emptyMap := map[string]int{}
	orderedEmptyMap := MapOrdFromStd(emptyMap)
	if orderedEmptyMap.Len() != 0 {
		t.Errorf("Expected ordered map of an empty map to be empty")
	}
}

func TestMapOrdEq(t *testing.T) {
	// Test case 1: Equal maps
	m1 := NewMapOrd[string, int]()
	m1.Set("a", 1)
	m1.Set("b", 2)
	m1.Set("c", 3)

	m2 := NewMapOrd[string, int]()
	m2.Set("a", 1)
	m2.Set("b", 2)
	m2.Set("c", 3)

	if !m1.Eq(m2) {
		t.Errorf("Expected maps to be equal")
	}

	// Test case 2: Unequal maps (different lengths)
	m3 := NewMapOrd[string, int]()
	m3.Set("a", 1)
	m3.Set("b", 2)

	if m1.Eq(m3) {
		t.Errorf("Expected maps to be unequal")
	}

	// Test case 3: Unequal maps (different values)
	m4 := NewMapOrd[string, int]()
	m4.Set("a", 1)
	m4.Set("b", 3)
	m4.Set("c", 3)

	if m1.Eq(m4) {
		t.Errorf("Expected maps to be unequal")
	}

	// Test case 4
	if !NewMapOrd[int, int]().Eq(NewMapOrd[int, int]()) {
		t.Errorf("Empty ordered maps should be considered equal")
	}
}

func TestMapOrdIterInspect(t *testing.T) {
	// Define an ordered map to iterate over
	mo := NewMapOrd[int, string]()
	mo.Set(1, "one")
	mo.Set(2, "two")
	mo.Set(3, "three")

	// Define a slice to store the inspected key-value pairs
	inspectedPairs := NewMapOrd[int, string]()

	// Create a new iterator with Inspect and collect the key-value pairs
	mo.Iter().Inspect(func(k int, v string) {
		inspectedPairs.Set(k, v)
	}).Collect()

	if mo.Len() != inspectedPairs.Len() {
		t.Errorf("Expected inspected pairs to have length %d, got %d", mo.Len(), inspectedPairs.Len())
	}

	if mo.Ne(inspectedPairs) {
		t.Errorf("Expected inspected pairs to be equal to the original map")
	}
}

func TestMapOrdIterChain(t *testing.T) {
	// Define the first ordered map to iterate over
	iter1 := NewMapOrd[int, string]()
	iter1.Set(1, "a")

	// Define the second ordered map to iterate over
	iter2 := NewMapOrd[int, string]()
	iter2.Set(2, "b")

	// Concatenate the iterators and collect the elements
	chainedIter := iter1.Iter().Chain(iter2.Iter())
	collected := chainedIter.Collect()

	// Verify the concatenated elements
	expected := NewMapOrd[int, string]()
	expected.Set(1, "a")
	expected.Set(2, "b")

	if !collected.Eq(expected) {
		t.Errorf("Concatenated map does not match expected map")
	}
}

func TestMapOrdIterCount(t *testing.T) {
	// Create a new ordered map
	seq := NewMapOrd[int, string]()
	seq.Set(1, "a")
	seq.Set(2, "b")
	seq.Set(3, "c")

	// Count the number of iterations
	count := seq.Iter().Count()

	// Verify the count
	expectedCount := Int(3) // Since there are 3 elements in the ordered map
	if count != expectedCount {
		t.Errorf("Expected count to be %d, but got %d", expectedCount, count)
	}
}

func TestMapOrdIterSkip(t *testing.T) {
	// Create a new ordered map
	seq := NewMapOrd[int, string]()
	seq.Set(1, "a")
	seq.Set(2, "b")
	seq.Set(3, "c")
	seq.Set(4, "d")

	// Skip the first two elements
	skipped := seq.Iter().Skip(2)

	// Collect the elements after skipping
	collected := skipped.Collect()

	// Verify the collected elements
	expected := NewMapOrd[int, string]()
	expected.Set(3, "c")
	expected.Set(4, "d")

	if !collected.Eq(expected) {
		t.Errorf("Expected %v, but got %v", expected, collected)
	}
}

func TestMapOrdIterExclude(t *testing.T) {
	// Create a new ordered map
	mo := NewMapOrd[int, int]()
	mo.Set(1, 1)
	mo.Set(2, 2)
	mo.Set(3, 3)
	mo.Set(4, 4)
	mo.Set(5, 5)

	// Exclude even values
	notEven := mo.Iter().Exclude(func(k, v int) bool {
		return v%2 == 0
	})

	// Collect the resulting elements
	collected := notEven.Collect()

	// Verify the collected elements
	expected := NewMapOrd[int, int]()
	expected.Set(1, 1)
	expected.Set(3, 3)
	expected.Set(5, 5)

	if !collected.Eq(expected) {
		t.Errorf("Expected %v, but got %v", expected, collected)
	}
}

func TestMapOrdIterFilter(t *testing.T) {
	// Create a new ordered map
	mo := NewMapOrd[int, int]()
	mo.Set(1, 1)
	mo.Set(2, 2)
	mo.Set(3, 3)
	mo.Set(4, 4)
	mo.Set(5, 5)

	// Filter even values
	even := mo.Iter().Filter(func(k, v int) bool {
		return v%2 == 0
	})

	// Collect the resulting elements
	collected := even.Collect()

	// Verify the collected elements
	expected := NewMapOrd[int, int]()
	expected.Set(2, 2)
	expected.Set(4, 4)

	if !collected.Eq(expected) {
		t.Errorf("Expected %v, but got %v", expected, collected)
	}
}

func TestMapOrdIterFind(t *testing.T) {
	// Create a new ordered map
	mo := NewMapOrd[int, int]()
	mo.Set(1, 1)
	mo.Set(2, 2)
	mo.Set(3, 3)
	mo.Set(4, 4)
	mo.Set(5, 5)

	// Find the first even value
	found := mo.Iter().Find(func(k, v int) bool {
		return v%2 == 0
	}).Some()

	// Verify the found element
	expected := Pair[int, int]{2, 2}

	if !reflect.DeepEqual(found, expected) {
		t.Errorf("Expected %v, but got %v", expected, found)
	}
}

func TestMapOrdIterMap(t *testing.T) {
	// Create a new ordered map
	mo := NewMapOrd[int, int]()
	mo.Set(1, 1)
	mo.Set(2, 2)
	mo.Set(3, 3)
	mo.Set(4, 4)
	mo.Set(5, 5)

	// Map each key-value pair to its square
	squared := mo.Iter().Map(func(k, v int) (int, int) {
		return k * k, v * v
	})

	// Collect the resulting elements
	collected := squared.Collect()

	// Verify the collected elements
	expected := NewMapOrd[int, int]()
	expected.Set(1, 1)
	expected.Set(4, 4)
	expected.Set(9, 9)
	expected.Set(16, 16)
	expected.Set(25, 25)

	if !collected.Eq(expected) {
		t.Errorf("Expected %v, but got %v", expected, collected)
	}
}

func TestMapOrdIterTake(t *testing.T) {
	// Create a new ordered map
	mo := NewMapOrd[int, int]()
	mo.Set(1, 1)
	mo.Set(2, 2)
	mo.Set(3, 3)
	mo.Set(4, 4)
	mo.Set(5, 5)

	// Take the first 3 elements
	taken := mo.Iter().Take(3)

	// Collect the resulting elements
	collected := taken.Collect()

	// Verify the collected elements
	expected := NewMapOrd[int, int]()
	expected.Set(1, 1)
	expected.Set(2, 2)
	expected.Set(3, 3)

	if !collected.Eq(expected) {
		t.Errorf("Expected %v, but got %v", expected, collected)
	}
}

func TestMapOrdIterToChannel(t *testing.T) {
	// Create a new ordered map
	mo := NewMapOrd[int, int]()
	mo.Set(1, 1)
	mo.Set(2, 2)
	mo.Set(3, 3)
	mo.Set(4, 4)
	mo.Set(5, 5)

	// Convert the iterator to a channel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation to avoid goroutine leaks.

	ch := mo.Iter().ToChan(ctx)

	// Collect elements from the channel
	collected := NewMapOrd[int, int]()
	for pair := range ch {
		collected.Set(pair.Key, pair.Value)
	}

	// Verify the collected elements
	expected := NewMapOrd[int, int]()
	expected.Set(1, 1)
	expected.Set(2, 2)
	expected.Set(3, 3)
	expected.Set(4, 4)
	expected.Set(5, 5)

	if !collected.Eq(expected) {
		t.Errorf("Expected %v, but got %v", expected, collected)
	}
}

func TestMapOrdIterUnzip(t *testing.T) {
	// Create a new ordered map
	mo := NewMapOrd[string, int]()
	mo.Set("a", 1)
	mo.Set("b", 2)
	mo.Set("c", 3)

	// Unzip the ordered map
	keys, values := mo.Iter().Unzip()

	// Verify the keys
	expectedKeys := SliceOf("a", "b", "c")
	if keys.Collect().Ne(expectedKeys) {
		t.Errorf("Expected keys %v, but got %v", expectedKeys, keys)
	}

	// Verify the values
	expectedValues := SliceOf(1, 2, 3)
	if values.Collect().Ne(expectedValues) {
		t.Errorf("Expected values %v, but got %v", expectedValues, values)
	}
}

func TestMapOrdShuffle(t *testing.T) {
	mo := NewMapOrd[int, int]()
	for i := 1; i <= 5; i++ {
		mo.Set(i, i)
	}

	clone := mo.Clone()
	mo.Shuffle()

	if mo.Eq(clone) {
		t.Errorf("The order of elements has not changed after shuffle")
	}
}

func TestMapOrdTransform(t *testing.T) {
	// Исходные данные
	original := MapOrd[string, int]{
		{"a", 1},
		{"b", 2},
	}

	addEntry := func(mo MapOrd[string, int]) MapOrd[string, int] {
		return append(mo, Pair[string, int]{"c", 3})
	}

	expected := MapOrd[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	result := original.Transform(addEntry)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Transform failed: expected %v, got %v", expected, result)
	}

	removeEntry := func(mo MapOrd[string, int]) MapOrd[string, int] {
		filtered := MapOrd[string, int]{}
		for _, pair := range mo {
			if pair.Key != "a" {
				filtered = append(filtered, pair)
			}
		}
		return filtered
	}

	expectedAfterRemoval := MapOrd[string, int]{
		{"b", 2},
		{"c", 3},
	}

	resultAfterRemoval := result.Transform(removeEntry)

	if !reflect.DeepEqual(resultAfterRemoval, expectedAfterRemoval) {
		t.Errorf("Transform with removal failed: expected %v, got %v", expectedAfterRemoval, resultAfterRemoval)
	}
}
