package g

import "fmt"

// SliceMap applies the given function to each element of a Slice and returns a new Slice
// containing the transformed values.
//
// Parameters:
// - sl: The input Slice.
// - fn: The function to apply to each element of the input Slice.
//
// Returns:
// A new Slice containing the results of applying the function to each element of the input Slice.
func SliceMap[T, E any](sl Slice[T], fn func(T) E) Slice[E] {
	result := NewSlice[E](0, sl.Len())
	sl.ForEach(func(t T) { result = result.Append(fn(t)) })

	return result
}

// SetMap applies the given function to each element of a Set and returns a new Set
// containing the transformed values.
//
// Parameters:
// - s: The input Set.
// - fn: The function to apply to each element of the input Set.
//
// Returns:
// A new Set containing the results of applying the function to each element of the input Set.
func SetMap[T, E comparable](s Set[T], fn func(T) E) Set[E] {
	result := NewSet[E](s.Len())
	s.ForEach(func(t T) { result.Add(fn(t)) })

	return result
}

// Sprintf formats according to a format specifier and returns the resulting String.
func Sprintf[T ~string](str T, a ...any) String { return NewString(fmt.Sprintf(string(str), a...)) }
