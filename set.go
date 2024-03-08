package g

import (
	"fmt"
	"strings"
)

// NewSet creates a new Set of the specified size or an empty Set if no size is provided.
func NewSet[T comparable](size ...int) Set[T] {
	if len(size) == 0 {
		return make(Set[T], 0)
	}

	return make(Set[T], size[0])
}

// SetOf creates a new generic set containing the provided elements.
func SetOf[T comparable](values ...T) Set[T] {
	set := NewSet[T](len(values))
	for _, v := range values {
		set.Add(v)
	}

	return set
}

// TransformSet applies the given function to each element of a Set and returns a new Set
// containing the transformed values.
//
// Parameters:
//
// - s: The input Set.
// - fn: The function to apply to each element of the input Set.
//
// Returns:
//
// A new Set containing the results of applying the function to each element of the input Set.
func TransformSet[T, U comparable](s Set[T], fn func(T) U) Set[U] {
	return mapSet(s.Iter(), fn).Collect()
}

// Iter returns an iterator (seqSet[T]) for the Set, allowing for sequential iteration
// over its elements. It is commonly used in combination with higher-order functions,
// such as 'ForEach' or 'TransformSet', to perform operations on each element of the Set.
//
// Returns:
//
// A pointer to a seqSet[T], which can be used for sequential iteration over the elements of the Set.
//
// Example usage:
//
//	iter := g.SetOf(1, 2, 3).Iter()
//	iter.ForEach(func(val T) {
//	    fmt.Println(val) // Replace this with the function logic you need.
//	})
//
// The 'Iter' method provides a convenient way to traverse the elements of a Set
// in a functional style, enabling operations like mapping or filtering.
func (s Set[T]) Iter() seqSet[T] { return liftSet(s) }

// Add adds the provided elements to the set and returns the modified set.
func (s Set[T]) Add(values ...T) Set[T] {
	for _, v := range values {
		s[v] = struct{}{}
	}

	return s
}

// Remove removes the specified values from the Set.
func (s Set[T]) Remove(values ...T) Set[T] {
	for _, v := range values {
		delete(s, v)
	}

	return s
}

// Len returns the number of values in the Set.
func (s Set[T]) Len() int { return len(s) }

// Contains checks if the Set contains the specified value.
func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// ContainsAny checks if the Set contains any element from another Set.
func (s Set[T]) ContainsAny(other Set[T]) bool {
	for v := range other {
		if s.Contains(v) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the Set contains all elements from another Set.
func (s Set[T]) ContainsAll(other Set[T]) bool {
	if s.Len() < other.Len() {
		return false
	}

	for v := range other {
		if !s.Contains(v) {
			return false
		}
	}

	return true
}

// Clone creates a new Set that is a copy of the original Set.
func (s Set[T]) Clone() Set[T] { return s.Iter().Collect() }

// ToSlice returns a new Slice with the same elements as the Set[T].
func (s Set[T]) ToSlice() Slice[T] {
	sl := NewSlice[T](0, s.Len())
	s.Iter().ForEach(func(v T) { sl = sl.Append(v) })

	return sl
}

// Intersection returns the intersection of the current set and another set, i.e., elements
// present in both sets.
//
// Parameters:
//
// - other Set[T]: The other set to calculate the intersection with.
//
// Returns:
//
// - Set[T]: A new Set containing the intersection of the two sets.
//
// Example usage:
//
//	s1 := g.SetOf(1, 2, 3, 4, 5)
//	s2 := g.SetOf(4, 5, 6, 7, 8)
//	intersection := s1.Intersection(s2)
//
// The resulting intersection will be: [4, 5].
func (s Set[T]) Intersection(other Set[T]) seqSet[T] {
	if s.Len() <= other.Len() {
		return intersectionS(s.Iter(), other)
	}

	return intersectionS(other.Iter(), s)
}

// Difference returns the difference between the current set and another set,
// i.e., elements present in the current set but not in the other set.
//
// Parameters:
//
// - other Set[T]: The other set to calculate the difference with.
//
// Returns:
//
// - Set[T]: A new Set containing the difference between the two sets.
//
// Example usage:
//
//	s1 := g.SetOf(1, 2, 3, 4, 5)
//	s2 := g.SetOf(4, 5, 6, 7, 8)
//	diff := s1.Difference(s2)
//
// The resulting diff will be: [1, 2, 3].
func (s Set[T]) Difference(other Set[T]) seqSet[T] { return differenceS(s.Iter(), other) }

// Union returns a new set containing the unique elements of the current set and the provided
// other set.
//
// Parameters:
//
// - other Set[T]: The other set to create the union with.
//
// Returns:
//
// - Set[T]: A new Set containing the unique elements of the current set and the provided
// other set.
//
// Example usage:
//
//	s1 := g.SetOf(1, 2, 3)
//	s2 := g.SetOf(3, 4, 5)
//	union := s1.Union(s2)
//
// The resulting union set will be: [1, 2, 3, 4, 5].
func (s Set[T]) Union(other Set[T]) seqSet[T] {
	if s.Len() >= other.Len() {
		return s.Iter().Chain(other.Difference(s))
	}

	return other.Iter().Chain(s.Difference(other))
}

// SymmetricDifference returns the symmetric difference between the current set and another
// set, i.e., elements present in either the current set or the other set but not in both.
//
// Parameters:
//
// - other Set[T]: The other set to calculate the symmetric difference with.
//
// Returns:
//
// - Set[T]: A new Set containing the symmetric difference between the two sets.
//
// Example usage:
//
//	s1 := g.SetOf(1, 2, 3, 4, 5)
//	s2 := g.SetOf(4, 5, 6, 7, 8)
//	symDiff := s1.SymmetricDifference(s2)
//
// The resulting symDiff will be: [1, 2, 3, 6, 7, 8].
func (s Set[T]) SymmetricDifference(other Set[T]) seqSet[T] {
	return s.Difference(other).Chain(other.Difference(s))
}

// Subset checks if the current set 's' is a subset of the provided 'other' set.
// A set 's' is a subset of 'other' if all elements of 's' are also elements of 'other'.
//
// Parameters:
//
// - other Set[T]: The other set to compare with.
//
// Returns:
//
// - bool: true if 's' is a subset of 'other', false otherwise.
//
// Example usage:
//
//	s1 := g.SetOf(1, 2, 3)
//	s2 := g.SetOf(1, 2, 3, 4, 5)
//	isSubset := s1.Subset(s2) // Returns true
func (s Set[T]) Subset(other Set[T]) bool { return other.ContainsAll(s) }

// Superset checks if the current set 's' is a superset of the provided 'other' set.
// A set 's' is a superset of 'other' if all elements of 'other' are also elements of 's'.
//
// Parameters:
//
// - other Set[T]: The other set to compare with.
//
// Returns:
//
// - bool: true if 's' is a superset of 'other', false otherwise.
//
// Example usage:
//
//	s1 := g.SetOf(1, 2, 3, 4, 5)
//	s2 := g.SetOf(1, 2, 3)
//	isSuperset := s1.Superset(s2) // Returns true
func (s Set[T]) Superset(other Set[T]) bool { return s.ContainsAll(other) }

// Eq checks if two Sets are equal.
func (s Set[T]) Eq(other Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}

	for v := range other {
		if !s.Contains(v) {
			return false
		}
	}

	return true
}

// Ne checks if two Sets are not equal.
func (s Set[T]) Ne(other Set[T]) bool { return !s.Eq(other) }

// Clear removes all values from the Set.
func (s Set[T]) Clear() Set[T] { return s.Remove(s.ToSlice()...) }

// Empty checks if the Set is empty.
func (s Set[T]) Empty() bool { return s.Len() == 0 }

// String returns a string representation of the Set.
func (s Set[T]) String() string {
	var builder strings.Builder

	s.Iter().ForEach(func(v T) { builder.WriteString(fmt.Sprintf("%v, ", v)) })

	return String(builder.String()).TrimRight(", ").Format("Set{%s}").Std()
}

// Print prints the elements of the Set to the standard output (console)
// and returns the Set unchanged.
func (s Set[T]) Print() Set[T] { fmt.Println(s); return s }
