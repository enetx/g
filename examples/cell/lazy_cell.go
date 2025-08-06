package main

import (
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cell"
)

func main() {
	fmt.Println("=== LazyCell Example ===")

	// Expensive computation that we want to defer
	expensive := cell.NewLazy(func() int {
		Println("Computing expensive value...")
		time.Sleep(100 * time.Millisecond)
		return 42
	})

	fmt.Println("LazyCell created, computation not run yet")

	// Check if initialized
	if val := expensive.Get(); val.IsSome() {
		Println("Already computed: {}", val.Some())
	} else {
		Println("Not computed yet")
	}

	// First call - triggers computation
	fmt.Println("First Force():")
	result1 := expensive.Force()
	Println("Result: {}", result1)

	// Second call - uses cached value
	fmt.Println("Second Force():")
	result2 := expensive.Force()
	Println("Result: {}", result2)

	// Check if initialized now
	if val := expensive.Get(); val.IsSome() {
		fmt.Println("Now cached:", val.Some())
	}

	Println("\n=== Transform Example ===")

	// Chain computations with Transform
	lazy1 := cell.NewLazy(func() int {
		Println("Computing base value...")
		return 10
	})

	lazy2 := cell.Transform(lazy1, func(i int) String {
		Println("Mapping to string...")
		return Format("Value: {}", i*2)
	})

	Println("Transformed lazy created")
	result := lazy2.Force()
	Println("Final result: {}", result)

	Println("\n=== Concurrent Access ===")

	concurrent := cell.NewLazy(func() int {
		Println("Computing in concurrent context...")
		time.Sleep(50 * time.Millisecond)

		return 99
	})

	// Multiple goroutines accessing the same lazy value
	done := make(chan int, 3)

	for i := range 3 {
		go func(id int) {
			val := concurrent.Force()
			Println("Goroutine {} got: {}", id, val)
			done <- val
		}(i)
	}

	// Wait for all goroutines
	for range 3 {
		<-done
	}

	Println("All goroutines completed")
}
