package g_test

import (
	"reflect"
	"testing"

	"gitlab.com/x0xO/g"
)

func TestMapOrdIterSortBy(t *testing.T) {
	// Sample data
	data := g.NewMapOrd[int, string]()
	data.
		Set(1, "d").
		Set(3, "b").
		Set(2, "c").
		Set(5, "e").
		Set(4, "a")

	// Expected result
	expected := g.NewMapOrd[int, string]()
	expected.
		Set(1, "d").
		Set(2, "c").
		Set(3, "b").
		Set(4, "a").
		Set(5, "e")

	sortedItems := data.Iter().SortBy(func(a, b g.Pair[int, string]) bool { return a.Key < b.Key }).Collect()

	// Check if the result matches the expected output
	if !reflect.DeepEqual(sortedItems, expected) {
		t.Errorf("Expected %v, got %v", expected, sortedItems)
	}

	expected = g.NewMapOrd[int, string]()
	expected.
		Set(4, "a").
		Set(3, "b").
		Set(2, "c").
		Set(1, "d").
		Set(5, "e")

	sortedItems = data.Iter().SortBy(func(a, b g.Pair[int, string]) bool { return a.Value < b.Value }).Collect()

	// Check if the result matches the expected output
	if !reflect.DeepEqual(sortedItems, expected) {
		t.Errorf("Expected %v, got %v", expected, sortedItems)
	}
}

func TestMapOrdIterStepBy(t *testing.T) {
	// Test case 1: StepBy with a step size of 2
	mapData := g.NewMapOrd[string, int]()
	mapData.
		Set("one", 1).
		Set("two", 2).
		Set("three", 3).
		Set("four", 4).
		Set("five", 5)

	expectedResult := g.NewMapOrd[string, int]()
	expectedResult.
		Set("one", 1).
		Set("three", 3).
		Set("five", 5)

	iter := mapData.Iter().StepBy(2)
	result := iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 2: StepBy with a step size of 3
	mapData = g.NewMapOrd[string, int]()
	mapData.
		Set("one", 1).
		Set("two", 2).
		Set("three", 3).
		Set("four", 4).
		Set("five", 5)

	expectedResult = g.NewMapOrd[string, int]()
	expectedResult.Set("one", 1).Set("four", 4)

	iter = mapData.Iter().StepBy(3)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 3: StepBy with a step size larger than the map length

	mapData = g.NewMapOrd[string, int]()
	mapData.Set("one", 1).
		Set("two", 2).
		Set("three", 3)

	expectedResult = g.NewMapOrd[string, int]()
	expectedResult.Set("one", 1)

	iter = mapData.Iter().StepBy(5)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}

	// Test case 4: StepBy with a step size of 1
	mapData = g.NewMapOrd[string, int]()
	mapData.Set("one", 1).
		Set("two", 2).
		Set("three", 3)

	expectedResult = g.NewMapOrd[string, int]()
	expectedResult.
		Set("one", 1).
		Set("two", 2).
		Set("three", 3)

	iter = mapData.Iter().StepBy(1)
	result = iter.Collect()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("StepBy failed. Expected %v, got %v", expectedResult, result)
	}
}

func TestMapOrdIterRange(t *testing.T) {
	// Test scenario: Function stops at a specific key-value pair
	t.Run("FunctionStopsAtSpecificPair", func(t *testing.T) {
		orderedMap := g.MapOrd[string, int]{
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

		orderedMap.Iter().Range(alwaysTrue)

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

		emptyMap.Iter().Range(anyFunc)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	})
}
