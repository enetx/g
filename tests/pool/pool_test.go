package g_test

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func TestGoPanic(t *testing.T) {
	t.Run("PanicWithoutCancelOnError", func(t *testing.T) {
		p := pool.New[int]()

		p.Go(func() Result[int] {
			panic("test panic")
		})

		results := p.Wait().Collect()

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].Err() == nil {
			t.Fatal("Expected error in result, got nil")
		}

		if !strings.Contains(results[0].Err().Error(), "panic: test panic") {
			t.Errorf("Expected error to contain 'panic: test panic', got %q", results[0].Err().Error())
		}

		if p.FailedTasks() != 1 {
			t.Errorf("Expected 1 failed task, got %d", p.FailedTasks())
		}

		if p.GetContext().Err() == nil {
			t.Error("Expected pool to be cancelled after Wait")
		}
	})

	t.Run("PanicWithCancelOnError", func(t *testing.T) {
		p := pool.New[int]()
		p.Limit(1)
		p.CancelOnError()

		p.Go(func() Result[int] {
			panic("test panic")
		})

		p.Go(func() Result[int] {
			return Ok(42)
		})

		results := p.Wait().Collect()

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].IsOk() {
			t.Fatal("Expected error in result, got nil")
		}

		if !strings.Contains(results[0].Err().Error(), "panic: test panic") {
			t.Errorf("Expected error to contain 'panic: test panic', got %q", results[0].Err().Error())
		}

		if p.FailedTasks() != 1 {
			t.Errorf("Expected 1 failed task, got %d", p.FailedTasks())
		}

		if !errors.Is(p.Cause(), results[0].Err()) {
			t.Errorf("Expected pool cancellation cause to match panic error, got %v", p.Cause())
		}
	})
}

func TestPool(t *testing.T) {
	p := pool.New[int]()

	successCount := int32(0)
	for i := range 5 {
		p.Go(func() Result[int] {
			if i%2 == 0 {
				atomic.AddInt32(&successCount, 1)
				return Ok(i)
			}
			return Err[int](nil)
		})
	}

	results := p.Wait().Collect()
	if p.ActiveTasks() != 0 {
		t.Errorf("expected no active tasks after Wait, got %d", p.ActiveTasks())
	}
	if p.TotalTasks() != 5 {
		t.Errorf("expected totalTasks=5, got %d", p.TotalTasks())
	}
	if p.FailedTasks() != 2 {
		t.Errorf("expected failedTasks=2, got %d", p.FailedTasks())
	}

	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}

	successFound := 0
	failFound := 0
	for _, r := range results {
		if r.IsErr() {
			failFound++
		} else {
			successFound++
		}
	}

	if successFound != 3 || failFound != 2 {
		t.Errorf("expected 3 successes and 2 fails, got %d successes and %d fails", successFound, failFound)
	}
}

func TestPoolLimit(t *testing.T) {
	p := pool.New[int]()
	p.Limit(2)

	activeGoroutines := int32(0)
	maxObserved := int32(0)

	for range 5 {
		p.Go(func() Result[int] {
			cur := atomic.AddInt32(&activeGoroutines, 1)
			if cur > atomic.LoadInt32(&maxObserved) {
				atomic.StoreInt32(&maxObserved, cur)
			}

			atomic.AddInt32(&activeGoroutines, -1)
			return Ok(0)
		})
	}

	p.Wait()

	if maxObserved > 2 {
		t.Errorf("observed concurrency %d, but limit was set to 2", maxObserved)
	}
}

func TestPoolReset(t *testing.T) {
	p := pool.New[int]()
	p.Go(func() Result[int] {
		return Ok(1)
	})
	p.Wait()

	if p.TotalTasks() != 1 {
		t.Errorf("expected totalTasks=1, got %d", p.TotalTasks())
	}

	p.Reset()

	if p.TotalTasks() != 0 {
		t.Errorf("expected totalTasks=0 after Reset, got %d", p.TotalTasks())
	}
	if p.FailedTasks() != 0 {
		t.Errorf("expected failedTasks=0 after Reset, got %d", p.FailedTasks())
	}

	p.Go(func() Result[int] {
		return Ok(2)
	})

	results := p.Wait().Collect()
	if len(results) != 1 {
		t.Errorf("expected 1 result after new task, got %d", len(results))
	}
}

func TestPoolCancel(t *testing.T) {
	p := pool.New[int]()
	p.Limit(1)

	ctx := context.Background()
	p.Context(ctx)
	p.Context(nil)

	for i := range 100 {
		p.Go(func() Result[int] {
			if i == 3 {
				p.Cancel()
			}
			return Ok(1)
		})
	}

	results := p.Wait().Collect()

	if len(results) != 4 {
		t.Errorf("expected 4 results, got %d", len(results))
	}

	t.Logf("Received %d results after calling pool.Cancel()", len(results))
}

func TestPoolCause(t *testing.T) {
	p := pool.New[int]()
	cancelErr := errors.New("custom cancellation reason")

	p.Cancel(cancelErr)

	if p.Cause() == nil {
		t.Errorf("expected Cause to return a non-nil error after cancellation")
	} else if !errors.Is(p.Cause(), cancelErr) {
		t.Errorf("expected Cause to return %v, got %v", cancelErr, p.Cause())
	}
}

func TestPoolCancelOnError(t *testing.T) {
	p := pool.New[int]().CancelOnError().Limit(1)

	p.Go(func() Result[int] {
		return Err[int](errors.New("task failed 1"))
	})

	p.Go(func() Result[int] {
		return Err[int](errors.New("task failed 2"))
	})

	p.Go(func() Result[int] {
		return Ok(42)
	})

	results := p.Wait().Collect()

	if len(results) != 1 {
		t.Errorf("Expected 1 results, got %d", len(results))
	}

	if !results[0].IsErr() {
		t.Errorf("Expected first task to fail, but it did not")
	}
}

func TestPoolResetWithActiveTasks(t *testing.T) {
	p := pool.New[int]().Limit(1)

	// Add a task but don't wait for it to complete
	p.Go(func() Result[int] {
		return Ok(1)
	})

	// Try to reset while tasks might be active
	// This should return an error since tasks are running
	if p.ActiveTasks() > 0 {
		err := p.Reset()
		if err == nil {
			t.Error("Expected error when resetting with active tasks")
		}
	}

	// Clean up
	p.Wait()

	// Now reset should work
	err := p.Reset()
	if err != nil {
		t.Errorf("Expected no error when resetting after wait, got: %v", err)
	}
}

func TestPoolLimitPanic(t *testing.T) {
	p := pool.New[int]().Limit(2)

	// Add some tasks to make the tokens channel have length > 0
	p.Go(func() Result[int] {
		return Ok(1)
	})

	// This should panic since we're trying to change limit while tasks are potentially running
	defer func() {
		if r := recover(); r != nil {
			// Expected panic
			if !strings.Contains(r.(string), "cannot change semaphore limit") {
				t.Errorf("Expected panic about semaphore limit, got: %v", r)
			}
		} else {
			// If no panic, wait for tasks and then check if changing limit works
			p.Wait()
			// After wait, changing limit should work fine
			p.Limit(3)
		}
	}()

	// Try to change limit while tasks might be running
	p.Limit(3)
}

func TestPoolGoWithNilFunction(t *testing.T) {
	p := pool.New[int]()

	p.Go(nil) // Pass nil function

	results := p.Wait().Collect()

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].IsOk() {
		t.Error("Expected error result for nil function, got success")
	}

	if !strings.Contains(results[0].Err().Error(), "nil function provided") {
		t.Errorf("Expected error about nil function, got: %v", results[0].Err())
	}
}

func TestPoolEpanicTypes(t *testing.T) {
	testCases := []struct {
		name      string
		panicVal  any
		expectStr string
	}{
		{
			name:      "string panic",
			panicVal:  "string panic test",
			expectStr: "panic: string panic test",
		},
		{
			name:      "error panic",
			panicVal:  errors.New("error panic test"),
			expectStr: "panic: error panic test",
		},
		{
			name:      "other type panic",
			panicVal:  123,
			expectStr: "panic: 123",
		},
		{
			name:      "nil panic",
			panicVal:  nil,
			expectStr: "panic: panic called with nil argument",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := pool.New[int]()

			p.Go(func() Result[int] {
				panic(tc.panicVal)
			})

			results := p.Wait().Collect()

			if len(results) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(results))
			}

			if results[0].IsOk() {
				t.Error("Expected error result for panic, got success")
			}

			if !strings.Contains(results[0].Err().Error(), tc.expectStr) {
				t.Errorf("Expected error to contain '%s', got: %v", tc.expectStr, results[0].Err())
			}
		})
	}
}

func TestPoolLimitEdgeCases(t *testing.T) {
	t.Run("negative limit", func(t *testing.T) {
		p := pool.New[int]().Limit(-5)

		// Should set tokens to nil for negative values
		p.Go(func() Result[int] {
			return Ok(1)
		})

		results := p.Wait().Collect()
		if len(results) != 1 {
			t.Errorf("Expected 1 result with negative limit, got %d", len(results))
		}
	})

	t.Run("zero limit", func(t *testing.T) {
		p := pool.New[int]().Limit(0)

		// Should set tokens to nil for zero
		p.Go(func() Result[int] {
			return Ok(1)
		})

		results := p.Wait().Collect()
		if len(results) != 1 {
			t.Errorf("Expected 1 result with zero limit, got %d", len(results))
		}
	})
}

func TestPoolContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := pool.New[int]().Context(ctx).Limit(1)

	// Cancel context immediately
	cancel()

	// Try to add tasks - they should not execute due to cancelled context
	p.Go(func() Result[int] {
		return Ok(1)
	})

	results := p.Wait().Collect()

	// Should get no results since context was cancelled
	if len(results) != 0 {
		t.Errorf("Expected 0 results with cancelled context, got %d", len(results))
	}
}

func TestPoolStream_Basic(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { return Ok(2) })
		p.Go(func() Result[int] { return Ok(3) })
	})

	var results []int
	for r := range ch {
		if r.IsErr() {
			t.Errorf("Unexpected error: %v", r.Err())
		}
		results = append(results, r.Ok())
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestPoolStream_ChannelCloses(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(42) })
	})

	count := 0
	for range ch {
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 result, got %d", count)
	}
}

func TestPoolStream_WithErrors(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { return Err[int](errors.New("fail")) })
		p.Go(func() Result[int] { return Ok(3) })
	})

	var oks, errs int
	for r := range ch {
		if r.IsOk() {
			oks++
		} else {
			errs++
		}
	}

	if oks != 2 {
		t.Errorf("Expected 2 ok results, got %d", oks)
	}
	if errs != 1 {
		t.Errorf("Expected 1 error result, got %d", errs)
	}
}

func TestPoolStream_WithPanic(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { panic("stream panic") })
	})

	var oks, errs int
	for r := range ch {
		if r.IsOk() {
			oks++
		} else {
			if !strings.Contains(r.Err().Error(), "panic: stream panic") {
				t.Errorf("Expected panic error, got: %v", r.Err())
			}
			errs++
		}
	}

	if oks != 1 {
		t.Errorf("Expected 1 ok result, got %d", oks)
	}
	if errs != 1 {
		t.Errorf("Expected 1 error result, got %d", errs)
	}
}

func TestPoolStream_WithNilFunction(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(nil)
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].IsOk() {
		t.Error("Expected error for nil function")
	}

	if !strings.Contains(results[0].Err().Error(), "nil function provided") {
		t.Errorf("Expected nil function error, got: %v", results[0].Err())
	}
}

func TestPoolStream_WithLimit(t *testing.T) {
	p := pool.New[int]().Limit(2)
	ch := p.Stream(func() {
		for i := range 10 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	var results []int
	for r := range ch {
		results = append(results, r.Ok())
	}

	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}
}

func TestPoolStream_Buffered(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		for i := range 5 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	}, 5)

	var results []int
	for r := range ch {
		results = append(results, r.Ok())
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}
}

func TestPoolStream_CancelOnError(t *testing.T) {
	p := pool.New[int]().CancelOnError().Limit(1)
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Err[int](errors.New("first fail")) })
		p.Go(func() Result[int] { return Ok(2) })
		p.Go(func() Result[int] { return Ok(3) })
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	if len(results) < 1 {
		t.Fatal("Expected at least 1 result (the error)")
	}

	hasError := false
	for _, r := range results {
		if r.IsErr() && strings.Contains(r.Err().Error(), "first fail") {
			hasError = true
		}
	}

	if !hasError {
		t.Error("Expected 'first fail' error in results")
	}

	if p.FailedTasks() < 1 {
		t.Errorf("Expected at least 1 failed task, got %d", p.FailedTasks())
	}
}

func TestPoolStream_CancelOnError_ErrorNotLost(t *testing.T) {
	for range 100 {
		p := pool.New[int]().CancelOnError().Limit(1)
		ch := p.Stream(func() {
			p.Go(func() Result[int] { return Err[int](errors.New("must arrive")) })
		}, 1)

		var gotError bool
		for r := range ch {
			if r.IsErr() && strings.Contains(r.Err().Error(), "must arrive") {
				gotError = true
			}
		}

		if !gotError {
			t.Fatal("Error result was lost — emit likely happened after Cancel")
		}
	}
}

func TestPoolStream_NoResults(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		// no tasks
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestPoolStream_DoesNotAccumulateResults(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { return Ok(2) })
	})

	for range ch {
	}

	if p.TotalTasks() != 2 {
		t.Errorf("Expected totalTasks=2, got %d", p.TotalTasks())
	}
}

func TestPoolStream_WithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := pool.New[int]().Context(ctx).Limit(1)
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] {
			cancel()
			return Ok(2)
		})
		p.Go(func() Result[int] { return Ok(3) })
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	if len(results) > 3 {
		t.Errorf("Expected at most 3 results, got %d", len(results))
	}

	t.Logf("Received %d results before context cancellation", len(results))
}

func TestPoolStream_PanicBeforeStream(t *testing.T) {
	p := pool.New[int]()

	p.Go(func() Result[int] { return Ok(1) })

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic when calling Stream after Go")
		}
		if !strings.Contains(r.(string), "Stream must be called before submitting tasks") {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()

	p.Stream(func() {})
}

func TestPoolStream_Concurrent(t *testing.T) {
	p := pool.New[int]().Limit(10)
	tasks := 1000
	ch := p.Stream(func() {
		for i := range tasks {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	var count atomic.Int32
	for range ch {
		count.Add(1)
	}

	if int(count.Load()) != tasks {
		t.Errorf("Expected %d results, got %d", tasks, count.Load())
	}

	if p.TotalTasks() != tasks {
		t.Errorf("Expected totalTasks=%d, got %d", tasks, p.TotalTasks())
	}
}

func TestPoolStream_ResetAfterStream(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
	})

	for range ch {
	}

	err := p.Reset()
	if err != nil {
		t.Errorf("Expected no error resetting after stream drained, got: %v", err)
	}

	p.Go(func() Result[int] { return Ok(2) })

	results := p.Wait().Collect()
	if len(results) != 1 {
		t.Errorf("Expected 1 result after reset, got %d", len(results))
	}

	if results[0].Ok() != 2 {
		t.Errorf("Expected 2, got %d", results[0].Ok())
	}
}

func TestPoolStream_Metrics(t *testing.T) {
	p := pool.New[int]()
	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { return Err[int](errors.New("fail")) })
		p.Go(func() Result[int] { return Ok(3) })
	})

	for range ch {
	}

	if p.TotalTasks() != 3 {
		t.Errorf("Expected totalTasks=3, got %d", p.TotalTasks())
	}
	if p.FailedTasks() != 1 {
		t.Errorf("Expected failedTasks=1, got %d", p.FailedTasks())
	}
	if p.ActiveTasks() != 0 {
		t.Errorf("Expected activeTasks=0, got %d", p.ActiveTasks())
	}
}

func TestPoolStream_WorkerCount(t *testing.T) {
	p := pool.New[int]().Limit(3)

	var maxConcurrent atomic.Int32
	var current atomic.Int32

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
				runtime.Gosched()
				current.Add(-1)
				return Ok(0)
			})
		}
	})

	for range ch {
	}

	if maxConcurrent.Load() > 3 {
		t.Errorf("Expected max concurrency <= 3, got %d", maxConcurrent.Load())
	}
}

func TestPoolStream_FIFO(t *testing.T) {
	p := pool.New[int]().Limit(1)
	ch := p.Stream(func() {
		for i := range 10 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	var results []int
	for r := range ch {
		results = append(results, r.Ok())
	}

	if len(results) != 10 {
		t.Fatalf("Expected 10 results, got %d", len(results))
	}

	for i, v := range results {
		if v != i {
			t.Errorf("Expected results[%d]=%d, got %d (FIFO violated)", i, i, v)
			break
		}
	}
}

func TestPoolStream_ConsumerCancel(t *testing.T) {
	p := pool.New[int]().Limit(1)
	ch := p.Stream(func() {
		for i := range 1000 {
			p.Go(func() Result[int] { return Ok(i * i) })
		}
	})

	var collected []int
	for r := range ch {
		collected = append(collected, r.Ok())
		if r.Ok() == 25 { // found 5*5
			p.Cancel()
			break
		}
	}

	if len(collected) < 1 {
		t.Error("Expected at least 1 result")
	}

	t.Logf("Collected %d results before cancel", len(collected))
}

func TestCancelOn_SuccessPredicate(t *testing.T) {
	p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
		return r.IsOk() && r.Ok() > 5
	})

	for i := range 100 {
		p.Go(func() Result[int] { return Ok(i) })
	}

	results := p.Wait().Collect()

	// Should stop early once a value > 5 is produced.
	if len(results) >= 100 {
		t.Errorf("Expected early cancellation, got all %d results", len(results))
	}

	// The triggering result (value > 5) must be present.
	found := false
	for _, r := range results {
		if r.IsOk() && r.Ok() > 5 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected at least one result > 5 (the trigger)")
	}

	t.Logf("Got %d results before cancellation", len(results))
}

func TestCancelOn_SpecificErrorType(t *testing.T) {
	errCritical := errors.New("critical")

	p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
		return r.IsErr() && errors.Is(r.Err(), errCritical)
	})

	p.Go(func() Result[int] { return Err[int](errors.New("minor")) })
	p.Go(func() Result[int] { return Err[int](errCritical) })
	p.Go(func() Result[int] { return Ok(42) })

	results := p.Wait().Collect()

	// Minor error should not trigger cancellation, critical should.
	hasCritical := false
	for _, r := range results {
		if r.IsErr() && errors.Is(r.Err(), errCritical) {
			hasCritical = true
		}
	}
	if !hasCritical {
		t.Error("Expected critical error in results")
	}

	if !errors.Is(p.Cause(), errCritical) {
		t.Errorf("Expected cause to be errCritical, got: %v", p.Cause())
	}
}

func TestCancelOn_NoMatch(t *testing.T) {
	p := pool.New[int]().Limit(2).CancelOn(func(r Result[int]) bool {
		return r.IsOk() && r.Ok() > 1000 // never true
	})

	for i := range 10 {
		p.Go(func() Result[int] { return Ok(i) })
	}

	results := p.Wait().Collect()

	if len(results) != 10 {
		t.Errorf("Expected all 10 results (predicate never matched), got %d", len(results))
	}
}

func TestCancelOn_PredicatePanic(t *testing.T) {
	p := pool.New[int]().Limit(1).CancelOn(func(Result[int]) bool {
		panic("predicate boom")
	})

	p.Go(func() Result[int] { return Ok(1) })

	results := p.Wait().Collect()

	// Pool should be cancelled due to predicate panic.
	cause := p.Cause()
	if cause == nil {
		t.Fatal("Expected pool to be cancelled after predicate panic")
	}

	if !strings.Contains(cause.Error(), "predicate panicked") {
		t.Errorf("Expected cause to mention predicate panic, got: %v", cause)
	}

	t.Logf("Got %d results, cause: %v", len(results), cause)
}

func TestCancelOn_PredicatePanicOnError(t *testing.T) {
	p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
		if r.IsErr() {
			panic("predicate panic on error path")
		}
		return false
	})

	p.Go(func() Result[int] { return Err[int](errors.New("task error")) })
	p.Go(func() Result[int] { return Ok(2) })

	results := p.Wait().Collect()

	cause := p.Cause()
	if cause == nil {
		t.Fatal("Expected cancellation after predicate panic")
	}
	if !strings.Contains(cause.Error(), "predicate panicked") {
		t.Errorf("Expected predicate panic cause, got: %v", cause)
	}

	// The error result should still be recorded (error() stores before checkCancel).
	hasErr := false
	for _, r := range results {
		if r.IsErr() {
			hasErr = true
		}
	}
	if !hasErr {
		t.Error("Expected the task error to be in results")
	}

	t.Logf("Got %d results, cause: %v", len(results), cause)
}

func TestCancelOn_CalledAfterGo_Panics(t *testing.T) {
	p := pool.New[int]()

	p.Go(func() Result[int] { return Ok(1) })

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic when calling CancelOn after Go")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "cannot set cancellation predicate") {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()

	p.CancelOn(func(Result[int]) bool { return true })
}

func TestCancelOn_Stream_SuccessPredicate(t *testing.T) {
	p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
		return r.IsOk() && r.Ok() == 5
	})

	ch := p.Stream(func() {
		for i := range 100 {
			p.Go(func() Result[int] { return Ok(i) })
		}
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	if len(results) >= 100 {
		t.Errorf("Expected early cancellation in stream, got all %d results", len(results))
	}

	found := false
	for _, r := range results {
		if r.IsOk() && r.Ok() == 5 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected trigger value (5) in stream results")
	}

	t.Logf("Stream got %d results before cancellation", len(results))
}

func TestCancelOn_Stream_ErrorPredicate(t *testing.T) {
	errFatal := errors.New("fatal")

	p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
		return r.IsErr() && errors.Is(r.Err(), errFatal)
	})

	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { return Err[int](errFatal) })
		p.Go(func() Result[int] { return Ok(3) })
	})

	var results []Result[int]
	for r := range ch {
		results = append(results, r)
	}

	hasFatal := false
	for _, r := range results {
		if r.IsErr() && errors.Is(r.Err(), errFatal) {
			hasFatal = true
		}
	}
	if !hasFatal {
		t.Error("Expected fatal error in stream results")
	}

	if !errors.Is(p.Cause(), errFatal) {
		t.Errorf("Expected cause to be errFatal, got: %v", p.Cause())
	}
}

func TestCancelOn_Stream_PredicatePanic(t *testing.T) {
	p := pool.New[int]().Limit(1).CancelOn(func(Result[int]) bool {
		panic("stream predicate boom")
	})

	ch := p.Stream(func() {
		p.Go(func() Result[int] { return Ok(1) })
		p.Go(func() Result[int] { return Ok(2) })
	})

	for range ch {
	}

	cause := p.Cause()
	if cause == nil {
		t.Fatal("Expected cancellation from predicate panic in stream")
	}
	if !strings.Contains(cause.Error(), "predicate panicked") {
		t.Errorf("Expected predicate panic cause in stream, got: %v", cause)
	}
}

func TestCancelOn_Stream_ErrorGuaranteed(t *testing.T) {
	// Same spirit as TestPoolStream_CancelOnError_ErrorNotLost but with CancelOn.
	for range 100 {
		p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
			return r.IsErr()
		})

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
			t.Fatal("Error result was lost with CancelOn in stream mode")
		}
	}
}

func TestCancelOn_WrappedError(t *testing.T) {
	errBase := errors.New("base")

	p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
		return r.IsErr() && errors.Is(r.Err(), errBase)
	})

	p.Go(func() Result[int] { return Ok(1) })
	p.Go(func() Result[int] {
		return Err[int](errors.Join(errors.New("wrapper"), errBase))
	})
	p.Go(func() Result[int] { return Ok(3) })

	results := p.Wait().Collect()

	if len(results) >= 3 {
		t.Errorf("Expected early cancellation on wrapped error, got %d results", len(results))
	}

	if !errors.Is(p.Cause(), errBase) {
		t.Errorf("Expected cause chain to contain errBase, got: %v", p.Cause())
	}
}

func TestCancelOn_CountingPredicate(t *testing.T) {
	// Cancel after N failures (using external counter).
	var failures atomic.Int32

	p := pool.New[int]().Limit(1).CancelOn(func(r Result[int]) bool {
		if r.IsErr() {
			return failures.Add(1) >= 3
		}
		return false
	})

	for i := range 20 {
		p.Go(func() Result[int] {
			if i%2 == 0 {
				return Err[int](errors.New("fail"))
			}
			return Ok(i)
		})
	}

	results := p.Wait().Collect()

	if len(results) >= 20 {
		t.Errorf("Expected cancellation after 3 failures, got all %d results", len(results))
	}

	if failures.Load() < 3 {
		t.Errorf("Expected at least 3 failures before cancel, got %d", failures.Load())
	}

	t.Logf("Got %d results, %d failures counted", len(results), failures.Load())
}

func TestCancelOn_ResetClearsPredicate(t *testing.T) {
	p := pool.New[int]().CancelOn(func(r Result[int]) bool {
		return r.IsErr()
	})

	p.Go(func() Result[int] { return Err[int](errors.New("fail")) })
	p.Wait()

	p.Reset()

	// After reset, predicate should be cleared — errors should not cancel.
	p.Go(func() Result[int] { return Err[int](errors.New("fail1")) })
	p.Go(func() Result[int] { return Err[int](errors.New("fail2")) })
	p.Go(func() Result[int] { return Ok(42) })

	results := p.Wait().Collect()

	if len(results) != 3 {
		t.Errorf("Expected 3 results after reset (no predicate), got %d", len(results))
	}
}
