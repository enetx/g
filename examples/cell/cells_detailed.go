package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cell"
)

func main() {
	Println("=== Detailed Cell Examples ===\n")

	cellExample()
	Println("")
	onceCellExample()
	Println("")
	lazyCellExample()
}

// Cell[T] - for safe value modification
func cellExample() {
	Println("Cell[T] - thread-safe container:")

	// Create Cell with number
	counter := cell.New(0)
	Println("Initial value: {}", counter.Get())

	// Change value
	counter.Set(10)
	Println("After Set(10): {}", counter.Get())

	// Replace and get old value
	old := counter.Replace(20)
	Println("Replace(20) returned: {}, new: {}", old, counter.Get())

	// Update through function
	counter.Update(func(current int) int {
		return current * 2
	})
	Println("After Update(*2): {}", counter.Get())

	// Example with text
	text := cell.New("Hello")
	text.Update(func(s string) string {
		return s + " World!"
	})
	Println("Text: {}", text.Get())
}

// OnceCell[T] - for setting value once
func onceCellExample() {
	Println("OnceCell[T] - set only once:")

	config := cell.NewOnce[string]()

	// Check if empty
	if config.Get().IsNone() {
		Println("OnceCell is empty")
	}

	// Set value first time
	result1 := config.Set("production")
	Println("First Set('production'): {}", result1.IsOk())

	// Try to set second time
	result2 := config.Set("development")
	Println("Second Set('development'): {}", result2.IsOk())
	if result2.IsErr() {
		Println("Error: {}", result2.Err())
	}

	// Get value
	if val := config.Get(); val.IsSome() {
		Println("Value: {}", val.Some())
	}

	// GetOrInit - sets if empty
	dbUrl := cell.NewOnce[string]()
	url := dbUrl.GetOrInit(func() string {
		Println("Initializing database URL...")
		return "postgres://localhost/myapp"
	})
	Println("Database URL: {}", url)

	// Second GetOrInit call doesn't execute function
	url2 := dbUrl.GetOrInit(func() string {
		Println("This function won't execute!")
		return "other URL"
	})
	Println("Database URL (again): {}", url2)

	// Take - extract and clear
	taken := config.Take()
	if taken.IsSome() {
		Println("Extracted: {}", taken.Some())
		Println("Now OnceCell is empty: {}", config.Get().IsNone())
	}
}

// LazyCell[T] - for lazy computations
func lazyCellExample() {
	Println("LazyCell[T] - lazy computations:")

	// Create LazyCell with expensive operation
	expensive := cell.NewLazy(func() int {
		Println("Performing expensive computation...")
		time.Sleep(100 * time.Millisecond)
		return 42
	})

	Println("LazyCell created, but function not called yet")

	// Check if computed already
	if expensive.Get().IsNone() {
		Println("Value not computed yet")
	}

	// First Force() call - executes computation
	Println("Calling Force() first time:")
	result1 := expensive.Force()
	Println("Result: {}", result1)

	// Second Force() call - returns cached result
	Println("Calling Force() second time:")
	result2 := expensive.Force()
	Println("Result (from cache): {}", result2)

	// Now Get() returns Some
	if val := expensive.Get(); val.IsSome() {
		Println("Get() now returns: {}", val.Some())
	}

	// Example with Option - LazyCell can return Option
	maybeValue := cell.NewLazy(func() Option[string] {
		Println("Loading config from file...")
		// Simulate file might or might not exist
		if time.Now().UnixNano()%2 == 0 {
			return Some("config loaded")
		}
		return None[string]()
	})

	config := maybeValue.Force()
	if config.IsSome() {
		Println("Config loaded: {}", config.Some())
	} else {
		Println("Config not found")
	}
}
