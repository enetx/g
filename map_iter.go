package g

import (
	"context"
	"math"

	"github.com/enetx/g/constraints"
	"github.com/enetx/iter"
)

// SeqMap is an iterator over sequences of pairs of values, most commonly key-value pairs.
type SeqMap[K comparable, V any] iter.Seq2[K, V]

// Pull converts the “push-style” iterator sequence seq
// into a “pull-style” iterator accessed by the two functions
// next and stop.
//
// Next returns the next pair in the sequence
// and a boolean indicating whether the pair is valid.
// When the sequence is over, next returns a pair of zero values and false.
// It is valid to call next after reaching the end of the sequence
// or after calling stop. These calls will continue
// to return a pair of zero values and false.
//
// Stop ends the iteration. It must be called when the caller is
// no longer interested in next values and next has not yet
// signaled that the sequence is over (with a false boolean return).
// It is valid to call stop multiple times and when next has
// already returned false.
//
// It is an error to call next or stop from multiple goroutines
// simultaneously.
func (seq SeqMap[K, V]) Pull() (func() (K, V, bool), func()) { return iter.Seq2[K, V](seq).Pull() }

// All checks whether all key-value pairs in the iterator satisfy the provided condition.
// This function is useful when you want to determine if all pairs in an iterator
// meet a specific criteria.
//
// Parameters:
// - fn (func(K, V) bool): A function that returns a boolean indicating whether the pair satisfies
// the condition.
//
// Returns:
// - bool: True if all pairs in the iterator satisfy the condition, false otherwise.
//
// Example usage:
//
//	m := g.Map[string, int]{"a": 1, "b": 2, "c": 3}
//	allPositive := m.Iter().All(func(_ string, v int) bool { return v > 0 })
//
// The resulting allPositive will be true if all values returned by the iterator are positive.
func (seq SeqMap[K, V]) All(fn func(K, V) bool) bool { return iter.Seq2[K, V](seq).All(fn) }

// Any checks whether any key-value pair in the iterator satisfies the provided condition.
// This function is useful when you want to determine if at least one pair in an iterator
// meets a specific criteria.
//
// Parameters:
// - fn (func(K, V) bool): A function that returns a boolean indicating whether the pair satisfies
// the condition.
//
// Returns:
// - bool: True if at least one pair in the iterator satisfies the condition, false otherwise.
//
// Example usage:
//
//	m := g.Map[string, int]{"a": 1, "b": 2, "c": 3}
//	anyEven := m.Iter().Any(func(_ string, v int) bool { return v%2 == 0 })
//
// The resulting anyEven will be true if at least one value returned by the iterator is even.
func (seq SeqMap[K, V]) Any(fn func(K, V) bool) bool { return iter.Seq2[K, V](seq).Any(fn) }

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
//
// Values of n larger than math.MaxInt are clamped to math.MaxInt so the
// conversion to int never wraps to a negative value (which would otherwise
// cause Take2 to yield nothing).
func (seq SeqMap[K, V]) Take(n uint) SeqMap[K, V] {
	if n > math.MaxInt {
		n = math.MaxInt
	}

	return SeqMap[K, V](iter.Seq2[K, V](seq).Take(int(n)))
}

// Nth returns the nth key-value pair (0-indexed) in the sequence.
func (seq SeqMap[K, V]) Nth(n Int) Option[Pair[K, V]] {
	key, value, found := iter.Seq2[K, V](seq).Nth(int(n))
	if found {
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
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
// - SeqMap[K, V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	m := g.NewMap[int, string]()
//	m.Insert(1, "a")
//	m.Insert(2, "b")
//	m.Insert(3, "c")
//	m.Insert(4, "d")
//
//	// Skipping the first two elements and collecting the rest.
//	m.Iter().Skip(2).Collect().Print()
//
// The resulting iterator will start after skipping the specified number of elements.
func (seq SeqMap[K, V]) Skip(n uint) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).Skip(int(n)))
}

// StepBy creates a new iterator that iterates over every N-th element of the original iterator.
// This function is useful when you want to skip a specific number of elements between each iteration.
//
// Parameters:
// - n uint: The step size, indicating how many elements to skip between each iteration.
//
// Returns:
// - SeqMap[K, V]: A new iterator that produces key-value pairs from the original iterator with a step size of N.
//
// Example usage:
//
//	m := g.NewMap[string, int]()
//	m.Insert("one", 1)
//	m.Insert("two", 2)
//	m.Insert("three", 3)
//
//	m.Iter().StepBy(2).Collect().Print()
//
// The resulting iterator will produce key-value pairs from the original iterator with a step size of N.
func (seq SeqMap[K, V]) StepBy(n uint) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).StepBy(int(n)))
}

// First returns the first key-value pair from the sequence.
func (seq SeqMap[K, V]) First() Option[Pair[K, V]] {
	if key, value, ok := iter.Seq2[K, V](seq).First(); ok {
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
}

// Keys returns an iterator containing all the keys in the ordered Map.
func (seq SeqMap[K, V]) Keys() SeqSlice[K] {
	return SeqSlice[K](iter.Seq2[K, V](seq).Keys())
}

// Values returns an iterator containing all the values in the ordered Map.
func (seq SeqMap[K, V]) Values() SeqSlice[V] {
	return SeqSlice[V](iter.Seq2[K, V](seq).Values())
}

// Chain creates a new iterator by concatenating the current iterator with other iterators.
//
// The function concatenates the key-value pairs from the current iterator with the key-value pairs from the provided iterators,
// producing a new iterator containing all concatenated elements.
//
// Params:
//
// - seqs ([]SeqMap[K, V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - SeqMap[K, V]: A new iterator containing elements from the current iterator and the provided iterators.
//
// Example usage:
//
//	m1 := g.NewMap[int, string]()
//	m1.Insert(1, "a")
//
//	m2 := g.NewMap[int, string]()
//	m2.Insert(2, "b")
//
//	// Concatenating iterators and collecting the result.
//	m1.Iter().Chain(m2.Iter()).Collect().Print()
//
// Output: Map{1:a, 2:b} // The output order may vary as Map is not ordered.
//
// The resulting iterator will contain elements from both iterators.
func (seq SeqMap[K, V]) Chain(seqs ...SeqMap[K, V]) SeqMap[K, V] {
	iterSeqs := make([]iter.Seq2[K, V], len(seqs))
	for i, s := range seqs {
		iterSeqs[i] = iter.Seq2[K, V](s)
	}

	return SeqMap[K, V](iter.Seq2[K, V](seq).Chain(iterSeqs...))
}

// Count consumes the iterator, counting the number of iterations and returning it.
func (seq SeqMap[K, V]) Count() Int { return Int(iter.Seq2[K, V](seq).Count()) }

// Collect collects all key-value pairs from the iterator and returns a Map.
func (seq SeqMap[K, V]) Collect() Map[K, V] {
	collection := NewMap[K, V]()

	seq(func(k K, v V) bool {
		collection[k] = v
		return true
	})

	return collection
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
// - SeqMap[K, V]: An iterator containing elements that satisfy the given function.
//
//	m := g.NewMap[int, int]()
//	m.Insert(1, 1)
//	m.Insert(2, 2)
//	m.Insert(3, 3)
//	m.Insert(4, 4)
//	m.Insert(5, 5)
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
func (seq SeqMap[K, V]) Filter(fn func(K, V) bool) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).Filter(fn))
}

// FilterByKey returns a new iterator lazily yielding only the pairs whose key
// satisfies the provided predicate; values are not inspected.
//
// It lifts a single-parameter predicate to the pair-wise Filter — composes
// with f.* factories:
//
//	m.Iter().FilterByKey(f.Eq("host"))
func (seq SeqMap[K, V]) FilterByKey(fn func(K) bool) SeqMap[K, V] {
	return seq.Filter(func(k K, _ V) bool { return fn(k) })
}

// FilterByValue returns a new iterator lazily yielding only the pairs whose
// value satisfies the provided predicate; keys are not inspected.
//
// It lifts a single-parameter predicate to the pair-wise Filter — composes
// with f.* factories:
//
//	m.Iter().FilterByValue(f.Gt(10))
func (seq SeqMap[K, V]) FilterByValue(fn func(V) bool) SeqMap[K, V] {
	return seq.Filter(func(_ K, v V) bool { return fn(v) })
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
// - SeqMap[K, V]: An iterator excluding elements that satisfy the given function.
//
// Example usage:
//
//	m := g.NewMap[int, int]()
//	m.Insert(1, 1)
//	m.Insert(2, 2)
//	m.Insert(3, 3)
//	m.Insert(4, 4)
//	m.Insert(5, 5)
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
func (seq SeqMap[K, V]) Exclude(fn func(K, V) bool) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).Exclude(fn))
}

// Find searches for an element in the iterator that satisfies the provided function.
//
// The function iterates through the elements of the iterator and returns the first element
// for which the provided function returns true.
//
// Params:
//
// - fn (func(K, V) bool): The function used to test elements for a condition.
//
// Returns:
//
// - Option[K, V]: An Option containing the first element that satisfies the condition; None if not found.
//
// Example usage:
//
//	m := g.NewMap[int, int]()
//	m.Insert(1, 1)
//	f := m.Iter().Find(func(_ int, v int) bool { return v == 1 })
//	if f.IsSome() {
//		print(f.Some().Key)
//	}
//
// The resulting Option may contain the first element that satisfies the condition, or None if not found.
func (seq SeqMap[K, V]) Find(fn func(k K, v V) bool) Option[Pair[K, V]] {
	key, value, found := iter.Seq2[K, V](seq).Find(fn)
	if found {
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
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
//	m := g.NewMap[int, int]()
//	m.Insert(1, 1)
//	m.Insert(2, 2)
//	m.Insert(3, 3)
//	m.Insert(4, 4)
//	m.Insert(5, 5)
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
func (seq SeqMap[K, V]) ForEach(fn func(k K, v V)) { iter.Seq2[K, V](seq).ForEach(fn) }

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each key-value pair as it passes through.
func (seq SeqMap[K, V]) Inspect(fn func(k K, v V)) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).Inspect(fn))
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
// - SeqMap[K, V]: A new iterator containing key-value pairs transformed by the provided function.
//
// Example usage:
//
//	m := g.NewMap[int, int]()
//	m.Insert(1, 1)
//	m.Insert(2, 2)
//	m.Insert(3, 3)
//	m.Insert(4, 4)
//	m.Insert(5, 5)
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
func (seq SeqMap[K, V]) Map[K2 comparable, V2 any](transform func(K, V) (K2, V2)) SeqMap[K2, V2] {
	return func(yield func(K2, V2) bool) {
		seq(func(k K, v V) bool { return yield(transform(k, v)) })
	}
}

// FilterMap applies a function to each key-value pair and filters out None results.
//
// The function transforms and filters pairs in a single pass. Pairs where the function
// returns None are filtered out, and pairs where it returns Some are unwrapped
// and included in the result.
//
// Params:
//
//   - fn (func(K, V) Option[Pair[K, V]]): The function that transforms and filters pairs.
//     Returns Some(Pair{key, value}) to include the transformed pair, or None to filter it out.
//
// Returns:
//
// - SeqMap[K, V]: A sequence containing only the successfully transformed pairs.
//
// Example usage:
//
//	configs := g.Map[string, string]{"host": "localhost", "port": "8080", "debug": "invalid"}
//	validConfigs := configs.Iter().FilterMap(func(k string, v string) Option[Pair[string, string]] {
//		if k == "port" || k == "host" {
//			return Some(Pair[string, string]{Key: k, Value: v + "_validated"})
//		}
//		return None[Pair[string, string]]()
//	})
//	// validConfigs will yield: {"host": "localhost_validated", "port": "8080_validated"}
//
//	users := g.Map[string, int]{"alice": 25, "bob": 17, "charlie": 30}
//	adults := users.Iter().FilterMap(func(name string, age int) Option[Pair[string, int]] {
//		if age >= 18 {
//			return Some(Pair[string, int]{Key: name, Value: age})
//		}
//		return None[Pair[string, int]]()
//	})
//	// adults will yield: {"alice": 25, "charlie": 30}
func (seq SeqMap[K, V]) FilterMap[K2 comparable, V2 any](fn func(K, V) Option[Pair[K2, V2]]) SeqMap[K2, V2] {
	return SeqMap[K2, V2](iter.Seq2[K, V](seq).FilterMap(func(k K, v V) (iter.Pair[K2, V2], bool) {
		return fn(k, v).Option()
	}))
}

// The iteration will stop when the provided function returns false for an element.
func (seq SeqMap[K, V]) Range(fn func(k K, v V) bool) { iter.Seq2[K, V](seq).Range(fn) }

// Context allows the iteration to be controlled with a context.Context.
func (seq SeqMap[K, V]) Context(ctx context.Context) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).Context(ctx))
}

// Next extracts the next key-value pair from the iterator and advances it.
//
// This method consumes the next key-value pair from the iterator and returns them wrapped in an Option.
// The iterator itself is modified to point to the remaining elements.
//
// Returns:
// - Option[Pair[K, V]]: Some(Pair{Key, Value}) if a pair exists, None if the iterator is exhausted.
func (seq *SeqMap[K, V]) Next() Option[Pair[K, V]] {
	if key, value, remaining, ok := iter.Seq2[K, V](*seq).Next(); ok {
		*seq = SeqMap[K, V](remaining)
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
}

// Last returns the final key-value pair of the sequence, or None if it is empty.
// Note: for an unordered Map the notion of "last" depends on iteration order.
func (seq SeqMap[K, V]) Last() Option[Pair[K, V]] {
	k, v, ok := iter.Seq2[K, V](seq).Last()
	return OptionOf(Pair[K, V]{Key: k, Value: v}, ok)
}

// Chan converts the sequence into a channel of Pair values, optionally bounded by a context.
func (seq SeqMap[K, V]) Chan(ctxs ...context.Context) chan Pair[K, V] {
	ctx := context.Background()
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	}

	return iter.Seq2[K, V](seq).ToChan(ctx)
}

// Fold reduces the sequence of key-value pairs to a single value using an accumulator.
// The accumulator type may differ from the key and value types.
func (seq SeqMap[K, V]) Fold[A any](init A, fn func(acc A, k K, v V) A) A {
	return iter.Seq2[K, V](seq).Fold(init, fn)
}

// SumBy maps each key-value pair to a numeric value via fn and returns the sum of those values.
// The fold order is undefined, so fn should be free of order-dependent side effects.
// An empty sequence yields the zero value of S.
func (seq SeqMap[K, V]) SumBy[S constraints.Number](fn func(K, V) S) S {
	var zero S
	return seq.Fold(zero, func(acc S, k K, v V) S { return acc + fn(k, v) })
}

// ProductBy maps each key-value pair to a numeric value via fn and returns their
// product. An empty sequence yields the multiplicative identity, one.
func (seq SeqMap[K, V]) ProductBy[S constraints.Number](fn func(K, V) S) S {
	return seq.Fold(S(1), func(acc S, k K, v V) S { return acc * fn(k, v) })
}

// FindMap applies fn to each key-value pair and returns the first Some result, or
// None if fn returns None for every pair. Iteration order is undefined.
func (seq SeqMap[K, V]) FindMap[U any](fn func(K, V) Option[U]) Option[U] {
	var result Option[U]

	seq(func(k K, v V) bool {
		if o := fn(k, v); o.IsSome() {
			result = o
			return false
		}

		return true
	})

	return result
}

// TryMap applies a fallible transform to each key-value pair and enters the
// Result pipeline, producing a SeqResult[U]. See SeqSlice.TryMap for the full
// contract (lazy, consumer-driven; terminals decide the Err policy).
func (seq SeqMap[K, V]) TryMap[U any](fn func(K, V) Result[U]) SeqResult[U] {
	return func(yield func(Result[U]) bool) {
		seq(func(k K, v V) bool {
			return yield(fn(k, v))
		})
	}
}

// TakeWhile yields key-value pairs while the predicate returns true, stopping at the first false.
func (seq SeqMap[K, V]) TakeWhile(fn func(K, V) bool) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).TakeWhile(fn))
}

// SkipWhile skips key-value pairs while the predicate returns true, then yields the rest.
func (seq SeqMap[K, V]) SkipWhile(fn func(K, V) bool) SeqMap[K, V] {
	return SeqMap[K, V](iter.Seq2[K, V](seq).SkipWhile(fn))
}
