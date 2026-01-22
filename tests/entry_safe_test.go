package g_test

import (
	"sync"
	"testing"

	. "github.com/enetx/g"
)

func TestMapSafeEntryOrInsert(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// Insert new key - returns value
	val := ms.Entry("a").OrInsert(1)
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	// Verify value in map
	if v := ms.Get("a"); v.IsNone() || v.Some() != 1 {
		t.Errorf("expected 1 in map, got %v", v)
	}

	// Key already exists - returns existing value, not new
	val = ms.Entry("a").OrInsert(999)
	if val != 1 {
		t.Errorf("expected 1 (existing), got %d", val)
	}
}

func TestMapSafeEntryOrInsertWith(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// Insert new key - fn should be called
	called := false
	val := ms.Entry("key").OrInsertWith(func() int {
		called = true
		return 42
	})

	if !called {
		t.Error("expected fn to be called for new key")
	}
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}

	// Key exists - fn should NOT be called
	called = false
	val = ms.Entry("key").OrInsertWith(func() int {
		called = true
		return 999
	})

	if called {
		t.Error("fn should not be called for existing key")
	}
	if val != 42 {
		t.Errorf("expected 42 (existing), got %d", val)
	}
}

func TestMapSafeEntryOrInsertWithKey(t *testing.T) {
	ms := NewMapSafe[string, int]()

	val := ms.Entry("hello").OrInsertWithKey(func(k string) int {
		return len(k)
	})

	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}

	// Key exists - fn should NOT be called
	called := false
	val = ms.Entry("hello").OrInsertWithKey(func(string) int {
		called = true
		return 999
	})

	if called {
		t.Error("fn should not be called for existing key")
	}
	if val != 5 {
		t.Errorf("expected 5 (existing), got %d", val)
	}
}

func TestMapSafeEntryOrDefault(t *testing.T) {
	ms := NewMapSafe[string, int]()

	val := ms.Entry("counter").OrDefault()
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	// Verify in map
	if v := ms.Get("counter"); v.IsNone() || v.Some() != 0 {
		t.Errorf("expected 0 in map, got %v", v)
	}
}

func TestMapSafeEntryAndModify(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// AndModify on non-existent key - should not panic
	ms.Entry("missing").AndModify(func(v *int) { *v += 100 })

	// Verify key still doesn't exist
	if ms.Contains("missing") {
		t.Error("key should not exist after AndModify on vacant")
	}

	// AndModify on existing key
	ms.Insert("counter", 10)
	ms.Entry("counter").AndModify(func(v *int) { *v += 5 })

	if v := ms.Get("counter"); v.IsNone() || v.Some() != 15 {
		t.Errorf("expected 15, got %v", v)
	}
}

func TestMapSafeEntryAndModifyOrInsert(t *testing.T) {
	ms := NewMapSafe[string, Int]()

	ms.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if ms.Get("counter").Some() != 1 {
		t.Errorf("expected 1, got %d", ms.Get("counter").Some())
	}

	ms.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if ms.Get("counter").Some() != 2 {
		t.Errorf("expected 2, got %d", ms.Get("counter").Some())
	}

	ms.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if ms.Get("counter").Some() != 3 {
		t.Errorf("expected 3, got %d", ms.Get("counter").Some())
	}
}

func TestMapSafeEntryKey(t *testing.T) {
	ms := NewMapSafe[string, int]()

	e := ms.Entry("mykey")
	if e.Key() != "mykey" {
		t.Errorf("expected 'mykey', got '%s'", e.Key())
	}
}

func TestMapSafeEntryWordFrequency(t *testing.T) {
	words := SliceOf("apple", "banana", "apple", "cherry", "banana", "apple")
	freq := NewMapSafe[string, Int]()

	words.Iter().ForEach(func(word string) {
		freq.Entry(word).AndModify(func(v *Int) { *v++ }).OrInsert(1)
	})

	if freq.Get("apple").Some() != 3 {
		t.Errorf("expected apple:3, got %d", freq.Get("apple").Some())
	}
	if freq.Get("banana").Some() != 2 {
		t.Errorf("expected banana:2, got %d", freq.Get("banana").Some())
	}
	if freq.Get("cherry").Some() != 1 {
		t.Errorf("expected cherry:1, got %d", freq.Get("cherry").Some())
	}
}

func TestMapSafeEntryGrouping(t *testing.T) {
	groups := NewMapSafe[int, Slice[int]]()

	for i := range 10 {
		groups.Entry(i % 3).
			AndModify(func(s *Slice[int]) { *s = s.Append(i) }).
			OrInsertWith(func() Slice[int] { return Slice[int]{i} })
	}

	if !groups.Get(0).Some().Eq(Slice[int]{0, 3, 6, 9}) {
		t.Errorf("expected [0,3,6,9], got %v", groups.Get(0).Some())
	}
	if !groups.Get(1).Some().Eq(Slice[int]{1, 4, 7}) {
		t.Errorf("expected [1,4,7], got %v", groups.Get(1).Some())
	}
	if !groups.Get(2).Some().Eq(Slice[int]{2, 5, 8}) {
		t.Errorf("expected [2,5,8], got %v", groups.Get(2).Some())
	}
}

func TestMapSafeEntryZeroValue(t *testing.T) {
	ms := NewMapSafe[string, int]()

	val := ms.Entry("zero").OrInsert(0)
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	if !ms.Contains("zero") {
		t.Error("should contain key with zero value")
	}
}

func TestMapSafeEntryEmptyStringKey(t *testing.T) {
	ms := NewMapSafe[string, int]()

	val := ms.Entry("").OrInsert(42)
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}

	if v := ms.Get(""); v.IsNone() || v.Some() != 42 {
		t.Errorf("expected 42 for empty key, got %v", v)
	}
}

func TestMapSafeEntryPatternMatch(t *testing.T) {
	ms := NewMapSafe[string, int]()

	// Vacant case
	switch e := ms.Entry("key").(type) {
	case OccupiedSafeEntry[string, int]:
		t.Error("expected VacantSafeEntry")
	case VacantSafeEntry[string, int]:
		e.Insert(42)
	}

	// Occupied case
	switch e := ms.Entry("key").(type) {
	case OccupiedSafeEntry[string, int]:
		if e.Get() != 42 {
			t.Errorf("expected 42, got %d", e.Get())
		}
	case VacantSafeEntry[string, int]:
		t.Error("expected OccupiedSafeEntry")
	}
}

func TestMapSafeEntryOccupiedInsert(t *testing.T) {
	ms := NewMapSafe[string, int]()
	ms.Insert("key", 10)

	if e, ok := ms.Entry("key").(OccupiedSafeEntry[string, int]); ok {
		old := e.Insert(20)
		if old != 10 {
			t.Errorf("expected old value 10, got %d", old)
		}
		if e.Get() != 20 {
			t.Errorf("expected new value 20, got %d", e.Get())
		}
	} else {
		t.Error("expected OccupiedSafeEntry")
	}
}

func TestMapSafeEntryOccupiedRemove(t *testing.T) {
	ms := NewMapSafe[string, int]()
	ms.Insert("key", 42)

	if e, ok := ms.Entry("key").(OccupiedSafeEntry[string, int]); ok {
		removed := e.Remove()
		if removed != 42 {
			t.Errorf("expected removed value 42, got %d", removed)
		}
		if ms.Contains("key") {
			t.Error("key should be removed")
		}
	} else {
		t.Error("expected OccupiedSafeEntry")
	}
}

func TestMapSafeEntryVacantInsert(t *testing.T) {
	ms := NewMapSafe[string, int]()

	if e, ok := ms.Entry("key").(VacantSafeEntry[string, int]); ok {
		val := e.Insert(42)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
		if v := ms.Get("key"); v.IsNone() || v.Some() != 42 {
			t.Errorf("expected 42 in map, got %v", v)
		}
	} else {
		t.Error("expected VacantSafeEntry")
	}
}

func TestMapSafeEntryChained(t *testing.T) {
	ms := NewMapSafe[string, Int]()

	// Chained operations
	val := ms.Entry("value").
		AndModify(func(v *Int) { *v += 100 }). // no effect, vacant
		OrInsert(10)

	if val != 10 {
		t.Errorf("expected 10, got %d", val)
	}

	// Now modify existing
	ms.Entry("value").AndModify(func(v *Int) { *v *= 2 })

	if ms.Get("value").Some() != 20 {
		t.Errorf("expected 20, got %d", ms.Get("value").Some())
	}
}

// Concurrent tests

func TestMapSafeEntryConcurrentAndModify(t *testing.T) {
	ms := NewMapSafe[string, int]()
	ms.Insert("counter", 0)

	var wg sync.WaitGroup

	// 100 goroutines each increment 10 times
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				ms.Entry("counter").AndModify(func(v *int) { *v++ })
			}
		}()
	}

	wg.Wait()

	if ms.Get("counter").Some() != 1000 {
		t.Errorf("expected 1000, got %d", ms.Get("counter").Some())
	}
}

func TestMapSafeEntryConcurrentWordFrequency(t *testing.T) {
	freq := NewMapSafe[string, Int]()
	words := SliceOf("apple", "banana", "apple", "cherry", "banana", "apple")

	var wg sync.WaitGroup

	// Multiple goroutines process same words
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			words.Iter().ForEach(func(word string) {
				freq.Entry(word).AndModify(func(v *Int) { *v++ }).OrInsert(1)
			})
		}()
	}

	wg.Wait()

	// apple: 3 * 10 = 30
	if freq.Get("apple").Some() != 30 {
		t.Errorf("expected apple:30, got %d", freq.Get("apple").Some())
	}
	// banana: 2 * 10 = 20
	if freq.Get("banana").Some() != 20 {
		t.Errorf("expected banana:20, got %d", freq.Get("banana").Some())
	}
	// cherry: 1 * 10 = 10
	if freq.Get("cherry").Some() != 10 {
		t.Errorf("expected cherry:10, got %d", freq.Get("cherry").Some())
	}
}

func TestMapSafeEntryConcurrentTrySet(t *testing.T) {
	ms := NewMapSafe[string, int]()
	var wg sync.WaitGroup

	insertCount := 0
	var mu sync.Mutex

	// Multiple goroutines try to insert same key using TrySet
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ms.TryInsert("key", 42).IsNone() {
				mu.Lock()
				insertCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Only one should have inserted
	if insertCount != 1 {
		t.Errorf("expected exactly 1 insert, got %d", insertCount)
	}

	if ms.Get("key").Some() != 42 {
		t.Errorf("expected 42, got %d", ms.Get("key").Some())
	}
}

func TestMapSafeEntryConcurrentUnique(t *testing.T) {
	seen := NewMapSafe[int, Unit]()
	var wg sync.WaitGroup

	results := NewMapSafe[int, int]()

	// Simulate concurrent unique filter using TrySet
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 10 {
				if seen.TryInsert(i, Unit{}).IsNone() {
					// First time seeing this value
					results.Insert(i, 1)
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
