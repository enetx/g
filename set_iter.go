package g

import "context"

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each element as it passes through.
func (iter *baseIterS[T]) Inspect(fn func(v T)) *inspectIterS[T] {
	return inspectS(iter, fn)
}

// Collect gathers all elements from the iterator into a Set.
func (iter *baseIterS[T]) Collect() Set[T] {
	set := NewSet[T]()

	for {
		next := iter.Next()
		if next.IsNone() {
			return set
		}

		set.Add(next.Some())
	}
}

// Chain concatenates the current iterator with other iterators, returning a new iterator.
//
// The function creates a new iterator that combines the elements of the current iterator
// with elements from the provided iterators in the order they are given.
//
// Params:
//
// - iterators ([]iterator[T]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - *chainIter[T]: A new iterator containing elements from the current iterator and the provided iterators.
//
// Example usage:
//
//	iter1 := g.SetOf(1, 2, 3).Iter()
//	iter2 := g.SetOf(4, 5, 6).Iter()
//	iter1.Chain(iter2).Collect().Print()
//
// Output: Set{3, 4, 5, 6, 1, 2} // The output order may vary as the Set type is not ordered.
//
// The resulting iterator will contain elements from both iterators.
func (iter *baseIterS[T]) Chain(iterators ...iteratorS[T]) *chainIterS[T] {
	return chainS(append([]iteratorS[T]{iter}, iterators...)...)
}

// ForEach iterates through all elements and applies the given function to each.
//
// The function applies the provided function to each element of the iterator.
//
// Params:
//
// - fn (func(T)): The function to apply to each element.
//
// Example usage:
//
//	iter := g.SetOf(1, 2, 3).Iter()
//	iter.ForEach(func(val T) {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	})
//
// The provided function will be applied to each element in the iterator.
func (iter *baseIterS[T]) ForEach(fn func(v T)) {
	for {
		next := iter.Next()
		if next.IsNone() {
			return
		}

		fn(next.Some())
	}
}

// The iteration will stop when the provided function returns false for an element.
func (iter *baseIterS[T]) Range(fn func(v T) bool) {
	for {
		next := iter.Next()
		if next.IsNone() || !fn(next.Some()) {
			return
		}
	}
}

// Skip returns a new iterator skipping the first n elements.
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
// - *skipIterS[T]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	iter := g.SetOf(1, 2, 3, 4, 5, 6).Iter()
//	iter.Skip(3).Collect().Print()
//
// Output: {4, 5, 6} // The output may vary as the Set type is not ordered.
//
// The resulting iterator will start after skipping the specified number of elements.
func (iter *baseIterS[T]) Skip(n uint) *skipIterS[T] {
	return skipS(iter, n)
}

// Filter returns a new iterator containing only the elements that satisfy the provided function.
//
// The function applies the provided function to each element of the iterator.
// If the function returns true for an element, that element is included in the resulting iterator.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the iterator
// to determine if it should be included in the result.
//
// Returns:
//
// - *filterIterS[T]: A new iterator containing the elements that satisfy the given condition.
//
// Example usage:
//
//	set := g.SetOf(1, 2, 3, 4, 5)
//	even := set.Iter().
//		Filter(
//			func(val int) bool {
//				return val%2 == 0
//			}).
//		Collect()
//	even.Print()
//
// Output: Set{2, 4} // The output order may vary as the Set type is not ordered.
//
// The resulting iterator will contain only the elements that satisfy the provided function.
func (iter *baseIterS[T]) Filter(fn func(v T) bool) *filterIterS[T] {
	return filterS(iter, fn)
}

// Exclude returns a new iterator excluding elements that satisfy the provided function.
//
// The function applies the provided function to each element of the iterator.
// If the function returns true for an element, that element is excluded from the resulting iterator.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the iterator
// to determine if it should be excluded from the result.
//
// Returns:
//
// - *filterIterS[T]: A new iterator containing the elements that do not satisfy the given condition.
//
// Example usage:
//
//	set := g.SetOf(1, 2, 3, 4, 5)
//	notEven := set.Iter().
//		Exclude(
//			func(val int) bool {
//				return val%2 == 0
//			}).
//		Collect()
//	notEven.Print()
//
// Output: Set{1, 3, 5} // The output order may vary as the Set type is not ordered.
//
// The resulting iterator will contain only the elements that do not satisfy the provided function.
func (iter *baseIterS[T]) Exclude(fn func(v T) bool) *filterIterS[T] {
	return excludeS(iter, fn)
}

// Map transforms each element in the iterator using the given function.
//
// The function creates a new iterator by applying the provided function to each element
// of the original iterator.
//
// Params:
//
// - fn (func(T) T): The function used to transform elements.
//
// Returns:
//
// - *mapIterS[T, T]: A new iterator containing elements transformed by the provided function.
//
// Example usage:
//
//	set := g.SetOf(1, 2, 3)
//	doubled := set.Iter().
//		Map(
//			func(val int) int {
//				return val * 2
//			}).
//		Collect()
//	doubled.Print()
//
// Output: Set{2, 4, 6} // The output order may vary as the Set type is not ordered.
//
// The resulting iterator will contain elements transformed by the provided function.
func (iter *baseIterS[T]) Map(fn func(v T) T) *mapIterS[T, T] {
	return transformS(iter, fn)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// liftS
type liftIterS[T comparable] struct {
	baseIterS[T]
	items  chan T
	cancel func()
}

func liftS[T comparable](hashmap map[T]struct{}) *liftIterS[T] {
	ctx, cancel := context.WithCancel(context.Background())

	iter := &liftIterS[T]{items: make(chan T), cancel: cancel}
	iter.baseIterS = baseIterS[T]{iter}

	go func() {
		defer close(iter.items)

		for k := range hashmap {
			select {
			case <-ctx.Done():
				return
			default:
				iter.items <- k
			}
		}
	}()

	return iter
}

func (iter *liftIterS[T]) Next() Option[T] {
	item, ok := <-iter.items
	if !ok {
		return None[T]()
	}

	return Some(item)
}

// Close stops the iteration and releases associated resources.
// It signals the iterator to stop processing items and waits for the
// completion of any ongoing operations. After calling Close, the iterator
// cannot be used for further iteration.
func (iter *liftIterS[T]) Close() {
	iter.cancel()
	<-iter.items
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inspect
type inspectIterS[T comparable] struct {
	baseIterS[T]
	iter      iterator[T]
	fn        func(T)
	exhausted bool
}

func inspectS[T comparable](iter iteratorS[T], fn func(T)) *inspectIterS[T] {
	iterator := &inspectIterS[T]{iter: iter, fn: fn}
	iterator.baseIterS = baseIterS[T]{iterator}

	return iterator
}

func (iter *inspectIterS[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
		return None[T]()
	}

	iter.fn(next.Some())

	return next
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// filter
type filterIterS[T comparable] struct {
	baseIterS[T]
	iter      iteratorS[T]
	fn        func(T) bool
	exhausted bool
}

func filterS[T comparable](iter iteratorS[T], fn func(T) bool) *filterIterS[T] {
	iterator := &filterIterS[T]{iter: iter, fn: fn}
	iterator.baseIterS = baseIterS[T]{iterator}

	return iterator
}

func (iter *filterIterS[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for {
		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true
			return None[T]()
		}

		if iter.fn(next.Some()) {
			return next
		}
	}
}

func excludeS[T comparable](iter iteratorS[T], fn func(T) bool) *filterIterS[T] {
	inverse := func(t T) bool { return !fn(t) }
	iterator := &filterIterS[T]{iter: iter, fn: inverse}
	iterator.baseIterS = baseIterS[T]{iterator}

	return iterator
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// map
type mapIterS[T, U comparable] struct {
	baseIterS[U]
	iter      iterator[T]
	fn        func(T) U
	exhausted bool
}

func mapiterS[T, U comparable](iter iterator[T], fn func(T) U) *mapIterS[T, U] {
	iterator := &mapIterS[T, U]{iter: iter, fn: fn}
	iterator.baseIterS = baseIterS[U]{iterator}

	return iterator
}

func transformS[T comparable](iter iterator[T], fn func(T) T) *mapIterS[T, T] {
	return mapiterS(iter, fn)
}

func (iter *mapIterS[T, U]) Next() Option[U] {
	if iter.exhausted {
		return None[U]()
	}

	next := iter.iter.Next()

	if next.IsNone() {
		iter.exhausted = true
		return None[U]()
	}

	return Some(iter.fn(next.Some()))
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// chain
type chainIterS[T comparable] struct {
	baseIterS[T]
	iterators     []iteratorS[T]
	iteratorIndex int
}

func chainS[T comparable](iterators ...iteratorS[T]) *chainIterS[T] {
	iter := &chainIterS[T]{iterators: iterators}
	iter.baseIterS = baseIterS[T]{iter}
	return iter
}

func (iter *chainIterS[T]) Next() Option[T] {
	for {
		if iter.iteratorIndex == len(iter.iterators) {
			return None[T]()
		}

		if next := iter.iterators[iter.iteratorIndex].Next(); next.IsSome() {
			return next
		}

		iter.iteratorIndex++
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// difference
type differenceIterS[T comparable] struct {
	baseIterS[T]
	other     Set[T]
	iter      iteratorS[T]
	exhausted bool
}

func differenceS[T comparable](iter iteratorS[T], other Set[T]) *differenceIterS[T] {
	iterator := &differenceIterS[T]{iter: iter, other: other}
	iterator.baseIterS = baseIterS[T]{iterator}

	return iterator
}

func (iter *differenceIterS[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for {
		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true
			return None[T]()
		}

		if !iter.other.Contains(next.Some()) {
			return next
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// intersection
type intersectionIterS[T comparable] struct {
	baseIterS[T]
	other     Set[T]
	iter      iteratorS[T]
	exhausted bool
}

func intersectionS[T comparable](iter iteratorS[T], other Set[T]) *intersectionIterS[T] {
	iterator := &intersectionIterS[T]{iter: iter, other: other}
	iterator.baseIterS = baseIterS[T]{iterator}

	return iterator
}

func (iter *intersectionIterS[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for {
		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true
			return None[T]()
		}

		if iter.other.Contains(next.Some()) {
			return next
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// skip
type skipIterS[T comparable] struct {
	baseIterS[T]
	iter      iteratorS[T]
	count     uint
	skipped   bool
	exhausted bool
}

func skipS[T comparable](iter iteratorS[T], count uint) *skipIterS[T] {
	iterator := &skipIterS[T]{iter: iter, count: count}
	iterator.baseIterS = baseIterS[T]{iterator}

	return iterator
}

func (iter *skipIterS[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	if !iter.skipped {
		iter.skipped = true

		for i := uint(0); i < iter.count; i++ {
			if iter.delegateNext().IsNone() {
				return None[T]()
			}
		}
	}

	return iter.delegateNext()
}

func (iter *skipIterS[T]) delegateNext() Option[T] {
	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
	}

	return next
}
