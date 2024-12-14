package main

import (
	"errors"

	. "github.com/enetx/g"
)

func main() {
	pool := NewPool[int]() // Create a new pool for managing tasks
	pool.Limit(1)          // Set the concurrency limit to 1, ensuring that only one task runs at a time

	// Launch 10 tasks in the pool
	for taskID := range 10 {
		pool.Go(func() Result[int] {
			// Simulate an error for task ID 2
			if taskID == 2 {
				return Err[int](errors.New("case 2"))
			}

			// Cancel the pool when task ID 7 is reached
			if taskID == 7 {
				pool.Cancel()
				return Err[int](errors.New("case 7"))
			}

			// Return the square of the task ID as the result
			return Ok(taskID * taskID)
		})
	}

	// Wait for all tasks to complete and print the results
	pool.Wait().Print()
}
