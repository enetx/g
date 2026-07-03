package g

import (
	"context"

	"github.com/enetx/g/cmp"
	"github.com/enetx/iter"
)

// SeqMapOrd is an iterator over sequences of ordered pairs of values, most commonly ordered key-value pairs.
type SeqMapOrd[K comparable, V any] iter.Seq2[K, V]

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
func (seq SeqMapOrd[K, V]) Pull() (func() (K, V, bool), func()) {
	return iter.Seq2[K, V](seq).Pull()
}

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
//	m := g.NewMapOrd[g.String, g.Int]()
//	m.Insert("a", 1)
//	m.Insert("b", 2)
//	allPositive := m.Iter().All(func(_ g.String, v g.Int) bool { return v > 0 })
//
// The resulting allPositive will be true if all values returned by the iterator are positive.
func (seq SeqMapOrd[K, V]) All(fn func(K, V) bool) bool { return iter.Seq2[K, V](seq).All(fn) }

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
//	m := g.NewMapOrd[g.String, g.Int]()
//	m.Insert("a", 1)
//	m.Insert("b", 2)
//	anyEven := m.Iter().Any(func(_ g.String, v g.Int) bool { return v%2 == 0 })
//
// The resulting anyEven will be true if at least one value returned by the iterator is even.
func (seq SeqMapOrd[K, V]) Any(fn func(K, V) bool) bool { return iter.Seq2[K, V](seq).Any(fn) }

// Keys returns an iterator containing all the keys in the ordered Map.
func (seq SeqMapOrd[K, V]) Keys() SeqSlice[K] {
	return SeqSlice[K](iter.Seq2[K, V](seq).Keys())
}

// Values returns an iterator containing all the values in the ordered Map.
func (seq SeqMapOrd[K, V]) Values() SeqSlice[V] {
	return SeqSlice[V](iter.Seq2[K, V](seq).Values())
}

// Unzip consumes the sequence and collects the keys and values
// of each pair into two separate slices.
func (seq SeqMapOrd[K, V]) Unzip() (Slice[K], Slice[V]) {
	var (
		keys   Slice[K]
		values Slice[V]
	)

	seq(func(k K, v V) bool {
		keys = append(keys, k)
		values = append(values, v)

		return true
	})

	return keys, values
}

// SortBy applies a custom sorting function to the elements in the iterator
// and returns a new iterator containing the sorted elements.
//
// The sorting function 'fn' should take two arguments, 'a' and 'b', of type Pair[K, V],
// and return true if 'a' should be ordered before 'b', and false otherwise.
//
// Example:
//
//	m := g.NewMapOrd[g.Int, g.String]()
//	m.Insert(6, "bb")
//	m.Insert(0, "dd")
//	m.Insert(1, "aa")
//	m.Insert(5, "xx")
//	m.Insert(2, "cc")
//	m.Insert(3, "ff")
//	m.Insert(4, "zz")
//
//	m.Iter().
//		SortBy(
//			func(a, b g.Pair[g.Int, g.String]) cmp.Ordering {
//				return a.Key.Cmp(b.Key)
//				// return a.Value.Cmp(b.Value)
//			}).
//		Collect().
//		Print()
//
// Output: MapOrd{0:dd, 1:aa, 2:cc, 3:ff, 4:zz, 5:xx, 6:bb}
//
// The returned iterator is of type SeqMapOrd[K, V], which implements the iterator
// interface for further iteration over the sorted elements.
func (seq SeqMapOrd[K, V]) SortBy(fn func(a, b Pair[K, V]) cmp.Ordering) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](
		iter.Seq2[K, V](seq).SortBy(func(a, b iter.Pair[K, V]) bool { return fn(a, b) == cmp.Less }),
	)
}

// SortByKey applies a custom sorting function to the keys in the iterator
// and returns a new iterator containing the sorted elements.
//
// The sorting function 'fn' should take two arguments, 'a' and 'b', of type K,
// and return true if 'a' should be ordered before 'b', and false otherwise.
//
// Example:
//
//	m := g.NewMapOrd[g.Int, g.String]()
//	m.Insert(6, "bb")
//	m.Insert(0, "dd")
//	m.Insert(1, "aa")
//	m.Insert(5, "xx")
//	m.Insert(2, "cc")
//	m.Insert(3, "ff")
//	m.Insert(4, "zz")
//
//	m.Iter().
//		SortByKey(g.Int.Cmp).
//		Collect().
//		Print()
//
// Output: MapOrd{0:dd, 1:aa, 2:cc, 3:ff, 4:zz, 5:xx, 6:bb}
func (seq SeqMapOrd[K, V]) SortByKey(fn func(a, b K) cmp.Ordering) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).SortByKey(func(a, b K) bool { return fn(a, b) == cmp.Less }))
}

// SortByValue applies a custom sorting function to the values in the iterator
// and returns a new iterator containing the sorted elements.
//
// The sorting function 'fn' should take two arguments, 'a' and 'b', of type V,
// and return true if 'a' should be ordered before 'b', and false otherwise.
//
// Example:
//
//	m := g.NewMapOrd[g.Int, g.String]()
//	m.Insert(6, "bb")
//	m.Insert(0, "dd")
//	m.Insert(1, "aa")
//	m.Insert(5, "xx")
//	m.Insert(2, "cc")
//	m.Insert(3, "ff")
//	m.Insert(4, "zz")
//
//	m.Iter().
//		SortByValue(g.String.Cmp).
//		Collect().
//		Print()
//
// Output: MapOrd{1:aa, 6:bb, 2:cc, 0:dd, 3:ff, 5:xx, 4:zz}
func (seq SeqMapOrd[K, V]) SortByValue(fn func(a, b V) cmp.Ordering) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).SortByValue(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each key-value pair as it passes through.
func (seq SeqMapOrd[K, V]) Inspect(fn func(k K, v V)) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).Inspect(fn))
}

// StepBy creates a new iterator that iterates over every N-th element of the original iterator.
// This function is useful when you want to skip a specific number of elements between each iteration.
//
// Parameters:
// - n int: The step size, indicating how many elements to skip between each iteration.
//
// Returns:
// - SeqMapOrd[K, V]: A new iterator that produces key-value pairs from the original iterator with a step size of N.
//
// Example usage:
//
//	mapIter := g.MapOrd[string, int]{{"one", 1}, {"two", 2}, {"three", 3}}.Iter()
//	iter := mapIter.StepBy(2)
//	result := iter.Collect()
//	result.Print()
//
// Output: MapOrd{one:1, three:3}
//
// The resulting iterator will produce key-value pairs from the original iterator with a step size of N.
func (seq SeqMapOrd[K, V]) StepBy(n uint) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).StepBy(int(n)))
}

// Chain concatenates the current iterator with other iterators, returning a new iterator.
//
// The function creates a new iterator that combines the elements of the current iterator
// with elements from the provided iterators in the order they are given.
//
// Params:
//
// - seqs ([]seqMapOrd[K, V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - SeqMapOrd[K, V]: A new iterator containing elements from the current iterator and the provided iterators.
//
// Example usage:
//
//	m1 := g.NewMapOrd[int, string]()
//	m1.Insert(1, "a")
//
//	m2 := g.NewMapOrd[int, string]()
//	m2.Insert(2, "b")
//
//	// Concatenating iterators and collecting the result.
//	m1.Iter().Chain(m2.Iter()).Collect().Print()
//
// Output: MapOrd{1:a, 2:b}
//
// The resulting iterator will contain elements from both iterators in the specified order.
func (seq SeqMapOrd[K, V]) Chain(seqs ...SeqMapOrd[K, V]) SeqMapOrd[K, V] {
	iterSeqs := make([]iter.Seq2[K, V], len(seqs))
	for i, s := range seqs {
		iterSeqs[i] = iter.Seq2[K, V](s)
	}

	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).Chain(iterSeqs...))
}

// Count consumes the iterator, counting the number of iterations and returning it.
func (seq SeqMapOrd[K, V]) Count() Int { return Int(iter.Seq2[K, V](seq).Count()) }

// Collect collects all key-value pairs from the iterator and returns a MapOrd.
//
// Duplicate keys keep their first-seen position, while the value is updated
// to the most recent one (last-write-wins).
func (seq SeqMapOrd[K, V]) Collect() MapOrd[K, V] {
	collection := NewMapOrd[K, V]()
	idx := make(map[K]int)

	seq(func(k K, v V) bool {
		if i, ok := idx[k]; ok {
			collection[i].Value = v
			return true
		}

		idx[k] = len(collection)
		collection = append(collection, Pair[K, V]{Key: k, Value: v})

		return true
	})

	return collection
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
// - SeqMapOrd[K, V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//

//	m := g.NewMapOrd[int, string]()
//	m.Insert(1, "a")
//	m.Insert(2, "b")
//	m.Insert(3, "c")
//	m.Insert(4, "d")
//
//	// Skipping the first two elements and collecting the rest.
//	m.Iter().Skip(2).Collect().Print()
//
// Output: MapOrd{3:c, 4:d}
//
// The resulting iterator will start after skipping the specified number of elements.
func (seq SeqMapOrd[K, V]) Skip(n uint) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).Skip(int(n)))
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
// - SeqMapOrd[K, V]: A new iterator excluding elements that satisfy the given condition.
//
// Example usage:
//
//	mo := g.NewMapOrd[int, int]()
//	mo.Insert(1, 1)
//	mo.Insert(2, 2)
//	mo.Insert(3, 3)
//	mo.Insert(4, 4)
//	mo.Insert(5, 5)
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
func (seq SeqMapOrd[K, V]) Exclude(fn func(K, V) bool) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).Exclude(fn))
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
// - SeqMapOrd[K, V]: A new iterator containing elements that satisfy the given condition.
//
// Example usage:
//
//	mo := g.NewMapOrd[int, int]()
//	mo.Insert(1, 1)
//	mo.Insert(2, 2)
//	mo.Insert(3, 3)
//	mo.Insert(4, 4)
//	mo.Insert(5, 5)
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
func (seq SeqMapOrd[K, V]) Filter(fn func(K, V) bool) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).Filter(fn))
}

// FilterByKey returns a new iterator lazily yielding only the pairs whose key
// satisfies the provided predicate; values are not inspected.
//
// It lifts a single-parameter predicate to the pair-wise Filter — composes
// with f.* factories:
//
//	mo.Iter().FilterByKey(f.Eq("host"))
func (seq SeqMapOrd[K, V]) FilterByKey(fn func(K) bool) SeqMapOrd[K, V] {
	return seq.Filter(func(k K, _ V) bool { return fn(k) })
}

// FilterByValue returns a new iterator lazily yielding only the pairs whose
// value satisfies the provided predicate; keys are not inspected.
//
// It lifts a single-parameter predicate to the pair-wise Filter — composes
// with f.* factories:
//
//	mo.Iter().FilterByValue(f.Gt(10))
func (seq SeqMapOrd[K, V]) FilterByValue(fn func(V) bool) SeqMapOrd[K, V] {
	return seq.Filter(func(_ K, v V) bool { return fn(v) })
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
//	m := g.NewMapOrd[int, int]()
//	m.Insert(1, 1)
//	f := m.Iter().Find(func(_ int, v int) bool { return v == 1 })
//	if f.IsSome() {
//		print(f.Some().Key)
//	}
//
// The resulting Option may contain the first element that satisfies the condition, or None if not found.
func (seq SeqMapOrd[K, V]) Find(fn func(k K, v V) bool) Option[Pair[K, V]] {
	key, value, found := iter.Seq2[K, V](seq).Find(fn)
	if found {
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
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
//	m := g.NewMapOrd[int, int]()
//	m.Insert(1, 1)
//	m.Insert(2, 2)
//	m.Insert(3, 3)
//	m.Insert(4, 4)
//	m.Insert(5, 5)
//
//	m.Iter().ForEach(func(key K, val V) {
//	    // Process key-value pair
//	})
//
// The provided function will be applied to each key-value pair in the iterator.
func (seq SeqMapOrd[K, V]) ForEach(fn func(k K, v V)) {
	iter.Seq2[K, V](seq).ForEach(fn)
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
// - SeqMapOrd[K, V]: A new iterator containing transformed key-value pairs.
//
// Example usage:
//
//	mo := g.NewMapOrd[int, int]()
//	mo.Insert(1, 1)
//	mo.Insert(2, 2)
//	mo.Insert(3, 3)
//	mo.Insert(4, 4)
//	mo.Insert(5, 5)
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
func (seq SeqMapOrd[K, V]) Map[K2 comparable, V2 any](transform func(K, V) (K2, V2)) SeqMapOrd[K2, V2] {
	return func(yield func(K2, V2) bool) {
		seq(func(k K, v V) bool { return yield(transform(k, v)) })
	}
}

// FilterMap applies a function to each key-value pair and filters out None results.
//
// Pairs where the function returns None are filtered out; pairs where it returns
// Some(Pair) are transformed and included in the result. Key and value types may differ
// from the input types.
func (seq SeqMapOrd[K, V]) FilterMap[K2 comparable, V2 any](fn func(K, V) Option[Pair[K2, V2]]) SeqMapOrd[K2, V2] {
	return SeqMapOrd[K2, V2](iter.Seq2[K, V](seq).FilterMap(func(k K, v V) (iter.Pair[K2, V2], bool) {
		return fn(k, v).Option()
	}))
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
//	m := g.NewMapOrd[int, int]()
//	m.Insert(1, 1)
//	m.Insert(2, 2)
//	m.Insert(3, 3)
//	m.Insert(4, 4)
//	m.Insert(5, 5)
//
//	m.Iter().Range(func(k, v int) bool {
//	    fmt.Println(v) // Replace this with the function logic you need.
//	    return v < 5 // Replace this with the condition for continuing iteration.
//	})
//
// The iteration will stop when the provided function returns false.
func (seq SeqMapOrd[K, V]) Range(fn func(k K, v V) bool) {
	iter.Seq2[K, V](seq).Range(fn)
}

// Context allows the iteration to be controlled with a context.Context.
func (seq SeqMapOrd[K, V]) Context(ctx context.Context) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).Context(ctx))
}

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (seq SeqMapOrd[K, V]) Take(n uint) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).Take(int(n)))
}

// First returns the first key-value pair from the sequence.
func (seq SeqMapOrd[K, V]) First() Option[Pair[K, V]] {
	if key, value, ok := iter.Seq2[K, V](seq).First(); ok {
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
}

// Last returns the last key-value pair from the sequence.
func (seq SeqMapOrd[K, V]) Last() Option[Pair[K, V]] {
	if key, value, ok := iter.Seq2[K, V](seq).Last(); ok {
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
}

// Nth returns the nth key-value pair (0-indexed) in the sequence.
func (seq SeqMapOrd[K, V]) Nth(n Int) Option[Pair[K, V]] {
	key, value, found := iter.Seq2[K, V](seq).Nth(int(n))
	if found {
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
}

// Chan converts the iterator into a channel, optionally with context(s).
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
// - chan Pair[K, V]: A channel emitting key-value pairs from the iterator.
//
// Example usage:
//
//	m := g.NewMapOrd[int, int]()
//	m.Insert(1, 1)
//	m.Insert(2, 2)
//	m.Insert(3, 3)
//	m.Insert(4, 4)
//	m.Insert(5, 5)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensure cancellation to avoid goroutine leaks.
//
//	ch := m.Iter().Chan(ctx)
//	for pair := range ch {
//	    // Process key-value pair from the channel
//	}
//
// The function converts the iterator into a channel to allow sequential or concurrent processing of key-value pairs.
func (seq SeqMapOrd[K, V]) Chan(ctxs ...context.Context) chan Pair[K, V] {
	ctx := context.Background()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	return iter.Seq2[K, V](seq).ToChan(ctx)
}

// Next extracts the next key-value pair from the iterator and advances it.
//
// This method consumes the next key-value pair from the iterator and returns them wrapped in an Option.
// The iterator itself is modified to point to the remaining elements.
//
// Returns:
// - Option[Pair[K, V]]: Some(Pair{Key, Value}) if a pair exists, None if the iterator is exhausted.
func (seq *SeqMapOrd[K, V]) Next() Option[Pair[K, V]] {
	if key, value, remaining, ok := iter.Seq2[K, V](*seq).Next(); ok {
		*seq = SeqMapOrd[K, V](remaining)
		return Some(Pair[K, V]{Key: key, Value: value})
	}

	return None[Pair[K, V]]()
}

// Fold reduces the sequence of key-value pairs to a single value using an accumulator.
// The accumulator type may differ from the key and value types.
func (seq SeqMapOrd[K, V]) Fold[A any](init A, fn func(acc A, k K, v V) A) A {
	return iter.Seq2[K, V](seq).Fold(init, fn)
}

// TakeWhile yields key-value pairs while the predicate returns true, stopping at the first false.
func (seq SeqMapOrd[K, V]) TakeWhile(fn func(K, V) bool) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).TakeWhile(fn))
}

// SkipWhile skips key-value pairs while the predicate returns true, then yields the rest.
func (seq SeqMapOrd[K, V]) SkipWhile(fn func(K, V) bool) SeqMapOrd[K, V] {
	return SeqMapOrd[K, V](iter.Seq2[K, V](seq).SkipWhile(fn))
}
