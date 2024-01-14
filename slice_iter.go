package g

import (
	"context"
	"reflect"
	"sort"
)

// Sort returns a new iterator containing the elements from the current iterator
// in sorted order. The elements must be of a comparable type.
//
// Example:
//
//	g.SliceOf(9, 8, 9, 8, 0, 1, 1, 1, 2, 7, 2, 2, 2, 3, 4, 5).
//		Iter().
//		Sort().
//		Collect().
//		Print()
//
// Output: Slice[0, 1, 1, 1, 2, 2, 2, 2, 3, 4, 5, 7, 8, 8, 9, 9]
//
// The returned iterator is of type *sortIter[T], which implements the iterator
// interface for further iteration over the sorted elements.
func (iter *baseIter[T]) Sort() *sortIter[T] {
	return sorti[T](iter)
}

// SortBy applies a custom sorting function to the elements in the iterator
// and returns a new iterator containing the sorted elements.
//
// The sorting function 'fn' should take two arguments, 'a' and 'b' of type T,
// and return true if 'a' should be ordered before 'b', and false otherwise.
//
// Example:
//
//	g.SliceOf("a", "c", "b").
//		Iter().
//		SortBy(func(a, b string) bool { return a > b }).
//		Collect().
//		Print()
//
// Output: Slice[c, b, a]
//
// The returned iterator is of type *sortIter[T], which implements the iterator
// interface for further iteration over the sorted elements.
func (iter *baseIter[T]) SortBy(fn func(a, b T) bool) *sortIter[T] {
	return sortBy[T](iter, fn)
}

// Dedup creates a new iterator that removes consecutive duplicate elements from the original iterator,
// leaving only one occurrence of each unique element. If the iterator is sorted, all elements will be unique.
//
// Parameters:
// - None
//
// Returns:
// - *dedupIter[T]: A new iterator with consecutive duplicates removed.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 2, 3, 4, 4, 4, 5}
//	iter := slice.Iter().Dedup()
//	result := iter.Collect()
//	result.Print()
//
// Output: [1 2 3 4 5]
//
// The resulting iterator will contain only unique elements, removing consecutive duplicates.
func (iter *baseIter[T]) Dedup() *dedupIter[T] {
	return dedup[T](iter)
}

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each element as it passes through.
func (iter *baseIter[T]) Inspect(fn func(T)) *inspectIter[T] {
	return inspect[T](iter, fn)
}

// StepBy creates a new iterator that iterates over every N-th element of the original iterator.
// This function is useful when you want to skip a specific number of elements between each iteration.
//
// Parameters:
// - n int: The step size, indicating how many elements to skip between each iteration.
//
// Returns:
// - *stepByIter[T]: A new iterator that produces elements from the original iterator with a step size of N.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	iter := slice.Iter().StepBy(3)
//	result := iter.Collect()
//	result.Print()
//
// Output: [1 4 7 10]
//
// The resulting iterator will produce elements from the original iterator with a step size of N.
func (iter *baseIter[T]) StepBy(n int) *stepByIter[T] {
	return stepBy[T](iter, n)
}

// All checks whether all elements in the iterator satisfy the provided condition.
// This function is useful when you want to determine if all elements in an iterator
// meet a specific criteria.
//
// Parameters:
// - fn func(T) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
// - bool: True if all elements in the iterator satisfy the condition, false otherwise.
//
// Example usage:
//
//	slice := g.SliceOf(1, 2, 3, 4, 5, 6, 7, -1, -2)
//	isPositive := func(num int) bool { return num > 0 }
//	allPositive := slice.Iter().All(isPositive)
//
// The resulting allPositive will be true if all elements returned by the iterator are positive.
func (iter *baseIter[T]) All(fn func(T) bool) bool {
	for {
		next := iter.Next()
		if next.IsNone() {
			return true
		}

		if !fn(next.Some()) {
			return false
		}
	}
}

// Any checks whether any element in the iterator satisfies the provided condition.
// This function is useful when you want to determine if at least one element in an iterator
// meets a specific criteria.
//
// Parameters:
// - fn func(T) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
// - bool: True if at least one element in the iterator satisfies the condition, false otherwise.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 3, 5, 7, 9}
//	isEven := func(num int) bool { return num%2 == 0 }
//	anyEven := slice.Iter().Any(isEven)
//
// The resulting anyEven will be true if at least one element returned by the iterator is even.
func (iter *baseIter[T]) Any(fn func(T) bool) bool {
	for {
		next := iter.Next()
		if next.IsNone() {
			return false
		}

		if fn(next.Some()) {
			return true
		}
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
//	iter1 := g.Slice[int]{1, 2, 3}.Iter()
//	iter2 := g.Slice[int]{4, 5, 6}.Iter()
//	iter1.Chain(iter2).Collect().Print()
//
// Output: [1, 2, 3, 4, 5, 6]
//
// The resulting iterator will contain elements from both iterators in the specified order.
func (iter *baseIter[T]) Chain(iterators ...iterator[T]) *chainIter[T] {
	return chain[T](append([]iterator[T]{iter}, iterators...)...)
}

// Collect gathers all elements from the iterator into a Slice.
func (iter *baseIter[T]) Collect() Slice[T] {
	values := make([]T, 0)

	for {
		next := iter.Next()
		if next.IsNone() {
			return values
		}

		values = append(values, next.Some())
	}
}

// Cycle returns an iterator that endlessly repeats the elements of the current iterator.
func (iter *baseIter[T]) Cycle() *cycleIter[T] {
	return cycle[T](iter)
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
// - *skipIter[T]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	iter := g.Slice[int]{1, 2, 3, 4, 5, 6}.Iter()
//	iter.Skip(3).Collect().Print()
//
// Output: [4, 5, 6]
//
// The resulting iterator will start after skipping the specified number of elements.
func (iter *baseIter[T]) Skip(n uint) *skipIter[T] {
	return skip[T](iter, n)
}

// Enumerate adds an index to each element in the iterator.
//
// Returns:
//
// - *enumerateIter[T]: An iterator with each element of type pair[uint, T], where the first
// element of the pair is the index and the second element is the original element from the
// iterator.
//
// Example usage:
//
//	pairs := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").
//		Iter().
//		Enumerate().
//		Collect()
//
//		ps := g.MapOrd[uint, g.String](pairs)
//		ps.Print()
//
// Output: MapOrd{0:bbb, 1:ddd, 2:xxx, 3:aaa, 4:ccc}
func (iter *baseIter[T]) Enumerate() *enumerateIter[T] {
	return enumerate[T](iter)
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
// - *filterIter[T]: A new iterator containing the elements that do not satisfy the given condition.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	notEven := slice.Iter().
//		Exclude(
//			func(val int) bool {
//				return val%2 == 0
//			}).
//		Collect()
//	notEven.Print()
//
// Output: [1, 3, 5]
//
// The resulting iterator will contain only the elements that do not satisfy the provided function.
func (iter *baseIter[T]) Exclude(fn func(T) bool) *filterIter[T] {
	return exclude[T](iter, fn)
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
// - *filterIter[T]: A new iterator containing the elements that satisfy the given condition.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	even := slice.Iter().
//		Filter(
//			func(val int) bool {
//				return val%2 == 0
//			}).
//		Collect()
//	even.Print()
//
// Output: [2 4].
//
// The resulting iterator will contain only the elements that satisfy the provided function.
func (iter *baseIter[T]) Filter(fn func(T) bool) *filterIter[T] {
	return filter[T](iter, fn)
}

// Find searches for an element in the iterator that satisfies the provided function.
//
// The function iterates through the elements of the iterator and returns the first element
// for which the provided function returns true.
//
// Params:
//
// - fn (func(T) bool): The function used to test elements for a condition.
//
// Returns:
//
// - Option[T]: An Option containing the first element that satisfies the condition; None if not found.
//
// Example usage:
//
//	iter := g.Slice[int]{1, 2, 3, 4, 5}.Iter()
//
//	found := iter.Find(
//		func(i int) bool {
//			return i == 2
//		})
//
//	if found.IsSome() {
//		fmt.Println("Found:", found.Some())
//	} else {
//		fmt.Println("Not found.")
//	}
//
// The resulting Option may contain the first element that satisfies the condition, or None if not found.
func (iter *baseIter[T]) Find(fn func(T) bool) Option[T] {
	for {
		next := iter.Next()
		if next.IsNone() {
			return None[T]()
		}

		if fn(next.Some()) {
			return next
		}
	}
}

// Flatten flattens an iterator of iterators into a single iterator.
//
// The function creates a new iterator that flattens a sequence of iterators,
// returning a single iterator containing elements from each iterator in sequence.
//
// Returns:
//
// - *flattenIter[T]: A single iterator containing elements from the sequence of iterators.
//
// Example usage:
//
//	nestedSlice := g.Slice[any]{
//		1,
//		g.SliceOf[any](2, 3),
//		"abc",
//		g.SliceOf[any]("def", "ghi"),
//		g.SliceOf[any](4.5, 6.7),
//	}
//
//	nestedSlice.Iter().Flatten().Collect().Print()
//
// Output: Slice[1, 2, 3, abc, def, ghi, 4.5, 6.7]
//
// The resulting iterator will contain elements from each iterator in sequence.
func (iter *baseIter[T]) Flatten() *flattenIter[T] {
	return flatten[T](iter)
}

// Fold accumulates values in the iterator using a function.
//
// The function iterates through the elements of the iterator, accumulating values
// using the provided function and an initial value.
//
// Params:
//
//   - init (T): The initial value for accumulation.
//   - fn (func(T, T) T): The function that accumulates values; it takes two arguments
//     of type T and returns a value of type T.
//
// Returns:
//
// - T: The accumulated value after applying the function to all elements.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	sum := slice.Iter().
//		Fold(0,
//			func(acc, val int) int {
//				return acc + val
//			})
//	fmt.Println(sum)
//
// Output: 15.
//
// The resulting value will be the accumulation of elements based on the provided function.
func (iter *baseIter[T]) Fold(init T, fn func(T, T) T) T {
	for {
		next := iter.Next()
		if next.IsNone() {
			return init
		}

		init = fn(init, next.Some())
	}
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
//	iter := g.Slice[int]{1, 2, 3, 4, 5}.Iter()
//	iter.ForEach(func(val T) {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	})
//
// The provided function will be applied to each element in the iterator.
func (iter *baseIter[T]) ForEach(fn func(T)) {
	for {
		next := iter.Next()
		if next.IsNone() {
			return
		}

		fn(next.Some())
	}
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
// - *mapIter[T, T]: A new iterator containing elements transformed by the provided function.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	doubled := slice.
//		Iter().
//		Map(
//			func(val int) int {
//				return val * 2
//			}).
//		Collect()
//	doubled.Print()
//
// Output: [2 4 6].
//
// The resulting iterator will contain elements transformed by the provided function.
func (iter *baseIter[T]) Map(fn func(T) T) *mapIter[T, T] {
	return transform[T](iter, fn)
}

// Range iterates through elements until the given function returns false.
//
// The function iterates through the elements of the iterator and applies the provided function
// to each element. It stops iteration when the function returns false for an element.
//
// Params:
//
// - fn (func(T) bool): The function that evaluates elements for continuation of iteration.
//
// Example usage:
//
//	iter := g.Slice[int]{1, 2, 3, 4, 5}.Iter()
//	iter.Range(func(val int) bool {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	    return val < 5 // Replace this with the condition for continuing iteration.
//	})
//
// The iteration will stop when the provided function returns false for an element.
func (iter *baseIter[T]) Range(fn func(T) bool) {
	for {
		next := iter.Next()
		if next.IsNone() || !fn(next.Some()) {
			return
		}
	}
}

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (iter *baseIter[T]) Take(n uint) *takeIter[T] {
	return take[T](iter, n)
}

// Unique returns an iterator with only unique elements.
//
// The function returns an iterator containing only the unique elements from the original iterator.
//
// Returns:
//
// - *uniqueIter[T]: An iterator containing unique elements from the original iterator.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 2, 4, 5, 3}
//	unique := slice.Iter().Unique().Collect()
//	unique.Print()
//
// Output: [1, 2, 3, 4, 5].
//
// The resulting iterator will contain only unique elements from the original iterator.
func (iter *baseIter[T]) Unique() *uniqueIter[T] {
	return unique[T](iter)
}

// Chunks returns an iterator that yields chunks of elements of the specified size.
//
// The function creates a new iterator that yields chunks of elements from the original iterator,
// with each chunk containing elements of the specified size.
//
// Params:
//
// - size (int): The size of each chunk.
//
// Returns:
//
// - *chunksIter[T]: An iterator yielding chunks of elements of the specified size.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6}
//	chunks := slice.Iter().Chunks(2).Collect()
//
// Output: [Slice[1, 2] Slice[3, 4] Slice[5, 6]]
//
// The resulting iterator will yield chunks of elements, each containing the specified number of elements.
func (iter *baseIter[T]) Chunks(size int) *chunksIter[T] {
	return chunks[T](iter, size)
}

// Windows returns an iterator that yields sliding windows of elements of the specified size.
//
// The function creates a new iterator that yields windows of elements from the original iterator,
// where each window is a slice containing elements of the specified size and moves one element at a time.
//
// Params:
//
// - size (int): The size of each window.
//
// Returns:
//
// - *windowsIter[T]: An iterator yielding sliding windows of elements of the specified size.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6}
//	windows := slice.Iter().Windows(3).Collect()
//
// Output: [Slice[1, 2, 3] Slice[2, 3, 4] Slice[3, 4, 5] Slice[4, 5, 6]]
//
// The resulting iterator will yield sliding windows of elements, each containing the specified number of elements.
func (iter *baseIter[T]) Windows(size int) *windowsIter[T] {
	return windows[T](iter, size)
}

// Permutations generates iterators of all permutations of elements.
//
// The function uses a recursive approach to generate all the permutations of the elements.
// If the iterator is empty or contains a single element, it returns the iterator itself
// wrapped in a single-element iterator.
//
// Returns:
//
// - *permutationsIter[T]: An iterator of iterators containing all possible permutations of the
// elements in the iterator.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	perms := slice.Iter().Permutations().Collect()
//	for _, perm := range perms {
//	    fmt.Println(perm)
//	}
//	// Output:
//	// [1 2 3]
//	// [1 3 2]
//	// [2 1 3]
//	// [2 3 1]
//	// [3 1 2]
//	// [3 2 1]
//
// The resulting iterator will contain iterators representing all possible permutations
// of the elements in the original iterator.
func (iter *baseIter[T]) Permutations() *permutationsIter[T] {
	return permutations[T](iter)
}

// Zip combines the elements of the given iterators with the current iterator into a new iterator
// of Slice[T] elements.
//
// The function combines the elements of the current iterator with the elements of the given
// iterators by index. The length of the resulting iterator is determined by the shortest
// input iterator.
//
// Params:
//
// - iterators: The iterators to be zipped with the current iterator.
//
// Returns:
//
// - *zipIter[T]: A new iterator of Slice[T] elements containing the zipped elements of the input
// iterators.
//
// Example usage:
//
//	iter1 := g.Slice[int]{1, 2, 3}.Iter()
//	iter2 := g.Slice[int]{4, 5, 6}.Iter()
//	iter3 := g.Slice[int]{7, 8, 9}.Iter()
//
//	zipped := iter1.Zip(iter2, iter3).Collect()
//
//	for _, v := range zipped {
//		v.Print()
//	}
//
//	// Output:
//	Slice[1, 4, 7]
//	Slice[2, 5, 8]
//	Slice[3, 6, 9]
func (iter *baseIter[T]) Zip(iterators ...iterator[T]) *zipIter[T] {
	return zip[T](append([]iterator[T]{iter}, iterators...)...)
}

// ToChannel converts the iterator into a channel, optionally with context(s).
//
// The function converts the elements of the iterator into a channel for streaming purposes.
// Optionally, it accepts context(s) to handle cancellation or timeout scenarios.
//
// Params:
//
// - ctxs (context.Context): Optional context(s) to control the channel behavior (e.g., cancellation).
//
// Returns:
//
// - chan T: A channel containing the elements from the iterator.
//
// Example usage:
//
//	iter := g.Slice[int]{1, 2, 3}.Iter()
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensure cancellation to avoid goroutine leaks.
//	ch := iter.ToChannel(ctx)
//	for val := range ch {
//	    fmt.Println(val)
//	}
//
// The resulting channel allows streaming elements from the iterator with optional context handling.
func (iter *baseIter[T]) ToChannel(ctxs ...context.Context) chan T {
	ch := make(chan T)

	ctx := context.Background()
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	}

	go func() {
		defer close(ch)

		for {
			next := iter.Next()
			if next.IsNone() {
				return
			}

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
// stepby
type stepByIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	n         int
	counter   uint
	exhausted bool
}

func stepBy[T any](iter iterator[T], n int) *stepByIter[T] {
	iterator := &stepByIter[T]{iter: iter, n: n}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *stepByIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for {
		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true
			return None[T]()
		}

		iter.counter++
		if (iter.counter-1)%uint(iter.n) == 0 {
			return next
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inspect
type inspectIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	fn        func(T)
	exhausted bool
}

func inspect[T any](iter iterator[T], fn func(T)) *inspectIter[T] {
	iterator := &inspectIter[T]{iter: iter, fn: fn}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *inspectIter[T]) Next() Option[T] {
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
// skip
type skipIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	count     uint
	skipped   bool
	exhausted bool
}

func skip[T any](iter iterator[T], count uint) *skipIter[T] {
	iterator := &skipIter[T]{iter: iter, count: count}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *skipIter[T]) Next() Option[T] {
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

func (iter *skipIter[T]) delegateNext() Option[T] {
	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
	}

	return next
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// enumerate
type enumerateIter[T any] struct {
	iter      iterator[T]
	counter   uint
	exhausted bool
}

func enumerate[T any](iter iterator[T]) *enumerateIter[T] {
	return &enumerateIter[T]{iter: iter}
}

func (iter *enumerateIter[T]) Next() Option[Pair[uint, T]] {
	if iter.exhausted {
		return None[Pair[uint, T]]()
	}

	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
		return None[Pair[uint, T]]()
	}

	enext := Pair[uint, T]{iter.counter, next.Some()}
	iter.counter++

	return Some(enext)
}

func (iter *enumerateIter[T]) Collect() []Pair[uint, T] {
	result := []Pair[uint, T]{}

	for {
		next := iter.Next()
		if next.IsNone() {
			return result
		}

		result = append(result, next.Some())
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// flatten
type flattenIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	innerIter *flattenIter[T]
	exhausted bool
}

func flatten[T any](iter iterator[T]) *flattenIter[T] {
	flattenIter := &flattenIter[T]{iter: iter}
	flattenIter.baseIter = baseIter[T]{flattenIter}

	return flattenIter
}

func (iter *flattenIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for {
		if iter.innerIter != nil {
			if next := iter.innerIter.Next(); next.IsSome() {
				return next
			}

			iter.innerIter = nil
		}

		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true
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

	for {
		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true
			return None[T]()
		}

		val := next.Some()
		if _, ok := iter.seen[val]; !ok {
			iter.seen[val] = struct{}{}
			return next
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// dedup
type dedupIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	current   Option[T]
	exhausted bool
}

func dedup[T any](iter iterator[T]) *dedupIter[T] {
	iterator := &dedupIter[T]{iter: iter, current: None[T]()}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *dedupIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for {
		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true
			return None[T]()
		}

		if !reflect.DeepEqual(iter.current, next) {
			iter.current = next
			return next
		}
	}
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

	result := make([]T, 0, iter.size)

	for i := 0; i < iter.size; i++ {
		next := iter.iter.Next()
		if next.IsNone() {
			iter.exhausted = true

			if len(result) == 0 {
				return None[Slice[T]]()
			}

			break
		}

		result = append(result, next.Some())
	}

	return Some(Slice[T](result))
}

func (iter *chunksIter[T]) Collect() []Slice[T] {
	result := make([]Slice[T], 0)

	for {
		next := iter.Next()
		if next.IsNone() {
			return result
		}

		result = append(result, next.Some())
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// windows
type windowsIter[T any] struct {
	iter      iterator[T]
	queue     []T
	size      int
	exhausted bool
}

func windows[T any](iter iterator[T], size int) *windowsIter[T] {
	return &windowsIter[T]{
		iter:  iter,
		size:  size,
		queue: make([]T, 0, size),
	}
}

func (iter *windowsIter[T]) Next() Option[Slice[T]] {
	if len(iter.queue) < iter.size && !iter.exhausted {
		for i := 0; i < iter.size; i++ {
			next := iter.iter.Next()
			if next.IsNone() {
				iter.exhausted = true
				break
			}

			iter.queue = append(iter.queue, next.Some())
		}
	}

	if len(iter.queue) < iter.size {
		return None[Slice[T]]()
	}

	window := iter.queue[:iter.size]
	iter.queue = iter.queue[1:]

	return Some(Slice[T](window))
}

func (iter *windowsIter[T]) Collect() []Slice[T] {
	result := make([]Slice[T], 0)

	for {
		next := iter.Next()
		if next.IsNone() {
			return result
		}

		result = append(result, next.Some())
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// zip
type zipIter[T any] struct {
	iterators []iterator[T]
	exhausted bool
}

func zip[T any](iterators ...iterator[T]) *zipIter[T] {
	return &zipIter[T]{iterators: iterators}
}

func (iter *zipIter[T]) Next() Option[Slice[T]] {
	if iter.exhausted {
		return Option[Slice[T]]{}
	}

	var values []T

	for _, it := range iter.iterators {
		next := it.Next()
		if next.IsNone() {
			iter.exhausted = true

			return None[Slice[T]]()
		}

		values = append(values, next.Some())
	}

	return Some(Slice[T](values))
}

func (iter *zipIter[T]) Collect() []Slice[T] {
	result := make([]Slice[T], 0)

	for {
		next := iter.Next()
		if next.IsNone() {
			return result
		}

		result = append(result, next.Some())
	}
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

	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}

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

	for {
		next := iter.Next()
		if next.IsNone() {
			return result
		}

		result = append(result, next.Some())
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// sort
type sortIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	items     Slice[T]
	index     int
	exhausted bool
}

func sorti[T any](iter iterator[T]) *sortIter[T] {
	iterator := &sortIter[T]{iter: iter, items: NewSlice[T]()}
	iterator.baseIter = baseIter[T]{iterator}
	iterator.collect(iter)
	iterator.items.Sort()

	return iterator
}

func sortBy[T any](iter iterator[T], fn func(a, b T) bool) *sortIter[T] {
	iterator := &sortIter[T]{iter: iter, items: NewSlice[T]()}
	iterator.baseIter = baseIter[T]{iterator}
	iterator.collect(iter)

	sort.Slice(iterator.items, func(i, j int) bool {
		return fn(iterator.items[i], iterator.items[j])
	})

	return iterator
}

func (iter *sortIter[T]) collect(inner iterator[T]) {
	for {
		next := inner.Next()
		if next.IsNone() {
			return
		}

		iter.items = append(iter.items, next.Some())
	}
}

func (iter *sortIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	if iter.index >= len(iter.items) {
		iter.exhausted = true
		return None[T]()
	}

	iter.index++

	return Some(iter.items[iter.index-1])
}
