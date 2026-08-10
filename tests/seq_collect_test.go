package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

// anyPairsSeq returns a Seq2[any, any] yielding the given pairs in order.
func anyPairsSeq(pairs ...Pair[any, any]) Seq2[any, any] {
	return func(yield func(any, any) bool) {
		for _, p := range pairs {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func TestSeqCollectSetTypeMismatchPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Collect().Set[int]() should panic on a non-int element")
		}
	}()

	_ = SliceOf[any](1, "x").Iter().Collect().Set[int]()
}

func TestSeq2CollectMapTypeMismatchPanics(t *testing.T) {
	t.Run("key", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Collect().Map[string, int]() should panic on a non-string key")
			}
		}()

		_ = anyPairsSeq(Pair[any, any]{Key: 1, Value: 2}).Collect().Map[string, int]()
	})

	t.Run("value", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Collect().Map[string, int]() should panic on a non-int value")
			}
		}()

		_ = anyPairsSeq(Pair[any, any]{Key: "a", Value: "b"}).Collect().Map[string, int]()
	})
}

func TestSeq2CollectMapSafeTypeMismatchPanics(t *testing.T) {
	t.Run("key", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Collect().MapSafe[string, int]() should panic on a non-string key")
			}
		}()

		_ = anyPairsSeq(Pair[any, any]{Key: 1, Value: 2}).Collect().MapSafe[string, int]()
	})

	t.Run("value", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Collect().MapSafe[string, int]() should panic on a non-int value")
			}
		}()

		_ = anyPairsSeq(Pair[any, any]{Key: "a", Value: "b"}).Collect().MapSafe[string, int]()
	})
}

func TestSeq2CollectMapOrdTypeMismatchPanics(t *testing.T) {
	t.Run("key", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Collect().MapOrd[string, int]() should panic on a non-string key")
			}
		}()

		_ = anyPairsSeq(Pair[any, any]{Key: 1, Value: 2}).Collect().MapOrd[string, int]()
	})

	t.Run("value", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Collect().MapOrd[string, int]() should panic on a non-int value")
			}
		}()

		_ = anyPairsSeq(Pair[any, any]{Key: "a", Value: "b"}).Collect().MapOrd[string, int]()
	})

	t.Run("value on repeated-key update", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Collect().MapOrd[string, int]() should panic on a non-int value for a repeated key")
			}
		}()

		_ = anyPairsSeq(
			Pair[any, any]{Key: "a", Value: 1},
			Pair[any, any]{Key: "a", Value: "boom"},
		).Collect().MapOrd[string, int]()
	})
}

func TestSeqCollectSetAnyElements(t *testing.T) {
	set := SliceOf[any](1, 2, 2, 3).Iter().Collect().Set[int]()

	if set.Len() != 3 {
		t.Errorf("expected 3 elements, got %d: %v", set.Len(), set)
	}

	for _, v := range []int{1, 2, 3} {
		if !set.Contains(v) {
			t.Errorf("expected set to contain %d, got %v", v, set)
		}
	}
}

func TestSeq2CollectMapAnyPairs(t *testing.T) {
	m := anyPairsSeq(
		Pair[any, any]{Key: "a", Value: 1},
		Pair[any, any]{Key: "b", Value: 2},
		Pair[any, any]{Key: "a", Value: 3},
	).Collect().Map[string, int]()

	if m.Len() != 2 {
		t.Errorf("expected 2 entries, got %d: %v", m.Len(), m)
	}

	if v := m.Get("a"); v.IsNone() || v.Some() != 3 {
		t.Errorf("expected a=3 (later key overwrites), got %v", v)
	}

	if v := m.Get("b"); v.IsNone() || v.Some() != 2 {
		t.Errorf("expected b=2, got %v", v)
	}
}

func TestSeq2CollectMapSafeAnyPairs(t *testing.T) {
	ms := anyPairsSeq(
		Pair[any, any]{Key: "a", Value: 1},
		Pair[any, any]{Key: "b", Value: 2},
		Pair[any, any]{Key: "a", Value: 3},
	).Collect().MapSafe[string, int]()

	if ms.Len() != 2 {
		t.Errorf("expected 2 entries, got %d: %v", ms.Len(), ms)
	}

	if v := ms.Get("a"); v.IsNone() || v.Some() != 3 {
		t.Errorf("expected a=3 (later key overwrites), got %v", v)
	}

	if v := ms.Get("b"); v.IsNone() || v.Some() != 2 {
		t.Errorf("expected b=2, got %v", v)
	}
}

func TestSeq2CollectMapOrdAnyPairs(t *testing.T) {
	mo := anyPairsSeq(
		Pair[any, any]{Key: "a", Value: 1},
		Pair[any, any]{Key: "b", Value: 2},
		Pair[any, any]{Key: "a", Value: 3},
	).Collect().MapOrd[string, int]()

	if mo.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", mo.Len(), mo)
	}

	if mo[0].Key != "a" || mo[0].Value != 3 {
		t.Errorf("expected first entry a=3 (repeated key updates in place), got %v=%v", mo[0].Key, mo[0].Value)
	}

	if mo[1].Key != "b" || mo[1].Value != 2 {
		t.Errorf("expected second entry b=2, got %v=%v", mo[1].Key, mo[1].Value)
	}
}
