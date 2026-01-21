package main

import (
	"fmt"
	"sync"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

// ==============================================================================
// MUTEX EXAMPLES - COMPLETE GUIDE
// ==============================================================================
//
// Mutex[T] is a typed mutual exclusion lock that protects a value of type T.
// Unlike sync.Mutex, it binds the protected data to the lock itself,
// making it impossible to access the data without holding the lock.
//
// Key Operations:
// - Lock(): Acquire lock, returns MutexGuard - blocks until acquired
// - TryLock(): Attempt to acquire without blocking - returns Option[MutexGuard]
// - Guard.Get(): Get a copy of the protected value
// - Guard.Set(): Replace the protected value
// - Guard.Deref(): Get pointer for direct manipulation
// - Guard.Unlock(): Release the lock (typically via defer)
//
// Advantages over sync.Mutex:
// - Data and lock are bound together - can't forget to lock
// - Impossible to access data without holding the lock
// - Type-safe - compiler ensures correct usage
// - Self-documenting - clear what data is protected
//
// Common Use Cases:
// - Shared counters
// - Thread-safe caches
// - Protected configuration
// - Concurrent data structures
// ==============================================================================

// Run all examples
func main() {
	BasicMutexOperations()
	CounterExample()
	SharedCacheExample()
	ProtectedConfigExample()
	TryLockExample()
	DerefExample()
	CompareWithSyncMutex()
	ConcurrentMapExample()
	MetricsCollector()
	ErrorHandlingExample()
}

// Example 1: Basic Mutex Operations
func BasicMutexOperations() {
	fmt.Println("=== Basic Mutex Operations ===")

	// Create a mutex protecting an integer
	counter := NewMutex(0)

	// Lock and modify
	guard := counter.Lock()
	fmt.Printf("Initial value: %d\n", guard.Get())

	guard.Set(42)
	fmt.Printf("After Set(42): %d\n", guard.Get())

	guard.Unlock()

	// Lock again to verify
	guard2 := counter.Lock()
	defer guard2.Unlock()
	fmt.Printf("Value persisted: %d\n", guard2.Get())
}

// Example 2: Thread-Safe Counter
func CounterExample() {
	fmt.Println("\n=== Thread-Safe Counter ===")

	counter := NewMutex(0)
	var wg sync.WaitGroup

	// 100 goroutines incrementing counter
	workers := 100
	incrementsPerWorker := 1000

	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for range incrementsPerWorker {
				guard := counter.Lock()
				guard.Set(guard.Get() + 1)
				guard.Unlock()
			}
		}()
	}

	wg.Wait()

	guard := counter.Lock()
	defer guard.Unlock()

	expected := workers * incrementsPerWorker
	fmt.Printf("Expected: %d, Got: %d\n", expected, guard.Get())
	fmt.Printf("Race-free: %t\n", guard.Get() == expected)
}

// Example 3: Shared Cache
func SharedCacheExample() {
	fmt.Println("\n=== Shared Cache ===")

	type CacheEntry struct {
		Value     string
		ExpiresAt time.Time
	}

	cache := NewMutex(NewMap[string, CacheEntry]())

	// Set cache entry
	set := func(key, value string, ttl time.Duration) {
		guard := cache.Lock()
		defer guard.Unlock()

		guard.Deref().Entry(key).OrInsert(CacheEntry{
			Value:     value,
			ExpiresAt: time.Now().Add(ttl),
		})

		fmt.Printf("Cached: %s = %s\n", key, value)
	}

	// Get cache entry
	get := func(key string) Option[string] {
		guard := cache.Lock()
		defer guard.Unlock()

		entry := guard.Deref().Get(key)
		if entry.IsNone() {
			return None[string]()
		}

		e := entry.Unwrap()
		if time.Now().After(e.ExpiresAt) {
			guard.Deref().Delete(key)
			return None[string]()
		}

		return Some(e.Value)
	}

	// Usage
	set("user:1", "Alice", time.Hour)
	set("user:2", "Bob", time.Hour)

	if val := get("user:1"); val.IsSome() {
		fmt.Printf("Got user:1 = %s\n", val.Unwrap())
	}

	if val := get("user:3"); val.IsNone() {
		fmt.Println("user:3 not found (expected)")
	}
}

// Example 4: Protected Configuration
func ProtectedConfigExample() {
	fmt.Println("\n=== Protected Configuration ===")

	type AppConfig struct {
		Debug    bool
		LogLevel string
		MaxConns int
	}

	config := NewMutex(AppConfig{
		Debug:    false,
		LogLevel: "info",
		MaxConns: 100,
	})

	// Read config
	readConfig := func() AppConfig {
		guard := config.Lock()
		defer guard.Unlock()
		return guard.Get()
	}

	// Update config
	updateConfig := func(fn func(*AppConfig)) {
		guard := config.Lock()
		defer guard.Unlock()
		fn(guard.Deref())
	}

	fmt.Printf("Initial config: %+v\n", readConfig())

	// Update in different goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		updateConfig(func(c *AppConfig) {
			c.Debug = true
			c.LogLevel = "debug"
		})
		fmt.Println("Enabled debug mode")
	}()

	go func() {
		defer wg.Done()
		updateConfig(func(c *AppConfig) {
			c.MaxConns = 200
		})
		fmt.Println("Increased max connections")
	}()

	wg.Wait()
	fmt.Printf("Final config: %+v\n", readConfig())
}

// Example 5: TryLock for Non-Blocking Access
func TryLockExample() {
	fmt.Println("\n=== TryLock (Non-Blocking) ===")

	resource := NewMutex("shared resource")

	// First lock succeeds
	guard1 := resource.Lock()
	fmt.Println("First lock acquired")

	// Try to acquire from another "thread"
	result := resource.TryLock()
	if result.IsNone() {
		fmt.Println("TryLock failed - resource busy (expected)")
	}

	guard1.Unlock()
	fmt.Println("First lock released")

	// Now TryLock should succeed
	if opt := resource.TryLock(); opt.IsSome() {
		guard := opt.Unwrap()
		fmt.Printf("TryLock succeeded, value: %s\n", guard.Get())
		guard.Unlock()
	}

	// Practical use: skip if busy
	processIfAvailable := func(m *Mutex[int]) bool {
		if opt := m.TryLock(); opt.IsSome() {
			guard := opt.Unwrap()
			defer guard.Unlock()
			// Process...
			guard.Set(guard.Get() + 1)
			return true
		}
		return false // Skip, resource busy
	}

	counter := NewMutex(0)
	processed := processIfAvailable(counter)
	fmt.Printf("Processed: %t\n", processed)
}

// Example 6: Deref for Direct Manipulation
func DerefExample() {
	fmt.Println("\n=== Deref for Direct Manipulation ===")

	// With slice
	numbers := NewMutex(SliceOf(1, 2, 3))

	guard := numbers.Lock()
	// Direct manipulation via pointer
	guard.Deref().Push(4, 5, 6)
	guard.Deref().SortBy(func(a, b int) cmp.Ordering {
		return cmp.Cmp(b, a) // Descending
	})
	fmt.Printf("Slice after Deref operations: %v\n", guard.Get())
	guard.Unlock()

	// With struct
	type Stats struct {
		Count int
		Sum   int
	}

	stats := NewMutex(Stats{})

	addValue := func(v int) {
		guard := stats.Lock()
		defer guard.Unlock()

		s := guard.Deref()
		s.Count++
		s.Sum += v
	}

	addValue(10)
	addValue(20)
	addValue(30)

	guard2 := stats.Lock()
	defer guard2.Unlock()
	s := guard2.Get()
	fmt.Printf("Stats: Count=%d, Sum=%d, Avg=%.1f\n", s.Count, s.Sum, float64(s.Sum)/float64(s.Count))
}

// Example 7: Compare with sync.Mutex
func CompareWithSyncMutex() {
	fmt.Println("\n=== Comparison: Mutex[T] vs sync.Mutex ===")

	// ❌ Traditional sync.Mutex - easy to make mistakes
	fmt.Println("\n--- Traditional sync.Mutex (problematic) ---")
	fmt.Println(`
type UserServiceOld struct {
    mu    sync.Mutex
    users map[string]User  // What protects this?
    cache map[string]Data  // Same mutex? Different? Unclear!
}

func (s *UserServiceOld) GetUser(id string) User {
    // Bug: forgot to lock!
    return s.users[id]
}

func (s *UserServiceOld) UpdateCache(id string, data Data) {
    // Bug: using wrong lock? No compile error!
    s.cache[id] = data
}
`)

	// With Mutex[T] - compile-time safety
	fmt.Println("--- With Mutex[T] (safe) ---")
	fmt.Println(`
type UserServiceNew struct {
    users *Mutex[Map[string, User]]
    cache *Mutex[Map[string, Data]]
}

func (s *UserServiceNew) GetUser(id string) Option[User] {
    guard := s.users.Lock()  // Must lock to access
    defer guard.Unlock()
    return guard.Deref().Get(id)
}

func (s *UserServiceNew) UpdateCache(id string, data Data) {
    guard := s.cache.Lock()  // Each data has its own lock
    defer guard.Unlock()
    guard.Deref().Entry(id).OrInsert(data)
}
`)

	fmt.Println("Key differences:")
	fmt.Println("1. Data and lock are bound - can't access data without lock")
	fmt.Println("2. Each field has explicit protection")
	fmt.Println("3. Compiler catches mistakes at compile time")
	fmt.Println("4. Self-documenting code")
}

// Example 8: Concurrent Map Operations
func ConcurrentMapExample() {
	fmt.Println("\n=== Concurrent Map Operations ===")

	userScores := NewMutex(NewMap[string, int]())
	var wg sync.WaitGroup

	// Multiple goroutines updating scores
	users := SliceOf("alice", "bob", "charlie", "diana")
	updates := 100

	wg.Add(users.Len().Std())

	for _, user := range users {
		go func(name string) {
			defer wg.Done()
			for range updates {
				guard := userScores.Lock()
				guard.Deref().Entry(name).AndModify(func(v *int) { *v++ }).OrInsert(1)
				guard.Unlock()
			}
		}(user)
	}

	wg.Wait()

	// Print final scores
	guard := userScores.Lock()
	defer guard.Unlock()

	fmt.Println("Final scores:")
	for key, value := range guard.Deref().Iter() {
		fmt.Printf("  %s: %d\n", key, value)
	}
}

// Example 9: Metrics Collector
func MetricsCollector() {
	fmt.Println("\n=== Metrics Collector ===")

	type Metrics struct {
		RequestCount int64
		ErrorCount   int64
		TotalLatency time.Duration
	}

	metrics := NewMutex(Metrics{})

	recordRequest := func(latency time.Duration, isError bool) {
		guard := metrics.Lock()
		defer guard.Unlock()

		m := guard.Deref()
		m.RequestCount++
		m.TotalLatency += latency

		if isError {
			m.ErrorCount++
		}
	}

	getMetrics := func() (requests, errors int64, avgLatency time.Duration) {
		guard := metrics.Lock()
		defer guard.Unlock()

		m := guard.Get()
		if m.RequestCount > 0 {
			avgLatency = m.TotalLatency / time.Duration(m.RequestCount)
		}
		return m.RequestCount, m.ErrorCount, avgLatency
	}

	// Simulate requests
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			latency := time.Duration(id) * time.Millisecond
			isError := id%10 == 0 // 10% error rate
			recordRequest(latency, isError)
		}(i)
	}

	wg.Wait()

	requests, errors, avgLatency := getMetrics()
	fmt.Printf("Requests: %d\n", requests)
	fmt.Printf("Errors: %d (%.1f%%)\n", errors, float64(errors)/float64(requests)*100)
	fmt.Printf("Avg Latency: %v\n", avgLatency)
}

// Example 10: Error Handling and Edge Cases
func ErrorHandlingExample() {
	fmt.Println("\n=== Error Handling and Edge Cases ===")

	// Zero value
	zeroMutex := NewMutex(0)
	guard := zeroMutex.Lock()
	fmt.Printf("Zero value works: %d\n", guard.Get())
	guard.Unlock()

	// Nil pointer
	var ptr *string
	nilMutex := NewMutex(ptr)
	guard2 := nilMutex.Lock()
	fmt.Printf("Nil pointer works: %v\n", guard2.Get() == nil)
	guard2.Unlock()

	// Empty collections
	emptySlice := NewMutex(NewSlice[int]())
	guard3 := emptySlice.Lock()
	fmt.Printf("Empty slice works: len=%d\n", guard3.Deref().Len())
	guard3.Unlock()

	// Proper cleanup with defer
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic, but lock was released!")
			}
		}()

		m := NewMutex(42)
		guard := m.Lock()
		defer guard.Unlock() // Always released, even on panic

		// Simulate work that might panic
		fmt.Println("Work completed successfully")
	}()

	fmt.Println("All edge cases handled correctly")
}
