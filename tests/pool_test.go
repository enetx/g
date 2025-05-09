package g_test

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	. "github.com/enetx/g"
)

func TestGoPanic(t *testing.T) {
	t.Run("PanicWithoutCancelOnError", func(t *testing.T) {
		pool := NewPool[int]()

		pool.Go(func() Result[int] {
			panic("test panic")
		})

		results := pool.Wait()

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].Err() == nil {
			t.Fatal("Expected error in result, got nil")
		}

		if !strings.Contains(results[0].Err().Error(), "panic: test panic") {
			t.Errorf("Expected error to contain 'panic: test panic', got %q", results[0].Err().Error())
		}

		if pool.FailedTasks() != 1 {
			t.Errorf("Expected 1 failed task, got %d", pool.FailedTasks())
		}

		if pool.GetContext().Err() == nil {
			t.Error("Expected pool to be cancelled after Wait")
		}
	})

	t.Run("PanicWithCancelOnError", func(t *testing.T) {
		pool := NewPool[int]()
		pool.Limit(1)
		pool.CancelOnError()

		pool.Go(func() Result[int] {
			panic("test panic")
		})

		pool.Go(func() Result[int] {
			return Ok(42)
		})

		results := pool.Wait()

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0].IsOk() {
			t.Fatal("Expected error in result, got nil")
		}

		if !strings.Contains(results[0].Err().Error(), "panic: test panic") {
			t.Errorf("Expected error to contain 'panic: test panic', got %q", results[0].Err().Error())
		}

		if pool.FailedTasks() != 1 {
			t.Errorf("Expected 1 failed task, got %d", pool.FailedTasks())
		}

		if !errors.Is(pool.Cause(), results[0].Err()) {
			t.Errorf("Expected pool cancellation cause to match panic error, got %v", pool.Cause())
		}
	})
}

func TestPool(t *testing.T) {
	pool := NewPool[int]()

	successCount := int32(0)
	for i := range 5 {
		pool.Go(func() Result[int] {
			if i%2 == 0 {
				atomic.AddInt32(&successCount, 1)
				return Ok(i)
			}
			return Err[int](nil)
		})
	}

	results := pool.Wait()
	if pool.ActiveTasks() != 0 {
		t.Errorf("expected no active tasks after Wait, got %d", pool.ActiveTasks())
	}
	if pool.TotalTasks() != 5 {
		t.Errorf("expected totalTasks=5, got %d", pool.TotalTasks())
	}
	if pool.FailedTasks() != 2 {
		t.Errorf("expected failedTasks=2, got %d", pool.FailedTasks())
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
	pool := NewPool[int]()
	pool.Limit(0)

	activeGoroutines := int32(0)
	maxObserved := int32(0)

	for range 5 {
		pool.Go(func() Result[int] {
			cur := atomic.AddInt32(&activeGoroutines, 1)
			if cur > atomic.LoadInt32(&maxObserved) {
				atomic.StoreInt32(&maxObserved, cur)
			}

			atomic.AddInt32(&activeGoroutines, -1)
			return Ok(0)
		})
	}

	pool.Wait()

	if maxObserved > 2 {
		t.Errorf("observed concurrency %d, but limit was set to 2", maxObserved)
	}
}

func TestPoolReset(t *testing.T) {
	pool := NewPool[int]()
	pool.Go(func() Result[int] {
		return Ok(1)
	})
	pool.Wait()

	if pool.TotalTasks() != 1 {
		t.Errorf("expected totalTasks=1, got %d", pool.TotalTasks())
	}

	pool.Reset()

	if pool.TotalTasks() != 0 {
		t.Errorf("expected totalTasks=0 after Reset, got %d", pool.TotalTasks())
	}
	if pool.FailedTasks() != 0 {
		t.Errorf("expected failedTasks=0 after Reset, got %d", pool.FailedTasks())
	}

	pool.Go(func() Result[int] {
		return Ok(2)
	})

	results := pool.Wait()
	if len(results) != 1 {
		t.Errorf("expected 1 result after new task, got %d", len(results))
	}
}

func TestPoolCancel(t *testing.T) {
	pool := NewPool[int]()
	pool.Limit(1)

	ctx := context.Background()
	pool.Context(ctx)
	pool.Context(nil)

	for i := range 100 {
		pool.Go(func() Result[int] {
			if i == 3 {
				pool.Cancel()
			}
			return Ok(1)
		})
	}

	results := pool.Wait()

	if len(results) != 4 {
		t.Errorf("expected 4 results, got %d", len(results))
	}

	t.Logf("Received %d results after calling pool.Cancel()", len(results))
}

func TestPoolCause(t *testing.T) {
	pool := NewPool[int]()
	cancelErr := errors.New("custom cancellation reason")

	pool.Cancel(cancelErr)

	if pool.Cause() == nil {
		t.Errorf("expected Cause to return a non-nil error after cancellation")
	} else if !errors.Is(pool.Cause(), cancelErr) {
		t.Errorf("expected Cause to return %v, got %v", cancelErr, pool.Cause())
	}
}

func TestPoolCancelOnError(t *testing.T) {
	pool := NewPool[int]().CancelOnError().Limit(1)

	pool.Go(func() Result[int] {
		return Err[int](errors.New("task failed 1"))
	})

	pool.Go(func() Result[int] {
		return Err[int](errors.New("task failed 2"))
	})

	pool.Go(func() Result[int] {
		return Ok(42)
	})

	results := pool.Wait()

	if len(results) != 1 {
		t.Errorf("Expected 1 results, got %d", len(results))
	}

	if !results[0].IsErr() {
		t.Errorf("Expected first task to fail, but it did not")
	}
}
