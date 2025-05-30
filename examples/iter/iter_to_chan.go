package main

import (
	"context"
	"fmt"
	"sync"

	. "github.com/enetx/g"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := SliceOf(1, 1, 1, 3, 4, 4, 8, 8, 9, 9).
		Iter().
		// Dedup().
		ToChan(ctx)

	// for job := range jobs {
	// 	fmt.Printf("job: %v\n", job)
	// }

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for job := range jobs {
			fmt.Println(job)
		}
	}()

	wg.Wait()
}
