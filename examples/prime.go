package main

import (
	"runtime"
	"time"

	. "github.com/enetx/g"
)

func main() {
	start := time.Now()

	limit := 10000
	numCPU := Int(runtime.NumCPU())

	sl := NewSlice[int](0, Int(limit))
	for i := range limit {
		sl.Push(i)
	}

	total := sl.Iter().Parallel(numCPU).Filter(isPrime).Count()
	Println("Computed prime({1}) = {2} in {3.Seconds} seconds.", limit, total, time.Since(start))
}

func isPrime(n int) bool {
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}
