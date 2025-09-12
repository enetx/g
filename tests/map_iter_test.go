package g_test

import (
	"context"
	"testing"

	"github.com/enetx/g"
)

func TestMap_Iter_Keys(t *testing.T) {
	m := g.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	keys := m.Iter().Keys().Collect()

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	keySet := g.NewSet[string]()
	for _, key := range keys {
		keySet.Insert(key)
	}

	if !keySet.Contains("one") || !keySet.Contains("two") || !keySet.Contains("three") {
		t.Error("Keys iterator should contain all map keys")
	}
}

func TestMap_Iter_Values(t *testing.T) {
	m := g.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	values := m.Iter().Values().Collect()

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	valueSet := g.NewSet[int]()
	for _, value := range values {
		valueSet.Insert(value)
	}

	if !valueSet.Contains(1) || !valueSet.Contains(2) || !valueSet.Contains(3) {
		t.Error("Values iterator should contain all map values")
	}
}

func TestMap_Iter_Collect(t *testing.T) {
	original := g.NewMap[string, int]()
	original.Set("a", 1)
	original.Set("b", 2)

	collected := original.Iter().Collect()

	if collected.Len() != 2 {
		t.Errorf("Expected collected map to have 2 entries, got %d", collected.Len())
	}

	if collected.Get("a").UnwrapOr(0) != 1 {
		t.Error("Collected map should contain key 'a' with value 1")
	}

	if collected.Get("b").UnwrapOr(0) != 2 {
		t.Error("Collected map should contain key 'b' with value 2")
	}
}

func TestMap_Iter_Filter(t *testing.T) {
	m := g.NewMap[string, int]()
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

	if !filtered.Contains("two") || !filtered.Contains("four") {
		t.Error("Filtered map should contain only even-valued entries")
	}
}

func TestMap_Iter_Count(t *testing.T) {
	m := g.NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	count := m.Iter().Count()

	if count != 3 {
		t.Errorf("Expected count of 3, got %d", count)
	}
}

func TestMap_Iter_Take(t *testing.T) {
	m := g.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)
	m.Set("four", 4)

	taken := m.Iter().Take(2).Collect()

	if taken.Len() != 2 {
		t.Errorf("Expected 2 entries after Take(2), got %d", taken.Len())
	}
}

func TestMap_Iter_EmptyMap(t *testing.T) {
	m := g.NewMap[string, int]()

	count := m.Iter().Count()
	if count != 0 {
		t.Errorf("Empty map iterator count should be 0, got %d", count)
	}

	keys := m.Iter().Keys().Collect()
	if len(keys) != 0 {
		t.Error("Empty map should have no keys")
	}

	values := m.Iter().Values().Collect()
	if len(values) != 0 {
		t.Error("Empty map should have no values")
	}
}

func TestMap_Iter_Pull(t *testing.T) {
	m := g.NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	iter := m.Iter()
	next, stop := iter.Pull()
	defer stop()

	count := 0
	for {
		key, value, ok := next()
		if !ok {
			break
		}
		count++

		// Verify the key and value are valid
		if len(key) == 0 || value <= 0 {
			t.Errorf("Invalid key-value pair: %s -> %d", key, value)
		}
	}

	if count != 3 {
		t.Errorf("Expected to pull 3 pairs, got %d", count)
	}
}

func TestMapIterContext(t *testing.T) {
	t.Run("context cancellation stops iteration", func(t *testing.T) {
		m := g.NewMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)
		m.Set("three", 3)
		m.Set("four", 4)
		m.Set("five", 5)

		ctx, cancel := context.WithCancel(context.Background())

		var collected []g.Pair[string, int]
		iter := m.Iter().Context(ctx)

		// Cancel context after processing 2 elements
		count := 0
		iter(func(k string, v int) bool {
			collected = append(collected, g.Pair[string, int]{Key: k, Value: v})
			count++
			if count == 2 {
				cancel()
			}
			return true
		})

		// Should have processed exactly 2 elements before cancellation
		if len(collected) != 2 {
			t.Errorf("Expected 2 elements, got %d: %v", len(collected), collected)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		m := g.NewMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)

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

func TestMapIterNth(t *testing.T) {
	t.Run("nth element exists", func(t *testing.T) {
		m := g.NewMap[string, int]()
		m.Set("first", 1)
		m.Set("second", 2)
		m.Set("third", 3)
		m.Set("fourth", 4)
		m.Set("fifth", 5)

		// Get the 2nd pair (0-indexed)
		nth := m.Iter().Nth(2)

		if nth.IsNone() {
			t.Error("Expected Some value, got None")
		} else {
			pair := nth.Some()
			// Verify the pair is from the original map
			if value := m.Get(pair.Key); value.IsNone() || value.Some() != pair.Value {
				t.Errorf("Nth pair {%s: %d} is not in original map", pair.Key, pair.Value)
			}
		}
	})

	t.Run("nth element out of bounds", func(t *testing.T) {
		m := g.NewMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)

		nth := m.Iter().Nth(5)

		if nth.IsSome() {
			t.Errorf("Expected None for out of bounds index, got Some(%v)", nth.Some())
		}
	})

	t.Run("negative index", func(t *testing.T) {
		m := g.NewMap[string, int]()
		m.Set("one", 1)
		m.Set("two", 2)

		nth := m.Iter().Nth(-1)

		if nth.IsSome() {
			t.Errorf("Expected None for negative index, got Some(%v)", nth.Some())
		}
	})

	t.Run("empty map", func(t *testing.T) {
		m := g.NewMap[string, int]()

		nth := m.Iter().Nth(0)

		if nth.IsSome() {
			t.Errorf("Expected None for empty map, got Some(%v)", nth.Some())
		}
	})
}

func TestSeqMapFilterMap(t *testing.T) {
	// Test FilterMap with config validation
	configs := g.NewMap[string, string]()
	configs.Set("host", "localhost")
	configs.Set("port", "8080")
	configs.Set("debug", "invalid")
	configs.Set("timeout", "30")

	validConfigs := configs.Iter().FilterMap(func(k, v string) g.Option[g.Pair[string, string]] {
		// Keep only port and host configs with validation suffix
		if k == "port" || k == "host" {
			return g.Some(g.Pair[string, string]{Key: k, Value: v + "_validated"})
		}
		return g.None[g.Pair[string, string]]()
	}).Collect()

	if validConfigs.Len() != 2 {
		t.Errorf("Expected 2 valid configs, got %d", validConfigs.Len())
	}

	if !validConfigs.Contains("host") || !validConfigs.Contains("port") {
		keys := validConfigs.Keys()
		keySlice := make([]string, 0)
		keys.Iter().ForEach(func(k string) { keySlice = append(keySlice, k) })
		t.Errorf("Expected host and port keys, got %v", keySlice)
	}

	if validConfigs.Get("host").UnwrapOr("") != "localhost_validated" {
		t.Errorf("Expected 'localhost_validated', got %v", validConfigs.Get("host"))
	}

	if validConfigs.Get("port").UnwrapOr("") != "8080_validated" {
		t.Errorf("Expected '8080_validated', got %v", validConfigs.Get("port"))
	}

	// Test FilterMap with age filtering
	users := g.NewMap[string, int]()
	users.Set("alice", 25)
	users.Set("bob", 17)
	users.Set("charlie", 30)
	users.Set("diana", 16)

	adults := users.Iter().FilterMap(func(name string, age int) g.Option[g.Pair[string, int]] {
		if age >= 18 {
			return g.Some(g.Pair[string, int]{Key: name, Value: age})
		}
		return g.None[g.Pair[string, int]]()
	}).Collect()

	if adults.Len() != 2 {
		t.Errorf("Expected 2 adults, got %d", adults.Len())
	}

	if !adults.Contains("alice") || !adults.Contains("charlie") {
		keys := adults.Keys()
		keySlice := make([]string, 0)
		keys.Iter().ForEach(func(k string) { keySlice = append(keySlice, k) })
		t.Errorf("Expected alice and charlie, got %v", keySlice)
	}

	// Test FilterMap that filters all elements
	allFiltered := users.Iter().FilterMap(func(name string, age int) g.Option[g.Pair[string, int]] {
		return g.None[g.Pair[string, int]]() // Filter all out
	}).Collect()

	if allFiltered.Len() != 0 {
		t.Errorf("Expected empty map, got %d elements", allFiltered.Len())
	}

	// Test FilterMap that keeps all elements with transformation
	doubled := users.Iter().FilterMap(func(name string, age int) g.Option[g.Pair[string, int]] {
		return g.Some(g.Pair[string, int]{Key: name + "_user", Value: age * 2})
	}).Collect()

	if doubled.Len() != users.Len() {
		t.Errorf("Expected %d elements, got %d", users.Len(), doubled.Len())
	}

	if doubled.Get("alice_user").UnwrapOr(0) != 50 {
		t.Errorf("Expected alice_user age 50, got %v", doubled.Get("alice_user"))
	}

	// Test FilterMap with empty input
	emptyMap := g.NewMap[string, int]()
	emptyResult := emptyMap.Iter().FilterMap(func(name string, age int) g.Option[g.Pair[string, int]] {
		return g.Some(g.Pair[string, int]{Key: name, Value: age * 2})
	}).Collect()

	if emptyResult.Len() != 0 {
		t.Errorf("FilterMap on empty map should return empty, got %d elements", emptyResult.Len())
	}
}

func TestSeqMapNext(t *testing.T) {
	t.Run("Next with non-empty iterator", func(t *testing.T) {
		m := g.Map[string, int]{"x": 100}
		iter := m.Iter()

		// Extract first pair
		first := iter.Next()
		if !first.IsSome() {
			t.Errorf("Expected Some(Pair), got None")
		}

		firstPair := first.Some()
		if firstPair.Key != "x" || firstPair.Value != 100 {
			t.Errorf("Expected {x: 100}, got {%s: %d}", firstPair.Key, firstPair.Value)
		}

		// Second call should return None
		second := iter.Next()
		if second.IsSome() {
			t.Errorf("Expected None after exhausting single-element map, got Some(%v)", second.Some())
		}
	})

	t.Run("Next with empty iterator", func(t *testing.T) {
		m := g.NewMap[string, int]()
		iter := m.Iter()

		result := iter.Next()
		if result.IsSome() {
			t.Errorf("Expected None, got Some(%v)", result.Some())
		}
	})

	t.Run("Next until exhausted", func(t *testing.T) {
		m := g.Map[string, int]{"a": 1, "b": 2}
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

		// Iterator should be empty now
		remaining := iter.Collect()
		if remaining.Len() != 0 {
			t.Errorf("Expected empty map, got length %d", remaining.Len())
		}
	})
}
