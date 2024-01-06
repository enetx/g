package g

import "context"

func collectMO[K comparable, V any](iter iteratorMO[K, V]) *MapOrd[K, V] {
	mp := NewMapOrd[K, V]()

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		mp.Set(next.Some().Key, next.Some().Value)
	}

	return mp
}

func foreachMO[K comparable, V any](iter iteratorMO[K, V], fn func(K, V)) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		fn(next.Some().Key, next.Some().Value)
	}
}

func rangeMO[K comparable, V any](iter iteratorMO[K, V], fn func(K, V) bool) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if !fn(next.Some().Key, next.Some().Value) {
			break
		}
	}
}

func tochannelMO[K comparable, V any](iter iteratorMO[K, V], ctxs ...context.Context) chan pair[K, V] {
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

// Collect collects all key-value pairs from the iterator and returns a MapOrd.
func (iter *baseIterMO[K, V]) Collect() *MapOrd[K, V] {
	return collectMO[K, V](iter.iteratorMO)
}

// Drop returns a new iterator that skips the first n elements.
func (iter *baseIterMO[K, V]) Drop(n uint) *dropIterMO[K, V] {
	return dropMO[K, V](iter, n)
}

// ForEach iterates through all elements and applies the given function to each key-value pair.
func (iter *baseIterMO[K, V]) ForEach(fn func(K, V)) {
	foreachMO[K, V](iter, fn)
}

// Map creates a new iterator by applying the given function to each key-value pair.
func (iter *baseIterMO[K, V]) Map(fn func(K, V) (K, V)) *mapIterMO[K, V] {
	return mapiterMO(iter, fn)
}

// Range iterates through elements until the given function returns false.
func (iter *baseIterMO[K, V]) Range(fn func(K, V) bool) {
	rangeMO[K, V](iter, fn)
}

// Take returns a new iterator with the first n elements.
func (iter *baseIterMO[K, V]) Take(n uint) *takeIterMO[K, V] {
	return takeMO[K, V](iter, n)
}

// Filter returns a new iterator containing only the elements that satisfy the provided function.
func (iter *baseIterMO[K, V]) Filter(fn func(K, V) bool) *filterIterMO[K, V] {
	return filterMO[K, V](iter, fn)
}

// Exclude returns a new iterator excluding elements that satisfy the provided function.
func (iter *baseIterMO[K, V]) Exclude(fn func(K, V) bool) *filterIterMO[K, V] {
	return excludeMO[K, V](iter, fn)
}

// ToChannel converts the iterator into a channel, optionally with context(s).
func (iter *baseIterMO[K, V]) ToChannel(ctxs ...context.Context) chan pair[K, V] {
	return tochannelMO[K, V](iter, ctxs...)
}

// Chain creates a new iterator by concatenating the current iterator with other iterators.
func (iter *baseIterMO[K, V]) Chain(iterators ...iteratorMO[K, V]) *chainIterMO[K, V] {
	return chainMO[K, V](append([]iteratorMO[K, V]{iter}, iterators...)...)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// lift
type liftIterMO[K comparable, V any] struct {
	baseIterMO[K, V]
	items []pair[K, V]
	index int
}

func liftMO[K comparable, V any](items []pair[K, V]) *liftIterMO[K, V] {
	iterator := &liftIterMO[K, V]{items: items}
	iterator.baseIterMO = baseIterMO[K, V]{iterator}

	return iterator
}

func (iter *liftIterMO[K, V]) Next() Option[pair[K, V]] {
	if iter.index >= len(iter.items) {
		return None[pair[K, V]]()
	}

	iter.index++

	return Some(iter.items[iter.index-1])
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// map
type mapIterMO[K comparable, V any] struct {
	baseIterMO[K, V]
	iter      iteratorMO[K, V]
	fn        func(K, V) (K, V)
	exhausted bool
}

func mapiterMO[K comparable, V any](iter iteratorMO[K, V], fn func(K, V) (K, V)) *mapIterMO[K, V] {
	iterator := &mapIterMO[K, V]{iter: iter, fn: fn}
	iterator.baseIterMO = baseIterMO[K, V]{iterator}

	return iterator
}

func (iter *mapIterMO[K, V]) Next() Option[pair[K, V]] {
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
type filterIterMO[K comparable, V any] struct {
	baseIterMO[K, V]
	iter      iteratorMO[K, V]
	fn        func(K, V) bool
	exhausted bool
}

func filterMO[K comparable, V any](iter iteratorMO[K, V], fn func(K, V) bool) *filterIterMO[K, V] {
	iterator := &filterIterMO[K, V]{iter: iter, fn: fn}
	iterator.baseIterMO = baseIterMO[K, V]{iterator}

	return iterator
}

func (iter *filterIterMO[K, V]) Next() Option[pair[K, V]] {
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

func excludeMO[K comparable, V any](iter iteratorMO[K, V], fn func(K, V) bool) *filterIterMO[K, V] {
	inverse := func(k K, v V) bool { return !fn(k, v) }
	iterator := &filterIterMO[K, V]{iter: iter, fn: inverse}
	iterator.baseIterMO = baseIterMO[K, V]{iterator}

	return iterator
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// chain
type chainIterMO[K comparable, V any] struct {
	baseIterMO[K, V]
	iterators     []iteratorMO[K, V]
	iteratorIndex int
}

func chainMO[K comparable, V any](iterators ...iteratorMO[K, V]) *chainIterMO[K, V] {
	iter := &chainIterMO[K, V]{iterators: iterators}
	iter.baseIterMO = baseIterMO[K, V]{iter}
	return iter
}

func (iter *chainIterMO[K, V]) Next() Option[pair[K, V]] {
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
// take
type takeIterMO[K comparable, V any] struct {
	baseIterMO[K, V]
	iter  iteratorMO[K, V]
	limit uint
}

func takeMO[K comparable, V any](iter iteratorMO[K, V], limit uint) *takeIterMO[K, V] {
	iterator := &takeIterMO[K, V]{iter: iter, limit: limit}
	iterator.baseIterMO = baseIterMO[K, V]{iterator}

	return iterator
}

func (iter *takeIterMO[K, V]) Next() Option[pair[K, V]] {
	if iter.limit == 0 {
		return None[pair[K, V]]()
	}

	next := iter.iter.Next()
	if next.IsNone() {
		iter.limit = 0
	} else {
		iter.limit--
	}

	return next
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// drop

type dropIterMO[K comparable, V any] struct {
	baseIterMO[K, V]
	iter      iteratorMO[K, V]
	count     uint
	dropped   bool
	exhausted bool
}

func dropMO[K comparable, V any](iter iteratorMO[K, V], count uint) *dropIterMO[K, V] {
	iterator := &dropIterMO[K, V]{iter: iter, count: count}
	iterator.baseIterMO = baseIterMO[K, V]{iterator}

	return iterator
}

func (iter *dropIterMO[K, V]) Next() Option[pair[K, V]] {
	if iter.exhausted {
		return None[pair[K, V]]()
	}

	if !iter.dropped {
		iter.dropped = true

		for i := uint(0); i < iter.count; i++ {
			if iter.delegateNextMO().IsNone() {
				return None[pair[K, V]]()
			}
		}
	}

	return iter.delegateNextMO()
}

func (iter *dropIterMO[K, V]) delegateNextMO() Option[pair[K, V]] {
	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
	}

	return next
}
