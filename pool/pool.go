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

// ErrAllTasksDone is the cancellation cause set when all tasks have completed
// and the pool shuts down normally via Wait or Stream.
var ErrAllTasksDone = errors.New("all tasks completed")

// Pool[T any] is a goroutine pool that allows parallel task execution.
//
// It supports two consumption modes:
//   - Wait: blocks until all tasks finish, returns collected results.
//     Each Go call spawns a goroutine, limited by a semaphore (tokens).
//   - Stream: returns a channel that emits results in real-time.
//     Uses a fixed worker pool — goroutine count equals the concurrency limit,
//     regardless of how many tasks are submitted.
//
// These modes are mutually exclusive — use one or the other per pool lifecycle.
type Pool[T any] struct {
	ctx           context.Context          // Context for controlling cancellation and timeouts
	cancel        context.CancelCauseFunc  // Function to cancel the context
	tokens        chan Unit                // Semaphore for Wait mode concurrency; nil means unlimited
	results       *MapSafe[int, Result[T]] // Stores task results (Wait mode only)
	wg            sync.WaitGroup           // Waits for goroutines (Wait mode) or workers (Stream mode)
	stream        chan Result[T]           // Real-time result channel (Stream mode only)
	jobs          chan func() Result[T]    // Task queue for worker pool (Stream mode only)
	streaming     atomic.Bool              // Whether the pool is in stream mode (race-free flag)
	totalTasks    int32                    // Total number of tasks submitted
	activeTasks   int32                    // Number of currently active tasks
	failedTasks   int32                    // Number of failed tasks
	cancelOnError bool                     // Cancels remaining tasks if any task fails
}

// New creates a new goroutine pool.
//
// Example (Wait mode):
//
//	p := pool.New[int]()
//	p.Go(func() Result[int] { return Ok(42) })
//	results := p.Wait().Collect()
//
// Example (Stream mode):
//
//	p := pool.New[int]().Limit(5)
//	ch := p.Stream(func() {
//	    for i := range 100 {
//	        p.Go(func() Result[int] { return Ok(i) })
//	    }
//	})
//	for r := range ch {
//	    fmt.Println(r.Ok())
//	}
func New[T any]() *Pool[T] {
	ctx, cancel := context.WithCancelCause(context.Background())

	return &Pool[T]{
		ctx:     ctx,
		cancel:  cancel,
		results: NewMapSafe[int, Result[T]](),
	}
}

// acquire blocks until a worker slot is available or the context is canceled.
// Used only in Wait mode.
func (p *Pool[T]) acquire() bool {
	if p.tokens == nil {
		return p.ctx.Err() == nil
	}

	select {
	case p.tokens <- Unit{}:
		return true
	case <-p.ctx.Done():
		return false
	}
}

// release frees a worker slot and signals task completion.
// Used only in Wait mode; called exactly once per acquired task.
func (p *Pool[T]) release() {
	defer atomic.AddInt32(&p.activeTasks, -1)

	if p.tokens != nil {
		<-p.tokens
	}

	p.wg.Done()
}

// error records a task failure. It delivers the error result
// BEFORE canceling the context, ensuring the first error is never lost
// when CancelOnError is enabled.
func (p *Pool[T]) error(index int32, err error) {
	atomic.AddInt32(&p.failedTasks, 1)

	result := Err[T](err)

	if p.streaming.Load() {
		p.emit(result)
	} else {
		p.results.Insert(int(index), result)
	}

	if p.cancelOnError {
		p.Cancel(err)
	}
}

// emit sends a result to the stream channel.
// If the context is canceled, the result is silently dropped —
// this is expected with CancelOnError, where only the first error
// is guaranteed to be delivered.
func (p *Pool[T]) emit(result Result[T]) {
	select {
	case p.stream <- result:
	case <-p.ctx.Done():
	}
}

// workers returns the number of workers for Stream mode.
// Uses Limit value if set, otherwise defaults to GOMAXPROCS.
func (p *Pool[T]) workers() int {
	if p.tokens != nil {
		return cap(p.tokens)
	}

	return runtime.GOMAXPROCS(0)
}

// worker is a long-lived goroutine that processes tasks from the jobs channel.
// It exits when the jobs channel is closed or the context is canceled.
func (p *Pool[T]) worker() {
	defer p.wg.Done()

	for {
		select {
		case fn, ok := <-p.jobs:
			if !ok {
				return
			}

			p.exec(fn)
		case <-p.ctx.Done():
			return
		}
	}
}

// exec runs a single task within a worker goroutine.
// Handles nil functions, panics, errors, and context cancellation.
func (p *Pool[T]) exec(fn func() Result[T]) {
	atomic.AddInt32(&p.activeTasks, 1)
	defer atomic.AddInt32(&p.activeTasks, -1)

	if fn == nil {
		p.error(-1, errors.New("nil function provided"))
		return
	}

	defer func() {
		if r := recover(); r != nil {
			p.error(-1, epanic(r))
		}
	}()

	if p.ctx.Err() != nil {
		return
	}

	result := fn()
	if result.IsErr() {
		p.error(-1, result.Err())
		return
	}

	p.emit(result)
}

// Go submits a task for execution.
//
// In Wait mode, Go blocks the caller until a worker slot is available
// (backpressure via semaphore), then spawns a goroutine to execute the task.
//
// In Stream mode, Go sends the task to the worker pool's job queue.
// It blocks if all workers are busy (backpressure via channel).
//
// If fn is nil, the task completes with an error.
// If fn panics, the panic is recovered and recorded as an error with a stack trace.
func (p *Pool[T]) Go(fn func() Result[T]) {
	if p.streaming.Load() {
		select {
		case p.jobs <- fn:
			atomic.AddInt32(&p.totalTasks, 1)
		case <-p.ctx.Done():
		}

		return
	}

	if !p.acquire() {
		return
	}

	index := atomic.AddInt32(&p.totalTasks, 1) - 1
	atomic.AddInt32(&p.activeTasks, 1)
	p.wg.Add(1)

	go func(index int32) {
		defer p.release()

		if fn == nil {
			p.error(index, errors.New("nil function provided"))
			return
		}

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

		p.results.Insert(int(index), result)
	}(index)
}

// Stream spawns a fixed worker pool, runs fn to submit tasks, and returns
// a channel that emits results as each task completes.
//
// fn is executed in a separate goroutine. Inside fn, call Go to submit tasks.
// When fn returns, the job queue is closed automatically and workers drain
// remaining tasks. The channel closes once all workers finish.
//
// The number of worker goroutines equals the Limit (or GOMAXPROCS if no limit
// is set). Memory usage is constant regardless of how many tasks are submitted.
//
// An optional buffer size prevents slow consumers from blocking workers.
// With CancelOnError, the first error is guaranteed to be delivered;
// subsequent results may be dropped after cancellation.
//
// Stream and Wait are mutually exclusive.
//
// Example:
//
//	p := pool.New[int]().Limit(5)
//	ch := p.Stream(func() {
//	    for i := range 100 {
//	        p.Go(func() Result[int] { return Ok(i * i) })
//	    }
//	})
//
//	for r := range ch {
//	    fmt.Println(r.Ok())
//	}
func (p *Pool[T]) Stream(fn func(), buffer ...int) <-chan Result[T] {
	if atomic.LoadInt32(&p.totalTasks) > 0 {
		panic("Stream must be called before submitting tasks with Go")
	}

	buf := 0
	if len(buffer) > 0 && buffer[0] > 0 {
		buf = buffer[0]
	}

	p.stream = make(chan Result[T], buf)
	p.jobs = make(chan func() Result[T])
	p.streaming.Store(true)

	workers := p.workers()
	for range workers {
		p.wg.Add(1)
		go p.worker()
	}

	// Producer: runs fn to submit tasks, then closes job queue.
	// When fn returns (all Go calls done), close(p.jobs) signals
	// workers to drain remaining tasks and exit.
	go func() {
		defer close(p.jobs)
		fn()
	}()

	// Closer: waits for all workers to finish, then closes stream channel.
	go func() {
		defer close(p.stream)
		defer p.Cancel(ErrAllTasksDone)
		p.wg.Wait()
	}()

	return p.stream
}

// Wait blocks until all submitted tasks finish and returns their results.
// Results are returned in an iterator; order is not guaranteed.
//
// Example:
//
//	p := pool.New[int]()
//	p.Go(func() Result[int] { return Ok(1) })
//	p.Go(func() Result[int] { return Ok(2) })
//	for r := range p.Wait() {
//	    fmt.Println(r.Ok())
//	}
func (p *Pool[T]) Wait() SeqResult[T] {
	defer p.Cancel(ErrAllTasksDone)

	done := make(chan Unit)

	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-p.ctx.Done():
	}

	return SeqResult[T](p.results.Iter().Values())
}

// Limit sets the maximum number of concurrently running tasks.
// In Wait mode, this controls the semaphore size.
// In Stream mode, this determines the number of worker goroutines.
// Zero or negative values remove the limit (unlimited in Wait mode,
// GOMAXPROCS workers in Stream mode).
// Cannot be changed while tasks are running.
//
// On non-Windows systems, the limit is capped by the process stack rlimit.
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

	p.tokens = make(chan Unit, workers)

	return p
}

// CancelOnError configures the pool to cancel all remaining tasks
// when any task returns an error or panics.
//
// In Stream mode, the triggering error is guaranteed to be delivered;
// subsequent results may be dropped.
func (p *Pool[T]) CancelOnError() *Pool[T] {
	p.cancelOnError = true
	return p
}

// Context replaces the pool's context with the provided context.
// If ctx is nil, context.Background() is used.
// The previous context is canceled before replacement.
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

// Cancel cancels all tasks in the pool. An optional error provides the
// cancellation cause, retrievable via Cause. Defaults to context.Canceled.
func (p *Pool[T]) Cancel(err ...error) {
	cause := context.Canceled
	if len(err) > 0 && err[0] != nil {
		cause = err[0]
	}

	p.cancel(cause)
}

// Cause returns the reason for the pool's cancellation.
// Returns nil if the pool has not been canceled.
func (p *Pool[T]) Cause() error { return context.Cause(p.ctx) }

// Reset restores the pool to its initial state for reuse.
// Returns an error if tasks are still running.
//
// Example:
//
//	p.Wait() // or drain Stream channel
//	p.Reset()
//	p.Go(func() Result[int] { return Ok(1) })
//	p.Wait()
func (p *Pool[T]) Reset() error {
	if p.ActiveTasks() > 0 {
		return errors.New("cannot reset while tasks are running")
	}

	p.Cancel()
	p.ClearMetrics()
	p.results.Clear()
	p.tokens = nil
	p.stream = nil
	p.jobs = nil
	p.streaming.Store(false)
	p.ctx, p.cancel = context.WithCancelCause(context.Background())

	return nil
}

// ClearMetrics resets total and failed task counters to zero.
func (p *Pool[T]) ClearMetrics() {
	atomic.StoreInt32(&p.totalTasks, 0)
	atomic.StoreInt32(&p.failedTasks, 0)
}

// TotalTasks returns the total number of tasks submitted.
func (p *Pool[T]) TotalTasks() int { return int(atomic.LoadInt32(&p.totalTasks)) }

// ActiveTasks returns the number of tasks currently running.
func (p *Pool[T]) ActiveTasks() int { return int(atomic.LoadInt32(&p.activeTasks)) }

// FailedTasks returns the number of tasks that completed with an error.
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
