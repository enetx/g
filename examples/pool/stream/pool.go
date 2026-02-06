package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func main() {
	errCritical := errors.New("critical task")

	// Create a context with a short timeout to demonstrate cancellation behavior.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Microsecond)
	defer cancel()

	p := pool.New[int]()
	p.
		Limit(1).    // Only one task runs at a time (FIFO order guaranteed)
		Context(ctx) // Pool respects the timeout - workers stop when context expires

	// Stream spawns workers and runs the provided function in a separate goroutine.
	// Tasks are submitted via Go and results arrive on the returned channel in real-time.
	// The channel closes automatically when all tasks finish or the context is canceled.
	ch := p.Stream(func() {
		for taskID := range 10 {
			p.Go(func() Result[int] {
				if taskID == 2 {
					return Err[int](errors.New("case 2"))
				}

				if taskID == 7 {
					return Err[int](errCritical)
				}

				return Ok(taskID * taskID)
			})
		}
	})

	// Consume results as they arrive.
	// The consumer decides when to stop — tasks don't need to know about pool lifecycle.
	for r := range ch {
		if r.IsErr() && r.ErrIs(errCritical) {
			p.Cancel(r.Err())
			break
		}

		fmt.Println(r)
	}

	// Cause reports why the pool was canceled:
	// - context.DeadlineExceeded if the timeout fired
	// - context.Canceled if we called p.Cancel()
	// - ErrAllTasksDone if all tasks completed normally
	if cause := p.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
