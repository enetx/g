package pool

import (
	"context"
	"sync"
	"time"

	"github.com/enetx/g"
)

// limiter implements a token-bucket rate limiter.
// Tokens refill at a fixed interval; burst controls the bucket size.
type limiter struct {
	tokens    chan g.Unit
	done      chan g.Unit
	interval  time.Duration // tick interval for the refill goroutine
	perTick   int           // tokens added per tick (>= 1)
	startOnce sync.Once     // lazily starts the refill goroutine on first wait()
	stopOnce  sync.Once     // closes done at most once
}

// newLimiter creates a rate limiter that allows n operations per duration d,
// with an initial burst capacity.
//
// Internally it uses a channel as a token bucket:
//   - The bucket (channel) is pre-filled with burst tokens, available immediately.
//   - A background goroutine refills tokens at interval d/n. It is started lazily
//     on the first wait() call, so a rate-configured pool that is never consumed
//     (no Wait/Stream/Go) does not leak a goroutine or a ticker.
//   - Consumers call wait() to take a token before proceeding.
//
// The tick interval is clamped to a 1µs minimum. When the requested rate exceeds
// one token per microsecond, perTick tokens are added each tick so the effective
// rate matches the request instead of silently capping at 1 token/µs.
func newLimiter(n int, d time.Duration, burst int) *limiter {
	if burst < 1 {
		burst = 1
	}

	l := &limiter{
		tokens: make(chan g.Unit, burst),
		done:   make(chan g.Unit),
	}

	// Pre-fill burst tokens — these are available immediately.
	for range burst {
		l.tokens <- g.Unit{}
	}

	// Compute the refill cadence. Refill goal is n tokens per duration d.
	interval := d / time.Duration(n)
	perTick := 1

	// Clamp the tick interval to a 1µs minimum to avoid a busy ticker. When the
	// requested rate is faster than one token per microsecond, add ceil(n / d_µs)
	// tokens per tick so the effective rate still matches the request instead of
	// silently capping at one token per microsecond.
	if interval < time.Microsecond {
		interval = time.Microsecond

		if dus := int64(d) / int64(time.Microsecond); dus > 0 {
			perTick = int((int64(n) + dus - 1) / dus) // ceil division
		}

		if perTick < 1 {
			perTick = 1
		}
	}

	l.interval = interval
	l.perTick = perTick

	return l
}

// start launches the refill goroutine. It is invoked at most once, on the first
// wait() call, via startOnce.
func (l *limiter) start() {
	go func() {
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				for range l.perTick {
					select {
					case l.tokens <- g.Unit{}:
					default: // bucket full, discard
					}
				}
			case <-l.done:
				return
			}
		}
	}()
}

// wait blocks until a token is available or the context is canceled.
// Returns true if a token was acquired, false on cancellation.
// The refill goroutine is started lazily on the first call.
func (l *limiter) wait(ctx context.Context) bool {
	l.startOnce.Do(l.start)

	select {
	case <-l.tokens:
		return true
	case <-ctx.Done():
		return false
	}
}

// stop shuts down the refill goroutine. Safe to call multiple times, and safe
// even if the goroutine was never started (the lazy start never fires after stop,
// and closing done is harmless to an unstarted goroutine).
func (l *limiter) stop() {
	l.stopOnce.Do(func() { close(l.done) })
}
