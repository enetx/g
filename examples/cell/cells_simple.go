package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cell"
)

func main() {
	// Cell[T] - thread-safe container for mutable values
	counter := cell.New(0)
	counter.Set(42)
	Println("Cell value: {}", counter.Get()) // 42

	// OnceCell[T] - value set only once
	config := cell.NewOnce[string]()
	config.Set("production")
	Println("OnceCell value: {}", config.Get().Some()) // production

	// LazyCell[T] - computation executed only on first access
	lazy := cell.NewLazy(func() string {
		time.Sleep(10 * time.Millisecond) // expensive operation
		return "computed result"
	})
	Println("LazyCell value: {}", lazy.Force()) // computed result
}

/*
When to use:

Cell[T]:
- Counters, caches, configuration
- When you need thread-safe access to mutable data
- counter.Update(func(old int) int { return old + 1 })

OnceCell[T]:
- Global configuration, singletons
- When value should be set only once
- config.GetOrInit(func() Config { return loadConfig() })

LazyCell[T]:
- Expensive computations, database connections
- When result might not be needed
- dbPool.Force() // connects only when needed
*/
