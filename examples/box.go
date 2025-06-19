package main

import (
	"sync"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/box"
)

// --- For Example 1 ---
type Config struct {
	DebugMode    bool
	AllowedHosts Slice[String]
}

// --- For Example 2 ---
type Cache struct {
	Data      String
	Timestamp time.Time
}

// --- For Example 3 ---
type Data struct {
	Counter Int
}

// --- For Example 4 ---
type WorkerState struct {
	Status String // e.g., "Idle", "Working", "Finished"
	TaskID Int
	Result String
}

func main() {
	example1()
	example2()
	example3()
	example4()
}

// =========================================================================
//
//	Example 1: Thread-Safe Configuration (Frequent Reads, Infrequent Writes)
//
// =========================================================================
func example1() {
	Println("=== Example 1: Thread-Safe Configuration ===")

	// Initial configuration
	initialConfig := &Config{
		DebugMode:    false,
		AllowedHosts: Slice[String]{"localhost", "127.0.0.1"},
	}

	configBox := box.New(initialConfig)
	var wg sync.WaitGroup

	// 100 "readers" constantly check the configuration
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Reading the configuration is completely safe and very fast (a single atomic operation)
			currentConfig := configBox.Load()
			_ = currentConfig.DebugMode // Simulate usage
		}()
	}

	// 1 "writer" atomically updates the configuration
	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(10 * time.Millisecond) // Let readers work with the old version for a bit
		Println("--> Updating configuration...")

		configBox.Update(func(current *Config) *Config {
			// 1. Create a shallow copy. This is mandatory!
			cp := *current
			// 2. Modify the copy.
			cp.DebugMode = true
			// 3. For slices/maps, you must also create a copy to avoid mutating the original.
			cp.AllowedHosts = NewSlice[String](current.AllowedHosts.Len())
			copy(cp.AllowedHosts, current.AllowedHosts)
			cp.AllowedHosts.Push("api.example.com")
			// 4. Return a pointer to the new, modified struct.
			return &cp
		})

		Println("--> Configuration updated.")
	}()

	wg.Wait()

	finalConfig := configBox.Load()

	Println("Final configuration: DebugMode={}, Hosts={}", finalConfig.DebugMode, finalConfig.AllowedHosts)
}

// =========================================================================
//
//	Example 2: Lazy Cache Update (to avoid "update stampede")
//
// =========================================================================
func example2() {
	Println("\n=== Example 2: Lazy Cache Update ===")

	cacheBox := box.New(&Cache{Data: "old data", Timestamp: time.Now().Add(-1 * time.Hour)})
	var wg sync.WaitGroup
	const cacheTTL = 50 * time.Millisecond

	// 10 goroutines simultaneously try to access the cache.
	// The cache is stale, so they will all attempt to update it.
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			currentCache := cacheBox.Load()

			// Check if the cache is stale
			if time.Since(currentCache.Timestamp) > cacheTTL {
				// The cache is stale, attempt to update it.
				// Only ONE goroutine will succeed thanks to CompareAndSwap inside Update.
				Println("Goroutine {}: Cache is stale, attempting to update.", id)
				cacheBox.Update(func(c *Cache) *Cache {
					// IMPORTANT: Re-check inside Update in case another goroutine
					// already updated the cache while we were waiting.
					if time.Since(c.Timestamp) < cacheTTL {
						Println("Goroutine {}: Someone else already updated the cache, aborting my attempt.", id)
						return c // Return the current pointer; no change will be made.
					}

					// We are the first! We get to update it.
					Println("Goroutine {}: SUCCESS! I am updating the cache.", id)
					return &Cache{
						Data:      Format("new data from goroutine {}", id),
						Timestamp: time.Now(),
					}
				})
			} else {
				Println("Goroutine {}: Cache is fresh, using it.", id)
			}
		}(i)
	}

	wg.Wait()

	Println("Final cache state: {}", cacheBox.Load().Data)
}

// =========================================================================
//
//	Example 3: Counter Example (demonstrating data race safety)
//
// =========================================================================
func example3() {
	Println("\n=== Example 3: Counter (Safe with Box) ===")

	b := box.New(&Data{Counter: 0})
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Update(func(d *Data) *Data {
				cp := *d
				cp.Counter++
				return &cp
			})
		}()
	}

	wg.Wait()
	Println("Final counter (Box): {}", b.Load().Counter)

	Println("\n=== ...and Unsafe without Box (data race) ===")
	data := &Data{Counter: 0}
	wg = sync.WaitGroup{}
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data.Counter++ // Not thread-safe!
		}()
	}
	wg.Wait()
	Println("Final counter (without Box): {}", data.Counter)
}

// =========================================================================
//
//	Example 4: Demonstrating Swap, CompareAndSwap, and UpdateAndGet
//
// =========================================================================
func example4() {
	Println("\n=== Example 4: State Management with New Methods ===")

	// Initial state: the worker is free
	initialState := &WorkerState{Status: "Idle", TaskID: 0}
	workerBox := box.New(initialState)
	Println("Initial state: [Status: {1.Status}, TaskID: {1.TaskID}]", workerBox.Load())

	// --- 1. Using CompareAndSwap to start a job ---
	// We try to start a job only if the worker is currently Idle.

	idleState := workerBox.Load() // Capture the "Idle" state
	workingState := &WorkerState{Status: "Working", TaskID: 101}

	// Atomically change Idle -> Working
	swapped := workerBox.CompareAndSwap(idleState, workingState)
	if swapped {
		Println("--> CompareAndSwap: Success! Worker has been transitioned to 'Working' state.")
		Println("Current state: [Status: {1.Status}, TaskID: {1.TaskID}]", workerBox.Load())
	} else {
		Println("--> CompareAndSwap: Failed! Worker was already busy.")
	}

	// A second attempt should fail because the state is no longer "Idle"
	swapped = workerBox.CompareAndSwap(idleState, &WorkerState{Status: "Working", TaskID: 102})
	if !swapped {
		Println("--> Second CompareAndSwap failed as expected.")
	}

	// --- 2. Using Swap to change tasks ---
	// Imagine a higher priority task arrives.
	// We atomically replace the current task with the new one and get the old one back to requeue it.

	highPriorityTask := &WorkerState{Status: "Working", TaskID: 999}
	previousState := workerBox.Swap(highPriorityTask)

	Println("\n--> Swap: Switched to a higher priority task.")
	Println("   - Old task (ID {}) was returned to the queue.", previousState.TaskID)
	Println("   - New state:  [Status: {1.Status}, TaskID: {1.TaskID}]", workerBox.Load())

	// --- 3. Using UpdateAndGet to finish the job ---
	// The worker finishes the task. We update its state and immediately get the final object.

	Println("\n--> UpdateAndGet: Finishing the job and retrieving the result...")

	finalState := workerBox.UpdateAndGet(func(current *WorkerState) *WorkerState {
		// We are working on task 999
		if current.Status == "Working" && current.TaskID == 999 {
			cp := *current
			cp.Status = "Finished"
			cp.Result = "All tasks completed successfully"
			return &cp
		}
		// If the state is unexpected, don't change anything
		return current
	})

	Println("Final state (from UpdateAndGet): [Status: {1.Status}, TaskID: {1.TaskID}, Result: {1.Result}]", finalState)
	// Verify that Load() also returns the same state
	Println(
		"Final state (verified via Load): [Status: {1.Status}, TaskID: {1.TaskID}, Result: {1.Result}]",
		workerBox.Load(),
	)
}
