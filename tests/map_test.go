package g_test

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestMapFromStd(t *testing.T) {
	// Test case 1: Test conversion of an empty standard map
	emptyStdMap := map[string]int{}
	emptyGenericMap := Map[string, int](emptyStdMap)
	if len(emptyGenericMap) != 0 {
		t.Errorf("Test case 1 failed: Expected empty generic map, got %v", emptyGenericMap)
	}

	// Test case 2: Test conversion of a standard map with elements
	stdMap := map[string]int{"a": 1, "b": 2, "c": 3}
	genericMap := Map[string, int](stdMap)
	for k, v := range stdMap {
		if genericMap[k] != v {
			t.Errorf("Test case 2 failed: Value mismatch for key %s. Expected %d, got %d", k, v, genericMap[k])
		}
	}
}

func TestMapClear(t *testing.T) {
	// Test case 1: Clearing an empty map
	emptyMap := Map[string, int]{}
	emptyMap.Clear()
	if !emptyMap.Empty() {
		t.Errorf("Test case 1 failed: Cleared empty map should be empty")
	}

	// Test case 2: Clearing a non-empty map
	testMap := Map[string, int]{"a": 1, "b": 2, "c": 3}
	testMap.Clear()
	if !testMap.Empty() {
		t.Errorf("Test case 2 failed: Cleared test map should be empty")
	}
}

func TestMapEmpty(t *testing.T) {
	// Test case 1: Empty map
	emptyMap := Map[string, int]{}
	if !emptyMap.Empty() {
		t.Errorf("Test case 1 failed: Empty map should be empty")
	}

	// Test case 2: Non-empty map
	testMap := Map[string, int]{"a": 1, "b": 2, "c": 3}
	if testMap.Empty() {
		t.Errorf("Test case 2 failed: Non-empty map should not be empty")
	}
}

func TestMapString(t *testing.T) {
	// Test case 1: Empty map
	emptyMap := Map[string, int]{}
	expectedEmptyMapString := "Map{}"
	emptyMapString := emptyMap.String()
	if emptyMapString != expectedEmptyMapString {
		t.Errorf("Test case 1 failed: Expected %q, got %q", expectedEmptyMapString, emptyMapString)
	}
	// Test case 2: Map with elements
	testMap := Map[string, int]{"a": 1}
	expectedTestMapString := "Map{a:1}"
	testMapString := testMap.String()
	if testMapString != expectedTestMapString {
		t.Errorf("Test case 2 failed: Expected %q, got %q", expectedTestMapString, testMapString)
	}
}

func TestMapKeys(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	keys := m.Keys()
	if keys.Len() != 3 {
		t.Errorf("Expected 3 keys, got %d", keys.Len())
	}

	if !keys.Contains("a") {
		t.Errorf("Expected key 'a'")
	}

	if !keys.Contains("b") {
		t.Errorf("Expected key 'b'")
	}

	if !keys.Contains("c") {
		t.Errorf("Expected key 'c'")
	}
}

func TestMapValues(t *testing.T) {
	m := NewMap[string, int]()

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	values := m.Values()

	if values.Len() != 3 {
		t.Errorf("Expected 3 values, got %d", values.Len())
	}

	if !values.Contains(1) {
		t.Errorf("Expected value '1'")
	}

	if !values.Contains(2) {
		t.Errorf("Expected value '2'")
	}

	if !values.Contains(3) {
		t.Errorf("Expected value '3'")
	}
}

func TestMapClone(t *testing.T) {
	m := NewMap[string, int]()
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3

	nm := m.Clone()

	if m.Len() != nm.Len() {
		t.Errorf("Clone failed: expected %d, got %d", m.Len(), nm.Len())
	}

	for k, v := range m {
		if nm[k] != v {
			t.Errorf("Clone failed: expected %d, got %d", v, nm[k])
		}
	}
}

func TestMapCopy(t *testing.T) {
	src := Map[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	dst := Map[string, int]{
		"d": 4,
		"e": 5,
		"a": 6,
	}

	dst.Copy(src)

	if dst.Len() != 5 {
		t.Errorf("Expected len(dst) to be 5, got %d", len(dst))
	}

	if dst["a"] != 1 {
		t.Errorf("Expected dst[\"a\"] to be 1, got %d", dst["a"])
	}

	if dst["b"] != 2 {
		t.Errorf("Expected dst[\"b\"] to be 2, got %d", dst["b"])
	}

	if dst["c"] != 3 {
		t.Errorf("Expected dst[\"c\"] to be 3, got %d", dst["c"])
	}
}

func TestMapAdd(t *testing.T) {
	m := Map[string, string]{}
	m.Set("key", "value")

	if m["key"] != "value" {
		t.Error("Expected value to be 'value'")
	}
}

func TestMapDelete(t *testing.T) {
	m := Map[string, int]{"a": 1, "b": 2, "c": 3}

	m.Delete("a", "b")

	if m.Len() != 1 {
		t.Errorf("Expected length of 1, got %d", m.Len())
	}

	if _, ok := m["a"]; ok {
		t.Errorf("Expected key 'a' to be deleted")
	}

	if _, ok := m["b"]; ok {
		t.Errorf("Expected key 'b' to be deleted")
	}

	if _, ok := m["c"]; !ok {
		t.Errorf("Expected key 'c' to be present")
	}
}

func TestMapEq(t *testing.T) {
	// Test case 1: Equal maps
	map1 := Map[string, int]{"a": 1, "b": 2, "c": 3}
	map2 := Map[string, int]{"a": 1, "b": 2, "c": 3}
	if !map1.Eq(map2) {
		t.Errorf("Test case 1 failed: Equal maps should be considered equal")
	}

	// Test case 2: Maps with different lengths
	map3 := Map[string, int]{"a": 1, "b": 2}
	if map1.Eq(map3) {
		t.Errorf("Test case 2 failed: Maps with different lengths should not be considered equal")
	}

	// Test case 3: Maps with different values
	map4 := Map[string, int]{"a": 1, "b": 2, "c": 4}
	if map1.Eq(map4) {
		t.Errorf("Test case 3 failed: Maps with different values should not be considered equal")
	}

	// Test case 4
	map5 := Map[string, []int]{"a": []int{1}, "b": []int{2}, "c": []int{4}}
	map6 := Map[string, []int]{"a": []int{1}, "b": []int{2}, "c": []int{4}}
	if map5.Ne(map6) {
		t.Errorf("Test case 4 failed: Equal maps should be considered equal")
	}

	// Test case 5
	map7 := Map[string, []int]{"a": []int{2}, "b": []int{5}, "c": []int{4}}
	if map5.Eq(map7) {
		t.Errorf("Test case 5 failed: Maps with different values should not be considered equal")
	}

	// Test case 6
	if !NewMap[int, int]().Eq(NewMap[int, int]()) {
		t.Errorf("Test case 6 failed: Empty maps should be considered equal")
	}
}

func TestMapToMap(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	nmap := m.Std()

	if len(nmap) != 3 {
		t.Errorf("Expected 3, got %d", len(nmap))
	}

	if nmap["a"] != 1 {
		t.Errorf("Expected 1, got %d", nmap["a"])
	}

	if nmap["b"] != 2 {
		t.Errorf("Expected 2, got %d", nmap["b"])
	}

	if nmap["c"] != 3 {
		t.Errorf("Expected 3, got %d", nmap["c"])
	}
}

func TestMapLen(t *testing.T) {
	m := Map[int, int]{}
	if m.Len() != 0 {
		t.Errorf("Expected 0, got %d", m.Len())
	}

	m[1] = 1
	if m.Len() != 1 {
		t.Errorf("Expected 1, got %d", m.Len())
	}

	m[2] = 2
	if m.Len() != 2 {
		t.Errorf("Expected 2, got %d", m.Len())
	}
}

func TestMapMap(t *testing.T) {
	m := NewMap[int, string](3)
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	expected := NewMap[int, string](3)
	expected.Set(2, "one")
	expected.Set(4, "two")
	expected.Set(6, "three")

	mapped := m.Iter().Map(func(k int, v string) (int, string) { return k * 2, v }).Collect()

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}

	expected = NewMap[int, string](3)
	expected.Set(1, "one_suffix")
	expected.Set(2, "two_suffix")
	expected.Set(3, "three_suffix")

	mapped = m.Iter().Map(func(k int, v string) (int, string) { return k, v + "_suffix" }).Collect()

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}

	expected = NewMap[int, string](3)
	expected.Set(0, "")
	expected.Set(1, "one")
	expected.Set(3, "three")

	mapped = m.Iter().Map(func(k int, v string) (int, string) {
		if k == 2 {
			return 0, ""
		}
		return k, v
	}).Collect()

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}
}

func TestMapFilter(t *testing.T) {
	m := NewMap[string, int](3)
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	expected := NewMap[string, int](1)
	expected.Set("two", 2)

	filtered := m.Iter().Filter(func(k string, v int) bool { return v%2 == 0 }).Collect()

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = NewMap[string, int](2)
	expected.Set("one", 1)
	expected.Set("three", 3)

	filtered = m.Iter().Filter(func(k string, v int) bool { return strings.Contains(k, "e") }).Collect()

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = NewMap[string, int](3)
	expected.Set("one", 1)
	expected.Set("two", 2)
	expected.Set("three", 3)

	filtered = m.Iter().Filter(func(k string, v int) bool { return true }).Collect()

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = NewMap[string, int](0)

	filtered = m.Iter().Filter(func(k string, v int) bool { return false }).Collect()

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}
}

func TestMapInvertValues(t *testing.T) {
	m := NewMap[int, string](0)
	inv := m.Invert()

	if inv.Len() != 0 {
		t.Errorf("Expected inverted map to have length 0, but got length %d", inv.Len())
	}

	m2 := NewMap[string, int](3)
	m2.Set("one", 1)
	m2.Set("two", 2)
	m2.Set("three", 3)

	inv2 := m2.Invert()

	if inv2.Len() != 3 {
		t.Errorf("Expected inverted map to have length 3, but got length %d", inv2.Len())
	}

	if inv2.Get(1).Some() != "one" {
		t.Errorf("Expected inverted map to map 1 to 'one', but got %s", inv2.Get(1).Some())
	}

	if inv2.Get(2).Some() != "two" {
		t.Errorf("Expected inverted map to map 2 to 'two', but got %s", inv2.Get(2).Some())
	}

	if inv2.Get(3).Some() != "three" {
		t.Errorf("Expected inverted map to map 3 to 'three', but got %s", inv2.Get(3).Some())
	}
}

func TestMapGet(t *testing.T) {
	// Test case 1: Get existing key
	map1 := Map[string, int]{"a": 1, "b": 2, "c": 3}
	key1 := "b"
	expectedValue1 := 2
	value1 := map1.Get(key1)
	if value1.Some() != expectedValue1 {
		t.Errorf("Test case 1 failed: Expected value for key '%s' is %v, got %v", key1, expectedValue1, value1)
	}

	// Test case 2: Get non-existing key
	key2 := "d"
	value2 := map1.Get(key2)
	if value2.IsSome() {
		t.Errorf("Test case 2 failed, got %v", value2.Some())
	}
}

func TestRandomMap(t *testing.T) {
	// Create a map for testing
	testMap := NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)
	testMap.Set("three", 3)

	// Randomly select a key-value pair
	randomResult := testMap.Iter().Take(1).Collect()

	// Check if the result is a map with a single key-value pair
	if randomResult.Len() != 1 {
		t.Errorf("Expected a map with a single key-value pair, but got length %d", randomResult.Len())
	}

	// Check if the selected key exists in the original map
	key := randomResult.Keys()[0]
	if !testMap.Contains(key) {
		t.Errorf("Randomly selected key not found in the original map")
	}

	// Check if the selected value matches the original map
	value := randomResult.Get(key).Some()
	originalValue := testMap.Get(key).Some()
	if value != originalValue {
		t.Errorf("Randomly selected value does not match the original value")
	}
}

func TestRandomEmptyMap(t *testing.T) {
	// Create an empty map for testing
	testMap := NewMap[string, int]()

	// Attempt to randomly select a key-value pair
	randomResult := testMap.Iter().Take(1).Collect()

	// Check if the result is an empty map
	if randomResult.Len() != 0 {
		t.Errorf("Expected an empty map, but got length %d", randomResult.Len())
	}
}

func TestRandomSampleMap(t *testing.T) {
	// Create a map for testing
	testMap := NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)
	testMap.Set("three", 3)

	// Randomly select a sample of key-value pairs
	randomResult := testMap.Iter().Take(2).Collect()

	// Check if the result is a map with the specified number of key-value pairs
	if randomResult.Len() != 2 {
		t.Errorf("Expected a map with 2 key-value pairs, but got length %d", randomResult.Len())
	}

	// Check if all selected keys exist in the original map
	keys := randomResult.Keys()
	for _, key := range keys {
		if !testMap.Contains(key) {
			t.Errorf("Randomly selected key '%s' not found in the original map", key)
		}
	}

	// Check if the selected values match the original map
	for _, key := range keys {
		value := randomResult.Get(key).Some()
		originalValue := testMap.Get(key).Some()
		if value != originalValue {
			t.Errorf("Randomly selected value for key '%s' does not match the original value", key)
		}
	}
}

func TestRandomSampleEmptyMap(t *testing.T) {
	// Create an empty map for testing
	testMap := NewMap[string, int]()

	// Attempt to randomly select a sample of key-value pairs
	randomResult := testMap.Iter().Take(3).Collect()

	// Check if the result is an empty map
	if randomResult.Len() != 0 {
		t.Errorf("Expected an empty map, but got length %d", randomResult.Len())
	}
}

func TestRandomSampleFullMap(t *testing.T) {
	// Create a map for testing
	testMap := NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)

	// Randomly select a sample of key-value pairs
	randomResult := testMap.Iter().Take(2).Collect()

	// Check if the result is the same as the original map
	if randomResult.Len() != 2 {
		t.Errorf("Expected a map with 2 key-value pairs, but got length %d", randomResult.Len())
	}

	keys := randomResult.Keys()
	for _, key := range keys {
		if !testMap.Contains(key) {
			t.Errorf("Randomly selected key '%s' not found in the original map", key)
		}
	}
}

func TestRandomRangeMapEmpty(t *testing.T) {
	// Create an empty map for testing
	testMap := NewMap[string, int]()

	// Attempt to randomize a range
	subrangeMap := testMap.Iter().Take(Int(3).RandomRange(5).UInt()).Collect()

	// Check if the result is an empty map
	if subrangeMap.Len() != 0 {
		t.Errorf("Expected an empty map, but got length %d", subrangeMap.Len())
	}
}

func TestRandomRangeMapInvalidRange(t *testing.T) {
	// Create a map for testing
	testMap := NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)

	// Test an invalid range

	subrangeMap := testMap.Iter().Take(Int(3).RandomRange(5).UInt()).Collect()

	// Check if the result is the same as the original map
	if subrangeMap.Len() != 2 {
		t.Errorf("Expected a map with 2 key-value pairs, but got length %d", subrangeMap.Len())
	}
}

func TestMapNe(t *testing.T) {
	// Test case 1: Maps are equal
	map1 := Map[string, int]{"a": 1, "b": 2}
	map2 := Map[string, int]{"a": 1, "b": 2}
	expectedResult1 := false
	result1 := map1.Ne(map2)
	if result1 != expectedResult1 {
		t.Errorf("Test case 1 failed: Expected result is %t, got %t", expectedResult1, result1)
	}

	// Test case 2: Maps are not equal
	map3 := Map[string, int]{"a": 1, "b": 2}
	map4 := Map[string, int]{"a": 1, "b": 3}
	expectedResult2 := true
	result2 := map3.Ne(map4)
	if result2 != expectedResult2 {
		t.Errorf("Test case 2 failed: Expected result is %t, got %t", expectedResult2, result2)
	}
}

func TestMapNotEmpty(t *testing.T) {
	// Test case 1: Map is not empty
	map1 := Map[string, int]{"a": 1, "b": 2}
	expectedResult1 := true
	result1 := map1.NotEmpty()
	if result1 != expectedResult1 {
		t.Errorf("Test case 1 failed: Expected result is %t, got %t", expectedResult1, result1)
	}

	// Test case 2: Map is empty
	map2 := Map[string, int]{}
	expectedResult2 := false
	result2 := map2.NotEmpty()
	if result2 != expectedResult2 {
		t.Errorf("Test case 2 failed: Expected result is %t, got %t", expectedResult2, result2)
	}
}

func TestMapIterChain(t *testing.T) {
	// Test case 1: Concatenate two iterators
	iter1 := NewMap[int, string]()
	iter1.Set(1, "a")

	iter2 := NewMap[int, string]()
	iter2.Set(2, "b")

	concatenated := iter1.Iter().Chain(iter2.Iter()).Collect()

	expected := NewMap[int, string]()
	expected.Set(1, "a")
	expected.Set(2, "b")

	if !reflect.DeepEqual(concatenated, expected) {
		t.Errorf("Expected concatenated map to be %v, got %v", expected, concatenated)
	}

	// Test case 2: Concatenate three iterators
	iter3 := NewMap[int, string]()
	iter3.Set(3, "c")

	concatenated2 := iter1.Iter().Chain(iter2.Iter(), iter3.Iter()).Collect()

	expected2 := NewMap[int, string]()
	expected2.Set(1, "a")
	expected2.Set(2, "b")
	expected2.Set(3, "c")
	if !reflect.DeepEqual(concatenated2, expected2) {
		t.Errorf("Expected concatenated map to be %v, got %v", expected2, concatenated2)
	}
}

func TestMapIterCount(t *testing.T) {
	// Test case 1: Count elements in a non-empty map
	iter := Map[int, string]{1: "a", 2: "b"}.Iter()

	count := iter.Count()

	expected := Int(2)
	if count != expected {
		t.Errorf("Expected count to be %d, got %d", expected, count)
	}

	// Test case 2: Count elements in an empty map
	emptyIter := NewMap[int, string]().Iter()

	emptyCount := emptyIter.Count()

	emptyExpected := Int(0)
	if emptyCount != emptyExpected {
		t.Errorf("Expected count to be %d, got %d", emptyExpected, emptyCount)
	}
}

func TestMapIterExclude(t *testing.T) {
	// Test case 1: Exclude even values
	m := NewMap[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(5, "e")

	notEven := m.Iter().
		Exclude(func(k int, v string) bool {
			return k%2 == 0
		}).
		Collect()

	expected := NewMap[int, string]()
	expected.Set(1, "a")
	expected.Set(3, "c")
	expected.Set(5, "e")

	if !notEven.Eq(expected) {
		t.Errorf("Excluded result incorrect, expected: %v, got: %v", expected, notEven)
	}

	// Test case 2: Exclude all elements
	empty := m.Iter().
		Exclude(func(k int, v string) bool {
			return true
		}).
		Collect()

	if !empty.Empty() {
		t.Errorf("Expected empty map after exclusion, got: %v", empty)
	}
}

func TestMapIterFind(t *testing.T) {
	// Test case 1: Find an existing element
	m := NewMap[int, string]()
	m.Set(1, "a")
	m.Set(2, "b")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(5, "e")

	found := m.Iter().
		Find(func(k int, v string) bool {
			return k == 3
		})

	if found.IsNone() {
		t.Errorf("Expected to find key-value pair, got None")
	} else {
		expected := Pair[int, string]{Key: 3, Value: "c"}
		if found.Some() != expected {
			t.Errorf("Found key-value pair incorrect, expected: %v, got: %v", expected, found.Some())
		}
	}

	// Test case 2: Find a non-existing element
	notFound := m.Iter().
		Find(func(k int, v string) bool {
			return k == 6
		})

	if notFound.IsSome() {
		t.Errorf("Expected not to find key-value pair, got: %v", notFound.Some())
	}
}

func TestMapIterRange(t *testing.T) {
	// Define a map to iterate over
	m := NewMap[int, string]()
	m.Set(1, "apple")
	m.Set(2, "banana")
	m.Set(3, "cherry")

	// Define a slice to collect the keys visited during iteration
	var keysVisited Slice[int]

	// Iterate over the map using Range
	m.Iter().Range(func(k int, v string) bool {
		keysVisited = append(keysVisited, k)
		// Continue iterating until all elements are visited
		return true
	})

	keysVisited.SortBy(func(a, b int) cmp.Ordering { return cmp.Cmp(a, b) })

	// Check if all keys were visited
	expectedKeys := Slice[int]{1, 2, 3}

	if !reflect.DeepEqual(keysVisited, expectedKeys) {
		t.Errorf("Expected keys visited to be %v, got %v", expectedKeys, keysVisited)
	}
}

func TestMapIterInspect(t *testing.T) {
	// Define a map to iterate over
	m := NewMap[int, string]()
	m.Set(1, "apple")
	m.Set(2, "banana")
	m.Set(3, "cherry")

	// Define a slice to store the inspected key-value pairs
	inspected := NewMap[int, string]()

	// Create a new iterator with Inspect and collect the pairs
	m.Iter().Inspect(func(k int, v string) {
		inspected.Set(k, v)
	}).Collect()

	if !inspected.Eq(m) {
		t.Errorf("Expected inspected map to be %v, got %v", m, inspected)
	}
}

func TestMapTransformMap(t *testing.T) {
	original := Map[string, int]{"a": 1, "b": 2}

	addEntry := func(m Map[string, int]) Map[string, int] {
		m["c"] = 3
		return m
	}

	expected := Map[string, int]{"a": 1, "b": 2, "c": 3}
	result := original.Transform(addEntry)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Transform failed: expected %v, got %v", expected, result)
	}

	removeEntry := func(m Map[string, int]) Map[string, int] {
		delete(m, "a")
		return m
	}

	expectedAfterRemoval := Map[string, int]{"b": 2, "c": 3}
	resultAfterRemoval := result.Transform(removeEntry)

	if !reflect.DeepEqual(resultAfterRemoval, expectedAfterRemoval) {
		t.Errorf("Transform with removal failed: expected %v, got %v", expectedAfterRemoval, resultAfterRemoval)
	}
}

// go test -bench=. -benchmem -count=4

func genM() Map[String, int] {
	mo := NewMap[String, int](10000)
	for i := range 10000 {
		mo.Set(Int(i).String(), i)
	}

	return mo
}

func BenchmarkMEq(b *testing.B) {
	m := genM()
	m2 := m.Clone()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = m.Eq(m2)
	}
}

func TestMapToMapSafe(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)

	safemap := m.ToMapSafe()

	if safemap.Len() != 2 {
		t.Errorf("ToMapSafe() should preserve length, expected 2, got %d", safemap.Len())
	}

	if val1 := safemap.Get("key1"); val1.IsNone() || val1.Unwrap() != 1 {
		t.Errorf("ToMapSafe() should preserve key1 value, got %v", val1)
	}

	if val2 := safemap.Get("key2"); val2.IsNone() || val2.Unwrap() != 2 {
		t.Errorf("ToMapSafe() should preserve key2 value, got %v", val2)
	}
}

func TestMapPrint(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("key", 42)
	result := m.Print()

	if result.Len() != m.Len() {
		t.Errorf("Print() should return original map unchanged")
	}
}

func TestMapPrintln(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("key", 42)
	result := m.Println()

	if result.Len() != m.Len() {
		t.Errorf("Println() should return original map unchanged")
	}
}

func TestMapEntryGet(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("existing", 42)

	// Test getting existing entry
	entry := m.Entry("existing")
	value := entry.Get()
	if value.IsNone() || value.Unwrap() != 42 {
		t.Errorf("Entry.Get() for existing key should return Some(42), got %v", value)
	}

	// Test getting non-existing entry
	nonEntry := m.Entry("nonexistent")
	value2 := nonEntry.Get()
	if value2.IsSome() {
		t.Errorf("Entry.Get() for non-existing key should return None, got %v", value2)
	}
}

func TestMapEntrySet(t *testing.T) {
	m := NewMap[string, int]()
	entry := m.Entry("newkey")

	result := entry.Set(123)
	if result.IsSome() {
		t.Errorf("Entry.Set() for new key should return None (no previous value), got %v", result)
	}

	// Verify the value was actually set
	if val := m.Get("newkey"); val.IsNone() || val.Unwrap() != 123 {
		t.Errorf("Entry.Set() should set the value in map, got %v", val)
	}
}

func TestMapEntryDelete(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("todelete", 99)

	entry := m.Entry("todelete")
	deleted := entry.Delete()

	if deleted.IsNone() || deleted.Unwrap() != 99 {
		t.Errorf("Entry.Delete() should return Some(deleted_value), got %v", deleted)
	}

	if val := m.Get("todelete"); val.IsSome() {
		t.Errorf("Key should be deleted from map")
	}

	// Test deleting non-existing key
	nonEntry := m.Entry("nonexistent")
	deleted2 := nonEntry.Delete()
	if deleted2.IsSome() {
		t.Errorf("Entry.Delete() should return None when deleting non-existing key, got %v", deleted2)
	}
}
