package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	nums := SliceOf(1, 2, 3, 4, 5, 6)

	start := time.Now()

	result := nums.
		Iter().Parallel(nums.Len()).
		Map(func(v int) int {
			time.Sleep(100 * time.Millisecond)
			return v * 2
		}).
		Filter(f.Ne(4)).
		Collect()

	result.Println()
	Println("Elapsed time: {}\n", time.Since(start))
}
