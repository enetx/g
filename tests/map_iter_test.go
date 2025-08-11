package g_test

import (
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
