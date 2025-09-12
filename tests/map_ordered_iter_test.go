package g_test

import (
	"context"
	"testing"

	"github.com/enetx/g"
)

func TestMapOrdered_Iter_Keys(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)

	keys := m.Iter().Keys().Collect()

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Check order is preserved
	expected := []string{"first", "second", "third"}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Key at position %d: expected %s, got %s", i, expected[i], key)
		}
	}
}

func TestMapOrdered_Iter_Values(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)

	values := m.Iter().Values().Collect()

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Check order is preserved
	expected := []int{1, 2, 3}
	for i, value := range values {
		if value != expected[i] {
			t.Errorf("Value at position %d: expected %d, got %d", i, expected[i], value)
		}
	}
}

func TestMapOrdered_Iter_Collect(t *testing.T) {
	original := g.NewMapOrd[string, int]()
	original.Set("a", 1)
	original.Set("b", 2)
	original.Set("c", 3)

	collected := original.Iter().Collect()

	if collected.Len() != 3 {
		t.Errorf("Expected collected map to have 3 entries, got %d", collected.Len())
	}

	// Check that order is preserved
	keys := collected.Keys()
	expectedKeys := []string{"a", "b", "c"}
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Collected map key order mismatch at position %d", i)
		}
	}
}

func TestMapOrdered_Iter_Filter(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)
	m.Set("four", 4)

	filtered := m.Iter().
		Filter(func(k string, v int) bool { return v%2 == 0 }).
		Collect()

	if filtered.Len() != 2 {
		t.Errorf("Expected 2 even values, got %d", filtered.Len())
	}

	// Check that original order is preserved for filtered items
	keys := filtered.Keys()
	if len(keys) != 2 || keys[0] != "two" || keys[1] != "four" {
		t.Error("Filtered ordered map should preserve original order")
	}
}

func TestMapOrdered_Iter_Count(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	count := m.Iter().Count()

	if count != 3 {
		t.Errorf("Expected count of 3, got %d", count)
	}
}

func TestMapOrdered_Iter_Take(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)
	m.Set("fourth", 4)

	taken := m.Iter().Take(2).Collect()

	if taken.Len() != 2 {
		t.Errorf("Expected 2 entries after Take(2), got %d", taken.Len())
	}

	// Check that first 2 entries are taken in order
	keys := taken.Keys()
	if len(keys) != 2 || keys[0] != "first" || keys[1] != "second" {
		t.Error("Take should preserve order and take first N entries")
	}
}

func TestMapOrdered_Iter_EmptyMap(t *testing.T) {
	m := g.NewMapOrd[string, int]()

	count := m.Iter().Count()
	if count != 0 {
		t.Errorf("Empty ordered map iterator count should be 0, got %d", count)
	}

	keys := m.Iter().Keys().Collect()
	if len(keys) != 0 {
		t.Error("Empty ordered map should have no keys")
	}

	values := m.Iter().Values().Collect()
	if len(values) != 0 {
		t.Error("Empty ordered map should have no values")
	}
}

func TestMapOrdered_Iter_Chain(t *testing.T) {
	m1 := g.NewMapOrd[string, int]()
	m1.Set("a", 1)
	m1.Set("b", 2)

	m2 := g.NewMapOrd[string, int]()
	m2.Set("c", 3)
	m2.Set("d", 4)

	chained := m1.Iter().Chain(m2.Iter()).Collect()

	if chained.Len() != 4 {
		t.Errorf("Expected 4 entries after chaining, got %d", chained.Len())
	}

	// Check that order is preserved across chain
	keys := chained.Keys()
	expected := []string{"a", "b", "c", "d"}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Chained map order mismatch at position %d", i)
		}
	}
}

func TestMapOrdered_Iter_Pull(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)

	iter := m.Iter()
	next, stop := iter.Pull()
	defer stop()

	count := 0
	expectedKeys := []string{"first", "second", "third"}
	expectedValues := []int{1, 2, 3}

	for {
		key, value, ok := next()
		if !ok {
			break
		}

		// Check order is preserved
		if count < len(expectedKeys) && key != expectedKeys[count] {
			t.Errorf("Expected key %s at position %d, got %s", expectedKeys[count], count, key)
		}
		if count < len(expectedValues) && value != expectedValues[count] {
			t.Errorf("Expected value %d at position %d, got %d", expectedValues[count], count, value)
		}

		count++
	}

	if count != 3 {
		t.Errorf("Expected to pull 3 pairs, got %d", count)
	}
}

func TestMapOrderedIterContext(t *testing.T) {
	t.Run("context cancellation stops iteration", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		m.Set("first", 1)
		m.Set("second", 2)
		m.Set("third", 3)
		m.Set("fourth", 4)
		m.Set("fifth", 5)

		ctx, cancel := context.WithCancel(context.Background())

		var collected []g.Pair[string, int]
		iter := m.Iter().Context(ctx)

		// Cancel context after processing 3 elements
		count := 0
		iter(func(k string, v int) bool {
			collected = append(collected, g.Pair[string, int]{Key: k, Value: v})
			count++
			if count == 3 {
				cancel()
			}
			return true
		})

		// Should have processed exactly 3 elements before cancellation
		if len(collected) != 3 {
			t.Errorf("Expected 3 elements, got %d: %v", len(collected), collected)
		}

		// Verify order is maintained
		expected := []g.Pair[string, int]{
			{Key: "first", Value: 1},
			{Key: "second", Value: 2},
			{Key: "third", Value: 3},
		}

		for i, pair := range collected {
			if pair.Key != expected[i].Key || pair.Value != expected[i].Value {
				t.Errorf("Order mismatch at index %d: got %v, want %v", i, pair, expected[i])
			}
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		m.Set("first", 1)
		m.Set("second", 2)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var collected []g.Pair[string, int]
		m.Iter().Context(ctx)(func(k string, v int) bool {
			collected = append(collected, g.Pair[string, int]{Key: k, Value: v})
			return true
		})

		// Should collect nothing due to immediate cancellation
		if len(collected) != 0 {
			t.Errorf("Expected 0 elements due to cancelled context, got %d: %v", len(collected), collected)
		}
	})
}

func TestMapOrderedIterNth(t *testing.T) {
	t.Run("nth element exists with order preservation", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		m.Set("first", 1)
		m.Set("second", 2)
		m.Set("third", 3)
		m.Set("fourth", 4)
		m.Set("fifth", 5)

		// Get the 2nd pair (0-indexed) - should be "third": 3
		nth := m.Iter().Nth(2)

		if nth.IsNone() {
			t.Error("Expected Some value, got None")
		} else {
			pair := nth.Some()
			if pair.Key != "third" || pair.Value != 3 {
				t.Errorf("Expected {third: 3}, got {%s: %d}", pair.Key, pair.Value)
			}
		}
	})

	t.Run("nth element maintains insertion order", func(t *testing.T) {
		m := g.NewMapOrd[int, string]()
		m.Set(10, "ten")
		m.Set(20, "twenty")
		m.Set(30, "thirty")
		m.Set(40, "forty")

		// Test each position
		testCases := []struct {
			index    g.Int
			expected g.Pair[int, string]
		}{
			{0, g.Pair[int, string]{Key: 10, Value: "ten"}},
			{1, g.Pair[int, string]{Key: 20, Value: "twenty"}},
			{2, g.Pair[int, string]{Key: 30, Value: "thirty"}},
			{3, g.Pair[int, string]{Key: 40, Value: "forty"}},
		}

		for _, tc := range testCases {
			nth := m.Iter().Nth(tc.index)
			if nth.IsNone() {
				t.Errorf("Expected Some value at index %d, got None", tc.index)
			} else {
				pair := nth.Some()
				if pair.Key != tc.expected.Key || pair.Value != tc.expected.Value {
					t.Errorf("At index %d: expected {%d: %s}, got {%d: %s}",
						tc.index, tc.expected.Key, tc.expected.Value, pair.Key, pair.Value)
				}
			}
		}
	})

	t.Run("nth element out of bounds", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)

		nth := m.Iter().Nth(5)

		if nth.IsSome() {
			t.Errorf("Expected None for out of bounds index, got Some(%v)", nth.Some())
		}
	})

	t.Run("negative index", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)

		nth := m.Iter().Nth(-1)

		if nth.IsSome() {
			t.Errorf("Expected None for negative index, got Some(%v)", nth.Some())
		}
	})

	t.Run("empty map", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()

		nth := m.Iter().Nth(0)

		if nth.IsSome() {
			t.Errorf("Expected None for empty map, got Some(%v)", nth.Some())
		}
	})
}

func TestSeqMapOrdNext(t *testing.T) {
	t.Run("Next with non-empty iterator maintains insertion order", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		m.Set("first", 1)
		m.Set("second", 2)
		m.Set("third", 3)
		iter := m.Iter()

		// First pair (should be "first" since it maintains insertion order)
		first := iter.Next()
		if !first.IsSome() {
			t.Errorf("Expected Some(Pair), got None")
		}

		firstPair := first.Some()
		if firstPair.Key != "first" || firstPair.Value != 1 {
			t.Errorf("Expected first pair {first: 1}, got {%s: %d}", firstPair.Key, firstPair.Value)
		}

		// Second pair (should be "second")
		second := iter.Next()
		if !second.IsSome() {
			t.Errorf("Expected Some(Pair), got None")
		}

		secondPair := second.Some()
		if secondPair.Key != "second" || secondPair.Value != 2 {
			t.Errorf("Expected second pair {second: 2}, got {%s: %d}", secondPair.Key, secondPair.Value)
		}

		// Remaining pairs
		remaining := iter.Collect()
		if remaining.Len() != 1 {
			t.Errorf("Expected 1 remaining pair, got %d", remaining.Len())
		}

		// Check the remaining pair
		remainingKeys := remaining.Keys()
		if len(remainingKeys) != 1 || remainingKeys[0] != "third" {
			t.Errorf("Expected remaining key 'third', got %v", remainingKeys)
		}
	})

	t.Run("Next with empty iterator", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		iter := m.Iter()

		result := iter.Next()
		if result.IsSome() {
			t.Errorf("Expected None, got Some(%v)", result.Some())
		}
	})

	t.Run("Next until exhausted", func(t *testing.T) {
		m := g.NewMapOrd[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		iter := m.Iter()

		// Extract all pairs
		first := iter.Next()
		second := iter.Next()
		third := iter.Next()

		if !first.IsSome() {
			t.Errorf("Expected first to be Some(Pair), got None")
		}
		if !second.IsSome() {
			t.Errorf("Expected second to be Some(Pair), got None")
		}
		if third.IsSome() {
			t.Errorf("Expected third to be None, got Some(%v)", third.Some())
		}

		// Verify order is maintained
		if first.Some().Key != "a" || first.Some().Value != 1 {
			t.Errorf("Expected first pair {a: 1}, got {%s: %d}", first.Some().Key, first.Some().Value)
		}
		if second.Some().Key != "b" || second.Some().Value != 2 {
			t.Errorf("Expected second pair {b: 2}, got {%s: %d}", second.Some().Key, second.Some().Value)
		}

		// Iterator should be empty now
		remaining := iter.Collect()
		if remaining.Len() != 0 {
			t.Errorf("Expected empty map, got length %d", remaining.Len())
		}
	})
}
