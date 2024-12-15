package main

import (
	"context"
	"errors"
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Microsecond)

	pool := NewPool[int]().Limit(1).Context(ctx)

	for taskID := range 10 {
		pool.Go(func() Result[int] {
			if taskID == 2 {
				return Err[int](errors.New("case 2"))
			}

			if taskID == 7 {
				// pool.Cancel()
				pool.Cancel(errors.New("case 7, cancel"))
			}

			return Ok(taskID * taskID)
		})
	}

	pool.Wait().Iter().
		ForEach(func(v Result[int]) {
			if v.IsErr() {
				var contextErr error
				if errors.As(v.Err(), &contextErr) {
					switch contextErr {
					case context.DeadlineExceeded:
						fmt.Println("Error: Context deadline exceeded")
					case context.Canceled:
						fmt.Println("Error: Context was canceled")
					default:
						fmt.Println("Error:", contextErr)
					}
					return
				}
			}
			fmt.Println(v)
		})
}
