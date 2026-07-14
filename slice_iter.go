package g

import (
	"context"
	"reflect"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/constraints"
	"github.com/enetx/iter"
)

// SeqSlice is an iterator over sequences of individual values.
type SeqSlice[V any] iter.Seq[V]

// SeqSlices is an iterator over slices of sequences of individual values.
type SeqSlices[V any] iter.Seq[[]V]

// Range returns a SeqSlice[T] yielding a sequence of integers of type T,
// starting at start, incrementing by step, and ending before stop (exclusive).
//
//   - If step is omitted, it defaults to 1.
//   - If step is 0, the sequence is empty.
//   - If step does not move toward stop (e.g., positive step with start > stop),
//     the sequence is empty.
//
// Examples:
//   - Range(0, 5) yields [0, 1, 2, 3, 4]
//   - Range(5, 0, -1) yields [5, 4, 3, 2, 1]
func Range[T constraints.Integer](start, stop T, step ...T) SeqSlice[T] {
	return SeqSlice[T](iter.Iota(start, stop, step...))
}

// RangeInclusive returns a SeqSlice[T] yielding a sequence of integers of type T,
// starting at start, incrementing by step, and ending at stop (inclusive).
//
//   - If step is omitted, it defaults to 1.
//   - If step is 0, the sequence is empty.
//   - If step does not move toward stop (e.g., positive step with start > stop),
//     the sequence is empty.
//
// Examples:
//   - RangeInclusive(0, 5) yields [0, 1, 2, 3, 4, 5]
//   - RangeInclusive(5, 0, -1) yields [5, 4, 3, 2, 1, 0]
func RangeInclusive[T constraints.Integer](start, stop T, step ...T) SeqSlice[T] {
	return SeqSlice[T](iter.IotaInclusive(start, stop, step...))
}

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
func (seq SeqSlice[V]) Pull() (func() (V, bool), func()) { return iter.Seq[V](seq).Pull() }

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
//	slice := g.SliceOf(1, 2, 3, 4, 5, 6, 7, -1, -2)
//	isPositive := func(num int) bool { return num > 0 }
//	allPositive := slice.Iter().All(isPositive)
//
// The resulting allPositive will be true if all elements returned by the iterator are positive.
func (seq SeqSlice[V]) All(fn func(v V) bool) bool { return iter.Seq[V](seq).All(fn) }

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
//	slice := g.Slice[int]{1, 3, 5, 7, 9}
//	isEven := func(num int) bool { return num%2 == 0 }
//	anyEven := slice.Iter().Any(isEven)
//
// The resulting anyEven will be true if at least one element returned by the iterator is even.
func (seq SeqSlice[V]) Any(fn func(V) bool) bool { return iter.Seq[V](seq).Any(fn) }

// Chain concatenates the current iterator with other iterators, returning a new iterator.
//
// The function creates a new iterator that combines the elements of the current iterator
// with elements from the provided iterators in the order they are given.
//
// Params:
//
// - seqs ([]SeqSlice[V]): Other iterators to be concatenated with the current iterator.
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
func (seq SeqSlice[V]) Chain(seqs ...SeqSlice[V]) SeqSlice[V] {
	iterSeqs := make([]iter.Seq[V], len(seqs))
	for i, s := range seqs {
		iterSeqs[i] = iter.Seq[V](s)
	}

	return SeqSlice[V](iter.Seq[V](seq).Chain(iterSeqs...))
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
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6}
//	chunks := slice.Iter().Chunks(2).Collect()
//
// Output: [Slice[1, 2] Slice[3, 4] Slice[5, 6]]
//
// The resulting iterator will yield chunks of elements, each containing the specified number of elements.
func (seq SeqSlice[V]) Chunks(n Int) SeqSlices[V] {
	return SeqSlices[V](iter.Chunks(iter.Seq[V](seq), int(n)))
}

// Collect gathers all elements from the iterator into a Slice.
func (seq SeqSlice[V]) Collect() Slice[V] { return iter.Seq[V](seq).ToSlice() }

// Collect gathers all elements from the iterator into a []Slice.
func (seqs SeqSlices[V]) Collect() []Slice[V] {
	collection := make([]Slice[V], 0)

	seqs(func(v []V) bool {
		chunk := make(Slice[V], len(v))
		copy(chunk, v)
		collection = append(collection, chunk)
		return true
	})

	return collection
}

// Map transforms each group (sub-slice) in the iterator using the given function.
//
// The function creates a new lazy iterator by applying the provided function to each
// group produced by the original iterator, preserving the streaming pipeline.
//
// Params:
//
// - fn (func(Slice[V]) Slice[V]): The function used to transform each group.
//
// Returns:
//
// - SeqSlices[V]: An iterator yielding the transformed groups.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4}
//	doubled := slice.Iter().
//		Chunks(2).
//		Map(func(chunk g.Slice[int]) g.Slice[int] {
//			return chunk.Iter().Map(func(v int) int { return v * 2 }).Collect()
//		}).
//		Collect()
//	// Output: [Slice[2, 4] Slice[6, 8]]
func (seqs SeqSlices[V]) Map[U any](fn func(Slice[V]) Slice[U]) SeqSlices[U] {
	return func(yield func([]U) bool) {
		seqs(func(v []V) bool {
			return yield(fn(Slice[V](v)))
		})
	}
}

// Filter returns a new iterator containing only the groups (sub-slices) that satisfy
// the provided function.
//
// The function applies the provided function to each group produced by the original
// iterator. If the function returns true for a group, that group is included in the
// resulting iterator.
//
// Params:
//
// - fn (func(Slice[V]) bool): The predicate applied to each group.
//
// Returns:
//
// - SeqSlices[V]: An iterator yielding the groups that satisfy the given condition.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6}
//	pairs := slice.Iter().
//		Chunks(2).
//		Filter(func(chunk g.Slice[int]) bool { return chunk.Len() == 2 }).
//		Collect()
func (seqs SeqSlices[V]) Filter(fn func(Slice[V]) bool) SeqSlices[V] {
	return func(yield func([]V) bool) {
		seqs(func(v []V) bool {
			if fn(Slice[V](v)) {
				return yield(v)
			}
			return true
		})
	}
}

// ForEach iterates through all groups (sub-slices) and applies the given function to each.
//
// Params:
//
// - fn (func(Slice[V])): The function to apply to each group.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4}
//	slice.Iter().Chunks(2).ForEach(func(chunk g.Slice[int]) {
//		fmt.Println(chunk)
//	})
func (seqs SeqSlices[V]) ForEach(fn func(s Slice[V])) {
	seqs(func(v []V) bool {
		fn(Slice[V](v))
		return true
	})
}

// Flatten flattens the iterator of groups (sub-slices) into a single SeqSlice[V],
// yielding the elements of each group in order.
//
// Returns:
//
// - SeqSlice[V]: A single iterator containing the elements from each group in sequence.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6}
//	flat := slice.Iter().Chunks(2).Flatten().Collect()
//	// Output: Slice[1, 2, 3, 4, 5, 6]
func (seqs SeqSlices[V]) Flatten() SeqSlice[V] {
	return func(yield func(V) bool) {
		seqs(func(v []V) bool {
			for _, item := range v {
				if !yield(item) {
					return false
				}
			}
			return true
		})
	}
}

// Count consumes the iterator, counting the number of iterations and returning it.
func (seq SeqSlice[V]) Count() Int { return Int(iter.Seq[V](seq).Count()) }

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
func (seq SeqSlice[V]) CounterBy[K comparable](fn func(V) K) SeqMapOrd[K, Int] {
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
//   - Each chunk is returned as a copy of the elements, since `SeqSlice` does not guarantee
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
//	slice := g.SliceOf(1, 1, 2, 3, 2, 3, 4)
//	chunks := slice.Iter().ChunkBy(func(a, b int) bool { return a <= b }).Collect()
//	// Output: [Slice[1, 1, 2, 3] Slice[2, 3, 4]]
//
// The resulting iterator will yield runs of consecutive elements according to the provided function.
func (seq SeqSlice[V]) ChunkBy(fn func(a, b V) bool) SeqSlices[V] {
	return SeqSlices[V](iter.GroupByAdjacent(iter.Seq[V](seq), fn))
}

// Combinations generates all combinations of length 'n' from the sequence.
func (seq SeqSlice[V]) Combinations(size Int) SeqSlices[V] {
	return SeqSlices[V](iter.Combinations(iter.Seq[V](seq), int(size)))
}

// Cycle returns an iterator that endlessly repeats the elements of the current sequence.
func (seq SeqSlice[V]) Cycle() SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Cycle())
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
//	ps := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").
//		Iter().
//		Enumerate().
//		Collect()
//
//	ps.Print()
//
// Output: MapOrd{0:bbb, 1:ddd, 2:xxx, 3:aaa, 4:ccc}
func (seq SeqSlice[V]) Enumerate() SeqMapOrd[Int, V] {
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
// - SeqSlice[V]: A new iterator with consecutive duplicates removed.
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
func (seq SeqSlice[V]) Dedup() SeqSlice[V] {
	if isValueComparable[V]() {
		return SeqSlice[V](iter.Seq[V](seq).DedupBy(func(a, b V) bool {
			return any(a) == any(b)
		}))
	}

	return SeqSlice[V](iter.Seq[V](seq).DedupBy(func(a, b V) bool {
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
// - SeqSlice[V]: A new iterator containing the elements that satisfy the given condition.
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
func (seq SeqSlice[V]) Filter(fn func(V) bool) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Filter(fn))
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
// - SeqSlice[V]: A new iterator containing the elements that do not satisfy the given condition.
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
func (seq SeqSlice[V]) Exclude(fn func(V) bool) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Exclude(fn))
}

// Fold accumulates values in the iterator using a function.
//
// The function iterates through the elements of the iterator, accumulating values
// using the provided function and an initial value.
//
// Params:
//
//   - init (A): The initial value for accumulation. The accumulator type may differ
//     from the element type.
//   - fn (func(A, V) A): The function that accumulates values; it takes the accumulator
//     and an element and returns the new accumulator.
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
func (seq SeqSlice[V]) Fold[A any](init A, fn func(acc A, val V) A) A {
	return iter.Seq[V](seq).Fold(init, fn)
}

// SumBy maps each element to a numeric value via fn and returns the sum of those values.
// An empty sequence yields the zero value of S. The result type S is chosen by fn,
// independent of the element type V.
//
// Params:
//   - fn (func(V) S): Projects an element to the numeric value to be summed.
//
// Returns:
//   - S: The sum of the projected values.
//
// Example usage:
//
//	words := g.SliceOf[g.String]("a", "bb", "ccc")
//	total := words.Iter().SumBy(func(s g.String) g.Int { return s.Len() })
//	fmt.Println(total) // 6
func (seq SeqSlice[V]) SumBy[S constraints.Number](fn func(V) S) S {
	var zero S
	return seq.Fold(zero, func(acc S, v V) S { return acc + fn(v) })
}

// ProductBy maps each element to a numeric value via fn and returns their product.
// An empty sequence yields the multiplicative identity, one.
func (seq SeqSlice[V]) ProductBy[S constraints.Number](fn func(V) S) S {
	return seq.Fold(S(1), func(acc S, v V) S { return acc * fn(v) })
}

// FindMap applies fn to each element and returns the first Some result, or None
// if fn returns None for every element.
func (seq SeqSlice[V]) FindMap[U any](fn func(V) Option[U]) Option[U] {
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
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	product := slice.Iter().Reduce(func(a, b int) int { return a * b })
//	if product.IsSome() {
//	    fmt.Println(product.Some()) // 120
//	} else {
//	    fmt.Println("empty")
//	}
func (seq SeqSlice[V]) Reduce(fn func(a, b V) V) Option[V] {
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
//	iter := g.Slice[int]{1, 2, 3, 4, 5}.Iter()
//	func(val V) {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	}.ForEach()
//
// The provided function will be applied to each element in the iterator.
func (seq SeqSlice[V]) ForEach(fn func(v V)) { iter.Seq[V](seq).ForEach(fn) }

// Flatten flattens an iterator of iterators into a single iterator.
//
// The function creates a new iterator that flattens a sequence of iterators,
// returning a single iterator containing elements from each iterator in sequence.
//
// Returns:
//
// - SeqSlice[V]: A single iterator containing elements from the sequence of iterators.
//
// Example usage:
//
//	nestedSlice := g.Slice[any]{
//		1,
//		g.SliceOf(2, 3),
//		"abc",
//		g.SliceOf("def", "ghi"),
//		g.SliceOf(4.5, 6.7),
//	}
//
//	nestedSlice.Iter().Flatten().Collect().Print()
//
// Output: Slice[1, 2, 3, abc, def, ghi, 4.5, 6.7]
//
// The resulting iterator will contain elements from each iterator in sequence.
func (seq SeqSlice[V]) Flatten() SeqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(item V) bool {
			return flattenValue(item, yield)
		})
	}
}

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each element as it passes through.
func (seq SeqSlice[V]) Inspect(fn func(v V)) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Inspect(fn))
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
// - SeqSlice[V]: An iterator containing elements with the separator interspersed.
//
// Example usage:
//
//	g.Slice[string]{"Hello", "World", "!"}.
//		Iter().
//		Intersperse(" ").
//		Collect().
//		Join().
//		Print()
//
// Output: "Hello World !".
//
// The resulting iterator will contain elements with the separator interspersed.
func (seq SeqSlice[V]) Intersperse(sep V) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Intersperse(sep))
}

// Map transforms each element in the iterator using the given function.
//
// The function creates a new iterator by applying the provided function to each element
// of the original iterator.
//
// Params:
//
//   - fn (func(V) U): The function used to transform elements. The result type may
//     differ from the element type.
//
// Returns:
//
// - SeqSlice[U]: A iterator containing elements transformed by the provided function.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	strs := slice.
//		Iter().
//		Map(
//			func(val int) g.String {
//				return g.Int(val * 2).String()
//			}).
//		Collect()
//	strs.Print()
//
// Output: Slice[2, 4, 6].
//
// The resulting iterator will contain elements transformed by the provided function.
func (seq SeqSlice[V]) Map[U any](transform func(V) U) SeqSlice[U] {
	return SeqSlice[U](iter.Seq[V](seq).Map(transform))
}

// FlatMap applies a function to each element that returns an iterator, then flattens the results.
//
// The function transforms each element into a sequence and then concatenates all sequences
// into a single flat sequence.
//
// Params:
//
// - fn (func(V) SeqSlice[V]): The function that transforms each element into a sequence.
//
// Returns:
//
// - SeqSlice[V]: A flattened sequence containing all elements from the transformed sequences.
//
// Example usage:
//
//	words := g.Slice[string]{"hello world", "foo bar"}.Iter()
//	chars := words.FlatMap(func(s string) SeqSlice[string] {
//		return g.String(s).Split("")
//	})
//	// chars will yield: "h", "e", "l", "l", "o", " ", "w", "o", "r", "l", "d", "f", "o", "o", " ", "b", "a", "r"
//
//	numbers := g.Slice[int]{1, 2, 3}.Iter()
//	expanded := numbers.FlatMap(func(n int) SeqSlice[int] {
//		return g.Slice[int]{n, n*10, n*100}.Iter()
//	})
//	// expanded will yield: 1, 10, 100, 2, 20, 200, 3, 30, 300
func (seq SeqSlice[V]) FlatMap[U any](fn func(V) SeqSlice[U]) SeqSlice[U] {
	mapped := iter.Seq[V](seq).Map(func(v V) iter.Seq[U] {
		return iter.Seq[U](fn(v))
	})

	return SeqSlice[U](iter.FlattenSeq(mapped))
}

// FilterMap applies a function to each element and filters out None results.
//
// The function transforms and filters in a single pass. Elements where the function
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
// - SeqSlice[V]: A sequence containing only the successfully transformed elements.
//
// Example usage:
//
//	strings := g.Slice[string]{"1", "2", "abc", "3", "xyz"}.Iter()
//	numbers := strings.FilterMap(func(s string) Option[int] {
//		if n, err := strconv.Atoi(s); err == nil {
//			return Some(n)
//		}
//		return None[int]()
//	})
//	// numbers will yield: 1, 2, 3
//
//	values := g.Slice[int]{1, -2, 3, -4, 5}.Iter()
//	positiveDoubled := values.FilterMap(func(n int) Option[int] {
//		if n > 0 {
//			return Some(n * 2)
//		}
//		return None[int]()
//	})
//	// positiveDoubled will yield: 2, 6, 10
func (seq SeqSlice[V]) FilterMap[U any](fn func(V) Option[U]) SeqSlice[U] {
	return SeqSlice[U](iter.Seq[V](seq).FilterMap(func(v V) (U, bool) {
		return fn(v).Option()
	}))
}

// TryMap applies a fallible transform to each element and enters the Result
// pipeline, producing a SeqResult[U]. It is the bridge from a plain sequence
// into SeqResult: map each element to a Result[U] and continue with the
// SeqResult terminals (TryCollect, SumBy, ...), which choose the Err policy.
//
// TryMap itself is lazy and consumer-driven: it yields fn(v) for each element
// and leaves the Err policy to the terminal — TryCollect / Fold / Reduce /
// SumBy / All / Any / First short-circuit on the first Err, while Collect and
// Count traverse every element.
//
// Example usage:
//
//	// "abc" fails to parse -> the whole batch short-circuits
//	res := g.SliceOf[g.String]("1", "2", "3").
//		Iter().
//		TryMap(g.String.TryInt).
//		TryCollect() // Ok(Slice[1, 2, 3])
//
//	sum := g.SliceOf[g.String]("1", "2", "3").
//		Iter().
//		TryMap(g.String.TryInt).
//		SumBy(f.Id) // Ok(6)
func (seq SeqSlice[V]) TryMap[U any](fn func(V) Result[U]) SeqResult[U] {
	return func(yield func(Result[U]) bool) {
		seq(func(v V) bool {
			return yield(fn(v))
		})
	}
}

// Partition divides the elements of the iterator into two separate slices based on a given predicate function.
//
// The function takes a predicate function 'fn', which should return true or false for each element in the iterator.
// Elements for which 'fn' returns true are collected into the left slice, while those for which 'fn' returns false
// are collected into the right slice.
//
// Params:
//
// - fn (func(V) bool): The predicate function used to determine the placement of elements.
//
// Returns:
//
// - (Slice[V], Slice[V]): Two slices representing elements that satisfy and don't satisfy the predicate, respectively.
//
// Example usage:
//
//	evens, odds := g.Slice[int]{1, 2, 3, 4, 5}.
//		Iter().
//		Partition(
//			func(v int) bool {
//				return v%2 == 0
//			})
//
//	fmt.Println("Even numbers:", evens) // Output: Even numbers: Slice[2, 4]
//	fmt.Println("Odd numbers:", odds)   // Output: Odd numbers: Slice[1, 3, 5]
//
// The resulting two slices will contain elements separated based on whether they satisfy the predicate or not.
func (seq SeqSlice[V]) Partition(fn func(v V) bool) (Slice[V], Slice[V]) {
	return iter.Seq[V](seq).Partition(fn)
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
//	slice := g.Slice[int]{1, 2, 3}
//	perms := slice.Iter().Permutations().Collect()
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
func (seq SeqSlice[V]) Permutations() SeqSlices[V] {
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
//	iter := g.Slice[int]{1, 2, 3, 4, 5}.Iter()
//	func(val int) bool {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	    return val < 5 // Replace this with the condition for continuing iteration.
//	}.Range()
//
// The iteration will stop when the provided function returns false for an element.
func (seq SeqSlice[V]) Range(fn func(v V) bool) { iter.Seq[V](seq).Range(fn) }

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
// - SeqSlice[V]: An iterator that starts after skipping the first n elements.
//
// Example usage:
//
//	iter := g.Slice[int]{1, 2, 3, 4, 5, 6}.Iter()
//	3.Skip().Collect().Print()
//
// Output: [4, 5, 6]
//
// The resulting iterator will start after skipping the specified number of elements.
func (seq SeqSlice[V]) Skip(n uint) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Skip(int(n)))
}

// StepBy creates a new iterator that iterates over every N-th element of the original iterator.
// This function is useful when you want to skip a specific number of elements between each iteration.
//
// Parameters:
// - n uint: The step size, indicating how many elements to skip between each iteration.
//
// Returns:
// - SeqSlice[V]: A new iterator that produces elements from the original iterator with a step size of N.
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
func (seq SeqSlice[V]) StepBy(n uint) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).StepBy(int(n)))
}

// SortBy applies a custom sorting function to the elements in the iterator
// and returns a new iterator containing the sorted elements.
//
// The sorting function 'fn' should take two arguments, 'a' and 'b' of type V,
// and return true if 'a' should be ordered before 'b', and false otherwise.
//
// Example:
//
//	g.SliceOf("a", "c", "b").
//		Iter().
//		SortBy(func(a, b string) cmp.Ordering { return b.Cmp(a) }).
//		Collect().
//		Print()
//
// Output: Slice[c, b, a]
//
// The returned iterator is of type SeqSlice[V], which implements the iterator
// interface for further iteration over the sorted elements.
func (seq SeqSlice[V]) SortBy(fn func(a, b V) cmp.Ordering) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).SortBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (seq SeqSlice[V]) Take(n uint) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Take(int(n)))
}

// First returns the first element from the sequence.
func (seq SeqSlice[V]) First() Option[V] {
	return OptionOf(iter.Seq[V](seq).First())
}

// Last returns the last element from the sequence.
func (seq SeqSlice[V]) Last() Option[V] {
	return OptionOf(iter.Seq[V](seq).Last())
}

// Nth returns the nth element (0-indexed) in the sequence.
func (seq SeqSlice[V]) Nth(n Int) Option[V] {
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
//	iter := g.Slice[int]{1, 2, 3}.Iter()
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensure cancellation to avoid goroutine leaks.
//	ch := iter.Chan(ctx)
//	for val := range ch {
//	    fmt.Println(val)
//	}
//
// The resulting channel allows streaming elements from the iterator with optional context handling.
func (seq SeqSlice[V]) Chan(ctxs ...context.Context) chan V {
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
// - SeqSlice[V]: An iterator containing unique elements from the original iterator.
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
func (seq SeqSlice[V]) Unique() SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Unique())
}

// Zip combines elements from the current sequence and another sequence into pairs.
// The element types of the two sequences may differ. Iteration stops when either
// sequence is exhausted.
func (seq SeqSlice[V]) Zip[U any](two SeqSlice[U]) SeqPairs[V, U] {
	return SeqPairs[V, U](iter.Seq[V](seq).Zip(iter.Seq[U](two)))
}

// Scan accumulates values of the iterator using a function, yielding all intermediate states.
//
// The function takes an initial accumulator value and a function that combines the accumulator
// with each element. It yields the initial value followed by each accumulated state.
//
// Params:
//
// - init (A): The initial accumulator value. The accumulator type may differ from the element type.
// - fn (func(acc A, val V) A): The function that combines the accumulator with each element.
//
// Returns:
//
// - SeqSlice[A]: A sequence of all intermediate accumulator states.
//
// Example usage:
//
//	numbers := g.Slice[int]{1, 2, 3, 4}.Iter()
//	sums := numbers.Scan(0, func(acc, val int) int {
//		return acc + val
//	})
//	// sums will yield: 0, 1, 3, 6, 10
//
//	words := g.Slice[string]{"a", "b", "c"}.Iter()
//	concatenated := words.Scan("", func(acc, val string) string {
//		return acc + val
//	})
//	// concatenated will yield: "", "a", "ab", "abc"
func (seq SeqSlice[V]) Scan[A any](init A, fn func(acc A, val V) A) SeqSlice[A] {
	return func(yield func(A) bool) {
		if !yield(init) {
			return
		}
		iter.Seq[V](seq).Scan(init, fn)(yield)
	}
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
//	iter := g.Slice[int]{1, 2, 3, 4, 5}.Iter()
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
func (seq SeqSlice[V]) Find(fn func(v V) bool) Option[V] {
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
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6}
//	windows := slice.Iter().Windows(3).Collect()
//
// Output: [Slice[1, 2, 3] Slice[2, 3, 4] Slice[3, 4, 5] Slice[4, 5, 6]]
//
// The resulting iterator will yield sliding windows of elements, each containing the specified number of elements.
func (seq SeqSlice[V]) Windows(n Int) SeqSlices[V] {
	return SeqSlices[V](iter.Windows(iter.Seq[V](seq), int(n)))
}

// Context allows the iteration to be controlled with a context.Context.
func (seq SeqSlice[V]) Context(ctx context.Context) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).Context(ctx))
}

// MaxBy returns the maximum element in the sequence using the provided comparison function.
func (seq SeqSlice[V]) MaxBy(fn func(V, V) cmp.Ordering) Option[V] {
	return OptionOf(iter.Seq[V](seq).MaxBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// Min returns the minimum element in the sequence using the provided comparison function.
func (seq SeqSlice[V]) MinBy(fn func(V, V) cmp.Ordering) Option[V] {
	return OptionOf(iter.Seq[V](seq).MinBy(func(a, b V) bool { return fn(a, b) == cmp.Less }))
}

// Next extracts the next element from the iterator and advances it.
//
// This method consumes the next element from the iterator and returns it wrapped in an Option.
// The iterator itself is modified to point to the remaining elements.
// This is similar to calling Pull() but more convenient for single-element extraction.
//
// Returns:
// - Option[V]: Some(value) if an element exists, None if the iterator is exhausted.
func (seq *SeqSlice[V]) Next() Option[V] {
	if value, remaining, ok := iter.Seq[V](*seq).Next(); ok {
		*seq = SeqSlice[V](remaining)
		return Some(value)
	}

	return None[V]()
}

// FromChan converts a channel into an iterator.
//
// This function takes a channel as input and converts its elements into an iterator,
// allowing seamless integration of channels into iterator-based processing pipelines.
// It continuously reads from the channel until it's closed,
// yielding each element to the provided yield function.
//
// Parameters:
// - ch (<-chan V): The input channel to convert into an iterator.
//
// Returns:
// - SeqSlice[V]: An iterator that yields elements from the channel.
//
// Example usage:
//
//	ch := make(chan int)
//	go func() {
//		defer close(ch)
//		for i := 1; i <= 5; i++ {
//			ch <- i
//		}
//	}()
//
//	// Convert the channel into an iterator and apply filtering and mapping operations.
//	g.FromChan(ch).
//		Filter(func(i int) bool { return i%2 == 0 }). // Filter even numbers.
//		Map(func(i int) int { return i * 2 }).        // Double each element.
//		Collect().                                    // Collect the results into a slice.
//		Print()                                       // Print the collected results.
//
// Output: Slice[4, 8]
//
// The resulting iterator will yield elements from the provided channel, filtering out odd numbers,
// doubling each even number, and finally collecting the results into a slice.
func FromChan[V any](ch <-chan V) SeqSlice[V] {
	return SeqSlice[V](iter.FromChan(ch))
}

// TakeWhile yields elements while the predicate returns true, stopping at the first false.
func (seq SeqSlice[V]) TakeWhile(fn func(V) bool) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).TakeWhile(fn))
}

// SkipWhile skips elements while the predicate returns true, then yields the rest.
func (seq SeqSlice[V]) SkipWhile(fn func(V) bool) SeqSlice[V] {
	return SeqSlice[V](iter.Seq[V](seq).SkipWhile(fn))
}
