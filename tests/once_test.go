package g_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/enetx/g/cell"
)

func TestOnceCellBasicOperations(t *testing.T) {
	once := cell.NewOnce[int]()

	// Initially empty
	if once.Get().IsSome() {
		t.Error("Expected OnceCell to be empty initially")
	}
	if val := once.Get(); val.IsSome() {
		t.Error("Expected Get() to return None initially")
	}

	// Set value
	result := once.Set(42)
	if result.IsErr() {
		t.Error("Expected Set() to succeed on empty cell")
	}
	if once.Get().IsNone() {
		t.Error("Expected OnceCell to be set after Set()")
	}

	// Get value
	if val := once.Get(); val.IsNone() {
		t.Error("Expected Get() to return Some after Set()")
	} else if val.Some() != 42 {
		t.Errorf("Expected 42, got %d", val.Some())
	}

	// Try to set again - should fail
	result2 := once.Set(100)
	if result2.IsOk() {
		t.Error("Expected second Set() to fail")
	}

	// Value should remain unchanged
	if val := once.Get(); val.Some() != 42 {
		t.Errorf("Expected value to remain 42, got %d", val.Some())
	}
}

func TestOnceCellGetOrInit(t *testing.T) {
	once := cell.NewOnce[string]()
	callCount := 0

	// First call should initialize
	value1 := once.GetOrInit(func() string {
		callCount++
		return "initialized"
	})
	if value1 != "initialized" {
		t.Errorf("Expected 'initialized', got %s", value1)
	}
	if callCount != 1 {
		t.Errorf("Expected init function called once, called %d times", callCount)
	}

	// Second call should return cached value without calling init
	value2 := once.GetOrInit(func() string {
		callCount++
		return "should not be called"
	})
	if value2 != "initialized" {
		t.Errorf("Expected 'initialized', got %s", value2)
	}
	if callCount != 1 {
		t.Errorf("Expected init function called once total, called %d times", callCount)
	}
}

func TestOnceCellConcurrentSet(t *testing.T) {
	once := cell.NewOnce[int]()
	var wg sync.WaitGroup
	var successCount int64

	const numGoroutines = 100

	// Launch many goroutines trying to set different values
	for i := range numGoroutines {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			if once.Set(value).IsOk() {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	// Exactly one should succeed
	if successCount != 1 {
		t.Errorf("Expected exactly 1 successful Set(), got %d", successCount)
	}

	// Cell should contain some value
	if once.Get().IsNone() {
		t.Error("Expected cell to be set after concurrent operations")
	}

	value := once.Get().Some()
	if value < 0 || value >= numGoroutines {
		t.Errorf("Expected value in range [0, %d), got %d", numGoroutines, value)
	}
}

func TestOnceCellConcurrentGetOrInit(t *testing.T) {
	once := cell.NewOnce[string]()
	var wg sync.WaitGroup
	var initCount int64

	const numGoroutines = 50

	// Launch many goroutines trying to initialize
	results := make([]string, numGoroutines)
	for i := range numGoroutines {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			result := once.GetOrInit(func() string {
				atomic.AddInt64(&initCount, 1)
				return "concurrent_init"
			})
			results[index] = result
		}(i)
	}

	wg.Wait()

	// Init function should be called exactly once
	if initCount != 1 {
		t.Errorf("Expected init function called once, called %d times", initCount)
	}

	// All results should be the same
	expected := "concurrent_init"
	for i, result := range results {
		if result != expected {
			t.Errorf("Result[%d] = %s, expected %s", i, result, expected)
		}
	}
}

func TestOnceCellMixedOperations(t *testing.T) {
	once := cell.NewOnce[int]()
	var wg sync.WaitGroup

	// Some goroutines try to Set
	for i := range 10 {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			once.Set(value + 100) // 100-109
		}(i)
	}

	// Some goroutines try to GetOrInit
	for i := range 10 {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			result := once.GetOrInit(func() int {
				return value + 200 // 200-209
			})
			_ = result
		}(i)
	}

	wg.Wait()

	// Cell should be set to some value
	if once.Get().IsNone() {
		t.Error("Expected cell to be set")
	}

	value := once.Get().Some()
	// Value should be either from Set (100-109) or GetOrInit (200-209)
	if (value < 100 || value > 109) && (value < 200 || value > 209) {
		t.Errorf("Unexpected value: %d", value)
	}
}

func TestOnceCellTakeOnce(t *testing.T) {
	once := cell.NewOnce[string]()

	// Set a value
	once.Set("test_value")
	if once.Get().IsNone() {
		t.Error("Expected cell to be set")
	}

	// Take the value
	taken := once.Take()
	if taken.IsNone() {
		t.Error("Expected Take to return Some")
	}
	if taken.Some() != "test_value" {
		t.Errorf("Expected 'test_value', got %s", taken.Some())
	}

	// Cell should now be empty
	if once.Get().IsSome() {
		t.Error("Expected cell to be empty after Take")
	}
	if once.Get().IsSome() {
		t.Error("Expected Get() to return None after Take")
	}

	// Take on empty cell should return None
	taken2 := once.Take()
	if taken2.IsSome() {
		t.Error("Expected Take on empty cell to return None")
	}
}

func TestOnceCellWithZeroValues(t *testing.T) {
	// Test with zero int
	onceInt := cell.NewOnce[int]()
	result := onceInt.Set(0)
	if result.IsErr() {
		t.Error("Expected Set(0) to succeed")
	}
	if onceInt.Get().IsNone() {
		t.Error("Expected cell with zero value to be set")
	}
	if val := onceInt.Get(); val.IsNone() || val.Some() != 0 {
		t.Error("Expected Get() to return Some(0)")
	}

	// Test with empty string
	onceStr := cell.NewOnce[string]()
	result = onceStr.Set("")
	if result.IsErr() {
		t.Error("Expected Set(\"\") to succeed")
	}
	if onceStr.Get().IsNone() {
		t.Error("Expected cell with empty string to be set")
	}
	if val := onceStr.Get(); val.IsNone() || val.Some() != "" {
		t.Error("Expected Get() to return Some(\"\")")
	}
}

func TestOnceCellGetOrInitWithZeroValue(t *testing.T) {
	once := cell.NewOnce[int]()

	// Initialize with zero value
	value := once.GetOrInit(func() int {
		return 0
	})
	if value != 0 {
		t.Errorf("Expected 0, got %d", value)
	}
	if once.Get().IsNone() {
		t.Error("Expected cell to be set after GetOrInit with zero value")
	}
}
