package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

type Response struct {
	ID       int
	Data     string
	Status   int
	Duration time.Duration
}

func callAPI(id int) (Response, error) {
	// Simulate network delay
	delay := time.Duration(50+rand.Intn(150)) * time.Millisecond
	time.Sleep(delay)

	// Simulate different scenarios based on ID
	switch {
	case id%20 == 0:
		// 5% chance of 500 error (will trigger cancellation)
		return Response{}, &APIError{Code: 500, Message: "Internal Server Error"}
	case id%15 == 0:
		// ~7% chance of 503 error (will trigger cancellation)
		return Response{}, &APIError{Code: 503, Message: "Service Unavailable"}
	case id%10 == 0:
		// 10% chance of 404 error (won't trigger cancellation)
		return Response{}, &APIError{Code: 404, Message: "Not Found"}
	case id%8 == 0:
		// ~12% chance of 429 error (won't trigger cancellation)
		return Response{}, &APIError{Code: 429, Message: "Too Many Requests"}
	default:
		// Success
		return Response{
			ID:       id,
			Data:     fmt.Sprintf("Response data for request %d", id),
			Status:   200,
			Duration: delay,
		}, nil
	}
}

func saveToDatabase(resp Response) {
	log.Printf("[V] Saved to DB: ID=%d, Status=%d, Duration=%v", resp.ID, resp.Status, resp.Duration)
}

func main() {
	log.Println("Starting API processing with pool (Wait mode)...")
	log.Println("Configuration:")
	log.Println("  - Workers: 10")
	log.Println("  - Mode: Wait (blocks until all tasks complete)")
	log.Println("  - Cancel on 5xx errors")
	log.Println()

	startTime := time.Now()

	p := pool.New[Response]().
		Limit(10).
		CancelOn(func(r Result[Response]) bool {
			if !r.IsErr() {
				return false
			}

			// Cancel on 5xx API errors
			var apiErr *APIError
			if r.ErrAs(&apiErr) && apiErr.Code >= 500 {
				log.Printf("[X] Critical error detected (Code %d) - cancelling all tasks", apiErr.Code)
				return true
			}

			// Don't cancel on 4xx errors (client errors)
			return false
		})

	// Submit all tasks
	log.Println("Submitting 100 tasks...")
	for id := range 100 {
		p.Go(func() Result[Response] { return ResultOf(callAPI(id)) })
	}

	// Wait for all tasks to complete (or cancellation)
	log.Println("Waiting for tasks to complete...")
	results := p.Wait()

	successful, failed := results.Partition()

	elapsed := time.Since(startTime)

	// Print summary
	log.Println()
	log.Println("========================================")
	log.Printf("Processing completed in %v", elapsed)
	log.Printf("Total tasks submitted: 100")
	log.Printf("Tasks actually executed: %d", p.TotalTasks())
	log.Printf("Active tasks (should be 0): %d", p.ActiveTasks())
	log.Printf("Successful: %d", successful.Len())
	log.Printf("Failed: %d", failed.Len())
	log.Println("========================================")
	log.Println()

	// Process successful responses
	if successful.Len() > 0 {
		log.Println("Saving successful responses to database...")
		successful.Iter().ForEach(saveToDatabase)
	} else {
		log.Println("No successful responses to save.")
	}

	log.Println()
	log.Println("Error summary:")

	// Categorize errors
	var (
		serverErrors int
		clientErrors int
		otherErrors  int
	)

	failed.Iter().ForEach(func(err error) {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			if apiErr.Code >= 500 {
				serverErrors++
				log.Printf("  [5xx] API error %d: %s", apiErr.Code, apiErr.Message)
			} else if apiErr.Code >= 400 {
				clientErrors++
				log.Printf("  [4xx] API error %d: %s", apiErr.Code, apiErr.Message)
			}
		} else {
			otherErrors++
			log.Printf("  [OTHER] Unknown error: %v", err)
		}
	})

	log.Println()
	log.Println("Error breakdown:")
	log.Printf("  - Server errors (5xx): %d", serverErrors)
	log.Printf("  - Client errors (4xx): %d", clientErrors)
	log.Printf("  - Other errors: %d", otherErrors)

	// Check cancellation reason
	if cause := p.Cause(); cause != nil && cause != pool.ErrAllTasksDone {
		log.Println()
		log.Printf("[X] Pool was cancelled: %v", cause)
		log.Printf("Note: Tasks submitted after cancellation were not executed")
	}

	log.Println()
	log.Println("Done!")
}
