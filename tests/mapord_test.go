package g_test

import (
	"reflect"
	"testing"

	"gitlab.com/x0xO/g"
)

func TestMapOrd_Range(t *testing.T) {
	// Test scenario: Function stops at a specific key-value pair
	t.Run("FunctionStopsAtSpecificPair", func(t *testing.T) {
		orderedMap := g.MapOrd[string, int]{
			{"a", 1},
			{"b", 2},
			{"c", 3},
		}
		expected := map[string]int{"a": 1, "b": 2}

		result := make(map[string]int)
		stopAtB := func(key string, val int) bool {
			result[key] = val
			return key != "b"
		}

		orderedMap.Range(stopAtB)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})

	// Test scenario: Function always returns true
	t.Run("FunctionAlwaysTrue", func(t *testing.T) {
		orderedMap := g.MapOrd[string, int]{
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

		orderedMap.Range(alwaysTrue)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})

	// Test scenario: Empty ordered map
	t.Run("EmptyMap", func(t *testing.T) {
		emptyMap := g.MapOrd[string, int]{}
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
