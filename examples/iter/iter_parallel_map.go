package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	m := NewMap[int, int]()
	for i := range 1000 {
		m.Set(i, i)
	}

	start := time.Now()
	m.Iter().Parallel(1000).
		ForEach(func(_, _ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed ForEach: {}", time.Since(start))

	start = time.Now()
	m.Iter().Parallel(1000).
		Range(func(_, _ int) bool {
			time.Sleep(100 * time.Millisecond)
			return true
		})

	Println("Elapsed Range: {}", time.Since(start))

	start = time.Now()
	m.Iter().Parallel(1000).
		Take(500).
		ForEach(func(_, _ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Take: {}", time.Since(start))

	start = time.Now()
	m.Iter().Parallel(1000).
		Skip(500).
		ForEach(func(_, _ int) {
			time.Sleep(100 * time.Millisecond)
		})

	Println("Elapsed Skip: {}", time.Since(start))

	start = time.Now()
	_ = m.Iter().Parallel(1000).
		Map(func(k, v int) (int, int) {
			time.Sleep(100 * time.Millisecond)
			return k, v * 2
		}).
		Collect()

	Println("Elapsed Map: {}", time.Since(start))

	start = time.Now()
	_ = m.Iter().Parallel(1000).
		Filter(func(k, _ int) bool {
			time.Sleep(100 * time.Millisecond)
			return k%2 == 0
		}).
		Collect()

	Println("Elapsed Filter: {}", time.Since(start))

	// Elapsed ForEach: 101.446958ms
	// Elapsed Range: 101.781209ms
	// Elapsed Take: 101.247833ms
	// Elapsed Skip: 101.68375ms
	// Elapsed Map: 101.378625ms
	// Elapsed Filter: 101.039584ms
}
