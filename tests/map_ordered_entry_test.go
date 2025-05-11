package g

import (
	"testing"

	. "github.com/enetx/g"
)

func TestMapOrdEntryOrSet(t *testing.T) {
	mo := NewMapOrd[string, int]()
	e := mo.Entry("foo").OrSet(10)
	opt := e.Get()
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
	e := mo.Entry("bar").OrSetBy(func() int { called = true; return 7 })
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
	mo.Entry("x").Transform(func(v *int) { *v = 5 })
	if mo.Get("x").IsSome() {
		t.Error("expected x to remain vacant after AndModify on empty")
	}
	mo.Entry("y").OrSet(2).Transform(func(v *int) { *v *= 3 })
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

func TestMapOrdEntryChaining(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Entry("chain").OrSet(1).Transform(func(v *int) { *v++ }).Set(100).Delete()
	if mo.Get("chain").IsSome() {
		t.Error("expected chain to be deleted after chained operations")
	}
}
