package g

import (
	"context"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/constraints"
	"github.com/enetx/iter"
)

// SeqSet is an iterator over sequences of unique values.
type SeqSet[V comparable] iter.Seq[V]

// Pull converts the “push-style” iterator sequence seq
// into a “pull-style” iterator accessed by the two functions
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
func (seq SeqSet[V]) Pull() (func() (V, bool), func()) { return iter.Seq[V](seq).Pull() }

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
//	set := g.SetOf(1, 2, 3, 4, 5, 6, 7)
//	isPositive := func(num int) bool { return num > 0 }
//	allPositive := set.Iter().All(isPositive)
//
// The resulting allPositive will be true if all elements returned by the iterator are positive.
func (seq SeqSet[V]) All(fn func(v V) bool) bool { return iter.Seq[V](seq).All(fn) }

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
//	set := g.SetOf(1, 3, 5, 7, 9)
//	isEven := func(num int) bool { return num%2 == 0 }
//	anyEven := set.Iter().Any(isEven)
//
// The resulting anyEven will be true if at least one element returned by the iterator is even.
func (seq SeqSet[V]) Any(fn func(V) bool) bool { return iter.Seq[V](seq).Any(fn) }

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each element as it passes through.
func (seq SeqSet[V]) Inspect(fn func(v V)) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).Inspect(fn))
}

// Collect gathers all elements from the iterator into a Set.
func (seq SeqSet[V]) Collect() Set[V] {
	collection := make(Set[V])

	seq(func(v V) bool {
		collection[v] = Unit{}
		return true
	})

	return collection
}

// Chain concatenates the current iterator with other iterators, returning a new iterator.
//
// The function creates a new iterator that combines the elements of the current iterator
// with elements from the provided iterators in the order they are given.
//
// Params:
//
// - seqs ([]SeqSet[V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - SeqSet[V]: A new iterator containing elements from the current iterator and the provided iterators.
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
func (seq SeqSet[V]) Chain(seqs ...SeqSet[V]) SeqSet[V] {
	iterSeqs := make([]iter.Seq[V], len(seqs))
	for i, s := range seqs {
		iterSeqs[i] = iter.Seq[V](s)
	}

	return SeqSet[V](iter.Seq[V](seq).Chain(iterSeqs...))
}

// Count consumes the iterator, counting the number of iterations and returning it.
func (seq SeqSet[V]) Count() Int { return Int(iter.Seq[V](seq).Count()) }

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
func (seq SeqSet[V]) CounterBy[K comparable](fn func(V) K) SeqMapOrd[K, Int] {
	return counterBy(iter.Seq[V](seq), fn)
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
// - V: The accumulated value after applying the function to all elements.
//
// Example usage:
//
//	set := g.SetOf(1, 2, 3, 4, 5)
//	sum := set.Iter().
//		Fold(0,
//			func(acc, val int) int {
//				return acc + val
//			})
//	fmt.Println(sum)
//
// Output: 15.
//
// The resulting value will be the accumulation of elements based on the provided function.
func (seq SeqSet[V]) Fold[A any](init A, fn func(acc A, val V) A) A {
	return iter.Seq[V](seq).Fold(init, fn)
}

// SumBy maps each element to a numeric value via fn and returns the sum of those values.
// The fold order is undefined, so fn should be free of order-dependent side effects.
// An empty sequence yields the zero value of S.
func (seq SeqSet[V]) SumBy[S constraints.Number](fn func(V) S) S {
	var zero S
	return seq.Fold(zero, func(acc S, v V) S { return acc + fn(v) })
}

// ProductBy maps each element to a numeric value via fn and returns their product.
// An empty sequence yields the multiplicative identity, one.
func (seq SeqSet[V]) ProductBy[S constraints.Number](fn func(V) S) S {
	return seq.Fold(S(1), func(acc S, v V) S { return acc * fn(v) })
}

// FindMap applies fn to each element and returns the first Some result, or None
// if fn returns None for every element. Iteration order over a set is undefined.
func (seq SeqSet[V]) FindMap[U any](fn func(V) Option[U]) Option[U] {
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
func (seq SeqSet[V]) TryMap[U any](fn func(V) Result[U]) SeqResult[U] {
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
//	set := g.SetOf(1, 2, 3, 4, 5)
//	product := set.Iter().Reduce(func(a, b int) int { return a * b })
//	if product.IsSome() {
//	    fmt.Println(product.Some()) // 120
//	} else {
//	    fmt.Println("empty")
//	}
func (seq SeqSet[V]) Reduce(fn func(a, b V) V) Option[V] {
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
//	iter := g.SetOf(1, 2, 3).Iter()
//	func(val V) {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	}.ForEach()
//
// The provided function will be applied to each element in the iterator.
func (seq SeqSet[V]) ForEach(fn func(v V)) { iter.Seq[V](seq).ForEach(fn) }

// Range iterates through elements until the given function returns false.
//
// The function iterates through the elements of the iterator and applies the provided function to each element.
// The iteration will stop when the provided function returns false for an element.
//
// Params:
// - fn (func(V) bool): The function that evaluates elements for continuation of iteration.
//
// Example usage:
//
// iter := g.SetOf(1, 2, 2, 3, 4, 5).Iter()
//
//	func(v int) bool {
//	    if v == 3 {
//	        return false
//	    }
//	    print(v)
//	    return true
//	}.Range()
func (seq SeqSet[V]) Range(fn func(v V) bool) { iter.Seq[V](seq).Range(fn) }

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
// - SeqSet[V]: A new iterator containing the elements that satisfy the given condition.
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
func (seq SeqSet[V]) Filter(fn func(V) bool) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).Filter(fn))
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
// - SeqSet[V]: A new iterator containing the elements that do not satisfy the given condition.
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
func (seq SeqSet[V]) Exclude(fn func(V) bool) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).Exclude(fn))
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
//     Returns Some(value) to include the transformed value, or None to filter it out.
//
// Returns:
//
// - SeqSet[V]: A sequence containing only the successfully transformed elements.
//
// Example usage:
//
//	set := g.SetOf(1, 2, 3, 4, 5)
//	result := set.Iter().FilterMap(func(n int) g.Option[int] {
//		if n%2 == 0 {
//			return g.Some(n * 10)
//		}
//		return g.None[int]()
//	}).Collect()
//	// result contains only even numbers multiplied by 10
func (seq SeqSet[V]) FilterMap[U comparable](fn func(V) Option[U]) SeqSet[U] {
	return SeqSet[U](iter.Seq[V](seq).FilterMap(func(v V) (U, bool) {
		return fn(v).Option()
	}))
}

// Map transforms each element in the iterator using the given function.
//
// The function creates a new iterator by applying the provided function to each element
// of the original iterator.
//
// Params:
//
//   - fn (func(V) U): The function used to transform elements. The result type may differ
//     from the element type but must be comparable.
//
// Returns:
//
// - SeqSet[U]: A new iterator containing elements transformed by the provided function.
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
func (seq SeqSet[V]) Map[U comparable](transform func(V) U) SeqSet[U] {
	return SeqSet[U](iter.Seq[V](seq).Map(transform))
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
//	iter := g.SetOf(1, 2, 3, 4, 5).Iter()
//
//	found := //		func(i int) bool {
//			return i == 2
//		}.Find()
//
//	if found.IsSome() {
//		fmt.Println("Found:", found.Some())
//	} else {
//		fmt.Println("Not found.")
//	}
//
// The resulting Option may contain the first element that satisfies the condition, or None if not found.
func (seq SeqSet[V]) Find(fn func(v V) bool) Option[V] {
	return OptionOf(iter.Seq[V](seq).Find(fn))
}

// Context allows the iteration to be controlled with a context.Context.
func (seq SeqSet[V]) Context(ctx context.Context) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).Context(ctx))
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
//	set := g.SetOf(1, 2, 3)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensure cancellation to avoid goroutine leaks.
//	ch := set.Iter().Chan(ctx)
//	for val := range ch {
//	    fmt.Println(val)
//	}
//
// The resulting channel allows streaming elements from the iterator with optional context handling.
func (seq SeqSet[V]) Chan(ctxs ...context.Context) chan V {
	ctx := context.Background()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	return iter.Seq[V](seq).ToChan(ctx)
}

// MaxBy returns the maximum element in the sequence using the provided comparison function.
func (seq SeqSet[V]) MaxBy(fn func(V, V) cmp.Ordering) Option[V] {
	return OptionOf(iter.Seq[V](seq).MaxBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// MinBy returns the minimum element in the sequence using the provided comparison function.
func (seq SeqSet[V]) MinBy(fn func(V, V) cmp.Ordering) Option[V] {
	return OptionOf(iter.Seq[V](seq).MinBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
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
// - SeqSet[V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	set := g.SetOf(1, 2, 3, 4, 5, 6)
//	set.Iter().Skip(3).Collect().Print()
//
// The resulting iterator will start after skipping the specified number of elements.
func (seq SeqSet[V]) Skip(n uint) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).Skip(int(n)))
}

// StepBy creates a new iterator that iterates over every N-th element of the original iterator.
// This function is useful when you want to skip a specific number of elements between each iteration.
// Note: for an unordered Set, which elements are selected depends on iteration order.
//
// Parameters:
// - n uint: The step size, indicating how many elements to skip between each iteration.
//
// Returns:
// - SeqSet[V]: A new iterator that produces elements from the original iterator with a step size of N.
//
// Example usage:
//
//	set := g.SetOf(1, 2, 3, 4, 5, 6)
//	set.Iter().StepBy(2).Collect().Print()
//
// The resulting iterator will produce elements from the original iterator with a step size of N.
func (seq SeqSet[V]) StepBy(n uint) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).StepBy(int(n)))
}

// First returns the first element from the sequence.
// Note: for an unordered Set the notion of "first" depends on iteration order.
func (seq SeqSet[V]) First() Option[V] {
	return OptionOf(iter.Seq[V](seq).First())
}

// Last returns the final element of the sequence, or None if it is empty.
// Note: for an unordered Set the notion of "last" depends on iteration order.
func (seq SeqSet[V]) Last() Option[V] {
	return OptionOf(iter.Seq[V](seq).Last())
}

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (seq SeqSet[V]) Take(n uint) SeqSet[V] { return SeqSet[V](iter.Seq[V](seq).Take(int(n))) }

// Nth returns the nth element (0-indexed) in the sequence.
func (seq SeqSet[V]) Nth(n Int) Option[V] {
	return OptionOf(iter.Seq[V](seq).Nth(int(n)))
}

func difference[V comparable](seq SeqSet[V], other Set[V]) SeqSet[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if !other.Contains(v) {
				return yield(v)
			}
			return true
		})
	}
}

// Next extracts the next element from the iterator and advances it.
//
// This method consumes the next element from the iterator and returns it wrapped in an Option.
// The iterator itself is modified to point to the remaining elements.
//
// Returns:
// - Option[V]: Some(value) if an element exists, None if the iterator is exhausted.
func (seq *SeqSet[V]) Next() Option[V] {
	if value, remaining, ok := iter.Seq[V](*seq).Next(); ok {
		*seq = SeqSet[V](remaining)
		return Some(value)
	}

	return None[V]()
}

func intersection[V comparable](seq SeqSet[V], other Set[V]) SeqSet[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if other.Contains(v) {
				return yield(v)
			}
			return true
		})
	}
}

// Difference returns a sequence containing the elements of the sequence
// that are not present in the other set.
func (seq SeqSet[V]) Difference(other Set[V]) SeqSet[V] { return difference(seq, other) }

// Intersection returns a sequence containing only the elements of the sequence
// that are also present in the other set.
func (seq SeqSet[V]) Intersection(other Set[V]) SeqSet[V] { return intersection(seq, other) }

// Partition consumes the sequence and splits its elements into two sets based on
// the predicate: elements for which fn returns true go into the first set, the
// rest into the second.
func (seq SeqSet[V]) Partition(fn func(V) bool) (Set[V], Set[V]) {
	left, right := NewSet[V](), NewSet[V]()

	seq(func(v V) bool {
		if fn(v) {
			left.Insert(v)
		} else {
			right.Insert(v)
		}

		return true
	})

	return left, right
}

// TakeWhile yields elements while the predicate returns true, stopping at the first false.
func (seq SeqSet[V]) TakeWhile(fn func(V) bool) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).TakeWhile(fn))
}

// SkipWhile skips elements while the predicate returns true, then yields the rest.
func (seq SeqSet[V]) SkipWhile(fn func(V) bool) SeqSet[V] {
	return SeqSet[V](iter.Seq[V](seq).SkipWhile(fn))
}
