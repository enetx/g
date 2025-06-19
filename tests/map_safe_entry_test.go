package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestMapSafeEntryOrSet(t *testing.T) {
	m := NewMapSafe[string, int]()
	entry := m.Entry("a")

	got := entry.OrSet(10)
	if got.IsSome() {
		t.Errorf("expected None, got Some(%v)", got.Some())
	}

	got = entry.OrSet(20)
	if got.IsNone() || got.Some() != 10 {
		t.Errorf("expected Some(10), got %v", got)
	}
}

func TestMapSafeEntryOrSetBy(t *testing.T) {
	m := NewMapSafe[string, int]()
	called := false
	entry := m.Entry("b")

	got := entry.OrSetBy(func() int {
		called = true
		return 42
	})

	if got.IsSome() {
		t.Errorf("expected None, got %v", got)
	}

	if !called {
		t.Errorf("expected fn to be called")
	}

	called = false
	got = entry.OrSetBy(func() int {
		called = true
		return 99
	})

	if got.IsNone() || got.Some() != 42 {
		t.Errorf("expected Some(42), got %v", got)
	}

	if called {
		t.Errorf("expected fn not to be called")
	}
}

func TestMapSafeEntryTransform(t *testing.T) {
	m := NewMapSafe[string, int]()
	m.Set("c", 5)
	entry := m.Entry("c")

	got := entry.Transform(func(v int) int {
		return v * 2
	})
	if got.IsNone() || got.Some() != 10 {
		t.Errorf("expected Some(10), got %v", got)
	}
}

func TestMapSafeEntryAndDelete(t *testing.T) {
	m := NewMapSafe[string, int]()
	m.Set("d", 7)
	entry := m.Entry("d")

	got := entry.Delete()
	if got.IsNone() || got.Some() != 7 {
		t.Errorf("expected Some(7), got %v", got)
	}

	got = entry.Delete()
	if got.IsSome() {
		t.Errorf("expected None, got %v", got)
	}
}

func TestMapSafeEntryGet(t *testing.T) {
	m := NewMapSafe[string, int]()
	m.Set("hello", 42)

	entry := m.Entry("hello")
	got := entry.Get()
	if !got.IsSome() || got.Unwrap() != 42 {
		t.Errorf("expected Some(42), got %v", got)
	}

	missing := m.Entry("missing").Get()
	if missing.IsSome() {
		t.Errorf("expected None, got %v", missing)
	}
}

func TestMapSafeEntryOrDefault(t *testing.T) {
	m := NewMapSafe[string, int]()
	entry := m.Entry("x")
	entry.OrDefault()

	if got := entry.Get(); !got.IsSome() || got.Unwrap() != 0 {
		t.Errorf("expected default 0, got %v", got)
	}
}

func TestMapSafeEntrySet(t *testing.T) {
	m := NewMapSafe[string, string]()
	entry := m.Entry("a")
	if prev := entry.Set("one"); prev.IsSome() {
		t.Errorf("expected None, got %v", prev)
	}

	if v := m.Get("a"); !v.IsSome() || v.Unwrap() != "one" {
		t.Errorf("expected one, got %v", v)
	}

	if prev := entry.Set("two"); prev.IsNone() || prev.Unwrap() != "one" {
		t.Errorf("expected Some(one), got %v", prev)
	}

	if v := m.Get("a"); !v.IsSome() || v.Unwrap() != "two" {
		t.Errorf("expected two, got %v", v)
	}
}
