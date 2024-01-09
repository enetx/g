package g

// Chain creates a new iterator by concatenating the current iterator with other iterators.
//
// The function concatenates the key-value pairs from the current iterator with the key-value pairs from the provided iterators,
// producing a new iterator containing all concatenated elements.
//
// Params:
//
// - iterators ([]iteratorM[K, V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - *chainIterM[K, V]: A new iterator containing elements from the current iterator and the provided iterators.
//
// Example usage:
//
//	iter1 := g.NewMap[int, string]().Set(1, "a").Iter()
//	iter2 := g.NewMap[int, string]().Set(2, "b").Iter()
//
//	// Concatenating iterators and collecting the result.
//	iter1.Chain(iter2).Collect().Print()
//
// Output: Map{1:a, 2:b} // The output order may vary as Map is not ordered.
//
// The resulting iterator will contain elements from both iterators.
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
//
// This function creates a new iterator starting from the (n+1)th key-value pair of the current iterator,
// excluding the first n key-value pairs.
//
// Params:
//
// - n (uint): The number of key-value pairs to skip from the beginning of the iterator.
//
// Returns:
//
// - *dropIterM[K, V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	iter := g.NewMap[int, string]().Set(1, "a").Set(2, "b").Set(3, "c").Set(4, "d").Iter()
//
//	// Skipping the first two elements and collecting the rest.
//	iter.Drop(2).Collect().Print()
//
// Output: Map{3:c, 4:d} // The output may vary as Map is not ordered.
//
// The resulting iterator will start after skipping the specified number of key-value pairs.
func (iter *baseIterM[K, V]) Drop(n uint) *dropIterM[K, V] {
	return dropM[K, V](iter, n)
}

// Exclude returns a new iterator excluding elements that satisfy the provided function.
//
// This function creates a new iterator excluding key-value pairs for which the provided function returns true.
// It iterates through the current iterator, applying the function to each key-value pair.
// If the function returns true for a key-value pair, it will be excluded from the resulting iterator.
//
// Params:
//
// - fn (func(K, V) bool): The function applied to each key-value pair to determine exclusion.
//
// Returns:
//
// - *filterIterM[K, V]: An iterator excluding elements that satisfy the given function.
//
// Example usage:
//
//	m := g.NewMap[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
//
//	notEven := m.Iter().
//		Exclude(
//			func(k, v int) bool {
//				return v%2 == 0
//			}).
//		Collect()
//	notEven.Print()
//
// Output: Map{1:1, 3:3, 5:5} // The output order may vary as Map is not ordered.
//
// The resulting iterator will exclude elements for which the function returns true.
func (iter *baseIterM[K, V]) Exclude(fn func(K, V) bool) *filterIterM[K, V] {
	return excludeM[K, V](iter, fn)
}

// Filter returns a new iterator containing only the elements that satisfy the provided function.
//
// This function creates a new iterator containing key-value pairs for which the provided function returns true.
// It iterates through the current iterator, applying the function to each key-value pair.
// If the function returns true for a key-value pair, it will be included in the resulting iterator.
//
// Params:
//
// - fn (func(K, V) bool): The function applied to each key-value pair to determine inclusion.
//
// Returns:
//
// - *filterIterM[K, V]: An iterator containing elements that satisfy the given function.
//
//	m := g.NewMap[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
//
//	even := m.Iter().
//		Filter(
//			func(k, v int) bool {
//				return v%2 == 0
//			}).
//		Collect()
//	even.Print()
//
// Output: Map{2:2, 4:4} // The output order may vary as Map is not ordered.
//
// The resulting iterator will contain elements for which the function returns true.
func (iter *baseIterM[K, V]) Filter(fn func(K, V) bool) *filterIterM[K, V] {
	return filterM[K, V](iter, fn)
}

// ForEach iterates through all elements and applies the given function to each key-value pair.
//
// This function traverses the entire iterator and applies the provided function to each key-value pair.
// It iterates through the current iterator, executing the function on each key-value pair.
//
// Params:
//
// - fn (func(K, V)): The function to be applied to each key-value pair in the iterator.
//
// Example usage:
//
//	m := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
//
//	mmap := m.Iter().
//		Map(
//			func(k, v int) (int, int) {
//				return k * k, v * v
//			}).
//		Collect()
//
//	mmap.Print()
//
// Output: Map{1:1, 4:4, 9:9, 16:16, 25:25} // The output order may vary as Map is not ordered.
//
// The function fn will be executed for each key-value pair in the iterator.
func (iter *baseIterM[K, V]) ForEach(fn func(K, V)) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		fn(next.Some().Key, next.Some().Value)
	}
}

// Map creates a new iterator by applying the given function to each key-value pair.
//
// This function generates a new iterator by traversing the current iterator and applying the provided
// function to each key-value pair. It transforms the key-value pairs according to the given function.
//
// Params:
//
//   - fn (func(K, V) (K, V)): The function to be applied to each key-value pair in the iterator.
//     It takes a key-value pair and returns a new transformed key-value pair.
//
// Returns:
//
// - *mapIterM[K, V]: A new iterator containing key-value pairs transformed by the provided function.
//
// Example usage:
//
//	m := g.NewMapOrd[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
//
//	mmap := m.Iter().
//		Map(
//			func(k, v int) (int, int) {
//				return k * k, v * v
//			}).
//		Collect()
//
//	mmap.Print()
//
// Output: Map{1:1, 4:4, 9:9, 16:16, 25:25} // The output order may vary as Map is not ordered.
//
// The resulting iterator will contain key-value pairs transformed by the given function.
func (iter *baseIterM[K, V]) Map(fn func(K, V) (K, V)) *mapIterM[K, V] {
	return mapiterM(iter, fn)
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
