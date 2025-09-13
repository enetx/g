package main

import (
	"errors"
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func main() {
	// Create a new p for managing tasks
	p := pool.New[int]().
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
					p.TotalTasks(), p.ActiveTasks(), p.FailedTasks())
			case <-p.GetContext().Done():
				fmt.Printf("\r\033[2KAll tasks completed. Total: %d, Failed: %d\n",
					p.TotalTasks(), p.FailedTasks())
				return
			}
		}
	}()

	// Launch 10 tasks in the pool
	for taskID := range 10 {
		p.Go(func() Result[int] {
			// Simulate task execution with a random delay
			time.Sleep(time.Duration(Int(500).RandomRange(1000)) * time.Millisecond)

			// Simulate an error for task ID 2
			if taskID == 2 {
				return Err[int](errors.New("case 2"))
			}

			// Cancel the pool when task ID 7 is reached
			if taskID == 7 {
				p.Cancel()
				return Err[int](errors.New("case 7"))
			}

			// Return the square of the task ID as the result
			return Ok(taskID * taskID)
		})
	}

	results := p.Wait().Collect() // Wait for all tasks to complete
	<-metricsDone                 // Wait for the metrics goroutine to finish
	results.Println()             // Print the results

	if cause := p.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
