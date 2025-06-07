package main

import (
	"errors"
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func main() {
	p := pool.New[int]().Limit(1).CancelOnError()

	for taskID := range 10 {
		p.Go(func() Result[int] {
			if taskID == 4 {
				return Err[int](errors.New("cancel on error"))
			}

			return Ok(taskID * taskID)
		})
	}

	p.Wait().Println()

	if cause := p.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
