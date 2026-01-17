package g_test

import (
	"sync"
	"testing"

	. "github.com/enetx/g"
)

func TestMapSafeTrySetNewKey(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// Insert new key
	opt := ms.TrySet("a", 1)
	if opt.IsSome() {
		t.Errorf("expected None for new key, got Some(%v)", opt.Some())
	}

	// Verify value was inserted
	if v := ms.Get("a"); v.IsNone() || v.Some() != 1 {
		t.Errorf("expected value 1, got %v", v)
	}
}

func TestMapSafeTrySetExistingKey(t *testing.T) {
	ms := NewMapSafe[string, int]()
	ms.Set("a", 1)

	// Try to insert existing key
	opt := ms.TrySet("a", 999)
	if opt.IsNone() || opt.Some() != 1 {
		t.Errorf("expected Some(1), got %v", opt)
	}

	// Verify value was NOT replaced
	if v := ms.Get("a"); v.IsNone() || v.Some() != 1 {
		t.Errorf("expected value 1 (unchanged), got %v", v)
	}
}

func TestMapSafeTrySetVsSet(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// TrySet on new key
	ms.TrySet("a", 1)
	if v := ms.Get("a").Some(); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	// TrySet on existing key - should NOT replace
	ms.TrySet("a", 100)
	if v := ms.Get("a").Some(); v != 1 {
		t.Errorf("expected 1 (unchanged), got %d", v)
	}

	// Set on existing key - should replace
	ms.Set("a", 100)
	if v := ms.Get("a").Some(); v != 100 {
		t.Errorf("expected 100 (replaced), got %d", v)
	}
}

func TestMapSafeTrySetZeroValue(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// Insert zero value
	opt := ms.TrySet("zero", 0)
	if opt.IsSome() {
		t.Error("expected None for new key")
	}

	// Try to insert again - should return existing zero
	opt = ms.TrySet("zero", 999)
	if opt.IsNone() || opt.Some() != 0 {
		t.Errorf("expected Some(0), got %v", opt)
	}
}

func TestMapSafeTrySetConcurrent(t *testing.T) {
	ms := NewMapSafe[string, int]()
	var wg sync.WaitGroup

	insertCount := 0
	var mu sync.Mutex

	// Multiple goroutines try to insert same key
	for i := range 100 {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			if ms.TrySet("key", val).IsNone() {
				mu.Lock()
				insertCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Only one should have inserted
	if insertCount != 1 {
		t.Errorf("expected exactly 1 insert, got %d", insertCount)
	}
}

func TestMapSafeTrySetConcurrentUnique(t *testing.T) {
	seen := NewMapSafe[int, Unit]()
	var wg sync.WaitGroup

	results := NewMapSafe[int, int]()

	// Simulate concurrent unique filter
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 10 {
				if seen.TrySet(i, Unit{}).IsNone() {
					// First time seeing this value
					results.Set(i, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Each value 0-9 should be "first" exactly once
	for i := range 10 {
		if v := results.Get(i); v.IsNone() || v.Some() != 1 {
			t.Errorf("value %d was not seen exactly once", i)
		}
	}
}

func TestMapSafeTrySetMultipleKeys(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// Insert multiple keys
	ms.TrySet("a", 1)
	ms.TrySet("b", 2)
	ms.TrySet("c", 3)

	// All should be inserted
	if ms.Len() != 3 {
		t.Errorf("expected 3 keys, got %d", ms.Len())
	}

	// Try to insert existing keys
	ms.TrySet("a", 100)
	ms.TrySet("b", 200)
	ms.TrySet("c", 300)

	// Values should be unchanged
	if ms.Get("a").Some() != 1 || ms.Get("b").Some() != 2 || ms.Get("c").Some() != 3 {
		t.Error("values should not be replaced by TrySet")
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
