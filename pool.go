package g

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"

	"github.com/enetx/g/internal/rlimit"
)

// NewPool[T any] creates a new goroutine pool.
func NewPool[T any]() *Pool[T] {
	ctx, cancel := context.WithCancelCause(context.Background())

	return &Pool[T]{
		ctx:       ctx,
		cancel:    cancel,
		semaphore: nil,
	}
}

func (p *Pool[T]) acquire() bool {
	if p.semaphore != nil {
		select {
		case p.semaphore <- struct{}{}:
			return true
		case <-p.ctx.Done():
			return false
		}
	}

	select {
	case <-p.ctx.Done():
		return false
	default:
		return true
	}
}

func (p *Pool[T]) done() {
	defer atomic.AddInt32(&p.activeTasks, -1)

	if p.semaphore != nil {
		<-p.semaphore
	}

	p.wg.Done()
}

// Go launches an asynchronous task fn() in its own goroutine.
func (p *Pool[T]) Go(fn func() Result[T]) {
	if !p.acquire() {
		return
	}

	index := atomic.AddInt32(&p.totalTasks, 1) - 1
	atomic.AddInt32(&p.activeTasks, 1)
	p.wg.Add(1)

	go func(index int32) {
		defer p.done()

		select {
		case <-p.ctx.Done():
		default:
			result := fn()
			if result.IsErr() {
				if p.cancelOnError {
					p.once.Do(func() { p.Cancel(errors.New("cancel on error")) })
				}

				atomic.AddInt32(&p.failedTasks, 1)
			}

			p.results.Store(int(index), result)
		}
	}(index)
}

// Wait waits for all submitted tasks in the pool to finish.
func (p *Pool[T]) Wait() Slice[Result[T]] {
	p.wg.Wait()
	p.Cancel()
	p.semaphore = nil

	var results []Result[T]

	p.results.Range(func(_, result any) bool {
		results = append(results, result.(Result[T]))
		return true
	})

	return results
}

// Limit sets the maximum number of concurrently running tasks.
func (p *Pool[T]) Limit(workers int) *Pool[T] {
	if workers <= 0 {
		p.semaphore = nil
		return p
	}

	if len(p.semaphore) > 0 {
		panic("cannot change semaphore limit while tasks are running")
	}

	if runtime.GOOS != "windows" {
		workers = rlimit.RlimitStack(workers)
	}

	p.semaphore = make(chan struct{}, workers)

	return p
}

// CancelOnError enables cancellation of remaining tasks on failure.
func (p *Pool[T]) CancelOnError() *Pool[T] {
	p.cancelOnError = true
	return p
}

// Context replaces the pool’s context with the provided context.
// If ctx is nil, context.Background() is used by default.
func (p *Pool[T]) Context(ctx context.Context) *Pool[T] {
	if ctx == nil {
		ctx = context.Background()
	}

	p.Cancel()
	p.ctx, p.cancel = context.WithCancelCause(ctx)

	return p
}

// GetContext returns the current context associated with the pool.
func (p *Pool[T]) GetContext() context.Context { return p.ctx }

// Cancel cancels all tasks in the pool.
func (p *Pool[T]) Cancel(err ...error) {
	if p.cancel != nil {
		cause := context.Canceled
		if len(err) != 0 {
			cause = err[0]
		}

		p.cancel(cause)
	}
}

// Cause returns the reason for the cancellation of the pool's context.
// It retrieves the underlying cause of the context's termination if the context has been canceled.
// If the pool's context is still active, it returns nil.
func (p *Pool[T]) Cause() error { return context.Cause(p.ctx) }

// Reset restores the pool to its initial state: cancels all tasks, clears results and metrics,
// and creates a new context. If there are any active tasks, it will panic.
func (p *Pool[T]) Reset() {
	if p.ActiveTasks() > 0 {
		panic("cannot reset while tasks are running")
	}

	p.Cancel()
	p.ClearResults()
	p.ClearMetrics()
	p.semaphore = nil
	p.ctx, p.cancel = context.WithCancelCause(context.Background())
}

// ClearResults removes all stored task results from the pool.
func (p *Pool[T]) ClearResults() {
	p.results.Range(func(key, _ any) bool {
		p.results.Delete(key)
		return true
	})
}

// ClearMetrics resets both total tasks and failed tasks counters to zero.
func (p *Pool[T]) ClearMetrics() {
	atomic.StoreInt32(&p.totalTasks, 0)
	atomic.StoreInt32(&p.failedTasks, 0)
}

// TotalTasks returns the total number of tasks that have been submitted.
func (p *Pool[T]) TotalTasks() int { return int(atomic.LoadInt32(&p.totalTasks)) }

// ActiveTasks returns the current number of tasks that are still running.
func (p *Pool[T]) ActiveTasks() int { return int(atomic.LoadInt32(&p.activeTasks)) }

// FailedTasks returns the number of tasks that have completed with an error.
func (p *Pool[T]) FailedTasks() int { return int(atomic.LoadInt32(&p.failedTasks)) }
