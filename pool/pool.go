package pool

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"

	. "github.com/enetx/g"
	"github.com/enetx/g/internal/rlimit"
)

// Pool[T any] is a goroutine pool that allows parallel task execution.
type Pool[T any] struct {
	ctx           context.Context          // Context for controlling cancellation and timeouts
	cancel        context.CancelCauseFunc  // Function to cancel the context
	tokens        chan struct{}            // Tokens for limiting concurrency
	results       *MapSafe[int, Result[T]] // Stores task results
	wg            sync.WaitGroup           // Waits for all tasks to complete
	totalTasks    int32                    // Total number of tasks submitted
	activeTasks   int32                    // Number of currently active tasks
	failedTasks   int32                    // Number of failed tasks
	cancelOnError bool                     // Cancels remaining tasks if any task fails
}

// New[T any] creates a new goroutine pool.
func New[T any]() *Pool[T] {
	ctx, cancel := context.WithCancelCause(context.Background())

	return &Pool[T]{
		ctx:     ctx,
		cancel:  cancel,
		tokens:  nil,
		results: NewMapSafe[int, Result[T]](),
	}
}

func (p *Pool[T]) acquire() bool {
	if p.tokens == nil {
		return p.ctx.Err() == nil
	}

	select {
	case p.tokens <- struct{}{}:
		return true
	case <-p.ctx.Done():
		return false
	}
}

func (p *Pool[T]) done() {
	defer atomic.AddInt32(&p.activeTasks, -1)

	if p.tokens != nil {
		<-p.tokens
	}

	p.wg.Done()
}

func (p *Pool[T]) error(index int32, err error) {
	atomic.AddInt32(&p.failedTasks, 1)
	if p.cancelOnError {
		p.Cancel(err)
	}

	p.results.Set(int(index), Err[T](err))
}

// Go launches an asynchronous task fn() in its own goroutine.
func (p *Pool[T]) Go(fn func() Result[T]) {
	if !p.acquire() {
		return
	}

	index := atomic.AddInt32(&p.totalTasks, 1) - 1

	if fn == nil {
		p.error(index, errors.New("nil function provided"))
		return
	}

	atomic.AddInt32(&p.activeTasks, 1)
	p.wg.Add(1)

	go func(index int32) {
		defer p.done()
		defer func() {
			if r := recover(); r != nil {
				p.error(index, epanic(r))
			}
		}()

		if p.ctx.Err() != nil {
			return
		}

		result := fn()
		if result.IsErr() {
			p.error(index, result.Err())
			return
		}

		p.results.Set(int(index), result)
	}(index)
}

// Wait waits for all submitted tasks in the pool to finish.
func (p *Pool[T]) Wait() SeqResult[T] {
	p.wg.Wait()
	p.Cancel()
	p.tokens = nil

	return SeqResult[T](p.results.Iter().Values())
}

// Limit sets the maximum number of concurrently running tasks.
func (p *Pool[T]) Limit(workers int) *Pool[T] {
	if p.tokens != nil && len(p.tokens) > 0 {
		panic("cannot change semaphore limit while tasks are running")
	}

	if workers <= 0 {
		p.tokens = nil
		return p
	}

	if runtime.GOOS != "windows" {
		workers = rlimit.RlimitStack(workers)
	}

	if workers < 1 {
		workers = 1
	}

	p.tokens = make(chan struct{}, workers)

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
	cause := context.Canceled
	if len(err) > 0 && err[0] != nil {
		cause = err[0]
	}

	p.cancel(cause)
}

// Cause returns the reason for the cancellation of the pool's context.
// It retrieves the underlying cause of the context's termination if the context has been canceled.
// If the pool's context is still active, it returns nil.
func (p *Pool[T]) Cause() error { return context.Cause(p.ctx) }

// Reset restores the pool to its initial state: cancels all tasks, clears results and metrics,
// and creates a new context.
func (p *Pool[T]) Reset() error {
	if p.ActiveTasks() > 0 {
		return errors.New("cannot reset while tasks are running")
	}

	p.Cancel()
	p.ClearMetrics()
	p.results.Clear()
	p.tokens = nil
	p.ctx, p.cancel = context.WithCancelCause(context.Background())

	return nil
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

func epanic(r any) error {
	stack := debug.Stack()
	switch x := r.(type) {
	case string:
		return fmt.Errorf("panic: %s\n%s", x, stack)
	case error:
		return fmt.Errorf("panic: %w\n%s", x, stack)
	default:
		return fmt.Errorf("panic: %v\n%s", x, stack)
	}
}
