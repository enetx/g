package g

import (
	"context"
	"iter"
	"reflect"
	"sort"
)

type seqSlice[V any] iter.Seq[V]

func (seq seqSlice[V]) pull() (func() (V, bool), func()) { return iter.Pull(iter.Seq[V](seq)) }

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
func (seq seqSlice[V]) All(fn func(v V) bool) bool {
	for v := range seq {
		if !fn(v) {
			return false
		}
	}

	return true
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
func (seq seqSlice[V]) Any(fn func(V) bool) bool {
	for v := range seq {
		if fn(v) {
			return true
		}
	}

	return false
}

// Chain concatenates the current iterator with other iterators, returning a new iterator.
//
// The function creates a new iterator that combines the elements of the current iterator
// with elements from the provided iterators in the order they are given.
//
// Params:
//
// - seqs ([]seqSlice[V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - sequence[V]: A new iterator containing elements from the current iterator and the provided iterators.
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
func (seq seqSlice[V]) Chain(seqs ...seqSlice[V]) seqSlice[V] {
	return chainSlice(append([]seqSlice[V]{seq}, seqs...)...)
}

// Collect gathers all elements from the iterator into a Slice.
func (seq seqSlice[V]) Collect() Slice[V] {
	collection := make([]V, 0)

	for v := range seq {
		collection = append(collection, v)
	}

	return collection
}

// Cycle returns an iterator that endlessly repeats the elements of the current iterator.
func (seq seqSlice[V]) Cycle() seqSlice[V] { return cycleSlice(seq) }

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
// - seqSlice[V]: A new iterator containing the elements that do not satisfy the given condition.
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
func (seq seqSlice[V]) Exclude(fn func(V) bool) seqSlice[V] { return excludeSlice(seq, fn) }

// Enumerate adds an index to each element in the iterator.
//
// Returns:
//
// - seqMapOrd[int, V] An iterator with each element of type Pair[int, V], where the first
// element of the pair is the index and the second element is the original element from the
// iterator.
//
// Example usage:
//
//	ps := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").
//		Iter().
//		Enumerate().
//		Collect()
//
//	ps.Print()
//
// Output: MapOrd{0:bbb, 1:ddd, 2:xxx, 3:aaa, 4:ccc}
func (seq seqSlice[V]) Enumerate() seqMapOrd[int, V] { return enumerate(seq) }

// Dedup creates a new iterator that removes consecutive duplicate elements from the original iterator,
// leaving only one occurrence of each unique element. If the iterator is sorted, all elements will be unique.
//
// Parameters:
// - None
//
// Returns:
// - seqSlice[V]: A new iterator with consecutive duplicates removed.
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
func (seq seqSlice[V]) Dedup() seqSlice[V] { return dedupSlice(seq) }

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
// - seqSlice[V]: A new iterator containing the elements that satisfy the given condition.
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
func (seq seqSlice[V]) Filter(fn func(V) bool) seqSlice[V] { return filterSlice(seq, fn) }

// Fold accumulates values in the iterator using a function.
//
// The function iterates through the elements of the iterator, accumulating values
// using the provided function and an initial value.
//
// Params:
//
//   - init (V): The initial value for accumulation.
//   - fn (func(V, V) V): The function that accumulates values; it takes two arguments
//     of type V and returns a value of type T.
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
func (seq seqSlice[V]) Fold(init V, fn func(acc, val V) V) V { return fold(seq, init, fn) }

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
func (seq seqSlice[V]) ForEach(fn func(v V)) {
	for v := range seq {
		fn(v)
	}
}

// Flatten flattens an iterator of iterators into a single iterator.
//
// The function creates a new iterator that flattens a sequence of iterators,
// returning a single iterator containing elements from each iterator in sequence.
//
// Returns:
//
// - seqSlice[V]: A single iterator containing elements from the sequence of iterators.
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
func (seq seqSlice[V]) Flatten() seqSlice[V] { return flatten(seq) }

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each element as it passes through.
func (seq seqSlice[V]) Inspect(fn func(v V)) seqSlice[V] { return inspectSlice(seq, fn) }

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
// - seqSlice[V]: A new iterator containing elements transformed by the provided function.
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
func (seq seqSlice[V]) Map(transform func(V) V) seqSlice[V] { return mapSlice(seq, transform) }

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
func (seq seqSlice[V]) Range(fn func(v V) bool) {
	for v := range seq {
		if !fn(v) {
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
// - seqSlice[V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	iter := g.Slice[int]{1, 2, 3, 4, 5, 6}.Iter()
//	iter.Skip(3).Collect().Print()
//
// Output: [4, 5, 6]
//
// The resulting iterator will start after skipping the specified number of elements.
func (seq seqSlice[V]) Skip(n uint) seqSlice[V] { return skipSlice(seq, n) }

// StepBy creates a new iterator that iterates over every N-th element of the original iterator.
// This function is useful when you want to skip a specific number of elements between each iteration.
//
// Parameters:
// - n uint: The step size, indicating how many elements to skip between each iteration.
//
// Returns:
// - seqSlice[V]: A new iterator that produces elements from the original iterator with a step size of N.
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
func (seq seqSlice[V]) StepBy(n uint) seqSlice[V] { return stepbySlice(seq, n) }

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
// The returned iterator is of type sequence[V], which implements the iterator
// interface for further iteration over the sorted elements.
func (seq seqSlice[V]) Sort() seqSlice[V] { return sortiSlice(seq) }

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
// The returned iterator is of type sequence[V], which implements the iterator
// interface for further iteration over the sorted elements.
func (seq seqSlice[V]) SortBy(fn func(a, b V) bool) seqSlice[V] { return sortbySlice(seq, fn) }

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (seq seqSlice[V]) Take(n uint) seqSlice[V] { return takeSlice(seq, n) }

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
// - chan V: A channel containing the elements from the iterator.
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
func (seq seqSlice[V]) ToChannel(ctxs ...context.Context) chan V {
	ch := make(chan V)

	ctx := context.Background()
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	}

	go func() {
		defer close(ch)

		for v := range seq {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- v
			}
		}
	}()

	return ch
}

// Unique returns an iterator with only unique elements.
//
// The function returns an iterator containing only the unique elements from the original iterator.
//
// Returns:
//
// - seqSlice[V]: An iterator containing unique elements from the original iterator.
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
func (seq seqSlice[V]) Unique() seqSlice[V] { return uniqueSlice(seq) }

// Zip combines elements from the current sequence and another sequence into pairs,
// creating an ordered map with identical keys and values of type V.
func (seq seqSlice[V]) Zip(two seqSlice[V]) seqMapOrd[V, V] { return zip(seq, two) }

func liftSlice[V any](slice []V) seqSlice[V] {
	return func(yield func(V) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
}

func chainSlice[V any](seqs ...seqSlice[V]) seqSlice[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			seq(func(v V) bool {
				return yield(v)
			})
		}
	}
}

func mapSlice[V, W any](seq seqSlice[V], fn func(V) W) seqSlice[W] {
	return func(yield func(W) bool) {
		seq(func(v V) bool {
			return yield(fn(v))
		})
	}
}

func filterSlice[V any](seq seqSlice[V], fn func(V) bool) seqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if fn(v) {
				return yield(v)
			}
			return true
		})
	}
}

func excludeSlice[V any](seq seqSlice[V], fn func(V) bool) seqSlice[V] {
	return filterSlice(seq, func(v V) bool { return !fn(v) })
}

func cycleSlice[V any](seq seqSlice[V]) seqSlice[V] {
	return func(yield func(V) bool) {
		var (
			saved []V
			i     int
		)

		for v := range seq {
			saved = append(saved, v)
			if !yield(v) {
				return
			}
		}

		for len(saved) > 0 {
			for ; i < len(saved); i++ {
				if !yield(saved[i]) {
					return
				}
			}
			i = 0
		}
	}
}

func stepbySlice[V any](seq seqSlice[V], n uint) seqSlice[V] {
	return func(yield func(V) bool) {
		i := uint(0)
		seq(func(v V) bool {
			i++
			if (i-1)%n == 0 {
				return yield(v)
			}
			return true
		})
	}
}

func takeSlice[V any](seq seqSlice[V], n uint) seqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if n == 0 {
				return false
			}
			n--
			return yield(v)
		})
	}
}

func uniqueSlice[V any](seq seqSlice[V]) seqSlice[V] {
	return func(yield func(V) bool) {
		seen := make(map[any]struct{})
		seq(func(v V) bool {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				return yield(v)
			}
			return true
		})
	}
}

func dedupSlice[V any](seq seqSlice[V]) seqSlice[V] {
	return func(yield func(V) bool) {
		var current V
		seq(func(v V) bool {
			if reflect.DeepEqual(current, v) {
				return true
			}
			current = v
			return yield(v)
		})
	}
}

func sortiSlice[V any](seq seqSlice[V]) seqSlice[V] {
	items := seq.Collect()
	items.Sort()

	return items.Iter()
}

func sortbySlice[V any](seq seqSlice[V], fn func(a, b V) bool) seqSlice[V] {
	items := seq.Collect()

	sort.Slice(items, func(i, j int) bool {
		return fn(items[i], items[j])
	})

	return items.Iter()
}

func skipSlice[V any](seq seqSlice[V], n uint) seqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if n > 0 {
				n--
				return true
			}
			return yield(v)
		})
	}
}

func inspectSlice[V any](seq seqSlice[V], fn func(V)) seqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			fn(v)
			return yield(v)
		})
	}
}

func enumerate[V any](seq seqSlice[V]) seqMapOrd[int, V] {
	return func(yield func(int, V) bool) {
		i := -1
		seq(func(v V) bool {
			i++
			return yield(i, v)
		})
	}
}

func zip[V, W any](one seqSlice[V], two seqSlice[W]) seqMapOrd[V, W] {
	return func(yield func(V, W) bool) {
		oneNext, oneStop := iter.Pull(iter.Seq[V](one))
		defer oneStop()

		twoNext, twoStop := iter.Pull(iter.Seq[W](two))
		defer twoStop()

		for {
			one, ok := oneNext()
			if !ok {
				return
			}

			two, ok := twoNext()
			if !ok {
				return
			}

			if !yield(one, two) {
				return
			}
		}
	}
}

func flatten[V any](seq seqSlice[V]) seqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if inner, ok := any(v).(Slice[V]); ok {
				flatten(inner.Iter())(func(v V) bool {
					return yield(v)
				})
				return true
			}
			return yield(v)
		})
	}
}

func fold[V any](seq seqSlice[V], init V, fn func(V, V) V) V {
	seq(func(v V) bool {
		init = fn(init, v)
		return true
	})
	return init
}
