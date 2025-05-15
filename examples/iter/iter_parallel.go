package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	nums := Range(0, 1000)

	start := time.Now()
	nums.Parallel(1000).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed ForEach: {}", time.Since(start))

	start = time.Now()
	nums.Parallel(1000).
		Take(500).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Take: {}", time.Since(start))

	start = time.Now()
	nums.Parallel(1000).
		Skip(500).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Skip: {}", time.Since(start))

	start = time.Now()
	nums.Parallel(1000).
		Range(func(_ int) bool {
			time.Sleep(100 * time.Millisecond)
			return true
		})

	Println("Elapsed Range: {}", time.Since(start))

	start = time.Now()
	_ = nums.Parallel(1000).
		Map(func(v int) int {
			time.Sleep(100 * time.Millisecond)
			return v * 2
		}).
		Filter(f.Ne(4)).
		Collect()

	Println("Elapsed Map: {}", time.Since(start))

	start = time.Now()
	_, _ = nums.Parallel(1000).
		Partition(func(i int) bool {
			time.Sleep(100 * time.Millisecond)
			return i%2 == 0
		})

	Println("Elapsed Partition: {}", time.Since(start))

	// Elapsed ForEach: 101.193875ms
	// Elapsed Take: 101.294792ms
	// Elapsed Skip: 101.651666ms
	// Elapsed Range: 101.499541ms
	// Elapsed Map: 101.807875ms
	// Elapsed Partition: 101.07325ms
}
