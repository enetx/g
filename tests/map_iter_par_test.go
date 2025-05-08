package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func assertMapContains(t *testing.T, m Map[int, int], expected map[int]int) {
	t.Helper()
	for k, v := range expected {
		opt := m.Get(k)
		if !opt.IsSome() {
			t.Errorf("expected key %d to be present", k)
			continue
		}
		if opt.Some() != v {
			t.Errorf("for key %d, expected value %d, got %d", k, v, opt.Some())
		}
	}
}

func TestSeqMapParCollectCount(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(1, 10)
	m.Set(2, 20)
	m.Set(3, 30)

	par := m.Iter().Parallel(2)
	col := par.Collect()
	expected := map[int]int{1: 10, 2: 20, 3: 30}
	assertMapContains(t, col, expected)
	if got := par.Count(); got != 3 {
		t.Errorf("Count: expected 3, got %d", got)
	}
}

func TestSeqMapParFilterMap(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(1, 1)
	m.Set(2, 2)
	m.Set(3, 3)

	par := m.Iter().Parallel(3).Filter(func(_, v int) bool { return v%2 == 1 })
	col := par.Collect()
	expected := map[int]int{1: 1, 3: 3}
	assertMapContains(t, col, expected)
	if got := par.Count(); got != 2 {
		t.Errorf("Filter Count: expected 2, got %d", got)
	}
}

func TestSeqMapParMapTransform(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(1, 2)
	m.Set(2, 3)

	par := m.Iter().Parallel(4).Map(func(k, v int) (int, int) { return k, v * v })
	col := par.Collect()
	expected := map[int]int{1: 4, 2: 9}
	assertMapContains(t, col, expected)
}

func TestSeqMapParTakeSkip(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(1, 100)
	m.Set(2, 200)
	m.Set(3, 300)

	takePar := m.Iter().Parallel(2).Take(2)
	if got := takePar.Count(); got != 2 {
		t.Errorf("Take Count: expected 2, got %d", got)
	}

	skipPar := m.Iter().Parallel(2).Skip(1)
	if got := skipPar.Count(); got != 2 {
		t.Errorf("Skip Count: expected 2, got %d", got)
	}
}

func TestSeqMapParChainAllAnyFind(t *testing.T) {
	m1 := NewMap[int, int]()
	m1.Set(1, 1)
	m1.Set(2, 2)

	m2 := NewMap[int, int]()
	m2.Set(3, 3)

	chain := m1.Iter().Parallel(2).Chain(m2.Iter().Parallel(2))
	if got := chain.Count(); got != 3 {
		t.Errorf("Chain Count: expected 3, got %d", got)
	}

	if !chain.All(func(_, v int) bool { return v > 0 }) {
		t.Error("All: expected all values > 0")
	}

	if !chain.Any(func(_, v int) bool { return v == 2 }) {
		t.Error("Any: expected to find value 2")
	}

	op := chain.Find(func(k, _ int) bool { return k == 3 })
	if !op.IsSome() {
		t.Error("Find: expected to find key 3")
	} else if pair := op.Some(); pair.Key != 3 || pair.Value != 3 {
		t.Errorf("Find: expected (3,3), got (%d,%d)", pair.Key, pair.Value)
	}
}
