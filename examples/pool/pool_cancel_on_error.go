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
			if taskID == 2 {
				return Err[int](errors.New("case 2"))
			}

			return Ok(taskID * taskID)
		})
	}

	pool.Wait().Print()

	if cause := pool.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
