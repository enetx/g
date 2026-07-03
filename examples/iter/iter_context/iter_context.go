package main

import (
	"context"
	"fmt"
	"time"

	. "github.com/enetx/g"
)

// Run all context examples
func main() {
	BasicContextCancellation()
	TimeoutProcessing()
	UserCancellation()
	PipelineWithContext()
	GracefulShutdown()
	ContextWithChannels()
	NestedContextOperations()
}

// Example 1: Basic Context Cancellation
func BasicContextCancellation() {
	fmt.Println("=== Basic Context Cancellation ===")

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Create a large sequence
	numbers := Range(1, 1000000).Context(ctx)

	// Cancel after processing a few elements
	go func() {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Cancelling context...")
		cancel()
	}()

	// Process elements until context is cancelled
	count := 0
	numbers.ForEach(func(int) {
		count++
		if count%1000 == 0 {
			fmt.Printf("Processed %d elements\n", count)
		}
		time.Sleep(time.Microsecond) // Simulate work
	})

	fmt.Printf("Total processed: %d elements\n", count)
	// Output: Much less than 1,000,000 due to cancellation
}

// Example 2: Timeout-based Processing
func TimeoutProcessing() {
	fmt.Println("\n=== Timeout-based Processing ===")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Create infinite sequence with context
	numbers := Range(1, 10000).
		Cycle(). // Infinite repetition
		Context(ctx)

	// Process until timeout
	processed := numbers.
		Map(func(n int) int {
			time.Sleep(time.Millisecond) // Simulate slow processing
			return n * 2
		}).
		Take(50). // Try to take 50, but timeout will stop us earlier
		Collect()

	fmt.Printf("Processed %d numbers before timeout: %v\n", len(processed), processed[:min(10, len(processed))])
}

// Example 3: User Cancellation in Long-running Operation
func UserCancellation() {
	fmt.Println("\n=== User Cancellation Example ===")

	ctx, cancel := context.WithCancel(context.Background())

	// Simulate user pressing Ctrl+C after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("User pressed Ctrl+C - cancelling operation...")
		cancel()
	}()

	// Expensive computation on large dataset
	result := Range(1, 1000000).
		Context(ctx).
		Filter(func(n int) bool {
			time.Sleep(time.Microsecond) // Simulate expensive filtering
			return n%1000 == 0
		}).
		Map(func(n int) int {
			time.Sleep(time.Microsecond) // Simulate expensive transformation
			return n * n
		}).
		Collect()

	fmt.Printf("Computation result (length: %d): %v...\n",
		len(result), result[:min(5, len(result))])
}

// Example 4: Pipeline with Multiple Context Checkpoints
func PipelineWithContext() {
	fmt.Println("\n=== Pipeline with Context Checkpoints ===")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Complex processing pipeline
	result := Range(1, 10000).
		Context(ctx). // Check context at input
		Filter(func(n int) bool {
			time.Sleep(time.Microsecond)
			return n%2 == 0
		}).
		Context(ctx). // Check context after filtering
		Map(func(n int) int {
			time.Sleep(time.Microsecond)
			return n * 3
		}).
		Context(ctx). // Check context after mapping
		Take(100).
		Collect()

	fmt.Printf("Pipeline result (length: %d): %v...\n",
		len(result), result[:min(5, len(result))])
}

// Example 5: Graceful Shutdown with Context
func GracefulShutdown() {
	fmt.Println("\n=== Graceful Shutdown Example ===")

	// Simulate server shutdown signal
	ctx, cancel := context.WithCancel(context.Background())

	// Start background processing
	go func() {
		time.Sleep(1500 * time.Millisecond)
		fmt.Println("Received shutdown signal...")
		cancel()
	}()

	// Process work items until shutdown
	workItems := Range(1, 100).
		Cycle(). // Infinite work
		Context(ctx)

	processedCount := 0
	workItems.ForEach(func(item int) {
		// Simulate work processing
		time.Sleep(50 * time.Millisecond)
		processedCount++
		fmt.Printf("Processed work item: %d\n", item)
	})

	fmt.Printf("Gracefully processed %d work items before shutdown\n", processedCount)
}

// Example 6: Context with Channel Integration
func ContextWithChannels() {
	fmt.Println("\n=== Context with Channels ===")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Create a channel with data
	ch := make(chan int, 100)
	go func() {
		defer close(ch)
		for i := 1; i <= 1000; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Channel producer cancelled")
				return
			case ch <- i:
				time.Sleep(time.Millisecond)
			}
		}
	}()

	// Process channel data with context
	result := FromChan(ch).
		Context(ctx).
		Filter(func(n int) bool { return n%10 == 0 }).
		Map(func(n int) int { return n * n }).
		Collect()

	fmt.Printf("Channel processing result: %v\n", result)
}

// Example 7: Context Propagation in Nested Operations
func NestedContextOperations() {
	fmt.Println("\n=== Nested Context Operations ===")

	// Parent context with timeout
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer parentCancel()

	// Child context that can be cancelled independently
	childCtx, childCancel := context.WithCancel(parentCtx)

	// Cancel child context early
	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Cancelling child context...")
		childCancel()
	}()

	// Nested processing with different contexts
	outerResult := Range(1, 100).
		Context(parentCtx).
		Map(func(n int) int {
			// Inner processing with child context
			innerResult := Range(n, n+10).
				Context(childCtx).
				Filter(func(x int) bool { return x%2 == 0 }).
				Count()

			time.Sleep(10 * time.Millisecond) // Simulate work
			return int(innerResult)
		}).
		Take(20).
		Collect()

	fmt.Printf("Nested processing result: %v\n", outerResult)
}
