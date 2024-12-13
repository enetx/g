package main

import (
	"errors"
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a new pool for managing tasks
	pool := NewPool[int]()

	// Set the concurrency limit to 1, ensuring only one task runs at a time
	pool.Limit(1)

	// Launch 10 tasks in the pool
	for i := range 10 {
		pool.Go(func() Result[int] {
			switch i {
			case 2:
				return Err[int](errors.New("case 2"))
			case 3:
				pool.Cancel()
			}

			return Ok(i * i)
		})
	}

	fmt.Printf("ActiveTasks: %d\n", pool.ActiveTasks())
	pool.Wait().Print()
}
