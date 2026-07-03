package main

import (
	"errors"
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func main() {
	errCritical := errors.New("critical task")

	p := pool.New[int]().
		Limit(1) // Single worker — tasks execute sequentially in FIFO order

	// Metrics goroutine prints live pool stats every 100ms.
	// It exits when the pool's context is canceled (via Cancel or ErrAllTasksDone).
	metricsDone := make(chan Unit)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
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

	// Tasks only produce results — they don't control pool lifecycle.
	// The consumer decides when to stop.
	ch := p.Stream(func() {
		for taskID := range 10 {
			p.Go(func() Result[int] {
				time.Sleep(time.Duration(Int(500).RandomRange(1000)) * time.Millisecond)

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

	// Consumer reads results and decides when to stop.
	// Cancel here is safe — the result is already in our hands,
	// nothing gets dropped.
	var results Slice[Result[int]]

	for r := range ch {
		results.Push(r)

		if r.IsErr() && r.ErrIs(errCritical) {
			p.Cancel(r.Err())
			break
		}
	}

	<-metricsDone // Wait for metrics goroutine to print final stats

	results.Println()

	if cause := p.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
