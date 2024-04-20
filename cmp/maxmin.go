package cmp

import "cmp"

// Max returns the maximum value among the given values. The values must be of a type that implements
// the cmp.Ordered interface for comparison.
func Max[T cmp.Ordered](a T, b ...T) T {
	m := a

	for _, v := range b {
		if v > m {
			m = v
		}
	}

	return m
}

// Min returns the minimum value among the given values. The values must be of a type that implements
// the cmp.Ordered interface for comparison.
func Min[T cmp.Ordered](a T, b ...T) T {
	m := a

	for _, v := range b {
		if v < m {
			m = v
		}
	}

	return m
}
