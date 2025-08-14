package g_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/enetx/g"
)

func TestMapOrdered_Iter_Keys(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)

	keys := m.Iter().Keys().Collect()

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Check order is preserved
	expected := []string{"first", "second", "third"}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Key at position %d: expected %s, got %s", i, expected[i], key)
		}
	}
}

func TestMapOrdered_Iter_Values(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)

	values := m.Iter().Values().Collect()

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Check order is preserved
	expected := []int{1, 2, 3}
	for i, value := range values {
		if value != expected[i] {
			t.Errorf("Value at position %d: expected %d, got %d", i, expected[i], value)
		}
	}
}

func TestMapOrdered_Iter_Collect(t *testing.T) {
	original := g.NewMapOrd[string, int]()
	original.Set("a", 1)
	original.Set("b", 2)
	original.Set("c", 3)

	collected := original.Iter().Collect()

	if collected.Len() != 3 {
		t.Errorf("Expected collected map to have 3 entries, got %d", collected.Len())
	}

	// Check that order is preserved
	keys := collected.Keys()
	expectedKeys := []string{"a", "b", "c"}
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Collected map key order mismatch at position %d", i)
		}
	}
}

func TestMapOrdered_Iter_Filter(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)
	m.Set("four", 4)

	filtered := m.Iter().
		Filter(func(k string, v int) bool { return v%2 == 0 }).
		Collect()

	if filtered.Len() != 2 {
		t.Errorf("Expected 2 even values, got %d", filtered.Len())
	}

	// Check that original order is preserved for filtered items
	keys := filtered.Keys()
	if len(keys) != 2 || keys[0] != "two" || keys[1] != "four" {
		t.Error("Filtered ordered map should preserve original order")
	}
}

func TestMapOrdered_Iter_Count(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	count := m.Iter().Count()

	if count != 3 {
		t.Errorf("Expected count of 3, got %d", count)
	}
}

func TestMapOrdered_Iter_Take(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)
	m.Set("fourth", 4)

	taken := m.Iter().Take(2).Collect()

	if taken.Len() != 2 {
		t.Errorf("Expected 2 entries after Take(2), got %d", taken.Len())
	}

	// Check that first 2 entries are taken in order
	keys := taken.Keys()
	if len(keys) != 2 || keys[0] != "first" || keys[1] != "second" {
		t.Error("Take should preserve order and take first N entries")
	}
}

func TestMapOrdered_Iter_EmptyMap(t *testing.T) {
	m := g.NewMapOrd[string, int]()

	count := m.Iter().Count()
	if count != 0 {
		t.Errorf("Empty ordered map iterator count should be 0, got %d", count)
	}

	keys := m.Iter().Keys().Collect()
	if len(keys) != 0 {
		t.Error("Empty ordered map should have no keys")
	}

	values := m.Iter().Values().Collect()
	if len(values) != 0 {
		t.Error("Empty ordered map should have no values")
	}
}

func TestMapOrdered_Iter_Chain(t *testing.T) {
	m1 := g.NewMapOrd[string, int]()
	m1.Set("a", 1)
	m1.Set("b", 2)

	m2 := g.NewMapOrd[string, int]()
	m2.Set("c", 3)
	m2.Set("d", 4)

	chained := m1.Iter().Chain(m2.Iter()).Collect()

	if chained.Len() != 4 {
		t.Errorf("Expected 4 entries after chaining, got %d", chained.Len())
	}

	// Check that order is preserved across chain
	keys := chained.Keys()
	expected := []string{"a", "b", "c", "d"}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Chained map order mismatch at position %d", i)
		}
	}
}

func TestMapOrdered_Iter_Pull(t *testing.T) {
	m := g.NewMapOrd[string, int]()
	m.Set("first", 1)
	m.Set("second", 2)
	m.Set("third", 3)

	iter := m.Iter()
	next, stop := iter.Pull()
	defer stop()

	count := 0
	expectedKeys := []string{"first", "second", "third"}
	expectedValues := []int{1, 2, 3}

	for {
		key, value, ok := next()
		if !ok {
			break
		}

		// Check order is preserved
		if count < len(expectedKeys) && key != expectedKeys[count] {
			t.Errorf("Expected key %s at position %d, got %s", expectedKeys[count], count, key)
		}
		if count < len(expectedValues) && value != expectedValues[count] {
			t.Errorf("Expected value %d at position %d, got %d", expectedValues[count], count, value)
		}

		count++
	}

	if count != 3 {
		t.Errorf("Expected to pull 3 pairs, got %d", count)
	}
}

// Enhanced map counter for tracking key-value processing
type mapCounter struct {
	current       int64
	max           int64
	total         int64
	sleepDuration time.Duration
	sequenceID    string
}

func (mc *mapCounter) Fn(k g.String, v int) {
	current := atomic.AddInt64(&mc.current, 1)
	atomic.AddInt64(&mc.total, 1)

	// Update max concurrency
	for {
		currentMax := atomic.LoadInt64(&mc.max)
		if current <= currentMax || atomic.CompareAndSwapInt64(&mc.max, currentMax, current) {
			break
		}
	}

	// Simulate work
	time.Sleep(mc.sleepDuration)
	atomic.AddInt64(&mc.current, -1)
}

func (mc *mapCounter) Max() int64 {
	return atomic.LoadInt64(&mc.max)
}

func (mc *mapCounter) Total() int64 {
	return atomic.LoadInt64(&mc.total)
}

// TestSeqMapParChainComprehensive tests Chain with multiple aspects for maps
func TestSeqMapParChainComprehensive(t *testing.T) {
	t.Run("BasicMapParallelism", func(t *testing.T) {
		// Test basic parallel execution for maps
		m1 := g.NewMap[g.String, int]()
		m1.Set("a", 1)
		m1.Set("b", 2)
		m1.Set("c", 3)
		m1.Set("d", 4)
		m1.Set("e", 5)

		m2 := g.NewMap[g.String, int]()
		m2.Set("x", 10)
		m2.Set("y", 20)
		m2.Set("z", 30)
		m2.Set("w", 40)
		m2.Set("v", 50)

		workers := g.Int(3)
		mc := &mapCounter{sleepDuration: 50 * time.Millisecond}

		start := time.Now()
		res := m1.
			Iter().
			Parallel(workers).
			Inspect(mc.Fn).
			Chain(
				m2.Iter().Parallel(workers).Inspect(mc.Fn),
			).
			Collect()
		duration := time.Since(start)

		// Verify results
		if res.Len() != 10 {
			t.Errorf("expected 10 entries, got %d", res.Len())
		}

		// Verify all keys are present
		keys := res.Keys()
		expectedKeys := []string{"a", "b", "c", "d", "e", "x", "y", "z", "w", "v"}
		if len(keys) != len(expectedKeys) {
			t.Errorf("key count mismatch: got %d, want %d", len(keys), len(expectedKeys))
		}

		// Verify parallelism
		if mc.Max() < 2 {
			t.Errorf("expected parallel execution, got max concurrency %d", mc.Max())
		}

		// Verify timing - should be faster than sequential
		expectedSequentialTime := time.Duration(10) * 50 * time.Millisecond
		if duration > expectedSequentialTime/2 {
			t.Errorf("execution too slow, might not be parallel: %v", duration)
		}

		t.Logf("Map Basic - Max concurrency: %d, Total processed: %d, Duration: %v",
			mc.Max(), mc.Total(), duration)
	})

	t.Run("MultipleMapSequencesChain", func(t *testing.T) {
		// Test chaining multiple map sequences
		m1 := g.NewMap[g.String, int]()
		m1.Set("a", 1)
		m1.Set("b", 2)

		m2 := g.NewMap[g.String, int]()
		m2.Set("x", 10)
		m2.Set("y", 20)

		m3 := g.NewMap[g.String, int]()
		m3.Set("p", 100)
		m3.Set("q", 200)

		m4 := g.NewMap[g.String, int]()
		m4.Set("z", 1000)
		m4.Set("w", 2000)

		workers := g.Int(4)
		mc := &mapCounter{sleepDuration: 30 * time.Millisecond}

		res := m1.
			Iter().
			Parallel(workers).
			Inspect(mc.Fn).
			Chain(
				m2.Iter().Parallel(workers).Inspect(mc.Fn),
				m3.Iter().Parallel(workers).Inspect(mc.Fn),
				m4.Iter().Parallel(workers).Inspect(mc.Fn),
			).
			Collect()

		// All entries should be present
		if res.Len() != 8 {
			t.Errorf("expected 8 entries, got %d", res.Len())
		}

		// Should achieve high concurrency with 4 map sequences
		if mc.Max() < 4 {
			t.Errorf("expected high concurrency with 4 map sequences, got %d", mc.Max())
		}

		t.Logf("Multiple map sequences - Max concurrency: %d", mc.Max())
	})

	t.Run("HeavyMapTransformationsParallel", func(t *testing.T) {
		// Test that heavy map transformations run in parallel
		m1 := g.NewMap[g.String, int]()
		for i := 1; i <= 8; i++ {
			m1.Set("key"+g.Int(i).String(), i)
		}

		m2 := g.NewMap[g.String, int]()
		for i := 10; i <= 17; i++ {
			m2.Set("key"+g.Int(i).String(), i)
		}

		workers1 := g.Int(3)
		workers2 := g.Int(5)

		mc1 := &mapCounter{sleepDuration: 40 * time.Millisecond, sequenceID: "map1"}
		mc2 := &mapCounter{sleepDuration: 40 * time.Millisecond, sequenceID: "map2"}

		heavyMapTransform := func(k g.String, v int) (g.String, int) {
			time.Sleep(20 * time.Millisecond) // Heavy work
			return k + "_transformed", v * 2
		}

		start := time.Now()
		res := m1.
			Iter().
			Parallel(workers1).
			Map(heavyMapTransform).
			Inspect(mc1.Fn).
			Chain(
				m2.
					Iter().
					Parallel(workers2).
					Map(heavyMapTransform).
					Inspect(mc2.Fn),
			).
			Collect()
		duration := time.Since(start)

		// Debug: Print actual count
		t.Logf("m1 size: %d, m2 size: %d, result size: %d", m1.Len(), m2.Len(), res.Len())

		// Verify both sequences achieved parallelism
		if mc1.Max() < 2 {
			t.Errorf("map1 not parallel enough, max concurrency: %d", mc1.Max())
		}
		if mc2.Max() < 2 {
			t.Errorf("map2 not parallel enough, max concurrency: %d", mc2.Max())
		}

		// Verify results are transformed correctly
		expectedLen := 16 // 8 + 8
		if len(res) != expectedLen {
			t.Errorf("expected %d entries, got %d (m1: %d, m2: %d)", expectedLen, res.Len(), m1.Len(), m2.Len())
		}

		// Check transformation applied
		foundTransformed := false

		res.Iter().Range(func(k g.String, v int) bool {
			if k.Contains("_transformed") && v%2 == 0 {
				foundTransformed = true
				return false
			}
			return true
		})
		if !foundTransformed {
			t.Error("transformations not applied correctly")
		}

		t.Logf("Heavy map transforms - Map1 concurrency: %d, Map2 concurrency: %d, Duration: %v",
			mc1.Max(), mc2.Max(), duration)
	})

	t.Run("MapEarlyTermination", func(t *testing.T) {
		// Test early termination with maps
		largeMap1 := g.NewMap[g.String, int]()
		largeMap2 := g.NewMap[g.String, int]()

		for i := range 500 {
			largeMap1.Set(g.Int(i).String(), i)
			largeMap2.Set(g.Int(i+1000).String(), i+1000)
		}

		workers := g.Int(4)
		var processedCount atomic.Int64

		start := time.Now()
		res := largeMap1.
			Iter().
			Parallel(workers).
			Inspect(func(k g.String, v int) {
				processedCount.Add(1)
				time.Sleep(1 * time.Millisecond)
			}).
			Chain(
				largeMap2.
					Iter().
					Parallel(workers).
					Inspect(func(k g.String, v int) {
						processedCount.Add(1)
						time.Sleep(1 * time.Millisecond)
					}),
			).
			Take(10). // Should stop early
			Collect()
		duration := time.Since(start)

		// Should only get 10 entries
		if res.Len() != 10 {
			t.Errorf("expected 10 entries with Take(10), got %d", res.Len())
		}

		// Should process significantly fewer than 1000 entries
		processed := processedCount.Load()
		if processed > 100 {
			t.Logf("Warning: processed %d entries, early termination might not be working optimally", processed)
		}

		// Should complete much faster than processing all entries
		maxExpectedDuration := 200 * time.Millisecond
		if duration > maxExpectedDuration {
			t.Errorf("early termination too slow: %v", duration)
		}

		t.Logf("Map early termination - Processed: %d entries, Duration: %v", processed, duration)
	})

	t.Run("DifferentMapWorkerCounts", func(t *testing.T) {
		// Test map sequences with different worker counts
		m1 := g.NewMap[g.String, int]()
		m2 := g.NewMap[g.String, int]()

		for i := range 20 {
			m1.Set(g.Int(i).String(), i)
			m2.Set(g.Int(i+100).String(), i+100)
		}

		workers1 := g.Int(2)
		workers2 := g.Int(8)

		mc1 := &mapCounter{sleepDuration: 25 * time.Millisecond}
		mc2 := &mapCounter{sleepDuration: 25 * time.Millisecond}

		res := m1.
			Iter().
			Parallel(workers1).
			Inspect(mc1.Fn).
			Chain(
				m2.
					Iter().
					Parallel(workers2).
					Inspect(mc2.Fn),
			).
			Collect()

		// Verify different concurrency levels
		if mc1.Max() > int64(workers1)+1 {
			t.Errorf("map1 exceeded expected concurrency: got %d, expected ~%d", mc1.Max(), workers1)
		}
		if mc2.Max() > int64(workers2)+1 {
			t.Errorf("map2 exceeded expected concurrency: got %d, expected ~%d", mc2.Max(), workers2)
		}

		// Both should achieve some level of parallelism
		if mc1.Max() < 1 || mc2.Max() < 2 {
			t.Errorf("map sequences didn't achieve expected parallelism: map1=%d, map2=%d", mc1.Max(), mc2.Max())
		}

		// All entries should be present
		if res.Len() != 40 {
			t.Errorf("expected 40 entries, got %d", res.Len())
		}

		t.Logf("Different map workers - Map1 (%d workers): %d concurrency, Map2 (%d workers): %d concurrency",
			workers1, mc1.Max(), workers2, mc2.Max())
	})

	t.Run("MapValueCorrectness", func(t *testing.T) {
		// Test that all key-value pairs are correctly preserved
		m1 := g.NewMap[string, int]()
		m1.Set("alpha", 100)
		m1.Set("beta", 200)
		m1.Set("gamma", 300)

		m2 := g.NewMap[string, int]()
		m2.Set("delta", 400)
		m2.Set("epsilon", 500)
		m2.Set("zeta", 600)

		workers := g.Int(3)

		res := m1.
			Iter().
			Parallel(workers).
			Chain(
				m2.Iter().Parallel(workers),
			).
			Collect()

		// Verify all key-value pairs
		expectedPairs := map[string]int{
			"alpha": 100, "beta": 200, "gamma": 300,
			"delta": 400, "epsilon": 500, "zeta": 600,
		}

		if len(res) != len(expectedPairs) {
			t.Errorf("expected %d pairs, got %d", len(expectedPairs), res.Len())
		}

		for expectedKey, expectedValue := range expectedPairs {
			if actualValue := res.Get(expectedKey); actualValue.IsNone() {
				t.Errorf("missing key: %s", expectedKey)
			} else if actualValue.Some() != expectedValue {
				t.Errorf("wrong value for key %s: got %d, want %d", expectedKey, actualValue.Some(), expectedValue)
			}
		}

		t.Logf("Map value correctness - All %d pairs verified", len(expectedPairs))
	})
}
