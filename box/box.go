package box

import "sync/atomic"

// Box is a lock-free, atomic, copy-on-write wrapper over *T.
// It enables safe concurrent access and updates via atomic pointer replacement.
type Box[T any] atomic.Pointer[T]

// NewBox creates a new Box wrapping the provided pointer.
//
// The value stored in the Box must follow an immutable pattern:
// it must never be mutated in-place after being stored.
//
// Instead, all updates must be performed by creating a new *T
// (typically via shallow or deep copy), modifying it, and replacing
// the pointer atomically using Update or Store.
//
//	Suitable field types for T:
//	 - primitives: int, bool, float64, string
//	 - arrays: [N]T
//	 - structs: only value-type fields
//	 - slices/maps: only if copied before modification
//
//	Unsafe field types (must NOT be mutated in-place):
//	 - slices: []T
//	 - maps: map[K]V
//	 - pointers: *U
//	 - channels, sync.Mutex, sync.Map, atomic.*
//
// Always copy such fields before modifying in Update or Store.
//
// Example:
//
//	type Config struct {
//	    Debug bool
//	    Port  int
//	}
//
//	box := NewBox(&Config{Port: 8080})
//
//	box.Update(func(c *Config) *Config {
//	    cp := *c       // shallow copy
//	    cp.Port = 9090 // safe update
//	    return &cp
//	})
//
// Never mutate the pointer returned by Load() or passed to Update() directly.
// Always copy before changing.
func New[T any](ptr *T) *Box[T] {
	b := new(atomic.Pointer[T])
	b.Store(ptr)
	return (*Box[T])(b)
}

// Load returns the current value stored in the Box.
//
// The returned pointer must not be mutated directly.
func (b *Box[T]) Load() *T { return (*atomic.Pointer[T])(b).Load() }

// Store replaces the current value atomically with the given pointer.
//
// The new value should ideally be a fresh copy and not shared elsewhere.
func (b *Box[T]) Store(value *T) { (*atomic.Pointer[T])(b).Store(value) }

// Update applies the given function to the current value and
// attempts to atomically replace it with the result.
//
// The apply function must return a new pointer (copied + modified).
// It must not mutate the original value.
//
// If another goroutine concurrently updates the value, Update will retry.
// If the returned pointer equals the current one, the update is skipped.
func (b *Box[T]) Update(apply func(current *T) *T) {
	for {
		if b.TryUpdate(apply) {
			break
		}
	}
}

// TryUpdate applies the given function to the current value and tries
// to replace it atomically using CompareAndSwap once.
//
// Returns true if the update succeeded, false if the value changed concurrently.
//
// Unlike Update, it does not retry on failure.
func (b *Box[T]) TryUpdate(apply func(current *T) *T) bool {
	p := (*atomic.Pointer[T])(b)

	current := p.Load()
	updated := apply(current)

	return updated == current || p.CompareAndSwap(current, updated)
}
