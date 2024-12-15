package main

import (
	"errors"
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pkg/ref"
)

func main() {
	pool := NewPool[int]().Limit(2)

	for taskID := range 10 {
		pool.Go(func() Result[int] {
			time.Sleep(time.Duration(Int(500).RandomRange(1000)) * time.Millisecond)

			if taskID == 2 {
				return Err[int](errors.New("case 2"))
			}

			if taskID == 7 {
				pool.Cancel()
				return Err[int](errors.New("case 7"))
			}

			return Ok(taskID * taskID)
		})
	}

	pool.Wait().Iter().
		ForEach(func(v Result[int]) {
			if v.IsErr() {
				if errors.As(v.Err(), ref.Of(&ErrorContext{})) {
					fmt.Println("Context Error:", v.Err())
					return
				}

				fmt.Println("Error:", v.Err())
				return
			}

			fmt.Println(v)
		})
}
