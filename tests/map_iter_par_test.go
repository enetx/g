package g_test

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/enetx/g"
)

// FnPair increments in-flight count, updates max, sleeps, then decrements.
func (cc *concurrentCounter) FnPair(k, v int) {
	cur := atomic.AddInt64(&cc.inFlight, 1)

	for {
		prev := atomic.LoadInt64(&cc.maxInFlight)
		if cur <= prev || atomic.CompareAndSwapInt64(&cc.maxInFlight, prev, cur) {
			break
		}
	}

	time.Sleep(cc.sleep)
	atomic.AddInt64(&cc.inFlight, -1)
}

// assertMapContains checks that map values match expected.
func assertMapContains(t *testing.T, m Map[int, int], expected map[int]int) {
	t.Helper()

	for k, v := range expected {
		opt := m.Get(k)
		if !opt.IsSome() {
			t.Errorf("expected key %d to be present", k)
			continue
		}
		if opt.Some() != v {
			t.Errorf("for key %d, expected %d, got %d", k, v, opt.Some())
		}
	}
}

// TestCollectCountParallel tests Collect and Count with parallelism.
func TestCollectCountParallel(t *testing.T) {
	m := NewMap[int, int]()
	m.Insert(1, 10)
	m.Insert(2, 20)
	m.Insert(3, 30)

	workers := Int(3)

	cc := &concurrentCounter{sleep: 20 * time.Millisecond}
	col := m.Iter().Parallel(workers).
		Inspect(cc.FnPair).
		Collect()

	assertMapContains(t, col, map[int]int{1: 10, 2: 20, 3: 30})

	if cc.Max() < 2 {
		t.Errorf("expected parallel Collect, got max %d", cc.Max())
	}

	cc2 := &concurrentCounter{sleep: 20 * time.Millisecond}
	cnt := m.Iter().Parallel(workers).
		Inspect(cc2.FnPair).
		Count()

	if cnt.Std() != 3 {
		t.Errorf("Count: expected 3, got %d", cnt.Std())
	}

	if cc2.Max() < 2 {
		t.Errorf("expected parallel Count, got max %d", cc2.Max())
	}
}

// TestFilterMapParallel tests Filter and Map with parallelism.
func TestFilterMapParallel(t *testing.T) {
	m := NewMap[int, int]()
	m.Insert(1, 1)
	m.Insert(2, 2)
	m.Insert(3, 3)

	workers := Int(2)
	cc := &concurrentCounter{sleep: 15 * time.Millisecond}

	res := m.Iter().Parallel(workers).
		Inspect(cc.FnPair).
		Filter(func(_, v int) bool { return v%2 == 1 }).
		Map(func(k, v int) (int, int) { return k, v * v }).
		Collect()

	assertMapContains(t, res, map[int]int{1: 1, 3: 9})
	if cc.Max() < 2 {
		t.Errorf("expected parallel Filter+Map, got max %d", cc.Max())
	}
}

// TestTakeSkipParallel tests Take and Skip with parallelism.
func TestTakeSkipParallel(t *testing.T) {
	m := NewMap[int, int]()
	m.Insert(1, 100)
	m.Insert(2, 200)
	m.Insert(3, 300)

	workers := Int(2)

	cc1 := &concurrentCounter{sleep: 10 * time.Millisecond}
	cnt1 := m.Iter().Parallel(workers).
		Inspect(cc1.FnPair).
		Take(2).Count()

	if cnt1.Std() != 2 {
		t.Errorf("Take: expected 2, got %d", cnt1.Std())
	}

	if cc1.Max() < 2 {
		t.Errorf("expected parallel Take, got max %d", cc1.Max())
	}

	cc2 := &concurrentCounter{sleep: 10 * time.Millisecond}
	cnt2 := m.Iter().Parallel(workers).
		Inspect(cc2.FnPair).
		Skip(1).Count()

	if cnt2.Std() != 2 {
		t.Errorf("Skip: expected 2, got %d", cnt2.Std())
	}

	if cc2.Max() < 2 {
		t.Errorf("expected parallel Skip, got max %d", cc2.Max())
	}
}

// TestExcludeForEachParallel tests Exclude and ForEach with parallelism.
func TestExcludeForEachParallel(t *testing.T) {
	m := NewMap[int, int]()
	m.Insert(1, 1)
	m.Insert(2, 2)
	m.Insert(3, 3)
	m.Insert(4, 4)

	workers := Int(2)

	// Test Exclude function
	cc := &concurrentCounter{sleep: 10 * time.Millisecond}
	excluded := m.Iter().Parallel(workers).
		Inspect(cc.FnPair).
		Exclude(func(_, v int) bool { return v%2 == 0 }). // Exclude even values
		Collect()

	// Should only contain odd values
	assertMapContains(t, excluded, map[int]int{1: 1, 3: 3})
	if excluded.Len() != 2 {
		t.Errorf("Exclude: expected 2 elements, got %d", excluded.Len())
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Exclude, got max %d", cc.Max())
	}

	// Test ForEach function
	cc2 := &concurrentCounter{sleep: 10 * time.Millisecond}
	visitedCount := int64(0)
	m.Iter().Parallel(workers).
		Inspect(cc2.FnPair).
		ForEach(func(k, v int) {
			atomic.AddInt64(&visitedCount, 1)
		})

	if visitedCount != 4 {
		t.Errorf("ForEach: expected to visit 4 elements, got %d", visitedCount)
	}

	if cc2.Max() < 2 {
		t.Errorf("expected parallel ForEach, got max %d", cc2.Max())
	}
}

// TestChainAllAnyFindParallel tests Chain, All, Any, Find with parallelism.
func TestChainAllAnyFindParallel(t *testing.T) {
	m1 := NewMap[int, int]()
	m1.Insert(1, 1)
	m1.Insert(2, 2)
	m2 := NewMap[int, int]()
	m2.Insert(3, 3)

	workers := Int(2)

	cc := &concurrentCounter{sleep: 10 * time.Millisecond}
	chain := m1.Iter().Parallel(workers).
		Inspect(cc.FnPair).
		Chain(m2.Iter().Parallel(workers).
			Inspect(cc.FnPair))

	if chain.Count().Std() != 3 {
		t.Errorf("Chain Count: expected 3, got %d", chain.Count().Std())
	}

	if cc.Max() < 2 {
		t.Errorf("expected parallel Chain, got max %d", cc.Max())
	}

	ccAll := &concurrentCounter{sleep: 5 * time.Millisecond}
	all := chain.
		Inspect(ccAll.FnPair).
		All(func(_, v int) bool { return v > 0 })

	if !all {
		t.Error("All: expected true")
	}

	if ccAll.Max() < 2 {
		t.Errorf("expected parallel All, got max %d", ccAll.Max())
	}

	if !chain.Any(func(_, v int) bool { return v == 2 }) {
		t.Error("Any: expected true for v==2")
	}

	ccFind := &concurrentCounter{sleep: 5 * time.Millisecond}
	opt := chain.
		Inspect(ccFind.FnPair).
		Find(func(k, _ int) bool { return k == 3 })

	if !opt.IsSome() || opt.Some().Key != 3 {
		t.Error("Find: expected key 3")
	}

	if ccFind.Max() < 2 {
		t.Errorf("expected parallel Find, got max %d", ccFind.Max())
	}
}
