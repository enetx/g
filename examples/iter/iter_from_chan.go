package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 5; i++ {
			ch <- i
		}
	}()

	// Convert the channel into an iterator and apply filtering and mapping operations.
	FromChan(ch).
		Filter(f.Even).                        // Filter even numbers
		Map(func(i int) int { return i * 2 }). // Double each element
		Collect().                             // Collect the results into a slice
		Print()                                // Print the collected results: Slice[4, 8]
}
