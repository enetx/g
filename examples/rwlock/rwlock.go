package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/enetx/g"
)

// ==============================================================================
// RWLOCK EXAMPLES - COMPLETE GUIDE
// ==============================================================================
//
// RwLock[T] is a typed reader-writer lock that protects a value of type T.
// It allows multiple readers OR a single writer at any point in time.
// Unlike sync.RWMutex, it binds the protected data to the lock itself.
//
// Key Operations:
// - Read(): Acquire read lock - multiple readers allowed
// - Write(): Acquire write lock - exclusive access
// - RWith(): Acquire read lock, call fn with copy of value, release automatically
// - With(): Acquire write lock, call fn with pointer to value, release automatically
// - TryRead(): Non-blocking read lock attempt
// - TryWrite(): Non-blocking write lock attempt
// - ReadGuard.Get()/Deref(): Access protected value (read-only)
// - WriteGuard.Get()/Set()/Deref(): Access and modify protected value
// - Guard.Unlock(): Release the lock (typically via defer)
//
// When to use RwLock vs Mutex:
// - RwLock: Many reads, few writes (config, caches, lookup tables)
// - Mutex: Frequent writes, or read/write ratio is balanced
//
// Common Use Cases:
// - Configuration management
// - Read-heavy caches
// - Shared lookup tables
// - Feature flags
// - Rate limiters
// ==============================================================================

// Run all examples
func main() {
	BasicRwLockOperations()
	MultipleReadersExample()
	ConfigurationManager()
	ReadHeavyCacheExample()
	FeatureFlagsExample()
	RateLimiterExample()
	TryLockExamples()
	ReaderWriterFairness()
	RwLockPerformanceComparison()
	ErrorHandlingExamples()
}

// Example 1: Basic RwLock Operations
func BasicRwLockOperations() {
	fmt.Println("=== Basic RwLock Operations ===")

	data := NewRwLock("initial value")

	// Read access with RWith
	data.RWith(func(v string) {
		fmt.Printf("Read value: %s\n", v)
	})

	// Write access with With
	data.With(func(v *string) {
		fmt.Printf("Before write: %s\n", *v)
		*v = "modified value"
		fmt.Printf("After write: %s\n", *v)
	})

	// Verify change
	data.RWith(func(v string) {
		fmt.Printf("Final value: %s\n", v)
	})
}

// Example 2: Multiple Readers Simultaneously
func MultipleReadersExample() {
	fmt.Println("\n=== Multiple Readers Simultaneously ===")

	data := NewRwLock(42)
	var wg sync.WaitGroup
	var readersActive atomic.Int32

	readers := 5
	wg.Add(readers)

	for i := range readers {
		go func(id int) {
			defer wg.Done()

			// Multiple readers hold read locks concurrently
			guard := data.Read()
			readersActive.Add(1)

			current := readersActive.Load()
			fmt.Printf("Reader %d: value=%d, concurrent readers=%d\n", id, guard.Get(), current)

			time.Sleep(50 * time.Millisecond) // Hold lock briefly
			readersActive.Add(-1)
			guard.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Println("All readers finished - they ran concurrently!")
}

// Example 3: Configuration Manager
func ConfigurationManager() {
	fmt.Println("\n=== Configuration Manager ===")

	type Config struct {
		DatabaseURL string
		MaxConns    int
		Timeout     time.Duration
		Features    Map[string, bool]
	}

	config := NewRwLock(Config{
		DatabaseURL: "postgres://localhost/db",
		MaxConns:    10,
		Timeout:     30 * time.Second,
		Features:    NewMap[string, bool](),
	})

	// Many goroutines reading config (needs return value — use Read/Guard)
	getConfig := func() Config {
		guard := config.Read()
		defer guard.Unlock()
		return guard.Get()
	}

	// Simulate usage
	var wg sync.WaitGroup

	// 10 readers
	wg.Add(10)
	for i := range 10 {
		go func(id int) {
			defer wg.Done()
			cfg := getConfig()
			fmt.Printf("Reader %d: MaxConns=%d\n", id, cfg.MaxConns)
		}(i)
	}

	// 1 writer (config reload) — With for clean scoped write
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		config.With(func(c *Config) {
			c.MaxConns = 20
			c.Timeout = 60 * time.Second
		})
		fmt.Println("Config updated!")
	}()

	wg.Wait()

	final := getConfig()
	fmt.Printf("Final config: MaxConns=%d, Timeout=%v\n", final.MaxConns, final.Timeout)
}

// Example 4: Read-Heavy Cache
func ReadHeavyCacheExample() {
	fmt.Println("\n=== Read-Heavy Cache ===")

	type CacheEntry struct {
		Data      string
		CreatedAt time.Time
	}

	cache := NewRwLock(NewMap[string, CacheEntry]())

	// Read (very frequent, needs return value — use Read/Guard)
	get := func(key string) Option[string] {
		guard := cache.Read()
		defer guard.Unlock()

		if entry := guard.Deref().Get(key); entry.IsSome() {
			return Some(entry.Unwrap().Data)
		}
		return None[string]()
	}

	// Write (infrequent) — With for clean scoped write
	set := func(key, value string) {
		cache.With(func(m *Map[string, CacheEntry]) {
			m.Entry(key).OrInsert(CacheEntry{
				Data:      value,
				CreatedAt: time.Now(),
			})
		})
	}

	// Populate cache
	set("user:1", "Alice")
	set("user:2", "Bob")
	set("user:3", "Charlie")

	// Simulate read-heavy workload
	var wg sync.WaitGroup
	var readOps, writeOps atomic.Int64

	// 100 readers
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				get("user:1")
				readOps.Add(1)
			}
		}()
	}

	// 1 writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			set(fmt.Sprintf("user:%d", i+10), fmt.Sprintf("User%d", i))
			writeOps.Add(1)
			time.Sleep(time.Millisecond)
		}
	}()

	wg.Wait()
	fmt.Printf("Read operations: %d\n", readOps.Load())
	fmt.Printf("Write operations: %d\n", writeOps.Load())
	fmt.Printf("Read/Write ratio: %.0f:1\n", float64(readOps.Load())/float64(writeOps.Load()))
}

// Example 5: Feature Flags
func FeatureFlagsExample() {
	fmt.Println("\n=== Feature Flags ===")

	flags := NewRwLock(NewMap[string, bool]())

	// Initialize flags
	flags.With(func(m *Map[string, bool]) {
		m.Entry("dark_mode").OrInsert(true)
		m.Entry("new_ui").OrInsert(false)
		m.Entry("beta_features").OrInsert(false)
	})

	// Check flag (very frequent, needs return value — use Read/Guard)
	isEnabled := func(flag string) bool {
		guard := flags.Read()
		defer guard.Unlock()
		return guard.Deref().Get(flag).UnwrapOr(false)
	}

	// Toggle flag (rare - admin action)
	toggle := func(flag string) {
		flags.With(func(m *Map[string, bool]) {
			current := m.Get(flag).UnwrapOr(false)
			m.Entry(flag).OrInsert(!current)
		})
	}

	fmt.Printf("dark_mode: %t\n", isEnabled("dark_mode"))
	fmt.Printf("new_ui: %t\n", isEnabled("new_ui"))

	toggle("new_ui")
	fmt.Printf("new_ui after toggle: %t\n", isEnabled("new_ui"))

	// Concurrent reads don't block each other
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			isEnabled("dark_mode")
		}()
	}
	wg.Wait()
	fmt.Println("100 concurrent flag checks completed")
}

// Example 6: Rate Limiter
func RateLimiterExample() {
	fmt.Println("\n=== Rate Limiter ===")

	type RateLimitState struct {
		Requests  int
		ResetTime time.Time
	}

	limiter := NewRwLock(NewMap[string, RateLimitState]())
	limit := 5
	window := 100 * time.Millisecond

	// Check if allowed (read-heavy, needs return value — use Read/Guard)
	isAllowed := func(clientID string) bool {
		guard := limiter.Read()
		defer guard.Unlock()

		state := guard.Deref().Get(clientID)
		if state.IsNone() {
			return true // New client
		}

		s := state.Unwrap()
		if time.Now().After(s.ResetTime) {
			return true // Window expired
		}

		return s.Requests < limit
	}

	// Record request (write) — With for clean scoped write
	recordRequest := func(clientID string) {
		limiter.With(func(m *Map[string, RateLimitState]) {
			now := time.Now()
			state := m.Get(clientID)

			var newState RateLimitState
			if state.IsNone() || now.After(state.Unwrap().ResetTime) {
				newState = RateLimitState{
					Requests:  1,
					ResetTime: now.Add(window),
				}
			} else {
				s := state.Unwrap()
				newState = RateLimitState{
					Requests:  s.Requests + 1,
					ResetTime: s.ResetTime,
				}
			}

			m.Entry(clientID).OrInsert(newState)
		})
	}

	// Simulate requests
	for i := range 10 {
		clientID := "client-1"
		allowed := isAllowed(clientID)
		if allowed {
			recordRequest(clientID)
			fmt.Printf("Request %d: ALLOWED\n", i+1)
		} else {
			fmt.Printf("Request %d: RATE LIMITED\n", i+1)
		}
	}

	// Wait for window to reset
	time.Sleep(window)

	if isAllowed("client-1") {
		fmt.Println("After window reset: ALLOWED")
	}
}

// Example 7: TryLock Examples
func TryLockExamples() {
	fmt.Println("\n=== TryLock Examples ===")

	data := NewRwLock(100)

	// TryRead when no locks held
	if opt := data.TryRead(); opt.IsSome() {
		guard := opt.Unwrap()
		fmt.Printf("TryRead succeeded: %d\n", guard.Get())
		guard.Unlock()
	}

	// TryWrite when no locks held
	if opt := data.TryWrite(); opt.IsSome() {
		guard := opt.Unwrap()
		guard.Set(200)
		fmt.Printf("TryWrite succeeded, new value: %d\n", guard.Get())
		guard.Unlock()
	}

	// TryWrite fails when read lock held
	readGuard := data.Read()
	if opt := data.TryWrite(); opt.IsNone() {
		fmt.Println("TryWrite failed while read lock held (expected)")
	}

	// TryRead succeeds when read lock held (multiple readers OK)
	if opt := data.TryRead(); opt.IsSome() {
		fmt.Println("TryRead succeeded while another read lock held (expected)")
		opt.Unwrap().Unlock()
	}
	readGuard.Unlock()

	// TryRead fails when write lock held
	writeGuard := data.Write()
	if opt := data.TryRead(); opt.IsNone() {
		fmt.Println("TryRead failed while write lock held (expected)")
	}
	writeGuard.Unlock()
}

// Example 8: Reader-Writer Fairness
func ReaderWriterFairness() {
	fmt.Println("\n=== Reader-Writer Behavior ===")

	data := NewRwLock(0)
	var wg sync.WaitGroup
	var log []string
	var logMu sync.Mutex

	logEvent := func(event string) {
		logMu.Lock()
		log = append(log, event)
		logMu.Unlock()
	}

	// Start readers (need to hold lock for a duration — use Read/Guard)
	for i := range 3 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			guard := data.Read()
			logEvent(fmt.Sprintf("Reader %d started", id))
			time.Sleep(20 * time.Millisecond)
			_ = guard.Get()
			logEvent(fmt.Sprintf("Reader %d finished", id))
			guard.Unlock()
		}(i)
	}

	// Start writer after small delay
	time.Sleep(5 * time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		logEvent("Writer waiting")
		data.With(func(v *int) {
			logEvent("Writer acquired lock")
			*v = 42
			time.Sleep(10 * time.Millisecond)
			logEvent("Writer finished")
		})
	}()

	wg.Wait()

	fmt.Println("Event log:")
	for _, event := range log {
		fmt.Printf("  %s\n", event)
	}
}

// Example 9: Performance Comparison (RwLock vs Mutex)
func RwLockPerformanceComparison() {
	fmt.Println("\n=== Performance: RwLock vs Mutex ===")

	readRatio := 95 // 95% reads, 5% writes
	ops := 10000

	// Test RwLock
	rwData := NewRwLock(0)
	start := time.Now()
	var wg sync.WaitGroup

	for i := range ops {
		wg.Add(1)
		go func(op int) {
			defer wg.Done()
			if op%100 < readRatio {
				rwData.RWith(func(v int) { _ = v })
			} else {
				rwData.With(func(v *int) { *v++ })
			}
		}(i)
	}
	wg.Wait()
	rwTime := time.Since(start)

	// Test Mutex
	muData := NewMutex(0)
	start = time.Now()

	for i := range ops {
		wg.Add(1)
		go func(op int) {
			defer wg.Done()
			if op%100 < readRatio {
				guard := muData.Lock()
				_ = guard.Get()
				guard.Unlock()
			} else {
				muData.With(func(v *int) { *v++ })
			}
		}(i)
	}
	wg.Wait()
	muTime := time.Since(start)

	fmt.Printf("Operations: %d (%d%% reads, %d%% writes)\n", ops, readRatio, 100-readRatio)
	fmt.Printf("RwLock time: %v\n", rwTime)
	fmt.Printf("Mutex time: %v\n", muTime)

	if rwTime < muTime {
		fmt.Printf("RwLock is %.1fx faster for read-heavy workload\n", float64(muTime)/float64(rwTime))
	} else {
		fmt.Println("Mutex was faster (overhead may dominate at this scale)")
	}
}

// Example 10: Error Handling
func ErrorHandlingExamples() {
	fmt.Println("\n=== Error Handling ===")

	// Zero value
	zeroLock := NewRwLock(0)
	zeroLock.RWith(func(v int) {
		fmt.Printf("Zero value works: %d\n", v)
	})

	// Nil pointer
	var ptr *string
	nilLock := NewRwLock(ptr)
	nilLock.RWith(func(v *string) {
		fmt.Printf("Nil pointer works: %v\n", v == nil)
	})

	// Empty collections
	emptyLock := NewRwLock(NewSlice[int]())
	emptyLock.RWith(func(sl Slice[int]) {
		fmt.Printf("Empty slice works: len=%d\n", sl.Len())
	})

	// With is panic-safe: defer inside With ensures unlock
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic, but lock was released!")
			}
		}()

		m := NewRwLock(42)
		m.With(func(*int) {
			fmt.Println("Work completed successfully")
		})
	}()

	fmt.Println("All edge cases handled correctly")
}
