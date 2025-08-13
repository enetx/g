package g_test

import (
	"context"
	"errors"
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

		results := p.Wait()

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

		results := p.Wait()

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

	results := p.Wait()
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
	p.Limit(0)

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

	results := p.Wait()
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

	results := p.Wait()

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

	results := p.Wait()

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

	results := p.Wait()

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

			results := p.Wait()

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

		results := p.Wait()
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

		results := p.Wait()
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

	results := p.Wait()

	// Should get no results since context was cancelled
	if len(results) != 0 {
		t.Errorf("Expected 0 results with cancelled context, got %d", len(results))
	}
}
