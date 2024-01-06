package g

import "context"

// Chain creates a new iterator by concatenating the current iterator with other iterators.
func (iter *baseIterM[K, V]) Chain(iterators ...iteratorM[K, V]) *chainIterM[K, V] {
	return chainM[K, V](append([]iteratorM[K, V]{iter}, iterators...)...)
}

// Collect collects all key-value pairs from the iterator and returns a Map.
func (iter *baseIterM[K, V]) Collect() Map[K, V] {
	mp := NewMap[K, V]()

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		mp.Set(next.Some().Key, next.Some().Value)
	}

	return mp
}

// Drop returns a new iterator that skips the first n elements.
func (iter *baseIterM[K, V]) Drop(n uint) *dropIterM[K, V] {
	return dropM[K, V](iter, n)
}

// Exclude returns a new iterator excluding elements that satisfy the provided function.
func (iter *baseIterM[K, V]) Exclude(fn func(K, V) bool) *filterIterM[K, V] {
	return excludeM[K, V](iter, fn)
}

// Filter returns a new iterator containing only the elements that satisfy the provided function.
func (iter *baseIterM[K, V]) Filter(fn func(K, V) bool) *filterIterM[K, V] {
	return filterM[K, V](iter, fn)
}

// ForEach iterates through all elements and applies the given function to each key-value pair.
func (iter *baseIterM[K, V]) ForEach(fn func(K, V)) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		fn(next.Some().Key, next.Some().Value)
	}
}

// Map creates a new iterator by applying the given function to each key-value pair.
func (iter *baseIterM[K, V]) Map(fn func(K, V) (K, V)) *mapIterM[K, V] {
	return mapiterM(iter, fn)
}

// ToChannel converts the iterator into a channel, optionally with context(s).
func (iter *baseIterM[K, V]) ToChannel(ctxs ...context.Context) chan pair[K, V] {
	ch := make(chan pair[K, V])

	ctx := context.Background()
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	}

	go func() {
		defer close(ch)

		for next := iter.Next(); next.IsSome(); next = iter.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- next.Some()
			}
		}
	}()

	return ch
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// lift
type liftIterM[K comparable, V any] struct {
	baseIterM[K, V]
	items chan pair[K, V]
}

func liftM[K comparable, V any](hashmap map[K]V) *liftIterM[K, V] {
	iter := &liftIterM[K, V]{items: make(chan pair[K, V])}
	iter.baseIterM = baseIterM[K, V]{iter}

	go func() {
		defer close(iter.items)

		for k, v := range hashmap {
			iter.items <- pair[K, V]{k, v}
		}
	}()

	return iter
}

func (iter *liftIterM[K, V]) Next() Option[pair[K, V]] {
	item, ok := <-iter.items
	if !ok {
		return None[pair[K, V]]()
	}

	return Some(item)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// map
type mapIterM[K comparable, V any] struct {
	baseIterM[K, V]
	iter      iteratorM[K, V]
	fn        func(K, V) (K, V)
	exhausted bool
}

func mapiterM[K comparable, V any](iter iteratorM[K, V], fn func(K, V) (K, V)) *mapIterM[K, V] {
	iterator := &mapIterM[K, V]{iter: iter, fn: fn}
	iterator.baseIterM = baseIterM[K, V]{iterator}

	return iterator
}

func (iter *mapIterM[K, V]) Next() Option[pair[K, V]] {
	if iter.exhausted {
		return None[pair[K, V]]()
	}

	next := iter.iter.Next()

	if next.IsNone() {
		iter.exhausted = true
		return None[pair[K, V]]()
	}

	key, value := iter.fn(next.Some().Key, next.Some().Value)

	return Some(pair[K, V]{Key: key, Value: value})
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// filter
type filterIterM[K comparable, V any] struct {
	baseIterM[K, V]
	iter      iteratorM[K, V]
	fn        func(K, V) bool
	exhausted bool
}

func filterM[K comparable, V any](iter iteratorM[K, V], fn func(K, V) bool) *filterIterM[K, V] {
	iterator := &filterIterM[K, V]{iter: iter, fn: fn}
	iterator.baseIterM = baseIterM[K, V]{iterator}

	return iterator
}

func (iter *filterIterM[K, V]) Next() Option[pair[K, V]] {
	if iter.exhausted {
		return None[pair[K, V]]()
	}

	for next := iter.iter.Next(); next.IsSome(); next = iter.iter.Next() {
		if iter.fn(next.Some().Key, next.Some().Value) {
			return next
		}
	}

	iter.exhausted = true

	return None[pair[K, V]]()
}

func excludeM[K comparable, V any](iter iteratorM[K, V], fn func(K, V) bool) *filterIterM[K, V] {
	inverse := func(k K, v V) bool { return !fn(k, v) }
	iterator := &filterIterM[K, V]{iter: iter, fn: inverse}
	iterator.baseIterM = baseIterM[K, V]{iterator}

	return iterator
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// chain
type chainIterM[K comparable, V any] struct {
	baseIterM[K, V]
	iterators     []iteratorM[K, V]
	iteratorIndex int
}

func chainM[K comparable, V any](iterators ...iteratorM[K, V]) *chainIterM[K, V] {
	iter := &chainIterM[K, V]{iterators: iterators}
	iter.baseIterM = baseIterM[K, V]{iter}
	return iter
}

func (iter *chainIterM[K, V]) Next() Option[pair[K, V]] {
	for {
		if iter.iteratorIndex == len(iter.iterators) {
			return None[pair[K, V]]()
		}

		if next := iter.iterators[iter.iteratorIndex].Next(); next.IsSome() {
			return next
		}

		iter.iteratorIndex++
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// drop
type dropIterM[K comparable, V any] struct {
	baseIterM[K, V]
	iter      iteratorM[K, V]
	count     uint
	dropped   bool
	exhausted bool
}

func dropM[K comparable, V any](iter iteratorM[K, V], count uint) *dropIterM[K, V] {
	iterator := &dropIterM[K, V]{iter: iter, count: count}
	iterator.baseIterM = baseIterM[K, V]{iterator}

	return iterator
}

func (iter *dropIterM[K, V]) Next() Option[pair[K, V]] {
	if iter.exhausted {
		return None[pair[K, V]]()
	}

	if !iter.dropped {
		iter.dropped = true

		for i := uint(0); i < iter.count; i++ {
			if iter.delegateNextM().IsNone() {
				return None[pair[K, V]]()
			}
		}
	}

	return iter.delegateNextM()
}

func (iter *dropIterM[K, V]) delegateNextM() Option[pair[K, V]] {
	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
	}

	return next
}
