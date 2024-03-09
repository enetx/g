package g

import "iter"

type seqSet[V comparable] iter.Seq[V]

func (seq seqSet[V]) pull() (func() (V, bool), func()) { return iter.Pull(iter.Seq[V](seq)) }

// Inspect creates a new iterator that wraps around the current iterator
// and allows inspecting each element as it passes through.
func (seq seqSet[V]) Inspect(fn func(v V)) seqSet[V] { return inspectSet(seq, fn) }

// Collect gathers all elements from the iterator into a Set.
func (seq seqSet[V]) Collect() Set[V] {
	collection := NewSet[V]()

	seq(func(v V) bool {
		collection.Add(v)
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
// - seqs ([]seqSet[V]): Other iterators to be concatenated with the current iterator.
//
// Returns:
//
// - seqSet[V]: A new iterator containing elements from the current iterator and the provided iterators.
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
func (seq seqSet[V]) Chain(seqs ...seqSet[V]) seqSet[V] {
	return chainSet(append([]seqSet[V]{seq}, seqs...)...)
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
//	iter := g.SetOf(1, 2, 3).Iter()
//	iter.ForEach(func(val T) {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	})
//
// The provided function will be applied to each element in the iterator.
func (seq seqSet[V]) ForEach(fn func(v V)) {
	seq(func(v V) bool {
		fn(v)
		return true
	})
}

// The iteration will stop when the provided function returns false for an element.
func (seq seqSet[V]) Range(fn func(v V) bool) {
	seq(func(v V) bool {
		return fn(v)
	})
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
// - seqSet[V]: A new iterator containing the elements that satisfy the given condition.
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
func (seq seqSet[V]) Filter(fn func(V) bool) seqSet[V] { return filterSet(seq, fn) }

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
// - seqSet[V]: A new iterator containing the elements that do not satisfy the given condition.
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
func (seq seqSet[V]) Exclude(fn func(V) bool) seqSet[V] { return excludeSet(seq, fn) }

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
// - seqSet[V]: A new iterator containing elements transformed by the provided function.
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
func (seq seqSet[V]) Map(transform func(V) V) seqSet[V] { return mapSet(seq, transform) }

func liftSet[V comparable](slice Set[V]) seqSet[V] {
	return func(yield func(V) bool) {
		for v := range slice {
			if !yield(v) {
				return
			}
		}
	}
}

func inspectSet[V comparable](seq seqSet[V], fn func(V)) seqSet[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			fn(v)
			return yield(v)
		})
	}
}

func chainSet[V comparable](seqs ...seqSet[V]) seqSet[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			seq(func(v V) bool {
				return yield(v)
			})
		}
	}
}

func mapSet[V, W comparable](seq seqSet[V], fn func(V) W) seqSet[W] {
	return func(yield func(W) bool) {
		seq(func(v V) bool {
			return yield(fn(v))
		})
	}
}

func filterSet[V comparable](seq seqSet[V], fn func(V) bool) seqSet[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if fn(v) {
				return yield(v)
			}
			return true
		})
	}
}

func excludeSet[V comparable](seq seqSet[V], fn func(V) bool) seqSet[V] {
	return filterSet(seq, func(v V) bool { return !fn(v) })
}

func differenceS[V comparable](seq seqSet[V], other Set[V]) seqSet[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if other.Contains(v) {
				return true
			}
			return yield(v)
		})
	}
}

func intersectionS[V comparable](seq seqSet[V], other Set[V]) seqSet[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if other.Contains(v) {
				return yield(v)
			}
			return true
		})
	}
}
