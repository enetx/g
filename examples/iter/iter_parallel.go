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
	FromChan(nums.ToChan()).
		Parallel(1000).
		ForEach(func(_ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Chan ForEach: {}", time.Since(start))

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

	start = time.Now()
	flatten := Slice[any]{
		[]int{1, 2, 3},
		[]int{4, 5, 6},
		[]int{7, 8, 9},
	}

	_ = flatten.Iter().Parallel(1000).
		Inspect(func(any) {
			time.Sleep(100 * time.Millisecond)
		}).
		Flatten().
		Collect()

	Println("Elapsed Flatten: {}", time.Since(start))

	// Elapsed ForEach: 101.397041ms
	// Elapsed Chan ForEach: 101.843583ms
	// Elapsed Take: 101.556208ms
	// Elapsed Skip: 100.986875ms
	// Elapsed Range: 101.816459ms
	// Elapsed Map: 102.723292ms
	// Elapsed Partition: 101.915125ms
	// Elapsed Flatten: 102.219167ms
}
