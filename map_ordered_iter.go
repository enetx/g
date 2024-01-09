package g

import "context"

// Chain concatenates the current iterator with other iterators, returning a new iterator.
//
// The function creates a new iterator that combines the elements of the current iterator
// with elements from the provided iterators in the order they are given.
//
// Params:
//
// - iterators ([]iteratorMO[K, V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - *chainIterMO[K, V]: A new iterator containing elements from the current iterator and the provided iterators.
//
// Example usage:
//
//	iter1 := g.NewMapOrd[int, string]().Set(1, "a").Iter()
//	iter2 := g.NewMapOrd[int, string]().Set(2, "b").Iter()
//
//	// Concatenating iterators and collecting the result.
//	iter1.Chain(iter2).Collect().Print()
//
// Output: MapOrd{1:a, 2:b}
//
// The resulting iterator will contain elements from both iterators in the specified order.
func (iter *baseIterMO[K, V]) Chain(iterators ...iteratorMO[K, V]) *chainIterMO[K, V] {
	return chainMO[K, V](append([]iteratorMO[K, V]{iter}, iterators...)...)
}

// Collect collects all key-value pairs from the iterator and returns a MapOrd.
func (iter *baseIterMO[K, V]) Collect() *MapOrd[K, V] {
	mp := NewMapOrd[K, V]()

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		mp.Set(next.Some().Key, next.Some().Value)
	}

	return mp
}

// Drop returns a new iterator skipping the first n elements.
//
// The function creates a new iterator that skips the first n elements of the current iterator
// and returns an iterator starting from the (n+1)th element.
//
// Params:
//
// - n (uint): The number of elements to skip from the beginning of the iterator.
//
// Returns:
//
// - *dropIterMO[K, V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	iter := g.NewMapOrd[int, string]().Set(1, "a").Set(2, "b").Set(3, "c").Set(4, "d").Iter()
//
//	// Skipping the first two elements and collecting the rest.
//	iter.Drop(2).Collect().Print()
//
// Output: MapOrd{3:c, 4:d}
//
// The resulting iterator will start after skipping the specified number of elements.
func (iter *baseIterMO[K, V]) Drop(n uint) *dropIterMO[K, V] {
	return dropMO[K, V](iter, n)
}

// Exclude returns a new iterator excluding elements that satisfy the provided function.
//
// The function creates a new iterator excluding elements from the current iterator
// for which the provided function returns true.
//
// Params:
//
// - fn (func(K, V) bool): The function used to determine exclusion criteria for elements.
//
// Returns:
//
// - *filterIterMO[K, V]: A new iterator excluding elements that satisfy the given condition.
//
// Example usage:
//
//	mo := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
//
//	notEven := mo.Iter().
//		Exclude(
//			func(k, v int) bool {
//				return v%2 == 0
//			}).
//		Collect()
//	notEven.Print()
//
// Output: MapOrd{1:1, 3:3, 5:5}
//
// The resulting iterator will exclude elements based on the provided condition.
func (iter *baseIterMO[K, V]) Exclude(fn func(K, V) bool) *filterIterMO[K, V] {
	return excludeMO[K, V](iter, fn)
}

// Filter returns a new iterator containing only the elements that satisfy the provided function.
//
// The function creates a new iterator including elements from the current iterator
// for which the provided function returns true.
//
// Params:
//
// - fn (func(K, V) bool): The function used to determine inclusion criteria for elements.
//
// Returns:
//
// - *filterIterMO[K, V]: A new iterator containing elements that satisfy the given condition.
//
// Example usage:
//
//	mo := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
//
//	even := mo.Iter().
//		Filter(
//			func(k, v int) bool {
//				return v%2 == 0
//			}).
//		Collect()
//	even.Print()
//
// Output: MapOrd{2:2, 4:4}
//
// The resulting iterator will include elements based on the provided condition.
func (iter *baseIterMO[K, V]) Filter(fn func(K, V) bool) *filterIterMO[K, V] {
	return filterMO[K, V](iter, fn)
}

// ForEach iterates through all elements and applies the given function to each key-value pair.
//
// The function applies the provided function to each key-value pair in the iterator.
//
// Params:
//
// - fn (func(K, V)): The function to be applied to each key-value pair in the iterator.
//
// Example usage:
//
//	iter := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5).
//		Iter()
//
//	iter.ForEach(func(key K, val V) {
//	    // Process key-value pair
//	})
//
// The provided function will be applied to each key-value pair in the iterator.
func (iter *baseIterMO[K, V]) ForEach(fn func(K, V)) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		fn(next.Some().Key, next.Some().Value)
	}
}

// Map creates a new iterator by applying the given function to each key-value pair.
//
// The function creates a new iterator by applying the provided function to each key-value pair in the iterator.
//
// Params:
//
// - fn (func(K, V) (K, V)): The function used to transform each key-value pair in the iterator.
//
// Returns:
//
// - *mapIterMO[K, V]: A new iterator containing transformed key-value pairs.
//
// Example usage:
//
//	mo := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
//
//	momap := mo.Iter().
//		Map(
//			func(k, v int) (int, int) {
//				return k * k, v * v
//			}).
//		Collect()
//
//	momap.Print()
//
// Output: MapOrd{1:1, 4:4, 9:9, 16:16, 25:25}
//
// The resulting iterator will contain transformed key-value pairs.
func (iter *baseIterMO[K, V]) Map(fn func(K, V) (K, V)) *mapIterMO[K, V] {
	return mapiterMO(iter, fn)
}

// Range iterates through elements until the given function returns false.
//
// The function iterates through the key-value pairs in the iterator, applying the provided function to each pair.
// It continues iterating until the function returns false.
//
// Params:
//
// - fn (func(K, V) bool): The function to be applied to each key-value pair in the iterator.
//
// Example usage:
//
//	iter := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5).
//		Iter()
//
//	iter.Range(func(k, v int) bool {
//	    fmt.Println(v) // Replace this with the function logic you need.
//	    return v < 5 // Replace this with the condition for continuing iteration.
//	})
//
// The iteration will stop when the provided function returns false.
func (iter *baseIterMO[K, V]) Range(fn func(K, V) bool) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if !fn(next.Some().Key, next.Some().Value) {
			break
		}
	}
}

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (iter *baseIterMO[K, V]) Take(n uint) *takeIterMO[K, V] {
	return takeMO[K, V](iter, n)
}

// ToChannel converts the iterator into a channel, optionally with context(s).
//
// The function converts the key-value pairs from the iterator into a channel, allowing iterative processing
// using channels. It can be used to stream key-value pairs for concurrent or asynchronous operations.
//
// Params:
//
// - ctxs (...context.Context): Optional context(s) that can be used to cancel or set deadlines for the operation.
//
// Returns:
//
// - chan pair[K, V]: A channel emitting key-value pairs from the iterator.
//
// Example usage:
//
//	iter := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5).
//		Iter()
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensure cancellation to avoid goroutine leaks.
//
//	ch := iter.ToChannel(ctx)
//	for pair := range ch {
//	    // Process key-value pair from the channel
//	}
//
// The function converts the iterator into a channel to allow sequential or concurrent processing of key-value pairs.
func (iter *baseIterMO[K, V]) ToChannel(ctxs ...context.Context) chan pair[K, V] {
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
