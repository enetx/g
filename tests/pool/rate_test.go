package g_test

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

// --- Rate enforcement ---

func TestRate_Wait_SmoothRate(t *testing.T) {
	// 5 tasks/sec, burst=1 (smooth). 5 tasks: 1 instant + 4 * 200ms = ~800ms.
	p := pool.New[int]().Rate(5, time.Second, 1)

	start := time.Now()
	for i := range 5 {
		p.Go(func() Result[int] { return Ok(i) })
	}
	p.Wait()

	elapsed := time.Since(start)
	if elapsed < 600*time.Millisecond {
		t.Errorf("5 tasks at 5/sec (burst=1) took %s, expected >= 600ms", elapsed)
	}
	if elapsed > 2*time.Second {
		t.Errorf("5 tasks at 5/sec took %s, too slow", elapsed)
	}
}

func TestRate_Wait_BurstThenSteady(t *testing.T) {
	// 5 tasks/sec, burst=3. First 3 instant, next 2 wait ~200ms each.
	// Timestamps recorded AFTER Go() returns — waitRate() blocks inside Go(),
	// so the timestamp reflects when the rate token was actually acquired.
	p := pool.New[int]().Rate(5, time.Second, 3)

	var timestamps [5]time.Duration
	start := time.Now()

	for i := range 5 {
		p.Go(func() Result[int] { return Ok(i) })
		timestamps[i] = time.Since(start)
	}
	p.Wait()

	// First 3 should be near-instant (burst tokens).
	for i := range 3 {
		if timestamps[i] > 100*time.Millisecond {
			t.Errorf("burst task %d returned at %s, expected near-instant", i, timestamps[i])
		}
	}

	// Task 3 should be delayed (first post-burst, ~200ms wait).
	if timestamps[3] < 100*time.Millisecond {
		t.Errorf("post-burst task 3 returned at %s, expected delay >= 100ms", timestamps[3])
	}

	t.Logf("submit timestamps: %v", timestamps)
}

func TestRate_Wait_DefaultBurst(t *testing.T) {
	// Rate(5, time.Second) — default burst equals n=5.
	// All 5 should be near-instant from burst.
	p := pool.New[int]().Rate(5, time.Second)

	start := time.Now()
	for i := range 5 {
		p.Go(func() Result[int] { return Ok(i) })
	}
	p.Wait()

	elapsed := time.Since(start)
	if elapsed > 200*time.Millisecond {
		t.Errorf("5 tasks with burst=5 took %s, expected near-instant", elapsed)
	}
}

func TestRate_Stream_SmoothRate(t *testing.T) {
	// 5 tasks/sec, burst=1 in stream mode.
	p := pool.New[int]().Limit(5).Rate(5, time.Second, 1)

	start := time.Now()
	ch := p.Stream(func() {
		for i := range 5 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	for range ch {
	}

	elapsed := time.Since(start)
	if elapsed < 600*time.Millisecond {
		t.Errorf("5 tasks at 5/sec (burst=1) streamed in %s, expected >= 600ms", elapsed)
	}
}

func TestRate_Stream_BurstThenSteady(t *testing.T) {
	p := pool.New[int]().Limit(5).Rate(5, time.Second, 3)

	var count atomic.Int32
	start := time.Now()

	ch := p.Stream(func() {
		for i := range 5 {
			p.Go(func() Result[int] {
				count.Add(1)
				return Ok(i)
			})
		}
	})

	for range ch {
	}

	elapsed := time.Since(start)
	if count.Load() != 5 {
		t.Errorf("expected 5 results, got %d", count.Load())
	}

	// 3 burst + 2 at 200ms each ≈ 400ms minimum.
	if elapsed < 300*time.Millisecond {
		t.Errorf("expected post-burst delay, completed in %s", elapsed)
	}
}

// --- Rate + Limit composition ---

func TestRate_Wait_WithLimit(t *testing.T) {
	// Rate 20/sec but only 2 concurrent.
	p := pool.New[int]().Limit(2).Rate(20, time.Second, 20)

	var maxConcurrent, current atomic.Int32

	for range 10 {
		p.Go(func() Result[int] {
			cur := current.Add(1)
			for {
				old := maxConcurrent.Load()
				if cur <= old || maxConcurrent.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			current.Add(-1)
			return Ok(0)
		})
	}

	p.Wait()

	if maxConcurrent.Load() > 2 {
		t.Errorf("expected max concurrency <= 2, got %d", maxConcurrent.Load())
	}
}

func TestRate_Stream_WithLimit(t *testing.T) {
	p := pool.New[int]().Limit(3).Rate(50, time.Second, 50)

	var maxConcurrent, current atomic.Int32

	ch := p.Stream(func() {
		for range 20 {
			p.Go(func() Result[int] {
				cur := current.Add(1)
				for {
					old := maxConcurrent.Load()
					if cur <= old || maxConcurrent.CompareAndSwap(old, cur) {
						break
					}
				}
				time.Sleep(10 * time.Millisecond)
				current.Add(-1)
				return Ok(0)
			})
		}
	})

	for range ch {
	}

	if maxConcurrent.Load() > 3 {
		t.Errorf("expected max concurrency <= 3 in stream, got %d", maxConcurrent.Load())
	}
}

func TestRate_Wait_OnlyRate_NoLimit(t *testing.T) {
	// Rate without Limit — unlimited concurrency, rate-limited starts.
	p := pool.New[int]().Rate(20, time.Second, 20)

	for i := range 10 {
		p.Go(func() Result[int] { return Ok(i) })
	}

	results := p.Wait().Collect()
	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}
}

func TestRate_Stream_OnlyRate_NoLimit(t *testing.T) {
	// Stream with Rate but no explicit Limit — uses GOMAXPROCS workers.
	p := pool.New[int]().Rate(50, time.Second, 50)

	ch := p.Stream(func() {
		for i := range 20 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	count := 0
	for range ch {
		count++
	}

	if count != 20 {
		t.Errorf("expected 20 results, got %d", count)
	}
}

// --- Rate + context ---

func TestRate_Wait_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// 1 task/sec, burst=1 — very slow.
	p := pool.New[int]().Rate(1, time.Second, 1).Context(ctx)

	for i := range 10 {
		p.Go(func() Result[int] { return Ok(i) })
	}

	results := p.Wait().Collect()

	if len(results) >= 5 {
		t.Errorf("expected few results before timeout, got %d", len(results))
	}

	t.Logf("got %d results before context timeout", len(results))
}

func TestRate_Stream_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	p := pool.New[int]().Limit(2).Rate(2, time.Second, 1).Context(ctx)

	ch := p.Stream(func() {
		for i := range 50 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	count := 0
	for range ch {
		count++
	}

	if count >= 50 {
		t.Errorf("expected early stop from timeout, got all %d", count)
	}

	t.Logf("stream got %d results before timeout", count)
}

func TestRate_Wait_ContextAlreadyCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	p := pool.New[int]().Rate(10, time.Second).Context(ctx)

	p.Go(func() Result[int] { return Ok(1) })

	results := p.Wait().Collect()

	if len(results) != 0 {
		t.Errorf("expected 0 results with pre-canceled context, got %d", len(results))
	}
}

// --- Rate + CancelOnError ---

func TestRate_Wait_CancelOnError(t *testing.T) {
	p := pool.New[int]().Limit(1).Rate(20, time.Second).CancelOnError()

	p.Go(func() Result[int] { return Ok(1) })
	p.Go(func() Result[int] { return Err[int](errors.New("boom")) })
	p.Go(func() Result[int] { return Ok(3) })

	results := p.Wait().Collect()

	if len(results) >= 3 {
		t.Errorf("expected early cancellation, got %d results", len(results))
	}

	hasErr := false
	for _, r := range results {
		if r.IsErr() {
			hasErr = true
		}
	}
	if !hasErr {
		t.Error("expected error in results")
	}
}

func TestRate_Stream_CancelOnError(t *testing.T) {
	p := pool.New[int]().Limit(1).Rate(20, time.Second).CancelOnError()

	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Err[int](errors.New("stream fail")) })
		p.Go(func() Result[int] { return Ok(2) })
		p.Go(func() Result[int] { return Ok(3) })
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	hasErr := false
	for _, r := range results {
		if r.IsErr() && strings.Contains(r.Err().Error(), "stream fail") {
			hasErr = true
		}
	}
	if !hasErr {
		t.Error("expected 'stream fail' error in results")
	}
}

func TestRate_Stream_CancelOnError_ErrorNotLost(t *testing.T) {
	for range 100 {
		p := pool.New[int]().Limit(1).Rate(100, time.Second, 100).CancelOnError()

		ch := p.Stream(func() {
			p.Go(func() Result[int] { return Err[int](errors.New("must arrive")) })
		}, 1)

		gotErr := false
		for r := range ch {
			if r.IsErr() && strings.Contains(r.Err().Error(), "must arrive") {
				gotErr = true
			}
		}

		if !gotErr {
			t.Fatal("error lost in rate-limited stream with CancelOnError")
		}
	}
}

// --- Rate + CancelOn predicate ---

func TestRate_Wait_CancelOn(t *testing.T) {
	p := pool.New[int]().Limit(1).Rate(20, time.Second, 20).
		CancelOn(func(r Result[int]) bool {
			return r.IsOk() && r.Ok() == 5
		})

	for i := range 100 {
		p.Go(func() Result[int] { return Ok(i) })
	}

	results := p.Wait().Collect()

	if len(results) >= 100 {
		t.Errorf("expected early cancellation, got all %d results", len(results))
	}

	found := false
	for _, r := range results {
		if r.IsOk() && r.Ok() == 5 {
			found = true
		}
	}
	if !found {
		t.Error("trigger value (5) not in results")
	}
}

func TestRate_Stream_CancelOn(t *testing.T) {
	p := pool.New[int]().Limit(1).Rate(20, time.Second, 20).
		CancelOn(func(r Result[int]) bool {
			return r.IsOk() && r.Ok() == 3
		})

	ch := p.Stream(func() {
		for i := range 50 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	if len(results) >= 50 {
		t.Errorf("expected early cancellation, got all %d results", len(results))
	}

	found := false
	for _, r := range results {
		if r.IsOk() && r.Ok() == 3 {
			found = true
		}
	}
	if !found {
		t.Error("trigger value (3) not in stream results")
	}
}

// --- Rate + error handling ---

func TestRate_Wait_NilFunction(t *testing.T) {
	p := pool.New[int]().Rate(10, time.Second, 10)

	p.Go(nil)

	results := p.Wait().Collect()

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].IsOk() {
		t.Error("expected error for nil function")
	}
	if !strings.Contains(results[0].Err().Error(), "nil function provided") {
		t.Errorf("expected nil function error, got: %v", results[0].Err())
	}
}

func TestRate_Wait_Panic(t *testing.T) {
	p := pool.New[int]().Rate(10, time.Second, 10)

	p.Go(func() Result[int] { panic("rate panic") })

	results := p.Wait().Collect()

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Err().Error(), "panic: rate panic") {
		t.Errorf("expected panic error, got: %v", results[0].Err())
	}
}

func TestRate_Stream_NilFunction(t *testing.T) {
	p := pool.New[int]().Limit(2).Rate(10, time.Second, 10)

	ch := p.Stream(func() {
		p.Go(nil)
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Err().Error(), "nil function provided") {
		t.Errorf("expected nil function error, got: %v", results[0].Err())
	}
}

func TestRate_Stream_Panic(t *testing.T) {
	p := pool.New[int]().Limit(2).Rate(10, time.Second, 10)

	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { panic("stream rate panic") })
		p.Go(func() Result[int] { return Ok(3) })
	})

	var oks, errs int
	for r := range ch {
		if r.IsOk() {
			oks++
		} else {
			if !strings.Contains(r.Err().Error(), "panic: stream rate panic") {
				t.Errorf("unexpected error: %v", r.Err())
			}
			errs++
		}
	}

	if errs != 1 {
		t.Errorf("expected 1 error, got %d", errs)
	}
}

// --- Rate lifecycle ---

func TestRate_PanicWhileRunning(t *testing.T) {
	p := pool.New[int]()
	p.Go(func() Result[int] { return Ok(1) })

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when calling Rate while tasks are running")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "cannot change rate limit") {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	p.Rate(5, time.Second)
}

func TestRate_ZeroRemovesLimit(t *testing.T) {
	p := pool.New[int]().Rate(5, time.Second).Rate(0, time.Second)

	start := time.Now()
	for i := range 10 {
		p.Go(func() Result[int] { return Ok(i) })
	}
	results := p.Wait().Collect()

	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Errorf("without rate limit, expected fast completion")
	}
}

func TestRate_NegativeRemovesLimit(t *testing.T) {
	p := pool.New[int]().Rate(5, time.Second).Rate(-1, time.Second)

	for i := range 5 {
		p.Go(func() Result[int] { return Ok(i) })
	}

	results := p.Wait().Collect()
	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
}

func TestRate_ReplaceStopsOldLimiter(t *testing.T) {
	p := pool.New[int]().
		Rate(1, time.Second, 1).    // very slow
		Rate(100, time.Second, 100) // fast — old should be stopped

	start := time.Now()
	for i := range 10 {
		p.Go(func() Result[int] { return Ok(i) })
	}
	results := p.Wait().Collect()

	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Errorf("old slow limiter may still be active, took %s", time.Since(start))
	}
}

func TestRate_ResetClearsRate(t *testing.T) {
	p := pool.New[int]().Rate(1, time.Second, 1)

	p.Go(func() Result[int] { return Ok(1) })
	p.Wait()

	p.Reset()

	// After reset, no rate limit — should be fast.
	start := time.Now()
	for i := range 10 {
		p.Go(func() Result[int] { return Ok(i) })
	}
	results := p.Wait().Collect()

	if len(results) != 10 {
		t.Errorf("expected 10 results after reset, got %d", len(results))
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Errorf("rate limit persisted after reset, took %s", time.Since(start))
	}
}

func TestRate_CancelStopsLimiter(t *testing.T) {
	p := pool.New[int]().Rate(1, time.Second, 1)

	p.Go(func() Result[int] { return Ok(1) })
	p.Wait()

	// Pool is canceled after Wait. Reset and verify no lingering delay.
	p.Reset()

	start := time.Now()
	for i := range 5 {
		p.Go(func() Result[int] { return Ok(i) })
	}
	p.Wait()

	if time.Since(start) > 500*time.Millisecond {
		t.Error("old limiter still active after cancel + reset")
	}
}

func TestRate_ResetThenStream(t *testing.T) {
	// Wait with rate → reset → Stream with rate.
	p := pool.New[int]().Rate(20, time.Second, 20)

	p.Go(func() Result[int] { return Ok(1) })
	results := p.Wait().Collect()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	p.Reset()
	p.Limit(2).Rate(20, time.Second, 20)

	ch := p.Stream(func() {
		for i := range 5 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	count := 0
	for range ch {
		count++
	}

	if count != 5 {
		t.Errorf("expected 5 results after reset+stream, got %d", count)
	}
}

// --- Rate metrics ---

func TestRate_Wait_Metrics(t *testing.T) {
	p := pool.New[int]().Rate(20, time.Second, 20)

	p.Go(func() Result[int] { return Ok(1) })
	p.Go(func() Result[int] { return Err[int](errors.New("fail")) })
	p.Go(func() Result[int] { return Ok(3) })

	p.Wait()

	if p.TotalTasks() != 3 {
		t.Errorf("expected totalTasks=3, got %d", p.TotalTasks())
	}
	if p.FailedTasks() != 1 {
		t.Errorf("expected failedTasks=1, got %d", p.FailedTasks())
	}
	if p.SuccessfulTasks() != 2 {
		t.Errorf("expected successfulTasks=2, got %d", p.SuccessfulTasks())
	}
	if p.ActiveTasks() != 0 {
		t.Errorf("expected activeTasks=0, got %d", p.ActiveTasks())
	}
}

func TestRate_Stream_Metrics(t *testing.T) {
	p := pool.New[int]().Limit(2).Rate(20, time.Second, 20)

	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { return Err[int](errors.New("fail")) })
		p.Go(func() Result[int] { return Ok(3) })
	})

	for range ch {
	}

	if p.TotalTasks() != 3 {
		t.Errorf("expected totalTasks=3, got %d", p.TotalTasks())
	}
	if p.FailedTasks() != 1 {
		t.Errorf("expected failedTasks=1, got %d", p.FailedTasks())
	}
	if p.ActiveTasks() != 0 {
		t.Errorf("expected activeTasks=0, got %d", p.ActiveTasks())
	}
}

// --- Minute-based rate ---

func TestRate_Wait_MinuteRate(t *testing.T) {
	// 60/min = 1/sec. Burst=1. Second task should wait ~1s.
	// Context with timeout prevents test from hanging if rate logic breaks.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	p := pool.New[int]().
		Rate(60, time.Minute, 1).
		Context(ctx) // set context BEFORE any Go calls to avoid race

	start := time.Now()

	// First: instant (burst).
	p.Go(func() Result[int] { return Ok(1) })

	// Second: should wait ~1s (rate = 1 token/sec after burst consumed).
	p.Go(func() Result[int] { return Ok(2) })

	p.Wait()

	elapsed := time.Since(start)
	if elapsed < 800*time.Millisecond {
		t.Errorf("expected ~1s for second task at 60/min, got %s", elapsed)
	}
}
