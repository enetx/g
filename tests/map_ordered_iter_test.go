package g_test

import (
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
