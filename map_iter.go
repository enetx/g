package g

import (
	"iter"
	"runtime"
)

// IterPar parallelizes the SeqMap using the specified number of workers.
func (seq SeqMap[K, V]) Parallel(workers ...Int) SeqMapPar[K, V] {
	numCPU := Int(runtime.NumCPU())
	count := Slice[Int](workers).Get(0).UnwrapOr(numCPU)

	if count.Lte(0) {
		count = numCPU
	}

	return SeqMapPar[K, V]{
		seq:     seq,
		workers: count,
		process: func(p Pair[K, V]) (Pair[K, V], bool) { return p, true },
	}
}

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
func (seq SeqMap[K, V]) Pull() (func() (K, V, bool), func()) { return iter.Pull2(iter.Seq2[K, V](seq)) }

// Take returns a new iterator with the first n elements.
// The function creates a new iterator containing the first n elements from the original iterator.
func (seq SeqMap[K, V]) Take(n uint) SeqMap[K, V] {
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

// Keys returns an iterator containing all the keys in the ordered Map.
func (seq SeqMap[K, V]) Keys() SeqSlice[K] {
	return func(yield func(K) bool) {
		seq(func(k K, _ V) bool {
			return yield(k)
		})
	}
}

// Values returns an iterator containing all the values in the ordered Map.
func (seq SeqMap[K, V]) Values() SeqSlice[V] {
	return func(yield func(V) bool) {
		seq(func(_ K, v V) bool {
			return yield(v)
		})
	}
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
//	iter1 := g.NewMap[int, string]().Set(1, "a").Iter()
//	iter2 := g.NewMap[int, string]().Set(2, "b").Iter()
//
//	// Concatenating iterators and collecting the result.
//	iter1.Chain(iter2).Collect().Print()
//
// Output: Map{1:a, 2:b} // The output order may vary as Map is not ordered.
//
// The resulting iterator will contain elements from both iterators.
func (seq SeqMap[K, V]) Chain(seqs ...SeqMap[K, V]) SeqMap[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range append([]SeqMap[K, V]{seq}, seqs...) {
			seq(func(k K, v V) bool {
				return yield(k, v)
			})
		}
	}
}

// Count consumes the iterator, counting the number of iterations and returning it.
func (seq SeqMap[K, V]) Count() Int {
	var counter Int
	seq(func(K, V) bool {
		counter++
		return true
	})

	return counter
}

// Collect collects all key-value pairs from the iterator and returns a Map.
func (seq SeqMap[K, V]) Collect() Map[K, V] {
	collection := NewMap[K, V]()

	seq(func(k K, v V) bool {
		collection.Set(k, v)
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
//	m := g.NewMap[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
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
func (seq SeqMap[K, V]) Filter(fn func(K, V) bool) SeqMap[K, V] { return filterMap(seq, fn) }

func filterMap[K comparable, V any](seq SeqMap[K, V], fn func(K, V) bool) SeqMap[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			if fn(k, v) {
				return yield(k, v)
			}
			return true
		})
	}
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
//	m := g.NewMap[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
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
	return filterMap(seq, func(k K, v V) bool { return !fn(k, v) })
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
//	m.Set(1, 1)
//	f := m.Iter().Find(func(_ int, v int) bool { return v == 1 })
//	if f.IsSome() {
//		print(f.Some().Key)
//	}
//
// The resulting Option may contain the first element that satisfies the condition, or None if not found.
func (seq SeqMap[K, V]) Find(fn func(k K, v V) bool) (r Option[Pair[K, V]]) {
	seq(func(k K, v V) bool {
		if !fn(k, v) {
			return true
		}
		r = Some(Pair[K, V]{k, v})
		return false
	})

	return r
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
//	m := g.NewMap[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
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
func (seq SeqMap[K, V]) ForEach(fn func(k K, v V)) {
	seq(func(k K, v V) bool {
		fn(k, v)
		return true
	})
}

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each key-value pair as it passes through.
func (seq SeqMap[K, V]) Inspect(fn func(k K, v V)) SeqMap[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			fn(k, v)
			return yield(k, v)
		})
	}
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
//	m := g.NewMap[int, int]().
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Set(4, 4).
//		Set(5, 5)
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
func (seq SeqMap[K, V]) Map(transform func(K, V) (K, V)) SeqMap[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			return yield(transform(k, v))
		})
	}
}

// The iteration will stop when the provided function returns false for an element.
func (seq SeqMap[K, V]) Range(fn func(k K, v V) bool) {
	seq(func(k K, v V) bool {
		return fn(k, v)
	})
}

func seqMap[K comparable, V any](hashmap map[K]V) SeqMap[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range hashmap {
			if !yield(k, v) {
				return
			}
		}
	}
}
