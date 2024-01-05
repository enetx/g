package g_test

import (
	"reflect"
	"strings"
	"testing"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/ref"
)

func TestMapKeys(t *testing.T) {
	m := g.NewMap[string, int]()
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
	m := g.NewMap[string, int]()

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
	m := g.NewMap[string, int]()
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
	src := g.Map[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	dst := g.Map[string, int]{
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
	m := g.Map[string, string]{}
	m = m.Set("key", "value")

	if m["key"] != "value" {
		t.Error("Expected value to be 'value'")
	}
}

func TestMapDelete(t *testing.T) {
	m := g.Map[string, int]{"a": 1, "b": 2, "c": 3}

	m = m.Delete("a", "b")

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

func TestMapEqual(t *testing.T) {
	m := g.NewMap[string, string]()
	m.Set("key", "value")

	other := g.NewMap[string, string]()
	other = other.Set("key", "value")

	if !m.Eq(other) {
		t.Error("m and other should be equal")
	}

	other = other.Set("key", "other value")

	if m.Eq(other) {
		t.Error("m and other should not be equal")
	}
}

func TestMapToMap(t *testing.T) {
	m := g.NewMap[string, int]()
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
	m := g.Map[int, int]{}
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
	m := g.NewMap[int, string](3)
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	expected := g.NewMap[int, string](3)
	expected.Set(2, "one")
	expected.Set(4, "two")
	expected.Set(6, "three")

	mapped := m.Map(func(k int, v string) (int, string) { return k * 2, v })

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}

	expected = g.NewMap[int, string](3)
	expected.Set(1, "one_suffix")
	expected.Set(2, "two_suffix")
	expected.Set(3, "three_suffix")

	mapped = m.Map(func(k int, v string) (int, string) { return k, v + "_suffix" })

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}

	expected = g.NewMap[int, string](3)
	expected.Set(0, "")
	expected.Set(1, "one")
	expected.Set(3, "three")

	mapped = m.Map(func(k int, v string) (int, string) {
		if k == 2 {
			return 0, ""
		}
		return k, v
	})

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}
}

func TestMapFilter(t *testing.T) {
	m := g.NewMap[string, int](3)
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	expected := g.NewMap[string, int](1)
	expected.Set("two", 2)

	filtered := m.Filter(func(k string, v int) bool { return v%2 == 0 })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = g.NewMap[string, int](2)
	expected.Set("one", 1)
	expected.Set("three", 3)

	filtered = m.Filter(func(k string, v int) bool { return strings.Contains(k, "e") })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = g.NewMap[string, int](3)
	expected.Set("one", 1)
	expected.Set("two", 2)
	expected.Set("three", 3)

	filtered = m.Filter(func(k string, v int) bool { return true })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = g.NewMap[string, int](0)

	filtered = m.Filter(func(k string, v int) bool { return false })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}
}


func TestMapInvertValues(t *testing.T) {
	m := g.NewMap[int, string](0)
	inv := m.Invert()

	if inv.Len() != 0 {
		t.Errorf("Expected inverted map to have length 0, but got length %d", inv.Len())
	}

	m2 := g.NewMap[string, int](3)
	m2.Set("one", 1)
	m2.Set("two", 2)
	m2.Set("three", 3)

	inv2 := m2.Invert()

	if inv2.Len() != 3 {
		t.Errorf("Expected inverted map to have length 3, but got length %d", inv2.Len())
	}

	if inv2.Get(1) != "one" {
		t.Errorf("Expected inverted map to map 1 to 'one', but got %s", inv2.Get(1))
	}

	if inv2.Get(2) != "two" {
		t.Errorf("Expected inverted map to map 2 to 'two', but got %s", inv2.Get(2))
	}

	if inv2.Get(3) != "three" {
		t.Errorf("Expected inverted map to map 3 to 'three', but got %s", inv2.Get(3))
	}
}

func TestGetOrSet(t *testing.T) {
	// Create a new ordered Map called "m" with string keys and integer pointers as values
	m := g.NewMap[string, *int]()

	// Use GetOrSet to set the value for the key "root" to 3 if it doesn't exist
	m.GetOrSet("root", ref.Of(3))

	// Check if the value for the key "root" is equal to 3
	value := m.Get("root")
	if *value != 3 {
		t.Errorf("Expected value 3 for key 'root', but got %v", *value)
	}

	// Use GetOrSet to retrieve the value for the key "root" (which is 3), multiply it by 2
	*m.GetOrSet("root", ref.Of(10)) *= 2

	// Check if the value for the key "root" is equal to 6
	value = m.Get("root")
	if *value != 6 {
		t.Errorf("Expected value 6 for key 'root', but got %v", *value)
	}
}

func TestRandomMap(t *testing.T) {
	// Create a map for testing
	testMap := g.NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)
	testMap.Set("three", 3)

	// Randomly select a key-value pair
	randomResult := testMap.Random()

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
	value := randomResult.Get(key)
	originalValue := testMap.Get(key)
	if value != originalValue {
		t.Errorf("Randomly selected value does not match the original value")
	}
}

func TestRandomEmptyMap(t *testing.T) {
	// Create an empty map for testing
	testMap := g.NewMap[string, int]()

	// Attempt to randomly select a key-value pair
	randomResult := testMap.Random()

	// Check if the result is an empty map
	if randomResult.Len() != 0 {
		t.Errorf("Expected an empty map, but got length %d", randomResult.Len())
	}
}

func TestRandomSampleMap(t *testing.T) {
	// Create a map for testing
	testMap := g.NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)
	testMap.Set("three", 3)

	// Randomly select a sample of key-value pairs
	randomResult := testMap.RandomSample(2)

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
		value := randomResult.Get(key)
		originalValue := testMap.Get(key)
		if value != originalValue {
			t.Errorf("Randomly selected value for key '%s' does not match the original value", key)
		}
	}
}

func TestRandomSampleEmptyMap(t *testing.T) {
	// Create an empty map for testing
	testMap := g.NewMap[string, int]()

	// Attempt to randomly select a sample of key-value pairs
	randomResult := testMap.RandomSample(3)

	// Check if the result is an empty map
	if randomResult.Len() != 0 {
		t.Errorf("Expected an empty map, but got length %d", randomResult.Len())
	}
}

func TestRandomSampleFullMap(t *testing.T) {
	// Create a map for testing
	testMap := g.NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)

	// Randomly select a sample of key-value pairs
	randomResult := testMap.RandomSample(2)

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
	testMap := g.NewMap[string, int]()

	// Attempt to randomize a range
	subrangeMap := testMap.RandomRange(1, 4)

	// Check if the result is an empty map
	if subrangeMap.Len() != 0 {
		t.Errorf("Expected an empty map, but got length %d", subrangeMap.Len())
	}
}

func TestRandomRangeMapInvalidRange(t *testing.T) {
	// Create a map for testing
	testMap := g.NewMap[string, int]()
	testMap.Set("one", 1)
	testMap.Set("two", 2)

	// Test an invalid range
	subrangeMap := testMap.RandomRange(3, 5)

	// Check if the result is an empty map
	if subrangeMap.Len() != 0 {
		t.Errorf("Expected an empty map, but got length %d", subrangeMap.Len())
	}
}

func TestMapRange(t *testing.T) {
	// Test scenario: Function always returns true
	t.Run("FunctionAlwaysTrue", func(t *testing.T) {
		m := g.Map[string, int]{"a": 1, "b": 2, "c": 3}
		expected := map[string]int{"a": 1, "b": 2, "c": 3}

		result := make(map[string]int)
		alwaysTrue := func(key string, val int) bool {
			result[key] = val
			return true
		}

		m.Range(alwaysTrue)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})

	// Test scenario: Empty map
	t.Run("EmptyMap", func(t *testing.T) {
		emptyMap := g.Map[string, int]{}
		expected := make(map[string]int)

		result := make(map[string]int)
		anyFunc := func(key string, val int) bool {
			result[key] = val
			return true
		}

		emptyMap.Range(anyFunc)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})
}
