package main

import (
	"errors"
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	pool := NewPool[int]().Limit(1).CancelOnError()

	for taskID := range 10 {
		pool.Go(func() Result[int] {
			if taskID == 4 {
				return Err[int](errors.New("cancel on error"))
			}

			return Ok(taskID * taskID)
		})
	}

	pool.Wait().Println()

	if cause := pool.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
