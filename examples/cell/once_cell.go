package main

import (
	"sync"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cell"
)

// Global configuration that should be initialized only once
var globalConfig = cell.NewOnce[*AppConfig]()

type AppConfig struct {
	DatabaseURL string
	APIKey      string
	Debug       bool
}

func main() {
	example1_BasicUsage()
	example2_GlobalConfiguration()
	example3_ConcurrentInitialization()
	example4_ConditionalInitialization()
}

// =========================================================================
//
//	Example 1: Basic OnceCell Usage
//
// =========================================================================
func example1_BasicUsage() {
	Println("=== Example 1: Basic OnceCell Usage ===")

	// Create empty OnceCell
	once := cell.NewOnce[string]()
	Println("Is set initially: {}", once.Get().IsSome())

	// Try to get value (should be None)
	if val := once.Get(); val.IsNone() {
		Println("Cell is empty, as expected")
	}

	// Set value for the first time
	result := once.Set("Hello OnceCell!")
	Println("First Set() succeeded: {}", result.IsOk())
	Println("Is set now: {}", once.Get().IsSome())

	// Get the value
	if val := once.Get(); val.IsSome() {
		Println("Current value: {}", val.Some())
	}

	// Try to set again (should fail)
	result2 := once.Set("This won't work")
	Println("Second Set() succeeded: {}", result2.IsOk())
	if result2.IsErr() {
		Println("Error: {}", result2.Err())
	}
	Println("Value remains: {}", once.Get().Some())
}

// =========================================================================
//
//	Example 2: Global Configuration Pattern
//
// =========================================================================
func example2_GlobalConfiguration() {
	Println("\n=== Example 2: Global Configuration ===")

	// Function to initialize config (expensive operation)
	initConfig := func() *AppConfig {
		Println("Initializing configuration... (this happens only once)")
		time.Sleep(10 * time.Millisecond) // Simulate expensive operation
		return &AppConfig{
			DatabaseURL: "postgres://localhost/myapp",
			APIKey:      "secret_api_key_12345",
			Debug:       true,
		}
	}

	// Multiple parts of the application try to get config
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// GetOrInit ensures initialization happens exactly once
			config := globalConfig.GetOrInit(initConfig)
			Println("Worker {}: Got config with DatabaseURL: {}", id, config.DatabaseURL)
		}(i)
	}

	wg.Wait()
	Println("All workers completed. Config initialized exactly once.")
}

// =========================================================================
//
//	Example 3: Concurrent Initialization Race
//
// =========================================================================
func example3_ConcurrentInitialization() {
	Println("\n=== Example 3: Concurrent Initialization Race ===")

	raceCell := cell.NewOnce[int]()
	var wg sync.WaitGroup
	const numWorkers = 10

	// Multiple goroutines try to set different values
	for i := range numWorkers {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()

			result := raceCell.Set(value)
			if result.IsOk() {
				Println("Worker {} successfully set value {}", value, value)
			} else {
				Println("Worker {} failed to set value {}", value, value)
			}
		}(i)
	}

	wg.Wait()

	finalValue := raceCell.Get().Some()
	Println("Final value: {} (one of 0-{})", finalValue, numWorkers-1)
}

// =========================================================================
//
//	Example 4: Conditional Initialization
//
// =========================================================================
func example4_ConditionalInitialization() {
	Println("\n=== Example 4: Conditional Initialization ===")

	// OnceCell for caching expensive computation result
	expensiveResult := cell.NewOnce[string]()

	// Function that might need the expensive result
	processData := func(id int) {
		// Only compute if we actually need it
		result := expensiveResult.GetOrInit(func() string {
			Println("Performing expensive computation... (only once)")
			time.Sleep(20 * time.Millisecond)
			return "EXPENSIVE_RESULT_" + time.Now().Format("15:04:05")
		})

		Println("Worker {}: Using result: {}", id, result)
	}

	// Simulate different workers that might or might not need the result
	var wg sync.WaitGroup

	// First batch - these will trigger initialization
	for i := range 3 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			processData(id)
		}(i)
	}

	wg.Wait()

	// Small delay to show the difference
	time.Sleep(50 * time.Millisecond)

	// Second batch - these will use cached result
	Println("--- Later workers use cached result ---")
	for i := 3; i < 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			processData(id)
		}(i)
	}

	wg.Wait()

	// Demonstrate Take
	Println("\n--- Demonstrating Take ---")
	taken := expensiveResult.Take()
	if taken.IsSome() {
		Println("Took value: {}", taken.Some())
		Println("Cell is now empty: {}", expensiveResult.Get().IsNone())
	}
}
