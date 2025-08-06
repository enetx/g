package main

import (
	"sync"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cell"
)

func main() {
	Println("=== Cell Types Examples ===\n")

	basicCellUsage()
	onceCellUsage()
	lazyCellUsage()
	practicalExamples()
}

// Cell[T] - thread-safe container for mutable values
func basicCellUsage() {
	Println("Cell[T] - thread-safe values:")

	// Simple usage
	counter := cell.New(0)
	counter.Set(100)
	Println("Value: {}", counter.Get()) // 100

	// Atomic update
	counter.Update(func(old int) int { return old + 50 })
	Println("After +50: {}", counter.Get()) // 150

	// Replace returns old value
	old := counter.Replace(999)
	Println("Replaced {} with {}", old, counter.Get()) // 150 with 999

	Println("")
}

// OnceCell[T] - value set only once
func onceCellUsage() {
	Println("OnceCell[T] - set only once:")

	apiKey := cell.NewOnce[string]()

	// First set
	if apiKey.Set("secret123").IsOk() {
		Println("API key set successfully")
	}

	// Second set fails
	if apiKey.Set("other_key").IsErr() {
		Println("Second set blocked")
	}

	// Get value
	if key := apiKey.Get(); key.IsSome() {
		Println("API key: {}", key.Some())
	}

	// GetOrInit - sets if not set
	dbConfig := cell.NewOnce[string]()
	config := dbConfig.GetOrInit(func() string {
		return "postgres://localhost/db"
	})
	Println("DB config: {}", config)

	Println("")
}

// LazyCell[T] - expensive computations that may not be needed
func lazyCellUsage() {
	Println("LazyCell[T] - lazy computations:")

	// Expensive computation executed only when Force() is called
	fibonacci := cell.NewLazy(func() int {
		Println("Computing 40th Fibonacci number...")
		return fib(40) // This takes time
	})

	Println("LazyCell created, but function not called yet")

	// Only here the computation happens
	result := fibonacci.Force()
	Println("Result: {}", result)

	// Second call uses cached value
	result2 := fibonacci.Force()
	Println("From cache: {}", result2)

	Println("")
}

// Practical examples
func practicalExamples() {
	Println("Practical Examples:")

	// 1. Global request counter
	requestCounter := cell.New(0)
	for range 5 {
		requestCounter.Update(func(count int) int { return count + 1 })
	}
	Println("Total requests: {}", requestCounter.Get())

	// 2. Application config (set once)
	appConfig := cell.NewOnce[string]()
	appConfig.Set("production")
	if cfg := appConfig.Get(); cfg.IsSome() {
		Println("App mode: {}", cfg.Some())
	}

	// 3. Lazy database loading
	database := cell.NewLazy(func() string {
		Println("Connecting to database...")
		time.Sleep(50 * time.Millisecond) // Simulate connection
		return "connected to PostgreSQL"
	})

	// Database connects only on first access
	Println("Database: {}", database.Force())

	// 4. Concurrent usage
	Println("\nConcurrent example:")
	sharedValue := cell.New(0)
	var wg sync.WaitGroup

	// 10 goroutines increment counter
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sharedValue.Update(func(v int) int { return v + 1 })
		}()
	}

	wg.Wait()
	Println("Final value: {}", sharedValue.Get()) // 10

	// 5. OnceCell with Option for nullable values
	Println("\nOnceCell with Option:")
	userSession := cell.NewOnce[Option[string]]()

	// User not logged in
	userSession.Set(None[string]())
	if session := userSession.Get(); session.IsSome() {
		sessionOpt := session.Some()
		if sessionOpt.IsNone() {
			Println("User not authorized")
		}
	}
}

// Simple Fibonacci function (inefficient, for LazyCell demo)
func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
