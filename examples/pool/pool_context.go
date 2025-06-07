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
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Microsecond)
	defer cancel()

	p := pool.New[int]().Limit(10).Context(ctx)

	for taskID := range 10 {
		p.Go(func() Result[int] {
			if taskID == 2 {
				return Err[int](errors.New("case 2"))
			}

			if taskID == 7 {
				p.Cancel(errors.New("case 7, cancel"))
			}

			return Ok(taskID * taskID)
		})
	}

	p.Wait().Println()

	if cause := p.Cause(); cause != nil {
		fmt.Println("Pool was canceled due to:", cause)
	}
}
