package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	. "github.com/enetx/g"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Microsecond)
	defer cancel()

	pool := NewPool[int]().Limit(10).Context(ctx)

	for taskID := range 10 {
		pool.Go(func() Result[int] {
			if taskID == 2 {
				return Err[int](errors.New("case 2"))
			}

			if taskID == 7 {
				pool.Cancel(errors.New("case 7, cancel"))
			}

			return Ok(taskID * taskID)
		})
	}

	pool.Wait().Println()

	if cause := pool.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
