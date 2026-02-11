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
	tokens chan g.Unit
	done   chan g.Unit
	once   sync.Once
}

// newLimiter creates a rate limiter that allows n operations per duration d,
// with an initial burst capacity.
//
// Internally it uses a channel as a token bucket:
//   - A background goroutine refills tokens at interval d/n
//   - The bucket (channel) is pre-filled with burst tokens
//   - Consumers call wait() to take a token before proceeding
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

	// Refill goroutine: adds one token every d/n.
	interval := d / time.Duration(n)
	if interval < time.Microsecond {
		interval = time.Microsecond
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				select {
				case l.tokens <- g.Unit{}:
				default: // bucket full, discard
				}
			case <-l.done:
				return
			}
		}
	}()

	return l
}

// wait blocks until a token is available or the context is canceled.
// Returns true if a token was acquired, false on cancellation.
func (l *limiter) wait(ctx context.Context) bool {
	select {
	case <-l.tokens:
		return true
	case <-ctx.Done():
		return false
	}
}

// stop shuts down the refill goroutine. Safe to call multiple times.
func (l *limiter) stop() {
	l.once.Do(func() { close(l.done) })
}
