package main

import (
	"context"
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func main() {
	exampleRateWithCancelOnError()
	fmt.Println()
	exampleRateWithTimeout()
	fmt.Println()
	exampleAPIRateLimit()
	fmt.Println()
	exampleSmoothVsBurst()
}

// exampleRateWithCancelOnError — rate limited pool that stops on first error.
func exampleRateWithCancelOnError() {
	fmt.Println("=== Rate + CancelOnError (Wait) ===")

	p := pool.New[int]().
		Limit(3).
		Rate(5, time.Second). // 5 tasks/sec
		CancelOnError()

	for i := range 20 {
		p.Go(func() Result[int] {
			if i == 7 {
				return Err[int](fmt.Errorf("task %d failed", i))
			}
			time.Sleep(50 * time.Millisecond)
			return Ok(i * i)
		})
	}

	results := p.Wait().Collect()
	fmt.Printf("Got %d results (some tasks were canceled)\n", results.Len())

	for _, r := range results {
		fmt.Println(" ", r)
	}

	if cause := p.Cause(); cause != nil {
		fmt.Println("Canceled due to:", cause)
	}
}

// exampleRateWithTimeout — rate limited pool with context deadline.
func exampleRateWithTimeout() {
	fmt.Println("=== Rate + Context Timeout (Stream) ===")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	p := pool.New[String]().
		Limit(2).
		Rate(3, time.Second, 1) // smooth 3/sec

	p.Context(ctx)

	count := 0
	ch := p.Stream(func() {
		for i := range 100 { // submit many, but timeout will cut it short
			p.Go(func() Result[String] {
				time.Sleep(100 * time.Millisecond)
				return Ok(Format("task {} done", i))
			})
		}
	})

	for r := range ch {
		fmt.Println(" ", r.Ok())
		count++
	}

	fmt.Printf("Completed %d tasks before timeout\n", count)
	fmt.Println("Cause:", p.Cause())
}

// exampleAPIRateLimit simulates calling an API with strict rate limits.
// 10 requests/sec with burst of 3, max 5 concurrent connections.
func exampleAPIRateLimit() {
	fmt.Println("=== API Rate Limiting Simulation (Stream) ===")

	start := time.Now()

	p := pool.New[String]().
		Limit(5).                 // max 5 concurrent "connections"
		Rate(10, time.Second, 3). // API allows 10 req/sec, burst 3
		CancelOnError()

	urls := NewSlice[String](30)
	for i := range urls {
		urls[i] = Format("https://api.example.com/items/{}", i+1)
	}

	ch := p.Stream(func() {
		for _, url := range urls {
			p.Go(func() Result[String] {
				// simulate API call
				time.Sleep(50 * time.Millisecond)

				// simulate occasional 429 Too Many Requests
				// (shouldn't happen with proper rate limiting)
				return Ok(Format("{} -> 200 OK", url))
			})
		}
	})

	var results Slice[String]

	for r := range ch {
		if r.IsOk() {
			results.Push(r.Ok())
		} else {
			fmt.Println("  ERROR:", r.Err())
		}
	}

	elapsed := time.Since(start).Truncate(time.Millisecond)
	fmt.Printf("Fetched %d URLs in %s\n", len(results), elapsed)
	fmt.Printf("Effective rate: %.1f req/sec\n", float64(len(results))/time.Since(start).Seconds())
}

// exampleSmoothVsBurst shows the difference between burst and smooth rate limiting.
func exampleSmoothVsBurst() {
	for _, mode := range []struct {
		name  string
		burst int
	}{
		{"burst (default)", 5},
		{"smooth (burst=1)", 1},
	} {
		fmt.Printf("\n=== %s ===\n", mode.name)
		start := time.Now()

		p := pool.New[int]().
			Limit(10).
			Rate(5, time.Second, mode.burst)

		ch := p.Stream(func() {
			for i := range 15 {
				p.Go(func() Result[int] {
					elapsed := time.Since(start).Truncate(time.Millisecond)
					fmt.Printf("  task %02d started at %s\n", i, elapsed)
					return Ok(i)
				})
			}
		})

		// drain
		for range ch {
		}

		fmt.Printf("  total: %s\n", time.Since(start).Truncate(time.Millisecond))
	}
}
