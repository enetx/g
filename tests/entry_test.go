package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestMapEntryOrInsert(t *testing.T) {
	m := Map[string, int]{}

	// Insert new key - returns value
	val := m.Entry("a").OrInsert(1)
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	// Verify value in map
	if v := m.Get("a"); v.IsNone() || v.Some() != 1 {
		t.Errorf("expected 1 in map, got %v", v)
	}

	// Key already exists - returns existing value, not new
	val = m.Entry("a").OrInsert(999)
	if val != 1 {
		t.Errorf("expected 1 (existing), got %d", val)
	}
}

func TestMapEntryOrInsertWith(t *testing.T) {
	m := Map[string, int]{}

	// Insert new key - fn should be called
	called := false
	val := m.Entry("key").OrInsertWith(func() int {
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
	val = m.Entry("key").OrInsertWith(func() int {
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

func TestMapEntryOrInsertWithKey(t *testing.T) {
	m := Map[string, int]{}

	// Insert new key using key in fn
	val := m.Entry("hello").OrInsertWithKey(func(k string) int {
		return len(k)
	})

	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}

	// Key exists - fn should NOT be called
	called := false
	val = m.Entry("hello").OrInsertWithKey(func(string) int {
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

func TestMapEntryOrDefault(t *testing.T) {
	m := Map[string, int]{}

	val := m.Entry("counter").OrDefault()
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	// Verify in map
	if v := m.Get("counter"); v.IsNone() || v.Some() != 0 {
		t.Errorf("expected 0 in map, got %v", v)
	}
}

func TestMapEntryAndModify(t *testing.T) {
	m := Map[string, int]{}

	// AndModify on non-existent key - should not panic, returns Entry
	e := m.Entry("missing").AndModify(func(v *int) { *v += 100 })
	_ = e // no value inserted yet

	// Verify key still doesn't exist
	if m.Contains("missing") {
		t.Error("key should not exist after AndModify on vacant")
	}

	// AndModify on existing key
	m.Set("counter", 10)
	m.Entry("counter").AndModify(func(v *int) { *v += 5 })

	if v := m.Get("counter"); v.IsNone() || v.Some() != 15 {
		t.Errorf("expected 15, got %v", v)
	}
}

func TestMapEntryAndModifyOrInsert(t *testing.T) {
	m := Map[string, Int]{}

	// Pattern: increment or init to 1
	m.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if m.Get("counter").Some() != 1 {
		t.Errorf("expected 1, got %d", m.Get("counter").Some())
	}

	m.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if m.Get("counter").Some() != 2 {
		t.Errorf("expected 2, got %d", m.Get("counter").Some())
	}

	m.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if m.Get("counter").Some() != 3 {
		t.Errorf("expected 3, got %d", m.Get("counter").Some())
	}
}

func TestMapEntryKey(t *testing.T) {
	m := Map[string, int]{}

	e := m.Entry("mykey")
	if e.Key() != "mykey" {
		t.Errorf("expected 'mykey', got '%s'", e.Key())
	}
}

func TestMapEntryWordFrequency(t *testing.T) {
	words := SliceOf("apple", "banana", "apple", "cherry", "banana", "apple")
	freq := Map[string, Int]{}

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

func TestMapEntryGrouping(t *testing.T) {
	groups := Map[int, Slice[int]]{}

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

func TestMapEntryZeroValue(t *testing.T) {
	m := Map[string, int]{}

	val := m.Entry("zero").OrInsert(0)
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	// Verify key exists with zero value
	if !m.Contains("zero") {
		t.Error("should contain key with zero value")
	}
}

func TestMapEntryEmptyStringKey(t *testing.T) {
	m := Map[string, int]{}

	val := m.Entry("").OrInsert(42)
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}

	// Verify empty string key works
	if v := m.Get(""); v.IsNone() || v.Some() != 42 {
		t.Errorf("expected 42 for empty key, got %v", v)
	}
}

func TestMapEntryPatternMatch(t *testing.T) {
	m := Map[string, int]{}

	// Vacant case
	switch e := m.Entry("key").(type) {
	case OccupiedEntry[string, int]:
		t.Error("expected VacantEntry")
	case VacantEntry[string, int]:
		e.Insert(42)
	}

	// Occupied case
	switch e := m.Entry("key").(type) {
	case OccupiedEntry[string, int]:
		if e.Get() != 42 {
			t.Errorf("expected 42, got %d", e.Get())
		}
	case VacantEntry[string, int]:
		t.Error("expected OccupiedEntry")
	}
}

func TestMapEntryOccupiedInsert(t *testing.T) {
	m := Map[string, int]{}
	m.Set("key", 10)

	if e, ok := m.Entry("key").(OccupiedEntry[string, int]); ok {
		old := e.Insert(20)
		if old != 10 {
			t.Errorf("expected old value 10, got %d", old)
		}
		if e.Get() != 20 {
			t.Errorf("expected new value 20, got %d", e.Get())
		}
	} else {
		t.Error("expected OccupiedEntry")
	}
}

func TestMapEntryOccupiedRemove(t *testing.T) {
	m := Map[string, int]{}
	m.Set("key", 42)

	if e, ok := m.Entry("key").(OccupiedEntry[string, int]); ok {
		removed := e.Remove()
		if removed != 42 {
			t.Errorf("expected removed value 42, got %d", removed)
		}
		if m.Contains("key") {
			t.Error("key should be removed")
		}
	} else {
		t.Error("expected OccupiedEntry")
	}
}

func TestMapEntryVacantInsert(t *testing.T) {
	m := Map[string, int]{}

	if e, ok := m.Entry("key").(VacantEntry[string, int]); ok {
		val := e.Insert(42)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
		if v := m.Get("key"); v.IsNone() || v.Some() != 42 {
			t.Errorf("expected 42 in map, got %v", v)
		}
	} else {
		t.Error("expected VacantEntry")
	}
}

func TestMapEntryChained(t *testing.T) {
	m := Map[string, Int]{}

	// Chained operations
	val := m.Entry("value").
		AndModify(func(v *Int) { *v += 100 }). // no effect, vacant
		OrInsert(10)

	if val != 10 {
		t.Errorf("expected 10, got %d", val)
	}

	// Now modify existing
	m.Entry("value").AndModify(func(v *Int) { *v *= 2 })

	if m.Get("value").Some() != 20 {
		t.Errorf("expected 20, got %d", m.Get("value").Some())
	}
}
