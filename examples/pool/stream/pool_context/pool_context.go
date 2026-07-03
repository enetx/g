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

	// Context with a short timeout — pool will stop if tasks take too long.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Microsecond)
	defer cancel()

	p := pool.New[int]().
		Limit(2).    // Two workers process tasks concurrently
		Context(ctx) // Pool respects the timeout — workers exit when context expires

	// Stream distributes tasks across workers and emits results in real-time.
	// With Limit(2), two tasks run in parallel; order is not guaranteed.
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

	// Consumer controls when to stop.
	// Cancel accepts a custom error — it becomes the pool's Cause.
	for r := range ch {
		if r.IsErr() && r.ErrIs(errCritical) {
			p.Cancel(r.Err())
			break
		}

		fmt.Println(r)
	}

	// Cause reports the reason for cancellation:
	// - "case 7, cancel" if consumer stopped the pool
	// - context.DeadlineExceeded if the 100μs timeout fired first
	// - ErrAllTasksDone if all 10 tasks completed before either
	if cause := p.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
