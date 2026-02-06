package main

import (
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

// ╔══════════════════════════════════════════════════════════════════════════════╗
// ║                           stdlib vs g/pool                                 ║
// ║                      Same problems. Less boilerplate.                      ║
// ╚══════════════════════════════════════════════════════════════════════════════╝

func main() {
	parallel()
	cancelOnError()
	stream()
	streamFromChannel()
	streamEarlyStop()
}

// ════════════════════════════════════════════════════════════════════════════════
// Example 1: Parallel tasks with concurrency limit
// ════════════════════════════════════════════════════════════════════════════════

func parallel() {
	fmt.Println("═══ Example 1: Parallel with Limit ═══")

	urls := []string{
		"https://api.example.com/users",
		"https://api.example.com/posts",
		"https://api.example.com/comments",
		"https://api.example.com/albums",
		"https://api.example.com/photos",
		"https://api.example.com/todos",
	}

	// ── stdlib ──────────────────────────────────────────────────────────────
	//
	//  var (
	//      mu      sync.Mutex
	//      wg      sync.WaitGroup
	//      results []string
	//      errs    []error
	//  )
	//
	//  sem := make(chan struct{}, 3)
	//
	//  for _, url := range urls {
	//      wg.Add(1)
	//      sem <- struct{}{}
	//      go func() {
	//          defer wg.Done()
	//          defer func() { <-sem }()
	//
	//          body, err := fetch(url)
	//          mu.Lock()
	//          defer mu.Unlock()
	//          if err != nil {
	//              errs = append(errs, err)
	//          } else {
	//              results = append(results, body)
	//          }
	//      }()
	//  }
	//
	//  wg.Wait()

	// ── g/pool ──────────────────────────────────────────────────────────────

	p := pool.New[String]().Limit(3)

	for _, url := range urls {
		p.Go(func() Result[String] {
			body, err := fetch(url)
			if err != nil {
				return Err[String](err)
			}
			return Ok(String(body))
		})
	}

	p.Wait().Collect().Println()
}

// ════════════════════════════════════════════════════════════════════════════════
// Example 2: Cancel everything on first error
// ════════════════════════════════════════════════════════════════════════════════

func cancelOnError() {
	fmt.Println("\n═══ Example 2: Cancel on Error ═══")

	// ── stdlib ──────────────────────────────────────────────────────────────
	//
	//  ctx, cancel := context.WithCancel(context.Background())
	//  defer cancel()
	//
	//  var (
	//      mu       sync.Mutex
	//      wg       sync.WaitGroup
	//      results  []int
	//      firstErr error
	//      errOnce  sync.Once
	//  )
	//
	//  sem := make(chan struct{}, 3)
	//
	//  for i := range 10 {
	//      select {
	//      case <-ctx.Done():
	//          break
	//      case sem <- struct{}{}:
	//      }
	//
	//      wg.Add(1)
	//      go func() {
	//          defer wg.Done()
	//          defer func() { <-sem }()
	//
	//          if ctx.Err() != nil {
	//              return
	//          }
	//
	//          result, err := process(i)
	//          if err != nil {
	//              errOnce.Do(func() {
	//                  firstErr = err
	//                  cancel()
	//              })
	//              return
	//          }
	//
	//          mu.Lock()
	//          results = append(results, result)
	//          mu.Unlock()
	//      }()
	//  }
	//
	//  wg.Wait()

	// ── g/pool ──────────────────────────────────────────────────────────────

	p := pool.New[Int]().Limit(3).CancelOnError()

	for i := range 10 {
		p.Go(func() Result[Int] {
			result, err := process(i)
			if err != nil {
				return Err[Int](err)
			}
			return Ok(Int(result))
		})
	}

	p.Wait().Collect().Println()

	if cause := p.Cause(); cause != nil {
		fmt.Println("stopped:", cause)
	}
}

// ════════════════════════════════════════════════════════════════════════════════
// Example 3: Stream results in real-time
// ════════════════════════════════════════════════════════════════════════════════

func stream() {
	fmt.Println("\n═══ Example 3: Stream ═══")

	// ── stdlib ──────────────────────────────────────────────────────────────
	//
	//  results := make(chan int, 10)
	//  sem := make(chan struct{}, 5)
	//  var wg sync.WaitGroup
	//
	//  go func() {
	//      for i := range 20 {
	//          wg.Add(1)
	//          sem <- struct{}{}
	//          go func() {
	//              defer wg.Done()
	//              defer func() { <-sem }()
	//              results <- i * i
	//          }()
	//      }
	//      wg.Wait()
	//      close(results)
	//  }()
	//
	//  for r := range results {
	//      fmt.Print(r, " ")
	//  }

	// ── g/pool ──────────────────────────────────────────────────────────────

	p := pool.New[Int]().Limit(5)

	ch := p.Stream(func() {
		for i := range 20 {
			p.Go(func() Result[Int] { return Ok(Int(i * i)) })
		}
	})

	for r := range ch {
		fmt.Print(r.Ok(), " ")
	}

	fmt.Println()
}

// ════════════════════════════════════════════════════════════════════════════════
// Example 4: Stream from an input channel
// ════════════════════════════════════════════════════════════════════════════════

func streamFromChannel() {
	fmt.Println("\n═══ Example 4: Stream from Channel ═══")

	// simulate an external source producing work
	jobs := make(chan string)

	go func() {
		defer close(jobs)

		for _, name := range []string{"alice", "bob", "charlie", "dave", "eve"} {
			jobs <- name
		}
	}()

	// ── stdlib ──────────────────────────────────────────────────────────────
	//
	//  results := make(chan string, 5)
	//  sem := make(chan struct{}, 2)
	//  var wg sync.WaitGroup
	//
	//  go func() {
	//      for name := range jobs {
	//          wg.Add(1)
	//          sem <- struct{}{}
	//          go func() {
	//              defer wg.Done()
	//              defer func() { <-sem }()
	//              results <- "Hello, " + name + "!"
	//          }()
	//      }
	//      wg.Wait()
	//      close(results)
	//  }()
	//
	//  for r := range results {
	//      fmt.Println(r)
	//  }

	// ── g/pool ──────────────────────────────────────────────────────────────

	p := pool.New[String]().Limit(2)

	ch := p.Stream(func() {
		for name := range jobs {
			p.Go(func() Result[String] {
				return Ok(String("Hello, " + name + "!"))
			})
		}
	})

	for r := range ch {
		fmt.Println(r.Ok())
	}
}

// ════════════════════════════════════════════════════════════════════════════════
// Example 5: Stream with early stop from consumer
// ════════════════════════════════════════════════════════════════════════════════

func streamEarlyStop() {
	fmt.Println("\n═══ Example 5: Early Stop ═══")

	// ── stdlib ──────────────────────────────────────────────────────────────
	//
	//  ctx, cancel := context.WithCancel(context.Background())
	//  defer cancel()
	//
	//  results := make(chan int, 5)
	//  sem := make(chan struct{}, 5)
	//  var wg sync.WaitGroup
	//
	//  go func() {
	//      for i := range 1_000_000 {
	//          select {
	//          case <-ctx.Done():
	//              break
	//          case sem <- struct{}{}:
	//          }
	//
	//          wg.Add(1)
	//          go func() {
	//              defer wg.Done()
	//              defer func() { <-sem }()
	//              if ctx.Err() != nil {
	//                  return
	//              }
	//              results <- i * i
	//          }()
	//      }
	//      wg.Wait()
	//      close(results)
	//  }()
	//
	//  for r := range results {
	//      if r > 100 {
	//          cancel()
	//          break
	//      }
	//      fmt.Print(r, " ")
	//  }

	// ── g/pool ──────────────────────────────────────────────────────────────

	p := pool.New[Int]().Limit(5)

	ch := p.Stream(func() {
		for i := range 1_000_000 {
			p.Go(func() Result[Int] { return Ok(Int(i * i)) })
		}
	})

	for r := range ch {
		if r.Ok() > 100 {
			p.Cancel()
			break
		}
		fmt.Print(r.Ok(), " ")
	}

	fmt.Println("\nstopped early")
}

// ════════════════════════════════════════════════════════════════════════════════
// Helpers
// ════════════════════════════════════════════════════════════════════════════════

func fetch(url string) (string, error) {
	time.Sleep(10 * time.Millisecond) // simulate network
	return fmt.Sprintf("data from %s", url), nil
}

func process(i int) (int, error) {
	time.Sleep(5 * time.Millisecond) // simulate work
	if i == 5 {
		return 0, fmt.Errorf("task %d failed", i)
	}
	return i * i, nil
}
