package g

import (
	"context"
	"reflect"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/constraints"
	"github.com/enetx/iter"
)

// SeqHeap is an iterator over sequences of Heap values.
type SeqHeap[V any] iter.Seq[V]

// Pull converts the "push-style" iterator sequence seq
// into a "pull-style" iterator accessed by the two functions
// next and stop.
//
// Next returns the next value in the sequence
// and a boolean indicating whether the value is valid.
// When the sequence is over, next returns the zero V and false.
// It is valid to call next after reaching the end of the sequence
// or after calling stop. These calls will continue
// to return the zero V and false.
//
// Stop ends the iteration. It must be called when the caller is
// no longer interested in next values and next has not yet
// signaled that the sequence is over (with a false boolean return).
// It is valid to call stop multiple times and when next has
// already returned false.
//
// It is an error to call next or stop from multiple goroutines
// simultaneously.
func (seq SeqHeap[V]) Pull() (func() (V, bool), func()) { return iter.Seq[V](seq).Pull() }

// All checks whether all elements in the iterator satisfy the provided condition.
// This function is useful when you want to determine if all elements in an iterator
// meet a specific criteria.
//
// Parameters:
// - fn func(V) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
// - bool: True if all elements in the iterator satisfy the condition, false otherwise.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5, 6, 7, -1, -2)
//	isPositive := func(num int) bool { return num > 0 }
//	allPositive := heap.Iter().All(isPositive)
//
// The resulting allPositive will be true if all elements returned by the iterator are positive.
func (seq SeqHeap[V]) All(fn func(v V) bool) bool { return iter.Seq[V](seq).All(fn) }

// Any checks whether any element in the iterator satisfies the provided condition.
// This function is useful when you want to determine if at least one element in an iterator
// meets a specific criteria.
//
// Parameters:
// - fn func(V) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
// - bool: True if at least one element in the iterator satisfies the condition, false otherwise.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 3, 5, 7, 9)
//	isEven := func(num int) bool { return num%2 == 0 }
//	anyEven := heap.Iter().Any(isEven)
//
// The resulting anyEven will be true if at least one element returned by the iterator is even.
func (seq SeqHeap[V]) Any(fn func(V) bool) bool { return iter.Seq[V](seq).Any(fn) }

// Chain concatenates the current iterator with other iterators, returning a new iterator.
//
// The function creates a new iterator that combines the elements of the current iterator
// with elements from the provided iterators in the order they are given.
//
// Params:
//
// - seqs ([]SeqHeap[V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - SeqHeap[V]: A new iterator containing elements from the current iterator and the provided iterators.
//
// Example usage:
//
//	heap1 := g.NewHeap(cmp.Cmp[int])
//	heap1.Push(1, 2, 3)
//	heap2 := g.NewHeap(cmp.Cmp[int])
//	heap2.Push(4, 5, 6)
//	heap1.Iter().Chain(heap2.Iter()).Collect(cmp.Cmp[int]) // Creates new heap with all elements
//
// The resulting iterator will contain elements from both iterators in the specified order.
func (seq SeqHeap[V]) Chain(seqs ...SeqHeap[V]) SeqHeap[V] {
	iterSeqs := make([]iter.Seq[V], len(seqs))
	for i, s := range seqs {
		iterSeqs[i] = iter.Seq[V](s)
	}

	return SeqHeap[V](iter.Seq[V](seq).Chain(iterSeqs...))
}

// Chunks returns an iterator that yields chunks of elements of the specified size.
//
// The function creates a new iterator that yields chunks of elements from the original iterator,
// with each chunk containing elements of the specified size.
//
// Params:
//
// - n (Int): The size of each chunk.
//
// Returns:
//
// - SeqSlices[V]: An iterator yielding chunks of elements of the specified size.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5, 6)
//	chunks := heap.Iter().Chunks(2).Collect()
//
// Output: [Slice[1, 2] Slice[3, 4] Slice[5, 6]]
//
// The resulting iterator will yield chunks of elements, each containing the specified number of elements.
func (seq SeqHeap[V]) Chunks(n Int) SeqSlices[V] {
	return SeqSlices[V](iter.Chunks(iter.Seq[V](seq), int(n)))
}

// Collect gathers all elements from the iterator into a new Heap with a custom comparison function.
//
// Note: the comparator argument is required because SeqHeap is a purely
// functional sequence type — it carries only the yielded elements, so the
// source heap's comparator cannot travel through adapters such as Map or
// Filter and must be supplied again when a new Heap is built. Partition
// shares the same characteristic (it takes leftCmp and rightCmp).
func (seq SeqHeap[V]) Collect(compareFn func(V, V) cmp.Ordering) *Heap[V] {
	result := NewHeap(compareFn)
	seq(func(v V) bool {
		result.Push(v)
		return true
	})

	return result
}

// Count consumes the iterator, counting the number of iterations and returning it.
func (seq SeqHeap[V]) Count() Int { return Int(iter.Seq[V](seq).Count()) }

// CounterBy consumes the sequence and returns an ordered map from fn(element) to
// the number of elements that produced that key: fn is applied to every element,
// and elements whose keys collide are merged into one bucket with their counts
// summed. Key order is first-seen. The key type must be comparable; for identity
// counting pass the identity function (func(v V) V { return v }).
//
// Example usage:
//
//	words.Iter().CounterBy(func(w String) Int { return w.Len() })
//	// MapOrd{5:2, 4:1} — counts by word length, in first-seen order
func (seq SeqHeap[V]) CounterBy[K comparable](fn func(V) K) SeqMapOrd[K, Int] {
	return counterBy(iter.Seq[V](seq), fn)
}

// ChunkBy groups CONSECUTIVE elements of the sequence into chunks based on a
// custom equality function. It is not an SQL-style
// GroupBy: elements are never reordered or bucketed by key, so equal elements
// that are not adjacent end up in different chunks.
//
// The provided function `fn` takes two consecutive elements `a` and `b` and returns `true`
// if they belong to the same chunk, or `false` if a new chunk should start.
// The function returns a `SeqSlices[V]`, where each `[]V` represents a run of consecutive
// elements that satisfy the provided equality condition.
//
// Notes:
//   - Each chunk is returned as a copy of the elements, since `SeqHeap` does not guarantee
//     that elements share the same backing array.
//
// Parameters:
//   - fn (func(a, b V) bool): Function that determines whether two consecutive elements belong to the same chunk.
//
// Returns:
//   - SeqSlices[V]: An iterator yielding slices, each containing one chunk.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 1, 2, 3, 2, 3, 4)
//	chunks := heap.Iter().ChunkBy(func(a, b int) bool { return a <= b }).Collect()
//	// Output: [Slice[1, 1, 2, 3] Slice[2, 3, 4]]
//
// The resulting iterator will yield runs of consecutive elements according to the provided function.
func (seq SeqHeap[V]) ChunkBy(fn func(a, b V) bool) SeqSlices[V] {
	return SeqSlices[V](iter.GroupByAdjacent(iter.Seq[V](seq), fn))
}

// Combinations generates all combinations of length 'n' from the sequence.
func (seq SeqHeap[V]) Combinations(size Int) SeqSlices[V] {
	return SeqSlices[V](iter.Combinations(iter.Seq[V](seq), int(size)))
}

// Cycle returns an iterator that endlessly repeats the elements of the current sequence.
func (seq SeqHeap[V]) Cycle() SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Cycle())
}

// Enumerate adds an index to each element in the iterator.
//
// Returns:
//
// - SeqMapOrd[Int, V] An iterator with each element of type Pair[Int, V], where the first
// element of the pair is the index and the second element is the original element from the
// iterator.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[g.String])
//	heap.Push("bbb", "ddd", "xxx", "aaa", "ccc")
//	ps := heap.Iter().
//		Enumerate().
//		Collect()
//
//	ps.Print()
//
// Output: MapOrd{0:aaa, 1:bbb, 2:ccc, 3:ddd, 4:xxx}
func (seq SeqHeap[V]) Enumerate() SeqMapOrd[Int, V] {
	return func(yield func(Int, V) bool) {
		iterEnum := iter.Seq[V](seq).Enumerate(0)
		iterEnum(func(i int, v V) bool {
			return yield(Int(i), v)
		})
	}
}

// Dedup creates a new iterator that removes consecutive duplicate elements from the original iterator,
// leaving only one occurrence of each unique element. If the iterator is sorted, all elements will be unique.
//
// Parameters:
// - None
//
// Returns:
// - SeqHeap[V]: A new iterator with consecutive duplicates removed.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 2, 3, 4, 4, 4, 5)
//	iter := heap.Iter().Dedup()
//	result := iter.Collect(cmp.Cmp[int])
//	result.Iter().ForEach(func(v int) { fmt.Print(v, " ") })
//
// Output: 1 2 3 4 5
//
// The resulting iterator will contain only unique elements, removing consecutive duplicates.
func (seq SeqHeap[V]) Dedup() SeqHeap[V] {
	if isValueComparable[V]() {
		return SeqHeap[V](iter.Seq[V](seq).DedupBy(func(a, b V) bool {
			return any(a) == any(b)
		}))
	}

	return SeqHeap[V](iter.Seq[V](seq).DedupBy(func(a, b V) bool {
		return reflect.DeepEqual(a, b)
	}))
}

// Filter returns a new iterator containing only the elements that satisfy the provided function.
//
// The function applies the provided function to each element of the iterator.
// If the function returns true for an element, that element is included in the resulting iterator.
//
// Parameters:
//
// - fn (func(V) bool): The function to be applied to each element of the iterator
// to determine if it should be included in the result.
//
// Returns:
//
// - SeqHeap[V]: A new iterator containing the elements that satisfy the given condition.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	even := heap.Iter().
//		Filter(
//			func(val int) bool {
//				return val%2 == 0
//			}).
//		Collect(cmp.Cmp[int])
//
// The resulting iterator will contain only the elements that satisfy the provided function.
func (seq SeqHeap[V]) Filter(fn func(V) bool) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Filter(fn))
}

// Exclude returns a new iterator excluding elements that satisfy the provided function.
//
// The function applies the provided function to each element of the iterator.
// If the function returns true for an element, that element is excluded from the resulting iterator.
//
// Parameters:
//
// - fn (func(V) bool): The function to be applied to each element of the iterator
// to determine if it should be excluded from the result.
//
// Returns:
//
// - SeqHeap[V]: A new iterator containing the elements that do not satisfy the given condition.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	notEven := heap.Iter().
//		Exclude(
//			func(val int) bool {
//				return val%2 == 0
//			}).
//		Collect(cmp.Cmp[int])
//
// The resulting iterator will contain only the elements that do not satisfy the provided function.
func (seq SeqHeap[V]) Exclude(fn func(V) bool) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Exclude(fn))
}

// Fold accumulates values in the iterator using a function.
//
// The function iterates through the elements of the iterator, accumulating values
// using the provided function and an initial value.
//
// Params:
//
//   - init (V): The initial value for accumulation.
//   - fn (func(V, V) V): The function that accumulates values; it takes two arguments
//     of type V and returns a value of type V.
//
// Returns:
//
// - T: The accumulated value after applying the function to all elements.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	sum := heap.Iter().
//		Fold(0,
//			func(acc, val int) int {
//				return acc + val
//			})
//	fmt.Println(sum)
//
// Output: 15.
//
// The resulting value will be the accumulation of elements based on the provided function.
func (seq SeqHeap[V]) Fold[A any](init A, fn func(acc A, val V) A) A {
	return iter.Seq[V](seq).Fold(init, fn)
}

// SumBy maps each element to a numeric value via fn and returns the sum of those values.
// An empty sequence yields the zero value of S.
func (seq SeqHeap[V]) SumBy[S constraints.Number](fn func(V) S) S {
	var zero S
	return seq.Fold(zero, func(acc S, v V) S { return acc + fn(v) })
}

// ProductBy maps each element to a numeric value via fn and returns their product.
// An empty sequence yields the multiplicative identity, one.
func (seq SeqHeap[V]) ProductBy[S constraints.Number](fn func(V) S) S {
	return seq.Fold(S(1), func(acc S, v V) S { return acc * fn(v) })
}

// FindMap applies fn to each element and returns the first Some result, or None
// if fn returns None for every element.
func (seq SeqHeap[V]) FindMap[U any](fn func(V) Option[U]) Option[U] {
	var result Option[U]

	seq(func(v V) bool {
		if o := fn(v); o.IsSome() {
			result = o
			return false
		}

		return true
	})

	return result
}

// TryMap applies a fallible transform to each element and enters the Result
// pipeline, producing a SeqResult[U]. See SeqSlice.TryMap for the full contract.
func (seq SeqHeap[V]) TryMap[U any](fn func(V) Result[U]) SeqResult[U] {
	return func(yield func(Result[U]) bool) {
		seq(func(v V) bool {
			return yield(fn(v))
		})
	}
}

// Reduce aggregates elements of the sequence using the provided function.
// The first element of the sequence is used as the initial accumulator value.
// If the sequence is empty, it returns None[V].
//
// Params:
//   - fn (func(V, V) V): Function that combines two values into one.
//
// Returns:
//   - Option[V]: The accumulated value wrapped in Some, or None if the sequence is empty.
//
// Example:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	product := heap.Iter().Reduce(func(a, b int) int { return a * b })
//	if product.IsSome() {
//	    fmt.Println(product.Some()) // 120
//	} else {
//	    fmt.Println("empty")
//	}
func (seq SeqHeap[V]) Reduce(fn func(a, b V) V) Option[V] {
	return OptionOf(iter.Seq[V](seq).Reduce(fn))
}

// ForEach iterates through all elements and applies the given function to each.
//
// The function applies the provided function to each element of the iterator.
//
// Params:
//
// - fn (func(V)): The function to apply to each element.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	heap.Iter().ForEach(func(val int) {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	})
//
// The provided function will be applied to each element in the iterator.
func (seq SeqHeap[V]) ForEach(fn func(v V)) { iter.Seq[V](seq).ForEach(fn) }

// Flatten flattens an iterator of iterators into a single iterator.
//
// The function creates a new iterator that flattens a sequence of iterators,
// returning a single iterator containing elements from each iterator in sequence.
//
// Returns:
//
// - SeqHeap[V]: A single iterator containing elements from the sequence of iterators.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[any])
//	heap.Push(
//		1,
//		g.SliceOf(2, 3),
//		"abc",
//		g.SliceOf("def", "ghi"),
//		g.SliceOf(4.5, 6.7),
//	)
//
//	heap.Iter().Flatten().ForEach(func(v any) { fmt.Print(v, " ") })
//
// Output: 1 2 3 abc def ghi 4.5 6.7
//
// The resulting iterator will contain elements from each iterator in sequence.
func (seq SeqHeap[V]) Flatten() SeqHeap[V] {
	return func(yield func(V) bool) {
		seq(func(item V) bool {
			return flattenValue(item, yield)
		})
	}
}

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each element as it passes through.
func (seq SeqHeap[V]) Inspect(fn func(v V)) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Inspect(fn))
}

// Intersperse inserts the provided separator between elements of the iterator.
//
// The function creates a new iterator that inserts the given separator between each
// consecutive pair of elements in the original iterator.
//
// Params:
//
// - sep (V): The separator to intersperse between elements.
//
// Returns:
//
// - SeqHeap[V]: An iterator containing elements with the separator interspersed.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[string])
//	heap.Push("Hello", "World", "!")
//	heap.Iter().
//		Intersperse(" ").
//		ForEach(func(s string) { fmt.Print(s) })
//
// Output: "! Hello World".
//
// The resulting iterator will contain elements with the separator interspersed.
func (seq SeqHeap[V]) Intersperse(sep V) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Intersperse(sep))
}

// Map transforms each element in the iterator using the given function.
//
// The function creates a new iterator by applying the provided function to each element
// of the original iterator.
//
// Params:
//
// - fn (func(V) V): The function used to transform elements.
//
// Returns:
//
// - SeqHeap[V]: A iterator containing elements transformed by the provided function.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3)
//	doubled := heap.
//		Iter().
//		Map(
//			func(val int) int {
//				return val * 2
//			}).
//		Collect(cmp.Cmp[int])
//
// The resulting iterator will contain elements transformed by the provided function.
func (seq SeqHeap[V]) Map[U any](transform func(V) U) SeqHeap[U] {
	return SeqHeap[U](iter.Seq[V](seq).Map(transform))
}

// Partition divides the elements of the iterator into two separate heaps with custom comparison functions.
// The comparator arguments are required for the same reason as in Collect:
// SeqHeap does not carry the source heap's comparator through adapters.
func (seq SeqHeap[V]) Partition(fn func(v V) bool, leftCmp, rightCmp func(V, V) cmp.Ordering) (*Heap[V], *Heap[V]) {
	left := NewHeap(leftCmp)
	right := NewHeap(rightCmp)

	seq(func(v V) bool {
		if fn(v) {
			left.Push(v)
		} else {
			right.Push(v)
		}
		return true
	})

	return left, right
}

// Permutations generates iterators of all permutations of elements.
//
// The function uses a recursive approach to generate all the permutations of the elements.
// If the iterator is empty or contains a single element, it returns the iterator itself
// wrapped in a single-element iterator.
//
// Returns:
//
// - SeqSlices[V]: An iterator of iterators containing all possible permutations of the
// elements in the iterator.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3)
//	perms := heap.Iter().Permutations().Collect()
//	for _, perm := range perms {
//	    fmt.Println(perm)
//	}
//
// Output:
// Slice[1, 2, 3]
// Slice[2, 1, 3]
// Slice[3, 1, 2]
// Slice[1, 3, 2]
// Slice[2, 3, 1]
// Slice[3, 2, 1]
//
// The resulting iterator will contain iterators representing all possible permutations
// of the elements in the original iterator.
func (seq SeqHeap[V]) Permutations() SeqSlices[V] {
	return SeqSlices[V](iter.Permutations(iter.Seq[V](seq)))
}

// Range iterates through elements until the given function returns false.
//
// The function iterates through the elements of the iterator and applies the provided function
// to each element. It stops iteration when the function returns false for an element.
//
// Params:
//
// - fn (func(V) bool): The function that evaluates elements for continuation of iteration.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	heap.Iter().Range(func(val int) bool {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	    return val < 5 // Replace this with the condition for continuing iteration.
//	})
//
// The iteration will stop when the provided function returns false for an element.
func (seq SeqHeap[V]) Range(fn func(v V) bool) { iter.Seq[V](seq).Range(fn) }

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
// - SeqHeap[V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5, 6)
//	heap.Iter().Skip(3).ForEach(func(v int) { fmt.Print(v, " ") })
//
// Output: 4 5 6
//
// The resulting iterator will start after skipping the specified number of elements.
func (seq SeqHeap[V]) Skip(n uint) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Skip(int(n)))
}

// StepBy creates a new iterator that iterates over every N-th element of the original iterator.
// This function is useful when you want to skip a specific number of elements between each iteration.
//
// Parameters:
// - n uint: The step size, indicating how many elements to skip between each iteration.
//
// Returns:
// - SeqHeap[V]: A new iterator that produces elements from the original iterator with a step size of N.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//	heap.Iter().StepBy(3).ForEach(func(v int) { fmt.Print(v, " ") })
//
// Output: 1 4 7 10
//
// The resulting iterator will produce elements from the original iterator with a step size of N.
func (seq SeqHeap[V]) StepBy(n uint) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).StepBy(int(n)))
}

// SortBy applies a custom sorting function to the elements in the iterator
// and returns a new iterator containing the sorted elements.
//
// The sorting function 'fn' should take two arguments, 'a' and 'b' of type V,
// and return the ordering between them.
//
// Example:
//
//	heap := g.NewHeap(cmp.Cmp[string])
//	heap.Push("a", "c", "b")
//	heap.Iter().
//		SortBy(func(a, b string) cmp.Ordering { return cmp.Cmp(b, a) }).
//		ForEach(func(s string) { fmt.Print(s, " ") })
//
// Output: c b a
//
// The returned iterator is of type SeqHeap[V], which implements the iterator
// interface for further iteration over the sorted elements.
func (seq SeqHeap[V]) SortBy(fn func(a, b V) cmp.Ordering) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).SortBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (seq SeqHeap[V]) Take(n uint) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Take(int(n)))
}

// First returns the first element from the sequence.
func (seq SeqHeap[V]) First() Option[V] {
	return OptionOf(iter.Seq[V](seq).First())
}

// Last returns the last element from the sequence.
func (seq SeqHeap[V]) Last() Option[V] {
	return OptionOf(iter.Seq[V](seq).Last())
}

// Nth returns the nth element (0-indexed) in the sequence.
func (seq SeqHeap[V]) Nth(n Int) Option[V] {
	return OptionOf(iter.Seq[V](seq).Nth(int(n)))
}

// Chan converts the iterator into a channel, optionally with context(s).
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
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensure cancellation to avoid goroutine leaks.
//	ch := heap.Iter().Chan(ctx)
//	for val := range ch {
//	    fmt.Println(val)
//	}
//
// The resulting channel allows streaming elements from the iterator with optional context handling.
func (seq SeqHeap[V]) Chan(ctxs ...context.Context) chan V {
	ctx := context.Background()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	return iter.Seq[V](seq).ToChan(ctx)
}

// Unique returns an iterator with only unique elements.
//
// The function returns an iterator containing only the unique elements from the original iterator.
//
// Returns:
//
// - SeqHeap[V]: An iterator containing unique elements from the original iterator.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 2, 4, 5, 3)
//	heap.Iter().Unique().ForEach(func(v int) { fmt.Print(v, " ") })
//
// Output: 1 2 3 4 5
//
// The resulting iterator will contain only unique elements from the original iterator.
func (seq SeqHeap[V]) Unique() SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Unique())
}

// Zip combines elements from the current sequence and another sequence into pairs.
// The element types of the two sequences may differ. Iteration stops when either
// sequence is exhausted.
func (seq SeqHeap[V]) Zip[U any](two SeqHeap[U]) SeqPairs[V, U] {
	return SeqPairs[V, U](iter.Seq[V](seq).Zip(iter.Seq[U](two)))
}

// Find searches for an element in the iterator that satisfies the provided function.
//
// The function iterates through the elements of the iterator and returns the first element
// for which the provided function returns true.
//
// Params:
//
// - fn (func(V) bool): The function used to test elements for a condition.
//
// Returns:
//
// - Option[V]: An Option containing the first element that satisfies the condition; None if not found.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	found := heap.Iter().Find(
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
func (seq SeqHeap[V]) Find(fn func(v V) bool) Option[V] {
	return OptionOf(iter.Seq[V](seq).Find(fn))
}

// Windows returns an iterator that yields sliding windows of elements of the specified size.
//
// The function creates a new iterator that yields windows of elements from the original iterator,
// where each window is a slice containing elements of the specified size and moves one element at a time.
//
// Params:
//
// - n (int): The size of each window.
//
// Returns:
//
// - SeqSlices[V]: An iterator yielding sliding windows of elements of the specified size.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5, 6)
//	windows := heap.Iter().Windows(3).Collect()
//
// Output: [Slice[1, 2, 3] Slice[2, 3, 4] Slice[3, 4, 5] Slice[4, 5, 6]]
//
// The resulting iterator will yield sliding windows of elements, each containing the specified number of elements.
func (seq SeqHeap[V]) Windows(n Int) SeqSlices[V] {
	return SeqSlices[V](iter.Windows(iter.Seq[V](seq), int(n)))
}

// Context allows the iteration to be controlled with a context.Context.
func (seq SeqHeap[V]) Context(ctx context.Context) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).Context(ctx))
}

// MaxBy returns the maximum element in the sequence using the provided comparison function.
func (seq SeqHeap[V]) MaxBy(fn func(V, V) cmp.Ordering) Option[V] {
	return OptionOf(iter.Seq[V](seq).MaxBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// MinBy returns the minimum element in the sequence using the provided comparison function.
func (seq SeqHeap[V]) MinBy(fn func(V, V) cmp.Ordering) Option[V] {
	return OptionOf(iter.Seq[V](seq).MinBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// FlatMap applies a function to each element and flattens the results into a single sequence.
//
// The function transforms each element into a new SeqHeap and then flattens all resulting
// sequences into a single sequence.
//
// Params:
//
//   - fn (func(V) SeqHeap[V]): The function that transforms each element into a SeqHeap.
//
// Returns:
//
// - SeqHeap[V]: A flattened sequence containing all elements from the transformed sequences.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3)
//	result := heap.Iter().FlatMap(func(n int) g.SeqHeap[int] {
//		subHeap := g.NewHeap(cmp.Cmp[int])
//		subHeap.Push(n, n*10)
//		return subHeap.Iter()
//	}).Collect(cmp.Cmp[int])
//	// result contains: 1, 10, 2, 20, 3, 30 (order depends on heap implementation)
func (seq SeqHeap[V]) FlatMap[U any](fn func(V) SeqHeap[U]) SeqHeap[U] {
	mapped := iter.Seq[V](seq).Map(func(v V) iter.Seq[U] {
		return iter.Seq[U](fn(v))
	})
	return SeqHeap[U](iter.FlattenSeq(mapped))
}

// FilterMap applies a function to each element and filters out None results.
//
// The function transforms and filters elements in a single pass. Elements where the function
// returns None are filtered out, and elements where it returns Some are unwrapped
// and included in the result.
//
// Params:
//
//   - fn (func(V) Option[V]): The function that transforms and filters elements.
//     Returns Some(value) to include the transformed element, or None to filter it out.
//
// Returns:
//
// - SeqHeap[V]: A sequence containing only the successfully transformed elements.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	result := heap.Iter().FilterMap(func(n int) g.Option[int] {
//		if n%2 == 0 {
//			return g.Some(n * 10)
//		}
//		return g.None[int]()
//	}).Collect(cmp.Cmp[int])
//	// result contains only even numbers multiplied by 10
func (seq SeqHeap[V]) FilterMap[U any](fn func(V) Option[U]) SeqHeap[U] {
	return SeqHeap[U](iter.Seq[V](seq).FilterMap(func(v V) (U, bool) {
		return fn(v).Option()
	}))
}

// Scan applies a function to each element and produces a sequence of successive accumulated results.
//
// The function takes an initial value and applies the provided function to each element along
// with the accumulated value, producing a new sequence where each element is the result of
// the accumulation. The initial value is included as the first element.
//
// Params:
//
//   - init (V): The initial value for the accumulation.
//   - fn (func(acc, val V) V): The function that combines the accumulator with each element.
//
// Returns:
//
// - SeqHeap[V]: A sequence containing the initial value and all accumulated results.
//
// Example usage:
//
//	heap := g.NewHeap(cmp.Cmp[int])
//	heap.Push(1, 2, 3, 4, 5)
//	result := heap.Iter().Scan(0, func(acc, val int) int {
//		return acc + val
//	}).Collect(cmp.Cmp[int])
//	// result contains: 0, plus cumulative sums of heap elements
func (seq SeqHeap[V]) Scan[A any](init A, fn func(acc A, val V) A) SeqHeap[A] {
	return func(yield func(A) bool) {
		if !yield(init) {
			return
		}
		iter.Seq[V](seq).Scan(init, fn)(yield)
	}
}

// Next extracts the next element from the iterator and advances it.
//
// This method consumes the next element from the iterator and returns it wrapped in an Option.
// The iterator itself is modified to point to the remaining elements.
//
// Returns:
// - Option[V]: Some(value) if an element exists, None if the iterator is exhausted.
func (seq *SeqHeap[V]) Next() Option[V] {
	if value, remaining, ok := iter.Seq[V](*seq).Next(); ok {
		*seq = SeqHeap[V](remaining)
		return Some(value)
	}

	return None[V]()
}

// TakeWhile yields elements while the predicate returns true, stopping at the first false.
func (seq SeqHeap[V]) TakeWhile(fn func(V) bool) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).TakeWhile(fn))
}

// SkipWhile skips elements while the predicate returns true, then yields the rest.
func (seq SeqHeap[V]) SkipWhile(fn func(V) bool) SeqHeap[V] {
	return SeqHeap[V](iter.Seq[V](seq).SkipWhile(fn))
}
