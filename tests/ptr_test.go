package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestInt_Ptr(t *testing.T) {
	i := Int(42)
	ptr := i.Ptr()
	if *ptr != i {
		t.Errorf("expected %v, got %v", i, *ptr)
	}
}

func TestFloat_Ptr(t *testing.T) {
	f := Float(3.14)
	ptr := f.Ptr()
	if *ptr != f {
		t.Errorf("expected %v, got %v", f, *ptr)
	}
}

func TestBytes_Ptr(t *testing.T) {
	bs := NewBytes("hello")
	ptr := bs.Ptr()
	if string(*ptr) != string(bs) {
		t.Errorf("expected %q, got %q", bs, *ptr)
	}
}

func TestString_Ptr(t *testing.T) {
	s := String("golang")
	ptr := s.Ptr()
	if *ptr != s {
		t.Errorf("expected %q, got %q", s, *ptr)
	}
}

func TestMap_Ptr(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	ptr := m.Ptr()
	if (*ptr)["a"] != 1 {
		t.Errorf("expected value at key 'a' to be 1, got %v", (*ptr)["a"])
	}
}

func TestMapOrd_Ptr(t *testing.T) {
	mo := NewMapOrd[string, int]()
	mo.Set("a", 1)
	ptr := mo.Ptr()
	if ptr.Get("a").Unwrap() != 1 {
		t.Errorf("expected value at key 'a' to be 1, got %v", ptr.Get("a"))
	}
}

func TestSet_Ptr(t *testing.T) {
	set := NewSet[string]()
	set.Insert("foo")
	ptr := set.Ptr()
	if !ptr.Contains("foo") {
		t.Errorf("expected set to contain 'foo'")
	}
}

func TestSlice_Ptr(t *testing.T) {
	sl := Slice[int]{1, 2, 3}
	ptr := sl.Ptr()
	if len(*ptr) != 3 || (*ptr)[0] != 1 || (*ptr)[2] != 3 {
		t.Errorf("expected slice to be [1 2 3], got %v", *ptr)
	}
}
