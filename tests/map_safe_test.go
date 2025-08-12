package g_test

import (
	"sync"
	"testing"

	. "github.com/enetx/g"
)

func TestMapSafeGetOrSetByNewKey(t *testing.T) {
	m := NewMapSafe[string, int]()

	called := false
	val := 0

	opt := m.Entry("key").OrSetBy(func() int {
		called = true
		val = 42
		return val
	})

	if called == false {
		t.Error("expected fn to be called for new key")
	}

	if opt.IsSome() {
		t.Errorf("expected None, got Some(%v)", opt.Some())
	}

	if v := m.Get("key"); !v.IsSome() || v.Some() != val {
		t.Errorf("expected value 42 in map, got %v", v)
	}
}

func TestMapSafeGetOrSetByExistingKey(t *testing.T) {
	m := NewMapSafe[string, int]()
	m.Set("key", 7)

	opt := m.Entry("key").OrSetBy(func() int {
		return 999
	})

	if !opt.IsSome() || opt.Some() != 7 {
		t.Errorf("expected Some(7), got %v", opt)
	}

	if v := m.Get("key"); !v.IsSome() || v.Some() != 7 {
		t.Errorf("expected value 7 in map, got %v", v)
	}
}

func TestMapSafeGetAndSet(t *testing.T) {
	ms := NewMapSafe[string, int]()

	prev := ms.Set("x", 10)

	if prev.IsSome() {
		t.Errorf("expected None, got Some(%v)", prev.Some())
	}

	if val := ms.Get("x"); val.IsNone() || val.Some() != 10 {
		t.Errorf("expected 10, got %v", val)
	}

	prev = ms.Set("x", 42)

	if prev.IsNone() || prev.Some() != 10 {
		t.Errorf("expected Some(10), got %v", prev)
	}

	if val := ms.Get("x"); val.IsNone() || val.Some() != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestMapSafeIntoIter(t *testing.T) {
	ms := NewMapSafe[string, int]()
	ms.Set("a", 1)
	ms.Set("b", 2)
	ms.Set("c", 3)

	if ms.Len() != 3 {
		t.Fatalf("expected Len = 3, got %d", ms.Len())
	}

	collected := make(map[string]int)

	ms.IntoIter().ForEach(func(k string, v int) {
		collected[k] = v
	})

	if ms.Len() != 0 {
		t.Fatalf("expected map to be empty after IntoIter, got Len = %d", ms.Len())
	}

	expected := map[string]int{"a": 1, "b": 2, "c": 3}
	if len(collected) != len(expected) {
		t.Fatalf("expected %d collected items, got %d", len(expected), len(collected))
	}

	for k, v := range expected {
		got, ok := collected[k]
		if !ok || got != v {
			t.Errorf("expected %q: %d, got %d (ok=%v)", k, v, got, ok)
		}
	}
}

func TestMapSafe(t *testing.T) {
	t.Run("TestNewMapSafe", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		if ms == nil {
			t.Fatal("NewMapSafe returned nil")
		}
		if ms.Len() != 0 {
			t.Fatal("Expected length 0, got", ms.Len())
		}
	})

	t.Run("TestSetAndGet", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		if v := ms.Get("key1"); v.IsNone() || v.Some() != 1 {
			t.Fatalf("Expected value 1, got %v", v)
		}
	})

	t.Run("TestDelete", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		ms.Delete("key1")
		if ms.Contains("key1") {
			t.Fatal("Key was not deleted")
		}
	})

	t.Run("TestContains", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		if !ms.Contains("key1") {
			t.Fatal("Key1 should exist")
		}
	})

	t.Run("TestKeysAndValues", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		ms.Set("key2", 2)
		keys := ms.Keys()
		values := ms.Values()
		if len(keys) != 2 || len(values) != 2 {
			t.Fatal("Expected 2 keys and values")
		}
	})

	t.Run("TestInvert", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		ms.Set("key2", 2)
		inverted := ms.Invert()
		if v := inverted.Get(1); v.IsNone() || v.Some() != "key1" {
			t.Fatal("Inversion failed for key1")
		}
	})

	t.Run("TestIter", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		ms.Set("key2", 2)
		ms.Set("key3", 2)

		ms.Iter().ForEach(func(k string, _ int) {
			if k == "key2" {
				ms.Set("key2", 44)
			}
			if k == "key1" {
				ms.Delete(k)
			}
		})

		if v := ms.Get("key2"); v.IsNone() || v.Some() != 44 {
			t.Fatal("Key 'key2' was not changed during iteration")
		}

		if v := ms.Get("key1"); v.IsSome() {
			t.Fatal("Key 'key1' was not deleted during iteration")
		}
	})

	t.Run("TestClone", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		cloned := ms.Clone()
		if !cloned.Contains("key1") {
			t.Fatal("Cloning failed")
		}
	})

	t.Run("TestClear", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		ms.Set("key1", 1)
		ms.Clear()
		if ms.Len() != 0 {
			t.Fatal("Clear failed")
		}
	})

	t.Run("TestCopy", func(t *testing.T) {
		src := NewMapSafe[string, int]()
		src.Set("key1", 1)
		src.Set("key2", 2)
		dest := NewMapSafe[string, int]()
		dest.Copy(src)
		if !dest.Contains("key1") || !dest.Contains("key2") {
			t.Fatal("Copy failed to copy all elements")
		}
	})

	t.Run("TestEq", func(t *testing.T) {
		ms1 := NewMapSafe[string, int]()
		ms1.Set("key1", 1)
		ms1.Set("key2", 2)
		ms2 := NewMapSafe[string, int]()
		ms2.Set("key1", 1)
		ms2.Set("key2", 2)
		if !ms1.Eq(ms2) {
			t.Fatal("Equality check failed for identical maps")
		}
	})

	t.Run("TestGetOrSet", func(t *testing.T) {
		ms := NewMapSafe[string, int]()
		value := ms.Entry("key1").OrSet(42)
		if value.IsSome() && value.Some() != 42 {
			t.Fatalf("Expected value 42, got %v", value)
		}

		value = ms.Entry("key1").OrSet(100)
		if value.IsSome() && value.Some() != 42 {
			t.Fatalf("Expected value 42 (existing), got %v", value)
		}
	})

	t.Run("TestString", func(t *testing.T) {
		ms := NewMapSafe[string, string]()
		ms.Set("key", "value")
		expected := "MapSafe{key:value}"
		if ms.String() != expected {
			t.Fatalf("Expected %v, got %v", expected, ms.String())
		}
	})

	t.Run("TestConcurrentAccess", func(t *testing.T) {
		ms := NewMapSafe[String, int]()
		wg := sync.WaitGroup{}
		for i := range 100 {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				ms.Set(Format("key{}", i), i)
			}(i)
		}
		wg.Wait()
		if ms.Len() != 100 {
			t.Fatal("Concurrent access failed")
		}
	})
}

func TestMapSafePrint(t *testing.T) {
	m := NewMapSafe[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)

	// Just test that Print() doesn't panic and returns the map
	result := m.Print()
	if !result.Eq(m) {
		t.Errorf("Print() should return the same map")
	}
}

func TestMapSafePrintln(t *testing.T) {
	m := NewMapSafe[string, int]()
	m.Set("x", 10)
	m.Set("y", 20)

	// Just test that Println() doesn't panic and returns the map
	result := m.Println()
	if !result.Eq(m) {
		t.Errorf("Println() should return the same map")
	}
}

func TestMapSafeNe(t *testing.T) {
	m1 := NewMapSafe[string, int]()
	m1.Set("a", 1)
	m1.Set("b", 2)

	m2 := NewMapSafe[string, int]()
	m2.Set("a", 1)
	m2.Set("b", 3) // Different value

	// Test Ne (not equal)
	if !m1.Ne(m2) {
		t.Errorf("m1 should not be equal to m2")
	}

	m3 := NewMapSafe[string, int]()
	m3.Set("a", 1)
	m3.Set("b", 2) // Same as m1

	if m1.Ne(m3) {
		t.Errorf("m1 should be equal to m3")
	}
}

func TestMapSafeNotEmpty(t *testing.T) {
	// Test non-empty map
	m := NewMapSafe[string, int]()
	m.Set("test", 42)

	if !m.NotEmpty() {
		t.Errorf("NotEmpty() should return true for non-empty map")
	}

	// Test empty map
	emptyMap := NewMapSafe[string, int]()
	if emptyMap.NotEmpty() {
		t.Errorf("NotEmpty() should return false for empty map")
	}
}

func TestMapSafeEmpty(t *testing.T) {
	// Test empty map
	m := NewMapSafe[string, int]()

	if !m.Empty() {
		t.Errorf("Empty() should return true for empty map")
	}

	// Test non-empty map
	m.Set("test", 42)
	if m.Empty() {
		t.Errorf("Empty() should return false for non-empty map")
	}
}
