package g

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// All returns true only if fn returns true for every element.
// It stops early on the first false.
func (p SeqSlicePar[V]) All(fn func(V) bool) bool {
	all := true
	p.Range(func(v V) bool {
		if !fn(v) {
			all = false
			return false
		}
		return true
	})

	return all
}

// Any returns true if fn returns true for any element.
// It stops early on the first true.
func (p SeqSlicePar[V]) Any(fn func(V) bool) bool {
	_any := false
	p.Range(func(v V) bool {
		if fn(v) {
			_any = true
			return false
		}
		return true
	})

	return _any
}

// Chain concatenates this SeqSlicePar with others, preserving process and worker count.
func (p SeqSlicePar[V]) Chain(others ...SeqSlicePar[V]) SeqSlicePar[V] {
	seq := func(yield func(V) bool) {
		p.src(yield)
		for _, o := range others {
			o.src(yield)
		}
	}

	return SeqSlicePar[V]{
		src:     seq,
		workers: p.workers,
		process: p.process,
	}
}

// Collect gathers all processed elements into a Slice.
func (p SeqSlicePar[V]) Collect() Slice[V] {
	var result []V
	p.Range(func(v V) bool {
		result = append(result, v)
		return true
	})

	return result
}

// Count returns the total number of elements processed.
func (p SeqSlicePar[V]) Count() Int {
	var count Int
	p.Range(func(V) bool {
		count++
		return true
	})

	return count
}

// Exclude removes elements for which fn returns true, in parallel.
func (p SeqSlicePar[V]) Exclude(fn func(V) bool) SeqSlicePar[V] {
	return p.Filter(func(v V) bool { return !fn(v) })
}

// Filter retains only elements where fn returns true.
func (p SeqSlicePar[V]) Filter(fn func(V) bool) SeqSlicePar[V] {
	prev := p.process

	return SeqSlicePar[V]{
		src:     p.src,
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
	var result Option[V]
	p.Range(func(v V) bool {
		if fn(v) {
			result = Some(v)
			return false
		}
		return true
	})

	return result
}

// Flatten unpacks nested slices or arrays in the source, returning a flat parallel sequence.
func (p SeqSlicePar[V]) Flatten() SeqSlicePar[V] {
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

		p.src(func(v V) bool {
			if mid, ok := p.process(v); ok {
				return recurse(mid)
			}
			return true
		})
	}

	return SeqSlicePar[V]{
		src:     seq,
		workers: p.workers,
		process: func(v V) (V, bool) { return v, true },
	}
}

// Fold reduces all elements into a single value, using fn to accumulate results.
func (p SeqSlicePar[V]) Fold(init V, fn func(acc, v V) V) V {
	acc := init
	p.Range(func(v V) bool {
		acc = fn(acc, v)
		return true
	})

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
		src:     p.src,
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
		src:     p.src,
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
	left, right := make([]V, 0), make([]V, 0)
	p.Range(func(v V) bool {
		if fn(v) {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
		return true
	})

	return left, right
}

// Range applies fn to each processed element in parallel, stopping on false.
func (p SeqSlicePar[V]) Range(fn func(V) bool) {
	in := make(chan V)
	out := make(chan V)
	done := make(chan struct{})
	defer close(done)

	go func() {
		defer close(in)

		p.src(func(v V) bool {
			select {
			case in <- v:
				return true
			case <-done:
				return false
			}
		})
	}()

	var wg sync.WaitGroup

	for range p.workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				if v2, ok := p.process(v); ok {
					select {
					case out <- v2:
					case <-done:
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for v := range out {
		if !fn(v) {
			return
		}
	}
}

// Skip skips the first n elements.
func (p SeqSlicePar[V]) Skip(n Int) SeqSlicePar[V] {
	prev := p.process
	var cnt int64

	return SeqSlicePar[V]{
		src:     p.src,
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
		src:     p.src,
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
	seen := NewMapSafe[any, any]()

	return SeqSlicePar[V]{
		src:     p.src,
		workers: p.workers,
		process: func(v V) (V, bool) {
			if mid, ok := prev(v); ok {
				if _, loaded := seen.GetOrSet(mid, struct{}{}); loaded {
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
