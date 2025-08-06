package main

import (
	"sync"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cell"
)

// Config represents application configuration
type Config struct {
	DebugMode    bool
	AllowedHosts Slice[String]
}

// Cache represents cached data with timestamp
type Cache struct {
	Data      String
	Timestamp time.Time
}

// Counter represents a simple counter
type Counter struct {
	Value Int
}

// TaskInfo represents task information for workers
type TaskInfo struct {
	ID     Int
	Status String
	Result String
}

func main() {
	example1_BasicOperations()
	example2_AtomicUpdates()
	example3_CacheWithUpdate()
	example4_SwapOperations()
}

// =========================================================================
//
//	Example 1: Basic Cell Operations - Get, Set, Replace
//
// =========================================================================
func example1_BasicOperations() {
	Println("=== Example 1: Basic Cell Operations ===")

	// Create a new Cell with initial configuration
	config := cell.New(&Config{
		DebugMode:    false,
		AllowedHosts: Slice[String]{"localhost", "127.0.0.1"},
	})

	// Get current value - thread-safe read
	current := config.Get()
	Println("Initial config: DebugMode={}, Hosts={}", current.DebugMode, current.AllowedHosts)

	// Set new value - replaces entire content
	newConfig := &Config{
		DebugMode:    true,
		AllowedHosts: Slice[String]{"localhost", "127.0.0.1", "api.example.com"},
	}

	config.Set(newConfig)

	// Replace returns the old value while setting the new one
	evenNewerConfig := &Config{
		DebugMode:    true,
		AllowedHosts: Slice[String]{"*"},
	}

	oldConfig := config.Replace(evenNewerConfig)

	Println("Old config returned by Replace: DebugMode={}, Hosts={}", oldConfig.DebugMode, oldConfig.AllowedHosts)
	Println("Current config after Replace: DebugMode={}, Hosts={}", config.Get().DebugMode, config.Get().AllowedHosts)
}

// =========================================================================
//
//	Example 2: Atomic Updates with Update method
//
// =========================================================================
func example2_AtomicUpdates() {
	Println("\n=== Example 2: Atomic Updates ===")

	counter := cell.New(&Counter{Value: 0})

	var wg sync.WaitGroup

	// Launch 1000 goroutines that increment the counter
	const numWorkers = 1000

	for range numWorkers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Update atomically modifies the value
			counter.Update(func(current *Counter) *Counter {
				// Always create a copy to avoid data races
				newCounter := *current
				newCounter.Value++
				return &newCounter
			})
		}()
	}

	wg.Wait()

	Println("Final counter value: {} (should be {})", counter.Get().Value, numWorkers)

	// Demonstrate unsafe increment without Cell for comparison
	Println("\n--- Unsafe version (data race) ---")

	unsafeCounter := &Counter{Value: 0}
	wg = sync.WaitGroup{}

	for range numWorkers {
		wg.Add(1)

		go func() {
			defer wg.Done()
			unsafeCounter.Value++ // DATA RACE!
		}()
	}

	wg.Wait()

	Println("Unsafe counter value: {} (probably less than {} due to data races)", unsafeCounter.Value, numWorkers)
}

// =========================================================================
//
//	Example 3: Cache with Conditional Updates
//
// =========================================================================
func example3_CacheWithUpdate() {
	Println("\n=== Example 3: Cache with Conditional Updates ===")

	// Create cache with stale data
	cache := cell.New(&Cache{
		Data:      "stale data",
		Timestamp: time.Now().Add(-1 * time.Hour), // 1 hour old
	})

	const cacheTTL = 50 * time.Millisecond
	var wg sync.WaitGroup

	// Multiple goroutines try to update stale cache simultaneously
	for i := range 10 {
		wg.Add(1)

		go func(goroutineID int) {
			defer wg.Done()

			current := cache.Get()
			if time.Since(current.Timestamp) > cacheTTL {
				Println("Goroutine {}: Cache is stale, attempting update...", goroutineID)

				// Update with conditional logic
				cache.Update(func(c *Cache) *Cache {
					// Double-check inside Update - another goroutine might have updated it
					if time.Since(c.Timestamp) <= cacheTTL {
						Println("Goroutine {}: Cache was already updated by another goroutine", goroutineID)
						return c // No change
					}

					// We're the first to update!
					Println("Goroutine {}: Successfully updating cache", goroutineID)
					return &Cache{
						Data:      Format("fresh data from goroutine {}", goroutineID),
						Timestamp: time.Now(),
					}
				})
			} else {
				Println("Goroutine {}: Cache is fresh", goroutineID)
			}
		}(i)
	}

	wg.Wait()

	finalCache := cache.Get()

	Println("Final cache: {}", finalCache.Data)
}

// =========================================================================
//
//	Example 4: Swap Operations and Task Management
//
// =========================================================================
func example4_SwapOperations() {
	Println("\n=== Example 4: Swap Operations ===")

	// Create two cells with different task info
	task1 := cell.New(&TaskInfo{ID: 1, Status: "pending", Result: ""})
	task2 := cell.New(&TaskInfo{ID: 2, Status: "in_progress", Result: ""})

	Println("Before swap:")
	Println("  Task1: ID={}, Status={}", task1.Get().ID, task1.Get().Status)
	Println("  Task2: ID={}, Status={}", task2.Get().ID, task2.Get().Status)

	// Swap values between two cells
	task1.Swap(task2)

	Println("After swap:")
	Println("  Task1: ID={}, Status={}", task1.Get().ID, task1.Get().Status)
	Println("  Task2: ID={}, Status={}", task2.Get().ID, task2.Get().Status)

	// Use Replace for task state transitions
	Println("\n--- Task State Transitions with Replace ---")

	// Start working on the task
	oldTask := task1.Replace(&TaskInfo{ID: task1.Get().ID, Status: "completed", Result: "success"})

	Println("Task {} changed from '{}' to '{}'", oldTask.ID, oldTask.Status, task1.Get().Status)

	// Demonstrate Update for more complex state transitions
	Println("\n--- Complex State Update ---")

	task1.Update(func(current *TaskInfo) *TaskInfo {
		if current.Status == "completed" {
			// Add timestamp to result
			updated := *current
			updated.Result = Format("{} at {}", current.Result, time.Now().Format("15:04:05"))
			return &updated
		}
		return current // No change if not completed
	})

	final := task1.Get()

	Println("Final task state: ID={}, Status={}, Result={}", final.ID, final.Status, final.Result)
}
