package g

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// All returns true only if fn returns true for every element.
// It stops early on the first false.
func (p SeqSlicePar[V]) All(fn func(V) bool) bool {
	var ok atomic.Bool
	ok.Store(true)

	p.Range(func(v V) bool {
		if !fn(v) {
			ok.Store(false)
			return false
		}
		return true
	})

	return ok.Load()
}

// Any returns true if fn returns true for any element.
// It stops early on the first true.
func (p SeqSlicePar[V]) Any(fn func(V) bool) bool {
	var ok atomic.Bool

	p.Range(func(v V) bool {
		if fn(v) {
			ok.Store(true)
			return false
		}
		return true
	})

	return ok.Load()
}

// Chain concatenates this SeqSlicePar with others, preserving process and worker count.
func (p SeqSlicePar[V]) Chain(others ...SeqSlicePar[V]) SeqSlicePar[V] {
	seq := func(yield func(V) bool) {
		p.seq(yield)
		for _, o := range others {
			o.seq(yield)
		}
	}

	return SeqSlicePar[V]{
		seq:     seq,
		workers: p.workers,
		process: p.process,
	}
}

// Collect gathers all processed elements into a Slice.
func (p SeqSlicePar[V]) Collect() Slice[V] {
	ch := make(chan V)

	go func() {
		defer close(ch)
		p.Range(func(v V) bool {
			ch <- v
			return true
		})
	}()

	var result []V
	for v := range ch {
		result = append(result, v)
	}

	return result
}

// Count returns the total number of elements processed.
func (p SeqSlicePar[V]) Count() Int {
	var count atomic.Int64
	p.Range(func(V) bool {
		count.Add(1)
		return true
	})

	return Int(count.Load())
}

// Exclude removes elements for which fn returns true, in parallel.
func (p SeqSlicePar[V]) Exclude(fn func(V) bool) SeqSlicePar[V] {
	return p.Filter(func(v V) bool { return !fn(v) })
}

// Filter retains only elements where fn returns true.
func (p SeqSlicePar[V]) Filter(fn func(V) bool) SeqSlicePar[V] {
	prev := p.process

	return SeqSlicePar[V]{
		seq:     p.seq,
		workers: p.workers,
		process: func(v V) (V, bool) {
			if mid, ok := prev(v); ok && fn(mid) {
				return mid, true
			}
			var zero V
			return zero, false
		},
	}
}

// Find returns the first element satisfying fn, or None if no such element exists.
func (p SeqSlicePar[V]) Find(fn func(V) bool) Option[V] {
	ch := make(chan V)

	go func() {
		defer close(ch)
		p.Range(func(v V) bool {
			if fn(v) {
				ch <- v
				return false
			}
			return true
		})
	}()

	if v, ok := <-ch; ok {
		return Some(v)
	}

	return None[V]()
}

// Flatten unpacks nested slices or arrays in the source, returning a flat parallel sequence.
func (p SeqSlicePar[V]) Flatten() SeqSlicePar[V] {
	prev := p.process

	seq := func(yield func(V) bool) {
		var recurse func(any) bool

		recurse = func(item any) bool {
			rv := reflect.ValueOf(item)
			switch rv.Kind() {
			case reflect.Slice, reflect.Array:
				for i := range rv.Len() {
					if !recurse(rv.Index(i).Interface()) {
						return false
					}
				}
			default:
				if v, ok := item.(V); ok {
					if !yield(v) {
						return false
					}
				}
			}
			return true
		}

		p.seq(func(v V) bool {
			if mid, ok := prev(v); ok {
				return recurse(mid)
			}
			return true
		})
	}

	return SeqSlicePar[V]{
		seq:     seq,
		workers: p.workers,
		process: func(v V) (V, bool) { return v, true },
	}
}

// Fold reduces all elements into a single value, using fn to accumulate results.
func (p SeqSlicePar[V]) Fold(init V, fn func(acc, v V) V) V {
	ch := make(chan V)

	go func() {
		defer close(ch)
		p.Range(func(v V) bool {
			ch <- v
			return true
		})
	}()

	acc := init
	for v := range ch {
		acc = fn(acc, v)
	}

	return acc
}

// ForEach applies fn to each element without early exit.
func (p SeqSlicePar[V]) ForEach(fn func(V)) {
	p.Range(func(v V) bool {
		fn(v)
		return true
	})
}

// Inspect invokes fn on each element without altering the resulting sequence.
func (p SeqSlicePar[V]) Inspect(fn func(V)) SeqSlicePar[V] {
	prev := p.process

	return SeqSlicePar[V]{
		seq:     p.seq,
		workers: p.workers,
		process: func(x V) (V, bool) {
			if mid, ok := prev(x); ok {
				fn(mid)
				return mid, true
			}
			var zero V
			return zero, false
		},
	}
}

// Map applies fn to each element.
func (p SeqSlicePar[V]) Map(fn func(V) V) SeqSlicePar[V] {
	prev := p.process

	return SeqSlicePar[V]{
		seq:     p.seq,
		workers: p.workers,
		process: func(v V) (V, bool) {
			if mid, ok := prev(v); ok {
				return fn(mid), true
			}
			var zero V
			return zero, false
		},
	}
}

// Partition splits elements into two slices: those satisfying fn, and the rest.
func (p SeqSlicePar[V]) Partition(fn func(V) bool) (Slice[V], Slice[V]) {
	leftCh := make(chan V)
	rightCh := make(chan V)

	go func() {
		defer close(leftCh)
		defer close(rightCh)
		p.Range(func(v V) bool {
			if fn(v) {
				leftCh <- v
			} else {
				rightCh <- v
			}
			return true
		})
	}()

	var left, right []V
	for leftCh != nil || rightCh != nil {
		select {
		case v, ok := <-leftCh:
			if !ok {
				leftCh = nil
				continue
			}
			left = append(left, v)
		case v, ok := <-rightCh:
			if !ok {
				rightCh = nil
				continue
			}
			right = append(right, v)
		}
	}

	return left, right
}

// Range applies fn to each processed element in parallel, stopping on false.
func (p SeqSlicePar[V]) Range(fn func(V) bool) {
	in := make(chan V, p.workers)
	var wg sync.WaitGroup
	var stop atomic.Bool

	go func() {
		defer close(in)
		p.seq(func(v V) bool {
			if stop.Load() {
				return false
			}
			in <- v
			return true
		})
	}()

	wg.Add(int(p.workers))
	for range p.workers {
		go func() {
			defer wg.Done()
			for v := range in {
				if mid, ok := p.process(v); ok {
					if !fn(mid) {
						stop.Store(true)
						return
					}
				}
			}
		}()
	}

	wg.Wait()
}

// Skip skips the first n elements.
func (p SeqSlicePar[V]) Skip(n Int) SeqSlicePar[V] {
	prev := p.process
	var cnt int64

	return SeqSlicePar[V]{
		seq:     p.seq,
		workers: p.workers,
		process: func(v V) (V, bool) {
			if mid, ok := prev(v); ok && atomic.AddInt64(&cnt, 1) > int64(n) {
				return mid, true
			}
			var zero V
			return zero, false
		},
	}
}

func (p SeqSlicePar[V]) Take(n Int) SeqSlicePar[V] {
	prev := p.process
	var cnt int64

	return SeqSlicePar[V]{
		seq:     p.seq,
		workers: p.workers,
		process: func(v V) (V, bool) {
			if mid, ok := prev(v); ok && atomic.AddInt64(&cnt, 1) <= int64(n) {
				return mid, true
			}
			var zero V
			return zero, false
		},
	}
}

// Unique removes duplicate elements, preserving the first occurrence.
func (p SeqSlicePar[V]) Unique() SeqSlicePar[V] {
	prev := p.process
	seen := NewMapSafe[any, struct{}]()

	return SeqSlicePar[V]{
		seq:     p.seq,
		workers: p.workers,
		process: func(v V) (V, bool) {
			if mid, ok := prev(v); ok {
				if loaded := seen.Entry(mid).OrSet(struct{}{}); loaded.IsSome() {
					var zero V
					return zero, false
				}

				return mid, true
			}

			var zero V
			return zero, false
		},
	}
}
