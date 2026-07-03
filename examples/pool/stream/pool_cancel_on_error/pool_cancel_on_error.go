package main

import (
	"errors"
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func main() {
	p := pool.New[int]().
		Limit(1).       // Single worker — tasks execute in FIFO order
		CancelOnError() // Stop all remaining tasks on first error

	// Stream runs tasks and emits results in real-time.
	// With Limit(1), results arrive in submission order.
	// On error, the triggering error is guaranteed to be delivered;
	// subsequent tasks are skipped.
	ch := p.Stream(func() {
		for taskID := range 10 {
			p.Go(func() Result[int] {
				if taskID == 4 {
					return Err[int](errors.New("cancel on error"))
				}
				return Ok(taskID * taskID)
			})
		}
	})

	// Collect all delivered results.
	// With CancelOnError, this will contain tasks 0..3 (Ok) + task 4 (Err).
	// Tasks 5..9 are never executed.
	FromChan(ch).Collect().Println()

	// Cause returns the error that triggered cancellation.
	if cause := p.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
