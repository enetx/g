package main

import (
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func main() {
	exampleWait()
	fmt.Println()
	exampleStream()
}

// exampleWait demonstrates Rate with Wait mode.
// 10 tasks, max 3 concurrent, max 2 starts per second.
func exampleWait() {
	fmt.Println("=== Wait mode with Rate ===")

	start := time.Now()

	p := pool.New[String]().
		Limit(3).               // max 3 goroutines at once
		Rate(2, time.Second, 2) // 2 tasks/sec, burst of 2

	for i := range 10 {
		p.Go(func() Result[String] {
			elapsed := time.Since(start).Truncate(time.Millisecond)
			msg := Format("task {} started at {}", i, elapsed)
			time.Sleep(200 * time.Millisecond) // simulate work
			return Ok(msg)
		})
	}

	for r := range p.Wait() {
		fmt.Println(" ", r.Ok())
	}

	fmt.Printf("Total time: %s\n", time.Since(start).Truncate(time.Millisecond))
	fmt.Printf("Tasks: %d total, %d successful, %d failed\n",
		p.TotalTasks(), p.SuccessfulTasks(), p.FailedTasks())
}

// exampleStream demonstrates Rate with Stream mode.
// 20 tasks, 5 workers, max 3 executions per second — smooth, no burst.
func exampleStream() {
	fmt.Println("=== Stream mode with Rate ===")

	start := time.Now()

	p := pool.New[String]().
		Limit(5).               // 5 worker goroutines
		Rate(3, time.Second, 1) // 3 tasks/sec, burst of 1 (smooth)

	ch := p.Stream(func() {
		for i := range 20 {
			p.Go(func() Result[String] {
				elapsed := time.Since(start).Truncate(time.Millisecond)
				msg := Format("task {} executed at {}", i, elapsed)
				time.Sleep(100 * time.Millisecond) // simulate fast work
				return Ok(msg)
			})
		}
	})

	for r := range ch {
		fmt.Println(" ", r.Ok())
	}

	fmt.Printf("Total time: %s\n", time.Since(start).Truncate(time.Millisecond))
	fmt.Printf("Tasks: %d total, %d successful, %d failed\n",
		p.TotalTasks(), p.SuccessfulTasks(), p.FailedTasks())
}
