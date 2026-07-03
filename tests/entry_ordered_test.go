package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestMapOrdEntryOrInsert(t *testing.T) {
	mo := NewMapOrd[string, int]()

	// Insert new key - returns value
	val := mo.Entry("a").OrInsert(1)
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	// Verify value in map
	if v := mo.Get("a"); v.IsNone() || v.Some() != 1 {
		t.Errorf("expected 1 in map, got %v", v)
	}

	// Key already exists - returns existing value, not new
	val = mo.Entry("a").OrInsert(999)
	if val != 1 {
		t.Errorf("expected 1 (existing), got %d", val)
	}
}

func TestMapOrdEntryOrInsertWith(t *testing.T) {
	mo := NewMapOrd[string, int]()

	// Insert new key - fn should be called
	called := false
	val := mo.Entry("key").OrInsertWith(func() int {
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
	val = mo.Entry("key").OrInsertWith(func() int {
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

func TestMapOrdEntryOrInsertWithKey(t *testing.T) {
	mo := NewMapOrd[string, int]()

	val := mo.Entry("hello").OrInsertWithKey(func(k string) int {
		return len(k)
	})

	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}

	// Key exists - fn should NOT be called
	called := false
	val = mo.Entry("hello").OrInsertWithKey(func(string) int {
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

func TestMapOrdEntryOrDefault(t *testing.T) {
	mo := NewMapOrd[string, int]()

	val := mo.Entry("counter").OrDefault()
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	// Verify in map
	if v := mo.Get("counter"); v.IsNone() || v.Some() != 0 {
		t.Errorf("expected 0 in map, got %v", v)
	}
}

func TestMapOrdEntryAndModify(t *testing.T) {
	mo := NewMapOrd[string, int]()

	// AndModify on non-existent key - should not panic
	mo.Entry("missing").AndModify(func(v *int) { *v += 100 })

	// Verify key still doesn't exist
	if mo.Contains("missing") {
		t.Error("key should not exist after AndModify on vacant")
	}

	// AndModify on existing key
	mo.Insert("counter", 10)
	mo.Entry("counter").AndModify(func(v *int) { *v += 5 })

	if v := mo.Get("counter"); v.IsNone() || v.Some() != 15 {
		t.Errorf("expected 15, got %v", v)
	}
}

func TestMapOrdEntryAndModifyOrInsert(t *testing.T) {
	mo := NewMapOrd[string, Int]()

	mo.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if mo.Get("counter").Some() != 1 {
		t.Errorf("expected 1, got %d", mo.Get("counter").Some())
	}

	mo.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if mo.Get("counter").Some() != 2 {
		t.Errorf("expected 2, got %d", mo.Get("counter").Some())
	}

	mo.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	if mo.Get("counter").Some() != 3 {
		t.Errorf("expected 3, got %d", mo.Get("counter").Some())
	}
}

func TestMapOrdEntryKey(t *testing.T) {
	mo := NewMapOrd[string, int]()

	e := mo.Entry("mykey")
	if e.Key() != "mykey" {
		t.Errorf("expected 'mykey', got '%s'", e.Key())
	}
}

func TestMapOrdEntryPreservesOrder(t *testing.T) {
	mo := NewMapOrd[string, int]()

	mo.Entry("c").OrInsert(3)
	mo.Entry("a").OrInsert(1)
	mo.Entry("b").OrInsert(2)

	keys := mo.Keys()
	expected := Slice[string]{"c", "a", "b"}

	if !keys.Eq(expected) {
		t.Errorf("expected order %v, got %v", expected, keys)
	}
}

func TestMapOrdEntryWordFrequency(t *testing.T) {
	words := SliceOf("apple", "banana", "apple", "cherry", "banana", "apple")
	freq := NewMapOrd[string, Int]()

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

	// Check order preserved
	keys := freq.Keys()
	expectedKeys := Slice[string]{"apple", "banana", "cherry"}
	if !keys.Eq(expectedKeys) {
		t.Errorf("expected order %v, got %v", expectedKeys, keys)
	}
}

func TestMapOrdEntryGrouping(t *testing.T) {
	groups := NewMapOrd[int, Slice[int]]()

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

func TestMapOrdEntryModifyPreservesOrder(t *testing.T) {
	mo := NewMapOrd[string, int]()

	mo.Entry("a").OrInsert(1)
	mo.Entry("b").OrInsert(2)
	mo.Entry("c").OrInsert(3)

	mo.Entry("b").AndModify(func(v *int) { *v *= 100 })

	keys := mo.Keys()
	expected := Slice[string]{"a", "b", "c"}

	if !keys.Eq(expected) {
		t.Errorf("expected order %v after modify, got %v", expected, keys)
	}

	if mo.Get("b").Some() != 200 {
		t.Errorf("expected b:200, got %d", mo.Get("b").Some())
	}
}

func TestMapOrdEntryZeroValue(t *testing.T) {
	mo := NewMapOrd[string, int]()

	val := mo.Entry("zero").OrInsert(0)
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	if !mo.Contains("zero") {
		t.Error("should contain key with zero value")
	}
}

func TestMapOrdEntryEmptyStringKey(t *testing.T) {
	mo := NewMapOrd[string, int]()

	val := mo.Entry("").OrInsert(42)
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}

	if v := mo.Get(""); v.IsNone() || v.Some() != 42 {
		t.Errorf("expected 42 for empty key, got %v", v)
	}
}

func TestMapOrdEntryPatternMatch(t *testing.T) {
	mo := NewMapOrd[string, int]()

	// Vacant case
	switch e := mo.Entry("key").(type) {
	case OccupiedOrdEntry[string, int]:
		t.Error("expected VacantOrdEntry")
	case VacantOrdEntry[string, int]:
		e.Insert(42)
	}

	// Occupied case
	switch e := mo.Entry("key").(type) {
	case OccupiedOrdEntry[string, int]:
		if e.Get() != 42 {
			t.Errorf("expected 42, got %d", e.Get())
		}
	case VacantOrdEntry[string, int]:
		t.Error("expected OccupiedOrdEntry")
	}
}

func TestMapOrdEntryOccupiedInsert(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Insert("key", 10)

	if e, ok := mo.Entry("key").(OccupiedOrdEntry[string, int]); ok {
		old := e.Insert(20)
		if old != 10 {
			t.Errorf("expected old value 10, got %d", old)
		}
		if e.Get() != 20 {
			t.Errorf("expected new value 20, got %d", e.Get())
		}
	} else {
		t.Error("expected OccupiedOrdEntry")
	}
}

func TestMapOrdEntryOccupiedRemove(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Insert("key", 42)

	if e, ok := mo.Entry("key").(OccupiedOrdEntry[string, int]); ok {
		removed := e.Remove()
		if removed != 42 {
			t.Errorf("expected removed value 42, got %d", removed)
		}
		if mo.Contains("key") {
			t.Error("key should be removed")
		}
	} else {
		t.Error("expected OccupiedOrdEntry")
	}
}

func TestMapOrdEntryVacantInsert(t *testing.T) {
	mo := NewMapOrd[string, int]()

	if e, ok := mo.Entry("key").(VacantOrdEntry[string, int]); ok {
		val := e.Insert(42)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
		if v := mo.Get("key"); v.IsNone() || v.Some() != 42 {
			t.Errorf("expected 42 in map, got %v", v)
		}
	} else {
		t.Error("expected VacantOrdEntry")
	}
}

func TestMapOrdEntryChained(t *testing.T) {
	mo := NewMapOrd[string, Int]()

	// Chained operations
	val := mo.Entry("value").
		AndModify(func(v *Int) { *v += 100 }). // no effect, vacant
		OrInsert(10)

	if val != 10 {
		t.Errorf("expected 10, got %d", val)
	}

	// Now modify existing
	mo.Entry("value").AndModify(func(v *Int) { *v *= 2 })

	if mo.Get("value").Some() != 20 {
		t.Errorf("expected 20, got %d", mo.Get("value").Some())
	}
}

// NOTE: OccupiedOrdEntry captures the key's index once, when the entry is
// created by Entry(), and every operation uses that index directly. Entries
// do NOT survive structural mutations (Remove, SortBy, Shuffle) of the
// MapOrd — a fresh entry must be obtained after such mutations. The tests
// below previously asserted per-operation re-resolution of a stale entry;
// they were adapted for the cached-index semantics.
func TestMapOrdEntryCachedIndexAfterRemove(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Insert("a", 1)
	mo.Insert("b", 2)
	mo.Insert("c", 3)

	// Structurally mutate the map FIRST: removing "a" shifts "c" from
	// index 2 to index 1. An entry obtained afterwards captures the new index.
	mo.Remove("a")

	e, ok := mo.Entry("c").(OccupiedOrdEntry[string, int])
	if !ok {
		t.Fatal("expected OccupiedOrdEntry for c")
	}

	if got := e.Get(); got != 3 {
		t.Errorf("Get: expected 3, got %d", got)
	}

	old := e.Insert(30)
	if old != 3 {
		t.Errorf("Insert: expected old 3, got %d", old)
	}
	if v := mo.Get("c"); v.IsNone() || v.Some() != 30 {
		t.Errorf("expected c=30 after Insert, got %v", v)
	}
	if v := mo.Get("b"); v.IsNone() || v.Some() != 2 {
		t.Errorf("Insert must not clobber b; expected 2, got %v", v)
	}

	e.AndModify(func(v *int) { *v += 5 })
	if v := mo.Get("c"); v.IsNone() || v.Some() != 35 {
		t.Errorf("AndModify: expected c=35, got %v", v)
	}
}

func TestMapOrdEntryCachedIndexAfterShuffle(t *testing.T) {
	mo := NewMapOrd[int, int]()
	for i := 0; i < 50; i++ {
		mo.Insert(i, i*10)
	}

	// Shuffle reorders the backing slice; an entry obtained afterwards
	// captures the key's post-shuffle index and stays valid.
	mo.Shuffle()

	e, ok := mo.Entry(7).(OccupiedOrdEntry[int, int])
	if !ok {
		t.Fatal("expected OccupiedOrdEntry for key 7")
	}

	if got := e.Get(); got != 70 {
		t.Errorf("Get after Shuffle: expected 70, got %d", got)
	}

	e.Insert(700)
	if v := mo.Get(7); v.IsNone() || v.Some() != 700 {
		t.Errorf("Insert after Shuffle: expected 700, got %v", v)
	}
}

func TestMapOrdEntryRemoveThenReuse(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Insert("x", 1)

	e, ok := mo.Entry("x").(OccupiedOrdEntry[string, int])
	if !ok {
		t.Fatal("expected OccupiedOrdEntry for x")
	}

	// Remove via the entry itself; after this the entry must not be reused.
	if got := e.Remove(); got != 1 {
		t.Errorf("Remove: expected 1, got %d", got)
	}

	// The key is gone; a fresh Entry must be vacant.
	ve, ok := mo.Entry("x").(VacantOrdEntry[string, int])
	if !ok {
		t.Fatal("expected VacantOrdEntry after removal")
	}

	// AndModify on a vacant entry does not call fn.
	called := false
	ve.AndModify(func(v *int) { called = true; *v++ })
	if called {
		t.Error("AndModify must not invoke fn for a vacant entry")
	}

	// Insert re-inserts the key at the end.
	if got := ve.Insert(99); got != 99 {
		t.Errorf("Insert: expected 99, got %d", got)
	}
	if v := mo.Get("x"); v.IsNone() || v.Some() != 99 {
		t.Errorf("expected x=99 after re-insert, got %v", v)
	}
}
