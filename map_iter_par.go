package g

import (
	"sync"
	"sync/atomic"
)

// Range applies fn to each processed pair in parallel, stopping early if fn returns false.
func (p SeqMapPar[K, V]) Range(fn func(K, V) bool) {
	in := make(chan Pair[K, V])
	out := make(chan Pair[K, V])
	done := make(chan struct{})
	defer close(done)

	go func() {
		defer close(in)

		p.src(func(k K, v V) bool {
			select {
			case in <- Pair[K, V]{k, v}:
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
			for pair := range in {
				if y, ok := p.process(pair); ok {
					select {
					case out <- y:
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

	for pair := range out {
		if !fn(pair.Key, pair.Value) {
			return
		}
	}
}

// Collect gathers all processed pairs into a Map.
func (p SeqMapPar[K, V]) Collect() Map[K, V] {
	m := NewMap[K, V]()
	p.Range(func(k K, v V) bool {
		m.Set(k, v)
		return true
	})

	return m
}

// Count returns the total number of processed pairs.
func (p SeqMapPar[K, V]) Count() Int {
	var cnt Int
	p.Range(func(_ K, _ V) bool {
		cnt++
		return true
	})

	return cnt
}

// Filter retains only pairs where fn returns true.
func (p SeqMapPar[K, V]) Filter(fn func(K, V) bool) SeqMapPar[K, V] {
	prev := p.process

	return SeqMapPar[K, V]{
		src:     p.src,
		workers: p.workers,
		process: func(pair Pair[K, V]) (Pair[K, V], bool) {
			if mid, ok := prev(pair); ok && fn(mid.Key, mid.Value) {
				return mid, true
			}
			return Pair[K, V]{}, false
		},
	}
}

// Exclude removes pairs where fn returns true.
func (p SeqMapPar[K, V]) Exclude(fn func(K, V) bool) SeqMapPar[K, V] {
	return p.Filter(func(k K, v V) bool { return !fn(k, v) })
}

// Map applies transform to each pair.
func (p SeqMapPar[K, V]) Map(transform func(K, V) (K, V)) SeqMapPar[K, V] {
	prev := p.process

	return SeqMapPar[K, V]{
		src:     p.src,
		workers: p.workers,
		process: func(pair Pair[K, V]) (Pair[K, V], bool) {
			if mid, ok := prev(pair); ok {
				k2, v2 := transform(mid.Key, mid.Value)
				return Pair[K, V]{k2, v2}, true
			}
			return Pair[K, V]{}, false
		},
	}
}

// Find returns the first pair matching fn, or a zero Option if none.
func (p SeqMapPar[K, V]) Find(fn func(K, V) bool) Option[Pair[K, V]] {
	var result Option[Pair[K, V]]
	p.Range(func(k K, v V) bool {
		if fn(k, v) {
			result = Some(Pair[K, V]{k, v})
			return false
		}
		return true
	})

	return result
}

// All returns true if fn returns true for every pair.
func (p SeqMapPar[K, V]) All(fn func(K, V) bool) bool {
	all := true
	p.Range(func(k K, v V) bool {
		if !fn(k, v) {
			all = false
			return false
		}
		return true
	})

	return all
}

// Any returns true if fn returns true for any pair.
func (p SeqMapPar[K, V]) Any(fn func(K, V) bool) bool {
	_any := false
	p.Range(func(k K, v V) bool {
		if fn(k, v) {
			_any = true
			return false
		}
		return true
	})

	return _any
}

// Chain concatenates this SeqMapPar with others, preserving process and workers.
func (p SeqMapPar[K, V]) Chain(others ...SeqMapPar[K, V]) SeqMapPar[K, V] {
	seq := func(yield func(K, V) bool) {
		p.src(yield)
		for _, o := range others {
			o.src(yield)
		}
	}

	return SeqMapPar[K, V]{
		src:     seq,
		workers: p.workers,
		process: p.process,
	}
}

// Skip drops the first n pairs.
func (p SeqMapPar[K, V]) Skip(n Int) SeqMapPar[K, V] {
	prev := p.process
	var cnt int64

	return SeqMapPar[K, V]{
		src:     p.src,
		workers: p.workers,
		process: func(pair Pair[K, V]) (Pair[K, V], bool) {
			if mid, ok := prev(pair); ok {
				if atomic.AddInt64(&cnt, 1) > int64(n) {
					return mid, true
				}
			}
			return Pair[K, V]{}, false
		},
	}
}

// Take yields at most n pairs.
func (p SeqMapPar[K, V]) Take(n Int) SeqMapPar[K, V] {
	prev := p.process
	var cnt int64

	return SeqMapPar[K, V]{
		src:     p.src,
		workers: p.workers,
		process: func(pair Pair[K, V]) (Pair[K, V], bool) {
			if mid, ok := prev(pair); ok {
				if atomic.AddInt64(&cnt, 1) <= int64(n) {
					return mid, true
				}
			}
			return Pair[K, V]{}, false
		},
	}
}
