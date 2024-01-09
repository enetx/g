package g

import (
	"context"
)

// All checks if all elements in the iterator satisfy the given predicate.
func (iter *baseIter[T]) All(fn func(T) bool) bool {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if !fn(next.Some()) {
			return false
		}
	}

	return true
}

// Any checks if any element in the iterator satisfies the given predicate.
func (iter *baseIter[T]) Any(fn func(T) bool) bool {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if fn(next.Some()) {
			return true
		}
	}

	return false
}

// Chain concatenates the current iterator with other iterators, returning a new iterator.
func (iter *baseIter[T]) Chain(iterators ...iterator[T]) *chainIter[T] {
	return chain[T](append([]iterator[T]{iter}, iterators...)...)
}

// Collect gathers all elements from the iterator into a Slice.
func (iter *baseIter[T]) Collect() Slice[T] {
	values := make([]T, 0)

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		values = append(values, next.Some())
	}

	return values
}

// Cycle returns an iterator that endlessly repeats the elements of the current iterator.
func (iter *baseIter[T]) Cycle() *cycleIter[T] {
	return cycle[T](iter)
}

// Drop returns a new iterator skipping the first n elements.
func (iter *baseIter[T]) Drop(n uint) *dropIter[T] {
	return drop[T](iter, n)
}

// Enumerate adds an index to each element in the iterator.
func (iter *baseIter[T]) Enumerate() *enumerateIter[T] {
	return enumerate[T](iter)
}

// Exclude returns a new iterator excluding elements that satisfy the provided function.
func (iter *baseIter[T]) Exclude(fn func(T) bool) *filterIter[T] {
	return exclude[T](iter, fn)
}

// Filter returns a new iterator containing only the elements that satisfy the provided function.
func (iter *baseIter[T]) Filter(fn func(T) bool) *filterIter[T] {
	return filter[T](iter, fn)
}

// Find searches for an element in the iterator that satisfies the provided function.
func (iter *baseIter[T]) Find(fn func(T) bool) Option[T] {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if fn(next.Some()) {
			return next
		}
	}

	return None[T]()
}

// Flatten flattens an iterator of iterators into a single iterator.
func (iter *baseIter[T]) Flatten() *flattenIter[T] {
	return flatten[T](iter)
}

// Fold accumulates values in the iterator using a function.
func (iter *baseIter[T]) Fold(init T, fn func(T, T) T) T {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		init = fn(init, next.Some())
	}

	return init
}

// ForEach iterates through all elements and applies the given function to each.
func (iter *baseIter[T]) ForEach(fn func(T)) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		fn(next.Some())
	}
}

// Map transforms each element in the iterator using the given function.
func (iter *baseIter[T]) Map(fn func(T) T) *mapIter[T, T] {
	return transform[T](iter, fn)
}

// Range iterates through elements until the given function returns false.
func (iter *baseIter[T]) Range(fn func(T) bool) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if !fn(next.Some()) {
			break
		}
	}
}

// Take returns a new iterator with the first n elements.
func (iter *baseIter[T]) Take(n uint) *takeIter[T] {
	return take[T](iter, n)
}

// Unique returns an iterator with only unique elements.
func (iter *baseIter[T]) Unique() *uniqueIter[T] {
	return unique[T](iter)
}

// Chunks returns an iterator that yields chunks of elements of the specified size.
func (iter *baseIter[T]) Chunks(size int) *chunksIter[T] {
	return chunks[T](iter, size)
}

// Permutations generates iterators of all permutations of elements.
func (iter *baseIter[T]) Permutations() *permutationsIter[T] {
	return permutations[T](iter)
}

// Zip combines multiple iterators into a single iterator of tuples.
func (iter *baseIter[T]) Zip(iterators ...iterator[T]) *zipIter[T] {
	return zip[T](append([]iterator[T]{iter}, iterators...)...)
}

// ToChannel converts the iterator into a channel, optionally with context(s).
func (iter *baseIter[T]) ToChannel(ctxs ...context.Context) chan T {
	ch := make(chan T)

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
type liftIter[T any] struct {
	baseIter[T]
	items []T
	index int
}

func lift[T any](items []T) *liftIter[T] {
	iterator := &liftIter[T]{items: items}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *liftIter[T]) Next() Option[T] {
	if iter.index >= len(iter.items) {
		return None[T]()
	}

	iter.index++

	return Some(iter.items[iter.index-1])
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// map
type mapIter[T, U any] struct {
	baseIter[U]
	iter      iterator[T]
	fn        func(T) U
	exhausted bool
}

func mapiter[T, U any](iter iterator[T], fn func(T) U) *mapIter[T, U] {
	iterator := &mapIter[T, U]{iter: iter, fn: fn}
	iterator.baseIter = baseIter[U]{iterator}

	return iterator
}

func transform[T any](iter iterator[T], fn func(T) T) *mapIter[T, T] {
	return mapiter[T, T](iter, fn)
}

func (iter *mapIter[T, U]) Next() Option[U] {
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
// filter
type filterIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	fn        func(T) bool
	exhausted bool
}

func filter[T any](iter iterator[T], fn func(T) bool) *filterIter[T] {
	iterator := &filterIter[T]{iter: iter, fn: fn}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *filterIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for next := iter.iter.Next(); next.IsSome(); next = iter.iter.Next() {
		if iter.fn(next.Some()) {
			return next
		}
	}

	iter.exhausted = true

	return None[T]()
}

func exclude[T any](iter iterator[T], fn func(T) bool) *filterIter[T] {
	inverse := func(t T) bool { return !fn(t) }
	iterator := &filterIter[T]{iter: iter, fn: inverse}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// chain
type chainIter[T any] struct {
	baseIter[T]
	iterators     []iterator[T]
	iteratorIndex int
}

func chain[T any](iterators ...iterator[T]) *chainIter[T] {
	iter := &chainIter[T]{iterators: iterators}
	iter.baseIter = baseIter[T]{iter}
	return iter
}

func (iter *chainIter[T]) Next() Option[T] {
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
// take
type takeIter[T any] struct {
	baseIter[T]
	iter  iterator[T]
	limit uint
}

func take[T any](iter iterator[T], limit uint) *takeIter[T] {
	iterator := &takeIter[T]{iter: iter, limit: limit}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *takeIter[T]) Next() Option[T] {
	if iter.limit == 0 {
		return None[T]()
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
type dropIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	count     uint
	dropped   bool
	exhausted bool
}

func drop[T any](iter iterator[T], count uint) *dropIter[T] {
	iterator := &dropIter[T]{iter: iter, count: count}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *dropIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	if !iter.dropped {
		iter.dropped = true

		for i := uint(0); i < iter.count; i++ {
			if iter.delegateNext().IsNone() {
				return None[T]()
			}
		}
	}

	return iter.delegateNext()
}

func (iter *dropIter[T]) delegateNext() Option[T] {
	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
	}

	return next
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// enumerate
type enumerateIter[T any] struct {
	baseIter[T]
	counter   uint
	exhausted bool
}

func enumerate[T any](iter iterator[T]) *enumerateIter[T] {
	return &enumerateIter[T]{baseIter: baseIter[T]{iter}}
}

func (iter *enumerateIter[T]) Next() Option[pair[uint, T]] {
	if iter.exhausted {
		return None[pair[uint, T]]()
	}

	next := iter.baseIter.Next()
	if next.IsNone() {
		iter.exhausted = true
		return None[pair[uint, T]]()
	}

	enext := pair[uint, T]{iter.counter, next.Some()}
	iter.counter++

	return Some(enext)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// flatten
type flattenIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	innerIter *flattenIter[T]
}

func flatten[T any](iter iterator[T]) *flattenIter[T] {
	flattenIter := &flattenIter[T]{iter: iter}
	flattenIter.baseIter = baseIter[T]{flattenIter}

	return flattenIter
}

func (iter *flattenIter[T]) Next() Option[T] {
	for {
		if iter.innerIter != nil {
			if next := iter.innerIter.Next(); next.IsSome() {
				return next
			}

			iter.innerIter = nil
		}

		next := iter.iter.Next()
		if next.IsNone() {
			return None[T]()
		}

		inner, ok := any(next.Some()).(Slice[T])
		if !ok {
			return next
		}

		iter.innerIter = flatten(inner.Iter())
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// cycle
type cycleIter[T any] struct {
	baseIter[T]
	iter  iterator[T]
	items []T
	index int
}

func cycle[T any](iter iterator[T]) *cycleIter[T] {
	iterator := &cycleIter[T]{iter: iter, items: make([]T, 0)}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *cycleIter[T]) Next() Option[T] {
	if iter.iter != nil {
		if next := iter.iter.Next(); next.IsSome() {
			iter.items = append(iter.items, next.Some())
			return next
		}

		iter.iter = nil
	}

	if len(iter.items) == 0 {
		return None[T]()
	}

	if iter.index == len(iter.items) {
		iter.index = 0
	}

	next := iter.items[iter.index]
	iter.index++

	return Some(next)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// unique
type uniqueIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	seen      map[any]struct{}
	exhausted bool
}

func unique[T any](iter iterator[T]) *uniqueIter[T] {
	iterator := &uniqueIter[T]{iter: iter}
	iterator.baseIter = baseIter[T]{iterator}
	iterator.seen = make(map[any]struct{})

	return iterator
}

func (iter *uniqueIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for next := iter.iter.Next(); next.IsSome(); next = iter.iter.Next() {
		val := next.Some()
		if _, ok := iter.seen[val]; !ok {
			iter.seen[val] = struct{}{}
			return Some(val)
		}
	}

	iter.exhausted = true

	return None[T]()
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// chunks
type chunksIter[T any] struct {
	iter      iterator[T]
	size      int
	exhausted bool
}

func chunks[T any](iter iterator[T], size int) *chunksIter[T] {
	return &chunksIter[T]{
		iter: iter,
		size: size,
	}
}

func (iter *chunksIter[T]) Next() Option[Slice[T]] {
	if iter.exhausted || iter.size <= 0 {
		return None[Slice[T]]()
	}

	chunkss := make([]T, 0, iter.size)

	for i := 0; i < iter.size; i++ {
		val := iter.iter.Next()
		if val.IsNone() {
			iter.exhausted = true

			if len(chunkss) == 0 {
				return None[Slice[T]]()
			}

			break
		}

		chunkss = append(chunkss, val.Some())
	}

	return Some(Slice[T](chunkss))
}

func (iter *chunksIter[T]) Collect() []Slice[T] {
	result := make([]Slice[T], 0)

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		result = append(result, next.Some())
	}

	return result
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// zip
type zipIter[T any] struct {
	iterators []iterator[T]
}

func zip[T any](iterators ...iterator[T]) *zipIter[T] {
	return &zipIter[T]{iterators: iterators}
}

func (iter *zipIter[T]) Next() Option[Slice[T]] {
	var values []T

	for _, it := range iter.iterators {
		next := it.Next()
		if next.IsNone() {
			return None[Slice[T]]()
		}

		values = append(values, next.Some())
	}

	return Some(Slice[T](values))
}

func (iter *zipIter[T]) Collect() []Slice[T] {
	result := make([]Slice[T], 0)

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		result = append(result, next.Some())
	}

	return result
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// permutations
type permutationsIter[T any] struct {
	data    []T
	indices []int
	first   bool
}

func permutations[T any](iter iterator[T]) *permutationsIter[T] {
	var data []T
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		data = append(data, next.Some())
	}

	return &permutationsIter[T]{
		data:    data,
		indices: nil,
		first:   true,
	}
}

func (iter *permutationsIter[T]) Next() Option[Slice[T]] {
	if iter.first {
		iter.first = false
		iter.indices = make([]int, len(iter.data))

		return Some(Slice[T](iter.data))
	}

	n := len(iter.data)

	for i := n - 2; i >= 0; i-- {
		if iter.indices[i] < n-i-1 {
			iter.indices[i]++
			for j := i + 1; j < n; j++ {
				iter.indices[j] = 0
			}

			result := make([]T, n)
			copy(result, iter.data)

			for i, idx := range iter.indices {
				result[i], result[i+idx] = result[i+idx], result[i]
			}

			return Some(Slice[T](result))
		}
	}

	return None[Slice[T]]()
}

func (iter *permutationsIter[T]) Collect() []Slice[T] {
	result := make([]Slice[T], 0)

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		result = append(result, next.Some())
	}

	return result
}
