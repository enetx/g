package main

import (
	"errors"
	"fmt"
	"time"

	. "github.com/enetx/g"
)

func main() {
	// Create a new pool for managing tasks
	pool := NewPool[int]().
		Limit(1) // Set the concurrency limit to 1, ensuring that only one task runs at a time

	metricsDone := make(chan struct{}) // Channel to synchronize the completion of the metrics goroutine

	// Goroutine to print live metrics about the pool's state
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond) // Timer to periodically update metrics

		defer func() {
			ticker.Stop()
			close(metricsDone)
		}()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("\r\033[2K[Metrics] Total: %d, Active: %d, Failed: %d",
					pool.TotalTasks(), pool.ActiveTasks(), pool.FailedTasks())
			case <-pool.GetContext().Done():
				fmt.Printf("\r\033[2KAll tasks completed. Total: %d, Failed: %d\n",
					pool.TotalTasks(), pool.FailedTasks())
				return
			}
		}
	}()

	// Launch 10 tasks in the pool
	for taskID := range 10 {
		pool.Go(func() Result[int] {
			// Simulate task execution with a random delay
			time.Sleep(time.Duration(Int(500).RandomRange(1000)) * time.Millisecond)

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

	results := pool.Wait() // Wait for all tasks to complete
	<-metricsDone          // Wait for the metrics goroutine to finish
	results.Println()      // Print the results

	if cause := pool.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
