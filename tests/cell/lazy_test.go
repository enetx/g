package g

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/enetx/g/cell"
)

func TestLazyBasicUsage(t *testing.T) {
	callCount := 0
	l := cell.NewLazy(func() int {
		callCount++
		return 42
	})

	// Initially not computed
	if val := l.Get(); val.IsSome() {
		t.Error("Expected None before computation, got Some")
	}

	// First Force() should compute
	result1 := l.Force()
	if result1 != 42 {
		t.Errorf("Expected 42, got %d", result1)
	}
	if callCount != 1 {
		t.Errorf("Expected function called once, called %d times", callCount)
	}

	// Get() should return Some after computation
	if val := l.Get(); val.IsNone() {
		t.Error("Expected Some after computation, got None")
	} else if val.Some() != 42 {
		t.Errorf("Expected 42, got %d", val.Some())
	}

	// Second Force() should return cached result
	result2 := l.Force()
	if result2 != 42 {
		t.Errorf("Expected 42, got %d", result2)
	}
	if callCount != 1 {
		t.Errorf("Expected function called once, called %d times", callCount)
	}
}

func TestLazyWithPointer(t *testing.T) {
	type User struct {
		Name string
	}

	callCount := 0
	l := cell.NewLazy(func() *User {
		callCount++
		return &User{Name: "John"}
	})

	// Test with pointer type
	result := l.Force()
	if result == nil {
		t.Error("Expected non-nil pointer")
	}
	if result.Name != "John" {
		t.Errorf("Expected John, got %s", result.Name)
	}
	if callCount != 1 {
		t.Errorf("Expected function called once, called %d times", callCount)
	}
}

func TestLazyWithNilPointer(t *testing.T) {
	type User struct {
		Name string
	}

	callCount := 0
	l := cell.NewLazy(func() *User {
		callCount++
		return nil // Return nil pointer
	})

	// Should work correctly with nil pointer
	result := l.Force()
	if result != nil {
		t.Error("Expected nil pointer")
	}
	if callCount != 1 {
		t.Errorf("Expected function called once, called %d times", callCount)
	}

	// Get() should still work
	if val := l.Get(); val.IsNone() {
		t.Error("Expected Some(nil), got None")
	} else if val.Some() != nil {
		t.Error("Expected nil value")
	}
}

func TestLazyWithSlice(t *testing.T) {
	callCount := 0
	l := cell.NewLazy(func() []string {
		callCount++
		return []string{"a", "b", "c"}
	})

	result := l.Force()
	if len(result) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(result))
	}
	if result[0] != "a" {
		t.Errorf("Expected 'a', got %s", result[0])
	}
	if callCount != 1 {
		t.Errorf("Expected function called once, called %d times", callCount)
	}
}

func TestLazyWithEmptySlice(t *testing.T) {
	callCount := 0
	l := cell.NewLazy(func() []string {
		callCount++
		return []string{} // Empty slice
	})

	result := l.Force()
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %v", result)
	}
	if callCount != 1 {
		t.Errorf("Expected function called once, called %d times", callCount)
	}

	// Get() should work with empty slice
	if val := l.Get(); val.IsNone() {
		t.Error("Expected Some(empty slice), got None")
	}
}

func TestLazyConcurrentAccess(t *testing.T) {
	var callCount int64
	l := cell.NewLazy(func() int {
		atomic.AddInt64(&callCount, 1)
		time.Sleep(10 * time.Millisecond) // Simulate work
		return 42
	})

	const numGoroutines = 10
	var wg sync.WaitGroup
	results := make(chan int, numGoroutines)

	// Launch multiple goroutines
	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := l.Force()
			results <- result
		}()
	}

	wg.Wait()
	close(results)

	// Check all results are the same
	expectedResult := 42
	resultCount := 0
	for result := range results {
		if result != expectedResult {
			t.Errorf("Expected %d, got %d", expectedResult, result)
		}
		resultCount++
	}

	if resultCount != numGoroutines {
		t.Errorf("Expected %d results, got %d", numGoroutines, resultCount)
	}

	// Function should be called exactly once
	if atomic.LoadInt64(&callCount) != 1 {
		t.Errorf("Expected function called once, called %d times", atomic.LoadInt64(&callCount))
	}
}

func TestLazyGetBeforeForce(t *testing.T) {
	l := cell.NewLazy(func() string {
		return "computed"
	})

	// Get() before Force() should return None
	val := l.Get()
	if val.IsSome() {
		t.Error("Expected None before Force(), got Some")
	}

	// Force() should compute
	result := l.Force()
	if result != "computed" {
		t.Errorf("Expected 'computed', got %s", result)
	}

	// Get() after Force() should return Some
	val = l.Get()
	if val.IsNone() {
		t.Error("Expected Some after Force(), got None")
	}
	if val.Some() != "computed" {
		t.Errorf("Expected 'computed', got %s", val.Some())
	}
}

func TestLazyWithZeroValues(t *testing.T) {
	// Test with zero int
	l1 := cell.NewLazy(func() int {
		return 0 // Zero value
	})

	result1 := l1.Force()
	if result1 != 0 {
		t.Errorf("Expected 0, got %d", result1)
	}

	// Get() should work with zero value
	if val := l1.Get(); val.IsNone() {
		t.Error("Expected Some(0), got None")
	} else if val.Some() != 0 {
		t.Errorf("Expected 0, got %d", val.Some())
	}

	// Test with empty string
	l2 := cell.NewLazy(func() string {
		return "" // Zero value
	})

	result2 := l2.Force()
	if result2 != "" {
		t.Errorf("Expected empty string, got %s", result2)
	}

	// Get() should work with empty string
	if val := l2.Get(); val.IsNone() {
		t.Error("Expected Some(''), got None")
	} else if val.Some() != "" {
		t.Errorf("Expected empty string, got %s", val.Some())
	}
}

// Benchmark tests
func BenchmarkLazyForce(b *testing.B) {
	l := cell.NewLazy(func() int {
		return 42
	})

	// First call to initialize
	l.Force()

	for b.Loop() {
		_ = l.Force()
	}
}

func BenchmarkLazyGet(b *testing.B) {
	l := cell.NewLazy(func() int {
		return 42
	})

	// Initialize
	l.Force()

	for b.Loop() {
		_ = l.Get()
	}
}

func BenchmarkLazyConcurrent(b *testing.B) {
	l := cell.NewLazy(func() int {
		time.Sleep(1 * time.Millisecond) // Simulate expensive computation
		return 42
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = l.Force()
		}
	})
}
