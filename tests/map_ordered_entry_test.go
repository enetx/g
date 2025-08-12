package g

import (
	"testing"

	. "github.com/enetx/g"
)

func TestMapOrdEntryOrSet(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Entry("foo").OrSet(10)
	opt := mo.Get("foo")
	if opt.IsNone() {
		t.Fatal("expected foo to be set, but Get returned None")
	}
	if val := opt.Unwrap(); val != 10 {
		t.Fatalf("expected foo=10, got %d", val)
	}
	mo.Entry("foo").OrSet(20)
	if got := mo.Get("foo").Unwrap(); got != 10 {
		t.Errorf("OrSet overwritten value: expected 10, got %d", got)
	}
}

func TestMapOrdEntryOrSetBy(t *testing.T) {
	mo := NewMapOrd[string, int]()
	called := false
	e := mo.Entry("bar")
	e.OrSetBy(func() int { called = true; return 7 })
	if !called {
		t.Error("expected OrSetBy fn to be called on vacant entry")
	}
	if mo.Get("bar").Unwrap() != 7 {
		t.Errorf("expected bar=7, got %d", mo.Get("bar").Unwrap())
	}
	called = false
	e.OrSetBy(func() int { called = true; return 9 })
	if called {
		t.Error("expected OrSetBy fn NOT to be called on occupied entry")
	}
}

func TestMapOrdEntryOrDefault(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Entry("baz").OrDefault()
	opt := mo.Get("baz")
	if opt.IsNone() {
		t.Fatal("expected baz to be set to default zero")
	}
	if val := opt.Unwrap(); val != 0 {
		t.Errorf("expected baz=0, got %d", val)
	}
}

func TestMapOrdEntryTransform(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Entry("x").Transform(func(int) int { return 5 })
	if mo.Get("x").IsSome() {
		t.Error("expected x to remain vacant after Transform on empty")
	}
	e := mo.Entry("y")
	e.OrSet(2)
	e.Transform(func(v int) int { return v * 3 })
	if got := mo.Get("y").Unwrap(); got != 6 {
		t.Errorf("expected y=6, got %d", got)
	}
}

func TestMapOrdEntrySetAndDelete(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Entry("a").Set(5)
	if mo.Get("a").Unwrap() != 5 {
		t.Errorf("expected a=5 after Set, got %d", mo.Get("a").Unwrap())
	}
	mo.Entry("a").Delete()
	if mo.Get("a").IsSome() {
		t.Error("expected a to be deleted, but Get returned Some")
	}
}

func TestMapOrdEntryGet(t *testing.T) {
	mo := NewMapOrd[string, int]()

	// Test Get on vacant entry
	entry := mo.Entry("vacant")
	value := entry.Get()
	if value.IsSome() {
		t.Errorf("expected Get on vacant entry to return None, got %v", value)
	}

	// Test Get on occupied entry
	mo.Set("occupied", 42)
	entryOccupied := mo.Entry("occupied")
	valueOccupied := entryOccupied.Get()
	if valueOccupied.IsNone() {
		t.Errorf("expected Get on occupied entry to return Some, got None")
	}
	if valueOccupied.Unwrap() != 42 {
		t.Errorf("expected Get to return 42, got %d", valueOccupied.Unwrap())
	}
}
