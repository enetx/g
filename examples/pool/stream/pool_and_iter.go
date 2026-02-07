package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/pool"
)

func main() {
	p := pool.New[int]().Limit(10)

	ch := p.Stream(func() {
		for taskID := range 100 {
			p.Go(func() Result[int] {
				if taskID%10 == 2 {
					return Err[int](fmt.Errorf("task %d failed", taskID))
				}
				return Ok(taskID * taskID)
			})
		}
	})

	// Convert channel to SeqResult and partition into successful/failed
	results := FromResultChan(ch)
	successful, failed := results.Partition()

	// Print statistics
	fmt.Printf("Statistics:\n")
	fmt.Printf("  Success: %d\n", successful.Len())
	fmt.Printf("  Errors:  %d\n", failed.Len())

	// Process and display first 10 successful results
	fmt.Printf("\nFirst 10 results:\n")
	successful.
		Iter().
		Map(func(x int) int { return x * x }).
		SortBy(cmp.Cmp).
		Take(10).
		ForEach(func(x int) {
			fmt.Printf("  %d\n", x)
		})

	// Display errors if any occurred
	if !failed.IsEmpty() {
		fmt.Printf("\nErrors:\n")
		failed.Iter().ForEach(func(err error) {
			fmt.Printf("  %v\n", err)
		})
	}
}
