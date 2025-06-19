package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestEntryOrSet(t *testing.T) {
	m := NewMap[string, int]()
	m.Entry("foo").OrSet(42)
	opt := m.Get("foo")
	if opt.IsNone() {
		t.Fatal("expected foo to be set, but Get returned None")
	}
	if val := opt.Unwrap(); val != 42 {
		t.Fatalf("expected foo=42, got %d", val)
	}
	m.Entry("foo").OrSet(100)
	if got := m.Get("foo").Unwrap(); got != 42 {
		t.Errorf("OrSet overwritten value: expected 42, got %d", got)
	}
}

func TestEntryOrSetBy(t *testing.T) {
	m := NewMap[string, int]()
	called := false
	e := m.Entry("bar")
	e.OrSetBy(func() int { called = true; return 7 })
	if !called {
		t.Error("expected OrSetBy fn to be called on vacant entry")
	}
	if m.Get("bar").Unwrap() != 7 {
		t.Errorf("expected bar=7, got %d", m.Get("bar").Unwrap())
	}
	called = false
	e.OrSetBy(func() int { called = true; return 9 })
	if called {
		t.Error("expected OrSetBy fn NOT to be called on occupied entry")
	}
}

func TestEntryOrDefault(t *testing.T) {
	m := NewMap[string, int]()
	m.Entry("baz").OrDefault()
	if !m.Get("baz").IsSome() {
		t.Fatal("expected baz to be set to default zero")
	}
	if val := m.Get("baz").Unwrap(); val != 0 {
		t.Errorf("expected baz=0, got %d", val)
	}
}

func TestEntryTransform(t *testing.T) {
	m := NewMap[string, int]()
	m.Entry("x").Transform(func(int) int { return 5 })
	if m.Get("x").IsSome() {
		t.Error("expected x to remain vacant after Transform on empty")
	}
	e := m.Entry("y")
	e.OrSet(2)
	e.Transform(func(v int) int { return v * 3 })
	if got := m.Get("y").Unwrap(); got != 6 {
		t.Errorf("expected y=6, got %d", got)
	}
}

func TestEntrySetAndDelete(t *testing.T) {
	m := NewMap[string, int]()
	m.Entry("a").Set(10)
	if m.Get("a").Unwrap() != 10 {
		t.Errorf("expected a=10 after Set, got %d", m.Get("a").Unwrap())
	}
	m.Entry("a").Delete()
	if m.Get("a").IsSome() {
		t.Error("expected a to be deleted, but Get returned Some")
	}
}
