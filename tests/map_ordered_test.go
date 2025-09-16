package g_test

import (
	"reflect"
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
		orderedMap := NewMapOrd[string, int]()
		orderedMap.Set("a", 1)
		orderedMap.Set("b", 2)
		orderedMap.Set("c", 3)
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
		orderedMap := NewMapOrd[string, int]()
		orderedMap.Set("a", 1)
		orderedMap.Set("b", 2)
		orderedMap.Set("c", 3)

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
		emptyMap := NewMapOrd[string, int]()
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

func TestMapOrdDelete(t *testing.T) {
	t.Run("Delete single existing key", func(t *testing.T) {
		m := NewMapOrd[int, string]()
		m.Set(1, "a")
		m.Set(2, "b")
		m.Set(3, "c")

		m.Delete(2)

		expected := NewMapOrd[int, string]()
		expected.Set(1, "a")
		expected.Set(3, "c")

		if m.Ne(expected) {
			t.Errorf("Expected %v, got %v", expected, m)
		}
	})

	t.Run("Delete multiple keys", func(t *testing.T) {
		m := NewMapOrd[int, string]()
		m.Set(1, "a")
		m.Set(2, "b")
		m.Set(3, "c")
		m.Set(4, "d")

		m.Delete(2, 4)

		expected := NewMapOrd[int, string]()
		expected.Set(1, "a")
		expected.Set(3, "c")

		if m.Ne(expected) {
			t.Errorf("Expected %v, got %v", expected, m)
		}
	})

	t.Run("Delete non-existing key", func(t *testing.T) {
		m := NewMapOrd[int, string]()
		m.Set(1, "a")
		m.Set(2, "b")

		m.Delete(3) // key 3 does not exist

		expected := NewMapOrd[int, string]()
		expected.Set(1, "a")
		expected.Set(2, "b")

		if m.Ne(expected) {
			t.Errorf("Expected %v, got %v", expected, m)
		}
	})

	t.Run("Delete all keys", func(t *testing.T) {
		m := NewMapOrd[int, string]()
		m.Set(1, "a")
		m.Set(2, "b")

		m.Delete(1, 2)

		expected := NewMapOrd[int, string]()

		if m.Ne(expected) {
			t.Errorf("Expected empty map, got %v", m)
		}
	})

	t.Run("Delete from empty map", func(t *testing.T) {
		m := NewMapOrd[int, string]()

		m.Delete(1, 2, 3)

		expected := NewMapOrd[int, string]()

		if m.Ne(expected) {
			t.Errorf("Expected empty map, got %v", m)
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

	// Test case 3: Map with string keys (hashable)
	m2 := NewMapOrd[string, []int]()
	m2.Set("key", []int{1})

	if !m2.Contains("key") {
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
	m := NewMapOrd[string, int]()
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("a", 1)
	m.SortBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Key, b.Key) })
	expectedKeyOrder := []string{"a", "b", "c"}
	i := 0
	m.Iter().ForEach(func(key string, _ int) {
		if key != expectedKeyOrder[i] {
			t.Errorf("Expected key at index %d to be %s, got %s", i, expectedKeyOrder[i], key)
		}
		i++
	})

	// Test case 2: Sort by value
	m2 := NewMapOrd[string, int]()
	m2.Set("a", 3)
	m2.Set("b", 1)
	m2.Set("c", 2)
	m2.SortBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Value, b.Value) })
	expectedValueOrder := []int{1, 2, 3}
	i = 0
	m2.Iter().ForEach(func(_ string, value int) {
		if value != expectedValueOrder[i] {
			t.Errorf("Expected value at index %d to be %d, got %d", i, expectedValueOrder[i], value)
		}
		i++
	})
}

func TestSortByKey(t *testing.T) {
	// Create a sample MapOrd to test sorting by keys
	mo := NewMapOrd[string, int]()
	mo.Set("b", 2)
	mo.Set("a", 1)
	mo.Set("c", 3)

	// Sort the MapOrd by keys using the custom comparison function
	mo.SortByKey(cmp.Cmp)

	// Expected sorted order by keys
	expected := NewMapOrd[string, int]()
	expected.Set("a", 1)
	expected.Set("b", 2)
	expected.Set("c", 3)

	// Check if the MapOrd is sorted as expected
	if !mo.Eq(expected) {
		t.Errorf("SortByKey failed: expected %v, got %v", expected, mo)
	}
}

func TestSortByValue(t *testing.T) {
	// Create a sample MapOrd to test sorting by values
	mo := NewMapOrd[string, int]()
	mo.Set("a", 2)
	mo.Set("b", 1)
	mo.Set("c", 3)

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
	expected := NewMapOrd[string, int]()
	expected.Set("b", 1)
	expected.Set("a", 2)
	expected.Set("c", 3)

	// Check if the MapOrd is sorted as expected
	if !mo.Eq(expected) {
		t.Errorf("SortByValue failed: expected %v, got %v", expected, mo)
	}
}

func TestSortIterByKey(t *testing.T) {
	// Create a sample MapOrd to test sorting by keys
	mo := NewMapOrd[string, int]()
	mo.Set("b", 2)
	mo.Set("a", 1)
	mo.Set("c", 3)

	// Sort the MapOrd by keys using the custom comparison function
	mo = mo.Iter().SortByKey(cmp.Cmp).Collect()

	// Expected sorted order by keys
	expected := NewMapOrd[string, int]()
	expected.Set("a", 1)
	expected.Set("b", 2)
	expected.Set("c", 3)

	// Check if the MapOrd is sorted as expected
	if !mo.Eq(expected) {
		t.Errorf("SortByKey failed: expected %v, got %v", expected, mo)
	}
}

func TestSortIterByValue(t *testing.T) {
	// Create a sample MapOrd to test sorting by values
	mo := NewMapOrd[string, int]()
	mo.Set("a", 2)
	mo.Set("b", 1)
	mo.Set("c", 3)

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
	expected := NewMapOrd[string, int]()
	expected.Set("b", 1)
	expected.Set("a", 2)
	expected.Set("c", 3)

	// Check if the MapOrd is sorted as expected
	if !mo.Eq(expected) {
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

	i := 0
	mapOrd.Iter().ForEach(func(key string, value int) {
		if i < len(expected) && (key != expected[i].Key || value != expected[i].Value) {
			t.Errorf("Expected mapOrd[%d] to be %v, got {%s, %d}", i, expected[i], key, value)
		}
		i++
	})

	// Test case 2: Empty Map
	m2 := NewMap[string, int]()
	mapOrd2 := m2.ToMapOrd()
	if mapOrd2.Len() != 0 {
		t.Errorf("Expected mapOrd2 to be empty")
	}
}

func TestMapOrdFromStd(t *testing.T) {
	// Test case 1: Map with elements
	inputMap := map[string]int{"a": 1, "b": 2, "c": 3}
	orderedMap := MapOrdFromStd(inputMap)
	if int(orderedMap.Len()) != len(inputMap) {
		t.Errorf("Expected ordered map to have length %d, got %d", len(inputMap), orderedMap.Len())
	}
	for key, value := range inputMap {
		if orderedMap.Get(key).Some() != value {
			t.Errorf("Expected ordered map to have key-value pair %s:%d", key, value)
		}
	}

	// Test case 2: Empty Map
	emptyMap := make(map[string]int)
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
	notEven := mo.Iter().Exclude(func(_, v int) bool {
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
	even := mo.Iter().Filter(func(_, v int) bool {
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
	found := mo.Iter().Find(func(_, v int) bool {
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
	ctx := t.Context()

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
	original := NewMapOrd[string, int]()
	original.Set("a", 1)
	original.Set("b", 2)

	addEntry := func(mo MapOrd[string, int]) MapOrd[string, int] {
		result := mo.Clone()
		// Add new entry
		result.Set("c", 3)
		return result
	}

	expected := NewMapOrd[string, int]()
	expected.Set("a", 1)
	expected.Set("b", 2)
	expected.Set("c", 3)

	result := original.Transform(addEntry)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Transform failed: expected %v, got %v", expected, result)
	}

	removeEntry := func(mo MapOrd[string, int]) MapOrd[string, int] {
		filtered := NewMapOrd[string, int]()
		mo.Iter().ForEach(func(key string, value int) {
			if key != "a" {
				filtered.Set(key, value)
			}
		})
		return filtered
	}

	expectedAfterRemoval := NewMapOrd[string, int]()
	expectedAfterRemoval.Set("b", 2)
	expectedAfterRemoval.Set("c", 3)

	resultAfterRemoval := result.Transform(removeEntry)

	if !reflect.DeepEqual(resultAfterRemoval, expectedAfterRemoval) {
		t.Errorf("Transform with removal failed: expected %v, got %v", expectedAfterRemoval, resultAfterRemoval)
	}
}

func TestSortByKeyIterator(t *testing.T) {
	// Test SeqMapOrd.SortByKey method
	mo := NewMapOrd[string, int]()
	mo.Set("c", 3)
	mo.Set("a", 1)
	mo.Set("b", 2)

	sorted := mo.Iter().SortByKey(cmp.Cmp).Collect()

	expected := NewMapOrd[string, int]()
	expected.Set("a", 1)
	expected.Set("b", 2)
	expected.Set("c", 3)

	if !sorted.Eq(expected) {
		t.Errorf("SortByKey iterator failed: expected %v, got %v", expected, sorted)
	}

	// Test with empty map
	empty := NewMapOrd[string, int]()
	emptySorted := empty.Iter().SortByKey(cmp.Cmp).Collect()
	if !emptySorted.Empty() {
		t.Errorf("SortByKey on empty map should return empty map")
	}
}

func TestSortByValueIterator(t *testing.T) {
	// Test SeqMapOrd.SortByValue method
	mo := NewMapOrd[string, int]()
	mo.Set("a", 3)
	mo.Set("b", 1)
	mo.Set("c", 2)

	sorted := mo.Iter().SortByValue(cmp.Cmp).Collect()

	expected := NewMapOrd[string, int]()
	expected.Set("b", 1)
	expected.Set("c", 2)
	expected.Set("a", 3)

	if !sorted.Eq(expected) {
		t.Errorf("SortByValue iterator failed: expected %v, got %v", expected, sorted)
	}

	// Test with empty map
	empty := NewMapOrd[string, int]()
	emptySorted := empty.Iter().SortByValue(cmp.Cmp).Collect()
	if !emptySorted.Empty() {
		t.Errorf("SortByValue on empty map should return empty map")
	}
}

func TestIsSortedBy(t *testing.T) {
	// Test empty map
	empty := NewMapOrd[string, int]()
	if !empty.IsSortedBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Key, b.Key) }) {
		t.Error("Empty map should be considered sorted")
	}

	// Test single element
	single := NewMapOrd[string, int]()
	single.Set("a", 1)
	if !single.IsSortedBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Key, b.Key) }) {
		t.Error("Single element map should be considered sorted")
	}

	// Test sorted map
	sorted := NewMapOrd[string, int]()
	sorted.Set("a", 1)
	sorted.Set("b", 2)
	sorted.Set("c", 3)
	if !sorted.IsSortedBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Key, b.Key) }) {
		t.Error("Sorted map should return true")
	}

	// Test unsorted map
	unsorted := NewMapOrd[string, int]()
	unsorted.Set("c", 3)
	unsorted.Set("a", 1)
	unsorted.Set("b", 2)
	if unsorted.IsSortedBy(func(a, b Pair[string, int]) cmp.Ordering { return cmp.Cmp(a.Key, b.Key) }) {
		t.Error("Unsorted map should return false")
	}
}

func TestIsSortedByKey(t *testing.T) {
	// Test empty map
	empty := NewMapOrd[string, int]()
	if !empty.IsSortedByKey(cmp.Cmp) {
		t.Error("Empty map should be considered sorted by key")
	}

	// Test single element
	single := NewMapOrd[string, int]()
	single.Set("a", 1)
	if !single.IsSortedByKey(cmp.Cmp) {
		t.Error("Single element map should be considered sorted by key")
	}

	// Test sorted by key map
	sorted := NewMapOrd[string, int]()
	sorted.Set("a", 3)
	sorted.Set("b", 1)
	sorted.Set("c", 2)
	if !sorted.IsSortedByKey(cmp.Cmp) {
		t.Error("Map sorted by key should return true")
	}

	// Test unsorted by key map
	unsorted := NewMapOrd[string, int]()
	unsorted.Set("c", 1)
	unsorted.Set("a", 2)
	unsorted.Set("b", 3)
	if unsorted.IsSortedByKey(cmp.Cmp) {
		t.Error("Map unsorted by key should return false")
	}

	// Test with a map that becomes sorted after SortByKey
	mo := NewMapOrd[string, int]()
	mo.Set("z", 1)
	mo.Set("a", 2)
	mo.Set("m", 3)
	if mo.IsSortedByKey(cmp.Cmp) {
		t.Error("Unsorted map should return false before sorting")
	}
	mo.SortByKey(cmp.Cmp)
	if !mo.IsSortedByKey(cmp.Cmp) {
		t.Error("Map should be sorted by key after SortByKey")
	}
}

func TestIsSortedByValue(t *testing.T) {
	// Test empty map
	empty := NewMapOrd[string, int]()
	if !empty.IsSortedByValue(cmp.Cmp) {
		t.Error("Empty map should be considered sorted by value")
	}

	// Test single element
	single := NewMapOrd[string, int]()
	single.Set("a", 1)
	if !single.IsSortedByValue(cmp.Cmp) {
		t.Error("Single element map should be considered sorted by value")
	}

	// Test sorted by value map
	sorted := NewMapOrd[string, int]()
	sorted.Set("c", 1)
	sorted.Set("a", 2)
	sorted.Set("b", 3)
	if !sorted.IsSortedByValue(cmp.Cmp) {
		t.Error("Map sorted by value should return true")
	}

	// Test unsorted by value map
	unsorted := NewMapOrd[string, int]()
	unsorted.Set("a", 3)
	unsorted.Set("b", 1)
	unsorted.Set("c", 2)
	if unsorted.IsSortedByValue(cmp.Cmp) {
		t.Error("Map unsorted by value should return false")
	}

	// Test with a map that becomes sorted after SortByValue
	mo := NewMapOrd[string, int]()
	mo.Set("a", 99)
	mo.Set("b", 1)
	mo.Set("c", 50)
	if mo.IsSortedByValue(cmp.Cmp) {
		t.Error("Unsorted map should return false before sorting")
	}
	mo.SortByValue(cmp.Cmp)
	if !mo.IsSortedByValue(cmp.Cmp) {
		t.Error("Map should be sorted by value after SortByValue")
	}
}

func TestMapOrdEntryDeleteEdgeCases(t *testing.T) {
	// Test Delete on non-existent key
	mo := NewMapOrd[string, int]()
	mo.Set("a", 1)
	mo.Set("b", 2)

	entry := mo.Entry("nonexistent")
	result := entry.Delete()

	if result.IsSome() {
		t.Errorf("Delete on non-existent key should return None")
	}

	// Verify map unchanged
	if mo.Len() != 2 {
		t.Errorf("Map should be unchanged after deleting non-existent key")
	}

	// Test Delete on existing key
	entry2 := mo.Entry("a")
	result2 := entry2.Delete()

	if !result2.IsSome() || result2.Some() != 1 {
		t.Errorf("Delete on existing key should return Some(1), got %v", result2)
	}

	if mo.Len() != 1 {
		t.Errorf("Map should have 1 element after deletion, got %d", mo.Len())
	}

	// Verify remaining element
	if !mo.Contains("b") || mo.Get("b").Some() != 2 {
		t.Errorf("Remaining element should be 'b': 2")
	}
}

func TestSeqMapOrdFindEdgeCases(t *testing.T) {
	// Test Find with no matching elements
	mo := NewMapOrd[string, int]()
	mo.Set("a", 1)
	mo.Set("b", 2)
	mo.Set("c", 3)

	notFound := mo.Iter().Find(func(_ string, v int) bool {
		return v > 10
	})

	if notFound.IsSome() {
		t.Errorf("Find should return None when no elements match")
	}

	// Test Find with matching element
	found := mo.Iter().Find(func(_ string, v int) bool {
		return v == 2
	})

	if !found.IsSome() || found.Some().Key != "b" || found.Some().Value != 2 {
		t.Errorf("Find should return Some(Pair{b, 2}), got %v", found)
	}

	// Test Find on empty map
	empty := NewMapOrd[string, int]()
	emptyResult := empty.Iter().Find(func(string, int) bool { return true })

	if emptyResult.IsSome() {
		t.Errorf("Find on empty map should return None")
	}
}

func TestIterReverseEmptyMap(t *testing.T) {
	// Test IterReverse on empty map
	empty := NewMapOrd[string, int]()

	count := 0
	empty.IterReverse().ForEach(func(string, int) {
		count++
	})

	if count != 0 {
		t.Errorf("IterReverse on empty map should iterate 0 times, got %d", count)
	}
}

func TestKeysEmptyMap(t *testing.T) {
	// Test Keys method on empty map
	empty := NewMapOrd[string, int]()
	keys := empty.Keys()

	if !keys.Empty() {
		t.Errorf("Keys on empty map should return empty slice")
	}

	if keys.Len() != 0 {
		t.Errorf("Keys on empty map should have length 0, got %d", keys.Len())
	}
}

func TestAsAnyEdgeCases(t *testing.T) {
	// Test AsAny on empty map
	empty := NewMapOrd[string, int]()
	anyEmpty := empty.AsAny()

	if !anyEmpty.Empty() {
		t.Errorf("AsAny on empty map should return empty map")
	}

	// Test AsAny with different types
	mo := NewMapOrd[int, string]()
	mo.Set(1, "one")
	mo.Set(2, "two")

	anyMo := mo.AsAny()

	if anyMo.Len() != 2 {
		t.Errorf("AsAny should preserve length, got %d", anyMo.Len())
	}

	// Verify values are preserved
	val1 := anyMo.Get(1)
	if !val1.IsSome() || val1.Some().(string) != "one" {
		t.Errorf("AsAny should preserve values, got %v", val1)
	}

	val2 := anyMo.Get(2)
	if !val2.IsSome() || val2.Some().(string) != "two" {
		t.Errorf("AsAny should preserve values, got %v", val2)
	}
}

// go test -bench=. -benchmem -count=4

func genMO() MapOrd[String, int] {
	mo := NewMapOrd[String, int](10000)
	for i := range 10000 {
		mo.Set(Int(i).String(), i)
	}

	return mo
}

func BenchmarkMoContains(b *testing.B) {
	mo := genMO()

	for b.Loop() {
		_ = mo.Contains("9999")
	}
}

func BenchmarkMoEq(b *testing.B) {
	mo := genMO()
	mo2 := mo.Clone()

	b.ResetTimer()

	for b.Loop() {
		_ = mo.Eq(mo2)
	}
}

func BenchmarkMoGet(b *testing.B) {
	mo := genMO()

	for b.Loop() {
		_ = mo.Get("9999")
	}
}

func TestMapOrdPrint(t *testing.T) {
	m := NewMapOrd[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)

	// Just test that Print() doesn't panic and returns the map
	result := m.Print()
	if !result.Eq(m) {
		t.Errorf("Print() should return the same map")
	}
}

func TestMapOrdPrintln(t *testing.T) {
	m := NewMapOrd[string, int]()
	m.Set("x", 10)
	m.Set("y", 20)

	// Just test that Println() doesn't panic and returns the map
	result := m.Println()
	if !result.Eq(m) {
		t.Errorf("Println() should return the same map")
	}
}

func TestMapOrdSet(t *testing.T) {
	t.Run("setting_new_key", func(t *testing.T) {
		m := NewMapOrd[string, int]()

		// Setting a new key should return None
		prev := m.Set("new_key", 100)
		if prev.IsSome() {
			t.Errorf("Expected None when setting new key, got Some(%v)", prev.Some())
		}

		// Verify the key was set
		if val := m.Get("new_key"); val.IsNone() || val.Some() != 100 {
			t.Errorf("Expected Some(100), got %v", val)
		}

		if m.Len() != 1 {
			t.Errorf("Expected length 1, got %v", m.Len())
		}
	})

	t.Run("updating_existing_key", func(t *testing.T) {
		m := NewMapOrd[string, int]()

		// First set
		prev := m.Set("existing_key", 100)
		if prev.IsSome() {
			t.Errorf("Expected None on first set, got Some(%v)", prev.Some())
		}

		// Update existing key should return previous value
		prev = m.Set("existing_key", 200)
		if prev.IsNone() {
			t.Errorf("Expected Some(100) when updating existing key, got None")
		}
		if prev.Some() != 100 {
			t.Errorf("Expected Some(100) when updating existing key, got Some(%v)", prev.Some())
		}

		// Verify the key was updated
		if val := m.Get("existing_key"); val.IsNone() || val.Some() != 200 {
			t.Errorf("Expected Some(200), got %v", val)
		}

		// Length should still be 1
		if m.Len() != 1 {
			t.Errorf("Expected length 1, got %v", m.Len())
		}
	})

	t.Run("multiple_updates_same_key", func(t *testing.T) {
		m := NewMapOrd[int, string]()

		// Chain of updates on same key
		prev1 := m.Set(1, "first")
		prev2 := m.Set(1, "second")
		prev3 := m.Set(1, "third")

		if prev1.IsSome() {
			t.Errorf("Expected None on first set, got Some(%v)", prev1.Some())
		}
		if prev2.IsNone() || prev2.Some() != "first" {
			t.Errorf("Expected Some('first'), got %v", prev2)
		}
		if prev3.IsNone() || prev3.Some() != "second" {
			t.Errorf("Expected Some('second'), got %v", prev3)
		}

		// Final value should be "third"
		if val := m.Get(1); val.IsNone() || val.Some() != "third" {
			t.Errorf("Expected Some('third'), got %v", val)
		}

		// Still only one entry
		if m.Len() != 1 {
			t.Errorf("Expected length 1, got %v", m.Len())
		}
	})

	t.Run("mixed_new_and_existing_keys", func(t *testing.T) {
		m := NewMapOrd[string, int]()

		// Set multiple new keys
		prev1 := m.Set("a", 1)
		prev2 := m.Set("b", 2)
		prev3 := m.Set("c", 3)

		// All should return None (new keys)
		if prev1.IsSome() || prev2.IsSome() || prev3.IsSome() {
			t.Errorf("Expected None for all new keys")
		}

		// Update existing keys
		prev4 := m.Set("b", 20) // Update middle
		prev5 := m.Set("a", 10) // Update first
		prev6 := m.Set("c", 30) // Update last

		// Should return previous values
		if prev4.IsNone() || prev4.Some() != 2 {
			t.Errorf("Expected Some(2), got %v", prev4)
		}
		if prev5.IsNone() || prev5.Some() != 1 {
			t.Errorf("Expected Some(1), got %v", prev5)
		}
		if prev6.IsNone() || prev6.Some() != 3 {
			t.Errorf("Expected Some(3), got %v", prev6)
		}

		// Verify final values
		if val := m.Get("a"); val.IsNone() || val.Some() != 10 {
			t.Errorf("Expected Some(10) for 'a', got %v", val)
		}
		if val := m.Get("b"); val.IsNone() || val.Some() != 20 {
			t.Errorf("Expected Some(20) for 'b', got %v", val)
		}
		if val := m.Get("c"); val.IsNone() || val.Some() != 30 {
			t.Errorf("Expected Some(30) for 'c', got %v", val)
		}

		// Should have 3 entries
		if m.Len() != 3 {
			t.Errorf("Expected length 3, got %v", m.Len())
		}
	})

	t.Run("edge_case_zero_values", func(t *testing.T) {
		m := NewMapOrd[int, int]()

		// Set zero value - should still work correctly
		prev := m.Set(0, 0)
		if prev.IsSome() {
			t.Errorf("Expected None when setting new key 0, got Some(%v)", prev.Some())
		}

		// Update with zero value
		prev = m.Set(0, 1)
		if prev.IsNone() || prev.Some() != 0 {
			t.Errorf("Expected Some(0) when updating key 0, got %v", prev)
		}

		// Update back to zero
		prev = m.Set(0, 0)
		if prev.IsNone() || prev.Some() != 1 {
			t.Errorf("Expected Some(1) when updating key 0, got %v", prev)
		}

		// Final value should be 0
		if val := m.Get(0); val.IsNone() || val.Some() != 0 {
			t.Errorf("Expected Some(0), got %v", val)
		}
	})
}

func TestMapOrdAsAny(t *testing.T) {
	m := NewMapOrd[string, int]()
	m.Set("test", 42)

	// Test AsAny() conversion - just verify it doesn't panic and returns something
	anyMap := m.AsAny()

	// Should return a MapOrd[any, any] type, just verify it's not nil and has expected length
	if anyMap.Len() != m.Len() {
		t.Errorf("AsAny() should preserve length: expected %d, got %d", m.Len(), anyMap.Len())
	}
}

func TestMapOrdNonHashableKeys(t *testing.T) {
	t.Run("FunctionKeys", func(t *testing.T) {
		mo := NewMapOrd[func(int) int, string]()

		fn1 := func(x int) int { return x + 1 }
		fn2 := func(x int) int { return x + 2 }
		fn3 := func(x int) int { return x + 3 }

		mo.Set(fn1, "function1")
		mo.Set(fn2, "function2")
		mo.Set(fn3, "function3")

		if mo.Len() != 3 {
			t.Errorf("Expected length 3, got %d", mo.Len())
		}

		if result := mo.Get(fn1); result.IsNone() || result.Some() != "function1" {
			t.Errorf("Expected to find fn1 -> function1, got %v", result)
		}

		if result := mo.Get(fn2); result.IsNone() || result.Some() != "function2" {
			t.Errorf("Expected to find fn2 -> function2, got %v", result)
		}

		if !mo.Contains(fn1) {
			t.Errorf("Expected to contain fn1")
		}

		if !mo.Contains(fn2) {
			t.Errorf("Expected to contain fn2")
		}

		mo.Delete(fn2)
		if mo.Len() != 2 {
			t.Errorf("Expected length 2 after delete, got %d", mo.Len())
		}

		if mo.Contains(fn2) {
			t.Errorf("Expected fn2 to be deleted")
		}

		if result := mo.Get(fn1); result.IsNone() || result.Some() != "function1" {
			t.Errorf("fn1 should still exist after deleting fn2, got %v", result)
		}
	})

	t.Run("SliceKeys", func(t *testing.T) {
		mo := NewMapOrd[[]int, string]()

		slice1 := []int{1, 2, 3}
		slice2 := []int{4, 5, 6}
		slice3 := []int{1, 2, 3}

		mo.Set(slice1, "slice1")
		mo.Set(slice2, "slice2")
		mo.Set(slice3, "slice3")

		if mo.Len() != 3 {
			t.Errorf("Expected length 3, got %d", mo.Len())
		}

		if result := mo.Get(slice1); result.IsNone() || result.Some() != "slice1" {
			t.Errorf("Expected to find slice1, got %v", result)
		}

		if result := mo.Get(slice2); result.IsNone() || result.Some() != "slice2" {
			t.Errorf("Expected to find slice2, got %v", result)
		}

		if result := mo.Get(slice3); result.IsNone() || result.Some() != "slice3" {
			t.Errorf("Expected to find slice3, got %v", result)
		}
	})

	t.Run("MapKeys", func(t *testing.T) {
		mo := NewMapOrd[map[string]int, string]()

		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"c": 3, "d": 4}

		mo.Set(map1, "map1")
		mo.Set(map2, "map2")

		if mo.Len() != 2 {
			t.Errorf("Expected length 2, got %d", mo.Len())
		}

		if result := mo.Get(map1); result.IsNone() || result.Some() != "map1" {
			t.Errorf("Expected to find map1, got %v", result)
		}

		if result := mo.Get(map2); result.IsNone() || result.Some() != "map2" {
			t.Errorf("Expected to find map2, got %v", result)
		}

		if !mo.Contains(map1) {
			t.Errorf("Expected to contain map1")
		}

		mo.Delete(map1)
		if mo.Len() != 1 {
			t.Errorf("Expected length 1 after delete, got %d", mo.Len())
		}

		if mo.Contains(map1) {
			t.Errorf("Expected map1 to be deleted")
		}
	})

	t.Run("MixedKeys", func(t *testing.T) {
		mo := NewMapOrd[any, string]()

		mo.Set("string_key", "string_value")
		mo.Set(42, "int_value")
		mo.Set(3.14, "float_value")

		fn := func() {}
		slice := []int{1, 2, 3}
		m := map[string]int{"x": 1}

		mo.Set(fn, "function_value")
		mo.Set(slice, "slice_value")
		mo.Set(m, "map_value")

		if mo.Len() != 6 {
			t.Errorf("Expected length 6, got %d", mo.Len())
		}

		if result := mo.Get("string_key"); result.IsNone() || result.Some() != "string_value" {
			t.Errorf("Expected to find string_key, got %v", result)
		}

		if result := mo.Get(42); result.IsNone() || result.Some() != "int_value" {
			t.Errorf("Expected to find int key 42, got %v", result)
		}

		if result := mo.Get(fn); result.IsNone() || result.Some() != "function_value" {
			t.Errorf("Expected to find function key, got %v", result)
		}

		if result := mo.Get(slice); result.IsNone() || result.Some() != "slice_value" {
			t.Errorf("Expected to find slice key, got %v", result)
		}

		if result := mo.Get(m); result.IsNone() || result.Some() != "map_value" {
			t.Errorf("Expected to find map key, got %v", result)
		}

		keys := mo.Keys()
		if keys.Len() != 6 {
			t.Errorf("Expected 6 keys in iteration, got %d", keys.Len())
		}

		mo.Delete("string_key", fn, slice)
		if mo.Len() != 3 {
			t.Errorf("Expected length 3 after mixed delete, got %d", mo.Len())
		}

		if mo.Contains("string_key") || mo.Contains(fn) || mo.Contains(slice) {
			t.Errorf("Deleted keys should not be found")
		}

		if !mo.Contains(42) || !mo.Contains(3.14) || !mo.Contains(m) {
			t.Errorf("Non-deleted keys should still be found")
		}
	})
}
