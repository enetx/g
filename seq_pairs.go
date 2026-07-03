package g

import "github.com/enetx/iter"

// SeqPairs is a lazy sequence of key-value pairs produced by zipping two sequences.
// Unlike SeqMapOrd, its first type parameter is not required to be comparable, so it can
// carry pairs of arbitrary types.
//
// Note: SeqPairs deliberately does not provide methods whose signatures would instantiate
// generic containers with Pair[K, V] as a type argument (e.g. Collect returning
// Slice[Pair[K, V]]); such signatures would create an instantiation cycle with the
// generic Zip methods. Collect therefore returns a plain []Pair[K, V], which callers
// can convert to a Slice with g.SliceOf(pairs...) at a concrete type.
type SeqPairs[K, V any] iter.Seq2[K, V]

// Keys returns a sequence of the first elements of each pair.
func (seq SeqPairs[K, V]) Keys() SeqSlice[K] {
	return func(yield func(K) bool) {
		seq(func(k K, _ V) bool { return yield(k) })
	}
}

// Values returns a sequence of the second elements of each pair.
func (seq SeqPairs[K, V]) Values() SeqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(_ K, v V) bool { return yield(v) })
	}
}

// Unzip consumes the sequence and collects the first and second elements
// of each pair into two separate slices.
func (seq SeqPairs[K, V]) Unzip() (Slice[K], Slice[V]) {
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

// Map transforms each pair into a single value using the given function,
// returning a sequence of the results.
func (seq SeqPairs[K, V]) Map[T any](fn func(K, V) T) SeqSlice[T] {
	return func(yield func(T) bool) {
		seq(func(k K, v V) bool { return yield(fn(k, v)) })
	}
}

// Filter returns a sequence containing only the pairs for which fn returns true.
func (seq SeqPairs[K, V]) Filter(fn func(K, V) bool) SeqPairs[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			if fn(k, v) {
				return yield(k, v)
			}

			return true
		})
	}
}

// FilterByKey returns a sequence lazily yielding only the pairs whose first
// element satisfies the provided predicate; second elements are not inspected.
//
// It lifts a single-parameter predicate to the pair-wise Filter — composes
// with f.* factories:
//
//	pairs.FilterByKey(f.Gt(10))
func (seq SeqPairs[K, V]) FilterByKey(fn func(K) bool) SeqPairs[K, V] {
	return seq.Filter(func(k K, _ V) bool { return fn(k) })
}

// FilterByValue returns a sequence lazily yielding only the pairs whose second
// element satisfies the provided predicate; first elements are not inspected.
//
// It lifts a single-parameter predicate to the pair-wise Filter — composes
// with f.* factories:
//
//	pairs.FilterByValue(f.Eq("ok"))
func (seq SeqPairs[K, V]) FilterByValue(fn func(V) bool) SeqPairs[K, V] {
	return seq.Filter(func(_ K, v V) bool { return fn(v) })
}

// Exclude returns a sequence containing only the pairs for which fn returns false.
func (seq SeqPairs[K, V]) Exclude(fn func(K, V) bool) SeqPairs[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			if fn(k, v) {
				return true
			}

			return yield(k, v)
		})
	}
}

// Take returns a sequence containing at most n leading pairs.
func (seq SeqPairs[K, V]) Take(n uint) SeqPairs[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			if n == 0 {
				return false
			}

			n--

			return yield(k, v)
		})
	}
}

// Skip returns a sequence that skips the first n pairs.
func (seq SeqPairs[K, V]) Skip(n uint) SeqPairs[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			if n > 0 {
				n--
				return true
			}

			return yield(k, v)
		})
	}
}

// TakeWhile yields pairs while the predicate returns true, stopping at the first false.
func (seq SeqPairs[K, V]) TakeWhile(fn func(K, V) bool) SeqPairs[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			if !fn(k, v) {
				return false
			}

			return yield(k, v)
		})
	}
}

// SkipWhile skips pairs while the predicate returns true, then yields the rest.
func (seq SeqPairs[K, V]) SkipWhile(fn func(K, V) bool) SeqPairs[K, V] {
	return func(yield func(K, V) bool) {
		skipping := true

		seq(func(k K, v V) bool {
			if skipping {
				if fn(k, v) {
					return true
				}

				skipping = false
			}

			return yield(k, v)
		})
	}
}

// ForEach applies fn to each pair in the sequence.
func (seq SeqPairs[K, V]) ForEach(fn func(K, V)) {
	seq(func(k K, v V) bool {
		fn(k, v)
		return true
	})
}

// Count consumes the sequence and returns the number of pairs.
func (seq SeqPairs[K, V]) Count() Int {
	var n Int

	seq(func(K, V) bool {
		n++
		return true
	})

	return n
}

// Fold reduces the sequence of pairs to a single value using an accumulator.
// The accumulator type may differ from the key and value types.
func (seq SeqPairs[K, V]) Fold[A any](init A, fn func(acc A, k K, v V) A) A {
	seq(func(k K, v V) bool {
		init = fn(init, k, v)
		return true
	})

	return init
}

// All returns true if fn returns true for every pair in the sequence.
// It stops at the first pair for which fn returns false.
func (seq SeqPairs[K, V]) All(fn func(K, V) bool) bool {
	result := true

	seq(func(k K, v V) bool {
		if !fn(k, v) {
			result = false
			return false
		}

		return true
	})

	return result
}

// Any returns true if fn returns true for at least one pair in the sequence.
// It stops at the first pair for which fn returns true.
func (seq SeqPairs[K, V]) Any(fn func(K, V) bool) bool {
	var result bool

	seq(func(k K, v V) bool {
		if fn(k, v) {
			result = true
			return false
		}

		return true
	})

	return result
}

// Find returns the first pair satisfying fn, or None if no pair matches.
func (seq SeqPairs[K, V]) Find(fn func(K, V) bool) Option[Pair[K, V]] {
	var (
		result Pair[K, V]
		found  bool
	)

	seq(func(k K, v V) bool {
		if fn(k, v) {
			result, found = Pair[K, V]{Key: k, Value: v}, true
			return false
		}

		return true
	})

	return OptionOf(result, found)
}

// Inspect calls fn on each pair as it passes through the sequence, without modifying it.
func (seq SeqPairs[K, V]) Inspect(fn func(K, V)) SeqPairs[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			fn(k, v)
			return yield(k, v)
		})
	}
}

// Collect consumes the sequence and returns all pairs as a plain []Pair slice.
// Convert to a g.Slice at a concrete type with g.SliceOf(pairs...) if chaining is needed.
func (seq SeqPairs[K, V]) Collect() []Pair[K, V] {
	var result []Pair[K, V]

	seq(func(k K, v V) bool {
		result = append(result, Pair[K, V]{Key: k, Value: v})
		return true
	})

	return result
}
